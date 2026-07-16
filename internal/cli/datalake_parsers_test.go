package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDatalakeParsersGetArgValidation(t *testing.T) {
	_, err := runCLI(t, "datalake", "parsers", "get")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected arg validation error, got %v", err)
	}
}

func TestDatalakeParsersDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "datalake", "parsers", "delete", "p-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete parser") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDatalakeParsersRequiresSubcommand(t *testing.T) {
	out, err := runCLI(t, "datalake", "parsers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Available Commands") {
		t.Fatalf("expected help output with subcommands, got %q", out)
	}
}

func TestDatalakeParsersCreateMissingFile(t *testing.T) {
	_, err := runCLI(t, "datalake", "parsers", "create", "--name", "test")
	if err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file error, got %v", err)
	}
}

func TestDatalakeParsersCreateMissingName(t *testing.T) {
	_, err := runCLI(t, "datalake", "parsers", "create", "--from-file", "x.txt")
	if err == nil || !strings.Contains(err.Error(), "--name is required") {
		t.Fatalf("expected --name error, got %v", err)
	}
}

func TestDatalakeParsersCreateDryRun(t *testing.T) {
	f := filepath.Join(t.TempDir(), "p.txt")
	os.WriteFile(f, []byte("parser content"), 0o644)
	out, err := runCLI(t, "datalake", "parsers", "create", "--name", "test", "--from-file", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would create parser") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}
