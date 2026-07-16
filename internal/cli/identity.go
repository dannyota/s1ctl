package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity",
		Short: "Manage Identity AD Service configuration, connectors, and ISPM",
		Long: `Identity AD Service and ISPM (Identity Security Posture Management).

Covers AD configuration management, connector operations, onboarding status,
and ISPM exposure management (skip/acknowledge). The existing ranger-ad group
covers posture reads (status, exposures, affected-objects) and assessment
triggers; this group covers the Identity AD Service configuration layer and
ISPM write operations.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newIdentityOnboardCmd())
	cmd.AddCommand(newIdentityConfigCmd())
	cmd.AddCommand(newIdentityConnectorCmd())
	cmd.AddCommand(newIdentityDomainsCmd())
	cmd.AddCommand(newIdentityFeaturesCmd())
	cmd.AddCommand(newIdentityTimezonesCmd())
	cmd.AddCommand(newIdentitySkipExposuresCmd())
	cmd.AddCommand(newIdentityAckExposuresCmd())
	return cmd
}

// --- onboard ---

func newIdentityOnboardCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "onboard",
		Short: "Show AD service onboarding status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.IdentityParams{
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
			}
			status, err := c.IdentityOnboardingStatus(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), status)
			}
			rows := [][]string{
				{"Status", string(status.Status)},
				{"Connector", string(status.ADConnector)},
				{"Domain", orDash(status.DomainName)},
				{"Features", strings.Join(status.FeatureSelected, ", ")},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

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

// redactADConfig returns a copy of the AD configuration with the username
// field blanked. The API returns the stored credential username — we redact
// it to avoid leaking bind credentials. Raw is also stripped because it
// contains the unredacted JSON.
func redactADConfig(cfg mgmt.ADConfiguration) mgmt.ADConfiguration {
	cfg.Username = ""
	cfg.Raw = nil
	return cfg
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
			params := &mgmt.IdentityParams{
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
			}
			configs, err := c.IdentityADConfigurations(cmd.Context(), params)
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
			rows := make([][]string, len(configs))
			for i, cfg := range configs {
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
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
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
				params := &mgmt.IdentityParams{
					SiteIDs:    strings.Join(siteIDs, ","),
					AccountIDs: strings.Join(accountIDs, ","),
				}
				if err := c.IdentityADConfigurationAdd(cmd.Context(), params, input); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "AD configuration added.")
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
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
				params := &mgmt.IdentityParams{
					SiteIDs:    strings.Join(siteIDs, ","),
					AccountIDs: strings.Join(accountIDs, ","),
				}
				if err := c.IdentityADConfigurationDelete(cmd.Context(), params, ids); err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted %d AD configuration(s).\n", len(ids))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return cmd
}

// --- connector ---

func newIdentityConnectorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connector",
		Short: "Manage AD connectors (Cloudlink agents)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newIdentityConnectorListCmd())
	cmd.AddCommand(newIdentityConnectorGetCmd())
	cmd.AddCommand(newIdentityConnectorReplaceCmd())
	cmd.AddCommand(newIdentityConnectorAgentsCmd())
	return cmd
}

func newIdentityConnectorListCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all AD connectors",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.IdentityParams{
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
			}
			connectors, err := c.IdentityConnectors(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), connectors)
			}
			headers := []string{"Cloudlink ID", "Computer", "Status", "Agent Type", "OS", "Version", "Domain", "IP"}
			rows := make([][]string, len(connectors))
			for i, cn := range connectors {
				rows[i] = []string{
					strconv.FormatInt(cn.CloudlinkID, 10),
					cn.ComputerName,
					string(cn.Status),
					cn.AgentType,
					truncate(cn.OSName, 30),
					cn.Version,
					cn.DomainName,
					cn.IPAddress,
				}
			}
			printTable(headers, rows)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

func newIdentityConnectorGetCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the current AD connector configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.IdentityParams{
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
			}
			cn, err := c.IdentityConnector(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), cn)
			}
			rows := [][]string{
				{"GUID", cn.GUID},
				{"Computer", cn.ComputerName},
				{"Status", string(cn.Status)},
				{"Agent Type", cn.AgentType},
				{"OS", cn.OSName},
				{"Version", cn.Version},
				{"Domain", cn.DomainName},
				{"IP", cn.IPAddress},
				{"Unified Agent", boolIcon(cn.IsUnifiedAgent)},
				{"Last Seen", orDash(cn.LastSeen)},
				{"Scope", orDash(cn.ScopePath)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

func newIdentityConnectorReplaceCmd() *cobra.Command {
	var yes bool
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "replace [agent-uuid]",
		Short: "Replace the AD connector with a different agent",
		Long: `Replace the AD connector (Cloudlink) with a new agent by UUID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "identity connector replace", "replace AD connector", args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				params := &mgmt.IdentityParams{
					SiteIDs:    strings.Join(siteIDs, ","),
					AccountIDs: strings.Join(accountIDs, ","),
				}
				if err := c.IdentityConnectorReplace(cmd.Context(), params, args[0]); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "AD connector replaced.")
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return cmd
}

