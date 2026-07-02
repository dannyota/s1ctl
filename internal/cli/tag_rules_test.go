package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTagRuleFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "rule.json")
	if err := os.WriteFile(path, []byte(`{"name":"Tag servers","conditions":{"op":"and"}}`), 0o600); err != nil {
		t.Fatalf("write tag rule file: %v", err)
	}
	return path
}

func TestTagRulesCreateDryRun(t *testing.T) {
	out, err := runCLI(t, "tag-rules", "create", "--from-file", writeTagRuleFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestTagRulesCreateRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "tag-rules", "create")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}

func TestTagRulesUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "tag-rules", "update", "TR1", "--from-file", writeTagRuleFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestTagRulesUpdateRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "tag-rules", "update", "TR1")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}

func TestTagRulesDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "tag-rules", "delete", "TR1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestTagRulesTestRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "tag-rules", "test")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}
