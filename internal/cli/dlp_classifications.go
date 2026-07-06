package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newDLPClassificationsListCmd() *cobra.Command {
	var types []string
	var search, scopeLevel, scopeID string
	var limit, page int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List DLP classifications",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			filter := dlpClassificationFilter(search, types)
			pageSize := limit
			if pageSize == 0 {
				pageSize = defaultPageSize
			}
			startPage := page
			if startPage == 0 {
				startPage = 1
			}

			var items []graphql.DLPClassification
			var total int
			if all {
				items, total, err = fetchAllDLP("dlp classification", func(p int) (*graphql.DLPConnection[graphql.DLPClassification], error) {
					return c.DLPClassificationsList(cmd.Context(), filter, scope, &graphql.DLPPage{Page: p, PageSize: pageSize})
				})
			} else {
				conn, connErr := c.DLPClassificationsList(cmd.Context(), filter, scope, &graphql.DLPPage{Page: startPage, PageSize: pageSize})
				if connErr != nil {
					return connErr
				}
				total = conn.PageInfo.TotalCount
				items = conn.Nodes
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Type", "System", "Used In Rules"}
			rows := make([][]string, len(items))
			for i, cl := range items {
				rows[i] = []string{
					cl.ID, truncate(orDash(cl.Name), 40), string(cl.Type),
					boolIcon(cl.SystemPolicy), fmt.Sprintf("%d", cl.UsedInRulesCount),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "dlp classification", all)
		},
	}
	cmd.Flags().StringSliceVar(&types, "type", nil, "filter by type (REGEX, SECRETS, SENSITIVE_DATA, CODE_DETECTOR, AI_CONTEXTUAL)")
	cmd.Flags().StringVar(&search, "search", "", "filter by classification name (prefix match)")
	cmd.Flags().IntVar(&limit, "limit", 0, "page size (default 50)")
	cmd.Flags().IntVar(&page, "page", 0, "page number (1-indexed)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	addDLPScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

func newDLPClassificationsGetCmd() *cobra.Command {
	var scopeLevel, scopeID string

	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get DLP classification details",
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
			cl, err := c.DLPClassificationGet(cmd.Context(), args[0], scope)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), cl)
			}
			rows := [][]string{
				{"ID", cl.ID},
				{"Name", orDash(cl.Name)},
				{"Description", orDash(cl.Description)},
				{"Type", string(cl.Type)},
				{"Classification Code", orDash(cl.ClassificationCode)},
				{"System", boolIcon(cl.SystemPolicy)},
				{"Used In Rules", fmt.Sprintf("%d", cl.UsedInRulesCount)},
				{"Scope", orDash(cl.Scope.Path)},
				{"Created", orDash(cl.CreatedAt)},
				{"Updated", orDash(cl.UpdatedAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	addDLPScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

func newDLPClassificationDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a DLP classification",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "dlp classifications delete", "delete dlp classification "+args[0], args[0], yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				ok, err := c.DLPClassificationDelete(cmd.Context(), args[0])
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{"id": args[0], "deleted": ok})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted DLP classification %s (%s)\n", args[0], boolIcon(ok))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply (default: dry-run)")
	return markJSON(cmd)
}

// dlpClassificationFilter builds a classification filter from flags, or nil.
func dlpClassificationFilter(search string, types []string) *graphql.DLPClassificationFilter {
	f := &graphql.DLPClassificationFilter{}
	set := false
	if search != "" {
		f.SearchName = search
		set = true
	}
	for _, t := range types {
		f.Type = append(f.Type, graphql.DLPClassificationType(strings.ToUpper(t)))
		set = true
	}
	if !set {
		return nil
	}
	return f
}
