package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The network group mutations are dry-run by default and never construct an API
// client without --yes, so these tests exercise registration, flag parsing, and
// guard wording fully offline.

func TestNetworkMutationsDryRun(t *testing.T) {
	cases := []struct {
		name string
		args []string
		want string
	}{
		{"delete", []string{"network", "delete", "1", "2"}, "Would delete 2 network quarantine rules"},
		{"enable", []string{"network", "enable", "1"}, "Would enable 1 network quarantine rule"},
		{"disable", []string{"network", "disable", "1", "2"}, "Would disable 2 network quarantine rules"},
		{"reorder", []string{"network", "reorder", "1:1", "2:2"}, "Would reorder 2 network quarantine rules"},
		{"copy", []string{"network", "copy", "--target-site-id", "225494730938493804"}, "Would copy network quarantine rules to site 225494730938493804"},
		{"move", []string{"network", "move", "1", "--target-group-id", "225494730938493904"}, "Would move 1 network quarantine rule to group 225494730938493904"},
		{"set-location", []string{"network", "set-location", "1", "--type", "all"}, "Would set location all on 1 network quarantine rule"},
		{"tags-add", []string{"network", "tags", "add", "1", "--tag-id", "t1"}, "Would add 1 tag to 1 network quarantine rule"},
		{"tags-remove", []string{"network", "tags", "remove", "1", "2", "--tag-id", "t1"}, "Would remove 1 tag from 2 network quarantine rules"},
		{"configuration-set", []string{"network", "configuration", "set", "--enabled"}, "Would update network quarantine configuration"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := runCLI(t, tc.args...)
			if err != nil {
				t.Fatalf("network %s: unexpected error: %v", tc.name, err)
			}
			if !strings.Contains(out, tc.want) {
				t.Fatalf("network %s: output %q does not contain %q", tc.name, out, tc.want)
			}
			if !strings.Contains(out, "Pass --yes to apply.") {
				t.Fatalf("network %s: output %q missing dry-run notice", tc.name, out)
			}
		})
	}
}

func TestNetworkImportDryRun(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.json")
	if err := os.WriteFile(path, []byte(`[{"name":"NQ Rule"}]`), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	out, err := runCLI(t, "network", "import", path)
	if err != nil {
		t.Fatalf("network import: unexpected error: %v", err)
	}
	if want := "Would import network quarantine rules from " + path; !strings.Contains(out, want) {
		t.Fatalf("network import: output %q does not contain %q", out, want)
	}
}

func TestNetworkCopyRequiresTarget(t *testing.T) {
	_, err := runCLI(t, "network", "copy")
	if err == nil {
		t.Fatal("network copy: expected error when no target scope is given")
	}
	if !strings.Contains(err.Error(), "--target-site-id") {
		t.Fatalf("network copy: error %q does not mention required target flags", err)
	}
}

func TestNetworkSetLocationSpecificRequiresID(t *testing.T) {
	_, err := runCLI(t, "network", "set-location", "1", "--type", "specific")
	if err == nil {
		t.Fatal("network set-location: expected error for specific type without --location-id")
	}
	if !strings.Contains(err.Error(), "--location-id is required") {
		t.Fatalf("network set-location: error %q does not name the missing flag", err)
	}
}

func TestNetworkSetLocationInvalidType(t *testing.T) {
	_, err := runCLI(t, "network", "set-location", "1", "--type", "bogus")
	if err == nil {
		t.Fatal("network set-location: expected error for invalid --type")
	}
	if !strings.Contains(err.Error(), "invalid --type") {
		t.Fatalf("network set-location: error %q does not describe the invalid type", err)
	}
}

func TestNetworkTagsRequiresTagID(t *testing.T) {
	_, err := runCLI(t, "network", "tags", "add", "1")
	if err == nil {
		t.Fatal("network tags add: expected error when no --tag-id is given")
	}
	if !strings.Contains(err.Error(), "--tag-id is required") {
		t.Fatalf("network tags add: error %q does not name the missing flag", err)
	}
}

func TestNetworkConfigurationSetRequiresChange(t *testing.T) {
	_, err := runCLI(t, "network", "configuration", "set")
	if err == nil {
		t.Fatal("network configuration set: expected error when nothing is changed")
	}
	if !strings.Contains(err.Error(), "nothing to update") {
		t.Fatalf("network configuration set: error %q does not describe the no-op", err)
	}
}

func TestNetworkPushMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := runCLI(t, "network", "push", "--dir", missing)
	if err == nil {
		t.Fatal("network push: expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("network push: error %q does not contain %q", err, "read "+missing)
	}
}

func TestNetworkPushEmptyDir(t *testing.T) {
	out, err := runCLI(t, "network", "push", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("network push: unexpected error: %v", err)
	}
	if want := "No network quarantine rule files found."; !strings.Contains(out, want) {
		t.Fatalf("network push: output %q does not contain %q", out, want)
	}
}
