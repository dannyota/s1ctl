package cli

import (
	"strings"
	"testing"
)

func TestSitesLifecycleDryRun(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"sites", "reactivate", "S1", "--unlimited"}, "Would reactivate site S1"},
		{[]string{"sites", "reactivate", "S1", "--expiration", "2027-01-01T00:00:00Z"}, "Would reactivate site S1"},
		{[]string{"sites", "expire", "S1"}, "Would expire site S1"},
		{[]string{"sites", "duplicate", "--name", "clone", "--source-site-id", "42"}, "Would duplicate site from 42 as \"clone\""},
		{[]string{"sites", "regenerate-key", "S1"}, "Would regenerate registration key for site S1"},
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

func TestSitesDuplicateValidation(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"sites", "duplicate"}, "--name is required"},
		{[]string{"sites", "duplicate", "--name", "clone"}, "--source-site-id is required"},
		{[]string{"sites", "duplicate", "--name", "clone", "--source-site-id", "abc"}, "must be numeric"},
		{[]string{"sites", "duplicate", "--name", "clone", "--source-site-id", "42", "--policy-source", "bogus"}, "--policy-source must be one of"},
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

func TestSitesReactivateFlagValidation(t *testing.T) {
	cases := [][]string{
		{"sites", "reactivate", "S1"},
		{"sites", "reactivate", "S1", "--unlimited", "--expiration", "2027-01-01T00:00:00Z"},
	}
	for _, args := range cases {
		_, err := runCLI(t, args...)
		if err == nil {
			t.Fatalf("%v: expected validation error", args)
		}
		if !strings.Contains(err.Error(), "exactly one of --unlimited or --expiration") {
			t.Fatalf("%v: unexpected error: %v", args, err)
		}
	}
}

func TestSitesLifecycleArgValidation(t *testing.T) {
	cases := [][]string{
		{"sites", "reactivate"},
		{"sites", "expire"},
		{"sites", "regenerate-key"},
		{"sites", "token"},
	}
	for _, args := range cases {
		if _, err := runCLI(t, args...); err == nil {
			t.Fatalf("%v: expected arg validation error", args)
		}
	}
}
