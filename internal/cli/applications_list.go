package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newApplicationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "applications",
		Aliases: []string{"apps"},
		Short:   "Application inventory",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newApplicationsListCmd())
	return cmd
}

func newApplicationsListCmd() *cobra.Command {
	var siteIDs, agentIDs []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed applications",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			apps, pag, err := c.ApplicationsList(cmd.Context(), &mgmt.ApplicationListParams{
				AgentIDs: agentIDs,
				SiteIDs:  siteIDs,
				Query:    query,
				Limit:    limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(apps)
			}
			var rows [][]string
			for _, a := range apps {
				rows = append(rows, []string{
					a.ID, truncate(a.Name, 40), a.Version,
					a.Publisher, a.OSType,
				})
			}
			printTable([]string{"ID", "Name", "Version", "Publisher", "OS"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "application"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&agentIDs, "agent-id", nil, "filter by agent ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}
