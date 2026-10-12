// Package code93 constructs logical barcodes for the pinned profile
// ANSI/AIM BC5-1995 (historical withdrawn target).
// Encode admits 1..80 ASCII bytes before the selected full-ASCII writer adds
// C/K checks. The local wrapper discards that writer error, so raw length alone
// does not prove successful expansion/capacity. The 2000 AIM publication is a
// separate review input, not this target (DEC-001/007).
// Output owns copied payload/geometry; calls are synchronous and borrow input
// only during encoding. It does not own rendering, scan devices, allocation
// policy for caller-sized margins or content authenticity. Exact identities and
// evidence are in specification/manifest.json and specification/evidence.tsv.
package code93

import (
	"errors"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/internal/linear"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

const (
	defaultQuietZone = 10
	defaultHeight    = 50
	maxHeight        = 4096
	maxPayloadBytes  = 80
)

// ErrInvalidInput reports data or options that cannot form a Code 93 symbol.
var ErrInvalidInput = errors.New("code93: invalid input")

// Options controls Code 93 quiet zone and bar height in modules.
type Options struct {
	// QuietZone is the per-side logical margin; zero selects 10 and smaller
	// positive values are rejected. There is no absolute positive upper bound here;
	// callers must bound dimensions before encoding untrusted requests.
	QuietZone int
	// Height is the logical bar height; zero selects 50, negative or above 4096 is rejected.
	Height int
}

// Encode requests a full-ASCII Code 93 symbol with mandatory C/K checks from
// the selected writer. It validates raw payload/options first, but currently
// discards writer and core-constructor errors; nil error alone does not prove a
// nonempty symbol if expansion exceeds peer capacity. Bound margins before use.
func Encode(payload []byte, options Options) (barcode.Symbol, error) {
	if len(payload) == 0 || len(payload) > maxPayloadBytes ||
		options.QuietZone < 0 || options.Height < 0 {
		return barcode.Symbol{}, ErrInvalidInput
	}
	for _, value := range payload {
		if value > 127 {
			return barcode.Symbol{}, ErrInvalidInput
		}
	}
	quietZone := options.QuietZone
	if quietZone == 0 {
		quietZone = defaultQuietZone
	}
	height := options.Height
	if height == 0 {
		height = defaultHeight
	}
	if quietZone < defaultQuietZone || height > maxHeight {
		return barcode.Symbol{}, ErrInvalidInput
	}

	bars, _ := linear.Encode(
		oned.NewCode93Writer(),
		gozxing.BarcodeFormat_CODE_93,
		string(payload),
		height,
		quietZone,
		quietZone,
		nil,
	)
	symbol, _ := barcode.NewSymbol(barcode.SymbolOptions{
		Format: barcode.Code93, Payload: payload, Bars: bars,
	})

	return symbol, nil
}
