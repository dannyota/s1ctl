package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// platformSurfaces are the engine-backed sites/groups/tags surfaces, keyed by
// CLI group with the noun their empty-directory message uses.
var platformSurfaces = map[string]string{
	"sites":  "No site files found.",
	"groups": "No group files found.",
	"tags":   "No tag files found.",
}

// TestPlatformPushMissingDir asserts a missing directory is a hard error naming
// it. The stat check runs before any API client is constructed, so this stays
// fully offline.
func TestPlatformPushMissingDir(t *testing.T) {
	for group := range platformSurfaces {
		missing := filepath.Join(t.TempDir(), "nope")
		_, err := runCLI(t, group, "push", "--dir", missing)
		if err == nil {
			t.Fatalf("%s push: expected error for missing dir", group)
		}
		if !strings.Contains(err.Error(), "read "+missing) {
			t.Fatalf("%s push: error %q does not contain %q", group, err, "read "+missing)
		}
	}
}

// TestPlatformPushEmptyDir asserts an empty (but present) directory returns
// cleanly with the "No <noun> files found." message before any API call.
func TestPlatformPushEmptyDir(t *testing.T) {
	for group, want := range platformSurfaces {
		out, err := runCLI(t, group, "push", "--dir", t.TempDir())
		if err != nil {
			t.Fatalf("%s push: unexpected error: %v", group, err)
		}
		if !strings.Contains(out, want) {
			t.Fatalf("%s push: output %q does not contain %q", group, out, want)
		}
	}
}
