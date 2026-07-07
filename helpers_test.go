package ixdtf_test

import (
	"strings"
	"testing"
	"time"

	"github.com/8beeeaaat/ixdtf"
)

func checkParseError(t *testing.T, err error, input string, strict bool, wantErr string) {
	if err == nil {
		t.Fatalf("Parse(%q, %v) expected error containing %q, got nil", input, strict, wantErr)
	}
	if !strings.Contains(err.Error(), wantErr) {
		t.Fatalf(
			"Parse(%q, %v) error = %q, want error containing %q",
			input,
			strict,
			err.Error(),
			wantErr,
		)
	}
}

func compareParseResults(
	t *testing.T,
	gotTime time.Time,
	gotExt *ixdtf.IXDTFExtensions,
	input string,
	strict bool,
	wantTime time.Time,
	wantExt *ixdtf.IXDTFExtensions,
) {
	if !gotTime.Equal(wantTime) {
		t.Errorf("Parse(%q, %v) time = %v, want %v", input, strict, gotTime, wantTime)
	}

	if !extensionsEqual(gotExt, wantExt) {
		t.Errorf("Parse(%q, %v) extensions = %+v, want %+v", input, strict, gotExt, wantExt)
	}
}

func extensionsEqual(a, b *ixdtf.IXDTFExtensions) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if (a.Location == nil) != (b.Location == nil) {
		return false
	}
	if a.Location != nil && b.Location != nil {
		if a.Location.String() != b.Location.String() {
			return false
		}
	}

	if a.CriticalLocation != b.CriticalLocation {
		return false
	}

	if len(a.Tags) != len(b.Tags) {
		return false
	}
	for k, v := range a.Tags {
		if b.Tags[k] != v {
			return false
		}
	}

	if len(a.Critical) != len(b.Critical) {
		return false
	}
	for k, v := range a.Critical {
		if b.Critical[k] != v {
			return false
		}
	}

	return true
}

func getTestTimezones() (*time.Location, *time.Location, *time.Location) {
	return time.FixedZone("Asia/Tokyo", 9*3600),
		time.FixedZone("Europe/Paris", 1*3600),
		time.FixedZone("CET", 1*3600)
}
