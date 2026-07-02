package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func addCloudPolicySyncCmds(parent *cobra.Command) {
	parent.AddCommand(newCloudPoliciesPullCmd())
	parent.AddCommand(newCloudPoliciesPushCmd())
}

func newCloudPoliciesPullCmd() *cobra.Command {
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull cloud security policies to a local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			policies, _, err := fetchAllGQL("cloud policies", func(after string) (*graphql.Connection[graphql.CloudPolicy], error) {
				return c.CloudPoliciesList(cmd.Context(), &graphql.ListParams{First: 100, After: after})
			})
			if err != nil {
				return err
			}
			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}
			path := filepath.Join(outDir, "cloud-policies.json")
			data, err := json.MarshalIndent(policies, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(path, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %d cloud policies to %s\n", len(policies), path)
			return nil
		},
	}
	cmd.Flags().StringVar(&outDir, "out", "samples", "output directory")
	return cmd
}

func newCloudPoliciesPushCmd() *cobra.Command {
	var inFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push --file <cloud-policies.json>",
		Short: "Apply enabled/disabled policy status from a local file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if inFile == "" {
				return fmt.Errorf("--file is required")
			}
			data, err := os.ReadFile(inFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", inFile, err)
			}
			var desired []graphql.CloudPolicy
			if err := json.Unmarshal(data, &desired); err != nil {
				return fmt.Errorf("parse %s: %w", inFile, err)
			}
			var enable, disable []string
			for _, p := range desired {
				switch strings.ToLower(p.Status) {
				case "enabled":
					enable = append(enable, p.ID)
				case "disabled":
					disable = append(disable, p.ID)
				default:
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: policy %s has unrecognized status %q, skipping\n", p.ID, p.Status)
				}
			}
			action := fmt.Sprintf("enable %d and disable %d cloud policies from %s", len(enable), len(disable), inFile)
			return guard(cmd.OutOrStdout(), "cloud-policies push", action, inFile, yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				if len(enable) > 0 {
					if _, err := c.CloudPoliciesEnable(cmd.Context(), enable); err != nil {
						return err
					}
				}
				if len(disable) > 0 {
					if _, err := c.CloudPoliciesDisable(cmd.Context(), disable); err != nil {
						return err
					}
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"enabled": len(enable), "disabled": len(disable)})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Enabled %d, disabled %d cloud policies\n", len(enable), len(disable))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&inFile, "file", "", "JSON file with an array of policies (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
