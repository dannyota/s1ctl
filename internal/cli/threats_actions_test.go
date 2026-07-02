package cli

import (
	"strings"
	"testing"
)

func TestThreatPlainActionsDryRun(t *testing.T) {
	for _, verb := range []string{"blacklist", "fetch-file"} {
		out, err := runCLI(t, "threats", verb, "T1")
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", verb, err)
		}
		if !strings.Contains(out, "Would") || !strings.Contains(out, "T1") {
			t.Fatalf("%s: expected dry-run message, got %q", verb, out)
		}
	}
}

func TestThreatAddToExclusionsDryRun(t *testing.T) {
	out, err := runCLI(t, "threats", "add-to-exclusions", "T1", "--scope", "site", "--type", "hash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "T1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestThreatAddToExclusionsRequiresScopeAndType(t *testing.T) {
	if _, err := runCLI(t, "threats", "add-to-exclusions", "T1", "--type", "hash"); err == nil || !strings.Contains(err.Error(), "--scope is required") {
		t.Fatalf("expected --scope required error, got %v", err)
	}
	if _, err := runCLI(t, "threats", "add-to-exclusions", "T1", "--scope", "site"); err == nil || !strings.Contains(err.Error(), "--type is required") {
		t.Fatalf("expected --type required error, got %v", err)
	}
}

func TestThreatMitigateAlertsDryRun(t *testing.T) {
	out, err := runCLI(t, "threats", "mitigate-alerts", "--agent-id", "A1", "--storyline", "S1", "--action", "quarantine")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "A1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestThreatMitigateAlertsRequiresAgentAndStoryline(t *testing.T) {
	if _, err := runCLI(t, "threats", "mitigate-alerts", "--storyline", "S1"); err == nil || !strings.Contains(err.Error(), "--agent-id is required") {
		t.Fatalf("expected --agent-id required error, got %v", err)
	}
	if _, err := runCLI(t, "threats", "mitigate-alerts", "--agent-id", "A1"); err == nil || !strings.Contains(err.Error(), "--storyline is required") {
		t.Fatalf("expected --storyline required error, got %v", err)
	}
}

func TestThreatSetTicketDryRun(t *testing.T) {
	out, err := runCLI(t, "threats", "set-ticket", "T1", "--ticket-id", "JIRA-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "JIRA-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestThreatSetTicketRequiresTicketID(t *testing.T) {
	if _, err := runCLI(t, "threats", "set-ticket", "T1"); err == nil || !strings.Contains(err.Error(), "--ticket-id is required") {
		t.Fatalf("expected --ticket-id required error, got %v", err)
	}
}

// The read commands need a live client; offline we can only assert they are
// registered and validate positional arguments (error must not be
// "unknown command").
func TestThreatReadCommandsArgValidation(t *testing.T) {
	for _, verb := range []string{"quarantined-files", "exclusion-options"} {
		_, err := runCLI(t, "threats", verb)
		if err == nil {
			t.Fatalf("%s: expected missing-argument error", verb)
		}
		if strings.Contains(err.Error(), "unknown command") {
			t.Fatalf("%s: command not registered: %v", verb, err)
		}
	}
}
