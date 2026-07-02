package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addGroupSyncCmds(parent *cobra.Command) {
	spec := groupsSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// groupFile is the YAML representation of a group on disk. It holds only the
// declarative fields; server-assigned IDs, rank, and timestamps are omitted.
type groupFile struct {
	Name        string `yaml:"name"`
	SiteID      string `yaml:"siteId"`
	Description string `yaml:"description,omitempty"`
}

func groupToFile(g mgmt.Group) groupFile {
	return groupFile{
		Name:        g.Name,
		SiteID:      g.SiteID,
		Description: g.Description,
	}
}

func (f groupFile) toCreate() mgmt.GroupCreate {
	return mgmt.GroupCreate{
		Name:        f.Name,
		SiteID:      f.SiteID,
		Description: f.Description,
	}
}

// toUpdate builds the pointer-field update body. SiteID is part of the identity
// (a moved group plans as create + live-only), so GroupUpdate carries only the
// mutable name and description.
func (f groupFile) toUpdate() mgmt.GroupUpdate {
	return mgmt.GroupUpdate{
		Name:        &f.Name,
		Description: &f.Description,
	}
}

// groupIdentity is the engine matching key for a group: site ID plus name, so
// the same group name under two sites stays distinct.
func groupIdentity(f groupFile) string {
	return f.SiteID + "/" + f.Name
}

// decodeGroup maps one local file to a canonical Object: it validates the group
// name and re-marshals through groupFile so bodies are byte-equal to the ones
// List produces. Identity is site ID + name.
func decodeGroup(data []byte) (reconcile.Object, error) {
	var f groupFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Name == "" {
		return reconcile.Object{}, fmt.Errorf("group has no name")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: groupIdentity(f), Body: body}, nil
}

// groupsSpec adapts groups to the shared sync builders. Identity is site ID +
// group name.
func groupsSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "group",
		Command:    "groups",
		DefaultDir: "groups",
		PullShort:  "Pull groups to local YAML files",
		PullLong: `Fetch all groups and write them as YAML files.

Each group produces one file. Server-only metadata (ID, rank, agent counts,
timestamps) is omitted so the files contain only the declarative definition.`,
		PushShort: "Push groups from local YAML files",
		PushLong: `Read group YAML files from a directory and sync them to SentinelOne.

Groups are matched by site ID + name: existing groups are updated, new groups
are created, and unchanged groups are skipped. A group file without a siteId
fails at create time. Dry-run by default — pass --yes to apply changes.`,
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
				Name:    "group",
				Command: "groups",
				Decode:  decodeGroup,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					// Pull filters by --site-id; push registers no scope flag,
					// so it lists every group (matching spans the tenant).
					params := &mgmt.GroupListParams{SiteIDs: scope.SiteIDs, Limit: 1000}
					groups, _, lErr := fetchAllREST("group", func(cur string) ([]mgmt.Group, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.GroupsList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(groups))
					for _, g := range groups {
						f := groupToFile(g)
						body, mErr := yaml.Marshal(f)
						if mErr != nil {
							return nil, fmt.Errorf("marshal group %s: %w", g.Name, mErr)
						}
						objs = append(objs, reconcile.Object{Name: groupIdentity(f), ID: g.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f groupFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					if f.SiteID == "" {
						return fmt.Errorf("group %q has no siteId", f.Name)
					}
					_, cErr := c.GroupsCreate(ctx, f.SiteID, f.toCreate())
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f groupFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, uErr := c.GroupsUpdate(ctx, id, f.toUpdate())
					return uErr
				},
			}, nil
		},
	}
}
