package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAgentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Manage endpoint agents",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAgentsListCmd())
	cmd.AddCommand(newAgentsGetCmd())
	cmd.AddCommand(newAgentsCountCmd())
	addAgentActions(cmd)
	return cmd
}

func newAgentsListCmd() *cobra.Command {
	var siteIDs, groupIDs, osTypes []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List agents",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			agents, pag, err := c.AgentsList(cmd.Context(), &mgmt.AgentListParams{
				SiteIDs:  siteIDs,
				GroupIDs: groupIDs,
				OSTypes:  osTypes,
				Query:    query,
				Limit:    limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(agents)
			}
			var rows [][]string
			for _, a := range agents {
				rows = append(rows, []string{
					a.ID, a.ComputerName, a.OSType, a.AgentVersion,
					a.NetworkStatus, boolIcon(a.IsActive), a.SiteName,
				})
			}
			printTable([]string{"ID", "Name", "OS", "Version", "Network", "Active", "Site"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "agent"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}

func newAgentsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <agent-id>",
		Short: "Get agent details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			agent, err := c.AgentsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(agent)
			}
			rows := [][]string{
				{"ID", agent.ID},
				{"Name", agent.ComputerName},
				{"OS", fmt.Sprintf("%s %s (%s)", agent.OSName, agent.OSArch, agent.OSType)},
				{"Version", agent.AgentVersion},
				{"Network", agent.NetworkStatus},
				{"Active", boolIcon(agent.IsActive)},
				{"Infected", boolIcon(agent.Infected)},
				{"Site", fmt.Sprintf("%s (%s)", agent.SiteName, agent.SiteID)},
				{"Group", fmt.Sprintf("%s (%s)", agent.GroupName, agent.GroupID)},
				{"External IP", agent.ExternalIP},
				{"Last Active", agent.LastActiveDate},
				{"Registered", agent.RegisteredAt},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}

func newAgentsCountCmd() *cobra.Command {
	var siteIDs []string

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count agents",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			count, err := c.AgentsCount(cmd.Context(), &mgmt.AgentListParams{SiteIDs: siteIDs})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(map[string]int{"count": count})
			}
			fmt.Fprintln(cmd.OutOrStdout(), count)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}

func mgmtClient() (*mgmt.Client, error) {
	consoleURL, token, err := resolveConfig()
	if err != nil {
		return nil, err
	}
	return mgmt.NewClient(consoleURL, token), nil
}
