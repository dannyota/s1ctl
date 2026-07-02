package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLocationsCreateDryRun(t *testing.T) {
	out, err := runCLI(t, "locations", "create", "--name", "HQ Office")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestLocationsCreateRequiresName(t *testing.T) {
	_, err := runCLI(t, "locations", "create")
	if err == nil {
		t.Fatal("expected validation error without --name")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Fatalf("expected %q, got %q", "--name is required", err.Error())
	}
}

func TestLocationsCreateInvalidOperator(t *testing.T) {
	_, err := runCLI(t, "locations", "create", "--name", "HQ", "--operator", "bogus")
	if err == nil {
		t.Fatal("expected validation error for bad operator")
	}
	if !strings.Contains(err.Error(), "--operator must be one of") {
		t.Fatalf("expected operator validation error, got %q", err.Error())
	}
}

func TestLocationsUpdateRequiresName(t *testing.T) {
	_, err := runCLI(t, "locations", "update", "L1")
	if err == nil {
		t.Fatal("expected validation error without --name")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Fatalf("expected %q, got %q", "--name is required", err.Error())
	}
}

func TestLocationsDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "locations", "delete", "L1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestLocationsPushMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := runCLI(t, "locations", "push", "--dir", missing)
	if err == nil {
		t.Fatal("expected error for missing directory")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("expected %q, got %q", "read "+missing, err.Error())
	}
}

func TestLocationsPushEmptyDir(t *testing.T) {
	dir := t.TempDir()
	out, err := runCLI(t, "locations", "push", "--dir", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No location files found.") {
		t.Fatalf("expected empty-dir message, got %q", out)
	}
}

// TestLocationFileRoundTrip verifies a location's declarative body survives the
// file -> data -> file round trip, including a raw detection-parameter group.
func TestLocationFileRoundTrip(t *testing.T) {
	body := "name: HQ Office\noperator: any\nipAddresses:\n    enabled: true\n"
	obj, err := decodeLocation([]byte(body))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if obj.Name != "HQ Office" {
		t.Fatalf("unexpected identity: %s", obj.Name)
	}
	var f locationFile
	if err := yaml.Unmarshal(obj.Body, &f); err != nil {
		t.Fatalf("reparse: %v", err)
	}
	data, err := f.toData()
	if err != nil {
		t.Fatalf("toData: %v", err)
	}
	if data.Name != "HQ Office" || string(data.Operator) != "any" {
		t.Fatalf("unexpected data: %+v", data)
	}
	if len(data.IPAddresses) == 0 {
		t.Fatal("expected ipAddresses to survive round trip")
	}
}
