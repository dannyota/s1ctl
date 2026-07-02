package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addNetworkSyncCmds(parent *cobra.Command) {
	spec := networkQuarantineSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// networkQuarantineSpec adapts network quarantine rules to the shared sync
// builders. Network quarantine rules share the firewall rule body shape, so it
// reuses the firewall file helpers (firewallFile, firewallToFile, decodeFirewall)
// and differs only in the API methods and default directory.
func networkQuarantineSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "network quarantine rule",
		Command:    "network",
		DefaultDir: "network-quarantine",
		PullShort:  "Pull network quarantine rules to local YAML files",
		PushShort:  "Push network quarantine rules from local YAML files",
		PushLong: `Read network quarantine rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created,
and unchanged rules are skipped. Dry-run by default — pass --yes to apply changes.`,
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
				Name:    "network quarantine rule",
				Command: "network",
				Decode:  decodeFirewall,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					params := &mgmt.FirewallRuleListParams{SiteIDs: scope.SiteIDs, Limit: 1000}
					rules, _, lErr := fetchAllREST("network quarantine rule", func(cur string) ([]mgmt.FirewallRule, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.NetworkQuarantineList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(rules))
					for _, r := range rules {
						body, mErr := yaml.Marshal(firewallToFile(r))
						if mErr != nil {
							return nil, fmt.Errorf("marshal network quarantine rule %s: %w", r.Name, mErr)
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
					_, cErr := c.NetworkQuarantineCreate(ctx, mgmt.FirewallRuleScope{SiteIDs: scope.SiteIDs}, ff.toCreate())
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
					_, uErr := c.NetworkQuarantineUpdate(ctx, id, ff.toCreate())
					return uErr
				},
			}, nil
		},
	}
}
