package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func printDVEventsJSON(t *testing.T, events []mgmt.DVEvent, total int) (stdout, stderr string) {
	t.Helper()
	prev := outputFormat
	outputFormat = "json"
	t.Cleanup(func() { outputFormat = prev })

	cmd := &cobra.Command{}
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	if err := printDVEvents(cmd, events, total); err != nil {
		t.Fatalf("printDVEvents: %v", err)
	}
	return out.String(), errOut.String()
}

func TestPrintDVEventsJSONEmptyIsArray(t *testing.T) {
	out, errOut := printDVEventsJSON(t, nil, 0)
	if strings.TrimSpace(out) != "[]" {
		t.Errorf("stdout = %q, want []", out)
	}
	if errOut != "" {
		t.Errorf("stderr = %q, want empty", errOut)
	}
}

func TestPrintDVEventsJSONTruncatedStaysArray(t *testing.T) {
	out, errOut := printDVEventsJSON(t, []mgmt.DVEvent{{EventType: "process"}}, 10)
	var arr []map[string]any
	if err := json.Unmarshal([]byte(out), &arr); err != nil {
		t.Fatalf("stdout is not a JSON array: %v\n%s", err, out)
	}
	if len(arr) != 1 {
		t.Fatalf("len = %d, want 1", len(arr))
	}
	if !strings.Contains(errOut, "Showing 1 of 10 events") {
		t.Errorf("stderr = %q, want truncation notice", errOut)
	}
}

func TestPrintDVEventsJSONCompleteNoNotice(t *testing.T) {
	_, errOut := printDVEventsJSON(t, []mgmt.DVEvent{{EventType: "process"}}, 1)
	if errOut != "" {
		t.Errorf("stderr = %q, want empty for complete results", errOut)
	}
}
