# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's documentation.
2. Run ordinary Go tests for the affected package before editing.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must update the
[decision register](docs/specification-decisions.md) and follow the
[specification governance contract](AGENTS.md#design), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Selected mutation gates must finish with zero surviving viable mutants.

Do not fork shared verification workflows, add permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for development and ordinary CI; the public minimum remains
Go 1.27.0. An ordinary standalone checkout is sufficient for tests and a
focused active fuzz smoke, without sibling repositories or goal ledgers:

```bash
go test ./...
go test ./imagedecode -run '^$' -fuzz '^FuzzDecodeBoundedImages$' -fuzztime=2s
```

Select the native fuzz target and package affected by your change. Ordinary
tests run deterministic seeds and saved corpora. For shared checks, install
the public `golib` v2.4.0 release declared in `.golib.yaml` and verify its
checksums:

```bash
make check
make ci
```

Both Make routes use bounded local checks, including one race run of the
existing concurrent encode/render/decode assertion. Base-free Make checks
run deterministic seeds/corpora; changed-source CI selects active fuzz smoke.
`make ci` checks the specification register offline. The sole root module
declares no services or test tags; pinned independent software readers and
writers do not require devices, hardware certification or sibling checkouts.

Optional mutation analysis can select one affected package:

```bash
golib mutation --module . --package ./imagedecode
```

For format, decoder or rendering changes, select the existing conformance
and reciprocal independent-reader/writer checks:

```bash
make -f verification/package.mk conformance
make -f verification/package.mk interoperability
```

Use `make check-full` for explicit online specification validation and full
assurance with all unchanged native fuzz budgets, dependency policy, API,
canonical render-fixture checks, conformance, interoperability and benchmarks.
Preserve pinned standard editions and the distinction between software profiles
and physical device certification. Select directly affected public consumers
and relevant release qualification, including `golib release check` and
`golib release dry-run --all`, when their contract changes. Full campaigns
remain available for selected scheduled/release risks; ordinary scheduled CI
stays bounded and separately monitors specification sources online. Report
unavailable or failing checks without claiming release readiness.
CI/configuration-only changes do not require a library release.

## Selected Monthly Assurance

The owner-level [selected-assurance caller](.github/workflows/assurance.yml) is
configured for main only, on day 10 of each month at 05:07 UTC. It selects one
of the 6 current declared root-module fuzz targets using the shared UTC-week
rotation, with 30 seconds of active exploration. The selected risks are:
payload/QR options, GS1 element strings and check digits, bounded image
decoding and logical matrix rendering.
These are software input/output checks, not physical device or printing
certification. Independent reader/writer interoperability, standard-edition
provenance and full conformance remain explicit.

Base-free ordinary checks execute deterministic seeds, not active exploration;
committed-range CI selects affected active smoke. Manual assurance requires a
declared selector from the root-module fuzz operation in `.golib.yaml`, for
example `./imagedecode::^FuzzDecodeBoundedImages$`. A bare anchored target is
accepted only when unique across declared packages. `.::^FuzzPayloadEncoders$`
selects the legitimate test-only root package without adding a production record.
Monthly selection samples one identity; it does not promise exhaustive rotation.
Empty, unknown and malformed selections refuse. The shared job bounds setup
and execution to ten minutes, active workers and CPU to two, owned memory to
2 GiB, native output to 1 MiB and private-cache lifetime. Failure, cancellation,
exhausted budgets and cleanup refusal remain attributable failures. This lane
does not run mutation or replace Required, original full fuzz budgets,
explicit `make check-full`, or relevant release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.

## Source Documentation

Install the checksum-pinned public verifier declared in `.golib.yaml`.
`make source-docs-check` validates objective Go comment, marker, link and
generated-producer structure; `make check`, `make check-full` and `make ci`
require it before their existing checks. Shared CI runs the same strict
contract and includes its result in Required. This is structural enforcement,
not a substitute for source-truth review or behavioral verification.
