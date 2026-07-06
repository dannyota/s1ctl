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

func newRulesValidateCmd() *cobra.Command {
	var inDir string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate rule YAML files without deploying",
		Long: `Read rule YAML files from a directory and check for errors:
missing required fields, invalid enum values, and empty queries.
No API calls are made — this is a local-only check.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, err := os.ReadDir(inDir)
			if err != nil {
				return fmt.Errorf("read %s: %w", inDir, err)
			}

			var files int
			var errs []string
			for _, e := range entries {
				if e.IsDir() || (!strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml")) {
					continue
				}
				files++
				data, rErr := os.ReadFile(filepath.Join(inDir, e.Name()))
				if rErr != nil {
					errs = append(errs, fmt.Sprintf("%s: %v", e.Name(), rErr))
					continue
				}
				var rf ruleFile
				if uErr := yaml.Unmarshal(data, &rf); uErr != nil {
					errs = append(errs, fmt.Sprintf("%s: invalid YAML: %v", e.Name(), uErr))
					continue
				}
				errs = append(errs, validateRule(e.Name(), rf)...)
			}

			if files == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No rule files found in %s\n", inDir)
				return nil
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]any{
					"files":  files,
					"errors": errs,
					"valid":  len(errs) == 0,
				})
			}

			if len(errs) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "%s validated, no errors\n", pluralize(files, "file"))
				return nil
			}
			for _, e := range errs {
				fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", e)
			}
			return fmt.Errorf("%d validation %s in %s",
				len(errs), pluralWord(len(errs), "error"), pluralize(files, "file"))
		},
	}
	cmd.Flags().StringVar(&inDir, "dir", "rules", "directory containing rule YAML files")
	return markJSON(cmd)
}

var validSeverities = map[mgmt.RuleSeverity]bool{
	mgmt.RuleSeverityInfo: true, mgmt.RuleSeverityLow: true,
	mgmt.RuleSeverityMedium: true, mgmt.RuleSeverityHigh: true,
	mgmt.RuleSeverityCritical: true,
}

var validQueryTypes = map[mgmt.RuleQueryType]bool{
	mgmt.RuleQueryTypeEvents: true, mgmt.RuleQueryTypeCorrelation: true,
	mgmt.RuleQueryTypeUEBAFirstSeen: true, mgmt.RuleQueryTypeScheduled: true,
}

var validExpirations = map[mgmt.RuleExpirationMode]bool{
	mgmt.RuleExpirationPermanent: true, mgmt.RuleExpirationTemporary: true,
}

func validateRule(filename string, rf ruleFile) []string {
	var errs []string
	add := func(msg string) { errs = append(errs, fmt.Sprintf("%s: %s", filename, msg)) }

	if rf.Name == "" {
		add("missing name")
	}
	if rf.S1QL == "" && rf.QueryType != mgmt.RuleQueryTypeCorrelation {
		add("missing s1ql query")
	}
	if rf.Severity != "" && !validSeverities[rf.Severity] {
		add(fmt.Sprintf("invalid severity %q", rf.Severity))
	}
	if rf.QueryType != "" && !validQueryTypes[rf.QueryType] {
		add(fmt.Sprintf("invalid queryType %q", rf.QueryType))
	}
	if rf.ExpirationMode != "" && !validExpirations[rf.ExpirationMode] {
		add(fmt.Sprintf("invalid expirationMode %q", rf.ExpirationMode))
	}
	if rf.ExpirationMode == mgmt.RuleExpirationTemporary && rf.Expiration == "" {
		add("temporary rule missing expiration date")
	}
	return errs
}

func pluralWord(n int, singular string) string {
	if n == 1 {
		return singular
	}
	return singular + "s"
}
