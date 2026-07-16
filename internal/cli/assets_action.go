package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAssetsActionCmd() *cobra.Command {
	var (
		assetType  string
		actionName string
		ids        []string
		yes        bool
	)

	cmd := &cobra.Command{
		Use:   "action",
		Short: "Perform an action on assets",
		Long: `Perform an action on one or more assets.

Action names are passed through to the API (e.g. mark_asset_criticality_high,
mark_asset_criticality_medium). Dry-run by default; pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if len(ids) == 0 {
				return fmt.Errorf("at least one --id is required")
			}
			target := strings.Join(ids, ",")
			return guard(cmd.OutOrStdout(), "assets action", actionName+" on "+pluralize(len(ids), "asset"), target, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				input := &mgmt.XDRAssetActionInput{
					ActionName: actionName,
					IDIn:       ids,
				}
				affected, err := c.XDRAssetAction(cmd.Context(), mgmt.AssetType(assetType), input)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", actionName, pluralize(affected, "asset"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&assetType, "type", "", "asset type slug (omit for cross-type action)")
	cmd.Flags().StringVar(&actionName, "action", "", "action name (required)")
	cmd.Flags().StringSliceVar(&ids, "id", nil, "asset ID(s) to act on (required, repeatable)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	_ = cmd.MarkFlagRequired("action")
	_ = cmd.MarkFlagRequired("id")
	return markJSON(cmd)
}
