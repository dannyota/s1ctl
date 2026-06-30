package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addThreatActions(parent *cobra.Command) {
	parent.AddCommand(newThreatMitigateCmd())
	parent.AddCommand(newThreatActionCmd("verdict", "Update analyst verdict on a threat",
		"--verdict", "analyst verdict", func(c *mgmt.Client, cmd *cobra.Command, val string, f mgmt.ActionFilter) (int, error) {
			return c.ThreatsUpdateVerdict(cmd.Context(), val, f)
		}))
	parent.AddCommand(newThreatActionCmd("status", "Update incident status on a threat",
		"--status", "incident status", func(c *mgmt.Client, cmd *cobra.Command, val string, f mgmt.ActionFilter) (int, error) {
			return c.ThreatsUpdateStatus(cmd.Context(), val, f)
		}))
}

type threatActionFn func(*mgmt.Client, *cobra.Command, string, mgmt.ActionFilter) (int, error)

func newThreatActionCmd(verb, short, flagName, flagDesc string, fn threatActionFn) *cobra.Command {
	var val string
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <threat-id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if val == "" {
				return fmt.Errorf("%s is required", flagName)
			}
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would set %s=%s on threat %s. Pass --yes to apply.\n", verb, val, args[0])
				return nil
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			affected, err := fn(c, cmd, val, mgmt.ActionFilter{IDs: []string{args[0]}})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(map[string]int{"affected": affected})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", verb, pluralize(affected, "threat"))
			return nil
		},
	}
	cmd.Flags().StringVar(&val, verb, "", flagDesc)
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newThreatMitigateCmd() *cobra.Command {
	var action string
	var yes bool

	cmd := &cobra.Command{
		Use:   "mitigate <threat-id>",
		Short: "Apply mitigation action to a threat",
		Long:  "Actions: kill, quarantine, remediate, rollback-remediation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if action == "" {
				return fmt.Errorf("--action is required (kill, quarantine, remediate, rollback-remediation)")
			}
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would %s threat %s. Pass --yes to apply.\n", action, args[0])
				return nil
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			affected, err := c.ThreatsMitigate(cmd.Context(), action, mgmt.ActionFilter{IDs: []string{args[0]}})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(map[string]int{"affected": affected})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", action, pluralize(affected, "threat"))
			return nil
		},
	}
	cmd.Flags().StringVar(&action, "action", "", "mitigation action (kill, quarantine, remediate, rollback-remediation)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
