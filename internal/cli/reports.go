package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newReportsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reports",
		Short: "Manage reports and report tasks",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newReportsListCmd())
	cmd.AddCommand(newReportTasksListCmd())
	cmd.AddCommand(newReportTypesCmd())
	cmd.AddCommand(newReportCreateCmd())
	cmd.AddCommand(newReportDownloadCmd())
	return cmd
}

func newReportsListCmd() *cobra.Command {
	var (
		siteIDs      []string
		scope        string
		frequency    string
		scheduleType string
		query        string
		cursor       string
		sortBy       string
		sortOrder    string
		limit        int
		all          bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List generated reports",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ReportListParams{
				SiteIDs:      siteIDs,
				Scope:        scope,
				Frequency:    frequency,
				ScheduleType: scheduleType,
				Query:        query,
				Limit:        limit,
				Cursor:       cursor,
				SortBy:       sortBy,
				SortOrder:    sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var reports []mgmt.Report
			var total int

			if all {
				reports, total, err = fetchAllREST("report", func(cur string) ([]mgmt.Report, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ReportsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				reports, pag, err = c.ReportsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Scope", "Type", "Status", "Created"}
			rows := make([][]string, len(reports))
			for i, r := range reports {
				rows[i] = []string{
					r.ID, truncate(r.Name, 40), r.Scope,
					r.ScheduleType, orDash(r.Status), orDash(r.CreatedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, reports, len(reports), total, "report", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&scope, "scope", "", "filter by scope (group, site, account, tenant)")
	cmd.Flags().StringVar(&frequency, "frequency", "", "filter by frequency (manually, weekly, monthly)")
	cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "filter by schedule type (manually, scheduled)")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. name, createdAt, status)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return cmd
}

func newReportTasksListCmd() *cobra.Command {
	var (
		siteIDs      []string
		scope        string
		frequency    string
		scheduleType string
		query        string
		cursor       string
		sortBy       string
		sortOrder    string
		limit        int
		all          bool
	)

	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "List report tasks and schedules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ReportTaskListParams{
				SiteIDs:      siteIDs,
				Scope:        scope,
				Frequency:    frequency,
				ScheduleType: scheduleType,
				Query:        query,
				Limit:        limit,
				Cursor:       cursor,
				SortBy:       sortBy,
				SortOrder:    sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var tasks []mgmt.ReportTask
			var total int

			if all {
				tasks, total, err = fetchAllREST("task", func(cur string) ([]mgmt.ReportTask, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ReportTasksList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				tasks, pag, err = c.ReportTasksList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Scope", "Frequency", "Type", "Day"}
			rows := make([][]string, len(tasks))
			for i, t := range tasks {
				rows[i] = []string{
					t.ID, truncate(t.Name, 40), t.Scope,
					t.Frequency, t.ScheduleType, orDash(t.Day),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, tasks, len(tasks), total, "task", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&scope, "scope", "", "filter by scope (group, site, account, tenant)")
	cmd.Flags().StringVar(&frequency, "frequency", "", "filter by frequency (manually, weekly, monthly)")
	cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "filter by schedule type (manually, scheduled)")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. name, frequency, scope)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return cmd
}

func newReportTypesCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string

	cmd := &cobra.Command{
		Use:   "types",
		Short: "List available report insight types",
		Long:  "List available report insight types. Output is always JSON because the schema is opaque.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.InsightTypesParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
			}
			types, err := c.ReportsInsightTypes(cmd.Context(), params)
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), types)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	return cmd
}

func newReportCreateCmd() *cobra.Command {
	var (
		name            string
		scheduleType    string
		insightTypesRaw string
		frequency       string
		day             string
		fromDate        string
		toDate          string
		attachmentTypes []string
		recipients      []string
		siteIDs         []string
		accountIDs      []string
		scope           string
		trend           bool
		yes             bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a report task",
		Long: `Create a new report task or schedule.

Schedule types: manually, scheduled
Frequencies: manually, weekly, monthly
Days (for weekly): sunday, monday, tuesday, wednesday, thursday, friday, saturday

Use "reports types" to list available insight types, then pass them
as a JSON array via --insight-types.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if scheduleType == "" {
				return fmt.Errorf("--schedule-type is required")
			}
			if insightTypesRaw == "" {
				return fmt.Errorf("--insight-types is required (JSON array from 'reports types')")
			}

			var insightTypes json.RawMessage
			if err := json.Unmarshal([]byte(insightTypesRaw), &insightTypes); err != nil {
				return fmt.Errorf("--insight-types must be valid JSON: %w", err)
			}

			task := mgmt.ReportTaskCreate{
				Name:            name,
				ScheduleType:    scheduleType,
				InsightTypes:    insightTypes,
				Frequency:       frequency,
				Day:             day,
				FromDate:        fromDate,
				ToDate:          toDate,
				AttachmentTypes: attachmentTypes,
				Recipients:      recipients,
			}
			if cmd.Flags().Changed("trend") {
				task.IsTrend = &trend
			}

			return guard(cmd.OutOrStdout(), "reports create", "create report task "+name+" ("+scheduleType+")", name, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.ReportTasksCreate(cmd.Context(), siteIDs, accountIDs, scope, task); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{"success": true, "name": name})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created report task %q\n", name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "report task name (required)")
	cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "schedule type: manually, scheduled (required)")
	cmd.Flags().StringVar(&insightTypesRaw, "insight-types", "", "insight types as JSON array (required; see 'reports types')")
	cmd.Flags().StringVar(&frequency, "frequency", "", "frequency: manually, weekly, monthly")
	cmd.Flags().StringVar(&day, "day", "", "day of week for weekly schedules")
	cmd.Flags().StringVar(&fromDate, "from-date", "", "report date range start (ISO timestamp)")
	cmd.Flags().StringVar(&toDate, "to-date", "", "report date range end (ISO timestamp)")
	cmd.Flags().StringSliceVar(&attachmentTypes, "attachment-type", nil, "attachment types (pdf, html)")
	cmd.Flags().StringSliceVar(&recipients, "recipient", nil, "email recipients")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "target site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "target account IDs")
	cmd.Flags().StringVar(&scope, "scope", "", "scope filter")
	cmd.Flags().BoolVar(&trend, "trend", false, "trend report (period = last month)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newReportDownloadCmd() *cobra.Command {
	var (
		format string
		output string
	)

	cmd := &cobra.Command{
		Use:   "download <report-id>",
		Short: "Download a generated report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reportID := args[0]
			if format != "pdf" && format != "html" {
				return fmt.Errorf("--format must be pdf or html")
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}
			data, err := c.ReportDownload(cmd.Context(), reportID, format)
			if err != nil {
				return err
			}

			if output == "" {
				output = fmt.Sprintf("report-%s.%s", reportID, format)
			}
			if err := os.WriteFile(output, data, 0o644); err != nil {
				return fmt.Errorf("write file: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Downloaded report to %s (%d bytes)\n", output, len(data))
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "pdf", "report format (pdf, html)")
	cmd.Flags().StringVar(&output, "output", "", "output file path (default: report-<id>.<format>)")
	return cmd
}
