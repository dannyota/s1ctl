package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newSitesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sites",
		Short: "Manage sites",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newSitesListCmd())
	cmd.AddCommand(newSitesCountCmd())
	cmd.AddCommand(newSitesGetCmd())
	return cmd
}

func newSitesListCmd() *cobra.Command {
	var accountIDs, states []string
	var query, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sites",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.SiteListParams{
				AccountIDs: accountIDs,
				States:     states,
				Query:      query,
				Limit:      limit,
				Cursor:     cursor,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var sites []mgmt.Site
			var total int

			if all {
				sites, total, err = fetchAllREST("site", func(cur string) ([]mgmt.Site, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.SitesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				sites, pag, err = c.SitesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "State", "Type", "Licenses"}
			rows := make([][]string, len(sites))
			for i, s := range sites {
				rows[i] = []string{
					s.ID, s.Name, s.State, s.SiteType,
					fmt.Sprintf("%d / %d", s.ActiveLicenses, s.TotalLicenses),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, sites, len(sites), total, "site", all)
		},
	}
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&states, "state", nil, "filter by state")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. name, state)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return cmd
}

func newSitesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <site-id>",
		Short: "Get site details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			s, err := c.SitesGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), s)
			}
			rows := [][]string{
				{"ID", s.ID},
				{"Name", s.Name},
				{"State", s.State},
				{"Type", s.SiteType},
				{"Account", fmt.Sprintf("%s (%s)", s.AccountName, s.AccountID)},
				{"Licenses", fmt.Sprintf("%d / %d", s.ActiveLicenses, s.TotalLicenses)},
				{"Expiration", s.Expiration},
				{"Created", s.CreatedAt},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}

func newSitesCountCmd() *cobra.Command {
	var accountIDs []string

	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count sites",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			count, err := c.SitesCount(cmd.Context(), &mgmt.SiteListParams{AccountIDs: accountIDs})
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
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return cmd
}
