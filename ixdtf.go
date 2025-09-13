// Copyright 2025 8beeeaaat. All rights reserved.
// Use of this source code is governed by a BSD-style.
// license that can be found in the LICENSE file.

// Package ixdtf implements RFC 9557 Internet Extended Date/Time Format (IXDTF).
// IXDTF extends RFC 3339 by adding optional suffix elements for timezone names.
// and additional metadata while maintaining full backward compatibility.
//
// See RFC 9557: https://datatracker.ietf.org/doc/rfc9557/
package ixdtf

import (
	"errors"
	"strings"
	"time"
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
	ErrCriticalExtension      = errors.New("critical extension cannot be processed")
	ErrExperimentalExtension  = errors.New("experimental extension cannot be processed")
	ErrPrivateExtension       = errors.New("private extension cannot be processed")
	ErrInvalidExtension       = errors.New("invalid extension format")
	ErrInvalidSuffix          = errors.New("invalid IXDTF suffix format")
	ErrInvalidTimezone        = errors.New("invalid timezone name")
	ErrTimezoneOffsetMismatch = errors.New("timezone offset does not match the specified timezone")
)

// IXDTFExtensions holds IXDTF suffix information that extends RFC 3339.
//
//nolint:revive // Keeping existing public API name for compatibility
type IXDTFExtensions struct {
	Location *time.Location

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
	// Skipped indicates if the consistency check was skipped (e.g., for Etc/GMT).
	Skipped bool
}

func (e *ParseError) Error() string {
	if e.Err == nil {
		return ""
	}
	return "parsing time \"" + e.Value + "\" as \"" + string(e.Layout) + "\": " + e.Err.Error()
}

// Format formats a time with IXDTF extensions using RFC 3339 format.
func Format(t time.Time, ext *IXDTFExtensions) (string, error) {
	// Validate extensions (including critical tag processing).
	if err := validateExtensions(ext); err != nil {
		return "", err
	}

	// Format the RFC 3339 portion.
	b := []byte(t.Format(time.RFC3339))
	b = appendSuffix(b, ext)

	return string(b), nil
}

// FormatNano formats a time with IXDTF extensions using RFC 3339 format with nanoseconds.
func FormatNano(t time.Time, ext *IXDTFExtensions) (string, error) {
	// Validate extensions (including critical tag processing).
	if err := validateExtensions(ext); err != nil {
		return "", err
	}

	b := []byte(t.Format(time.RFC3339Nano))
	b = appendSuffix(b, ext)

	return string(b), nil
}

// NewIXDTFExtensionsArgs contains the arguments for creating IXDTFExtensions.
type NewIXDTFExtensionsArgs struct {
	Location *time.Location
	Tags     map[string]string
	Critical map[string]bool
}

// NewIXDTFExtensions creates a new IXDTFExtensions with initialized maps.
func NewIXDTFExtensions(args *NewIXDTFExtensionsArgs) *IXDTFExtensions {
	ext := &IXDTFExtensions{
		Location: time.UTC,
		Tags:     make(map[string]string),
		Critical: make(map[string]bool),
	}

	if args != nil {
		if args.Location != nil {
			ext.Location = args.Location
		}
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
		if ext, err = parseSuffix(suffixPortion); err != nil {
			return time.Time{}, nil, newParseError(LayoutRFC3339Extended, s, err)
		}
	}

	// Apply the timezone to the timestamp if provided
	if ext.Location != nil {
		result, checkErr := checkTimezoneConsistency(t, ext.Location, strict)
		if checkErr != nil {
			return time.Time{}, nil, newParseError(
				LayoutRFC3339NanoExtended,
				s,
				checkErr,
			)
		}
		if result.Location != nil {
			t = t.In(result.Location)
		}
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
		if ext, err = parseSuffix(suffixPortion); err != nil {
			return newParseError(LayoutRFC3339Extended, s, err)
		}
	}

	// Check timezone consistency if timezone is provided
	if ext != nil && ext.Location != nil {
		_, tzErr := checkTimezoneConsistency(t, ext.Location, strict)
		if tzErr != nil {
			return newParseError(LayoutRFC3339NanoExtended, s, tzErr)
		}
		// Note: Inconsistencies are not treated as validation errors per RFC 9557
		// but invalid timezone names will cause errors
	}

	// Optional: check if the complete string matches the ABNF pattern.
	// This serves as an additional validation layer, but only for well-formed strings.
	// Skip for empty strings and malformed inputs to avoid double error reporting.
	if rfc3339Portion != "" && (rfc3339End >= len(s) || s[rfc3339End:] != "") {
		if abnfErr := AbnfDateTimeExt.Validate(s); abnfErr != nil {
			// Only report ABNF mismatch if basic validation passes.
			// This prevents confusing error messages for clearly invalid input.
			return newParseError(LayoutRFC3339Extended, s, abnfErr)
		}
	}

	return nil
}

