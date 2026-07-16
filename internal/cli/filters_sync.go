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

func addFilterSyncCmds(parent *cobra.Command) {
	spec := filtersSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// filterSyncFile is the YAML representation of a saved filter on disk. It holds
// only the declarative fields; server-assigned ID, scope, and timestamps are
// omitted. filterFields is kept as a generic value so it round-trips through
// YAML natively (the SDK models it as raw JSON).
type filterSyncFile struct {
	Name         string `yaml:"name"`
	FilterFields any    `yaml:"filterFields,omitempty"`
}

func filterToSyncFile(f mgmt.Filter) (filterSyncFile, error) {
	sf := filterSyncFile{Name: f.Name}
	if len(f.FilterFields) > 0 {
		var v any
		if err := json.Unmarshal(f.FilterFields, &v); err != nil {
			return sf, fmt.Errorf("filter %s filterFields: %w", f.Name, err)
		}
		sf.FilterFields = v
	}
	return sf, nil
}

func (f filterSyncFile) toData() (mgmt.FilterData, error) {
	d := mgmt.FilterData{Name: f.Name}
	if f.FilterFields != nil {
		b, err := json.Marshal(f.FilterFields)
		if err != nil {
			return d, fmt.Errorf("filter %s filterFields: %w", f.Name, err)
		}
		d.FilterFields = b
	}
	return d, nil
}

// decodeFilter maps one local file to a canonical Object. Identity is the
// filter name; the body is re-marshalled so it is byte-equal to what List
// produces.
func decodeFilter(data []byte) (reconcile.Object, error) {
	var f filterSyncFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Name == "" {
		return reconcile.Object{}, fmt.Errorf("filter has no name")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: f.Name, Body: body}, nil
}

// filtersSpec adapts saved filters to the shared sync builders. Identity is the
// filter name.
func filtersSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "filter",
		Command:    "filters",
		DefaultDir: "filters",
		PullShort:  "Pull saved filters to local YAML files",
		PullLong: `Fetch all saved filters and write them as YAML files.

Each filter produces one file named by its sanitized name. Server-only metadata
(ID, scope, timestamps) is omitted so the files contain only the declarative
definition: the filter name and its filterFields criteria set.`,
		PushShort: "Push saved filters from local YAML files",
		PushLong: `Read filter YAML files from a directory and sync them to SentinelOne.

Filters are matched by name: existing filters are updated, new ones are created,
and unchanged ones are skipped. Dry-run by default — pass --yes to apply changes.
New filters are created at the scope given by --site-id (default: global/tenant).`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
			cmd.Flags().StringSliceVar(&scope.AccountIDs, "account-id", nil, "filter by account ID")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "scope for new filters (default: global/tenant)")
			cmd.Flags().StringSliceVar(&scope.AccountIDs, "account-id", nil, "scope for new filters")
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
				Name:    "filter",
				Command: "filters",
				Decode:  decodeFilter,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					params := &mgmt.FilterListParams{Limit: 1000}
					if !scope.push {
						params.SiteIDs = scope.SiteIDs
						params.AccountIDs = scope.AccountIDs
					}
					filters, _, lErr := fetchAllREST("filter", func(cur string) ([]mgmt.Filter, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.FiltersList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(filters))
					for _, f := range filters {
						sf, fErr := filterToSyncFile(f)
						if fErr != nil {
							return nil, fErr
						}
						body, mErr := yaml.Marshal(sf)
						if mErr != nil {
							return nil, fmt.Errorf("marshal filter %s: %w", f.Name, mErr)
						}
						objs = append(objs, reconcile.Object{Name: f.Name, ID: f.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f filterSyncFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					data, dErr := f.toData()
					if dErr != nil {
						return dErr
					}
					body := mgmt.FilterCreate{Data: data}
					if len(scope.SiteIDs) > 0 || len(scope.AccountIDs) > 0 {
						body.Filter = &mgmt.FilterScope{
							SiteIDs:    scope.SiteIDs,
							AccountIDs: scope.AccountIDs,
						}
					}
					_, cErr := c.FiltersCreate(ctx, body)
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f filterSyncFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					data, dErr := f.toData()
					if dErr != nil {
						return dErr
					}
					_, uErr := c.FiltersUpdate(ctx, id, mgmt.FilterUpdate{Data: data})
					return uErr
				},
			}, nil
		},
	}
}
