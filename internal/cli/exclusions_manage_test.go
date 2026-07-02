package cli

import (
	"strings"
	"testing"
)

func TestExclusionsUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "exclusions", "update", "E1", "--type", "path", "--value", "/tmp/x", "--os-type", "linux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestExclusionsUpdateValidation(t *testing.T) {
	if _, err := runCLI(t, "exclusions", "update", "E1"); err == nil {
		t.Fatal("expected validation error without --type/--value/--os-type")
	}
}
