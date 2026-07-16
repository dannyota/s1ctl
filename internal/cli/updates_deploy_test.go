package cli

import (
	"strings"
	"testing"
)

func TestDeployCreateGroupDryRun(t *testing.T) {
	out, err := runCLI(t, "updates", "deploy", "create-group",
		"--group-name", "test-group",
		"--group-passphrase", "enc-pass",
		"--scope-id", "225494730938493804",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDeployCreateGroupRequiresName(t *testing.T) {
	_, err := runCLI(t, "updates", "deploy", "create-group",
		"--group-passphrase", "enc",
		"--scope-id", "123",
	)
	if err == nil {
		t.Fatal("expected validation error without --group-name")
	}
	if !strings.Contains(err.Error(), "--group-name is required") {
		t.Fatalf("expected %q, got %q", "--group-name is required", err.Error())
	}
}

func TestDeployCreateGroupRequiresPassphrase(t *testing.T) {
	_, err := runCLI(t, "updates", "deploy", "create-group",
		"--group-name", "x",
		"--scope-id", "123",
	)
	if err == nil {
		t.Fatal("expected validation error without --group-passphrase")
	}
	if !strings.Contains(err.Error(), "--group-passphrase is required") {
		t.Fatalf("expected %q, got %q", "--group-passphrase is required", err.Error())
	}
}

func TestDeployCreateGroupRequiresScopeID(t *testing.T) {
	_, err := runCLI(t, "updates", "deploy", "create-group",
		"--group-name", "x",
		"--group-passphrase", "enc",
	)
	if err == nil {
		t.Fatal("expected validation error without --scope-id")
	}
	if !strings.Contains(err.Error(), "--scope-id is required") {
		t.Fatalf("expected %q, got %q", "--scope-id is required", err.Error())
	}
}

func TestDeployDeleteGroupDryRun(t *testing.T) {
	out, err := runCLI(t, "updates", "deploy", "delete-group", "G1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDeployAddDetailDryRun(t *testing.T) {
	out, err := runCLI(t, "updates", "deploy", "add-detail",
		"--cred-group-id", "G1",
		"--title", "Admin",
		"--cred-type", "User/Password",
		"--encrypted-key", "key1",
		"--encrypted-cred", "cred1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDeployAddDetailRequiresFields(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{"missing cred-group-id", []string{"--title", "x", "--cred-type", "y", "--encrypted-key", "k", "--encrypted-cred", "c"}, "--cred-group-id is required"},
		{"missing title", []string{"--cred-group-id", "G1", "--cred-type", "y", "--encrypted-key", "k", "--encrypted-cred", "c"}, "--title is required"},
		{"missing cred-type", []string{"--cred-group-id", "G1", "--title", "x", "--encrypted-key", "k", "--encrypted-cred", "c"}, "--cred-type is required"},
		{"missing encrypted-key", []string{"--cred-group-id", "G1", "--title", "x", "--cred-type", "y", "--encrypted-cred", "c"}, "--encrypted-key is required"},
		{"missing encrypted-cred", []string{"--cred-group-id", "G1", "--title", "x", "--cred-type", "y", "--encrypted-key", "k"}, "--encrypted-cred is required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"updates", "deploy", "add-detail"}, tt.args...)
			_, err := runCLI(t, args...)
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestDeployUpdateDetailDryRun(t *testing.T) {
	out, err := runCLI(t, "updates", "deploy", "update-detail", "D1",
		"--title", "Updated",
		"--cred-type", "User/Password",
		"--encrypted-key", "k",
		"--encrypted-cred", "c",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestDeployDeleteDetailDryRun(t *testing.T) {
	out, err := runCLI(t, "updates", "deploy", "delete-detail", "D1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}
