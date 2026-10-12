// Package pdf417 constructs copied logical PDF417 matrices via Encode, with
// explicit compaction, correction, layout, ECI and caller-managed macro metadata.
// The target is ISO/IEC 15438:2015 edition 3; no full licensed-clause/device or
// macro assembly claim is made (BARCODE-DEC-010). NOTICE.md records the locally
// maintained zxinggo v0.1.0 encoder derivation. It owns no external I/O lifecycle.
package pdf417

import (
	"errors"
	"fmt"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/faustbrian/go-barcode/v2/internal/pdf417encoder"
)

const (
	defaultQuietZone = 2
	maxPayloadBytes  = 4096
)

var (
	// ErrInvalidInput and ErrUnsupported classify PDF417 option failures.
	ErrInvalidInput = errors.New("pdf417: invalid input")
	// ErrUnsupported is a reserved unsupported-option classifier. Current Encode
	// validation/encoder failures use ErrInvalidInput; no path returns it directly.
	ErrUnsupported = errors.New("pdf417: unsupported option")
)

// Compaction selects a PDF417 high-level encoding strategy.
type Compaction uint8

const (
	// Automatic and the other constants select PDF417 compaction modes.
	Automatic Compaction = iota
	// Text forces text compaction.
	Text
	// Byte forces byte compaction.
	Byte
	// Numeric forces numeric compaction.
	Numeric
)

// ErrorCorrection selects one of the nine PDF417 correction levels.
type ErrorCorrection uint8

const (
	// DefaultErrorCorrection and the other constants select correction levels.
	DefaultErrorCorrection ErrorCorrection = iota
	// Level0 adds two correction codewords.
	Level0
	// Level1 adds four correction codewords.
	Level1
	// Level2 adds eight correction codewords.
	Level2
	// Level3 adds sixteen correction codewords.
	Level3
	// Level4 adds thirty-two correction codewords.
	Level4
	// Level5 adds sixty-four correction codewords.
	Level5
	// Level6 adds 128 correction codewords.
	Level6
	// Level7 adds 256 correction codewords.
	Level7
	// Level8 adds 512 correction codewords.
	Level8
)

// Options controls PDF417 compaction, correction, layout, and quiet zone.
type Options struct {
	// Compaction selects Automatic (zero), Text, Byte or Numeric.
	Compaction Compaction
	// ErrorCorrection is zero for Level2; explicit levels 0..8 add 2..512 parity
	// codewords. The enum Level0 has value 1, not the zero/default enum.
	ErrorCorrection ErrorCorrection
	// Compact requests truncated right-side PDF417 layout.
	Compact bool
	// MinRows is zero for 3; a nonzero bound must be positive and at most 90.
	MinRows int
	// MaxRows is zero for 90; explicit bounds must agree with MinRows.
	MaxRows int
	// MinColumns is zero for 1; a nonzero bound must be positive and at most 30.
	MinColumns int
	// MaxColumns is zero for 30; explicit bounds must agree with MinColumns.
	MaxColumns int
	// QuietZone is zero for 2 modules, otherwise positive and at most 256.
	QuietZone int
	// ECI is zero for no header, otherwise below 811800; emits metadata, not transcoding.
	ECI int
	// Macro is optional; scalar pointers are read during encoding and deeply
	// cloned into Symbol metadata. Bound caller-supplied text/FileID before invoking
	// this API: the 4096-byte payload cap does not bound macro fields.
	Macro *Macro
}

// Macro describes a Macro PDF417 control block. FileID contains one or more
// three-character decimal groups parsed as codewords 0..899. The current parser
// delegates each group to strconv.Atoi, so callers should supply unsigned ASCII
// digit groups rather than rely on stricter sign rejection. Other fields are
// declarative metadata, not file-content or cross-segment consistency checks.
type Macro struct {
	// SegmentIndex is 0..99998; the caller owns sequence identity/order.
	SegmentIndex int
	// FileID is nonempty three-character groups parsed as 0..899 codewords.
	// It is caller-managed metadata, not a file path or authenticity claim.
	FileID string
	// LastSegment emits the macro terminator; it does not assemble other parts.
	LastSegment bool
	// FileName optionally encodes text-compaction metadata, never opens a file.
	FileName string
	// SegmentCount optionally supplies a nonnegative declared count;
	// no cross-part consistency or index/count relation is validated.
	SegmentCount *int
	// Timestamp optionally supplies a nonnegative numeric value;
	// no timezone/unit conversion is performed by this encoder.
	Timestamp *int64
	// Sender optionally encodes text-compaction metadata, not authenticated identity.
	Sender string
	// Addressee optionally encodes text-compaction metadata, not routing or delivery.
	Addressee string
	// FileSize optionally supplies a nonnegative number, not a payload-size validation.
	FileSize *int64
	// Checksum optionally supplies a nonnegative macro number, not a checksum
	// calculated or verified against the assembled file.
	Checksum *int
}

// Symbol combines an immutable logical symbol with PDF417 metadata.
type Symbol struct {
	logical         barcode.Symbol
	compaction      Compaction
	errorCorrection ErrorCorrection
	compact         bool
	eci             int
	macro           *Macro
}

// Logical returns the immutable format-neutral symbol.
func (symbol Symbol) Logical() barcode.Symbol { return symbol.logical }

// Compaction returns the requested high-level encoding strategy.
func (symbol Symbol) Compaction() Compaction { return symbol.compaction }

