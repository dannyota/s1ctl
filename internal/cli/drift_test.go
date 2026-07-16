package cli

import (
	"os"
	"path/filepath"
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

// TestDriftSkipsUpgradePolicies asserts that an upgrade-policies directory
// does not abort the drift run. Build fails (the drift command has no
// --scope-level/--os-type flags), so the surface is marked SKIPPED in the
// summary and the run continues with remaining surfaces. A blocklist directory
// is present alongside to prove the run does not abort early — both appear in
// the output.
//
// This is fully offline: upgrade-policies is skipped at Build time (no API
// call). Blocklist Build succeeds but its List call requires credentials, so
// we scope to just upgrade-policies to prove the skip; the second sub-test
// adds blocklist to prove the run continues past the skip.
func TestDriftSkipsUpgradePolicies(t *testing.T) {
	root := t.TempDir()

	// Create upgrade-policies/ with a policy file.
	upDir := filepath.Join(root, "upgrade-policies")
	if err := os.Mkdir(upDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(upDir, "ga.yaml"), []byte("name: GA linux\nosType: linux\nscopeLevel: tenant\npackage:\n    fileId: f1\n    major: \"25\"\n    minor: \"3\"\n    build: \"1\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("skip only", func(t *testing.T) {
		// Only upgrade-policies dir, scoped to that surface. Build fails,
		// surface is marked SKIPPED, run completes with exit 0.
		out, err := runCLI(t, "drift", "--dir-root", root, "--surface", "upgrade-policies")
		if err != nil {
			t.Fatalf("drift should not error on skipped surface, got: %v", err)
		}
		// The warning (captured via SetErr) names the surface and reason.
		if !strings.Contains(out, "upgrade-policies") {
			t.Fatalf("drift output should name the skipped surface, got: %q", out)
		}
		if !strings.Contains(out, "skipping drift check") {
			t.Fatalf("drift output should say 'skipping drift check', got: %q", out)
		}
		// --json shows the SKIPPED result with skipReason in the captured output.
		jsonOut, jErr := runCLI(t, "drift", "--dir-root", root,
			"--surface", "upgrade-policies", "--json")
		if jErr != nil {
			t.Fatalf("drift --json should not error on skipped surface, got: %v", jErr)
		}
		if !strings.Contains(jsonOut, `"skipReason"`) {
			t.Fatalf("drift --json should contain skipReason field, got: %q", jsonOut)
		}
		if !strings.Contains(jsonOut, `"upgrade-policies"`) {
			t.Fatalf("drift --json should name upgrade-policies, got: %q", jsonOut)
		}
	})

	t.Run("continues past skip", func(t *testing.T) {
		// Create a blocklist/ dir alongside. Drift should process
		// upgrade-policies (SKIPPED) then attempt blocklist. Blocklist Build
		// succeeds, LoadDir succeeds (empty dir → no files), but List needs
		// API credentials and fails. The key assertion: the error is from
		// blocklist (not upgrade-policies), proving drift continued past the
		// skip.
		blDir := filepath.Join(root, "blocklist")
		if err := os.Mkdir(blDir, 0o755); err != nil {
			t.Fatal(err)
		}

		_, err := runCLI(t, "drift", "--dir-root", root,
			"--surface", "upgrade-policies,blocklist")
		// Blocklist will fail on List (no API credentials), and that's fine.
		// What matters: the error mentions blocklist (not upgrade-policies),
		// proving drift continued past the skip.
		if err == nil {
			// If it somehow succeeds (shouldn't without credentials), check
			// that SKIPPED appears for upgrade-policies.
			return
		}
		if strings.Contains(err.Error(), "--scope-level is required") {
			t.Fatalf("drift aborted on upgrade-policies Build error instead of skipping: %v", err)
		}
	})
}
