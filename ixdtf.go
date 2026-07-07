// Package ixdtf implements RFC 9557 Internet Extended Date/Time Format (IXDTF).
// IXDTF extends RFC 3339 by adding optional suffix elements for timezone names.
// and additional metadata while maintaining full backward compatibility.
//
// See RFC 9557: https://datatracker.ietf.org/doc/rfc9557/
package ixdtf

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/8beeeaaat/ixdtf/abnf"
)

// Layout represents time layout strings.
type Layout string

const (
	// LayoutRFC3339 represents the RFC 3339 layout.
	LayoutRFC3339 Layout = time.RFC3339

	// LayoutRFC3339Nano represents the RFC 3339 layout with nanoseconds.
	LayoutRFC3339Nano Layout = time.RFC3339Nano

	// LayoutRFC3339Extended represents the RFC 9557(IXDTF) layout with optional timezone and extensions.
	// https://www.rfc-editor.org/rfc/rfc9557.html#section-3
	LayoutRFC3339Extended Layout = time.RFC3339 + "*([time-zone-name][tags])"

	// LayoutRFC3339NanoExtended represents the RFC 9557(IXDTF) layout with nanoseconds and optional timezone and extensions.
	// https://www.rfc-editor.org/rfc/rfc9557.html#section-3
	LayoutRFC3339NanoExtended Layout = time.RFC3339Nano + "*([time-zone-name][tags])"

	// ExtensionUnicodeCalendar is tag key for Unicode calendar extension.
	// https://www.rfc-editor.org/rfc/rfc9557.html#section-5
	ExtensionUnicodeCalendar = "u-ca"
)

// Common parsing errors.
var (
	ErrCriticalExtension            = errors.New("critical extension cannot be processed")
	ErrExperimentalExtension        = abnf.ErrExperimentalExtension
	ErrInvalidExtension             = errors.New("invalid extension format")
	ErrInvalidSuffix                = errors.New("invalid IXDTF suffix format")
	ErrInvalidTagCalendarIdentifier = errors.New("invalid calendar tag identifier")
	ErrInvalidTimezone              = errors.New("invalid timezone name")
	ErrPrivateExtension             = abnf.ErrPrivateExtension
	ErrTimezoneOffsetMismatch       = errors.New("timezone offset does not match the specified timezone")
)

// IXDTFExtensions holds IXDTF suffix information that extends RFC 3339.
//
//nolint:revive // Keeping existing public API name for compatibility
type IXDTFExtensions struct {
	Location *time.Location

	// CriticalLocation reports whether the time-zone annotation carried a
	// critical "!" flag (e.g. "[!Europe/London]"). Per RFC 9557 Section 3.3 a
	// critical annotation MUST be processable, and per Section 3.4 an
	// application MUST act on any inconsistency it introduces.
	CriticalLocation bool

	// Tags contains extension tags as key-value pairs.
	// Example: map[ExtensionUnicodeCalendar]"japanese".
	Tags map[string]string

	// Critical indicates which tags are marked as critical (must be processed).
	// Critical tags are marked with "!" prefix in the IXDTF string.
	Critical map[string]bool
}

// ParseError represents an error that occurred during IXDTF parsing.
type ParseError struct {
	Err    error
	Layout Layout
	Value  string
}

func (e *ParseError) Error() string {
	if e.Err == nil {
		return ""
	}
	return "IXDTFE parsing time \"" + e.Value + "\" as \"" + string(e.Layout) + "\": " + e.Err.Error()
}

// Unwrap returns the underlying error, so callers can match sentinel errors
// with errors.Is (e.g. errors.Is(err, ErrTimezoneOffsetMismatch)).
func (e *ParseError) Unwrap() error {
	return e.Err
}

