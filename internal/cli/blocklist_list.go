package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newBlocklistCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocklist",
		Short: "Manage the blocklist (blocked file hashes)",
		Long: `Manage the SentinelOne blocklist (restrictions).

The blocklist holds SHA1/SHA256 hashes that agents block from executing. Items
are scoped globally (tenant) or to accounts, sites, or groups.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newBlocklistListCmd())
	cmd.AddCommand(newBlocklistCreateCmd())
	cmd.AddCommand(newBlocklistUpdateCmd())
	cmd.AddCommand(newBlocklistDeleteCmd())
	cmd.AddCommand(newBlocklistValidateCmd())
	cmd.AddCommand(newBlocklistExportCmd())
	addBlocklistSyncCmds(cmd)
	return cmd
}

func newBlocklistListCmd() *cobra.Command {
	var siteIDs, groupIDs, accountIDs, osTypes []string
	var query, value, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List blocklist items",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.BlocklistListParams{
				SiteIDs:    siteIDs,
				GroupIDs:   groupIDs,
				AccountIDs: accountIDs,
				OSTypes:    osTypes,
				Query:      query,
				Value:      value,
				Limit:      limit,
				Cursor:     cursor,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.BlocklistItem
			var total int

			if all {
				items, total, err = fetchAllREST("blocklist item", func(cur string) ([]mgmt.BlocklistItem, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.BlocklistList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				items, pag, err = c.BlocklistList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Value (SHA1)", "SHA256", "OS", "Source", "Description"}
			rows := make([][]string, len(items))
			for i, b := range items {
				rows[i] = []string{
					b.ID, truncate(b.Value, 40), truncate(b.SHA256Value, 20),
					b.OSType, b.Source, truncate(b.Description, 40),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "blocklist item", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type (windows, linux, macos, windows_legacy)")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().StringVar(&value, "value", "", "filter by hash value")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. createdAt, osType)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return markJSON(cmd)
}

func newBlocklistExportCmd() *cobra.Command {
	var siteIDs, groupIDs, accountIDs, osTypes []string
	var tenant bool
	var outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export blocklist items as CSV",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.BlocklistListParams{
				SiteIDs:    siteIDs,
				GroupIDs:   groupIDs,
				AccountIDs: accountIDs,
				OSTypes:    osTypes,
			}
			if tenant {
				params.Tenant = &tenant
			}

			data, err := c.BlocklistExport(cmd.Context(), params)
			if err != nil {
				return err
			}

			if outFile != "" {
				if err := os.WriteFile(outFile, data, 0o644); err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Exported to %s\n", outFile)
				return nil
			}
			_, err = cmd.OutOrStdout().Write(data)
			return err
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type")
	cmd.Flags().BoolVar(&tenant, "tenant", false, "export the global (tenant) blocklist")
	cmd.Flags().StringVar(&outFile, "out", "", "write export to file (default: stdout)")
	return cmd
}
