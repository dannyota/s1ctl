package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// scopeRoleInput builds a single scope-role assignment from CLI flags, or nil
// when none was requested. Multi-role assignments require the API directly.
func scopeRoleInput(scopeID, roleID, roleName string) []mgmt.ServiceUserScopeRoleInput {
	if scopeID == "" && roleID == "" && roleName == "" {
		return nil
	}
	return []mgmt.ServiceUserScopeRoleInput{{
		ID:       scopeID,
		RoleID:   roleID,
		RoleName: roleName,
	}}
}

// printGeneratedToken writes a freshly minted API token to stdout exactly once.
// The token is the deliverable; it is intentionally not routed through the
// audit log (guard records action strings only).
func printGeneratedToken(w io.Writer, token, expiresAt string) {
	fmt.Fprintln(w, "API token (shown once — store it securely):")
	fmt.Fprintln(w, token)
	if expiresAt != "" {
		fmt.Fprintf(w, "Expires: %s\n", expiresAt)
	}
}

func newServiceUsersCreateCmd() *cobra.Command {
	var (
		name, description, scope, expiration string
		scopeID, roleID, roleName            string
		yes                                  bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a service user (generates an API token)",
		Long: `Create a service user. Creation mints an API token that is shown once.

Scope must be one of: tenant, account, site. For account/site scope, pass
--scope-id (the account/site ID) and a role via --role-id or --role-name.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if scope == "" {
				return fmt.Errorf("--scope is required (tenant, account, site)")
			}
			if expiration == "" {
				return fmt.Errorf("--expiration is required (RFC3339, e.g. 2026-01-01T00:00:00Z)")
			}
			data := mgmt.ServiceUserCreate{
				Name:           name,
				Description:    description,
				ExpirationDate: expiration,
				Scope:          mgmt.ServiceUserScope(scope),
				ScopeRoles:     scopeRoleInput(scopeID, roleID, roleName),
			}
			action := fmt.Sprintf("create %s service user %q", scope, name)
			return guard(cmd.OutOrStdout(), "service-users create", action, name, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.ServiceUsersCreate(cmd.Context(), data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created service user %s (%s)\n", created.ID, created.Name)
				if created.APIToken.Value != "" {
					printGeneratedToken(cmd.OutOrStdout(), created.APIToken.Value, created.APIToken.ExpiresAt)
				}
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "service user name (required)")
	cmd.Flags().StringVar(&description, "description", "", "description")
	cmd.Flags().StringVar(&scope, "scope", "", "scope: tenant, account, site (required)")
	cmd.Flags().StringVar(&expiration, "expiration", "", "token expiration, RFC3339 (required)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "account/site ID for account/site scope")
	cmd.Flags().StringVar(&roleID, "role-id", "", "RBAC role ID to assign")
	cmd.Flags().StringVar(&roleName, "role-name", "", "predefined role name to assign")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newServiceUsersUpdateCmd() *cobra.Command {
	var (
		description, scope        string
		scopeID, roleID, roleName string
		yes                       bool
	)

	cmd := &cobra.Command{
		Use:   "update <service-user-id>",
		Short: "Update a service user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data := mgmt.ServiceUserUpdate{
				Description: description,
				Scope:       mgmt.ServiceUserScope(scope),
				ScopeRoles:  scopeRoleInput(scopeID, roleID, roleName),
			}
			return guard(cmd.OutOrStdout(), "service-users update", "update service user "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				updated, err := c.ServiceUsersUpdate(cmd.Context(), args[0], data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), updated)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated service user %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&description, "description", "", "new description")
	cmd.Flags().StringVar(&scope, "scope", "", "new scope: tenant, account, site")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "account/site ID for account/site scope")
	cmd.Flags().StringVar(&roleID, "role-id", "", "RBAC role ID to assign")
	cmd.Flags().StringVar(&roleName, "role-name", "", "predefined role name to assign")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newServiceUsersDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <service-user-id>",
		Short: "Delete a service user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "service-users delete", "delete service user "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.ServiceUsersDelete(cmd.Context(), args[0]); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted service user %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newServiceUsersBulkDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "bulk-delete <service-user-id>...",
		Short: "Delete multiple service users by ID",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			action := fmt.Sprintf("delete %s", pluralize(len(args), "service user"))
			return guard(cmd.OutOrStdout(), "service-users bulk-delete", action, fmt.Sprintf("%v", args), yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.ServiceUsersBulkDelete(cmd.Context(), args)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s\n", pluralize(affected, "service user"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newServiceUsersGenerateTokenCmd() *cobra.Command {
	var expiration string
	var yes bool

	cmd := &cobra.Command{
		Use:   "generate-token <service-user-id>",
		Short: "Regenerate a service user's API token",
		Long: `Regenerate the API token for a service user. The new token is shown once
and replaces any existing token.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if expiration == "" {
				return fmt.Errorf("--expiration is required (RFC3339, e.g. 2026-01-01T00:00:00Z)")
			}
			action := "generate an API token for service user " + args[0]
			return guard(cmd.OutOrStdout(), "service-users generate-token", action, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				tok, err := c.ServiceUsersGenerateToken(cmd.Context(), args[0], expiration)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), tok)
				}
				printGeneratedToken(cmd.OutOrStdout(), tok.Token, tok.ExpiresAt)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&expiration, "expiration", "", "token expiration, RFC3339 (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
