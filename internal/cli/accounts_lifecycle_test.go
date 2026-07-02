package cli

import (
	"strings"
	"testing"
)

func TestAccountsLifecycleDryRun(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"accounts", "reactivate", "A1", "--unlimited"}, "Would reactivate account A1"},
		{[]string{"accounts", "reactivate", "A1", "--expiration", "2027-01-01T00:00:00Z"}, "Would reactivate account A1"},
		{[]string{"accounts", "expire", "A1"}, "Would expire account A1"},
		{[]string{"accounts", "uninstall-password", "generate", "A1", "--expiration", "2027-01-01"}, "Would generate uninstall password for account A1"},
		{[]string{"accounts", "uninstall-password", "revoke", "A1"}, "Would revoke uninstall password for account A1"},
	}
	for _, tc := range cases {
		out, err := runCLI(t, tc.args...)
		if err != nil {
			t.Fatalf("%v: unexpected error: %v", tc.args, err)
		}
		if !strings.Contains(out, tc.want) {
			t.Fatalf("%v: expected %q, got %q", tc.args, tc.want, out)
		}
	}
}

func TestAccountsLifecycleFlagValidation(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"accounts", "reactivate", "A1"}, "exactly one of --unlimited or --expiration"},
		{[]string{"accounts", "reactivate", "A1", "--unlimited", "--expiration", "2027-01-01T00:00:00Z"}, "exactly one of --unlimited or --expiration"},
		{[]string{"accounts", "uninstall-password", "generate", "A1"}, "--expiration is required"},
	}
	for _, tc := range cases {
		_, err := runCLI(t, tc.args...)
		if err == nil {
			t.Fatalf("%v: expected validation error", tc.args)
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("%v: expected error %q, got %v", tc.args, tc.want, err)
		}
	}
}

func TestAccountsLifecycleArgValidation(t *testing.T) {
	cases := [][]string{
		{"accounts", "reactivate"},
		{"accounts", "expire"},
		{"accounts", "uninstall-password", "show"},
		{"accounts", "uninstall-password", "generate"},
		{"accounts", "uninstall-password", "revoke"},
	}
	for _, args := range cases {
		if _, err := runCLI(t, args...); err == nil {
			t.Fatalf("%v: expected arg validation error", args)
		}
	}
}

// TestUninstallPasswordShowIsRead confirms the read-side `show` command is not
// gated by the mutation guard (no "Would" dry-run wording); it still requires an
// account ID argument. Actual output is exercised via SDK tests to avoid a live
// client.
func TestUninstallPasswordShowRequiresArg(t *testing.T) {
	if _, err := runCLI(t, "accounts", "uninstall-password", "show"); err == nil {
		t.Fatal("expected arg validation error for show without id")
	}
}
