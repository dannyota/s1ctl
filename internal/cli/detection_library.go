package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newDetectionLibraryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "detection-library",
		Aliases: []string{"dl"},
		Short:   "Manage platform detection rules (detection library)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDLListCmd())
	cmd.AddCommand(newDLSurfacesCmd())
	cmd.AddCommand(newDLDataSourcesCmd())
	cmd.AddCommand(newDLEnableCmd())
	cmd.AddCommand(newDLDisableCmd())
	return cmd
}

func newDLListCmd() *cobra.Command {
	var severities, statuses, attackSurfaces, sources, categories, tags []string
	var nameContains, scopeLevel, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List platform detection rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.PlatformRuleListParams{
				Severities:     severities,
				Statuses:       statuses,
				AttackSurfaces: attackSurfaces,
				Sources:        sources,
				Categories:     categories,
				Tags:           tags,
				NameContains:   nameContains,
				ScopeLevel:     scopeLevel,
				Limit:          limit,
				Cursor:         cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var rules []mgmt.PlatformRule
			var total int

			if all {
				rules, total, err = fetchAllREST("rule", func(cur string) ([]mgmt.PlatformRule, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.PlatformRulesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				rules, pag, err = c.PlatformRulesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Status", "Severity", "Scope", "Surfaces", "Alerts"}
			rows := make([][]string, len(rules))
			for i, r := range rules {
				rows[i] = []string{
					r.ID,
					truncate(r.Name, 50),
					string(r.Status),
					string(r.Severity),
					string(r.ScopeLevel),
					joinTruncate(r.AttackSurfaces, 30),
					fmt.Sprintf("%d", r.GeneratedAlerts),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, rules, len(rules), total, "rule", all)
		},
	}
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (Info, Low, Medium, High, Critical)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by status (Active, Disabled, Activating, Disabling)")
	cmd.Flags().StringSliceVar(&attackSurfaces, "surface", nil, "filter by attack surface")
	cmd.Flags().StringSliceVar(&sources, "source", nil, "filter by data source")
	cmd.Flags().StringSliceVar(&categories, "category", nil, "filter by category (Events, Correlation, UEBAFirstSeen, Scheduled)")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "filter by tag")
	cmd.Flags().StringVar(&nameContains, "name", "", "filter by rule name (substring match)")
	cmd.Flags().StringVar(&scopeLevel, "scope", "", "filter by scope level (global, account, site, group)")
	cmd.Flags().IntVar(&limit, "limit", 0, fmt.Sprintf("max results per page (default %d)", defaultPageSize))
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}

func newDLSurfacesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "surfaces",
		Short: "List available detection surfaces",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			surfaces, err := c.DetectionSurfacesList(cmd.Context())
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), surfaces)
			}
			headers := []string{"Key", "Title"}
			rows := make([][]string, len(surfaces))
			for i, s := range surfaces {
				rows[i] = []string{s.Key, s.Title}
			}
			printTable(headers, rows)
			return nil
		},
	}
}

func newDLDataSourcesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "data-sources",
		Short: "List available detection data sources",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			sources, err := c.DetectionDataSourcesList(cmd.Context())
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), sources)
			}
			headers := []string{"Key", "Title"}
			rows := make([][]string, len(sources))
			for i, s := range sources {
				rows[i] = []string{s.Key, s.Title}
			}
			printTable(headers, rows)
			return nil
		},
	}
}

func newDLEnableCmd() *cobra.Command {
	var yes bool
	var scopeID, scopeLevel string

	cmd := &cobra.Command{
		Use:   "enable <rule-id>...",
		Short: "Enable platform detection rules",
		Long: `Enable one or more platform detection rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would enable %s. Pass --yes to apply.\n",
					pluralize(len(args), "rule"))
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			filter := mgmt.PlatformRuleActionFilter{
				PlatformRuleIDs: args,
				ScopeID:         scopeID,
				ScopeLevel:      scopeLevel,
			}
			affected, err := c.PlatformRulesEnable(cmd.Context(), filter)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Enabled %s\n", pluralize(affected, "rule"))
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "account, site, or group ID for scoped enable")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (global, account, site, group)")
	return cmd
}

func newDLDisableCmd() *cobra.Command {
	var yes bool
	var scopeID, scopeLevel string

	cmd := &cobra.Command{
		Use:   "disable <rule-id>...",
		Short: "Disable platform detection rules",
		Long: `Disable one or more platform detection rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would disable %s. Pass --yes to apply.\n",
					pluralize(len(args), "rule"))
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			filter := mgmt.PlatformRuleActionFilter{
				PlatformRuleIDs: args,
				ScopeID:         scopeID,
				ScopeLevel:      scopeLevel,
			}
			affected, err := c.PlatformRulesDisable(cmd.Context(), filter)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Disabled %s\n", pluralize(affected, "rule"))
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "account, site, or group ID for scoped disable")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (global, account, site, group)")
	return cmd
}

// joinTruncate joins string slice items with commas, truncating if too long.
func joinTruncate(items []string, max int) string {
	if len(items) == 0 {
		return "-"
	}
	result := items[0]
	for _, s := range items[1:] {
		next := result + ", " + s
		if len(next) > max {
			return result + ", ..."
		}
		result = next
	}
	return result
}
