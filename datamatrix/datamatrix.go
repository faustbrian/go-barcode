// Package datamatrix constructs logical ECC200 matrices using selected ZXing
// encoding and a controlled Base256 path. Encode/EncodeGS1 return copied logical
// values; no decoder, device or multi-symbol assembly lifecycle is owned.
// The profile targets ISO/IEC 16022:2024 edition 3; full sequence assembly remains
// unadvertised (BARCODE-DEC-009). See the local evidence/manifest inventories.
package datamatrix

import (
	"errors"
	"fmt"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/gs1"
	"github.com/makiuchi-d/gozxing"
	zxingdatamatrix "github.com/makiuchi-d/gozxing/datamatrix"
	"github.com/makiuchi-d/gozxing/datamatrix/encoder"
)

const (
	defaultQuietZone = 1
	maxPayloadBytes  = 3116
)

// ErrInvalidInput reports invalid Data Matrix payloads and options.
var ErrInvalidInput = errors.New("datamatrix: invalid input")

// Shape constrains Data Matrix ECC 200 symbol geometry.
type Shape uint8

const (
	// Automatic and the other constants select Data Matrix shape constraints.
	Automatic Shape = iota
	// Square forces a square Data Matrix symbol.
	Square
	// Rectangle forces a rectangular Data Matrix symbol.
	Rectangle
)

// Macro selects a Data Matrix Macro 05 or Macro 06 control header.
type Macro uint8

const (
	// MacroNone encodes ordinary Data Matrix data.
	MacroNone Macro = iota
	// Macro05 emits the ISO/IEC 16022 Macro 05 control codeword.
	Macro05
	// Macro06 emits the ISO/IEC 16022 Macro 06 control codeword.
	Macro06
)

// StructuredAppend identifies one part of a two-to-sixteen-symbol sequence.
type StructuredAppend struct {
	// Index is zero-based and must be less than Total.
	Index int
	// Total is 2..16; sequence assembly is not implemented.
	Total int
	// FileID is 1..254 and identifies the caller-managed sequence.
	FileID byte
}

// Options controls Data Matrix shape, dimensions, and quiet zone.
type Options struct {
	// Shape selects automatic, square or rectangular symbol lookup.
	Shape Shape
	// QuietZone is zero for 1 module, otherwise 1..256.
	QuietZone int
	// MinWidth constrains the symbol before quiet zone; zero requires MinHeight zero.
	MinWidth int
	// MinHeight pairs with MinWidth and must be nonnegative.
	MinHeight int
	// MaxWidth constrains the symbol before quiet zone; zero requires MaxHeight zero.
	MaxWidth int
	// MaxHeight pairs with MaxWidth and cannot be below a nonzero minimum.
	MaxHeight int
	// GS1 validates raw element syntax and inserts FNC1; it conflicts with Macro.
	GS1 bool
	// ECI is zero for absence, otherwise 1..999999; it emits a control, not transcoding.
	ECI int
	// Macro selects 05/06 envelope control and conflicts with GS1/StructuredAppend.
	Macro Macro
	// StructuredAppend supplies copied part metadata and conflicts with Macro;
	// the caller owns the other parts and later assembly.
	StructuredAppend *StructuredAppend
}

// Symbol combines an immutable logical symbol with Data Matrix metadata.
type Symbol struct {
	logical          barcode.Symbol
	shape            Shape
	gs1              bool
	eci              int
	macro            Macro
	structuredAppend *StructuredAppend
}

// Logical returns the immutable format-neutral symbol.
func (symbol Symbol) Logical() barcode.Symbol { return symbol.logical }

// Shape returns the requested Data Matrix shape constraint.
func (symbol Symbol) Shape() Shape { return symbol.shape }

// GS1 reports whether the symbol begins with the FNC1 control codeword.
func (symbol Symbol) GS1() bool { return symbol.gs1 }

// ECI returns the requested ECI assignment, or zero when absent.
func (symbol Symbol) ECI() int { return symbol.eci }

// Macro returns the requested Macro 05 or Macro 06 control block.
func (symbol Symbol) Macro() Macro { return symbol.macro }

// StructuredAppend returns a defensive copy of the sequence header.
func (symbol Symbol) StructuredAppend() (StructuredAppend, bool) {
	if symbol.structuredAppend == nil {
		return StructuredAppend{}, false
	}

	return *symbol.structuredAppend, true
}

