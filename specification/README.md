# Standards and fixture inventory

This repository identifies standards but does not redistribute restricted ISO
or AIM publications. Implementers need licensed access to the listed editions
for normative review. Public metadata links are recorded in the `barcode`
capability registry.

The GS1 Barcode Syntax Dictionary is redistributed under Apache-2.0 from the
GS1 repository. It is pinned to release `2026-01-27` and SHA-256
`6461ece7c03420fb27a3f18163ef9a08694aa4bb2491dea25f63bcdf22451a99`.
Run `scripts/sync-gs1-dictionary.sh` to reproduce it.

`normative.tsv` maps each required format behavior to its pinned standard.
`evidence.tsv` maps the same stable IDs to executable or explicitly partial
evidence. They contain descriptive metadata only and do not reproduce
restricted standards text. Every normative ID must have exactly one evidence
row; the documentation gate verifies that relationship.

`manifest.json` also records authoritative catalogue links for every governing
edition. ISO and AIM publications are licensed and are not redistributed, so
their exact document identities and catalogue records are the provenance
boundary rather than a fabricated checksum of a mutable product page. Hashes
are recorded for every redistributed fixture and independent implementation
archive. The package remains pinned to ISO/IEC 15420:2009 and ISO/IEC
24778:2008; their 2025 and 2024 replacements are tracked but are not claimed
until a licensed edition-difference review and matching conformance evidence
are complete.

Material interpretations and package-owned policy are recorded in
[`docs/specification-decisions.md`](../docs/specification-decisions.md).

`render-fixtures.tsv` inventories deterministic logical, PNG, and SVG goldens
for every implemented format. The checked-in files in `render-fixtures/` are
canonical software rendering fixtures. `make docs` verifies their hashes.

For Code 93, ANSI's 2006 catalogue and withdrawal notice identify
`ANSI/AIM BC5-1995` as a withdrawn historical publication. AIM's current
publisher catalogue separately identifies `USS - Code 93` with publication
date 2000 and no BC5 designation. Public metadata does not establish that the
2000 product replaces, revises, or is normatively equivalent to BC5-1995, so
the capability remains explicitly pinned to the withdrawn 1995 compatibility
target. The 2000 identity is recorded for review, not claimed as implemented.

The independent writer tests use `github.com/speedata/barcode` v1.1.1 under
the MIT license. Its module archive hash is recorded in `manifest.json` and its
module checksum is enforced by `go.sum`. The tests generate QR, Code 128,
Code 39, Code 93, EAN-8, EAN-13, ITF, Codabar, Data Matrix, and Aztec symbols
through that separate implementation before decoding them through this
library. They also generate GS1-128 and ITF-14 payloads through the independent
Code 128 and ITF writers. UPC-A and UPC-E use the Apache-2.0 ZXing writer.
PDF417 uses `github.com/ruudk/golang-pdf417` at commit `a7e3863a1245` under the
MIT license because speedata/barcode v1.1.1 emits an invalid PDF417
error-correction checksum for this fixture.

The reciprocal reader tests use `github.com/ericlevine/zxinggo` v0.1.0 under
Apache-2.0. They render every format from this library and decode it through
that separate reader implementation. Its module archive hash is also recorded
in `manifest.json` and its module checksum is enforced by `go.sum`.

No ISO or AIM example is represented as an official conformance fixture unless
its redistribution terms independently permit that use. Independently derived
logical vectors identify their derivation in the corresponding test.

`monitoring.json` records ISO and ANSI publication identities as restricted sources and
checks the catalogue's exact `403` denial without hashing its body. The public
AIM source pin identifies AIM's standards program, not the licensed Codabar
publication bytes. AIM public-review notices and GS1 change notices are
separate content-pinned change surfaces. None of those public hashes is a hash
of restricted normative text, and no exact restricted clause is claimed where
the publication was unavailable for review.

## Decision bindings

[Specification decisions](../docs/specification-decisions.md) are bound to
executable evidence below.

