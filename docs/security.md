# Security and resource limits

## Threat model (v1, 2026-09-28)

The assets are service availability, barcode payload confidentiality, and the
integrity of encoded and decoded symbol data. Callers may pass attacker-chosen
payloads, options, images, encoded streams, and output writers. The library
does not authenticate users or interpret decoded URLs; applications own those
authorization and content-use decisions.

| Boundary | Main threat | Owned control |
| --- | --- | --- |
| Encoders and GS1 parser | Oversized or malformed payloads consume resources or escape through errors | Input limits, validation, and classified errors |
| `DecodeEncoded` stream and image decoders | Compressed-image amplification, malformed data, and error-text disclosure | Encoded-byte, dimension, pixel, memory, and attempt limits; redacted invalid-image errors |
| `Decode` image and third-party symbol readers | Expensive candidates, dependency panics, and overlong decoded data | Pre-conversion geometry limits, bounded attempts, payload limits, and documented two-dimensional panic containment |
| Logical symbols and rendering | Dimension overflow or excessive raster allocation | Checked products and configurable pixel limits |

The decoder and encoder dependencies and the Go image codecs are trusted code
processing untrusted data. Dependency review and vulnerability scans remain
release evidence, not a substitute for hostile-input tests. No API here opens
network connections, files, processes, or background workers. Caller-provided
readers and writers may do so; their permissions and lifetimes remain with the
caller.

**Open risk — blocking readers:** `DecodeEncoded` cannot interrupt an arbitrary
`io.Reader` while its `Read` call is blocked, so `MaxDuration` is checked after
that call returns, not enforced within it. The go-barcode maintainers own a
future reader-ownership design review before claiming end-to-end cancellation.
Until then, callers handling hostile streams must use a reader with its own
deadline or cancellation behavior. This risk is not accepted as resolved.

Barcode payloads and decoded URLs are untrusted data. This library never
executes, fetches, redirects to, or automatically follows decoded content.

`imagedecode.Limits` bounds encoded bytes, width, height, total pixels,
candidate attempts, payload bytes, rotations, memory, corrections, and time
before or between expensive work. Cancellation and the derived time deadline
are checked before conversion and between decoder attempts. Set stricter
application-specific limits for public uploads.

`MaxCorrections` limits the exact corrected-error count reported by QR, Data
Matrix, Aztec, and PDF417 readers. Linear formats report zero because their
readers validate checksums rather than applying error correction.

`DecodeEncoded` reports invalid reader and image data with `ErrInvalidImage`
without returning lower-level error text or wrapping lower-level errors, which
may contain sensitive input. Context cancellation remains distinguishable.
PDF417 text-compaction failures report the invalid character position without
echoing the payload character.

Use `DecodeEncoded` for untrusted PNG, JPEG, or GIF streams. It limits input
bytes, inspects encoded dimensions, applies pixel and memory budgets before
full decompression, rejects truncation, and checks cancellation between stages.
Callers using `Decode` directly remain responsible for how their `image.Image`
was allocated.

The additive QR, Data Matrix, and Aztec readers isolate their dependencies at
the adapter boundary. A panic from malformed symbol data in those decoder
dependencies is converted to a failed candidate and cannot escape `Decode`;
another bounded candidate may still be attempted. PDF417 does not make this
panic-containment guarantee.

Matrix and renderer constructors reject invalid dimensions, multiplication
overflow, non-positive scales, and allocation budgets. Errors avoid including
full payloads by default. Report suspected vulnerabilities through the private
process in [`SECURITY.md`](../SECURITY.md).
