// Package render converts logical symbols to exact integer-module raster, PNG
// or SVG output through Image/PNG/SVG. It adds no quiet zone or antialiasing;
// margins belong to the supplied symbol. Inputs remain caller-owned and writers
// are never closed. Returned pixels are mutable caller-owned image storage.
// Geometry limits do not bound all copies/callbacks/output I/O or guarantee
// physical readability; see BARCODE-DEC-005 and rendering/golden fixtures.
package render

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"strconv"

	"github.com/faustbrian/go-barcode/v2/barcode"
)

const (
	defaultMaxDimension = 32768
	defaultMaxPixels    = 64 * 1024 * 1024
)

var (
	// ErrInvalidSymbol and ErrLimitExceeded classify rendering failures.
	ErrInvalidSymbol = errors.New("render: invalid symbol")
	// ErrLimitExceeded reports a rendering dimension or allocation violation.
	ErrLimitExceeded = errors.New("render: limit exceeded")
)

// Limits bounds final geometry before the new raster or SVG writes. Logical
// accessors make defensive copies before this admission step; upstream symbol
// storage, color callbacks, PNG compression and writer allocations are not capped.
type Limits struct {
	// MaxDimension caps either final output axis; zero selects 32768, negative fails.
	MaxDimension int
	// MaxPixels caps final width*height; zero selects 64Mi pixels, negative fails.
	// This is geometry admission, not a total temporary-memory or output-byte budget.
	MaxPixels int
}

// Options controls integer scaling, colors, and allocation limits.
type Options struct {
	// Scale is an integer module multiplier; zero selects 1, negative fails.
	Scale int
	// Foreground is borrowed Color input; nil selects opaque black.
	// Its RGBA callback must be prompt/safe; it is not isolated from panics.
	Foreground color.Color
	// Background is borrowed Color input; nil selects opaque white.
	// Arbitrary color choices do not guarantee scanner contrast.
	Background color.Color
	// Limits constrains final pixel geometry after obtaining logical defensive copies.
	Limits Limits
}

// Image renders each logical module as an exact integer square into a new
// mutable NRGBA image. Invalid/overflowing/excess geometry fails before raster
// allocation, but logical defensive copies precede these checks. Borrowed Color
// callbacks can block/panic and are not canceled; output contrast is caller policy.
func Image(symbol barcode.Symbol, options Options) (image.Image, error) {
	geometry, err := prepare(symbol, options)
	if err != nil {
		return nil, err
	}
	output := image.NewNRGBA(image.Rect(0, 0, geometry.width, geometry.height))
	draw.Draw(output, output.Bounds(), &image.Uniform{C: geometry.background}, image.Point{}, draw.Src)

	if bars, ok := symbol.Bars(); ok {
		x := 0
		for _, run := range bars.Runs() {
			if run.Dark {
				draw.Draw(output, image.Rect(
					x*geometry.scale,
					0,
					(x+run.Width)*geometry.scale,
					geometry.height,
				), &image.Uniform{C: geometry.foreground}, image.Point{}, draw.Src)
			}
			x += run.Width
		}
		return output, nil
	}

	matrix := symbol.Matrix()
	for y := 0; y < matrix.Height(); y++ {
		for x := 0; x < matrix.Width(); x++ {
			if matrix.At(x, y) {
				draw.Draw(output, image.Rect(
					x*geometry.scale,
					y*geometry.scale,
					(x+1)*geometry.scale,
					(y+1)*geometry.scale,
				), &image.Uniform{C: geometry.foreground}, image.Point{}, draw.Src)
			}
		}
	}

	return output, nil
}

// PNG constructs the raster and writes its PNG encoding to a borrowed writer.
// Nil writer fails; writer errors are wrapped and partial output may remain.
// It does not close the writer, retry, cap encoded bytes or impose an I/O deadline.
// Canonical fixtures describe selected bytes, not stability across all future codecs.
func PNG(writer io.Writer, symbol barcode.Symbol, options Options) error {
	if writer == nil {
		return ErrInvalidSymbol
	}
	raster, err := Image(symbol, options)
	if err != nil {
		return err
	}
	if err := png.Encode(writer, raster); err != nil {
		return fmt.Errorf("render PNG: %w", err)
	}

	return nil
}

