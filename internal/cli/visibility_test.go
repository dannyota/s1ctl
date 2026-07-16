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

// dvEventsEnvelope is the expected JSON shape from --json output.
type dvEventsEnvelope struct {
	Data     []map[string]any `json:"data"`
	Returned int              `json:"returned"`
	Total    int              `json:"total"`
}

func TestPrintDVEventsJSONEmptyIsEnvelope(t *testing.T) {
	out, errOut := printDVEventsJSON(t, nil, 0)
	var env dvEventsEnvelope
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("stdout is not a JSON envelope: %v\n%s", err, out)
	}
	if env.Data == nil {
		t.Error("data must be [] not null")
	}
	if len(env.Data) != 0 {
		t.Errorf("data len = %d, want 0", len(env.Data))
	}
	if env.Returned != 0 {
		t.Errorf("returned = %d, want 0", env.Returned)
	}
	if env.Total != 0 {
		t.Errorf("total = %d, want 0", env.Total)
	}
	if errOut != "" {
		t.Errorf("stderr = %q, want empty", errOut)
	}
}

func TestPrintDVEventsJSONTruncatedEnvelope(t *testing.T) {
	out, errOut := printDVEventsJSON(t, []mgmt.DVEvent{{EventType: "process"}}, 10)
	var env dvEventsEnvelope
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("stdout is not a JSON envelope: %v\n%s", err, out)
	}
	if len(env.Data) != 1 {
		t.Fatalf("data len = %d, want 1", len(env.Data))
	}
	if env.Returned != 1 {
		t.Errorf("returned = %d, want 1", env.Returned)
	}
	if env.Total != 10 {
		t.Errorf("total = %d, want 10", env.Total)
	}
	if !strings.Contains(errOut, "Showing 1 of 10 events") {
		t.Errorf("stderr = %q, want truncation notice", errOut)
	}
}

func TestPrintDVEventsJSONCompleteNoNotice(t *testing.T) {
	out, errOut := printDVEventsJSON(t, []mgmt.DVEvent{{EventType: "process"}}, 1)
	var env dvEventsEnvelope
	if err := json.Unmarshal([]byte(out), &env); err != nil {
		t.Fatalf("stdout is not a JSON envelope: %v\n%s", err, out)
	}
	if env.Returned != 1 {
		t.Errorf("returned = %d, want 1", env.Returned)
	}
	if env.Total != 1 {
		t.Errorf("total = %d, want 1", env.Total)
	}
	if errOut != "" {
		t.Errorf("stderr = %q, want empty for complete results", errOut)
	}
}