// TimezoneConsistencyResult holds information about timezone offset consistency.
type TimezoneConsistencyResult struct {
	// Location is the loaded timezone location.
	Location *time.Location
	// IsConsistent indicates whether the offset matches the timezone.
	IsConsistent bool
	// OriginalOffset is the offset from the timestamp.
	OriginalOffset int
	// ExpectedOffset is the expected offset for the timezone.
	ExpectedOffset int
	// Timezone is the timezone identifier that was checked.
	Timezone string
	// Skipped indicates if the consistency check was skipped.
	//
	// Deprecated: no checks are skipped anymore; the field is retained for
	// API compatibility and is always false.
	Skipped bool
}

// Format formats a time with IXDTF extensions using RFC 3339 format.
// The time-zone annotation is emitted with a leading "!" when
// ext.CriticalLocation is set.
func Format(t time.Time, ext *IXDTFExtensions) (string, error) {
	// Validate extensions (including critical tag processing).
	if err := validateExtensions(ext); err != nil {
		return "", err
	}
	if err := validateCriticalLocation(t, ext); err != nil {
		return "", err
	}
	b := appendSuffix(t, ext, time.RFC3339)

	return string(b), nil
}

// FormatNano formats a time with IXDTF extensions using RFC 3339 format with nanoseconds.
// The time-zone annotation is emitted with a leading "!" when
// ext.CriticalLocation is set.
func FormatNano(t time.Time, ext *IXDTFExtensions) (string, error) {
	// Validate extensions (including critical tag processing).
	if err := validateExtensions(ext); err != nil {
		return "", err
	}
	if err := validateCriticalLocation(t, ext); err != nil {
		return "", err
	}
	b := appendSuffix(t, ext, time.RFC3339Nano)

	return string(b), nil
}

// NewIXDTFExtensionsArgs contains the arguments for creating IXDTFExtensions.
type NewIXDTFExtensionsArgs struct {
	Location *time.Location
	// CriticalLocation marks the time-zone annotation as critical ("!");
	// see IXDTFExtensions.CriticalLocation.
	CriticalLocation bool
	Tags             map[string]string
	Critical         map[string]bool
}

// NewIXDTFExtensions creates a new IXDTFExtensions with initialized maps.
func NewIXDTFExtensions(args *NewIXDTFExtensionsArgs) *IXDTFExtensions {
	ext := &IXDTFExtensions{
		Location:         nil,
		CriticalLocation: false,
		Tags:             make(map[string]string),
		Critical:         make(map[string]bool),
	}

	if args != nil {
		if args.Location != nil {
			ext.Location = args.Location
		}
		ext.CriticalLocation = args.CriticalLocation
		if args.Tags != nil {
			ext.Tags = args.Tags
		}
		if args.Critical != nil {
			ext.Critical = args.Critical
		}
	}

	return ext
}

// Parse parses an IXDTF string and returns the time and extension information.
func Parse(s string, strict bool) (time.Time, *IXDTFExtensions, error) {
	rfc3339End := findRFC3339End(s)
	rfc3339Portion := s[:rfc3339End]

	t, err := parseRFC3339Portion(rfc3339Portion)
	if err != nil {
		return time.Time{}, nil, newParseError(LayoutRFC3339, s, err)
	}

	ext := NewIXDTFExtensions(nil)
	if rfc3339End < len(s) {
		suffixPortion := s[rfc3339End:]
		if ext, err = parseSuffix(suffixPortion, strict); err != nil {
			return time.Time{}, nil, newParseError(LayoutRFC3339Extended, s, err)
		}
	}

	if err = validateExtensionsStrict(ext, strict); err != nil {
		return time.Time{}, nil, newParseError(LayoutRFC3339Extended, s, err)
	}

	// Check timezone consistency if timezone is provided
	if ext.Location != nil {
		offsetUnknown := hasUnknownLocalOffset(rfc3339Portion)
		// A critical time zone must be acted upon, so an inconsistency is an
		// error even in non-strict mode (RFC 9557 Section 3.4).
		result, checkErr := checkTimezoneConsistency(t, ext.Location, strict || ext.CriticalLocation, offsetUnknown)
		if checkErr != nil {
			return time.Time{}, nil, newParseError(
				LayoutRFC3339NanoExtended,
				s,
				checkErr,
			)
		}

		// Per RFC 9557: In non-strict mode with inconsistent timezone,
		// preserve the original timestamp and only apply timezone if consistent
		if result.IsConsistent && result.Location != nil {
			t = t.In(result.Location)
		}
		// In non-strict mode with inconsistency, keep original timestamp as-is
	}

	return t, ext, nil
}

