package cli

import (
	"strings"
	"testing"
)

func TestAlertsNoteUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "alerts", "note-update", "note-1", "--text", "revised")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would update note note-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAlertsNoteUpdateRequiresText(t *testing.T) {
	_, err := runCLI(t, "alerts", "note-update", "note-1")
	if err == nil || !strings.Contains(err.Error(), "--text is required") {
		t.Fatalf("expected --text is required, got %v", err)
	}
}

func TestAlertsNoteDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "alerts", "note-delete", "note-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete note note-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAlertsNotesRequiresAlertID(t *testing.T) {
	if _, err := runCLI(t, "alerts", "notes"); err == nil {
		t.Fatal("expected error for missing alert id, got nil")
	}
}

func TestAlertsTimelineRequiresAlertID(t *testing.T) {
	if _, err := runCLI(t, "alerts", "timeline"); err == nil {
		t.Fatal("expected error for missing alert id, got nil")
	}
}

func TestAlertsCountsRequiresField(t *testing.T) {
	_, err := runCLI(t, "alerts", "counts")
	if err == nil || !strings.Contains(err.Error(), "--field is required") {
		t.Fatalf("expected --field is required, got %v", err)
	}
}

func TestAlertsExportRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "alerts", "export", "--scope-level", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}

func TestAlertsCountsRejectsBadScopeLevel(t *testing.T) {
	_, err := runCLI(t, "alerts", "counts", "--field", "severity", "--scope-level", "bogus")
	if err == nil || !strings.Contains(err.Error(), "invalid --scope-level") {
		t.Fatalf("expected invalid --scope-level, got %v", err)
	}
}
