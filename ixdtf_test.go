package ixdtf_test

import (
	"errors"
	"testing"
	"time"

	"github.com/8beeeaaat/ixdtf"
)

func TestFormat(t *testing.T) {
	tokyo, paris, cet := getTestTimezones()

	tests := []struct {
		name    string
		tm      time.Time
		ext     *ixdtf.IXDTFExtensions
		want    string
		wantErr error
	}{
		{
			name: "no extensions",
			tm:   time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
			ext:  ixdtf.NewIXDTFExtensions(nil),
			want: "2025-01-02T03:04:05Z",
		},
		{
			name: "use location as timezone",
			tm:   time.Date(2025, 2, 3, 4, 5, 6, 0, tokyo),
			ext:  ixdtf.NewIXDTFExtensions(nil),
			want: "2025-02-03T04:05:06+09:00[Asia/Tokyo]",
		},
		{
			name: "use specified timezone",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: tokyo,
			}),
			want: "2025-01-01T00:00:00Z[Asia/Tokyo]",
		},
		{
			name: "tags sorting and critical",
			tm:   time.Date(2025, 3, 4, 5, 6, 7, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{
					"b-tag": "2",
					"a-tag": "1",
				},
				Critical: map[string]bool{
					"b-tag": true,
				},
			}),
			want: "2025-03-04T05:06:07Z[a-tag=1][!b-tag=2]",
		},
		{
			name: "timezone and tags",
			tm:   time.Date(2025, 6, 7, 8, 9, 10, 0, cet),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: paris,
				Tags: map[string]string{
					"u-ca": "gregory",
				},
				Critical: map[string]bool{
					"u-ca": true,
				},
			}),
			want: "2025-06-07T08:09:10+01:00[Europe/Paris][!u-ca=gregory]",
		},
		{
			name: "private extension error",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{
					"x-demo": "yes",
				},
			}),
			wantErr: ixdtf.ErrPrivateExtension,
		},
		{
			name: "invalid timezone",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: time.FixedZone("No/SuchZone", 1*3600),
			}),
			wantErr: ixdtf.ErrInvalidTimezone,
		},
		{
			name: "missing critical tag",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Critical: map[string]bool{"u-ca": true},
			}),
			wantErr: ixdtf.ErrCriticalExtension,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ixdtf.Format(tc.tm, tc.ext)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.want {
				t.Fatalf("want %q got %q", tc.want, got)
			}
		})
	}
}

