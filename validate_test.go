package ixdtf_test

import (
	"sort"
	"testing"

	"github.com/8beeeaaat/ixdtf"
)

func TestValidate(t *testing.T) {
	t.Parallel()
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
			// Validate must accept a critical zone exactly like Parse does
			// (RFC 9557 Section 3.4, Figure 2: Z carries no local offset).
			name:   "valid with critical timezone",
			input:  "2022-07-08T00:14:07Z[!Europe/London]",
			strict: false,
		},
		{
			name:   "valid with critical numeric offset",
			input:  "2025-01-01T00:00:00+09:00[!+09:00]",
			strict: false,
		},
		{
			name:   "valid with invalid calendar tag (non-critical)",
			input:  "2025-03-04T05:06:07Z[u-ca=hoge]",
			strict: false,
		},
		{
			name:    "invalid calendar tag (non-critical, strict)",
			input:   "2025-03-04T05:06:07Z[u-ca=hoge]",
			strict:  true,
			wantErr: "IXDTFE parsing time \"2025-03-04T05:06:07Z[u-ca=hoge]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid calendar tag identifier",
		},
		{
			name:    "invalid calendar tag (critical)",
			input:   "2025-03-04T05:06:07Z[!u-ca=hoge]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-03-04T05:06:07Z[!u-ca=hoge]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid calendar tag identifier",
		},
		{
			name:   "valid with timezone and tags",
			input:  "2025-06-07T08:09:10+01:00[Europe/Paris][!u-ca=gregory]",
			strict: false,
		},
		{
			name:   "invalid timezone name in non-strict mode",
			input:  "2025-01-01T00:00:00Z[No/SuchZone]",
			strict: false,
		},
		{
			name:    "invalid timezone name in strict mode",
			input:   "2025-01-01T00:00:00Z[No/SuchZone]",
			strict:  true,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[No/SuchZone]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid timezone name",
		},
		{
			name:   "invalid extension format in non-strict mode",
			input:  "2025-01-01T00:00:00Z[invalid]",
			strict: false,
		},
		{
			name:    "invalid extension format in strict mode",
			input:   "2025-01-01T00:00:00Z[invalid]",
			strict:  true,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[invalid]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid timezone name",
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
			name:   "timezone offset mismatch in non-strict mode",
			input:  "2025-06-01T12:00:00+09:00[America/New_York]",
			strict: false,
		},
		{
			name:    "timezone offset mismatch in strict mode",
			input:   "2025-06-01T12:00:00+09:00[America/New_York]",
			strict:  true,
			wantErr: "IXDTFE parsing time \"2025-06-01T12:00:00+09:00[America/New_York]\" as \"2006-01-02T15:04:05.999999999Z07:00*([time-zone-name][tags])\": timezone offset does not match the specified timezone",
		},
		{
			name:   "timezone content with u- prefix - non-strict mode",
			input:  "2025-01-01T00:00:00Z[u-invalid-timezone]",
			strict: false,
		},
		{
			name:   "timezone content with x- prefix - non-strict mode",
			input:  "2025-01-01T00:00:00Z[x-invalid-timezone]",
			strict: false,
		},
		{
			name:    "timezone content with u- prefix - strict mode",
			input:   "2025-01-01T00:00:00Z[u-invalid-timezone]",
			strict:  true,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[u-invalid-timezone]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid timezone name",
		},
		{
			name:    "timezone content with x- prefix - strict mode",
			input:   "2025-01-01T00:00:00Z[x-invalid-timezone]",
			strict:  true,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[x-invalid-timezone]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid timezone name",
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
			name:    "critical indicator without timezone content",
			input:   "2025-01-01T00:00:00Z[!]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[!]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid IXDTF suffix format",
		},
		{
			name:   "invalid numeric offset ignored in non-strict mode",
			input:  "2025-01-01T00:00:00Z[+24:00]",
			strict: false,
		},
		{
			name:    "invalid numeric offset rejected in strict mode",
			input:   "2025-01-01T00:00:00Z[+24:00]",
			strict:  true,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[+24:00]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid timezone name",
		},
		{
			name:    "invalid timezone syntax rejected by ABNF",
			input:   "2025-01-01T00:00:00Z[invalid@zone]",
			strict:  false,
			wantErr: "IXDTFE parsing time \"2025-01-01T00:00:00Z[invalid@zone]\" as \"2006-01-02T15:04:05Z07:00*([time-zone-name][tags])\": invalid extension format",
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

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ixdtf.Validate(tc.input, tc.strict)
			if (err != nil && err.Error() != tc.wantErr) || (err == nil && tc.wantErr != "") {
				t.Fatalf("Validate(%q, %v) error = %v, wantErr %v", tc.input, tc.strict, err, tc.wantErr)
			}
		})
	}
}
