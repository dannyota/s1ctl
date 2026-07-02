package cli

import (
	"strings"
	"testing"

	"danny.vn/s1/mgmt"
)

func TestUsersGenerateTokenDryRun(t *testing.T) {
	out, err := runCLI(t, "users", "generate-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would generate API token for current user") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestUsersRevokeTokenDryRun(t *testing.T) {
	out, err := runCLI(t, "users", "revoke-token", "U1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would revoke API token for user U1") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestUsersTokenDetailsArgValidation(t *testing.T) {
	if _, err := runCLI(t, "users", "token-details", "U1", "U2"); err == nil {
		t.Fatal("expected error for too many args, got nil")
	}
}

// redactUserTokenDetails must blank any secret token value before output. The
// details endpoints normally return timestamps only; this guards against an API
// that echoes the secret.
func TestRedactUserTokenDetails(t *testing.T) {
	d := &mgmt.UserTokenDetails{
		CreatedAt: "2026-01-01T00:00:00Z",
		ExpiresAt: "2027-01-01T00:00:00Z",
		Token:     "super-secret-token",
		Raw:       []byte(`{"token":"super-secret-token"}`),
	}
	got := redactUserTokenDetails(d)
	if got.Token != "" {
		t.Fatalf("expected token to be blanked, got %q", got.Token)
	}
	if got.Raw != nil {
		t.Fatalf("expected Raw to be nil, got %q", got.Raw)
	}
	if got.CreatedAt != d.CreatedAt || got.ExpiresAt != d.ExpiresAt {
		t.Fatal("timestamps must be preserved")
	}
	if d.Token == "" {
		t.Fatal("original must not be mutated")
	}
}
