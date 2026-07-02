package cli

import (
	"strings"
	"testing"
)

func TestSitesManageDryRun(t *testing.T) {
	cases := [][]string{
		{"sites", "create", "--name", "n1", "--account-id", "AC1"},
		{"sites", "update", "S1", "--name", "n2"},
		{"sites", "delete", "S1"},
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

func TestSitesManageValidation(t *testing.T) {
	cases := [][]string{
		{"sites", "create"},
		{"sites", "create", "--name", "n1"},
		{"sites", "update", "S1"}, // no fields to change
	}
	for _, args := range cases {
		if _, err := runCLI(t, args...); err == nil {
			t.Fatalf("%v: expected validation error", args)
		}
	}
}
