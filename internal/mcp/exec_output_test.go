package mcp

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestCapOutputShortUnchanged(t *testing.T) {
	if got := capOutput("hello"); got != "hello" {
		t.Errorf("capOutput = %q, want unchanged", got)
	}
}

func TestCapOutputTruncatesOnRuneBoundary(t *testing.T) {
	big := strings.Repeat("é", maxOutputBytes/2+16)
	got := capOutput(big)
	if len(got) > maxOutputBytes+64 {
		t.Errorf("len = %d, want capped near %d", len(got), maxOutputBytes)
	}
	if !strings.Contains(got, "[output truncated:") {
		t.Error("missing truncation marker")
	}
	if !utf8.ValidString(got) {
		t.Error("truncated output is not valid UTF-8")
	}
}

// TestStderrFallbackIsJSON verifies that when a subprocess exits 0 with empty
// stdout and non-empty stderr, the result is a JSON envelope with output and
// warnings keys rather than raw prose.
func TestStderrFallbackIsJSON(t *testing.T) {
	got := stderrFallback("  warning: something happened  ")
	var env struct {
		Output   string `json:"output"`
		Warnings string `json:"warnings"`
	}
	if err := json.Unmarshal([]byte(got), &env); err != nil {
		t.Fatalf("stderrFallback result is not JSON: %v\n%s", err, got)
	}
	if env.Output != "" {
		t.Errorf("output = %q, want empty", env.Output)
	}
	if env.Warnings != "warning: something happened" {
		t.Errorf("warnings = %q, want trimmed stderr", env.Warnings)
	}
}

// TestStderrFallbackEmpty verifies that blank stderr returns empty string.
func TestStderrFallbackEmpty(t *testing.T) {
	got := stderrFallback("   ")
	if got != "" {
		t.Errorf("stderrFallback(%q) = %q, want empty", "   ", got)
	}
}

func TestSpillOutputRoundTripAndSweep(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())

	payload := bytes.Repeat([]byte("x"), maxOutputBytes+1)
	res, err := spillOutput(payload)
	if err != nil {
		t.Fatalf("spillOutput: %v", err)
	}

	var meta struct {
		File    string `json:"file"`
		Bytes   int    `json:"bytes"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(res), &meta); err != nil {
		t.Fatalf("result is not JSON: %v\n%s", err, res)
	}
	if meta.Bytes != len(payload) {
		t.Errorf("bytes = %d, want %d", meta.Bytes, len(payload))
	}
	data, err := os.ReadFile(meta.File)
	if err != nil {
		t.Fatalf("read spill file: %v", err)
	}
	if !bytes.Equal(data, payload) {
		t.Error("spill file content differs from output")
	}

	// A file older than the retention window is swept; a fresh one stays.
	old := time.Now().Add(-2 * spillMaxAge)
	if err := os.Chtimes(meta.File, old, old); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	fresh, err := spillOutput(payload)
	if err != nil {
		t.Fatalf("second spillOutput: %v", err)
	}
	if _, err := os.Stat(meta.File); !os.IsNotExist(err) {
		t.Error("expected aged spill file to be swept")
	}
	if err := json.Unmarshal([]byte(fresh), &meta); err != nil {
		t.Fatalf("second result is not JSON: %v", err)
	}
	if _, err := os.Stat(meta.File); err != nil {
		t.Errorf("fresh spill file should remain: %v", err)
	}
}
