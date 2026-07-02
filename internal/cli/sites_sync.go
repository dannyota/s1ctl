package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addSiteSyncCmds(parent *cobra.Command) {
	parent.AddCommand(newSitesPullCmd())
	parent.AddCommand(newSitesPushCmd())
}

func newSitesPullCmd() *cobra.Command {
	var accountIDs []string
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull sites to a local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			sites, _, err := fetchAllREST("sites", func(cursor string) ([]mgmt.Site, *mgmt.Pagination, error) {
				return c.SitesList(cmd.Context(), &mgmt.SiteListParams{AccountIDs: accountIDs, Limit: 1000, Cursor: cursor})
			})
			if err != nil {
				return err
			}
			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}
			path := filepath.Join(outDir, "sites.json")
			data, err := json.MarshalIndent(sites, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(len(sites), "site"), path)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringVar(&outDir, "out", "samples", "output directory")
	return cmd
}

func newSitesPushCmd() *cobra.Command {
	var inFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push --file <sites.json>",
		Short: "Create sites from a local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if inFile == "" {
				return fmt.Errorf("--file is required")
			}
			data, err := os.ReadFile(inFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", inFile, err)
			}
			var sites []mgmt.SiteCreate
			if err := json.Unmarshal(data, &sites); err != nil {
				return fmt.Errorf("parse %s: %w", inFile, err)
			}
			return guard(cmd.OutOrStdout(), "sites push", fmt.Sprintf("create %s from %s", pluralize(len(sites), "site"), inFile), inFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				var created int
				for _, s := range sites {
					if _, cErr := c.SitesCreate(cmd.Context(), s); cErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", cErr)
						continue
					}
					created++
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"created": created})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", pluralize(created, "site"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&inFile, "file", "", "JSON file with an array of sites (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
