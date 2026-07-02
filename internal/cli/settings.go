package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Manage platform settings",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newSettingsListCmd())
	cmd.AddCommand(newSettingsGetCmd())
	cmd.AddCommand(newSettingsTestCmd())
	return cmd
}

func newSettingsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List settings categories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			categories := []map[string]string{
				{"category": "notifications", "description": "Notification preferences and alert routing"},
				{"category": "sso", "description": "SSO/SAML authentication configuration"},
				{"category": "smtp", "description": "SMTP email server configuration"},
				{"category": "syslog", "description": "Syslog forwarding configuration"},
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), categories)
			}

			headers := []string{"Category", "Description"}
			rows := make([][]string, len(categories))
			for i, cat := range categories {
				rows[i] = []string{cat["category"], cat["description"]}
			}
			printTable(headers, rows)
			return nil
		},
	}
}

func newSettingsGetCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "get <type>",
		Short: "Get settings configuration",
		Long: `Get configuration for a specific settings type.

Types: notifications, sso, smtp, syslog`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.SettingsParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
			}

			switch args[0] {
			case "notifications":
				return settingsGetNotifications(cmd, c, params)
			case "sso":
				return settingsGetSSO(cmd, c, params)
			case "smtp":
				return settingsGetSMTP(cmd, c, params)
			case "syslog":
				return settingsGetSyslog(cmd, c, params)
			default:
				return fmt.Errorf("unknown settings type %q (valid: notifications, sso, smtp, syslog)", args[0])
			}
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return cmd
}

func settingsGetNotifications(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	s, err := c.SettingsNotificationsGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), s)
	}

	emailStatus := s.Configurations.Email
	if emailStatus == "" {
		emailStatus = "configured"
	}
	syslogStatus := s.Configurations.Syslog
	if syslogStatus == "" {
		syslogStatus = "configured"
	}

	rows := [][]string{
		{"Email Config", emailStatus},
		{"Syslog Config", syslogStatus},
		{"Last Modified", s.LastModified.UpdatedAt},
		{"Modified By", s.LastModified.UpdatedBy},
	}
	printTable([]string{"Field", "Value"}, rows)
	return nil
}

func settingsGetSSO(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	s, err := c.SettingsSSOGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), s)
	}

	rows := [][]string{
		{"Enabled", boolIcon(s.Enabled)},
		{"IDP SSO URL", s.IDPSsoURL},
		{"IDP Entity ID", s.IDPEntityID},
		{"IDP Certificate", s.IDPCertName},
		{"SP ACS URL", s.SPAcsURL},
		{"SP Entity ID", s.SPEntityID},
		{"Default Role", s.DefaultUserRole},
		{"Auto Provisioning", boolIcon(s.AutoProvisioning)},
		{"Domains", strings.Join(s.Domains, ", ")},
		{"Sign Request", boolIcon(s.SignRequest)},
	}
	printTable([]string{"Field", "Value"}, rows)
	return nil
}

func settingsGetSMTP(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	s, err := c.SettingsSMTPGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		redacted := *s
		redacted.Password = ""
		redacted.Raw = nil
		return printJSON(cmd.OutOrStdout(), redacted)
	}

	rows := [][]string{
		{"Enabled", boolIcon(s.Enabled)},
		{"Inherits", boolIcon(s.Inherits)},
		{"Host", s.Host},
		{"Port", strconv.Itoa(s.Port)},
		{"Encryption", orDash(s.Encryption)},
		{"Username", orDash(s.Username)},
		{"No-Reply Email", orDash(s.NoReplyEmail)},
	}
	printTable([]string{"Field", "Value"}, rows)
	return nil
}

func settingsGetSyslog(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	s, err := c.SettingsSyslogGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		redacted := *s
		redacted.Token = ""
		redacted.ClientKeyContent = ""
		redacted.ClientCertContent = ""
		redacted.ServerCertContent = ""
		redacted.Raw = nil
		return printJSON(cmd.OutOrStdout(), redacted)
	}

	token := "-"
	if s.Token != "" {
		token = redactToken(s.Token)
	}

	rows := [][]string{
		{"Enabled", boolIcon(s.Enabled)},
		{"Host", s.Host},
		{"Port", strconv.Itoa(s.Port)},
		{"SSL", boolIcon(s.SSL)},
		{"Format", orDash(s.Format)},
		{"Token", token},
	}
	printTable([]string{"Field", "Value"}, rows)
	return nil
}

func newSettingsTestCmd() *cobra.Command {
	var siteIDs, accountIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "test <type>",
		Short: "Test settings connectivity",
		Long: `Test connectivity for SMTP or syslog settings.

Types: smtp, syslog`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsType := args[0]
			if settingsType != "smtp" && settingsType != "syslog" {
				return fmt.Errorf("unknown settings type %q (valid: smtp, syslog)", settingsType)
			}

			return guard(cmd.OutOrStdout(), "settings test", "test "+settingsType+" connectivity", settingsType, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				params := &mgmt.SettingsParams{
					SiteIDs:    siteIDs,
					AccountIDs: accountIDs,
				}

				switch settingsType {
				case "smtp":
					return settingsTestSMTP(cmd, c, params)
				case "syslog":
					return settingsTestSyslog(cmd, c, params)
				}
				return nil
			})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func settingsTestSMTP(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	current, err := c.SettingsSMTPGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	result, err := c.SettingsSMTPTest(cmd.Context(), params, *current)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), result)
	}
	if result.Status {
		fmt.Fprintf(cmd.OutOrStdout(), "Test passed: smtp connectivity OK\n")
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Test failed: smtp connectivity check returned failure\n")
	}
	return nil
}

func settingsTestSyslog(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	current, err := c.SettingsSyslogGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	result, err := c.SettingsSyslogTest(cmd.Context(), params, *current)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), result)
	}
	if result.Status {
		fmt.Fprintf(cmd.OutOrStdout(), "Test passed: syslog connectivity OK\n")
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Test failed: syslog connectivity check returned failure\n")
	}
	return nil
}
