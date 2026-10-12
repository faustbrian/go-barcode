# API and format options

The `barcode` package owns typed format identifiers, immutable `Matrix` and
`Bars` values, `Symbol`, `DecodeResult`, checksum state, diagnostics, and the
capability registry. Constructors copy caller slices; getters return defensive
copies where a slice could otherwise alias internal state.

Each encoder has a focused package and typed options:

- `qr`: mode, version, mask, ECI, FNC1, structured append, correction, and
  quiet zone;
- `code128`: automatic or forced code sets, validated GS1-128, and an explicit
  raw-FNC1 compatibility mode for legacy symbols;
- `code39` and `code93`: full-ASCII/checksum behavior supported by the format;
- `ean` and `upc`: check digits and two- or five-digit supplements;
- `itf`: ITF and ITF-14 dimension/bearer controls;
- `codabar`: explicit validated guards;
- `datamatrix`: ECC 200 shape/dimensions, GS1, ECI, Macro 05/06, structured
  append, and controlled Base 256;
- `pdf417`: compaction, correction, compact layout, row/column constraints,
  ECI, and Macro PDF417 control-block metadata (not sequence assembly);
- `aztec`: compact/full layers, correction percentage, GS1, and ECI.

Encoders include profile-specific quiet zones in logical symbols; these
defaults do not qualify physical printing, scanning, or every licensed-standard
requirement. `render` converts them to raster or SVG at exact integer scale.
`imagedecode.Decode` accepts an `image.Image`, requested formats, behavior
flags, and explicit resource limits, including an optional elapsed-time budget.
`DecodeEncoded` accepts a caller-owned byte slice and bounds compressed PNG,
JPEG, or GIF bytes before full decoding. It does not read from an arbitrary
`io.Reader`. Content is returned as untrusted bytes and text; the library never
executes or follows it.

Errors expose package-level sentinels; use `errors.Is`, not string matching.
Some errors wrap diagnostic causes, but `DecodeEncoded` returns bare
`imagedecode.ErrInvalidImage` for malformed bytes and image-decoder failures so
untrusted error text cannot disclose payloads. Cancellation remains
distinguishable as `context.Canceled` or `context.DeadlineExceeded`.

ECI support depends on the encoder path and assignment. Declaring an ECI is not
general-purpose transcoding. Controlled Data Matrix decoding preserves raw
bytes for a syntactically valid unknown ECI and reports both ECI_ASSIGNMENT and
ECI_UNSUPPORTED diagnostics. GS1-aware encoders validate element syntax and
wrap the parser's classified cause; this is not allocation or authenticity
proof. Reader checksum failures discard candidates, but returned status is a
format classifier, not independent GS1 or ITF-14 application validation.

Core constructors copy slices and validate representation/metadata shape. They
impose no general payload or diagnostic length cap, do not redact content, and
do not prove that logical modules agree with the declared format. Confidence
uses range comparisons rather than explicit finite-number validation. Treat
decoded payloads, raw bytes, diagnostics, and reader causes as untrusted
content.

Image memory limits use a four-bytes-per-pixel admission estimate, not peak
heap measurement. Payload and correction caps are checked after reader work.
Cancellation is observed between phases/attempts, not inside borrowed image
methods, codecs, or synchronous readers. An application requiring a hard
deadline or heap bound needs its own isolation boundary.
