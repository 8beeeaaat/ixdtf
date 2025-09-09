// Copyright 2024 8beeeaaat. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ixdtf

import (
	"testing"
)

// TestIsValidSuffixValueEdgeCases tests specific edge cases for isValidSuffixValue
func TestIsValidSuffixValueEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		// Test length 1 strings
		{"single hyphen", "-", false},
		{"single valid char", "a", true},
		{"single digit", "1", true},

		// Test boundary conditions for the checks
		{"two hyphens", "--", false},
		{"hyphen at start", "-a", false},
		{"hyphen at end", "a-", false},
		{"valid with hyphen", "a-b", true},

		// Test consecutive hyphen detection (different positions)
		{"consecutive in middle", "abc--def", false},
		{"triple hyphen", "a---b", false},

		// Test edge case with length 2 starting with hyphen
		{"length 2 start hyphen", "-x", false},

		// Test the loop boundary for consecutive hyphens
		{"exactly length 2 valid", "ab", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidSuffixValue(tt.value)
			if got != tt.valid {
				t.Errorf("isValidSuffixValue(%q) = %v, want %v", tt.value, got, tt.valid)
			}
		})
	}
}

// TestIsValidSuffixValueCoverage tests edge cases to reach 100% coverage
func TestIsValidSuffixValueCoverage(t *testing.T) {
	// Test the boundary conditions that might not be covered
	tests := []struct {
		value string
		valid bool
	}{
		{"-", false}, // Single hyphen
	}

	for _, test := range tests {
		got := isValidSuffixValue(test.value)
		if got != test.valid {
			t.Errorf("isValidSuffixValue(%q) = %v, want %v", test.value, got, test.valid)
		}
	}
}

// TestSplitOnHyphenEdgeCases tests edge cases for splitOnHyphen
func TestSplitOnHyphenEdgeCases(t *testing.T) {
	// Test with string that starts with hyphen
	parts := splitOnHyphen("-test")
	if len(parts) != 1 || parts[0] != "test" {
		t.Errorf("splitOnHyphen(\"-test\") = %v, want [\"test\"]", parts)
	}

	// Test with string that ends with hyphen
	parts = splitOnHyphen("test-")
	if len(parts) != 1 || parts[0] != "test" {
		t.Errorf("splitOnHyphen(\"test-\") = %v, want [\"test\"]", parts)
	}
}

// TestParseSuffixEdgeCases tests specific edge cases for parseSuffix
func TestParseSuffixEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		suffix  string
		wantErr bool
	}{
		{"starts with non-bracket", "invalid", true},
		{"unclosed bracket", "[test", true},
		{"empty bracket", "[]", true},
		{"valid timezone", "[UTC]", false},
		{"valid extension", "[u-ca=japanese]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseSuffix(tt.suffix)
			if tt.wantErr && err == nil {
				t.Errorf("parseSuffix(%q) expected error but got none", tt.suffix)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("parseSuffix(%q) unexpected error: %v", tt.suffix, err)
			}
		})
	}
}

// TestAppendSuffixEdgeCases tests edge cases in appendSuffix function
func TestAppendSuffixEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		ext  IXDTFExtensions
		want string
	}{
		{
			name: "Multiple tags in sorted order",
			ext: IXDTFExtensions{
				TimeZone: "UTC",
				Tags:     map[string]string{"z-test": "value", "a-test": "value", "m-test": "value"},
				Critical: map[string]bool{"a-test": true},
			},
			want: "[UTC][!a-test=value][m-test=value][z-test=value]",
		},
		{
			name: "Empty timezone with tags",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"test": "value"},
				Critical: map[string]bool{},
			},
			want: "[test=value]",
		},
		{
			name: "Only timezone",
			ext: IXDTFExtensions{
				TimeZone: "America/New_York",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
			want: "[America/New_York]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := []byte("2006-01-02T15:04:05Z")
			result := appendSuffix(b, tt.ext)
			suffix := string(result[len(b):])
			if suffix != tt.want {
				t.Errorf("appendSuffix() suffix = %v, want %v", suffix, tt.want)
			}
		})
	}
}

// TestAppendSuffixCoverageEdgeCases tests specific edge cases for appendSuffix
func TestAppendSuffixCoverageEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		ext  IXDTFExtensions
		want string
	}{
		{
			name: "empty extensions",
			ext:  NewIXDTFExtensions(),
			want: "",
		},
		{
			name: "only timezone no tags",
			ext: IXDTFExtensions{
				TimeZone: "UTC",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
			want: "[UTC]",
		},
		{
			name: "only tags no timezone",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{},
			},
			want: "[u-ca=japanese]",
		},
		{
			name: "empty tags map but has critical",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{},
				Critical: map[string]bool{"test": true}, // This should not appear since tag doesn't exist
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := []byte("2006-01-02T15:04:05Z")
			result := appendSuffix(base, tt.ext)
			suffix := string(result[len(base):])
			if suffix != tt.want {
				t.Errorf("appendSuffix() suffix = %v, want %v", suffix, tt.want)
			}
		})
	}
}

