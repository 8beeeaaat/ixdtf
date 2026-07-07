package ixdtf_test

import (
	"sort"
	"testing"
	"time"

	"github.com/8beeeaaat/ixdtf"
)

func TestParse(t *testing.T) {
	t.Parallel()
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
			// RFC 9557 Section 3.4 (Figure 2): Z is an unknown local offset, so
			// pairing it with a timezone is consistent even in strict mode.
			name:     "Z with IANA timezone is consistent in strict mode",
			input:    "2022-07-08T00:14:07Z[Europe/London]",
			strict:   true,
			wantTime: time.Date(2022, 7, 8, 0, 14, 7, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: time.FixedZone("Europe/London", 0),
			}),
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
			name:     "with invalid calendar tag (non-critical)",
			input:    "2025-03-04T05:06:07Z[u-ca=hoge]",
			strict:   false,
			wantTime: time.Date(2025, 3, 4, 5, 6, 7, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{"u-ca": "hoge"},
			}),
		},
		{
			name:    "with invalid calendar tag (non-critical, strict)",
			input:   "2025-03-04T05:06:07Z[u-ca=hoge]",
			strict:  true,
			wantErr: "invalid calendar tag identifier",
		},
		{
			name:    "with invalid calendar tag (critical)",
			input:   "2025-03-04T05:06:07Z[!u-ca=hoge]",
			strict:  false,
			wantErr: "invalid calendar tag identifier",
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
			name:     "multiple tags with different prefixes",
			input:    "2025-03-04T05:06:07Z[t-calendar=japanese][u-ca=gregory]",
			strict:   false,
			wantTime: time.Date(2025, 3, 4, 5, 6, 7, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{
					"t-calendar": "japanese",
					"u-ca":       "gregory",
				},
			}),
		},
		{
			name:     "multiple tags with critical flag",
			input:    "2025-03-04T05:06:07Z[!t-format=iso][u-ca=hebrew]",
			strict:   false,
			wantTime: time.Date(2025, 3, 4, 5, 6, 7, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags: map[string]string{
					"t-format": "iso",
					"u-ca":     "hebrew",
				},
				Critical: map[string]bool{
					"t-format": true,
				},
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
			name:     "invalid timezone",
			input:    "2025-01-01T00:00:00Z[Foo/Bar]",
			strict:   false,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: nil,
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			}),
		},
		{
			name:    "invalid timezone in strict mode",
			input:   "2025-01-01T00:00:00Z[Foo/Bar]",
			strict:  true,
			wantErr: "invalid timezone",
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
			name:     "timezone offset mismatch in non-strict mode",
			input:    "2025-06-01T12:00:00+09:00[America/New_York]",
			strict:   false,
			wantTime: time.Date(2025, 6, 1, 12, 0, 0, 0, time.FixedZone("+0900", 9*3600)),
			wantExt: ixdtf.NewIXDTFExtensions(
				&ixdtf.NewIXDTFExtensionsArgs{Location: time.FixedZone("America/New_York", -4*3600)},
			),
		},
		{
			name:    "timezone offset mismatch in strict mode",
			input:   "2025-06-01T12:00:00+09:00[America/New_York]",
			strict:  true,
			wantErr: "timezone offset does not match",
		},
		{
			name:     "timezone with Etc/GMT pattern in non-strict mode",
			input:    "2025-01-01T00:00:00+05:00[Etc/GMT-5]",
			strict:   false,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.FixedZone("Etc/GMT-5", 5*3600)),
			wantExt: ixdtf.NewIXDTFExtensions(
				&ixdtf.NewIXDTFExtensionsArgs{Location: time.FixedZone("Etc/GMT-5", 5*3600)},
			),
		},
		{
			// Etc/GMT+3 is UTC-03:00 (POSIX-inverted sign); a concrete +05:00
			// conflicts with it, and the critical flag forces the error even
			// in non-strict mode (RFC 9557 Section 3.4).
			name:    "critical Etc/GMT zone with mismatching offset errors in non-strict",
			input:   "2022-07-08T00:14:07+05:00[!Etc/GMT+3]",
			strict:  false,
			wantErr: "timezone offset does not match",
		},
		{
			name:    "Etc/GMT zone with mismatching offset errors in strict mode",
			input:   "2022-07-08T00:14:07+05:00[Etc/GMT+3]",
			strict:  true,
			wantErr: "timezone offset does not match",
		},
		{
			name:     "timezone with Etc/GMT pattern in strict mode",
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
			// "[t-invalid]" has no "=" so it can only be a time-zone
			// annotation, which the RFC 9557 Section 4.1 grammar
			// ("suffix = [time-zone] *suffix-tag") forbids after a suffix tag.
			name:    "timezone-shaped element after tag - non-strict mode",
			input:   "2025-01-01T00:00:00Z[!u-ca=gregory][t-invalid]",
			strict:  false,
			wantErr: "invalid IXDTF suffix format",
		},
		{
			name:    "timezone-shaped element after tag - strict mode",
			input:   "2025-01-01T00:00:00Z[!u-ca=gregory][t-invalid]",
			strict:  true,
			wantErr: "invalid IXDTF suffix format",
		},
		{
			name:    "timezone annotation after suffix tag rejected in strict",
			input:   "2025-03-04T05:06:07Z[u-ca=hebrew][Asia/Tokyo]",
			strict:  true,
			wantErr: "invalid IXDTF suffix format",
		},
		{
			name:    "timezone annotation after suffix tag rejected in non-strict",
			input:   "2025-03-04T05:06:07Z[u-ca=hebrew][Asia/Tokyo]",
			strict:  false,
			wantErr: "invalid IXDTF suffix format",
		},
		{
			// Even when the first zone is unknown and ignored by a non-strict
			// parse, a second time-zone annotation stays a grammar violation;
			// it must not be silently applied in place of the first.
			name:    "second timezone annotation after ignored unknown zone rejected",
			input:   "2022-07-08T00:14:07Z[Foo/Bar][Asia/Tokyo]",
			strict:  false,
			wantErr: "invalid IXDTF suffix format",
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
		{
			// RFC 9557 Section 4.1 permits a critical "!" flag on a time zone;
			// combined with Z (unknown offset) it is consistent (Figure 2).
			name:     "critical timezone with Z is accepted",
			input:    "2025-01-01T00:00:00Z[!Asia/Tokyo]",
			strict:   false,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location:         time.FixedZone("Asia/Tokyo", 9*3600),
				CriticalLocation: true,
			}),
		},
		{
			// Critical + concrete offset that conflicts with the zone: an
			// application MUST act on the inconsistency even in non-strict
			// mode (RFC 9557 Section 3.4).
			name:    "critical timezone inconsistent offset errors in non-strict",
			input:   "2022-07-08T00:14:07+00:00[!Europe/London]",
			strict:  false,
			wantErr: "timezone offset does not match",
		},
		{
			// "-00:00" is the other unknown-local-offset form (RFC 9557
			// Section 2.2), so a critical zone pairs with it consistently.
			name:     "critical timezone with -00:00 is accepted",
			input:    "2022-07-08T00:14:07-00:00[!Europe/London]",
			strict:   false,
			wantTime: time.Date(2022, 7, 8, 0, 14, 7, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location:         time.FixedZone("Europe/London", 1*3600),
				CriticalLocation: true,
			}),
		},
		{
			// A critical annotation MUST be processable (Section 3.3), so an
			// unknown name errors even in non-strict mode.
			name:    "critical unknown timezone errors in non-strict",
			input:   "2025-01-01T00:00:00Z[!Foo/Bar]",
			strict:  false,
			wantErr: "invalid timezone",
		},
		{
			// The grammar allows at most one time-zone annotation (RFC 9557
			// Section 4.1); a second one would overwrite the zone and its
			// critical flag, hiding a mandatory Section 3.4 inconsistency error.
			name:    "second timezone annotation rejected in non-strict",
			input:   "2022-07-08T00:14:07+00:00[!Europe/London][Asia/Tokyo]",
			strict:  false,
			wantErr: "invalid IXDTF suffix format",
		},
		{
			name:    "second timezone annotation rejected in strict mode",
			input:   "2025-01-01T00:00:00Z[Asia/Tokyo][Europe/Paris]",
			strict:  true,
			wantErr: "invalid IXDTF suffix format",
		},
		{
			name:     "timezone numeric offset suffix",
			input:    "2025-01-01T00:00:00Z[+09:00]",
			strict:   false,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: time.FixedZone("+09:00", 9*3600),
				Tags:     map[string]string{},
				Critical: map[string]bool{},
			}),
		},
		{
			// RFC 9557 Section 4.1 permits a critical numeric-offset
			// annotation; a matching offset is fully processable
			// (Section 3.3) and consistent (Section 3.4).
			name:     "critical matching numeric offset is accepted",
			input:    "2025-01-01T00:00:00+09:00[!+09:00]",
			strict:   false,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.FixedZone("+09:00", 9*3600)),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location:         time.FixedZone("+09:00", 9*3600),
				CriticalLocation: true,
			}),
		},
		{
			name:    "critical mismatching numeric offset errors in non-strict",
			input:   "2025-01-01T00:00:00+09:00[!+05:30]",
			strict:  false,
			wantErr: "timezone offset does not match",
		},
		{
			// An offset-derived FixedZone must not be resolved via the
			// timezone database, so a matching offset annotation is valid in
			// strict mode too.
			name:     "matching numeric offset is accepted in strict mode",
			input:    "2025-01-01T00:00:00+09:00[+09:00]",
			strict:   true,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.FixedZone("+09:00", 9*3600)),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Location: time.FixedZone("+09:00", 9*3600),
			}),
		},
		{
			name:     "invalid numeric offset ignored in non-strict mode",
			input:    "2025-01-01T00:00:00Z[+24:00]",
			strict:   false,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantExt:  ixdtf.NewIXDTFExtensions(nil),
		},
		{
			name:    "invalid numeric offset rejected in strict mode",
			input:   "2025-01-01T00:00:00Z[+24:00]",
			strict:  true,
			wantErr: "invalid timezone name",
		},
		{
			name:     "duplicate tag first wins",
			input:    "2025-01-01T00:00:00Z[key=one][key=two]",
			strict:   false,
			wantTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{Tags: map[string]string{
				"key": "one",
			}}),
		},
		{
			// RFC 9557 Section 3.3 example: a duplicate suffix key where the
			// first occurrence is critical MUST be treated as erroneous.
			name:    "critical duplicate tag rejected (first critical)",
			input:   "2022-07-08T00:14:07Z[!u-ca=chinese][u-ca=japanese]",
			strict:  false,
			wantErr: "critical extension cannot be processed",
		},
		{
			// RFC 9557 Section 3.3 example: the critical flag on the second
			// occurrence must not be silently discarded either.
			name:    "critical duplicate tag rejected (second critical)",
			input:   "2022-07-08T00:14:07Z[u-ca=chinese][!u-ca=japanese]",
			strict:  false,
			wantErr: "critical extension cannot be processed",
		},
		{
			name:    "critical duplicate tag rejected in strict mode",
			input:   "2022-07-08T00:14:07Z[!u-ca=chinese][u-ca=japanese]",
			strict:  true,
			wantErr: "critical extension cannot be processed",
		},
		{
			// RFC 9557 Section 3.3 example "[!knort=blargel]": in strict mode
			// this library is the recipient, understands only "u-ca", and MUST
			// treat an unrecognized critical key as erroneous.
			name:    "unknown critical suffix key rejected in strict mode",
			input:   "2022-07-08T00:14:07Z[!knort=blargel]",
			strict:  true,
			wantErr: "critical extension cannot be processed",
		},
		{
			// In non-strict mode processing of unrecognized critical keys is
			// delegated to the caller via the Critical map.
			name:     "unknown critical suffix key delegated in non-strict mode",
			input:    "2022-07-08T00:14:07Z[!knort=blargel]",
			strict:   false,
			wantTime: time.Date(2022, 7, 8, 0, 14, 7, 0, time.UTC),
			wantExt: ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
				Tags:     map[string]string{"knort": "blargel"},
				Critical: map[string]bool{"knort": true},
			}),
		},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
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

