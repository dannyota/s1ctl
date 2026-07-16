package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloudOnboardingOnboardDryRun(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "onboard.json")
	if err := os.WriteFile(file, []byte(`{"onBoardingType":"INDIVIDUAL","cloudProvider":"AWS","products":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	out, err := runCLI(t, "cloud-onboarding", "onboard", "--from-file", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would onboard cloud entity from") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestCloudOnboardingOnboardRequiresFile(t *testing.T) {
	_, err := runCLI(t, "cloud-onboarding", "onboard")
	if err == nil || !strings.Contains(err.Error(), "--from-file is required") {
		t.Fatalf("expected --from-file is required, got %v", err)
	}
}

func TestCloudOnboardingDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "cloud-onboarding", "delete", "acc-1", "acc-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete 2 cloud entities") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestCloudOnboardingDeleteSingular(t *testing.T) {
	out, err := runCLI(t, "cloud-onboarding", "delete", "acc-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would delete 1 cloud entity") {
		t.Fatalf("expected dry-run singular, got %q", out)
	}
}

func TestCloudOnboardingDeleteRequiresID(t *testing.T) {
	_, err := runCLI(t, "cloud-onboarding", "delete")
	if err == nil || !strings.Contains(err.Error(), "at least 1") {
		t.Fatalf("expected 'at least 1' arg error, got %v", err)
	}
}

func TestCloudOnboardingGetRequiresID(t *testing.T) {
	_, err := runCLI(t, "cloud-onboarding", "get")
	if err == nil {
		t.Fatal("expected error for missing account id, got nil")
	}
}
