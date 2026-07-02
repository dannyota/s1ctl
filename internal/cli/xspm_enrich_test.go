package cli

import (
	"strings"
	"testing"
)

// Misconfigurations enrichment CLI.

func TestMisconfigurationsNoteAddDryRun(t *testing.T) {
	out, err := runCLI(t, "misconfigurations", "add-note", "m-1", "--text", "check this")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would add note to misconfiguration m-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMisconfigurationsNoteAddRequiresText(t *testing.T) {
	_, err := runCLI(t, "misconfigurations", "add-note", "m-1")
	if err == nil || !strings.Contains(err.Error(), "--text is required") {
		t.Fatalf("expected --text is required, got %v", err)
	}
}

func TestMisconfigurationsNoteUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "misconfigurations", "update-note", "note-1", "--text", "revised")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would update note note-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMisconfigurationsNoteUpdateRequiresText(t *testing.T) {
	_, err := runCLI(t, "misconfigurations", "update-note", "note-1")
	if err == nil || !strings.Contains(err.Error(), "--text is required") {
		t.Fatalf("expected --text is required, got %v", err)
	}
}

func TestMisconfigurationsNoteDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "misconfigurations", "delete-note", "note-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete note note-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMisconfigurationsAssignDryRun(t *testing.T) {
	out, err := runCLI(t, "misconfigurations", "assign", "m-1", "--user-id", "u-9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would assign misconfiguration m-1 to user u-9") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMisconfigurationsAssignRequiresUserID(t *testing.T) {
	_, err := runCLI(t, "misconfigurations", "assign", "m-1")
	if err == nil || !strings.Contains(err.Error(), "--user-id is required") {
		t.Fatalf("expected --user-id is required, got %v", err)
	}
}

func TestMisconfigurationsNotesRequiresID(t *testing.T) {
	if _, err := runCLI(t, "misconfigurations", "notes"); err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
}

func TestMisconfigurationsHistoryRequiresID(t *testing.T) {
	if _, err := runCLI(t, "misconfigurations", "history"); err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
}

func TestMisconfigurationsRelatedAssetsRequiresID(t *testing.T) {
	if _, err := runCLI(t, "misconfigurations", "related-assets"); err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
}

func TestMisconfigurationsExportRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "misconfigurations", "export", "--scope-level", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}

// Vulnerabilities enrichment CLI.

func TestVulnerabilitiesNoteAddDryRun(t *testing.T) {
	out, err := runCLI(t, "vulnerabilities", "add-note", "v-1", "--text", "patch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would add note to vulnerability v-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestVulnerabilitiesNoteAddRequiresText(t *testing.T) {
	_, err := runCLI(t, "vulnerabilities", "add-note", "v-1")
	if err == nil || !strings.Contains(err.Error(), "--text is required") {
		t.Fatalf("expected --text is required, got %v", err)
	}
}

func TestVulnerabilitiesNoteUpdateRequiresText(t *testing.T) {
	_, err := runCLI(t, "vulnerabilities", "update-note", "note-1")
	if err == nil || !strings.Contains(err.Error(), "--text is required") {
		t.Fatalf("expected --text is required, got %v", err)
	}
}

func TestVulnerabilitiesNoteDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "vulnerabilities", "delete-note", "note-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete note note-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestVulnerabilitiesAssignDryRun(t *testing.T) {
	out, err := runCLI(t, "vulnerabilities", "assign", "v-1", "--user-id", "u-9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would assign vulnerability v-1 to user u-9") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestVulnerabilitiesAssignRequiresUserID(t *testing.T) {
	_, err := runCLI(t, "vulnerabilities", "assign", "v-1")
	if err == nil || !strings.Contains(err.Error(), "--user-id is required") {
		t.Fatalf("expected --user-id is required, got %v", err)
	}
}

func TestVulnerabilitiesCveRequiresID(t *testing.T) {
	if _, err := runCLI(t, "vulnerabilities", "cve"); err == nil {
		t.Fatal("expected error for missing cve id, got nil")
	}
}

func TestVulnerabilitiesHistoryRequiresID(t *testing.T) {
	if _, err := runCLI(t, "vulnerabilities", "history"); err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
}

func TestVulnerabilitiesRelatedAssetsRequiresID(t *testing.T) {
	if _, err := runCLI(t, "vulnerabilities", "related-assets"); err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
}

func TestVulnerabilitiesStatsRejectsBadTop(t *testing.T) {
	_, err := runCLI(t, "vulnerabilities", "stats", "--top", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --top") {
		t.Fatalf("expected invalid --top, got %v", err)
	}
}

func TestVulnerabilitiesStatsRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "vulnerabilities", "stats", "--scope-level", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}

func TestVulnerabilitiesExportRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "vulnerabilities", "export", "--scope-level", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}