// TestParseUnknownLocalOffsetAppliesTimezone verifies RFC 9557 Section 3.4
// (Figure 2): when the RFC 3339 part uses an unknown local offset ("Z" or
// "-00:00"), the time-zone annotation is applied to resolve local time, and
// this is never treated as an inconsistency — in either strict or non-strict
// mode. See https://github.com/8beeeaaat/ixdtf/issues/22.
func TestParseUnknownLocalOffsetAppliesTimezone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		wantOffset int // seconds east of UTC after applying the timezone
		wantHour   int // local wall-clock hour after applying the timezone
	}{
		{
			name:       "Z with London in summer applies BST (+01:00)",
			input:      "2022-07-08T00:14:07Z[Europe/London]",
			wantOffset: 1 * 3600,
			wantHour:   1,
		},
		{
			name:       "Z with London in winter applies GMT (+00:00)",
			input:      "2022-01-08T00:14:07Z[Europe/London]",
			wantOffset: 0,
			wantHour:   0,
		},
		{
			name:       "negative-zero offset with London applies BST",
			input:      "2022-07-08T00:14:07-00:00[Europe/London]",
			wantOffset: 1 * 3600,
			wantHour:   1,
		},
		{
			name:       "Z with New York in summer applies EDT (-04:00)",
			input:      "2022-07-08T12:00:00Z[America/New_York]",
			wantOffset: -4 * 3600,
			wantHour:   8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, strict := range []bool{false, true} {
				got, ext, err := ixdtf.Parse(tc.input, strict)
				if err != nil {
					t.Fatalf("Parse(%q, %v) unexpected error: %v", tc.input, strict, err)
				}
				if _, off := got.Zone(); off != tc.wantOffset {
					t.Errorf("Parse(%q, %v) offset = %d, want %d", tc.input, strict, off, tc.wantOffset)
				}
				if got.Hour() != tc.wantHour {
					t.Errorf("Parse(%q, %v) local hour = %d, want %d", tc.input, strict, got.Hour(), tc.wantHour)
				}
				if ext.Location == nil {
					t.Errorf("Parse(%q, %v) expected location to be set", tc.input, strict)
				}
			}
		})
	}
}

