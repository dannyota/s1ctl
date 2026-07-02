package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDatalakeIngestEventsDryRun(t *testing.T) {
	f := filepath.Join(t.TempDir(), "events.json")
	if err := os.WriteFile(f, []byte(`[{"ts":"1","attrs":{"message":"hi"}}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, "datalake", "ingest", "events", "--file", f, "--session", "s-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDatalakeIngestLogsDryRun(t *testing.T) {
	f := filepath.Join(t.TempDir(), "app.log")
	if err := os.WriteFile(f, []byte("line one\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, "datalake", "ingest", "logs", "--file", f, "--parser", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDatalakeIngestValidation(t *testing.T) {
	if _, err := runCLI(t, "datalake", "ingest", "events", "--session", "s"); err == nil || !strings.Contains(err.Error(), "--file is required") {
		t.Fatalf("expected --file validation error, got %v", err)
	}
	if _, err := runCLI(t, "datalake", "ingest", "logs"); err == nil || !strings.Contains(err.Error(), "--file is required") {
		t.Fatalf("expected --file validation error, got %v", err)
	}
}
