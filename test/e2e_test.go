package e2e_test

import (
	"testing"
	"time"

	"github.com/8beeeaaat/ixdtf"
)

func TestE2E_RoundTrip(t *testing.T) {
	tests := []testRoundTripArgs{
		{
			name:                                     "basic UTC time",
			input:                                    "2025-01-02T03:04:05Z",
			inputHasErrorInNonStrictMode:             false,
			inputHasErrorInStrictMode:                false,
			overrideWithNyExt:                        "2025-01-02T03:04:05Z[America/New_York]",
			overrideWithNyExtHasErrorInNonStrictMode: false,
			overrideWithNyExtHasErrorInStrictMode:    true,
		},
		{
			name:                         "basic UTC time with invalid critical time zone tag",
			input:                        "2025-01-02T03:04:05Z[!America/New_York]",
			inputHasErrorInNonStrictMode: true,
			inputHasErrorInStrictMode:    true,
		},
		{
			name:                         "invalid offset time with critical time zone tag",
			input:                        "2025-01-02T03:04:05+09:00[!America/New_York]",
			inputHasErrorInNonStrictMode: true,
			inputHasErrorInStrictMode:    true,
		},
		{
			name:                                     "with unknown time zone tag",
			input:                                    "2025-01-02T03:04:05+09:00[Foo/Bar]",
			inputHasErrorInNonStrictMode:             false,
			inputHasErrorInStrictMode:                true,
			inputFormatted:                           "2025-01-02T03:04:05+09:00",
			overrideWithNyExt:                        "2025-01-02T03:04:05+09:00[America/New_York]",
			overrideWithNyExtHasErrorInNonStrictMode: false,
			overrideWithNyExtHasErrorInStrictMode:    true,
		},
		{
			name:                                     "offset time",
			input:                                    "2025-02-03T04:05:06+09:00",
			inputHasErrorInNonStrictMode:             false,
			inputHasErrorInStrictMode:                false,
			overrideWithNyExt:                        "2025-02-03T04:05:06+09:00[America/New_York]",
			overrideWithNyExtHasErrorInNonStrictMode: false,
			overrideWithNyExtHasErrorInStrictMode:    true,
		},
		{
			name:                                     "timezone with offset - New York",
			input:                                    "2025-02-03T04:05:06-05:00[America/New_York]",
			inputHasErrorInNonStrictMode:             false,
			inputHasErrorInStrictMode:                false,
			overrideWithNyExt:                        "2025-02-03T04:05:06-05:00[America/New_York]",
			overrideWithNyExtHasErrorInNonStrictMode: false,
			overrideWithNyExtHasErrorInStrictMode:    false,
		},
		{
			name:                                     "timezone with offset - Tokyo",
			input:                                    "2025-02-03T04:05:06+09:00[Asia/Tokyo]",
			inputHasErrorInNonStrictMode:             false,
			inputHasErrorInStrictMode:                false,
			overrideWithNyExt:                        "2025-02-03T04:05:06+09:00[America/New_York]",
			overrideWithNyExtHasErrorInNonStrictMode: false,
			overrideWithNyExtHasErrorInStrictMode:    true,
		},
		{
			name:                                     "with tags",
			input:                                    "2025-03-04T05:06:07Z[UTC][u-ca=gregory]",
			inputHasErrorInNonStrictMode:             false,
			inputHasErrorInStrictMode:                false,
			overrideWithNyExt:                        "2025-03-04T05:06:07Z[America/New_York][u-ca=gregory]",
			overrideWithNyExtHasErrorInNonStrictMode: false,
			overrideWithNyExtHasErrorInStrictMode:    true,
		},
		{
			name:                                     "with critical tags",
			input:                                    "2025-02-03T04:05:06+09:00[Asia/Tokyo][!u-ca=gregory]",
			inputHasErrorInNonStrictMode:             false,
			inputHasErrorInStrictMode:                false,
			overrideWithNyExt:                        "2025-02-03T04:05:06+09:00[America/New_York][!u-ca=gregory]",
			overrideWithNyExtHasErrorInNonStrictMode: false,
			overrideWithNyExtHasErrorInStrictMode:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Validate input in both modes
			validateInput(t, tc)

			if tc.inputHasErrorInNonStrictMode {
				return // Skip further tests if input is expected to fail
			}

			parsedTime, ext := parseAndValidateRoundTrip(t, tc)
			testNewYorkConversion(t, tc, parsedTime, ext)
		})
	}
}

