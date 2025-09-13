// ABNF patterns and validation for RFC 9557 Internet Extended Date/Time Format (IXDTF).
// https://www.rfc-editor.org/rfc/rfc9557.html#section-4.1
//
// This file contains compiled regular expression patterns that correspond to the.
// ABNF (Augmented Backus-Naur Form) definitions specified in RFC 9557.
//
// Note: Due to limitations in Go's regexp engine, some ABNF constraints cannot.
// be fully expressed in regular expressions alone and require additional validation logic:
//   - Negative lookaheads (e.g., "but not '.' or '..'")
//   - Context-dependent rules
//   - Semantic validation (timezone existence, date validity)
//
// For complete validation, use these patterns in combination with the parsing.
// functions provided in the main ixdtf package.

// Package abnf provides ABNF patterns and validation for RFC 9557 Internet Extended Date/Time Format (IXDTF).
// https://www.rfc-editor.org/rfc/rfc9557.html#section-4.1
package abnf

import (
	"errors"
	"regexp"
	"time"
)

// Common validation errors.
var (
	ErrPrivateExtension      = errors.New("private extension cannot be processed")
	ErrExperimentalExtension = errors.New("experimental extension cannot be processed")
)

type Abnf struct {
	regexp *regexp.Regexp
}

func newAbnf(pattern string) *Abnf {
	return &Abnf{regexp: regexp.MustCompile(pattern)}
}

// Syntax Extensions to RFC 3339.
// Note: Some ABNF constraints (like negative lookaheads and context-dependent rules).
// cannot be fully expressed in Go's regexp engine and require additional validation logic.
//
//nolint:gochecknoglobals // ABNF patterns are constants used for validation
var (
	AbnfTimezoneName = newAbnf(`^[A-Za-z_][A-Za-z._0-9+-]*(/[A-Za-z_][A-Za-z._0-9+-]*)*$`)
	AbnfTimezone     = newAbnf(
		`^\[!?[A-Za-z_][A-Za-z._0-9+-]*(/[A-Za-z_][A-Za-z._0-9+-]*)*\]$|^\[!?[+-][0-9]{2}:[0-9]{2}\]$`,
	)

	AbnfSuffixKey    = newAbnf(`^[a-z_][a-z_0-9-]*$`)
	AbnfSuffixValues = newAbnf(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`)

	AbnfDateTimeExt = newAbnf(
		`^[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])T([01][0-9]|2[0-3]):[0-9]{2}:[0-9]{2}(\.[0-9]+)?(Z|[+-][0-9]{2}:[0-9]{2})(?:\[\]|\[!?[A-Za-z._0-9+/:-]+\]|\[!?[a-z_][a-z_0-9-]*=[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*\])*$`,
	)
)

// IsTimezoneNameSyntax returns true if the input matches the lexical pattern of a timezone name.
// (ABNF for tz name) without performing existence (time.LoadLocation) validation.
func IsTimezoneNameSyntax(input string) bool {
	return AbnfTimezoneName.regexp.MatchString(input)
}

// Validate validates a string against a specific ABNF pattern.
func (a *Abnf) Validate(input string) error {
	if !a.regexp.MatchString(input) {
		return errors.New("invalid extension format")
	}
	// additional validation for patterns that cannot be fully expressed with regex
	switch a {
	case AbnfSuffixKey:
		return validateSuffixKey(input)
	case AbnfTimezoneName:
		return validateTimezoneName(input)
	default:
		return nil
	}
}

func validateTimezoneName(name string) error {
	if !IsTimezoneNameSyntax(name) {
		return errors.New("invalid timezone name syntax")
	}
	_, err := time.LoadLocation(name)
	return err
}

func validateSuffixKey(key string) error {
	switch {
	case len(key) >= 2 && (key[:2] == "x-" || key[:2] == "X-"): // Private use extensions. https://www.rfc-editor.org/info/rfc6648
		return ErrPrivateExtension
	case len(key) >= 1 && key[0] == '_': // Experimental extensions. https://www.rfc-editor.org/info/rfc6648
		return ErrExperimentalExtension
	}
	return nil
}
