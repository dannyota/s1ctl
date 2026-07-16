package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addTagRuleSyncCmds(parent *cobra.Command) {
	spec := tagRulesSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// tagRuleSyncFile is the YAML representation of a tag rule on disk. It holds
// only the declarative fields; server-assigned ID, scope IDs, audit fields, and
// timestamps are omitted. The nested condition/scope/tag structures are kept as
// generic values so they round-trip through YAML natively (the SDK models them
// as raw JSON).
type tagRuleSyncFile struct {
	Name           string `yaml:"name"`
	Description    string `yaml:"description,omitempty"`
	Status         string `yaml:"status,omitempty"`
	Conditions     any    `yaml:"conditions,omitempty"`
	Scopes         any    `yaml:"scopes,omitempty"`
	Tags           any    `yaml:"tags,omitempty"`
	ExcludedAssets any    `yaml:"excludedAssets,omitempty"`
}

func tagRuleToSyncFile(r mgmt.TagRule) (tagRuleSyncFile, error) {
	f := tagRuleSyncFile{
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
	}
	blobs := []struct {
		raw json.RawMessage
		dst *any
	}{
		{r.Conditions, &f.Conditions},
		{r.Scopes, &f.Scopes},
		{r.Tags, &f.Tags},
		{r.ExcludedAssets, &f.ExcludedAssets},
	}
	for _, b := range blobs {
		if len(b.raw) == 0 {
			continue
		}
		var v any
		if err := json.Unmarshal(b.raw, &v); err != nil {
			return f, fmt.Errorf("tag rule %s: %w", r.Name, err)
		}
		*b.dst = v
	}
	return f, nil
}

func (f tagRuleSyncFile) toWrite() (mgmt.TagRuleWrite, error) {
	w := mgmt.TagRuleWrite{
		Name:        f.Name,
		Description: f.Description,
		Status:      f.Status,
	}
	fields := []struct {
		src any
		dst *json.RawMessage
	}{
		{f.Conditions, &w.Conditions},
		{f.Scopes, &w.Scopes},
		{f.Tags, &w.Tags},
		{f.ExcludedAssets, &w.ExcludedAssets},
	}
	for _, p := range fields {
		if p.src == nil {
			continue
		}
		b, err := json.Marshal(p.src)
		if err != nil {
			return w, fmt.Errorf("tag rule %s: %w", f.Name, err)
		}
		*p.dst = b
	}
	return w, nil
}

// decodeTagRule maps one local file to a canonical Object. Identity is the rule
// name; the body is re-marshalled so it is byte-equal to what List produces.
func decodeTagRule(data []byte) (reconcile.Object, error) {
	var f tagRuleSyncFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Name == "" {
		return reconcile.Object{}, fmt.Errorf("tag rule has no name")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: f.Name, Body: body}, nil
}

// tagRulesSpec adapts dynamic tag rules to the shared sync builders. Identity is
// the rule name.
func tagRulesSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "tag rule",
		Command:    "tag-rules",
		DefaultDir: "tag-rules",
		PullShort:  "Pull dynamic tag rules to local YAML files",
		PullLong: `Fetch all dynamic tag rules and write them as YAML files.

Each rule produces one file named by its sanitized name. Server-only metadata
(ID, scope IDs, audit fields, timestamps) is omitted so the files contain only
the declarative definition: name, status, conditions, scopes, tags, and excluded
assets.`,
		PushShort: "Push dynamic tag rules from local YAML files",
		PushLong: `Read tag rule YAML files from a directory and sync them to SentinelOne.

Rules are matched by name: existing rules are updated, new ones are created, and
unchanged ones are skipped. Dry-run by default — pass --yes to apply changes.`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
			cmd.Flags().StringSliceVar(&scope.AccountIDs, "account-id", nil, "filter by account ID")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
			cmd.Flags().StringSliceVar(&scope.AccountIDs, "account-id", nil, "filter by account ID")
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
				Name:    "tag rule",
				Command: "tag-rules",
				Decode:  decodeTagRule,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					params := &mgmt.TagRuleListParams{Limit: 1000}
					if !scope.push {
						params.SiteIDs = scope.SiteIDs
						params.AccountIDs = scope.AccountIDs
					}
					rules, _, lErr := fetchAllREST("tag rule", func(cur string) ([]mgmt.TagRule, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.TagRulesList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(rules))
					for _, r := range rules {
						f, fErr := tagRuleToSyncFile(r)
						if fErr != nil {
							return nil, fErr
						}
						body, mErr := yaml.Marshal(f)
						if mErr != nil {
							return nil, fmt.Errorf("marshal tag rule %s: %w", r.Name, mErr)
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
					var f tagRuleSyncFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					w, wErr := f.toWrite()
					if wErr != nil {
						return wErr
					}
					_, cErr := c.TagRulesCreate(ctx, w)
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f tagRuleSyncFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					w, wErr := f.toWrite()
					if wErr != nil {
						return wErr
					}
					w.ID = id
					_, uErr := c.TagRulesUpdate(ctx, w)
					return uErr
				},
			}, nil
		},
	}
}
