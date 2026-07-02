package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloudPoliciesPushDryRun(t *testing.T) {
	f := filepath.Join(t.TempDir(), "cloud-policies.json")
	payload := `[{"id":"P1","status":"enabled"},{"id":"P2","status":"disabled"}]`
	if err := os.WriteFile(f, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, "cloud-policies", "push", "--file", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}
