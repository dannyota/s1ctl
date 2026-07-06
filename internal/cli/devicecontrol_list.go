package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newDeviceControlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devicecontrol",
		Short: "Manage device control rules",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDeviceControlListCmd())
	cmd.AddCommand(newDeviceControlGetCmd())
	cmd.AddCommand(newDeviceControlDeleteCmd())
	cmd.AddCommand(newDeviceControlEnableCmd())
	cmd.AddCommand(newDeviceControlDisableCmd())
	cmd.AddCommand(newDeviceControlReorderCmd())
	cmd.AddCommand(newDeviceControlCopyCmd())
	cmd.AddCommand(newDeviceControlEventsCmd())
	addDeviceControlSyncCmds(cmd)
	return cmd
}

func newDeviceControlGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <rule-id>",
		Short: "Get a device control rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			r, err := c.DeviceRulesGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), r)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"ID", r.ID},
				{"RuleName", r.RuleName},
				{"Status", string(r.Status)},
				{"Action", string(r.Action)},
				{"Interface", string(r.Interface)},
				{"RuleType", string(r.RuleType)},
				{"AccessPermission", string(r.AccessPermission)},
				{"DeviceClass", orDash(r.DeviceClass)},
			})
			return nil
		},
	}
}

func newDeviceControlListCmd() *cobra.Command {
	var siteIDs []string
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List device control rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.DeviceRuleListParams{
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
				Cursor:  cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var rules []mgmt.DeviceRule
			var total int

			if all {
				rules, total, err = fetchAllREST("device rule", func(cur string) ([]mgmt.DeviceRule, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.DeviceRulesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				rules, pag, err = c.DeviceRulesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Class", "Action", "Status"}
			rows := make([][]string, len(rules))
			for i, r := range rules {
				rows[i] = []string{
					r.ID, r.RuleName, r.DeviceClass,
					string(r.Action), string(r.Status),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, rules, len(rules), total, "device rule", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}
