package cli

import (
	"strings"
	"testing"
)

func TestAssetsListHelp(t *testing.T) {
	out, err := runCLI(t, "assets", "list", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, flag := range []string{"--type", "--filter", "--limit", "--sort-by", "--site-id", "--all"} {
		if !strings.Contains(out, flag) {
			t.Errorf("expected %s in help output", flag)
		}
	}
}

func TestAssetsExportHelp(t *testing.T) {
	out, err := runCLI(t, "assets", "export", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "--output-file") {
		t.Fatal("expected --output-file in help output")
	}
}

func TestAssetsNotesAddDryRun(t *testing.T) {
	out, err := runCLI(t, "assets", "notes", "add", "--asset-id", "A1", "--note", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "A1") {
		t.Fatalf("expected dry-run message with asset ID, got %q", out)
	}
}

func TestAssetsNotesDeleteDryRun(t *testing.T) {
	out, err := runCLI(t, "assets", "notes", "delete", "--note-id", "N1", "--asset-id", "A1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") || !strings.Contains(out, "N1") {
		t.Fatalf("expected dry-run message with note ID, got %q", out)
	}
}

func TestAssetsActionDryRun(t *testing.T) {
	out, err := runCLI(t, "assets", "action", "--action", "mark_asset_criticality_high", "--id", "D1,D2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Would") {
		t.Fatalf("expected dry-run message, got %q", out)
	}
}

func TestAssetsActionRequiresFlags(t *testing.T) {
	_, err := runCLI(t, "assets", "action")
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

func TestAssetsNotesRequiresFlags(t *testing.T) {
	_, err := runCLI(t, "assets", "notes", "add")
	if err == nil {
		t.Fatal("expected error for missing required flags")
	}
}

func TestAssetsFilterOptionsRequiresType(t *testing.T) {
	_, err := runCLI(t, "assets", "filter-options")
	if err == nil {
		t.Fatal("expected error for missing --type flag")
	}
}

func TestAssetsSubcommandHelp(t *testing.T) {
	out, err := runCLI(t, "assets", "list", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "--type") {
		t.Fatalf("expected --type in help, got %q", out)
	}
}

func TestParseFilterKV(t *testing.T) {
	vals, err := parseFilterKV([]string{"key1=val1", "key2=val2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vals.Get("key1") != "val1" {
		t.Fatalf("expected key1=val1, got %s", vals.Get("key1"))
	}
	if vals.Get("key2") != "val2" {
		t.Fatalf("expected key2=val2, got %s", vals.Get("key2"))
	}
}

func TestParseFilterKVInvalid(t *testing.T) {
	_, err := parseFilterKV([]string{"noequals"})
	if err == nil {
		t.Fatal("expected error for invalid filter")
	}
}
