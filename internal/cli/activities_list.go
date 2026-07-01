package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newActivitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activities",
		Short: "View activity log",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newActivitiesListCmd())
	return cmd
}

func newActivitiesListCmd() *cobra.Command {
	var siteIDs []string
	var cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List activities",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ActivityListParams{
				SiteIDs: siteIDs,
				Limit:   limit,
				Cursor:  cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var activities []mgmt.Activity
			var total int

			if all {
				activities, total, err = fetchAllREST("activity", func(cur string) ([]mgmt.Activity, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ActivitiesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				activities, pag, err = c.ActivitiesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Description", "Site", "Created"}
			rows := make([][]string, len(activities))
			for i, a := range activities {
				rows[i] = []string{
					a.ID, truncate(a.PrimaryDesc, 60), a.SiteName, a.CreatedAt,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, activities, len(activities), total, "activity", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}
