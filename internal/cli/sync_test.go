package cli

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestSyncSurfaceSpecs asserts the registry the drift command will iterate.
func TestSyncSurfaceSpecs(t *testing.T) {
	want := map[string]string{
		"blocklist":      "blocklist",
		"cloud-policies": "cloud-policies",
		"devicecontrol":  "devicecontrol",
		"exclusions":     "exclusions",
		"firewall":       "firewall",
		"network":        "network-quarantine",
		"groups":         "groups",
		"rules":          "rules",
		"sites":          "sites",
		"tags":           "tags",
	}
	specs := syncSurfaceSpecs()
	if len(specs) != len(want) {
		t.Fatalf("got %d specs, want %d", len(specs), len(want))
	}
	for _, s := range specs {
		dir, ok := want[s.Command]
		if !ok {
			t.Errorf("unexpected spec command %q", s.Command)
			continue
		}
		if s.DefaultDir != dir {
			t.Errorf("%s: DefaultDir = %q, want %q", s.Command, s.DefaultDir, dir)
		}
		if s.Noun == "" {
			t.Errorf("%s: empty Noun", s.Command)
		}
		if s.Build == nil {
			t.Errorf("%s: nil Build", s.Command)
		}
	}
}

// TestSyncPushMissingDir asserts a missing directory is a hard error naming the
// directory (offline: the check runs before any API client is constructed).
func TestSyncPushMissingDir(t *testing.T) {
	for _, group := range []string{"devicecontrol", "firewall", "rules"} {
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

// TestSyncPushEmptyDir asserts an empty (but present) directory returns cleanly
// with the legacy "No <noun> files found." message (offline: no API call).
func TestSyncPushEmptyDir(t *testing.T) {
	cases := map[string]string{
		"devicecontrol": "No device rule files found.",
		"firewall":      "No firewall rule files found.",
		"rules":         "No rule files found.",
	}
	for group, want := range cases {
		out, err := runCLI(t, group, "push", "--dir", t.TempDir())
		if err != nil {
			t.Fatalf("%s push: unexpected error: %v", group, err)
		}
		if !strings.Contains(out, want) {
			t.Fatalf("%s push: output %q does not contain %q", group, out, want)
		}
	}
}
