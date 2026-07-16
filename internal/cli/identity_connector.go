package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// --- connector ---

func newIdentityConnectorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connector",
		Short: "Manage AD connectors (Cloudlink agents)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newIdentityConnectorListCmd())
	cmd.AddCommand(newIdentityConnectorGetCmd())
	cmd.AddCommand(newIdentityConnectorReplaceCmd())
	cmd.AddCommand(newIdentityConnectorAgentsCmd())
	return cmd
}

func newIdentityConnectorListCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all AD connectors",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			connectors, err := c.IdentityConnectors(cmd.Context(), identityParams(siteIDs, accountIDs))
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), connectors)
			}
			headers := []string{"Cloudlink ID", "Computer", "Status", "Agent Type", "OS", "Version", "Domain", "IP"}
			rows := make([][]string, len(connectors))
			for i, cn := range connectors {
				rows[i] = []string{
					strconv.FormatInt(cn.CloudlinkID, 10),
					cn.ComputerName,
					string(cn.Status),
					cn.AgentType,
					truncate(cn.OSName, 30),
					cn.Version,
					cn.DomainName,
					cn.IPAddress,
				}
			}
			printTable(headers, rows)
			return nil
		},
	}
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	return markJSON(cmd)
}

func newIdentityConnectorGetCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the current AD connector configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			cn, err := c.IdentityConnector(cmd.Context(), identityParams(siteIDs, accountIDs))
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), cn)
			}
			rows := [][]string{
				{"GUID", cn.GUID},
				{"Computer", cn.ComputerName},
				{"Status", string(cn.Status)},
				{"Agent Type", cn.AgentType},
				{"OS", cn.OSName},
				{"Version", cn.Version},
				{"Domain", cn.DomainName},
				{"IP", cn.IPAddress},
				{"Unified Agent", boolIcon(cn.IsUnifiedAgent)},
				{"Last Seen", orDash(cn.LastSeen)},
				{"Scope", orDash(cn.ScopePath)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	return markJSON(cmd)
}

func newIdentityConnectorReplaceCmd() *cobra.Command {
	var yes bool
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "replace [agent-uuid]",
		Short: "Replace the AD connector with a different agent",
		Long: `Replace the AD connector (Cloudlink) with a new agent by UUID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "identity connector replace", "replace AD connector", args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.IdentityConnectorReplace(cmd.Context(), identityParams(siteIDs, accountIDs), args[0]); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "AD connector replaced.")
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	return cmd
}

func newIdentityConnectorAgentsCmd() *cobra.Command {
	var siteIDs, accountIDs []string
	var filterInput string

	cmd := &cobra.Command{
		Use:   "agents",
		Short: "List Windows agents available as connectors",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.WindowsAgentParams{
				SiteIDs:     strings.Join(siteIDs, ","),
				AccountIDs:  strings.Join(accountIDs, ","),
				FilterInput: filterInput,
			}
			agents, err := c.IdentityWindowsAgents(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), agents)
			}
			headers := []string{"UUID", "Host", "OS", "Version", "Status", "Domain", "IP"}
			rows := make([][]string, len(agents))
			for i, a := range agents {
				rows[i] = []string{
					a.UUID,
					a.HostName,
					truncate(a.OSName, 25),
					a.AgentVersion,
					a.Status,
					a.DomainName,
					a.IPAddress,
				}
			}
			printTable(headers, rows)
			return nil
		},
	}
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	cmd.Flags().StringVar(&filterInput, "filter", "", "filter agents by name")
	return markJSON(cmd)
}
