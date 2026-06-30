package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Manage accounts",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAccountsListCmd())
	cmd.AddCommand(newAccountsGetCmd())
	return cmd
}

func newAccountsListCmd() *cobra.Command {
	var states []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			accounts, pag, err := c.AccountsList(cmd.Context(), &mgmt.AccountListParams{
				States: states,
				Query:  query,
				Limit:  limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(accounts)
			}
			var rows [][]string
			for _, a := range accounts {
				rows = append(rows, []string{
					a.ID, a.Name, a.State, a.AccountType,
					fmt.Sprintf("%d / %d", a.ActiveLicenses, a.TotalLicenses),
					fmt.Sprint(a.NumberOfSites),
				})
			}
			printTable([]string{"ID", "Name", "State", "Type", "Licenses", "Sites"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "account"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&states, "state", nil, "filter by state")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}

func newAccountsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <account-id>",
		Short: "Get account details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			a, err := c.AccountsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(a)
			}
			rows := [][]string{
				{"ID", a.ID},
				{"Name", a.Name},
				{"State", a.State},
				{"Type", a.AccountType},
				{"Licenses", fmt.Sprintf("%d / %d", a.ActiveLicenses, a.TotalLicenses)},
				{"Sites", fmt.Sprint(a.NumberOfSites)},
				{"Expiration", a.Expiration},
				{"Created", a.CreatedAt},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}
