package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/mgmt"
)

func addPolicySyncCmds(parent *cobra.Command) {
	parent.AddCommand(newPoliciesPullCmd())
	parent.AddCommand(newPoliciesPushCmd())
}

// policyFile is the YAML representation of a policy on disk.
type policyFile struct {
	Scope     string `yaml:"scope"`
	SiteID    string `yaml:"siteId,omitempty"`
	SiteName  string `yaml:"siteName,omitempty"`
	AccountID string `yaml:"accountId,omitempty"`
	GroupID   string `yaml:"groupId,omitempty"`
	GroupName string `yaml:"groupName,omitempty"`

	MitigationMode           string `yaml:"mitigationMode"`
	MitigationModeSuspicious string `yaml:"mitigationModeSuspicious"`
	AntiTamperingOn          bool   `yaml:"antiTamperingOn"`
	NetworkQuarantineOn      bool   `yaml:"networkQuarantineOn"`
	SnapshotsOn              bool   `yaml:"snapshotsOn"`
	AllowRemoteShell         bool   `yaml:"allowRemoteShell"`
	ScanNewAgents            bool   `yaml:"scanNewAgents"`
	AutoDecommissionOn       bool   `yaml:"autoDecommissionOn"`
	AutoDecommissionDays     int    `yaml:"autoDecommissionDays"`
	Ioc                      bool   `yaml:"ioc"`
}

func sitePolicyToFile(siteID, siteName string, p mgmt.Policy) policyFile {
	return policyFile{
		Scope:                    "site",
		SiteID:                   siteID,
		SiteName:                 siteName,
		MitigationMode:           p.MitigationMode,
		MitigationModeSuspicious: p.MitigationModeSuspicious,
		AntiTamperingOn:          p.AntiTamperingOn,
		NetworkQuarantineOn:      p.NetworkQuarantineOn,
		SnapshotsOn:              p.SnapshotsOn,
		AllowRemoteShell:         p.AllowRemoteShell,
		ScanNewAgents:            p.ScanNewAgents,
		AutoDecommissionOn:       p.AutoDecommissionOn,
		AutoDecommissionDays:     p.AutoDecommissionDays,
		Ioc:                      p.Ioc,
	}
}

func accountPolicyToFile(accountID string, p mgmt.Policy) policyFile {
	return policyFile{
		Scope:                    "account",
		AccountID:                accountID,
		MitigationMode:           p.MitigationMode,
		MitigationModeSuspicious: p.MitigationModeSuspicious,
		AntiTamperingOn:          p.AntiTamperingOn,
		NetworkQuarantineOn:      p.NetworkQuarantineOn,
		SnapshotsOn:              p.SnapshotsOn,
		AllowRemoteShell:         p.AllowRemoteShell,
		ScanNewAgents:            p.ScanNewAgents,
		AutoDecommissionOn:       p.AutoDecommissionOn,
		AutoDecommissionDays:     p.AutoDecommissionDays,
		Ioc:                      p.Ioc,
	}
}

func groupPolicyToFile(siteID, groupID, groupName string, p mgmt.Policy) policyFile {
	return policyFile{
		Scope:                    "group",
		SiteID:                   siteID,
		GroupID:                  groupID,
		GroupName:                groupName,
		MitigationMode:           p.MitigationMode,
		MitigationModeSuspicious: p.MitigationModeSuspicious,
		AntiTamperingOn:          p.AntiTamperingOn,
		NetworkQuarantineOn:      p.NetworkQuarantineOn,
		SnapshotsOn:              p.SnapshotsOn,
		AllowRemoteShell:         p.AllowRemoteShell,
		ScanNewAgents:            p.ScanNewAgents,
		AutoDecommissionOn:       p.AutoDecommissionOn,
		AutoDecommissionDays:     p.AutoDecommissionDays,
		Ioc:                      p.Ioc,
	}
}

// policyUpdatePayload returns the JSON body for a policy PUT, containing
// only the mutable fields from the YAML file.
func policyUpdatePayload(pf policyFile) (json.RawMessage, error) {
	payload := map[string]any{
		"mitigationMode":           pf.MitigationMode,
		"mitigationModeSuspicious": pf.MitigationModeSuspicious,
		"antiTamperingOn":          pf.AntiTamperingOn,
		"networkQuarantineOn":      pf.NetworkQuarantineOn,
		"snapshotsOn":              pf.SnapshotsOn,
		"allowRemoteShell":         pf.AllowRemoteShell,
		"scanNewAgents":            pf.ScanNewAgents,
		"autoDecommissionOn":       pf.AutoDecommissionOn,
		"autoDecommissionDays":     pf.AutoDecommissionDays,
		"ioc":                      pf.Ioc,
	}
	return json.Marshal(payload)
}

