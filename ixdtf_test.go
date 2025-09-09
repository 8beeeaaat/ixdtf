// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ixdtf

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTime string
		wantExt  IXDTFExtensions
	}{
		{
			name:     "RFC 3339 without suffixes",
			input:    "2006-01-02T15:04:05Z",
			wantTime: "2006-01-02T15:04:05Z",
			wantExt:  NewIXDTFExtensions(),
		},
		{
			name:     "RFC 3339 with timezone suffix",
			input:    "2006-01-02T15:04:05Z[UTC]",
			wantTime: "2006-01-02T15:04:05Z",
			wantExt: IXDTFExtensions{
				TimeZone: "UTC",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
		},
		{
			name:     "RFC 3339 with offset and timezone suffix",
			input:    "2006-01-02T15:04:05+09:00[Asia/Tokyo]",
			wantTime: "2006-01-02T15:04:05+09:00",
			wantExt: IXDTFExtensions{
				TimeZone: "Asia/Tokyo",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
		},
		{
			name:     "Single extension tag",
			input:    "2006-01-02T15:04:05Z[u-ca=japanese]",
			wantTime: "2006-01-02T15:04:05Z",
			wantExt: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{},
			},
		},
		{
			name:     "Critical extension tag",
			input:    "2006-01-02T15:04:05Z[!u-ca=japanese]",
			wantTime: "2006-01-02T15:04:05Z",
			wantExt: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{"u-ca": true},
			},
		},
		{
			name:     "Timezone with extension",
			input:    "2006-01-02T15:04:05+09:00[Asia/Tokyo][u-ca=japanese]",
			wantTime: "2006-01-02T15:04:05+09:00",
			wantExt: IXDTFExtensions{
				TimeZone: "Asia/Tokyo",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, gotExt, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) error = %v", tt.input, err)
				return
			}

			// Compare time strings since time comparison can be tricky with locations
			if got := gotTime.Format(time.RFC3339); got != tt.wantTime {
				t.Errorf("Parse(%q) time = %v, want %v", tt.input, got, tt.wantTime)
			}

			if gotExt.TimeZone != tt.wantExt.TimeZone {
				t.Errorf("Parse(%q) timezone = %v, want %v", tt.input, gotExt.TimeZone, tt.wantExt.TimeZone)
			}

			if len(gotExt.Tags) != len(tt.wantExt.Tags) {
				t.Errorf("Parse(%q) tags count = %v, want %v", tt.input, len(gotExt.Tags), len(tt.wantExt.Tags))
			}

			for k, v := range tt.wantExt.Tags {
				if gotExt.Tags[k] != v {
					t.Errorf("Parse(%q) tag[%s] = %v, want %v", tt.input, k, gotExt.Tags[k], v)
				}
			}

			for k, v := range tt.wantExt.Critical {
				if gotExt.Critical[k] != v {
					t.Errorf("Parse(%q) critical[%s] = %v, want %v", tt.input, k, gotExt.Critical[k], v)
				}
			}
		})
	}
}

func TestFormat(t *testing.T) {
	baseTime := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

	tests := []struct {
		name string
		time time.Time
		ext  IXDTFExtensions
		want string
	}{
		{
			name: "Basic RFC 3339 without extensions",
			time: baseTime,
			ext:  NewIXDTFExtensions(),
			want: "2006-01-02T15:04:05Z",
		},
		{
			name: "With timezone suffix",
			time: baseTime,
			ext: IXDTFExtensions{
				TimeZone: "UTC",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
			want: "2006-01-02T15:04:05Z[UTC]",
		},
		{
			name: "With extension tag",
			time: baseTime,
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{},
			},
			want: "2006-01-02T15:04:05Z[u-ca=japanese]",
		},
		{
			name: "With critical extension",
			time: baseTime,
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{"u-ca": true},
			},
			want: "2006-01-02T15:04:05Z[!u-ca=japanese]",
		},
		{
			name: "With timezone and extension",
			time: baseTime,
			ext: IXDTFExtensions{
				TimeZone: "UTC",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{},
			},
			want: "2006-01-02T15:04:05Z[UTC][u-ca=japanese]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Format(tt.time, tt.ext)
			if got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Invalid RFC 3339", "invalid-time[UTC]"},
		{"Invalid suffix format", "2006-01-02T15:04:05Z["},
		{"Empty suffix", "2006-01-02T15:04:05Z[]"},
		{"Invalid extension format", "2006-01-02T15:04:05Z[u-ca]"},
		{"Critical timezone", "2006-01-02T15:04:05Z[!UTC]"},
		{"Empty extension key", "2006-01-02T15:04:05Z[=value]"},
		{"Empty extension value", "2006-01-02T15:04:05Z[key=]"},
		{"Multiple timezone suffixes", "2006-01-02T15:04:05Z[UTC][GMT]"},
		{"Extension tag looks like timezone", "2006-01-02T15:04:05Z[x-test]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := Parse(tt.input)
			if err == nil {
				t.Errorf("Parse(%q) expected error but got none", tt.input)
			}
		})
	}
}

// TestParseErrorDetails tests the ParseError struct specifically
func TestParseErrorDetails(t *testing.T) {
	_, _, err := Parse("invalid-time")
	if err == nil {
		t.Fatal("Expected error but got none")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("Expected ParseError, got %T", err)
	}

	// Test Error() method with message
	errorStr := parseErr.Error()
	if errorStr == "" {
		t.Error("Error() should return non-empty string")
	}

	// Test Error() method without message
	parseErr2 := &ParseError{
		Layout: "test-layout",
		Value:  "test-value",
		Err:    nil,
	}
	errorStr2 := parseErr2.Error()
	expected := ""
	if errorStr2 != expected {
		t.Errorf("Error() = %q, want %q", errorStr2, expected)
	}
}

