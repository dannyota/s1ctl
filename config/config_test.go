package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("S1_CONSOLE_URL", "https://test.sentinelone.net")
	t.Setenv("S1_TOKEN", "test-token")

	inst, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if inst.ConsoleURL != "https://test.sentinelone.net" {
		t.Errorf("ConsoleURL = %q", inst.ConsoleURL)
	}
	if inst.Token != "test-token" {
		t.Errorf("Token = %q", inst.Token)
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Setenv("S1_CONSOLE_URL", "")
	t.Setenv("S1_TOKEN", "")

	path := filepath.Join(t.TempDir(), "config.yaml")
	os.WriteFile(path, []byte("console_url: https://file.sentinelone.net\ntoken: file-token\n"), 0o600)

	inst, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if inst.ConsoleURL != "https://file.sentinelone.net" {
		t.Errorf("ConsoleURL = %q", inst.ConsoleURL)
	}
	if inst.Token != "file-token" {
		t.Errorf("Token = %q", inst.Token)
	}
	if inst.source != path {
		t.Errorf("source = %q, want %q", inst.source, path)
	}
}

func TestEnvOverridesFile(t *testing.T) {
	t.Setenv("S1_CONSOLE_URL", "https://env.sentinelone.net")
	t.Setenv("S1_TOKEN", "")

	path := filepath.Join(t.TempDir(), "config.yaml")
	os.WriteFile(path, []byte("console_url: https://file.sentinelone.net\ntoken: file-token\n"), 0o600)

	inst, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if inst.ConsoleURL != "https://env.sentinelone.net" {
		t.Errorf("ConsoleURL = %q, want env override", inst.ConsoleURL)
	}
	if inst.Token != "file-token" {
		t.Errorf("Token = %q, want file value", inst.Token)
	}
}

func TestTrailingSlashTrimmed(t *testing.T) {
	t.Setenv("S1_CONSOLE_URL", "https://test.sentinelone.net/")
	t.Setenv("S1_TOKEN", "t")

	inst, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if inst.ConsoleURL != "https://test.sentinelone.net" {
		t.Errorf("trailing slash not trimmed: %q", inst.ConsoleURL)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		inst    Instance
		wantErr bool
	}{
		{"empty", Instance{}, true},
		{"no token", Instance{ConsoleURL: "https://x.sentinelone.net"}, true},
		{"no url", Instance{Token: "t"}, true},
		{"valid", Instance{ConsoleURL: "https://x.sentinelone.net", Token: "t"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.inst.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	t.Setenv("S1_CONSOLE_URL", "")
	t.Setenv("S1_TOKEN", "")

	path := filepath.Join(t.TempDir(), "sub", "config.yaml")
	inst := &Instance{ConsoleURL: "https://test.sentinelone.net", Token: "test-token"}
	if err := Save(path, inst); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ConsoleURL != inst.ConsoleURL || loaded.Token != inst.Token {
		t.Errorf("round-trip mismatch: got %+v", loaded)
	}
}

func TestResolvedSourceMissing(t *testing.T) {
	if got := ResolvedSource("/nonexistent/path.yaml"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