// Private functions in alphabetical order.
func appendSuffix(b []byte, ext *IXDTFExtensions) []byte {
	if ext.Location != nil {
		b = append(b, '[')
		b = append(b, ext.Location.String()...)
		b = append(b, ']')
	}

	keys := make([]string, 0, len(ext.Tags))
	for key := range ext.Tags {
		keys = append(keys, key)
	}

	for i := 1; i < len(keys); i++ {
		key := keys[i]
		j := i - 1
		for j >= 0 && keys[j] > key {
			keys[j+1] = keys[j]
			j--
		}
		keys[j+1] = key
	}

	for _, key := range keys {
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
// If strict is true, returns an error when offsets don't match (except for Etc/GMT patterns).
// Returns consistency information and an error for timezone loading failures or strict mode mismatches.
func checkTimezoneConsistency(
	timestamp time.Time,
	location *time.Location,
	strict bool,
) (*TimezoneConsistencyResult, error) {
	result := &TimezoneConsistencyResult{
		Location: location,
	}

	if location == nil {
		result.IsConsistent = true // No timezone means no inconsistency
		return result, nil
	}

	loc, err := time.LoadLocation(location.String())
	if err != nil {
		return result, err
	}
	result.Location = loc

	// Get the timezone offset for the given timestamp.
	expectedTimestamp := timestamp.In(loc)
	_, originalOffset := timestamp.Zone()
	_, expectedOffset := expectedTimestamp.Zone()

	result.OriginalOffset = originalOffset
	result.ExpectedOffset = expectedOffset

	// For certain timezone patterns (like Etc/GMT-X), skip consistency checks.
	// These have POSIX-inverted offset semantics that cause false inconsistency reports.
	if strings.HasPrefix(loc.String(), "Etc/GMT") {
		result.Skipped = true
		result.IsConsistent = true // Assume consistent for Etc/GMT patterns
		return result, nil
	}

	// Check if offsets match (allowing for some flexibility with DST transitions).
	result.IsConsistent = (originalOffset == expectedOffset)

	// In strict mode, return an error for inconsistencies
	if strict && !result.IsConsistent {
		return nil, ErrTimezoneOffsetMismatch
	}

	// Per RFC 9557, inconsistencies should be detectable but not cause failures.
	// This allows applications to handle inconsistencies as needed.
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

func isValidSuffixKeyRange(s string, start, end int) error {
	if start >= end {
		return ErrInvalidExtension
	}
	return AbnfSuffixKey.Validate(s[start:end])
}

func isValidSuffixValue(value string) error {
	if value == "" {
		return nil
	}
	// Use ABNF pattern for basic validation, then check additional constraints
	if err := AbnfSuffixValues.Validate(value); err != nil {
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
func validateCriticalExtension(_, value string) error {
	if value == "" { // empty value not allowed for critical
		return ErrCriticalExtension
	}
	// Known-key specific validation hooks could be added here in future.
	return nil
}

func isValidTimezoneContent(s string) bool {
	if s == "" {
		return false
	}

	// Check if this looks like an incomplete extension tag
	// Extension tags typically contain hyphens with specific patterns
	if strings.IndexByte(s, '-') >= 0 {
		// If has hyphen and looks like extension prefix, it's invalid
		if len(s) >= 2 && s[1] == '-' {
			prefix := s[0]
			if prefix == 'u' || prefix == 'x' || prefix == 't' {
				return false // Likely an extension tag prefix.
			}
		}
	}

	return true
}

func newParseError(layout Layout, value string, err error) *ParseError {
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

func parseSuffix(s string) (*IXDTFExtensions, error) {
	ext := NewIXDTFExtensions(nil)

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
		if err := parseSuffixElement(content, ext); err != nil {
			return ext, err
		}

		i = j + 1
	}

	return ext, nil
}

func parseSuffixElement(content string, ext *IXDTFExtensions) error {
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
		return handleExtensionTag(content, critical, startIdx, ext)
	}

	// Timezone name handling.
	if critical {
		return ErrInvalidTimezone // Timezone cannot be critical
	}
	tzContent := content[startIdx:]
	if tzContent != "" {
		// Validate using ABNF pattern first - check both timezone name and offset patterns
		if err := AbnfTimezoneName.Validate(tzContent); err != nil {
			// If not a timezone name, check if it's a valid timezone offset (+HH:MM, -HH:MM)
			offsetPattern := "[" + tzContent + "]"
			if offsetErr := AbnfTimezone.Validate(offsetPattern); offsetErr != nil {
				return ErrInvalidExtension
			}
		}
		if !isValidTimezoneContent(tzContent) {
			return ErrInvalidExtension
		}
		ext.Location = time.FixedZone(tzContent, 0)
	}
	return nil
}

// handleExtensionTag processes an extension tag element (key=value pair).
func handleExtensionTag(content string, critical bool, startIdx int, ext *IXDTFExtensions) error {
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
	// Respect RFC 9557: first occurrence wins.
	if _, exists := ext.Tags[key]; exists {
		return nil
	}
	value := content[equalIndex+1:]
	ext.Tags[key] = value
	if critical {
		ext.Critical[key] = true
	}
	return nil
}

// validateExtensions validates IXDTF extensions for correctness and processes critical extensions.
// This function checks if all critical extensions can be handled and returns an error if not.
func validateExtensions(ext *IXDTFExtensions) error {
	if ext == nil {
		return nil
	}
	if ext.Location != nil {
		if _, err := time.LoadLocation(ext.Location.String()); err != nil {
			return ErrInvalidTimezone
		}
	}

	// Basic tag key validation (syntactic). Value validation is already handled when creating tags.
	for key := range ext.Tags {
		if err := AbnfSuffixKey.Validate(key); err != nil {
			return err
		}
	}

	// Critical tag processing:
	// * If a key is marked critical but missing in Tags -> error.
	// * If present but value is empty -> error.
	for key, isCritical := range ext.Critical {
		if !isCritical {
			continue
		}
		value, exists := ext.Tags[key]
		if !exists { // missing critical tag
			return ErrCriticalExtension
		}
		if err := validateCriticalExtension(key, value); err != nil {
			return err
		}
	}

	return nil
}
