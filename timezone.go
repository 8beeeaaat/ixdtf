package ixdtf

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/8beeeaaat/ixdtf/abnf"
)

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
	loc, err := resolveLocation(location)
	if err != nil {
		// In non-strict mode, ignore unknown timezone errors per RFC 9557
		if !strict {
			result.IsConsistent = true // Can't check consistency for unknown timezone
			return result, nil
		}
		return result, err
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

// resolveLocation maps a location to its authoritative form. A numeric-offset
// zone name in the RFC 3339 serialization form (e.g. "+09:00", produced when
// parsing "[+09:00]") is authoritative as-is and has no timezone-database
// entry, so the location is returned unchanged. Any other name — including a
// placeholder FixedZone constructed by a caller — resolves through the
// timezone-database cache; an unknown name returns the load error.
func resolveLocation(location *time.Location) (*time.Location, error) {
	name := location.String()
	if isOffsetLocationName(name) {
		return location, nil
	}
	return loadLocationCached(name)
}

// resolveZoneAnnotation resolves the body of a time-zone annotation
// (RFC 9557 Section 4.1) to a location. An IANA name loads through the
// timezone-database cache. A numeric offset becomes a FixedZone that keeps
// the RFC 3339 serialization form ("+09:00") as the zone name so Format
// round-trips the annotation per RFC 9557 Section 1.2 and the Section 4.1
// time-numoffset grammar. An unknown or invalid name returns
// ErrInvalidTimezone; whether that is fatal is the caller's decision.
func resolveZoneAnnotation(name string) (*time.Location, error) {
	if loc, ok := tryLoadTimezone(name); ok {
		return loc, nil
	}
	if offset, err := parseNumericOffset(name); err == nil {
		return time.FixedZone(name, offset), nil
	}
	return nil, ErrInvalidTimezone
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
// and returns the offset in seconds. The shape of the string is defined by
// isOffsetLocationName; this function adds the hour and minute range checks.
func parseNumericOffset(s string) (int, error) {
	if !isOffsetLocationName(s) {
		return 0, errors.New("invalid offset format")
	}

	hours, err := strconv.Atoi(s[1:3])
	if err != nil || hours > 23 {
		return 0, errors.New("invalid hours")
	}

	minutes, err := strconv.Atoi(s[4:6])
	if err != nil || minutes > 59 {
		return 0, errors.New("invalid minutes")
	}

	sign := 1
	if s[0] == '-' {
		sign = -1
	}

	return sign * (hours*3600 + minutes*60), nil
}

// isOffsetLocationName reports whether name is a numeric-offset zone name in
// the RFC 3339 serialization form used for offset time-zone annotations
// (e.g. "+09:00", "-03:30"), as produced when parsing "[+09:00]". This is the
// single definition of the offset-name shape; parseNumericOffset builds on it.
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
