// Copyright 2011 ZXing authors. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Ported from Java ZXing library.

// Package pdf417encoder implements ISO/IEC 15438 PDF417 encoding. It derives
// from zxinggo v0.1.0 under Apache-2.0 with locally maintained correctness fixes
// (NOTICE.md). The wrapper targets ISO/IEC 15438:2015; historical 2001 section
// references describe port provenance, not independently reviewed 2015 clauses.
// Builders are mutable and not concurrency-safe; validated public pdf417 callers
// own admission. This internal API does not cap arbitrary dimensions or I/O.
package pdf417encoder

// BarcodeRow is an owned mutable row builder. Coordinates/appended widths must
// fit its allocated capacity; no growth, bounds normalization or synchronization
// is supplied. Invalid dimensions/indexes can panic at this trusted internal seam.
type BarcodeRow struct {
	row             []byte
	currentLocation int
}

// NewBarcodeRow allocates width bytes; callers must supply a bounded nonnegative
// width. Negative width panics, and there is no independent allocation budget.
func NewBarcodeRow(width int) *BarcodeRow {
	return &BarcodeRow{
		row:             make([]byte, width),
		currentLocation: 0,
	}
}

// Set sets a specific location in the bar.
func (br *BarcodeRow) Set(x int, value byte) {
	br.row[x] = value
}

// AddBar adds a bar (black or white) of the given width at the current location.
func (br *BarcodeRow) AddBar(black bool, width int) {
	var value byte
	if black {
		value = 1
	}
	for i := 0; i < width; i++ {
		br.row[br.currentLocation] = value
		br.currentLocation++
	}
}

// GetScaledRow allocates a copied row with integer x scaling. Callers must bound
// the positive scale/product; zero yields empty storage and overflow/negative
// allocation or invalid inputs can panic. The original row remains mutable.
func (br *BarcodeRow) GetScaledRow(scale int) []byte {
	output := make([]byte, len(br.row)*scale)
	for i := range output {
		output[i] = br.row[i/scale]
	}
	return output
}

// BarcodeMatrix holds all of the information for a barcode in a format where
// it can be easily accessible. Rows and cursor are mutable, unsynchronized state;
// obtain it only through a bounded caller/encoder, not arbitrary hostile dimensions.
type BarcodeMatrix struct {
	matrix     []*BarcodeRow
	currentRow int
	height     int
}

// NewBarcodeMatrix creates a new BarcodeMatrix with the given height (rows)
// and width (cols).
func NewBarcodeMatrix(height, width int) *BarcodeMatrix {
	m := &BarcodeMatrix{
		matrix:     make([]*BarcodeRow, height),
		currentRow: -1,
		height:     height,
	}
	for i := range m.matrix {
		m.matrix[i] = NewBarcodeRow((width+4)*17 + 1)
	}
	return m
}

// Set sets a specific location in the matrix.
func (bm *BarcodeMatrix) Set(x, y int, value byte) {
	bm.matrix[y].Set(x, value)
}

// StartRow increments the current row counter.
func (bm *BarcodeMatrix) StartRow() {
	bm.currentRow++
}

// CurrentRow lends the mutable selected row. StartRow must have selected an
// in-range row first; subsequent mutations affect the same matrix builder.
func (bm *BarcodeMatrix) CurrentRow() *BarcodeRow {
	return bm.matrix[bm.currentRow]
}

// Matrix returns newly copied unscaled rows in reversed vertical order, matching
// ScaledMatrix(1,1). It is not the builder's internal row slice.
func (bm *BarcodeMatrix) Matrix() [][]byte {
	return bm.ScaledMatrix(1, 1)
}

// ScaledMatrix returns newly allocated rows with integer scaling and vertical
// reversal. Positive scales and products must be bounded by the caller; no
// overflow, allocation budget or concurrent mutation protection is added here.
func (bm *BarcodeMatrix) ScaledMatrix(xScale, yScale int) [][]byte {
	matrixOut := make([][]byte, bm.height*yScale)
	yMax := bm.height * yScale
	for i := 0; i < yMax; i++ {
		matrixOut[yMax-i-1] = bm.matrix[i/yScale].GetScaledRow(xScale)
	}
	return matrixOut
}
