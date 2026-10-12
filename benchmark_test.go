package barcode_test

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"testing"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/code128"
	"github.com/faustbrian/go-barcode/v2/imagedecode"
	"github.com/faustbrian/go-barcode/v2/qr"
	"github.com/faustbrian/go-barcode/v2/render"
)

// BenchmarkEncodeQR measures selected logical QR payloads under fixed correction;
// compare identical payload/options/toolchains, not different symbology capacities.
func BenchmarkEncodeQR(b *testing.B) {
	for _, size := range []int{16, 128, 1024} {
		b.Run(fmt.Sprintf("%d_bytes", size), func(b *testing.B) {
			payload := bytes.Repeat([]byte{'A'}, size)
			b.ReportAllocs()
			b.SetBytes(int64(size))
			for b.Loop() {
				if _, err := qr.Encode(payload, qr.Options{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkEncodeCode128 measures one fixed logical payload, not rendering or
// image acquisition. Compare the same selected dependency/toolchain profile.
func BenchmarkEncodeCode128(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := code128.Encode([]byte("ORDER-0123456789"), code128.Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRenderQR excludes symbol setup and measures raster creation at fixed
// integer scale; it does not measure PNG I/O or scanner readability.
func BenchmarkRenderQR(b *testing.B) {
	symbol, err := qr.Encode([]byte("render benchmark"), qr.Options{})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := render.Image(symbol.Logical(), render.Options{Scale: 4}); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDecodeQR excludes source encode/render setup and measures repeated
// image decoding of the same fixed QR geometry and selected format.
func BenchmarkDecodeQR(b *testing.B) {
	symbol, err := qr.Encode([]byte("decode benchmark"), qr.Options{})
	if err != nil {
		b.Fatal(err)
	}
	input, err := render.Image(symbol.Logical(), render.Options{Scale: 4})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := imagedecode.Decode(context.Background(), input, imagedecode.Options{
			Formats: []barcode.Format{barcode.QRCode},
		}); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkDetectNoSymbol measures failed detection on the fixed no-symbol
// image, not all hostile candidate or resource-limit boundaries.
func BenchmarkDetectNoSymbol(b *testing.B) {
	input := image.NewGray(image.Rect(0, 0, 64, 64))
	b.ReportAllocs()
	for b.Loop() {
		_, _ = imagedecode.Decode(context.Background(), input, imagedecode.Options{
			Formats: []barcode.Format{barcode.QRCode},
		})
	}
}

// BenchmarkDecodeEncodedQR excludes fixture encode/render/PNG setup and measures
// full byte-to-image-to-symbol decode. External time sampling includes process
// overhead; it is not a library-only peak-allocation guarantee.
func BenchmarkDecodeEncodedQR(b *testing.B) {
	symbol, err := qr.Encode([]byte("encoded image decode benchmark"), qr.Options{})
	if err != nil {
		b.Fatal(err)
	}
	var encoded bytes.Buffer
	if err := render.PNG(&encoded, symbol.Logical(), render.Options{Scale: 4}); err != nil {
		b.Fatal(err)
	}
	input := encoded.Bytes()
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	for b.Loop() {
		if _, err := imagedecode.DecodeEncoded(
			context.Background(),
			input,
			imagedecode.Options{Formats: []barcode.Format{barcode.QRCode}},
		); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRejectMalformedImage measures rejection of one fixed malformed byte
// fixture, not every codec failure or worst-case amplification pattern.
func BenchmarkRejectMalformedImage(b *testing.B) {
	input := bytes.Repeat([]byte{0xff}, 4096)
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	for b.Loop() {
		_, _ = imagedecode.DecodeEncoded(
			context.Background(),
			input,
			imagedecode.Options{},
		)
	}
}
