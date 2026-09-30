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

func TestDecodeEncodedBytesRejectsOversizeAndMalformedInput(t *testing.T) {
	const secret = "private-payload-7f3a9d"
	for _, test := range []struct {
		name   string
		input  []byte
		limits imagedecode.Limits
		want   error
	}{
		{name: "nil", want: imagedecode.ErrInvalidImage},
		{name: "oversize", input: []byte(secret), limits: imagedecode.Limits{MaxEncodedBytes: 8}, want: imagedecode.ErrLimitExceeded},
		{name: "malformed", input: []byte(secret), want: imagedecode.ErrInvalidImage},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := imagedecode.DecodeEncoded(context.Background(), test.input, imagedecode.Options{Limits: test.limits})
			if !errors.Is(err, test.want) {
				t.Fatalf("DecodeEncoded() error = %v, want %v", err, test.want)
			}
			if strings.Contains(err.Error(), secret) {
				t.Fatalf("DecodeEncoded() exposed input: %v", err)
			}
		})
	}
}

func TestDecodeEncodedBytesHonorsPreexistingCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := imagedecode.DecodeEncoded(ctx, []byte("malformed"), imagedecode.Options{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("DecodeEncoded() error = %v, want canceled", err)
	}
}

func TestDecodeEncodedBytesStopsBetweenConfigurationAndDecode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	decodeCalled := false
	image.RegisterFormat("cancel-after-config-test", "CAC1", func(io.Reader) (image.Image, error) {
		decodeCalled = true
		return image.NewGray(image.Rect(0, 0, 1, 1)), nil
	}, func(io.Reader) (image.Config, error) {
		cancel()
		return image.Config{ColorModel: color.GrayModel, Width: 1, Height: 1}, nil
	})

	_, err := imagedecode.DecodeEncoded(ctx, []byte("CAC1"), imagedecode.Options{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("DecodeEncoded() error = %v, want canceled", err)
	}
	if decodeCalled {
		t.Fatal("DecodeEncoded() entered full image decode after cancellation")
	}
}
