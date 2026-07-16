package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addFirewallSyncCmds(parent *cobra.Command) {
	spec := firewallSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// firewallHostFile is the YAML representation of a host matcher.
type firewallHostFile struct {
	Type   string   `yaml:"type"`
	Values []string `yaml:"values,omitempty"`
}

// firewallPortFile is the YAML representation of a port matcher.
type firewallPortFile struct {
	Type   string   `yaml:"type"`
	Values []string `yaml:"values,omitempty"`
}

// firewallLocationFile is the YAML representation of a location matcher.
type firewallLocationFile struct {
	Type   string   `yaml:"type"`
	Values []string `yaml:"values,omitempty"`
}

// firewallAppFile is the YAML representation of an application matcher.
type firewallAppFile struct {
	Type   string   `yaml:"type"`
	Values []string `yaml:"values,omitempty"`
}

// firewallFile is the YAML representation of a firewall rule on disk.
type firewallFile struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description,omitempty"`
	Direction   mgmt.FirewallDirection `yaml:"direction"`
	Protocol    string                 `yaml:"protocol,omitempty"`
	OSTypes     []string               `yaml:"osTypes,omitempty"`
	Action      mgmt.FirewallAction    `yaml:"action"`
	Status      mgmt.FirewallStatus    `yaml:"status"`
	Application *firewallAppFile       `yaml:"application,omitempty"`
	LocalHost   *firewallHostFile      `yaml:"localHost,omitempty"`
	LocalPort   *firewallPortFile      `yaml:"localPort,omitempty"`
	RemoteHosts []firewallHostFile     `yaml:"remoteHosts,omitempty"`
	RemotePort  *firewallPortFile      `yaml:"remotePort,omitempty"`
	Location    *firewallLocationFile  `yaml:"location,omitempty"`
	TagIDs      []string               `yaml:"tagIds,omitempty"`
}

func firewallToFile(r mgmt.FirewallRule) firewallFile {
	f := firewallFile{
		Name:        r.Name,
		Description: r.Description,
		Direction:   r.Direction,
		Protocol:    r.Protocol,
		OSTypes:     r.OSTypes,
		Action:      r.Action,
		Status:      r.Status,
		TagIDs:      r.TagIDs,
	}
	if r.Application != nil {
		f.Application = &firewallAppFile{Type: string(r.Application.Type), Values: r.Application.Values}
	}
	if r.LocalHost != nil {
		f.LocalHost = &firewallHostFile{Type: string(r.LocalHost.Type), Values: r.LocalHost.Values}
	}
	if r.LocalPort != nil {
		f.LocalPort = &firewallPortFile{Type: string(r.LocalPort.Type), Values: r.LocalPort.Values}
	}
	for _, rh := range r.RemoteHosts {
		f.RemoteHosts = append(f.RemoteHosts, firewallHostFile{Type: string(rh.Type), Values: rh.Values})
	}
	if r.RemotePort != nil {
		f.RemotePort = &firewallPortFile{Type: string(r.RemotePort.Type), Values: r.RemotePort.Values}
	}
	if r.Location != nil {
		f.Location = &firewallLocationFile{Type: string(r.Location.Type), Values: r.Location.Values}
	}
	return f
}

func (ff firewallFile) toCreate() mgmt.FirewallRuleCreate {
	c := mgmt.FirewallRuleCreate{
		Name:        ff.Name,
		Description: ff.Description,
		Direction:   ff.Direction,
		Protocol:    ff.Protocol,
		OSTypes:     ff.OSTypes,
		Action:      ff.Action,
		Status:      ff.Status,
		TagIDs:      ff.TagIDs,
	}
	if ff.Application != nil {
		c.Application = &mgmt.FirewallApplication{Type: mgmt.FirewallAppType(ff.Application.Type), Values: ff.Application.Values}
	}
	if ff.LocalHost != nil {
		c.LocalHost = &mgmt.FirewallHost{Type: mgmt.FirewallHostType(ff.LocalHost.Type), Values: ff.LocalHost.Values}
	}
	if ff.LocalPort != nil {
		c.LocalPort = &mgmt.FirewallPort{Type: mgmt.FirewallPortType(ff.LocalPort.Type), Values: ff.LocalPort.Values}
	}
	for _, rh := range ff.RemoteHosts {
		c.RemoteHosts = append(c.RemoteHosts, mgmt.FirewallHost{Type: mgmt.FirewallHostType(rh.Type), Values: rh.Values})
	}
	if ff.RemotePort != nil {
		c.RemotePort = &mgmt.FirewallPort{Type: mgmt.FirewallPortType(ff.RemotePort.Type), Values: ff.RemotePort.Values}
	}
	if ff.Location != nil {
		c.Location = &mgmt.FirewallLocation{Type: mgmt.FirewallLocationType(ff.Location.Type), Values: ff.Location.Values}
	}
	return c
}

// decodeFirewall maps one local file to a canonical Object: it validates the
// rule name and re-marshals through firewallFile so bodies are byte-equal to
// the ones List produces.
func decodeFirewall(data []byte) (reconcile.Object, error) {
	var ff firewallFile
	if err := yaml.Unmarshal(data, &ff); err != nil {
		return reconcile.Object{}, err
	}
	if ff.Name == "" {
		return reconcile.Object{}, fmt.Errorf("firewall rule has no name")
	}
	body, err := yaml.Marshal(ff)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: ff.Name, Body: body}, nil
}

// firewallSpec adapts firewall rules to the shared sync builders.
func firewallSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "firewall rule",
		Command:    "firewall",
		DefaultDir: "firewall",
		PullShort:  "Pull firewall rules to local YAML files",
		PushShort:  "Push firewall rules from local YAML files",
		PushLong: `Read firewall rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created,
and unchanged rules are skipped. Dry-run by default — pass --yes to apply changes.

Note: the plan is built against the list scoped by --site-id. Without --site-id
the list is unscoped, which may match rules from other sites and produce an
incorrect plan. Always pass --site-id when pushing to a specific site.`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "target site IDs")
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
				Name:    "firewall rule",
				Command: "firewall",
				Decode:  decodeFirewall,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					// Both pull and push scope the list by --site-id (legacy).
					params := &mgmt.FirewallRuleListParams{SiteIDs: scope.SiteIDs, Limit: 1000}
					rules, _, lErr := fetchAllREST("firewall rule", func(cur string) ([]mgmt.FirewallRule, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.FirewallRulesList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(rules))
					for _, r := range rules {
						body, mErr := yaml.Marshal(firewallToFile(r))
						if mErr != nil {
							return nil, fmt.Errorf("marshal firewall rule %s: %w", r.Name, mErr)
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
					var ff firewallFile
					if uErr := yaml.Unmarshal(local.Body, &ff); uErr != nil {
						return uErr
					}
					_, cErr := c.FirewallRulesCreate(ctx, mgmt.FirewallRuleScope{SiteIDs: scope.SiteIDs}, ff.toCreate())
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var ff firewallFile
					if uErr := yaml.Unmarshal(local.Body, &ff); uErr != nil {
						return uErr
					}
					_, uErr := c.FirewallRulesUpdate(ctx, id, ff.toCreate())
					return uErr
				},
			}, nil
		},
	}
}
