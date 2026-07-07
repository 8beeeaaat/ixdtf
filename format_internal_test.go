package ixdtf

import (
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
