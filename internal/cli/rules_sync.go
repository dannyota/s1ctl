package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/mgmt"
)

func addRuleSyncCmds(parent *cobra.Command) {
	parent.AddCommand(newRulesPullCmd())
	parent.AddCommand(newRulesPushCmd())
}

// ruleFile is the YAML representation of a custom detection rule on disk.
type ruleFile struct {
	Name              string                  `yaml:"name"`
	Description       string                  `yaml:"description,omitempty"`
	S1QL              string                  `yaml:"s1ql"`
	Severity          mgmt.RuleSeverity       `yaml:"severity"`
	Status            mgmt.RuleStatus         `yaml:"status"`
	QueryType         mgmt.RuleQueryType      `yaml:"queryType"`
	QueryLang         string                  `yaml:"queryLang,omitempty"`
	Scope             mgmt.RuleScope          `yaml:"scope,omitempty"`
	ExpirationMode    mgmt.RuleExpirationMode `yaml:"expirationMode"`
	Expiration        string                  `yaml:"expiration,omitempty"`
	TreatAsThreat     mgmt.RuleTreatAsThreat  `yaml:"treatAsThreat"`
	NetworkQuarantine bool                    `yaml:"networkQuarantine,omitempty"`
}

func ruleToFile(r mgmt.Rule) ruleFile {
	return ruleFile{
		Name:              r.Name,
		Description:       r.Description,
		S1QL:              r.S1QL,
		Severity:          r.Severity,
		Status:            r.Status,
		QueryType:         r.QueryType,
		QueryLang:         r.QueryLang,
		Scope:             r.Scope,
		ExpirationMode:    r.ExpirationMode,
		Expiration:        r.Expiration,
		TreatAsThreat:     r.TreatAsThreat,
		NetworkQuarantine: r.NetworkQuarantine,
	}
}

func (rf ruleFile) toCreate() mgmt.RuleCreate {
	return mgmt.RuleCreate{
		Name:              rf.Name,
		Description:       rf.Description,
		S1QL:              rf.S1QL,
		Severity:          rf.Severity,
		Status:            rf.Status,
		QueryType:         rf.QueryType,
		QueryLang:         rf.QueryLang,
		ExpirationMode:    rf.ExpirationMode,
		Expiration:        rf.Expiration,
		TreatAsThreat:     rf.TreatAsThreat,
		NetworkQuarantine: rf.NetworkQuarantine,
	}
}

var unsafeChars = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

// sanitizeFilename converts a rule name into a safe filename stem.
func sanitizeFilename(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = unsafeChars.ReplaceAllString(s, "")
	if s == "" {
		s = "rule"
	}
	return s
}

func newRulesPullCmd() *cobra.Command {
	var siteIDs []string
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull custom detection rules to local YAML files",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.RuleListParams{
				SiteIDs: siteIDs,
				Limit:   1000,
			}
			rules, _, err := fetchAllREST("rule", func(cur string) ([]mgmt.Rule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.RulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}

			// Track used filenames to avoid collisions.
			used := make(map[string]int)
			for _, r := range rules {
				stem := sanitizeFilename(r.Name)
				if n := used[stem]; n > 0 {
					stem = fmt.Sprintf("%s-%d", stem, n)
				}
				used[sanitizeFilename(r.Name)]++

				rf := ruleToFile(r)
				data, mErr := yaml.Marshal(rf)
				if mErr != nil {
					return fmt.Errorf("marshal rule %s: %w", r.Name, mErr)
				}
				path := filepath.Join(outDir, stem+".yaml")
				if wErr := os.WriteFile(path, data, 0o644); wErr != nil {
					return wErr
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(len(rules), "rule"), outDir)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&outDir, "out", "rules", "output directory")
	return cmd
}

func newRulesPushCmd() *cobra.Command {
	var inDir string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push custom detection rules from local YAML files",
		Long: `Read rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created.
Dry-run by default — pass --yes to apply changes.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, err := os.ReadDir(inDir)
			if err != nil {
				return fmt.Errorf("read %s: %w", inDir, err)
			}

			var localRules []ruleFile
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
				if rf.Name == "" {
					return fmt.Errorf("rule in %s has no name", e.Name())
				}
				localRules = append(localRules, rf)
			}
			if len(localRules) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No rule files found.")
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			// Fetch all remote rules to match by name.
			params := &mgmt.RuleListParams{Limit: 1000}
			remoteRules, _, err := fetchAllREST("rule", func(cur string) ([]mgmt.Rule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.RulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}
			byName := make(map[string]mgmt.Rule, len(remoteRules))
			for _, r := range remoteRules {
				byName[r.Name] = r
			}

			var toCreate, toUpdate []ruleFile
			var updateIDs []string
			for _, lr := range localRules {
				if existing, ok := byName[lr.Name]; ok {
					toUpdate = append(toUpdate, lr)
					updateIDs = append(updateIDs, existing.ID)
				} else {
					toCreate = append(toCreate, lr)
				}
			}

			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would create %s, update %s from %s. Pass --yes to apply.\n",
					pluralize(len(toCreate), "rule"),
					pluralize(len(toUpdate), "rule"),
					inDir)
				return nil
			}

			var created, updated int
			for _, lr := range toCreate {
				if _, cErr := c.RulesCreate(cmd.Context(), lr.toCreate()); cErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: create %s: %v\n", lr.Name, cErr)
					continue
				}
				created++
			}
			for i, lr := range toUpdate {
				if _, uErr := c.RulesUpdate(cmd.Context(), updateIDs[i], lr.toCreate()); uErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: update %s: %v\n", lr.Name, uErr)
					continue
				}
				updated++
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created %s, updated %s\n",
				pluralize(created, "rule"),
				pluralize(updated, "rule"))
			return nil
		},
	}
	cmd.Flags().StringVar(&inDir, "dir", "rules", "directory containing rule YAML files")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}
