# Performance and benchmarks

Barcode performance depends on payload length, symbology, correction level,
logical dimensions, render scale, source-image dimensions, requested formats,
and enabled decode transformations. Compare results only when those inputs and
the Go version, platform, and commit are the same.

The executable benchmarks in [`benchmark_test.go`](../benchmark_test.go) cover
QR and Code 128 encoding, QR image and encoded-stream decoding, raster
rendering, no-symbol detection, and malformed encoded-image rejection. They
report allocations, and selected byte-oriented cases report processed bytes.

Run the repository-owned benchmark gate:

```sh
make -f verification/package.mk benchmark
```

The gate also records peak resident memory for a representative bounded QR
decode. Treat the output as comparable evidence for its exact environment, not
as a universal throughput guarantee. [Software interoperability](interoperability.md)
is a separate correctness boundary; a faster result does not establish that a
format is advertised or independently interoperable.

Applications should benchmark their own payload distribution, rendering
dimensions, image transformations, and `imagedecode.Limits`. Do not disable
resource limits to improve a benchmark that is intended to model untrusted
input.
