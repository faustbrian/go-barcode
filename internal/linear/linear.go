// Package linear adapts selected ZXing writers to alternating logical runs.
// It is an internal trusted-wrapper seam: writers/hints are borrowed, no input
// or dimension admission is supplied, and no symbology/device qualification is
// implied. Hints are mutated to zero dependency margin so caller quiet zones
// can be added exactly once. See TestEncodeMergesTrailingLightRunWithQuietZone.
package linear

import (
	"fmt"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/makiuchi-d/gozxing"
)

// Encode returns exact logical runs with asymmetric quiet zones added after
// the dependency has produced its unscaled core modules. The provided hints map
// is mutated (MARGIN=0), not copied. writer must return a nonempty matrix and
// dimensions/margins must already be bounded; callbacks can block or panic.
// It returns writer errors unchanged and wraps bar-shape errors. No cancellation
// or protection against margin-sum overflow is added by this trusted seam.
func Encode(
	writer gozxing.Writer,
	format gozxing.BarcodeFormat,
	contents string,
	height int,
	quietLeft int,
	quietRight int,
	hints map[gozxing.EncodeHintType]interface{},
) (barcode.Bars, error) {
	if hints == nil {
		hints = make(map[gozxing.EncodeHintType]interface{})
	}
	hints[gozxing.EncodeHintType_MARGIN] = 0
	matrix, err := writer.Encode(contents, format, 0, 1, hints)
	if err != nil {
		return barcode.Bars{}, err
	}

	var runs []barcode.Bar
	runs = append(runs, barcode.Bar{Width: quietLeft})
	dark := matrix.Get(0, 0)
	width := 1
	for offset := range matrix.GetWidth() - 1 {
		x := offset + 1
		if matrix.Get(x, 0) == dark {
			width++
			continue
		}
		runs = append(runs, barcode.Bar{Dark: dark, Width: width})
		dark = !dark
		width = 1
	}
	runs = append(runs, barcode.Bar{Dark: dark, Width: width})
	if dark {
		runs = append(runs, barcode.Bar{Width: quietRight})
	} else {
		runs[len(runs)-1].Width += quietRight
	}

	bars, err := barcode.NewBars(height, runs)
	if err != nil {
		return barcode.Bars{}, fmt.Errorf("logical bars: %w", err)
	}

	return bars, nil
}
