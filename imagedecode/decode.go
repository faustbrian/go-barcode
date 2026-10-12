// Package imagedecode returns the first successful software barcode candidate
// from Decode or bounded-byte DecodeEncoded. It owns neither byte acquisition,
// image callbacks, cameras/devices nor decoded sequence assembly. Inputs are
// borrowed during the call and output uses copied barcode.DecodeResult values.
// Geometry/attempt limits and phase deadlines are not peak-memory or hard-time
// guarantees; only QR/DataMatrix/Aztec adapter calls contain dependency panics.
// See BARCODE-DEC-012 through 014 and docs/security.md for the exact boundaries.
package imagedecode

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  //nolint:revive // Register GIF for the generic image decoder.
	_ "image/jpeg" //nolint:revive // Register JPEG for the generic image decoder.
	_ "image/png"  //nolint:revive // Register PNG for the generic image decoder.
	"math"
	"strings"
	"time"

	"github.com/faustbrian/go-barcode/v2/barcode"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

const (
	defaultMaxDimension    = 8192
	defaultMaxPixels       = 16 * 1024 * 1024
	defaultMaxCandidates   = 64
	defaultMaxPayload      = 4096
	defaultMaxRotations    = 4
	defaultMaxMemoryBytes  = 64 * 1024 * 1024
	defaultMaxEncodedBytes = 16 * 1024 * 1024
)

var (
	// ErrLimitExceeded reports a caller-provided decode budget violation.
	ErrLimitExceeded = errors.New("imagedecode: limit exceeded")
	// ErrUnsupportedFormat reports a format absent from the image decoder.
	ErrUnsupportedFormat = errors.New("imagedecode: unsupported format")
	// ErrInvalidImage reports an image that cannot be converted for decoding.
	ErrInvalidImage = errors.New("imagedecode: invalid image")
	// ErrNotFound reports that no requested barcode was detected.
	ErrNotFound = errors.New("imagedecode: barcode not found")
)

// Limits admit input geometry and candidate attempts; returned payload and
// correction caps apply after reader work. They are not peak-memory or hard-time bounds.
// MaxRotations includes the original orientation and is limited to four.
type Limits struct {
	// MaxWidth is an inclusive positive input width cap; zero selects 8192.
	MaxWidth int
	// MaxHeight is the corresponding height cap; zero selects 8192.
	MaxHeight int
	// MaxPixels bounds width*height before conversion; zero selects 16Mi pixels.
	MaxPixels int
	// MaxCandidates budgets formats*rotations*inversion variants; zero selects 64.
	// Duplicate format entries count too; it is not detected-symbol count.
	MaxCandidates int
	// MaxPayloadBytes caps returned text bytes and raw bytes separately;
	// zero selects 4096. It is checked after the selected reader has decoded/allocated.
	MaxPayloadBytes int
	// MaxRotations includes the original orientation; zero selects 4, maximum 4.
	MaxRotations int
	// MaxMemoryBytes admits geometry at a 4-byte-per-pixel estimate;
	// zero selects 64MiB. It does not measure or cap all codec/reader/rotation allocations.
	MaxMemoryBytes int
	// MaxEncodedBytes caps DecodeEncoded input; zero selects 16MiB.
	// Decode ignores it because its image has already been acquired.
	MaxEncodedBytes int
	// MaxCorrections is zero for no extra correction check. Positive values
	// reject a reported correction count above the cap after decoding, not parity work.
	MaxCorrections int
	// MaxDuration is zero for no additional deadline; positive derives a timeout.
	// Checks occur between phases/attempts and do not preempt codecs/image callbacks
	// or an in-progress dependency reader; success is not universally postchecked.
	MaxDuration time.Duration
}

