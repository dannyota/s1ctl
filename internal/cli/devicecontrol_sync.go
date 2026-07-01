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

func addDeviceControlSyncCmds(parent *cobra.Command) {
	parent.AddCommand(newDeviceControlPullCmd())
	parent.AddCommand(newDeviceControlPushCmd())
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

func newDeviceControlPullCmd() *cobra.Command {
	var siteIDs, accountIDs []string
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull device control rules to local YAML files",
		Long: `Fetch all device control rules and write them as YAML files.

Each rule produces one file named by its sanitized rule name (e.g. block-usb-storage.yaml).
Server-only metadata (ID, scope, timestamps) is omitted from the YAML so the files
contain only the declarative rule definition.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.DeviceRuleListParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				Limit:      1000,
			}
			rules, _, err := fetchAllREST("device rule", func(cur string) ([]mgmt.DeviceRule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.DeviceRulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}

			used := make(map[string]int)
			for _, r := range rules {
				stem := sanitizeFilename(r.RuleName)
				if n := used[stem]; n > 0 {
					stem = fmt.Sprintf("%s-%d", stem, n)
				}
				used[sanitizeFilename(r.RuleName)]++

				rf := deviceRuleToFile(r)
				data, mErr := yaml.Marshal(rf)
				if mErr != nil {
					return fmt.Errorf("marshal rule %s: %w", r.RuleName, mErr)
				}
				path := filepath.Join(outDir, stem+".yaml")
				if wErr := os.WriteFile(path, data, 0o644); wErr != nil {
					return wErr
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n",
				pluralize(len(rules), "device rule"), outDir)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringVar(&outDir, "out", "devicecontrol", "output directory")
	return cmd
}

func newDeviceControlPushCmd() *cobra.Command {
	var inDir string
	var siteIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push device control rules from local YAML files",
		Long: `Read device control rule YAML files from a directory and sync them to SentinelOne.

Rules are matched by name: existing rules are updated, new rules are created.
Dry-run by default — pass --yes to apply changes.

New rules are created at the scope specified by --site-id. If no scope flag
is given, new rules are created at the global (tenant) scope.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, err := os.ReadDir(inDir)
			if err != nil {
				return fmt.Errorf("read %s: %w", inDir, err)
			}

			var localRules []deviceRuleFile
			for _, e := range entries {
				if e.IsDir() || (!strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml")) {
					continue
				}
				data, rErr := os.ReadFile(filepath.Join(inDir, e.Name()))
				if rErr != nil {
					return fmt.Errorf("read %s: %w", e.Name(), rErr)
				}
				var rf deviceRuleFile
				if uErr := yaml.Unmarshal(data, &rf); uErr != nil {
					return fmt.Errorf("parse %s: %w", e.Name(), uErr)
				}
				if rf.RuleName == "" {
					return fmt.Errorf("rule in %s has no ruleName", e.Name())
				}
				localRules = append(localRules, rf)
			}
			if len(localRules) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No device rule files found.")
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			// Fetch all remote rules to match by name.
			params := &mgmt.DeviceRuleListParams{Limit: 1000}
			remoteRules, _, err := fetchAllREST("device rule", func(cur string) ([]mgmt.DeviceRule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.DeviceRulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}
			byName := make(map[string]mgmt.DeviceRule, len(remoteRules))
			for _, r := range remoteRules {
				byName[r.RuleName] = r
			}

			var toCreate, toUpdate []deviceRuleFile
			var updateIDs []string
			for _, lr := range localRules {
				if existing, ok := byName[lr.RuleName]; ok {
					toUpdate = append(toUpdate, lr)
					updateIDs = append(updateIDs, existing.ID)
				} else {
					toCreate = append(toCreate, lr)
				}
			}

			action := fmt.Sprintf("create %s, update %s from %s",
				pluralize(len(toCreate), "device rule"),
				pluralize(len(toUpdate), "device rule"),
				inDir)
			return guard(cmd.OutOrStdout(), "devicecontrol push", action, inDir, yes, func() error {
				// Build scope filter for new rules.
				filter := mgmt.DeviceRuleScopeFilter{
					SiteIDs: siteIDs,
				}
				if len(siteIDs) == 0 {
					t := true
					filter.Tenant = &t
				}
				var created, updated int
				for _, lr := range toCreate {
					if _, cErr := c.DeviceRulesCreate(cmd.Context(), lr.toCreate(), filter); cErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: create %s: %v\n", lr.RuleName, cErr)
						continue
					}
					created++
				}
				for i, lr := range toUpdate {
					if _, uErr := c.DeviceRulesUpdate(cmd.Context(), updateIDs[i], lr.toCreate()); uErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "warning: update %s: %v\n", lr.RuleName, uErr)
						continue
					}
					updated++
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s, updated %s\n",
					pluralize(created, "device rule"),
					pluralize(updated, "device rule"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&inDir, "dir", "devicecontrol", "directory containing device rule YAML files")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope for new rules (default: global/tenant)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}
