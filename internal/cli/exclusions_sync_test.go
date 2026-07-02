package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestExclusionsPushMissingDir asserts a missing directory is a hard error
// naming it. The stat check runs before any API client is constructed, so this
// stays fully offline.
func TestExclusionsPushMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := runCLI(t, "exclusions", "push", "--dir", missing)
	if err == nil {
		t.Fatal("exclusions push: expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("exclusions push: error %q does not contain %q", err, "read "+missing)
	}
}

// TestExclusionsPushEmptyDir asserts an empty (but present) directory returns
// cleanly with the "No exclusion files found." message before any API call.
func TestExclusionsPushEmptyDir(t *testing.T) {
	out, err := runCLI(t, "exclusions", "push", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("exclusions push: unexpected error: %v", err)
	}
	if want := "No exclusion files found."; !strings.Contains(out, want) {
		t.Fatalf("exclusions push: output %q does not contain %q", out, want)
	}
}
