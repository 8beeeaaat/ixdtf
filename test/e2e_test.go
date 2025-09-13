package e2e_test

import (
	"testing"

	"github.com/8beeeaaat/ixdtf"
)

func TestE2E_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "basic UTC time",
			input: "2025-01-02T03:04:05Z[UTC]",
		},
		{
			name:  "timezone with offset",
			input: "2025-02-03T04:05:06+09:00[Asia/Tokyo]",
		},
		{
			name:  "with tags",
			input: "2025-03-04T05:06:07Z[UTC][u-ca=gregory]",
		},
		{
			name:  "with critical tags",
			input: "2025-02-03T04:05:06+09:00[Asia/Tokyo][!u-ca=gregory]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := ixdtf.Validate(tc.input, true); err != nil {
				t.Fatalf("validation failed for %q: %v", tc.input, err)
			}

			parsedTime, ext, err := ixdtf.Parse(tc.input, true)
			if err != nil {
				t.Fatalf("failed to parse %q: %v", tc.input, err)
			}

			formatted, err := ixdtf.Format(parsedTime, ext)
			if err != nil {
				t.Fatalf("failed to format: %v", err)
			}

			if formatted != tc.input {
				t.Fatalf("round trip failed: input %q, got %q", tc.input, formatted)
			}
		})
	}
}

func TestE2E_TimezoneConversion(t *testing.T) {
	// Test that times in different timezones represent the same instant
	utcTime := "2025-01-01T12:00:00Z[UTC]"
	tokyoTime := "2025-01-01T21:00:00+09:00[Asia/Tokyo]"

	parsedUTC, _, err := ixdtf.Parse(utcTime, true)
	if err != nil {
		t.Fatalf("failed to parse UTC time: %v", err)
	}

	parsedTokyo, _, err := ixdtf.Parse(tokyoTime, true)
	if err != nil {
		t.Fatalf("failed to parse Tokyo time: %v", err)
	}

	if !parsedUTC.Equal(parsedTokyo) {
		t.Fatalf("times should be equal: UTC=%v, Tokyo=%v", parsedUTC, parsedTokyo)
	}
}
