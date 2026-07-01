package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addExclusionSyncCmds(parent *cobra.Command) {
	parent.AddCommand(newExclusionsPullCmd())
	parent.AddCommand(newExclusionsPushCmd())
}

func newExclusionsPullCmd() *cobra.Command {
	var siteIDs []string
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull exclusions to local files",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			var all []mgmt.Exclusion
			var cursor string
			for {
				exclusions, pag, lErr := c.ExclusionsList(cmd.Context(), &mgmt.ExclusionListParams{
					SiteIDs: siteIDs,
					Limit:   1000,
					Cursor:  cursor,
				})
				if lErr != nil {
					return lErr
				}
				all = append(all, exclusions...)
				if pag.NextCursor == "" {
					break
				}
				cursor = pag.NextCursor
			}

			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}
			path := filepath.Join(outDir, "exclusions.json")
			data, err := json.MarshalIndent(all, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(len(all), "exclusion"), path)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&outDir, "out", "samples", "output directory")
	return cmd
}

func newExclusionsPushCmd() *cobra.Command {
	var inFile string
	var siteIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push exclusions from local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, err := os.ReadFile(inFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", inFile, err)
			}
			var exclusions []mgmt.ExclusionCreate
			if err := json.Unmarshal(data, &exclusions); err != nil {
				return fmt.Errorf("parse %s: %w", inFile, err)
			}
			return guard(cmd.OutOrStdout(), "exclusions push", fmt.Sprintf("push %s from %s", pluralize(len(exclusions), "exclusion"), inFile), inFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				var created int
				for _, excl := range exclusions {
					if _, cErr := c.ExclusionsCreate(cmd.Context(), siteIDs, excl); cErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", cErr)
						continue
					}
					created++
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", pluralize(created, "exclusion"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&inFile, "file", "samples/exclusions.json", "input file")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "target site IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}
