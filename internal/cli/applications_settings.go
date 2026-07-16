package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAppControlSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Manage application control settings",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAppControlSettingsGetCmd())
	cmd.AddCommand(newAppControlSettingsUpdateCmd())
	return cmd
}

func newAppControlSettingsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get application control settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			s, err := c.AppControlSettingsGet(cmd.Context())
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), s)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"Fallback Behavior", string(s.FallbackBehavior)},
				{"Application Control Enabled", fmt.Sprintf("%t", s.EnableApplicationControl)},
				{"Inherit Settings", fmt.Sprintf("%t", s.InheritApplicationControl)},
			})
			return nil
		},
	}
	return markJSON(cmd)
}

func newAppControlSettingsUpdateCmd() *cobra.Command {
	var fallbackBehavior, scopeType string
	var scopeIDs []string
	var yes bool

	// Use string flags and convert to *bool manually since cobra does not
	// natively support optional boolean flags with nil-default.
	var enableStr, inheritStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update application control settings",
		Long: `Update application control (NAC) settings.
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fallbackBehavior == "" && enableStr == "" && inheritStr == "" {
				return fmt.Errorf("at least one of --fallback-behavior, --enable, or --inherit is required")
			}

			var enableAppControl, inheritAppControl *bool
			if enableStr != "" {
				b, pErr := strconv.ParseBool(enableStr)
				if pErr != nil {
					return fmt.Errorf("invalid --enable value %q: expected true or false", enableStr)
				}
				enableAppControl = &b
			}
			if inheritStr != "" {
				b, pErr := strconv.ParseBool(inheritStr)
				if pErr != nil {
					return fmt.Errorf("invalid --inherit value %q: expected true or false", inheritStr)
				}
				inheritAppControl = &b
			}

			return guard(cmd.OutOrStdout(), "applications settings update",
				"update application control settings",
				"settings", yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					input := mgmt.AppControlSettingsInput{
						EnableApplicationControl:  enableAppControl,
						InheritApplicationControl: inheritAppControl,
					}
					if fallbackBehavior != "" {
						input.FallbackBehavior = mgmt.AppControlBehavior(strings.ToUpper(fallbackBehavior))
					}
					if scopeType != "" && len(scopeIDs) > 0 {
						input.Scope = &mgmt.AppControlScope{
							ScopeType: mgmt.AppControlScopeLevel(strings.ToUpper(scopeType)),
							ScopeIDs:  scopeIDs,
						}
					}
					resp, err := c.AppControlSettingsUpdate(cmd.Context(), input)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), resp)
					}
					fmt.Fprintln(cmd.OutOrStdout(), "Updated application control settings")
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&fallbackBehavior, "fallback-behavior", "", "default behavior: allow, monitor, block")
	cmd.Flags().StringVar(&enableStr, "enable", "", "enable application control (true/false)")
	cmd.Flags().StringVar(&inheritStr, "inherit", "", "inherit settings from parent (true/false)")
	cmd.Flags().StringVar(&scopeType, "scope-type", "", "scope type: account, site, group")
	cmd.Flags().StringSliceVar(&scopeIDs, "scope-id", nil, "scope IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

func newAppMgmtSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mgmt-settings",
		Short: "Manage application management settings (scan schedule, extensive scan)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAppMgmtSettingsGetCmd())
	cmd.AddCommand(newAppMgmtSettingsUpdateCmd())
	return cmd
}

func newAppMgmtSettingsGetCmd() *cobra.Command {
	var siteIDs, groupIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get application management settings",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AppMgmtSettingsListParams{
				SiteIDs:    siteIDs,
				GroupIDs:   groupIDs,
				AccountIDs: accountIDs,
			}
			s, err := c.AppMgmtSettingsGet(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), s)
			}
			schedule := "--"
			if s.ScanSchedule != nil {
				schedule = fmt.Sprintf("every %d weeks on %s at %s (%s)",
					s.ScanSchedule.ScanEvery, s.ScanSchedule.RepeatOn,
					s.ScanSchedule.Time, s.ScanSchedule.Timezone)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"Extensive Scan", fmt.Sprintf("%t", s.ExtensiveScanEnabled)},
				{"Default Policy", fmt.Sprintf("%t", s.IsDefaultPolicy)},
				{"Scan Schedule", schedule},
				{"Breaking Inheritance", fmt.Sprintf("%t", s.HasBreakingInheritance)},
			})
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

func newAppMgmtSettingsUpdateCmd() *cobra.Command {
	var siteIDs, groupIDs, accountIDs []string
	var extensiveScanStr, defaultPolicyStr string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update application management settings",
		Long: `Update application management settings (scan schedule, extensive scan).
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if extensiveScanStr == "" && defaultPolicyStr == "" {
				return fmt.Errorf("at least one of --extensive-scan or --default-policy is required")
			}

			data := mgmt.AppMgmtSettingsUpdateData{}
			if extensiveScanStr != "" {
				b, pErr := strconv.ParseBool(extensiveScanStr)
				if pErr != nil {
					return fmt.Errorf("invalid --extensive-scan value %q: expected true or false", extensiveScanStr)
				}
				data.ExtensiveScanEnabled = &b
			}
			if defaultPolicyStr != "" {
				b, pErr := strconv.ParseBool(defaultPolicyStr)
				if pErr != nil {
					return fmt.Errorf("invalid --default-policy value %q: expected true or false", defaultPolicyStr)
				}
				data.IsDefaultPolicy = &b
			}

			return guard(cmd.OutOrStdout(), "applications mgmt-settings update",
				"update application management settings",
				"mgmt-settings", yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					scope := mgmt.AppMgmtSettingsScope{
						SiteIDs:    siteIDs,
						GroupIDs:   groupIDs,
						AccountIDs: accountIDs,
					}
					if len(siteIDs) == 0 && len(groupIDs) == 0 && len(accountIDs) == 0 {
						scope.Tenant = true
					}
					s, err := c.AppMgmtSettingsUpdate(cmd.Context(), scope, data)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), s)
					}
					fmt.Fprintln(cmd.OutOrStdout(), "Updated application management settings")
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope: site IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope: group IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope: account IDs")
	cmd.Flags().StringVar(&extensiveScanStr, "extensive-scan", "", "enable extensive scan (true/false)")
	cmd.Flags().StringVar(&defaultPolicyStr, "default-policy", "", "use default policy (true/false)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

func newAppControlLabelsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "labels",
		Short: "Manage application control labels",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAppControlLabelsListCmd())
	return cmd
}

func newAppControlLabelsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List application control labels",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			labels, err := c.AppControlLabelsList(cmd.Context())
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name"}
			rows := make([][]string, len(labels))
			for i, l := range labels {
				rows[i] = []string{l.ID, l.LabelName}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, labels, len(labels), len(labels), "label", false)
		},
	}
	return markJSON(cmd)
}
