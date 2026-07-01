package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRulesEnableCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "enable <rule-id>...",
		Short: "Enable custom detection rules",
		Long: `Activate one or more custom detection rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would enable %s. Pass --yes to apply.\n",
					pluralize(len(args), "rule"))
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			filter := mgmt.RuleActionFilter{IDs: args}
			affected, err := c.RulesEnable(cmd.Context(), filter)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Enabled %s\n", pluralize(affected, "rule"))
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

func newRulesDisableCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "disable <rule-id>...",
		Short: "Disable custom detection rules",
		Long: `Deactivate one or more custom detection rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would disable %s. Pass --yes to apply.\n",
					pluralize(len(args), "rule"))
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			filter := mgmt.RuleActionFilter{IDs: args}
			affected, err := c.RulesDisable(cmd.Context(), filter)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Disabled %s\n", pluralize(affected, "rule"))
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}
