package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newSettingsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update settings from a JSON file (pull with 'settings get', edit, push back)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newSettingsUpdateSubCmd("notifications", "Update notification settings", (*mgmt.Client).SettingsNotificationsUpdate))
	cmd.AddCommand(newSettingsUpdateSubCmd("sso", "Update SSO settings", (*mgmt.Client).SettingsSSOUpdate))
	cmd.AddCommand(newSettingsUpdateSubCmd("smtp", "Update SMTP settings", (*mgmt.Client).SettingsSMTPUpdate))
	cmd.AddCommand(newSettingsUpdateSubCmd("syslog", "Update syslog settings", (*mgmt.Client).SettingsSyslogUpdate))
	cmd.AddCommand(newSettingsUpdateSubCmd("sms", "Update SMS settings", (*mgmt.Client).SettingsSMSUpdate))
	cmd.AddCommand(newSettingsUpdateSubCmd("recipients", "Set a notification recipient", (*mgmt.Client).SettingsRecipientsUpdate))
	cmd.AddCommand(newSettingsUpdateSubCmd("ad", "Update Active Directory settings", (*mgmt.Client).SettingsADUpdate))
	cmd.AddCommand(newSettingsUpdateSubCmd("ad-scope-mapping", "Update Active Directory scope mapping", (*mgmt.Client).SettingsADScopeMappingUpdate))
	return cmd
}

func newSettingsUpdateSubCmd[T any](name, short string, update func(*mgmt.Client, context.Context, *mgmt.SettingsParams, T) (*T, error)) *cobra.Command {
	var fromFile string
	var siteIDs, accountIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   name + " --from-file <settings.json>",
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			raw, err := os.ReadFile(fromFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", fromFile, err)
			}
			var data T
			if err := json.Unmarshal(raw, &data); err != nil {
				return fmt.Errorf("parse %s: %w", fromFile, err)
			}
			action := fmt.Sprintf("update %s settings from %s", name, fromFile)
			return guard(cmd.OutOrStdout(), "settings update "+name, action, fromFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if _, err := update(c, cmd.Context(), &mgmt.SettingsParams{SiteIDs: siteIDs, AccountIDs: accountIDs}, data); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "settings": name})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated %s settings\n", name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "JSON file with the settings payload (required)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope to site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope to account IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
