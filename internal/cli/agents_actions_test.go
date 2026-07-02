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

func TestAgentExtraPlainActionsDryRun(t *testing.T) {
	verbs := []string{"reset-passphrase", "fetch-installed-apps", "fetch-firewall-rules"}
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

func TestAgentExtraParameterizedActionsDryRun(t *testing.T) {
	cases := [][]string{
		{"agents", "broadcast", "A1", "--message", "heads up"},
		{"agents", "fetch-files", "A1", "--path", "/etc/hosts"},
		{"agents", "fetch-files", "A1", "--path", "/etc/hosts", "--password", "pw-placeholder"},
		{"agents", "ranger", "A1", "--state", "on"},
		{"agents", "ranger", "A1", "--state", "off"},
		{"agents", "local-upgrade", "A1", "--state", "off"},
		{"agents", "local-upgrade", "A1", "--state", "on", "--until", "2030-01-01T00:00:00Z"},
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

// TestAgentFetchFilesPasswordNotInDryRun confirms the file password never
// appears in the guard action string (which is what the audit log records).
func TestAgentFetchFilesPasswordNotInDryRun(t *testing.T) {
	out, err := runCLI(t, "agents", "fetch-files", "A1", "--path", "/etc/hosts", "--password", "pw-placeholder")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "pw-placeholder") {
		t.Fatalf("password leaked into dry-run/action output: %q", out)
	}
}

func TestAgentExtraActionsValidation(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"agents", "broadcast", "A1"}, "--message is required"},
		{[]string{"agents", "fetch-files", "A1"}, "--path is required"},
		{[]string{"agents", "ranger", "A1"}, `--state must be "on" or "off"`},
		{[]string{"agents", "ranger", "A1", "--state", "sideways"}, `--state must be "on" or "off"`},
		{[]string{"agents", "local-upgrade", "A1"}, `--state must be "on" or "off"`},
		{[]string{"agents", "local-upgrade", "A1", "--state", "on"}, "--until is required when --state is on"},
	}
	for _, tc := range cases {
		_, err := runCLI(t, tc.args...)
		if err == nil {
			t.Fatalf("%v: expected validation error", tc.args)
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("%v: expected error %q, got %v", tc.args, tc.want, err)
		}
	}
}

// TestAgentReadsRequireArg confirms the read commands validate their argument
// and are not gated by the mutation guard (no "Would" wording). Output is
// exercised via SDK tests to avoid a live client.
func TestAgentReadsRequireArg(t *testing.T) {
	if _, err := runCLI(t, "agents", "local-upgrade-status"); err == nil {
		t.Fatal("expected arg validation error for local-upgrade-status without id")
	}
}
