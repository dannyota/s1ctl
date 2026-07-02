package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newMisconfigurationsNotesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "notes <id>",
		Short: "List investigation notes on a misconfiguration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			notes, err := c.MisconfigurationsNotes(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			headers := []string{"ID", "Author", "Created", "Updated", "Text"}
			rows := make([][]string, len(notes))
			for i, n := range notes {
				rows[i] = []string{
					n.ID, orDash(n.AuthorName()), orDash(n.CreatedAt),
					orDash(n.UpdatedAt), truncate(orDash(n.Text), 50),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, notes, len(notes), len(notes), "note", true)
		},
	}
}

func newMisconfigurationsNoteAddCmd() *cobra.Command {
	var yes bool
	var text string

	cmd := &cobra.Command{
		Use:   "note-add <id> --text <text>",
		Short: "Add an investigation note to a misconfiguration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if text == "" {
				return fmt.Errorf("--text is required")
			}
			id := args[0]
			return guard(cmd.OutOrStdout(), "misconfigurations note-add", "add note to misconfiguration "+id, id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.MisconfigurationsAddNote(cmd.Context(), []string{id}, text); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "noted", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "note: added to misconfiguration %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&text, "text", "", "note text (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newMisconfigurationsNoteUpdateCmd() *cobra.Command {
	var yes bool
	var text string

	cmd := &cobra.Command{
		Use:   "note-update <note-id> --text <text>",
		Short: "Update the text of a misconfiguration note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if text == "" {
				return fmt.Errorf("--text is required")
			}
			id := args[0]
			return guard(cmd.OutOrStdout(), "misconfigurations note-update", "update note "+id, id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.MisconfigurationsUpdateNote(cmd.Context(), id, text); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "id": id})
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

func newMisconfigurationsNoteDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "note-delete <note-id>",
		Short: "Delete a misconfiguration note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return guard(cmd.OutOrStdout(), "misconfigurations note-delete", "delete note "+id, id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.MisconfigurationsDeleteNote(cmd.Context(), id); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "note: deleted %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newMisconfigurationsAssignCmd() *cobra.Command {
	var yes bool
	var userID string

	cmd := &cobra.Command{
		Use:   "assign <id> --user-id <user-id>",
		Short: "Assign a misconfiguration to a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if userID == "" {
				return fmt.Errorf("--user-id is required")
			}
			id := args[0]
			return guard(cmd.OutOrStdout(), "misconfigurations assign", fmt.Sprintf("assign misconfiguration %s to user %s", id, userID), id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.MisconfigurationsAssign(cmd.Context(), []string{id}, userID); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "assigned", "id": id, "userId": userID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "assign: misconfiguration %s -> user %s\n", id, userID)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&userID, "user-id", "", "assignee user ID (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newMisconfigurationsHistoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "history <id>",
		Short: "Show the history of a misconfiguration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			items, err := c.MisconfigurationsHistory(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			headers := []string{"Created", "Type", "Event"}
			rows := make([][]string, len(items))
			for i, h := range items {
				rows[i] = []string{orDash(h.CreatedAt), orDash(h.EventType), truncate(orDash(h.EventText), 60)}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), len(items), "history item", true)
		},
	}
}

func newMisconfigurationsRelatedAssetsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "related-assets <id>",
		Short: "List assets related to a misconfiguration",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			assets, err := c.MisconfigurationsRelatedAssets(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			headers := []string{"Asset ID", "Name", "Type", "OS", "Organization"}
			rows := make([][]string, len(assets))
			for i, a := range assets {
				rows[i] = []string{
					a.Asset.ID, orDash(a.Asset.Name), orDash(a.Asset.Type),
					orDash(a.Asset.OsType), orDash(a.Organization),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, assets, len(assets), len(assets), "related asset", true)
		},
	}
}

func newMisconfigurationsExportCmd() *cobra.Command {
	var severities, statuses []string
	var scopeLevel, scopeID, outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export misconfigurations to a CSV file",
		Long: `Export misconfigurations matching the filters as CSV via
misconfigurationsExportToCsv. The API returns the full CSV inline; it is written
to --out, or to stdout when --out is omitted.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			data, err := c.MisconfigurationsExport(cmd.Context(), alertsFilters(severities, statuses, nil), scope)
			if err != nil {
				return err
			}
			return writeExport(cmd, outFile, data, "misconfigurations")
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (account, site, group)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "account, site, or group ID")
	cmd.Flags().StringVar(&outFile, "out", "", "output file (default: stdout)")
	return cmd
}

// writeExport writes CSV export data to a file or stdout.
func writeExport(cmd *cobra.Command, outFile, data, resource string) error {
	if outFile == "" || outFile == "-" {
		_, err := cmd.OutOrStdout().Write([]byte(data))
		return err
	}
	if err := os.WriteFile(outFile, []byte(data), 0o644); err != nil { //nolint:gosec
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Exported %s to %s\n", resource, outFile)
	return nil
}
