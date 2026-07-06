package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newUsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newUsersListCmd())
	cmd.AddCommand(newUsersGetCmd())
	cmd.AddCommand(newUsersUpdateCmd())
	cmd.AddCommand(newUsersDeleteCmd())
	cmd.AddCommand(newUsersGenerateTokenCmd())
	cmd.AddCommand(newUsersRevokeTokenCmd())
	cmd.AddCommand(newUsersTokenDetailsCmd())
	cmd.AddCommand(newUsers2FACmd())
	return cmd
}

func newUsersListCmd() *cobra.Command {
	var query, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.UserListParams{
				Query:     query,
				Limit:     limit,
				Cursor:    cursor,
				SortBy:    sortBy,
				SortOrder: sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var users []mgmt.User
			var total int

			if all {
				users, total, err = fetchAllREST("user", func(cur string) ([]mgmt.User, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.UsersList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				users, pag, err = c.UsersList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Email", "Scope", "Source"}
			rows := make([][]string, len(users))
			for i, u := range users {
				rows[i] = []string{
					u.ID, u.FullName, u.Email, u.Scope, u.Source,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, users, len(users), total, "user", all)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. fullName, email)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return markJSON(cmd)
}

func newUsersGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <user-id>",
		Short: "Get user details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			u, err := c.UsersGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), u)
			}
			rows := [][]string{
				{"ID", u.ID},
				{"Name", u.FullName},
				{"Email", u.Email},
				{"Scope", u.Scope},
				{"Source", u.Source},
				{"2FA", boolIcon(u.TwoFaEnabled)},
				{"Joined", u.DateJoined},
				{"Last Login", u.LastLogin},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}
