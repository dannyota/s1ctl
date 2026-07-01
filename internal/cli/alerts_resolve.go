package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newAlertsResolveCmd() *cobra.Command {
	var yes bool
	var name string
	var severities, sources []string

	cmd := &cobra.Command{
		Use:   "resolve [id...]",
		Short: "Resolve alerts by ID or filter",
		Long: `Set status to "RESOLVED" on one or more alerts.

Specify alert IDs directly, or use --name/--severity/--source to match alerts.
Filter flags only match alerts with status NEW. Dry-run by default.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			hasFilters := name != "" || len(severities) > 0 || len(sources) > 0
			if len(args) == 0 && !hasFilters {
				return fmt.Errorf("specify alert IDs or use --name/--severity/--source to match")
			}

			var ids []string

			if hasFilters {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				params := &graphql.ListParams{First: 1000}
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "status",
					StringIn: &graphql.InStr{Values: []string{"NEW"}},
				})
				if len(severities) > 0 {
					params.Filters = append(params.Filters, graphql.Filter{
						FieldID:  "severity",
						StringIn: &graphql.InStr{Values: severities},
					})
				}

				alerts, _, err := fetchAllGQL("alert", func(cur string) (*graphql.Connection[graphql.Alert], error) {
					params.After = cur
					return c.AlertsList(cmd.Context(), params)
				})
				if err != nil {
					return err
				}

				nameUpper := strings.ToUpper(name)
				var sourceSet map[string]bool
				if len(sources) > 0 {
					sourceSet = make(map[string]bool, len(sources))
					for _, s := range sources {
						sourceSet[strings.ToUpper(s)] = true
					}
				}

				for _, a := range alerts {
					if name != "" && !strings.Contains(strings.ToUpper(a.Name), nameUpper) {
						continue
					}
					if sourceSet != nil && !sourceSet[strings.ToUpper(a.DetectionSource.Product)] {
						continue
					}
					ids = append(ids, a.ID)
				}
			}

			ids = append(ids, args...)

			if len(ids) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No matching alerts found.")
				return nil
			}

			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would resolve %s. Pass --yes to apply.\n",
					pluralize(len(ids), "alert"))
				return nil
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			if err := c.AlertsUpdateStatus(cmd.Context(), ids, "RESOLVED"); err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]any{
					"status":   "resolved",
					"affected": len(ids),
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "resolve: %s affected\n", pluralize(len(ids), "alert"))
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	cmd.Flags().StringVar(&name, "name", "", "match alerts by name (contains, case-insensitive)")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL)")
	cmd.Flags().StringSliceVar(&sources, "source", nil, "filter by detection source (STAR, EDR, CWS)")
	return cmd
}
