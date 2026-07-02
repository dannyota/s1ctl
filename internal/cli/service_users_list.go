package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newServiceUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service-users",
		Short:   "Manage service users (API-token identities)",
		Aliases: []string{"service-user"},
		Long: `Manage SentinelOne service users.

A service user is a non-interactive identity that authenticates with an API
token instead of a password. Tokens are shown only once, at creation or when
regenerated with generate-token.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newServiceUsersListCmd())
	cmd.AddCommand(newServiceUsersGetCmd())
	cmd.AddCommand(newServiceUsersExportCmd())
	cmd.AddCommand(newServiceUsersCreateCmd())
	cmd.AddCommand(newServiceUsersUpdateCmd())
	cmd.AddCommand(newServiceUsersDeleteCmd())
	cmd.AddCommand(newServiceUsersBulkDeleteCmd())
	cmd.AddCommand(newServiceUsersGenerateTokenCmd())
	return cmd
}

func newServiceUsersListCmd() *cobra.Command {
	var siteIDs, accountIDs, roleIDs []string
	var query, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List service users",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ServiceUserListParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				RoleIDs:    roleIDs,
				Query:      query,
				Limit:      limit,
				Cursor:     cursor,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.ServiceUser
			var total int

			if all {
				items, total, err = fetchAllREST("service user", func(cur string) ([]mgmt.ServiceUser, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ServiceUsersList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				items, pag, err = c.ServiceUsersList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Scope", "Description", "Expires", "Created"}
			rows := make([][]string, len(items))
			for i, s := range items {
				rows[i] = []string{
					s.ID, s.Name, string(s.Scope), truncate(s.Description, 40),
					orDash(s.APIToken.ExpiresAt), s.CreatedAt,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "service user", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&roleIDs, "role-id", nil, "filter by RBAC role ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search (name, description)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. id, name)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return cmd
}

func newServiceUsersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <service-user-id>",
		Short: "Get service user details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			s, err := c.ServiceUsersGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), s)
			}
			rows := [][]string{
				{"ID", s.ID},
				{"Name", s.Name},
				{"Description", orDash(s.Description)},
				{"Scope", string(s.Scope)},
				{"Created At", orDash(s.CreatedAt)},
				{"Last Activation", orDash(s.LastActivation)},
				{"Token Expires", orDash(s.APIToken.ExpiresAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}

func newServiceUsersExportCmd() *cobra.Command {
	var siteIDs, accountIDs, roleIDs []string
	var query, outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export service users",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ServiceUserListParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				RoleIDs:    roleIDs,
				Query:      query,
			}
			data, err := c.ServiceUsersExport(cmd.Context(), params)
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
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&roleIDs, "role-id", nil, "filter by RBAC role ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search (name, description)")
	cmd.Flags().StringVar(&outFile, "out", "", "write export to file (default: stdout)")
	return cmd
}
