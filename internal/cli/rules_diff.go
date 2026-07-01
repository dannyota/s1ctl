package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRulesDiffCmd() *cobra.Command {
	var siteIDs []string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Show rules with hit counts and status breakdown",
		Long: `Fetch all custom detection rules and show which have fired
(generatedAlerts > 0) vs dormant (zero alerts). Helps identify
rules worth tuning or disabling.`,
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

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), rules)
			}

			headers := []string{"Name", "Status", "Severity", "Alerts", "Scope", "Response"}
			rows := make([][]string, len(rules))
			for i, r := range rules {
				response := "-"
				if r.TreatAsThreat != "" && r.TreatAsThreat != mgmt.RuleTreatUndefined {
					response = string(r.TreatAsThreat)
				}
				rows[i] = []string{
					truncate(r.Name, 40),
					string(r.Status),
					string(r.Severity),
					fmt.Sprintf("%d", r.GeneratedAlerts),
					string(r.Scope),
					response,
				}
			}
			printTable(headers, rows)

			var active, disabled, fired, dormant int
			for _, r := range rules {
				switch r.Status {
				case mgmt.RuleStatusActive:
					active++
				case mgmt.RuleStatusDisabled:
					disabled++
				}
				if r.GeneratedAlerts > 0 {
					fired++
				} else {
					dormant++
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s total: %d active, %d disabled, %d fired, %d dormant\n",
				pluralize(len(rules), "rule"), active, disabled, fired, dormant)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}
