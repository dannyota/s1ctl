package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPoliciesRevertCmd() *cobra.Command {
	var scope, id string
	var yes bool

	cmd := &cobra.Command{
		Use:   "revert",
		Short: "Revert a policy to its parent inherited values",
		Long: `Reset an endpoint policy to the values inherited from its parent scope.

Site policies revert to their account's policy, group policies revert to their
site's policy, and account policies revert to global defaults.

Specify the scope with --scope (site, account, or group) and the target with --id.
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if id == "" {
				return fmt.Errorf("--id is required")
			}

			switch scope {
			case "site":
				return revertSitePolicy(cmd, id, yes)
			case "account":
				return revertAccountPolicy(cmd, id, yes)
			case "group":
				return revertGroupPolicy(cmd, id, yes)
			default:
				return fmt.Errorf("invalid --scope %q: must be site, account, or group", scope)
			}
		},
	}
	cmd.Flags().StringVar(&scope, "scope", "site", "policy scope: site, account, or group")
	cmd.Flags().StringVar(&id, "id", "", "target scope ID (site, account, or group ID)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the revert (default: dry-run)")
	return markJSON(cmd)
}

func revertSitePolicy(cmd *cobra.Command, siteID string, yes bool) error {
	return guard(cmd.OutOrStdout(), "policies revert", "revert policy for site "+siteID+" to account inherited values", siteID, yes, func() error {
		c, err := mgmtClient()
		if err != nil {
			return err
		}
		if err := c.PolicyRevertSite(cmd.Context(), siteID); err != nil {
			return fmt.Errorf("revert site %s policy: %w", siteID, err)
		}
		if outputFormat == "json" {
			return printJSON(cmd.OutOrStdout(), map[string]any{"success": true, "scope": "site", "id": siteID})
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Reverted policy for site %s\n", siteID)
		return nil
	})
}

func revertAccountPolicy(cmd *cobra.Command, accountID string, yes bool) error {
	return guard(cmd.OutOrStdout(), "policies revert", "revert policy for account "+accountID+" to global inherited values", accountID, yes, func() error {
		c, err := mgmtClient()
		if err != nil {
			return err
		}
		if err := c.PolicyRevertAccount(cmd.Context(), accountID); err != nil {
			return fmt.Errorf("revert account %s policy: %w", accountID, err)
		}
		if outputFormat == "json" {
			return printJSON(cmd.OutOrStdout(), map[string]any{"success": true, "scope": "account", "id": accountID})
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Reverted policy for account %s\n", accountID)
		return nil
	})
}

func revertGroupPolicy(cmd *cobra.Command, groupID string, yes bool) error {
	return guard(cmd.OutOrStdout(), "policies revert", "revert policy for group "+groupID+" to site inherited values", groupID, yes, func() error {
		c, err := mgmtClient()
		if err != nil {
			return err
		}
		if err := c.PolicyRevertGroup(cmd.Context(), groupID); err != nil {
			return fmt.Errorf("revert group %s policy: %w", groupID, err)
		}
		if outputFormat == "json" {
			return printJSON(cmd.OutOrStdout(), map[string]any{"success": true, "scope": "group", "id": groupID})
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Reverted policy for group %s\n", groupID)
		return nil
	})
}
