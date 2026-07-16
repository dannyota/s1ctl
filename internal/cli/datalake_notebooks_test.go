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

func TestDatalakeNotebooksCreateMissingName(t *testing.T) {
	_, err := runCLI(t, "datalake", "notebooks", "create")
	if err == nil || !strings.Contains(err.Error(), "--name is required") {
		t.Fatalf("expected --name error, got %v", err)
	}
}

func TestDatalakeNotebooksCreateDryRun(t *testing.T) {
	out, err := runCLI(t, "datalake", "notebooks", "create", "--name", "test-nb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would create notebook") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDatalakeNotebooksUpdateArgValidation(t *testing.T) {
	_, err := runCLI(t, "datalake", "notebooks", "update")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected arg validation error, got %v", err)
	}
}

func TestDatalakeNotebooksUpdateNoFlags(t *testing.T) {
	_, err := runCLI(t, "datalake", "notebooks", "update", "nb-1")
	if err == nil || !strings.Contains(err.Error(), "at least one of --name or --description") {
		t.Fatalf("expected flag validation error, got %v", err)
	}
}

func TestDatalakeNotebooksUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "datalake", "notebooks", "update", "nb-1", "--name", "new-name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would update notebook") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}
