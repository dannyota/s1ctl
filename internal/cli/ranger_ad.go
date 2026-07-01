package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRangerADCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ranger-ad",
		Aliases: []string{"rad"},
		Short:   "Ranger AD exposure assessment (ISPM)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newRADStatusCmd())
	cmd.AddCommand(newRADExposuresCmd())
	cmd.AddCommand(newRADAffectedObjectsCmd())
	cmd.AddCommand(newRADAssessCmd())
	return cmd
}

func newRADStatusCmd() *cobra.Command {
	var siteIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show AD assessment status",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ADAssessmentStatusParams{
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
			}
			status, err := c.RangerADAssessmentStatus(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), status)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Assessment Status: %s\n", status.Status)

			if len(status.Domains) > 0 {
				fmt.Fprintln(cmd.OutOrStdout())
				headers := []string{"Domain", "Forest", "Completed", "Total Jobs", "Done"}
				rows := make([][]string, len(status.Domains))
				for i, d := range status.Domains {
					rows[i] = []string{
						d.DomainName,
						orDash(d.ForestName),
						boolIcon(d.DomainCompleted),
						strconv.Itoa(d.TotalJobs),
						strconv.Itoa(d.CompletedJobs),
					}
				}
				printTable(headers, rows)
			}

			if len(status.Tenants) > 0 {
				fmt.Fprintln(cmd.OutOrStdout())
				headers := []string{"Tenant ID", "Completed", "Total Jobs", "Done"}
				rows := make([][]string, len(status.Tenants))
				for i, t := range status.Tenants {
					rows[i] = []string{
						t.TenantID,
						boolIcon(t.TenantCompleted),
						strconv.Itoa(t.TotalJobs),
						strconv.Itoa(t.CompletedJobs),
					}
				}
				printTable(headers, rows)
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return cmd
}

func newRADExposuresCmd() *cobra.Command {
	var severity, status, source, detectionName, domainName []string
	var siteIDs, accountIDs []string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "exposures",
		Short: "List AD exposures",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ADExposureListParams{
				Limit:      limit,
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
				Filter: mgmt.ADExposureFilter{
					Severity:        severity,
					DetectionStatus: status,
					Source:          source,
					DetectionName:   detectionName,
					DomainName:      domainName,
				},
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var exposures []mgmt.ADExposure
			var total int

			if all {
				for {
					items, pag, fetchErr := c.RangerADExposures(cmd.Context(), params)
					if fetchErr != nil {
						return fetchErr
					}
					exposures = append(exposures, items...)
					if pag != nil {
						total = pag.TotalItems
					}
					printProgress("exposure", len(exposures), total)
					if len(items) == 0 || len(exposures) >= total {
						break
					}
					params.Skip = len(exposures)
				}
				clearProgress()
			} else {
				var pag *mgmt.Pagination
				exposures, pag, err = c.RangerADExposures(cmd.Context(), params)
				if err != nil {
					return err
				}
				if pag != nil {
					total = pag.TotalItems
				}
			}

			headers := []string{"ID", "Detection", "Status", "Severity", "Source", "Domain", "Vulnerable"}
			rows := make([][]string, len(exposures))
			for i, e := range exposures {
				rows[i] = []string{
					e.ID,
					truncate(e.DetectionName, 40),
					string(e.DetectionStatus),
					string(e.Severity),
					string(e.Source),
					e.DomainName,
					strconv.Itoa(e.VulnerableCount),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, exposures, len(exposures), total, "exposure", all)
		},
	}
	cmd.Flags().StringSliceVar(&severity, "severity", nil, "filter by severity (Critical, High, Medium, Low)")
	cmd.Flags().StringSliceVar(&status, "status", nil, "filter by detection status (Vulnerable, Not_Vulnerable, Skipped, ...)")
	cmd.Flags().StringSliceVar(&source, "source", nil, "filter by source (OnPremAD, AzureAD)")
	cmd.Flags().StringSliceVar(&detectionName, "detection", nil, "filter by detection name")
	cmd.Flags().StringSliceVar(&domainName, "domain", nil, "filter by domain name")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().IntVar(&limit, "limit", 0, fmt.Sprintf("max results per page (default %d)", defaultPageSize))
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	return cmd
}

