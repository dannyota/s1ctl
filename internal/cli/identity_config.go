package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// --- config (get/add/delete) ---

func newIdentityConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage AD configurations",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newIdentityConfigGetCmd())
	cmd.AddCommand(newIdentityConfigAddCmd())
	cmd.AddCommand(newIdentityConfigDeleteCmd())
	return cmd
}

func newIdentityConfigGetCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "List AD configurations (credentials redacted)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			configs, err := c.IdentityADConfigurations(cmd.Context(), identityParams(siteIDs, accountIDs))
			if err != nil {
				return err
			}

			redacted := make([]mgmt.ADConfiguration, len(configs))
			for i := range configs {
				redacted[i] = redactADConfig(configs[i])
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), redacted)
			}
			headers := []string{"ID", "Domain", "DC FQDN", "Encryption", "Status", "Connected", "Features"}
			rows := make([][]string, len(redacted))
			for i, cfg := range redacted {
				rows[i] = []string{
					strconv.FormatInt(cfg.ID, 10),
					cfg.DomainName,
					truncate(cfg.DomainControllerFqdn, 40),
					string(cfg.EncryptionMethod),
					string(cfg.AssessmentStatus),
					boolIcon(cfg.IsConnected),
					strings.Join(cfg.FeaturesOpted, ", "),
				}
			}
			printTable(headers, rows)
			return nil
		},
	}
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	return markJSON(cmd)
}

func newIdentityConfigAddCmd() *cobra.Command {
	var (
		yes             bool
		siteIDs         []string
		accountIDs      []string
		domainName      string
		dcFQDN          string
		userName        string
		password        string
		encryption      string
		features        []string
		threatDetection bool
		ldapReferral    bool
		winRMSSL        bool
		adSync          bool
		assessOther     bool
	)

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new AD configuration",
		Long: `Add a new AD configuration with domain and credential details.
Dry-run by default — pass --yes to apply. Credentials (--user, --password) are
sent to the API but never echoed in output.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if domainName == "" || dcFQDN == "" || userName == "" || password == "" {
				return fmt.Errorf("--domain, --dc-fqdn, --user, and --password are required")
			}
			enc := mgmt.EncryptionMethod(strings.ToUpper(encryption))
			if enc != mgmt.EncryptionMethodLDAP && enc != mgmt.EncryptionMethodLDAPS {
				return fmt.Errorf("--encryption must be LDAP or LDAPS")
			}
			return guard(cmd.OutOrStdout(), "identity config add", "add AD configuration", domainName, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				input := mgmt.ADConfigurationInput{
					DomainName:           domainName,
					DomainControllerFqdn: dcFQDN,
					UserName:             userName,
					Password:             password,
					EncryptionMethod:     enc,
					FeaturesOpted:        features,
				}
				if threatDetection {
					input.EnableThreatDetection = &threatDetection
				}
				if ldapReferral {
					input.LDAPReferral = &ldapReferral
				}
				if winRMSSL {
					input.UseWinRmOverSSL = &winRMSSL
				}
				if adSync {
					input.ADSync = &adSync
				}
				if assessOther {
					input.AssessOtherDomainsInForest = &assessOther
				}
				if err := c.IdentityADConfigurationAdd(cmd.Context(), identityParams(siteIDs, accountIDs), input); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "AD configuration added.")
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	cmd.Flags().StringVar(&domainName, "domain", "", "AD domain name (required)")
	cmd.Flags().StringVar(&dcFQDN, "dc-fqdn", "", "domain controller FQDN (required)")
	cmd.Flags().StringVar(&userName, "user", "", "bind username (required, secret)")
	cmd.Flags().StringVar(&password, "password", "", "bind password (required, secret)")
	cmd.Flags().StringVar(&encryption, "encryption", "LDAPS", "LDAP or LDAPS")
	cmd.Flags().StringSliceVar(&features, "feature", nil, "features to enable (RANGER_AD, SINGULARITY_IDENTITY, RANGER_AD_PROTECT)")
	cmd.Flags().BoolVar(&threatDetection, "threat-detection", false, "enable threat detection")
	cmd.Flags().BoolVar(&ldapReferral, "ldap-referral", false, "enable LDAP referral")
	cmd.Flags().BoolVar(&winRMSSL, "winrm-ssl", false, "use WinRM over SSL")
	cmd.Flags().BoolVar(&adSync, "ad-sync", false, "enable AD sync")
	cmd.Flags().BoolVar(&assessOther, "assess-other-domains", false, "assess other domains in forest")
	return cmd
}

func newIdentityConfigDeleteCmd() *cobra.Command {
	var yes bool
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "delete [id...]",
		Short: "Delete AD configurations by ID",
		Long: `Delete one or more AD configurations by their numeric ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := make([]int64, len(args))
			for i, a := range args {
				id, err := strconv.ParseInt(a, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid config ID %q: %w", a, err)
				}
				ids[i] = id
			}
			target := strings.Join(args, ", ")
			return guard(cmd.OutOrStdout(), "identity config delete", "delete AD configuration(s)", target, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.IdentityADConfigurationDelete(cmd.Context(), identityParams(siteIDs, accountIDs), ids); err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted %d AD configuration(s).\n", len(ids))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	return cmd
}
