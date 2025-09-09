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
	"time"
)

// Layout constants for IXDTF formatting.
const (
	// IXDTF is the layout for RFC 9557 format with timezone suffix.
	// Example: 2006-01-02T15:04:05Z[UTC]
	IXDTF = "2006-01-02T15:04:05Z07:00[time-zone]"

	// IXDTFNano is like IXDTF but with nanoseconds.
	// Example: 2006-01-02T15:04:05.999999999Z[UTC]
	IXDTFNano = "2006-01-02T15:04:05.999999999Z07:00[time-zone]"
)

// IXDTFExtensions holds IXDTF suffix information that extends RFC 3339.
type IXDTFExtensions struct {
	// TimeZone contains the IANA timezone name from the timezone suffix.
	// Example: "Asia/Tokyo", "UTC", "America/New_York"
	TimeZone string

	// Tags contains extension tags as key-value pairs.
	// Example: map["u-ca"]"japanese" for calendar extension
	Tags map[string]string

	// Critical indicates which tags are marked as critical (must be processed).
	// Critical tags are marked with "!" prefix in the IXDTF string.
	Critical map[string]bool
}

// NewIXDTFExtensions creates a new IXDTFExtensions with initialized maps.
func NewIXDTFExtensions() IXDTFExtensions {
	return IXDTFExtensions{
		Tags:     make(map[string]string),
		Critical: make(map[string]bool),
	}
}

// ParseError represents an error that occurred during IXDTF parsing.
type ParseError struct {
	Layout string
	Value  string
	Msg    string
}

func (e *ParseError) Error() string {
	if e.Msg == "" {
		return "parsing time \"" + e.Value + "\" as \"" + e.Layout + "\": cannot parse"
	}
	return "parsing time \"" + e.Value + "\" as \"" + e.Layout + "\": " + e.Msg
}

// Common parsing errors
var (
	ErrInvalidSuffix     = errors.New("invalid IXDTF suffix format")
	ErrInvalidTimezone   = errors.New("invalid timezone name")
	ErrInvalidExtension  = errors.New("invalid extension format")
	ErrCriticalExtension = errors.New("critical extension cannot be processed")
)

func findRFC3339End(s string) int {
	for i, c := range s {
		if c == '[' {
			return i
		}
	}
	return len(s)
}

// parseSuffix parses IXDTF suffix elements from the string starting at the given position.
// The suffix format is: [timezone][extension1][extension2]...
// Returns the parsed extensions and an error if parsing fails.
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

// parseSuffixElement parses a single suffix element (content between [ and ]).
func parseSuffixElement(content string, ext *IXDTFExtensions) error {
	if len(content) == 0 {
		return ErrInvalidSuffix
	}

	// Check for critical flag
	critical := false
	if content[0] == '!' {
		critical = true
		content = content[1:]
		if len(content) == 0 {
			return ErrInvalidSuffix
		}
	}

	// Check if this is an extension tag (contains '=')
	if equalIndex := findEqual(content); equalIndex >= 0 {
		key := content[:equalIndex]
		value := content[equalIndex+1:]

		if key == "" || value == "" {
			return ErrInvalidExtension
		}

		// Validate key format according to RFC 9557 ABNF grammar
		if !isValidSuffixKey(key) {
			return ErrInvalidExtension
		}

		// Validate value format according to RFC 9557 ABNF grammar
		if !isValidSuffixValue(value) {
			return ErrInvalidExtension
		}

		ext.Tags[key] = value
		if critical {
			ext.Critical[key] = true
		}
	} else {
		// This is a timezone name, but validate it's not an invalid extension format
		if critical {
			return ErrInvalidTimezone // Timezone cannot be critical
		}

		// Check if this looks like an incomplete extension tag
		// Extension tags typically contain hyphens or specific patterns
		if containsHyphen(content) && !isValidTimezone(content) {
			return ErrInvalidExtension // Looks like incomplete extension tag
		}

		if ext.TimeZone != "" {
			return ErrInvalidTimezone // Multiple timezone suffixes not allowed
		}
		ext.TimeZone = content
	}

	return nil
}

func findEqual(s string) int {
	for i, c := range s {
		if c == '=' {
			return i
		}
	}
	return -1
}

