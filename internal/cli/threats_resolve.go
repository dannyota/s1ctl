package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newThreatsResolveCmd() *cobra.Command {
	var siteIDs, classifications, verdicts, mitigationStatuses []string
	var name, query string
	var yes bool

	cmd := &cobra.Command{
		Use:   "resolve [threat-id...]",
		Short: "Resolve threats (bulk)",
		Long: `Set incident status to "resolved" on one or more threats.

Specify threat IDs as arguments, or use filter flags to match threats.
Filter flags only match unresolved threats. Dry-run by default.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			hasFilters := name != "" || len(classifications) > 0 || len(verdicts) > 0 ||
				len(mitigationStatuses) > 0 || len(siteIDs) > 0 || query != ""
			if len(args) == 0 && !hasFilters {
				return fmt.Errorf("specify threat IDs or use filter flags")
			}

			var ids []string

			if hasFilters {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				params := &mgmt.ThreatListParams{
					SiteIDs:            siteIDs,
					Classifications:    classifications,
					AnalystVerdicts:    verdicts,
					MitigationStatuses: mitigationStatuses,
					IncidentStatuses:   []string{"unresolved"},
					Query:              query,
					Limit:              1000,
				}
				threats, _, err := fetchAllREST("threat", func(cur string) ([]mgmt.Threat, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ThreatsList(cmd.Context(), params)
				})
				if err != nil {
					return err
				}
				nameUpper := strings.ToUpper(name)
				for _, t := range threats {
					if name != "" && !strings.Contains(strings.ToUpper(t.ThreatName), nameUpper) {
						continue
					}
					ids = append(ids, t.ID)
				}
			}

			ids = append(ids, args...)

			if len(ids) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No matching threats found.")
				return nil
			}

			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would resolve %s. Pass --yes to apply.\n",
					pluralize(len(ids), "threat"))
				return nil
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			filter := mgmt.ActionFilter{IDs: ids}
			affected, err := c.ThreatsUpdateStatus(cmd.Context(), "resolved", filter)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "resolve: %s affected\n", pluralize(affected, "threat"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&name, "name", "", "match threats by name (contains, case-insensitive)")
	cmd.Flags().StringSliceVar(&classifications, "classification", nil, "filter by classification (e.g. Malware, PUP)")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter by analyst verdict")
	cmd.Flags().StringSliceVar(&mitigationStatuses, "mitigation-status", nil, "filter by mitigation status")
	cmd.Flags().StringVar(&query, "query", "", "free text search filter")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
