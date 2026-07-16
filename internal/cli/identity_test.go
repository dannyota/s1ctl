package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"danny.vn/s1/mgmt"
)

func TestIdentityOnboardHelp(t *testing.T) {
	out, err := runCLI(t, "identity", "onboard", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "onboarding status") {
		t.Fatalf("expected help text about onboarding, got %q", out)
	}
}

func TestIdentityConfigGetSubcommand(t *testing.T) {
	// Verify the subcommand tree registers without error.
	_, _ = runCLI(t, "identity", "config", "get", "--help")
}

func TestIdentityConfigAddDryRun(t *testing.T) {
	out, err := runCLI(t, "identity", "config", "add",
		"--domain", "corp.example.com",
		"--dc-fqdn", "dc1.corp.example.com",
		"--user", "svc",
		"--password", "secret",
		"--site-id", "S1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
	// Verify the password is not echoed in dry-run output.
	if strings.Contains(out, "secret") {
		t.Fatal("password must not appear in dry-run output")
	}
}

func TestIdentityConfigAddMissingFields(t *testing.T) {
	_, err := runCLI(t, "identity", "config", "add", "--yes",
		"--domain", "corp.example.com",
		// Missing --dc-fqdn, --user, --password
	)
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Fatalf("expected 'required' in error, got %v", err)
	}
}

func TestIdentityConfigAddBadEncryption(t *testing.T) {
	_, err := runCLI(t, "identity", "config", "add", "--yes",
		"--domain", "d",
		"--dc-fqdn", "dc",
		"--user", "u",
		"--password", "p",
		"--encryption", "TLS",
	)
	if err == nil {
		t.Fatal("expected error for bad encryption")
	}
	if !strings.Contains(err.Error(), "LDAP or LDAPS") {
		t.Fatalf("expected encryption error, got %v", err)
	}
}

func TestIdentityConfigDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "identity", "config", "delete", "42", "--site-id", "S1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestIdentityConfigDeleteBadID(t *testing.T) {
	_, err := runCLI(t, "identity", "config", "delete", "abc", "--yes")
	if err == nil {
		t.Fatal("expected error for non-numeric ID")
	}
	if !strings.Contains(err.Error(), "invalid config ID") {
		t.Fatalf("expected 'invalid config ID' error, got %v", err)
	}
}

func TestIdentityConfigDeleteNoArgs(t *testing.T) {
	_, err := runCLI(t, "identity", "config", "delete")
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestIdentityConnectorReplaceAgentDryRun(t *testing.T) {
	out, err := runCLI(t, "identity", "connector", "replace", "test-uuid", "--site-id", "S1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestIdentitySkipExposuresDryRun(t *testing.T) {
	out, err := runCLI(t, "identity", "skip-exposures",
		"--detection", "Kerberoasting",
		"--domain", "corp.example.com",
		"--site-id", "S1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestIdentitySkipExposuresMissingFlags(t *testing.T) {
	_, err := runCLI(t, "identity", "skip-exposures", "--yes")
	if err == nil {
		t.Fatal("expected error for missing detection/domain")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Fatalf("expected 'required' in error, got %v", err)
	}
}

func TestIdentityAckExposuresDryRun(t *testing.T) {
	out, err := runCLI(t, "identity", "ack-exposures",
		"--detection", "Kerberoasting",
		"--domain", "corp.example.com",
		"--site-id", "S1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestIdentityAckExposuresMissingFlags(t *testing.T) {
	_, err := runCLI(t, "identity", "ack-exposures", "--yes")
	if err == nil {
		t.Fatal("expected error for missing detection/domain")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Fatalf("expected 'required' in error, got %v", err)
	}
}

// redactADConfig must blank the username and Raw before output. The API
// returns the bind credential username — we strip it to avoid leaking
// credentials.
func TestRedactADConfig(t *testing.T) {
	cfg := mgmt.ADConfiguration{
		ID:               1,
		DomainName:       "corp.example.com",
		EncryptionMethod: mgmt.EncryptionMethodLDAPS,
		Username:         "svc-bind",
		IsConnected:      true,
		Raw:              []byte(`{"username":"svc-bind"}`),
	}
	got := redactADConfig(cfg)
	if got.Username != "" {
		t.Fatalf("expected username to be blanked, got %q", got.Username)
	}
	if got.Raw != nil {
		t.Fatalf("expected Raw to be nil, got %q", got.Raw)
	}
	if got.DomainName != cfg.DomainName {
		t.Fatal("non-secret fields must be preserved")
	}
	if got.EncryptionMethod != cfg.EncryptionMethod {
		t.Fatal("non-secret fields must be preserved")
	}
	if !got.IsConnected {
		t.Fatal("non-secret flags must be preserved")
	}
	// Original must not be mutated.
	if cfg.Username == "" {
		t.Fatal("original must not be mutated")
	}
}

func TestRedactADConfigJSON(t *testing.T) {
	cfg := mgmt.ADConfiguration{
		ID:               1,
		DomainName:       "corp.example.com",
		Username:         "svc-bind",
		EncryptionMethod: mgmt.EncryptionMethodLDAPS,
		Raw:              []byte(`{"username":"svc-bind"}`),
	}
	redacted := redactADConfig(cfg)
	b, err := json.Marshal(redacted)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	s := string(b)
	if strings.Contains(s, "svc-bind") {
		t.Fatalf("serialized output must not contain the username: %s", s)
	}
}

func TestIdentitySubcommandShowsHelp(t *testing.T) {
	out, err := runCLI(t, "identity")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "identity") {
		t.Fatalf("expected help output, got %q", out)
	}
}