func containsHyphen(s string) bool {
	for _, c := range s {
		if c == '-' {
			return true
		}
	}
	return false
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

// isValidSuffixKey validates suffix keys according to RFC 9557 ABNF grammar.
// suffix-key = key-initial *key-char
// key-initial = lcalpha / "_"
// key-char = key-initial / DIGIT / "-"
func isValidSuffixKey(key string) bool {
	if len(key) == 0 {
		return false
	}

	// Check first character (key-initial)
	first := key[0]
	if !((first >= 'a' && first <= 'z') || first == '_') {
		return false
	}

	// Check remaining characters (key-char)
	for i := 1; i < len(key); i++ {
		char := key[i]
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return false
		}
	}

	return true
}

// isValidSuffixValue validates suffix values according to RFC 9557 ABNF grammar.
// suffix-value = 1*alphanum
// suffix-values = suffix-value *("-" suffix-value)
// alphanum = ALPHA / DIGIT
func isValidSuffixValue(value string) bool {
	if len(value) == 0 {
		return false
	}

	// Check for leading or trailing hyphens
	if value[0] == '-' || value[len(value)-1] == '-' {
		return false
	}

	// Check for consecutive hyphens
	for i := 0; i < len(value)-1; i++ {
		if value[i] == '-' && value[i+1] == '-' {
			return false
		}
	}

	// Split on hyphens to validate each part
	parts := splitOnHyphen(value)
	if len(parts) == 0 {
		return false
	}

	for _, part := range parts {
		if len(part) == 0 {
			return false // Empty parts not allowed
		}

		// Each part must be 1*alphanum
		for _, char := range part {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
				return false
			}
		}
	}

	return true
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

// Parse parses an IXDTF string and returns the time and extension information.
func Parse(s string) (time.Time, IXDTFExtensions, error) {
	rfc3339End := findRFC3339End(s)

	rfc3339Portion := s[:rfc3339End]

	var t time.Time
	var err error

	if t, err = time.Parse(time.RFC3339Nano, rfc3339Portion); err != nil {
		if t, err = time.Parse(time.RFC3339, rfc3339Portion); err != nil {
			return time.Time{}, IXDTFExtensions{}, &ParseError{
				Layout: "RFC3339",
				Value:  s,
				Msg:    "invalid RFC 3339 portion: " + err.Error(),
			}
		}
	}

	ext := NewIXDTFExtensions()
	if rfc3339End < len(s) {
		suffixPortion := s[rfc3339End:]
		var err error
		if ext, err = parseSuffix(suffixPortion); err != nil {
			return time.Time{}, IXDTFExtensions{}, &ParseError{
				Layout: "RFC9557",
				Value:  s,
				Msg:    err.Error(),
			}
		}
	}

	if ext.TimeZone != "" {
		if loc, err := time.LoadLocation(ext.TimeZone); err == nil {
			t = t.In(loc)
		}
	}

	return t, ext, nil
}

// Format formats a time with IXDTF extensions using RFC 3339 format.
func Format(t time.Time, ext IXDTFExtensions) string {
	// Format the RFC 3339 portion
	b := []byte(t.Format(time.RFC3339))
	b = appendSuffix(b, ext)

	return string(b)
}

// FormatNano formats a time with IXDTF extensions using RFC 3339 format with nanoseconds.
func FormatNano(t time.Time, ext IXDTFExtensions) string {
	b := []byte(t.Format(time.RFC3339Nano))
	b = appendSuffix(b, ext)

	return string(b)
}

// Validate validates an IXDTF string for format correctness without parsing the time component.
// This is useful for quick validation without the overhead of time parsing.
func Validate(s string) error {
	rfc3339End := findRFC3339End(s)

	rfc3339Portion := s[:rfc3339End]
	if len(rfc3339Portion) == 0 {
		return &ParseError{
			Layout: "RFC3339",
			Value:  s,
			Msg:    "empty datetime string",
		}
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
	}

	var lastErr error
	for _, layout := range layouts {
		if _, err := time.Parse(layout, rfc3339Portion); err == nil {
			break
		} else {
			lastErr = err
		}
	}

	if lastErr != nil {
		return &ParseError{
			Layout: "RFC3339",
			Value:  s,
			Msg:    "invalid RFC 3339 portion: " + lastErr.Error(),
		}
	}

	if rfc3339End < len(s) {
		suffixPortion := s[rfc3339End:]
		if _, err := parseSuffix(suffixPortion); err != nil {
			return &ParseError{
				Layout: "RFC9557",
				Value:  s,
				Msg:    err.Error(),
			}
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
	default:
		return ErrCriticalExtension
	}

	return nil
}

// IsValidIXDTF performs a quick check if a string could be a valid IXDTF format.
// This is a lightweight check that doesn't parse the full datetime.
func IsValidIXDTF(s string) bool {
	return Validate(s) == nil
}
