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
	cmd.AddCommand(newSitesGetCmd())
	return cmd
}

func newSitesListCmd() *cobra.Command {
	var accountIDs, states []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sites",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			sites, pag, err := c.SitesList(cmd.Context(), &mgmt.SiteListParams{
				AccountIDs: accountIDs,
				States:     states,
				Query:      query,
				Limit:      limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(sites)
			}
			var rows [][]string
			for _, s := range sites {
				rows = append(rows, []string{
					s.ID, s.Name, s.State, s.SiteType,
					fmt.Sprintf("%d / %d", s.ActiveLicenses, s.TotalLicenses),
				})
			}
			printTable([]string{"ID", "Name", "State", "Type", "Licenses"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "site"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&states, "state", nil, "filter by state")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
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
			if jsonOutput {
				return printJSON(s)
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
