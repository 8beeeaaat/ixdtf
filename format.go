package ixdtf

import (
	"sort"
	"time"

	"github.com/8beeeaaat/ixdtf/abnf"
)

// Format formats a time with IXDTF extensions using RFC 3339 format.
// The time-zone annotation is emitted with a leading "!" when
// ext.CriticalLocation is set.
func Format(t time.Time, ext *IXDTFExtensions) (string, error) {
	return format(t, ext, time.RFC3339)
}

// FormatNano formats a time with IXDTF extensions using RFC 3339 format with nanoseconds.
// The time-zone annotation is emitted with a leading "!" when
// ext.CriticalLocation is set.
func FormatNano(t time.Time, ext *IXDTFExtensions) (string, error) {
	return format(t, ext, time.RFC3339Nano)
}

// format validates the extensions and serializes the timestamp with its IXDTF
// suffix. Formatting always validates strictly: the producer of a string must
// only emit annotations it can process (RFC 9557 Section 3.3).
func format(t time.Time, ext *IXDTFExtensions, layout string) (string, error) {
	if err := validateExtensionsStrict(ext, true); err != nil {
		return "", err
	}
	if err := validateCriticalLocation(t, ext); err != nil {
		return "", err
	}
	return string(appendSuffix(t, ext, layout)), nil
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
	b := t.AppendFormat(nil, format)

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
