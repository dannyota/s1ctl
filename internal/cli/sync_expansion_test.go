package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
)

// --- Filters sync tests ---

func TestFiltersPushMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := runCLI(t, "filters", "push", "--dir", missing)
	if err == nil {
		t.Fatal("filters push: expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("filters push: error %q does not contain %q", err, "read "+missing)
	}
}

func TestFiltersPushEmptyDir(t *testing.T) {
	out, err := runCLI(t, "filters", "push", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("filters push: unexpected error: %v", err)
	}
	if want := "No filter files found."; !strings.Contains(out, want) {
		t.Fatalf("filters push: output %q does not contain %q", out, want)
	}
}

func TestDecodeFilterRoundTrip(t *testing.T) {
	input := `name: Infected
filterFields:
    infected: true
`
	obj, err := decodeFilter([]byte(input))
	if err != nil {
		t.Fatalf("decodeFilter: %v", err)
	}
	if obj.Name != "Infected" {
		t.Fatalf("Name = %q, want %q", obj.Name, "Infected")
	}
	// Re-decode must produce identical body.
	obj2, err := decodeFilter(obj.Body)
	if err != nil {
		t.Fatalf("re-decode: %v", err)
	}
	if string(obj.Body) != string(obj2.Body) {
		t.Fatalf("round-trip mismatch:\n%s\nvs\n%s", obj.Body, obj2.Body)
	}
}

func TestDecodeFilterMissingName(t *testing.T) {
	input := `filterFields:
    infected: true
`
	_, err := decodeFilter([]byte(input))
	if err == nil {
		t.Fatal("expected error for filter without name")
	}
	if !strings.Contains(err.Error(), "no name") {
		t.Fatalf("error %q does not mention missing name", err)
	}
}

// --- Tag rules sync tests ---

func TestTagRulesPushMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := runCLI(t, "tag-rules", "push", "--dir", missing)
	if err == nil {
		t.Fatal("tag-rules push: expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("tag-rules push: error %q does not contain %q", err, "read "+missing)
	}
}

func TestTagRulesPushEmptyDir(t *testing.T) {
	out, err := runCLI(t, "tag-rules", "push", "--dir", t.TempDir())
	if err != nil {
		t.Fatalf("tag-rules push: unexpected error: %v", err)
	}
	if want := "No tag rule files found."; !strings.Contains(out, want) {
		t.Fatalf("tag-rules push: output %q does not contain %q", out, want)
	}
}

func TestDecodeTagRuleRoundTrip(t *testing.T) {
	input := `name: Tag servers
description: Auto-tag production servers
status: enabled
conditions:
    op: and
    items:
        - field: hostname
          op: contains
          value: prod
tags:
    - key: env
      value: production
`
	obj, err := decodeTagRule([]byte(input))
	if err != nil {
		t.Fatalf("decodeTagRule: %v", err)
	}
	if obj.Name != "Tag servers" {
		t.Fatalf("Name = %q, want %q", obj.Name, "Tag servers")
	}
	obj2, err := decodeTagRule(obj.Body)
	if err != nil {
		t.Fatalf("re-decode: %v", err)
	}
	if string(obj.Body) != string(obj2.Body) {
		t.Fatalf("round-trip mismatch:\n%s\nvs\n%s", obj.Body, obj2.Body)
	}
}

func TestDecodeTagRuleMissingName(t *testing.T) {
	input := `status: enabled
conditions:
    op: and
`
	_, err := decodeTagRule([]byte(input))
	if err == nil {
		t.Fatal("expected error for tag rule without name")
	}
	if !strings.Contains(err.Error(), "no name") {
		t.Fatalf("error %q does not mention missing name", err)
	}
}

// --- Upgrade policies sync tests ---

func TestUpgradePoliciesPushRequiresScope(t *testing.T) {
	_, err := runCLI(t, "upgrade-policies", "push", "--dir", t.TempDir(), "--os-type", "linux")
	if err == nil {
		t.Fatal("expected error without --scope-level")
	}
	if !strings.Contains(err.Error(), "--scope-level is required") {
		t.Fatalf("error %q does not mention --scope-level", err)
	}
}

