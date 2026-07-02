package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSettingsUpdateDryRun(t *testing.T) {
	f := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(f, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, sub := range []string{"notifications", "sso", "smtp", "syslog", "sms", "recipients", "ad", "ad-scope-mapping"} {
		out, err := runCLI(t, "settings", "update", sub, "--from-file", f, "--site-id", "S1")
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", sub, err)
		}
		if !strings.Contains(out, "Would") {
			t.Fatalf("%s: expected dry-run message, got %q", sub, out)
		}
	}
}

func TestSettingsUpdateValidation(t *testing.T) {
	if _, err := runCLI(t, "settings", "update", "sso"); err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file validation error, got %v", err)
	}
}
