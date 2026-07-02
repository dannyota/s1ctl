package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// upgradePolicyFlags holds the shared flag set for creating and updating
// upgrade policies. Both commands register the same flags, run the same
// validation, and assemble the same mgmt.UpgradePolicyCreate body.
type upgradePolicyFlags struct {
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
}

func (f *upgradePolicyFlags) register(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.name, "name", "", "policy name (required)")
	cmd.Flags().StringVar(&f.description, "description", "", "policy description")
	cmd.Flags().StringVar(&f.osType, "os-type", "", "OS type: linux, macos, windows (required)")
	cmd.Flags().StringVar(&f.scopeLevel, "scope-level", "", "scope level: account, group, site, tenant (required)")
	cmd.Flags().StringVar(&f.scopeID, "scope-id", "", "scope ID")
	cmd.Flags().BoolVar(&f.isActive, "active", false, "activate the policy immediately")
	cmd.Flags().BoolVar(&f.allEndpoints, "all-endpoints", true, "apply to all endpoints (set false with tags)")
	cmd.Flags().IntVar(&f.maxRetries, "max-retries", 5, "max upgrade retries on failure")
	cmd.Flags().StringVar(&f.fileID, "file-id", "", "package file ID (required; see 'upgrade-policies packages')")
	cmd.Flags().StringVar(&f.major, "major", "", "package major version")
	cmd.Flags().StringVar(&f.minor, "minor", "", "package minor version")
	cmd.Flags().StringVar(&f.build, "build", "", "package build version")
	cmd.Flags().StringSliceVar(&f.tags, "tag", nil, "endpoint tags (when --all-endpoints=false)")
	cmd.Flags().BoolVar(&f.yes, "yes", false, "apply the action (default: dry-run)")
}

func (f *upgradePolicyFlags) validate() error {
	if f.name == "" {
		return fmt.Errorf("--name is required")
	}
	if f.osType == "" {
		return fmt.Errorf("--os-type is required (linux, macos, windows)")
	}
	if f.scopeLevel == "" {
		return fmt.Errorf("--scope-level is required (account, group, site, tenant)")
	}
	if f.fileID == "" {
		return fmt.Errorf("--file-id is required (use 'upgrade-policies packages' to find IDs)")
	}
	return nil
}

func (f *upgradePolicyFlags) data() mgmt.UpgradePolicyCreate {
	return mgmt.UpgradePolicyCreate{
		Name:         f.name,
		Description:  f.description,
		OSType:       mgmt.UpgradePolicyOSType(f.osType),
		ScopeLevel:   mgmt.UpgradePolicyScopeLevel(f.scopeLevel),
		ScopeID:      f.scopeID,
		IsActive:     f.isActive,
		AllEndpoints: f.allEndpoints,
		MaxRetries:   f.maxRetries,
		Package: mgmt.UpgradePolicyPkg{
			FileID: f.fileID,
			Major:  f.major,
			Minor:  f.minor,
			Build:  f.build,
		},
		Tags: f.tags,
	}
}

func newUpgradePoliciesCreateCmd() *cobra.Command {
	var f upgradePolicyFlags

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an upgrade policy",
		Long: `Create a new agent auto-upgrade policy.

Scope levels: account, group, site, tenant
OS types: linux, macos, windows

Use "upgrade-policies packages" to find available package versions and file IDs.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := f.validate(); err != nil {
				return err
			}

			data := f.data()

			return guard(cmd.OutOrStdout(), "upgrade-policies create", "create upgrade policy "+f.name+" ("+f.osType+", "+f.scopeLevel+")", f.name, f.yes, func() error {
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
						"name":   f.name,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created upgrade policy %q\n", f.name)
				return nil
			})
		},
	}
	f.register(cmd)
	return cmd
}

func newUpgradePoliciesUpdateCmd() *cobra.Command {
	var f upgradePolicyFlags

	cmd := &cobra.Command{
		Use:   "update <policy-id>",
		Short: "Update an upgrade policy",
		Long: `Update an existing agent auto-upgrade policy.

The full policy body is sent, so provide every flag as with "create".

Scope levels: account, group, site, tenant
OS types: linux, macos, windows

Use "upgrade-policies packages" to find available package versions and file IDs.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := f.validate(); err != nil {
				return err
			}

			data := f.data()

			return guard(cmd.OutOrStdout(), "upgrade-policies update", "update upgrade policy "+args[0], args[0], f.yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.UpgradePoliciesUpdate(cmd.Context(), args[0], data); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated upgrade policy %s\n", args[0])
				return nil
			})
		},
	}
	f.register(cmd)
	return cmd
}

func newUpgradePolicyToggleCmd(verb string, call func(*mgmt.Client, context.Context, string) error) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <policy-id>",
		Short: strings.ToUpper(verb[:1]) + verb[1:] + " an upgrade policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "upgrade-policies "+verb, verb+" upgrade policy "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := call(c, cmd.Context(), args[0]); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": verb + "d", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%sd upgrade policy %s\n", strings.ToUpper(verb[:1])+verb[1:], args[0])
				return nil
			})
		},
	}
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
