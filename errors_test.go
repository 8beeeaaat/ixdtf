package ixdtf_test

import (
	"errors"
	"testing"

	"github.com/8beeeaaat/ixdtf"
)

// TestParseErrorUnwrap verifies that ParseError supports errors.Is matching
// against the package's sentinel errors via Unwrap.
func TestParseErrorUnwrap(t *testing.T) {
	t.Parallel()

	_, _, err := ixdtf.Parse("2025-06-01T12:00:00+09:00[America/New_York]", true)
	if !errors.Is(err, ixdtf.ErrTimezoneOffsetMismatch) {
		t.Fatalf("expected errors.Is(err, ErrTimezoneOffsetMismatch), got %v", err)
	}

	_, _, err = ixdtf.Parse("2022-07-08T00:14:07Z[Asia/Tokyo][Europe/Paris]", true)
	if !errors.Is(err, ixdtf.ErrInvalidSuffix) {
		t.Fatalf("expected errors.Is(err, ErrInvalidSuffix), got %v", err)
	}
}