// SVG streams integer-coordinate rectangles from the logical symbol after
// geometry admission. It returns wrapped writer failures; partial documents are
// possible and no transactional publication, writer closure or deadline is added.
// Colors become fixed numeric NRGBA attributes, never arbitrary XML content.
func SVG(writer io.Writer, symbol barcode.Symbol, options Options) error {
	if writer == nil {
		return ErrInvalidSymbol
	}
	geometry, err := prepare(symbol, options)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(writer,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" shape-rendering="crispEdges">`,
		geometry.width, geometry.height, geometry.width, geometry.height,
	); err != nil {
		return fmt.Errorf("render SVG: %w", err)
	}
	if err = writeRect(writer, 0, 0, geometry.width, geometry.height, geometry.background); err != nil {
		return err
	}

	if bars, ok := symbol.Bars(); ok {
		x := 0
		for _, run := range bars.Runs() {
			if run.Dark {
				if err = writeRect(writer, x*geometry.scale, 0, run.Width*geometry.scale, geometry.height, geometry.foreground); err != nil {
					return err
				}
			}
			x += run.Width
		}
	} else {
		matrix := symbol.Matrix()
		for y := 0; y < matrix.Height(); y++ {
			for x := 0; x < matrix.Width(); x++ {
				if matrix.At(x, y) {
					if err = writeRect(writer, x*geometry.scale, y*geometry.scale, geometry.scale, geometry.scale, geometry.foreground); err != nil {
						return err
					}
				}
			}
		}
	}
	if _, err = io.WriteString(writer, "</svg>\n"); err != nil {
		return fmt.Errorf("render SVG: %w", err)
	}

	return nil
}

type prepared struct {
	width      int
	height     int
	scale      int
	foreground color.NRGBA
	background color.NRGBA
}

func prepare(symbol barcode.Symbol, options Options) (prepared, error) {
	logicalWidth, logicalHeight := 0, 0
	if bars, ok := symbol.Bars(); ok {
		logicalWidth, logicalHeight = bars.Width(), bars.Height()
	} else {
		matrix := symbol.Matrix()
		logicalWidth, logicalHeight = matrix.Width(), matrix.Height()
	}
	if min(logicalWidth, logicalHeight) == 0 {
		return prepared{}, ErrInvalidSymbol
	}
	scale := options.Scale
	if scale < 0 {
		return prepared{}, ErrLimitExceeded
	}
	if scale == 0 {
		scale = 1
	}
	width, widthOK := boundedProduct(logicalWidth, scale, math.MaxInt)
	height, heightOK := boundedProduct(logicalHeight, scale, math.MaxInt)
	if !widthOK {
		return prepared{}, ErrLimitExceeded
	}
	if !heightOK {
		return prepared{}, ErrLimitExceeded
	}
	maxDimension := options.Limits.MaxDimension
	if maxDimension == 0 {
		maxDimension = defaultMaxDimension
	}
	maxPixels := options.Limits.MaxPixels
	if maxPixels == 0 {
		maxPixels = defaultMaxPixels
	}
	if maxDimension < 1 {
		return prepared{}, ErrLimitExceeded
	}
	if maxPixels < 1 {
		return prepared{}, ErrLimitExceeded
	}
	if width > maxDimension {
		return prepared{}, ErrLimitExceeded
	}
	if height > maxDimension {
		return prepared{}, ErrLimitExceeded
	}
	if _, ok := boundedProduct(width, height, maxPixels); !ok {
		return prepared{}, ErrLimitExceeded
	}

	foreground := color.NRGBA{A: 255}
	if options.Foreground != nil {
		foreground = color.NRGBAModel.Convert(options.Foreground).(color.NRGBA)
	}
	background := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	if options.Background != nil {
		background = color.NRGBAModel.Convert(options.Background).(color.NRGBA)
	}

	return prepared{
		width: width, height: height, scale: scale,
		foreground: foreground, background: background,
	}, nil
}

// Positive operands are established by prepare; division avoids overflowing
// multiplication before comparing the pixel budget. Do not call with zero right.
func boundedProduct(left, right, limit int) (int, bool) {
	if left > math.MaxInt/right {
		return 0, false
	}
	product := left * right

	return product, product <= limit
}

func writeRect(writer io.Writer, x, y, width, height int, value color.NRGBA) error {
	opacity := ""
	if value.A != 255 {
		opacity = ` fill-opacity="` + strconv.FormatFloat(float64(value.A)/255, 'f', 6, 64) + `"`
	}
	_, err := fmt.Fprintf(writer,
		`<rect x="%d" y="%d" width="%d" height="%d" fill="#%02x%02x%02x"%s/>`,
		x, y, width, height, value.R, value.G, value.B, opacity,
	)
	if err != nil {
		return fmt.Errorf("render SVG: %w", err)
	}

	return nil
}