// Validate validates an IXDTF string for format correctness without parsing the time component.
func Validate(s string, strict bool) error {
	rfc3339End := findRFC3339End(s)
	rfc3339Portion := s[:rfc3339End]

	if rfc3339Portion == "" {
		return newParseError(LayoutRFC3339, s, errors.New("empty datetime string"))
	}

	// Parse the RFC3339 portion to validate format and get the timestamp
	t, err := parseRFC3339Portion(rfc3339Portion)
	if err != nil {
		return newParseError(LayoutRFC3339, s, errors.New("invalid portion: "+err.Error()))
	}

	var ext *IXDTFExtensions
	if rfc3339End < len(s) {
		suffixPortion := s[rfc3339End:]
		if ext, err = parseSuffix(suffixPortion, strict); err != nil {
			return newParseError(LayoutRFC3339Extended, s, err)
		}
	}

	if err = validateExtensionsStrict(ext, strict); err != nil {
		return newParseError(LayoutRFC3339Extended, s, err)
	}

	// Check timezone consistency if timezone is provided
	if ext != nil && ext.Location != nil {
		offsetUnknown := hasUnknownLocalOffset(rfc3339Portion)
		// A critical time zone must be acted upon, so an inconsistency is an
		// error even in non-strict mode (RFC 9557 Section 3.4).
		_, tzErr := checkTimezoneConsistency(t, ext.Location, strict || ext.CriticalLocation, offsetUnknown)
		if tzErr != nil {
			return newParseError(LayoutRFC3339NanoExtended, s, tzErr)
		}
		// Non-critical inconsistencies are not validation errors in
		// non-strict mode; strict mode or a critical zone turns them into
		// errors above.
	}

	// Optional: check if the complete string matches the ABNF pattern.
	// This serves as an additional validation layer, but only for well-formed strings.
	// Skip for empty strings and malformed inputs to avoid double error reporting.
	if rfc3339Portion != "" && (rfc3339End >= len(s) || s[rfc3339End:] != "") {
		if abnfErr := abnf.AbnfDateTimeExt.ValidateDateTimeExt(s); abnfErr != nil {
			// Only report ABNF mismatch if basic validation passes.
			// This prevents confusing error messages for clearly invalid input.
			return newParseError(LayoutRFC3339Extended, s, abnfErr)
		}
	}

	return nil
}

// formatLocation returns the location whose name is emitted as the time-zone
// annotation: ext.Location when set, otherwise the timestamp's own named zone.
// When falling back to the timestamp's zone, UTC, Local, and unnamed zones
// produce no annotation, so nil is returned.
func formatLocation(t time.Time, ext *IXDTFExtensions) *time.Location {
	loc := ext.Location
	if loc == nil {
		loc = t.Location()
		if loc == time.UTC || loc.String() == "Local" {
			return nil
		}
	}
	if loc.String() == "" {
		return nil
	}
	return loc
}

// validateCriticalLocation rejects a critical time-zone flag that has no zone
// to attach to — neither ext.Location nor the timestamp's own named zone.
// Emitting output that silently drops the "!" would misrepresent the caller's
// declared critical intent (RFC 9557 Section 3.3).
func validateCriticalLocation(t time.Time, ext *IXDTFExtensions) error {
	if ext != nil && ext.CriticalLocation && formatLocation(t, ext) == nil {
		return ErrCriticalExtension
	}
	return nil
}

