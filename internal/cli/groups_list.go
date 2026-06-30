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
	cmd.AddCommand(newGroupsGetCmd())
	return cmd
}

func newGroupsListCmd() *cobra.Command {
	var siteIDs []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List groups",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			groups, pag, err := c.GroupsList(cmd.Context(), &mgmt.GroupListParams{
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(groups)
			}
			var rows [][]string
			for _, g := range groups {
				rows = append(rows, []string{
					g.ID, g.Name, g.Type, strconv.Itoa(g.TotalAgents),
					boolIcon(g.IsDefault), g.SiteID,
				})
			}
			printTable([]string{"ID", "Name", "Type", "Agents", "Default", "Site"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "group"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}

func newGroupsGetCmd() *cobra.Command {
	return &cobra.Command{
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
			if jsonOutput {
				return printJSON(g)
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
}
