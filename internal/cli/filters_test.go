package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFilterFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "filter.json")
	if err := os.WriteFile(path, []byte(`{"name":"Infected","filterFields":{"infected":true}}`), 0o600); err != nil {
		t.Fatalf("write filter file: %v", err)
	}
	return path
}

func TestFiltersCreateDryRun(t *testing.T) {
	out, err := runCLI(t, "filters", "create", "--from-file", writeFilterFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestFiltersCreateRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "filters", "create")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}

func TestFiltersCreateMissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.json")
	_, err := runCLI(t, "filters", "create", "--from-file", missing)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("expected %q, got %q", "read "+missing, err.Error())
	}
}

func TestFiltersUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "filters", "update", "F1", "--from-file", writeFilterFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestFiltersUpdateRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "filters", "update", "F1")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}

func TestFiltersDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "filters", "delete", "F1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}