func TestFormatNano(t *testing.T) {
	tokyo, paris, cet := getTestTimezones()

	tests := []struct {
		name    string
		tm      time.Time
		ext     *ixdtf.IXDTFExtensions
		want    string
		wantErr string
	}{
		{
			name: "basic with nanoseconds",
			tm:   time.Date(2025, 1, 2, 3, 4, 5, 123456789, time.UTC),
			ext:  ixdtf.NewIXDTFExtensions(nil),
			want: "2025-01-02T03:04:05.123456789Z",
		},
		{
			name: "with timezone and nanoseconds",
			tm:   time.Date(2025, 2, 3, 4, 5, 6, 789000000, tokyo),
			ext:  ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{Location: tokyo}),
			want: "2025-02-03T04:05:06.789+09:00[Asia/Tokyo]",
		},
		{
			name: "with tags and nanoseconds",
			tm:   time.Date(2025, 3, 4, 5, 6, 7, 500000000, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{"u-ca": "gregory"},
			}),
			want: "2025-03-04T05:06:07.5Z[u-ca=gregory]",
		},
		{
			name: "with timezone, tags and nanoseconds",
			tm:   time.Date(2025, 6, 7, 8, 9, 10, 250000000, paris),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: paris,
				Tags:     map[string]string{"u-ca": "gregory"},
				Critical: map[string]bool{"u-ca": true},
			}),
			want: "2025-06-07T08:09:10.25+01:00[Europe/Paris][!u-ca=gregory]",
		},
		{
			name: "CET timezone with nanoseconds",
			tm:   time.Date(2025, 12, 25, 15, 30, 45, 999000000, cet),
			ext:  ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{Location: cet}),
			want: "2025-12-25T15:30:45.999+01:00[CET]",
		},
		{
			name: "zero nanoseconds",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext:  ixdtf.NewIXDTFExtensions(nil),
			want: "2025-01-01T00:00:00Z",
		},
		{
			name: "invalid critical extension - empty value",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags:     map[string]string{"u-ca": ""},
				Critical: map[string]bool{"u-ca": true},
			}),
			wantErr: "critical extension cannot be processed",
		},
		{
			name: "missing critical extension",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags:     map[string]string{},
				Critical: map[string]bool{"u-ca": true},
			}),
			wantErr: "critical extension cannot be processed",
		},
		{
			name: "invalid tag key format",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{"INVALID-KEY": "value"},
			}),
			wantErr: "invalid extension format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ixdtf.FormatNano(tc.tm, tc.ext)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("FormatNano(%v, %+v) expected error containing %q, got nil", tc.tm, tc.ext, tc.wantErr)
				}
				if !containsString(err.Error(), tc.wantErr) {
					t.Fatalf(
						"FormatNano(%v, %+v) error = %q, want error containing %q",
						tc.tm,
						tc.ext,
						err.Error(),
						tc.wantErr,
					)
				}
				return
			}

			if err != nil {
				t.Fatalf("FormatNano(%v, %+v) unexpected error: %v", tc.tm, tc.ext, err)
			}

			if got != tc.want {
				t.Errorf("FormatNano(%v, %+v) = %q, want %q", tc.tm, tc.ext, got, tc.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tokyo, paris, cet := getTestTimezones()

	tests := []struct {
		name     string
		input    string
		strict   bool
		wantTime time.Time
		wantExt  *ixdtf.IXDTFExtensions
		wantErr  string
	}{
		{
			name:     "basic RFC3339",
			input:    "2025-01-02T03:04:05Z",
			strict:   false,
			wantTime: time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC),
			wantExt:  ixdtf.NewIXDTFExtensions(nil),
		},
		{
			name:     "RFC3339 with nanoseconds",
			input:    "2025-01-02T03:04:05.123456789Z",
			strict:   false,
			wantTime: time.Date(2025, 1, 2, 3, 4, 5, 123456789, time.UTC),
			wantExt:  ixdtf.NewIXDTFExtensions(nil),
		},
		{
			name:     "with timezone",
			input:    "2025-02-03T04:05:06+09:00[Asia/Tokyo]",
			strict:   false,
			wantTime: time.Date(2025, 2, 3, 4, 5, 6, 0, tokyo),
			wantExt:  ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{Location: tokyo}),
		},
		{
			name:     "with tags",
			input:    "2025-03-04T05:06:07Z[u-ca=gregory]",
			strict:   false,
			wantTime: time.Date(2025, 3, 4, 5, 6, 7, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{"u-ca": "gregory"},
			}),
		},
		{
			name:     "with critical tag",
			input:    "2025-03-04T05:06:07Z[!u-ca=gregory]",
			strict:   false,
			wantTime: time.Date(2025, 3, 4, 5, 6, 7, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags:     map[string]string{"u-ca": "gregory"},
				Critical: map[string]bool{"u-ca": true},
			}),
		},
		{
			name:     "with timezone and tags",
			input:    "2025-06-07T08:09:10+01:00[Europe/Paris][!u-ca=gregory]",
			strict:   false,
			wantTime: time.Date(2025, 6, 7, 8, 9, 10, 0, paris),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: paris,
				Tags:     map[string]string{"u-ca": "gregory"},
				Critical: map[string]bool{"u-ca": true},
			}),
		},
		{
			name:     "timezone with CET",
			input:    "2025-12-25T15:30:45+01:00[CET]",
			strict:   false,
			wantTime: time.Date(2025, 12, 25, 15, 30, 45, 0, cet),
			wantExt:  ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{Location: cet}),
		},
		{
			name:    "invalid timezone",
			input:   "2025-01-01T00:00:00Z[Invalid/Zone]",
			strict:  false,
			wantErr: "IXDTFE parsing time",
		},
		{
			name:    "malformed suffix",
			input:   "2025-01-01T00:00:00Z[unclosed",
			strict:  false,
			wantErr: "IXDTFE parsing time",
		},
		{
			name:    "empty string",
			input:   "",
			strict:  false,
			wantErr: "IXDTFE parsing time",
		},
		{
			name:    "invalid RFC3339",
			input:   "not-a-date",
			strict:  false,
			wantErr: "IXDTFE parsing time",
		},
		{
			name:     "timezone offset mismatch - non-strict",
			input:    "2025-06-01T12:00:00+09:00[America/New_York]",
			strict:   false,
			wantTime: time.Date(2025, 6, 1, 12, 0, 0, 0, time.FixedZone("+0900", 9*3600)),
			wantExt: ixdtf.NewIXDTFExtensions(
				&ixdtf.NewIXDTFExtensionsArgs{Location: time.FixedZone("America/New_York", -4*3600)},
			),
		},
		{
			name:    "timezone offset mismatch - strict",
			input:   "2025-06-01T12:00:00+09:00[America/New_York]",
			strict:  true,
			wantErr: "timezone offset does not match",
		},
		{
			name:     "timezone with Etc/GMT pattern - non-strict",
			input:    "2025-01-01T00:00:00+05:00[Etc/GMT-5]",
			strict:   false,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.FixedZone("Etc/GMT-5", 5*3600)),
			wantExt: ixdtf.NewIXDTFExtensions(
				&ixdtf.NewIXDTFExtensionsArgs{Location: time.FixedZone("Etc/GMT-5", 5*3600)},
			),
		},
		{
			name:     "timezone with Etc/GMT pattern - strict",
			input:    "2025-01-01T00:00:00+05:00[Etc/GMT-5]",
			strict:   true,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.FixedZone("Etc/GMT-5", 5*3600)),
			wantExt: ixdtf.NewIXDTFExtensions(
				&ixdtf.NewIXDTFExtensionsArgs{Location: time.FixedZone("Etc/GMT-5", 5*3600)},
			),
		},
		{
			name:    "invalid suffix key format in parsing",
			input:   "2025-01-01T00:00:00Z[INVALID-KEY=value]",
			strict:  false,
			wantErr: "IXDTFE parsing time",
		},
		{
			name:    "suffix with extension-like patterns",
			input:   "2025-01-01T00:00:00Z[u-ca=gregory][t-invalid]",
			strict:  false,
			wantErr: "IXDTFE parsing time",
		},
		{
			name:    "suffix with private extension",
			input:   "2025-01-01T00:00:00Z[x-private=test]",
			strict:  false,
			wantErr: "private extension cannot be processed",
		},
		{
			name:    "suffix with experimental extension",
			input:   "2025-01-01T00:00:00Z[_experiment=test]",
			strict:  false,
			wantErr: "experimental extension cannot be processed",
		},
		{
			name:    "critical private extension",
			input:   "2025-01-01T00:00:00Z[!x-private=test]",
			strict:  false,
			wantErr: "private extension cannot be processed",
		},
		{
			name:    "critical experimental extension",
			input:   "2025-01-01T00:00:00Z[!_experiment=test]",
			strict:  false,
			wantErr: "experimental extension cannot be processed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTime, gotExt, err := ixdtf.Parse(tc.input, tc.strict)

			if tc.wantErr != "" {
				checkParseError(t, err, tc.input, tc.strict, tc.wantErr)
				return
			}

			if err != nil {
				t.Fatalf("Parse(%q, %v) unexpected error: %v", tc.input, tc.strict, err)
			}

			compareParseResults(t, gotTime, gotExt, tc.input, tc.strict, tc.wantTime, tc.wantExt)
		})
	}
}

