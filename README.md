# barcode

[![CI](https://github.com/faustbrian/go-barcode/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-barcode/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-barcode/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-barcode.svg)](https://pkg.go.dev/github.com/faustbrian/go-barcode)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-barcode?sort=semver)](https://github.com/faustbrian/go-barcode/releases)
[![Go](https://img.shields.io/badge/go-1.27.0-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`barcode` is a standards-driven Go library for validating, encoding,
rendering, and decoding common one-dimensional and two-dimensional barcodes.
Its core values are immutable logical modules; PNG, SVG, and `image.Image`
output are derived views rather than the source of truth.

Import paths use `github.com/faustbrian/go-barcode`. The image decoder is an
additive package; encoders and logical rendering do not require callers to use
image detection.

## Status and platform

The module is stable at v1 and requires Go 1.26.6. Its public packages are
portable Go and require no operating-system service, device driver, camera,
printer, scanner, or external runtime backend. Physical-device control,
print-quality certification, and multi-symbol sequence assembly remain outside
the module's scope.

## Install

```sh
go get github.com/faustbrian/go-barcode
```

## Quick start

```go
package main

import (
    "fmt"

    "github.com/faustbrian/go-barcode/qr"
)

func main() {
    symbol, err := qr.Encode([]byte("https://example.invalid/parcel/42"),
        qr.Options{ErrorCorrection: qr.Quartile})
    if err != nil {
        panic(err)
    }
    matrix := symbol.Logical().Matrix()
    fmt.Printf("%s %dx%d\n", symbol.Logical().Format(), matrix.Width(), matrix.Height())
}
```

The checked-in [`Example_qrCode`](example_test.go) is the executable version of
this flow. It is compiled and run by the Go example test gate.

## When to use it

Use this module when an application needs immutable logical barcode symbols,
strict format controls, bounded software image decoding, deterministic
rendering, or auditable standards and interoperability evidence.

Do not use it to control barcode hardware, certify physical print or scan
quality, discover every symbol in one image, or assemble Data Matrix structured
append and Macro PDF417 sequences. Applications own those workflows and any
side effects triggered by decoded content.

## Packages

- `barcode` owns shared immutable symbols, format identifiers, capabilities,
  decode results, and classified core errors.
- `qr`, `code128`, `code39`, `code93`, `ean`, `upc`, `itf`, and `codabar`
  provide focused one- and two-dimensional encoders.
- `datamatrix`, `pdf417`, and `aztec` provide bounded two-dimensional format
  controls and encoding.
- `gs1` parses and validates bounded GS1 element strings and check digits.
- `render` derives `image.Image`, PNG, or SVG output from logical symbols.
- `imagedecode` performs caller-bounded software image detection and decoding.
- `specification` exposes the embedded standards catalogue and conformance
  metadata.

See the [API map](docs/api.md) for format-specific options and limitations.

## Lifecycle, ownership, and concurrency

Encoding, parsing, rendering, and decoding are synchronous operations. They
start no background work and acquire no resources requiring `Close` or
`Shutdown`. Logical symbols and results are immutable; constructors copy
caller byte and slice inputs, and accessors return defensive copies.

Public functions are safe to call concurrently when callers do not mutate an
input object or use the same `io.Writer` concurrently. `render.PNG` and
`render.SVG` borrow the supplied writer for the duration of the call and never
close it. `imagedecode.DecodeEncoded` reads but never closes its caller-owned
reader. `Decode` borrows its caller-owned `image.Image` for the call.

Only image-decoding operations accept `context.Context`. They do not retain the
context and check cancellation between bounded stages and candidate attempts;
timeouts cannot interrupt a blocking caller-owned reader or image method, or
third-party decoder code while that code is actively running. Configure
`imagedecode.Limits` for public or otherwise hostile input. Use `errors.Is`
with package sentinels rather than matching error text.

## Capability status

No format is marked advertised until its software encoding, decoding,
validation, metadata, and independent interoperability evidence are complete.
Call `barcode.Formats()` and `barcode.CapabilityFor(format)` for the
machine-readable per-format registry. The overview below groups related
formats, and its GS1 column reports an explicit GS1 encoding entry point rather
than serializing the registry's format identity. Hardware APIs, device control,
and physical device certification are outside this library's scope.

| Format | Encode | Image decode | GS1 encode | Advertised | Current software limitation |
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

- [Documentation index](docs/README.md)
- [API and format options](docs/api.md)
- [Executable examples](example_test.go)
- [GS1 recipes](docs/gs1.md)
- [Rendering and image decoding](docs/rendering-and-scanning.md)
- [Error-correction tradeoffs](docs/error-correction.md)
- [Performance and benchmarks](docs/performance.md)
- [Software interoperability](docs/interoperability.md)
- [Security and resource limits](docs/security.md)
- [Adoption guide](docs/adoption.md)
- [Comparison with established libraries](docs/comparison.md)
- [FAQ](docs/faq.md)
- [Troubleshooting](docs/troubleshooting.md)
- [Standards and licensing](specification/README.md)
- [Specification decisions](docs/specification-decisions.md)
- [Conformance scope](docs/conformance.md)
- [Compatibility](COMPATIBILITY.md)
- [Changelog](CHANGELOG.md)
- [Support](SUPPORT.md)
- [Security reporting](SECURITY.md)
- [Contribution guide](CONTRIBUTING.md)
- [License](LICENSE)

See the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.5.0/docs/ecosystem/README.md)
and [Domain utilities family guidance](https://github.com/faustbrian/go-library-tools/blob/v1.5.0/docs/ecosystem/design-language.md#package-families-and-selection)
for the shared design language this module follows.

Run `make ci` for the complete local repository contract, including
configuration, inventory, cohesion, repository, workflow, and online
specification validation. Use `make check` for implementation gates alone.
Coverage below meaningful 100%, missing conformance evidence, or an unsupported
advertised control remains a release blocker.

Use `make cohesion` to validate the module's ecosystem classification,
documentation links, and package-selection metadata locally.
