package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newAlertsStatsCmd() *cobra.Command {
	var groupBy string
	var severities, statuses []string

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show alert volume grouped by field",
		Long: `Show alert counts grouped by a specified field using the GraphQL alertGroups query.

Common group-by fields: severity, status, analystVerdict, classification,
detectionSource.product, assets.name.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			params := &graphql.ListParams{First: 100}
			if len(severities) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "severity",
					StringIn: &graphql.InStr{Values: severities},
				})
			}
			if len(statuses) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "status",
					StringIn: &graphql.InStr{Values: statuses},
				})
			}
			conn, err := c.AlertGroups(cmd.Context(), groupBy, params)
			if err != nil {
				return err
			}
			var groups []graphql.AlertGroup
			for _, edge := range conn.Edges {
				groups = append(groups, edge.Node)
			}
			headers := []string{"Value", "Count"}
			rows := make([][]string, len(groups))
			for i, g := range groups {
				val := orDash(g.Value)
				if g.Label != "" && g.Label != g.Value {
					val = g.Label
				}
				rows[i] = []string{val, strconv.FormatInt(g.Count, 10)}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, groups, len(groups), len(groups), groupBy+" group", true)
		},
	}
	cmd.Flags().StringVar(&groupBy, "group-by", "severity", "field to group by (e.g. severity, status, analystVerdict)")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status (NEW, RESOLVED, etc.)")
	return markJSON(cmd)
}