func TestParseErrorError(t *testing.T) {
	tests := []struct {
		name     string
		pe       *ixdtf.ParseError
		expected string
	}{
		{
			name:     "nil underlying error returns empty string",
			pe:       &ixdtf.ParseError{Err: nil, Layout: ixdtf.LayoutRFC3339, Value: "2025-01-01T00:00:00Z"},
			expected: "",
		},
		{
			name:     "with underlying error and known layout",
			pe:       &ixdtf.ParseError{Err: errors.New("invalid time"), Layout: ixdtf.LayoutRFC3339, Value: "bad"},
			expected: "IXDTFE parsing time \"bad\" as \"" + string(ixdtf.LayoutRFC3339) + "\": invalid time",
		},
		{
			name:     "with underlying error and empty layout",
			pe:       &ixdtf.ParseError{Err: errors.New("parse fail"), Layout: "", Value: "foo"},
			expected: "IXDTFE parsing time \"foo\" as \"\": parse fail",
		},
		{
			name:     "with custom layout value",
			pe:       &ixdtf.ParseError{Err: errors.New("boom"), Layout: ixdtf.Layout("CUSTOM"), Value: "x"},
			expected: "IXDTFE parsing time \"x\" as \"CUSTOM\": boom",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.pe.Error()
			if got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		strict  bool
		wantErr string
	}{
		{
			name:   "valid RFC3339",
			input:  "2025-01-02T03:04:05Z",
			strict: false,
		},
		{
			name:   "valid RFC3339 with nanoseconds",
			input:  "2025-01-02T03:04:05.123456789Z",
			strict: false,
		},
		{
			name:   "valid with timezone",
			input:  "2025-02-03T04:05:06Z[Asia/Tokyo]",
			strict: false,
		},
		{
			name:   "valid with offset and timezone",
			input:  "2025-02-03T04:05:06+09:00[Asia/Tokyo]",
			strict: false,
		},
		{
			name:   "valid with tags",
			input:  "2025-03-04T05:06:07Z[u-ca=gregory]",
			strict: false,
		},
		{
			name:   "valid with critical tag",
			input:  "2025-03-04T05:06:07Z[!u-ca=gregory]",
			strict: false,
		},
		{
			name:   "valid with timezone and tags",
			input:  "2025-06-07T08:09:10+01:00[Europe/Paris][!u-ca=gregory]",
			strict: false,
		},
		{
			name:    "invalid timezone name",
			input:   "2025-01-01T00:00:00Z[No/SuchZone]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[No/SuchZone]\" as \"2006-01-02T15:04:05.999999999Z07:00*([time-zone-name][tags])\": unknown time zone No/SuchZone",
		},
		{
			name:    "invalid extension format",
			input:   "2025-01-01T00:00:00Z[invalid]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[invalid]\" as \"2006-01-02T15:04:05.999999999Z07:00*([time-zone-name][tags])\": unknown time zone invalid",
		},
		{
			name:    "private extension",
			input:   "2025-01-01T00:00:00Z[x-demo=yes]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[x-demo=yes]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": private extension cannot be processed",
		},
		{
			name:    "experimental extension",
			input:   "2025-01-01T00:00:00Z[_test=value]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[_test=value]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": experimental extension cannot be processed",
		},
		{
			name:    "empty string",
			input:   "",
			strict:  false,
			wantErr: "IXDTFE parsing time \"\" as \"2006-01-02T15:04:05Z07:00\": empty datetime string",
		},
		{
			name:    "invalid RFC3339",
			input:   "not-a-date",
			strict:  false,
			wantErr: "IXDTFE parsing time \"not-a-date\" as \"2006-01-02T15:04:05Z07:00\": invalid portion: parsing time \"not-a-date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"not-a-date\" as \"2006\"",
		},
		{
			name:    "malformed suffix",
			input:   "2025-01-01T00:00:00Z[unclosed",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[unclosed\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid IXDTF suffix format",
		},
		{
			name:   "timezone offset mismatch - non-strict",
			input:  "2025-06-01T12:00:00+09:00[America/New_York]",
			strict: false,
		},
		{
			name:    "timezone offset mismatch - strict",
			input:   "2025-06-01T12:00:00+09:00[America/New_York]",
			strict:  true,
			wantErr: "IXDTFE parsing time \"2025-06-01T12:00:00+09:00[America/New_York]\" as \"2006-01-02T15:04:05.999999999Z07:00*([time-zone-name][tags])\": timezone offset does not match the specified timezone",
		},
		{
			name:    "timezone content with invalid extension prefix u-",
			input:   "2025-01-01T00:00:00Z[u-invalid-timezone]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[u-invalid-timezone]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "timezone content with invalid extension prefix x-",
			input:   "2025-01-01T00:00:00Z[x-invalid-timezone]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[x-invalid-timezone]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "timezone content with invalid extension prefix t-",
			input:   "2025-01-01T00:00:00Z[t-invalid-timezone]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[t-invalid-timezone]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "suffix with invalid characters in key",
			input:   "2025-01-01T00:00:00Z[invalid@key=value]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[invalid@key=value]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "suffix with invalid characters in value range",
			input:   "2025-01-01T00:00:00Z[key=invalid@value]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[key=invalid@value]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "critical extension with invalid value range",
			input:   "2025-01-01T00:00:00Z[!key=invalid@value]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[!key=invalid@value]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:   "valid suffix with hyphenated key",
			input:  "2025-01-01T00:00:00Z[key-with-hyphens=value]",
			strict: false,
		},
		{
			name:   "valid suffix with hyphenated value",
			input:  "2025-01-01T00:00:00Z[key=value-with-hyphens]",
			strict: false,
		},
		{
			name:   "valid suffix with numeric value",
			input:  "2025-01-01T00:00:00Z[key=123]",
			strict: false,
		},
		{
			name:   "valid suffix with mixed alphanumeric",
			input:  "2025-01-01T00:00:00Z[key=abc123def]",
			strict: false,
		},
		{
			name:    "empty suffix key",
			input:   "2025-01-01T00:00:00Z[=value]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[=value]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "empty suffix value",
			input:   "2025-01-01T00:00:00Z[key=]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[key=]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "empty timezone brackets",
			input:   "2025-01-01T00:00:00Z[]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid IXDTF suffix format",
		},
		{
			name:    "suffix key starting with number",
			input:   "2025-01-01T00:00:00Z[123key=value]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[123key=value]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "suffix key with spaces",
			input:   "2025-01-01T00:00:00Z[key with spaces=value]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[key with spaces=value]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
		{
			name:    "suffix value with underscore",
			input:   "2025-01-01T00:00:00Z[key=val_ue]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[key=val_ue]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ixdtf.Validate(tc.input, tc.strict)
			if (err != nil && err.Error() != tc.wantErr) || (err == nil && tc.wantErr != "") {
				t.Fatalf("Validate(%q, %v) error = %v, wantErr %v", tc.input, tc.strict, err, tc.wantErr)
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr) && contains(s, substr)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func checkParseError(t *testing.T, err error, input string, strict bool, wantErr string) {
	if err == nil {
		t.Fatalf("Parse(%q, %v) expected error containing %q, got nil", input, strict, wantErr)
	}
	if !containsString(err.Error(), wantErr) {
		t.Fatalf(
			"Parse(%q, %v) error = %q, want error containing %q",
			input,
			strict,
			err.Error(),
			wantErr,
		)
	}
}

func compareParseResults(
	t *testing.T,
	gotTime time.Time,
	gotExt *ixdtf.IXDTFExtensions,
	input string,
	strict bool,
	wantTime time.Time,
	wantExt *ixdtf.IXDTFExtensions,
) {
	// Compare times
	if !gotTime.Equal(wantTime) {
		t.Errorf("Parse(%q, %v) time = %v, want %v", input, strict, gotTime, wantTime)
	}

	// Compare extensions
	if !extensionsEqual(gotExt, wantExt) {
		t.Errorf("Parse(%q, %v) extensions = %+v, want %+v", input, strict, gotExt, wantExt)
	}
}

func extensionsEqual(a, b *ixdtf.IXDTFExtensions) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare locations
	if (a.Location == nil) != (b.Location == nil) {
		return false
	}
	if a.Location != nil && b.Location != nil {
		if a.Location.String() != b.Location.String() {
			return false
		}
	}

	// Compare tags
	if len(a.Tags) != len(b.Tags) {
		return false
	}
	for k, v := range a.Tags {
		if b.Tags[k] != v {
			return false
		}
	}

	// Compare critical
	if len(a.Critical) != len(b.Critical) {
		return false
	}
	for k, v := range a.Critical {
		if b.Critical[k] != v {
			return false
		}
	}

	return true
}

func getTestTimezones() (*time.Location, *time.Location, *time.Location) {
	return time.FixedZone("Asia/Tokyo", 9*3600),
		time.FixedZone("Europe/Paris", 1*3600),
		time.FixedZone("CET", 1*3600)
}
