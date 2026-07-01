package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newThreatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "threats",
		Short: "Manage threats",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newThreatsListCmd())
	cmd.AddCommand(newThreatsCountCmd())
	cmd.AddCommand(newThreatsGetCmd())
	cmd.AddCommand(newThreatNotesCmd())
	cmd.AddCommand(newThreatAddNoteCmd())
	addThreatActions(cmd)
	return cmd
}

func newThreatsListCmd() *cobra.Command {
	var siteIDs, classifications, statuses, verdicts []string
	var query, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List threats",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ThreatListParams{
				SiteIDs:          siteIDs,
				Classifications:  classifications,
				IncidentStatuses: statuses,
				AnalystVerdicts:  verdicts,
				Query:            query,
				Limit:            limit,
				Cursor:           cursor,
				SortBy:           sortBy,
				SortOrder:        sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var threats []mgmt.Threat
			var total int

			if all {
				threats, total, err = fetchAllREST("threat", func(cur string) ([]mgmt.Threat, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ThreatsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				threats, pag, err = c.ThreatsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Class", "Mitigation", "Verdict", "Status"}
			rows := make([][]string, len(threats))
			for i, t := range threats {
				rows[i] = []string{
					t.ID, truncate(t.ThreatName, 40), t.Classification,
					t.MitigationStatus, t.AnalystVerdict, t.IncidentStatus,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, threats, len(threats), total, "threat", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&classifications, "classification", nil, "filter by classification")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by incident status")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter by analyst verdict")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. createdAt, classification)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return cmd
}

func newThreatsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <threat-id>",
		Short: "Get threat details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			t, err := c.ThreatsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), t)
			}
			rows := [][]string{
				{"ID", t.ID},
				{"Name", t.ThreatName},
				{"Classification", t.Classification},
				{"Confidence", t.ConfidenceLevel},
				{"Mitigation", t.MitigationStatus},
				{"Verdict", t.AnalystVerdict},
				{"Status", t.IncidentStatus},
				{"Agent", t.AgentID},
				{"Created", t.CreatedAt},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}

func newThreatsCountCmd() *cobra.Command {
	var siteIDs []string

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count threats",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			count, err := c.ThreatsCount(cmd.Context(), &mgmt.ThreatListParams{SiteIDs: siteIDs})
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]int{"count": count})
			}
			fmt.Fprintln(cmd.OutOrStdout(), count)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}
