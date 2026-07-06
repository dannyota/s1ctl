package mcp

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Tool struct {
	Name        string
	Description string
	InputSchema map[string]any
	Run         func(args map[string]any) (string, error)
}

var skipCommands = map[string]bool{
	"mcp": true, "completion": true, "docs": true, "help": true,
}

func ToolsFromCobra(root *cobra.Command) []Tool {
	var tools []Tool
	walkCommands(root, nil, &tools)
	return tools
}

func walkCommands(cmd *cobra.Command, path []string, tools *[]Tool) {
	for _, c := range cmd.Commands() {
		if c.Hidden || (len(path) == 0 && skipCommands[c.Name()]) {
			continue
		}
		cur := append(append([]string(nil), path...), c.Name())

		if c.HasSubCommands() {
			walkCommands(c, cur, tools)
			continue
		}

		if c.RunE == nil && c.Run == nil {
			continue
		}

		t := buildTool(c, cur)
		*tools = append(*tools, t)
	}
}

var skipFlags = map[string]bool{
	"help": true, "json": true, "output": true,
	"verbose": true, "no-progress": true, "config": true,
}

func buildTool(cmd *cobra.Command, path []string) Tool {
	name := strings.Join(path, "_")

	desc := cmd.Short
	if hasMutationFlag(cmd) {
		desc += " [mutation: requires --yes to apply, dry-run by default]"
	}

	schema := buildInputSchema(cmd)

	return Tool{
		Name:        name,
		Description: desc,
		InputSchema: schema,
		Run:         makeRunner(cmd, path),
	}
}

func buildInputSchema(cmd *cobra.Command) map[string]any {
	properties := map[string]any{}
	var required []string

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if skipFlags[f.Name] {
			return
		}
		if f.Hidden {
			return
		}

		prop := flagToProperty(f)
		properties[f.Name] = prop

		ann := cmd.Flag(f.Name)
		if ann != nil {
			if _, ok := ann.Annotations[cobra.BashCompOneRequiredFlag]; ok {
				required = append(required, f.Name)
			}
		}
	})

	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if skipFlags[f.Name] {
			return
		}
		if f.Hidden {
			return
		}
		if f.Name == "read-only" {
			properties[f.Name] = flagToProperty(f)
		}
	})

	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func flagToProperty(f *pflag.Flag) map[string]any {
	prop := map[string]any{}

	switch f.Value.Type() {
	case "bool":
		prop["type"] = "boolean"
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		prop["type"] = "integer"
	case "float32", "float64":
		prop["type"] = "number"
	case "stringSlice", "stringArray":
		prop["type"] = "array"
		prop["items"] = map[string]any{"type": "string"}
	default:
		prop["type"] = "string"
	}

	if f.Usage != "" {
		prop["description"] = f.Usage
	}
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
		prop["default"] = f.DefValue
	}

	return prop
}

func hasMutationFlag(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup("yes") != nil
}

func makeRunner(cmd *cobra.Command, path []string) func(map[string]any) (string, error) {
	return func(args map[string]any) (string, error) {
		root := cmd.Root()

		cliArgs := append([]string(nil), path...)
		cliArgs = append(cliArgs, "--json", "--no-progress")

		for k, v := range args {
			flag := cmd.Flags().Lookup(k)
			if flag == nil {
				flag = cmd.InheritedFlags().Lookup(k)
			}
			if flag == nil {
				continue
			}

			switch flag.Value.Type() {
			case "bool":
				b, _ := toBool(v)
				if b {
					cliArgs = append(cliArgs, "--"+k)
				}
			case "stringSlice", "stringArray":
				if arr, ok := v.([]any); ok {
					for _, item := range arr {
						cliArgs = append(cliArgs, "--"+k, fmt.Sprint(item))
					}
				}
			default:
				cliArgs = append(cliArgs, "--"+k, fmt.Sprint(v))
			}
		}

		// Capture both cobra output and raw os.Stdout writes.
		// Many commands write to os.Stdout directly via fmt.Printf;
		// os.Pipe intercepts those writes too.
		origStdout := os.Stdout
		pr, pw, err := os.Pipe()
		if err != nil {
			return "", fmt.Errorf("create pipe: %w", err)
		}
		os.Stdout = pw

		var stderr bytes.Buffer
		root.SetArgs(cliArgs)
		root.SetOut(pw)
		root.SetErr(&stderr)

		execErr := root.Execute()

		_ = pw.Close()
		os.Stdout = origStdout

		var stdout bytes.Buffer
		_, _ = stdout.ReadFrom(pr)
		_ = pr.Close()

		out := stdout.String()
		if errOut := stderr.String(); errOut != "" && out == "" {
			out = errOut
		}
		if out == "" && execErr != nil {
			out = execErr.Error()
		}

		return out, execErr
	}
}

