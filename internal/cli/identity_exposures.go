package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

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
			domains, err := c.IdentityDomains(cmd.Context(), identityParams(siteIDs, accountIDs))
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
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
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
			features, err := c.IdentityAvailableFeatures(cmd.Context(), identityParams(siteIDs, accountIDs))
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
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
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
			tzs, err := c.IdentityTimezones(cmd.Context(), identityParams(siteIDs, accountIDs))
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
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
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
			actionVerb := "skip"
			actionPast := "skipped"
			if unskip {
				actionVerb = "unskip"
				actionPast = "unskipped"
			}
			target := strings.Join(detectionName, ",") + " in " + strings.Join(domainName, ",")
			return guard(cmd.OutOrStdout(), "identity skip-exposures", actionVerb+" exposures", target, yes, func() error {
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
				fmt.Fprintf(cmd.OutOrStdout(), "Exposures %s: %s\n", actionPast, msg)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().BoolVar(&unskip, "unskip", false, "reverse a previous skip")
	cmd.Flags().StringVar(&reason, "reason", "", "reason for skipping")
	cmd.Flags().StringSliceVar(&detectionName, "detection", nil, "detection name(s) (required)")
	cmd.Flags().StringSliceVar(&domainName, "domain", nil, "domain name(s) (required)")
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
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
			actionVerb := "acknowledge"
			actionPast := "acknowledged"
			if unack {
				actionVerb = "unacknowledge"
				actionPast = "unacknowledged"
			}
			target := strings.Join(detectionName, ",") + " in " + strings.Join(domainName, ",")
			return guard(cmd.OutOrStdout(), "identity ack-exposures", actionVerb+" exposures", target, yes, func() error {
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
				fmt.Fprintf(cmd.OutOrStdout(), "Exposures %s: %s\n", actionPast, msg)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().BoolVar(&unack, "unack", false, "reverse acknowledgement")
	cmd.Flags().StringSliceVar(&detectionName, "detection", nil, "detection name(s) (required)")
	cmd.Flags().StringSliceVar(&domainName, "domain", nil, "domain name(s) (required)")
	addIdentityScopeFlags(cmd, &siteIDs, &accountIDs)
	return markJSON(cmd)
}
