package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAlertsNotesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "notes <alert-id>",
		Short: "List investigation notes on an alert",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			notes, err := c.AlertNotes(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			headers := []string{"ID", "Author", "Type", "Updated", "Text"}
			rows := make([][]string, len(notes))
			for i, n := range notes {
				rows[i] = []string{
					n.ID, orDash(n.AuthorName()), orDash(n.Type),
					orDash(n.UpdatedAt), truncate(orDash(n.Text), 50),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, notes, len(notes), len(notes), "note", true)
		},
	}
}

func newAlertsNoteUpdateCmd() *cobra.Command {
	var yes bool
	var text string

	cmd := &cobra.Command{
		Use:   "note-update <note-id> --text <text>",
		Short: "Update the text of an alert note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if text == "" {
				return fmt.Errorf("--text is required")
			}
			id := args[0]
			return guard(cmd.OutOrStdout(), "alerts note-update", "update note "+id, id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				notes, err := c.AlertsUpdateNote(cmd.Context(), id, text)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), notes)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "note: updated %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&text, "text", "", "new note text (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newAlertsNoteDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "note-delete <note-id>",
		Short: "Delete an alert note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return guard(cmd.OutOrStdout(), "alerts note-delete", "delete note "+id, id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				notes, err := c.AlertsDeleteNote(cmd.Context(), id)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), notes)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "note: deleted %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

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
