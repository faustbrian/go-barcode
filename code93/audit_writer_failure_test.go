package code93_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/code93"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

func TestCode93PreservesActualExtendedCapacityWriterFailure(t *testing.T) {
	// Forty-one lower-case ASCII bytes expand into eighty-two Code93 symbols.
	// Invoke the selected real writer to establish its failure, not a fake writer.
	payload := bytes.Repeat([]byte("a"), 41)
	_, writerErr := oned.NewCode93Writer().Encode(string(payload), gozxing.BarcodeFormat_CODE_93, 0, 1, nil)
	var original gozxing.WriterException
	if !errors.As(writerErr, &original) {
		t.Fatalf("expected actual writer capacity failure: %v", writerErr)
	}
	symbol, err := code93.Encode(payload, code93.Options{})
	var retained gozxing.WriterException
	if !errors.Is(err, code93.ErrInvalidInput) || !errors.As(err, &retained) || !reflect.DeepEqual(symbol, barcode.Symbol{}) {
		t.Fatalf("writer failure lost classification/cause or published partial symbol: %v", err)
	}
	for _, accepted := range [][]byte{
		bytes.Repeat([]byte("a"), 40),
		bytes.Repeat([]byte("A"), 80),
		append(bytes.Repeat([]byte("A"), 78), byte(127)),
	} {
		symbol, err := code93.Encode(accepted, code93.Options{})
		bars, ok := symbol.Bars()
		if err != nil || !ok || bars.Width() <= 0 || !bytes.Equal(symbol.Payload(), accepted) {
			t.Fatalf("accepted expanded-capacity boundary lost actual symbol: %v", err)
		}
	}
}
