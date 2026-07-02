package cli

import (
	"strings"
	"testing"
)

func TestDatalakeFacetValidation(t *testing.T) {
	if _, err := runCLI(t, "datalake", "facet", "--start", "1h"); err == nil || !strings.Contains(err.Error(), "--field is required") {
		t.Fatalf("expected --field validation error, got %v", err)
	}
	if _, err := runCLI(t, "datalake", "facet", "--field", "serverHost"); err == nil || !strings.Contains(err.Error(), "--start is required") {
		t.Fatalf("expected --start validation error, got %v", err)
	}
}

func TestDatalakeTimeseriesValidation(t *testing.T) {
	if _, err := runCLI(t, "datalake", "timeseries", "--start", "1h"); err == nil || !strings.Contains(err.Error(), "--filter is required") {
		t.Fatalf("expected --filter validation error, got %v", err)
	}
	if _, err := runCLI(t, "datalake", "timeseries", "--filter", "x"); err == nil || !strings.Contains(err.Error(), "--start is required") {
		t.Fatalf("expected --start validation error, got %v", err)
	}
}
