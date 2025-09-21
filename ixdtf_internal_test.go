package ixdtf

import (
	"errors"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestAppendSuffix(t *testing.T) {
	t.Parallel()
	t.Run("appendSuffix skips invalid tag keys", func(t *testing.T) {
		t.Parallel()
		ext := NewIXDTFExtensions(nil)
		ext.Tags["InvalidKey"] = "1"
		ext.Tags["valid"] = "ok"
		ext.Critical["valid"] = true

		formatted := string(appendSuffix(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), ext, time.RFC3339))
		if strings.Contains(formatted, "InvalidKey") {
			t.Fatalf("expected invalid key to be skipped, got %q", formatted)
		}
		if !strings.Contains(formatted, "[!valid=ok]") {
			t.Fatalf("expected valid tag to be rendered, got %q", formatted)
		}
	})
}

func TestCheckTimeZoneConsistency(t *testing.T) {
	t.Parallel()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("fallback load", func(t *testing.T) {
		t.Parallel()
		timezoneCache.Delete("Asia/Tokyo")
		placeholder := time.FixedZone("Asia/Tokyo", 9*3600)
		res, err := checkTimeZoneConsistency(now, placeholder, false)
		if err != nil {
			t.Fatalf("expected fallback load to succeed, got %v", err)
		}
		if res.Location == nil || res.Location.String() != "Asia/Tokyo" {
			t.Fatalf("expected real location to be resolved, got %+v", res)
		}
	})

	t.Run("nil location", func(t *testing.T) {
		t.Parallel()
		if res, err := checkTimeZoneConsistency(now, nil, false); err != nil || !res.IsConsistent {
			t.Fatalf("expected nil location to be consistent, got res=%+v err=%v", res, err)
		}
	})

	t.Run("strict mode with unknown fixed zone", func(t *testing.T) {
		t.Parallel()
		if _, err := checkTimeZoneConsistency(now, time.FixedZone("+0900", 9*3600), true); err == nil {
			t.Fatalf("expected strict mode to fail for unknown fixed zone")
		}
	})
}

func TestFormatOffsetName(t *testing.T) {
	t.Parallel()
	t.Run("colon variant compression", func(t *testing.T) {
		t.Parallel()
		if got := formatOffsetName("+09:00"); got != "+0900" {
			t.Fatalf("expected colon variant to compress, got %q", got)
		}
	})

	t.Run("non-colon offset", func(t *testing.T) {
		t.Parallel()
		if got := formatOffsetName("+0900"); got != "+0900" {
			t.Fatalf("expected formatOffsetName to return original value for non-colon offset, got %q", got)
		}
	})
}

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

func TestParseNumericOffset(t *testing.T) {
	t.Parallel()
	t.Run("invalid minutes", func(t *testing.T) {
		t.Parallel()
		if _, err := parseNumericOffset("+09:61"); err == nil {
			t.Fatalf("expected parseNumericOffset to reject invalid minutes")
		}
	})

	t.Run("malformed sign", func(t *testing.T) {
		t.Parallel()
		if _, err := parseNumericOffset("a09:00"); err == nil {
			t.Fatalf("expected parseNumericOffset to reject malformed sign")
		}
	})

	t.Run("missing colon", func(t *testing.T) {
		t.Parallel()
		if _, err := parseNumericOffset("+2400"); err == nil {
			t.Fatalf("expected parseNumericOffset to reject missing colon")
		}
	})
}

func TestParseSuffix(t *testing.T) {
	t.Parallel()
	t.Run("critical timezone", func(t *testing.T) {
		t.Parallel()
		ext := NewIXDTFExtensions(nil)
		if err := parseSuffixElement("!Asia/Tokyo", ext, false); !errors.Is(err, ErrInvalidTimeZone) {
			t.Fatalf("expected ErrInvalidTimeZone for critical timezone, got %v", err)
		}
	})

	t.Run("missing brackets", func(t *testing.T) {
		t.Parallel()
		if _, err := parseSuffix("invalid", false); !errors.Is(err, ErrInvalidSuffix) {
			t.Fatalf("parseSuffix should fail for missing brackets, got %v", err)
		}
	})
}

func TestSuffixValidation(t *testing.T) {
	t.Parallel()
	t.Run("empty key range", func(t *testing.T) {
		t.Parallel()
		if err := isValidSuffixKeyRange("key=val", 4, 4); !errors.Is(err, ErrInvalidExtension) {
			t.Fatalf("expected ErrInvalidExtension for empty key range, got %v", err)
		}
	})

	t.Run("empty suffix value", func(t *testing.T) {
		t.Parallel()
		if err := isValidSuffixValue(""); err != nil {
			t.Fatalf("expected empty suffix value to be valid, got %v", err)
		}
	})

	t.Run("empty value range", func(t *testing.T) {
		t.Parallel()
		if err := isValidSuffixValueRange("key=val", 6, 6); !errors.Is(err, ErrInvalidExtension) {
			t.Fatalf("expected ErrInvalidExtension for empty value range, got %v", err)
		}
	})
}

func TestValidateInternal(t *testing.T) {
	t.Parallel()
	t.Run("nil extensions", func(t *testing.T) {
		t.Parallel()
		if err := validateExtensionsStrict(nil, false); err != nil {
			t.Fatalf("expected nil extensions to validate, got %v", err)
		}
	})

	t.Run("non-critical tags", func(t *testing.T) {
		t.Parallel()
		if err := validateCriticalTags(map[string]string{"k": "v"}, map[string]bool{"k": false}); err != nil {
			t.Fatalf("expected non-critical tags to pass, got %v", err)
		}
	})

	t.Run("non-strict location validation", func(t *testing.T) {
		t.Parallel()
		if err := validateLocationStrict(time.FixedZone("No/SuchZone", 0), false); err != nil {
			t.Fatalf("expected non-strict validation to ignore unknown zone, got %v", err)
		}
	})

	t.Run("strict location validation", func(t *testing.T) {
		t.Parallel()
		if err := validateLocationStrict(time.FixedZone("No/SuchZone", 0), true); !errors.Is(err, ErrInvalidTimeZone) {
			t.Fatalf("expected strict validation to fail, got %v", err)
		}
	})
}
