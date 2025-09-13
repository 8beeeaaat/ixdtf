package abnf_test

import (
	"testing"

	"github.com/8beeeaaat/ixdtf/abnf"
)

func TestAbnf_IsTimezoneNameSyntax(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// Valid timezone names
		{"simple zone", "UTC", true},
		{"zone with slash", "Asia/Tokyo", true},
		{"zone with underscore", "America/New_York", true},
		{"zone with multiple slashes", "America/North_Dakota/New_Salem", true},
		{"zone with numbers", "Etc/GMT", true},
		{"zone with plus", "Etc/GMT+5", true},
		{"zone with minus", "Etc/GMT-5", true},
		{"zone with dots", "Antarctica/DumontDUrville", true},
		{"zone with mixed chars", "A_b/C.d+e-1", true},
		{"underscore start", "_test/zone", true},
		{"letter start", "test", true},

		// Invalid timezone names
		{"empty string", "", false},
		{"starts with number", "1Asia/Tokyo", false},
		{"starts with slash", "/Tokyo", false},
		{"ends with slash", "Asia/", false},
		{"double slash", "Asia//Tokyo", false},
		{"only slash", "/", false},
		{"starts with dash", "-zone", false},
		{"starts with plus", "+zone", false},
		{"starts with dot", ".zone", false},
		{"only numbers", "123", false},
		{"invalid chars", "Asia/Tokyo@", false},
		{"space in name", "Asia Tokyo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := abnf.IsTimezoneNameSyntax(tt.input); got != tt.want {
				t.Errorf("IsTimezoneNameSyntax(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestAbnf_Validate(t *testing.T) {
	tests := []struct {
		name     string
		pat      *abnf.Abnf
		valids   []string
		invalids []string
	}{
		{
			name: "TimezoneName",
			pat:  abnf.AbnfTimezoneName,
			valids: []string{
				"Asia/Tokyo",
				"Europe/Paris",
				"America/New_York",
				"America/North_Dakota/New_Salem",
				"UTC",
				"Etc/GMT",
				"Etc/GMT-5",
				"Antarctica/DumontDUrville",
			},
			invalids: []string{
				"A/B", "Zone/SubZone", "Asia/Tokyo@", "A/B/C", "A_b/C.d+e-1", "A/", "/A", "A//B", "A/B/", "A//", "1A/B",
			},
		},
		{
			name: "Timezone",
			pat:  abnf.AbnfTimezone,
			valids: []string{
				"[Europe/London]",
				"[!Europe/Paris]",
				"[America/North_Dakota/New_Salem]",
				"[+09:00]",
				"[!-05:30]",
				"[UTC]",
				"[!UTC]",
				"[Etc/GMT]",
				"[!Etc/GMT-5]",
			},
			invalids: []string{
				"Europe/London",
				"[Europe//London]",
				"[!+9:00]",
				"[+0900]",
				"[++09:00]",
			},
		},
		{
			name:     "SuffixKey",
			pat:      abnf.AbnfSuffixKey,
			valids:   []string{"a", "a--", "ab", "a_b", "a-b", "a--b"},
			invalids: []string{"A", "-a", "a.", "", "_k9"},
		},
		{
			name:     "SuffixValues",
			pat:      abnf.AbnfSuffixValues,
			valids:   []string{"a", "abc-123", "A1-b2-C3"},
			invalids: []string{"-a", "a-", "a--b", "", "a_b"},
		},
		{
			name: "DateTimeExt",
			pat:  abnf.AbnfDateTimeExt,
			valids: []string{
				"2025-01-02T03:04:05Z",
				"2025-12-31T23:59:59+09:00",
				"2025-06-07T08:09:10.123456-05:30",
				"2025-06-07T08:09:10Z[Europe/Paris]",
				"2025-06-07T08:09:10Z[!Europe/Paris][key=value]",
				"2025-06-07T08:09:10.1+00:00[Abc_Def][!k1=v-1-2]",
			},
			invalids: []string{
				"2025-6-07T08:09:10Z",
				"2025-06-07 08:09:10Z",
				"2025-06-07T08:09:10",
				"2025-13-07T08:09:10Z",
				"2025-00-07T08:09:10Z",
				"2025-06-32T08:09:10Z",
				"2025-06-07T24:00:00Z",
				"2025-06-07T08:09:10+9:00",
				"2025-06-07T08:09:10+0900",
				"2025-06-07T08:09:10Z[Key=val]",
				"2025-06-07T08:09:10Zextra",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, v := range tc.valids {
				if err := tc.pat.Validate(v); err != nil {
					t.Errorf("expected valid %q for %s: %v", v, tc.name, err)
				}
			}
			for _, iv := range tc.invalids {
				if err := tc.pat.Validate(iv); err == nil {
					t.Errorf("expected invalid %q for %s", iv, tc.name)
				}
			}
		})
	}
}
