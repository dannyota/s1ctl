package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addExclusionSyncCmds(parent *cobra.Command) {
	spec := exclusionsSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// exclusionFile is the YAML representation of an exclusion on disk. It holds
// only the declarative fields; server-assigned IDs, scope, source, and
// timestamps are omitted.
type exclusionFile struct {
	Type              string `yaml:"type"`
	Value             string `yaml:"value"`
	OSType            string `yaml:"osType"`
	Mode              string `yaml:"mode,omitempty"`
	Description       string `yaml:"description,omitempty"`
	PathExclusionType string `yaml:"pathExclusionType,omitempty"`
}

func exclusionToFile(e mgmt.Exclusion) exclusionFile {
	return exclusionFile{
		Type:              e.Type,
		Value:             e.Value,
		OSType:            e.OSType,
		Mode:              e.Mode,
		Description:       e.Description,
		PathExclusionType: e.PathExclusionType,
	}
}

func (f exclusionFile) toCreate() mgmt.ExclusionCreate {
	return mgmt.ExclusionCreate{
		Type:              f.Type,
		Value:             f.Value,
		OSType:            f.OSType,
		Mode:              f.Mode,
		Description:       f.Description,
		PathExclusionType: f.PathExclusionType,
	}
}

// exclusionIdentity is the engine matching key: type + OS + value, so the same
// value excluded on two operating systems stays a distinct object.
func exclusionIdentity(f exclusionFile) string {
	return f.Type + "/" + f.OSType + "/" + f.Value
}

// decodeExclusion maps one local file to a canonical Object: it validates the
// identity fields and re-marshals through exclusionFile so bodies are byte-equal
// to the ones List produces. Identity is type + OS + value.
func decodeExclusion(data []byte) (reconcile.Object, error) {
	var f exclusionFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Type == "" || f.Value == "" || f.OSType == "" {
		return reconcile.Object{}, fmt.Errorf("exclusion in file has no type/value/osType")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: exclusionIdentity(f), Body: body}, nil
}

// exclusionsSpec adapts exclusions to the shared sync builders. Identity is
// type + OS + value.
func exclusionsSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "exclusion",
		Command:    "exclusions",
		DefaultDir: "exclusions",
		PullShort:  "Pull exclusions to local YAML files",
		PullLong: `Fetch all exclusions and write them as YAML files.

Each exclusion produces one file. Server-only metadata (ID, scope, source,
timestamps) is omitted so the files contain only the declarative definition.`,
		PushShort: "Push exclusions from local YAML files",
		PushLong: `Read exclusion YAML files from a directory and sync them to SentinelOne.

Exclusions are matched by type + OS + value: existing exclusions are updated,
new exclusions are created, and unchanged exclusions are skipped. Dry-run by
default — pass --yes to apply changes.

New exclusions are created at the scope specified by --site-id. If no scope
flag is given, they are created at the global (tenant) scope.`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "scope for new exclusions (default: global/tenant)")
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
				Name:    "exclusion",
				Command: "exclusions",
				Decode:  decodeExclusion,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					// Pull filters by --site-id; push lists every exclusion so
					// matching spans the tenant (--site-id is the create scope).
					params := &mgmt.ExclusionListParams{Limit: 1000}
					if !scope.push {
						params.SiteIDs = scope.SiteIDs
					}
					exclusions, _, lErr := fetchAllREST("exclusion", func(cur string) ([]mgmt.Exclusion, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.ExclusionsList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(exclusions))
					for _, e := range exclusions {
						f := exclusionToFile(e)
						body, mErr := yaml.Marshal(f)
						if mErr != nil {
							return nil, fmt.Errorf("marshal exclusion %s: %w", exclusionIdentity(f), mErr)
						}
						objs = append(objs, reconcile.Object{Name: exclusionIdentity(f), ID: e.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f exclusionFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					// Creates are scoped by --site-id (preserving the legacy
					// push, which passed the site IDs to every create).
					_, cErr := c.ExclusionsCreate(ctx, scope.SiteIDs, f.toCreate())
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f exclusionFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, uErr := c.ExclusionsUpdate(ctx, id, f.toCreate())
					return uErr
				},
			}, nil
		},
	}
}
