package mcp

import (
	"encoding/json"
	"fmt"
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
	Annotations *toolAnnotations
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
	name := strings.ReplaceAll(strings.Join(path, "_"), "-", "_")

	mutation := hasMutationFlag(cmd)
	desc := cmd.Short
	if mutation {
		desc += " [mutation: requires --yes to apply, dry-run by default]"
	}
	if hasJSONAnnotation(cmd) {
		desc += " [supports --json output]"
	}

	schema := buildInputSchema(cmd)

	roHint := !mutation
	destHint := mutation
	return Tool{
		Name:        name,
		Description: desc,
		InputSchema: schema,
		Annotations: &toolAnnotations{ReadOnlyHint: &roHint, DestructiveHint: &destHint},
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

		return execSubprocess(cliArgs, nil)
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
		s.buildUsageTool(),
		s.buildFocusTool(),
		s.buildUnfocusTool(),
	}
}

func (s *Server) buildRunTool() Tool {
	desc := "Run any s1ctl command. Pass the full command (e.g. 'agents list --site-id 123'). For filter expressions use shell quoting: --filter 'event.type = \"Login\"'. Prefer focus + typed tools for complex filters."
	var ann *toolAnnotations
	if s.readOnly {
		desc += " [read-only mode: mutations are blocked]"
		ro := true
		ann = &toolAnnotations{ReadOnlyHint: &ro}
	}
	return Tool{
		Name:        "run",
		Description: desc,
		Annotations: ann,
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
			return s.execCommand(splitCommand(cmdStr))
		},
	}
}

func readOnlyAnnotation() *toolAnnotations {
	ro, dest := true, false
	return &toolAnnotations{ReadOnlyHint: &ro, DestructiveHint: &dest}
}

func (s *Server) buildHelpTool() Tool {
	return Tool{
		Name:        "help",
		Description: "List command groups, or subcommands within a group.",
		Annotations: readOnlyAnnotation(),
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
		Annotations: readOnlyAnnotation(),
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
			s.mu.Lock()
			s.focused[group] = tools
			s.rebuildToolList()
			s.mu.Unlock()

			// Report only the tools that survive read-only filtering so
			// the agent never sees names absent from tools/list.
			visible := tools
			if s.readOnly {
				visible = filterReadOnly(tools)
			}
			names := make([]string, len(visible))
			for i, t := range visible {
				names[i] = t.Name
			}
			return fmt.Sprintf("Loaded %d tools for %s: %s\nTools are available on the next turn.",
				len(visible), group, strings.Join(names, ", ")), nil
		},
	}
}

func (s *Server) buildUnfocusTool() Tool {
	return Tool{
		Name:        "unfocus",
		Description: "Unload a command group's tools to free context space. Omit group to unload all.",
		Annotations: readOnlyAnnotation(),
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
				s.mu.Lock()
				count := len(s.focused)
				s.focused = make(map[string][]Tool)
				s.rebuildToolList()
				s.mu.Unlock()
				return fmt.Sprintf("Unloaded all %d groups.", count), nil
			}
			s.mu.Lock()
			if _, ok := s.focused[group]; !ok {
				s.mu.Unlock()
				return fmt.Sprintf("Group %q is not focused.", group), nil
			}
			delete(s.focused, group)
			s.rebuildToolList()
			s.mu.Unlock()
			return fmt.Sprintf("Unloaded %s.", group), nil
		},
	}
}

func (s *Server) buildUsageTool() Tool {
	return Tool{
		Name:        "usage",
		Description: "Show flags, args, and description for one command. Use before run to learn a command's interface.",
		Annotations: readOnlyAnnotation(),
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "command path (e.g. 'agents list', 'threats mitigate')",
				},
			},
			"required": []string{"command"},
		},
		Run: func(args map[string]any) (string, error) {
			cmd, _ := args["command"].(string)
			if cmd == "" {
				return "", fmt.Errorf("command is required")
			}
			return s.usageOutput(cmd)
		},
	}
}

