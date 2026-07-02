package cli

import (
	"github.com/spf13/cobra"
)

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
	Kind  string `json:"kind"`
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

// commandKind classifies a command as a mutation iff it declares a --yes
// flag; every mutation registers one via the guard pattern.
func commandKind(cmd *cobra.Command) string {
	if cmd.Flags().Lookup("yes") != nil {
		return "mutation"
	}
	return "read"
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
			entries = append(entries, cmdEntry{
				Name:  name,
				Short: c.Short,
				Kind:  commandKind(c),
			})
		}
	}
	return entries
}