// TestIsValidSuffixKey tests private suffix key validation function
func TestIsValidSuffixKey(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		valid bool
	}{
		// Valid keys
		{"lowercase start", "u-ca", true},
		{"underscore start", "_custom", true},
		{"with digits", "x-test123", true},
		{"with hyphens", "u-ca-variant", true},
		{"underscore and hyphen", "_test-123", true},
		
		// Invalid keys
		{"empty", "", false},
		{"uppercase start", "U-ca", false},
		{"digit start", "1-test", false},
		{"special char start", "@test", false},
		{"special char in middle", "u@ca", false},
		{"space", "u ca", false},
		{"dot", "u.ca", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidSuffixKey(tt.key)
			if got != tt.valid {
				t.Errorf("isValidSuffixKey(%q) = %v, want %v", tt.key, got, tt.valid)
			}
		})
	}
}

// TestSplitOnHyphen tests private string splitting function
func TestSplitOnHyphen(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty string", "", []string{}},
		{"no hyphens", "test", []string{"test"}},
		{"single hyphen", "a-b", []string{"a", "b"}},
		{"multiple hyphens", "a-b-c-d", []string{"a", "b", "c", "d"}},
		{"complex", "ja-JP-u-ca-japanese", []string{"ja", "JP", "u", "ca", "japanese"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitOnHyphen(tt.input)
			if len(got) != len(tt.expected) {
				t.Errorf("splitOnHyphen(%q) length = %v, want %v", tt.input, len(got), len(tt.expected))
				return
			}
			for i, part := range got {
				if part != tt.expected[i] {
					t.Errorf("splitOnHyphen(%q)[%d] = %v, want %v", tt.input, i, part, tt.expected[i])
				}
			}
		})
	}
}

// TestValidateCriticalExtension tests private critical extension validation
func TestValidateCriticalExtension(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{"Valid u-ca", "u-ca", "japanese", false},
		{"Valid u-nu", "u-nu", "latn", false},
		{"Valid private", "x-custom", "value", false},
		{"Empty u-ca value", "u-ca", "", true},
		{"Empty u-nu value", "u-nu", "", true},
		{"Empty private value", "x-private", "", true},
		{"Unknown extension", "unknown", "value", true},
		{"Unknown u- extension", "u-unknown", "value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCriticalExtension(tt.key, tt.value)
			if tt.wantErr && err == nil {
				t.Errorf("validateCriticalExtension(%q, %q) expected error but got none", tt.key, tt.value)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("validateCriticalExtension(%q, %q) unexpected error: %v", tt.key, tt.value, err)
			}
		})
	}
}

// TestIsValidTimezone covers all branches of isValidTimezone
func TestIsValidTimezone(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"length < 2", "Z", true},          // len(s) < 2 -> true
		{"no hyphen at index 1", "Etc/GMT", true}, // s[1] != '-' -> true
		{"u- prefix (extension-like)", "u-test", false},
		{"x- prefix (extension-like)", "x-foo", false},
		{"t- prefix (extension-like)", "t-bar", false},
		{"other prefix with hyphen", "a-test", true}, // s[1]=='-' but prefix not in set -> true
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidTimezone(tt.s)
			if got != tt.want {
				t.Errorf("isValidTimezone(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

// TestParse_TimezoneWithHyphen ensures timezone names containing hyphens are accepted (e.g., Etc/GMT-1)
func TestParse_TimezoneWithHyphen(t *testing.T) {
	input := "2006-01-02T15:04:05Z[Etc/GMT-1]"
	_, ext, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", input, err)
	}
	if ext.TimeZone != "Etc/GMT-1" {
		t.Errorf("timezone = %q, want %q", ext.TimeZone, "Etc/GMT-1")
	}
}

// TestIsValidSuffixValueMore adds missing edge cases for isValidSuffixValue
func TestIsValidSuffixValueMore(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"", false},                  // len==0
		{"Abc-DEF-123", true},        // mixed case alphanum parts
		{"a@b", false},               // non-alphanum character
	}
	for _, tt := range tests {
		got := isValidSuffixValue(tt.value)
		if got != tt.want {
			t.Errorf("isValidSuffixValue(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}