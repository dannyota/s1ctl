package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const jsonAnnotation = "s1ctl_json"

func markJSON(cmd *cobra.Command) *cobra.Command {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[jsonAnnotation] = "true"
	return cmd
}

func newCommandsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commands [group] [command]",
		Short: "List command groups, subcommands, or flag detail",
		Long: `List available commands in a compact, progressive format.

  commands              List all command groups with counts
  commands agents       List subcommands in the agents group
  commands agents list  Show full flag detail for agents list
  commands --all        Flat list of every command (JSON only)`,
		Args: cobra.MaximumNArgs(2),
		RunE: runCommands,
	}
	cmd.Flags().Bool("all", false, "list every command (flat, JSON-only)")
	return cmd
}

type cmdEntry struct {
	Name  string `json:"name"`
	Short string `json:"short"`
	Kind  string `json:"kind"`
	JSON  bool   `json:"json"`
}

type groupEntry struct {
	Name  string   `json:"name"`
	Short string   `json:"short"`
	Reads int      `json:"reads"`
	Muts  int      `json:"mutations"`
	Cmds  []string `json:"commands"`
}

func runCommands(cmd *cobra.Command, args []string) error {
	allFlag, _ := cmd.Flags().GetBool("all")
	if allFlag {
		entries := collectCommands(cmd.Root(), "")
		return printJSON(cmd.OutOrStdout(), entries)
	}

	switch len(args) {
	case 0:
		return commandsGroups(cmd)
	case 1:
		return commandsGroup(cmd, args[0])
	default:
		return commandsDetail(cmd, args[0], args[1])
	}
}

func commandsGroups(cmd *cobra.Command) error {
	root := cmd.Root()
	var groups []groupEntry
	for _, c := range root.Commands() {
		if c.Hidden || isSkippedGroup(c.Name()) {
			continue
		}
		reads, muts := 0, 0
		var names []string
		collectLeafNames(c, &reads, &muts, &names)
		groups = append(groups, groupEntry{
			Name: c.Name(), Short: c.Short,
			Reads: reads, Muts: muts, Cmds: names,
		})
	}

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), groups)
	}

	var rows [][]string
	for _, g := range groups {
		count := fmt.Sprintf("%dr/%dm", g.Reads, g.Muts)
		cmds := strings.Join(g.Cmds, ", ")
		rows = append(rows, []string{g.Name, g.Short, cmds, count})
	}
	printTable([]string{"Group", "Description", "Commands", "Count"}, rows)
	total := 0
	for _, g := range groups {
		total += g.Reads + g.Muts
	}
	fmt.Fprintf(cmd.OutOrStdout(), "\n%d groups, %d commands (* = mutation). Use: commands <group> for details.\n", len(groups), total)
	return nil
}

func commandsGroup(cmd *cobra.Command, group string) error {
	root := cmd.Root()
	groupCmd := findGroupCmd(root, group)
	if groupCmd == nil {
		return fmt.Errorf("unknown group: %s", group)
	}

	var entries []cmdEntry
	collectGroupEntries(groupCmd, group, &entries)

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), entries)
	}

	var rows [][]string
	for _, e := range entries {
		rows = append(rows, []string{e.Name, e.Short, e.Kind})
	}
	printTable([]string{"Command", "Description", "Kind"}, rows)
	return nil
}

var boilerplateFlags = map[string]bool{
	"yes": true, "json": true,
}

type flagInfo struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Default  string   `json:"default,omitempty"`
	Required bool     `json:"required,omitempty"`
	Enum     []string `json:"enum,omitempty"`
	Usage    string   `json:"usage"`
}

func commandsDetail(cmd *cobra.Command, group, name string) error {
	root := cmd.Root()
	groupCmd := findGroupCmd(root, group)
	if groupCmd == nil {
		return fmt.Errorf("unknown group: %s", group)
	}

	sub := findSub(groupCmd, name)
	if sub == nil {
		return fmt.Errorf("unknown command: %s %s", group, name)
	}

	flags := localFlagInfos(sub)

	if outputFormat == "json" {
		detail := map[string]any{
			"name":  group + " " + name,
			"short": sub.Short,
			"kind":  commandKind(sub),
			"json":  sub.Annotations[jsonAnnotation] == "true",
		}
		if spec := positionalSpec(sub); spec != "" {
			detail["args"] = spec
		}
		if len(flags) > 0 {
			detail["flags"] = flags
		}
		return printJSON(cmd.OutOrStdout(), detail)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s %s — %s", group, name, sub.Short)
	if commandKind(sub) == "mutation" {
		fmt.Fprint(cmd.OutOrStdout(), " [mutation]")
	}
	fmt.Fprintln(cmd.OutOrStdout())
	if spec := positionalSpec(sub); spec != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "  args: %s\n", spec)
	}

	if len(flags) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "\nFlags:")
		var rows [][]string
		for _, f := range flags {
			extra := ""
			if f.Required {
				extra = " [required]"
			}
			if len(f.Enum) > 0 {
				extra += " {" + strings.Join(f.Enum, "|") + "}"
			}
			rows = append(rows, []string{"--" + f.Name, f.Type, f.Usage + extra, f.Default})
		}
		printTable([]string{"Flag", "Type", "Description", "Default"}, rows)
		fmt.Fprintln(cmd.OutOrStdout(), "\nStandard flags (--yes, --json) omitted.")
	}
	return nil
}

