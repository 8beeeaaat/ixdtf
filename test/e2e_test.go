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
			name:  "offset time",
			input: "2025-02-03T04:05:06+09:00",
		},
		{
			name:  "timezone with offset - New York",
			input: "2025-02-03T04:05:06-05:00[America/New_York]",
		},
		{
			name:  "timezone with offset - Tokyo",
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
				t.Fatalf("round trip failed: input %q, formatted %q", tc.input, formatted)
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

func TestE2E_TimezoneInconsistency(t *testing.T) {
	// Test RFC 9557 compliant handling of timezone offset mismatches
	// Input has +09:00 offset but [America/New_York] timezone (which should be -4/-5 hours)
	input := "2025-06-01T12:00:00+09:00[America/New_York]"

	// In non-strict mode, preserve original timestamp and detect inconsistency
	parsedTime, ext, err := ixdtf.Parse(input, false)
	if err != nil {
		t.Fatalf("failed to parse inconsistent timezone in non-strict mode: %v", err)
	}

	// Should preserve original +09:00 offset as per RFC 9557
	expectedTime := "2025-06-01T12:00:00+09:00"
	if parsedTime.Format("2006-01-02T15:04:05Z07:00") != expectedTime {
		t.Fatalf("non-strict mode should preserve original timestamp: got %s, want %s",
			parsedTime.Format("2006-01-02T15:04:05Z07:00"), expectedTime)
	}

	// Extension should still contain the parsed timezone information
	if ext.Location == nil {
		t.Fatalf("extension should contain location information")
	}
	if ext.Location.String() != "America/New_York" {
		t.Fatalf("extension location should be America/New_York, got %s", ext.Location.String())
	}

	// In strict mode, should return error
	_, _, err = ixdtf.Parse(input, true)
	if err == nil {
		t.Fatalf("strict mode should return error for timezone offset mismatch")
	}
}
