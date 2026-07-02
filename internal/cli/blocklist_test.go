package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestBlocklistCreateDryRun(t *testing.T) {
	out, err := runCLI(t, "blocklist", "create", "--value", "ffffffffffffffffffffffffffffffffffffffff", "--os-type", "linux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestBlocklistCreateRequiresValue(t *testing.T) {
	_, err := runCLI(t, "blocklist", "create", "--os-type", "linux")
	if err == nil {
		t.Fatal("expected validation error without --value")
	}
	if !strings.Contains(err.Error(), "--value is required") {
		t.Fatalf("expected %q, got %q", "--value is required", err.Error())
	}
}

func TestBlocklistCreateRequiresOSType(t *testing.T) {
	_, err := runCLI(t, "blocklist", "create", "--value", "ff")
	if err == nil {
		t.Fatal("expected validation error without --os-type")
	}
	if !strings.Contains(err.Error(), "--os-type is required") {
		t.Fatalf("expected %q, got %q", "--os-type is required", err.Error())
	}
}

func TestBlocklistUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "blocklist", "update", "B1", "--value", "ff", "--os-type", "linux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestBlocklistUpdateValidation(t *testing.T) {
	if _, err := runCLI(t, "blocklist", "update", "B1"); err == nil {
		t.Fatal("expected validation error without --value/--os-type")
	}
}

func TestBlocklistDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "blocklist", "delete", "B1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

// TestBlocklistPushMissingDir asserts a missing directory is a hard error naming
// it, before any API client is constructed (fully offline).
func TestBlocklistPushMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := runCLI(t, "blocklist", "push", "--dir", missing)
	if err == nil {
		t.Fatal("blocklist push: expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("blocklist push: error %q does not contain %q", err, "read "+missing)
	}
}

// TestBlocklistPushEmptyDir asserts an empty (but present) directory returns
// cleanly with the "No blocklist files found." message before any API call.
func TestBlocklistPushEmptyDir(t *testing.T) {
	out, err := runCLI(t, "blocklist", "push", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("blocklist push: unexpected error: %v", err)
	}
	if want := "No blocklist files found."; !strings.Contains(out, want) {
		t.Fatalf("blocklist push: output %q does not contain %q", out, want)
	}
}
