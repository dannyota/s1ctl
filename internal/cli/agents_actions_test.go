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

func TestAgentParameterizedActionsDryRun(t *testing.T) {
	cases := [][]string{
		{"agents", "move-to-site", "A1", "--site-id", "S1"},
		{"agents", "set-external-id", "A1", "--external-id", "X1"},
		{"agents", "firewall-logging", "A1", "--state", "on"},
	}
	for _, args := range cases {
		out, err := runCLI(t, args...)
		if err != nil {
			t.Fatalf("%v: unexpected error: %v", args, err)
		}
		if !strings.Contains(out, "Would") {
			t.Fatalf("%v: expected dry-run message, got %q", args, out)
		}
	}
}

func TestAgentParameterizedActionsRequireFlags(t *testing.T) {
	cases := [][]string{
		{"agents", "move-to-site", "A1"},
		{"agents", "set-external-id", "A1"},
		{"agents", "firewall-logging", "A1"},
		{"agents", "firewall-logging", "A1", "--state", "sideways"},
	}
	for _, args := range cases {
		if _, err := runCLI(t, args...); err == nil {
			t.Fatalf("%v: expected flag validation error", args)
		}
	}
}
