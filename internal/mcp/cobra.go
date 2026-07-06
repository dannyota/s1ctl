package mcp

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
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
	if hasJSONAnnotation(cmd) {
		desc += " [supports --json output]"
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

	if vals := mcpEnumFromUsage(f.Usage); len(vals) > 0 {
		enumAny := make([]any, len(vals))
		for i, v := range vals {
			enumAny[i] = v
		}
		prop["enum"] = enumAny
	}

	return prop
}

func hasMutationFlag(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup("yes") != nil
}

func hasJSONAnnotation(cmd *cobra.Command) bool {
	return cmd.Annotations != nil && cmd.Annotations["s1ctl_json"] == "true"
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
		Description: "List command groups, subcommands within a group, or full flag detail for one command.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"group": map[string]any{
					"type":        "string",
					"description": "group name (e.g. 'agents'). Omit to list all groups.",
				},
				"command": map[string]any{
					"type":        "string",
					"description": "subcommand name within the group (e.g. 'isolate'). Returns full flag detail for that command.",
				},
			},
		},
		Run: func(args map[string]any) (string, error) {
			group, _ := args["group"].(string)
			command, _ := args["command"].(string)
			return s.helpOutput(group, command)
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

func (s *Server) helpOutput(group, command string) (string, error) {
	if s.root == nil {
		return "no command tree available", nil
	}

	if group == "" {
		return s.helpGroups(), nil
	}

	groupCmd := findGroup(s.root, group)
	if groupCmd == nil {
		return "", fmt.Errorf("unknown group: %s", group)
	}

	if command != "" {
		return s.helpCommand(groupCmd, group, command)
	}
	return s.helpGroup(groupCmd, group), nil
}

func findGroup(root *cobra.Command, name string) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Name() == name && !c.Hidden {
			if _, ok := skipCommands[c.Name()]; !ok {
				return c
			}
		}
	}
	return nil
}

func (s *Server) helpGroups() string {
	var b strings.Builder
	b.WriteString("Command groups (help {group} for subcommands, help {group} {command} for flags):\n\n")
	for _, c := range s.root.Commands() {
		if c.Hidden {
			continue
		}
		if _, ok := skipCommands[c.Name()]; ok {
			continue
		}
		reads, muts := 0, 0
		var names []string
		collectLeafNamesForMCP(c, &reads, &muts, &names)
		status := ""
		if _, focused := s.focused[c.Name()]; focused {
			status = " *"
		}
		fmt.Fprintf(&b, "  %-22s %s (%dr/%dm)%s\n", c.Name(), c.Short, reads, muts, status)
		fmt.Fprintf(&b, "    %s\n", strings.Join(names, ", "))
	}
	return b.String()
}

func (s *Server) helpGroup(cmd *cobra.Command, group string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s — %s\n\n", group, cmd.Short)

	if cmd.RunE != nil || cmd.Run != nil {
		writeHelpLine(&b, cmd, "  ")
	}
	for _, sub := range cmd.Commands() {
		if sub.Hidden {
			continue
		}
		writeHelpLine(&b, sub, "  ")
		for _, subsub := range sub.Commands() {
			if !subsub.Hidden {
				writeHelpLine(&b, subsub, "    ")
			}
		}
	}

	if _, focused := s.focused[group]; focused {
		b.WriteString("\n[focused]")
	} else {
		fmt.Fprintf(&b, "\nUse help {group} {command} for flag detail, or focus to load typed tools.")
	}
	return b.String()
}

func (s *Server) helpCommand(groupCmd *cobra.Command, group, command string) (string, error) {
	cmd := findSubcommand(groupCmd, command)
	if cmd == nil {
		return "", fmt.Errorf("unknown command: %s %s", group, command)
	}

	var b strings.Builder
	fullName := group + " " + command
	fmt.Fprintf(&b, "## %s\n%s", fullName, cmd.Short)
	if hasMutationFlag(cmd) {
		b.WriteString(" [mutation: dry-run by default, --yes to apply]")
	}
	if hasJSONAnnotation(cmd) {
		b.WriteString(" [json]")
	}
	b.WriteString("\n")

	if spec := positionalSpec(cmd); spec != "" {
		fmt.Fprintf(&b, "  args: %s\n", spec)
	}

	var flags []flagDetail
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if skipFlags[f.Name] || f.Hidden || f.Name == "yes" {
			return
		}
		flags = append(flags, describeFlag(f))
	})
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if skipFlags[f.Name] || f.Hidden {
			return
		}
		if f.Name == "site-id" || f.Name == "read-only" {
			flags = append(flags, describeFlag(f))
		}
	})

	if len(flags) > 0 {
		b.WriteString("\nFlags:\n")
		for _, fd := range flags {
			fmt.Fprintf(&b, "  --%-20s %-8s %s", fd.Name, fd.Type, fd.Usage)
			if fd.Required {
				b.WriteString(" [required]")
			}
			if len(fd.Enum) > 0 {
				b.WriteString(" {" + strings.Join(fd.Enum, "|") + "}")
			}
			if fd.Default != "" {
				fmt.Fprintf(&b, " (default: %s)", fd.Default)
			}
			b.WriteString("\n")
		}
	}
	return b.String(), nil
}

