package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/mgmt"
)

func newRolesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles",
		Short: "Manage RBAC roles",
		Long: `Manage SentinelOne Role-Based Access Control (RBAC) roles.

A role bundles a set of console permissions at a scope (tenant, account, site,
or group). Custom roles can be created, updated, and deleted; predefined
(system) roles are read-only.

Role bodies are large permission sets, so create and update read a declarative
role file (YAML or JSON) rather than taking many flags. Start from 'roles
template' to obtain the permission tree, then edit and apply it.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newRolesListCmd())
	cmd.AddCommand(newRolesGetCmd())
	cmd.AddCommand(newRolesTemplateCmd())
	cmd.AddCommand(newRolesCreateCmd())
	cmd.AddCommand(newRolesUpdateCmd())
	cmd.AddCommand(newRolesDeleteCmd())
	return cmd
}

// roleFile is the declarative representation of an RBAC role read from disk by
// the create and update commands. It holds the role identity (name),
// description, and the flat permission-ID set. Scope is set from CLI flags on
// create (--site-id/--account-id/--group-id/--tenant), not from the file;
// server-managed fields (ID, timestamps, predefined flag, user counts) are not
// part of the write shape.
type roleFile struct {
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	PermissionIDs []string `yaml:"permissionIds,omitempty"`
}

func (f roleFile) toData() mgmt.RoleData {
	return mgmt.RoleData{
		Name:          f.Name,
		Description:   f.Description,
		PermissionIDs: f.PermissionIDs,
	}
}

// roleScopeFilter builds the create scope from CLI flags: the given scope IDs,
// or the global (tenant) scope when no scope flag is set (a role must be
// created somewhere).
func roleScopeFilter(siteIDs, accountIDs, groupIDs []string, tenant bool) mgmt.RoleScopeFilter {
	if len(siteIDs) == 0 && len(accountIDs) == 0 && len(groupIDs) == 0 {
		tenant = true
	}
	return mgmt.RoleScopeFilter{
		SiteIDs:    siteIDs,
		AccountIDs: accountIDs,
		GroupIDs:   groupIDs,
		Tenant:     tenant,
	}
}

func newRolesListCmd() *cobra.Command {
	var accountIDs, siteIDs, groupIDs []string
	var query, cursor, sortBy, sortOrder string
	var predefined bool
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List RBAC roles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.RoleListParams{
				AccountIDs: accountIDs,
				SiteIDs:    siteIDs,
				GroupIDs:   groupIDs,
				Query:      query,
				Limit:      limit,
				Cursor:     cursor,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
			}
			if cmd.Flags().Changed("predefined") {
				params.PredefinedRole = &predefined
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var roles []mgmt.Role
			var total int

			if all {
				roles, total, err = fetchAllREST("role", func(cur string) ([]mgmt.Role, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.RolesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				roles, pag, err = c.RolesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Scope", "Predefined", "Users", "Description"}
			rows := make([][]string, len(roles))
			for i, r := range roles {
				rows[i] = []string{
					r.ID, r.Name, string(r.Scope), boolIcon(r.PredefinedRole),
					fmt.Sprintf("%d", r.UsersInRoles), truncate(r.Description, 40),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, roles, len(roles), total, "role", all)
		},
	}
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search (name, description)")
	cmd.Flags().BoolVar(&predefined, "predefined", false, "filter by predefined roles (true) or custom roles (false)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (e.g. name, createdAt)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return cmd
}

func newRolesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <role-id>",
		Short: "Get a role definition, including its permission tree",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			r, err := c.RoleGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), r)
			}
			rows := [][]string{
				{"ID", r.ID},
				{"Name", r.Name},
				{"Description", orDash(r.Description)},
				{"Scope", string(r.Scope)},
				{"Predefined", boolIcon(r.PredefinedRole)},
				{"Users in role", fmt.Sprintf("%d", r.UsersInRoles)},
				{"Created At", orDash(r.CreatedAt)},
				{"Updated At", orDash(r.UpdatedAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}

func newRolesTemplateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "template",
		Short: "Print the blank role template for editing",
		Long: `Fetch the blank role template (description and the full permission tree with
default values) and print it as JSON. Use it as a starting point for a new
role: edit the values, then create the role with 'roles create --from-file'.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			tmpl, err := c.RoleTemplate(cmd.Context())
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), tmpl)
		},
	}
}

// readRoleFile reads and parses a declarative role file (YAML or JSON, since
// JSON is a subset of YAML). It surfaces the read error with the file name so
// a missing file reports "read <file>".
func readRoleFile(path string) (roleFile, error) {
	var f roleFile
	raw, err := os.ReadFile(path)
	if err != nil {
		return f, fmt.Errorf("read %s: %w", path, err)
	}
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return f, fmt.Errorf("parse %s: %w", path, err)
	}
	if f.Name == "" {
		return f, fmt.Errorf("role file %s has no name", path)
	}
	return f, nil
}

func newRolesCreateCmd() *cobra.Command {
	var fromFile string
	var siteIDs, accountIDs, groupIDs []string
	var tenant, yes bool

	cmd := &cobra.Command{
		Use:   "create --from-file <role.yaml>",
		Short: "Create a role from a declarative role file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			f, err := readRoleFile(fromFile)
			if err != nil {
				return err
			}
			filter := roleScopeFilter(siteIDs, accountIDs, groupIDs, tenant)

			action := fmt.Sprintf("create role %q from %s", f.Name, fromFile)
			return guard(cmd.OutOrStdout(), "roles create", action, fromFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.RoleCreate(cmd.Context(), mgmt.RoleCreate{Data: f.toData(), Filter: filter})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created role %s (%s)\n", created.ID, created.Name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "role definition file, YAML or JSON (required)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "create in these site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "create in these account IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "create in these group IDs")
	cmd.Flags().BoolVar(&tenant, "tenant", false, "create at the global (tenant) scope")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newRolesUpdateCmd() *cobra.Command {
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update <role-id> --from-file <role.yaml>",
		Short: "Update a role from a declarative role file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			f, err := readRoleFile(fromFile)
			if err != nil {
				return err
			}
			if len(f.PermissionIDs) == 0 {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning: no permissionIds in file; the API may replace the role's permissions with an empty set — include permissionIds to set them explicitly")
			}

			action := fmt.Sprintf("update role %s from %s", args[0], fromFile)
			return guard(cmd.OutOrStdout(), "roles update", action, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				updated, err := c.RoleUpdate(cmd.Context(), args[0], mgmt.RoleUpdate{Data: f.toData()})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), updated)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated role %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "role definition file, YAML or JSON (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newRolesDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <role-id>",
		Short: "Delete a role",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "roles delete",
				"delete role "+args[0], args[0], yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					if err := c.RoleDelete(cmd.Context(), args[0]); err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": args[0]})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Deleted role %s\n", args[0])
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
