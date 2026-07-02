package cli

import (
	"strings"
	"testing"

	"danny.vn/s1/mgmt"
)

// redactADSettings must blank the bind password before output. The AD GET
// endpoint does not normally echo the password; this guards against an API that
// returns it anyway.
func TestRedactADSettings(t *testing.T) {
	s := &mgmt.ADSettings{
		Enabled:  true,
		Host:     "ad.example.com",
		Port:     636,
		Username: "svc-bind",
		RootDN:   "dc=example,dc=com",
		SSL:      true,
		Password: "bind-secret-placeholder",
		Raw:      []byte(`{"password":"bind-secret-placeholder"}`),
	}
	got := redactADSettings(s)
	if got.Password != "" {
		t.Fatalf("expected password to be blanked, got %q", got.Password)
	}
	if got.Raw != nil {
		t.Fatalf("expected Raw to be nil, got %q", got.Raw)
	}
	if got.Host != s.Host || got.Port != s.Port || got.Username != s.Username || got.RootDN != s.RootDN {
		t.Fatal("non-secret fields must be preserved")
	}
	if !got.Enabled || !got.SSL {
		t.Fatal("non-secret flags must be preserved")
	}
	if s.Password == "" {
		t.Fatal("original must not be mutated")
	}
}

func TestSettingsTestADDryRun(t *testing.T) {
	out, err := runCLI(t, "settings", "test", "ad", "--site-id", "S1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestSettingsTestUnknownType(t *testing.T) {
	if _, err := runCLI(t, "settings", "test", "bogus"); err == nil || !strings.Contains(err.Error(), "unknown settings type") {
		t.Fatalf("expected unknown type error, got %v", err)
	}
}

func TestSettingsCancelPendingEmailsDryRun(t *testing.T) {
	out, err := runCLI(t, "settings", "cancel-pending-emails", "--site-id", "S1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestSettingsDeleteRecipientDryRun(t *testing.T) {
	out, err := runCLI(t, "settings", "delete-recipient", "999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestSettingsDeleteRecipientArgValidation(t *testing.T) {
	if _, err := runCLI(t, "settings", "delete-recipient"); err == nil {
		t.Fatal("expected error for missing recipient id, got nil")
	}
}
