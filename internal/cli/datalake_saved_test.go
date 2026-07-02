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
