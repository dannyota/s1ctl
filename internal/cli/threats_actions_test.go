package cli

import (
	"strings"
	"testing"
)

func TestThreatPlainActionsDryRun(t *testing.T) {
	for _, verb := range []string{"blacklist", "fetch-file"} {
		out, err := runCLI(t, "threats", verb, "T1")
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", verb, err)
		}
		if !strings.Contains(out, "Would") || !strings.Contains(out, "T1") {
			t.Fatalf("%s: expected dry-run message, got %q", verb, out)
		}
	}
}
