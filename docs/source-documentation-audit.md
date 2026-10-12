# Source documentation audit

## Scope and acceptance

Baseline `a5798b5606e278b36d18e185c67a3cb94f7f18fb`; one module,
`github.com/faustbrian/go-barcode/v2`, 76 tracked Go files and 19 catalog source
packages, plus the test-only root package. Tier A source contracts only:
preserve executable AST/tokens/literals/directives, dependency/API baselines,
standards identities, embedded dictionary, independent peers and frozen fixtures.
Owner: barcode maintainers. Consumers: symbology encoders, format-neutral
symbols, GS1 parsing, raster/vector writers, bounded image readers and internal
encoder adapters. No runtime defect repair or physical/standards qualification.

Smallest proof: changed Go formatting, all-file structural identity and public
comment/header checks, parsed comment rendering, native package Go-doc output,
local docs/links, offline specification/history validation, unchanged fixture and
module inventories, strong changed/full/staged secrets and one complete final
diff inspection. Executable examples are unchanged, so runtime/examples,
conformance/readers/writers, race, active fuzz, benchmarks, online authorities,
fixture regeneration, dependency matrices and releases are not selected.

## Missing documentation versus investigated constraints

Missing: public configuration fields and internal mutable encoder obligations;
complete package ownership/non-goals; QR StructuredAppend's obsolete rejection
comment; borrowed callbacks and writer partial-output behavior; actual parser,
image-time/memory and logical-copy limits. Aztec had duplicate canonical package
comments; the primary header now owns intent and derivative attribution.
No generated Go header, build tag,
platform variant or TODO/FIXME marker was found. Existing conversion/security and
embed directives are preserved exactly. Long lookup tables and local ports
are hand-maintained derivatives, not generated Go code.

Source-backed limitations, not repair scope: core constructors validate shape
and copy slices but do not impose payload/diagnostic caps or prove a symbol's
format/representation agreement; confidence uses range comparisons, not an
explicit finite-number check. Raster limits precede new pixel allocation, not
all logical defensive copies, arbitrary Color callbacks or writer allocations.
Image memory is a four-bytes-per-pixel admission estimate, not peak heap usage;
context is checked between phases/attempts, not inside image codecs/readers or
caller image callbacks. Success does not necessarily recheck a deadline.
Reader-reported checksum metadata is format classification, not an independent
GS1 or ITF-14 application validator. GS1 string slices can retain input backing
storage; field syntax is not authentication or registry-allocation proof.
Some linear quiet zones have only lower bounds, and full-ASCII expansion occurs
before final encoded-length admission. Those caller-sized/overflow/allocation
boundaries must not be described as globally bounded. Internal encoders trust
validated public wrappers and remain mutable, not independently safe parsers.

QR explicit structured-append Index has only a nonnegative local check before
four-bit byte packing. Callers must supply Index below Total and within 0..15;
original-integer upper validation is a runtime follow-up, not repaired here.
Forced-control claims are limited to local admission and the selected peer path,
not inferred universal validation. Existing conversion/security markers remain
unchanged; their historical assumptions do not prove the original Index bound.
Restricted
ISO/AIM exact clauses are not available here: cite the pinned edition and local
decision/evidence inventory, not invented normative clause claims. No proof of
new-edition or hardware/optical conformance will be inferred from comments.

## Native policy, provenance and generation

`.golib.yaml` selects apidiff `api/baseline.txt`, not a prose Go-doc snapshot.
Inspected `Makefile`, `verification/package.mk` and all six scripts; native docs
also run unchanged examples/render goldens, which are outside this prose pass.
`scripts/sync-gs1-dictionary.sh` authenticates release 2026-01-27 and SHA256
6461ece7c03420fb27a3f18163ef9a08694aa4bb2491dea25f63bcdf22451a99 before
replacing the embedded dictionary; it was not executed here. The embed
producer is explicit and no frozen source is regenerated for comment checks.

