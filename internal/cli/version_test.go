package cli

import (
	"strings"
	"testing"
)

func TestResolveBuildInfoFallback(t *testing.T) {
	bi := resolveBuildInfo()
	if bi.Version == "" {
		t.Error("version must never be empty")
	}
	if bi.Commit == "" {
		t.Error("commit must never be empty")
	}
	if bi.GoVersion == "" || bi.OS == "" || bi.Arch == "" {
		t.Errorf("runtime fields unset: %+v", bi)
	}
}

func TestShortCommit(t *testing.T) {
	if got := shortCommit("0123456789abcdef0123"); got != "0123456789ab" {
		t.Errorf("shortCommit long = %q, want 12 chars", got)
	}
	if got := shortCommit("abc"); got != "abc" {
		t.Errorf("shortCommit short = %q, want unchanged", got)
	}
}

func TestVersionLine(t *testing.T) {
	l := versionLine()
	if !strings.HasPrefix(l, "s1ctl ") || strings.Contains(l, "\n") {
		t.Errorf("versionLine malformed: %q", l)
	}
}
