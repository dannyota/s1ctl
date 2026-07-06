package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newDLPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dlp",
		Short: "Manage Data Loss Prevention (DLP) rules and classifications",
		Long: `Manage Data Loss Prevention (DLP) via the cloud security GraphQL API.

This surface reads data protection rules and DLP classifications, toggles and
deletes rules, deletes classifications, and shows engine settings. Rule and
classification bodies are large; creating and updating them is not yet exposed.

DLP list queries are page-based (--page/--limit), not cursor-based.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDLPRulesCmd())
	cmd.AddCommand(newDLPClassificationsCmd())
	cmd.AddCommand(newDLPSettingsCmd())
	return cmd
}

func newDLPRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage data protection rules",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDLPRulesListCmd())
	cmd.AddCommand(newDLPRulesGetCmd())
	addDLPRuleActions(cmd)
	return cmd
}

func newDLPClassificationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "classifications",
		Aliases: []string{"class"},
		Short:   "Manage DLP classifications",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDLPClassificationsListCmd())
	cmd.AddCommand(newDLPClassificationsGetCmd())
	cmd.AddCommand(newDLPClassificationDeleteCmd())
	return cmd
}

func newDLPSettingsCmd() *cobra.Command {
	var scopeLevel, scopeID string

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Show DLP engine settings for a scope",
		Long: `Show DLP engine settings. A scope is required by the API: pass both
--scope-level and --scope-id.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			if scope == nil {
				return fmt.Errorf("--scope-level and --scope-id are required for dlp settings")
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			s, err := c.DLPEngineSettings(cmd.Context(), scope)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), s)
			}
			rows := [][]string{
				{"Publishing Enabled", boolIcon(s.PublishingEnabled)},
				{"Prevent Action", orDash(s.PreventAction)},
				{"Character Inspection", orDash(s.CharacterInspectionDepth)},
				{"Classifications To Inspect", fmt.Sprintf("%d", s.ClassificationsToInspect)},
				{"OCR Enabled", boolIcon(s.EnableOCR)},
				{"Mask Evidence", boolIcon(s.MaskEvidence)},
				{"Block Encrypted Archive", boolIcon(s.BlockEncryptedArchive)},
				{"Block USB Modifications", boolIcon(s.BlockUsbModifications)},
				{"Max Archive Levels", fmt.Sprintf("%d", s.MaxArchiveLevels)},
				{"Inspection Size Limit", fmt.Sprintf("%d", s.InspectionSizeLimit)},
				{"Max Inspected File Size", fmt.Sprintf("%d", s.MaxInspectedFileSize)},
				{"Ignore Keywords", joinOrDash(s.IgnoreKeywords)},
				{"Ignore Regexes", joinOrDash(s.IgnoreRegexes)},
				{"Notification Message", orDash(s.NotificationMessage)},
				{"Scope", orDash(s.Scope.Path)},
				{"Updated", orDash(s.UpdatedAt)},
				{"Updated By", orDash(s.UpdatedBy)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	addDLPScopeFlags(cmd, &scopeLevel, &scopeID)
	return markJSON(cmd)
}

// addDLPScopeFlags registers the shared --scope-level/--scope-id flags.
func addDLPScopeFlags(cmd *cobra.Command, level, id *string) {
	cmd.Flags().StringVar(level, "scope-level", "", "scope level (account, site, group)")
	cmd.Flags().StringVar(id, "scope-id", "", "account, site, or group ID")
}

// fetchAllDLP walks every numbered page of a DLP list query and returns the
// accumulated nodes plus the reported total. DLP paginates by page number, so
// this increments the page while the connection reports another page.
func fetchAllDLP[T any](resource string, fn func(page int) (*graphql.DLPConnection[T], error)) ([]T, int, error) {
	var all []T
	var total int
	for page := 1; ; page++ {
		conn, err := fn(page)
		if err != nil {
			clearProgress()
			return nil, 0, err
		}
		total = conn.PageInfo.TotalCount
		all = append(all, conn.Nodes...)
		printProgress(resource, len(all), total)
		if !conn.PageInfo.HasNextPage {
			break
		}
	}
	clearProgress()
	return all, total, nil
}
