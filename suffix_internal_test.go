package ixdtf

import (
	"errors"
	"testing"
)

func TestParseSuffix(t *testing.T) {
	t.Parallel()
	t.Run("critical timezone is accepted", func(t *testing.T) {
		t.Parallel()
		// RFC 9557 Section 4.1 permits a "!" flag on a time-zone annotation.
		ext := NewIXDTFExtensions(nil)
		if err := parseSuffixElement("!Asia/Tokyo", ext, false, &suffixParseState{}); err != nil {
			t.Fatalf("expected critical timezone to be accepted, got %v", err)
		}
		if ext.Location == nil || ext.Location.String() != "Asia/Tokyo" {
			t.Fatalf("expected Asia/Tokyo location, got %+v", ext.Location)
		}
		if !ext.CriticalLocation {
			t.Fatalf("expected CriticalLocation to be true")
		}
	})

	t.Run("critical unknown timezone is rejected in non-strict mode", func(t *testing.T) {
		t.Parallel()
		// A critical annotation MUST be processable (Section 3.3), so an
		// unknown name is an error even in non-strict mode.
		ext := NewIXDTFExtensions(nil)
		if err := parseSuffixElement("!Foo/Bar", ext, false, &suffixParseState{}); !errors.Is(err, ErrInvalidTimezone) {
			t.Fatalf("expected ErrInvalidTimezone for critical unknown timezone, got %v", err)
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
	t.Run("empty suffix value", func(t *testing.T) {
		t.Parallel()
		if err := isValidSuffixValue(""); err != nil {
			t.Fatalf("expected empty suffix value to be valid, got %v", err)
		}
	})

	t.Run("empty key", func(t *testing.T) {
		t.Parallel()
		ext := NewIXDTFExtensions(nil)
		if err := parseSuffixElement("=val", ext, false, &suffixParseState{}); !errors.Is(err, ErrInvalidExtension) {
			t.Fatalf("expected ErrInvalidExtension for empty key, got %v", err)
		}
	})

	t.Run("empty value", func(t *testing.T) {
		t.Parallel()
		ext := NewIXDTFExtensions(nil)
		if err := parseSuffixElement("key=", ext, false, &suffixParseState{}); !errors.Is(err, ErrInvalidExtension) {
			t.Fatalf("expected ErrInvalidExtension for empty value, got %v", err)
		}
	})
}