// DecodeEncoded bounds and decodes PNG, JPEG, or GIF bytes before scanning for
// a barcode. The caller retains ownership of encoded and must not mutate it
// during this call. Encoded length and header geometry are admitted before the
// full codec decode, but this is not a peak-memory bound. Registered image codecs
// and config/decode calls do not accept ctx; checks only surround those phases.
// The package registers GIF/JPEG/PNG and generic image registration can add formats.
// It owns no stream/file acquisition and returns categorical codec errors.
func DecodeEncoded(ctx context.Context, encoded []byte, options Options) (barcode.DecodeResult, error) {
	if err := ctx.Err(); err != nil {
		return barcode.DecodeResult{}, err
	}
	if encoded == nil {
		return barcode.DecodeResult{}, ErrInvalidImage
	}
	limits, err := normalizeLimits(options.Limits)
	if err != nil {
		return barcode.DecodeResult{}, err
	}
	if limits.MaxDuration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, limits.MaxDuration)
		defer cancel()
	}
	if len(encoded) > limits.MaxEncodedBytes {
		return barcode.DecodeResult{}, ErrLimitExceeded
	}
	if err := ctx.Err(); err != nil {
		return barcode.DecodeResult{}, err
	}
	configuration, _, err := image.DecodeConfig(bytes.NewReader(encoded))
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return barcode.DecodeResult{}, ctxErr
		}
		return barcode.DecodeResult{}, ErrInvalidImage
	}
	if err := validateImageSize(configuration.Width, configuration.Height, limits); err != nil {
		return barcode.DecodeResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return barcode.DecodeResult{}, err
	}
	decoded, _, err := image.Decode(bytes.NewReader(encoded))
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return barcode.DecodeResult{}, ctxErr
		}
		return barcode.DecodeResult{}, ErrInvalidImage
	}
	if err := ctx.Err(); err != nil {
		return barcode.DecodeResult{}, err
	}
	options.Limits = limits

	return Decode(ctx, decoded, options)
}

// Options controls requested formats, scan behavior, and resource limits.
type Options struct {
	// Formats is borrowed during the call; empty selects the fixed default list,
	// including implemented-but-unadvertised DataMatrix/PDF417, but not ITF14.
	// Requested semantic labels select readers; no general GS1/ITF14 revalidation occurs.
	Formats []barcode.Format
	// AllowInverted adds an inverted candidate per rotation and doubles admission count.
	AllowInverted bool
	// TryHarder forwards a reader hint, not a readability/latency guarantee.
	TryHarder bool
	// AssumeCode39Checksum asks the reader to validate/remove Code 39 checks;
	// returned checksum status still follows the fixed format classifier.
	AssumeCode39Checksum bool
	// Limits supplies positive budgets/defaults; negative values fail.
	Limits Limits
}

type candidate struct {
	format barcode.Format
	reader gozxing.Reader
	gs1    bool
}

var defaultFormats = [...]barcode.Format{
	barcode.QRCode, barcode.Code128, barcode.GS1128, barcode.Code39,
	barcode.Code93, barcode.EAN8, barcode.EAN13, barcode.UPCA,
	barcode.UPCE, barcode.ITF, barcode.Codabar, barcode.DataMatrix,
	barcode.Aztec, barcode.PDF417,
}

