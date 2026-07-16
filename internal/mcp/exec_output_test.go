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
		Preview string `json:"preview"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(res), &meta); err != nil {
		t.Fatalf("result is not JSON: %v\n%s", err, res)
	}
	if meta.Bytes != len(payload) {
		t.Errorf("bytes = %d, want %d", meta.Bytes, len(payload))
	}
	if len(meta.Preview) != previewBytes {
		t.Errorf("preview len = %d, want %d", len(meta.Preview), previewBytes)
	}
	if meta.Preview != strings.Repeat("x", previewBytes) {
		t.Error("preview content mismatch")
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

func TestSpillPreviewRuneSafe(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())

	// Build payload where byte 2048 falls in the middle of a multi-byte rune.
	// U+00E9 (e-acute) is 2 bytes in UTF-8.
	payload := make([]byte, 0, maxOutputBytes+1)
	chunk := []byte("e\xc3\xa9") // "e" + "e-acute" = 3 bytes
	for len(payload) < maxOutputBytes+1 {
		payload = append(payload, chunk...)
	}

	res, err := spillOutput(payload)
	if err != nil {
		t.Fatalf("spillOutput: %v", err)
	}

	var meta struct {
		Preview string `json:"preview"`
	}
	if err := json.Unmarshal([]byte(res), &meta); err != nil {
		t.Fatalf("result is not JSON: %v", err)
	}
	if !utf8.ValidString(meta.Preview) {
		t.Error("preview is not valid UTF-8")
	}
	if len(meta.Preview) > previewBytes {
		t.Errorf("preview len = %d, exceeds %d", len(meta.Preview), previewBytes)
	}
}

func TestMergeWarnings(t *testing.T) {
	got := mergeWarnings(`{"data":[]}`, "warning: duplicate stem")
	var env struct {
		Output   string `json:"output"`
		Warnings string `json:"warnings"`
	}
	if err := json.Unmarshal([]byte(got), &env); err != nil {
		t.Fatalf("not JSON: %v\n%s", err, got)
	}
	if env.Output != `{"data":[]}` {
		t.Errorf("output = %q, want original stdout", env.Output)
	}
	if env.Warnings != "warning: duplicate stem" {
		t.Errorf("warnings = %q, want stderr", env.Warnings)
	}
}

func TestMergeWarningsStdoutOnlyUnchanged(t *testing.T) {
	// When mergeWarnings is not called (no stderr), stdout is returned raw.
	// This test verifies the helper itself always wraps both.
	got := mergeWarnings(`plain output`, "warn")
	if !strings.Contains(got, `"output"`) {
		t.Errorf("expected JSON envelope, got %s", got)
	}
}

func TestRunePrefixShortInput(t *testing.T) {
	got := runePrefix([]byte("short"), 2048)
	if got != "short" {
		t.Errorf("runePrefix = %q, want %q", got, "short")
	}
}
