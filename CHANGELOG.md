# Changelog

All notable changes are documented here. The project follows semantic
versioning.

## Unreleased

### Changed

- Adopt the pinned `go-library-tools` v1.2.0 CLI and reusable workflow so CI
  enforces specification decisions, conformance bindings, source monitoring,
  and change control while retaining package-owned policy and verification
  evidence.
- Adopt the checksum-verified `go-library-tools` v1.3.0 CLI, schema-v2 cohesion
  metadata, and repository-local cohesion gate while retaining package-owned
  source and evidence.
- Adopt the checksum-verified `go-library-tools` v1.4.0 CLI and immutable
  workflow. The complete local contract covers configuration, inventory,
  cohesion, repository, workflow, online specification, and implementation
  checks; hosted CI enforces the repository, module, and online specification
  contracts. Published modules now resolve from the public Go proxy before the
  immutable bootstrap fallback, preventing fallback bytes from shadowing public
  releases.
- Adopt the checksum-verified `go-library-tools` v1.5.0 CLI and immutable
  workflow so conditional source monitoring accepts only an explicitly reviewed
  publisher-edge denial while hosted repository checks also enforce workflow
  validation.
- Pin official ISO per-standard RSS metadata and public ANSI catalogue records
  by content. Conditional monitoring verifies the ANSI content when available
  and accepts only its reviewed HTTP 403 edge denial otherwise. The licensed
  normative publications remain outside the repository and are not claimed as
  reviewed.
- Classify ANSI/AIM BC5-1995 as the withdrawn historical Code 93 compatibility
  target and the separately catalogued 2000 AIM publication as unclaimed until
  licensed review establishes its normative relationship.
- Correct Codabar provenance and public capability metadata to the withdrawn
  ANSI/AIM BC3-1995 compatibility target while leaving wire behavior and the
  unsupported optional-checksum policy unchanged.
- Pin AIM public-review monitoring to the publisher's canonical JSON record so
  generated HTML nonces do not obscure substantive authority changes.

### Documentation

- Replace the archived monorepo link with package-owned documentation.
- Link the public package suite to the immutable v1.3.0 Golib ecosystem and
  package-family selection guidance.
- Link the public package suite to the immutable v1.4.0 Golib ecosystem and
  Domain utilities family guidance.
- Link the public package suite to the immutable v1.5.0 Golib ecosystem and
  Domain utilities family guidance.

### Specification Decisions

- BARCODE-DEC-001 sha256:f11a763d3fbf349801f7dd939c319bc48cda1371855404983a5b963e849150a0:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-001 sha256:82a8bec27704f43f2af57f159a8706789f46c8f21375196512bfd731bc647176:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-001 sha256:b6afd03f1fd8a4229e3ffe66a2f37827f9321829b019577fe95cc2267162500f:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-002 sha256:4dd3716c04a9146f2892f0f151f829de5a30ca5e74c2b55a2efb9f314a8acce5:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-003 sha256:0ca6140d4b9e09de91ca55713a10899ac7766b03fcdb62fe86f6b7e94c25e8e5:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-004 sha256:e85527ffd2fbee6c5f48a296efd98c1dd5fe5b62e80297da28370725e0e806a2:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-005 sha256:5509d9b90142b2fca8cb8320822095112f57baedb3c11fceb1a35d41f28d5798:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-006 sha256:c572c740bf280c45bdd512820f8ce6917332b76c817f07bdd425316a08076e32:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-007 sha256:ab3e1ddc7e019756dd86648f115241fc7d8287800e7be885a438711471173602:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-007 sha256:7b7141a13ebd83094462b95885b81716e5e988163fc002c399ec1915e5649d98:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-007 sha256:a79d5c4b4a910e962766aae4fcfe31ad1190a816624653fbfcbf4a5da3c32d3d:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-007 sha256:bf469f172668ace09189add762f0f6667e367efea13434dcad9fbad58d3c96f9:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-008 sha256:4f87656d572006036d118439f8fe950aba97ed0fdea0befefb98615b816408f2:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-009 sha256:11169c5c9775e7e2870637141166796cbd7ac83948ce88b3e7fcf28c5a5eb283:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-010 sha256:c2e0275df87a31c7c810729a79ed7e11439a34a465a348989946105f0c42416e:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-011 sha256:2d9880efc00d803422a80d562d29860c7d3ba5f67513c6051a6f4e028ce08537:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-012 sha256:7981dd70757727139bdabab90bb8609d8e89ec3acf0504719d046799b77f9226:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-013 sha256:d238d62376d4fbc5714414bf3628b6949f3d93ffb2e3979b457cef0b47d2747a:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-014 sha256:2d5a3129516cab975b90111f12f8cfff19e22222f41af8e67dd348d4738f4fc8:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-014 sha256:fb3ef64be8a6f098cb66185c13e94cc24d94a4840e4046c2bae028bb3b4175a8:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-015 sha256:aa729a65e1365d254737da972cd96195e8aee4ff5ceb7cb32376ba827c45bb56:
  [Decision register](docs/specification-decisions.md).