// Decode returns the first requested format found within the configured
// bounds. The search is rotation-major, then configured formats, then ordinary/
// inverted variants; it does not enumerate every detected symbol. Input image
// and Formats must remain safe to read during the call. Codecs, image callbacks,
// linear/PDF417 readers are not panic-contained or preempted by context checks.
// Limits admit estimated geometry/attempts and postdecode payload/correction data,
// not all allocations. Metadata reflects the requested format/reader profile,
// not an independent application checksum or GS1 validator. Result diagnostics
// can contain decoded supplement/macro metadata: treat them as content. No content
// is executed/followed. Errors classify invalid image, unsupported format, limits,
// cancellation/deadline or exhaustion; there is no partial result on these errors.
func Decode(ctx context.Context, input image.Image, options Options) (barcode.DecodeResult, error) {
	started := time.Now()
	if err := ctx.Err(); err != nil {
		return barcode.DecodeResult{}, err
	}
	if input == nil {
		return barcode.DecodeResult{}, ErrInvalidImage
	}
	limits, err := normalizeLimits(options.Limits)
	if err != nil {
		return barcode.DecodeResult{}, err
	}
	if limits.MaxDuration > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, limits.MaxDuration)
		defer cancel()
		if err := contextError(ctx, started, limits.MaxDuration); err != nil {
			return barcode.DecodeResult{}, err
		}
	}
	bounds := input.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if err := validateImageSize(width, height, limits); err != nil {
		return barcode.DecodeResult{}, err
	}
	inversionAttempts := 1
	if options.AllowInverted {
		inversionAttempts = 2
	}
	candidateCount := len(options.Formats)
	if candidateCount == 0 {
		candidateCount = len(defaultFormats)
	}
	attemptsPerCandidate := limits.MaxRotations * inversionAttempts
	if candidateCount > limits.MaxCandidates/attemptsPerCandidate {
		return barcode.DecodeResult{}, ErrLimitExceeded
	}
	candidates, err := candidatesFor(options.Formats, options.AssumeCode39Checksum)
	if err != nil {
		return barcode.DecodeResult{}, err
	}

	source := gozxing.NewLuminanceSourceFromImage(input)
	bitmap, _ := gozxing.NewBinaryBitmap(gozxing.NewHybridBinarizer(source))
	var inverted *gozxing.BinaryBitmap
	if options.AllowInverted {
		inverted, _ = gozxing.NewBinaryBitmap(gozxing.NewHybridBinarizer(source.Invert()))
	}
	hints := map[gozxing.DecodeHintType]interface{}{}
	if options.AllowInverted {
		hints[gozxing.DecodeHintType_ALSO_INVERTED] = true
	}
	if options.TryHarder {
		hints[gozxing.DecodeHintType_TRY_HARDER] = true
	}

	current := bitmap
	currentInverted := inverted
	for rotation := range limits.MaxRotations {
		if rotation > 0 {
			current, _ = current.RotateCounterClockwise()
			if currentInverted != nil {
				currentInverted, _ = currentInverted.RotateCounterClockwise()
			}
		}
		for _, item := range candidates {
			variants := []*gozxing.BinaryBitmap{current}
			if currentInverted != nil {
				variants = append(variants, currentInverted)
			}
			for _, variant := range variants {
				if err := contextError(ctx, started, limits.MaxDuration); err != nil {
					return barcode.DecodeResult{}, err
				}
				attemptHints := cloneHints(hints)
				if item.gs1 {
					attemptHints[gozxing.DecodeHintType_ASSUME_GS1] = true
				}
				result, decodeErr := item.reader.Decode(variant, attemptHints)
				item.reader.Reset()
				if decodeErr != nil {
					continue
				}
				if exceedsCorrectionBudget(correctionsFor(result), limits.MaxCorrections) {
					return barcode.DecodeResult{}, ErrLimitExceeded
				}
				if exceedsPayloadBudget(result, limits.MaxPayloadBytes) {
					return barcode.DecodeResult{}, ErrLimitExceeded
				}
				payload := result.GetText()
				if item.gs1 {
					if strings.HasPrefix(payload, "]C1") {
						payload = strings.TrimPrefix(payload, "]C1")
						result.PutMetadata(gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER, "]C1")
					}
				}
				// Preserve the requested semantic format. Some symbologies share a
				// physical reader, notably ITF and GS1 ITF-14, so the reader's broad
				// format cannot always represent the caller's requested contract.
				format := item.format
				orientation := orientationFor(result, rotation)
				decoded, _ := barcode.NewDecodeResult(barcode.DecodeResultOptions{
					Format:      format,
					Payload:     []byte(payload),
					RawBytes:    result.GetRawBytes(),
					Orientation: orientation,
					Checksum:    checksumStatus(format),
					Diagnostics: metadataDiagnostics(result),
				})
				return decoded, nil
			}
		}
	}
	return barcode.DecodeResult{}, decodeFailure(ctx, started, limits.MaxDuration)
}