func newRADAffectedObjectsCmd() *cobra.Command {
	var detectionName, domainName, objectType []string
	var siteIDs, accountIDs []string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "affected-objects",
		Short: "List objects affected by an AD exposure",
		Long: `List Active Directory objects affected by a specific detection.
Requires --detection and --domain flags to identify the exposure.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if len(detectionName) == 0 || len(domainName) == 0 {
				return fmt.Errorf("--detection and --domain are required")
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ADAffectedObjectListParams{
				Limit:      limit,
				SiteIDs:    strings.Join(siteIDs, ","),
				AccountIDs: strings.Join(accountIDs, ","),
				Filter: mgmt.ADAffectedObjectFilter{
					DetectionName: detectionName,
					DomainName:    domainName,
					ObjectType:    objectType,
				},
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var objects []mgmt.ADAffectedObject
			var total int

			if all {
				for {
					items, pag, fetchErr := c.RangerADAffectedObjects(cmd.Context(), params)
					if fetchErr != nil {
						return fetchErr
					}
					objects = append(objects, items...)
					if pag != nil {
						total = pag.TotalItems
					}
					printProgress("object", len(objects), total)
					if len(items) == 0 || len(objects) >= total {
						break
					}
					params.Skip = len(objects)
				}
				clearProgress()
			} else {
				var pag *mgmt.Pagination
				objects, pag, err = c.RangerADAffectedObjects(cmd.Context(), params)
				if err != nil {
					return err
				}
				if pag != nil {
					total = pag.TotalItems
				}
			}

			headers := []string{"ID", "Display Name", "Type", "DN", "Account Status"}
			rows := make([][]string, len(objects))
			for i, o := range objects {
				rows[i] = []string{
					strconv.Itoa(o.ID),
					orDash(deref(o.DisplayName)),
					orDash(deref(o.ObjectType)),
					truncate(orDash(deref(o.DN)), 60),
					orDash(deref(o.AccountStatus)),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, objects, len(objects), total, "object", all)
		},
	}
	cmd.Flags().StringSliceVar(&detectionName, "detection", nil, "detection name (required)")
	cmd.Flags().StringSliceVar(&domainName, "domain", nil, "domain name (required)")
	cmd.Flags().StringSliceVar(&objectType, "object-type", nil, "filter by object type (Computer, User, Group, ...)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().IntVar(&limit, "limit", 0, fmt.Sprintf("max results per page (default %d)", defaultPageSize))
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	return cmd
}

func newRADAssessCmd() *cobra.Command {
	var yes, fullScan bool
	var scanSource string
	var siteIDs, accountIDs []string
	var domainName []string

	cmd := &cobra.Command{
		Use:   "assess",
		Short: "Trigger a new AD assessment",
		Long: `Trigger a Ranger AD assessment scan.
Use --full-scan for a complete scan, or omit for a targeted reassessment.
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			scanType := "targeted assessment"
			if fullScan {
				scanType = "full scan"
			}
			target := strings.Join(domainName, ",")
			if target == "" {
				target = "all"
			}
			return guard(cmd.OutOrStdout(), "ranger-ad assess", "trigger "+scanType, target, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}

				filter := mgmt.ADTriggerAssessmentFilter{
					IsFullScan: fullScan,
					DomainName: domainName,
				}
				if scanSource != "" {
					filter.ScanSource = &scanSource
				}

				params := &mgmt.ADTriggerAssessmentParams{
					SiteIDs:    strings.Join(siteIDs, ","),
					AccountIDs: strings.Join(accountIDs, ","),
					Filter:     filter,
				}
				success, msg, err := c.RangerADTriggerAssessment(cmd.Context(), params)
				if err != nil {
					return err
				}
				if !success {
					return fmt.Errorf("assessment trigger failed: %s", msg)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Assessment triggered: %s\n", msg)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	cmd.Flags().BoolVar(&fullScan, "full-scan", false, "perform a full scan (default: targeted)")
	cmd.Flags().StringSliceVar(&domainName, "domain", nil, "domain names to scan")
	cmd.Flags().StringVar(&scanSource, "scan-source", "", "scan source (AD, Azure)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return cmd
}

// deref returns the value of a string pointer, or empty string if nil.
func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
