// Package specification exposes the immutable embedded GS1 syntax dictionary,
// release 2026-01-27. It redistributes no restricted ISO/AIM text and performs no
// runtime downloads. Read GS1SyntaxDictionary; authenticated regeneration uses
// scripts/sync-gs1-dictionary.sh, not a generated Go-code producer.
// Exact source/license/digest provenance lives in specification/README.md.
package specification

import _ "embed"

//go:embed gs1/gs1-syntax-dictionary.txt
var gs1SyntaxDictionary string

// GS1SyntaxDictionary returns the pinned GS1 syntax dictionary bytes.
func GS1SyntaxDictionary() string { return gs1SyntaxDictionary }
