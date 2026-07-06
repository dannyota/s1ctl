package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newCloudRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cloud-rules",
		Aliases: []string{"cns"},
		Short:   "Manage CNS custom cloud rules (Cloud Native Security)",
		Long: `Manage CNS (Cloud Native Security) custom cloud rules.

CNS rules are Rego- or graph-query-based cloud policies. This surface supports
full lifecycle: list, get, create/update from a rule JSON file, enable/disable/
delete, evaluate a Rego query against asset JSON, and inspect supported rule
types. Rule bodies (Rego plus config) are supplied via files, not flags.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newCloudRulesListCmd())
	cmd.AddCommand(newCloudRulesGetCmd())
	cmd.AddCommand(newCloudRulesTypesCmd())
	addCloudRuleMutations(cmd)
	return cmd
}

func newCloudRulesListCmd() *cobra.Command {
	var severities, statuses []string
	var scopeLevel, scopeID, after string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List CNS custom cloud rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}

			var filters []graphql.Filter
			if len(severities) > 0 {
				filters = append(filters, graphql.Filter{
					FieldID:  "severity",
					StringIn: &graphql.InStr{Values: severities},
				})
			}
			if len(statuses) > 0 {
				filters = append(filters, graphql.Filter{
					FieldID:  "status",
					StringIn: &graphql.InStr{Values: statuses},
				})
			}

			page := &graphql.ListParams{First: limit, After: after}
			if page.First == 0 {
				page.First = defaultPageSize
			}

			var items []graphql.CNSRule
			var total int64

			if all {
				items, total, err = fetchAllGQL("cns rule", func(cur string) (*graphql.Connection[graphql.CNSRule], error) {
					page.After = cur
					return c.CNSRulesList(cmd.Context(), filters, scope, page)
				})
			} else {
				conn, connErr := c.CNSRulesList(cmd.Context(), filters, scope, page)
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
			for i, r := range items {
				rows[i] = []string{
					r.ID, truncate(orDash(r.Name), 40), r.Severity,
					r.Status, orDash(r.Type), orDash(r.Category),
					joinOrDash(r.Providers),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), int(total), "cns rule", all)
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (LOW, MEDIUM, HIGH, CRITICAL)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&after, "after", "", "pagination cursor")
	addCloudRuleScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

func newCloudRulesGetCmd() *cobra.Command {
	var scopeLevel, scopeID string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get CNS custom cloud rule details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			r, err := c.CNSRuleGet(cmd.Context(), args[0], scope)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), r)
			}
			rows := [][]string{
				{"ID", r.ID},
				{"Name", orDash(r.Name)},
				{"Description", orDash(r.Description)},
				{"Severity", r.Severity},
				{"Status", r.Status},
				{"Type", orDash(r.Type)},
				{"Policy Code", orDash(r.PolicyCode)},
				{"Category", orDash(r.Category)},
				{"Sub-Category", orDash(r.SubCategory)},
				{"Resource Type", orDash(r.ResourceType)},
				{"Providers", joinOrDash(r.Providers)},
				{"Query Type", orDash(r.QueryType)},
				{"System", boolIcon(r.IsSystem)},
				{"Enforcement", orDash(r.EnforcementAction)},
				{"Impact", orDash(r.Impact)},
				{"Recommended Action", orDash(r.RecommendedAction)},
				{"Issue Message", orDash(r.IssueMessage)},
				{"Reference", orDash(r.Reference)},
				{"Scope", orDash(r.Scope.Path)},
				{"Created", orDash(r.CreatedAt)},
				{"Updated", orDash(r.UpdatedAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	addCloudRuleScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

func newCloudRulesTypesCmd() *cobra.Command {
	var scopeLevel, scopeID string

	cmd := &cobra.Command{
		Use:   "types",
		Short: "List supported CNS rule types",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			types, err := c.CNSRuleTypes(cmd.Context(), scope)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), types)
			}
			rows := make([][]string, len(types))
			for i, t := range types {
				rows[i] = []string{t.Key, orDash(t.Title)}
			}
			printTable([]string{"Type", "Title"}, rows)
			return nil
		},
	}
	addCloudRuleScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

// addCloudRuleScopeFlags registers the shared --scope-level/--scope-id flags.
func addCloudRuleScopeFlags(cmd *cobra.Command, level, id *string) {
	cmd.Flags().StringVar(level, "scope-level", "", "scope level (account, site, group)")
	cmd.Flags().StringVar(id, "scope-id", "", "account, site, or group ID")
}
