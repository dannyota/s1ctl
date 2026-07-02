package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestCloudPoliciesPushMissingDir asserts a missing directory is a hard error
// naming it. The stat check runs before any GraphQL client is constructed, so
// this stays fully offline.
func TestCloudPoliciesPushMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := runCLI(t, "cloud-policies", "push", "--dir", missing)
	if err == nil {
		t.Fatal("cloud-policies push: expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("cloud-policies push: error %q does not contain %q", err, "read "+missing)
	}
}

// TestCloudPoliciesPushEmptyDir asserts an empty (but present) directory returns
// cleanly with the "No cloud policy files found." message before any API call.
func TestCloudPoliciesPushEmptyDir(t *testing.T) {
	out, err := runCLI(t, "cloud-policies", "push", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("cloud-policies push: unexpected error: %v", err)
	}
	if want := "No cloud policy files found."; !strings.Contains(out, want) {
		t.Fatalf("cloud-policies push: output %q does not contain %q", out, want)
	}
}
