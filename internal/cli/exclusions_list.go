package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newExclusionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exclusions",
		Short: "Manage exclusions and blocklist",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newExclusionsListCmd())
	cmd.AddCommand(newExclusionsGetCmd())
	cmd.AddCommand(newExclusionsCreateCmd())
	addExclusionSyncCmds(cmd)
	return cmd
}

func newExclusionsListCmd() *cobra.Command {
	var siteIDs, types, osTypes []string
	var query, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List exclusions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ExclusionListParams{
				SiteIDs:   siteIDs,
				Types:     types,
				OSTypes:   osTypes,
				Query:     query,
				Limit:     limit,
				Cursor:    cursor,
				SortBy:    sortBy,
				SortOrder: sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var exclusions []mgmt.Exclusion
			var total int

			if all {
				exclusions, total, err = fetchAllREST("exclusion", func(cur string) ([]mgmt.Exclusion, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ExclusionsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				exclusions, pag, err = c.ExclusionsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Type", "Value", "OS", "Mode"}
			rows := make([][]string, len(exclusions))
			for i, e := range exclusions {
				rows[i] = []string{
					e.ID, e.Type, truncate(e.Value, 50), e.OSType, e.Mode,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, exclusions, len(exclusions), total, "exclusion", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&types, "type", nil, "filter by exclusion type")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. type, osType)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return cmd
}

func newExclusionsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <exclusion-id>",
		Short: "Get exclusion details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			e, err := c.ExclusionsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), e)
			}
			rows := [][]string{
				{"ID", e.ID},
				{"Type", e.Type},
				{"Value", e.Value},
				{"OS", e.OSType},
				{"Mode", e.Mode},
				{"Description", e.Description},
				{"Scope", e.ScopeName},
				{"User", e.UserName},
				{"Created", e.CreatedAt},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}
