package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newNetworkDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <rule-id>...",
		Short: "Delete network quarantine rules",
		Long: `Delete one or more network quarantine rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "network delete",
				"delete "+pluralize(len(args), "network quarantine rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.NetworkQuarantineDelete(cmd.Context(), args)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s\n", pluralize(affected, "network quarantine rule"))
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

func newNetworkEnableCmd() *cobra.Command {
	return newNetworkStatusCmd("enable", "Enable", mgmt.FirewallStatusEnabled)
}

func newNetworkDisableCmd() *cobra.Command {
	return newNetworkStatusCmd("disable", "Disable", mgmt.FirewallStatusDisabled)
}

// newNetworkStatusCmd builds the enable/disable commands, which differ only by
// the target status and the verb used in output.
func newNetworkStatusCmd(verb, title string, status mgmt.FirewallStatus) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <rule-id>...",
		Short: title + " network quarantine rules",
		Long: title + ` one or more network quarantine rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "network "+verb,
				verb+" "+pluralize(len(args), "network quarantine rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.NetworkQuarantineSetStatus(cmd.Context(), args, status)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "%sd %s\n", title, pluralize(affected, "network quarantine rule"))
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

func newNetworkReorderCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "reorder <id:order>...",
		Short: "Reorder network quarantine rules",
		Long: `Change the evaluation order of network quarantine rules.

Each argument is an id:order pair, for example:

  s1ctl network reorder 123:1 456:2 789:3

The order determines rule evaluation priority (1 = first).
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orders, err := parseRuleOrders(args)
			if err != nil {
				return err
			}
			return guard(cmd.OutOrStdout(), "network reorder",
				"reorder "+pluralize(len(orders), "network quarantine rule"),
				fmt.Sprintf("%d rules", len(orders)), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					filter := mgmt.FirewallRuleReorderFilter{
						AccountIDs: accountIDs,
						SiteIDs:    siteIDs,
						GroupIDs:   groupIDs,
					}
					if len(siteIDs) == 0 && len(accountIDs) == 0 && len(groupIDs) == 0 {
						t := true
						filter.Tenant = &t
					}
					if err := c.NetworkQuarantineReorder(cmd.Context(), orders, filter); err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": len(orders)})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Reordered %s\n", pluralize(len(orders), "network quarantine rule"))
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope: site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope: account IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope: group IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

func newNetworkCopyCmd() *cobra.Command {
	var sourceSiteIDs, sourceAccountIDs []string
	var targetSiteID, targetAccountID, targetGroupID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy network quarantine rules between scopes",
		Long: `Copy network quarantine rules from a source scope to a target scope.

Use --source-site-id or --source-account-id to define the source, and
--target-site-id, --target-account-id, or --target-group-id for the destination.
At least one target flag is required.
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			targetDesc, err := describeCopyTarget(targetSiteID, targetAccountID, targetGroupID)
			if err != nil {
				return err
			}
			return guard(cmd.OutOrStdout(), "network copy",
				"copy network quarantine rules to "+targetDesc,
				targetDesc, yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					filter := mgmt.FirewallRuleReorderFilter{
						SiteIDs:    sourceSiteIDs,
						AccountIDs: sourceAccountIDs,
					}
					if len(sourceSiteIDs) == 0 && len(sourceAccountIDs) == 0 {
						t := true
						filter.Tenant = &t
					}
					target := buildCopyTarget(targetSiteID, targetAccountID, targetGroupID)
					affected, err := c.NetworkQuarantineCopy(cmd.Context(), filter, []mgmt.FirewallRuleCopyTarget{target})
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Copied %s\n", pluralize(affected, "network quarantine rule"))
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&sourceSiteIDs, "source-site-id", nil, "source site IDs")
	cmd.Flags().StringSliceVar(&sourceAccountIDs, "source-account-id", nil, "source account IDs")
	cmd.Flags().StringVar(&targetSiteID, "target-site-id", "", "target site ID")
	cmd.Flags().StringVar(&targetAccountID, "target-account-id", "", "target account ID")
	cmd.Flags().StringVar(&targetGroupID, "target-group-id", "", "target group ID")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

// describeCopyTarget renders a human description of the copy/move destination
// and errors if no target scope is given.
func describeCopyTarget(siteID, accountID, groupID string) (string, error) {
	var targets []string
	if siteID != "" {
		targets = append(targets, "site "+siteID)
	}
	if accountID != "" {
		targets = append(targets, "account "+accountID)
	}
	if groupID != "" {
		targets = append(targets, "group "+groupID)
	}
	if len(targets) == 0 {
		return "", fmt.Errorf("at least one of --target-site-id, --target-account-id, or --target-group-id is required")
	}
	return strings.Join(targets, ", "), nil
}

func buildCopyTarget(siteID, accountID, groupID string) mgmt.FirewallRuleCopyTarget {
	target := mgmt.FirewallRuleCopyTarget{}
	if siteID != "" {
		target.SiteID = &siteID
	}
	if accountID != "" {
		target.AccountID = &accountID
	}
	if groupID != "" {
		target.GroupID = &groupID
	}
	return target
}