func newPoliciesPullCmd() *cobra.Command {
	var accountIDs, siteIDs []string
	var outDir, scope string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull policies to local YAML files",
		Long: `Fetch endpoint policies and write them as YAML files.

By default pulls site-level policies. Use --scope to select account or group level.
Each policy produces one YAML file. The YAML includes scope metadata for push matching,
plus the key policy fields: mitigationMode, antiTamperingOn, networkQuarantineOn, etc.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			switch scope {
			case "site":
				return pullSitePolicies(cmd, accountIDs, siteIDs, outDir)
			case "account":
				return pullAccountPolicies(cmd, accountIDs, outDir)
			case "group":
				return pullGroupPolicies(cmd, siteIDs, outDir)
			default:
				return fmt.Errorf("invalid --scope %q: must be site, account, or group", scope)
			}
		},
	}
	cmd.Flags().StringVar(&scope, "scope", "site", "policy scope: site, account, or group")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&outDir, "out", "policies", "output directory")
	return cmd
}

func pullSitePolicies(cmd *cobra.Command, accountIDs, siteIDs []string, outDir string) error {
	c, err := mgmtClient()
	if err != nil {
		return err
	}

	sites, err := fetchSitesForPolicies(cmd, c, accountIDs, siteIDs)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return err
	}

	used := make(map[string]int)
	var pulled int
	for i, s := range sites {
		printProgress("policy", i, len(sites))

		p, pErr := c.PolicyGetSite(cmd.Context(), s.ID)
		if pErr != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: site %s (%s): %v\n", s.Name, s.ID, pErr)
			continue
		}

		pf := sitePolicyToFile(s.ID, s.Name, *p)
		data, mErr := yaml.Marshal(pf)
		if mErr != nil {
			return fmt.Errorf("marshal policy for %s: %w", s.Name, mErr)
		}

		stem := sanitizeFilename(s.Name)
		if n := used[stem]; n > 0 {
			stem = fmt.Sprintf("%s-%d", stem, n)
		}
		used[sanitizeFilename(s.Name)]++

		path := filepath.Join(outDir, stem+".yaml")
		if wErr := os.WriteFile(path, data, 0o644); wErr != nil {
			return wErr
		}
		pulled++
	}
	clearProgress()
	fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(pulled, "policy"), outDir)
	return nil
}

func pullAccountPolicies(cmd *cobra.Command, accountIDs []string, outDir string) error {
	c, err := mgmtClient()
	if err != nil {
		return err
	}

	accounts, err := fetchAccountsForPolicies(cmd, c, accountIDs)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return err
	}

	var pulled int
	for i, a := range accounts {
		printProgress("policy", i, len(accounts))

		p, pErr := c.PolicyGetAccount(cmd.Context(), a.ID)
		if pErr != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: account %s (%s): %v\n", a.Name, a.ID, pErr)
			continue
		}

		pf := accountPolicyToFile(a.ID, *p)
		data, mErr := yaml.Marshal(pf)
		if mErr != nil {
			return fmt.Errorf("marshal policy for account %s: %w", a.Name, mErr)
		}

		fname := fmt.Sprintf("policy_account_%s.yaml", a.ID)
		path := filepath.Join(outDir, fname)
		if wErr := os.WriteFile(path, data, 0o644); wErr != nil {
			return wErr
		}
		pulled++
	}
	clearProgress()
	fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(pulled, "policy"), outDir)
	return nil
}

func pullGroupPolicies(cmd *cobra.Command, siteIDs []string, outDir string) error {
	c, err := mgmtClient()
	if err != nil {
		return err
	}

	if len(siteIDs) == 0 {
		return fmt.Errorf("--site-id is required when --scope=group")
	}

	groups, err := fetchGroupsForPolicies(cmd, c, siteIDs)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return err
	}

	var pulled int
	for i, g := range groups {
		printProgress("policy", i, len(groups))

		p, pErr := c.PolicyGetGroup(cmd.Context(), g.ID)
		if pErr != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: group %s (%s): %v\n", g.Name, g.ID, pErr)
			continue
		}

		pf := groupPolicyToFile(g.SiteID, g.ID, g.Name, *p)
		data, mErr := yaml.Marshal(pf)
		if mErr != nil {
			return fmt.Errorf("marshal policy for group %s: %w", g.Name, mErr)
		}

		fname := fmt.Sprintf("policy_group_%s.yaml", g.ID)
		path := filepath.Join(outDir, fname)
		if wErr := os.WriteFile(path, data, 0o644); wErr != nil {
			return wErr
		}
		pulled++
	}
	clearProgress()
	fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n", pluralize(pulled, "policy"), outDir)
	return nil
}

func fetchAccountsForPolicies(cmd *cobra.Command, c *mgmt.Client, accountIDs []string) ([]mgmt.Account, error) {
	if len(accountIDs) > 0 {
		var accounts []mgmt.Account
		for _, id := range accountIDs {
			a, err := c.AccountsGet(cmd.Context(), id)
			if err != nil {
				return nil, fmt.Errorf("get account %s: %w", id, err)
			}
			accounts = append(accounts, *a)
		}
		return accounts, nil
	}
	params := &mgmt.AccountListParams{Limit: defaultPageSize}
	accounts, _, err := fetchAllREST("account", func(cur string) ([]mgmt.Account, *mgmt.Pagination, error) {
		params.Cursor = cur
		return c.AccountsList(cmd.Context(), params)
	})
	return accounts, err
}

func fetchGroupsForPolicies(cmd *cobra.Command, c *mgmt.Client, siteIDs []string) ([]mgmt.Group, error) {
	params := &mgmt.GroupListParams{
		SiteIDs: siteIDs,
		Limit:   defaultPageSize,
	}
	groups, _, err := fetchAllREST("group", func(cur string) ([]mgmt.Group, *mgmt.Pagination, error) {
		params.Cursor = cur
		return c.GroupsList(cmd.Context(), params)
	})
	return groups, err
}

func newPoliciesPushCmd() *cobra.Command {
	var inDir string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push policies from local YAML files",
		Long: `Read policy YAML files from a directory and update the corresponding policies.

Each file must contain a scope field (site, account, or group) and the matching
scope ID to identify the target. The command fetches the current policy, diffs it
against the desired state, and applies changes.
Dry-run by default — pass --yes to apply changes.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries, err := os.ReadDir(inDir)
			if err != nil {
				return fmt.Errorf("read %s: %w", inDir, err)
			}

			var files []policyFile
			for _, e := range entries {
				if e.IsDir() || (!strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml")) {
					continue
				}
				data, rErr := os.ReadFile(filepath.Join(inDir, e.Name()))
				if rErr != nil {
					return fmt.Errorf("read %s: %w", e.Name(), rErr)
				}
				var pf policyFile
				if uErr := yaml.Unmarshal(data, &pf); uErr != nil {
					return fmt.Errorf("parse %s: %w", e.Name(), uErr)
				}
				if pf.Scope == "" {
					pf.Scope = "site"
				}
				switch pf.Scope {
				case "site":
					if pf.SiteID == "" {
						return fmt.Errorf("policy in %s has scope=site but no siteId", e.Name())
					}
				case "account":
					if pf.AccountID == "" {
						return fmt.Errorf("policy in %s has scope=account but no accountId", e.Name())
					}
				case "group":
					if pf.GroupID == "" || pf.SiteID == "" {
						return fmt.Errorf("policy in %s has scope=group but missing groupId or siteId", e.Name())
					}
				default:
					return fmt.Errorf("policy in %s has invalid scope %q", e.Name(), pf.Scope)
				}
				files = append(files, pf)
			}
			if len(files) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No policy files found.")
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			type policyDiff struct {
				file    policyFile
				changes []string
			}
			var diffs []policyDiff

			for i, pf := range files {
				printProgress("policy", i, len(files))

				current, gErr := getPolicyForScope(cmd, c, pf)
				if gErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: get policy for %s: %v\n", policyScopeLabel(pf), gErr)
					continue
				}

				changes := diffPolicyFields(pf, current)
				if len(changes) > 0 {
					diffs = append(diffs, policyDiff{file: pf, changes: changes})
				}
			}
			clearProgress()

			if len(diffs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "All policies are up to date.")
				return nil
			}

			// In JSON mode the summary goes to stderr so stdout stays valid JSON.
			summaryW := cmd.OutOrStdout()
			if outputFormat == "json" {
				summaryW = cmd.ErrOrStderr()
			}
			for _, d := range diffs {
				fmt.Fprintf(summaryW, "%s:\n", policyScopeLabel(d.file))
				for _, ch := range d.changes {
					fmt.Fprintln(summaryW, ch)
				}
			}

			return guard(cmd.OutOrStdout(), "policies push", "update "+pluralize(len(diffs), "policy"), inDir, yes, func() error {
				var updated int
				for _, d := range diffs {
					payload, mErr := policyUpdatePayload(d.file)
					if mErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Warning: marshal policy for %s: %v\n", policyScopeLabel(d.file), mErr)
						continue
					}
					if uErr := updatePolicyForScope(cmd, c, d.file, payload); uErr != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Warning: update %s: %v\n", policyScopeLabel(d.file), uErr)
						continue
					}
					updated++
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"updated": updated})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated %s\n", pluralize(updated, "policy"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&inDir, "dir", "policies", "directory containing policy YAML files")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

