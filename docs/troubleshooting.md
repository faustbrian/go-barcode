# Troubleshooting

## Encoding returns a validation error

Check the format-specific payload alphabet, length, checksum policy, quiet-zone
minimum, correction level, and structured-control fields in the
[API map](api.md). Use `errors.Is` with the package sentinel; error strings are
diagnostic text, not a classification contract.

## Decoding returns `imagedecode.ErrNotFound`

Confirm the requested format, preserve the complete quiet zone, and avoid
non-integer image resampling. The zero value tries up to four rotations; lower
`Limits.MaxRotations` when fewer attempts are needed, and enable inversion only
when the input requires it. A checksum-failing candidate is deliberately not
returned as a successful decode.

## Decoding returns `imagedecode.ErrLimitExceeded`

Inspect the encoded-byte, image-dimension, pixel, memory, candidate, rotation,
payload, correction, and duration limits. Raise only the bound supported by the
application's threat model; do not remove all bounds for public uploads.

## A generated symbol does not scan on a physical device

First verify the logical symbol through the pinned software interoperability
tests. Then inspect integer module scale, contrast, quiet zone, printer
resolution, substrate, and camera conditions. This repository does not certify
physical devices or print quality.

## A format is implemented but not advertised

Use `barcode.CapabilityFor` to inspect its explicit limitations. Encoding and
image decoding alone do not establish complete validation, metadata, sequence
assembly, or independent interoperability support.

For unresolved reproducible problems, follow the [support policy](../SUPPORT.md).
Report suspected vulnerabilities through the private process in
[`SECURITY.md`](../SECURITY.md).
