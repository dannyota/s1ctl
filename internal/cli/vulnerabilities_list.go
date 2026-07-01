package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newVulnerabilitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vulnerabilities",
		Aliases: []string{"vulns"},
		Short:   "Manage xSPM vulnerabilities",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newVulnerabilitiesListCmd())
	cmd.AddCommand(newVulnerabilitiesGetCmd())
	cmd.AddCommand(newVulnerabilitiesStatusCmd())
	cmd.AddCommand(newVulnerabilitiesVerdictCmd())
	cmd.AddCommand(newVulnerabilitiesHealthCmd())
	return cmd
}

func newVulnerabilitiesListCmd() *cobra.Command {
	var severities, statuses []string
	var after string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List vulnerabilities",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			params := &graphql.ListParams{First: limit, After: after}
			if params.First == 0 {
				params.First = defaultPageSize
			}
			if len(severities) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "severity",
					StringIn: &graphql.InStr{Values: severities},
				})
			}
			if len(statuses) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "status",
					StringIn: &graphql.InStr{Values: statuses},
				})
			}

			var items []graphql.Vulnerability
			var total int64

			if all {
				items, total, err = fetchAllGQL("vulnerability", func(cur string) (*graphql.Connection[graphql.Vulnerability], error) {
					params.After = cur
					return c.VulnerabilitiesList(cmd.Context(), params)
				})
			} else {
				conn, connErr := c.VulnerabilitiesList(cmd.Context(), params)
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

			headers := []string{"ID", "Name", "Severity", "Status", "CVE", "EPSS", "Asset", "Site"}
			rows := make([][]string, len(items))
			for i, v := range items {
				rows[i] = []string{
					v.ID, truncate(orDash(v.Name), 40), v.Severity,
					v.Status, orDash(v.CVE.ID),
					fmt.Sprintf("%.4f", v.CVE.EPSSScore),
					orDash(v.Asset.Name), orDash(v.Scope.Site.Name),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), int(total), "vulnerability", all)
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&after, "after", "", "pagination cursor")
	return cmd
}

func newVulnerabilitiesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get vulnerability details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			v, err := c.VulnerabilitiesGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), v)
			}
			rows := [][]string{
				{"ID", v.ID},
				{"Name", orDash(v.Name)},
				{"Severity", v.Severity},
				{"Status", v.Status},
				{"Verdict", orDash(v.AnalystVerdict)},
				{"CVE", orDash(v.CVE.ID)},
				{"EPSS", fmt.Sprintf("%.4f", v.CVE.EPSSScore)},
				{"NVD Score", fmt.Sprintf("%.1f", v.CVE.NVDBaseScore)},
				{"Exploit Maturity", orDash(v.CVE.ExploitMaturity)},
				{"Exploited in Wild", boolIcon(v.CVE.ExploitedInWild)},
				{"Software", orDash(v.Software.Name)},
				{"Version", orDash(v.Software.Version)},
				{"Fix Version", orDash(v.Software.FixVersion)},
				{"Product", orDash(v.Product)},
				{"Vendor", orDash(v.Vendor)},
				{"Asset", orDash(v.Asset.Name)},
				{"Site", orDash(v.Scope.Site.Name)},
				{"Detected", orDash(v.DetectedAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}

func newVulnerabilitiesStatusCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "status <id> <status>",
		Short: "Update vulnerability status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, status := args[0], args[1]
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would set status=%s on vulnerability %s. Pass --yes to apply.\n", status, id)
				return nil
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			if err := c.VulnerabilitiesUpdateStatus(cmd.Context(), []string{id}, status); err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "id": id})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "status: updated vulnerability %s\n", id)
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newVulnerabilitiesVerdictCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "verdict <id> <verdict>",
		Short: "Update vulnerability analyst verdict",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, verdict := args[0], args[1]
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would set verdict=%s on vulnerability %s. Pass --yes to apply.\n", verdict, id)
				return nil
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			if err := c.VulnerabilitiesUpdateVerdict(cmd.Context(), []string{id}, verdict); err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]string{"verdict": "updated", "id": id})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "verdict: updated vulnerability %s\n", id)
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
