// Package aztec constructs copied logical Aztec symbols via Encode/EncodeGS1.
// It owns no decoder/device or sequence-assembly lifecycle. Explicit correction,
// compact/full layers, GS1 and ECI target ISO/IEC 24778:2008 edition 1, not its 2024
// replacement or hardware certification (BARCODE-DEC-002/011). NOTICE.md records
// the maintained local zxinggo v0.1.0 encoder derivation.
package aztec

import (
	"errors"
	"fmt"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/gs1"
	aztecencoder "github.com/faustbrian/go-barcode/v2/internal/aztecencoder"
)

const (
	defaultErrorCorrectionPercent = 33
	defaultQuietZone              = 1
	maxPayloadBytes               = 4096
)

// ErrInvalidInput reports data or options that cannot form an Aztec symbol.
var ErrInvalidInput = errors.New("aztec: invalid input")

// Options controls Aztec compact mode, layers, correction, and quiet zone.
type Options struct {
	// Compact restricts output to compact layers; false allows auto compact/full choice.
	Compact bool
	// Layers is zero for automatic selection; positive forces 1..32 full or 1..4 compact.
	Layers int
	// ErrorCorrectionPercent is zero for 33; positive values 1..100 are requested
	// parity sizing, not a guarantee of recoverable image-damage percentage.
	ErrorCorrectionPercent int
	// QuietZone is zero for 1 module, otherwise 1..256; included in logical output.
	QuietZone int
	// GS1 validates raw GS1 elements and emits FLG(0); no authenticity proof is supplied.
	GS1 bool
	// ECI is zero for absence, otherwise 1..999999; a control prefix, not transcoding.
	ECI int
}

// Symbol combines an immutable logical symbol with Aztec metadata.
type Symbol struct {
	logical                barcode.Symbol
	compact                bool
	layers                 int
	errorCorrectionPercent int
	gs1                    bool
	eci                    int
}

// Logical returns the immutable format-neutral symbol.
func (symbol Symbol) Logical() barcode.Symbol { return symbol.logical }

// Compact reports whether the compact Aztec form was selected.
func (symbol Symbol) Compact() bool { return symbol.compact }

// Layers returns the encoded Aztec layer count.
func (symbol Symbol) Layers() int { return symbol.layers }

// ErrorCorrectionPercent returns the effective requested percentage (33 by default).
func (symbol Symbol) ErrorCorrectionPercent() int { return symbol.errorCorrectionPercent }

// GS1 reports whether the symbol begins with the FLG(0) FNC1 control.
func (symbol Symbol) GS1() bool { return symbol.gs1 }

// ECI returns the requested ECI assignment, or zero when absent.
func (symbol Symbol) ECI() int { return symbol.eci }

// Encode returns an Aztec matrix. Positive forced layers select a full symbol;
// Compact with forced layers selects one of the four compact sizes. Payload must
// be nonempty and at most 4096 bytes; actual capacity can be lower after word
// stuffing/parity. Input is borrowed synchronously, output copies payload/modules.
// Errors wrap ErrInvalidInput and may join the selected encoder or GS1 cause.
// This is the pinned 2008 software profile, not automatic 2024-edition adoption.
func Encode(payload []byte, options Options) (Symbol, error) {
	if len(payload) == 0 || len(payload) > maxPayloadBytes || options.Layers < 0 ||
		options.Layers > 32 || options.Compact && options.Layers > 4 ||
		options.ErrorCorrectionPercent < 0 || options.ErrorCorrectionPercent > 100 ||
		options.QuietZone < 0 || options.ECI < 0 || options.ECI >= 1_000_000 {
		return Symbol{}, ErrInvalidInput
	}
	if options.GS1 {
		if _, err := gs1.ParseRaw(string(payload), gs1.ParseLimits{MaxInputBytes: maxPayloadBytes}); err != nil {
			return Symbol{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
		}
	}
	errorCorrection := options.ErrorCorrectionPercent
	if errorCorrection == 0 {
		errorCorrection = defaultErrorCorrectionPercent
	}
	quietZone := options.QuietZone
	if quietZone == 0 {
		quietZone = defaultQuietZone
	}
	if quietZone > 256 {
		return Symbol{}, ErrInvalidInput
	}
	forcedLayers := options.Layers
	if options.Compact {
		forcedLayers = -forcedLayers
	}
	code, err := aztecencoder.EncodeWithControls(
		payload, errorCorrection, forcedLayers, options.GS1, options.ECI,
	)
	if options.Compact {
		switch options.Layers {
		case 0:
			for layer := 1; layer <= 4; layer++ {
				code, err = aztecencoder.EncodeWithControls(
					payload, errorCorrection, -layer, options.GS1, options.ECI,
				)
				if err == nil {
					break
				}
			}
		}
	}
	if err != nil {
		return Symbol{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}
	width := code.Matrix.Width() + 2*quietZone
	modules := make([]bool, width*width)
	for y := range code.Matrix.Height() {
		for x := range code.Matrix.Width() {
			modules[(y+quietZone)*width+x+quietZone] = code.Matrix.Get(x, y)
		}
	}
	matrix, _ := barcode.NewMatrix(width, width, modules)
	logical, _ := barcode.NewSymbol(barcode.SymbolOptions{
		Format: barcode.Aztec, Payload: payload, Matrix: matrix,
	})

	return Symbol{
		logical: logical, compact: code.Compact, layers: code.Layers,
		errorCorrectionPercent: errorCorrection, gs1: options.GS1, eci: options.ECI,
	}, nil
}

// EncodeGS1 serializes validated GS1 elements and emits FLG(0) FNC1.
func EncodeGS1(elements gs1.ElementString, options Options) (Symbol, error) {
	options.GS1 = true

	return Encode([]byte(elements.Raw()), options)
}
