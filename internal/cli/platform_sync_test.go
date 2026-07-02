package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlatformPushDryRun(t *testing.T) {
	dir := t.TempDir()
	cases := []struct {
		resource, payload, count string
	}{
		{"sites", `[{"name":"s1","accountId":"AC1"}]`, "1 site"},
		{"groups", `[{"name":"g1","siteId":"S1"}]`, "1 group"},
		{"tags", `[{"key":"env","value":"prod"}]`, "1 tag"},
	}
	for _, tc := range cases {
		f := filepath.Join(dir, tc.resource+".json")
		if err := os.WriteFile(f, []byte(tc.payload), 0o644); err != nil {
			t.Fatal(err)
		}
		out, err := runCLI(t, tc.resource, "push", "--file", f)
		if err != nil {
			t.Fatalf("%s push: unexpected error: %v", tc.resource, err)
		}
		if !strings.Contains(out, "Would") {
			t.Fatalf("%s push: expected dry-run message, got %q", tc.resource, out)
		}
		if !strings.Contains(out, tc.count) {
			t.Fatalf("%s push: expected count %q in output, got %q", tc.resource, tc.count, out)
		}
	}
}
