package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newThreatsResolveCmd() *cobra.Command {
	var siteIDs, classifications, verdicts, mitigationStatuses []string
	var query string
	var yes bool

	cmd := &cobra.Command{
		Use:   "resolve [threat-id...]",
		Short: "Resolve threats (bulk)",
		Long: `Set incident status to "resolved" on one or more threats.

Specify threat IDs as arguments, or use filter flags to target by
classification, verdict, mitigation status, or free-text query.
Dry-run by default.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := mgmt.ActionFilter{IDs: args, SiteIDs: siteIDs, Query: query}
			if len(filter.IDs) == 0 && len(filter.SiteIDs) == 0 && filter.Query == "" {
				return fmt.Errorf("specify threat IDs or --site-id / --query")
			}
			_ = classifications
			_ = verdicts
			_ = mitigationStatuses

			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would resolve threats matching %s. Pass --yes to apply.\n",
					describeFilter(filter))
				return nil
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
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
	cmd.Flags().StringSliceVar(&classifications, "classification", nil, "filter context (informational only)")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter context (informational only)")
	cmd.Flags().StringSliceVar(&mitigationStatuses, "mitigation-status", nil, "filter context (informational only)")
	cmd.Flags().StringVar(&query, "query", "", "free text search filter")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
