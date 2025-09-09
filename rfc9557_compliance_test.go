// Copyright 2024 8beeeaaat. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ixdtf

import (
	"testing"
)

// TestRFC9557ABNFCompliance tests compliance with RFC 9557 ABNF grammar
func TestRFC9557ABNFCompliance(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		reason    string
	}{
		// Valid suffix keys (key-initial = lcalpha / "_", key-char = key-initial / DIGIT / "-")
		{
			name:  "Valid key starting with lowercase letter",
			input: "2006-01-02T15:04:05Z[u-ca=japanese]",
		},
		{
			name:  "Valid key starting with underscore",
			input: "2006-01-02T15:04:05Z[_custom=value]",
		},
		{
			name:  "Valid key with digits and hyphens",
			input: "2006-01-02T15:04:05Z[x-test123-ext=value]",
		},
		
		// Invalid suffix keys
		{
			name:      "Invalid key starting with uppercase",
			input:     "2006-01-02T15:04:05Z[U-ca=japanese]",
			wantError: true,
			reason:    "key-initial must be lcalpha or underscore",
		},
		{
			name:      "Invalid key starting with digit",
			input:     "2006-01-02T15:04:05Z[1-test=value]",
			wantError: true,
			reason:    "key-initial cannot be digit",
		},
		{
			name:      "Invalid key with special characters",
			input:     "2006-01-02T15:04:05Z[u@ca=japanese]",
			wantError: true,
			reason:    "key-char limited to lcalpha/digit/hyphen/underscore",
		},
		
		// Valid suffix values (suffix-value = 1*alphanum, suffix-values = suffix-value *("-" suffix-value))
		{
			name:  "Valid simple alphanumeric value",
			input: "2006-01-02T15:04:05Z[u-ca=japanese123]",
		},
		{
			name:  "Valid hyphenated values",
			input: "2006-01-02T15:04:05Z[u-ca=ja-JP-u-ca-japanese]",
		},
		{
			name:  "Valid mixed case value",
			input: "2006-01-02T15:04:05Z[x-test=TestValue123]",
		},
		
		// Invalid suffix values  
		{
			name:      "Invalid value with underscore",
			input:     "2006-01-02T15:04:05Z[u-ca=japanese_variant]",
			wantError: true,
			reason:    "suffix-value-char cannot include underscore",
		},
		{
			name:      "Invalid value with special characters",
			input:     "2006-01-02T15:04:05Z[u-ca=japanese@variant]",
			wantError: true,
			reason:    "suffix-value-char limited to alphanum",
		},
		{
			name:      "Invalid empty value part",
			input:     "2006-01-02T15:04:05Z[u-ca=japanese--variant]",
			wantError: true,
			reason:    "empty value parts not allowed in hyphenated values",
		},
		
		// Critical extensions
		{
			name:  "Valid critical extension",
			input: "2006-01-02T15:04:05Z[!u-ca=japanese]",
		},
		{
			name:      "Invalid critical timezone",
			input:     "2006-01-02T15:04:05Z[!UTC]",
			wantError: true,
			reason:    "timezones cannot be critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := Parse(tt.input)
			if tt.wantError && err == nil {
				t.Errorf("Parse(%q) expected error (%s) but got none", tt.input, tt.reason)
			}
			if !tt.wantError && err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}


// TestRFC9557ExampleFormats tests example formats from RFC 9557
func TestRFC9557ExampleFormats(t *testing.T) {
	validExamples := []string{
		// Basic RFC 3339
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+09:00",
		
		// With timezone suffix
		"2006-01-02T15:04:05Z[UTC]",
		"2006-01-02T15:04:05+09:00[Asia/Tokyo]",
		
		// With extensions
		"2006-01-02T15:04:05Z[u-ca=japanese]",
		"2006-01-02T15:04:05Z[UTC][u-ca=japanese]",
		
		// With critical extensions
		"2006-01-02T15:04:05Z[!u-ca=japanese]",
		"2006-01-02T15:04:05Z[Asia/Tokyo][!u-ca=japanese][u-nu=latn]",
	}

	for _, example := range validExamples {
		t.Run(example, func(t *testing.T) {
			_, _, err := Parse(example)
			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", example, err)
			}
		})
	}
}