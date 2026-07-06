package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newGroupsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "Manage groups",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newGroupsListCmd())
	cmd.AddCommand(newGroupsCountCmd())
	cmd.AddCommand(newGroupsGetCmd())
	cmd.AddCommand(newGroupsCreateCmd())
	cmd.AddCommand(newGroupsUpdateCmd())
	cmd.AddCommand(newGroupsDeleteCmd())
	addGroupSyncCmds(cmd)
	return cmd
}

func newGroupsListCmd() *cobra.Command {
	var siteIDs []string
	var query, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.GroupListParams{
				SiteIDs:   siteIDs,
				Query:     query,
				Limit:     limit,
				Cursor:    cursor,
				SortBy:    sortBy,
				SortOrder: sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var groups []mgmt.Group
			var total int

			if all {
				groups, total, err = fetchAllREST("group", func(cur string) ([]mgmt.Group, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.GroupsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				groups, pag, err = c.GroupsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Type", "Agents", "Default", "Site"}
			rows := make([][]string, len(groups))
			for i, g := range groups {
				rows[i] = []string{
					g.ID, g.Name, g.Type, strconv.Itoa(g.TotalAgents),
					boolIcon(g.IsDefault), g.SiteID,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, groups, len(groups), total, "group", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. name, type)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return markJSON(cmd)
}

func newGroupsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get group details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			g, err := c.GroupsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), g)
			}
			rows := [][]string{
				{"ID", g.ID},
				{"Name", g.Name},
				{"Type", g.Type},
				{"Agents", strconv.Itoa(g.TotalAgents)},
				{"Default", boolIcon(g.IsDefault)},
				{"Site", g.SiteID},
				{"Created", g.CreatedAt},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}

func newGroupsCountCmd() *cobra.Command {
	var siteIDs []string

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count groups",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			count, err := c.GroupsCount(cmd.Context(), &mgmt.GroupListParams{SiteIDs: siteIDs})
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
	return markJSON(cmd)
}
