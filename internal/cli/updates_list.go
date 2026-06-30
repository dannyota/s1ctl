package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newUpdatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updates",
		Short: "Agent update packages",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newUpdatesListCmd())
	return cmd
}

func newUpdatesListCmd() *cobra.Command {
	var siteIDs []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List update packages",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			pkgs, pag, err := c.UpdatesList(cmd.Context(), &mgmt.UpdateListParams{
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(pkgs)
			}
			var rows [][]string
			for _, p := range pkgs {
				rows = append(rows, []string{
					p.ID, p.FileName, p.Version,
					p.OSType, p.Status,
				})
			}
			printTable([]string{"ID", "File", "Version", "OS", "Status"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "package"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}