// TestFormatNano tests nanosecond formatting
func TestFormatNano(t *testing.T) {
	baseTime := time.Date(2006, 1, 2, 15, 4, 5, 123456789, time.UTC)

	tests := []struct {
		name string
		time time.Time
		ext  IXDTFExtensions
		want string
	}{
		{
			name: "Basic nanosecond formatting",
			time: baseTime,
			ext:  NewIXDTFExtensions(),
			want: "2006-01-02T15:04:05.123456789Z",
		},
		{
			name: "Nanosecond with timezone",
			time: baseTime,
			ext: IXDTFExtensions{
				TimeZone: "UTC",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
			want: "2006-01-02T15:04:05.123456789Z[UTC]",
		},
		{
			name: "Nanosecond with extensions",
			time: baseTime,
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{"u-ca": true},
			},
			want: "2006-01-02T15:04:05.123456789Z[!u-ca=japanese]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := FormatNano(tt.time, tt.ext)
			if got != tt.want {
				t.Errorf("FormatNano() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidate tests the Validate function and ValidateExtensions function
func TestValidate(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		ext          *IXDTFExtensions // nil if testing string validation
		wantError    bool
		wantErrorMsg string
	}{
		// String validation tests
		{
			name:      "Valid RFC 3339",
			input:     "2006-01-02T15:04:05Z",
			wantError: false,
		},
		{
			name:      "Valid with timezone",
			input:     "2006-01-02T15:04:05Z[UTC]",
			wantError: false,
		},
		{
			name:      "Valid with extensions",
			input:     "2006-01-02T15:04:05Z[UTC][u-ca=japanese]",
			wantError: false,
		},
		{
			name:         "Empty string",
			input:        "",
			wantError:    true,
			wantErrorMsg: "parsing time \"\" as \"2006-01-02T15:04:05Z07:00\": empty datetime string",
		},
		{
			name:         "Invalid suffix format - missing bracket",
			input:        "2006-01-02T15:04:05Z[UTC",
			wantError:    true,
			wantErrorMsg: "parsing time \"2006-01-02T15:04:05Z[UTC\" as \"2006-01-02T15:04:05Z07:00[time-zone]\": invalid IXDTF suffix format",
		},
		{
			name:         "Empty suffix",
			input:        "2006-01-02T15:04:05Z[]",
			wantError:    true,
			wantErrorMsg: "parsing time \"2006-01-02T15:04:05Z[]\" as \"2006-01-02T15:04:05Z07:00[time-zone]\": invalid IXDTF suffix format",
		},
		{
			name:         "Invalid extension format",
			input:        "2006-01-02T15:04:05Z[u-ca]",
			wantError:    true,
			wantErrorMsg: "parsing time \"2006-01-02T15:04:05Z[u-ca]\" as \"2006-01-02T15:04:05Z07:00[time-zone]\": invalid extension format",
		},
		{
			name:         "Critical timezone",
			input:        "2006-01-02T15:04:05Z[!UTC]",
			wantError:    true,
			wantErrorMsg: "parsing time \"2006-01-02T15:04:05Z[!UTC]\" as \"2006-01-02T15:04:05Z07:00[time-zone]\": invalid timezone name",
		},
		// Extension validation tests
		{
			name:      "Valid empty extensions",
			ext:       func() *IXDTFExtensions { ext := NewIXDTFExtensions(); return &ext }(),
			wantError: false,
		},
		{
			name: "Valid timezone",
			ext: &IXDTFExtensions{
				TimeZone: "UTC",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
			wantError: false,
		},
		{
			name: "Invalid timezone",
			ext: &IXDTFExtensions{
				TimeZone: "Invalid/Timezone",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
			wantError: true,
		},
		{
			name: "Valid critical extension - u-ca",
			ext: &IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{"u-ca": true},
			},
			wantError: false,
		},
		{
			name: "Valid critical extension - u-nu",
			ext: &IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-nu": "latn"},
				Critical: map[string]bool{"u-nu": true},
			},
			wantError: false,
		},
		{
			name: "Valid critical extension - private",
			ext: &IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"x-custom": "value"},
				Critical: map[string]bool{"x-custom": true},
			},
			wantError: false,
		},
		{
			name: "Invalid critical extension - empty value",
			ext: &IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": ""},
				Critical: map[string]bool{"u-ca": true},
			},
			wantError: true,
		},
		{
			name: "Invalid critical extension - missing tag",
			ext: &IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{},
				Critical: map[string]bool{"u-ca": true},
			},
			wantError: true,
		},
		{
			name: "Invalid critical extension - unknown",
			ext: &IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"unknown": "value"},
				Critical: map[string]bool{"unknown": true},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			if tt.ext != nil {
				// Test ValidateExtensions
				err = ValidateExtensions(*tt.ext)
			} else {
				// Test Validate (string validation)
				err = Validate(tt.input)
			}

			if tt.wantError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tt.wantErrorMsg != "" && err != nil && err.Error() != tt.wantErrorMsg {
				t.Errorf("Error = %v, want %v", err.Error(), tt.wantErrorMsg)
			}
		})
	}
}
