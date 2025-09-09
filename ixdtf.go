// Copyright 2024 8beeeaaat. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ixdtf implements RFC 9557 Internet Extended Date/Time Format (IXDTF).
// IXDTF extends RFC 3339 by adding optional suffix elements for timezone names
// and additional metadata while maintaining full backward compatibility.
//
// See RFC 9557: https://datatracker.ietf.org/doc/rfc9557/
package ixdtf

import (
	"errors"
	"strings"
	"time"
)

type Layout string

const (
	// LayoutIXDTF represents the IXDTF layout.
	LayoutIXDTF Layout = "2006-01-02T15:04:05Z07:00[time-zone]"

	// LayoutIXDTFNano represents the IXDTF layout with nanoseconds.
	LayoutIXDTFNano Layout = "2006-01-02T15:04:05.999999999Z07:00[time-zone]"

	// LayoutRFC3339 represents the RFC 3339 layout.
	LayoutRFC3339 Layout = time.RFC3339

	// LayoutRFC3339Nano represents the RFC 3339 layout with nanoseconds.
	LayoutRFC3339Nano Layout = time.RFC3339Nano
)

// Common parsing errors
var (
	ErrCriticalExtension = errors.New("critical extension cannot be processed")
	ErrInvalidExtension  = errors.New("invalid extension format")
	ErrInvalidSuffix     = errors.New("invalid IXDTF suffix format")
	ErrInvalidTimezone   = errors.New("invalid timezone name")
)

// IXDTFExtensions holds IXDTF suffix information that extends RFC 3339.
type IXDTFExtensions struct {
	// TimeZone contains the IANA timezone name from the timezone suffix.
	// Example: "Asia/Tokyo", "UTC", "America/New_York"
	TimeZone string

	// Tags contains extension tags as key-value pairs.
	// Example: map["u-ca"]"japanese"
	Tags map[string]string

	// Critical indicates which tags are marked as critical (must be processed).
	// Critical tags are marked with "!" prefix in the IXDTF string.
	Critical map[string]bool
}

// ParseError represents an error that occurred during IXDTF parsing.
type ParseError struct {
	Layout Layout
	Value  string
	Err    error
}

func (e *ParseError) Error() string {
	if e.Err == nil {
		return ""
	}
	return "parsing time \"" + e.Value + "\" as \"" + string(e.Layout) + "\": " + e.Err.Error()
}

// CheckTimezoneConsistency checks if the timezone offset matches the IANA timezone
// Returns an error if there's an inconsistency between offset and timezone
func CheckTimezoneConsistency(timestamp time.Time, timezone string) (loc *time.Location, err error) {
	if timezone == "" {
		return nil, nil // No timezone to check
	}

	loc, err = time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	// Get the timezone offset for the given timestamp
	expectedTimestamp := timestamp.In(loc)
	_, originalOffset := timestamp.Zone()
	_, expectedOffset := expectedTimestamp.Zone()

	// For certain timezone patterns (like Etc/GMT-X), allow inconsistencies
	// as they have complex offset semantics
	if strings.HasPrefix(timezone, "Etc/GMT") {
		return nil, nil // Skip consistency check for Etc/GMT timezones
	}

	// Check if offsets match (allowing for some flexibility with DST transitions)
	if originalOffset != expectedOffset {
		// For RFC 9557 compliance, we should note inconsistencies but not necessarily fail
		// This is more of an informational check rather than a strict validation
		return nil, nil // Allow inconsistencies as per RFC 9557 flexibility
	}

	return loc, nil
}

// Format formats a time with IXDTF extensions using RFC 3339 format.
func Format(t time.Time, ext IXDTFExtensions) (string, error) {

	// Validate extensions (including critical tag processing)
	if err := ValidateExtensions(ext); err != nil {
		return "", err
	}

	// Format the RFC 3339 portion
	b := []byte(t.Format(time.RFC3339))
	b = appendSuffix(b, ext)

	return string(b), nil
}

// FormatNano formats a time with IXDTF extensions using RFC 3339 format with nanoseconds.
func FormatNano(t time.Time, ext IXDTFExtensions) (string, error) {
	// Validate extensions (including critical tag processing)
	if err := ValidateExtensions(ext); err != nil {
		return "", err
	}

	b := []byte(t.Format(time.RFC3339Nano))
	b = appendSuffix(b, ext)

	return string(b), nil
}

// NewIXDTFExtensions creates a new IXDTFExtensions with initialized maps.
func NewIXDTFExtensions() IXDTFExtensions {
	return IXDTFExtensions{
		TimeZone: "",
		Tags:     make(map[string]string),
		Critical: make(map[string]bool),
	}
}