type flagDetail struct {
	Name     string
	Type     string
	Usage    string
	Default  string
	Required bool
	Enum     []string
}

func describeFlag(f *pflag.Flag) flagDetail {
	_, required := f.Annotations[cobra.BashCompOneRequiredFlag]
	fd := flagDetail{
		Name:     f.Name,
		Type:     f.Value.Type(),
		Usage:    f.Usage,
		Required: required,
		Enum:     mcpEnumFromUsage(f.Usage),
	}
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
		fd.Default = f.DefValue
	}
	return fd
}

func findSubcommand(parent *cobra.Command, name string) *cobra.Command {
	for _, c := range parent.Commands() {
		if c.Hidden {
			continue
		}
		if c.Name() == name {
			return c
		}
		for _, sub := range c.Commands() {
			if !sub.Hidden && sub.Name() == name {
				return sub
			}
		}
	}
	return nil
}

func collectLeafNamesForMCP(cmd *cobra.Command, reads, muts *int, names *[]string) {
	if !cmd.HasSubCommands() {
		if cmd.RunE != nil || cmd.Run != nil {
			n := cmd.Name()
			if hasMutationFlag(cmd) {
				*muts++
				n += "*"
			} else {
				*reads++
			}
			if hasJSONAnnotation(cmd) {
				n += "[j]"
			}
			*names = append(*names, n)
		}
		return
	}
	for _, c := range cmd.Commands() {
		if !c.Hidden {
			collectLeafNamesForMCP(c, reads, muts, names)
		}
	}
}

func writeHelpLine(b *strings.Builder, cmd *cobra.Command, indent string) {
	line := indent + cmd.Name()
	if spec := positionalSpec(cmd); spec != "" {
		line += " " + spec
	}
	line += "  " + cmd.Short
	if hasMutationFlag(cmd) {
		line += " [mutation]"
	}
	if hasJSONAnnotation(cmd) {
		line += " [json]"
	}
	if hints := flagHints(cmd); hints != "" {
		line += "  " + hints
	}
	fmt.Fprintln(b, line)
}

func positionalSpec(c *cobra.Command) string {
	use := strings.TrimSpace(c.Use)
	i := strings.IndexAny(use, " \t")
	if i < 0 {
		return ""
	}
	return strings.TrimSpace(use[i+1:])
}

func flagHints(cmd *cobra.Command) string {
	seen := map[string]bool{}
	var flags []string
	add := func(f *pflag.Flag) {
		if skipFlags[f.Name] || f.Hidden || f.Name == "yes" || seen[f.Name] {
			return
		}
		seen[f.Name] = true
		hint := "--" + f.Name
		if _, req := f.Annotations[cobra.BashCompOneRequiredFlag]; req {
			hint += "*"
		}
		if vals := mcpEnumFromUsage(f.Usage); len(vals) > 0 {
			hint += "={" + strings.Join(vals, "|") + "}"
		}
		flags = append(flags, hint)
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

var (
	mcpEnumPattern   = regexp.MustCompile(`[A-Za-z][\w-]+(?:\s*\|\s*[A-Za-z][\w-]+)+`)
	mcpPlaceholderRE = regexp.MustCompile(`<[^>]*>`)
)

func mcpEnumFromUsage(usage string) []string {
	run := mcpEnumPattern.FindString(mcpPlaceholderRE.ReplaceAllString(usage, ""))
	if run == "" {
		return nil
	}
	parts := strings.Split(run, "|")
	vals := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			vals = append(vals, p)
		}
	}
	if len(vals) < 2 {
		return nil
	}
	return vals
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