func TestE2E_TimeZoneConversion(t *testing.T) {
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

func TestE2E_TimeZoneInconsistency(t *testing.T) {
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

type testRoundTripArgs struct {
	name                                     string
	input                                    string
	inputHasErrorInNonStrictMode             bool
	inputHasErrorInStrictMode                bool
	inputFormatted                           string
	overrideWithNyExt                        string
	overrideWithNyExtHasErrorInNonStrictMode bool
	overrideWithNyExtHasErrorInStrictMode    bool
}

// validateInput validates the input string in both strict and non-strict modes.
func validateInput(t *testing.T, tc testRoundTripArgs) {
	t.Helper()

	// Test non-strict mode
	err := ixdtf.Validate(tc.input, false)
	if err != nil && !tc.inputHasErrorInNonStrictMode {
		t.Fatalf("validation failed for %q: %v in non-strict mode", tc.input, err)
	}
	if err == nil && tc.inputHasErrorInNonStrictMode {
		t.Fatalf("expected validation error for %q in non-strict mode, but got none", tc.input)
	}

	// Test strict mode
	err = ixdtf.Validate(tc.input, true)
	if err != nil && !tc.inputHasErrorInStrictMode {
		t.Fatalf("validation failed for %q: %v in strict mode", tc.input, err)
	}
	if err == nil && tc.inputHasErrorInStrictMode {
		t.Fatalf("expected validation error for %q in strict mode, but got none", tc.input)
	}
}

// parseAndValidateRoundTrip parses input and validates round-trip formatting.
func parseAndValidateRoundTrip(t *testing.T, tc testRoundTripArgs) (time.Time, *ixdtf.IXDTFExtensions) {
	t.Helper()

	parsedTime, ext, err := ixdtf.Parse(tc.input, false)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", tc.input, err)
	}

	formatted, err := ixdtf.Format(parsedTime, ext)
	if err != nil {
		t.Fatalf("failed to format: %v", err)
	}

	if tc.inputFormatted != "" && formatted != tc.inputFormatted {
		t.Fatalf("round trip failed: inputFormatted %q, formatted %q", tc.inputFormatted, formatted)
	}

	if tc.inputFormatted == "" && formatted != tc.input {
		t.Fatalf("round trip failed: should be same as input %q, formatted %q", tc.input, formatted)
	}

	return parsedTime, ext
}

// testNewYorkConversion tests timezone conversion to New York and validates the result.
func testNewYorkConversion(t *testing.T, tc testRoundTripArgs, parsedTime time.Time, ext *ixdtf.IXDTFExtensions) {
	t.Helper()

	// Convert to New York timezone and format
	nyLoc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("failed to load New York location: %v", err)
	}
	nyExt := ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
		Location: nyLoc,
		Tags:     ext.Tags,
		Critical: ext.Critical,
	})

	overrideWithNyExt, err := ixdtf.Format(parsedTime, nyExt)
	if err != nil {
		t.Fatalf("failed to format in New York timezone: %v", err)
	}
	if overrideWithNyExt != tc.overrideWithNyExt {
		t.Fatalf("New York format mismatch: got %q, want %q", overrideWithNyExt, tc.overrideWithNyExt)
	}

	// Validate the New York formatted string in non-strict mode
	err = ixdtf.Validate(overrideWithNyExt, false)
	if err != nil && !tc.overrideWithNyExtHasErrorInNonStrictMode {
		t.Fatalf("validation failed for New York format %q: %v in non-strict mode", overrideWithNyExt, err)
	}
	if err == nil && tc.overrideWithNyExtHasErrorInNonStrictMode {
		t.Fatalf(
			"expected validation error for New York format %q in non-strict mode, but got none",
			overrideWithNyExt,
		)
	}

	// Validate the New York formatted string in strict mode
	err = ixdtf.Validate(overrideWithNyExt, true)
	if err != nil && !tc.overrideWithNyExtHasErrorInStrictMode {
		t.Fatalf("unexpected validation error in strict mode for %q: %v", overrideWithNyExt, err)
	}
	if err == nil && tc.overrideWithNyExtHasErrorInStrictMode {
		t.Fatalf("expected validation error in strict mode for %q, but got none", overrideWithNyExt)
	}

	// Expected pass, with located time
	nyTime := parsedTime.In(nyLoc)
	overrideWithNyExtWithLocatedTime, err := ixdtf.Format(nyTime, nyExt)
	if err != nil {
		t.Fatalf("failed to format with located time: %v", err)
	}
	if err = ixdtf.Validate(overrideWithNyExtWithLocatedTime, true); err != nil {
		t.Fatalf("failed to validate New York format with located time: %v", err)
	}
}
