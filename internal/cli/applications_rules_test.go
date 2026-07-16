package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestAppControlRulesCreateDryRun(t *testing.T) {
	out, err := runCLI(t, "applications", "rules", "create",
		"--name", "Block malware",
		"--behavior", "block",
		"--os-type", "windows")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAppControlRulesCreateRequiresName(t *testing.T) {
	_, err := runCLI(t, "applications", "rules", "create",
		"--behavior", "block")
	if err == nil {
		t.Fatal("expected validation error without --name")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Fatalf("expected %q, got %q", "--name is required", err.Error())
	}
}

func TestAppControlRulesCreateRequiresBehavior(t *testing.T) {
	_, err := runCLI(t, "applications", "rules", "create",
		"--name", "Test")
	if err == nil {
		t.Fatal("expected validation error without --behavior")
	}
	if !strings.Contains(err.Error(), "--behavior is required") {
		t.Fatalf("expected %q, got %q", "--behavior is required", err.Error())
	}
}

func TestAppControlRulesUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "applications", "rules", "update", "12345",
		"--name", "Updated rule")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAppControlRulesUpdateRequiresField(t *testing.T) {
	_, err := runCLI(t, "applications", "rules", "update", "12345")
	if err == nil {
		t.Fatal("expected validation error without any field flag")
	}
	if !strings.Contains(err.Error(), "at least one of") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestAppControlRulesDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "applications", "rules", "delete", "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAppControlRulesDeleteMultiple(t *testing.T) {
	out, err := runCLI(t, "applications", "rules", "delete", "a", "b", "c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "3 application control rules") {
		t.Fatalf("expected plural count in dry-run, got %q", out)
	}
}

func TestAppControlRulesDeleteRequiresArgs(t *testing.T) {
	_, err := runCLI(t, "applications", "rules", "delete")
	if err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestAppControlRulesGetRequiresArg(t *testing.T) {
	_, err := runCLI(t, "applications", "rules", "get")
	if err == nil {
		t.Fatal("expected error with no arg")
	}
}

func TestAppControlSettingsUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "applications", "settings", "update",
		"--fallback-behavior", "block")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAppControlSettingsUpdateRequiresField(t *testing.T) {
	_, err := runCLI(t, "applications", "settings", "update")
	if err == nil {
		t.Fatal("expected error without any setting flag")
	}
	if !strings.Contains(err.Error(), "at least one of") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestAppControlRulesPushMissingDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nonexistent")
	_, err := runCLI(t, "applications", "rules", "push", "--dir", dir)
	if err == nil {
		t.Fatal("expected error for missing directory")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Fatalf("expected directory name in error, got %q", err.Error())
	}
}

func TestAppControlSettingsUpdateInvalidBool(t *testing.T) {
	_, err := runCLI(t, "applications", "settings", "update",
		"--enable", "maybe")
	if err == nil {
		t.Fatal("expected error for invalid --enable value")
	}
	if !strings.Contains(err.Error(), "invalid --enable value") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestAppMgmtSettingsUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "applications", "mgmt-settings", "update",
		"--extensive-scan", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAppMgmtSettingsUpdateRequiresField(t *testing.T) {
	_, err := runCLI(t, "applications", "mgmt-settings", "update")
	if err == nil {
		t.Fatal("expected error without any setting flag")
	}
	if !strings.Contains(err.Error(), "at least one of") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestAppMgmtSettingsUpdateInvalidBool(t *testing.T) {
	_, err := runCLI(t, "applications", "mgmt-settings", "update",
		"--extensive-scan", "nope")
	if err == nil {
		t.Fatal("expected error for invalid --extensive-scan value")
	}
	if !strings.Contains(err.Error(), "invalid --extensive-scan value") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestAppControlRulesPushScopeIdWithoutScopeType(t *testing.T) {
	dir := t.TempDir()
	_, err := runCLI(t, "applications", "rules", "push",
		"--dir", dir, "--scope-id", "000000", "--scope-type", "")
	if err == nil {
		t.Fatal("expected error when --scope-id given without --scope-type on push")
	}
	if !strings.Contains(err.Error(), "--scope-type is required") {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestAppControlRulesPullHasScopeTypeFlag(t *testing.T) {
	out, err := runCLI(t, "applications", "rules", "pull", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "--scope-type") {
		t.Fatalf("expected --scope-type in pull help, got %q", out)
	}
}

func TestAppControlRulesPushHasScopeTypeFlag(t *testing.T) {
	out, err := runCLI(t, "applications", "rules", "push", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "--scope-type") {
		t.Fatalf("expected --scope-type in push help, got %q", out)
	}
}

func TestAppControlSubcommandRegistration(t *testing.T) {
	// Verify the command tree is wired correctly.
	for _, args := range [][]string{
		{"applications", "rules", "--help"},
		{"applications", "settings", "--help"},
		{"applications", "labels", "--help"},
		{"applications", "rules", "pull", "--help"},
		{"applications", "rules", "push", "--help"},
		{"applications", "mgmt-settings", "--help"},
		{"applications", "mgmt-settings", "get", "--help"},
		{"applications", "mgmt-settings", "update", "--help"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			out, err := runCLI(t, args...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out == "" {
				t.Fatal("expected help output")
			}
		})
	}
}

// TestAppControlRulesUpdateMergePreservesFields runs update against a mock
// server and verifies the PUT body keeps every field the user did not change.
func TestAppControlRulesUpdateMergePreservesFields(t *testing.T) {
	var putBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			json.NewEncoder(w).Encode(map[string]any{
				"id":          "12345",
				"ruleName":    "Block unsigned apps",
				"description": "keep me",
				"behavior":    "MONITOR",
				"osType":      []string{"WINDOWS"},
				"propagation": true,
				"parameters":  map[string]any{"path": "C:\\apps\\*"},
			})
		case http.MethodPut:
			if err := json.NewDecoder(r.Body).Decode(&putBody); err != nil {
				t.Fatalf("decode PUT body: %v", err)
			}
			json.NewEncoder(w).Encode(map[string]any{"success": true})
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer srv.Close()
	t.Setenv("S1_CONSOLE_URL", srv.URL)
	t.Setenv("S1_TOKEN", "test-token")

	if _, err := runCLI(t, "applications", "rules", "update", "12345",
		"--behavior", "block", "--yes"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if putBody == nil {
		t.Fatal("PUT body never captured")
	}
	if got := putBody["behavior"]; got != "BLOCK" {
		t.Errorf("behavior = %v, want BLOCK", got)
	}
	if got := putBody["ruleName"]; got != "Block unsigned apps" {
		t.Errorf("ruleName = %v, want preserved original", got)
	}
	if got := putBody["description"]; got != "keep me" {
		t.Errorf("description = %v, want preserved original", got)
	}
	if got := putBody["propagation"]; got != true {
		t.Errorf("propagation = %v, want preserved true", got)
	}
	params, _ := putBody["parameters"].(map[string]any)
	if params == nil || params["path"] != "C:\\apps\\*" {
		t.Errorf("parameters = %v, want preserved path condition", putBody["parameters"])
	}
}
