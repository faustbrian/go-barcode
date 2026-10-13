# Software interoperability

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API and format options](api.md)

Automated writer interoperability uses the independently maintained
`github.com/speedata/barcode` v1.1.1 implementation for all formats except
PDF417, which uses `github.com/ruudk/golang-pdf417`. EAN/UPC logical encoders
are also decoded through ZXing, and logical-module golden vectors cover UPC-E
and ITF-14. Self-round-trip tests remain insufficient because matching defects
can agree. Reciprocal reader evidence renders every format from this library
and decodes it with `github.com/ericlevine/zxinggo` v0.1.0. Provenance and
archive checksums are recorded in
`specification/manifest.json`.

Hardware APIs, camera capture, printer control, scanner control, and physical
device certification are outside this library's scope. Interoperability claims
refer only to repeatable software-to-software tests with pinned independent
implementations.

Canonical PNG and SVG inputs for every implemented format live under
`specification/render-fixtures/`. Their logical-module, PNG, and SVG hashes,
payload classes, scales, and correction settings are pinned in
`specification/render-fixtures.tsv`. Regenerate them deliberately with:

```console
go test . -run TestRenderFixtureGoldensCoverEveryFormat -count=1 \
  -args -update-render-fixtures
```

Normal tests reconstruct every logical symbol and render it, then compare the
bytes with checked-in artifacts. These fixtures prove deterministic software
output and provide inputs for independent software decoders.

`make -f verification/package.mk benchmark` records isolated logical encode,
render, no-symbol detection, decoded-image scanning, encoded-image decoding,
malformed-input rejection,
bytes per operation, allocations per operation, and peak resident memory for
a representative encoded QR decode. This selected performance-risk route and
`make check-full` retain the configured benchmark campaign; ordinary local and
PR checks do not run it. Compare results with an exact source, Go version,
workload and duration.