// TestParseCriticalTimezone verifies RFC 9557 handling of a critical "!" flag
// on a time-zone annotation (Section 4.1 grammar; Section 3.3/3.4 semantics).
// See the follow-up to https://github.com/8beeeaaat/ixdtf/issues/22.
func TestParseCriticalTimezone(t *testing.T) {
	t.Parallel()

	t.Run("Z with critical zone is consistent in both modes", func(t *testing.T) {
		t.Parallel()
		const input = "2022-07-08T00:14:07Z[!Europe/London]" // Figure 2: no inconsistency
		for _, strict := range []bool{false, true} {
			got, ext, err := ixdtf.Parse(input, strict)
			if err != nil {
				t.Fatalf("Parse(%q, %v) unexpected error: %v", input, strict, err)
			}
			if !ext.CriticalLocation {
				t.Errorf("Parse(%q, %v) expected CriticalLocation to be true", input, strict)
			}
			if _, off := got.Zone(); off != 1*3600 { // London is BST (+01:00) in July
				t.Errorf("Parse(%q, %v) offset = %d, want %d", input, strict, off, 1*3600)
			}
		}
	})

	t.Run("critical zone round-trips with the ! flag", func(t *testing.T) {
		t.Parallel()
		const input = "2025-01-01T00:00:00+09:00[!Asia/Tokyo]"
		got, ext, err := ixdtf.Parse(input, true)
		if err != nil {
			t.Fatalf("Parse(%q) unexpected error: %v", input, err)
		}
		formatted, err := ixdtf.Format(got, ext)
		if err != nil {
			t.Fatalf("Format unexpected error: %v", err)
		}
		if formatted != input {
			t.Errorf("round trip failed: got %q, want %q", formatted, input)
		}
	})

	t.Run("critical inconsistency errors in non-strict mode", func(t *testing.T) {
		t.Parallel()
		// +00:00 asserts a zero offset that conflicts with London's +01:00 in
		// July; the critical flag forces the error even in non-strict mode.
		const input = "2022-07-08T00:14:07+00:00[!Europe/London]"
		_, _, err := ixdtf.Parse(input, false)
		checkParseError(t, err, input, false, "timezone offset does not match")
	})
}

// TestParseNumericZeroOffsetStillInconsistent verifies the contrast with Z: a
// concrete "+00:00" asserts a zero offset, so it conflicts with a timezone
// whose offset differs at that instant (RFC 9557 Section 3.4, Figure 1). This
// behavior must remain unchanged by the Z fix.
func TestParseNumericZeroOffsetStillInconsistent(t *testing.T) {
	t.Parallel()

	const input = "2022-07-08T00:14:07+00:00[Europe/London]" // London is BST (+01:00) in July

	// strict mode: the inconsistency is reported as an error.
	if _, _, err := ixdtf.Parse(input, true); err == nil {
		t.Fatalf("Parse(%q, true) expected inconsistency error, got nil", input)
	}

	// non-strict mode: the original +00:00 offset is preserved (timezone not applied).
	got, _, err := ixdtf.Parse(input, false)
	if err != nil {
		t.Fatalf("Parse(%q, false) unexpected error: %v", input, err)
	}
	if _, off := got.Zone(); off != 0 {
		t.Errorf("Parse(%q, false) offset = %d, want 0 (original preserved)", input, off)
	}
}
