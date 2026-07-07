package ixdtf

import (
	"testing"
	"time"
)

func TestCheckTimezoneConsistency(t *testing.T) {
	t.Parallel()
	now := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("fallback load", func(t *testing.T) {
		t.Parallel()
		timezoneCache.Delete("Asia/Tokyo")
		placeholder := time.FixedZone("Asia/Tokyo", 9*3600)
		res, err := checkTimezoneConsistency(now, placeholder, false, false)
		if err != nil {
			t.Fatalf("expected fallback load to succeed, got %v", err)
		}
		if res.Location == nil || res.Location.String() != "Asia/Tokyo" {
			t.Fatalf("expected real location to be resolved, got %+v", res)
		}
	})

	t.Run("nil location", func(t *testing.T) {
		t.Parallel()
		if res, err := checkTimezoneConsistency(now, nil, false, false); err != nil || !res.IsConsistent {
			t.Fatalf("expected nil location to be consistent, got res=%+v err=%v", res, err)
		}
	})

	t.Run("strict mode with unknown fixed zone", func(t *testing.T) {
		t.Parallel()
		if _, err := checkTimezoneConsistency(now, time.FixedZone("+0900", 9*3600), true, false); err == nil {
			t.Fatalf("expected strict mode to fail for unknown fixed zone")
		}
	})

	t.Run("unknown local offset is consistent even in strict mode", func(t *testing.T) {
		t.Parallel()
		london, err := time.LoadLocation("Europe/London")
		if err != nil {
			t.Skipf("Europe/London unavailable: %v", err)
		}
		// July: London is BST (+01:00), which differs from the Z offset (0),
		// yet an unknown local offset must never be flagged as inconsistent.
		summer := time.Date(2022, 7, 8, 0, 14, 7, 0, time.UTC)
		res, err := checkTimezoneConsistency(summer, london, true, true)
		if err != nil {
			t.Fatalf("unknown offset should not error in strict mode, got %v", err)
		}
		if !res.IsConsistent {
			t.Fatalf("unknown offset should be consistent, got %+v", res)
		}
	})
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

	t.Run("invalid hours", func(t *testing.T) {
		t.Parallel()
		if _, err := parseNumericOffset("+24:00"); err == nil {
			t.Fatalf("expected parseNumericOffset to reject invalid hours")
		}
	})
}

func TestIsOffsetLocationName(t *testing.T) {
	t.Parallel()
	cases := map[string]bool{
		"+09:00": true,
		"-03:30": true,
		"":       false,
		"+09:0":  false,
		"+0900":  false,
		"09:000": false,
		"+09:a0": false,
		"+09.00": false,
	}
	for name, want := range cases {
		if got := isOffsetLocationName(name); got != want {
			t.Errorf("isOffsetLocationName(%q) = %v, want %v", name, got, want)
		}
	}
}

// TestCheckTimezoneConsistencyUnloadableZone covers the resolution fallback
// for a location that is neither cached, offset-derived, nor loadable from
// the timezone database.
func TestCheckTimezoneConsistencyUnloadableZone(t *testing.T) {
	t.Parallel()
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	unknown := time.FixedZone("Not/AZone", 0)

	t.Run("non-strict treats unloadable zone as consistent", func(t *testing.T) {
		t.Parallel()
		result, err := checkTimezoneConsistency(ts, unknown, false, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.IsConsistent {
			t.Errorf("expected IsConsistent for unloadable zone in non-strict mode")
		}
	})

	t.Run("strict errors on unloadable zone", func(t *testing.T) {
		t.Parallel()
		if _, err := checkTimezoneConsistency(ts, unknown, true, false); err == nil {
			t.Fatalf("expected error for unloadable zone in strict mode")
		}
	})
}

func TestResolveLocationOffsetZone(t *testing.T) {
	t.Parallel()
	// A numeric-offset zone name is authoritative as-is and must not be
	// resolved through the timezone database.
	offset := time.FixedZone("+09:00", 9*3600)
	got, err := resolveLocation(offset)
	if err != nil {
		t.Fatalf("resolveLocation(+09:00) returned error: %v", err)
	}
	if got != offset {
		t.Fatalf("resolveLocation(+09:00) = %v, want the offset zone unchanged", got)
	}
}
