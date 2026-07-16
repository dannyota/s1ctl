package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newMarketplaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "marketplace",
		Short: "Manage Singularity Marketplace applications",
		Long: `Browse the Singularity Marketplace catalog, install integrations,
configure installed applications, and manage their lifecycle.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newMarketplaceCatalogCmd())
	cmd.AddCommand(newMarketplaceCatalogConfigCmd())
	cmd.AddCommand(newMarketplaceListCmd())
	cmd.AddCommand(newMarketplaceConfigCmd())
	cmd.AddCommand(newMarketplaceLogCmd())
	cmd.AddCommand(newMarketplaceInstallCmd())
	cmd.AddCommand(newMarketplaceUpdateCmd())
	cmd.AddCommand(newMarketplaceDeleteCmd())
	cmd.AddCommand(newMarketplaceEnableCmd())
	cmd.AddCommand(newMarketplaceDisableCmd())
	return cmd
}

// --- catalog ---

func newMarketplaceCatalogCmd() *cobra.Command {
	var nameContains, categoryContains, query, sortBy, sortOrder string
	var categoryIDs []string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "List marketplace catalog applications",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.MarketplaceCatalogListParams{
				NameContains:     nameContains,
				CategoryContains: categoryContains,
				Query:            query,
				CategoryIDs:      categoryIDs,
				Limit:            limit,
				SortBy:           sortBy,
				SortOrder:        sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.MarketplaceCatalogItem
			var total int

			if all {
				for {
					page, pag, err := c.MarketplaceCatalogList(cmd.Context(), params)
					if err != nil {
						return err
					}
					items = append(items, page...)
					if pag != nil {
						total = pag.TotalItems
						params.Cursor = pag.NextCursor
					}
					if len(page) < params.Limit || params.Cursor == "" {
						break
					}
				}
			} else {
				page, pag, err := c.MarketplaceCatalogList(cmd.Context(), params)
				if err != nil {
					return err
				}
				items = page
				if pag != nil {
					total = pag.TotalItems
				}
			}

			headers := []string{"ID", "Name", "Category", "Installed"}
			rows := make([][]string, len(items))
			for i, item := range items {
				rows[i] = []string{item.ID, item.Name, item.Category, boolIcon(item.Installed)}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "catalog app", all)
		},
	}
	cmd.Flags().StringVar(&nameContains, "name", "", "filter by name (contains)")
	cmd.Flags().StringVar(&categoryContains, "category", "", "filter by category (contains)")
	cmd.Flags().StringVar(&query, "query", "", "free-text search")
	cmd.Flags().StringSliceVar(&categoryIDs, "category-id", nil, "filter by category ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order (asc, desc)")
	return markJSON(cmd)
}

// --- catalog-config ---

func newMarketplaceCatalogConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog-config CATALOG_ID",
		Short: "Show configuration fields for a catalog application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			data, err := c.MarketplaceCatalogConfig(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), data)
		},
	}
	return markJSON(cmd)
}

// --- list ---

func newMarketplaceListCmd() *cobra.Command {
	var nameContains, creatorContains, query, catalogID, sortBy, sortOrder string
	var siteIDs, accountIDs []string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed marketplace applications",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.MarketplaceAppListParams{
				ApplicationCatalogID: catalogID,
				NameContains:         nameContains,
				CreatorContains:      creatorContains,
				Query:                query,
				SiteIDs:              siteIDs,
				AccountIDs:           accountIDs,
				Limit:                limit,
				SortBy:               sortBy,
				SortOrder:            sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.MarketplaceApp
			var total int

			if all {
				for {
					page, pag, err := c.MarketplaceAppList(cmd.Context(), params)
					if err != nil {
						return err
					}
					items = append(items, page...)
					if pag != nil {
						total = pag.TotalItems
						params.Cursor = pag.NextCursor
					}
					if len(page) < params.Limit || params.Cursor == "" {
						break
					}
				}
			} else {
				page, pag, err := c.MarketplaceAppList(cmd.Context(), params)
				if err != nil {
					return err
				}
				items = page
				if pag != nil {
					total = pag.TotalItems
				}
			}

			headers := []string{"ID", "InstanceName", "CatalogID", "Status", "Installed"}
			var rows [][]string
			for _, item := range items {
				if len(item.Scopes) == 0 {
					rows = append(rows, []string{"", "", item.ApplicationCatalogID, "", orDash(item.LastInstalledAt)})
					continue
				}
				for _, sc := range item.Scopes {
					rows = append(rows, []string{sc.ID, sc.ApplicationInstanceName, item.ApplicationCatalogID, sc.Status, orDash(item.LastInstalledAt)})
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "installed app", all)
		},
	}
	cmd.Flags().StringVar(&catalogID, "catalog-id", "", "filter by catalog application ID")
	cmd.Flags().StringVar(&nameContains, "name", "", "filter by name (contains)")
	cmd.Flags().StringVar(&creatorContains, "creator", "", "filter by creator (contains)")
	cmd.Flags().StringVar(&query, "query", "", "free-text search")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order (asc, desc)")
	return markJSON(cmd)
}

// --- config ---

func newMarketplaceConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config APP_ID",
		Short: "Show configuration for an installed application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			data, err := c.MarketplaceAppConfig(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), data)
		},
	}
	return markJSON(cmd)
}

// --- log ---

func newMarketplaceLogCmd() *cobra.Command {
	var onlyErrors bool

	cmd := &cobra.Command{
		Use:   "log APP_ID",
		Short: "Show log entries for an installed application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			var errPtr *bool
			if cmd.Flags().Changed("only-errors") {
				errPtr = &onlyErrors
			}
			entries, err := c.MarketplaceAppLog(cmd.Context(), args[0], errPtr)
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), entries)
		},
	}
	cmd.Flags().BoolVar(&onlyErrors, "only-errors", false, "show only error entries")
	return markJSON(cmd)
}

// --- install ---

func newMarketplaceInstallCmd() *cobra.Command {
	var catalogID, name string
	var configs, siteIDs, accountIDs, groupIDs []string
	var tenant, yes bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a marketplace application",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if catalogID == "" {
				return fmt.Errorf("--catalog-id is required")
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			cfgs, err := parseMarketplaceConfigs(configs)
			if err != nil {
				return err
			}
			input := &mgmt.MarketplaceInstallInput{}
			input.Data.Name = name
			input.Data.Configurations = cfgs
			input.Filter.ApplicationCatalogID = catalogID
			input.Filter.SiteIDs = siteIDs
			input.Filter.AccountIDs = accountIDs
			input.Filter.GroupIDs = groupIDs
			if tenant {
				input.Filter.Tenant = &tenant
			}

			return guard(cmd.OutOrStdout(), "marketplace install",
				fmt.Sprintf("install %q from catalog %s", name, catalogID),
				catalogID, yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					return c.MarketplaceInstall(cmd.Context(), input)
				})
		},
	}
	cmd.Flags().StringVar(&catalogID, "catalog-id", "", "catalog application ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "instance name (required)")
	cmd.Flags().StringSliceVar(&configs, "config", nil, "configuration (id=value, repeatable)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope to site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope to account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope to group ID")
	cmd.Flags().BoolVar(&tenant, "tenant", false, "scope to tenant")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the change (default: dry-run)")
	return markJSON(cmd)
}

// --- update ---

func newMarketplaceUpdateCmd() *cobra.Command {
	var id, name string
	var configs, siteIDs, accountIDs, groupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an installed marketplace application",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if id == "" {
				return fmt.Errorf("--id is required")
			}
			cfgs, err := parseMarketplaceConfigs(configs)
			if err != nil {
				return err
			}
			input := &mgmt.MarketplaceUpdateInput{}
			if name != "" {
				input.Data.NameMap = map[string]string{id: name}
			}
			input.Data.Configurations = cfgs
			input.Filter.IDs = []string{id}
			input.Filter.SiteIDs = siteIDs
			input.Filter.AccountIDs = accountIDs
			input.Filter.GroupIDs = groupIDs

			return guard(cmd.OutOrStdout(), "marketplace update",
				fmt.Sprintf("update application %s", id),
				id, yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					return c.MarketplaceUpdate(cmd.Context(), input)
				})
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "application ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "new instance name")
	cmd.Flags().StringSliceVar(&configs, "config", nil, "configuration (id=value, repeatable)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope to site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope to account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope to group ID")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the change (default: dry-run)")
	return markJSON(cmd)
}

// --- delete ---

func newMarketplaceDeleteCmd() *cobra.Command {
	var id string
	var siteIDs, accountIDs, groupIDs []string
	var tenant, yes bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an installed marketplace application",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if id == "" {
				return fmt.Errorf("--id is required")
			}
			filter := &mgmt.MarketplaceDeleteFilter{
				ID:         []string{id},
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
			}
			if tenant {
				filter.Tenant = &tenant
			}
			return guard(cmd.OutOrStdout(), "marketplace delete",
				fmt.Sprintf("delete application %s", id),
				id, yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					return c.MarketplaceDelete(cmd.Context(), filter)
				})
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "application ID (required)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope to site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope to account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope to group ID")
	cmd.Flags().BoolVar(&tenant, "tenant", false, "scope to tenant")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the change (default: dry-run)")
	return markJSON(cmd)
}

// --- enable ---

func newMarketplaceEnableCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "enable APP_ID",
		Short: "Enable an installed marketplace application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := &mgmt.MarketplaceScopeFilter{
				ApplicationID: args[0],
				SiteIDs:       siteIDs,
				AccountIDs:    accountIDs,
				GroupIDs:      groupIDs,
			}
			return guard(cmd.OutOrStdout(), "marketplace enable",
				fmt.Sprintf("enable application %s", args[0]),
				args[0], yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					return c.MarketplaceSetMode(cmd.Context(), "enable", filter)
				})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope to site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope to account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope to group ID")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the change (default: dry-run)")
	return markJSON(cmd)
}

// --- disable ---

func newMarketplaceDisableCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "disable APP_ID",
		Short: "Disable an installed marketplace application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := &mgmt.MarketplaceScopeFilter{
				ApplicationID: args[0],
				SiteIDs:       siteIDs,
				AccountIDs:    accountIDs,
				GroupIDs:      groupIDs,
			}
			return guard(cmd.OutOrStdout(), "marketplace disable",
				fmt.Sprintf("disable application %s", args[0]),
				args[0], yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					return c.MarketplaceSetMode(cmd.Context(), "disable", filter)
				})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope to site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope to account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope to group ID")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the change (default: dry-run)")
	return markJSON(cmd)
}

// --- helpers ---

func parseMarketplaceConfigs(raw []string) ([]mgmt.MarketplaceConfig, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	cfgs := make([]mgmt.MarketplaceConfig, len(raw))
	for i, s := range raw {
		k, v, ok := strings.Cut(s, "=")
		if !ok {
			return nil, fmt.Errorf("invalid config format %q: expected id=value", s)
		}
		cfgs[i] = mgmt.MarketplaceConfig{ID: k, Value: v}
	}
	return cfgs, nil
}
