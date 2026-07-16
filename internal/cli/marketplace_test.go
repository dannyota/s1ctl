package cli

import (
	"strings"
	"testing"
)

func TestMarketplaceCatalogHelp(t *testing.T) {
	out, err := runCLI(t, "marketplace", "catalog", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "catalog") {
		t.Fatalf("expected catalog in help, got %q", out)
	}
}

func TestMarketplaceListHelp(t *testing.T) {
	out, err := runCLI(t, "marketplace", "list", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "installed") {
		t.Fatalf("expected installed in help, got %q", out)
	}
}

func TestMarketplaceCatalogConfigRequiresArg(t *testing.T) {
	_, err := runCLI(t, "marketplace", "catalog-config")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestMarketplaceConfigRequiresArg(t *testing.T) {
	_, err := runCLI(t, "marketplace", "config")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestMarketplaceLogRequiresArg(t *testing.T) {
	_, err := runCLI(t, "marketplace", "log")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestMarketplaceInstallRequiresCatalogID(t *testing.T) {
	_, err := runCLI(t, "marketplace", "install", "--name", "test")
	if err == nil || !strings.Contains(err.Error(), "--catalog-id is required") {
		t.Fatalf("expected --catalog-id required error, got %v", err)
	}
}

func TestMarketplaceInstallRequiresName(t *testing.T) {
	_, err := runCLI(t, "marketplace", "install", "--catalog-id", "cat-1")
	if err == nil || !strings.Contains(err.Error(), "--name is required") {
		t.Fatalf("expected --name required error, got %v", err)
	}
}

func TestMarketplaceInstallDryRun(t *testing.T) {
	out, err := runCLI(t, "marketplace", "install", "--catalog-id", "cat-1", "--name", "My App")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would install") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
	if !strings.Contains(out, "My App") {
		t.Fatalf("expected app name in message, got %q", out)
	}
}

func TestMarketplaceInstallDryRunJSON(t *testing.T) {
	out, err := runCLI(t, "marketplace", "install", "--catalog-id", "cat-1", "--name", "My App", "--json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dryRun") {
		t.Fatalf("expected dryRun in JSON output, got %q", out)
	}
}

func TestMarketplaceInstallBadConfig(t *testing.T) {
	_, err := runCLI(t, "marketplace", "install", "--catalog-id", "cat-1", "--name", "Test", "--config", "noequalssign")
	if err == nil || !strings.Contains(err.Error(), "expected id=value") {
		t.Fatalf("expected config parse error, got %v", err)
	}
}

func TestMarketplaceUpdateRequiresID(t *testing.T) {
	_, err := runCLI(t, "marketplace", "update")
	if err == nil || !strings.Contains(err.Error(), "--id is required") {
		t.Fatalf("expected --id required error, got %v", err)
	}
}

func TestMarketplaceUpdateDryRun(t *testing.T) {
	out, err := runCLI(t, "marketplace", "update", "--id", "app-1", "--name", "New Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would update application app-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMarketplaceDeleteRequiresID(t *testing.T) {
	_, err := runCLI(t, "marketplace", "delete")
	if err == nil || !strings.Contains(err.Error(), "--id is required") {
		t.Fatalf("expected --id required error, got %v", err)
	}
}

func TestMarketplaceDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "marketplace", "delete", "--id", "app-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete application app-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMarketplaceEnableRequiresArg(t *testing.T) {
	_, err := runCLI(t, "marketplace", "enable")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestMarketplaceEnableDryRun(t *testing.T) {
	out, err := runCLI(t, "marketplace", "enable", "app-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would enable application app-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMarketplaceDisableRequiresArg(t *testing.T) {
	_, err := runCLI(t, "marketplace", "disable")
	if err == nil || !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Fatalf("expected 1-arg error, got %v", err)
	}
}

func TestMarketplaceDisableDryRun(t *testing.T) {
	out, err := runCLI(t, "marketplace", "disable", "app-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would disable application app-1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestMarketplaceRequiresSubcommand(t *testing.T) {
	out, err := runCLI(t, "marketplace")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, sub := range []string{"catalog", "list", "install", "delete", "enable", "disable"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected help to list %q subcommand, got %q", sub, out)
		}
	}
}
