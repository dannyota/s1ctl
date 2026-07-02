package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addAgentActions(parent *cobra.Command) {
	parent.AddCommand(newAgentsIsolateCmd())
	parent.AddCommand(newAgentsReconnectCmd())
	parent.AddCommand(newAgentMoveCmd())
	plain := []struct {
		verb, short string
		call        func(*mgmt.Client, context.Context, mgmt.ActionFilter) (int, error)
	}{
		{"scan", "Start full disk scan", (*mgmt.Client).AgentsInitiateScan},
		{"abort-scan", "Abort a running disk scan", (*mgmt.Client).AgentsAbortScan},
		{"decommission", "Decommission an agent", (*mgmt.Client).AgentsDecommission},
		{"uninstall", "Uninstall an agent", (*mgmt.Client).AgentsUninstall},
		{"shutdown", "Shut down the endpoint", (*mgmt.Client).AgentsShutdown},
		{"restart", "Restart the endpoint", (*mgmt.Client).AgentsRestartMachine},
		{"fetch-logs", "Fetch agent logs to the console", (*mgmt.Client).AgentsFetchLogs},
		{"enable", "Enable a disabled agent", (*mgmt.Client).AgentsEnableAgent},
		{"disable", "Disable an agent", (*mgmt.Client).AgentsDisableAgent},
		{"reset-config", "Reset agent local configuration", (*mgmt.Client).AgentsResetLocalConfig},
		{"approve-uninstall", "Approve a pending uninstall request", (*mgmt.Client).AgentsApproveUninstall},
		{"reject-uninstall", "Reject a pending uninstall request", (*mgmt.Client).AgentsRejectUninstall},
		{"mark-up-to-date", "Mark an agent as up to date", (*mgmt.Client).AgentsMarkUpToDate},
		{"randomize-uuid", "Randomize the agent UUID", (*mgmt.Client).AgentsRandomizeUUID},
	}
	for _, a := range plain {
		parent.AddCommand(newAgentActionCmd(a.verb, a.short, func(c *mgmt.Client, cmd *cobra.Command, f mgmt.ActionFilter) (int, error) {
			return a.call(c, cmd.Context(), f)
		}))
	}
}

type agentActionFn func(*mgmt.Client, *cobra.Command, mgmt.ActionFilter) (int, error)

func newAgentActionCmd(verb, short string, fn agentActionFn) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <agent-id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "agents "+verb, verb+" agent "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := fn(c, cmd, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", verb, pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newAgentMoveCmd() *cobra.Command {
	var (
		groupID string
		yes     bool
	)

	cmd := &cobra.Command{
		Use:   "move <agent-id> --group-id <target-group-id>",
		Short: "Move an agent to a different group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if groupID == "" {
				return fmt.Errorf("--group-id is required")
			}
			return guard(cmd.OutOrStdout(), "agents move", "move agent "+args[0]+" to group "+groupID, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsMoveToGroup(cmd.Context(), groupID, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "move: %s affected\n", pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&groupID, "group-id", "", "target group ID (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