func TestUpgradePoliciesPushRequiresOSType(t *testing.T) {
	_, err := runCLI(t, "upgrade-policies", "push", "--dir", t.TempDir(), "--scope-level", "site")
	if err == nil {
		t.Fatal("expected error without --os-type")
	}
	if !strings.Contains(err.Error(), "--os-type is required") {
		t.Fatalf("error %q does not mention --os-type", err)
	}
}

func TestUpgradePoliciesPushMissingDir(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := runCLI(t, "upgrade-policies", "push", "--dir", missing,
		"--scope-level", "tenant", "--os-type", "linux")
	if err == nil {
		t.Fatal("upgrade-policies push: expected error for missing dir")
	}
	if !strings.Contains(err.Error(), "read "+missing) {
		t.Fatalf("upgrade-policies push: error %q does not contain %q", err, "read "+missing)
	}
}

func TestUpgradePoliciesPushEmptyDir(t *testing.T) {
	out, err := runCLI(t, "upgrade-policies", "push", "--dir", t.TempDir(),
		"--scope-level", "tenant", "--os-type", "linux")
	if err != nil {
		t.Fatalf("upgrade-policies push: unexpected error: %v", err)
	}
	if want := "No upgrade policy files found."; !strings.Contains(out, want) {
		t.Fatalf("upgrade-policies push: output %q does not contain %q", out, want)
	}
}

func TestDecodeUpgradePolicyRoundTrip(t *testing.T) {
	f := upgradePolicyFile{
		Name:         "Auto upgrade linux",
		OSType:       "linux",
		ScopeLevel:   "tenant",
		IsActive:     true,
		AllEndpoints: true,
		MaxRetries:   5,
		Package: upgradePolicyPkgFile{
			FileID: "pkg-123",
			Major:  "25",
			Minor:  "3",
			Build:  "1",
		},
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	obj, err := decodeUpgradePolicy(data)
	if err != nil {
		t.Fatalf("decodeUpgradePolicy: %v", err)
	}
	if obj.Name != "Auto upgrade linux" {
		t.Fatalf("Name = %q, want %q", obj.Name, "Auto upgrade linux")
	}
	obj2, err := decodeUpgradePolicy(obj.Body)
	if err != nil {
		t.Fatalf("re-decode: %v", err)
	}
	if string(obj.Body) != string(obj2.Body) {
		t.Fatalf("round-trip mismatch:\n%s\nvs\n%s", obj.Body, obj2.Body)
	}
}

func TestDecodeUpgradePolicyMissingName(t *testing.T) {
	input := `osType: linux
scopeLevel: tenant
`
	_, err := decodeUpgradePolicy([]byte(input))
	if err == nil {
		t.Fatal("expected error for upgrade policy without name")
	}
	if !strings.Contains(err.Error(), "no name") {
		t.Fatalf("error %q does not mention missing name", err)
	}
}

// --- LoadDir integration tests ---

func TestFilterLoadDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "infected.yaml"), []byte("name: Infected\nfilterFields:\n    infected: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "servers.yaml"), []byte("name: Servers\nfilterFields:\n    machineType: server\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	objs, err := reconcile.LoadDir(dir, decodeFilter)
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 2 {
		t.Fatalf("loaded %d objects, want 2", len(objs))
	}
}

func TestTagRuleLoadDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "prod.yaml"), []byte("name: Tag prod\nconditions:\n    op: and\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	objs, err := reconcile.LoadDir(dir, decodeTagRule)
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 1 {
		t.Fatalf("loaded %d objects, want 1", len(objs))
	}
	if objs[0].Name != "Tag prod" {
		t.Fatalf("Name = %q, want %q", objs[0].Name, "Tag prod")
	}
}

func TestUpgradePolicyLoadDir(t *testing.T) {
	dir := t.TempDir()
	f := upgradePolicyFile{
		Name:       "GA linux",
		OSType:     "linux",
		ScopeLevel: "tenant",
		MaxRetries: 3,
		Package:    upgradePolicyPkgFile{FileID: "f1", Major: "25", Minor: "3", Build: "1"},
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ga-linux.yaml"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	objs, err := reconcile.LoadDir(dir, decodeUpgradePolicy)
	if err != nil {
		t.Fatal(err)
	}
	if len(objs) != 1 {
		t.Fatalf("loaded %d objects, want 1", len(objs))
	}
	if objs[0].Name != "GA linux" {
		t.Fatalf("Name = %q, want %q", objs[0].Name, "GA linux")
	}
}
