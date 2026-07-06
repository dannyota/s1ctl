package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newUpgradePoliciesPackagesCmd() *cobra.Command {
	var (
		scopeLevel string
		scopeID    string
		osType     string
		query      string
	)

	cmd := &cobra.Command{
		Use:   "packages",
		Short: "List available upgrade packages",
		Long: `List agent packages available for upgrade policies.

Scope levels: account, group, site, tenant
OS types: linux, macos, windows

Each package may include multiple file variants. Use the file ID
when creating an upgrade policy (--file-id).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if scopeLevel == "" {
				return fmt.Errorf("--scope-level is required (account, group, site, tenant)")
			}
			if osType == "" {
				return fmt.Errorf("--os-type is required (linux, macos, windows)")
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.UpgradePackageListParams{
				ScopeLevel:          scopeLevel,
				ScopeID:             scopeID,
				OSType:              osType,
				DisplayNameContains: query,
			}
			pkgs, err := c.UpgradePackagesList(cmd.Context(), params)
			if err != nil {
				return err
			}

			headers := []string{"Version", "Display Name", "File ID", "File Name"}
			var rows [][]string
			for _, p := range pkgs {
				ver := p.Major + "." + p.Minor + "." + p.Build
				if len(p.FileNames) == 0 {
					rows = append(rows, []string{ver, p.DisplayName, "-", "-"})
					continue
				}
				for _, f := range p.FileNames {
					rows = append(rows, []string{ver, p.DisplayName, f.ID, truncate(f.Name, 50)})
				}
			}
			if outputFormat == "table" && len(rows) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No packages found.")
				return nil
			}

			if err := printOutput(cmd.OutOrStdout(), headers, rows, pkgs, len(pkgs), len(pkgs), "package", true); err != nil {
				return err
			}

			if outputFormat == "table" {
				var fileIDs []string
				for _, p := range pkgs {
					for _, f := range p.FileNames {
						fileIDs = append(fileIDs, f.ID)
					}
				}
				if len(fileIDs) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "\nUse --file-id with 'upgrade-policies create' to target a package.\nExample: s1ctl upgrade-policies create --file-id %s ...\n",
						strings.Split(fileIDs[0], ",")[0])
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (account, group, site, tenant) [required]")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope ID")
	cmd.Flags().StringVar(&osType, "os-type", "", "OS type (linux, macos, windows) [required]")
	cmd.Flags().StringVar(&query, "query", "", "filter by display name (partial match)")
	return markJSON(cmd)
}
