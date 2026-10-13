package qr_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/faustbrian/go-barcode/v2/qr"
)

func TestStructuredAppendRejectsOriginalIntegerAliases(t *testing.T) {
	for total := 2; total <= 16; total++ {
		for _, index := range []int{-1, total, 16, 17, 256, 257, int(^uint(0) >> 1)} {
			symbol, err := qr.Encode([]byte("A"), qr.Options{StructuredAppend: &qr.StructuredAppend{Index: index, Total: total, Parity: 0x5a}})
			if !errors.Is(err, qr.ErrInvalidInput) || !reflect.DeepEqual(symbol, qr.Symbol{}) {
				t.Fatalf("original Index=%d Total=%d accepted or published a symbol: %v", index, total, err)
			}
		}
		for index := range total {
			header := qr.StructuredAppend{Index: index, Total: total, Parity: 0x5a}
			symbol, err := qr.Encode([]byte("A"), qr.Options{StructuredAppend: &header})
			if err != nil {
				t.Fatalf("valid Index=%d Total=%d: %v", index, total, err)
			}
			retained, ok := symbol.StructuredAppend()
			if !ok || retained != header {
				t.Fatal("valid explicit sequence identity changed")
			}
		}
	}
}

// FuzzStructuredAppendAdmission checks original integer admission, not packed aliases.
func FuzzStructuredAppendAdmission(f *testing.F) {
	for _, pair := range [][2]int{{0, 2}, {15, 16}, {16, 16}, {256, 2}, {-1, 2}, {0, 1}, {0, 17}, {int(^uint(0) >> 1), 16}} {
		f.Add(pair[0], pair[1])
	}
	f.Fuzz(func(t *testing.T, index, total int) {
		header := qr.StructuredAppend{Index: index, Total: total, Parity: 0x5a}
		symbol, err := qr.Encode([]byte("A"), qr.Options{StructuredAppend: &header})
		admitted := total >= 2 && total <= 16 && index >= 0 && index < total
		if !admitted {
			if !errors.Is(err, qr.ErrInvalidInput) || !reflect.DeepEqual(symbol, qr.Symbol{}) {
				t.Fatalf("unadmitted header published: Index=%d Total=%d error=%v", index, total, err)
			}
			return
		}
		retained, ok := symbol.StructuredAppend()
		if err != nil || !ok || retained != header {
			t.Fatalf("admitted header lost: Index=%d Total=%d error=%v", index, total, err)
		}
	})
}
