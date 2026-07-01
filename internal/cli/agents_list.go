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
	cmd.AddCommand(newAgentsUpgradeCmd())
	cmd.AddCommand(newAgentsOutdatedCmd())
	cmd.AddCommand(newAgentsVersionsCmd())
	addAgentActions(cmd)
	return cmd
}

func newAgentsListCmd() *cobra.Command {
	var siteIDs, groupIDs, osTypes, networkStatuses, machineTypes []string
	var query, cursor, sortBy, sortOrder string
	var limit int
	var all, infected, active bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List agents",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AgentListParams{
				SiteIDs:         siteIDs,
				GroupIDs:        groupIDs,
				OSTypes:         osTypes,
				NetworkStatuses: networkStatuses,
				MachineTypes:    machineTypes,
				Query:           query,
				Limit:           limit,
				Cursor:          cursor,
				SortBy:          sortBy,
				SortOrder:       sortOrder,
			}
			if cmd.Flags().Changed("infected") {
				params.Infected = &infected
			}
			if cmd.Flags().Changed("active") {
				params.IsActive = &active
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var agents []mgmt.Agent
			var total int

			if all {
				agents, total, err = fetchAllREST("agent", func(cur string) ([]mgmt.Agent, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.AgentsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				agents, pag, err = c.AgentsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "OS", "Version", "Network", "Active", "Site"}
			rows := make([][]string, len(agents))
			for i, a := range agents {
				rows[i] = []string{
					a.ID, a.ComputerName, a.OSType, a.AgentVersion,
					a.NetworkStatus, boolIcon(a.IsActive), a.SiteName,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, agents, len(agents), total, "agent", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type")
	cmd.Flags().StringSliceVar(&networkStatuses, "network-status", nil, "filter by network status (connected, disconnected)")
	cmd.Flags().StringSliceVar(&machineTypes, "machine-type", nil, "filter by machine type (server, desktop, laptop)")
	cmd.Flags().BoolVar(&infected, "infected", false, "filter by infection status")
	cmd.Flags().BoolVar(&active, "active", false, "filter by active status")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. computerName, lastActiveDate)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
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
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), agent)
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
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]int{"count": count})
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