func appendSuffix(t time.Time, ext *IXDTFExtensions, format string) []byte {
	if ext == nil {
		ext = NewIXDTFExtensions(nil)
	}
	b := []byte(t.Format(format))

	// Add timezone if we have a valid location to display
	if loc := formatLocation(t, ext); loc != nil {
		b = append(b, '[')
		if ext.CriticalLocation {
			b = append(b, '!')
		}
		b = append(b, loc.String()...)
		b = append(b, ']')
	}

	// set tags
	keys := make([]string, 0, len(ext.Tags))
	for key := range ext.Tags {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	// Append tags in sorted order for consistency
	for _, key := range keys {
		if err := abnf.AbnfSuffixKey.ValidateSuffixKey(key); err != nil {
			continue
		}
		value := ext.Tags[key]
		b = append(b, '[')
		if ext.Critical[key] {
			b = append(b, '!')
		}
		b = append(b, key...)
		b = append(b, '=')
		b = append(b, value...)
		b = append(b, ']')
	}

	return b
}

// checkTimezoneConsistency checks if the timezone offset matches the IANA timezone.
// If strict is true, returns an error when offsets don't match.
// Returns consistency information and an error for timezone loading failures or strict mode mismatches.
func checkTimezoneConsistency(
	timestamp time.Time,
	location *time.Location,
	strict bool,
	offsetUnknown bool,
) (*TimezoneConsistencyResult, error) {
	result := &TimezoneConsistencyResult{
		Location: location,
	}

	if location == nil {
		result.IsConsistent = true // No timezone means no inconsistency
		return result, nil
	}
	loc := ensureRealLocation(location)
	switch {
	case loc != nil:
	case isOffsetLocationName(location.String()):
		// A numeric-offset annotation (e.g. "[+09:00]") produced this
		// FixedZone; its offset is authoritative, so use it directly.
		loc = location
	default:
		loaded, err := loadLocationCached(location.String())
		if err != nil {
			// In non-strict mode, ignore unknown timezone errors per RFC 9557
			if !strict {
				result.IsConsistent = true // Can't check consistency for unknown timezone
				return result, nil
			}
			return result, err
		}
		loc = loaded
	}
	result.Location = loc

	// Per RFC 3339 Section 4.3 (as updated by RFC 9557 Section 2.2) and
	// RFC 9557 Section 3.4 (Figure 2): when the RFC 3339 part carries an
	// unknown local offset ("Z" or "-00:00"), the instant in UTC is known but
	// no specific local offset is asserted. Pairing it with a time-zone
	// annotation is therefore never an inconsistency; resolve local time using
	// the annotation's rules rather than comparing offsets.
	// See https://www.rfc-editor.org/rfc/rfc9557#section-3.4 (Figure 2).
	if offsetUnknown {
		_, result.OriginalOffset = timestamp.Zone()
		_, result.ExpectedOffset = timestamp.In(loc).Zone()
		result.IsConsistent = true
		return result, nil
	}

	// Get the timezone offset for the given timestamp.
	expectedTimestamp := timestamp.In(loc)
	_, originalOffset := timestamp.Zone()
	_, expectedOffset := expectedTimestamp.Zone()

	result.OriginalOffset = originalOffset
	result.ExpectedOffset = expectedOffset

	// Check if offsets match (allowing for some flexibility with DST transitions).
	// Etc/GMT zones need no special casing: their POSIX-inverted sign only
	// affects the name, and Go resolves the actual offset correctly.
	result.IsConsistent = (originalOffset == expectedOffset)

	// In strict mode, return an error for inconsistencies
	if strict && !result.IsConsistent {
		return nil, ErrTimezoneOffsetMismatch
	}

	// In non-strict mode (callers escalate critical zones to strict per
	// RFC 9557 Section 3.4), inconsistencies are recorded in the result but
	// do not cause failures.
	return result, nil
}

func findRFC3339End(s string) int {
	for i, c := range s {
		if c == '[' {
			return i
		}
	}
	return len(s)
}

// hasUnknownLocalOffset reports whether the RFC 3339 portion uses the
// "unknown local offset" designator defined in RFC 3339 Section 4.3 and
// updated by RFC 9557 Section 2.2: a "Z" or a negative-zero offset "-00:00".
//
// In that case the instant in UTC is known but no specific local offset is
// asserted, so pairing it with a time-zone annotation is never an
// inconsistency (RFC 9557 Section 3.4, Figure 2); the annotation's rules are
// applied to resolve local time. This differs from a concrete "+00:00", which
// asserts a zero offset and can therefore conflict with the annotation.
//
// References:
//   - RFC 3339 Section 4.3: https://www.rfc-editor.org/rfc/rfc3339#section-4.3
//   - RFC 9557 Section 2.2: https://www.rfc-editor.org/rfc/rfc9557#section-2.2
//   - RFC 9557 Section 3.4: https://www.rfc-editor.org/rfc/rfc9557#section-3.4
func hasUnknownLocalOffset(rfc3339Portion string) bool {
	if n := len(rfc3339Portion); n > 0 {
		if c := rfc3339Portion[n-1]; c == 'Z' || c == 'z' {
			return true
		}
	}
	return strings.HasSuffix(rfc3339Portion, "-00:00")
}

func isValidSuffixKeyRange(s string, start, end int) error {
	if start >= end {
		return ErrInvalidExtension
	}
	return abnf.AbnfSuffixKey.ValidateSuffixKey(s[start:end])
}

func isValidSuffixValue(value string) error {
	if value == "" {
		return nil
	}
	// Use ABNF pattern for basic validation, then check additional constraints
	if err := abnf.AbnfSuffixValues.ValidateSuffixValues(value); err != nil {
		return err
	}
	// Additional validation: no leading/trailing hyphens, no consecutive hyphens
	if value[0] == '-' || value[len(value)-1] == '-' {
		return ErrInvalidExtension
	}
	for i := 1; i < len(value); i++ {
		if value[i] == '-' && value[i-1] == '-' {
			return ErrInvalidExtension
		}
	}
	return nil
}

func isValidSuffixValueRange(s string, start, end int) error {
	if start >= end {
		return ErrInvalidExtension
	}
	return isValidSuffixValue(s[start:end])
}

// validateCriticalExtension enforces critical extension processing rules.
func validateCriticalExtension(key, value string) error {
	if value == "" { // empty value not allowed for critical
		return ErrCriticalExtension
	}
	if key == ExtensionUnicodeCalendar && !isUnicodeCalendarIdentifier(value) {
		return ErrInvalidTagCalendarIdentifier
	}
	return nil
}

func isUnicodeCalendarIdentifier(value string) bool {
	if value == "" {
		return false
	}
	switch strings.ToLower(value) {
	case "buddhist",
		"chinese",
		"coptic",
		"dangi",
		"ethioaa",
		"ethiopic-amete-alem",
		"ethiopic",
		"gregory",
		"gregorian",
		"hebrew",
		"indian",
		"islamic",
		"islamic-umalqura",
		"islamic-tbla",
		"islamic-civil",
		"islamic-rgsa",
		"iso8601",
		"japanese",
		"persian",
		"roc",
		"islamicc":
		return true
	default:
		return false
	}
}

func newParseError(layout Layout, value string, err error) error {
	return &ParseError{
		Layout: layout,
		Value:  value,
		Err:    err,
	}
}

func parseRFC3339Portion(rfc3339Portion string) (time.Time, error) {
	layouts := []string{time.RFC3339Nano, time.RFC3339}
	var lastErr error

	for _, layout := range layouts {
		t, err := time.Parse(layout, rfc3339Portion)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}

// suffixParseState tracks which element kinds have been seen while parsing a
// suffix, enforcing the RFC 9557 Section 4.1 grammar
// "suffix = [time-zone] *suffix-tag": at most one time-zone annotation, and
// it must precede all suffix tags. The state is deliberately independent of
// ext.Location, which stays nil when a non-strict parse ignores an unknown
// zone.
type suffixParseState struct {
	seenTimezone bool
	seenTag      bool
}

func parseSuffix(s string, strict bool) (*IXDTFExtensions, error) {
	ext := NewIXDTFExtensions(nil)
	state := &suffixParseState{}

	i := 0
	for i < len(s) {
		if s[i] != '[' {
			return ext, ErrInvalidSuffix
		}

		// Find the matching ']'
		j := i + 1
		for j < len(s) && s[j] != ']' {
			j++
		}
		if j >= len(s) {
			return ext, ErrInvalidSuffix
		}

		// Parse the content between '[' and ']'
		content := s[i+1 : j]
		if err := parseSuffixElement(content, ext, strict, state); err != nil {
			return ext, err
		}

		i = j + 1
	}

	return ext, nil
}

func parseSuffixElement(content string, ext *IXDTFExtensions, strict bool, state *suffixParseState) error {
	if content == "" {
		return ErrInvalidSuffix
	}

	critical := false
	startIdx := 0
	if content[0] == '!' {
		if len(content) == 1 {
			return ErrInvalidSuffix
		}
		critical = true
		startIdx = 1
	}

	// Extension tag (has '=') vs timezone name.
	if strings.IndexByte(content[startIdx:], '=') >= 0 {
		state.seenTag = true
		return handleExtensionTag(content, critical, startIdx, ext, strict)
	}

	// Timezone name handling.
	//
	// RFC 9557 Section 4.1 permits a critical flag ("!") on a time-zone
	// annotation, e.g. "[!Europe/London]" (Figures 1 and 2 in Section 3.4).
	// A critical annotation MUST be processable (Section 3.3), so an unknown
	// or invalid name is rejected even in non-strict mode.
	//
	// The Section 4.1 grammar ("suffix = [time-zone] *suffix-tag") allows at
	// most one time-zone annotation, placed before any suffix tags. A second
	// one would overwrite the zone and its critical flag, hiding a mandatory
	// Section 3.4 inconsistency error — even when the first zone was unknown
	// and ignored by a non-strict parse.
	if state.seenTimezone || state.seenTag {
		return ErrInvalidSuffix
	}
	state.seenTimezone = true
	tzContent := content[startIdx:]
	if tzContent == "" {
		return nil
	}
	tzStrict := strict || critical

	if loc, ok := tryLoadTimezone(tzContent); ok {
		ext.Location = loc
		ext.CriticalLocation = critical
		return nil
	}

	// Check if it's a numeric offset pattern (+HH:MM / -HH:MM)
	offsetPattern := "[" + tzContent + "]"
	if offsetErr := abnf.AbnfTimezoneTag.ValidateTimezoneTag(offsetPattern, false); offsetErr == nil {
		// Parse numeric offset
		if offset, err := parseNumericOffset(tzContent); err == nil {
			// Keep the RFC 3339 serialization form ("+09:00") as the zone
			// name so Format round-trips the annotation per RFC 9557
			// Section 1.2 and the Section 4.1 time-numoffset grammar.
			ext.Location = time.FixedZone(tzContent, offset)
			ext.CriticalLocation = critical
			return nil
		}
	}

	// Try to validate timezone existence
	if err := abnf.AbnfTimezone.ValidateTimezone(tzContent, tzStrict); err != nil {
		// In non-strict mode, ignore unknown timezone names per RFC 9557
		if !tzStrict {
			return nil
		}
		return ErrInvalidTimezone
	}

	// Additional check: try to load the timezone to see if it actually exists
	if _, err := time.LoadLocation(tzContent); err != nil {
		// In non-strict mode, ignore unknown timezone names per RFC 9557
		if !tzStrict {
			return nil
		}
		return ErrInvalidTimezone
	}

	// Timezone exists - set the location
	ext.Location = time.FixedZone(tzContent, 0)
	ext.CriticalLocation = critical
	return nil
}

// handleExtensionTag processes an extension tag element (key=value pair).
func handleExtensionTag(content string, critical bool, startIdx int, ext *IXDTFExtensions, strict bool) error {
	equalIndex := strings.IndexByte(content[startIdx:], '=')
	if equalIndex < 0 {
		return ErrInvalidExtension
	}
	equalIndex += startIdx
	if equalIndex == startIdx || equalIndex == len(content)-1 {
		return ErrInvalidExtension // empty key or value
	}
	if err := isValidSuffixKeyRange(content, startIdx, equalIndex); err != nil {
		return err
	}
	if err := isValidSuffixValueRange(content, equalIndex+1, len(content)); err != nil {
		return err
	}

	key := content[startIdx:equalIndex]

	// RFC 9557 Section 3.3: for elective duplicates the first occurrence
	// wins, but a duplicate suffix key involving a critical flag on either
	// occurrence MUST be treated as erroneous — in both modes.
	if _, exists := ext.Tags[key]; exists {
		if critical || ext.Critical[key] {
			return ErrCriticalExtension
		}
		return nil
	}
	value := content[equalIndex+1:]
	if critical {
		if err := validateCriticalExtension(key, value); err != nil {
			return err
		}
		// RFC 9557 Section 3.3: a recipient MUST treat the string as
		// erroneous when it cannot process a critical suffix key. In strict
		// mode this library acts as the recipient and only understands
		// "u-ca"; in non-strict mode processing is delegated to the caller
		// via the Critical map.
		if strict && key != ExtensionUnicodeCalendar {
			return ErrCriticalExtension
		}
	}
	ext.Tags[key] = value
	if critical {
		ext.Critical[key] = true
	}
	return nil
}

// validateExtensions validates IXDTF extensions for correctness and processes critical extensions.
// This function checks if all critical extensions can be handled and returns an error if not.
func validateExtensions(ext *IXDTFExtensions) error {
	return validateExtensionsStrict(ext, true)
}

func validateExtensionsStrict(ext *IXDTFExtensions, strict bool) error {
	if ext == nil {
		return nil
	}

	if err := validateLocationStrict(ext.Location, strict); err != nil {
		return err
	}

	if err := validateTagKeys(ext.Tags); err != nil {
		return err
	}

	if err := validateCriticalTags(ext.Tags, ext.Critical); err != nil {
		return err
	}

	if strict {
		if err := validateTagValuesStrict(ext.Tags); err != nil {
			return err
		}
	}

	return nil
}

func validateLocationStrict(location *time.Location, strict bool) error {
	if location == nil {
		return nil
	}
	if ensureRealLocation(location) == nil {
		// An offset-derived FixedZone (e.g. from "[+09:00]") is valid as-is
		// and never resolvable via the timezone database.
		if isOffsetLocationName(location.String()) {
			return nil
		}
		if _, err := loadLocationCached(location.String()); err != nil {
			// In non-strict mode, ignore unknown timezone errors per RFC 9557
			if !strict {
				return nil
			}
			return ErrInvalidTimezone
		}
	}
	return nil
}

func validateTagKeys(tags map[string]string) error {
	// Basic tag key validation (syntactic). Value validation is already handled when creating tags.
	for key := range tags {
		if err := abnf.AbnfSuffixKey.ValidateSuffixKey(key); err != nil {
			return err
		}
	}
	return nil
}

func validateCriticalTags(tags map[string]string, critical map[string]bool) error {
	// Critical tag processing:
	// * If a key is marked critical but missing in Tags -> error.
	// * If present but value is empty -> error.
	for key, isCritical := range critical {
		if !isCritical {
			continue
		}
		value, exists := tags[key]
		if !exists { // missing critical tag
			return ErrCriticalExtension
		}
		if err := validateCriticalExtension(key, value); err != nil {
			return err
		}
	}
	return nil
}

func validateTagValuesStrict(tags map[string]string) error {
	for key, value := range tags {
		if key == ExtensionUnicodeCalendar && !isUnicodeCalendarIdentifier(value) {
			return ErrInvalidTagCalendarIdentifier
		}
	}
	return nil
}

// timezoneCache stores successfully loaded *time.Location by name.
//
//nolint:gochecknoglobals // Package-level cache avoids repeated time.LoadLocation cost; safe read-mostly structure.
var timezoneCache sync.Map // map[string]*time.Location

// loadLocationCached loads a timezone using cache.
func loadLocationCached(name string) (*time.Location, error) {
	if v, ok := timezoneCache.Load(name); ok {
		if loc, ok2 := v.(*time.Location); ok2 {
			return loc, nil
		}
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}
	timezoneCache.Store(name, loc)
	return loc, nil
}

// tryLoadTimezone attempts to treat s as a timezone name (not numeric offset) and load it.
// Returns (location, true) if loaded, (nil, false) otherwise. Errors are treated as non-match.
func tryLoadTimezone(s string) (*time.Location, bool) {
	if s == "" || !abnf.IsTimezoneSyntax(s) {
		return nil, false
	}
	loc, err := loadLocationCached(s)
	if err != nil {
		return nil, false
	}
	return loc, true
}

// parseNumericOffset parses a numeric timezone offset string (e.g., "+09:00", "-05:00")
// and returns the offset in seconds.
func parseNumericOffset(s string) (int, error) {
	const expectedOffsetLength = 6
	if len(s) != expectedOffsetLength {
		return 0, errors.New("invalid offset format")
	}

	sign := 1
	if s[0] == '-' {
		sign = -1
	} else if s[0] != '+' {
		return 0, errors.New("invalid offset sign")
	}

	if s[3] != ':' {
		return 0, errors.New("invalid offset format")
	}

	var err error

	hours, err := strconv.Atoi(s[1:3])
	if err != nil || hours < 0 || hours > 23 {
		return 0, errors.New("invalid hours")
	}

	minutes, err := strconv.Atoi(s[4:6])
	if err != nil || minutes < 0 || minutes > 59 {
		return 0, errors.New("invalid minutes")
	}

	return sign * (hours*3600 + minutes*60), nil
}

// isOffsetLocationName reports whether name is a numeric-offset zone name in
// the RFC 3339 serialization form used for offset time-zone annotations
// (e.g. "+09:00", "-03:30"), as produced when parsing "[+09:00]". Such a
// location's offset is authoritative and has no timezone-database entry, so
// it must not be resolved via time.LoadLocation (the lookup would always
// fail).
func isOffsetLocationName(name string) bool {
	const offsetNameLength = 6 // len("+09:00")
	if len(name) != offsetNameLength || (name[0] != '+' && name[0] != '-') || name[3] != ':' {
		return false
	}
	for _, i := range [...]int{1, 2, 4, 5} {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}
	return true
}

// ensureRealLocation returns nil if the given location appears to be a placeholder FixedZone
// that should be resolved (we created FixedZone(tzName,0) when name didn't load earlier).
// If location is already a real load (cache hit or standard lib) return it.
func ensureRealLocation(loc *time.Location) *time.Location {
	if loc == nil {
		return nil
	}
	name := loc.String()
	// If we already cached a real one, return it.
	if v, ok := timezoneCache.Load(name); ok {
		if realLoc, ok2 := v.(*time.Location); ok2 {
			return realLoc
		}
	}
	// Heuristic: FixedZone created here will have zero offset but arbitrary name; try to resolve.
	if loaded, err := loadLocationCached(name); err == nil {
		return loaded
	}
	return nil
}
