package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newAlertsHistoryCmd() *cobra.Command {
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "history <alert-id>",
		Short: "Show audit trail for an alert",
		Args:  cobra.ExactArgs(1),
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

			var items []graphql.AlertHistoryItem
			var total int64

			if all {
				items, total, err = fetchAllGQL("history item", func(cur string) (*graphql.Connection[graphql.AlertHistoryItem], error) {
					return c.AlertHistory(cmd.Context(), alertID, pageSize, cur)
				})
			} else {
				conn, connErr := c.AlertHistory(cmd.Context(), alertID, pageSize, "")
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

			headers := []string{"Timestamp", "Action", "Text", "Actor"}
			rows := make([][]string, len(items))
			for i, h := range items {
				rows[i] = []string{
					orDash(h.CreatedAt),
					orDash(h.EventType),
					truncate(orDash(h.EventText), 60),
					orDash(h.ActorName()),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), int(total), "history item", all)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	return cmd
}
