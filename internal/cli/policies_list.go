package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func newPoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policies",
		Short: "View endpoint policies",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newPoliciesGetCmd())
	return cmd
}

func newPoliciesGetCmd() *cobra.Command {
	var siteID, accountID, groupID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get policy for a scope (site, account, or group)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			var raw json.RawMessage
			switch {
			case groupID != "" && siteID != "":
				p, pErr := c.PolicyGetGroup(cmd.Context(), siteID, groupID)
				if pErr != nil {
					return pErr
				}
				raw = p.Raw
			case siteID != "":
				p, pErr := c.PolicyGetSite(cmd.Context(), siteID)
				if pErr != nil {
					return pErr
				}
				raw = p.Raw
			case accountID != "":
				p, pErr := c.PolicyGetAccount(cmd.Context(), accountID)
				if pErr != nil {
					return pErr
				}
				raw = p.Raw
			default:
				return fmt.Errorf("specify --site-id, --account-id, or both --site-id and --group-id")
			}
			return printJSON(cmd.OutOrStdout(), raw)
		},
	}
	cmd.Flags().StringVar(&siteID, "site-id", "", "site ID")
	cmd.Flags().StringVar(&accountID, "account-id", "", "account ID")
	cmd.Flags().StringVar(&groupID, "group-id", "", "group ID (requires --site-id)")
	return cmd
}
