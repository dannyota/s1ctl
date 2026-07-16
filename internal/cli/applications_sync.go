package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addAppControlSyncCmds(parent *cobra.Command) {
	spec := appControlRulesSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// appControlRuleFile is the YAML representation of an application control rule
// on disk. Only declarative fields; server-assigned IDs, scope metadata, and
// timestamps are omitted.
type appControlRuleFile struct {
	RuleName    string                     `yaml:"ruleName"`
	Description string                     `yaml:"description,omitempty"`
	Behavior    string                     `yaml:"behavior"`
	OSType      []string                   `yaml:"osType,omitempty"`
	Propagation bool                       `yaml:"propagation"`
	Parameters  *appControlConditionsFile  `yaml:"parameters,omitempty"`
	Exceptions  []appControlConditionsFile `yaml:"exceptions,omitempty"`
}

// appControlConditionsFile is the YAML representation of rule conditions.
// Only writable fields are included; applicationVersion is read-only in the
// API and excluded so pull/push round-trips cleanly.
type appControlConditionsFile struct {
	Publisher     string `yaml:"publisher,omitempty"`
	Path          string `yaml:"path,omitempty"`
	Signer        string `yaml:"signer,omitempty"`
	SHA256        string `yaml:"sha256,omitempty"`
	Process       string `yaml:"process,omitempty"`
	ParentProcess string `yaml:"parentProcess,omitempty"`
}

func appControlRuleToFile(r mgmt.AppControlRule) appControlRuleFile {
	f := appControlRuleFile{
		RuleName:    r.RuleName,
		Description: r.Description,
		Behavior:    string(r.Behavior),
		Propagation: r.Propagation,
	}
	for _, o := range r.OSType {
		f.OSType = append(f.OSType, string(o))
	}
	if r.Parameters != nil {
		p := appControlConditionsFile{
			Publisher:     r.Parameters.Publisher,
			Path:          r.Parameters.Path,
			Signer:        r.Parameters.Signer,
			SHA256:        r.Parameters.SHA256,
			Process:       r.Parameters.Process,
			ParentProcess: r.Parameters.ParentProcess,
		}
		f.Parameters = &p
	}
	for _, e := range r.Exceptions {
		f.Exceptions = append(f.Exceptions, appControlConditionsFile{
			Publisher:     e.Publisher,
			Path:          e.Path,
			Signer:        e.Signer,
			SHA256:        e.SHA256,
			Process:       e.Process,
			ParentProcess: e.ParentProcess,
		})
	}
	return f
}

func (f appControlRuleFile) toInput(scopeType string, scopeIDs []string) mgmt.AppControlRuleInput {
	desc := f.Description
	input := mgmt.AppControlRuleInput{
		RuleName:    f.RuleName,
		Description: &desc,
		Behavior:    mgmt.AppControlBehavior(f.Behavior),
		Propagation: &f.Propagation,
	}
	for _, o := range f.OSType {
		input.OSType = append(input.OSType, mgmt.AppControlOSType(o))
	}
	if f.Parameters != nil {
		input.Parameters = &mgmt.AppControlConditionsInput{
			Publisher:     f.Parameters.Publisher,
			Path:          f.Parameters.Path,
			Signer:        f.Parameters.Signer,
			SHA256:        f.Parameters.SHA256,
			Process:       f.Parameters.Process,
			ParentProcess: f.Parameters.ParentProcess,
		}
	}
	for _, e := range f.Exceptions {
		input.Exceptions = append(input.Exceptions, mgmt.AppControlConditionsInput{
			Publisher:     e.Publisher,
			Path:          e.Path,
			Signer:        e.Signer,
			SHA256:        e.SHA256,
			Process:       e.Process,
			ParentProcess: e.ParentProcess,
		})
	}
	if scopeType != "" && len(scopeIDs) > 0 {
		input.Scope = &mgmt.AppControlScope{
			ScopeType: mgmt.AppControlScopeLevel(strings.ToUpper(scopeType)),
			ScopeIDs:  scopeIDs,
		}
	}
	return input
}

// decodeAppControlRule maps one local file to a canonical Object.
func decodeAppControlRule(data []byte) (reconcile.Object, error) {
	var f appControlRuleFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.RuleName == "" {
		return reconcile.Object{}, fmt.Errorf("application control rule has no ruleName")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: f.RuleName, Body: body}, nil
}

// appControlRulesSpec adapts application control rules to the shared sync
// builders. Identity is the rule name.
func appControlRulesSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "application control rule",
		Command:    "applications-rules",
		DefaultDir: "appcontrol-rules",
		PullShort:  "Pull application control rules to local YAML files",
		PullLong: `Fetch all application control rules and write them as YAML files.

Each rule produces one file. Server-only metadata (ID, scope, timestamps)
is omitted so the files contain only the declarative definition.`,
		PushShort: "Push application control rules from local YAML files",
		PushLong: `Read application control rule YAML files from a directory and sync them to SentinelOne.

Rules are matched by name: existing rules are updated, new rules are created,
and unchanged rules are skipped. Dry-run by default — pass --yes to apply.`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().String("scope-type", "site", "scope type (account, site, group)")
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "scope-id", nil, "scope IDs to filter by")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().String("scope-type", "", "scope type (account, site, group) [required with --scope-id]")
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "scope-id", nil, "scope IDs for new rules")
		},
		Build: func(cmd *cobra.Command, scope scopeFlags) (reconcile.Surface, error) {
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

			scopeType := ""
			if f := cmd.Flags().Lookup("scope-type"); f != nil {
				scopeType = f.Value.String()
			}

			// On push, require --scope-type when --scope-id is given.
			if scope.push && len(scope.SiteIDs) > 0 && scopeType == "" {
				return reconcile.Surface{}, fmt.Errorf("--scope-type is required when --scope-id is set on push")
			}

			return reconcile.Surface{
				Name:    "application control rule",
				Command: "applications-rules",
				Decode:  decodeAppControlRule,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					params := &mgmt.AppControlQueryParams{
						IncludeParents: false,
						PageSize:       100,
					}
					if len(scope.SiteIDs) > 0 {
						params.ScopeType = mgmt.AppControlScopeLevel(strings.ToUpper(scopeType))
						if params.ScopeType == "" {
							params.ScopeType = mgmt.AppControlScopeSite
						}
						params.ScopeIDs = scope.SiteIDs
					}

					var allRules []mgmt.AppControlRule
					for {
						page, cursor, _, lErr := c.AppControlRulesList(ctx, params)
						if lErr != nil {
							return nil, lErr
						}
						allRules = append(allRules, page...)
						if cursor == "" || len(page) == 0 {
							break
						}
						params.Cursor = cursor
					}

					objs := make([]reconcile.Object, 0, len(allRules))
					for _, r := range allRules {
						f := appControlRuleToFile(r)
						body, mErr := yaml.Marshal(f)
						if mErr != nil {
							return nil, fmt.Errorf("marshal application control rule %s: %w", r.RuleName, mErr)
						}
						objs = append(objs, reconcile.Object{Name: r.RuleName, ID: r.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f appControlRuleFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, cErr := c.AppControlRulesCreate(ctx, f.toInput(scopeType, scope.SiteIDs))
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f appControlRuleFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, uErr := c.AppControlRulesUpdate(ctx, id, f.toInput("", nil))
					return uErr
				},
			}, nil
		},
	}
}
