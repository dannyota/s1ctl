package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

func (s *Server) execCommand(parts []string) (string, error) {
	cliArgs := make([]string, 0, len(parts)+2)
	cliArgs = append(cliArgs, parts...)
	cliArgs = append(cliArgs, "--json", "--no-progress")
	var extraEnv []string
	if s.readOnly {
		extraEnv = append(extraEnv, "S1_READONLY=1")
	}
	return execSubprocess(cliArgs, extraEnv)
}

const (
	maxOutputBytes    = 4 << 20 // 4 MiB
	spillMaxAge       = 24 * time.Hour
	subprocessTimeout = 5 * time.Minute
)

func execSubprocess(args, extraEnv []string) (string, error) {
	self, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("find executable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), subprocessTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, self, args...) //nolint:gosec // self is os.Executable, args from tool schema
	cmd.Env = append(os.Environ(), extraEnv...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	execErr := cmd.Run()

	if execErr != nil {
		cause := execErr.Error()
		if ctx.Err() == context.DeadlineExceeded {
			cause = fmt.Sprintf("command timed out after %s", subprocessTimeout)
		}
		parts := []string{cause}
		if errOut := strings.TrimSpace(stderr.String()); errOut != "" {
			parts = append(parts, errOut)
		}
		if out := strings.TrimSpace(stdout.String()); out != "" {
			parts = append(parts, out)
		}
		return "", errors.New(capOutput(strings.Join(parts, "\n")))
	}

	out := stdout.Bytes()
	stderrTrimmed := strings.TrimSpace(stderr.String())

	if len(out) > maxOutputBytes {
		return spillOutputWithWarnings(out, stderrTrimmed)
	}
	if len(bytes.TrimSpace(out)) == 0 {
		// A successful command may write only diagnostics (e.g. sync
		// warnings) to stderr; surface them as a JSON envelope so the
		// always-JSON MCP channel stays parseable.
		return stderrFallback(stderr.String()), nil
	}
	if stderrTrimmed != "" {
		return mergeWarnings(string(out), stderrTrimmed), nil
	}
	return string(out), nil
}

// mergeWarnings wraps stdout and stderr into a JSON envelope when both are
// non-empty on a successful subprocess exit.
func mergeWarnings(stdout, stderr string) string {
	b, _ := json.Marshal(struct {
		Output   string `json:"output"`
		Warnings string `json:"warnings"`
	}{
		Output:   stdout,
		Warnings: stderr,
	})
	return string(b)
}

// stderrFallback returns a JSON envelope when stderr is non-empty, or an
// empty string when it is blank. This keeps the MCP channel always-JSON.
func stderrFallback(stderr string) string {
	trimmed := strings.TrimSpace(stderr)
	if trimmed == "" {
		return ""
	}
	b, _ := json.Marshal(struct {
		Output   string `json:"output"`
		Warnings string `json:"warnings"`
	}{
		Output:   "",
		Warnings: trimmed,
	})
	return string(b)
}

// capOutput bounds s to maxOutputBytes so oversized failure output cannot
// flood the client, cutting on a rune boundary.
func capOutput(s string) string {
	if len(s) <= maxOutputBytes {
		return s
	}
	cut := maxOutputBytes
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return fmt.Sprintf("%s\n[output truncated: %d of %d bytes shown]", s[:cut], cut, len(s))
}

func spillDir() string {
	return filepath.Join(os.TempDir(), "s1ctl-mcp")
}

const previewBytes = 2048

// runePrefix returns the first n bytes of s, cut on a rune boundary.
func runePrefix(s []byte, n int) string {
	if len(s) <= n {
		return string(s)
	}
	cut := n
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return string(s[:cut])
}

// spillOutputWithWarnings delegates to spillOutput and attaches stderr
// warnings when non-empty.
func spillOutputWithWarnings(out []byte, stderrTrimmed string) (string, error) {
	result, err := spillOutput(out)
	if err != nil || stderrTrimmed == "" {
		return result, err
	}
	// Inject warnings into the spill envelope.
	var env map[string]any
	if jsonErr := json.Unmarshal([]byte(result), &env); jsonErr == nil {
		env["warnings"] = stderrTrimmed
		b, _ := json.Marshal(env)
		return string(b), nil
	}
	return result, nil
}

// spillOutput writes oversized output to a spill file and returns a JSON
// pointer to it instead of the raw bytes.
func spillOutput(out []byte) (string, error) {
	dir := spillDir()
	sweepSpillFiles(dir, spillMaxAge)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", spillError(len(out), err)
	}
	f, err := os.CreateTemp(dir, "s1ctl-mcp-*.json")
	if err != nil {
		return "", spillError(len(out), err)
	}
	_, err = f.Write(out)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(f.Name())
		return "", spillError(len(out), err)
	}
	result, _ := json.Marshal(map[string]any{
		"file":    f.Name(),
		"bytes":   len(out),
		"preview": runePrefix(out, previewBytes),
		"message": fmt.Sprintf(
			"Output exceeded %d MiB limit. Results saved to a temporary file (removed after %dh). Read the file to analyze, or use --max-results or narrower filters to reduce output.",
			maxOutputBytes>>20, int(spillMaxAge.Hours())),
	})
	return string(result), nil
}

func spillError(size int, err error) error {
	return fmt.Errorf("output too large (%d bytes, limit %d) and failed to write temp file: %w", size, maxOutputBytes, err)
}

// sweepSpillFiles removes spill files older than maxAge so query output does
// not accumulate in the temp dir across sessions.
func sweepSpillFiles(dir string, maxAge time.Duration) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	cutoff := time.Now().Add(-maxAge)
	for _, e := range entries {
		info, err := e.Info()
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(dir, e.Name()))
		}
	}
}
