package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newSettingsOverridesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overrides",
		Short: "Manage agent config overrides",
		Long: `Manage configuration overrides that change agent behavior.

Config overrides are powerful: they override the agent's configuration at
the selected scope (tenant, account, site, or group). Use with care.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newSettingsOverridesListCmd())
	cmd.AddCommand(newSettingsOverridesGetCmd())
	cmd.AddCommand(newSettingsOverridesCreateCmd())
	cmd.AddCommand(newSettingsOverridesUpdateCmd())
	cmd.AddCommand(newSettingsOverridesDeleteCmd())
	return cmd
}

func newSettingsOverridesListCmd() *cobra.Command {
	var (
		siteIDs    []string
		accountIDs []string
		groupIDs   []string
		osTypes    []string
		query      string
		cursor     string
		sortBy     string
		sortOrder  string
		limit      int
		all        bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List config overrides",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ConfigOverrideListParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
				OSTypes:    osTypes,
				Query:      query,
				Cursor:     cursor,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
				Limit:      limit,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.ConfigOverride
			var total int

			if all {
				items, total, err = fetchAllREST("config override", func(cur string) ([]mgmt.ConfigOverride, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ConfigOverrideList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				items, pag, err = c.ConfigOverrideList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "OS", "Scope", "Version", "Created"}
			rows := make([][]string, len(items))
			for i, o := range items {
				rows[i] = []string{
					o.ID, o.Name, string(o.OSType),
					string(o.Scope), orDash(o.AgentVersion), o.CreatedAt,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "config override", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type (linux|macos|windows|windows_legacy)")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (id|createdAt|updatedAt|name|scope|osType)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc|desc)")
	return markJSON(cmd)
}

func newSettingsOverridesGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <override-id>",
		Short: "Get a config override by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			o, err := c.ConfigOverrideGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), o)
			}
			cfgBytes, _ := json.MarshalIndent(o.Config, "", "  ")
			rows := [][]string{
				{"ID", o.ID},
				{"Name", o.Name},
				{"Description", orDash(o.Description)},
				{"OS Type", string(o.OSType)},
				{"Scope", string(o.Scope)},
				{"Version Option", string(o.VersionOption)},
				{"Agent Version", orDash(o.AgentVersion)},
				{"Created", orDash(o.CreatedAt)},
				{"Updated", orDash(o.UpdatedAt)},
				{"Config", string(cfgBytes)},
			}
			if o.Site != nil {
				rows = append(rows, []string{"Site", o.Site.ID + " (" + o.Site.Name + ")"})
			}
			if o.Group != nil {
				rows = append(rows, []string{"Group", o.Group.ID + " (" + o.Group.Name + ")"})
			}
			if o.Account != nil {
				rows = append(rows, []string{"Account", o.Account.ID + " (" + o.Account.Name + ")"})
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}

func newSettingsOverridesCreateCmd() *cobra.Command {
	var (
		name          string
		description   string
		osType        string
		config        string
		scope         string
		scopeID       string
		agentVersion  string
		versionOption string
		yes           bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a config override",
		Long: `Create an agent configuration override at the specified scope.

Config overrides are powerful: they change agent behavior for all agents
matching the scope. The --config flag accepts a JSON object of configuration
keys to override.

