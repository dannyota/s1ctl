package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newApplicationsRisksCmd() *cobra.Command {
	var siteIDs, accountIDs, severities []string
	var vendor, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "risks",
		Short: "List application risks (CVE vulnerabilities per endpoint)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ApplicationRiskListParams{
				SiteIDs:           siteIDs,
				AccountIDs:        accountIDs,
				Severities:        severities,
				ApplicationVendor: vendor,
				Limit:             limit,
				Cursor:            cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var risks []mgmt.ApplicationRisk
			var total int

			if all {
				risks, total, err = fetchAllREST("risk", func(cur string) ([]mgmt.ApplicationRisk, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ApplicationRisksList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				risks, pag, err = c.ApplicationRisksList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"Application", "Version", "CVE", "Severity", "Risk Score", "Status", "Endpoint"}
			rows := make([][]string, len(risks))
			for i, r := range risks {
				rows[i] = []string{
					truncate(orDash(r.ApplicationName), 30),
					orDash(r.ApplicationVersion),
					orDash(r.CveID),
					orDash(r.Severity),
					orDash(r.RiskScore),
					orDash(r.MitigationStatus),
					truncate(orDash(r.EndpointName), 25),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, risks, len(risks), total, "risk", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (CRITICAL, HIGH, MEDIUM, LOW)")
	cmd.Flags().StringVar(&vendor, "vendor", "", "filter by vendor (contains)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return markJSON(cmd)
}
