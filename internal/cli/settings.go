package cli

import (
	"fmt"
	"os"
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
	cmd.AddCommand(newSettingsUpdateCmd())
	cmd.AddCommand(newSettingsTestCmd())
	cmd.AddCommand(newSettingsSSOCertCmd())
	cmd.AddCommand(newSettingsCancelPendingEmailsCmd())
	cmd.AddCommand(newSettingsDeleteRecipientCmd())
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
				{"category": "sms", "description": "SMS notification service configuration"},
				{"category": "recipients", "description": "Notification recipient list"},
				{"category": "active-directory", "description": "Active Directory integration and scope mapping"},
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

Types: notifications, sso, smtp, syslog, sms, recipients, ad, ad-scope-mapping`,
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
			case "sms":
				return settingsGetSMS(cmd, c, params)
			case "recipients":
				return settingsGetRecipients(cmd, c, params)
			case "ad":
				return settingsGetAD(cmd, c, params)
			case "ad-scope-mapping":
				return settingsGetADScopeMapping(cmd, c, params)
			default:
				return fmt.Errorf("unknown settings type %q (valid: notifications, sso, smtp, syslog, sms, recipients, ad, ad-scope-mapping)", args[0])
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

func settingsGetSMS(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	s, err := c.SettingsSMSGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), s)
	}
	printTable([]string{"Field", "Value"}, [][]string{{"Enabled", boolIcon(s.Enabled)}})
	return nil
}

func settingsGetRecipients(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	recipients, err := c.SettingsRecipientsGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), recipients)
	}
	rows := make([][]string, len(recipients))
	for i, r := range recipients {
		rows[i] = []string{r.ID, orDash(r.Name), orDash(r.Email), orDash(r.SMS)}
	}
	printTable([]string{"ID", "Name", "Email", "SMS"}, rows)
	return nil
}

// redactADSettings returns a copy of the Active Directory settings with the
// bind password and raw JSON stripped, safe to print. The GET endpoint does not
// echo the password; this defends against an API that returns it anyway.
func redactADSettings(s *mgmt.ADSettings) mgmt.ADSettings {
	redacted := *s
	redacted.Password = ""
	redacted.Raw = nil
	return redacted
}

func settingsGetAD(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	s, err := c.SettingsADGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), redactADSettings(s))
	}
	rows := [][]string{
		{"Enabled", boolIcon(s.Enabled)},
		{"Host", orDash(s.Host)},
		{"Port", strconv.Itoa(s.Port)},
		{"Username", orDash(s.Username)},
		{"Root DN", orDash(s.RootDN)},
		{"SSL", boolIcon(s.SSL)},
	}
	printTable([]string{"Field", "Value"}, rows)
	return nil
}

func settingsGetADScopeMapping(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	s, err := c.SettingsADScopeMappingGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), s)
	}
	rows := [][]string{
		{"Admin", orDash(strings.Join(s.Admin, ", "))},
		{"Viewer", orDash(strings.Join(s.Viewer, ", "))},
	}
	printTable([]string{"Scope", "Groups"}, rows)
	return nil
}

func newSettingsTestCmd() *cobra.Command {
	var siteIDs, accountIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "test <type>",
		Short: "Test settings connectivity",
		Long: `Test connectivity for SMTP, syslog, or Active Directory settings.

Types: smtp, syslog, ad`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			settingsType := args[0]
			if settingsType != "smtp" && settingsType != "syslog" && settingsType != "ad" {
				return fmt.Errorf("unknown settings type %q (valid: smtp, syslog, ad)", settingsType)
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
				case "ad":
					return settingsTestAD(cmd, c, params)
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

func settingsTestAD(cmd *cobra.Command, c *mgmt.Client, params *mgmt.SettingsParams) error {
	current, err := c.SettingsADGet(cmd.Context(), params)
	if err != nil {
		return err
	}
	result, err := c.SettingsADTest(cmd.Context(), params, *current)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), result)
	}
	if result.Status {
		fmt.Fprintf(cmd.OutOrStdout(), "Test passed: active directory connectivity OK\n")
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Test failed: active directory connectivity check returned failure\n")
	}
	return nil
}

func newSettingsSSOCertCmd() *cobra.Command {
	var out string
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "sso-cert",
		Short: "Show or download the SSO service-provider signing certificate",
		Long: `Show the SAML service-provider signing certificate. The certificate is
public key material, not a secret. With --out, download the raw certificate
file to disk; otherwise print its metadata and PEM.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.SettingsParams{SiteIDs: siteIDs, AccountIDs: accountIDs}

			if out != "" {
				data, err := c.SettingsSSOCertDownload(cmd.Context(), params)
				if err != nil {
					return err
				}
				if err := os.WriteFile(out, data, 0o644); err != nil { //nolint:gosec
					return fmt.Errorf("write %s: %w", out, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Wrote certificate to %s (%d bytes)\n", out, len(data))
				return nil
			}

			cert, err := c.SettingsSSOCert(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), cert)
			}
			rows := [][]string{
				{"File Name", orDash(cert.FileName)},
				{"Issued At", orDash(cert.IssuedAt)},
				{"Expires At", orDash(cert.ExpiresAt)},
			}
			printTable([]string{"Field", "Value"}, rows)
			if cert.PEM != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", cert.PEM)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&out, "out", "", "write the downloaded certificate to this file")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return cmd
}

func newSettingsCancelPendingEmailsCmd() *cobra.Command {
	var siteIDs, accountIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "cancel-pending-emails",
		Short: "Cancel queued pending email notifications",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			target := strings.Join(append(append([]string{}, siteIDs...), accountIDs...), ",")
			if target == "" {
				target = "all"
			}
			return guard(cmd.OutOrStdout(), "settings cancel-pending-emails", "cancel pending email notifications", target, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				params := &mgmt.SettingsParams{SiteIDs: siteIDs, AccountIDs: accountIDs}
				result, err := c.SettingsCancelPendingEmails(cmd.Context(), params)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), result)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Canceled %d pending email notification(s)\n", result.Canceled)
				return nil
			})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope to site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope to account IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newSettingsDeleteRecipientCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete-recipient <id>",
		Short: "Delete a notification recipient",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return guard(cmd.OutOrStdout(), "settings delete-recipient", "delete notification recipient "+id, id, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.SettingsRecipientDelete(cmd.Context(), id); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted recipient %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
