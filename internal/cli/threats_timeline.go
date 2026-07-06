package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newThreatsTimelineCmd() *cobra.Command {
	var query, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "timeline <threat-id>",
		Short: "Show activity timeline for a threat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ThreatTimelineParams{
				Query:     query,
				Limit:     limit,
				Cursor:    cursor,
				SortBy:    sortBy,
				SortOrder: sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var entries []mgmt.ThreatTimelineEntry
			var total int

			if all {
				entries, total, err = fetchAllREST("event", func(cur string) ([]mgmt.ThreatTimelineEntry, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ThreatTimeline(cmd.Context(), args[0], params)
				})
			} else {
				var pag *mgmt.Pagination
				entries, pag, err = c.ThreatTimeline(cmd.Context(), args[0], params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"Timestamp", "Type", "Description", "Actor"}
			rows := make([][]string, len(entries))
			for i, e := range entries {
				rows[i] = []string{
					orDash(e.CreatedAt),
					fmt.Sprintf("%d", e.ActivityType),
					truncate(orDash(e.PrimaryDescription), 60),
					orDash(e.UserID),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, entries, len(entries), total, "event", all)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return markJSON(cmd)
}
