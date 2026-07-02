package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeRoleFile writes a minimal valid role file and returns its path.
func writeRoleFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "role.yaml")
	if err := os.WriteFile(path, []byte("name: Custom IT\ndescription: IT operators\n"), 0o600); err != nil {
		t.Fatalf("write role file: %v", err)
	}
	return path
}

func TestRolesCreateDryRun(t *testing.T) {
	out, err := runCLI(t, "roles", "create", "--from-file", writeRoleFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestRolesCreateRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "roles", "create")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}

func TestRolesCreateMissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.yaml")
	_, err := runCLI(t, "roles", "create", "--from-file", missing)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("expected %q, got %q", "read "+missing, err.Error())
	}
}

func TestRolesUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "roles", "update", "R1", "--from-file", writeRoleFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestRolesUpdateRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "roles", "update", "R1")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}

func TestRolesDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "roles", "delete", "R1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}
