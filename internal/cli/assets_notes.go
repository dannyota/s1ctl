package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAssetsNotesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notes",
		Short: "Manage asset notes",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAssetsNotesAddCmd())
	cmd.AddCommand(newAssetsNotesDeleteCmd())
	return cmd
}

func newAssetsNotesAddCmd() *cobra.Command {
	var (
		assetID string
		note    string
		yes     bool
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a note to an asset",
		Long: `Add or update a note on an asset.

Dry-run by default; pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return guard(cmd.OutOrStdout(), "assets notes add", "add note to asset "+assetID, assetID, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				input := &mgmt.XDRAssetNoteInput{
					ResourceID: assetID,
					Note:       note,
				}
				if err := c.XDRAssetNoteCreate(cmd.Context(), input); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "ok", "assetId": assetID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Note added to asset %s\n", assetID)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&assetID, "asset-id", "", "asset ID (required)")
	cmd.Flags().StringVar(&note, "note", "", "note text (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	_ = cmd.MarkFlagRequired("asset-id")
	_ = cmd.MarkFlagRequired("note")
	return markJSON(cmd)
}

func newAssetsNotesDeleteCmd() *cobra.Command {
	var (
		noteID  string
		assetID string
		yes     bool
	)

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a note from an asset",
		Long: `Delete a note from an asset.

Dry-run by default; pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return guard(cmd.OutOrStdout(), "assets notes delete", "delete note "+noteID+" from asset "+assetID, noteID, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				input := &mgmt.XDRAssetNoteInput{
					ID:         noteID,
					ResourceID: assetID,
				}
				if err := c.XDRAssetNoteDelete(cmd.Context(), input); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "ok", "noteId": noteID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Note %s deleted from asset %s\n", noteID, assetID)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&noteID, "note-id", "", "note ID (required)")
	cmd.Flags().StringVar(&assetID, "asset-id", "", "asset ID (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	_ = cmd.MarkFlagRequired("note-id")
	_ = cmd.MarkFlagRequired("asset-id")
	return markJSON(cmd)
}
