package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newDeviceControlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devicecontrol",
		Short: "Device control rules",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDeviceControlListCmd())
	return cmd
}

func newDeviceControlListCmd() *cobra.Command {
	var siteIDs []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List device control rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			rules, pag, err := c.DeviceRulesList(cmd.Context(), &mgmt.DeviceRuleListParams{
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
					r.ID, r.RuleName, r.DeviceClass,
					r.Action, r.Status,
				})
			}
			printTable([]string{"ID", "Name", "Class", "Action", "Status"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "rule"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}
