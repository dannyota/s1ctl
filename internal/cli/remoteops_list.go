package cli

import (
	"fmt"
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
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List remote scripts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			scripts, pag, err := c.RemoteScriptsList(cmd.Context(), &mgmt.RemoteScriptListParams{
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(scripts)
			}
			var rows [][]string
			for _, s := range scripts {
				rows = append(rows, []string{
					s.ID, s.FileName, s.ScriptType,
					strings.Join(s.OSTypes, ","), s.CreatorName,
				})
			}
			printTable([]string{"ID", "File", "Type", "OS", "Creator"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "script"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}
