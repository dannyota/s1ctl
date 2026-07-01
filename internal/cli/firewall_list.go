package cli

import (
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
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List firewall rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.FirewallRuleListParams{
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
				Cursor:  cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var rules []mgmt.FirewallRule
			var total int

			if all {
				rules, total, err = fetchAllREST("rule", func(cur string) ([]mgmt.FirewallRule, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.FirewallRulesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				rules, pag, err = c.FirewallRulesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Direction", "Protocol", "Action", "Status"}
			rows := make([][]string, len(rules))
			for i, r := range rules {
				rows[i] = []string{
					r.ID, r.Name, r.Direction,
					r.Protocol, r.Action, r.Status,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, rules, len(rules), total, "rule", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}