func getPolicyForScope(cmd *cobra.Command, c *mgmt.Client, pf policyFile) (*mgmt.Policy, error) {
	switch pf.Scope {
	case "account":
		return c.PolicyGetAccount(cmd.Context(), pf.AccountID)
	case "group":
		return c.PolicyGetGroup(cmd.Context(), pf.GroupID)
	default:
		return c.PolicyGetSite(cmd.Context(), pf.SiteID)
	}
}

func updatePolicyForScope(cmd *cobra.Command, c *mgmt.Client, pf policyFile, payload json.RawMessage) error {
	var err error
	switch pf.Scope {
	case "account":
		_, err = c.PolicyUpdateAccount(cmd.Context(), pf.AccountID, payload)
	case "group":
		_, err = c.PolicyUpdateGroup(cmd.Context(), pf.GroupID, payload)
	default:
		_, err = c.PolicyUpdateSite(cmd.Context(), pf.SiteID, payload)
	}
	return err
}

func policyScopeLabel(pf policyFile) string {
	switch pf.Scope {
	case "account":
		return fmt.Sprintf("Account %s", pf.AccountID)
	case "group":
		name := pf.GroupName
		if name == "" {
			name = pf.GroupID
		}
		return fmt.Sprintf("Group %s (%s)", name, pf.GroupID)
	default:
		name := pf.SiteName
		if name == "" {
			name = pf.SiteID
		}
		return fmt.Sprintf("Site %s (%s)", name, pf.SiteID)
	}
}

