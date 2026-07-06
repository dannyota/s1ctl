package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

// alertsScope builds a GraphQL scope selector from the --scope-level and
// --scope-id flags. Returns (nil, nil) when neither is set (whole-tenant).
func alertsScope(level, id string) (*graphql.Scope, error) {
	if level == "" && id == "" {
		return nil, nil
	}
	lv := strings.ToUpper(level)
	if lv == "" {
		lv = "SITE"
	}
	switch lv {
	case "ACCOUNT", "SITE", "GROUP":
	default:
		return nil, fmt.Errorf("invalid --scope-level %q (want account, site, or group)", level)
	}
	if id == "" {
		return nil, fmt.Errorf("--scope-id is required with --scope-level")
	}
	return &graphql.Scope{ScopeIDs: []string{id}, ScopeType: lv}, nil
}

func alertsFilters(severities, statuses, verdicts []string) []graphql.Filter {
	var filters []graphql.Filter
	if len(severities) > 0 {
		filters = append(filters, graphql.Filter{FieldID: "severity", StringIn: &graphql.InStr{Values: severities}})
	}
	if len(statuses) > 0 {
		filters = append(filters, graphql.Filter{FieldID: "status", StringIn: &graphql.InStr{Values: statuses}})
	}
	if len(verdicts) > 0 {
		filters = append(filters, graphql.Filter{FieldID: "analystVerdict", StringIn: &graphql.InStr{Values: verdicts}})
	}
	return filters
}

func newAlertsCountsCmd() *cobra.Command {
	var fields, severities, statuses []string
	var scopeLevel, scopeID string
	var groupBy bool

	cmd := &cobra.Command{
		Use:   "counts --field <fieldId> [--field ...]",
		Short: "Count alert values per field (filter counts or group-by counts)",
		Long: `Return per-field value counts for the current alert selection.

By default uses alertFiltersCount (distinct filterable values and their
cardinality). Pass --group-by to use the deprecated alertGroupByCount query
instead; for grouped alert volume prefer "alerts stats" (alertGroups).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if len(fields) == 0 {
				return fmt.Errorf("--field is required")
			}
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			filters := alertsFilters(severities, statuses, nil)

			var counts []graphql.AlertFieldCount
			if groupBy {
				counts, err = c.AlertsGroupByCount(cmd.Context(), fields, filters, scope)
			} else {
				counts, err = c.AlertsFiltersCount(cmd.Context(), fields, filters, scope)
			}
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), counts)
			}
			headers := []string{"Field", "Value", "Count"}
			var rows [][]string
			for _, f := range counts {
				for _, v := range f.Values {
					label := v.Value
					if v.Label != "" {
						label = v.Label
					}
					rows = append(rows, []string{f.FieldID, orDash(label), strconv.FormatInt(v.Count, 10)})
				}
			}
			printTable(headers, rows)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&fields, "field", nil, "field ID to count (repeatable, required)")
	cmd.Flags().BoolVar(&groupBy, "group-by", false, "use the deprecated alertGroupByCount query")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status (NEW, IN_PROGRESS, RESOLVED)")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (account, site, group)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "account, site, or group ID")
	return markJSON(cmd)
}

func newAlertsExportCmd() *cobra.Command {
	var severities, statuses, verdicts []string
	var scopeLevel, scopeID, view, outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export alerts to a CSV file",
		Long: `Export alerts matching the filters as CSV via alertsCsvExport.

The API returns the full CSV inline. It is written to --out, or to stdout when
--out is omitted.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			filters := alertsFilters(severities, statuses, verdicts)
			data, err := c.AlertsExport(cmd.Context(), filters, scope, graphql.ViewType(strings.ToUpper(view)))
			if err != nil {
				return err
			}
			if outFile == "" || outFile == "-" {
				_, err = cmd.OutOrStdout().Write([]byte(data))
				return err
			}
			if err := os.WriteFile(outFile, []byte(data), 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Exported alerts to %s\n", outFile)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status (NEW, IN_PROGRESS, RESOLVED)")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter by analyst verdict")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (account, site, group)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "account, site, or group ID")
	cmd.Flags().StringVar(&view, "view", "", "predefined view (ALL, CLOUD, ENDPOINT, IDENTITY, CUSTOM_ALERTS, THIRD_PARTY)")
	cmd.Flags().StringVar(&outFile, "out", "", "output file (default: stdout)")
	return cmd
}
