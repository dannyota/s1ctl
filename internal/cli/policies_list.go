package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newPoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policies",
		Short: "View endpoint policies",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newPoliciesListCmd())
	cmd.AddCommand(newPoliciesGetCmd())
	cmd.AddCommand(newPoliciesDiffCmd())
	cmd.AddCommand(newPoliciesRevertCmd())
	addPolicySyncCmds(cmd)
	return cmd
}

func newPoliciesGetCmd() *cobra.Command {
	var siteID, accountID, groupID string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get policy for a scope (site, account, or group)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			var raw json.RawMessage
			switch {
			case groupID != "" && siteID != "":
				p, pErr := c.PolicyGetGroup(cmd.Context(), siteID, groupID)
				if pErr != nil {
					return pErr
				}
				raw = p.Raw
			case siteID != "":
				p, pErr := c.PolicyGetSite(cmd.Context(), siteID)
				if pErr != nil {
					return pErr
				}
				raw = p.Raw
			case accountID != "":
				p, pErr := c.PolicyGetAccount(cmd.Context(), accountID)
				if pErr != nil {
					return pErr
				}
				raw = p.Raw
			default:
				return fmt.Errorf("specify --site-id, --account-id, or both --site-id and --group-id")
			}
			return printJSON(cmd.OutOrStdout(), raw)
		},
	}
	cmd.Flags().StringVar(&siteID, "site-id", "", "site ID")
	cmd.Flags().StringVar(&accountID, "account-id", "", "account ID")
	cmd.Flags().StringVar(&groupID, "group-id", "", "group ID (requires --site-id)")
	return cmd
}

// policyEntry combines scope info with a policy for list output.
type policyEntry struct {
	Scope    string      `json:"scope"`
	ScopeID  string      `json:"scopeId"`
	SiteName string      `json:"siteName"`
	Policy   mgmt.Policy `json:"policy"`
}

func newPoliciesListCmd() *cobra.Command {
	var accountIDs, siteIDs []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List policies across sites",
		Long: `List endpoint policies across all sites (or a filtered subset).

The SentinelOne API returns one policy per scope. This command fetches sites
and retrieves each site's policy, presenting them side by side for comparison.

Use --account-id or --site-id to narrow the scope.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			sites, err := fetchSitesForPolicies(cmd, c, accountIDs, siteIDs)
			if err != nil {
				return err
			}

			entries, errs := fetchPoliciesForSites(cmd, c, sites)
			for _, e := range errs {
				fmt.Fprintf(cmd.ErrOrStderr(), "Warning: %v\n", e)
			}

			headers := []string{
				"Site", "Site ID", "Mitigation", "Suspicious",
				"Anti-Tamper", "Quarantine", "Snapshots", "Remote Shell",
				"Inherited",
			}
			rows := make([][]string, len(entries))
			for i, e := range entries {
				rows[i] = []string{
					truncate(e.SiteName, 30),
					e.ScopeID,
					orDash(e.Policy.MitigationMode),
					orDash(e.Policy.MitigationModeSuspicious),
					boolIcon(e.Policy.AntiTamperingOn),
					boolIcon(e.Policy.NetworkQuarantineOn),
					boolIcon(e.Policy.SnapshotsOn),
					boolIcon(e.Policy.AllowRemoteShell),
					orDash(e.Policy.InheritedFrom),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, entries, len(entries), len(entries), "policy", true)
		},
	}
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}

// fetchSitesForPolicies returns sites filtered by the given account/site IDs.
// If no filters are given, all sites are returned.
func fetchSitesForPolicies(cmd *cobra.Command, c *mgmt.Client, accountIDs, siteIDs []string) ([]mgmt.Site, error) {
	if len(siteIDs) > 0 {
		var sites []mgmt.Site
		for _, id := range siteIDs {
			s, err := c.SitesGet(cmd.Context(), id)
			if err != nil {
				return nil, fmt.Errorf("get site %s: %w", id, err)
			}
			sites = append(sites, *s)
		}
		return sites, nil
	}

	params := &mgmt.SiteListParams{
		AccountIDs: accountIDs,
		Limit:      defaultPageSize,
	}
	sites, _, err := fetchAllREST("site", func(cur string) ([]mgmt.Site, *mgmt.Pagination, error) {
		params.Cursor = cur
		return c.SitesList(cmd.Context(), params)
	})
	return sites, err
}

// fetchPoliciesForSites fetches the policy for each site, collecting errors
// for sites that fail without aborting the entire operation.
func fetchPoliciesForSites(cmd *cobra.Command, c *mgmt.Client, sites []mgmt.Site) ([]policyEntry, []error) {
	var entries []policyEntry
	var errs []error

	for i, s := range sites {
		printProgress("policy", i, len(sites))
		p, err := c.PolicyGetSite(cmd.Context(), s.ID)
		if err != nil {
			errs = append(errs, fmt.Errorf("site %s (%s): %w", s.Name, s.ID, err))
			continue
		}
		entries = append(entries, policyEntry{
			Scope:    "site",
			ScopeID:  s.ID,
			SiteName: s.Name,
			Policy:   *p,
		})
	}
	clearProgress()
	return entries, errs
}
