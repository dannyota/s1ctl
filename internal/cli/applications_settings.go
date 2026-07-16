package cli

import (
	"fmt"
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
				b := enableStr == "true"
				enableAppControl = &b
			}
			if inheritStr != "" {
				b := inheritStr == "true"
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
