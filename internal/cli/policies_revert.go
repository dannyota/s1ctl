package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newPoliciesRevertCmd() *cobra.Command {
	var scope, id, siteID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "revert",
		Short: "Revert a policy to its parent inherited values",
		Long: `Reset an endpoint policy to the values inherited from its parent scope.

Site policies revert to their account's policy, group policies revert to their
site's policy, and account policies revert to global defaults.

Specify the scope with --scope (site, account, or group) and the target with --id.
For group scope, --site-id is also required.
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
				if siteID == "" {
					return fmt.Errorf("--site-id is required when --scope=group")
				}
				return revertGroupPolicy(cmd, siteID, id, yes)
			default:
				return fmt.Errorf("invalid --scope %q: must be site, account, or group", scope)
			}
		},
	}
	cmd.Flags().StringVar(&scope, "scope", "site", "policy scope: site, account, or group")
	cmd.Flags().StringVar(&id, "id", "", "target scope ID (site, account, or group ID)")
	cmd.Flags().StringVar(&siteID, "site-id", "", "site ID (required for group scope)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the revert (default: dry-run)")
	return cmd
}

func revertSitePolicy(cmd *cobra.Command, siteID string, yes bool) error {
	if !yes {
		fmt.Fprintf(cmd.OutOrStdout(), "Would revert policy for site %s to account inherited values. Pass --yes to apply.\n", siteID)
		return nil
	}
	c, err := mgmtClient()
	if err != nil {
		return err
	}
	if err := c.PolicyRevertSite(cmd.Context(), siteID); err != nil {
		return fmt.Errorf("revert site %s policy: %w", siteID, err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Reverted policy for site %s\n", siteID)
	return nil
}

func revertAccountPolicy(cmd *cobra.Command, accountID string, yes bool) error {
	if !yes {
		fmt.Fprintf(cmd.OutOrStdout(), "Would revert policy for account %s to global inherited values. Pass --yes to apply.\n", accountID)
		return nil
	}
	c, err := mgmtClient()
	if err != nil {
		return err
	}
	if err := c.PolicyRevertAccount(cmd.Context(), accountID); err != nil {
		return fmt.Errorf("revert account %s policy: %w", accountID, err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Reverted policy for account %s\n", accountID)
	return nil
}

func revertGroupPolicy(cmd *cobra.Command, siteID, groupID string, yes bool) error {
	if !yes {
		fmt.Fprintf(cmd.OutOrStdout(), "Would revert policy for group %s to site inherited values. Pass --yes to apply.\n", groupID)
		return nil
	}
	c, err := mgmtClient()
	if err != nil {
		return err
	}
	if err := c.PolicyRevertGroup(cmd.Context(), siteID, groupID); err != nil {
		return fmt.Errorf("revert group %s policy: %w", groupID, err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Reverted policy for group %s\n", groupID)
	return nil
}