func newIdentityConnectorAgentsCmd() *cobra.Command {
	var siteIDs, accountIDs []string
	var filterInput string

	cmd := &cobra.Command{
		Use:   "agents",
		Short: "List Windows agents available as connectors",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.WindowsAgentParams{
				SiteIDs:     strings.Join(siteIDs, ","),
				AccountIDs:  strings.Join(accountIDs, ","),
				FilterInput: filterInput,
			}
			agents, err := c.IdentityWindowsAgents(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), agents)
			}
			headers := []string{"UUID", "Host", "OS", "Version", "Status", "Domain", "IP"}
			rows := make([][]string, len(agents))
			for i, a := range agents {
				rows[i] = []string{
					a.UUID,
					a.HostName,
					truncate(a.OSName, 25),
					a.AgentVersion,
					a.Status,
					a.DomainName,
					a.IPAddress,
				}
			}
			printTable(headers, rows)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringVar(&filterInput, "filter", "", "filter agents by name")
	return markJSON(cmd)
}

// --- domains, features, timezones ---

func newIdentityDomainsCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "domains",
		Short: "List AD domains",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.IdentityParams{
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
			}
			domains, err := c.IdentityDomains(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), domains)
			}
			headers := []string{"Domain", "Parent Domain", "Root"}
			rows := make([][]string, len(domains))
			for i, d := range domains {
				rows[i] = []string{d.Domain, orDash(d.ParentDomain), boolIcon(d.Root)}
			}
			printTable(headers, rows)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

func newIdentityFeaturesCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "features",
		Short: "List available AD features",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.IdentityParams{
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
			}
			features, err := c.IdentityAvailableFeatures(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), features)
			}
			headers := []string{"Feature", "Available"}
			rows := make([][]string, len(features))
			for i, f := range features {
				rows[i] = []string{string(f.FeatureName), boolIcon(f.Available)}
			}
			printTable(headers, rows)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

func newIdentityTimezonesCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "timezones",
		Short: "List available timezones for AD configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.IdentityParams{
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
			}
			tzs, err := c.IdentityTimezones(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), tzs)
			}
			headers := []string{"Timezone ID", "Display Name"}
			rows := make([][]string, len(tzs))
			for i, tz := range tzs {
				rows[i] = []string{tz.TimeZoneID, tz.DisplayName}
			}
			printTable(headers, rows)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

// --- ISPM mutations (skip/ack exposures) ---

func newIdentitySkipExposuresCmd() *cobra.Command {
	var yes bool
	var siteIDs, accountIDs, detectionName, domainName []string
	var unskip bool
	var reason string

	cmd := &cobra.Command{
		Use:   "skip-exposures",
		Short: "Skip or unskip ISPM exposures",
		Long: `Set exposures as skipped (accepted risk) or unskip previously skipped exposures.
Requires --detection and --domain. Dry-run by default — pass --yes to apply.
Use --unskip to reverse a previous skip.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if len(detectionName) == 0 || len(domainName) == 0 {
				return fmt.Errorf("--detection and --domain are required")
			}
			action := "skip"
			if unskip {
				action = "unskip"
			}
			target := strings.Join(detectionName, ",") + " in " + strings.Join(domainName, ",")
			return guard(cmd.OutOrStdout(), "identity skip-exposures", action+" exposures", target, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				params := &mgmt.ADSkipExposuresParams{
					SiteIDs:    strings.Join(siteIDs, ","),
					AccountIDs: strings.Join(accountIDs, ","),
					Filter: mgmt.ADSkipExposuresFilter{
						DetectionName: detectionName,
						DomainName:    domainName,
						Skip:          !unskip,
						SkipReason:    reason,
					},
				}
				ok, msg, err := c.RangerADSetSkippedExposures(cmd.Context(), params)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("skip-exposures failed: %s", msg)
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{"success": true, "message": msg})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Exposures %sed: %s\n", action, msg)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().BoolVar(&unskip, "unskip", false, "reverse a previous skip")
	cmd.Flags().StringVar(&reason, "reason", "", "reason for skipping")
	cmd.Flags().StringSliceVar(&detectionName, "detection", nil, "detection name(s) (required)")
	cmd.Flags().StringSliceVar(&domainName, "domain", nil, "domain name(s) (required)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

func newIdentityAckExposuresCmd() *cobra.Command {
	var yes bool
	var siteIDs, accountIDs, detectionName, domainName []string
	var unack bool

	cmd := &cobra.Command{
		Use:   "ack-exposures",
		Short: "Acknowledge or unacknowledge ISPM exposures",
		Long: `Set the acknowledged status on exposures. Requires --detection and --domain.
Dry-run by default — pass --yes to apply. Use --unack to reverse.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if len(detectionName) == 0 || len(domainName) == 0 {
				return fmt.Errorf("--detection and --domain are required")
			}
			action := "acknowledge"
			if unack {
				action = "unacknowledge"
			}
			target := strings.Join(detectionName, ",") + " in " + strings.Join(domainName, ",")
			return guard(cmd.OutOrStdout(), "identity ack-exposures", action+" exposures", target, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				params := &mgmt.ADAckExposuresParams{
					SiteIDs:    strings.Join(siteIDs, ","),
					AccountIDs: strings.Join(accountIDs, ","),
					Filter: mgmt.ADAckExposuresFilter{
						DetectionName: detectionName,
						DomainName:    domainName,
						Acknowledged:  !unack,
					},
				}
				ok, msg, err := c.RangerADSetAckStatus(cmd.Context(), params)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("ack-exposures failed: %s", msg)
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{"success": true, "message": msg})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Exposures %sd: %s\n", action, msg)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().BoolVar(&unack, "unack", false, "reverse acknowledgement")
	cmd.Flags().StringSliceVar(&detectionName, "detection", nil, "detection name(s) (required)")
	cmd.Flags().StringSliceVar(&domainName, "domain", nil, "domain name(s) (required)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}
