package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newUpdatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updates",
		Short: "Manage agent update packages",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newUpdatesListCmd())
	cmd.AddCommand(newUpdatesGetCmd())
	cmd.AddCommand(newUpdatesDeployCmd())
	return cmd
}

func newUpdatesGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <package-id>",
		Short: "Get an update package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			r, err := c.UpdatesGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), r)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"ID", r.ID},
				{"FileName", r.FileName},
				{"Version", orDash(r.Version)},
				{"OSType", orDash(r.OSType)},
				{"Status", orDash(r.Status)},
				{"FileSize", strconv.FormatInt(r.FileSize, 10)},
				{"ScopeName", orDash(r.ScopeName)},
			})
			return nil
		},
	}
	return markJSON(cmd)
}

func newUpdatesListCmd() *cobra.Command {
	var siteIDs []string
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List update packages",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.UpdateListParams{
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
				Cursor:  cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var pkgs []mgmt.UpdatePackage
			var total int

			if all {
				pkgs, total, err = fetchAllREST("update", func(cur string) ([]mgmt.UpdatePackage, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.UpdatesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				pkgs, pag, err = c.UpdatesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "File", "Version", "OS", "Status"}
			rows := make([][]string, len(pkgs))
			for i, p := range pkgs {
				rows[i] = []string{
					p.ID, p.FileName, p.Version,
					p.OSType, p.Status,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, pkgs, len(pkgs), total, "update", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return markJSON(cmd)
}
