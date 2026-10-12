// Package code128 constructs logical barcodes for the pinned profile
// ISO/IEC 15417:2007 edition 2 + Amendment 1:2026; GS1 release 26.0.
// Encode accepts 1..80 ASCII bytes and adds the mandatory modulo-103 checksum.
// RawFNC1 is a legacy unvalidated-GS1 mode; see BARCODE-DEC-007/008.
// Output owns copied payload/geometry; calls are synchronous and borrow input
// only during encoding. It does not own rendering, scan devices, allocation
// policy for caller-sized margins or content authenticity. Exact identities and
// evidence are in specification/manifest.json and specification/evidence.tsv.
package code128

import (
	"errors"
	"fmt"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/gs1"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

const (
	defaultQuietZone = 10
	defaultHeight    = 50
	maxHeight        = 4096
	maxPayload       = 80
)

// ErrInvalidInput reports a payload or option that cannot produce a
// selected Code 128 software profile; it is not a physical standards certificate.
var ErrInvalidInput = errors.New("code128: invalid input")

// CodeSet selects an ISO/IEC 15417 character set. Auto switches sets to
// follow the selected writer's switching policy; this is not a global optimality
// or performance guarantee.
type CodeSet uint8

const (
	// CodeSetAuto and the other constants select Code 128 character sets.
	CodeSetAuto CodeSet = iota
	// CodeSetA forces Code Set A.
	CodeSetA
	// CodeSetB forces Code Set B.
	CodeSetB
	// CodeSetC forces numeric-pair Code Set C.
	CodeSetC
)

// Options controls logical Code 128 encoding in module units.
type Options struct {
	// CodeSet selects automatic switching or a forced A/B/C repertoire;
	// invalid bytes/pairs in a forced set can fail at the selected writer.
	CodeSet CodeSet
	// GS1 requires a valid raw GS1 element string and prepends FNC1.
	// It conflicts with RawFNC1; application identifiers are not authenticated.
	GS1 bool
	// RawFNC1 prepends a leading FNC1 without validating the payload as a GS1
	// element string. It exists only for compatibility with legacy symbols;
	// new GS1-128 integrations should use GS1 or EncodeGS1.
	RawFNC1 bool
	// QuietZone is the per-side logical margin; zero selects 10 and smaller
	// positive values are rejected. There is no absolute positive upper bound here;
	// callers must bound dimensions before encoding untrusted requests.
	QuietZone int
	// Height is the logical bar height; zero selects 50, negative or above 4096 is rejected.
	Height int
}

// Encode validates payload and options and returns an immutable logical
// symbol. QuietZone defaults to ten modules on each side and cannot be made
// smaller.
func Encode(payload []byte, options Options) (barcode.Symbol, error) {
	if len(payload) == 0 || len(payload) > maxPayload || options.CodeSet > CodeSetC ||
		options.GS1 && options.RawFNC1 ||
		options.QuietZone < 0 || options.Height < 0 {
		return barcode.Symbol{}, ErrInvalidInput
	}
	for _, value := range payload {
		if value > 127 {
			return barcode.Symbol{}, ErrInvalidInput
		}
	}
	if options.GS1 {
		if _, err := gs1.ParseRaw(string(payload), gs1.ParseLimits{MaxInputBytes: maxPayload}); err != nil {
			return barcode.Symbol{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
		}
	}

	quietZone := options.QuietZone
	if quietZone == 0 {
		quietZone = defaultQuietZone
	}
	if quietZone < defaultQuietZone {
		return barcode.Symbol{}, ErrInvalidInput
	}
	height := options.Height
	if height == 0 {
		height = defaultHeight
	}
	if height > maxHeight {
		return barcode.Symbol{}, ErrInvalidInput
	}

	contents := string(payload)
	if options.GS1 || options.RawFNC1 {
		contents = "\u00f1" + contents
	}
	hints := map[gozxing.EncodeHintType]interface{}{
		gozxing.EncodeHintType_MARGIN: quietZone * 2,
	}
	if options.CodeSet != CodeSetAuto {
		hints[gozxing.EncodeHintType_FORCE_CODE_SET] = codeSetName(options.CodeSet)
	}

	matrix, err := oned.NewCode128Writer().Encode(
		contents,
		gozxing.BarcodeFormat_CODE_128,
		0,
		height,
		hints,
	)
	if err != nil {
		return barcode.Symbol{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}

	runs := matrixRuns(matrix)
	bars, _ := barcode.NewBars(height, runs)
	format := barcode.Code128
	if options.GS1 || options.RawFNC1 {
		format = barcode.GS1128
	}

	symbol, _ := barcode.NewSymbol(barcode.SymbolOptions{
		Format:  format,
		Payload: payload,
		Bars:    bars,
	})

	return symbol, nil
}

// EncodeGS1 serializes validated structured elements and inserts the required
// leading FNC1 marker.
func EncodeGS1(elements gs1.ElementString, options Options) (barcode.Symbol, error) {
	options.GS1 = true

	return Encode([]byte(elements.Raw()), options)
}

func codeSetName(codeSet CodeSet) string {
	switch codeSet {
	case CodeSetA:
		return "A"
	case CodeSetB:
		return "B"
	case CodeSetC:
		return "C"
	case CodeSetAuto:
		return ""
	default:
		return ""
	}
}

func matrixRuns(matrix *gozxing.BitMatrix) []barcode.Bar {
	// Count transitions so quiet-zone whitespace cannot inflate run storage.
	count := 1
	previous := matrix.Get(0, 0)
	for x := range matrix.GetWidth() {
		dark := matrix.Get(x, 0)
		if dark != previous {
			count++
			previous = dark
		}
	}
	runs := make([]barcode.Bar, 0, count)
	dark := matrix.Get(0, 0)
	width := 1
	for x := range matrix.GetWidth() {
		if x == 0 {
			continue
		}
		if matrix.Get(x, 0) == dark {
			width++
			continue
		}
		runs = append(runs, barcode.Bar{Dark: dark, Width: width})
		dark = !dark
		width = 1
	}

	return append(runs, barcode.Bar{Dark: dark, Width: width})
}
