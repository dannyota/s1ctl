package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage custom detection rules (STAR)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newRulesListCmd())
	cmd.AddCommand(newRulesGetCmd())
	cmd.AddCommand(newRulesHealthCmd())
	cmd.AddCommand(newRulesTrendsCmd())
	cmd.AddCommand(newRulesDetectionsCmd())
	cmd.AddCommand(newRulesDiffCmd())
	cmd.AddCommand(newRulesValidateCmd())
	cmd.AddCommand(newRulesEnableCmd())
	cmd.AddCommand(newRulesDisableCmd())
	addRuleSyncCmds(cmd)
	return cmd
}

func newRulesListCmd() *cobra.Command {
	var siteIDs, status, severity, scopes, queryType []string
	var nameContains, query, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List custom detection rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.RuleListParams{
				SiteIDs:      siteIDs,
				Status:       status,
				Severity:     severity,
				Scopes:       scopes,
				QueryType:    queryType,
				NameContains: nameContains,
				Query:        query,
				Limit:        limit,
				Cursor:       cursor,
				SortBy:       sortBy,
				SortOrder:    sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var rules []mgmt.Rule
			var total int

			if all {
				rules, total, err = fetchAllREST("rule", func(cur string) ([]mgmt.Rule, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.RulesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				rules, pag, err = c.RulesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Status", "Severity", "Scope", "Created"}
			rows := make([][]string, len(rules))
			for i, r := range rules {
				rows[i] = []string{
					r.ID,
					truncate(r.Name, 40),
					string(r.Status),
					string(r.Severity),
					string(r.Scope),
					r.CreatedAt,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, rules, len(rules), total, "rule", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&status, "status", nil, "filter by status (Draft, Active, Disabled, ...)")
	cmd.Flags().StringSliceVar(&severity, "severity", nil, "filter by severity (Info, Low, Medium, High, Critical)")
	cmd.Flags().StringSliceVar(&scopes, "scope", nil, "filter by scope (global, account, site, group)")
	cmd.Flags().StringSliceVar(&queryType, "query-type", nil, "filter by query type (events, correlation, scheduled)")
	cmd.Flags().StringVar(&nameContains, "name", "", "filter by rule name (substring match)")
	cmd.Flags().StringVar(&query, "query", "", "free text search on S1QL")
	cmd.Flags().IntVar(&limit, "limit", 0, fmt.Sprintf("max results per page (default %d)", defaultPageSize))
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. name, severity, createdAt)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return markJSON(cmd)
}

func newRulesGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <rule-id>",
		Short: "Get custom detection rule details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			r, err := c.RulesGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), r)
			}
			rows := [][]string{
				{"ID", r.ID},
				{"Name", r.Name},
				{"Description", orDash(r.Description)},
				{"Status", string(r.Status)},
				{"Status Reason", orDash(r.StatusReason)},
				{"Severity", string(r.Severity)},
				{"Query Type", string(r.QueryType)},
				{"S1QL", truncate(r.S1QL, 80)},
				{"Scope", string(r.Scope)},
				{"Expiration", orDash(r.Expiration)},
				{"Treat As Threat", string(r.TreatAsThreat)},
				{"Active Response", boolIcon(r.ActiveResponse)},
				{"Alerts", fmt.Sprintf("%d", r.GeneratedAlerts)},
				{"Creator", orDash(r.Creator)},
				{"Account", orDash(r.AccountName)},
				{"Site", orDash(r.SiteName)},
				{"Created", r.CreatedAt},
				{"Updated", r.UpdatedAt},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}
