package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addTagSyncCmds(parent *cobra.Command) {
	spec := tagsSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// tagFile is the YAML representation of a tag on disk. It holds only the
// declarative fields; server-assigned IDs and timestamps are omitted.
type tagFile struct {
	Key         string `yaml:"key"`
	Value       string `yaml:"value"`
	Description string `yaml:"description,omitempty"`
	Scope       string `yaml:"scope,omitempty"`
	ScopeID     string `yaml:"scopeId,omitempty"`
}

func tagToFile(t mgmt.Tag) tagFile {
	return tagFile{
		Key:         t.Key,
		Value:       t.Value,
		Description: t.Description,
		Scope:       t.Scope,
		ScopeID:     t.ScopeID,
	}
}

func (f tagFile) toCreate() mgmt.TagCreate {
	return mgmt.TagCreate{
		Key:         f.Key,
		Value:       f.Value,
		Description: f.Description,
		Scope:       f.Scope,
		ScopeID:     f.ScopeID,
	}
}

// toUpdate builds the pointer-field update body. Scope and ScopeID are set only
// at creation time and are not part of TagUpdate.
func (f tagFile) toUpdate() mgmt.TagUpdate {
	return mgmt.TagUpdate{
		Key:         &f.Key,
		Value:       &f.Value,
		Description: &f.Description,
	}
}

// decodeTag maps one local file to a canonical Object: it validates the tag key
// and re-marshals through tagFile so bodies are byte-equal to the ones List
// produces. Identity is the tag key.
func decodeTag(data []byte) (reconcile.Object, error) {
	var f tagFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Key == "" {
		return reconcile.Object{}, fmt.Errorf("tag has no key")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: f.Key, Body: body}, nil
}

// tagsSpec adapts tags to the shared sync builders. Identity is the tag key.
func tagsSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "tag",
		Command:    "tags",
		DefaultDir: "tags",
		PullShort:  "Pull tags to local YAML files",
		PullLong: `Fetch all tags and write them as YAML files.

Each tag produces one file named by its sanitized key. Server-only metadata
(ID, timestamps) is omitted so the files contain only the declarative
definition. Pull writes every live tag, including duplicates of the same key
as suffixed files; push and drift match by key, so resolve duplicate keys
before pushing.`,
		PushShort: "Push tags from local YAML files",
		PushLong: `Read tag YAML files from a directory and sync them to SentinelOne.

Tags are matched by key: existing tags are updated, new tags are created,
and unchanged tags are skipped. Dry-run by default — pass --yes to apply changes.`,
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
				Name:    "tag",
				Command: "tags",
				Decode:  decodeTag,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					// Pull filters by --site-id; push registers no scope flag,
					// so it lists every tag (matching spans the tenant).
					params := &mgmt.TagListParams{SiteIDs: scope.SiteIDs, Limit: 1000}
					tags, _, lErr := fetchAllREST("tag", func(cur string) ([]mgmt.Tag, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.TagsList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(tags))
					for _, t := range tags {
						body, mErr := yaml.Marshal(tagToFile(t))
						if mErr != nil {
							return nil, fmt.Errorf("marshal tag %s: %w", t.Key, mErr)
						}
						objs = append(objs, reconcile.Object{Name: t.Key, ID: t.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f tagFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, cErr := c.TagsCreate(ctx, f.toCreate())
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f tagFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, uErr := c.TagsUpdate(ctx, id, f.toUpdate())
					return uErr
				},
			}, nil
		},
	}
}
