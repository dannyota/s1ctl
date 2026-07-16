package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAutomationListDryRun(t *testing.T) {
	// list is a read command; it requires an API client. Just verify
	// the command exists and parses flags without error.
	_, err := runCLI(t, "automation", "list", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAutomationGetRequiresArgs(t *testing.T) {
	_, err := runCLI(t, "automation", "get")
	if err == nil || !strings.Contains(err.Error(), "accepts 2 arg") {
		t.Fatalf("expected 2-arg error, got %v", err)
	}
}

func TestAutomationVersionsRequiresArg(t *testing.T) {
	_, err := runCLI(t, "automation", "versions")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestAutomationExportRequiresArgs(t *testing.T) {
	_, err := runCLI(t, "automation", "export")
	if err == nil || !strings.Contains(err.Error(), "accepts 2 arg") {
		t.Fatalf("expected 2-arg error, got %v", err)
	}
}

func TestAutomationCreateRequiresFile(t *testing.T) {
	_, err := runCLI(t, "automation", "create")
	if err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file required error, got %v", err)
	}
}

func TestAutomationCreateRequiresName(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "wf.json")
	if err := os.WriteFile(file, []byte(`{"description":"no name"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := runCLI(t, "automation", "create", "--from-file", file)
	if err == nil || !strings.Contains(err.Error(), "has no name") {
		t.Fatalf("expected no-name error, got %v", err)
	}
}

func TestAutomationCreateDryRun(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "wf.json")
	if err := os.WriteFile(file, []byte(`{"name":"Alert triage","actions":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, "automation", "create", "--from-file", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would import workflow") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
	if !strings.Contains(out, "Alert triage") {
		t.Fatalf("expected workflow name in message, got %q", out)
	}
}

func TestAutomationRunRequiresArgs(t *testing.T) {
	_, err := runCLI(t, "automation", "run")
	if err == nil || !strings.Contains(err.Error(), "accepts 2 arg") {
		t.Fatalf("expected 2-arg error, got %v", err)
	}
}

func TestAutomationRunDryRun(t *testing.T) {
	out, err := runCLI(t, "automation", "run", "wf-123", "v-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would run workflow wf-123") {
		t.Fatalf("expected dry-run message naming workflow, got %q", out)
	}
}

func TestAutomationActivateDryRun(t *testing.T) {
	out, err := runCLI(t, "automation", "activate", "wf-1", "v-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would activate workflow wf-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAutomationActivateRequiresArgs(t *testing.T) {
	_, err := runCLI(t, "automation", "activate")
	if err == nil || !strings.Contains(err.Error(), "accepts 2 arg") {
		t.Fatalf("expected 2-arg error, got %v", err)
	}
}

func TestAutomationDeactivateDryRun(t *testing.T) {
	out, err := runCLI(t, "automation", "deactivate", "wf-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would deactivate workflow wf-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAutomationDeactivateRequiresArg(t *testing.T) {
	_, err := runCLI(t, "automation", "deactivate")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestAutomationExecutionGetRequiresArg(t *testing.T) {
	_, err := runCLI(t, "automation", "execution-get")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestAutomationExecutionOutputRequiresArg(t *testing.T) {
	_, err := runCLI(t, "automation", "execution-output")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestAutomationSubcommandRequired(t *testing.T) {
	out, err := runCLI(t, "automation")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "automation [command]") {
		t.Fatalf("expected usage output, got %q", out)
	}
}
