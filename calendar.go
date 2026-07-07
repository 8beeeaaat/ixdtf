package ixdtf

import "strings"

// ExtensionUnicodeCalendar is tag key for Unicode calendar extension.
// https://www.rfc-editor.org/rfc/rfc9557.html#section-5
const ExtensionUnicodeCalendar = "u-ca"

// validateTagValue enforces value rules for registered suffix keys
// (RFC 9557 Section 5). Unregistered keys have no value constraints.
func validateTagValue(key, value string) error {
	if key == ExtensionUnicodeCalendar && !isUnicodeCalendarIdentifier(value) {
		return ErrInvalidTagCalendarIdentifier
	}
	return nil
}

func isUnicodeCalendarIdentifier(value string) bool {
	if value == "" {
		return false
	}
	switch strings.ToLower(value) {
	case "buddhist",
		"chinese",
		"coptic",
		"dangi",
		"ethioaa",
		"ethiopic-amete-alem",
		"ethiopic",
		"gregory",
		"gregorian",
		"hebrew",
		"indian",
		"islamic",
		"islamic-umalqura",
		"islamic-tbla",
		"islamic-civil",
		"islamic-rgsa",
		"iso8601",
		"japanese",
		"persian",
		"roc",
		"islamicc":
		return true
	default:
		return false
	}
}