func validateImageSize(width, height int, limits Limits) error {
	if width <= 0 {
		return ErrLimitExceeded
	}
	if height <= 0 {
		return ErrLimitExceeded
	}
	if width > limits.MaxWidth {
		return ErrLimitExceeded
	}
	if height > limits.MaxHeight {
		return ErrLimitExceeded
	}
	if width > limits.MaxPixels/height {
		return ErrLimitExceeded
	}
	const bytesPerPixel = 4
	if width > (limits.MaxMemoryBytes/bytesPerPixel)/height {
		return ErrLimitExceeded
	}

	return nil
}

func decodeFailure(ctx context.Context, started time.Time, maximum time.Duration) error {
	if err := contextError(ctx, started, maximum); err != nil {
		return err
	}

	return ErrNotFound
}

func contextError(ctx context.Context, started time.Time, maximum time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if maximum <= 0 {
		return nil
	}
	if time.Now().Before(started.Add(maximum)) {
		return nil
	}
	return context.DeadlineExceeded
}

func normalizeLimits(limits Limits) (Limits, error) {
	if limits.MaxWidth < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxHeight < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxPixels < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxCandidates < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxPayloadBytes < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxRotations < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxMemoryBytes < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxEncodedBytes < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxCorrections < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxDuration < 0 {
		return Limits{}, ErrLimitExceeded
	}
	if limits.MaxWidth == 0 {
		limits.MaxWidth = defaultMaxDimension
	}
	if limits.MaxHeight == 0 {
		limits.MaxHeight = defaultMaxDimension
	}
	if limits.MaxPixels == 0 {
		limits.MaxPixels = defaultMaxPixels
	}
	if limits.MaxCandidates == 0 {
		limits.MaxCandidates = defaultMaxCandidates
	}
	if limits.MaxPayloadBytes == 0 {
		limits.MaxPayloadBytes = defaultMaxPayload
	}
	if limits.MaxRotations == 0 {
		limits.MaxRotations = defaultMaxRotations
	}
	if limits.MaxMemoryBytes == 0 {
		limits.MaxMemoryBytes = defaultMaxMemoryBytes
	}
	if limits.MaxEncodedBytes == 0 {
		limits.MaxEncodedBytes = defaultMaxEncodedBytes
	}
	if limits.MaxRotations > 4 {
		return Limits{}, ErrLimitExceeded
	}

	return limits, nil
}

func candidatesFor(formats []barcode.Format, assumeCode39Checksum bool) ([]candidate, error) {
	if len(formats) == 0 {
		formats = defaultFormats[:]
	}
	result := make([]candidate, 0, len(formats))
	for _, format := range formats {
		var item candidate
		switch format {
		case barcode.QRCode:
			item = candidate{format: format, reader: newQRCodeReader()}
		case barcode.DataMatrix:
			item = candidate{format: format, reader: newDataMatrixReader()}
		case barcode.Aztec:
			item = candidate{format: format, reader: newAztecReader()}
		case barcode.PDF417:
			item = candidate{format: format, reader: pdf417Reader{}}
		case barcode.Code128:
			item = candidate{format: format, reader: oned.NewCode128Reader()}
		case barcode.GS1128:
			item = candidate{format: format, reader: oned.NewCode128Reader(), gs1: true}
		case barcode.Code39:
			item = candidate{format: format, reader: oned.NewCode39ReaderWithFlags(assumeCode39Checksum, true)}
		case barcode.Code93:
			item = candidate{format: format, reader: oned.NewCode93Reader()}
		case barcode.EAN8:
			item = candidate{format: format, reader: oned.NewEAN8Reader()}
		case barcode.EAN13:
			item = candidate{format: format, reader: oned.NewEAN13Reader()}
		case barcode.UPCA:
			item = candidate{format: format, reader: oned.NewUPCAReader()}
		case barcode.UPCE:
			item = candidate{format: format, reader: oned.NewUPCEReader()}
		case barcode.ITF, barcode.ITF14:
			item = candidate{format: format, reader: oned.NewITFReader()}
		case barcode.Codabar:
			item = candidate{format: format, reader: oned.NewCodaBarReader()}
		default:
			return nil, fmt.Errorf("%w: %q", ErrUnsupportedFormat, format)
		}
		result = append(result, item)
	}

	return result, nil
}

