package barcode_test

import (
	"testing"

	"github.com/faustbrian/go-barcode/aztec"
	"github.com/faustbrian/go-barcode/codabar"
	"github.com/faustbrian/go-barcode/code128"
	"github.com/faustbrian/go-barcode/code39"
	"github.com/faustbrian/go-barcode/code93"
	"github.com/faustbrian/go-barcode/datamatrix"
	"github.com/faustbrian/go-barcode/ean"
	"github.com/faustbrian/go-barcode/itf"
	"github.com/faustbrian/go-barcode/pdf417"
	"github.com/faustbrian/go-barcode/qr"
	"github.com/faustbrian/go-barcode/upc"
)

func FuzzPayloadEncoders(f *testing.F) {
	f.Add([]byte("ABC123"))
	f.Add([]byte("0123456789"))
	f.Add([]byte{0, 29, 127, 255})

	f.Fuzz(func(t *testing.T, payload []byte) {
		if len(payload) == 0 || len(payload) > 128 {
			t.Skip()
		}
		_, _ = qr.Encode(payload, qr.Options{})
		_, _ = code128.Encode(payload, code128.Options{})
		_, _ = code39.Encode(payload, code39.Options{})
		_, _ = code93.Encode(payload, code93.Options{})
		_, _ = codabar.Encode(payload, codabar.Options{})
		_, _ = itf.Encode(string(payload), itf.Options{})
		_, _ = datamatrix.Encode(payload, datamatrix.Options{})
		_, _ = pdf417.Encode(payload, pdf417.Options{})
		_, _ = aztec.Encode(payload, aztec.Options{})
		_, _ = ean.Encode8(string(payload), ean.Options{})
		_, _ = ean.Encode13(string(payload), ean.Options{})
		_, _ = upc.EncodeA(string(payload), upc.Options{})
		_, _ = upc.EncodeE(string(payload), upc.Options{})
	})
}