| Decision | Executable evidence |
| --- | --- |
| BARCODE-DEC-001 | `TestCapabilityMetadataIdentifiesExactGoverningEditions`, `TestGS1SyntaxDictionaryIsEmbedded`, `TestRenderFixtureGoldensCoverEveryFormat` |
| BARCODE-DEC-002 | `TestCapabilitiesAreExplicitForEveryKnownFormat`, `TestCapabilityMetadataIdentifiesExactGoverningEditions`, `TestSymbolsDecodeWithIndependentReaders` |
| BARCODE-DEC-003 | `TestCapabilitiesAreExplicitForEveryKnownFormat`, `TestCapabilitiesReflectSoftwareScope`, `TestDecodeEverySupportedFormat` |
| BARCODE-DEC-004 | `TestMatrixAndSymbolDoNotAliasCallerData`, `TestBarsAndDecodeResultDoNotAliasCallerData`, `TestSymbolAcceptsExactlyOneLogicalRepresentation` |
| BARCODE-DEC-005 | `TestImageUsesIntegerModuleScalingAndExactColors`, `TestRenderedOutputsMatchGoldenChecksums`, `TestRenderRejectsOverflowAndExplicitLimitViolations`, `FuzzRenderLogicalMatrices` |
| BARCODE-DEC-006 | `TestEncodeHonorsMaskECIAndGS1Controls`, `TestEncodeStructuredValidatesSequenceOptions`, `TestEncodeRejectsModeCharsetAndCapacityMismatches`, `FuzzEncodeOptions` |
| BARCODE-DEC-007 | `TestEncodeACalculatesAndValidatesCheckDigit`, `TestEncode14CalculatesCheckDigitAndAddsBearerBars`, `TestChecksumMatchesModulo43Vector`, `TestEncodeAddsMandatoryChecksumsAndSupportsFullASCII` |
| BARCODE-DEC-008 | `TestParseRawElementStringUsesPredefinedLengthsAndFNC1`, `TestParserEnforcesRequiredAndExcludedAssociations`, `TestEncodeGS1AcceptsValidatedStructuredElements`, `TestEncodeRawFNC1SupportsLegacyPayloadsWithoutGS1Validation`, `FuzzParseElementStrings` |
| BARCODE-DEC-009 | `TestEncodeSupportsStructuredAppendBoundaries`, `TestEncodeSupportsECIAssignmentWidths`, `TestEncodeSupportsMacro05And06`, `TestDecodeDataMatrixControls` |
| BARCODE-DEC-010 | `TestEncodeSupportsECIAndMacroControlBlocks`, `TestEncodeSupportsAllCorrectionLevelsAndQuietZones`, `TestEncodeDistinguishesPayloadSafetyFromSymbolCapacity`, `TestPDF417ReaderImplementsReaderLifecycle` |
| BARCODE-DEC-011 | `TestEncodeSupportsAutomaticAndForcedLayers`, `TestEncodeSelectsSmallestAutomaticCompactLayer`, `TestEncodeSupportsGS1FNC1`, `TestEncodeSupportsECI` |
| BARCODE-DEC-012 | `TestDecodeEnforcesBoundsBeforeImageAllocation`, `TestDecodeEnforcesCallerTimeBudget`, `TestTwoDReaderContainsDependencyPanics`, `FuzzDecodeBoundedImages` |
| BARCODE-DEC-013 | `TestDecodeQRDocumentedImageDegradationThresholds`, `TestDecodeSupportsInvertedImages`, `TestDecodeMultipleSymbolsReturnsOneDecodableCandidate`, `TestDecodeEnforcesCorrectionBudget` |
| BARCODE-DEC-014 | `TestDecodeReportsRotationAndPayloadLimits`, `TestOrientationChecksumAndFormatMappings`, `TestInvalidInputErrorsAreClassifiedAndPayloadRedacted`, `TestDecodeRejectsUnsupportedFormatsAndCandidateLimits` |
| BARCODE-DEC-015 | `TestDecodeIndependentWriters`, `TestSymbolsDecodeWithIndependentReaders`, `TestEncodeMatchesPinnedZXingImplementation`, `TestEncoderMatchesPinnedZXingImplementation` |
| BARCODE-DEC-016 | `TestCapabilitiesReflectSoftwareScope`, `TestDecodeResultValidatesMetadataAndReturnsDefensiveValues`, `TestInvalidInputErrorsAreClassifiedAndPayloadRedacted` |
