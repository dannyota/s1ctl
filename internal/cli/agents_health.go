package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAgentsHealthCmd() *cobra.Command {
	var siteIDs []string

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Classify agents by operational state",
		Long: `Fetch all agents and classify them as active, offline (disconnected),
decommissioned, or infected. Helps identify endpoints that need attention.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.AgentListParams{SiteIDs: siteIDs, Limit: 1000}
			agents, _, err := fetchAllREST("agent", func(cur string) ([]mgmt.Agent, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.AgentsList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			type classified struct {
				Agent mgmt.Agent
				State string
			}
			var items []classified
			var active, offline, decommissioned, infected int

			for _, a := range agents {
				var state string
				switch {
				case a.Infected:
					state = "infected"
					infected++
				case a.IsDecommissioned:
					state = "decommissioned"
					decommissioned++
				case a.NetworkStatus == "disconnected" || !a.IsActive:
					state = "offline"
					offline++
				default:
					state = "active"
					active++
				}
				items = append(items, classified{Agent: a, State: state})
			}

			sort.Slice(items, func(i, j int) bool {
				order := map[string]int{"infected": 0, "offline": 1, "decommissioned": 2, "active": 3}
				oi := order[items[i].State]
				oj := order[items[j].State]
				if oi != oj {
					return oi < oj
				}
				return items[i].Agent.ComputerName < items[j].Agent.ComputerName
			})

			if outputFormat == "json" {
				type jsonItem struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					State   string `json:"state"`
					OS      string `json:"os"`
					Version string `json:"version"`
					Site    string `json:"site"`
				}
				out := make([]jsonItem, len(items))
				for i, it := range items {
					out[i] = jsonItem{
						ID:      it.Agent.ID,
						Name:    it.Agent.ComputerName,
						State:   it.State,
						OS:      it.Agent.OSType,
						Version: it.Agent.AgentVersion,
						Site:    it.Agent.SiteName,
					}
				}
				return printJSON(cmd.OutOrStdout(), out)
			}

			headers := []string{"Name", "State", "OS", "Version", "Network", "Site"}
			rows := make([][]string, len(items))
			for i, it := range items {
				rows[i] = []string{
					truncate(it.Agent.ComputerName, 30),
					it.State,
					it.Agent.OSType,
					it.Agent.AgentVersion,
					it.Agent.NetworkStatus,
					it.Agent.SiteName,
				}
			}
			printTable(headers, rows)

			fmt.Fprintf(cmd.OutOrStdout(), "\n%s: %d active, %d offline, %d decommissioned, %d infected\n",
				pluralize(len(agents), "agent"), active, offline, decommissioned, infected)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}
