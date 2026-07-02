package cli

import (
	"bytes"
	"testing"
)

// runCLI executes the full command tree offline and returns combined output.
// Dry-run (no --yes) mutations never construct an API client, so these
// tests exercise registration, flag parsing, and guard wording only.
func runCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()
	t.Setenv("S1_READONLY", "")
	root := newRootCmd()
	registerCommands(root)
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}
