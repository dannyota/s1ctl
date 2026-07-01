package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newCloudPoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cloud-policies",
		Aliases: []string{"cloud"},
		Short:   "Manage cloud security policies (CNS rules)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newCloudPoliciesListCmd())
	cmd.AddCommand(newCloudPoliciesGetCmd())
	return cmd
}

func newCloudPoliciesListCmd() *cobra.Command {
	var severities, statuses []string
	var after string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List cloud security policies",
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

			var items []graphql.CloudPolicy
			var total int64

			if all {
				items, total, err = fetchAllGQL("cloud policy", func(cur string) (*graphql.Connection[graphql.CloudPolicy], error) {
					params.After = cur
					return c.CloudPoliciesList(cmd.Context(), params)
				})
			} else {
				conn, connErr := c.CloudPoliciesList(cmd.Context(), params)
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

			headers := []string{"ID", "Name", "Severity", "Status", "Type", "Category", "Providers"}
			rows := make([][]string, len(items))
			for i, p := range items {
				rows[i] = []string{
					p.ID, truncate(orDash(p.Name), 40), p.Severity,
					p.Status, orDash(p.Type), orDash(p.Category),
					joinOrDash(p.Providers),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), int(total), "cloud policy", all)
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (HIGH, CRITICAL, etc.)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&after, "after", "", "pagination cursor")
	return cmd
}

func newCloudPoliciesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get cloud policy details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			p, err := c.CloudPoliciesGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), p)
			}
			rows := [][]string{
				{"ID", p.ID},
				{"Name", orDash(p.Name)},
				{"Description", orDash(p.Description)},
				{"Severity", p.Severity},
				{"Status", p.Status},
				{"Type", orDash(p.Type)},
				{"Policy Code", orDash(p.PolicyCode)},
				{"Category", orDash(p.Category)},
				{"Sub-Category", orDash(p.SubCategory)},
				{"Resource Type", orDash(p.ResourceType)},
				{"Providers", joinOrDash(p.Providers)},
				{"System", boolIcon(p.IsSystem)},
				{"Impact", orDash(p.Impact)},
				{"Recommended Action", orDash(p.RecommendedAction)},
				{"Issue Message", orDash(p.IssueMessage)},
				{"Reference", orDash(p.Reference)},
				{"Scope", orDash(p.Scope.Path)},
				{"Created", orDash(p.CreatedAt)},
				{"Updated", orDash(p.UpdatedAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}

// joinOrDash joins a string slice with commas, returning "-" if empty.
func joinOrDash(ss []string) string {
	if len(ss) == 0 {
		return "-"
	}
	return strings.Join(ss, ", ")
}
