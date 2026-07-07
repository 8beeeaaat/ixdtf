package ixdtf_test

import (
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/8beeeaaat/ixdtf"
)

func TestFormat(t *testing.T) {
	t.Parallel()
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
			// A critical time-zone flag without a zone cannot be honored;
			// silently dropping the "!" would violate RFC 9557 Section 3.3.
			name: "critical location flag without location errors",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				CriticalLocation: true,
			}),
			wantErr: ixdtf.ErrCriticalExtension,
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
		{
			name: "invalid critical calendar value",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags:     map[string]string{"u-ca": "hoge"},
				Critical: map[string]bool{"u-ca": true},
			}),
			wantErr: ixdtf.ErrInvalidTagCalendarIdentifier,
		},
		{
			name: "empty timezone name should not add brackets",
			tm:   time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: time.FixedZone("", 0),
			}),
			want: "2025-01-01T12:00:00Z",
		},
		{
			name: "nil extensions behaves like empty extensions",
			tm:   time.Date(2025, 2, 3, 4, 5, 6, 0, tokyo),
			ext:  nil,
			want: "2025-02-03T04:05:06+09:00[Asia/Tokyo]",
		},
		{
			// RFC 9557 Section 1.2: an offset time zone is serialized using
			// the same numeric form as the RFC 3339 offset ("+09:00", not the
			// Go zone-name convention "+0900").
			name: "offset time zone serialized in RFC form",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.FixedZone("+09:00", 9*3600)),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: time.FixedZone("+09:00", 9*3600),
			}),
			want: "2025-01-01T00:00:00+09:00[+09:00]",
		},
		{
			name: "critical offset time zone keeps flag and RFC form",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.FixedZone("+09:00", 9*3600)),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location:         time.FixedZone("+09:00", 9*3600),
				CriticalLocation: true,
			}),
			want: "2025-01-01T00:00:00+09:00[!+09:00]",
		},
		{
			// The critical flag also applies to the timestamp's own named
			// zone when ext.Location is unset — the same fallback used for
			// non-critical output, so the "!" is not silently dropped.
			name: "critical flag applies to fallback zone from timestamp",
			tm:   time.Date(2025, 2, 3, 4, 5, 6, 0, tokyo),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				CriticalLocation: true,
			}),
			want: "2025-02-03T04:05:06+09:00[!Asia/Tokyo]",
		},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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

// TestOffsetTimezoneRoundTrip verifies that offset time-zone annotations
// survive a Parse -> Format -> Parse round trip in the RFC 9557 Section 1.2
// serialization form, including the Section 1.2 example "+08:45[+08:45]" and
// the critical variant.
func TestOffsetTimezoneRoundTrip(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"2025-01-01T00:00:00+09:00[+09:00]",
		"2022-07-08T00:14:07+08:45[+08:45]", // RFC 9557 Section 1.2 example
		"2025-01-01T00:00:00+09:00[!+09:00]",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			tm, ext, err := ixdtf.Parse(input, true)
			if err != nil {
				t.Fatalf("Parse(%q, true) unexpected error: %v", input, err)
			}

			formatted, err := ixdtf.Format(tm, ext)
			if err != nil {
				t.Fatalf("Format after Parse(%q) unexpected error: %v", input, err)
			}
			if formatted != input {
				t.Fatalf("round trip = %q, want %q", formatted, input)
			}

			if _, _, reparseErr := ixdtf.Parse(formatted, true); reparseErr != nil {
				t.Fatalf("re-Parse(%q, true) unexpected error: %v", formatted, reparseErr)
			}
			if validateErr := ixdtf.Validate(formatted, true); validateErr != nil {
				t.Fatalf("Validate(%q, true) unexpected error: %v", formatted, validateErr)
			}
		})
	}
}

func TestFormatNano(t *testing.T) {
	t.Parallel()
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
			name: "multiple tags with different prefixes and nanoseconds",
			tm:   time.Date(2025, 3, 4, 5, 6, 7, 500000000, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{
					"t-calendar": "japanese",
					"u-ca":       "gregory",
				},
			}),
			want: "2025-03-04T05:06:07.5Z[t-calendar=japanese][u-ca=gregory]",
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
			name: "invalid critical calendar value",
			tm:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			ext: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags:     map[string]string{"u-ca": "hoge"},
				Critical: map[string]bool{"u-ca": true},
			}),
			wantErr: "invalid calendar tag identifier",
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

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ixdtf.FormatNano(tc.tm, tc.ext)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("FormatNano(%v, %+v) expected error containing %q, got nil", tc.tm, tc.ext, tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
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