- BARCODE-DEC-016 sha256:ff06cc03def08bf7bb539b6a6aca6b9abc635e143319d92f8e335895b6b7b431:
  [Decision register](docs/specification-decisions.md).

## 1.0.0 - 2026-08-25

### Changed

- Upgrade `golang.org/x/text` to v0.41.0 so the dependency graph no longer
  contains GO-2026-5970.
- Delegate local mutation checks to the canonical exact-100 repository runner
  instead of configurable package-local thresholds.
- Activate the complete barcode gate in hosted CI with path filtering,
  scheduled fuzzing, benchmark artifacts, and manual dispatch.
- Shard hosted mutation testing by package directory while retaining the full
  local mutation command and thresholds.
- Scope capability advertising and interoperability evidence to independently
  tested software behavior, excluding hardware integration and certification.
- Refresh the reviewed zero-mutant identity for the extracted specification
  package without weakening the exact mutation contract.

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-barcode` identity while preserving its documented API and behavior.

### Documentation

- Link the package README to package-owned documentation.
- Document classified failures for unsupported ECI, invalid GS1, and
  checksum-failing image candidates.

### Added

- Add deterministic logical, PNG, and SVG rendering fixtures for every
  implemented format, with reproducible checksum validation.
- Link contribution guidance directly to the specification decision register
  and its provenance and evidence requirements.
- Add an auditable specification-decision register, exact standards
  provenance, and a conformance gate that binds format policy to executable
  evidence.

### Fixed

- Preserve the GS1-128 AIM symbology identifier as a decode diagnostic instead
  of including it in the decoded GS1 payload.
- Contain panics from additive two-dimensional decoder dependencies at the
  candidate boundary.
- Correct the public ISO edition metadata for Code 39 and Data Matrix and the
  authoritative ISO catalogue URLs for ITF and PDF417.
- Reject invalid PDF417 text controls and over-capacity symbols with classified
  errors instead of allowing low-level encoding panics.
- Keep GS1 parsing, image decoding, and barcode layout boundary processing
  bounded and deterministic for hostile or extreme inputs.

### Compatibility

- Added an explicit raw-FNC1 Code 128 compatibility mode for legacy payloads
  that are not valid GS1 element strings.
- Added a pinned module export baseline so incompatible public API changes
  fail the canonical repository gate.

- Added immutable barcode core models and bounded raster/SVG rendering.
- Added complete phase-one linear formats and QR control surfaces.
- Added Data Matrix, PDF417, and Aztec encoding, decoding, controls, and
  reciprocal interoperability evidence.
- Added pinned GS1 parsing with component and AI association validation.
- Added offline GS1 allocation, check-pair, and coupon content validation.
- Added reciprocal GS1-128, UPC, supplement, and ITF-14 interoperability
  evidence, including explicit ITF-14 decode-format preservation.
- Preserved Data Matrix structured-append sequence and file identifiers in
  decode diagnostics.
- Added bounded raw and encoded image decoding, hostile-input fuzzing,
  deterministic render goldens, benchmarks, and reproducible CI gates.
- Added normative-to-executable evidence inventories for every format.
