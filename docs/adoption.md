# Adoption guide

1. Inventory formats, payload lengths, GS1 use, render dimensions, and decode
   inputs in the existing system.
2. Query `barcode.CapabilityFor` and accept the documented limitations for each
   selected format. Do not convert an unadvertised format into a product claim.
3. Compare logical modules and decoded metadata against existing production
   fixtures. Pixel snapshots alone are insufficient.
4. Render at integer scale with the old and new pipelines, then decode the same
   image corpus with identical transformations and resource limits.
5. Set explicit decode limits from the application threat model.
6. Roll out per format and retain the previous encoder/decoder until observed
   decode and error rates meet the acceptance criteria.

Migration should preserve canonical payload bytes, check-digit ownership, GS1
separator behavior, quiet zones, and requested correction settings. Treat any
silent normalization difference as a compatibility defect.

For v2, change imports to `github.com/faustbrian/go-barcode/v2`. Pass an
already acquired `[]byte` to `imagedecode.DecodeEncoded` instead of an
`io.Reader`. Applications own stream acquisition, its deadline and cancellation,
and its byte cap; the decoder enforces `MaxEncodedBytes` before image parsing.
Do not start a goroutine solely to wrap a blocking reader, because cancellation
cannot guarantee that goroutine exits. Version 1 remains available at its
unchanged import path for callers that cannot yet migrate.
