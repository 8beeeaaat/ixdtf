package ixdtf

import (
	"errors"
	"sort"
	"testing"
)

func TestParseErrorError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		pe       *ParseError
		expected string
	}{
		{
			name:     "with custom layout value",
			pe:       &ParseError{Err: errors.New("boom"), Layout: Layout("CUSTOM"), Value: "x"},
			expected: "IXDTFE parsing time \"x\" as \"CUSTOM\": boom",
		},
		{
			name:     "with underlying error and empty layout",
			pe:       &ParseError{Err: errors.New("parse fail"), Layout: "", Value: "foo"},
			expected: "IXDTFE parsing time \"foo\" as \"\": parse fail",
		},
		{
			name:     "with underlying error and known layout",
			pe:       &ParseError{Err: errors.New("invalid time"), Layout: LayoutRFC3339, Value: "bad"},
			expected: "IXDTFE parsing time \"bad\" as \"" + string(LayoutRFC3339) + "\": invalid time",
		},
		{
			name:     "nil underlying error returns empty string",
			pe:       &ParseError{Err: nil, Layout: LayoutRFC3339, Value: "2025-01-01T00:00:00Z"},
			expected: "",
		},
	}

	sort.Slice(tests, func(i, j int) bool { return tests[i].name < tests[j].name })

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.pe.Error(); got != tc.expected {
				t.Fatalf("ParseError.Error() = %q, want %q", got, tc.expected)
			}
		})
	}
}
