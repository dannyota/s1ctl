package cli

import (
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
