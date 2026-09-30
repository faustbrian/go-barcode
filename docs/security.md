# Security and resource limits

## Threat model (v2 candidate, 2026-09-30)

The assets are service availability, barcode payload confidentiality, and the
integrity of encoded and decoded symbol data. Callers may pass attacker-chosen
payloads, options, images, encoded bytes, and output writers. The library
does not authenticate users or interpret decoded URLs; applications own those
authorization and content-use decisions.

| Boundary | Main threat | Owned control |
| --- | --- | --- |
| Encoders and GS1 parser | Oversized or malformed payloads consume resources or escape through errors | Input limits, validation, and classified errors |
| `DecodeEncoded` bytes and image decoders | Compressed-image amplification, malformed data, and error-text disclosure | Encoded-byte, dimension, pixel, memory, and attempt limits; redacted invalid-image errors |
| `Decode` image and third-party symbol readers | Expensive candidates, dependency panics, and overlong decoded data | Pre-conversion geometry limits, bounded attempts, payload limits, and documented two-dimensional panic containment |
| Logical symbols and rendering | Dimension overflow or excessive raster allocation | Checked products and configurable pixel limits |

The decoder and encoder dependencies and the Go image codecs are trusted code
processing untrusted data. Dependency review and vulnerability scans remain
release evidence, not a substitute for hostile-input tests. No API here opens
network connections, files, processes, or background workers. Caller-provided
image methods and output writers may do so; their permissions and lifetimes
remain with the caller.

The v1 `DecodeEncoded(io.Reader)` boundary could remain blocked inside an
arbitrary caller-owned `Read` after cancellation. V2 removes that reader
operation: `DecodeEncoded` accepts already acquired `[]byte`, checks
`MaxEncodedBytes` before parsing, and never starts a read goroutine. Callers
own stream acquisition, including its deadline, cancellation, and byte cap.
The v1 import path and its documented risk remain unchanged.

**Residual risk — in-process decoding:** `MaxDuration` and context cancellation
are checked between decoding stages and candidates, not inside a trusted image
codec, third-party symbol decoder, or caller-owned `image.Image` method. The
go-barcode maintainers own this bounded-process availability risk. The rationale
is that Go cannot safely preempt those synchronous calls; the mitigation is
encoded-byte, dimension, pixel, memory, correction, and candidate limits, plus
caller isolation when a hard wall-clock deadline is required. Review this
disposition when a codec changes, a bounded hostile image exceeds the service's
time budget, or an uncaught decoder panic is reproduced. No end-to-end
wall-clock cancellation guarantee is claimed.

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

`DecodeEncoded` reports invalid encoded and image data with `ErrInvalidImage`
without returning lower-level error text or wrapping lower-level errors, which
may contain sensitive input. Context cancellation remains distinguishable.
PDF417 text-compaction failures report the invalid character position without
echoing the payload character.

Use `DecodeEncoded` for bounded PNG, JPEG, or GIF bytes after acquiring them
under an application-owned stream deadline and size cap. It rejects oversized
input before parsing, inspects encoded dimensions, applies pixel and memory
budgets before full decompression, rejects truncation, and checks cancellation
between stages.
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
