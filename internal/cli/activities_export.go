package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newActivitiesExportCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var activityTypes []int
	var start, end, output string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export activities as CSV",
		Long:  "Bulk export the activity log as CSV. Output goes to stdout by default, or to a file with --out.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ActivityExportParams{
				SiteIDs:       siteIDs,
				AccountIDs:    accountIDs,
				GroupIDs:      groupIDs,
				ActivityTypes: activityTypes,
				CreatedAtGt:   start,
				CreatedAtLt:   end,
			}
			data, err := c.ActivitiesExport(cmd.Context(), params)
			if err != nil {
				return err
			}

			if output != "" {
				if err := os.WriteFile(output, data, 0o644); err != nil {
					return fmt.Errorf("write %s: %w", output, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Exported %d bytes to %s\n", len(data), output)
				return nil
			}

			_, err = cmd.OutOrStdout().Write(data)
			return err
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().IntSliceVar(&activityTypes, "activity-type", nil, "filter by activity type ID")
	cmd.Flags().StringVar(&start, "start", "", "activities after this timestamp (ISO 8601)")
	cmd.Flags().StringVar(&end, "end", "", "activities before this timestamp (ISO 8601)")
	cmd.Flags().StringVar(&output, "out", "", "write to file instead of stdout")
	return cmd
}
