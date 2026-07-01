package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/mgmt"
)

func addFirewallSyncCmds(parent *cobra.Command) {
	parent.AddCommand(newFirewallPullCmd())
	parent.AddCommand(newFirewallPushCmd())
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

func newFirewallPullCmd() *cobra.Command {
	var siteIDs []string
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull firewall rules to local YAML files",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.FirewallRuleListParams{
				SiteIDs: siteIDs,
				Limit:   1000,
			}
			rules, _, err := fetchAllREST("firewall rule", func(cur string) ([]mgmt.FirewallRule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.FirewallRulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}

			// Track used filenames to avoid collisions.
			used := make(map[string]int)
			for _, r := range rules {
				stem := sanitizeFilename(r.Name)
				if n := used[stem]; n > 0 {
					stem = fmt.Sprintf("%s-%d", stem, n)
				}
				used[sanitizeFilename(r.Name)]++

				ff := firewallToFile(r)
				data, mErr := yaml.Marshal(ff)
				if mErr != nil {
					return fmt.Errorf("marshal rule %s: %w", r.Name, mErr)
				}
				path := filepath.Join(outDir, stem+".yaml")
				if wErr := os.WriteFile(path, data, 0o644); wErr != nil {
					return wErr
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(len(rules), "firewall rule"), outDir)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&outDir, "out", "firewall", "output directory")
	return cmd
}

func newFirewallPushCmd() *cobra.Command {
	var inDir string
	var siteIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push firewall rules from local YAML files",
		Long: `Read firewall rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created.
Dry-run by default — pass --yes to apply changes.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, err := os.ReadDir(inDir)
			if err != nil {
				return fmt.Errorf("read %s: %w", inDir, err)
			}

			var localRules []firewallFile
			for _, e := range entries {
				if e.IsDir() || (!strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml")) {
					continue
				}
				data, rErr := os.ReadFile(filepath.Join(inDir, e.Name()))
				if rErr != nil {
					return fmt.Errorf("read %s: %w", e.Name(), rErr)
				}
				var ff firewallFile
				if uErr := yaml.Unmarshal(data, &ff); uErr != nil {
					return fmt.Errorf("parse %s: %w", e.Name(), uErr)
				}
				if ff.Name == "" {
					return fmt.Errorf("rule in %s has no name", e.Name())
				}
				localRules = append(localRules, ff)
			}
			if len(localRules) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No firewall rule files found.")
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			// Fetch all remote rules to match by name.
			params := &mgmt.FirewallRuleListParams{SiteIDs: siteIDs, Limit: 1000}
			remoteRules, _, err := fetchAllREST("firewall rule", func(cur string) ([]mgmt.FirewallRule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.FirewallRulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}
			byName := make(map[string]mgmt.FirewallRule, len(remoteRules))
			for _, r := range remoteRules {
				byName[r.Name] = r
			}

			var toCreate, toUpdate []firewallFile
			var updateIDs []string
			for _, lr := range localRules {
				if existing, ok := byName[lr.Name]; ok {
					toUpdate = append(toUpdate, lr)
					updateIDs = append(updateIDs, existing.ID)
				} else {
					toCreate = append(toCreate, lr)
				}
			}

			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would create %s, update %s from %s. Pass --yes to apply.\n",
					pluralize(len(toCreate), "firewall rule"),
					pluralize(len(toUpdate), "firewall rule"),
					inDir)
				return nil
			}

			scope := mgmt.FirewallRuleScope{SiteIDs: siteIDs}
			var created, updated int
			for _, lr := range toCreate {
				if _, cErr := c.FirewallRulesCreate(cmd.Context(), scope, lr.toCreate()); cErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: create %s: %v\n", lr.Name, cErr)
					continue
				}
				created++
			}
			for i, lr := range toUpdate {
				if _, uErr := c.FirewallRulesUpdate(cmd.Context(), updateIDs[i], lr.toCreate()); uErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: update %s: %v\n", lr.Name, uErr)
					continue
				}
				updated++
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created %s, updated %s\n",
				pluralize(created, "firewall rule"),
				pluralize(updated, "firewall rule"))
			return nil
		},
	}
	cmd.Flags().StringVar(&inDir, "dir", "firewall", "directory containing firewall rule YAML files")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "target site IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}
