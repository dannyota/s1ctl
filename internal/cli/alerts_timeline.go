package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newAlertsTimelineCmd() *cobra.Command {
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "timeline <alert-id>",
		Short: "Show the timeline for an alert",
		Long: `Show the alert timeline: notes, activities, enrichments, indicators,
asset operations, mitigation actions, and related alerts, newest first.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alertID := args[0]
			c, err := gqlClient()
			if err != nil {
				return err
			}

			pageSize := limit
			if pageSize == 0 {
				pageSize = defaultPageSize
			}

			var items []graphql.AlertTimelineEntry
			var total int64

			if all {
				items, total, err = fetchAllGQL("timeline item", func(cur string) (*graphql.Connection[graphql.AlertTimelineEntry], error) {
					return c.AlertTimeline(cmd.Context(), alertID, pageSize, cur)
				})
			} else {
				conn, connErr := c.AlertTimeline(cmd.Context(), alertID, pageSize, "")
				if connErr != nil {
					return connErr
				}
				total = conn.TotalCount
				for _, edge := range conn.Edges {
					items = append(items, edge.Node)
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"Timestamp", "Event", "Text", "Actor"}
			rows := make([][]string, len(items))
			for i, e := range items {
				rows[i] = []string{
					orDash(e.CreatedAt),
					orDash(e.EventType),
					truncate(orDash(e.EventText), 60),
					orDash(e.ActorName()),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), int(total), "timeline item", all)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	return markJSON(cmd)
}
