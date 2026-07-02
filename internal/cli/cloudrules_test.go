package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloudRulesEnableDryRun(t *testing.T) {
	out, err := runCLI(t, "cloud-rules", "enable", "rule-1", "rule-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would enable 2 cns rules") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestCloudRulesDisableDryRun(t *testing.T) {
	out, err := runCLI(t, "cloud-rules", "disable", "rule-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would disable 1 cns rule") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestCloudRulesDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "cloud-rules", "delete", "rule-9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete 1 cns rule") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

// TestCloudRulesActionRequiresID mirrors the SDK empty-ids guard: the CLI must
// reject enable/disable/delete with no rule IDs before any API call.
func TestCloudRulesActionRequiresID(t *testing.T) {
	for _, verb := range []string{"enable", "disable", "delete"} {
		_, err := runCLI(t, "cloud-rules", verb)
		if err == nil || !strings.Contains(err.Error(), "at least 1") {
			t.Fatalf("%s: expected 'at least 1' arg error, got %v", verb, err)
		}
	}
}

func TestCloudRulesCreateRequiresFile(t *testing.T) {
	_, err := runCLI(t, "cloud-rules", "create")
	if err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file is required, got %v", err)
	}
}

func TestCloudRulesUpdateRequiresFile(t *testing.T) {
	_, err := runCLI(t, "cloud-rules", "update", "rule-1")
	if err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file is required, got %v", err)
	}
}

func TestCloudRulesUpdateRequiresID(t *testing.T) {
	_, err := runCLI(t, "cloud-rules", "update", "--from-file", "x.json")
	if err == nil {
		t.Fatal("expected error for missing rule id, got nil")
	}
}

func TestCloudRulesCreateDryRun(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "rule.json")
	if err := os.WriteFile(file, []byte(`{"name":"r","queryType":"rego","severity":"HIGH","type":"CloudMisconfiguration"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, "cloud-rules", "create", "--from-file", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would create CNS rule") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestCloudRulesEvaluateRequiresResource(t *testing.T) {
	_, err := runCLI(t, "cloud-rules", "evaluate", "--query", "package s1")
	if err == nil || !strings.Contains(err.Error(), "--resource is required") {
		t.Fatalf("expected --resource is required, got %v", err)
	}
}

func TestCloudRulesEvaluateRequiresRuleOrQuery(t *testing.T) {
	_, err := runCLI(t, "cloud-rules", "evaluate", "--resource", "asset.json")
	if err == nil || !strings.Contains(err.Error(), "--rule or --query is required") {
		t.Fatalf("expected --rule or --query is required, got %v", err)
	}
}

func TestCloudRulesGetRequiresID(t *testing.T) {
	if _, err := runCLI(t, "cloud-rules", "get"); err == nil {
		t.Fatal("expected error for missing rule id, got nil")
	}
}

func TestCloudRulesListRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "cloud-rules", "list", "--scope-level", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}

func TestCloudRulesTypesRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "cloud-rules", "types", "--scope-level", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}