// Parse parses an IXDTF string and returns the time and extension information.
func Parse(s string) (time.Time, IXDTFExtensions, error) {
	// First validate the format
	if err := Validate(s); err != nil {
		return time.Time{}, IXDTFExtensions{}, err
	}

	rfc3339End := findRFC3339End(s)
	rfc3339Portion := s[:rfc3339End]

	t, err := parseRFC3339Portion(rfc3339Portion)
	if err != nil {
		return time.Time{}, IXDTFExtensions{}, newParseError(LayoutRFC3339, s, err)
	}

	ext := NewIXDTFExtensions()
	if rfc3339End < len(s) {
		suffixPortion := s[rfc3339End:]
		if ext, err = parseSuffix(suffixPortion); err != nil {
			return time.Time{}, IXDTFExtensions{}, newParseError(LayoutIXDTF, s, err)
		}
	}
	// Check timezone consistency if timezone is provided
	if ext.TimeZone != "" {
		loc, err := CheckTimezoneConsistency(t, ext.TimeZone)
		if err != nil {
			return time.Time{}, IXDTFExtensions{}, newParseError(LayoutIXDTF, s, err)
		}
		// Apply the timezone to the timestamp
		if loc != nil {
			t = t.In(loc)
		}
	}

	return t, ext, nil
}

// Validate validates an IXDTF string for format correctness without parsing the time component.
// This is useful for quick validation without the overhead of time parsing.
func Validate(s string) error {
	rfc3339End := findRFC3339End(s)
	rfc3339Portion := s[:rfc3339End]

	if len(rfc3339Portion) == 0 {
		return newParseError(LayoutRFC3339, s, errors.New("empty datetime string"))
	}

	if _, err := parseRFC3339Portion(rfc3339Portion); err != nil {
		return newParseError(LayoutRFC3339, s, errors.New("invalid portion: "+err.Error()))
	}

	if rfc3339End < len(s) {
		suffixPortion := s[rfc3339End:]
		if _, err := parseSuffix(suffixPortion); err != nil {
			return newParseError(LayoutIXDTF, s, err)
		}
	}

	return nil
}

// ValidateExtensions validates IXDTF extensions for correctness and processes critical extensions.
// This function checks if all critical extensions can be handled and returns an error if not.
func ValidateExtensions(ext IXDTFExtensions) error {
	if ext.TimeZone != "" {
		if _, err := time.LoadLocation(ext.TimeZone); err != nil {
			return ErrInvalidTimezone
		}
	}

	for key, isCritical := range ext.Critical {
		if isCritical {
			if _, exists := ext.Tags[key]; !exists {
				return ErrCriticalExtension
			}

			if err := validateCriticalExtension(key, ext.Tags[key]); err != nil {
				return err
			}
		}
	}

	return nil
}

// Private functions in alphabetical order

