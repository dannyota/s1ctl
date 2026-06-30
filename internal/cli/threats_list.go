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
	cmd.AddCommand(newThreatsGetCmd())
	addThreatActions(cmd)
	return cmd
}

func newThreatsListCmd() *cobra.Command {
	var siteIDs, classifications, statuses, verdicts []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List threats",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			threats, pag, err := c.ThreatsList(cmd.Context(), &mgmt.ThreatListParams{
				SiteIDs:          siteIDs,
				Classifications:  classifications,
				IncidentStatuses: statuses,
				AnalystVerdicts:  verdicts,
				Query:            query,
				Limit:            limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(threats)
			}
			var rows [][]string
			for _, t := range threats {
				rows = append(rows, []string{
					t.ID, truncate(t.ThreatName, 40), t.Classification,
					t.MitigationStatus, t.AnalystVerdict, t.IncidentStatus,
				})
			}
			printTable([]string{"ID", "Name", "Class", "Mitigation", "Verdict", "Status"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "threat"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&classifications, "classification", nil, "filter by classification")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by incident status")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter by analyst verdict")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
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
			if jsonOutput {
				return printJSON(t)
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
