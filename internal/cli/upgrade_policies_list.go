package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newUpgradePoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade-policies",
		Aliases: []string{"up"},
		Short:   "Agent auto-upgrade policies",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newUpgradePoliciesListCmd())
	cmd.AddCommand(newUpgradePoliciesGetCmd())
	cmd.AddCommand(newUpgradePoliciesCreateCmd())
	cmd.AddCommand(newUpgradePoliciesDeleteCmd())
	cmd.AddCommand(newUpgradePoliciesPackagesCmd())
	return cmd
}

func newUpgradePoliciesListCmd() *cobra.Command {
	var (
		scopeLevel string
		scopeID    string
		osType     string
		limit      int
		skip       int
		sortBy     string
		sortOrder  string
		all        bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List upgrade policies",
		Long: `List agent auto-upgrade policies for a given scope and OS type.

Scope levels: account, group, site, tenant
OS types: linux, macos, windows`,
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
			if limit == 0 {
				limit = defaultPageSize
			}
			if sortBy == "" {
				sortBy = "priority"
			}
			if sortOrder == "" {
				sortOrder = "asc"
			}

			params := &mgmt.UpgradePolicyListParams{
				ScopeLevel: scopeLevel,
				ScopeID:    scopeID,
				OSType:     osType,
				Limit:      limit,
				Skip:       skip,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
			}

			var policies []mgmt.UpgradePolicy
			var total int

			if all {
				policies, total, err = fetchAllUpgradePolicies(cmd, c, params)
			} else {
				policies, total, err = c.UpgradePoliciesList(cmd.Context(), params)
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "OS", "Active", "Package", "Priority", "Scope"}
			rows := make([][]string, len(policies))
			for i, p := range policies {
				pkg := p.Package.Major + "." + p.Package.Minor + "." + p.Package.Build
				rows[i] = []string{
					p.ID, truncate(p.Name, 40), p.OSType,
					boolIcon(p.IsActive), pkg,
					strconv.Itoa(p.Priority), p.ScopeLevel,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, policies, len(policies), total, "upgrade policy", all)
		},
	}
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (account, group, site, tenant) [required]")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope ID")
	cmd.Flags().StringVar(&osType, "os-type", "", "OS type (linux, macos, windows) [required]")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().IntVar(&skip, "skip", 0, "skip first N results")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (default: priority)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc; default: asc)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	return cmd
}

func newUpgradePoliciesGetCmd() *cobra.Command {
	var (
		scopeLevel string
		scopeID    string
		osType     string
	)

	cmd := &cobra.Command{
		Use:   "get <policy-id>",
		Short: "Get upgrade policy details",
		Long: `Get details for a single upgrade policy by ID.

The API requires scope and OS filters even for a single lookup.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			targetID := args[0]

			// The API has no GET-by-ID endpoint; list and filter client-side.
			policies, _, err := c.UpgradePoliciesList(cmd.Context(), &mgmt.UpgradePolicyListParams{
				ScopeLevel: scopeLevel,
				ScopeID:    scopeID,
				OSType:     osType,
				Limit:      200,
				SortBy:     "priority",
				SortOrder:  "asc",
			})
			if err != nil {
				return err
			}

			var found *mgmt.UpgradePolicy
			for i := range policies {
				if policies[i].ID == targetID {
					found = &policies[i]
					break
				}
			}
			if found == nil {
				return fmt.Errorf("upgrade policy %s not found", targetID)
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), found)
			}

			pkg := found.Package.Major + "." + found.Package.Minor + "." + found.Package.Build
			tags := "-"
			if len(found.Tags) > 0 {
				tags = fmt.Sprintf("%v", found.Tags)
			}
			rows := [][]string{
				{"ID", found.ID},
				{"Name", found.Name},
				{"Description", orDash(found.Description)},
				{"OS", found.OSType},
				{"Active", boolIcon(found.IsActive)},
				{"Scheduled", boolIcon(found.IsScheduled)},
				{"All Endpoints", boolIcon(found.AllEndpoints)},
				{"Package", pkg},
				{"Package File", orDash(found.Package.FileID)},
				{"Priority", strconv.Itoa(found.Priority)},
				{"Max Retries", strconv.Itoa(found.MaxRetries)},
				{"Scope", found.ScopeLevel},
				{"Scope ID", orDash(found.ScopeID)},
				{"Tags", tags},
				{"Activated", orDash(found.ActivatedAt)},
				{"Created", orDash(found.CreatedAt)},
				{"Updated", orDash(found.UpdatedAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (account, group, site, tenant) [required]")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope ID")
	cmd.Flags().StringVar(&osType, "os-type", "", "OS type (linux, macos, windows) [required]")
	return cmd
}

// fetchAllUpgradePolicies pages through all upgrade policies using skip/limit.
func fetchAllUpgradePolicies(cmd *cobra.Command, c *mgmt.Client, params *mgmt.UpgradePolicyListParams) ([]mgmt.UpgradePolicy, int, error) {
	var all []mgmt.UpgradePolicy
	var total int
	params.Skip = 0
	for {
		items, t, err := c.UpgradePoliciesList(cmd.Context(), params)
		if err != nil {
			clearProgress()
			return nil, 0, err
		}
		total = t
		all = append(all, items...)
		printProgress("upgrade policy", len(all), total)
		if len(all) >= total || len(items) == 0 {
			break
		}
		params.Skip = len(all)
	}
	clearProgress()
	return all, total, nil
}
