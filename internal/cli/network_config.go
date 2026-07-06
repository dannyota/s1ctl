package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newNetworkConfigurationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configuration",
		Short: "Get or set network quarantine control configuration",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newNetworkConfigurationGetCmd())
	cmd.AddCommand(newNetworkConfigurationSetCmd())
	return cmd
}

// networkConfigScope resolves the account/site/group/tenant scope shared by the
// configuration get and set commands. Tenant is implied when nothing else is set.
func networkConfigScope(siteIDs, accountIDs, groupIDs []string) mgmt.FirewallConfigScope {
	scope := mgmt.FirewallConfigScope{
		AccountIDs: accountIDs,
		SiteIDs:    siteIDs,
		GroupIDs:   groupIDs,
	}
	if len(siteIDs) == 0 && len(accountIDs) == 0 && len(groupIDs) == 0 {
		t := true
		scope.Tenant = &t
	}
	return scope
}

func newNetworkConfigurationGetCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the network quarantine configuration for a scope",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			cfg, err := c.NetworkQuarantineConfigurationGet(cmd.Context(),
				networkConfigScope(siteIDs, accountIDs, groupIDs))
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), cfg)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"Enabled", fmt.Sprintf("%t", cfg.Enabled)},
				{"Location aware", fmt.Sprintf("%t", cfg.LocationAware)},
				{"Report blocked", fmt.Sprintf("%t", cfg.ReportBlocked)},
				{"Inherits", fmt.Sprintf("%t", cfg.Inherits)},
				{"Inherited from", orDash(cfg.InheritedFrom)},
				{"Inherit settings", fmt.Sprintf("%t", cfg.InheritSettings)},
				{"Inherit all rules", fmt.Sprintf("%t", cfg.InheritAllFirewallRules)},
				{"Selected tags", fmt.Sprintf("%d", len(cfg.SelectedTags))},
			})
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope: site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope: account IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope: group IDs")
	return markJSON(cmd)
}

func newNetworkConfigurationSetCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var enabled, locationAware, reportBlocked bool
	var selectedTags []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Update the network quarantine configuration for a scope",
		Long: `Update network quarantine control configuration for a scope.

Only flags you pass are changed; boolean toggles are sent only when set.
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var data mgmt.FirewallConfigurationUpdate
			var changes []string
			if cmd.Flags().Changed("enabled") {
				data.Enabled = &enabled
				changes = append(changes, fmt.Sprintf("enabled=%t", enabled))
			}
			if cmd.Flags().Changed("location-aware") {
				data.LocationAware = &locationAware
				changes = append(changes, fmt.Sprintf("locationAware=%t", locationAware))
			}
			if cmd.Flags().Changed("report-blocked") {
				data.ReportBlocked = &reportBlocked
				changes = append(changes, fmt.Sprintf("reportBlocked=%t", reportBlocked))
			}
			if cmd.Flags().Changed("selected-tag") {
				data.SelectedTags = selectedTags
				changes = append(changes, fmt.Sprintf("selectedTags=%d", len(selectedTags)))
			}
			if len(changes) == 0 {
				return fmt.Errorf("nothing to update: set at least one of --enabled, --location-aware, --report-blocked, or --selected-tag")
			}

			scope := networkConfigScope(siteIDs, accountIDs, groupIDs)
			return guard(cmd.OutOrStdout(), "network configuration set",
				"update network quarantine configuration",
				fmt.Sprintf("%v", changes), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					cfg, err := c.NetworkQuarantineConfigurationUpdate(cmd.Context(), scope, data)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), cfg)
					}
					fmt.Fprintln(cmd.OutOrStdout(), "Updated network quarantine configuration")
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope: site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope: account IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope: group IDs")
	cmd.Flags().BoolVar(&enabled, "enabled", false, "enable network quarantine control for the scope")
	cmd.Flags().BoolVar(&locationAware, "location-aware", false, "enable location awareness for the scope")
	cmd.Flags().BoolVar(&reportBlocked, "report-blocked", false, "report blocked events")
	cmd.Flags().StringSliceVar(&selectedTags, "selected-tag", nil, "selected tag IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}
