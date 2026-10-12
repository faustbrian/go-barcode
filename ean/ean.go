// Package ean constructs logical barcodes for the pinned profile
// ISO/IEC 15420:2009 edition 2.
// Encode8/Encode13 calculate missing modulo-10 or validate complete digits.
// Default margins are 7/7 (EAN8) and 11/7 (EAN13). No 2025-edition claim (DEC-002/007).
// Output owns copied payload/geometry; calls are synchronous and borrow input
// only during encoding. It does not own rendering, scan devices, allocation
// policy for caller-sized margins or content authenticity. Exact identities and
// evidence are in specification/manifest.json and specification/evidence.tsv.
package ean

import (
	"errors"
	"strings"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/gs1"
)

const (
	defaultHeight = 50
	maxHeight     = 4096
)

// ErrInvalidInput reports invalid EAN payloads and options.
var ErrInvalidInput = errors.New("ean: invalid input")

// Options controls symbol dimensions in module units and an optional EAN-2
// or EAN-5 supplement.
type Options struct {
	// Supplement optionally carries exactly 2 or 5 ASCII digits; it is encoded
	// beside the main symbol and is not appended to Symbol.Payload.
	Supplement string
	// QuietZoneLeft is a logical margin; zero selects the format-specific minimum.
	// Negative/below-minimum values fail; no absolute positive upper bound is enforced.
	QuietZoneLeft int
	// QuietZoneRight follows the corresponding right-margin minimum/default.
	// Callers must cap positive margins to bound temporary strings and allocations.
	QuietZoneRight int
	// Height is the logical bar height; zero selects 50, negative or above 4096 is rejected.
	Height int
}

var (
	leftPatterns = [...]string{
		"0001101", "0011001", "0010011", "0111101", "0100011",
		"0110001", "0101111", "0111011", "0110111", "0001011",
	}
	guardedPatterns = [...]string{
		"0100111", "0110011", "0011011", "0100001", "0011101",
		"0111001", "0000101", "0010001", "0001001", "0010111",
	}
	rightPatterns = [...]string{
		"1110010", "1100110", "1101100", "1000010", "1011100",
		"1001110", "1010000", "1000100", "1001000", "1110100",
	}
	ean13Parity = [...]string{
		"LLLLLL", "LLGLGG", "LLGGLG", "LLGGGL", "LGLLGG",
		"LGGLLG", "LGGGLL", "LGLGLG", "LGLGGL", "LGGLGL",
	}
	ean5Parity = [...]string{
		"GGLLL", "GLGLL", "GLLGL", "GLLLG", "LGGLL",
		"LLGGL", "LLLGG", "LGLGL", "LGLLG", "LLGLG",
	}
)

// Encode8 accepts a seven-digit body or a complete eight-digit EAN-8 value.
func Encode8(value string, options Options) (barcode.Symbol, error) {
	complete, err := completeValue(value, 7)
	if err != nil {
		return barcode.Symbol{}, err
	}

	core := "101"
	for index := 0; index < 4; index++ {
		core += leftPatterns[complete[index]-'0']
	}
	core += "01010"
	for index := 4; index < 8; index++ {
		core += rightPatterns[complete[index]-'0']
	}
	core += "101"

	return makeSymbol(barcode.EAN8, complete, core, 7, 7, options)
}

// Encode13 accepts a twelve-digit body or a complete thirteen-digit EAN-13
// value.
func Encode13(value string, options Options) (barcode.Symbol, error) {
	complete, err := completeValue(value, 12)
	if err != nil {
		return barcode.Symbol{}, err
	}

	parity := ean13Parity[complete[0]-'0']
	core := "101"
	for index := 1; index <= 6; index++ {
		patterns := leftPatterns[:]
		if parity[index-1] == 'G' {
			patterns = guardedPatterns[:]
		}
		core += patterns[complete[index]-'0']
	}
	core += "01010"
	for index := 7; index < 13; index++ {
		core += rightPatterns[complete[index]-'0']
	}
	core += "101"

	return makeSymbol(barcode.EAN13, complete, core, 11, 7, options)
}

func completeValue(value string, bodyLength int) (string, error) {
	if len(value) == bodyLength {
		check, err := gs1.CalculateCheckDigit(value)
		if err != nil {
			return "", ErrInvalidInput
		}
		return value + string(check), nil
	}
	if len(value) != bodyLength+1 || gs1.ValidateCheckDigit(value) != nil {
		return "", ErrInvalidInput
	}

	return value, nil
}

func makeSymbol(
	format barcode.Format,
	payload string,
	core string,
	minimumLeft int,
	minimumRight int,
	options Options,
) (barcode.Symbol, error) {
	if options.QuietZoneLeft < 0 {
		return barcode.Symbol{}, ErrInvalidInput
	}
	if options.QuietZoneRight < 0 {
		return barcode.Symbol{}, ErrInvalidInput
	}
	if options.Height < 0 {
		return barcode.Symbol{}, ErrInvalidInput
	}
	left := options.QuietZoneLeft
	if left == 0 {
		left = minimumLeft
	}
	right := options.QuietZoneRight
	if right == 0 {
		right = minimumRight
	}
	if left < minimumLeft || right < minimumRight {
		return barcode.Symbol{}, ErrInvalidInput
	}
	height := options.Height
	if height == 0 {
		height = defaultHeight
	}
	if height > maxHeight {
		return barcode.Symbol{}, ErrInvalidInput
	}

	modules := strings.Repeat("0", left) + core
	if options.Supplement != "" {
		supplement, err := encodeSupplement(options.Supplement)
		if err != nil {
			return barcode.Symbol{}, err
		}
		modules += strings.Repeat("0", 9) + supplement
	}
	modules += strings.Repeat("0", right)
	bars, _ := barcode.NewBars(height, runs(modules))

	symbol, _ := barcode.NewSymbol(barcode.SymbolOptions{
		Format:  format,
		Payload: []byte(payload),
		Bars:    bars,
	})

	return symbol, nil
}

func encodeSupplement(value string) (string, error) {
	if len(value) != 2 && len(value) != 5 {
		return "", ErrInvalidInput
	}
	for index := range value {
		if value[index] < '0' || value[index] > '9' {
			return "", ErrInvalidInput
		}
	}

	parity := ""
	if len(value) == 2 {
		first := strings.IndexByte("0123456789", value[0])
		second := strings.IndexByte("0123456789", value[1])
		parity = [...]string{"LL", "LG", "GL", "GG"}[(first*10+second)%4]
	} else {
		checksum := (3*int(value[0]-'0') + 9*int(value[1]-'0') +
			3*int(value[2]-'0') + 9*int(value[3]-'0') + 3*int(value[4]-'0')) % 10
		parity = ean5Parity[checksum]
	}

	result := "1011"
	for index := range value {
		if index > 0 {
			result += "01"
		}
		if parity[index] == 'G' {
			result += guardedPatterns[value[index]-'0']
		} else {
			result += leftPatterns[value[index]-'0']
		}
	}

	return result, nil
}

func runs(modules string) []barcode.Bar {
	result := make([]barcode.Bar, 0)
	dark := modules[0] == '1'
	width := 1
	for index := 1; index < len(modules); index++ {
		if (modules[index] == '1') == dark {
			width++
			continue
		}
		result = append(result, barcode.Bar{Dark: dark, Width: width})
		dark = !dark
		width = 1
	}

	return append(result, barcode.Bar{Dark: dark, Width: width})
}
