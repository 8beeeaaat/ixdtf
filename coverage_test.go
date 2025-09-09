// Copyright 2024 8beeeaaat. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ixdtf

import (
	"errors"
	"testing"
)

// TestErrorTypes tests specific error type handling
func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   string
		checkType bool
	}{
		{
			name:      "Invalid suffix format",
			input:     "2006-01-02T15:04:05Z[",
			wantErr:   "parsing time \"2006-01-02T15:04:05Z[\" as \"2006-01-02T15:04:05Z07:00[time-zone]\": invalid IXDTF suffix format",
			checkType: true,
		},
		{
			name:      "Invalid extension format",
			input:     "2006-01-02T15:04:05Z[key=]",
			wantErr:   "parsing time \"2006-01-02T15:04:05Z[key=]\" as \"2006-01-02T15:04:05Z07:00[time-zone]\": invalid extension format",
			checkType: true,
		},
		{
			name:      "Critical timezone",
			input:     "2006-01-02T15:04:05Z[!UTC]",
			wantErr:   "parsing time \"2006-01-02T15:04:05Z[!UTC]\" as \"2006-01-02T15:04:05Z07:00[time-zone]\": invalid timezone name",
			checkType: true,
		},
		{
			name:      "Multiple timezones",
			input:     "2006-01-02T15:04:05Z[UTC][GMT]",
			wantErr:   "parsing time \"2006-01-02T15:04:05Z[UTC][GMT]\" as \"2006-01-02T15:04:05Z07:00[time-zone]\": invalid timezone name",
			checkType: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%q) expected error but got none", tt.input)
				return
			}

			// Check if the error is wrapped in a ParseError
			var parseErr *ParseError
			if errors.As(err, &parseErr) {
					expectedMsg := tt.wantErr
					if parseErr.Error() != expectedMsg {
						t.Errorf("Parse(%q) error message = %q, want %q", tt.input, parseErr.Error(), expectedMsg)
					}
			} else {
				t.Errorf("Parse(%q) expected ParseError, got %T", tt.input, err)
			}
		})
	}
}

// TestParseRFC3339Fallback tests RFC3339 parsing fallback behavior
func TestParseRFC3339Fallback(t *testing.T) {
	// Test a time format that RFC3339Nano can't parse but RFC3339 can
	tests := []struct {
		name  string
		input string
	}{
		// These should test the fallback from RFC3339Nano to RFC3339
		{"basic rfc3339", "2006-01-02T15:04:05Z"},
		{"with offset", "2006-01-02T15:04:05+09:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

// TestParseInvalidRFC3339 tests parsing of invalid RFC 3339 formats
func TestParseInvalidRFC3339(t *testing.T) {
	tests := []struct {
		name  string
		input string
		wantError string
	}{
		// These should fail both RFC3339Nano and RFC3339 parsing
		{"completely invalid", "invalid-time-format","parsing time \"invalid-time-format\" as \"2006-01-02T15:04:05Z07:00\": parsing time \"invalid-time-format\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid-time-format\" as \"2006\"", },
		{"invalid month", "2006-13-02T15:04:05Z","parsing time \"2006-13-02T15:04:05Z\" as \"2006-01-02T15:04:05Z07:00\": parsing time \"2006-13-02T15:04:05Z\": month out of range",},
		{"invalid hour", "2006-01-02T25:04:05Z","parsing time \"2006-01-02T25:04:05Z\" as \"2006-01-02T15:04:05Z07:00\": parsing time \"2006-01-02T25:04:05Z\": hour out of range",},
		{"malformed", "not-a-date-at-all","parsing time \"not-a-date-at-all\" as \"2006-01-02T15:04:05Z07:00\": parsing time \"not-a-date-at-all\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"not-a-date-at-all\" as \"2006\"",},
		{"partial date", "2006-01-02","parsing time \"2006-01-02\" as \"2006-01-02T15:04:05Z07:00\": parsing time \"2006-01-02\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"T\"" },
		{"missing T separator", "2006-01-02 15:04:05Z","parsing time \"2006-01-02 15:04:05Z\" as \"2006-01-02T15:04:05Z07:00\": parsing time \"2006-01-02 15:04:05Z\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \" 15:04:05Z\" as \"T\"",},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%q) expected error but got none", tt.input)
				return
			}

			// Check if it's the specific RFC3339 error we're looking for
			var parseErr *ParseError
			if !errors.As(err, &parseErr) {
				t.Errorf("Parse(%q) expected ParseError, got %T", tt.input, err)
				return
			}

			if err.Error() != tt.wantError {
				t.Errorf("Parse(%q) error = %q, want %q", tt.input, err.Error(), tt.wantError)
			}
		})
	}
}
