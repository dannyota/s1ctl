package cli

import (
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

// --- shared scope-param helpers ---

// addIdentityScopeFlags registers the --site-id and --account-id flags used by
// every identity subcommand.
func addIdentityScopeFlags(cmd *cobra.Command, siteIDs, accountIDs *[]string) {
	cmd.Flags().StringSliceVar(siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(accountIDs, "account-id", nil, "filter by account ID")
}

// identityParams builds an IdentityParams from the scope flag values.
func identityParams(siteIDs, accountIDs []string) *mgmt.IdentityParams {
	return &mgmt.IdentityParams{
		SiteIDs:    strings.Join(siteIDs, ","),
		AccountIDs: strings.Join(accountIDs, ","),
	}
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
			status, err := c.IdentityOnboardingStatus(cmd.Context(), identityParams(siteIDs, accountIDs))
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
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	return markJSON(cmd)
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
