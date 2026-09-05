# Security policy

Do not open a public issue for a suspected vulnerability. Report it privately
through [GitHub Security Advisories for
`faustbrian/go-barcode`](https://github.com/faustbrian/go-barcode/security/advisories/new).
Provide the affected package and version, a minimized input when safe, expected
and actual resource bounds, and whether the issue can panic, hang,
over-allocate, or bypass validation.

Do not include sensitive payloads, private barcode images, credentials, or
customer data. Decoded content is always untrusted and must not be executed or
followed.

The latest released v1 minor version receives security fixes. The `main` branch
contains unreleased development and is not a substitute for a supported
release. See [security and resource limits](docs/security.md) for the complete
threat boundary.
