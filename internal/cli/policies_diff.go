package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newPoliciesDiffCmd() *cobra.Command {
	var accountIDs, siteIDs []string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Compare policies across sites",
		Long: `Fetch policies for all sites (or a filtered subset) and highlight
fields that differ between them. Useful for spotting inconsistencies
like one site in detect mode while others are in protect mode.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			sites, err := fetchSitesForPolicies(cmd, c, accountIDs, siteIDs)
			if err != nil {
				return err
			}
			if len(sites) < 2 {
				return fmt.Errorf("need at least 2 sites to diff (found %d)", len(sites))
			}

			entries, errs := fetchPoliciesForSites(cmd, c, sites)
			for _, e := range errs {
				fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %v\n", e)
			}
			if len(entries) < 2 {
				return fmt.Errorf("need at least 2 policies to diff (fetched %d)", len(entries))
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), entries)
			}

			type field struct {
				name string
				get  func(mgmt.Policy) string
			}
			fields := []field{
				{"mitigationMode", func(p mgmt.Policy) string { return p.MitigationMode }},
				{"mitigationModeSuspicious", func(p mgmt.Policy) string { return p.MitigationModeSuspicious }},
				{"antiTamperingOn", func(p mgmt.Policy) string { return fmt.Sprint(p.AntiTamperingOn) }},
				{"networkQuarantineOn", func(p mgmt.Policy) string { return fmt.Sprint(p.NetworkQuarantineOn) }},
				{"snapshotsOn", func(p mgmt.Policy) string { return fmt.Sprint(p.SnapshotsOn) }},
				{"allowRemoteShell", func(p mgmt.Policy) string { return fmt.Sprint(p.AllowRemoteShell) }},
				{"scanNewAgents", func(p mgmt.Policy) string { return fmt.Sprint(p.ScanNewAgents) }},
				{"autoDecommissionOn", func(p mgmt.Policy) string { return fmt.Sprint(p.AutoDecommissionOn) }},
				{"autoDecommissionDays", func(p mgmt.Policy) string { return fmt.Sprintf("%d", p.AutoDecommissionDays) }},
				{"ioc", func(p mgmt.Policy) string { return fmt.Sprint(p.Ioc) }},
			}

			w := cmd.OutOrStdout()
			var diffs int
			for _, f := range fields {
				vals := make(map[string][]string)
				for _, e := range entries {
					v := f.get(e.Policy)
					vals[v] = append(vals[v], e.SiteName)
				}
				if len(vals) <= 1 {
					continue
				}
				diffs++
				fmt.Fprintf(w, "%s:\n", f.name)
				for v, names := range vals {
					for _, name := range names {
						fmt.Fprintf(w, "  %-30s %s\n", name, v)
					}
				}
			}
			if diffs == 0 {
				fmt.Fprintln(w, "All policies are identical.")
			} else {
				fmt.Fprintf(w, "\n%s differ across %s\n",
					pluralize(diffs, "field"), pluralize(len(entries), "site"))
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}
