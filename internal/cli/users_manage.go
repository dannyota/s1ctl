package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newUsersUpdateCmd() *cobra.Command {
	var fullName, email, scope string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update <user-id>",
		Short: "Update a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fullName == "" && email == "" && scope == "" {
				return fmt.Errorf("nothing to update: pass --full-name, --email, or --scope")
			}
			data := mgmt.UserUpdate{
				FullName: fullName,
				Email:    email,
				Scope:    scope,
			}
			return guard(cmd.OutOrStdout(), "users update", "update user "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				updated, err := c.UsersUpdate(cmd.Context(), args[0], data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), updated)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated user %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fullName, "full-name", "", "new full name")
	cmd.Flags().StringVar(&email, "email", "", "new email address")
	cmd.Flags().StringVar(&scope, "scope", "", "new scope")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newUsers2FACmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "2fa",
		Short: "Enable or disable two-factor authentication for a user",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newUsers2FAToggleCmd("enable"))
	cmd.AddCommand(newUsers2FAToggleCmd("disable"))
	return cmd
}

func newUsers2FAToggleCmd(mode string) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   mode + " <user-id>",
		Short: mode + " two-factor authentication for a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "users 2fa "+mode, mode+" 2FA for user "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if mode == "enable" {
					err = c.Users2FAEnable(cmd.Context(), args[0])
				} else {
					err = c.Users2FADisable(cmd.Context(), args[0])
				}
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": mode + "d", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "2FA %sd for user %s\n", mode, args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newUsersGenerateTokenCmd() *cobra.Command {
	var forceLegacy, yes bool

	cmd := &cobra.Command{
		Use:   "generate-token",
		Short: "Generate an API token for the current user",
		Long: `Generate an API token for the authenticated user. The token is shown
once and replaces any existing token for that user.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Action string is intentionally generic — the generated token must
			// never reach the audit log, which records the action string only.
			return guard(cmd.OutOrStdout(), "users generate-token", "generate API token for current user", "current-user", yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				token, err := c.UsersGenerateToken(cmd.Context(), forceLegacy)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"token": token})
				}
				printGeneratedToken(cmd.OutOrStdout(), token, "")
				noteSensitiveOutput(cmd.ErrOrStderr())
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&forceLegacy, "force-legacy", false, "request a legacy token even when auth-tokens is enabled")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newUsersRevokeTokenCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "revoke-token <user-id>",
		Short: "Revoke a user's API token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "users revoke-token", "revoke API token for user "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.UsersRevokeToken(cmd.Context(), args[0]); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "revoked", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Revoked API token for user %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

// redactUserTokenDetails returns a copy of the token metadata with the secret
// token value and raw JSON stripped, safe to print. The details endpoints
// normally return only timestamps; this defends against an API that echoes the
// secret back to the caller.
func redactUserTokenDetails(d *mgmt.UserTokenDetails) mgmt.UserTokenDetails {
	redacted := *d
	redacted.Token = ""
	redacted.Raw = nil
	return redacted
}

func newUsersTokenDetailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-details [<user-id>]",
		Short: "Show API-token metadata (created/expires) for a user",
		Long: `Show API-token metadata for a user. With no argument, reports the
authenticated user's token; with a user ID, reports that user's token. Only
timestamps are shown — any secret value is redacted.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			var d *mgmt.UserTokenDetails
			if len(args) == 0 {
				d, err = c.UsersTokenDetails(cmd.Context())
			} else {
				d, err = c.UsersTokenDetailsByID(cmd.Context(), args[0])
			}
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), redactUserTokenDetails(d))
			}
			rows := [][]string{
				{"Created", orDash(d.CreatedAt)},
				{"Expires", orDash(d.ExpiresAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}

func newUsersDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <user-id>",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "users delete", "delete user "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.UsersDelete(cmd.Context(), args[0]); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted user %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}
