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
	cmd.AddCommand(newAccountsCountCmd())
	cmd.AddCommand(newAccountsGetCmd())
	return cmd
}

func newAccountsListCmd() *cobra.Command {
	var states []string
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List accounts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AccountListParams{
				States: states,
				Query:  query,
				Limit:  limit,
				Cursor: cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var accounts []mgmt.Account
			var total int

			if all {
				accounts, total, err = fetchAllREST("account", func(cur string) ([]mgmt.Account, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.AccountsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				accounts, pag, err = c.AccountsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "State", "Type", "Licenses", "Sites"}
			rows := make([][]string, len(accounts))
			for i, a := range accounts {
				rows[i] = []string{
					a.ID, a.Name, a.State, a.AccountType,
					fmt.Sprintf("%d / %d", a.ActiveLicenses, a.TotalLicenses),
					fmt.Sprint(a.NumberOfSites),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, accounts, len(accounts), total, "account", all)
		},
	}
	cmd.Flags().StringSliceVar(&states, "state", nil, "filter by state")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
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
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), a)
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

func newAccountsCountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Count accounts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			count, err := c.AccountsCount(cmd.Context(), &mgmt.AccountListParams{})
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
	return cmd
}
