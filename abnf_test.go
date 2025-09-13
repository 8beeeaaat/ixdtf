package ixdtf_test

import (
	"testing"

	"github.com/8beeeaaat/ixdtf"
)

func TestAbnf_Validate(t *testing.T) {
	tests := []struct {
		name     string
		pat      *ixdtf.Abnf
		valids   []string
		invalids []string
	}{
		{
			name: "TimezoneName",
			pat:  ixdtf.AbnfTimezoneName,
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
				"A/B", "Zone/SubZone", "A/B/C", "A_b/C.d+e-1", "A/", "/A", "A//B", "A/B/", "A//", "1A/B",
			},
		},
		{
			name: "Timezone",
			pat:  ixdtf.AbnfTimezone,
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
			pat:      ixdtf.AbnfSuffixKey,
			valids:   []string{"a", "a--", "ab", "a_b", "a-b", "a--b"},
			invalids: []string{"A", "-a", "a.", "", "_k9"},
		},
		{
			name:     "SuffixValues",
			pat:      ixdtf.AbnfSuffixValues,
			valids:   []string{"a", "abc-123", "A1-b2-C3"},
			invalids: []string{"-a", "a-", "a--b", "", "a_b"},
		},
		{
			name: "DateTimeExt",
			pat:  ixdtf.AbnfDateTimeExt,
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
