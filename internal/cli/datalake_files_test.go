package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDatalakeFilesPutDryRun(t *testing.T) {
	f := filepath.Join(t.TempDir(), "conf.json")
	if err := os.WriteFile(f, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, "datalake", "files", "put", "/config/x", "--from-file", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
	out, err = runCLI(t, "datalake", "files", "put", "/config/x", "--delete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDatalakeFilesPutValidation(t *testing.T) {
	if _, err := runCLI(t, "datalake", "files", "put", "/config/x"); err == nil || !strings.Contains(err.Error(), "exactly one of --from-file or --delete") {
		t.Fatalf("expected put validation error, got %v", err)
	}
}