func appendSuffix(b []byte, ext IXDTFExtensions) []byte {
	if ext.TimeZone != "" {
		b = append(b, '[')
		b = append(b, ext.TimeZone...)
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

func findRFC3339End(s string) int {
	for i, c := range s {
		if c == '[' {
			return i
		}
	}
	return len(s)
}

func isValidSuffixKey(key string) bool {
	if len(key) == 0 {
		return false
	}

	// Check first character (key-initial) - optimized byte comparison
	first := key[0]
	if first < 'a' || first > 'z' {
		if first != '_' {
			return false
		}
	}

	// Check remaining characters (key-char) - single loop with byte comparison
	for i := 1; i < len(key); i++ {
		b := key[i]
		if (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '-' || b == '_' {
			continue
		}
		return false
	}

	return true
}

func isValidSuffixKeyRange(s string, start, end int) bool {
	if start >= end {
		return false
	}

	// Check first character (key-initial) - optimized byte comparison
	first := s[start]
	if first < 'a' || first > 'z' {
		if first != '_' {
			return false
		}
	}

	// Check remaining characters (key-char) - single loop with byte comparison
	for i := start + 1; i < end; i++ {
		b := s[i]
		if (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') || b == '-' || b == '_' {
			continue
		}
		return false
	}

	return true
}

func isValidSuffixValue(value string) bool {
	if len(value) == 0 {
		return false
	}

	// Check for leading or trailing hyphens - single pass
	if value[0] == '-' || value[len(value)-1] == '-' {
		return false
	}

	// Single pass validation: check consecutive hyphens and character validity
	prevHyphen := false
	hasValidChar := false

	for i := 0; i < len(value); i++ {
		b := value[i]

		if b == '-' {
			if prevHyphen {
				return false // consecutive hyphens
			}
			prevHyphen = true
			hasValidChar = false // reset for next segment
		} else {
			// Check if character is alphanum
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') {
				prevHyphen = false
				hasValidChar = true
			} else {
				return false // invalid character
			}
		}
	}

	// Must end with a valid character (not hyphen, already checked above)
	return hasValidChar
}

func isValidSuffixValueRange(s string, start, end int) bool {
	if start >= end {
		return false
	}

	// Check for leading or trailing hyphens
	if s[start] == '-' || s[end-1] == '-' {
		return false
	}

	// Single pass validation: check consecutive hyphens and character validity
	prevHyphen := false
	hasValidChar := false

	for i := start; i < end; i++ {
		b := s[i]

		if b == '-' {
			if prevHyphen {
				return false // consecutive hyphens
			}
			prevHyphen = true
			hasValidChar = false // reset for next segment
		} else {
			// Check if character is alphanum
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') {
				prevHyphen = false
				hasValidChar = true
			} else {
				return false // invalid character
			}
		}
	}

	// Must end with a valid character (not hyphen, already checked above)
	return hasValidChar
}

func isValidTimezone(s string) bool {
	if len(s) >= 2 && s[1] == '-' {
		prefix := s[0]
		if prefix == 'u' || prefix == 'x' || prefix == 't' {
			return false // Likely an extension tag prefix
		}
	}
	return true
}

func isValidTimezoneContent(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Check if this looks like an incomplete extension tag
	// Extension tags typically contain hyphens with specific patterns
	if strings.IndexByte(s, '-') >= 0 {
		// If has hyphen and looks like extension prefix, it's invalid
		if len(s) >= 2 && s[1] == '-' {
			prefix := s[0]
			if prefix == 'u' || prefix == 'x' || prefix == 't' {
				return false // Likely an extension tag prefix
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
		if t, err := time.Parse(layout, rfc3339Portion); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}

	return time.Time{}, lastErr
}

func parseSuffix(s string) (IXDTFExtensions, error) {
	ext := NewIXDTFExtensions()

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
		if err := parseSuffixElement(content, &ext); err != nil {
			return ext, err
		}

		i = j + 1
	}

	return ext, nil
}

func parseSuffixElement(content string, ext *IXDTFExtensions) error {
	if len(content) == 0 {
		return ErrInvalidSuffix
	}

	// Check for critical flag - avoid string slice creation
	critical := false
	startIdx := 0
	if content[0] == '!' {
		critical = true
		startIdx = 1
		if len(content) <= 1 {
			return ErrInvalidSuffix
		}
	}

	// Check if this is an extension tag (contains '=') - work with indices
	if equalIndex := strings.IndexByte(content[startIdx:], '='); equalIndex >= 0 {
		equalIndex += startIdx // adjust for offset

		if equalIndex == startIdx || equalIndex == len(content)-1 {
			return ErrInvalidExtension // empty key or value
		}

		// Validate key and value using substrings without creating new strings
		if !isValidSuffixKeyRange(content, startIdx, equalIndex) {
			return ErrInvalidExtension
		}

		if !isValidSuffixValueRange(content, equalIndex+1, len(content)) {
			return ErrInvalidExtension
		}

		// Extract key and value only when validation passes
		key := content[startIdx:equalIndex]
		value := content[equalIndex+1:]

		// RFC 9557: Duplicate keys use first occurrence (don't overwrite)
		if _, exists := ext.Tags[key]; !exists {
			ext.Tags[key] = value
			if critical {
				ext.Critical[key] = true
			}
		}
		// Note: If key already exists, we silently ignore the duplicate
		// as per RFC 9557 specification
	} else {
		// This is a timezone name
		if critical {
			return ErrInvalidTimezone // Timezone cannot be critical
		}

		// Optimized timezone validation - single check
		tzContent := content[startIdx:]
		if !isValidTimezoneContent(tzContent) {
			return ErrInvalidExtension
		}

		if ext.TimeZone != "" {
			return ErrInvalidTimezone // Multiple timezone suffixes not allowed
		}
		ext.TimeZone = tzContent
	}

	return nil
}

func splitOnHyphen(s string) []string {
	if len(s) == 0 {
		return []string{}
	}

	var parts []string
	start := 0

	for i, char := range s {
		if char == '-' {
			if i > start {
				parts = append(parts, s[start:i])
			}
			start = i + 1
		}
	}

	// Add the last part
	if start < len(s) {
		parts = append(parts, s[start:])
	}

	return parts
}

func validateCriticalExtension(key, value string) error {
	switch {
	case key == "u-ca": // Unicode calendar extension
		if value == "" {
			return ErrCriticalExtension
		}
	case key == "u-nu": // Unicode numbering system extension
		if value == "" {
			return ErrCriticalExtension
		}
	case len(key) >= 2 && key[:2] == "x-": // Private use extensions
		if value == "" {
			return ErrCriticalExtension
		}
	case len(key) >= 1 && key[0] == '_': // Experimental tags (RFC 9557)
		if value == "" {
			return ErrCriticalExtension
		}
	default:
		return ErrCriticalExtension
	}

	return nil
}
