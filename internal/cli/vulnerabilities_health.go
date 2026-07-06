package cli

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newVulnerabilitiesHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Summarize vulnerabilities by severity and status",
		Long: `Show a breakdown of vulnerability counts by severity and open/resolved status.
Uses count queries — no bulk data fetch needed.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}

			sevs := []string{"CRITICAL", "HIGH", "MEDIUM", "LOW"}
			type bucket struct {
				Severity string `json:"severity"`
				Open     int    `json:"open"`
				Resolved int    `json:"resolved"`
				Total    int    `json:"total"`
			}
			results := make([]bucket, len(sevs))
			var mu sync.Mutex
			var wg sync.WaitGroup
			var firstErr error

			ctx := cmd.Context()
			for i, sev := range sevs {
				wg.Add(1)
				go func(idx int, severity string) {
					defer wg.Done()
					total, e := countVulns(c, ctx, severity, "")
					if e != nil {
						mu.Lock()
						if firstErr == nil {
							firstErr = e
						}
						mu.Unlock()
						return
					}
					resolved, e := countVulns(c, ctx, severity, "RESOLVED")
					if e != nil {
						mu.Lock()
						if firstErr == nil {
							firstErr = e
						}
						mu.Unlock()
						return
					}
					mu.Lock()
					results[idx] = bucket{
						Severity: severity,
						Open:     total - resolved,
						Resolved: resolved,
						Total:    total,
					}
					mu.Unlock()
				}(i, sev)
			}
			wg.Wait()
			if firstErr != nil {
				return firstErr
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), results)
			}

			headers := []string{"Severity", "Open", "Resolved", "Total"}
			rows := make([][]string, len(results))
			var totalOpen, totalResolved, totalAll int
			for i, b := range results {
				rows[i] = []string{
					b.Severity,
					fmt.Sprintf("%d", b.Open),
					fmt.Sprintf("%d", b.Resolved),
					fmt.Sprintf("%d", b.Total),
				}
				totalOpen += b.Open
				totalResolved += b.Resolved
				totalAll += b.Total
			}
			printTable(headers, rows)

			fmt.Fprintf(cmd.OutOrStdout(), "\n%s: %d open, %d resolved\n",
				pluralize(totalAll, "vulnerability"), totalOpen, totalResolved)
			return nil
		},
	}
	return markJSON(cmd)
}

func countVulns(c *graphql.Client, ctx context.Context, severity, status string) (int, error) {
	params := &graphql.ListParams{First: 1}
	params.Filters = append(params.Filters, graphql.Filter{
		FieldID:  "severity",
		StringIn: &graphql.InStr{Values: []string{severity}},
	})
	if status != "" {
		params.Filters = append(params.Filters, graphql.Filter{
			FieldID:  "status",
			StringIn: &graphql.InStr{Values: []string{status}},
		})
	}
	conn, err := c.VulnerabilitiesList(ctx, params)
	if err != nil {
		return 0, err
	}
	return int(conn.TotalCount), nil
}
