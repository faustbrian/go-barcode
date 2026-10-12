// Package code39 constructs logical barcodes for the pinned profile
// ISO/IEC 16388:2023 edition 3.
// Encode expands full ASCII, optionally adding modulo-43; Checksum accepts
// only the base alphabet. See BARCODE-DEC-007 and character-boundary tests.
// Output owns copied payload/geometry; calls are synchronous and borrow input
// only during encoding. It does not own rendering, scan devices, allocation
// policy for caller-sized margins or content authenticity. Exact identities and
// evidence are in specification/manifest.json and specification/evidence.tsv.
package code39

import (
	"errors"
	"strings"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/internal/linear"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

const (
	alphabet         = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ-. $/+%"
	defaultQuietZone = 10
	defaultHeight    = 50
	maxHeight        = 4096
	maxEncodedLength = 80
)

// ErrInvalidInput reports invalid Code 39 payloads and options.
var ErrInvalidInput = errors.New("code39: invalid input")

// Options controls Code 39 checksum, quiet zone, and bar height.
type Options struct {
	// Checksum appends modulo-43 to full-ASCII-expanded data, not the raw payload.
	// It consumes one of the 80 encoded-character slots.
	Checksum bool
	// QuietZone is the per-side logical margin; zero selects 10 and smaller
	// positive values are rejected. There is no absolute positive upper bound here;
	// callers must bound dimensions before encoding untrusted requests.
	QuietZone int
	// Height is the logical bar height; zero selects 50, negative or above 4096 is rejected.
	Height int
}

// Checksum returns the optional Code 39 modulo-43 check character. The input
// must already use the base Code 39 character set and be nonempty. This helper
// has no length budget and accumulates into int; bound input at the application
// boundary. Encode limits the expanded content passed to it.
func Checksum(payload []byte) (byte, error) {
	if len(payload) == 0 {
		return 0, ErrInvalidInput
	}
	sum := 0
	for _, value := range payload {
		index := strings.IndexByte(alphabet, value)
		if index < 0 {
			return 0, ErrInvalidInput
		}
		sum += index
	}

	return alphabet[sum%43], nil
}

// Encode returns a full-ASCII Code 39 symbol with optional modulo-43 check
// character. Payload must be nonempty ASCII; its expanded representation must
// fit 80 characters including the optional checksum. Expansion occurs before
// length rejection, so caller input length must be bounded before this API.
func Encode(payload []byte, options Options) (barcode.Symbol, error) {
	if len(payload) == 0 || options.QuietZone < 0 || options.Height < 0 {
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

	contents, err := fullASCII(payload)
	if err != nil || len(contents) > maxEncodedLength ||
		options.Checksum && len(contents) == maxEncodedLength {
		return barcode.Symbol{}, ErrInvalidInput
	}
	if options.Checksum {
		check, _ := Checksum(contents)
		contents = append(contents, check)
	}
	bars, _ := linear.Encode(
		oned.NewCode39Writer(),
		gozxing.BarcodeFormat_CODE_39,
		string(contents),
		height,
		quietZone,
		quietZone,
		nil,
	)
	symbol, _ := barcode.NewSymbol(barcode.SymbolOptions{
		Format: barcode.Code39, Payload: payload, Bars: bars,
	})

	return symbol, nil
}

func fullASCII(payload []byte) ([]byte, error) {
	firstExtended := strings.IndexFunc(string(payload), func(character rune) bool {
		return !strings.ContainsRune(alphabet, character)
	})
	if firstExtended == -1 {
		return append([]byte(nil), payload...), nil
	}

	encoded := make([]byte, 0)
	for _, character := range payload {
		switch character {
		case 0:
			encoded = append(encoded, '%', 'U')
		case ' ', '-', '.':
			encoded = append(encoded, character)
		case '@':
			encoded = append(encoded, '%', 'V')
		case '`':
			encoded = append(encoded, '%', 'W')
		default:
			switch {
			case character <= 26:
				encoded = append(encoded, '$', 'A'+character-1)
			case character <= 31:
				encoded = append(encoded, '%', 'A'+character-27)
			case character <= ',' || character == '/' || character == ':':
				encoded = append(encoded, '/', 'A'+character-33)
			case character <= '9':
				encoded = append(encoded, character)
			case character <= '?':
				encoded = append(encoded, '%', 'F'+character-59)
			case character <= 'Z':
				encoded = append(encoded, character)
			case character <= '_':
				encoded = append(encoded, '%', 'K'+character-91)
			case character <= 'z':
				encoded = append(encoded, '+', 'A'+character-97)
			case character <= 127:
				encoded = append(encoded, '%', 'P'+character-123)
			default:
				return nil, ErrInvalidInput
			}
		}
	}

	return encoded, nil
}
