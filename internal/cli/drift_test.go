package cli

import (
	"strings"
	"testing"
)

// TestDriftUnknownSurface asserts an unrecognized --surface value is a hard
// error that names the bad surface and lists the valid ones (offline: the
// registry check runs before any directory read or API client).
func TestDriftUnknownSurface(t *testing.T) {
	_, err := runCLI(t, "drift", "--surface", "bogus")
	if err == nil {
		t.Fatal("drift --surface bogus: expected error, got nil")
	}
	if !strings.Contains(err.Error(), `unknown surface "bogus"`) {
		t.Fatalf("drift: error %q does not contain %q", err, `unknown surface "bogus"`)
	}
	// The error is helpful: it enumerates valid surfaces (registry-driven).
	if !strings.Contains(err.Error(), "sites") {
		t.Fatalf("drift: error %q does not list valid surfaces (e.g. \"sites\")", err)
	}
}

// TestDriftNoLocalDirs asserts that when no surface directory exists under
// --dir-root, drift reports it and exits 0 without constructing any client
// (an empty t.TempDir() has no per-surface subdirectories, so every surface is
// skipped and List is never reached — no credentials required).
func TestDriftNoLocalDirs(t *testing.T) {
	out, err := runCLI(t, "drift", "--dir-root", t.TempDir())
	if err != nil {
		t.Fatalf("drift: unexpected error: %v", err)
	}
	if !strings.Contains(out, "no local surface directories found") {
		t.Fatalf("drift: output %q does not contain %q", out, "no local surface directories found")
	}
}
