package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRulesTrendsCmd() *cobra.Command {
	var siteIDs []string
	var top int

	cmd := &cobra.Command{
		Use:   "trends",
		Short: "Show noisiest rules by detection count",
		Long: `Fetch all custom detection rules and sort by generated alert count
(descending). Helps identify alert fatigue candidates for tuning.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.RuleListParams{SiteIDs: siteIDs, Limit: 1000}
			rules, _, err := fetchAllREST("rule", func(cur string) ([]mgmt.Rule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.RulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			sort.Slice(rules, func(i, j int) bool {
				return rules[i].GeneratedAlerts > rules[j].GeneratedAlerts
			})

			if top > 0 && top < len(rules) {
				rules = rules[:top]
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), rules)
			}

			headers := []string{"Name", "Alerts", "Status", "Severity", "Scope", "Response"}
			rows := make([][]string, len(rules))
			for i, r := range rules {
				response := "-"
				if r.TreatAsThreat != "" && r.TreatAsThreat != mgmt.RuleTreatUndefined {
					response = string(r.TreatAsThreat)
				}
				rows[i] = []string{
					truncate(r.Name, 40),
					fmt.Sprintf("%d", r.GeneratedAlerts),
					string(r.Status),
					string(r.Severity),
					string(r.Scope),
					response,
				}
			}
			printTable(headers, rows)

			var totalAlerts int
			for _, r := range rules {
				totalAlerts += r.GeneratedAlerts
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s, %d total alerts\n",
				pluralize(len(rules), "rule"), totalAlerts)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().IntVar(&top, "top", 0, "show only top N rules (default: all)")
	return cmd
}
