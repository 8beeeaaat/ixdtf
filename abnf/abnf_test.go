package abnf_test

import (
	"sort"
	"testing"

	"github.com/8beeeaaat/ixdtf/abnf"
)

func TestAbnf_IsTimeZoneSyntax(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"double slash", "Asia//Tokyo", false},
		{"empty string", "", false},
		{"ends with slash", "Asia/", false},
		{"invalid chars", "Asia/Tokyo@", false},
		{"letter start", "test", true},
		{"only numbers", "123", false},
		{"only slash", "/", false},
		{"simple zone", "UTC", true},
		{"space in name", "Asia Tokyo", false},
		{"starts with dash", "-zone", false},
		{"starts with dot", ".zone", false},
		{"starts with number", "1Asia/Tokyo", false},
		{"starts with plus", "+zone", false},
		{"starts with slash", "/Tokyo", false},
		{"underscore start", "_test/zone", true},
		{"zone with dots", "Antarctica/DumontDUrville", true},
		{"zone with minus", "Etc/GMT-5", true},
		{"zone with mixed chars", "A_b/C.d+e-1", true},
		{"zone with multiple slashes", "America/North_Dakota/New_Salem", true},
		{"zone with numbers", "Etc/GMT", true},
		{"zone with plus", "Etc/GMT+5", true},
		{"zone with slash", "Asia/Tokyo", true},
		{"zone with underscore", "America/New_York", true},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := abnf.IsTimeZoneSyntax(tc.input); got != tc.want {
				t.Errorf("IsTimeZoneSyntax(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestAbnf_ValidateDateTimeExt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		pat      *abnf.Abnf
		valids   []string
		invalids []string
	}{
		{
			name: "DateTimeExt",
			pat:  abnf.AbnfDateTimeExt,
			valids: []string{
				"2023-01-01T00:00:00Z",
				"2023-06-15T12:30:45+09:00",
				"2023-06-15T12:30:45-05:00",
				"2023-06-15T12:30:45.123+00:00",
				"2023-12-31T23:59:59.999999999Z",
			},
			invalids: []string{
				"",
				"2023-01-01T25:00:00Z",
				"2023-01-32T00:00:00Z",
				"2023-13-01T00:00:00Z",
				"invalid-format",
			},
		},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, v := range tc.valids {
				if err := tc.pat.ValidateDateTimeExt(v); err != nil {
					t.Errorf("expected valid %q for %s: %v", v, tc.name, err)
				}
			}
			for _, iv := range tc.invalids {
				if err := tc.pat.ValidateDateTimeExt(iv); err == nil {
					t.Errorf("expected invalid %q for %s", iv, tc.name)
				}
			}

			if err := abnf.AbnfSuffixKey.ValidateDateTimeExt("2023-01-01T00:00:00Z"); err == nil {
				t.Error("expected error for mismatched ABNF type in ValidateDateTimeExt")
			}
		})
	}
}

func TestAbnf_ValidateSuffixKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		pat      *abnf.Abnf
		valids   []string
		invalids []string
	}{
		{
			name:     "SuffixKey",
			pat:      abnf.AbnfSuffixKey,
			valids:   []string{"a", "a--", "a-b", "a--b", "a_b", "ab"},
			invalids: []string{"", "-a", "A", "a.", "_k9", "x-private"},
		},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, v := range tc.valids {
				if err := tc.pat.ValidateSuffixKey(v); err != nil {
					t.Errorf("expected valid %q for %s: %v", v, tc.name, err)
				}
			}
			for _, iv := range tc.invalids {
				if err := tc.pat.ValidateSuffixKey(iv); err == nil {
					t.Errorf("expected invalid %q for %s", iv, tc.name)
				}
			}

			if err := abnf.AbnfDateTimeExt.ValidateSuffixKey("demo"); err == nil {
				t.Error("expected error for mismatched ABNF type in ValidateSuffixKey")
			}
		})
	}
}

