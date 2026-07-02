package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newUpgradePoliciesCreateCmd() *cobra.Command {
	var (
		name         string
		description  string
		osType       string
		scopeLevel   string
		scopeID      string
		isActive     bool
		allEndpoints bool
		maxRetries   int
		fileID       string
		major        string
		minor        string
		build        string
		tags         []string
		yes          bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an upgrade policy",
		Long: `Create a new agent auto-upgrade policy.

Scope levels: account, group, site, tenant
OS types: linux, macos, windows

Use "upgrade-policies packages" to find available package versions and file IDs.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if osType == "" {
				return fmt.Errorf("--os-type is required (linux, macos, windows)")
			}
			if scopeLevel == "" {
				return fmt.Errorf("--scope-level is required (account, group, site, tenant)")
			}
			if fileID == "" {
				return fmt.Errorf("--file-id is required (use 'upgrade-policies packages' to find IDs)")
			}

			data := mgmt.UpgradePolicyCreate{
				Name:         name,
				Description:  description,
				OSType:       mgmt.UpgradePolicyOSType(osType),
				ScopeLevel:   mgmt.UpgradePolicyScopeLevel(scopeLevel),
				ScopeID:      scopeID,
				IsActive:     isActive,
				AllEndpoints: allEndpoints,
				MaxRetries:   maxRetries,
				Package: mgmt.UpgradePolicyPkg{
					FileID: fileID,
					Major:  major,
					Minor:  minor,
					Build:  build,
				},
				Tags: tags,
			}

			return guard(cmd.OutOrStdout(), "upgrade-policies create", "create upgrade policy "+name+" ("+osType+", "+scopeLevel+")", name, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.UpgradePoliciesCreate(cmd.Context(), data); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{
						"status": "created",
						"name":   name,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created upgrade policy %q\n", name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "policy name (required)")
	cmd.Flags().StringVar(&description, "description", "", "policy description")
	cmd.Flags().StringVar(&osType, "os-type", "", "OS type: linux, macos, windows (required)")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level: account, group, site, tenant (required)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope ID")
	cmd.Flags().BoolVar(&isActive, "active", false, "activate the policy immediately")
	cmd.Flags().BoolVar(&allEndpoints, "all-endpoints", true, "apply to all endpoints (set false with tags)")
	cmd.Flags().IntVar(&maxRetries, "max-retries", 5, "max upgrade retries on failure")
	cmd.Flags().StringVar(&fileID, "file-id", "", "package file ID (required; see 'upgrade-policies packages')")
	cmd.Flags().StringVar(&major, "major", "", "package major version")
	cmd.Flags().StringVar(&minor, "minor", "", "package minor version")
	cmd.Flags().StringVar(&build, "build", "", "package build version")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "endpoint tags (when --all-endpoints=false)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newUpgradePoliciesDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <policy-id>",
		Short: "Delete an upgrade policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return guard(cmd.OutOrStdout(), "upgrade-policies delete", "delete upgrade policy "+id, id, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.UpgradePoliciesDelete(cmd.Context(), id); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted upgrade policy %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
