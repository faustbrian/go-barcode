package code128

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
)

func TestMatrixRunsAllocationBudget(t *testing.T) {
	matrix, err := gozxing.NewBitMatrix(128, 1)
	if err != nil {
		t.Fatal(err)
	}
	for x := 0; x < matrix.GetWidth(); x += 2 {
		matrix.Set(x, 0)
	}
	allocations := testing.AllocsPerRun(100, func() {
		runs := matrixRuns(matrix)
		if len(runs) != matrix.GetWidth() {
			t.Fatalf("run count = %d, want %d", len(runs), matrix.GetWidth())
		}
		for x, run := range runs {
			if run.Width != 1 || run.Dark != (x%2 == 0) {
				t.Fatalf("run %d = %+v", x, run)
			}
		}
	})
	if allocations > 2 {
		t.Fatalf("matrix conversion allocated %.0f times, want at most two", allocations)
	}
	benchmark := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if runs := matrixRuns(matrix); len(runs) != matrix.GetWidth() {
				b.Fatal("matrix conversion lost runs")
			}
		}
	})
	if bytes := benchmark.AllocedBytesPerOp(); bytes > int64(matrix.GetWidth()*24) {
		t.Fatalf("matrix conversion allocated %d bytes, exceeding its pixel budget", bytes)
	}
}

func TestMatrixRunsQuietZoneAllocationBudget(t *testing.T) {
	matrix, err := gozxing.NewBitMatrix(8192, 1)
	if err != nil {
		t.Fatal(err)
	}
	matrix.Set(4096, 0)
	benchmark := testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			runs := matrixRuns(matrix)
			if len(runs) != 3 || runs[0].Dark || !runs[1].Dark || runs[2].Dark ||
				runs[0].Width != 4096 || runs[1].Width != 1 || runs[2].Width != 4095 {
				b.Fatalf("quiet-zone runs = %+v", runs)
			}
		}
	})
	if bytes := benchmark.AllocedBytesPerOp(); bytes > 128 {
		t.Fatalf("quiet zones inflated run storage to %d bytes", bytes)
	}
}

func TestCodeSetNames(t *testing.T) {
	tests := map[CodeSet]string{
		CodeSetAuto: "", CodeSetA: "A", CodeSetB: "B", CodeSetC: "C", CodeSet(99): "",
	}
	for codeSet, want := range tests {
		if got := codeSetName(codeSet); got != want {
			t.Fatalf("codeSetName(%v) = %q, want %q", codeSet, got, want)
		}
	}
}
