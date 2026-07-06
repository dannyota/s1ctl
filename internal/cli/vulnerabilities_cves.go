package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newVulnerabilitiesCvesCmd() *cobra.Command {
	var after string
	var limit int
	var all bool
	var minCVSS float64

	cmd := &cobra.Command{
		Use:   "cves",
		Short: "List CVEs",
		Long: `List CVEs via the cves query.

The cves server-side filter (CveFilterInput) supports only datetime-range
filtering, so --min-cvss is applied client-side against each CVE's NVD base
score after fetching. It only sees the fetched page, so pair it with --all to
filter the full result set rather than a single page.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			first := limit
			if first == 0 {
				first = defaultPageSize
			}

			var items []graphql.Cve
			var total int64
			if all {
				items, total, err = fetchAllGQL("cve", func(cur string) (*graphql.Connection[graphql.Cve], error) {
					return c.CvesList(cmd.Context(), nil, nil, first, cur)
				})
			} else {
				conn, connErr := c.CvesList(cmd.Context(), nil, nil, first, after)
				if connErr != nil {
					return connErr
				}
				total = conn.TotalCount
				for _, edge := range conn.Edges {
					items = append(items, edge.Node)
				}
			}
			if err != nil {
				return err
			}

			if minCVSS > 0 {
				filtered := items[:0:0]
				for _, cve := range items {
					if cve.NVDBaseScore >= minCVSS {
						filtered = append(filtered, cve)
					}
				}
				items = filtered
			}

			headers := []string{"CVE", "NVD", "Risk", "EPSS", "Exploited", "Published"}
			rows := make([][]string, len(items))
			for i, cve := range items {
				rows[i] = []string{
					cve.ID,
					fmt.Sprintf("%.1f", cve.NVDBaseScore),
					fmt.Sprintf("%.1f", cve.RiskScore),
					fmt.Sprintf("%.4f", cve.EPSSScore),
					boolIcon(cve.ExploitedInTheWild),
					orDash(cve.PublishedDate),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), int(total), "CVE", all || minCVSS > 0)
		},
	}
	cmd.Flags().Float64Var(&minCVSS, "min-cvss", 0, "only show CVEs with NVD base score >= this value (client-side)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&after, "after", "", "pagination cursor")
	return markJSON(cmd)
}

func newVulnerabilitiesCveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cve <id>",
		Short: "Get CVE details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			cve, err := c.CveGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), cve)
			}
			rows := [][]string{
				{"CVE", cve.ID},
				{"NVD Base Score", fmt.Sprintf("%.1f", cve.NVDBaseScore)},
				{"Risk Score", fmt.Sprintf("%.1f", cve.RiskScore)},
				{"EPSS Score", fmt.Sprintf("%.4f", cve.EPSSScore)},
				{"EPSS Percentile", fmt.Sprintf("%.4f", cve.EPSSPercentile)},
				{"Exploit Maturity", orDash(cve.ExploitMaturity)},
				{"Exploited in Wild", boolIcon(cve.ExploitedInTheWild)},
				{"KEV Available", boolIcon(cve.KevAvailable)},
				{"Remediation Level", orDash(cve.RemediationLevel)},
				{"Published", orDash(cve.PublishedDate)},
				{"NVD Reference", orDash(cve.NVDReferenceURL)},
				{"Description", truncate(orDash(cve.Description), 80)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}

func newVulnerabilitiesStatsCmd() *cobra.Command {
	var severities []string
	var scopeLevel, scopeID, top string
	var limit int

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Summarize vulnerability posture (unique CVEs + top vulnerable applications/assets/OS)",
		Long: `Summarize vulnerability posture.

By default reports the unique CVE count plus the top vulnerable applications,
assets, and OS types. Pass --top applications|assets|os to show only one list.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch top {
			case "", "applications", "assets", "os":
			default:
				return fmt.Errorf("invalid --top %q (want applications, assets, or os)", top)
			}
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			if limit == 0 {
				limit = 10
			}
			filters := alertsFilters(severities, nil, nil)
			ctx := cmd.Context()

			switch top {
			case "applications":
				apps, e := c.TopVulnerableApplications(ctx, filters, scope, limit)
				if e != nil {
					return e
				}
				return printApplicationStats(cmd, apps)
			case "assets":
				assets, e := c.TopVulnerableAssets(ctx, filters, scope, limit)
				if e != nil {
					return e
				}
				return printAssetStats(cmd, assets)
			case "os":
				osTypes, e := c.TopVulnerableOsTypes(ctx, filters, scope, limit)
				if e != nil {
					return e
				}
				return printOsTypeStats(cmd, osTypes)
			}

			count, err := c.UniqueCveCount(ctx, filters, scope)
			if err != nil {
				return err
			}
			apps, err := c.TopVulnerableApplications(ctx, filters, scope, limit)
			if err != nil {
				return err
			}
			assets, err := c.TopVulnerableAssets(ctx, filters, scope, limit)
			if err != nil {
				return err
			}
			osTypes, err := c.TopVulnerableOsTypes(ctx, filters, scope, limit)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]any{
					"uniqueCveCount":            count,
					"topVulnerableApplications": apps,
					"topVulnerableAssets":       assets,
					"topVulnerableOsTypes":      osTypes,
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Unique CVEs: %d\n\nTop vulnerable applications:\n", count)
			if err := printApplicationStats(cmd, apps); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "\nTop vulnerable assets:")
			if err := printAssetStats(cmd, assets); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "\nTop vulnerable OS types:")
			return printOsTypeStats(cmd, osTypes)
		},
	}
	cmd.Flags().StringVar(&top, "top", "", "show only one list: applications, assets, or os")
	cmd.Flags().IntVar(&limit, "limit", 0, "number of top entries per list (default 10)")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (account, site, group)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "account, site, or group ID")
	return markJSON(cmd)
}

func printApplicationStats(cmd *cobra.Command, apps []graphql.ApplicationStats) error {
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), apps)
	}
	headers := []string{"Application", "Version", "Assets", "CVEs", "Highest Risk"}
	rows := make([][]string, len(apps))
	for i, a := range apps {
		rows[i] = []string{
			orDash(a.Name), orDash(a.Version),
			fmt.Sprintf("%d", a.AssetCount), fmt.Sprintf("%d", a.CveCount),
			fmt.Sprintf("%.1f", a.HighestRiskScore),
		}
	}
	printTable(headers, rows)
	return nil
}

func printAssetStats(cmd *cobra.Command, assets []graphql.AssetStats) error {
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), assets)
	}
	headers := []string{"Asset", "Scope", "CVEs", "Highest Risk"}
	rows := make([][]string, len(assets))
	for i, a := range assets {
		rows[i] = []string{
			orDash(a.Name), orDash(a.ScopeName),
			fmt.Sprintf("%d", a.CveCount), fmt.Sprintf("%.1f", a.HighestRiskScore),
		}
	}
	printTable(headers, rows)
	return nil
}

func printOsTypeStats(cmd *cobra.Command, osTypes []graphql.OsTypeStats) error {
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), osTypes)
	}
	headers := []string{"OS Type", "Version", "Assets", "CVEs", "Avg Risk"}
	rows := make([][]string, len(osTypes))
	for i, o := range osTypes {
		rows[i] = []string{
			orDash(o.Name), orDash(o.Version),
			fmt.Sprintf("%d", o.AssetCount), fmt.Sprintf("%d", o.CveCount),
			fmt.Sprintf("%.1f", o.AverageRiskScore),
		}
	}
	printTable(headers, rows)
	return nil
}
