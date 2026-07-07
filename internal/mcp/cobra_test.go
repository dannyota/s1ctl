package mcp

import (
	"testing"

	"github.com/spf13/cobra"
)

func testCobraTree() *cobra.Command {
	root := &cobra.Command{Use: "test"}

	agents := &cobra.Command{Use: "agents", Short: "Manage agents"}
	agents.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List agents",
		RunE:  func(_ *cobra.Command, _ []string) error { return nil },
	})

	isolate := &cobra.Command{
		Use:   "isolate",
		Short: "Isolate an agent",
		RunE:  func(_ *cobra.Command, _ []string) error { return nil },
	}
	isolate.Flags().Bool("yes", false, "apply the mutation")
	isolate.Flags().String("agent-id", "", "agent ID")
	agents.AddCommand(isolate)

	root.AddCommand(agents)

	// hidden command should be skipped
	root.AddCommand(&cobra.Command{
		Use:    "internal",
		Hidden: true,
		RunE:   func(_ *cobra.Command, _ []string) error { return nil },
	})

	// skipped top-level command
	root.AddCommand(&cobra.Command{
		Use:  "completion",
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	})

	return root
}

func TestToolsFromCobra(t *testing.T) {
	root := testCobraTree()
	tools := ToolsFromCobra(root)

	names := map[string]bool{}
	for _, tool := range tools {
		names[tool.Name] = true
	}

	if !names["agents_list"] {
		t.Error("missing agents_list tool")
	}
	if !names["agents_isolate"] {
		t.Error("missing agents_isolate tool")
	}
	if names["internal"] {
		t.Error("hidden command should not generate a tool")
	}
	if names["completion"] {
		t.Error("completion should be skipped")
	}
}

func TestMutationDescription(t *testing.T) {
	root := testCobraTree()
	tools := ToolsFromCobra(root)

	for _, tool := range tools {
		if tool.Name != "agents_isolate" {
			continue
		}
		want := "[mutation: requires --yes to apply, dry-run by default]"
		if len(tool.Description) < len(want) {
			t.Fatalf("description too short: %q", tool.Description)
		}
		suffix := tool.Description[len(tool.Description)-len(want):]
		if suffix != want {
			t.Errorf("mutation suffix = %q, want %q", suffix, want)
		}
		return
	}
	t.Fatal("agents_isolate tool not found")
}

func TestInputSchemaFlags(t *testing.T) {
	root := testCobraTree()
	tools := ToolsFromCobra(root)

	for _, tool := range tools {
		if tool.Name != "agents_isolate" {
			continue
		}
		props, _ := tool.InputSchema["properties"].(map[string]any)
		if props == nil {
			t.Fatal("no properties in schema")
		}

		yes, _ := props["yes"].(map[string]any)
		if yes == nil {
			t.Fatal("missing --yes in schema")
		}
		if yes["type"] != "boolean" {
			t.Errorf("--yes type = %v, want boolean", yes["type"])
		}

		agentID, _ := props["agent-id"].(map[string]any)
		if agentID == nil {
			t.Fatal("missing --agent-id in schema")
		}
		if agentID["type"] != "string" {
			t.Errorf("--agent-id type = %v, want string", agentID["type"])
		}

		if props["help"] != nil {
			t.Error("--help should be excluded from schema")
		}
		return
	}
	t.Fatal("agents_isolate tool not found")
}

func TestGroupTools(t *testing.T) {
	root := testCobraTree()
	tools, err := GroupTools(root, "agents")
	if err != nil {
		t.Fatal(err)
	}

	names := map[string]bool{}
	for _, tool := range tools {
		names[tool.Name] = true
	}
	if !names["agents_list"] {
		t.Error("missing agents_list")
	}
	if !names["agents_isolate"] {
		t.Error("missing agents_isolate")
	}
	if len(tools) != 2 {
		t.Errorf("got %d tools, want 2", len(tools))
	}
}

func TestGroupToolsUnknown(t *testing.T) {
	root := testCobraTree()
	_, err := GroupTools(root, "nonexistent")
	if err == nil {
		t.Error("expected error for unknown group")
	}
}

func TestShortDescriptionOnly(t *testing.T) {
	root := &cobra.Command{Use: "test"}
	leaf := &cobra.Command{
		Use:   "thing",
		Short: "Do the thing",
		Long:  "This is a very long description that should not appear in the tool.",
		RunE:  func(_ *cobra.Command, _ []string) error { return nil },
	}
	root.AddCommand(leaf)

	tools := ToolsFromCobra(root)
	if len(tools) != 1 {
		t.Fatalf("got %d tools, want 1", len(tools))
	}
	if tools[0].Description != "Do the thing" {
		t.Errorf("description = %q, want %q", tools[0].Description, "Do the thing")
	}
}

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"agents list", []string{"agents", "list"}},
		{`datalake facet --filter 'event.type = "Login"'`, []string{"datalake", "facet", "--filter", `event.type = "Login"`}},
		{`datalake facet --filter "event.type = 'Login'"`, []string{"datalake", "facet", "--filter", "event.type = 'Login'"}},
		{`--filter "A AND B" --field x`, []string{"--filter", "A AND B", "--field", "x"}},
		{`say hello\ world`, []string{"say", "hello world"}},
		{"  spaces  everywhere  ", []string{"spaces", "everywhere"}},
		{"", nil},
	}
	for _, tt := range tests {
		got := splitCommand(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("splitCommand(%q) = %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitCommand(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestSkippedGlobalFlags(t *testing.T) {
	root := &cobra.Command{Use: "test"}
	root.PersistentFlags().String("output", "table", "output format")
	root.PersistentFlags().Bool("json", false, "json output")
	root.PersistentFlags().Bool("verbose", false, "verbose")
	root.PersistentFlags().Bool("no-progress", false, "no progress")
	root.PersistentFlags().String("config", "", "config")

	leaf := &cobra.Command{
		Use:  "run",
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	}
	leaf.Flags().String("name", "", "resource name")
	root.AddCommand(leaf)

	tools := ToolsFromCobra(root)
	if len(tools) != 1 {
		t.Fatalf("got %d tools, want 1", len(tools))
	}

	props, _ := tools[0].InputSchema["properties"].(map[string]any)
	for _, skip := range []string{"output", "json", "verbose", "no-progress", "config", "help"} {
		if props[skip] != nil {
			t.Errorf("global flag --%s should be excluded", skip)
		}
	}
	if props["name"] == nil {
		t.Error("local flag --name should be included")
	}
}
