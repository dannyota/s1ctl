package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addBlocklistSyncCmds(parent *cobra.Command) {
	spec := blocklistSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// blocklistFile is the YAML representation of a blocklist item on disk. It holds
// only the declarative fields; server-assigned IDs, scope, source, and
// timestamps are omitted.
type blocklistFile struct {
	Type        string `yaml:"type"`
	Value       string `yaml:"value"`
	OSType      string `yaml:"osType"`
	SHA256Value string `yaml:"sha256Value,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func blocklistToFile(b mgmt.BlocklistItem) blocklistFile {
	return blocklistFile{
		Type:        b.Type,
		Value:       b.Value,
		OSType:      b.OSType,
		SHA256Value: b.SHA256Value,
		Description: b.Description,
	}
}

func (f blocklistFile) toCreate() mgmt.BlocklistCreate {
	return mgmt.BlocklistCreate{
		Type:        mgmt.BlocklistType(f.Type),
		OSType:      mgmt.BlocklistOSType(f.OSType),
		Value:       f.Value,
		SHA256Value: f.SHA256Value,
		Description: f.Description,
	}
}

// blocklistIdentity is the engine matching key: type + OS + value, so the same
// hash blocked on two operating systems stays a distinct object.
func blocklistIdentity(f blocklistFile) string {
	return f.Type + "/" + f.OSType + "/" + f.Value
}

// decodeBlocklist maps one local file to a canonical Object: it validates the
// identity fields and re-marshals through blocklistFile so bodies are byte-equal
// to the ones List produces. Identity is type + OS + value.
func decodeBlocklist(data []byte) (reconcile.Object, error) {
	var f blocklistFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Type == "" || f.Value == "" || f.OSType == "" {
		return reconcile.Object{}, fmt.Errorf("blocklist item in file has no type/value/osType")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: blocklistIdentity(f), Body: body}, nil
}

// blocklistSpec adapts the blocklist to the shared sync builders. Identity is
// type + OS + value, mirroring the exclusions surface.
func blocklistSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "blocklist",
		Command:    "blocklist",
		DefaultDir: "blocklist",
		PullShort:  "Pull blocklist items to local YAML files",
		PullLong: `Fetch all blocklist items and write them as YAML files.

Each item produces one file. Server-only metadata (ID, scope, source,
timestamps) is omitted so the files contain only the declarative definition.`,
		PushShort: "Push blocklist items from local YAML files",
		PushLong: `Read blocklist YAML files from a directory and sync them to SentinelOne.

Items are matched by type + OS + value: existing items are updated, new items
are created, and unchanged items are skipped. Dry-run by default — pass --yes
to apply changes.

New items are created at the scope specified by --site-id. If no scope flag is
given, they are created at the global (tenant) scope.`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "scope for new items (default: global/tenant)")
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
				Name:    "blocklist",
				Command: "blocklist",
				Decode:  decodeBlocklist,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					// Pull filters by --site-id; push lists every item so
					// matching spans the tenant (--site-id is the create scope).
					params := &mgmt.BlocklistListParams{Limit: 1000}
					if !scope.push {
						params.SiteIDs = scope.SiteIDs
					}
					items, _, lErr := fetchAllREST("blocklist item", func(cur string) ([]mgmt.BlocklistItem, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.BlocklistList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(items))
					for _, b := range items {
						f := blocklistToFile(b)
						body, mErr := yaml.Marshal(f)
						if mErr != nil {
							return nil, fmt.Errorf("marshal blocklist item %s: %w", blocklistIdentity(f), mErr)
						}
						objs = append(objs, reconcile.Object{Name: blocklistIdentity(f), ID: b.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f blocklistFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, cErr := c.BlocklistCreate(ctx, blocklistScope(scope.SiteIDs, nil, nil), f.toCreate())
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f blocklistFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, uErr := c.BlocklistUpdate(ctx, id, mgmt.BlocklistScope{}, f.toCreate())
					return uErr
				},
			}, nil
		},
	}
}
