package cli

import (
	"strings"
	"testing"
)

func TestUsersDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "users", "delete", "U1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}
