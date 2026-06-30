package cli

import (
	"fmt"

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
	return cmd
}

func newUsersListCmd() *cobra.Command {
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			users, pag, err := c.UsersList(cmd.Context(), &mgmt.UserListParams{
				Query: query,
				Limit: limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(users)
			}
			var rows [][]string
			for _, u := range users {
				rows = append(rows, []string{
					u.ID, u.FullName, u.Email, u.Scope, u.Source,
				})
			}
			printTable([]string{"ID", "Name", "Email", "Scope", "Source"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "user"))
			return nil
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}

func newUsersGetCmd() *cobra.Command {
	return &cobra.Command{
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
			if jsonOutput {
				return printJSON(u)
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
}