`specification/manifest.json`, normative.tsv/evidence.tsv, decision records and
NOTICE.md own exact editions and derivative licenses: QR ISO/IEC 18004:2024,
Code128 15417:2007 + Amendment1:2026, Code39 16388:2023, EAN/UPC15420:2009,
ITF16390:2007, DataMatrix16022:2024, PDF41715438:2015, Aztec24778:2008,
historical Code93BC5-1995/CodabarBC3-1995 and GS1 release26.0/January2026.
Replacement editions/current AIM publications are review inputs, not implemented
claims. Local Aztec/PDF417 ports derive from zxinggo v0.1.0; GS1 linter tables
from reference commit 2043848ad254657237f40aabef4ee66461fabc4b. Preserve pins,
copyrights and history prefixes; append native typed digests only for necessary
current-record corrections. DataMatrix/PDF417 sequence assembly stays unadvertised.

## High-risk evidence map

Core copy/shape/metadata tests protect defensive symbol values; render overflow,
color/golden and fuzz tests protect integer-module geometry. Each symbology's
boundary/charset/checksum tests document its actual subset. QR tests distinguish
forced/automatic controls, ECI width, split parity and mask/version selection;
DataMatrix tests cover Base256 length/position randomization and control conflicts;
PDF417 controls cover macro/ECI layout and Unicode codeword error correction;
Aztec tests cover word stuffing, compact capacity and flags. GS1 dictionary/linter
fixtures cover longest-AI matching, predefined length versus variable separators,
associations, allocation tables and bounded decimal checks. Image tests cover
rotation/inversion, diagnostics, correction admission, payload redaction and
independent software peers. Tests/helpers/fuzz/bench/examples are inventoried
below; inspected evidence is not a newly executed campaign.

## Reproducible complete inventory

Reproduce with `git ls-files '*.go' 'go.mod' 'go.sum'`, JSON catalogs and scans
for package clauses, declarations, entrypoints, directives and generated markers.

