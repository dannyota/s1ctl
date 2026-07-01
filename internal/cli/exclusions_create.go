package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func init() {
	// Registered via addExclusionCreateCmd in exclusions_list.go
}

func newExclusionsCreateCmd() *cobra.Command {
	var (
		exclType    string
		value       string
		osType      string
		mode        string
		description string
		pathType    string
		siteIDs     []string
		groupIDs    []string
		yes         bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an exclusion",
		Long: `Create a new exclusion entry.

Types: path, file_type, white_hash, browser, certificate, document_type
OS types: windows, linux, macos, windows_legacy
Modes: suppress, suppress_dynamic_only, suppress_app_control

For path exclusions, --path-type specifies the match type:
  subfolders (default), file, glob`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if exclType == "" {
				return fmt.Errorf("--type is required")
			}
			if value == "" {
				return fmt.Errorf("--value is required")
			}
			if osType == "" {
				return fmt.Errorf("--os-type is required")
			}

			excl := mgmt.ExclusionCreate{
				Type:              exclType,
				Value:             value,
				OSType:            osType,
				Mode:              mode,
				Description:       description,
				PathExclusionType: pathType,
				GroupIDs:          groupIDs,
			}

			action := fmt.Sprintf("create %s exclusion for %q (%s)", exclType, value, osType)
			return guard(cmd.OutOrStdout(), "exclusions create", action, value, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.ExclusionsCreate(cmd.Context(), siteIDs, excl)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created exclusion %s (%s: %s)\n", created.ID, created.Type, created.Value)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&exclType, "type", "", "exclusion type (path, file_type, white_hash, browser, certificate, document_type)")
	cmd.Flags().StringVar(&value, "value", "", "exclusion value (path, hash, extension, etc.)")
	cmd.Flags().StringVar(&osType, "os-type", "", "target OS (windows, linux, macos)")
	cmd.Flags().StringVar(&mode, "mode", "suppress", "exclusion mode (suppress, suppress_dynamic_only, suppress_app_control)")
	cmd.Flags().StringVar(&description, "description", "", "exclusion description")
	cmd.Flags().StringVar(&pathType, "path-type", "", "path exclusion type (subfolders, file, glob)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "target site IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "target group IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
