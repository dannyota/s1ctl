package cli

import (
	"strings"
	"testing"
)

func TestAgentPlainActionsDryRun(t *testing.T) {
	verbs := []string{
		"scan", "abort-scan", "decommission", "uninstall",
		"shutdown", "restart", "fetch-logs", "enable", "disable",
		"reset-config", "approve-uninstall", "reject-uninstall",
		"mark-up-to-date", "randomize-uuid",
	}
	for _, verb := range verbs {
		out, err := runCLI(t, "agents", verb, "A1")
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", verb, err)
		}
		if !strings.Contains(out, "Would") || !strings.Contains(out, "A1") {
			t.Fatalf("%s: expected dry-run message, got %q", verb, out)
		}
	}
}