func cloneHints(source map[gozxing.DecodeHintType]interface{}) map[gozxing.DecodeHintType]interface{} {
	result := make(map[gozxing.DecodeHintType]interface{}, len(source))
	for key, value := range source {
		result[key] = value
	}

	return result
}

func orientationFor(result *gozxing.Result, rotation int) barcode.Orientation {
	degrees := (360 - rotation*90) % 360
	if value, ok := result.GetResultMetadata()[gozxing.ResultMetadataType_ORIENTATION].(int); ok {
		degrees = (degrees + value) % 360
	} else if result.GetBarcodeFormat() == gozxing.BarcodeFormat_QR_CODE {
		degrees = (degrees + qrOrientation(result.GetResultPoints())) % 360
	}
	switch degrees {
	case 90:
		return barcode.Orientation90
	case 180:
		return barcode.Orientation180
	case 270:
		return barcode.Orientation270
	default:
		return barcode.Orientation0
	}
}

func qrOrientation(points []gozxing.ResultPoint) int {
	if len(points) < 2 {
		return 0
	}
	horizontal := points[1].GetX() - points[0].GetX()
	vertical := points[1].GetY() - points[0].GetY()
	dx := math.Abs(horizontal)
	dy := math.Abs(vertical)
	if dx > dy {
		if math.Signbit(horizontal) {
			return 90
		}
		return 270
	}
	if vertical > 0 {
		return 180
	}

	return 0
}

// This is a fixed format/reader classification, not recomputation over decoded
// payload: ITF14 and GS1 labels alone do not establish their application checks.
func checksumStatus(format barcode.Format) barcode.ChecksumStatus {
	switch format {
	case barcode.QRCode, barcode.Codabar:
		return barcode.ChecksumNotApplicable
	case barcode.Code39:
		return barcode.ChecksumUnknown
	case barcode.Code128, barcode.GS1128, barcode.Code93, barcode.EAN8,
		barcode.EAN13, barcode.UPCA, barcode.UPCE, barcode.ITF, barcode.ITF14,
		barcode.DataMatrix, barcode.PDF417, barcode.Aztec:
		return barcode.ChecksumValid
	default:
		return barcode.ChecksumValid
	}
}

func metadataDiagnostics(result *gozxing.Result) []string {
	metadata := result.GetResultMetadata()
	var diagnostics []string
	for _, key := range []gozxing.ResultMetadataType{
		gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL,
		gozxing.ResultMetadataType_UPC_EAN_EXTENSION,
		gozxing.ResultMetadataType_SYMBOLOGY_IDENTIFIER,
		gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE,
		gozxing.ResultMetadataType_STRUCTURED_APPEND_PARITY,
		gozxing.ResultMetadataType_PDF417_EXTRA_METADATA,
	} {
		if value, ok := metadata[key]; ok {
			diagnostics = append(diagnostics, key.String()+"="+fmt.Sprint(value))
		}
	}
	if corrections := correctionsFor(result); corrections > 0 {
		diagnostics = append(diagnostics, fmt.Sprintf("ERRORS_CORRECTED=%d", corrections))
	}
	if eci, ok := metadata[resultMetadataECI].(eciMetadata); ok {
		diagnostics = append(diagnostics, fmt.Sprintf("ECI_ASSIGNMENT=%d", eci.assignment))
		if !eci.supported {
			diagnostics = append(diagnostics, fmt.Sprintf("ECI_UNSUPPORTED=%d", eci.assignment))
		}
	}

	return diagnostics
}

func correctionsFor(result *gozxing.Result) int {
	corrections, _ := result.GetResultMetadata()[gozxing.ResultMetadataType_OTHER].(int)

	return corrections
}

func exceedsCorrectionBudget(corrections, maximum int) bool {
	if maximum <= 0 {
		return false
	}
	return corrections > maximum
}

func exceedsPayloadBudget(result *gozxing.Result, maximum int) bool {
	if len(result.GetText()) > maximum {
		return true
	}
	return len(result.GetRawBytes()) > maximum
}