func toBool(v any) (bool, bool) {
	switch b := v.(type) {
	case bool:
		return b, true
	case string:
		val, err := strconv.ParseBool(b)
		return val, err == nil
	default:
		return false, false
	}
}

// GroupTools returns tools for a single top-level command group.
func GroupTools(root *cobra.Command, group string) ([]Tool, error) {
	for _, c := range root.Commands() {
		if c.Name() == group && !c.Hidden {
			var tools []Tool
			if c.RunE != nil || c.Run != nil {
				tools = append(tools, buildTool(c, []string{group}))
			}
			walkCommands(c, []string{group}, &tools)
			return tools, nil
		}
	}
	return nil, fmt.Errorf("unknown group: %s", group)
}

func (s *Server) buildMetaTools() []Tool {
	return []Tool{
		s.buildRunTool(),
		s.buildHelpTool(),
		s.buildFocusTool(),
		s.buildUnfocusTool(),
	}
}

func (s *Server) buildRunTool() Tool {
	return Tool{
		Name:        "run",
		Description: "Run any s1ctl command. Pass the full command (e.g. 'agents list --site-id 123').",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "s1ctl command to run (without the 's1ctl' prefix)",
				},
			},
			"required": []string{"command"},
		},
		Run: func(args map[string]any) (string, error) {
			cmdStr, _ := args["command"].(string)
			if cmdStr == "" {
				return "", fmt.Errorf("command is required")
			}
			return s.execCommand(strings.Fields(cmdStr))
		},
	}
}

func (s *Server) buildHelpTool() Tool {
	return Tool{
		Name:        "help",
		Description: "List available command groups, or subcommands within a group.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"group": map[string]any{
					"type":        "string",
					"description": "group name (e.g. 'agents'). Omit to list all groups.",
				},
			},
		},
		Run: func(args map[string]any) (string, error) {
			group, _ := args["group"].(string)
			return s.helpOutput(group)
		},
	}
}

func (s *Server) buildFocusTool() Tool {
	return Tool{
		Name:        "focus",
		Description: "Load typed tools for a command group (enables full schemas). Call help first to see groups.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"group": map[string]any{
					"type":        "string",
					"description": "group to load (e.g. 'agents', 'threats')",
				},
			},
			"required": []string{"group"},
		},
		Run: func(args map[string]any) (string, error) {
			group, _ := args["group"].(string)
			if group == "" {
				return "", fmt.Errorf("group is required")
			}
			tools, err := GroupTools(s.root, group)
			if err != nil {
				return "", err
			}
			s.focused[group] = tools
			s.rebuildToolList()

			names := make([]string, len(tools))
			for i, t := range tools {
				names[i] = t.Name
			}
			return fmt.Sprintf("Loaded %d tools for %s: %s\nTools are available on the next turn.",
				len(tools), group, strings.Join(names, ", ")), nil
		},
	}
}

func (s *Server) buildUnfocusTool() Tool {
	return Tool{
		Name:        "unfocus",
		Description: "Unload a command group's tools to free context space. Omit group to unload all.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"group": map[string]any{
					"type":        "string",
					"description": "group to unload. Omit to unload all focused groups.",
				},
			},
		},
		Run: func(args map[string]any) (string, error) {
			group, _ := args["group"].(string)
			if group == "" {
				count := len(s.focused)
				s.focused = make(map[string][]Tool)
				s.rebuildToolList()
				return fmt.Sprintf("Unloaded all %d groups.", count), nil
			}
			if _, ok := s.focused[group]; !ok {
				return fmt.Sprintf("Group %q is not focused.", group), nil
			}
			delete(s.focused, group)
			s.rebuildToolList()
			return fmt.Sprintf("Unloaded %s.", group), nil
		},
	}
}

