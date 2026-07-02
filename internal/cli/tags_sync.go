package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addTagSyncCmds(parent *cobra.Command) {
	parent.AddCommand(newTagsPullCmd())
	parent.AddCommand(newTagsPushCmd())
}

func newTagsPullCmd() *cobra.Command {
	var siteIDs []string
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull tags to a local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			tags, _, err := fetchAllREST("tags", func(cursor string) ([]mgmt.Tag, *mgmt.Pagination, error) {
				return c.TagsList(cmd.Context(), &mgmt.TagListParams{SiteIDs: siteIDs, Limit: 1000, Cursor: cursor})
			})
			if err != nil {
				return err
			}
			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}
			path := filepath.Join(outDir, "tags.json")
			data, err := json.MarshalIndent(tags, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(len(tags), "tag"), path)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&outDir, "out", "samples", "output directory")
	return cmd
}

func newTagsPushCmd() *cobra.Command {
	var inFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push --file <tags.json>",
		Short: "Create tags from a local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if inFile == "" {
				return fmt.Errorf("--file is required")
			}
			data, err := os.ReadFile(inFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", inFile, err)
			}
			var tags []mgmt.TagCreate
			if err := json.Unmarshal(data, &tags); err != nil {
				return fmt.Errorf("parse %s: %w", inFile, err)
			}
			return guard(cmd.OutOrStdout(), "tags push", fmt.Sprintf("create %s from %s", pluralize(len(tags), "tag"), inFile), inFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				var created int
				for _, tag := range tags {
					if _, cErr := c.TagsCreate(cmd.Context(), tag); cErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", cErr)
						continue
					}
					created++
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"created": created})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", pluralize(created, "tag"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&inFile, "file", "", "JSON file with an array of tags (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
