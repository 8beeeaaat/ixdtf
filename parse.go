package ixdtf

import (
	"errors"
	"strings"
	"time"

	"github.com/8beeeaaat/ixdtf/abnf"
)

// Parse parses an IXDTF string and returns the time and extension information.
func Parse(s string, strict bool) (time.Time, *IXDTFExtensions, error) {
	rfc3339End := findRFC3339End(s)

	t, err := parseRFC3339Portion(s[:rfc3339End])
	if err != nil {
		return time.Time{}, nil, newParseError(LayoutRFC3339, s, err)
	}

	ext, result, err := parseExtensions(s, rfc3339End, t, strict)
	if err != nil {
		return time.Time{}, nil, err
	}

	// Per RFC 9557: In non-strict mode with inconsistent timezone,
	// preserve the original timestamp and only apply timezone if consistent
	if result != nil && result.IsConsistent && result.Location != nil {
		t = t.In(result.Location)
	}
	// In non-strict mode with inconsistency, keep original timestamp as-is

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

	if _, _, err = parseExtensions(s, rfc3339End, t, strict); err != nil {
		return err
	}

	// Check the complete string against the ABNF pattern as an additional
	// validation layer on top of the structural parse above.
	if abnfErr := abnf.AbnfDateTimeExt.ValidateDateTimeExt(s); abnfErr != nil {
		return newParseError(LayoutRFC3339Extended, s, abnfErr)
	}

	return nil
}

// parseExtensions parses the optional IXDTF suffix and enforces the RFC 9557
// semantics shared by Parse and Validate: suffix grammar (Section 4.1),
// extension validation (Section 3.3), and time-zone consistency
// (Section 3.4, escalated to strict for a critical zone). The returned
// consistency result is nil when no time-zone annotation applies.
func parseExtensions(
	s string,
	rfc3339End int,
	t time.Time,
	strict bool,
) (*IXDTFExtensions, *TimezoneConsistencyResult, error) {
	var ext *IXDTFExtensions
	if rfc3339End < len(s) {
		var err error
		if ext, err = parseSuffix(s[rfc3339End:], strict); err != nil {
			return nil, nil, newParseError(LayoutRFC3339Extended, s, err)
		}
	} else {
		ext = NewIXDTFExtensions(nil)
	}

	if err := validateExtensionsStrict(ext, strict); err != nil {
		return nil, nil, newParseError(LayoutRFC3339Extended, s, err)
	}

	if ext.Location == nil {
		return ext, nil, nil
	}

	offsetUnknown := hasUnknownLocalOffset(s[:rfc3339End])
	// A critical time zone must be acted upon, so an inconsistency is an
	// error even in non-strict mode (RFC 9557 Section 3.4).
	result, err := checkTimezoneConsistency(t, ext.Location, strict || ext.CriticalLocation, offsetUnknown)
	if err != nil {
		return nil, nil, newParseError(LayoutRFC3339NanoExtended, s, err)
	}
	return ext, result, nil
}

func findRFC3339End(s string) int {
	if i := strings.IndexByte(s, '['); i >= 0 {
		return i
	}
	return len(s)
}

// parseRFC3339Portion parses the RFC 3339 date-time part. The RFC 3339 layout
// also accepts fractional seconds, so no separate nanosecond layout is needed.
func parseRFC3339Portion(rfc3339Portion string) (time.Time, error) {
	return time.Parse(time.RFC3339, rfc3339Portion)
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
