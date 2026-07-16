package cli

import (
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
		t.Fatal("expected validation error without --name or --behavior")
	}
	if !strings.Contains(err.Error(), "at least --name or --behavior is required") {
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

func TestAppControlSubcommandRegistration(t *testing.T) {
	// Verify the command tree is wired correctly.
	for _, args := range [][]string{
		{"applications", "rules", "--help"},
		{"applications", "settings", "--help"},
		{"applications", "labels", "--help"},
		{"applications", "rules", "pull", "--help"},
		{"applications", "rules", "push", "--help"},
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
