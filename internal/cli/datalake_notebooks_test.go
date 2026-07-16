package cli

import (
	"strings"
	"testing"
)

func TestDatalakeNotebooksGetArgValidation(t *testing.T) {
	_, err := runCLI(t, "datalake", "notebooks", "get")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected arg validation error, got %v", err)
	}
}

func TestDatalakeNotebooksDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "datalake", "notebooks", "delete", "nb-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete notebook") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDatalakeNotebooksRequiresSubcommand(t *testing.T) {
	out, err := runCLI(t, "datalake", "notebooks")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Available Commands") {
		t.Fatalf("expected help output with subcommands, got %q", out)
	}
}
