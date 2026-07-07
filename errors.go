package ixdtf

import (
	"errors"

	"github.com/8beeeaaat/ixdtf/abnf"
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

func newParseError(layout Layout, value string, err error) error {
	return &ParseError{
		Layout: layout,
		Value:  value,
		Err:    err,
	}
}
