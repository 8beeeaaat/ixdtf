package ixdtf_test

import (
	"testing"

	"github.com/8beeeaaat/ixdtf"
)

func TestValidateUnicodeCalendarIdentifiers(t *testing.T) {
	t.Parallel()

	validIdentifiers := []string{
		"buddhist",
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
		"islamicc",
	}

	for _, identifier := range validIdentifiers {
		t.Run(identifier, func(t *testing.T) {
			t.Parallel()
			input := "2025-03-04T05:06:07Z[!u-ca=" + identifier + "]"
			if err := ixdtf.Validate(input, false); err != nil {
				t.Fatalf("Validate(%q) unexpected error: %v", input, err)
			}
		})
	}
}
