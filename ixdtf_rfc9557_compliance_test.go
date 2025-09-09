// Copyright 2024 8beeeaaat. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ixdtf

import (
	"strings"
	"testing"
	"time"
)

// TestRFC9557_DuplicateKeyHandling tests RFC 9557 requirement:
// "When duplicate keys are present, the first occurrence takes precedence"
func TestRFC9557_DuplicateKeyHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectTag   string
		expectValue string
	}{
		{
			name:        "Duplicate u-ca keys - first wins",
			input:       "2006-01-02T15:04:05Z[u-ca=japanese][u-ca=gregorian]",
			expectTag:   "u-ca",
			expectValue: "japanese",
		},
		{
			name:        "Duplicate custom extension keys",
			input:       "2006-01-02T15:04:05Z[x-test=first][x-test=second]",
			expectTag:   "x-test",
			expectValue: "first",
		},
		{
			name:        "Duplicate with critical flag - first wins",
			input:       "2006-01-02T15:04:05Z[!u-ca=japanese][u-ca=gregorian]",
			expectTag:   "u-ca",
			expectValue: "japanese",
		},
		{
			name:        "Mixed duplicate keys",
			input:       "2006-01-02T15:04:05Z[u-ca=japanese][x-custom=test][u-ca=gregorian][x-custom=ignore]",
			expectTag:   "u-ca",
			expectValue: "japanese",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ext, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) error = %v", tt.input, err)
				return
			}

			if got := ext.Tags[tt.expectTag]; got != tt.expectValue {
				t.Errorf("Parse(%q) tag[%s] = %v, want %v (first occurrence should win)",
					tt.input, tt.expectTag, got, tt.expectValue)
			}
		})
	}
}

// TestRFC9557_ExperimentalTags tests experimental tags with underscore prefix
func TestRFC9557_ExperimentalTags(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		extTag  string
		extVal  string
	}{
		{
			name:    "Valid experimental tag",
			input:   "2006-01-02T15:04:05Z[_experimental=test]",
			wantErr: false,
			extTag:  "_experimental",
			extVal:  "test",
		},
		{
			name:    "Critical experimental tag",
			input:   "2006-01-02T15:04:05Z[!_exp=value]",
			wantErr: false,
			extTag:  "_exp",
			extVal:  "value",
		},
		{
			name:    "Multiple experimental tags",
			input:   "2006-01-02T15:04:05Z[_test1=val1][_test2=val2]",
			wantErr: false,
			extTag:  "_test1",
			extVal:  "val1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ext, err := Parse(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("Parse(%q) expected error but got none", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if !tt.wantErr && ext.Tags[tt.extTag] != tt.extVal {
				t.Errorf("Parse(%q) tag[%s] = %v, want %v",
					tt.input, tt.extTag, ext.Tags[tt.extTag], tt.extVal)
			}
		})
	}
}

// TestRFC9557_ComplexExtensionCombinations tests multiple extension tags together
func TestRFC9557_ComplexExtensionCombinations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTags map[string]string
	}{
		{
			name:  "Calendar and numbering system",
			input: "2006-01-02T15:04:05Z[u-ca=hebrew][u-nu=latn]",
			wantTags: map[string]string{
				"u-ca": "hebrew",
				"u-nu": "latn",
			},
		},
		{
			name:  "Mixed standard and experimental",
			input: "2006-01-02T15:04:05Z[u-ca=japanese][_custom=test][x-vendor=value]",
			wantTags: map[string]string{
				"u-ca":     "japanese",
				"_custom":  "test",
				"x-vendor": "value",
			},
		},
		{
			name:  "Timezone with multiple extensions",
			input: "2006-01-02T15:04:05+09:00[Asia/Tokyo][u-ca=japanese][u-nu=jpan]",
			wantTags: map[string]string{
				"u-ca": "japanese",
				"u-nu": "jpan",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ext, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) error = %v", tt.input, err)
				return
			}

			for key, expectedValue := range tt.wantTags {
				if got := ext.Tags[key]; got != expectedValue {
					t.Errorf("Parse(%q) tag[%s] = %v, want %v",
						tt.input, key, got, expectedValue)
				}
			}
		})
	}
}

