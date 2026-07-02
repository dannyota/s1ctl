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

func addLocationSyncCmds(parent *cobra.Command) {
	spec := locationsSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// locationFile is the YAML representation of a location on disk. It holds the
// declarative definition only; server-assigned IDs, scope, counters, and
// timestamps are omitted. The six detection-parameter groups are kept as
// generic values so they round-trip through YAML natively (the SDK models them
// as raw JSON).
type locationFile struct {
	Name               string `yaml:"name"`
	Description        string `yaml:"description,omitempty"`
	Operator           string `yaml:"operator"`
	DNSLookup          any    `yaml:"dnsLookup,omitempty"`
	DNSServers         any    `yaml:"dnsServers,omitempty"`
	RegistryKeys       any    `yaml:"registryKeys,omitempty"`
	ServerConnectivity any    `yaml:"serverConnectivity,omitempty"`
	NetworkInterfaces  any    `yaml:"networkInterfaces,omitempty"`
	IPAddresses        any    `yaml:"ipAddresses,omitempty"`
}

// locationToFile projects a live location onto its declarative file shape,
// decoding each raw detection-parameter blob into a generic value.
func locationToFile(l mgmt.Location) (locationFile, error) {
	f := locationFile{Name: l.Name, Description: l.Description, Operator: string(l.Operator)}
	params := []struct {
		raw json.RawMessage
		dst *any
	}{
		{l.DNSLookup, &f.DNSLookup},
		{l.DNSServers, &f.DNSServers},
		{l.RegistryKeys, &f.RegistryKeys},
		{l.ServerConnectivity, &f.ServerConnectivity},
		{l.NetworkInterfaces, &f.NetworkInterfaces},
		{l.IPAddresses, &f.IPAddresses},
	}
	for _, p := range params {
		if len(p.raw) == 0 {
			continue
		}
		var v any
		if err := yaml.Unmarshal(p.raw, &v); err != nil {
			return f, fmt.Errorf("location %s: %w", l.Name, err)
		}
		*p.dst = v
	}
	return f, nil
}

// toData converts a file back into the SDK write payload, re-encoding each
// detection-parameter group as raw JSON.
func (f locationFile) toData() (mgmt.LocationData, error) {
	d := mgmt.LocationData{
		Name:        f.Name,
		Description: f.Description,
		Operator:    mgmt.LocationOperator(f.Operator),
	}
	params := []struct {
		src any
		dst *json.RawMessage
	}{
		{f.DNSLookup, &d.DNSLookup},
		{f.DNSServers, &d.DNSServers},
		{f.RegistryKeys, &d.RegistryKeys},
		{f.ServerConnectivity, &d.ServerConnectivity},
		{f.NetworkInterfaces, &d.NetworkInterfaces},
		{f.IPAddresses, &d.IPAddresses},
	}
	for _, p := range params {
		if p.src == nil {
			continue
		}
		b, err := json.Marshal(p.src)
		if err != nil {
			return d, fmt.Errorf("location %s: %w", f.Name, err)
		}
		*p.dst = b
	}
	return d, nil
}

// decodeLocation maps one local file to a canonical Object. Identity is the
// location name; the body is re-marshalled so it is byte-equal to what List
// produces.
func decodeLocation(data []byte) (reconcile.Object, error) {
	var f locationFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Name == "" {
		return reconcile.Object{}, fmt.Errorf("location has no name")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: f.Name, Body: body}, nil
}

// locationsSpec adapts locations to the shared sync builders. Identity is the
// location name.
func locationsSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "location",
		Command:    "locations",
		DefaultDir: "locations",
		PullShort:  "Pull locations to local YAML files",
		PullLong: `Fetch all locations and write them as YAML files.

Each location produces one file named by its sanitized name. Server-only
metadata (ID, scope, counters, timestamps) is omitted so the files contain only
the declarative definition, including the detection parameters.`,
		PushShort: "Push locations from local YAML files",
		PushLong: `Read location YAML files from a directory and sync them to SentinelOne.

Locations are matched by name: existing locations are updated, new ones are
created, and unchanged ones are skipped. Dry-run by default — pass --yes to
apply. New locations are created at the scope given by --site-id (default:
global/tenant).`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "scope for new locations (default: global/tenant)")
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
				Name:    "location",
				Command: "locations",
				Decode:  decodeLocation,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					params := &mgmt.LocationListParams{Limit: 1000}
					if !scope.push {
						params.SiteIDs = scope.SiteIDs
					}
					locs, _, lErr := fetchAllREST("location", func(cur string) ([]mgmt.Location, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.LocationsList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(locs))
					for _, l := range locs {
						f, fErr := locationToFile(l)
						if fErr != nil {
							return nil, fErr
						}
						body, mErr := yaml.Marshal(f)
						if mErr != nil {
							return nil, fmt.Errorf("marshal location %s: %w", l.Name, mErr)
						}
						objs = append(objs, reconcile.Object{Name: l.Name, ID: l.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f locationFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					data, dErr := f.toData()
					if dErr != nil {
						return dErr
					}
					_, cErr := c.LocationsCreate(ctx, mgmt.LocationCreate{
						Data:   data,
						Filter: mgmt.LocationScope{SiteIDs: scope.SiteIDs},
					})
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f locationFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					data, dErr := f.toData()
					if dErr != nil {
						return dErr
					}
					_, uErr := c.LocationsUpdate(ctx, id, mgmt.LocationUpdate{Data: data})
					return uErr
				},
			}, nil
		},
	}
}
