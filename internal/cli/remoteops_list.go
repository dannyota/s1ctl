package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRemoteOpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remoteops",
		Short: "Remote operations and scripts",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newRemoteOpsListCmd())
	return cmd
}

func newRemoteOpsListCmd() *cobra.Command {
	var siteIDs []string
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List remote scripts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.RemoteScriptListParams{
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
				Cursor:  cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var scripts []mgmt.RemoteScript
			var total int

			if all {
				scripts, total, err = fetchAllREST("script", func(cur string) ([]mgmt.RemoteScript, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.RemoteScriptsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				scripts, pag, err = c.RemoteScriptsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "File", "Type", "OS", "Creator"}
			rows := make([][]string, len(scripts))
			for i, s := range scripts {
				rows[i] = []string{
					s.ID, s.FileName, s.ScriptType,
					strings.Join(s.OSTypes, ","), s.CreatorName,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, scripts, len(scripts), total, "script", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}
