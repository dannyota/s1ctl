package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addGroupSyncCmds(parent *cobra.Command) {
	parent.AddCommand(newGroupsPullCmd())
	parent.AddCommand(newGroupsPushCmd())
}

func newGroupsPullCmd() *cobra.Command {
	var siteIDs []string
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull groups to a local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			groups, _, err := fetchAllREST("groups", func(cursor string) ([]mgmt.Group, *mgmt.Pagination, error) {
				return c.GroupsList(cmd.Context(), &mgmt.GroupListParams{SiteIDs: siteIDs, Limit: 1000, Cursor: cursor})
			})
			if err != nil {
				return err
			}
			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}
			path := filepath.Join(outDir, "groups.json")
			data, err := json.MarshalIndent(groups, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(len(groups), "group"), path)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&outDir, "out", "samples", "output directory")
	return cmd
}

func newGroupsPushCmd() *cobra.Command {
	var inFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push --file <groups.json>",
		Short: "Create groups from a local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if inFile == "" {
				return fmt.Errorf("--file is required")
			}
			data, err := os.ReadFile(inFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", inFile, err)
			}
			var groups []mgmt.GroupCreate
			if err := json.Unmarshal(data, &groups); err != nil {
				return fmt.Errorf("parse %s: %w", inFile, err)
			}
			return guard(cmd.OutOrStdout(), "groups push", fmt.Sprintf("create %s from %s", pluralize(len(groups), "group"), inFile), inFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				var created int
				for _, g := range groups {
					if g.SiteID == "" {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: skipping group %q: missing siteId\n", g.Name)
						continue
					}
					if _, cErr := c.GroupsCreate(cmd.Context(), g.SiteID, g); cErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", cErr)
						continue
					}
					created++
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"created": created})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", pluralize(created, "group"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&inFile, "file", "", "JSON file with an array of groups (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
