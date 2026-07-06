package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newDLPRulesListCmd() *cobra.Command {
	var statuses []string
	var search, scopeLevel, scopeID string
	var limit, page int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List data protection rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			filter := dlpRuleFilter(search, statuses)
			pageSize := limit
			if pageSize == 0 {
				pageSize = defaultPageSize
			}
			startPage := page
			if startPage == 0 {
				startPage = 1
			}

			var items []graphql.DLPRule
			var total int
			if all {
				items, total, err = fetchAllDLP("data protection rule", func(p int) (*graphql.DLPConnection[graphql.DLPRule], error) {
					return c.DLPRulesList(cmd.Context(), filter, scope, &graphql.DLPPage{Page: p, PageSize: pageSize})
				})
			} else {
				conn, connErr := c.DLPRulesList(cmd.Context(), filter, scope, &graphql.DLPPage{Page: startPage, PageSize: pageSize})
				if connErr != nil {
					return connErr
				}
				total = conn.PageInfo.TotalCount
				items = conn.Nodes
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Status", "Rank", "System", "Classifications"}
			rows := make([][]string, len(items))
			for i, r := range items {
				rows[i] = []string{
					r.ID, truncate(orDash(r.Name), 40), string(r.Status),
					fmt.Sprintf("%d", r.Rank), boolIcon(r.SystemPolicy),
					truncate(dlpClassificationNames(r.Classifications), 40),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "data protection rule", all)
		},
	}
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status (ENABLED, DISABLED)")
	cmd.Flags().StringVar(&search, "search", "", "filter by rule name (prefix match)")
	cmd.Flags().IntVar(&limit, "limit", 0, "page size (default 50)")
	cmd.Flags().IntVar(&page, "page", 0, "page number (1-indexed)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	addDLPScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

func newDLPRulesGetCmd() *cobra.Command {
	var scopeLevel, scopeID string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get data protection rule details",
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
			r, err := c.DLPRuleGet(cmd.Context(), args[0], scope)
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
				{"Status", string(r.Status)},
				{"Rank", fmt.Sprintf("%d", r.Rank)},
				{"Rule Code", orDash(r.RuleCode)},
				{"System", boolIcon(r.SystemPolicy)},
				{"Classifications", orDash(dlpClassificationNames(r.Classifications))},
				{"Scope", orDash(r.Scope.Path)},
				{"Created", orDash(r.CreatedAt)},
				{"Created By", orDash(r.CreatedBy)},
				{"Updated", orDash(r.UpdatedAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	addDLPScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

// addDLPRuleActions registers the guarded enable/disable/delete commands. Each
// applies its verb to one or more rule IDs in a single bulk request.
func addDLPRuleActions(parent *cobra.Command) {
	parent.AddCommand(newDLPRuleActionCmd("enable", "Enable data protection rules"))
	parent.AddCommand(newDLPRuleActionCmd("disable", "Disable data protection rules"))
	parent.AddCommand(newDLPRuleActionCmd("delete", "Delete data protection rules"))
}

func newDLPRuleActionCmd(verb, short string) *cobra.Command {
	var yes bool
	var scopeLevel, scopeID string

	cmd := &cobra.Command{
		Use:   verb + " <id> [id...]",
		Short: short,
		Long:  short + ". One or more IDs are applied together in a single bulk request.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			action := verb + " " + pluralize(len(args), "data protection rule")
			return guard(cmd.OutOrStdout(), "dlp rules "+verb, action, strings.Join(args, ","), yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				affected, err := runDLPRuleAction(cmd.Context(), c, verb, args, scope)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{
						"action":   verb,
						"affected": affected,
						"ids":      args,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", verb, pluralize(affected, "data protection rule"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	addDLPScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

// runDLPRuleAction dispatches to the bulk enable/disable/delete SDK methods and
// reports how many rules were affected. Enable/disable return the affected
// rules; delete returns a boolean, so its count is the input count on success.
func runDLPRuleAction(ctx context.Context, c *graphql.Client, verb string, ids []string, scope *graphql.Scope) (int, error) {
	switch verb {
	case "enable":
		rules, err := c.DLPRulesBulkEnable(ctx, ids, scope)
		return len(rules), err
	case "disable":
		rules, err := c.DLPRulesBulkDisable(ctx, ids, scope)
		return len(rules), err
	case "delete":
		ok, err := c.DLPRulesBulkDelete(ctx, ids, scope)
		if err != nil {
			return 0, err
		}
		if ok {
			return len(ids), nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("unknown action %q", verb)
}

// dlpRuleFilter builds a rule filter from flags, or nil when none are set.
func dlpRuleFilter(search string, statuses []string) *graphql.DLPRuleFilter {
	f := &graphql.DLPRuleFilter{}
	set := false
	if search != "" {
		f.SearchName = search
		set = true
	}
	for _, s := range statuses {
		f.Status = append(f.Status, graphql.DLPRuleStatus(strings.ToUpper(s)))
		set = true
	}
	if !set {
		return nil
	}
	return f
}

// dlpClassificationNames joins the names of a rule's associated classifications.
func dlpClassificationNames(cs []graphql.DLPClassificationSummary) string {
	names := make([]string, len(cs))
	for i, c := range cs {
		names[i] = c.Name
	}
	return joinOrDash(names)
}
