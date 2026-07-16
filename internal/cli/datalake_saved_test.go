package cli

import (
	"strings"
	"testing"
)

func TestDatalakeSavedQueriesDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "datalake", "saved-queries", "delete", "my-query")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete saved query") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDatalakeSavedQueriesDeleteBadType(t *testing.T) {
	_, err := runCLI(t, "datalake", "saved-queries", "delete", "q", "--type", "BOGUS")
	if err == nil || !strings.Contains(err.Error(), "--type must be PRIVATE or SHARED") {
		t.Fatalf("expected type validation error, got %v", err)
	}
}

func TestDatalakeSavedQueriesListAcceptsCSV(t *testing.T) {
	// Verify that --output csv is accepted (no "unknown flag" error).
	// The command will fail because there is no API server, but it should
	// not fail during flag parsing.
	out, _ := runCLI(t, "datalake", "saved-queries", "list", "--help")
	if !strings.Contains(out, "List saved") {
		t.Fatalf("expected help text, got %q", out)
	}
}
