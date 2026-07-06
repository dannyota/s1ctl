package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRemoteOpsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remoteops",
		Short: "Manage remote operations and scripts",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newRemoteOpsListCmd())
	cmd.AddCommand(newRemoteOpsGetCmd())
	cmd.AddCommand(newRemoteOpsRunCmd())
	cmd.AddCommand(newRemoteOpsResultsCmd())
	cmd.AddCommand(newRemoteOpsUpdateCmd())
	cmd.AddCommand(newRemoteOpsContentCmd())
	cmd.AddCommand(newRemoteOpsUploadLimitsCmd())
	cmd.AddCommand(newRemoteOpsPendingCmd())
	cmd.AddCommand(newRemoteOpsGuardrailsCmd())
	return cmd
}

func newRemoteOpsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <script-id>",
		Short: "Get a remote script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			r, err := c.RemoteScriptsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), r)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"ID", r.ID},
				{"FileName", r.FileName},
				{"FileType", orDash(r.FileType)},
				{"ScriptType", orDash(r.ScriptType)},
				{"OSTypes", orDash(strings.Join(r.OSTypes, ", "))},
				{"ScopeLevel", orDash(r.ScopeLevel)},
				{"ScopeID", orDash(r.ScopeID)},
				{"CreatedAt", orDash(r.CreatedAt)},
			})
			return nil
		},
	}
	return markJSON(cmd)
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
	return markJSON(cmd)
}
