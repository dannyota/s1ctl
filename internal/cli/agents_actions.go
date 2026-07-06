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
	parent.AddCommand(newAgentMoveToSiteCmd())
	parent.AddCommand(newAgentSetExternalIDCmd())
	parent.AddCommand(newAgentFirewallLoggingCmd())
	parent.AddCommand(newAgentBroadcastCmd())
	parent.AddCommand(newAgentFetchFilesCmd())
	parent.AddCommand(newAgentRangerCmd())
	parent.AddCommand(newAgentLocalUpgradeCmd())
	parent.AddCommand(newAgentLocalUpgradeStatusCmd())
	parent.AddCommand(newAgentsPassphrasesCmd())
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
		{"reset-passphrase", "Reset the agent maintenance passphrase", (*mgmt.Client).AgentsResetPassphrase},
		{"fetch-installed-apps", "Fetch the installed-applications inventory", (*mgmt.Client).AgentsFetchInstalledApps},
		{"fetch-firewall-rules", "Fetch the current firewall-rules inventory", (*mgmt.Client).AgentsFetchFirewallRules},
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
	return markJSON(cmd)
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
	return markJSON(cmd)
}

func newAgentMoveToSiteCmd() *cobra.Command {
	var siteID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "move-to-site <agent-id> --site-id <target-site-id>",
		Short: "Move an agent to a different site",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if siteID == "" {
				return fmt.Errorf("--site-id is required")
			}
			return guard(cmd.OutOrStdout(), "agents move-to-site", "move agent "+args[0]+" to site "+siteID, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsMoveToSite(cmd.Context(), siteID, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "move-to-site: %s affected\n", pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&siteID, "site-id", "", "target site ID (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newAgentSetExternalIDCmd() *cobra.Command {
	var externalID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "set-external-id <agent-id> --external-id <value>",
		Short: "Set the external ID on an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if externalID == "" {
				return fmt.Errorf("--external-id is required")
			}
			return guard(cmd.OutOrStdout(), "agents set-external-id", "set external ID "+externalID+" on agent "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsSetExternalID(cmd.Context(), externalID, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "set-external-id: %s affected\n", pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&externalID, "external-id", "", "external ID value (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newAgentFirewallLoggingCmd() *cobra.Command {
	var state string
	var yes bool

	cmd := &cobra.Command{
		Use:   "firewall-logging <agent-id> --state on|off",
		Short: "Enable or disable firewall logging on an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if state != "on" && state != "off" {
				return fmt.Errorf(`--state must be "on" or "off"`)
			}
			return guard(cmd.OutOrStdout(), "agents firewall-logging", "turn firewall logging "+state+" for agent "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsFirewallLogging(cmd.Context(), state == "on", mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "firewall-logging %s: %s affected\n", state, pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&state, "state", "", `"on" or "off" (required)`)
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newAgentBroadcastCmd() *cobra.Command {
	var message string
	var yes bool

	cmd := &cobra.Command{
		Use:   "broadcast <agent-id> --message <text>",
		Short: "Display a broadcast message on an agent's endpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if message == "" {
				return fmt.Errorf("--message is required")
			}
			return guard(cmd.OutOrStdout(), "agents broadcast", "broadcast a message to agent "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsBroadcast(cmd.Context(), message, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "broadcast: %s affected\n", pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&message, "message", "", "message text to broadcast (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newAgentFetchFilesCmd() *cobra.Command {
	var paths []string
	var password string
	var yes bool

	cmd := &cobra.Command{
		Use:   "fetch-files <agent-id> --path <file> [--path <file>...] [--password <pw>]",
		Short: "Fetch specific files from an agent to the console",
		Long: `Fetch up to 10 files from a single agent. The files are uploaded to the
console encrypted with --password (required by the platform to open the
resulting archive). The password is never written to the audit log.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(paths) == 0 {
				return fmt.Errorf("--path is required")
			}
			// The action string is deliberately generic: the file password
			// must never reach the audit log.
			return guard(cmd.OutOrStdout(), "agents fetch-files", "fetch files from agent "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				ok, err := c.AgentsFetchFiles(cmd.Context(), args[0], paths, password)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]bool{"success": ok})
				}
				status := "request accepted"
				if !ok {
					status = "request not accepted"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "fetch-files: %s\n", status)
				return nil
			})
		},
	}
	cmd.Flags().StringArrayVar(&paths, "path", nil, "absolute file path to fetch (repeatable, up to 10) (required)")
	cmd.Flags().StringVar(&password, "password", "", "archive encryption password")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newAgentRangerCmd() *cobra.Command {
	var state string
	var yes bool

	cmd := &cobra.Command{
		Use:   "ranger <agent-id> --state on|off",
		Short: "Enable or disable Ranger network discovery on an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if state != "on" && state != "off" {
				return fmt.Errorf(`--state must be "on" or "off"`)
			}
			return guard(cmd.OutOrStdout(), "agents ranger", "turn Ranger "+state+" for agent "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsRanger(cmd.Context(), state == "on", mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "ranger %s: %s affected\n", state, pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&state, "state", "", `"on" or "off" (required)`)
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newAgentLocalUpgradeCmd() *cobra.Command {
	var state, until string
	var yes bool

	cmd := &cobra.Command{
		Use:   "local-upgrade <agent-id> --state on|off [--until <timestamp>]",
		Short: "Authorize or revoke local upgrade/downgrade on an agent",
		Long: `Set an agent's local upgrade/downgrade authorization.

--state on authorizes local upgrades until the --until expiration timestamp
(RFC3339, e.g. 2030-01-01T00:00:00Z), which is required. --state off revokes
the authorization.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if state != "on" && state != "off" {
				return fmt.Errorf(`--state must be "on" or "off"`)
			}
			if state == "on" && until == "" {
				return fmt.Errorf("--until is required when --state is on")
			}
			authorization := ""
			if state == "on" {
				authorization = until
			}
			return guard(cmd.OutOrStdout(), "agents local-upgrade", "set local upgrade authorization "+state+" for agent "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsLocalUpgradeAuthorization(cmd.Context(), mgmt.ActionFilter{IDs: []string{args[0]}}, authorization)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "local-upgrade %s: %s affected\n", state, pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&state, "state", "", `"on" or "off" (required)`)
	cmd.Flags().StringVar(&until, "until", "", "authorization expiration timestamp (RFC3339); required with --state on")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}