func (s *Server) usageOutput(cmd string) (string, error) {
	toolName := strings.ReplaceAll(strings.ReplaceAll(cmd, " ", "_"), "-", "_")

	tool, ok := s.allToolIndex[toolName]
	if !ok {
		return "", fmt.Errorf("unknown command: %s. Use help to see available commands", cmd)
	}

	b, _ := json.MarshalIndent(map[string]any{
		"name":        tool.Name,
		"description": tool.Description,
		"inputSchema": tool.InputSchema,
	}, "", "  ")
	return string(b), nil
}

func (s *Server) helpOutput(group string) (string, error) {
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
	s.mu.Lock()
	focused := make(map[string]bool, len(s.focused))
	for g := range s.focused {
		focused[g] = true
	}
	s.mu.Unlock()

	var b strings.Builder
	b.WriteString("Command groups (help {group} for subcommands):\n\n")
	for _, c := range s.root.Commands() {
		if c.Hidden {
			continue
		}
		if _, ok := skipCommands[c.Name()]; ok {
			continue
		}
		reads, muts := countLeaves(c)
		if s.readOnly {
			muts = 0
		}
		status := ""
		if focused[c.Name()] {
			status = " *"
		}
		fmt.Fprintf(&b, "  %-22s %s (%dr/%dm)%s\n", c.Name(), c.Short, reads, muts, status)
	}
	return b.String()
}

func (s *Server) helpGroup(cmd *cobra.Command, group string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s — %s\n\n", group, cmd.Short)

	writeCmd := func(c *cobra.Command, indent string) {
		if s.readOnly && hasMutationFlag(c) {
			return
		}
		line := indent + c.Name()
		if hasMutationFlag(c) {
			line += " [mutation]"
		}
		if hasJSONAnnotation(c) {
			line += " [json]"
		}
		line += "  " + c.Short
		fmt.Fprintln(&b, line)
	}

	if cmd.RunE != nil || cmd.Run != nil {
		writeCmd(cmd, "  ")
	}
	for _, sub := range cmd.Commands() {
		if sub.Hidden {
			continue
		}
		writeCmd(sub, "  ")
		for _, subsub := range sub.Commands() {
			if !subsub.Hidden {
				writeCmd(subsub, "    ")
			}
		}
	}

	s.mu.Lock()
	_, isFocused := s.focused[group]
	s.mu.Unlock()
	if isFocused {
		b.WriteString("\n[focused]")
	} else {
		b.WriteString("\nUse usage {command} for flags, or focus to load typed tools.")
	}
	return b.String()
}

func countLeaves(cmd *cobra.Command) (int, int) {
	if !cmd.HasSubCommands() {
		if (cmd.RunE != nil || cmd.Run != nil) && hasMutationFlag(cmd) {
			return 0, 1
		}
		if cmd.RunE != nil || cmd.Run != nil {
			return 1, 0
		}
		return 0, 0
	}
	reads, muts := 0, 0
	for _, c := range cmd.Commands() {
		if !c.Hidden {
			r, m := countLeaves(c)
			reads += r
			muts += m
		}
	}
	return reads, muts
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

// splitCommand tokenises a command string with shell-style quoting.
// Single-quoted, double-quoted, and backslash-escaped characters are
// kept as single tokens so that filter expressions like
//
//	--filter 'event.type = "Login"'
//
// survive as one argument.
func splitCommand(s string) []string {
	var args []string
	var cur []byte
	inSingle, inDouble := false, false

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\\' && !inSingle && i+1 < len(s):
			i++
			cur = append(cur, s[i])
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == '"' && !inSingle:
			inDouble = !inDouble
		case (c == ' ' || c == '\t') && !inSingle && !inDouble:
			if len(cur) > 0 {
				args = append(args, string(cur))
				cur = cur[:0]
			}
		default:
			cur = append(cur, c)
		}
	}
	if len(cur) > 0 {
		args = append(args, string(cur))
	}
	return args
}
