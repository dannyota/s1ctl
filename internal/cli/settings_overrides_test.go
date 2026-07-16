package cli

import (
	"strings"
	"testing"
)

func TestSettingsOverridesCreateDryRun(t *testing.T) {
	out, err := runCLI(t, "settings", "overrides", "create",
		"--name", "test-override",
		"--os-type", "linux",
		"--config", `{"key":"value"}`,
		"--scope", "site",
		"--scope-id", "225494730938493804",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
	if !strings.Contains(out, `{"key":"value"}`) {
		t.Fatalf("expected config in dry-run output, got %q", out)
	}
}

func TestSettingsOverridesCreateRequiresName(t *testing.T) {
	_, err := runCLI(t, "settings", "overrides", "create",
		"--os-type", "linux",
		"--config", `{}`,
		"--scope", "site",
		"--scope-id", "123",
	)
	if err == nil {
		t.Fatal("expected validation error without --name")
	}
	if !strings.Contains(err.Error(), "--name is required") {
		t.Fatalf("expected %q, got %q", "--name is required", err.Error())
	}
}

func TestSettingsOverridesCreateRequiresOSType(t *testing.T) {
	_, err := runCLI(t, "settings", "overrides", "create",
		"--name", "x",
		"--config", `{}`,
		"--scope", "site",
		"--scope-id", "123",
	)
	if err == nil {
		t.Fatal("expected validation error without --os-type")
	}
	if !strings.Contains(err.Error(), "--os-type is required") {
		t.Fatalf("expected %q, got %q", "--os-type is required", err.Error())
	}
}

func TestSettingsOverridesCreateRequiresConfig(t *testing.T) {
	_, err := runCLI(t, "settings", "overrides", "create",
		"--name", "x",
		"--os-type", "linux",
		"--scope", "site",
		"--scope-id", "123",
	)
	if err == nil {
		t.Fatal("expected validation error without --config")
	}
	if !strings.Contains(err.Error(), "--config is required") {
		t.Fatalf("expected %q, got %q", "--config is required", err.Error())
	}
}

func TestSettingsOverridesCreateRequiresScope(t *testing.T) {
	_, err := runCLI(t, "settings", "overrides", "create",
		"--name", "x",
		"--os-type", "linux",
		"--config", `{}`,
	)
	if err == nil {
		t.Fatal("expected validation error without --scope")
	}
	if !strings.Contains(err.Error(), "--scope is required") {
		t.Fatalf("expected %q, got %q", "--scope is required", err.Error())
	}
}

func TestSettingsOverridesCreateRequiresScopeID(t *testing.T) {
	_, err := runCLI(t, "settings", "overrides", "create",
		"--name", "x",
		"--os-type", "linux",
		"--config", `{}`,
		"--scope", "site",
	)
	if err == nil {
		t.Fatal("expected validation error without --scope-id for site scope")
	}
	if !strings.Contains(err.Error(), "--scope-id is required") {
		t.Fatalf("expected %q, got %q", "--scope-id is required", err.Error())
	}
}

func TestSettingsOverridesCreateTenantNoScopeID(t *testing.T) {
	out, err := runCLI(t, "settings", "overrides", "create",
		"--name", "x",
		"--os-type", "linux",
		"--config", `{}`,
		"--scope", "tenant",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestSettingsOverridesCreateInvalidJSON(t *testing.T) {
	_, err := runCLI(t, "settings", "overrides", "create",
		"--name", "x",
		"--os-type", "linux",
		"--config", "not-json",
		"--scope", "tenant",
	)
	if err == nil {
		t.Fatal("expected validation error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "--config must be valid JSON") {
		t.Fatalf("expected JSON error, got %q", err.Error())
	}
}

func TestSettingsOverridesUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "settings", "overrides", "update", "O1",
		"--name", "updated",
		"--config", `{"new":"val"}`,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
	if !strings.Contains(out, `{"new":"val"}`) {
		t.Fatalf("expected config in dry-run output, got %q", out)
	}
}

func TestSettingsOverridesDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "settings", "overrides", "delete", "O1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}
