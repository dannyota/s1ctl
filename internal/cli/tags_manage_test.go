package cli

import (
	"strings"
	"testing"
)

func TestTagsManageDryRun(t *testing.T) {
	cases := [][]string{
		{"tags", "create", "--key", "env", "--value", "prod"},
		{"tags", "update", "TG1", "--value", "dev"},
		{"tags", "delete", "TG1"},
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

func TestTagsManageValidation(t *testing.T) {
	cases := [][]string{
		{"tags", "create"},
		{"tags", "create", "--key", "env"},
		{"tags", "update", "TG1"},
		{"tags", "get"},
	}
	for _, args := range cases {
		if _, err := runCLI(t, args...); err == nil {
			t.Fatalf("%v: expected validation error", args)
		}
	}
}
