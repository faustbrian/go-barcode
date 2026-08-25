# barcode

[![CI](https://github.com/faustbrian/go-barcode/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-barcode/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-barcode/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-barcode.svg)](https://pkg.go.dev/github.com/faustbrian/go-barcode)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-barcode?sort=semver)](https://github.com/faustbrian/go-barcode/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`barcode` is a standards-driven Go library for validating, encoding,
rendering, and decoding common one-dimensional and two-dimensional barcodes.
Its core values are immutable logical modules; PNG, SVG, and `image.Image`
output are derived views rather than the source of truth.

```go
symbol, err := qr.Encode([]byte("https://example.invalid/parcel/42"), qr.Options{
    ErrorCorrection: qr.Quartile,
})
if err != nil {
    return err
}
return render.PNG(output, symbol.Logical(), render.Options{Scale: 4})
```

Import paths use `github.com/faustbrian/go-barcode`. The image decoder is an
additive package; encoders and logical rendering do not require callers to use
image detection.

## Capability status

No format is marked advertised until its software encoding, decoding,
validation, metadata, and independent interoperability evidence are complete.
Call `barcode.Formats()` and `barcode.CapabilityFor(format)` for the same
machine-readable status. Hardware APIs, device control, and physical device
certification are outside this library's scope.

| Format | Encode | Image decode | GS1 | Advertised | Current software limitation |
|---|---:|---:|---:|---:|---|
| QR Code | yes | yes | yes | yes | none |
| Code 128 / GS1-128 | yes | yes | yes | yes | none |
| Code 39 / Code 93 | yes | yes | no | yes | none |
| EAN-8 / EAN-13 | yes | yes | yes | yes | none |
| UPC-A / UPC-E | yes | yes | yes | yes | none |
| ITF / ITF-14 | yes | yes | ITF-14 | yes | none |
| Codabar | yes | yes | no | yes | optional checksum profiles unavailable |
| Data Matrix ECC 200 | yes | yes | yes | no | structured-append sequence assembly incomplete |
| PDF417 | yes | yes | no | no | macro sequence assembly incomplete |
| Aztec | yes | yes | yes | yes | none |

## Documentation

- [API and format options](docs/api.md)
- [GS1 recipes](docs/gs1.md)
- [Rendering and image decoding](docs/rendering-and-scanning.md)
- [Error-correction tradeoffs](docs/error-correction.md)
- [Software interoperability](docs/interoperability.md)
- [Security and resource limits](docs/security.md)
- [Adoption guide](docs/adoption.md)
- [Comparison with established libraries](docs/comparison.md)
- [FAQ](docs/faq.md)
- [Standards and licensing](specification/README.md)
- [Specification decisions](docs/specification-decisions.md)
- [Conformance scope](docs/conformance.md)

Run `make check` for every blocking local gate. Coverage below meaningful 100%,
missing conformance evidence, or an unsupported advertised control remains a
release blocker.

## Ecosystem

Use the [Golib documentation portal](https://github.com/faustbrian/golib/blob/main/docs/index.md)
to choose companion packages, supported stacks, recipes, and operations guidance.
