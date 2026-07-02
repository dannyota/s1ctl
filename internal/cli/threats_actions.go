package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addThreatActions(parent *cobra.Command) {
	parent.AddCommand(newThreatMitigateCmd())
	parent.AddCommand(newThreatActionCmd("verdict", "Update analyst verdict on a threat",
		"--verdict", "analyst verdict (true_positive, false_positive, suspicious, undefined)", func(c *mgmt.Client, cmd *cobra.Command, val string, f mgmt.ActionFilter) (int, error) {
			return c.ThreatsUpdateVerdict(cmd.Context(), val, f)
		}))
	parent.AddCommand(newThreatActionCmd("status", "Update incident status on a threat",
		"--status", "incident status (unresolved, in_progress, resolved)", func(c *mgmt.Client, cmd *cobra.Command, val string, f mgmt.ActionFilter) (int, error) {
			return c.ThreatsUpdateStatus(cmd.Context(), val, f)
		}))
	parent.AddCommand(newThreatPlainActionCmd("blacklist", "Add the threat file hash to the blacklist", (*mgmt.Client).ThreatsAddToBlacklist))
	parent.AddCommand(newThreatPlainActionCmd("fetch-file", "Fetch the threat file from the endpoint to the console", (*mgmt.Client).ThreatsFetchFile))
}

func newThreatPlainActionCmd(verb, short string, call func(*mgmt.Client, context.Context, mgmt.ActionFilter) (int, error)) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <threat-id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "threats "+verb, verb+" threat "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := call(c, cmd.Context(), mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", verb, pluralize(affected, "threat"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
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
			return guard(cmd.OutOrStdout(), "threats "+verb, fmt.Sprintf("set %s=%s on threat %s", verb, val, args[0]), args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := fn(c, cmd, val, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", verb, pluralize(affected, "threat"))
				return nil
			})
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
			return guard(cmd.OutOrStdout(), "threats mitigate", action+" threat "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.ThreatsMitigate(cmd.Context(), action, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", action, pluralize(affected, "threat"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&action, "action", "", "mitigation action (kill, quarantine, remediate, rollback-remediation)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
