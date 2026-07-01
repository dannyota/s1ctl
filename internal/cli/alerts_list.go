package cli

import (
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
	var severities, verdicts []string
	var after string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List alerts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			params := &graphql.ListParams{First: limit, After: after}
			if params.First == 0 {
				params.First = defaultPageSize
			}
			if len(severities) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "severity",
					StringIn: &graphql.InStr{Values: severities},
				})
			}
			if len(verdicts) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "analystVerdict",
					StringIn: &graphql.InStr{Values: verdicts},
				})
			}

			var alerts []graphql.Alert
			var total int64

			if all {
				alerts, total, err = fetchAllGQL("alert", func(cur string) (*graphql.Connection[graphql.Alert], error) {
					params.After = cur
					return c.AlertsList(cmd.Context(), params)
				})
			} else {
				conn, connErr := c.AlertsList(cmd.Context(), params)
				if connErr != nil {
					return connErr
				}
				total = conn.TotalCount
				for _, edge := range conn.Edges {
					alerts = append(alerts, edge.Node)
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Severity", "Status", "Verdict", "Detected"}
			rows := make([][]string, len(alerts))
			for i, a := range alerts {
				rows[i] = []string{
					a.ID, truncate(orDash(a.Name), 40), a.Severity,
					a.Status, a.AnalystVerdict, orDash(a.DetectedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, alerts, len(alerts), int(total), "alert", all)
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter by analyst verdict")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&after, "after", "", "pagination cursor")
	return cmd
}

func gqlClient() (*graphql.Client, error) {
	consoleURL, token, err := resolveConfig()
	if err != nil {
		return nil, err
	}
	return graphql.NewClient(consoleURL, token), nil
}
