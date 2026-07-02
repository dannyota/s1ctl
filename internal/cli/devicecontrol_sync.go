package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addDeviceControlSyncCmds(parent *cobra.Command) {
	spec := deviceControlSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// deviceRuleFile is the YAML representation of a device control rule on disk.
type deviceRuleFile struct {
	RuleName         string   `yaml:"ruleName"`
	Interface        string   `yaml:"interface"`
	RuleType         string   `yaml:"ruleType"`
	Action           string   `yaml:"action"`
	Status           string   `yaml:"status"`
	AccessPermission string   `yaml:"accessPermission"`
	DeviceClass      string   `yaml:"deviceClass,omitempty"`
	DeviceID         string   `yaml:"deviceId,omitempty"`
	VendorID         string   `yaml:"vendorId,omitempty"`
	ProductID        string   `yaml:"productId,omitempty"`
	UID              string   `yaml:"uid,omitempty"`
	Version          string   `yaml:"version,omitempty"`
	MinorClasses     []string `yaml:"minorClasses,omitempty"`
	BluetoothAddress string   `yaml:"bluetoothAddress,omitempty"`
	GattService      []string `yaml:"gattService,omitempty"`
	ManufacturerName string   `yaml:"manufacturerName,omitempty"`
	DeviceName       string   `yaml:"deviceName,omitempty"`
}

func deviceRuleToFile(r mgmt.DeviceRule) deviceRuleFile {
	return deviceRuleFile{
		RuleName:         r.RuleName,
		Interface:        string(r.Interface),
		RuleType:         string(r.RuleType),
		Action:           string(r.Action),
		Status:           string(r.Status),
		AccessPermission: string(r.AccessPermission),
		DeviceClass:      r.DeviceClass,
		DeviceID:         r.DeviceID,
		VendorID:         r.VendorID,
		ProductID:        r.ProductID,
		UID:              r.UID,
		Version:          r.Version,
		MinorClasses:     r.MinorClasses,
		BluetoothAddress: r.BluetoothAddress,
		GattService:      r.GattService,
		ManufacturerName: r.ManufacturerName,
		DeviceName:       r.DeviceName,
	}
}

func (rf deviceRuleFile) toCreate() mgmt.DeviceRuleCreate {
	return mgmt.DeviceRuleCreate{
		RuleName:         rf.RuleName,
		Interface:        mgmt.DeviceRuleInterface(rf.Interface),
		RuleType:         mgmt.DeviceRuleType(rf.RuleType),
		Action:           mgmt.DeviceRuleAction(rf.Action),
		Status:           mgmt.DeviceRuleStatus(rf.Status),
		AccessPermission: mgmt.DeviceRuleAccessPermission(rf.AccessPermission),
		DeviceClass:      rf.DeviceClass,
		DeviceID:         rf.DeviceID,
		VendorID:         rf.VendorID,
		ProductID:        rf.ProductID,
		UID:              rf.UID,
		Version:          rf.Version,
		MinorClasses:     rf.MinorClasses,
		BluetoothAddress: rf.BluetoothAddress,
		GattService:      rf.GattService,
		ManufacturerName: rf.ManufacturerName,
		DeviceName:       rf.DeviceName,
	}
}

// decodeDeviceRule maps one local file to a canonical Object: it validates the
// rule name and re-marshals through deviceRuleFile so bodies are byte-equal to
// the ones List produces.
func decodeDeviceRule(data []byte) (reconcile.Object, error) {
	var rf deviceRuleFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return reconcile.Object{}, err
	}
	if rf.RuleName == "" {
		return reconcile.Object{}, fmt.Errorf("device rule has no ruleName")
	}
	body, err := yaml.Marshal(rf)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: rf.RuleName, Body: body}, nil
}

// deviceControlSpec adapts device control rules to the shared sync builders.
func deviceControlSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "device rule",
		Command:    "devicecontrol",
		DefaultDir: "devicecontrol",
		PullShort:  "Pull device control rules to local YAML files",
		PullLong: `Fetch all device control rules and write them as YAML files.

Each rule produces one file named by its sanitized rule name (e.g. block-usb-storage.yaml).
Server-only metadata (ID, scope, timestamps) is omitted from the YAML so the files
contain only the declarative rule definition.`,
		PushShort: "Push device control rules from local YAML files",
		PushLong: `Read device control rule YAML files from a directory and sync them to SentinelOne.

Rules are matched by name: existing rules are updated, new rules are created,
and unchanged rules are skipped. Dry-run by default — pass --yes to apply changes.

New rules are created at the scope specified by --site-id. If no scope flag
is given, new rules are created at the global (tenant) scope.`,
		RegisterPullFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "filter by site ID")
			cmd.Flags().StringSliceVar(&scope.AccountIDs, "account-id", nil, "filter by account ID")
		},
		RegisterPushFlags: func(cmd *cobra.Command, scope *scopeFlags) {
			cmd.Flags().StringSliceVar(&scope.SiteIDs, "site-id", nil, "scope for new rules (default: global/tenant)")
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
				Name:    "device rule",
				Command: "devicecontrol",
				Decode:  decodeDeviceRule,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					// Pull filters by scope; push lists every rule so matching
					// spans the whole tenant (--site-id is the create scope).
					params := &mgmt.DeviceRuleListParams{Limit: 1000}
					if !scope.push {
						params.SiteIDs = scope.SiteIDs
						params.AccountIDs = scope.AccountIDs
					}
					rules, _, lErr := fetchAllREST("device rule", func(cur string) ([]mgmt.DeviceRule, *mgmt.Pagination, error) {
						params.Cursor = cur
						return c.DeviceRulesList(ctx, params)
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(rules))
					for _, r := range rules {
						body, mErr := yaml.Marshal(deviceRuleToFile(r))
						if mErr != nil {
							return nil, fmt.Errorf("marshal device rule %s: %w", r.RuleName, mErr)
						}
						objs = append(objs, reconcile.Object{Name: r.RuleName, ID: r.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var rf deviceRuleFile
					if uErr := yaml.Unmarshal(local.Body, &rf); uErr != nil {
						return uErr
					}
					filter := mgmt.DeviceRuleScopeFilter{SiteIDs: scope.SiteIDs}
					if len(scope.SiteIDs) == 0 {
						tenant := true
						filter.Tenant = &tenant
					}
					_, cErr := c.DeviceRulesCreate(ctx, rf.toCreate(), filter)
					return cErr
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var rf deviceRuleFile
					if uErr := yaml.Unmarshal(local.Body, &rf); uErr != nil {
						return uErr
					}
					_, uErr := c.DeviceRulesUpdate(ctx, id, rf.toCreate())
					return uErr
				},
			}, nil
		},
	}
}
