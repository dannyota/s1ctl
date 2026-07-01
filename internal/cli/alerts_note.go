package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAlertsAddNoteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "add-note <alert-id> <text>",
		Short: "Add an investigation note to an alert",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, text := args[0], args[1]
			return guard(cmd.OutOrStdout(), "alerts add-note", "add note to alert "+id, id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.AlertsAddNote(cmd.Context(), []string{id}, text); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "noted", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "note: added to alert %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
