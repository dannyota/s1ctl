package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRemoteOpsResultsCmd() *cobra.Command {
	var (
		status []string
		cursor string
		limit  int
		all    bool
	)

	cmd := &cobra.Command{
		Use:   "results <parent-task-id>",
		Short: "Get remote script execution results",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.RemoteScriptsStatusParams{
				ParentTaskID: args[0],
				Status:       status,
				Limit:        limit,
				Cursor:       cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var tasks []mgmt.RemoteScriptTask
			var total int

			if all {
				tasks, total, err = fetchAllREST("task", func(cur string) ([]mgmt.RemoteScriptTask, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.RemoteScriptsStatus(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				tasks, pag, err = c.RemoteScriptsStatus(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Agent", "Status", "Detailed Status", "Updated"}
			rows := make([][]string, len(tasks))
			for i, t := range tasks {
				rows[i] = []string{
					t.ID, orDash(t.AgentComputerName), t.Status,
					orDash(t.DetailedStatus), orDash(t.UpdatedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, tasks, len(tasks), total, "task", all)
		},
	}
	cmd.Flags().StringSliceVar(&status, "status", nil, "filter by status (created, pending, in_progress, completed, failed, canceled, expired)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return markJSON(cmd)
}
