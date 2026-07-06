package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newAlertsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alerts",
		Short: "Manage unified alerts (GraphQL UAM)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAlertsListCmd())
	cmd.AddCommand(newAlertsGetCmd())
	cmd.AddCommand(newAlertsCountCmd())
	cmd.AddCommand(newAlertsResolveCmd())
	cmd.AddCommand(newAlertsStatusCmd())
	cmd.AddCommand(newAlertsVerdictCmd())
	cmd.AddCommand(newAlertsAddNoteCmd())
	cmd.AddCommand(newAlertsNotesCmd())
	cmd.AddCommand(newAlertsNoteUpdateCmd())
	cmd.AddCommand(newAlertsNoteDeleteCmd())
	cmd.AddCommand(newAlertsTimelineCmd())
	cmd.AddCommand(newAlertsCountsCmd())
	cmd.AddCommand(newAlertsExportCmd())
	cmd.AddCommand(newAlertsStatsCmd())
	cmd.AddCommand(newAlertsHistoryCmd())
	return cmd
}

func newAlertsListCmd() *cobra.Command {
	var severities, statuses, verdicts, sources []string
	var after, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List alerts",
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
			if len(verdicts) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "analystVerdict",
					StringIn: &graphql.InStr{Values: verdicts},
				})
			}
			if sortBy != "" {
				order := "DESC"
				if sortOrder != "" {
					order = strings.ToUpper(sortOrder)
				}
				params.Sort = &graphql.SortInput{By: sortBy, Order: order}
			}

			var alerts []graphql.Alert
			var total int64

			if all {
				alerts, total, err = fetchAllGQL("alert", func(cur string) (*graphql.Connection[graphql.Alert], error) {
					params.After = cur
					return c.AlertsList(cmd.Context(), params)
				})
			} else {
				conn, connErr := c.AlertsList(cmd.Context(), params)
				if connErr != nil {
					return connErr
				}
				total = conn.TotalCount
				for _, edge := range conn.Edges {
					alerts = append(alerts, edge.Node)
				}
			}
			if err != nil {
				return err
			}

			if len(sources) > 0 {
				sourceSet := make(map[string]bool, len(sources))
				for _, s := range sources {
					sourceSet[strings.ToUpper(s)] = true
				}
				filtered := alerts[:0]
				for _, a := range alerts {
					if sourceSet[strings.ToUpper(a.DetectionSource.Product)] {
						filtered = append(filtered, a)
					}
				}
				alerts = filtered
			}

			headers := []string{"ID", "Name", "Agent", "Severity", "Source", "Status", "Detected"}
			rows := make([][]string, len(alerts))
			for i, a := range alerts {
				rows[i] = []string{
					a.ID, truncate(orDash(a.Name), 35),
					truncate(orDash(a.AgentName()), 20),
					a.Severity, orDash(a.DetectionSource.Product),
					a.Status, orDash(a.DetectedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, alerts, len(alerts), int(total), "alert", all)
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status (NEW, IN_PROGRESS, RESOLVED)")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter by analyst verdict (e.g. FALSE_POSITIVE_BENIGN, TRUE_POSITIVE_MALWARE; see 'enums' cmd)")
	cmd.Flags().StringSliceVar(&sources, "source", nil, "filter by detection source (STAR, EDR, CWS)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&after, "after", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. detectedAt, severity)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (ASC, DESC)")
	return markJSON(cmd)
}

func newAlertsCountCmd() *cobra.Command {
	var severities, statuses, verdicts []string

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count alerts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			params := &graphql.ListParams{First: 1}
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
			if len(verdicts) > 0 {
				params.Filters = append(params.Filters, graphql.Filter{
					FieldID:  "analystVerdict",
					StringIn: &graphql.InStr{Values: verdicts},
				})
			}
			conn, err := c.AlertsList(cmd.Context(), params)
			if err != nil {
				return err
			}
			count := int(conn.TotalCount)
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]int{"count": count})
			}
			fmt.Fprintln(cmd.OutOrStdout(), count)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status (NEW, IN_PROGRESS, RESOLVED)")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter by analyst verdict (e.g. FALSE_POSITIVE_BENIGN, TRUE_POSITIVE_MALWARE; see 'enums' cmd)")
	return markJSON(cmd)
}

func newAlertsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get alert details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			a, err := c.AlertsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), a)
			}

			analyticsUID := "-"
			if a.Analytics != nil {
				analyticsUID = orDash(a.Analytics.UID)
			}

			rows := [][]string{
				{"ID", a.ID},
				{"Name", orDash(a.Name)},
				{"Description", orDash(a.Description)},
				{"Severity", a.Severity},
				{"Status", a.Status},
				{"Classification", orDash(a.Classification)},
				{"Confidence", orDash(a.ConfidenceLevel)},
				{"Verdict", orDash(a.AnalystVerdict)},
				{"Agent", orDash(a.AgentName())},
				{"Detected", orDash(a.DetectedAt)},
				{"Created", orDash(a.CreatedAt)},
				{"Updated", orDash(a.UpdatedAt)},
				{"Storyline ID", orDash(a.StorylineID)},
				{"Source", orDash(a.DetectionSource.Product)},
				{"Vendor", orDash(a.DetectionSource.Vendor)},
				{"Analytics UID", analyticsUID},
				{"Account", orDash(a.RealTime.Scope.Account.Name)},
				{"Site", orDash(a.RealTime.Scope.Site.Name)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}

func newAlertsStatusCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "status <id> <status>",
		Short: "Update alert status (NEW, IN_PROGRESS, RESOLVED)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, status := args[0], args[1]
			return guard(cmd.OutOrStdout(), "alerts status", fmt.Sprintf("set status=%s on alert %s", status, id), id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.AlertsUpdateStatus(cmd.Context(), []string{id}, status); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "status: updated alert %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newAlertsVerdictCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "verdict <id> <verdict>",
		Short: "Update alert analyst verdict (e.g. FALSE_POSITIVE_BENIGN, TRUE_POSITIVE_MALWARE)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, verdict := args[0], args[1]
			return guard(cmd.OutOrStdout(), "alerts verdict", fmt.Sprintf("set verdict=%s on alert %s", verdict, id), id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.AlertsUpdateVerdict(cmd.Context(), []string{id}, verdict); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"verdict": "updated", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "verdict: updated alert %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func gqlClient() (*graphql.Client, error) {
	consoleURL, token, err := resolveConfig()
	if err != nil {
		return nil, err
	}
	return graphql.NewClient(consoleURL, token), nil
}
