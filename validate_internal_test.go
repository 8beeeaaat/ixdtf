package ixdtf

import (
	"errors"
	"testing"
	"time"
)

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
		if err := validateLocationStrict(time.FixedZone("No/SuchZone", 0), true); !errors.Is(err, ErrInvalidTimezone) {
			t.Fatalf("expected strict validation to fail, got %v", err)
		}
	})
}