| Tracked Go file | Package | Reviewed role / selected action |
| --- | --- | --- |
| `aztec/aztec.go` | `aztec` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `aztec/aztec_test.go` | `aztec_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeSupportsAutomaticAndForcedLayers, TestEncodeRoundTripsThroughImageDecoder, TestEncodeRejectsUnsafeOptions, TestEncodeAcceptsExactLimits, TestEncodeReturnsStableValidationError, TestEncodeSelectsSmallestAutomaticCompactLayer, TestEncodeAppliesConfiguredQuietZone, TestEncodeSupportsGS1FNC1, TestEncodeSupportsECI) |
| `barcode/core.go` | `barcode` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `barcode/core_test.go` | `barcode_test` | Existing test/helper evidence, unchanged executable tokens (TestMatrixAndSymbolDoNotAliasCallerData, TestDecodeResultValidatesMetadataAndReturnsDefensiveValues, TestSymbolAndCapabilityRejectUnknownOrMissingRepresentations, TestMatrixRejectsInvalidDimensions, TestMatrixDistinguishesMaximumSafeProductFromOverflow, TestCapabilitiesAreExplicitForEveryKnownFormat, TestCapabilityMetadataIdentifiesExactGoverningEditions, TestCapabilitiesReflectSoftwareScope, TestBarsAndDecodeResultDoNotAliasCallerData, TestBarsRejectInvalidRunsAndOverflow, TestBarsAcceptMaximumSafeWidth, TestDecodeResultAcceptsExactMetadataLimits, TestSymbolAcceptsExactlyOneLogicalRepresentation) |
| `benchmark_test.go` | `barcode_test` | Existing test/helper evidence, unchanged executable tokens (BenchmarkEncodeQR, BenchmarkEncodeCode128, BenchmarkRenderQR, BenchmarkDecodeQR, BenchmarkDetectNoSymbol, BenchmarkDecodeEncodedQR, BenchmarkRejectMalformedImage) |
| `codabar/codabar.go` | `codabar` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `codabar/codabar_test.go` | `codabar_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeUsesExplicitStartAndStopCharacters, TestEncodeRejectsInvalidPayloadGuardsAndDimensions, TestEncodePreservesExplicitDimensions, TestEncodeAcceptsPayloadAndOptionBoundaries) |
| `code128/code128.go` | `code128` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `code128/code128_internal_test.go` | `code128` | Existing test/helper evidence, unchanged executable tokens (TestMatrixRunsAllocationBudget, TestMatrixRunsQuietZoneAllocationBudget, TestCodeSetNames) |
| `code128/code128_test.go` | `code128_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeCodeSetBMatchesIndependentLogicalVector, TestEncodeDefaultsAreStandardsSafeAndAutoSwitchCodeSets, TestEncodeSupportsEachForcedCodeSet, TestEncodeGS1AddsFNC1AndReportsGS1Format, TestEncodeRawFNC1SupportsLegacyPayloadsWithoutGS1Validation, TestEncodeGS1AcceptsValidatedStructuredElements, TestEncodeRejectsUnsafeOptionsAndPayloads, TestEncodeReturnsDirectValidationErrors, TestEncodeAcceptsPayloadCharacterAndOptionBoundaries) |
| `code39/code39.go` | `code39` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `code39/code39_internal_test.go` | `code39` | Existing test/helper evidence, unchanged executable tokens (TestFullASCIIExpansionCoversEveryASCIIByte, TestChecksumRejectsEmptyAndNonBaseCharacters, TestFullASCIIExpansionMatchesUpperPunctuationBoundary, TestFullASCIIExpansionMatchesRangeBoundaries) |
| `code39/code39_test.go` | `code39_test` | Existing test/helper evidence, unchanged executable tokens (TestChecksumMatchesModulo43Vector, TestChecksumAcceptsAlphabetBoundary, TestEncodeSupportsChecksumsAndFullASCII, TestEncodeRejectsInvalidInput, TestEncodeAcceptsCapacityCharacterAndOptionBoundaries) |
| `code93/code93.go` | `code93` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `code93/code93_test.go` | `code93_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeAddsMandatoryChecksumsAndSupportsFullASCII, TestEncodeRejectsInvalidInput, TestEncodePreservesExplicitDimensions, TestEncodeAcceptsPayloadAndOptionBoundaries) |
| `datamatrix/datamatrix.go` | `datamatrix` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `datamatrix/datamatrix_internal_test.go` | `datamatrix` | Existing test/helper evidence, unchanged executable tokens (TestECICodewordsMatchAssignmentWidthBoundaries, TestPaddingRandomizationMatchesECC200Vectors, TestControlledEncodingDistinguishesBase256LengthLimit, TestStructuredAppendSequenceCodeword, TestControlledMatrixGoldenModules) |
| `datamatrix/datamatrix_test.go` | `datamatrix_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeCreatesECC200Shapes, TestEncodeHonorsDimensionAndQuietZoneConstraints, TestEncodeRoundTripsThroughImageDecoder, TestEncodeRejectsUnsafeOptions, TestEncodeReturnsDirectValidationErrors, TestEncodeAcceptsExactDimensionAndQuietZoneLimits, TestEncodeSupportsStructuredAppendBoundaries, TestEncodeAcceptsECC200CapacityBoundaries, TestEncodeGS1EmitsFNC1ControlCodeword, TestEncodeSupportsECIAssignmentWidths, TestEncodeECIRoundTripsKnownCharacterSets, TestEncodeSupportsMacro05And06) |
| `ean/ean.go` | `ean` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `ean/ean_internal_test.go` | `ean` | Existing test/helper evidence, unchanged executable tokens (TestSupplementParityAndRunConversion) |
| `ean/ean_test.go` | `ean_test` | Existing test/helper evidence, unchanged executable tokens (TestEncode13MatchesIndependentLogicalVector, TestEncodeCalculatesMissingCheckDigit, TestEncodeSupportsTwoAndFiveDigitSupplements, TestEncodeRejectsInvalidDataAndUnsafeQuietZones, TestEncodeReturnsDirectOptionValidationErrors, TestEncodeAcceptsDimensionAndSupplementBoundaries) |
| `encoder_fuzz_test.go` | `barcode_test` | Existing test/helper evidence, unchanged executable tokens (FuzzPayloadEncoders) |
| `example_test.go` | `barcode_test` | Existing test/helper evidence, unchanged executable tokens (Example_qrCode, Example_gs1128, Example_imageDecode, Example_svg) |
| `gs1/checkdigit.go` | `gs1` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `gs1/checkdigit_test.go` | `gs1_test` | Existing test/helper evidence, unchanged executable tokens (TestCalculateAndValidateCheckDigit, TestCheckDigitValidationIsStrict) |
| `gs1/element.go` | `gs1` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `gs1/element_fuzz_test.go` | `gs1_test` | Existing test/helper evidence, unchanged executable tokens (FuzzParseElementStrings, FuzzCheckDigits) |
| `gs1/element_internal_test.go` | `gs1` | Existing test/helper evidence, unchanged executable tokens (TestDefinitionLoadFailuresRemainObservable, TestDictionaryGrammarSupportsRangesFlagsAndOptionalComponents, TestDictionaryGrammarSumsRequiredComponents, TestDictionaryGrammarSupportsAssociationAttributes, TestDictionaryGrammarRejectsMalformedEntries, TestDictionaryTokenClassificationAndComponents, TestValueValidationCoversSupportedCharacterClasses, TestDictionaryHelpersCoverExactBoundaries, TestDefinitionMatchingPrefersLongestAI) |
| `gs1/element_test.go` | `gs1_test` | Existing test/helper evidence, unchanged executable tokens (TestParseBracketedElementString, TestParseRawElementStringUsesPredefinedLengthsAndFNC1, TestElementStringSerializesRawAndBracketedSyntax, TestParserUsesPinnedDictionaryRanges, TestParserEnforcesRequiredAndExcludedAssociations, TestParserRejectsUnknownInvalidAndOverLimitData, TestParsersRejectMalformedStructureAndLimits, TestParserAcceptsExactInputLimit) |
| `gs1/linter.go` | `gs1` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `gs1/linter_data.go` | `gs1` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `gs1/linter_internal_test.go` | `gs1` | Existing test/helper evidence, unchanged executable tokens (TestDictionaryLinterDispatch, TestLinterHelperMalformedInputs, TestAlphaCheckPairBoundaries, TestCouponPositiveOfferBoundaries, TestCouponCodeGrammar, TestLinterExactBoundaries, TestLinterMaximumLengthChecks) |
| `imagedecode/decode.go` | `imagedecode` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `imagedecode/decode_bytes_test.go` | `imagedecode_test` | Existing test/helper evidence, unchanged executable tokens (TestDecodeEncodedBytesRejectsOversizeAndMalformedInput, TestDecodeEncodedBytesHonorsPreexistingCancellation, TestDecodeEncodedBytesStopsBetweenConfigurationAndDecode) |
| `imagedecode/decode_fuzz_test.go` | `imagedecode_test` | Existing test/helper evidence, unchanged executable tokens (FuzzDecodeBoundedImages) |
| `imagedecode/decode_test.go` | `imagedecode_test` | Existing test/helper evidence, unchanged executable tokens (TestDecodeQRAndCode128Images, TestDecodeCode39FullASCIIWithChecksum, TestDecodeEverySupportedFormat, TestDecodeDataMatrixControls, TestDecodeEANAndUPCSupplementsReportsExtensionMetadata, TestDecodeReportsRotationAndPayloadLimits, TestDecodeRejectsInvalidImagesAndLimits, TestDecodeEncodedBoundsCompressedImages, TestDecodeSupportsInvertedImages, TestDecodeEnforcesBoundsBeforeImageAllocation, TestDecodeEnforcesCallerTimeBudget, TestDecodeRejectsUnsupportedFormatsAndCandidateLimits) |
| `imagedecode/independent_reader_test.go` | `imagedecode_test` | Existing test/helper evidence, unchanged executable tokens (TestSymbolsDecodeWithIndependentReaders) |
| `imagedecode/independent_writer_test.go` | `imagedecode_test` | Existing test/helper evidence, unchanged executable tokens (TestDecodeIndependentWriters) |
| `imagedecode/internal_test.go` | `imagedecode` | Existing test/helper evidence, unchanged executable tokens (TestTwoDReaderDecodeWithoutHints, TestTwoDReaderContainsDependencyPanics, TestControlledDataMatrixECIWidths, TestDataMatrixECIExactBoundaries, TestControlledDataMatrixSequenceMacroAndFallbackECI, TestTwoDReaderFallbackControls, TestReadersPropagateBlackMatrixErrors, TestTwoDMetadataCopiesSequenceAndByteSegments, TestDecodeChecksDeadlineBetweenAttempts, TestDecodeRotatesInvertedInput, TestControlledDataMatrixRejectsMalformedCodewords, TestBase256LengthFormsAndFailures, TestOrientationChecksumAndFormatMappings, TestImageAndLimitExactBoundaries) |
| `imagedecode/pdf417_reader.go` | `imagedecode` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `imagedecode/pdf417_reader_test.go` | `imagedecode` | Existing test/helper evidence, unchanged executable tokens (TestPDF417ReaderImplementsReaderLifecycle) |
| `imagedecode/redaction_test.go` | `imagedecode_test` | Existing test/helper evidence, unchanged executable tokens (TestDecodeEncodedDoesNotExposeInputErrors) |
| `imagedecode/robustness_test.go` | `imagedecode_test` | Existing test/helper evidence, unchanged executable tokens (TestDecodeQRDocumentedImageDegradationThresholds, TestDecodeMultipleSymbolsReturnsOneDecodableCandidate, TestDecodeEnforcesCorrectionBudget) |
| `imagedecode/twod_reader.go` | `imagedecode` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `internal/aztecencoder/encoder.go` | `encoder` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `internal/aztecencoder/encoder_test.go` | `encoder` | Existing test/helper evidence, unchanged executable tokens (TestEncodeMatchesPinnedZXingImplementation, TestEncodeMatchesPinnedZXingCapacityBoundaries, TestEncodeExercisesAutomaticCompactAndFullLayers, TestEncodeRejectsInvalidLayersAndCapacity, TestHighLevelEncodingExercisesModesPairsAndBinaryShift, TestInternalModeAndStuffingEdges, TestHighLevelExactTablesFlagsAndBinaryLimits, TestLayerFitAndFullMatrixExactBoundaries, TestUnsupportedWordSizePanics, TestEveryLatchSequenceAndShiftDecision) |
| `internal/aztecencoder/highlevel.go` | `encoder` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `internal/linear/linear.go` | `linear` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `internal/linear/linear_test.go` | `linear` | Existing test/helper evidence, unchanged executable tokens (TestEncodeBuildsExactAlternatingRuns, TestEncodePropagatesWriterAndBarValidationErrors, TestEncodeMergesTrailingLightRunWithQuietZone) |
| `internal/pdf417encoder/barcode_matrix.go` | `pdf417encoder` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `internal/pdf417encoder/controls_test.go` | `pdf417encoder` | Existing test/helper evidence, unchanged executable tokens (TestEncodeECIRanges, TestEncodeMacroControlBlock, TestEncodeMacroAcceptsFieldBoundaries, TestEncodeMacroRejectsInvalidFields) |
| `internal/pdf417encoder/encoder_test.go` | `pdf417encoder` | Existing test/helper evidence, unchanged executable tokens (TestEncoderMatchesPinnedZXingImplementation, TestBarcodeRowsAndMatrices, TestErrorCorrectionContracts, TestHighLevelCompactionPaths, TestHighLevelCharacterAndRunBoundaries, TestBinaryCompactionCodewordBoundaries, TestSymbolCodewordCapacityBoundaries, TestEncodeCharPatterns, TestEncoderLayoutAndFailurePaths) |
| `internal/pdf417encoder/errorcorrection.go` | `pdf417encoder` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `internal/pdf417encoder/highlevel.go` | `pdf417encoder` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `internal/pdf417encoder/pdf417enc.go` | `pdf417encoder` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `itf/itf.go` | `itf` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `itf/itf_internal_test.go` | `itf` | Existing test/helper evidence, unchanged executable tokens (TestEncodeBarsPropagatesWriterErrors) |
| `itf/itf_test.go` | `itf_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeInterleavesNumericPairs, TestEncode14CalculatesCheckDigitAndAddsBearerBars, TestEncode14SupportsCompleteValueAndFramedBearerBars, TestEncodeRejectsInvalidInput, TestEncodeAcceptsPayloadAndDimensionBoundaries, TestEncode14BearerGeometryMatchesRequestedStyle, TestEncodePreservesExplicitDimensions) |
| `pdf417/pdf417.go` | `pdf417` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `pdf417/pdf417_test.go` | `pdf417_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeSupportsErrorCorrection, TestEncodeSupportsECIAndMacroControlBlocks, TestEncodeSupportsAllCorrectionLevelsAndQuietZones, TestEncodeSupportsCompactionAndLayoutControls, TestEncodeRoundTripsThroughImageDecoder, TestEncodeRejectsUnsafeOptions, TestEncodeAcceptsOptionBoundaries, TestEncodeDistinguishesPayloadSafetyFromSymbolCapacity, TestEncodeRejectsNegativeMacroOptionalNumbers) |
| `qr/qr.go` | `qr` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `qr/qr_fuzz_test.go` | `qr_test` | Existing test/helper evidence, unchanged executable tokens (FuzzEncodeOptions) |
| `qr/qr_internal_test.go` | `qr` | Existing test/helper evidence, unchanged executable tokens (TestModeAndCorrectionMappings, TestOptimizedModeClassifiesSegmentPlans, TestECIAssignmentRejectsNegativeValue, TestECIAssignmentEncodingMatchesWidthBoundaries, TestStructuredPayloadExcludesControlSegments, TestExtendedSegmentsEncodeControlMetadata, TestEncodeSegmentsRejectsInvalidVersionAndCapacity) |
| `qr/qr_test.go` | `qr_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeSupportsEveryErrorCorrectionLevel, TestEncodeHonorsMaskECIAndGS1Controls, TestEncodeSupportsSecondPositionFNC1, TestEncodeSupportsStructuredAppendHeaders, TestEncodeStructuredSplitsAndCalculatesParity, TestEncodeStructuredAcceptsExactlyTwoParts, TestEncodeKanjiUsesKanjiMode, TestEncodeAutoOptimizesMixedModePayloads, TestEncodeGS1AcceptsValidatedElements, TestEncodeRejectsUnsafeOrUnsupportedOptions, TestEncodeSupportsAutomaticModesAndECIAssignmentWidths, TestEncodeRejectsModeCharsetAndCapacityMismatches, TestEncodeExtendedControlsPreserveRequestedMetadata, TestEncodeStructuredValidatesSequenceOptions, TestEncodeStructuredSupportsFNC1Modes, TestEncodeAcceptsExactOptionBoundaries, TestEncodeFirstPositionFNC1WithForcedMode, TestEncodeStructuredAcceptsMaximumPayloadBoundary) |
| `reliability_test.go` | `barcode_test` | Existing test/helper evidence, unchanged executable tokens (TestConcurrentEncodeRenderAndDecodeIsDeterministic) |
| `render/render.go` | `render` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `render/render_internal_test.go` | `render` | Existing test/helper evidence, unchanged executable tokens (TestBoundedProduct) |
| `render/render_test.go` | `render_test` | Existing test/helper evidence, unchanged executable tokens (TestImageUsesIntegerModuleScalingAndExactColors, TestRenderMatrixSVGTransparencyAndDefaults, TestRenderRejectsInvalidSymbolsAndWriterFailures, TestSVGPropagatesEveryWriterFailurePosition, TestPNGAndSVGRenderMatrixAndBars, TestRenderedOutputsMatchGoldenChecksums, TestRenderRejectsOverflowAndExplicitLimitViolations, TestRenderAcceptsSinglePixelAndRejectsAsymmetricOverflow, FuzzRenderLogicalMatrices) |
| `render_fixtures_test.go` | `barcode_test` | Existing test/helper evidence, unchanged executable tokens (TestRenderFixtureGoldensCoverEveryFormat) |
| `security_test.go` | `barcode_test` | Existing test/helper evidence, unchanged executable tokens (TestInvalidInputErrorsAreClassifiedAndPayloadRedacted, TestPDF417TextErrorDoesNotExposeUnsupportedPayload) |
| `specification/embed.go` | `specification` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `specification/embed_test.go` | `specification_test` | Existing test/helper evidence, unchanged executable tokens (TestGS1SyntaxDictionaryIsEmbedded) |
| `specification/manifest_test.go` | `specification_test` | Existing test/helper evidence, unchanged executable tokens (TestManifestProvenanceIsWellFormed) |
| `upc/upc.go` | `upc` | Source algorithms/contracts, public fields/methods and canonical package intent |
| `upc/upc_internal_test.go` | `upc` | Existing test/helper evidence, unchanged executable tokens (TestSupplementParityAndRunConversion) |
| `upc/upc_test.go` | `upc_test` | Existing test/helper evidence, unchanged executable tokens (TestEncodeACalculatesAndValidatesCheckDigit, TestEncodeECalculatesCheckDigitAndExpandsToUPCA, TestExpandECoversEachCompressionRule, TestUPCSupplementsExtendLogicalSymbol, TestUPCFiveDigitSupplementExtendsLogicalSymbol, TestUPCRejectsInvalidInput, TestUPCAcceptsNumberSystemAndDimensionBoundaries) |

