package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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
}

type groupEntry struct {
	Name  string `json:"name"`
	Short string `json:"short"`
	Reads int    `json:"reads"`
	Muts  int    `json:"mutations"`
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
		countKinds(c, &reads, &muts)
		groups = append(groups, groupEntry{
			Name: c.Name(), Short: c.Short,
			Reads: reads, Muts: muts,
		})
	}

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), groups)
	}

	var rows [][]string
	for _, g := range groups {
		count := fmt.Sprintf("%dr/%dm", g.Reads, g.Muts)
		rows = append(rows, []string{g.Name, g.Short, count})
	}
	printTable([]string{"Group", "Description", "Commands"}, rows)
	total := 0
	for _, g := range groups {
		total += g.Reads + g.Muts
	}
	fmt.Fprintf(cmd.OutOrStdout(), "\n%d groups, %d commands. Use: commands <group> for details.\n", len(groups), total)
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

	type flagInfo struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Usage   string `json:"usage"`
		Default string `json:"default,omitempty"`
	}

	var flags []flagInfo
	sub.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" {
			return
		}
		fi := flagInfo{Name: f.Name, Type: f.Value.Type(), Usage: f.Usage}
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
			fi.Default = f.DefValue
		}
		flags = append(flags, fi)
	})

	if outputFormat == "json" {
		detail := map[string]any{
			"name":  group + " " + name,
			"short": sub.Short,
			"kind":  commandKind(sub),
			"flags": flags,
		}
		return printJSON(cmd.OutOrStdout(), detail)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s %s — %s", group, name, sub.Short)
	if commandKind(sub) == "mutation" {
		fmt.Fprint(cmd.OutOrStdout(), " [mutation]")
	}
	fmt.Fprintln(cmd.OutOrStdout())

	if len(flags) > 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "\nFlags:")
		var rows [][]string
		for _, f := range flags {
			def := ""
			if f.Default != "" {
				def = f.Default
			}
			rows = append(rows, []string{"--" + f.Name, f.Type, f.Usage, def})
		}
		printTable([]string{"Flag", "Type", "Description", "Default"}, rows)
	}
	return nil
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

func countKinds(cmd *cobra.Command, reads, muts *int) {
	if !cmd.HasSubCommands() {
		if cmd.RunE != nil || cmd.Run != nil {
			if cmd.Flags().Lookup("yes") != nil {
				*muts++
			} else {
				*reads++
			}
		}
		return
	}
	for _, c := range cmd.Commands() {
		if !c.Hidden {
			countKinds(c, reads, muts)
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
			})
		}
	}
	return entries
}
