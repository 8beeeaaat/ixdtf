package ixdtf

import (
	"strings"

	"github.com/8beeeaaat/ixdtf/abnf"
)

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
	if eq := strings.IndexByte(content[startIdx:], '='); eq >= 0 {
		state.seenTag = true
		return handleExtensionTag(content, critical, startIdx, startIdx+eq, ext, strict)
	}

	// Time-zone annotation.
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

	loc, err := resolveZoneAnnotation(content[startIdx:])
	if err != nil {
		// RFC 9557 Section 4.1 permits a critical flag ("!") on a time-zone
		// annotation, e.g. "[!Europe/London]" (Figures 1 and 2 in Section
		// 3.4). A critical annotation MUST be processable (Section 3.3), so
		// an unknown or invalid name is rejected even in non-strict mode;
		// otherwise a non-strict parse ignores the annotation per RFC 9557.
		if strict || critical {
			return err
		}
		return nil
	}
	ext.Location = loc
	ext.CriticalLocation = critical
	return nil
}

// handleExtensionTag processes an extension tag element (key=value pair).
// equalIndex is the position of '=' within content, at or after startIdx.
func handleExtensionTag(
	content string,
	critical bool,
	startIdx, equalIndex int,
	ext *IXDTFExtensions,
	strict bool,
) error {
	if equalIndex == startIdx || equalIndex == len(content)-1 {
		return ErrInvalidExtension // empty key or value
	}
	if err := abnf.AbnfSuffixKey.ValidateSuffixKey(content[startIdx:equalIndex]); err != nil {
		return err
	}
	if err := isValidSuffixValue(content[equalIndex+1:]); err != nil {
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
