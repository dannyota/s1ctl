package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newFirewallEnableCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "enable <rule-id>...",
		Short: "Enable firewall rules",
		Long: `Enable one or more firewall rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "firewall enable",
				"enable "+pluralize(len(args), "firewall rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.FirewallRulesSetStatus(cmd.Context(), args, mgmt.FirewallStatusEnabled)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Enabled %s\n", pluralize(affected, "firewall rule"))
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

func newFirewallDisableCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "disable <rule-id>...",
		Short: "Disable firewall rules",
		Long: `Disable one or more firewall rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "firewall disable",
				"disable "+pluralize(len(args), "firewall rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.FirewallRulesSetStatus(cmd.Context(), args, mgmt.FirewallStatusDisabled)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Disabled %s\n", pluralize(affected, "firewall rule"))
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

func newFirewallReorderCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "reorder <id:order>...",
		Short: "Reorder firewall rules",
		Long: `Change the evaluation order of firewall rules.

Each argument is an id:order pair, for example:

  s1ctl firewall reorder 123:1 456:2 789:3

The order determines rule evaluation priority (1 = first).
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orders, err := parseRuleOrders(args)
			if err != nil {
				return err
			}

			return guard(cmd.OutOrStdout(), "firewall reorder",
				"reorder "+pluralize(len(orders), "firewall rule"),
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
					if err := c.FirewallRulesReorder(cmd.Context(), orders, filter); err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": len(orders)})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Reordered %s\n", pluralize(len(orders), "firewall rule"))
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

func newFirewallCopyCmd() *cobra.Command {
	var sourceSiteIDs, sourceAccountIDs []string
	var targetSiteID, targetAccountID, targetGroupID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy firewall rules between scopes",
		Long: `Copy firewall rules from a source scope to a target scope.

Use --source-site-id or --source-account-id to define the source, and
--target-site-id, --target-account-id, or --target-group-id for the destination.
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return guard(cmd.OutOrStdout(), "firewall copy",
				"copy firewall rules to target scope",
				"target scope", yes, func() error {
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
					target := mgmt.FirewallRuleCopyTarget{}
					if targetSiteID != "" {
						target.SiteID = &targetSiteID
					}
					if targetAccountID != "" {
						target.AccountID = &targetAccountID
					}
					if targetGroupID != "" {
						target.GroupID = &targetGroupID
					}
					affected, err := c.FirewallRulesCopy(cmd.Context(), filter, []mgmt.FirewallRuleCopyTarget{target})
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Copied %s\n", pluralize(affected, "firewall rule"))
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
