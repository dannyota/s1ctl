package cli

import (
	"strings"
	"testing"
)

func TestDatalakeDashboardsGetArgValidation(t *testing.T) {
	_, err := runCLI(t, "datalake", "dashboards", "get")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected arg validation error, got %v", err)
	}
}
