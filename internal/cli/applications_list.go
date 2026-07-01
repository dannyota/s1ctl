package cli

import (
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
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed applications",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ApplicationListParams{
				AgentIDs: agentIDs,
				SiteIDs:  siteIDs,
				Query:    query,
				Limit:    limit,
				Cursor:   cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var apps []mgmt.Application
			var total int

			if all {
				apps, total, err = fetchAllREST("application", func(cur string) ([]mgmt.Application, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ApplicationsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				apps, pag, err = c.ApplicationsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Version", "Publisher", "OS"}
			rows := make([][]string, len(apps))
			for i, a := range apps {
				rows[i] = []string{
					a.ID, truncate(a.Name, 40), a.Version,
					a.Publisher, a.OSType,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, apps, len(apps), total, "application", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&agentIDs, "agent-id", nil, "filter by agent ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}
