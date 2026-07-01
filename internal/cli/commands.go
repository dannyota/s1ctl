package cli

import "github.com/spf13/cobra"

func newCommandsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "commands",
		Short: "List all available commands",
		Args:  cobra.NoArgs,
		RunE:  runCommands,
	}
}

type cmdEntry struct {
	Name  string `json:"name"`
	Short string `json:"short"`
}

func runCommands(cmd *cobra.Command, _ []string) error {
	entries := collectCommands(cmd.Root(), "")

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), entries)
	}

	var rows [][]string
	for _, e := range entries {
		rows = append(rows, []string{e.Name, e.Short})
	}
	printTable([]string{"Command", "Description"}, rows)
	return nil
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
			entries = append(entries, cmdEntry{Name: name, Short: c.Short})
		}
	}
	return entries
}