func (s *Server) helpOutput(group string) (string, error) {
	if s.root == nil {
		return "no command tree available", nil
	}

	if group != "" {
		for _, c := range s.root.Commands() {
			if c.Name() != group || c.Hidden {
				continue
			}
			if _, ok := skipCommands[c.Name()]; ok {
				continue
			}
			var b strings.Builder
			fmt.Fprintf(&b, "## %s\n%s\n\n", c.Name(), c.Short)
			if c.RunE != nil || c.Run != nil {
				fmt.Fprintf(&b, "  %s (root command)\n", c.Name())
			}
			for _, sub := range c.Commands() {
				if sub.Hidden {
					continue
				}
				writeHelpLine(&b, sub, "  ")
				for _, subsub := range sub.Commands() {
					if subsub.Hidden {
						continue
					}
					writeHelpLine(&b, subsub, "    ")
				}
			}
			if _, focused := s.focused[group]; focused {
				fmt.Fprintf(&b, "\n[focused — typed tools loaded]")
			}
			return b.String(), nil
		}
		return "", fmt.Errorf("unknown group: %s", group)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "Available command groups (use focus to load typed tools):\n\n")
	for _, c := range s.root.Commands() {
		if c.Hidden {
			continue
		}
		if _, ok := skipCommands[c.Name()]; ok {
			continue
		}
		n := 0
		countLeaves(c, &n)
		status := ""
		if _, focused := s.focused[c.Name()]; focused {
			status = " [focused]"
		}
		fmt.Fprintf(&b, "  %-24s %s (%d commands)%s\n", c.Name(), c.Short, n, status)
	}
	return b.String(), nil
}

func writeHelpLine(b *strings.Builder, cmd *cobra.Command, indent string) {
	line := indent + cmd.Name() + "  " + cmd.Short
	if hasMutationFlag(cmd) {
		line += " [mutation]"
	}
	if hints := flagHints(cmd); hints != "" {
		line += "  " + hints
	}
	fmt.Fprintln(b, line)
}

func flagHints(cmd *cobra.Command) string {
	seen := map[string]bool{}
	var flags []string
	add := func(f *pflag.Flag) {
		if skipFlags[f.Name] || f.Hidden || f.Name == "yes" || seen[f.Name] {
			return
		}
		seen[f.Name] = true
		flags = append(flags, "--"+f.Name)
	}
	cmd.Flags().VisitAll(add)
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if skipFlags[f.Name] || f.Hidden || f.Name == "yes" || f.Name == "read-only" || f.Name == "site-id" {
			return
		}
		add(f)
	})
	if len(flags) == 0 {
		return ""
	}
	return "(" + strings.Join(flags, ", ") + ")"
}

func countLeaves(cmd *cobra.Command, n *int) {
	if !cmd.HasSubCommands() {
		if cmd.RunE != nil || cmd.Run != nil {
			*n++
		}
		return
	}
	for _, c := range cmd.Commands() {
		if !c.Hidden {
			countLeaves(c, n)
		}
	}
}

func (s *Server) execCommand(parts []string) (string, error) {
	cliArgs := make([]string, 0, len(parts)+2)
	cliArgs = append(cliArgs, parts...)
	cliArgs = append(cliArgs, "--json", "--no-progress")

	origStdout := os.Stdout
	pr, pw, err := os.Pipe()
	if err != nil {
		return "", fmt.Errorf("create pipe: %w", err)
	}
	os.Stdout = pw

	var stderr bytes.Buffer
	s.root.SetArgs(cliArgs)
	s.root.SetOut(pw)
	s.root.SetErr(&stderr)

	execErr := s.root.Execute()

	_ = pw.Close()
	os.Stdout = origStdout

	var stdout bytes.Buffer
	_, _ = stdout.ReadFrom(pr)
	_ = pr.Close()

	out := stdout.String()
	if errOut := stderr.String(); errOut != "" && out == "" {
		out = errOut
	}
	if out == "" && execErr != nil {
		out = execErr.Error()
	}
	return out, execErr
}
