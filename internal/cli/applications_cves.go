package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newApplicationsCVEsCmd() *cobra.Command {
	var siteIDs, accountIDs, severities []string
	var appName, vendor, cveID, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "cves",
		Short: "List CVEs across applications",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ApplicationCVEListParams{
				SiteIDs:           siteIDs,
				AccountIDs:        accountIDs,
				Severities:        severities,
				ApplicationName:   appName,
				ApplicationVendor: vendor,
				CveID:             cveID,
				Limit:             limit,
				Cursor:            cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var cves []mgmt.ApplicationCVE
			var total int

			if all {
				cves, total, err = fetchAllREST("CVE", func(cur string) ([]mgmt.ApplicationCVE, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ApplicationCVEsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				cves, pag, err = c.ApplicationCVEsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"CVE ID", "Severity", "Base Score", "Risk Score", "Exploited", "Published"}
			rows := make([][]string, len(cves))
			for i, c := range cves {
				rows[i] = []string{
					orDash(c.CveID),
					orDash(c.Severity),
					orDash(c.NvdBaseScore),
					orDash(c.RiskScore),
					orDash(c.ExploitedInTheWild),
					orDash(c.PublishedDate),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, cves, len(cves), total, "CVE", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (CRITICAL, HIGH, MEDIUM, LOW)")
	cmd.Flags().StringVar(&appName, "app-name", "", "filter by application name")
	cmd.Flags().StringVar(&vendor, "vendor", "", "filter by vendor")
	cmd.Flags().StringVar(&cveID, "cve-id", "", "filter by CVE ID (contains)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}
