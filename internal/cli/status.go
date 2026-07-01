package cli

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
	"danny.vn/s1/mgmt"
)

func newStatusCmd() *cobra.Command {
	var siteIDs []string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show environment health summary",
		Long: `One-shot dashboard: agent count and health, unresolved threats,
NEW alerts, and site count.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			mc, err := mgmtClient()
			if err != nil {
				return err
			}
			gc, err := gqlClient()
			if err != nil {
				return err
			}

			type counts struct {
				agents         int
				activeAgents   int
				outdatedAgents int
				infectedAgents int
				threats        int
				unresolvedT    int
				sites          int
				groups         int
				alerts         int
				newAlerts      int
				criticalAlerts int
			}

			var c counts
			var mu sync.Mutex
			var wg sync.WaitGroup
			var firstErr error
			setErr := func(e error) {
				mu.Lock()
				if firstErr == nil {
					firstErr = e
				}
				mu.Unlock()
			}

			ctx := cmd.Context()
			agentParams := &mgmt.AgentListParams{SiteIDs: siteIDs}

			wg.Add(7)
			go func() {
				defer wg.Done()
				n, e := mc.AgentsCount(ctx, agentParams)
				if e != nil {
					setErr(e)
					return
				}
				mu.Lock()
				c.agents = n
				mu.Unlock()
			}()
			go func() {
				defer wg.Done()
				active := true
				p := &mgmt.AgentListParams{SiteIDs: siteIDs, IsActive: &active}
				n, e := mc.AgentsCount(ctx, p)
				if e != nil {
					setErr(e)
					return
				}
				mu.Lock()
				c.activeAgents = n
				mu.Unlock()
			}()
			go func() {
				defer wg.Done()
				upToDate := false
				p := &mgmt.AgentListParams{SiteIDs: siteIDs, IsUpToDate: &upToDate}
				n, e := mc.AgentsCount(ctx, p)
				if e != nil {
					setErr(e)
					return
				}
				mu.Lock()
				c.outdatedAgents = n
				mu.Unlock()
			}()
			go func() {
				defer wg.Done()
				infected := true
				p := &mgmt.AgentListParams{SiteIDs: siteIDs, Infected: &infected}
				n, e := mc.AgentsCount(ctx, p)
				if e != nil {
					setErr(e)
					return
				}
				mu.Lock()
				c.infectedAgents = n
				mu.Unlock()
			}()
			go func() {
				defer wg.Done()
				n, e := mc.ThreatsCount(ctx, &mgmt.ThreatListParams{SiteIDs: siteIDs})
				if e != nil {
					setErr(e)
					return
				}
				mu.Lock()
				c.threats = n
				mu.Unlock()

				nu, e2 := mc.ThreatsCount(ctx, &mgmt.ThreatListParams{
					SiteIDs:          siteIDs,
					IncidentStatuses: []string{"unresolved"},
				})
				if e2 != nil {
					setErr(e2)
					return
				}
				mu.Lock()
				c.unresolvedT = nu
				mu.Unlock()
			}()
			go func() {
				defer wg.Done()
				n, e := mc.SitesCount(ctx, &mgmt.SiteListParams{})
				if e != nil {
					setErr(e)
					return
				}
				mu.Lock()
				c.sites = n
				mu.Unlock()

				ng, e2 := mc.GroupsCount(ctx, &mgmt.GroupListParams{SiteIDs: siteIDs})
				if e2 != nil {
					setErr(e2)
					return
				}
				mu.Lock()
				c.groups = ng
				mu.Unlock()
			}()
			go func() {
				defer wg.Done()
				conn, e := gc.AlertsList(ctx, &graphql.ListParams{First: 1})
				if e != nil {
					setErr(e)
					return
				}
				mu.Lock()
				c.alerts = int(conn.TotalCount)
				mu.Unlock()

				connNew, e2 := gc.AlertsList(ctx, &graphql.ListParams{
					First: 1,
					Filters: []graphql.Filter{{
						FieldID:  "status",
						StringIn: &graphql.InStr{Values: []string{"NEW"}},
					}},
				})
				if e2 != nil {
					setErr(e2)
					return
				}
				mu.Lock()
				c.newAlerts = int(connNew.TotalCount)
				mu.Unlock()

				connCrit, e3 := gc.AlertsList(ctx, &graphql.ListParams{
					First: 1,
					Filters: []graphql.Filter{{
						FieldID:  "severity",
						StringIn: &graphql.InStr{Values: []string{"CRITICAL"}},
					}},
				})
				if e3 != nil {
					setErr(e3)
					return
				}
				mu.Lock()
				c.criticalAlerts = int(connCrit.TotalCount)
				mu.Unlock()
			}()
			wg.Wait()

			if firstErr != nil {
				return firstErr
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]int{
					"agents":          c.agents,
					"active_agents":   c.activeAgents,
					"outdated_agents": c.outdatedAgents,
					"infected_agents": c.infectedAgents,
					"threats":         c.threats,
					"unresolved":      c.unresolvedT,
					"sites":           c.sites,
					"groups":          c.groups,
					"alerts":          c.alerts,
					"new_alerts":      c.newAlerts,
					"critical_alerts": c.criticalAlerts,
				})
			}

			w := cmd.OutOrStdout()
			fmt.Fprintln(w, "Agents")
			fmt.Fprintf(w, "  Total:    %d\n", c.agents)
			fmt.Fprintf(w, "  Active:   %d\n", c.activeAgents)
			fmt.Fprintf(w, "  Outdated: %d\n", c.outdatedAgents)
			fmt.Fprintf(w, "  Infected: %d\n", c.infectedAgents)
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Threats")
			fmt.Fprintf(w, "  Total:      %d\n", c.threats)
			fmt.Fprintf(w, "  Unresolved: %d\n", c.unresolvedT)
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Alerts")
			fmt.Fprintf(w, "  Total:    %d\n", c.alerts)
			fmt.Fprintf(w, "  NEW:      %d\n", c.newAlerts)
			fmt.Fprintf(w, "  CRITICAL: %d\n", c.criticalAlerts)
			fmt.Fprintln(w)
			fmt.Fprintf(w, "Sites: %d  Groups: %d\n", c.sites, c.groups)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}
