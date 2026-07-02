package cli

import (
	"strings"
	"testing"
)

func TestGroupsUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "groups", "update", "G1", "--name", "n2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestGroupsUpdateValidation(t *testing.T) {
	if _, err := runCLI(t, "groups", "update", "G1"); err == nil {
		t.Fatal("expected validation error with no field flags")
	}
}