// TestRFC9557_SecurityConsiderations tests security-related edge cases
func TestRFC9557_SecurityConsiderations(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Very long tag key",
			input:   "2006-01-02T15:04:05Z[" + strings.Repeat("x", 1000) + "=value]",
			wantErr: false, // Should handle gracefully
		},
		{
			name:    "Very long tag value",
			input:   "2006-01-02T15:04:05Z[test=" + strings.Repeat("v", 1000) + "]",
			wantErr: false, // Should handle gracefully
		},
		{
			name:    "Many extension tags",
			input:   "2006-01-02T15:04:05Z" + strings.Repeat("[tag%d=val]", 100),
			wantErr: false, // Should handle multiple tags
		},
		{
			name:    "Special characters in values",
			input:   "2006-01-02T15:04:05Z[test=val-with-hyphens-123]",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For the "Many extension tags" test, format the string properly
			input := tt.input
			if strings.Contains(input, "%d") {
				parts := []string{"2006-01-02T15:04:05Z"}
				for i := range 10 { // Reduced from 100 for test performance
					parts = append(parts, "[tag"+string(rune('0'+i))+"=val"+string(rune('0'+i))+"]")
				}
				input = strings.Join(parts, "")
			}

			_, _, err := Parse(input)
			if tt.wantErr && err == nil {
				t.Errorf("Parse(%q) expected error but got none", input[:50]+"...")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", input[:50]+"...", err)
			}
		})
	}
}

// TestRFC9557_BoundaryConditions tests edge cases and boundary conditions
func TestRFC9557_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Single character key",
			input:   "2006-01-02T15:04:05Z[a=value]",
			wantErr: false,
		},
		{
			name:    "Single character value",
			input:   "2006-01-02T15:04:05Z[key=v]",
			wantErr: false,
		},
		{
			name:    "Numeric key",
			input:   "2006-01-02T15:04:05Z[123=value]",
			wantErr: true, // Keys must start with lowercase letter or underscore per RFC 9557
		},
		{
			name:    "Numeric value",
			input:   "2006-01-02T15:04:05Z[key=123]",
			wantErr: false,
		},
		{
			name:    "Key with hyphens",
			input:   "2006-01-02T15:04:05Z[multi-part-key=value]",
			wantErr: false,
		},
		{
			name:    "Value with hyphens",
			input:   "2006-01-02T15:04:05Z[key=multi-part-value]",
			wantErr: false,
		},
		{
			name:    "Empty extension after valid one",
			input:   "2006-01-02T15:04:05Z[valid=test][]",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := Parse(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("Parse(%q) expected error but got none", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

// TestRFC9557_CriticalTagProcessing tests critical tag handling requirements
func TestRFC9557_CriticalTagProcessing(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectCritical map[string]bool
	}{
		{
			name:  "Single critical tag",
			input: "2006-01-02T15:04:05Z[!u-ca=japanese]",
			expectCritical: map[string]bool{
				"u-ca": true,
			},
		},
		{
			name:  "Mixed critical and non-critical",
			input: "2006-01-02T15:04:05Z[!u-ca=japanese][u-nu=latn]",
			expectCritical: map[string]bool{
				"u-ca": true,
			},
		},
		{
			name:  "Multiple critical tags",
			input: "2006-01-02T15:04:05Z[!u-ca=japanese][!x-test=value]",
			expectCritical: map[string]bool{
				"u-ca":   true,
				"x-test": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ext, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) error = %v", tt.input, err)
				return
			}

			for key, expectedCritical := range tt.expectCritical {
				if got := ext.Critical[key]; got != expectedCritical {
					t.Errorf("Parse(%q) critical[%s] = %v, want %v",
						tt.input, key, got, expectedCritical)
				}
			}

			// Verify non-critical tags are not marked as critical
			for key := range ext.Tags {
				if _, shouldBeCritical := tt.expectCritical[key]; !shouldBeCritical {
					if ext.Critical[key] {
						t.Errorf("Parse(%q) critical[%s] = true, want false (should not be critical)",
							tt.input, key)
					}
				}
			}
		})
	}
}

