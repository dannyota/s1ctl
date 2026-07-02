package cli

import (
	"strings"
	"testing"
)

func TestDLPRulesEnableDryRun(t *testing.T) {
	out, err := runCLI(t, "dlp", "rules", "enable", "dlp-1", "dlp-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would enable 2 data protection rules") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDLPRulesDisableDryRun(t *testing.T) {
	out, err := runCLI(t, "dlp", "rules", "disable", "dlp-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would disable 1 data protection rule") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDLPRulesDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "dlp", "rules", "delete", "dlp-9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete 1 data protection rule") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

// TestDLPRulesActionRequiresID mirrors the SDK empty-ids guard: the CLI must
// reject enable/disable/delete with no rule IDs before any API call.
func TestDLPRulesActionRequiresID(t *testing.T) {
	for _, verb := range []string{"enable", "disable", "delete"} {
		_, err := runCLI(t, "dlp", "rules", verb)
		if err == nil || !strings.Contains(err.Error(), "at least 1") {
			t.Fatalf("%s: expected 'at least 1' arg error, got %v", verb, err)
		}
	}
}

func TestDLPRulesGetRequiresID(t *testing.T) {
	if _, err := runCLI(t, "dlp", "rules", "get"); err == nil {
		t.Fatal("expected error for missing rule id, got nil")
	}
}

func TestDLPRulesListRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "dlp", "rules", "list", "--scope-level", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}

func TestDLPClassificationsDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "dlp", "classifications", "delete", "c-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete dlp classification c-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDLPClassificationsDeleteRequiresID(t *testing.T) {
	if _, err := runCLI(t, "dlp", "classifications", "delete"); err == nil {
		t.Fatal("expected error for missing classification id, got nil")
	}
}

func TestDLPClassificationsGetRequiresID(t *testing.T) {
	if _, err := runCLI(t, "dlp", "classifications", "get"); err == nil {
		t.Fatal("expected error for missing classification id, got nil")
	}
}

func TestDLPSettingsRequiresScope(t *testing.T) {
	_, err := runCLI(t, "dlp", "settings")
	if err == nil || !strings.Contains(err.Error(), "required for dlp settings") {
		t.Fatalf("expected scope-required error, got %v", err)
	}
}

func TestDLPSettingsRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "dlp", "settings", "--scope-level", "bogus", "--scope-id", "x")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}

func TestDLPRequiresSubcommand(t *testing.T) {
	out, err := runCLI(t, "dlp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"rules", "classifications", "settings"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected help to list %q subcommand, got %q", sub, out)
		}
	}
}