func TestAbnf_ValidateSuffixValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		pat      *abnf.Abnf
		valids   []string
		invalids []string
	}{
		{
			name: "SuffixValues",
			pat:  abnf.AbnfSuffixValues,
			valids: []string{
				"A1B2C3",
				"a",
				"gregorian",
				"japanese",
				"val-ue",
			},
			invalids: []string{
				"",
				"-val",
				"val.",
				"val@",
				"val_ue",
			},
		},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, v := range tc.valids {
				if err := tc.pat.ValidateSuffixValues(v); err != nil {
					t.Errorf("expected valid %q for %s: %v", v, tc.name, err)
				}
			}
			for _, iv := range tc.invalids {
				if err := tc.pat.ValidateSuffixValues(iv); err == nil {
					t.Errorf("expected invalid %q for %s", iv, tc.name)
				}
			}

			if err := abnf.AbnfDateTimeExt.ValidateSuffixValues("value"); err == nil {
				t.Error("expected error for mismatched ABNF type in ValidateSuffixValues")
			}
		})
	}
}

func TestAbnf_ValidateTimeZone(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		pat      *abnf.Abnf
		strict   bool
		valids   []string
		invalids []string
	}{
		{
			name:   "non-strict mode",
			pat:    abnf.AbnfTimeZone,
			strict: false,
			valids: []string{
				"A/B",
				"A/B/C",
				"A_b/C.d+e-1",
				"America/New_York",
				"America/North_Dakota/New_Salem",
				"Antarctica/DumontDUrville",
				"Asia/Tokyo",
				"Etc/GMT",
				"Etc/GMT-5",
				"Europe/Paris",
				"UTC",
				"Zone/SubZone",
			},
			invalids: []string{"A/", "A//", "A//B", "A/B/", "Asia/Tokyo@", "/A", "1A/B"},
		},
		{
			name:   "strict mode",
			pat:    abnf.AbnfTimeZone,
			strict: true,
			valids: []string{
				"America/New_York",
				"America/North_Dakota/New_Salem",
				"Antarctica/DumontDUrville",
				"Asia/Tokyo",
				"Etc/GMT",
				"Etc/GMT-5",
				"Europe/Paris",
				"UTC",
			},
			invalids: []string{
				"A/B",
				"A/B/C",
				"A_b/C.d+e-1",
				"A/",
				"A//",
				"A//B",
				"A/B/",
				"Zone/SubZone",
				"Asia/Tokyo@",
				"/A",
				"1A/B",
			},
		},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, v := range tc.valids {
				if err := tc.pat.ValidateTimeZone(v, tc.strict); err != nil {
					t.Errorf("expected valid %q for %s: %v", v, tc.name, err)
				}
			}
			for _, iv := range tc.invalids {
				if err := tc.pat.ValidateTimeZone(iv, tc.strict); err == nil {
					t.Errorf("expected invalid %q for %s", iv, tc.name)
				}
			}

			if err := abnf.AbnfSuffixKey.ValidateTimeZone("Asia/Tokyo", tc.strict); err == nil {
				t.Error("expected error for mismatched ABNF type in ValidateTimeZone")
			}
		})
	}
}

func TestAbnf_ValidateTimeZoneTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		pat      *abnf.Abnf
		strict   bool
		valids   []string
		invalids []string
	}{
		{
			name:     "TimeZoneTag non-strict",
			pat:      abnf.AbnfTimeZoneTag,
			strict:   false,
			valids:   []string{"[!America/New_York]", "[-05:00]", "[+09:00]", "[Asia/Tokyo]", "[UTC]"},
			invalids: []string{"", "[=value]", "[TAG=VALUE]", "invalid", "[tag@value]", "[tag=]", "u-ca=gregorian"},
		},
		{
			name:   "TimeZoneTag strict",
			pat:    abnf.AbnfTimeZoneTag,
			strict: true,
			valids: []string{"[!America/New_York]", "[-05:00]", "[+09:00]", "[Asia/Tokyo]", "[UTC]"},
			invalids: []string{
				"",
				"[=value]",
				"[TAG=VALUE]",
				"invalid",
				"[tag@value]",
				"[tag=]",
				"[unknown/timezone]",
				"u-ca=gregorian",
			},
		},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, v := range tc.valids {
				if err := tc.pat.ValidateTimeZoneTag(v, tc.strict); err != nil {
					t.Errorf("expected valid %q for %s: %v", v, tc.name, err)
				}
			}
			for _, iv := range tc.invalids {
				if err := tc.pat.ValidateTimeZoneTag(iv, tc.strict); err == nil {
					t.Errorf("expected invalid %q for %s", iv, tc.name)
				}
			}

			if err := abnf.AbnfSuffixKey.ValidateTimeZoneTag("[Asia/Tokyo]", tc.strict); err == nil {
				t.Error("expected error for mismatched ABNF type in ValidateTimeZoneTag")
			}
		})
	}
}
