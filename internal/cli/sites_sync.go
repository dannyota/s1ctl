package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addSiteSyncCmds(parent *cobra.Command) {
	spec := sitesSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// siteFile is the YAML representation of a site on disk. It holds only the
// declarative fields; server-assigned IDs, state, and timestamps are omitted.
type siteFile struct {
	Name              string `yaml:"name"`
	AccountID         string `yaml:"accountId"`
	SiteType          string `yaml:"siteType,omitempty"`
	Description       string `yaml:"description,omitempty"`
	Expiration        string `yaml:"expiration,omitempty"`
	UnlimitedLicenses bool   `yaml:"unlimitedLicenses"`
	TotalLicenses     int    `yaml:"totalLicenses"`
}

func siteToFile(s mgmt.Site) siteFile {
	return siteFile{
		Name:              s.Name,
		AccountID:         s.AccountID,
		SiteType:          s.SiteType,
		Description:       s.Description,
		Expiration:        s.Expiration,
		UnlimitedLicenses: s.UnlimitedLicenses,
		TotalLicenses:     s.TotalLicenses,
	}
}

func (f siteFile) toCreate() mgmt.SiteCreate {
	return mgmt.SiteCreate{
		Name:              f.Name,
		AccountID:         f.AccountID,
		SiteType:          f.SiteType,
		Description:       f.Description,
		Expiration:        f.Expiration,
		UnlimitedLicenses: f.UnlimitedLicenses,
		TotalLicenses:     f.TotalLicenses,
	}
}

// toUpdate builds the pointer-field update body, setting every declaratively
// updatable field (a full replacement of the declared fields). AccountID and
// SiteType are set only at creation time and are not part of SiteUpdate.
func (f siteFile) toUpdate() mgmt.SiteUpdate {
	return mgmt.SiteUpdate{
		Name:              &f.Name,
		Description:       &f.Description,
		Expiration:        &f.Expiration,
		UnlimitedLicenses: &f.UnlimitedLicenses,
		TotalLicenses:     &f.TotalLicenses,
	}
}

// decodeSite maps one local file to a canonical Object: it validates the site
// name and re-marshals through siteFile so bodies are byte-equal to the ones
// List produces.
func decodeSite(data []byte) (reconcile.Object, error) {
	var f siteFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Name == "" {
		return reconcile.Object{}, fmt.Errorf("site has no name")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: f.Name, Body: body}, nil
}

// sitesSpec adapts sites to the shared sync builders. Identity is the site name.
func sitesSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "site",
		Command:    "sites",
		DefaultDir: "sites",
		PullShort:  "Pull sites to local YAML files",
		PullLong: `Fetch all sites and write them as YAML files.

Each site produces one file named by its sanitized name. Server-only metadata
(ID, state, licenses in use, timestamps) is omitted so the files contain only
the declarative site definition.`,
		PushShort: "Push sites from local YAML files",
		PushLong: `Read site YAML files from a directory and sync them to SentinelOne.

Sites are matched by name: existing sites are updated, new sites are created,
and unchanged sites are skipped. Dry-run by default — pass --yes to apply changes.`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
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
				Name:    "site",
				Command: "sites",
				Decode:  decodeSite,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					// Pull filters by --account-id; push registers no scope
					// flag, so it lists every site (matching spans the tenant).
					params := &mgmt.SiteListParams{AccountIDs: scope.AccountIDs, Limit: 1000}
					sites, _, lErr := fetchAllREST("site", func(cur string) ([]mgmt.Site, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.SitesList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(sites))
					for _, s := range sites {
						body, mErr := yaml.Marshal(siteToFile(s))
						if mErr != nil {
							return nil, fmt.Errorf("marshal site %s: %w", s.Name, mErr)
						}
						objs = append(objs, reconcile.Object{Name: s.Name, ID: s.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f siteFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, cErr := c.SitesCreate(ctx, f.toCreate())
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f siteFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					_, uErr := c.SitesUpdate(ctx, id, f.toUpdate())
					return uErr
				},
			}, nil
		},
	}
}
