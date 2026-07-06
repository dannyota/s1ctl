package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newSitesLicensesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "licenses",
		Short: "Show license utilization across sites",
		Long: `Aggregate license health view across all sites.
Each site shows active vs total licenses, utilization percentage,
expiration date, and a status indicator (OK, WARNING, CRITICAL).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.SiteListParams{Limit: 1000}
			sites, _, err := fetchAllREST("site", func(cur string) ([]mgmt.Site, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.SitesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			now := time.Now()
			type licenseEntry struct {
				ID          string  `json:"id"`
				Name        string  `json:"name"`
				Active      int     `json:"active"`
				Total       int     `json:"total"`
				Unlimited   bool    `json:"unlimited"`
				Utilization float64 `json:"utilization"`
				Expiration  string  `json:"expiration"`
				Status      string  `json:"status"`
			}

			var entries []licenseEntry
			var totalActive, totalCapacity int
			var hasUnlimited bool

			for _, s := range sites {
				var util float64
				if s.UnlimitedLicenses {
					hasUnlimited = true
				} else if s.TotalLicenses > 0 {
					util = float64(s.ActiveLicenses) / float64(s.TotalLicenses) * 100
				}
				totalActive += s.ActiveLicenses
				if !s.UnlimitedLicenses {
					totalCapacity += s.TotalLicenses
				}

				status := licenseStatus(util, s.Expiration, s.UnlimitedLicenses, now)
				entries = append(entries, licenseEntry{
					ID:          s.ID,
					Name:        s.Name,
					Active:      s.ActiveLicenses,
					Total:       s.TotalLicenses,
					Unlimited:   s.UnlimitedLicenses,
					Utilization: util,
					Expiration:  s.Expiration,
					Status:      status,
				})
			}

			if outputFormat == "json" {
				type jsonOutput struct {
					Sites         []licenseEntry `json:"sites"`
					TotalActive   int            `json:"totalActive"`
					TotalCapacity int            `json:"totalCapacity"`
					Utilization   float64        `json:"utilization"`
				}
				var overallUtil float64
				if totalCapacity > 0 {
					overallUtil = float64(totalActive) / float64(totalCapacity) * 100
				}
				return printJSON(cmd.OutOrStdout(), jsonOutput{
					Sites:         entries,
					TotalActive:   totalActive,
					TotalCapacity: totalCapacity,
					Utilization:   overallUtil,
				})
			}

			headers := []string{"Name", "Active", "Total", "Utilization", "Expiration", "Status"}
			rows := make([][]string, len(entries))
			for i, e := range entries {
				totalStr := fmt.Sprintf("%d", e.Total)
				utilStr := fmt.Sprintf("%.1f%%", e.Utilization)
				if e.Unlimited {
					totalStr = "unlimited"
					utilStr = "-"
				}
				rows[i] = []string{
					truncate(e.Name, 30),
					fmt.Sprintf("%d", e.Active),
					totalStr,
					utilStr,
					formatExpiration(e.Expiration),
					e.Status,
				}
			}
			printTable(headers, rows)

			var overallUtil float64
			if totalCapacity > 0 {
				overallUtil = float64(totalActive) / float64(totalCapacity) * 100
			}
			capacityStr := fmt.Sprintf("%d", totalCapacity)
			if hasUnlimited {
				capacityStr += "+ (some sites unlimited)"
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\n%d active licenses / %s capacity (%.1f%% utilization)\n",
				totalActive, capacityStr, overallUtil)
			return nil
		},
	}
	return markJSON(cmd)
}

func licenseStatus(utilization float64, expiration string, unlimited bool, now time.Time) string {
	daysUntilExpiry := daysUntil(expiration, now)

	if daysUntilExpiry >= 0 && daysUntilExpiry <= 7 {
		return "CRITICAL"
	}
	if !unlimited && utilization > 95 {
		return "CRITICAL"
	}
	if daysUntilExpiry >= 0 && daysUntilExpiry <= 30 {
		return "WARNING"
	}
	if !unlimited && utilization > 80 {
		return "WARNING"
	}
	return "OK"
}

func daysUntil(expiration string, now time.Time) int {
	if expiration == "" {
		return -1
	}
	t, err := time.Parse(time.RFC3339, expiration)
	if err != nil {
		return -1
	}
	return int(t.Sub(now).Hours() / 24)
}

func formatExpiration(s string) string {
	if s == "" {
		return "-"
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return s
	}
	return t.Format("2006-01-02")
}
