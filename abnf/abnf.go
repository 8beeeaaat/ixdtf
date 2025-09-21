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

var (
	errInvalidExtensionFormat = errors.New("invalid extension format")
	errUnknownDateTimeExt     = errors.New("unknown ABNF for date-time extension")
	errUnknownSuffixKey       = errors.New("unknown ABNF for suffix key")
	errUnknownSuffixValues    = errors.New("unknown ABNF for suffix values")
	errUnknownTimezone        = errors.New("unknown ABNF for timezone name")
	errUnknownTimezoneTag     = errors.New("unknown ABNF for timezone tag")
)

func (a *Abnf) ensure(expected *Abnf, err error) error {
	if a != expected {
		return err
	}
	return nil
}

func (a *Abnf) ensurePattern(input string) error {
	if !a.regexp.MatchString(input) {
		return errInvalidExtensionFormat
	}
	return nil
}

// Syntax Extensions to RFC 3339.
// Note: Some ABNF constraints (like negative lookaheads and context-dependent rules).
// cannot be fully expressed in Go's regexp engine and require additional validation logic.
//
//nolint:gochecknoglobals // ABNF patterns are constants used for validation
var (
	AbnfTimezone    = newAbnf(`^[A-Za-z_][A-Za-z._0-9+-]*(/[A-Za-z_][A-Za-z._0-9+-]*)*$`)
	AbnfTimezoneTag = newAbnf(
		`^\[!?[A-Za-z_][A-Za-z._0-9+-]*(/[A-Za-z_][A-Za-z._0-9+-]*)*\]$|^\[!?[+-][0-9]{2}:[0-9]{2}\]$`,
	)

	AbnfSuffixKey    = newAbnf(`^[a-z_][a-z_0-9-]*$`)
	AbnfSuffixValues = newAbnf(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`)

	AbnfDateTimeExt = newAbnf(
		`^[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])T([01][0-9]|2[0-3]):[0-9]{2}:[0-9]{2}(\.[0-9]+)?(Z|[+-][0-9]{2}:[0-9]{2})(?:\[\]|\[!?[A-Za-z._0-9+/:-]+\]|\[!?[a-z_][a-z_0-9-]*=[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*\])*$`,
	)
)

// IsTimezoneSyntax returns true if the input matches the lexical pattern of a timezone name.
// (ABNF for tz name) without performing existence (time.LoadLocation) validation.
func IsTimezoneSyntax(input string) bool {
	return AbnfTimezone.regexp.MatchString(input)
}

// ValidateDateTimeExt validates a date-time string with extensions according to the ABNF and additional rules.
func (a *Abnf) ValidateDateTimeExt(input string) error {
	if err := a.ensure(AbnfDateTimeExt, errUnknownDateTimeExt); err != nil {
		return err
	}
	return a.ensurePattern(input)
}

// ValidateSuffixKey validates a suffix key according to the ABNF and additional rules.
func (a *Abnf) ValidateSuffixKey(input string) error {
	if err := a.ensure(AbnfSuffixKey, errUnknownSuffixKey); err != nil {
		return err
	}
	if err := a.ensurePattern(input); err != nil {
		return err
	}
	return validateSuffixKey(input)
}

// ValidateSuffixValues validates suffix values according to the ABNF and additional rules.
func (a *Abnf) ValidateSuffixValues(input string) error {
	if err := a.ensure(AbnfSuffixValues, errUnknownSuffixValues); err != nil {
		return err
	}
	return a.ensurePattern(input)
}

func (a *Abnf) ValidateTimezone(input string, strict bool) error {
	if err := a.ensure(AbnfTimezone, errUnknownTimezone); err != nil {
		return err
	}
	if err := a.ensurePattern(input); err != nil {
		return err
	}
	return validateTimezone(input, strict)
}

func (a *Abnf) ValidateTimezoneTag(input string, strict bool) error {
	if err := a.ensure(AbnfTimezoneTag, errUnknownTimezoneTag); err != nil {
		return err
	}
	if err := a.ensurePattern(input); err != nil {
		return err
	}
	return validateTimezoneTag(input, strict)
}

func validateTimezone(name string, strict bool) error {
	if !strict {
		// In non-strict mode, skip time.LoadLocation validation
		return nil
	}
	_, err := time.LoadLocation(name)
	return err
}

func validateTimezoneTag(name string, strict bool) error {
	if !strict {
		// In non-strict mode, skip time.LoadLocation validation
		return nil
	}
	// Remove leading '[', trailing ']', and optional leading '!' for strict validation
	cleanName := name[1 : len(name)-1]
	if cleanName[0] == '!' {
		cleanName = cleanName[1:]
	}
	if cleanName[0] == '+' || cleanName[0] == '-' {
		// Offset format, no need to validate with time.LoadLocation
		return nil
	}
	_, err := time.LoadLocation(cleanName)
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
