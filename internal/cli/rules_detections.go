package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRulesDetectionsCmd() *cobra.Command {
	var siteIDs, severity, status []string
	var since, cursor, sortBy, sortOrder, groupBy string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "detections <rule-name>",
		Short: "List recent detections for a rule",
		Long: `Fetch cloud detection alerts (STAR alerts) filtered by rule name.
Shows what a specific rule is catching.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.CDAlertListParams{
				SiteIDs:          siteIDs,
				RuleNameContains: []string{args[0]},
				Severity:         severity,
				IncidentStatus:   status,
				ReportedAtGt:     since,
				Limit:            limit,
				Cursor:           cursor,
				SortBy:           sortBy,
				SortOrder:        sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}
			if params.SortBy == "" {
				params.SortBy = "id"
				params.SortOrder = "desc"
			}

			var alerts []mgmt.CloudDetectionAlert
			var total int

			if all {
				alerts, total, err = fetchAllREST("detection", func(cur string) ([]mgmt.CloudDetectionAlert, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.CloudDetectionAlertsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				alerts, pag, err = c.CloudDetectionAlertsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			if groupBy == "agent" {
				return printDetectionsByAgent(cmd, alerts, total)
			}

			headers := []string{"Alert ID", "Agent", "Event", "Severity", "Status", "Reported"}
			rows := make([][]string, len(alerts))
			for i, a := range alerts {
				rows[i] = []string{
					a.AlertInfo.AlertID,
					truncate(orDash(a.AgentDetectionInfo.Name), 25),
					orDash(a.AlertInfo.EventType),
					orDash(a.RuleInfo.Severity),
					orDash(a.AlertInfo.IncidentStatus),
					orDash(a.AlertInfo.ReportedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, alerts, len(alerts), total, "detection", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&severity, "severity", nil, "filter by severity")
	cmd.Flags().StringSliceVar(&status, "status", nil, "filter by incident status")
	cmd.Flags().StringVar(&since, "since", "", "show detections after this time (RFC3339)")
	cmd.Flags().IntVar(&limit, "limit", 0, fmt.Sprintf("max results per page (default %d)", defaultPageSize))
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (default: id)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (default: desc)")
	cmd.Flags().StringVar(&groupBy, "group-by", "", "group results (agent)")
	return cmd
}

func printDetectionsByAgent(cmd *cobra.Command, alerts []mgmt.CloudDetectionAlert, total int) error {
	type agentSummary struct {
		Name  string `json:"name"`
		OS    string `json:"os"`
		Count int    `json:"count"`
	}
	counts := map[string]*agentSummary{}
	for _, a := range alerts {
		name := orDash(a.AgentDetectionInfo.Name)
		s, ok := counts[name]
		if !ok {
			s = &agentSummary{Name: name, OS: orDash(a.AgentDetectionInfo.OSFamily)}
			counts[name] = s
		}
		s.Count++
	}
	sorted := make([]*agentSummary, 0, len(counts))
	for _, s := range counts {
		sorted = append(sorted, s)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Count > sorted[j].Count })

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), sorted)
	}

	headers := []string{"Agent", "OS", "Detections"}
	rows := make([][]string, len(sorted))
	for i, s := range sorted {
		rows[i] = []string{truncate(s.Name, 30), s.OS, fmt.Sprintf("%d", s.Count)}
	}
	return printOutput(cmd.OutOrStdout(), headers, rows, sorted, len(sorted), total, "agent", false)
}
