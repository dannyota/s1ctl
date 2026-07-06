package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/mgmt"
)

func newRulesDiffCmd() *cobra.Command {
	var inDir string

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Compare local rule YAML files against live rules",
		Long: `Read rule YAML files from a directory, fetch corresponding live rules
by name, and show what differs. Helps review changes before pushing.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, err := os.ReadDir(inDir)
			if err != nil {
				return fmt.Errorf("read %s: %w", inDir, err)
			}

			localByName := make(map[string]ruleFile)
			for _, e := range entries {
				if e.IsDir() || (!strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml")) {
					continue
				}
				data, rErr := os.ReadFile(filepath.Join(inDir, e.Name()))
				if rErr != nil {
					return fmt.Errorf("read %s: %w", e.Name(), rErr)
				}
				var rf ruleFile
				if uErr := yaml.Unmarshal(data, &rf); uErr != nil {
					return fmt.Errorf("parse %s: %w", e.Name(), uErr)
				}
				if rf.Name != "" {
					localByName[rf.Name] = rf
				}
			}

			if len(localByName) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No rule files found in %s\n", inDir)
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.RuleListParams{Limit: 1000}
			remoteRules, _, err := fetchAllREST("rule", func(cur string) ([]mgmt.Rule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.RulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			remoteByName := make(map[string]mgmt.Rule, len(remoteRules))
			for _, r := range remoteRules {
				remoteByName[r.Name] = r
			}

			type diffResult struct {
				Name   string   `json:"name"`
				Status string   `json:"status"`
				Diffs  []string `json:"diffs,omitempty"`
			}

			var results []diffResult
			var added, modified, unchanged int

			for name, local := range localByName {
				remote, exists := remoteByName[name]
				if !exists {
					results = append(results, diffResult{Name: name, Status: "new"})
					added++
					continue
				}

				var diffs []string
				if local.S1QL != remote.S1QL {
					diffs = append(diffs, "s1ql")
				}
				if local.Severity != remote.Severity {
					diffs = append(diffs, fmt.Sprintf("severity: %s → %s", remote.Severity, local.Severity))
				}
				if local.Status != remote.Status {
					diffs = append(diffs, fmt.Sprintf("status: %s → %s", remote.Status, local.Status))
				}
				if local.Description != remote.Description {
					diffs = append(diffs, "description")
				}
				if local.QueryType != remote.QueryType {
					diffs = append(diffs, fmt.Sprintf("queryType: %s → %s", remote.QueryType, local.QueryType))
				}
				if local.ExpirationMode != remote.ExpirationMode {
					diffs = append(diffs, fmt.Sprintf("expirationMode: %s → %s", remote.ExpirationMode, local.ExpirationMode))
				}
				if local.TreatAsThreat != remote.TreatAsThreat {
					diffs = append(diffs, fmt.Sprintf("treatAsThreat: %s → %s", remote.TreatAsThreat, local.TreatAsThreat))
				}
				if local.NetworkQuarantine != remote.NetworkQuarantine {
					diffs = append(diffs, fmt.Sprintf("networkQuarantine: %v → %v", remote.NetworkQuarantine, local.NetworkQuarantine))
				}

				if len(diffs) > 0 {
					results = append(results, diffResult{Name: name, Status: "modified", Diffs: diffs})
					modified++
				} else {
					results = append(results, diffResult{Name: name, Status: "unchanged"})
					unchanged++
				}
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), results)
			}

			headers := []string{"Rule", "Status", "Changes"}
			rows := make([][]string, len(results))
			for i, r := range results {
				changes := "-"
				if len(r.Diffs) > 0 {
					changes = strings.Join(r.Diffs, ", ")
				}
				rows[i] = []string{truncate(r.Name, 40), r.Status, truncate(changes, 50)}
			}
			printTable(headers, rows)

			fmt.Fprintf(cmd.OutOrStdout(), "\n%s: %d new, %d modified, %d unchanged\n",
				pluralize(len(results), "rule"), added, modified, unchanged)
			return nil
		},
	}
	cmd.Flags().StringVar(&inDir, "dir", "rules", "directory containing rule YAML files")
	return markJSON(cmd)
}
