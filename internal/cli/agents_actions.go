package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addAgentActions(parent *cobra.Command) {
	parent.AddCommand(newAgentActionCmd("isolate", "Network-isolate an agent", func(c *mgmt.Client, cmd *cobra.Command, f mgmt.ActionFilter) (int, error) {
		return c.AgentsDisconnect(cmd.Context(), f)
	}))
	parent.AddCommand(newAgentActionCmd("connect", "Reconnect an isolated agent", func(c *mgmt.Client, cmd *cobra.Command, f mgmt.ActionFilter) (int, error) {
		return c.AgentsConnect(cmd.Context(), f)
	}))
	parent.AddCommand(newAgentActionCmd("scan", "Start full disk scan", func(c *mgmt.Client, cmd *cobra.Command, f mgmt.ActionFilter) (int, error) {
		return c.AgentsInitiateScan(cmd.Context(), f)
	}))
	parent.AddCommand(newAgentActionCmd("decommission", "Decommission an agent", func(c *mgmt.Client, cmd *cobra.Command, f mgmt.ActionFilter) (int, error) {
		return c.AgentsDecommission(cmd.Context(), f)
	}))
}

type agentActionFn func(*mgmt.Client, *cobra.Command, mgmt.ActionFilter) (int, error)

func newAgentActionCmd(verb, short string, fn agentActionFn) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <agent-id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would %s agent %s. Pass --yes to apply.\n", verb, args[0])
				return nil
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			affected, err := fn(c, cmd, mgmt.ActionFilter{IDs: []string{args[0]}})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(map[string]int{"affected": affected})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", verb, pluralize(affected, "agent"))
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
