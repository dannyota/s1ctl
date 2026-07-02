package cli

import (
	"strings"
	"testing"
)

func TestUpgradePolicyLifecycleDryRun(t *testing.T) {
	cases := [][]string{
		{"upgrade-policies", "update", "UP1", "--name", "n", "--os-type", "linux", "--scope-level", "site", "--file-id", "F1"},
		{"upgrade-policies", "activate", "UP1"},
		{"upgrade-policies", "deactivate", "UP1"},
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
