package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newFirewallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firewall",
		Short: "Firewall control rules",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newFirewallListCmd())
	return cmd
}

func newFirewallListCmd() *cobra.Command {
	var siteIDs []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List firewall rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			rules, pag, err := c.FirewallRulesList(cmd.Context(), &mgmt.FirewallRuleListParams{
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(rules)
			}
			var rows [][]string
			for _, r := range rules {
				rows = append(rows, []string{
					r.ID, r.Name, r.Direction,
					r.Protocol, r.Action, r.Status,
				})
			}
			printTable([]string{"ID", "Name", "Direction", "Protocol", "Action", "Status"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "rule"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}
