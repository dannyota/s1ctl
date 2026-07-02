package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeMaintenanceDataFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "data.json")
	body := `{"maxConcurrent":5,"timezoneGmt":"GMT+00:00","maintenanceWindowsByDay":{"monday":[]}}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write maintenance data file: %v", err)
	}
	return path
}

func TestMaintenanceSetDryRun(t *testing.T) {
	out, err := runCLI(t, "maintenance", "set", "--task-type", "agents_upgrade", "--from-file", writeMaintenanceDataFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMaintenanceSetRequiresTaskType(t *testing.T) {
	_, err := runCLI(t, "maintenance", "set", "--from-file", writeMaintenanceDataFile(t))
	if err == nil {
		t.Fatal("expected validation error without --task-type")
	}
	if !strings.Contains(err.Error(), "--task-type is required") {
		t.Fatalf("expected %q, got %q", "--task-type is required", err.Error())
	}
}

func TestMaintenanceSetRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "maintenance", "set", "--task-type", "agents_upgrade")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}

func TestMaintenanceGetRequiresTaskType(t *testing.T) {
	_, err := runCLI(t, "maintenance", "get")
	if err == nil {
		t.Fatal("expected validation error without --task-type")
	}
	if !strings.Contains(err.Error(), "--task-type is required") {
		t.Fatalf("expected %q, got %q", "--task-type is required", err.Error())
	}
}

func TestMaintenanceSetFlexibleDryRun(t *testing.T) {
	path := filepath.Join(t.TempDir(), "body.json")
	if err := os.WriteFile(path, []byte(`{"data":{},"filter":{"taskType":"agents_upgrade","tenant":true}}`), 0o600); err != nil {
		t.Fatalf("write flexible body: %v", err)
	}
	out, err := runCLI(t, "maintenance", "set-flexible", "--from-file", path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMaintenanceSetFlexibleRequiresFromFile(t *testing.T) {
	_, err := runCLI(t, "maintenance", "set-flexible")
	if err == nil {
		t.Fatal("expected validation error without --from-file")
	}
	if !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected %q, got %q", "--from-file is required", err.Error())
	}
}
