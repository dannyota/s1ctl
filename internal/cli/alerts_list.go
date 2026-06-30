package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newAlertsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alerts",
		Short: "Query unified alerts (GraphQL UAM)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAlertsListCmd())
	return cmd
}

func newAlertsListCmd() *cobra.Command {
	var severities, statuses []string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List alerts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			params := &graphql.AlertsListParams{First: limit}
			if len(severities) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "severity",
					StringIn: &graphql.InStr{Values: severities},
				})
			}
			if len(statuses) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "analystVerdict",
					StringIn: &graphql.InStr{Values: statuses},
				})
			}
			conn, err := c.AlertsList(cmd.Context(), params)
			if err != nil {
				return err
			}
			if jsonOutput {
				var alerts []graphql.Alert
				for _, edge := range conn.Edges {
					alerts = append(alerts, edge.Node)
				}
				return printJSON(alerts)
			}
			var rows [][]string
			for _, edge := range conn.Edges {
				a := edge.Node
				rows = append(rows, []string{
					a.ID, truncate(orDash(a.Name), 40), a.Severity,
					a.Status, a.AnalystVerdict, orDash(a.DetectedAt),
				})
			}
			printTable([]string{"ID", "Name", "Severity", "Status", "Verdict", "Detected"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(int(conn.TotalCount), "alert"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "verdict", nil, "filter by analyst verdict")
	cmd.Flags().IntVar(&limit, "limit", 25, "max results")
	return cmd
}

func gqlClient() (*graphql.Client, error) {
	consoleURL, token, err := resolveConfig()
	if err != nil {
		return nil, err
	}
	return graphql.NewClient(consoleURL, token), nil
}
