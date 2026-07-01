package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newFirewallExportCmd() *cobra.Command {
	var siteIDs []string
	var outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export firewall rules to a JSON file",
		Long: `Export firewall rules from a scope to a JSON file.
The exported file can be imported into another scope with "firewall import".`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.FirewallRuleListParams{
				SiteIDs: siteIDs,
			}

			data, err := c.FirewallRulesExport(cmd.Context(), params)
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
			fmt.Fprintf(cmd.OutOrStdout(), "Exported firewall rules to %s\n", outFile)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope: site IDs")
	cmd.Flags().StringVar(&outFile, "out", "firewall-rules.json", "output file (use - for stdout)")
	return cmd
}

func newFirewallImportCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import firewall rules from a JSON file",
		Long: `Import firewall rules from a previously exported JSON file into a scope.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]

			data, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("read %s: %w", filename, err)
			}

			return guard(cmd.OutOrStdout(), "firewall import",
				"import firewall rules from "+filename,
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
					if err := c.FirewallRulesImport(cmd.Context(), scope, filename, data); err != nil {
						return err
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Imported firewall rules from %s\n", filename)
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "target site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "target account IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "target group IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}
