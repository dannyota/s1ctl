package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newMisconfigurationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "misconfigurations",
		Aliases: []string{"misconfigs"},
		Short:   "Manage xSPM misconfigurations",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newMisconfigurationsListCmd())
	cmd.AddCommand(newMisconfigurationsGetCmd())
	cmd.AddCommand(newMisconfigurationsStatusCmd())
	cmd.AddCommand(newMisconfigurationsVerdictCmd())
	cmd.AddCommand(newMisconfigurationsNotesCmd())
	cmd.AddCommand(newMisconfigurationsNoteAddCmd())
	cmd.AddCommand(newMisconfigurationsNoteUpdateCmd())
	cmd.AddCommand(newMisconfigurationsNoteDeleteCmd())
	cmd.AddCommand(newMisconfigurationsAssignCmd())
	cmd.AddCommand(newMisconfigurationsHistoryCmd())
	cmd.AddCommand(newMisconfigurationsRelatedAssetsCmd())
	cmd.AddCommand(newMisconfigurationsExportCmd())
	return cmd
}

func newMisconfigurationsListCmd() *cobra.Command {
	var severities, statuses []string
	var after string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List misconfigurations",
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

			var items []graphql.Misconfiguration
			var total int64

			if all {
				items, total, err = fetchAllGQL("misconfiguration", func(cur string) (*graphql.Connection[graphql.Misconfiguration], error) {
					params.After = cur
					return c.MisconfigurationsList(cmd.Context(), params)
				})
			} else {
				conn, connErr := c.MisconfigurationsList(cmd.Context(), params)
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

			headers := []string{"ID", "Name", "Severity", "Status", "Environment", "Product", "Site"}
			rows := make([][]string, len(items))
			for i, m := range items {
				rows[i] = []string{
					m.ID, truncate(orDash(m.Name), 40), m.Severity,
					m.Status, orDash(m.Environment), orDash(m.Product),
					orDash(m.Scope.Site.Name),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), int(total), "misconfiguration", all)
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&after, "after", "", "pagination cursor")
	return markJSON(cmd)
}

func newMisconfigurationsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get misconfiguration details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			m, err := c.MisconfigurationsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), m)
			}
			rows := [][]string{
				{"ID", m.ID},
				{"Name", orDash(m.Name)},
				{"Severity", m.Severity},
				{"Status", m.Status},
				{"Verdict", orDash(m.AnalystVerdict)},
				{"Environment", orDash(m.Environment)},
				{"Product", orDash(m.Product)},
				{"Vendor", orDash(m.Vendor)},
				{"Asset", orDash(m.Asset.Name)},
				{"Site", orDash(m.Scope.Site.Name)},
				{"Detected", orDash(m.DetectedAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}

func newMisconfigurationsStatusCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "status <id> <status>",
		Short: "Update misconfiguration status",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, status := args[0], args[1]
			return guard(cmd.OutOrStdout(), "misconfigurations status", fmt.Sprintf("set status=%s on misconfiguration %s", status, id), id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.MisconfigurationsUpdateStatus(cmd.Context(), []string{id}, status); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "status: updated misconfiguration %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newMisconfigurationsVerdictCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "verdict <id> <verdict>",
		Short: "Update misconfiguration analyst verdict (TRUE_POSITIVE, FALSE_POSITIVE)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, verdict := args[0], args[1]
			return guard(cmd.OutOrStdout(), "misconfigurations verdict", fmt.Sprintf("set verdict=%s on misconfiguration %s", verdict, id), id, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if err := c.MisconfigurationsUpdateVerdict(cmd.Context(), []string{id}, verdict); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"verdict": "updated", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "verdict: updated misconfiguration %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}
