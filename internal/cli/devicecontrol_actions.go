package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newDeviceControlEnableCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "enable <rule-id>...",
		Short: "Enable device control rules",
		Long: `Enable one or more device control rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "devicecontrol enable",
				"enable "+pluralize(len(args), "device rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.DeviceRulesSetStatus(cmd.Context(), args, mgmt.DeviceRuleStatusEnabled)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Enabled %s\n", pluralize(affected, "device rule"))
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

func newDeviceControlDisableCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "disable <rule-id>...",
		Short: "Disable device control rules",
		Long: `Disable one or more device control rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "devicecontrol disable",
				"disable "+pluralize(len(args), "device rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.DeviceRulesSetStatus(cmd.Context(), args, mgmt.DeviceRuleStatusDisabled)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Disabled %s\n", pluralize(affected, "device rule"))
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

func newDeviceControlReorderCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "reorder <id:order>...",
		Short: "Reorder device control rules",
		Long: `Change the evaluation order of device control rules.

Each argument is an id:order pair, for example:

  s1ctl devicecontrol reorder 123:1 456:2 789:3

The order determines rule evaluation priority (1 = first).
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			orders, err := parseRuleOrders(args)
			if err != nil {
				return err
			}

			return guard(cmd.OutOrStdout(), "devicecontrol reorder",
				"reorder "+pluralize(len(orders), "device rule"),
				fmt.Sprintf("%d rules", len(orders)), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					filter := mgmt.DeviceRuleReorderFilter{
						AccountIDs: accountIDs,
						SiteIDs:    siteIDs,
						GroupIDs:   groupIDs,
					}
					if len(siteIDs) == 0 && len(accountIDs) == 0 && len(groupIDs) == 0 {
						t := true
						filter.Tenant = &t
					}
					if err := c.DeviceRulesReorder(cmd.Context(), orders, filter); err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": len(orders)})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Reordered %s\n", pluralize(len(orders), "device rule"))
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope: site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope: account IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope: group IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

func newDeviceControlCopyCmd() *cobra.Command {
	var sourceSiteIDs, sourceAccountIDs []string
	var targetSiteID, targetAccountID string
	var targetGroupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy device control rules between scopes",
		Long: `Copy device control rules from a source scope to a target scope.

Use --source-site-id or --source-account-id to define the source, and
--target-site-id, --target-account-id, or --target-group-id for the destination.
At least one target flag is required.
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var targets []string
			if targetSiteID != "" {
				targets = append(targets, "site "+targetSiteID)
			}
			if targetAccountID != "" {
				targets = append(targets, "account "+targetAccountID)
			}
			if len(targetGroupIDs) > 0 {
				targets = append(targets, "group "+strings.Join(targetGroupIDs, ","))
			}
			if len(targets) == 0 {
				return fmt.Errorf("at least one of --target-site-id, --target-account-id, or --target-group-id is required")
			}
			targetDesc := strings.Join(targets, ", ")

			return guard(cmd.OutOrStdout(), "devicecontrol copy",
				"copy device rules to "+targetDesc,
				targetDesc, yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					filter := mgmt.DeviceRuleScopeFilter{
						SiteIDs:    sourceSiteIDs,
						AccountIDs: sourceAccountIDs,
					}
					if len(sourceSiteIDs) == 0 && len(sourceAccountIDs) == 0 {
						t := true
						filter.Tenant = &t
					}
					target := mgmt.DeviceRuleCopyTarget{
						GroupIDs: targetGroupIDs,
					}
					if targetSiteID != "" {
						target.SiteID = &targetSiteID
					}
					if targetAccountID != "" {
						target.AccountID = &targetAccountID
					}
					affected, err := c.DeviceRulesCopy(cmd.Context(), filter, []mgmt.DeviceRuleCopyTarget{target})
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Copied %s\n", pluralize(affected, "device rule"))
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&sourceSiteIDs, "source-site-id", nil, "source site IDs")
	cmd.Flags().StringSliceVar(&sourceAccountIDs, "source-account-id", nil, "source account IDs")
	cmd.Flags().StringVar(&targetSiteID, "target-site-id", "", "target site ID")
	cmd.Flags().StringVar(&targetAccountID, "target-account-id", "", "target account ID")
	cmd.Flags().StringSliceVar(&targetGroupIDs, "target-group-id", nil, "target group IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

// parseRuleOrders parses "id:order" pairs from command-line arguments.
func parseRuleOrders(args []string) ([]mgmt.RuleOrder, error) {
	orders := make([]mgmt.RuleOrder, len(args))
	for i, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format %q: expected id:order", arg)
		}
		order, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid order in %q: %w", arg, err)
		}
		orders[i] = mgmt.RuleOrder{ID: parts[0], Order: order}
	}
	return orders, nil
}