Scopes: tenant, account, site, group
For non-tenant scopes, --scope-id identifies the target.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if osType == "" {
				return fmt.Errorf("--os-type is required")
			}
			if config == "" {
				return fmt.Errorf("--config is required (JSON object)")
			}
			if scope == "" {
				return fmt.Errorf("--scope is required (tenant|account|site|group)")
			}
			if scope != "tenant" && scopeID == "" {
				return fmt.Errorf("--scope-id is required for %s scope", scope)
			}

			var cfgRaw json.RawMessage
			if err := json.Unmarshal([]byte(config), &cfgRaw); err != nil {
				return fmt.Errorf("--config must be valid JSON: %w", err)
			}

			input := mgmt.ConfigOverrideCreateInput{
				Name:   name,
				OSType: mgmt.ConfigOverrideOSType(osType),
				Config: cfgRaw,
				Scope:  mgmt.ConfigOverrideScope(scope),
			}
			if description != "" {
				input.Description = &description
			}
			if agentVersion != "" {
				input.AgentVersion = &agentVersion
			}
			if versionOption != "" {
				vo := mgmt.ConfigOverrideVersionOption(versionOption)
				input.VersionOption = &vo
			}

			switch scope {
			case "site":
				input.Site = &mgmt.ConfigOverrideScopeRef{ID: scopeID}
			case "group":
				input.Group = &mgmt.ConfigOverrideScopeRef{ID: scopeID}
			case "account":
				input.Account = &mgmt.ConfigOverrideScopeRef{ID: scopeID}
			}

			detail := fmt.Sprintf("create config override %q (%s, scope=%s)", name, osType, scope)
			detail += "\n  Config: " + config
			return guard(cmd.OutOrStdout(), "settings overrides create", detail, name, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.ConfigOverrideCreate(cmd.Context(), input)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created config override %s (%s)\n", created.ID, created.Name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "override name")
	cmd.Flags().StringVar(&description, "description", "", "override description")
	cmd.Flags().StringVar(&osType, "os-type", "", "target OS (linux|macos|windows|windows_legacy)")
	cmd.Flags().StringVar(&config, "config", "", "config override JSON object")
	cmd.Flags().StringVar(&scope, "scope", "", "scope level (tenant|account|site|group)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope target ID (required for non-tenant scopes)")
	cmd.Flags().StringVar(&agentVersion, "agent-version", "", "target agent version")
	cmd.Flags().StringVar(&versionOption, "version-option", "", "version option (ALL|SPECIFIC)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newSettingsOverridesUpdateCmd() *cobra.Command {
	var (
		name          string
		description   string
		osType        string
		config        string
		scope         string
		scopeID       string
		agentVersion  string
		versionOption string
		yes           bool
	)

	cmd := &cobra.Command{
		Use:   "update <override-id>",
		Short: "Update a config override",
		Long: `Update an existing config override. Only provided fields are changed.

Config overrides are powerful: changes take effect on all agents matching
the scope. The --config flag accepts a JSON object of configuration keys.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			input := mgmt.ConfigOverrideUpdateInput{}
			parts := []string{}

			if name != "" {
				input.Name = &name
				parts = append(parts, "name="+name)
			}
			if description != "" {
				input.Description = &description
			}
			if osType != "" {
				ot := mgmt.ConfigOverrideOSType(osType)
				input.OSType = &ot
				parts = append(parts, "osType="+osType)
			}
			if config != "" {
				var cfgRaw json.RawMessage
				if err := json.Unmarshal([]byte(config), &cfgRaw); err != nil {
					return fmt.Errorf("--config must be valid JSON: %w", err)
				}
				input.Config = cfgRaw
				parts = append(parts, "config="+config)
			}
			if scope != "" {
				s := mgmt.ConfigOverrideScope(scope)
				input.Scope = &s
			}
			if agentVersion != "" {
				input.AgentVersion = &agentVersion
			}
			if versionOption != "" {
				vo := mgmt.ConfigOverrideVersionOption(versionOption)
				input.VersionOption = &vo
			}

			if scope != "" && scope != "tenant" && scopeID != "" {
				switch scope {
				case "site":
					input.Site = &mgmt.ConfigOverrideScopeRef{ID: scopeID}
				case "group":
					input.Group = &mgmt.ConfigOverrideScopeRef{ID: scopeID}
				case "account":
					input.Account = &mgmt.ConfigOverrideScopeRef{ID: scopeID}
				}
			}

			detail := fmt.Sprintf("update config override %s", id)
			if len(parts) > 0 {
				detail += " (" + strings.Join(parts, ", ") + ")"
			}
			if config != "" {
				detail += "\n  Config: " + config
			}
			return guard(cmd.OutOrStdout(), "settings overrides update", detail, id, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				updated, err := c.ConfigOverrideUpdate(cmd.Context(), id, input)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), updated)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated config override %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "override name")
	cmd.Flags().StringVar(&description, "description", "", "override description")
	cmd.Flags().StringVar(&osType, "os-type", "", "target OS (linux|macos|windows|windows_legacy)")
	cmd.Flags().StringVar(&config, "config", "", "config override JSON object")
	cmd.Flags().StringVar(&scope, "scope", "", "scope level (tenant|account|site|group)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope target ID")
	cmd.Flags().StringVar(&agentVersion, "agent-version", "", "target agent version")
	cmd.Flags().StringVar(&versionOption, "version-option", "", "version option (ALL|SPECIFIC)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newSettingsOverridesDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <override-id>",
		Short: "Delete a config override",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return guard(cmd.OutOrStdout(), "settings overrides delete", "delete config override "+id, id, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.ConfigOverrideDelete(cmd.Context(), id); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted config override %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}
