package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeJSONFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestRemoteOpsUpdateDryRun(t *testing.T) {
	f := writeJSONFile(t, "script.json", `{"scriptName":"Collect Logs","scriptType":"dataCollection","osTypes":["linux"],"inputRequired":false,"inputExample":"-","inputInstructions":"-","scriptRuntimeTimeoutSeconds":3600}`)
	out, err := runCLI(t, "remoteops", "update", "S1", "--from-file", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestRemoteOpsUpdateRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "remoteops", "update", "S1")
	if err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file required, got %v", err)
	}
}

func TestRemoteOpsUpdateMissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.json")
	_, err := runCLI(t, "remoteops", "update", "S1", "--from-file", missing)
	if err == nil || !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("expected read error, got %v", err)
	}
}

func TestRemoteOpsContentRequiresArg(t *testing.T) {
	_, err := runCLI(t, "remoteops", "content")
	if err == nil {
		t.Fatal("expected error without script-id arg")
	}
}

func TestRemoteOpsUploadLimitsNoArgs(t *testing.T) {
	_, err := runCLI(t, "remoteops", "upload-limits", "extra")
	if err == nil {
		t.Fatal("expected error with unexpected positional arg")
	}
}

func TestRemoteOpsPendingApproveDryRun(t *testing.T) {
	out, err := runCLI(t, "remoteops", "pending", "approve", "P1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "approve") {
		t.Fatalf("expected approve dry-run message, got %q", out)
	}
}

func TestRemoteOpsPendingDeclineDryRun(t *testing.T) {
	out, err := runCLI(t, "remoteops", "pending", "decline", "P1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "decline") {
		t.Fatalf("expected decline dry-run message, got %q", out)
	}
}

func TestRemoteOpsPendingApproveRequiresArg(t *testing.T) {
	_, err := runCLI(t, "remoteops", "pending", "approve")
	if err == nil {
		t.Fatal("expected error without pending-execution-id arg")
	}
}

func TestRemoteOpsPendingRequiresSubcommand(t *testing.T) {
	out, err := runCLI(t, "remoteops", "pending")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "approve") || !strings.Contains(out, "list") {
		t.Fatalf("expected help listing subcommands, got %q", out)
	}
}

func TestRemoteOpsGuardrailsGetRequiresScope(t *testing.T) {
	_, err := runCLI(t, "remoteops", "guardrails", "get")
	if err == nil || !strings.Contains(err.Error(), "--scope-id is required") {
		t.Fatalf("expected scope-id required, got %v", err)
	}
	_, err = runCLI(t, "remoteops", "guardrails", "get", "--scope-id", "000000000000000000")
	if err == nil || !strings.Contains(err.Error(), "--scope-level is required") {
		t.Fatalf("expected scope-level required, got %v", err)
	}
}

func TestRemoteOpsGuardrailsSetDryRun(t *testing.T) {
	f := writeJSONFile(t, "gr.json", `{"scopeId":"000000000000000000","scopeLevel":"site","endpointsQuantity":100,"scriptTypes":["action"],"enabled":true}`)
	out, err := runCLI(t, "remoteops", "guardrails", "set", "--from-file", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "set guardrail") {
		t.Fatalf("expected set dry-run message, got %q", out)
	}
}

func TestRemoteOpsGuardrailsSetRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "remoteops", "guardrails", "set")
	if err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file required, got %v", err)
	}
}

func TestRemoteOpsGuardrailsDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "remoteops", "guardrails", "delete", "--scope-id", "000000000000000000", "--scope-level", "site")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "delete guardrail") {
		t.Fatalf("expected delete dry-run message, got %q", out)
	}
}

func TestRemoteOpsGuardrailsDeleteRequiresScope(t *testing.T) {
	_, err := runCLI(t, "remoteops", "guardrails", "delete")
	if err == nil || !strings.Contains(err.Error(), "--scope-id is required") {
		t.Fatalf("expected scope-id required, got %v", err)
	}
}

func TestRemoteOpsGuardrailsCheckRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "remoteops", "guardrails", "check")
	if err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file required, got %v", err)
	}
}
