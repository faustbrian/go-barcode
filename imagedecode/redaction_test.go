package imagedecode_test

import (
	"context"
	"errors"
	"image"
	"image/color"
	"io"
	"strings"
	"testing"

	"github.com/faustbrian/go-barcode/v2/imagedecode"
)

func TestDecodeEncodedDoesNotExposeInputErrors(t *testing.T) {
	const secret = "private-barcode-payload-7f3a9d"
	sensitive := errors.New(secret)
	image.RegisterFormat("redaction-config-test", "RDC1", func(io.Reader) (image.Image, error) {
		return nil, nil
	}, func(io.Reader) (image.Config, error) {
		return image.Config{}, sensitive
	})
	image.RegisterFormat("redaction-pixels-test", "RDP1", func(io.Reader) (image.Image, error) {
		return nil, sensitive
	}, func(io.Reader) (image.Config, error) {
		return image.Config{ColorModel: color.GrayModel, Width: 1, Height: 1}, nil
	})

	for _, test := range []struct {
		name  string
		input []byte
	}{
		{name: "image configuration", input: []byte("RDC1")},
		{name: "image pixels", input: []byte("RDP1")},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := imagedecode.DecodeEncoded(context.Background(), test.input, imagedecode.Options{})
			if !errors.Is(err, imagedecode.ErrInvalidImage) {
				t.Fatalf("DecodeEncoded() error = %v, want invalid image", err)
			}
			if strings.Contains(err.Error(), secret) || errors.Is(err, sensitive) {
				t.Fatalf("DecodeEncoded() exposed input error: %v", err)
			}
		})
	}
}
