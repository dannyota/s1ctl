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
type appControlConditionsFile struct {
	Publisher          string `yaml:"publisher,omitempty"`
	Path               string `yaml:"path,omitempty"`
	Signer             string `yaml:"signer,omitempty"`
	SHA256             string `yaml:"sha256,omitempty"`
	Process            string `yaml:"process,omitempty"`
	ParentProcess      string `yaml:"parentProcess,omitempty"`
	ApplicationVersion string `yaml:"applicationVersion,omitempty"`
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
			Publisher:          r.Parameters.Publisher,
			Path:               r.Parameters.Path,
			Signer:             r.Parameters.Signer,
			SHA256:             r.Parameters.SHA256,
			Process:            r.Parameters.Process,
			ParentProcess:      r.Parameters.ParentProcess,
			ApplicationVersion: r.Parameters.ApplicationVersion,
		}
		f.Parameters = &p
	}
	for _, e := range r.Exceptions {
		f.Exceptions = append(f.Exceptions, appControlConditionsFile{
			Publisher:          e.Publisher,
			Path:               e.Path,
			Signer:             e.Signer,
			SHA256:             e.SHA256,
			Process:            e.Process,
			ParentProcess:      e.ParentProcess,
			ApplicationVersion: e.ApplicationVersion,
		})
	}
	return f
}

func (f appControlRuleFile) toInput(scopeType string, scopeIDs []string) mgmt.AppControlRuleInput {
	input := mgmt.AppControlRuleInput{
		RuleName:    f.RuleName,
		Description: f.Description,
		Behavior:    mgmt.AppControlBehavior(f.Behavior),
		Propagation: &f.Propagation,
	}
	for _, o := range f.OSType {
		input.OSType = append(input.OSType, mgmt.AppControlOSType(o))
	}
	if f.Parameters != nil {
		input.Parameters = &mgmt.AppControlConditions{
			Publisher:          f.Parameters.Publisher,
			Path:               f.Parameters.Path,
			Signer:             f.Parameters.Signer,
			SHA256:             f.Parameters.SHA256,
			Process:            f.Parameters.Process,
			ParentProcess:      f.Parameters.ParentProcess,
			ApplicationVersion: f.Parameters.ApplicationVersion,
		}
	}
	for _, e := range f.Exceptions {
		input.Exceptions = append(input.Exceptions, mgmt.AppControlConditions{
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
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "scope-id", nil, "scope IDs to filter by")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
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

			// Resolve the scope type flag from the parent.
			scopeType := ""
			if f := cmd.Flag("scope-type"); f != nil {
				scopeType = f.Value.String()
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
					if !scope.push && len(scope.SiteIDs) > 0 {
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