// Encode returns an ECC 200 logical matrix with a one-module quiet zone by
// default. Dimension constraints refer to the symbol before its quiet zone.
// Payload is nonempty and capped at 3116 bytes; the controlled Base256 path has
// an additional 1555-byte payload ceiling and symbol capacity can be smaller.
// Input is borrowed during encoding, then copied; ECI does not transcode it.
// Errors wrap ErrInvalidInput and may include peer/GS1 causes. Sequence assembly
// is not advertised; see BARCODE-DEC-009 and datamatrix control-boundary tests.
func Encode(payload []byte, options Options) (Symbol, error) {
	if len(payload) == 0 || len(payload) > maxPayloadBytes || options.Shape > Rectangle ||
		options.QuietZone < 0 || options.MinWidth < 0 || options.MinHeight < 0 ||
		options.MaxWidth < 0 || options.MaxHeight < 0 ||
		(options.MinWidth == 0) != (options.MinHeight == 0) ||
		(options.MaxWidth == 0) != (options.MaxHeight == 0) ||
		options.ECI < 0 || options.ECI >= 1_000_000 || options.Macro > Macro06 ||
		options.GS1 && options.Macro != MacroNone ||
		options.StructuredAppend != nil && (options.StructuredAppend.Total < 2 ||
			options.StructuredAppend.Total > 16 || options.StructuredAppend.Index < 0 ||
			options.StructuredAppend.Index >= options.StructuredAppend.Total ||
			options.StructuredAppend.FileID == 0 || options.StructuredAppend.FileID == 255 ||
			options.Macro != MacroNone) ||
		options.MaxWidth != 0 && options.MinWidth > options.MaxWidth ||
		options.MaxHeight != 0 && options.MinHeight > options.MaxHeight {
		return Symbol{}, ErrInvalidInput
	}
	quietZone := options.QuietZone
	if quietZone == 0 {
		quietZone = defaultQuietZone
	}
	if quietZone > 256 {
		return Symbol{}, ErrInvalidInput
	}
	if options.GS1 {
		if _, err := gs1.ParseRaw(string(payload), gs1.ParseLimits{MaxInputBytes: maxPayloadBytes}); err != nil {
			return Symbol{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
		}
	}
	hints := make(map[gozxing.EncodeHintType]interface{}, 3)
	shape := encoder.SymbolShapeHint_FORCE_NONE
	switch options.Shape {
	case Automatic:
	case Square:
		shape = encoder.SymbolShapeHint_FORCE_SQUARE
		hints[gozxing.EncodeHintType_DATA_MATRIX_SHAPE] = shape
	case Rectangle:
		shape = encoder.SymbolShapeHint_FORCE_RECTANGLE
		hints[gozxing.EncodeHintType_DATA_MATRIX_SHAPE] = shape
	}
	var minSize, maxSize *gozxing.Dimension
	if options.MinWidth != 0 {
		dimension, _ := gozxing.NewDimension(options.MinWidth, options.MinHeight)
		minSize = dimension
		hints[gozxing.EncodeHintType_MIN_SIZE] = dimension
	}
	if options.MaxWidth != 0 {
		dimension, _ := gozxing.NewDimension(options.MaxWidth, options.MaxHeight)
		maxSize = dimension
		hints[gozxing.EncodeHintType_MAX_SIZE] = dimension
	}
	var matrix *gozxing.BitMatrix
	var err error
	if options.GS1 || options.ECI != 0 || options.Macro != MacroNone ||
		options.StructuredAppend != nil {
		matrix, err = encodeControlled(payload, options, shape, minSize, maxSize)
	} else {
		matrix, err = zxingdatamatrix.NewDataMatrixWriter().Encode(
			string(payload), gozxing.BarcodeFormat_DATA_MATRIX, 1, 1, hints,
		)
	}
	if err != nil {
		return Symbol{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}
	logicalMatrix := copyMatrix(matrix, quietZone)
	logical, _ := barcode.NewSymbol(barcode.SymbolOptions{
		Format: barcode.DataMatrix, Payload: payload, Matrix: logicalMatrix,
	})

	var structuredAppend *StructuredAppend
	if options.StructuredAppend != nil {
		headerCopy := *options.StructuredAppend
		structuredAppend = &headerCopy
	}
	return Symbol{
		logical: logical, shape: options.Shape, gs1: options.GS1,
		eci: options.ECI, macro: options.Macro, structuredAppend: structuredAppend,
	}, nil
}

// EncodeGS1 serializes validated GS1 elements and emits first-position FNC1.
func EncodeGS1(elements gs1.ElementString, options Options) (Symbol, error) {
	options.GS1 = true

	return Encode([]byte(elements.Raw()), options)
}

func encodeControlled(
	payload []byte,
	options Options,
	shape encoder.SymbolShapeHint,
	minSize, maxSize *gozxing.Dimension,
) (*gozxing.BitMatrix, error) {
	codewords := make([]byte, 0, len(payload)+8)
	if options.StructuredAppend != nil {
		codewords = append(
			codewords,
			233,
			structuredAppendSequence(
				options.StructuredAppend.Index,
				options.StructuredAppend.Total,
			),
			options.StructuredAppend.FileID,
		)
	}
	switch options.Macro {
	case MacroNone:
	case Macro05:
		codewords = append(codewords, encoder.HighLevelEncoder_MACRO_05)
	case Macro06:
		codewords = append(codewords, encoder.HighLevelEncoder_MACRO_06)
	}
	if options.GS1 {
		codewords = append(codewords, 232)
	}
	if options.ECI != 0 {
		codewords = append(codewords, 241)
		codewords = append(codewords, eciCodewords(options.ECI)...)
	}
	codewords = append(codewords, encoder.HighLevelEncoder_LATCH_TO_BASE256)
	length := len(payload)
	if length <= 249 {
		codewords = appendBase256(codewords, byte(length))
	} else {
		if length > 1555 {
			return nil, ErrInvalidInput
		}
		codewords = appendBase256(codewords, byte(length/250+249))
		codewords = appendBase256(codewords, byte(length%250))
	}
	for _, value := range payload {
		codewords = appendBase256(codewords, value)
	}

	info, err := encoder.SymbolInfo_Lookup(len(codewords), shape, minSize, maxSize, true)
	if err != nil {
		return nil, err
	}
	if len(codewords) < info.GetDataCapacity() {
		codewords = append(codewords, encoder.HighLevelEncoder_PAD)
	}
	for len(codewords) < info.GetDataCapacity() {
		codewords = append(codewords, randomize253(len(codewords)+1))
	}
	withCorrection, _ := encoder.ErrorCorrection_EncodeECC200(codewords, info)
	placement := encoder.NewDefaultPlacement(
		withCorrection, info.GetSymbolDataWidth(), info.GetSymbolDataHeight(),
	)
	placement.Place()

	return placementMatrix(placement, info), nil
}

func structuredAppendSequence(index, total int) byte {
	// #nosec G115 -- validation bounds the expression to one codeword.
	return byte(index*16 + 17 - total)
}

// Base256 randomization depends on the absolute one-based codeword position,
// including control prefixes and length bytes. Resetting the position for payload
// would change the wire symbol; datamatrix_internal_test.go protects this mapping.
func appendBase256(codewords []byte, value byte) []byte {
	position := len(codewords) + 1
	pseudoRandom := (149*position)%255 + 1

	// #nosec G115 -- conversion implements the ECC 200 modulo-256 sum.
	return append(codewords, byte(int(value)+pseudoRandom))
}

func randomize253(position int) byte {
	pseudoRandom := (149*position)%253 + 1
	randomized := encoder.HighLevelEncoder_PAD + pseudoRandom
	if randomized > 254 {
		randomized -= 254
	}

	// #nosec G115 -- modulo arithmetic bounds randomized to one byte.
	return byte(randomized)
}

func eciCodewords(assignment int) []byte {
	switch {
	case assignment <= 126:
		// #nosec G115 -- the branch bounds assignment+1 to one byte.
		return []byte{byte(assignment + 1)}
	case assignment <= 16_382:
		value := assignment - 127
		return []byte{byte(value/254) | 0x80, byte(value%254 + 1)}
	default:
		value := assignment - 16_383
		// #nosec G115 -- validation bounds every ECI component to one byte.
		return []byte{
			byte(value/64_516 + 192),
			byte(value/254%254 + 1),
			byte(value%254 + 1),
		}
	}
}

func placementMatrix(
	placement *encoder.DefaultPlacement,
	info *encoder.SymbolInfo,
) *gozxing.BitMatrix {
	matrix, _ := gozxing.NewBitMatrix(info.GetSymbolWidth(), info.GetSymbolHeight())
	matrixY := 0
	for y := 0; y < info.GetSymbolDataHeight(); y++ {
		matrixX := 0
		if y%info.GetMatrixHeight() == 0 {
			for x := range info.GetSymbolWidth() {
				if x%2 == 0 {
					matrix.Set(matrixX, matrixY)
				}
				matrixX++
			}
			matrixY++
		}
		matrixX = 0
		for x := 0; x < info.GetSymbolDataWidth(); x++ {
			if x%info.GetMatrixWidth() == 0 {
				matrix.Set(matrixX, matrixY)
				matrixX++
			}
			if placement.GetBit(x, y) {
				matrix.Set(matrixX, matrixY)
			}
			matrixX++
			if x%info.GetMatrixWidth() == info.GetMatrixWidth()-1 {
				if y%2 == 0 {
					matrix.Set(matrixX, matrixY)
				}
				matrixX++
			}
		}
		matrixY++
		if y%info.GetMatrixHeight() == info.GetMatrixHeight()-1 {
			for x := 0; x < info.GetSymbolWidth(); x++ {
				matrix.Set(x, matrixY)
			}
			matrixY++
		}
	}

	return matrix
}

func copyMatrix(source *gozxing.BitMatrix, quietZone int) barcode.Matrix {
	width, height := source.GetWidth()+2*quietZone, source.GetHeight()+2*quietZone
	modules := make([]bool, width*height)
	for y := range source.GetHeight() {
		for x := range source.GetWidth() {
			modules[(y+quietZone)*width+x+quietZone] = source.Get(x, y)
		}
	}
	matrix, _ := barcode.NewMatrix(width, height, modules)

	return matrix
}
