// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ixdtf

import (
	"testing"
	"time"
)

func TestParse_BasicRFC3339(t *testing.T) {
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
		})
	}
}

func TestParse_Extensions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantExt IXDTFExtensions
	}{
		{
			name:  "Single extension tag",
			input: "2006-01-02T15:04:05Z[u-ca=japanese]",
			wantExt: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{},
			},
		},
		{
			name:  "Critical extension tag",
			input: "2006-01-02T15:04:05Z[!u-ca=japanese]",
			wantExt: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{"u-ca": true},
			},
		},
		{
			name:  "Timezone with extension",
			input: "2006-01-02T15:04:05+09:00[Asia/Tokyo][u-ca=japanese]",
			wantExt: IXDTFExtensions{
				TimeZone: "Asia/Tokyo",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotExt, err := Parse(tt.input)
			if err != nil {
				t.Errorf("Parse(%q) error = %v", tt.input, err)
				return
			}

			if gotExt.TimeZone != tt.wantExt.TimeZone {
				t.Errorf("Parse(%q) timezone = %v, want %v", tt.input, gotExt.TimeZone, tt.wantExt.TimeZone)
			}

			if len(gotExt.Tags) != len(tt.wantExt.Tags) {
				t.Errorf("Parse(%q) tags = %v, want %v", tt.input, gotExt.Tags, tt.wantExt.Tags)
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
			got := Format(tt.time, tt.ext)
			if got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05Z[UTC]",
		"2006-01-02T15:04:05+09:00[Asia/Tokyo]",
		"2006-01-02T15:04:05Z[u-ca=japanese]",
		"2006-01-02T15:04:05Z[UTC][u-ca=japanese]",
		"2006-01-02T15:04:05Z[!u-ca=japanese]",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			parsedTime, parsedExt, err := Parse(tt)
			if err != nil {
				t.Errorf("Parse(%q) error = %v", tt, err)
				return
			}

			formatted := Format(parsedTime, parsedExt)

			// For timezone cases, we might get a different but equivalent representation
			// due to timezone loading, so we parse again to compare
			reparsedTime, reparsedExt, err := Parse(formatted)
			if err != nil {
				t.Errorf("Re-parsing formatted result failed: %v", err)
				return
			}

			// Times should be equal (within timezone differences)
			if !parsedTime.Equal(reparsedTime) {
				t.Errorf("Round trip time mismatch: original=%v, reparsed=%v", parsedTime, reparsedTime)
			}

			// Extensions should be equal
			if parsedExt.TimeZone != reparsedExt.TimeZone {
				t.Errorf("Round trip timezone mismatch: original=%v, reparsed=%v", parsedExt.TimeZone, reparsedExt.TimeZone)
			}

			if len(parsedExt.Tags) != len(reparsedExt.Tags) {
				t.Errorf("Round trip tags count mismatch: original=%d, reparsed=%d", len(parsedExt.Tags), len(reparsedExt.Tags))
			}

			for k, v := range parsedExt.Tags {
				if reparsedExt.Tags[k] != v {
					t.Errorf("Round trip tag[%s] mismatch: original=%v, reparsed=%v", k, v, reparsedExt.Tags[k])
				}
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
		Msg:    "",
	}
	errorStr2 := parseErr2.Error()
	expected := "parsing time \"test-value\" as \"test-layout\": cannot parse"
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
			got := FormatNano(tt.time, tt.ext)
			if got != tt.want {
				t.Errorf("FormatNano() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidate tests the Validate function
func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
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
			name:      "Empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "Invalid suffix format - missing bracket",
			input:     "2006-01-02T15:04:05Z[UTC",
			wantError: true,
		},
		{
			name:      "Empty suffix",
			input:     "2006-01-02T15:04:05Z[]",
			wantError: true,
		},
		{
			name:      "Invalid extension format",
			input:     "2006-01-02T15:04:05Z[u-ca]",
			wantError: true,
		},
		{
			name:      "Critical timezone",
			input:     "2006-01-02T15:04:05Z[!UTC]",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if tt.wantError && err == nil {
				t.Errorf("Validate(%q) expected error but got none", tt.input)
			}
			if !tt.wantError && err != nil {
				t.Errorf("Validate(%q) unexpected error: %v", tt.input, err)
			}
		})
	}
}

// TestIsValidIXDTF tests the IsValidIXDTF function
func TestIsValidIXDTF(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"Valid RFC 3339", "2006-01-02T15:04:05Z", true},
		{"Valid with timezone", "2006-01-02T15:04:05Z[UTC]", true},
		{"Invalid format", "invalid", false},
		{"Empty suffix", "2006-01-02T15:04:05Z[]", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidIXDTF(tt.input)
			if got != tt.want {
				t.Errorf("IsValidIXDTF(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestValidateExtensions tests the ValidateExtensions function
func TestValidateExtensions(t *testing.T) {
	tests := []struct {
		name      string
		ext       IXDTFExtensions
		wantError bool
	}{
		{
			name: "Valid empty extensions",
			ext:  NewIXDTFExtensions(),
		},
		{
			name: "Valid timezone",
			ext: IXDTFExtensions{
				TimeZone: "UTC",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
		},
		{
			name: "Invalid timezone",
			ext: IXDTFExtensions{
				TimeZone: "Invalid/Timezone",
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			},
			wantError: true,
		},
		{
			name: "Valid critical extension - u-ca",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": "japanese"},
				Critical: map[string]bool{"u-ca": true},
			},
		},
		{
			name: "Valid critical extension - u-nu",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-nu": "latn"},
				Critical: map[string]bool{"u-nu": true},
			},
		},
		{
			name: "Valid critical extension - private",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"x-custom": "value"},
				Critical: map[string]bool{"x-custom": true},
			},
		},
		{
			name: "Invalid critical extension - empty value",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"u-ca": ""},
				Critical: map[string]bool{"u-ca": true},
			},
			wantError: true,
		},
		{
			name: "Invalid critical extension - missing tag",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{},
				Critical: map[string]bool{"u-ca": true},
			},
			wantError: true,
		},
		{
			name: "Invalid critical extension - unknown",
			ext: IXDTFExtensions{
				TimeZone: "",
				Tags:     map[string]string{"unknown": "value"},
				Critical: map[string]bool{"unknown": true},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExtensions(tt.ext)
			if tt.wantError && err == nil {
				t.Errorf("ValidateExtensions() expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("ValidateExtensions() unexpected error: %v", err)
			}
		})
	}
}



