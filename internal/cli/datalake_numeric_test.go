package cli

import (
	"strings"
	"testing"
)

func TestDatalakeNumericValidation(t *testing.T) {
	if _, err := runCLI(t, "datalake", "numeric"); err == nil || !strings.Contains(err.Error(), "--start is required") {
		t.Fatalf("expected --start validation error, got %v", err)
	}
}
