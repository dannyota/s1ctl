package cli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addRuleSyncCmds(parent *cobra.Command) {
	spec := rulesSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
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

// decodeRule maps one local file to a canonical Object: it validates the rule
// name and re-marshals through ruleFile so bodies are byte-equal to the ones
// List produces.
func decodeRule(data []byte) (reconcile.Object, error) {
	var rf ruleFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return reconcile.Object{}, err
	}
	if rf.Name == "" {
		return reconcile.Object{}, fmt.Errorf("rule has no name")
	}
	body, err := yaml.Marshal(rf)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: rf.Name, Body: body}, nil
}

// rulesSpec adapts custom detection rules to the shared sync builders. Rules
// have no create-time scope flag; the file's Scope field is informational only
// (ruleFile.toCreate omits it), so create and update send the same payload.
func rulesSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "rule",
		Command:    "rules",
		DefaultDir: "rules",
		PullShort:  "Pull custom detection rules to local YAML files",
		PushShort:  "Push custom detection rules from local YAML files",
		PushLong: `Read rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created,
and unchanged rules are skipped. Dry-run by default — pass --yes to apply changes.`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
		},
		Build: func(_ *cobra.Command, scope scopeFlags) (reconcile.Surface, error) {
			var client *mgmt.Client
			getClient := func() (*mgmt.Client, error) {
				if client == nil {
					c, err := mgmtClient()
					if err != nil {
						return nil, err
					}
					client = c
				}
				return client, nil
			}

			return reconcile.Surface{
				Name:    "rule",
				Command: "rules",
				Decode:  decodeRule,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					// Pull filters by --site-id; push lists every rule so
					// matching spans the whole tenant.
					params := &mgmt.RuleListParams{Limit: 1000}
					if !scope.push {
						params.SiteIDs = scope.SiteIDs
					}
					rules, _, lErr := fetchAllREST("rule", func(cur string) ([]mgmt.Rule, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.RulesList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(rules))
					for _, r := range rules {
						body, mErr := yaml.Marshal(ruleToFile(r))
						if mErr != nil {
							return nil, fmt.Errorf("marshal rule %s: %w", r.Name, mErr)
						}
						objs = append(objs, reconcile.Object{Name: r.Name, ID: r.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var rf ruleFile
					if uErr := yaml.Unmarshal(local.Body, &rf); uErr != nil {
						return uErr
					}
					_, cErr := c.RulesCreate(ctx, rf.toCreate())
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var rf ruleFile
					if uErr := yaml.Unmarshal(local.Body, &rf); uErr != nil {
						return uErr
					}
					_, uErr := c.RulesUpdate(ctx, id, rf.toCreate())
					return uErr
				},
			}, nil
		},
	}
}
