package cli

import (
	"strings"
	"testing"
)

func TestFirewallPushHelpMentionsSiteID(t *testing.T) {
	out, err := runCLI(t, "firewall", "push", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Always pass --site-id") {
		t.Fatalf("expected site-id caveat in push help, got %q", out)
	}
}