## Verification disposition

Executed locally with Go 1.27.2: formatting on 33 changed Go paths; parsed all
76 tracked Go files and proved exact noncomment scanner tokens (including
literals and implicit semicolons), position/comment-free ASTs and directives;
403 name-led public declarations/fields/methods, 19 canonical source package
comments and 644 parsed HTML comment groups. The root example package remains
test-only and its four executable example bodies are unchanged. Native
`go doc -all ./<catalog-directory>` rendered all 19 source packages. The
final Code 93 wording was formatted/rendered again and all-file structural/HTML
proof plus all 16 typed digests and native source/docs/specification checks
were refreshed; the other 18 unchanged renders were reused. The selected
apidiff baseline, module/dependency/catalog bytes and all frozen inputs remain
unchanged; no prose snapshot or source artifact is regenerated.

The shared native source checker reports zero findings for 76 files/one module;
local documentation and offline specification checks pass. Native typed Go
canonicalization reproduces all 16 baseline/current decision digests; only
BARCODE-DEC-005/006/012/014 append one actual current digest each, retaining all
historical prefixes, IDs, pins, evidence paths and other typed fields. Existing
restricted-source qualifications and software test evidence are retained, not
represented as newly executed conformance.

Runtime/example, peer/service, race, fuzz, benchmark, online authority, consumer
and release campaigns are not selected: executable behavior, examples and all
fixtures are unchanged. Changed-content and staged secrets, affected Markdown
links/frozen-file identity, and one ordinary complete final-diff truth/usefulness
inspection complete the Tier A verification; their result belongs in the local
handoff, not a claim of physical-device or released-software qualification.
All three disposable Go caches, helper files and rendered copies were removed
immediately after the bounded session, including checked absence.

Remaining owner follow-ups are runtime work, not documentation blockers: original
QR structured-append Index upper validation before four-bit packing; caller-sized
linear margin arithmetic/allocation (and Code 39 pre-admission expansion); core
finite-confidence and content/diagnostic admission if applications require them.
Code 93 also discards its selected writer error locally, so full-ASCII expansion
capacity deserves focused runtime investigation rather than a promised classified
failure for every admitted input. PDF417 FileID sign acceptance follows Atoi;
applications needing digit-only identifiers must impose that constraint. No new
standards edition, hardware readability, peak-heap, hard-time, universal checksum
or payload-redaction guarantee is claimed.