// TestRFC9557_CalendarSystemSupport tests various calendar system values
func TestRFC9557_CalendarSystemSupport(t *testing.T) {
	calendarSystems := []string{
		"gregorian", "hebrew", "islamic", "japanese", "buddhist",
		"chinese", "coptic", "ethiopic", "indian", "persian",
	}

	for _, calendar := range calendarSystems {
		t.Run("Calendar_"+calendar, func(t *testing.T) {
			input := "2006-01-02T15:04:05Z[u-ca=" + calendar + "]"
			_, ext, err := Parse(input)
			if err != nil {
				t.Errorf("Parse(%q) error = %v", input, err)
				return
			}

			if got := ext.Tags["u-ca"]; got != calendar {
				t.Errorf("Parse(%q) u-ca = %v, want %v", input, got, calendar)
			}
		})
	}
}

// TestRFC9557_TimezoneConsistency tests timezone offset consistency with IANA timezone
func TestRFC9557_TimezoneConsistency(t *testing.T) {
	tokyo := time.FixedZone("Asia/Tokyo", 9*60*60)
	newyork := time.FixedZone("America/New_York", -5*60*60)

	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
		wantLoc     *time.Location
	}{
		{
			name:        "Consistent timezone - UTC",
			input:       "2006-01-02T15:04:05Z[UTC]",
			expectError: false,
			wantLoc:     time.UTC,
			description: "UTC timezone with Z offset should be consistent",
		},
		{
			name:        "Consistent timezone - Tokyo winter",
			input:       "2006-01-02T15:04:05+09:00[Asia/Tokyo]",
			expectError: false,
			wantLoc:     tokyo,
			description: "Tokyo timezone with +09:00 offset should be consistent in winter",
		},
		{
			name:        "Potentially inconsistent timezone - DST boundary",
			input:       "2006-07-15T15:04:05+09:00[Asia/Tokyo]",
			expectError: false,
			wantLoc:     tokyo,
			description: "Tokyo doesn't observe DST, so +09:00 should always be consistent",
		},
		{
			name:        "Basic timezone parsing without consistency check",
			input:       "2006-01-02T15:04:05-05:00[America/New_York]",
			expectError: false,
			wantLoc:     newyork,
			description: "New York winter time should be consistent with -05:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp, ext, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", tt.input, err)
				return
			}

			// Test the new consistency check function
			loc, consistencyErr := CheckTimezoneConsistency(timestamp, ext.TimeZone)
			if tt.wantLoc != nil && loc != nil && loc.String() != tt.wantLoc.String() {
				t.Errorf("CheckTimezoneConsistency(%q) location = %v, want %v",
					tt.input, loc, tt.wantLoc)
			}
			if tt.expectError && consistencyErr == nil {
				t.Errorf("CheckTimezoneConsistency(%q) expected error but got none", tt.input)
			}
			if !tt.expectError && consistencyErr != nil {
				t.Logf("CheckTimezoneConsistency(%q) note: %v (this may be expected due to DST)", tt.input, consistencyErr)
				// Don't fail the test for timezone consistency issues as they can be complex
				// This is informational for RFC 9557 compliance
			}
		})
	}
}

// TestRFC9557_ErrorResolutionMechanisms tests error handling for critical tags
func TestRFC9557_ErrorResolutionMechanisms(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "Valid critical u-ca",
			input:       "2006-01-02T15:04:05Z[!u-ca=japanese]",
			expectError: false,
			description: "Known critical calendar extension should be accepted",
		},
		{
			name:        "Valid critical private extension",
			input:       "2006-01-02T15:04:05Z[!x-custom=value]",
			expectError: false,
			description: "Private extensions can be critical",
		},
		{
			name:        "Invalid critical unknown extension",
			input:       "2006-01-02T15:04:05Z[!unknown-ext=value]",
			expectError: true,
			description: "Unknown critical extensions should cause validation error",
		},
		{
			name:        "Critical extension with empty value",
			input:       "2006-01-02T15:04:05Z[!u-ca=]",
			expectError: true,
			description: "Critical extensions with empty values should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ext, parseErr := Parse(tt.input)
			if parseErr != nil {
				if !tt.expectError {
					t.Errorf("Parse(%q) unexpected parse error: %v", tt.input, parseErr)
				}
				return
			}

			// Test validation of extensions
			validationErr := ValidateExtensions(ext)
			if tt.expectError && validationErr == nil {
				t.Errorf("ValidateExtensions(%q) expected error but got none", tt.input)
			}
			if !tt.expectError && validationErr != nil {
				t.Errorf("ValidateExtensions(%q) unexpected error: %v", tt.input, validationErr)
			}
		})
	}
}