func diffPolicyFields(pf policyFile, current *mgmt.Policy) []string {
	var changes []string
	if pf.MitigationMode != current.MitigationMode {
		changes = append(changes, fmt.Sprintf("  mitigationMode: %s -> %s", current.MitigationMode, pf.MitigationMode))
	}
	if pf.MitigationModeSuspicious != current.MitigationModeSuspicious {
		changes = append(changes, fmt.Sprintf("  mitigationModeSuspicious: %s -> %s", current.MitigationModeSuspicious, pf.MitigationModeSuspicious))
	}
	if pf.AntiTamperingOn != current.AntiTamperingOn {
		changes = append(changes, fmt.Sprintf("  antiTamperingOn: %v -> %v", current.AntiTamperingOn, pf.AntiTamperingOn))
	}
	if pf.NetworkQuarantineOn != current.NetworkQuarantineOn {
		changes = append(changes, fmt.Sprintf("  networkQuarantineOn: %v -> %v", current.NetworkQuarantineOn, pf.NetworkQuarantineOn))
	}
	if pf.SnapshotsOn != current.SnapshotsOn {
		changes = append(changes, fmt.Sprintf("  snapshotsOn: %v -> %v", current.SnapshotsOn, pf.SnapshotsOn))
	}
	if pf.AllowRemoteShell != current.AllowRemoteShell {
		changes = append(changes, fmt.Sprintf("  allowRemoteShell: %v -> %v", current.AllowRemoteShell, pf.AllowRemoteShell))
	}
	if pf.ScanNewAgents != current.ScanNewAgents {
		changes = append(changes, fmt.Sprintf("  scanNewAgents: %v -> %v", current.ScanNewAgents, pf.ScanNewAgents))
	}
	if pf.AutoDecommissionOn != current.AutoDecommissionOn {
		changes = append(changes, fmt.Sprintf("  autoDecommissionOn: %v -> %v", current.AutoDecommissionOn, pf.AutoDecommissionOn))
	}
	if pf.AutoDecommissionDays != current.AutoDecommissionDays {
		changes = append(changes, fmt.Sprintf("  autoDecommissionDays: %d -> %d", current.AutoDecommissionDays, pf.AutoDecommissionDays))
	}
	if pf.Ioc != current.Ioc {
		changes = append(changes, fmt.Sprintf("  ioc: %v -> %v", current.Ioc, pf.Ioc))
	}
	return changes
}
