package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

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
	var types, severities, sources []string
	var value, cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List threat intelligence IOCs",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.IOCListParams{
				Types:      types,
				Severities: severities,
				Sources:    sources,
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

			headers := []string{"ID", "Type", "Value", "Severity", "Source", "Created"}
			rows := make([][]string, len(iocs))
			for i, ioc := range iocs {
				rows[i] = []string{
					ioc.ID,
					string(ioc.Type),
					truncate(ioc.Value, 50),
					string(ioc.Severity),
					orDash(ioc.Source),
					orDash(ioc.CreationTime),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, iocs, len(iocs), total, "IOC", all)
		},
	}
	cmd.Flags().StringSliceVar(&types, "type", nil, "filter by IOC type (DNS, IPV4, IPV6, MD5, SHA1, SHA256, URL)")
	cmd.Flags().StringSliceVar(&severities, "severity", nil, "filter by severity (Low, Medium, High)")
	cmd.Flags().StringSliceVar(&sources, "source", nil, "filter by source")
	cmd.Flags().StringVar(&value, "value", "", "filter by IOC value")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return cmd
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
Severities: Low, Medium, High

Dry-run by default; pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if iocType == "" {
				return fmt.Errorf("--type is required")
			}
			if value == "" {
				return fmt.Errorf("--value is required")
			}

			ioc := mgmt.IOCCreateInput{
				Type:        mgmt.IOCType(iocType),
				Value:       value,
				Source:      source,
				Severity:    mgmt.IOCSeverity(severity),
				Method:      method,
				Name:        name,
				Description: description,
				ExternalID:  externalID,
				ValidUntil:  validUntil,
			}

			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would create %s IOC for %q. Pass --yes to apply.\n",
					iocType, value)
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}
			affected, err := c.IOCsCreate(cmd.Context(), []mgmt.IOCCreateInput{ioc})
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", pluralize(affected, "IOC"))
			return nil
		},
	}
	cmd.Flags().StringVar(&iocType, "type", "", "IOC type (DNS, IPV4, IPV6, MD5, SHA1, SHA256, URL)")
	cmd.Flags().StringVar(&value, "value", "", "indicator value")
	cmd.Flags().StringVar(&source, "source", "", "intelligence source")
	cmd.Flags().StringVar(&severity, "severity", "", "severity (Low, Medium, High)")
	cmd.Flags().StringVar(&method, "method", "", "detection method")
	cmd.Flags().StringVar(&name, "name", "", "IOC name")
	cmd.Flags().StringVar(&description, "description", "", "IOC description")
	cmd.Flags().StringVar(&externalID, "external-id", "", "external reference ID")
	cmd.Flags().StringVar(&validUntil, "valid-until", "", "expiration date (ISO 8601)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newIOCsDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <ioc-id...>",
		Short: "Delete threat intelligence IOCs",
		Long: `Delete one or more threat intelligence IOCs by ID.

Dry-run by default; pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would delete %s. Pass --yes to apply.\n",
					pluralize(len(args), "IOC"))
				return nil
			}

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
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newIOCsConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show threat intelligence configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			cfg, err := c.ThreatIntelConfig(cmd.Context())
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), cfg)
			}
			rows := [][]string{
				{"Total IOCs", fmt.Sprintf("%d", cfg.TotalIOCs)},
				{"Max IOCs", fmt.Sprintf("%d", cfg.MaxIOCs)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}