func localFlagInfos(c *cobra.Command) []flagInfo {
	var infos []flagInfo
	c.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" || boilerplateFlags[f.Name] {
			return
		}
		_, required := f.Annotations[cobra.BashCompOneRequiredFlag]
		fi := flagInfo{
			Name:     f.Name,
			Type:     f.Value.Type(),
			Required: required,
			Enum:     enumFromUsage(f.Usage),
			Usage:    f.Usage,
		}
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
			fi.Default = f.DefValue
		}
		infos = append(infos, fi)
	})
	return infos
}

func findGroupCmd(root *cobra.Command, name string) *cobra.Command {
	for _, c := range root.Commands() {
		if c.Name() == name && !c.Hidden {
			return c
		}
	}
	return nil
}

func findSub(parent *cobra.Command, name string) *cobra.Command {
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

func collectGroupEntries(cmd *cobra.Command, prefix string, entries *[]cmdEntry) {
	if cmd.RunE != nil || cmd.Run != nil {
		if cmd.HasSubCommands() {
			*entries = append(*entries, cmdEntry{
				Name: prefix, Short: cmd.Short, Kind: commandKind(cmd),
				JSON: cmd.Annotations[jsonAnnotation] == "true",
			})
		}
	}
	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}
		name := prefix + " " + c.Name()
		if c.HasSubCommands() {
			collectGroupEntries(c, name, entries)
		} else if c.RunE != nil || c.Run != nil {
			*entries = append(*entries, cmdEntry{
				Name: name, Short: c.Short, Kind: commandKind(c),
				JSON: c.Annotations[jsonAnnotation] == "true",
			})
		}
	}
}

func isSkippedGroup(name string) bool {
	switch name {
	case "help", "completion":
		return true
	}
	return false
}

func commandKind(cmd *cobra.Command) string {
	if cmd.Flags().Lookup("yes") != nil {
		return "mutation"
	}
	return "read"
}

func collectLeafNames(cmd *cobra.Command, reads, muts *int, names *[]string) {
	if !cmd.HasSubCommands() {
		if cmd.RunE != nil || cmd.Run != nil {
			n := cmd.Name()
			if cmd.Flags().Lookup("yes") != nil {
				*muts++
				n += "*"
			} else {
				*reads++
			}
			*names = append(*names, n)
		}
		return
	}
	for _, c := range cmd.Commands() {
		if !c.Hidden {
			collectLeafNames(c, reads, muts, names)
		}
	}
}

func collectCommands(cmd *cobra.Command, prefix string) []cmdEntry {
	var entries []cmdEntry
	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}
		name := prefix + c.Name()
		if c.HasSubCommands() {
			entries = append(entries, collectCommands(c, name+" ")...)
		} else {
			kind := "read"
			if c.Flags().Lookup("yes") != nil {
				kind = "mutation"
			}
			entries = append(entries, cmdEntry{
				Name:  name,
				Short: c.Short,
				Kind:  kind,
				JSON:  c.Annotations[jsonAnnotation] == "true",
			})
		}
	}
	return entries
}

func positionalSpec(c *cobra.Command) string {
	use := strings.TrimSpace(c.Use)
	i := strings.IndexAny(use, " \t")
	if i < 0 {
		return ""
	}
	return strings.TrimSpace(use[i+1:])
}

var (
	enumPattern   = regexp.MustCompile(`[A-Za-z][\w-]+(?:\s*\|\s*[A-Za-z][\w-]+)+`)
	placeholderRE = regexp.MustCompile(`<[^>]*>`)
)

func enumFromUsage(usage string) []string {
	run := enumPattern.FindString(placeholderRE.ReplaceAllString(usage, ""))
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
