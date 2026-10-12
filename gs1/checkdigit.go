// Package gs1 parses GS1 release 26.0/January 2026 element syntax and bounded
// modulo-10 numeric keys. ParseRaw/ParseBracketed use the embedded dictionary
// release 2026-01-27, loaded once and then read-only; Elements returns copied
// entries but strings can retain the original input's backing storage.
//
// Syntax/linter/check-digit validation is not identity authentication or a live
// company-prefix registry lookup. No barcode control codewords, transport,
// device or cancellation lifecycle are owned here. See BARCODE-DEC-008,
// specification/README.md and NOTICE.md for the exact dictionary/reference pins.
package gs1

import "errors"

// ErrInvalidCheckDigit reports malformed input or a check-digit mismatch.
var ErrInvalidCheckDigit = errors.New("gs1: invalid check digit")

// CalculateCheckDigit returns the GS1 modulo-10 check digit for a numeric key
// body. GS1 keys contain at most 18 digits including the check digit.
func CalculateCheckDigit(body string) (byte, error) {
	if len(body) == 0 || len(body) > 17 {
		return 0, ErrInvalidCheckDigit
	}

	sum := 0
	weight := 3
	for index := len(body) - 1; index >= 0; index-- {
		digit := body[index]
		if digit < '0' || digit > '9' {
			return 0, ErrInvalidCheckDigit
		}
		sum += int(digit-'0') * weight
		weight = 4 - weight
	}

	return byte('0' + (10-byte(sum%10))%10), nil
}

// ValidateCheckDigit strictly validates a complete GS1 numeric key.
func ValidateCheckDigit(value string) error {
	if len(value) < 2 || len(value) > 18 {
		return ErrInvalidCheckDigit
	}

	want, err := CalculateCheckDigit(value[:len(value)-1])
	if err != nil || value[len(value)-1] != want {
		return ErrInvalidCheckDigit
	}

	return nil
}
