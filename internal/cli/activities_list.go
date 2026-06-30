package cli

import (
	"fmt"

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
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List activities",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			activities, pag, err := c.ActivitiesList(cmd.Context(), &mgmt.ActivityListParams{
				SiteIDs: siteIDs,
				Limit:   limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(activities)
			}
			var rows [][]string
			for _, a := range activities {
				rows = append(rows, []string{
					a.ID, truncate(a.PrimaryDesc, 60), a.SiteName, a.CreatedAt,
				})
			}
			printTable([]string{"ID", "Description", "Site", "Created"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "activity"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}
