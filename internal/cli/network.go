package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// The network group operates SentinelOne Network Quarantine rules. Network
// Quarantine shares the firewall-control endpoint family with the firewall
// group, so the command surface mirrors firewall's verbs plus the operations
// only exposed for network quarantine (configuration, set-location, move, tags).
func newNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage network quarantine rules",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newNetworkListCmd())
	cmd.AddCommand(newNetworkGetCmd())
	cmd.AddCommand(newNetworkDeleteCmd())
	cmd.AddCommand(newNetworkEnableCmd())
	cmd.AddCommand(newNetworkDisableCmd())
	cmd.AddCommand(newNetworkReorderCmd())
	cmd.AddCommand(newNetworkCopyCmd())
	cmd.AddCommand(newNetworkMoveCmd())
	cmd.AddCommand(newNetworkSetLocationCmd())
	cmd.AddCommand(newNetworkTagsCmd())
	cmd.AddCommand(newNetworkConfigurationCmd())
	cmd.AddCommand(newNetworkProtocolsCmd())
	cmd.AddCommand(newNetworkExportCmd())
	cmd.AddCommand(newNetworkImportCmd())
	addNetworkSyncCmds(cmd)
	return cmd
}

func newNetworkGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <rule-id>",
		Short: "Get a network quarantine rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			r, err := c.NetworkQuarantineGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), r)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"ID", r.ID},
				{"Name", r.Name},
				{"Status", string(r.Status)},
				{"Action", string(r.Action)},
				{"Direction", string(r.Direction)},
				{"Protocol", orDash(r.Protocol)},
				{"OS", orDash(r.OSType)},
				{"Description", orDash(r.Description)},
			})
			return nil
		},
	}
	return markJSON(cmd)
}

func newNetworkListCmd() *cobra.Command {
	var siteIDs []string
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List network quarantine rules",
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
				rules, total, err = fetchAllREST("network quarantine rule", func(cur string) ([]mgmt.FirewallRule, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.NetworkQuarantineList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				rules, pag, err = c.NetworkQuarantineList(cmd.Context(), params)
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
					r.ID, r.Name, string(r.Direction),
					r.Protocol, string(r.Action), string(r.Status),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, rules, len(rules), total, "network quarantine rule", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return markJSON(cmd)
}

func newNetworkProtocolsCmd() *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:   "protocols",
		Short: "List available network quarantine protocols",
		Long:  `Show protocols that can be used in network quarantine rules.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.FirewallProtocolListParams{Query: query, Limit: 1000}
			protocols, _, err := c.NetworkQuarantineProtocolsList(cmd.Context(), params)
			if err != nil {
				return err
			}
			headers := []string{"Value", "Name"}
			rows := make([][]string, len(protocols))
			for i, p := range protocols {
				rows[i] = []string{p.Value, p.Name}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, protocols, len(protocols), len(protocols), "protocol", true)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search protocols")
	return markJSON(cmd)
}

func newNetworkExportCmd() *cobra.Command {
	var siteIDs []string
	var outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export network quarantine rules to a JSON file",
		Long: `Export network quarantine rules from a scope to a JSON file.
The exported file can be imported into another scope with "network import".`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			data, err := c.NetworkQuarantineExport(cmd.Context(), &mgmt.FirewallRuleListParams{SiteIDs: siteIDs})
			if err != nil {
				return err
			}
			if outFile == "-" {
				_, err = cmd.OutOrStdout().Write(data)
				return err
			}
			if err := os.WriteFile(outFile, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Exported network quarantine rules to %s\n", outFile)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope: site IDs")
	cmd.Flags().StringVar(&outFile, "out", "network-quarantine-rules.json", "output file (use - for stdout)")
	return cmd
}

func newNetworkImportCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import network quarantine rules from a JSON file",
		Long: `Import network quarantine rules from a previously exported JSON file into a scope.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]
			data, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("read %s: %w", filename, err)
			}
			return guard(cmd.OutOrStdout(), "network import",
				"import network quarantine rules from "+filename,
				filename, yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					scope := mgmt.FirewallImportScope{
						SiteIDs:    siteIDs,
						AccountIDs: accountIDs,
						GroupIDs:   groupIDs,
					}
					if len(siteIDs) == 0 && len(accountIDs) == 0 && len(groupIDs) == 0 {
						scope.Tenant = true
					}
					if err := c.NetworkQuarantineImport(cmd.Context(), scope, filename, data); err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]any{"success": true, "file": filename})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Imported network quarantine rules from %s\n", filename)
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "target site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "target account IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "target group IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}
