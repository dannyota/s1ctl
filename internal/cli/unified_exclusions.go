package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newUnifiedExclusionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unified-exclusions",
		Short: "Manage unified exclusions",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newUnifiedExclusionsListCmd())
	cmd.AddCommand(newUnifiedExclusionsCreateCmd())
	cmd.AddCommand(newUnifiedExclusionsExportCmd())
	return cmd
}

func newUnifiedExclusionsListCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs, osTypes, source, modeType, threatType, engines []string
	var nameContains, valueContains []string
	var cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List unified exclusions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.UnifiedExclusionListParams{
				SiteIDs:       siteIDs,
				AccountIDs:    accountIDs,
				GroupIDs:      groupIDs,
				OSTypes:       osTypes,
				Source:        source,
				ModeType:      modeType,
				ThreatType:    threatType,
				Engines:       engines,
				NameContains:  nameContains,
				ValueContains: valueContains,
				Limit:         limit,
				Cursor:        cursor,
				SortBy:        sortBy,
				SortOrder:     sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var exclusions []mgmt.UnifiedExclusion
			var total int

			if all {
				exclusions, total, err = fetchAllREST("exclusion", func(cur string) ([]mgmt.UnifiedExclusion, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.UnifiedExclusionsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				exclusions, pag, err = c.UnifiedExclusionsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "OS", "Threat", "Mode", "Type", "Scope", "Source"}
			rows := make([][]string, len(exclusions))
			for i, e := range exclusions {
				rows[i] = []string{
					e.ID, truncate(e.ExclusionName, 40), e.OSType,
					e.ThreatType, e.ModeType, e.Type,
					e.ScopePath, e.Source,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, exclusions, len(exclusions), total, "exclusion", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type")
	cmd.Flags().StringSliceVar(&source, "source", nil, "filter by source")
	cmd.Flags().StringSliceVar(&modeType, "mode-type", nil, "filter by mode type")
	cmd.Flags().StringSliceVar(&threatType, "threat-type", nil, "filter by threat type")
	cmd.Flags().StringSliceVar(&engines, "engines", nil, "filter by engines")
	cmd.Flags().StringSliceVar(&nameContains, "name", nil, "filter by name (contains)")
	cmd.Flags().StringSliceVar(&valueContains, "value", nil, "filter by value (contains)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return markJSON(cmd)
}

func newUnifiedExclusionsCreateCmd() *cobra.Command {
	var (
		name             string
		osType           string
		threatType       string
		modeType         string
		reason           string
		scopeLevel       string
		scopeID          string
		exclType         string
		description      string
		interactionLevel string
		value            string
		pathType         string
		engines          string
		source           string
		yes              bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a unified exclusion",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if osType == "" {
				return fmt.Errorf("--os-type is required")
			}
			if threatType == "" {
				return fmt.Errorf("--threat-type is required")
			}
			if modeType == "" {
				return fmt.Errorf("--mode-type is required")
			}
			if reason == "" {
				return fmt.Errorf("--reason is required")
			}
			if scopeLevel == "" {
				return fmt.Errorf("--scope-level is required")
			}
			var scopeLevelID *int64
			if scopeID != "" {
				n, err := strconv.ParseInt(scopeID, 10, 64)
				if err != nil {
					return fmt.Errorf("--scope-id must be numeric: %w", err)
				}
				scopeLevelID = &n
			}

			return guard(cmd.OutOrStdout(), "unified-exclusions create", "create unified exclusion "+name+" ("+osType+", "+threatType+")", name, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}

				data := mgmt.UnifiedExclusionCreate{
					ExclusionName:     name,
					OSType:            mgmt.UnifiedExclusionOSType(osType),
					ThreatType:        mgmt.UnifiedExclusionThreatType(threatType),
					ModeType:          mgmt.UnifiedExclusionModeType(modeType),
					Reason:            reason,
					Type:              mgmt.UnifiedExclusionType(exclType),
					Description:       description,
					InteractionLevel:  mgmt.UnifiedExclusionInteractionLevel(interactionLevel),
					PathExclusionType: mgmt.UnifiedExclusionPathType(pathType),
					Engines:           engines,
					Source:            mgmt.UnifiedExclusionSource(source),
				}
				if value != "" {
					data.Value = value
				}

				scope := mgmt.UnifiedExclusionScope{
					ScopeLevel:   mgmt.UnifiedExclusionScopeLevel(scopeLevel),
					ScopeLevelID: scopeLevelID,
				}

				created, err := c.UnifiedExclusionsCreate(cmd.Context(), scope, data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created unified exclusion %s (%s)\n", created.ExclusionName, created.ID)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "exclusion name (required)")
	cmd.Flags().StringVar(&osType, "os-type", "", "target OS type (required)")
	cmd.Flags().StringVar(&threatType, "threat-type", "", "threat type (required)")
	cmd.Flags().StringVar(&modeType, "mode-type", "", "mode type (required)")
	cmd.Flags().StringVar(&reason, "reason", "", "exclusion reason (required)")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level (required)")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope level ID")
	cmd.Flags().StringVar(&exclType, "type", "", "exclusion type")
	cmd.Flags().StringVar(&description, "description", "", "exclusion description")
	cmd.Flags().StringVar(&interactionLevel, "interaction-level", "", "interaction level")
	cmd.Flags().StringVar(&value, "value", "", "exclusion value")
	cmd.Flags().StringVar(&pathType, "path-type", "", "path exclusion type")
	cmd.Flags().StringVar(&engines, "engines", "", "engines")
	cmd.Flags().StringVar(&source, "source", "", "exclusion source")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newUnifiedExclusionsExportCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs, osTypes, source, modeType, threatType []string
	var outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export unified exclusions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.UnifiedExclusionListParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
				OSTypes:    osTypes,
				Source:     source,
				ModeType:   modeType,
				ThreatType: threatType,
			}

			data, err := c.UnifiedExclusionsExport(cmd.Context(), params)
			if err != nil {
				return err
			}

			if outFile != "" {
				if err := os.WriteFile(outFile, data, 0o644); err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Exported to %s\n", outFile)
				return nil
			}
			_, err = cmd.OutOrStdout().Write(data)
			return err
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type")
	cmd.Flags().StringSliceVar(&source, "source", nil, "filter by source")
	cmd.Flags().StringSliceVar(&modeType, "mode-type", nil, "filter by mode type")
	cmd.Flags().StringSliceVar(&threatType, "threat-type", nil, "filter by threat type")
	cmd.Flags().StringVar(&outFile, "out", "", "write export to file (default: stdout)")
	return cmd
}
