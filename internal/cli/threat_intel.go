package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// iocSeverityNames maps friendly severity names to OCSF scores (0-7).
var iocSeverityNames = map[string]mgmt.IOCSeverity{
	"unknown":       mgmt.IOCSeverityUnknown,
	"informational": mgmt.IOCSeverityInformational,
	"low":           mgmt.IOCSeverityLow,
	"medium":        mgmt.IOCSeverityMedium,
	"high":          mgmt.IOCSeverityHigh,
	"critical":      mgmt.IOCSeverityCritical,
	"fatal":         mgmt.IOCSeverityFatal,
}

// parseIOCSeverity accepts a severity name (case-insensitive) or a raw OCSF
// score (0-7) and returns the typed severity.
func parseIOCSeverity(s string) (mgmt.IOCSeverity, error) {
	if sev, ok := iocSeverityNames[strings.ToLower(s)]; ok {
		return sev, nil
	}
	if n, err := strconv.Atoi(s); err == nil && n >= 0 && n <= 7 {
		return mgmt.IOCSeverity(n), nil
	}
	return 0, fmt.Errorf("invalid severity %q (use Unknown, Informational, Low, Medium, High, Critical, Fatal, or 0-7)", s)
}

func parseIOCSeverities(vals []string) ([]mgmt.IOCSeverity, error) {
	sevs := make([]mgmt.IOCSeverity, 0, len(vals))
	for _, s := range vals {
		sev, err := parseIOCSeverity(s)
		if err != nil {
			return nil, err
		}
		sevs = append(sevs, sev)
	}
	return sevs, nil
}

func newIOCsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iocs",
		Short: "Manage threat intelligence IOCs",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newIOCsListCmd())
	cmd.AddCommand(newIOCsCreateCmd())
	cmd.AddCommand(newIOCsDeleteCmd())
	cmd.AddCommand(newIOCsConfigCmd())
	return cmd
}

func newIOCsListCmd() *cobra.Command {
	var severities, sources, creators []string
	var iocType, value, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List threat intelligence IOCs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			sevs, err := parseIOCSeverities(severities)
			if err != nil {
				return err
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.IOCListParams{
				Type:       mgmt.IOCType(iocType),
				Severities: sevs,
				Sources:    sources,
				Creators:   creators,
				Value:      value,
				Limit:      limit,
				Cursor:     cursor,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var iocs []mgmt.IOC
			var total int

			if all {
				iocs, total, err = fetchAllREST("IOC", func(cur string) ([]mgmt.IOC, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.IOCsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				iocs, pag, err = c.IOCsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"UUID", "Type", "Value", "Severity", "Source", "Created"}
			rows := make([][]string, len(iocs))
			for i, ioc := range iocs {
				rows[i] = []string{
					ioc.UUID,
					string(ioc.Type),
					truncate(ioc.Value, 50),
					ioc.Severity.String(),
					orDash(ioc.Source),
					orDash(ioc.CreationTime),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, iocs, len(iocs), total, "IOC", all)
		},
	}
	cmd.Flags().StringVar(&iocType, "type", "", "filter by IOC type (DNS, IPV4, IPV6, MD5, SHA1, SHA256, URL)")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (Unknown, Informational, Low, Medium, High, Critical, Fatal, or 0-7)")
	cmd.Flags().StringSliceVar(&sources, "source", nil, "filter by source")
	cmd.Flags().StringSliceVar(&creators, "creator", nil, "filter by creator (substring match)")
	cmd.Flags().StringVar(&value, "value", "", "filter by IOC value")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (id, creationTime, uploadTime, updatedAt, source, type)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return markJSON(cmd)
}

func newIOCsCreateCmd() *cobra.Command {
	var (
		iocType     string
		value       string
		source      string
		severity    string
		method      string
		name        string
		description string
		externalID  string
		validUntil  string
		yes         bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a threat intelligence IOC",
		Long: `Create a new threat intelligence indicator of compromise.

Types: DNS, IPV4, IPV6, MD5, SHA1, SHA256, URL
Severities: Unknown, Informational, Low, Medium, High, Critical, Fatal (OCSF scores 0-7)

Dry-run by default; pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if iocType == "" {
				return fmt.Errorf("--type is required")
			}
			if value == "" {
				return fmt.Errorf("--value is required")
			}
			if source == "" {
				return fmt.Errorf("--source is required")
			}

			ioc := mgmt.IOCCreateInput{
				Type:        mgmt.IOCType(iocType),
				Value:       value,
				Source:      source,
				Method:      method,
				Name:        name,
				Description: description,
				ExternalID:  externalID,
				ValidUntil:  validUntil,
			}
			if severity != "" {
				sev, err := parseIOCSeverity(severity)
				if err != nil {
					return err
				}
				ioc.Severity = &sev
			}

			return guard(cmd.OutOrStdout(), "iocs create", fmt.Sprintf("create %s IOC for %q", iocType, value), value, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.IOCsCreate(cmd.Context(), []mgmt.IOCCreateInput{ioc})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", pluralize(len(created), "IOC"))
				for _, ind := range created {
					fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", ind.UUID)
				}
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&iocType, "type", "", "IOC type (DNS, IPV4, IPV6, MD5, SHA1, SHA256, URL)")
	cmd.Flags().StringVar(&value, "value", "", "indicator value")
	cmd.Flags().StringVar(&source, "source", "", "intelligence source")
	cmd.Flags().StringVar(&severity, "severity", "", "severity (Unknown, Informational, Low, Medium, High, Critical, Fatal, or 0-7)")
	cmd.Flags().StringVar(&method, "method", "", "comparison method (EQUALS; server default when empty)")
	cmd.Flags().StringVar(&name, "name", "", "IOC name")
	cmd.Flags().StringVar(&description, "description", "", "IOC description")
	cmd.Flags().StringVar(&externalID, "external-id", "", "external reference ID")
	cmd.Flags().StringVar(&validUntil, "valid-until", "", "expiration date (ISO 8601)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newIOCsDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <ioc-uuid...>",
		Short: "Delete threat intelligence IOCs",
		Long: `Delete one or more threat intelligence IOCs by UUID.

Dry-run by default; pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "iocs delete", "delete "+pluralize(len(args), "IOC"), strings.Join(args, ","), yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.IOCsDelete(cmd.Context(), args)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s\n", pluralize(affected, "IOC"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newIOCsConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Show threat intelligence configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			cfgs, err := c.ThreatIntelConfigs(cmd.Context())
			if err != nil {
				return err
			}
			headers := []string{"Scope", "Scope ID", "Min Score", "Threat Disabled", "RetroHunt Disabled", "XDR Matching", "Updated"}
			rows := make([][]string, len(cfgs))
			for i, cfg := range cfgs {
				rows[i] = []string{
					orDash(string(cfg.ScopeLevel)),
					orDash(cfg.ScopeID),
					strconv.Itoa(cfg.ThreatMinScore),
					boolIcon(cfg.DisableThreat),
					boolIcon(cfg.DisableRH),
					boolIcon(cfg.EnableXDRMatching),
					orDash(cfg.UpdatedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, cfgs, len(cfgs), len(cfgs), "config", false)
		},
	}
	return markJSON(cmd)
}