// ErrorCorrection returns the selected correction level.
func (symbol Symbol) ErrorCorrection() ErrorCorrection { return symbol.errorCorrection }

// Compact reports whether compact PDF417 layout was selected.
func (symbol Symbol) Compact() bool { return symbol.compact }

// ECI returns the requested ECI assignment, or zero when absent.
func (symbol Symbol) ECI() int { return symbol.eci }

// Macro returns a defensive copy of the Macro PDF417 control block.
func (symbol Symbol) Macro() (Macro, bool) {
	if symbol.macro == nil {
		return Macro{}, false
	}

	return cloneMacro(*symbol.macro), true
}

// Encode returns a PDF417 matrix with explicit compaction, correction, and
// row/column bounds. Each codeword row is rendered four modules high. Nonempty
// payload is capped at 4096 bytes but available codeword capacity can be smaller;
// macro text/FileID has no separate preencoding admission cap. Payload/options
// must not mutate during the synchronous call; output clones payload and macro
// pointer values. Error causes may be wrapped with ErrInvalidInput. Macro assembly
// and physical scan qualification are outside this unadvertised profile (DEC-010).
func Encode(payload []byte, options Options) (Symbol, error) {
	if len(payload) == 0 {
		return Symbol{}, ErrInvalidInput
	}
	if len(payload) > maxPayloadBytes {
		return Symbol{}, ErrInvalidInput
	}
	if options.Compaction > Numeric {
		return Symbol{}, ErrInvalidInput
	}
	if options.QuietZone < 0 {
		return Symbol{}, ErrInvalidInput
	}
	if !validDimension(options.MinRows, options.MaxRows, 90) {
		return Symbol{}, ErrInvalidInput
	}
	if !validDimension(options.MinColumns, options.MaxColumns, 30) {
		return Symbol{}, ErrInvalidInput
	}
	quietZone := options.QuietZone
	if quietZone == 0 {
		quietZone = defaultQuietZone
	}
	if quietZone > 256 {
		return Symbol{}, ErrInvalidInput
	}
	errorCorrection := options.ErrorCorrection
	if errorCorrection == DefaultErrorCorrection {
		errorCorrection = Level2
	}
	encoder := pdf417encoder.NewPDF417Encoder()
	encoder.SetCompaction(pdf417encoder.Compaction(options.Compaction))
	encoder.SetCompact(options.Compact)
	minRows, maxRows := dimensions(options.MinRows, options.MaxRows, 3, 90)
	minColumns, maxColumns := dimensions(options.MinColumns, options.MaxColumns, 1, 30)
	encoder.SetDimensions(maxColumns, minColumns, maxRows, minRows)
	if err := encoder.GenerateBarcodeLogicWithControls(
		string(payload), int(errorCorrection-1), options.ECI, internalMacro(options.Macro),
	); err != nil {
		return Symbol{}, fmt.Errorf("%w: %w", ErrInvalidInput, err)
	}
	source := encoder.BarcodeMatrix().ScaledMatrix(1, 4)
	width, height := len(source[0])+2*quietZone, len(source)+2*quietZone
	modules := make([]bool, width*height)
	for y, row := range source {
		outputY := len(source) - y - 1 + quietZone
		for x, value := range row {
			modules[outputY*width+x+quietZone] = value == 1
		}
	}
	matrix, _ := barcode.NewMatrix(width, height, modules)
	logical, _ := barcode.NewSymbol(barcode.SymbolOptions{
		Format: barcode.PDF417, Payload: payload, Matrix: matrix,
	})

	return Symbol{
		logical: logical, compaction: options.Compaction,
		errorCorrection: errorCorrection, compact: options.Compact, eci: options.ECI,
		macro: cloneMacroPointer(options.Macro),
	}, nil
}

func validDimension(minimum, maximum, limit int) bool {
	switch {
	case minimum < 0:
		return false
	case maximum < 0:
		return false
	case minimum > limit:
		return false
	case maximum > limit:
		return false
	case maximum != 0 && minimum > maximum:
		return false
	default:
		return true
	}
}

func dimensions(minimum, maximum, defaultMinimum, defaultMaximum int) (int, int) {
	if minimum == 0 {
		minimum = defaultMinimum
	}
	if maximum == 0 {
		maximum = defaultMaximum
	}
	return minimum, maximum
}

func internalMacro(macro *Macro) *pdf417encoder.Macro {
	if macro == nil {
		return nil
	}
	return &pdf417encoder.Macro{
		SegmentIndex: macro.SegmentIndex, FileID: macro.FileID, LastSegment: macro.LastSegment,
		FileName: macro.FileName, SegmentCount: macro.SegmentCount, Timestamp: macro.Timestamp,
		Sender: macro.Sender, Addressee: macro.Addressee, FileSize: macro.FileSize,
		Checksum: macro.Checksum,
	}
}

func cloneMacroPointer(macro *Macro) *Macro {
	if macro == nil {
		return nil
	}
	cloned := cloneMacro(*macro)
	return &cloned
}

func cloneMacro(macro Macro) Macro {
	macro.SegmentCount = clonePointer(macro.SegmentCount)
	macro.Timestamp = clonePointer(macro.Timestamp)
	macro.FileSize = clonePointer(macro.FileSize)
	macro.Checksum = clonePointer(macro.Checksum)
	return macro
}

func clonePointer[T any](value *T) *T {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}
