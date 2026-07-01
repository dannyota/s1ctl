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

// policyFile is the YAML representation of a site policy on disk.
type policyFile struct {
	SiteID   string `yaml:"siteId"`
	SiteName string `yaml:"siteName"`

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

func policyToFile(siteID, siteName string, p mgmt.Policy) policyFile {
	return policyFile{
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
	var outDir string

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull site policies to local YAML files",
		Long: `Fetch endpoint policies for each site and write them as YAML files.

Each site produces one file named by its sanitized site name (e.g. production.yaml).
The YAML includes site metadata (siteId, siteName) for push matching, plus the key
policy fields: mitigationMode, antiTamperingOn, networkQuarantineOn, etc.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
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

				pf := policyToFile(s.ID, s.Name, *p)
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
		},
	}
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&outDir, "out", "policies", "output directory")
	return cmd
}

func newPoliciesPushCmd() *cobra.Command {
	var inDir string
	var yes bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push site policies from local YAML files",
		Long: `Read policy YAML files from a directory and update the corresponding site policies.

Each file must contain a siteId field to identify the target site. The command
fetches the current policy, diffs it against the desired state, and applies changes.
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
				if pf.SiteID == "" {
					return fmt.Errorf("policy in %s has no siteId", e.Name())
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

			// Fetch current policies and diff against desired state.
			type policyDiff struct {
				file    policyFile
				changes []string
			}
			var diffs []policyDiff

			for i, pf := range files {
				printProgress("policy", i, len(files))

				current, gErr := c.PolicyGetSite(cmd.Context(), pf.SiteID)
				if gErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: get policy for site %s: %v\n", pf.SiteID, gErr)
					continue
				}

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

				if len(changes) > 0 {
					diffs = append(diffs, policyDiff{file: pf, changes: changes})
				}
			}
			clearProgress()

			if len(diffs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "All policies are up to date.")
				return nil
			}

			// Print diff summary.
			for _, d := range diffs {
				name := d.file.SiteName
				if name == "" {
					name = d.file.SiteID
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Site %s (%s):\n", name, d.file.SiteID)
				for _, ch := range d.changes {
					fmt.Fprintln(cmd.OutOrStdout(), ch)
				}
			}

			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "\nWould update %s. Pass --yes to apply.\n",
					pluralize(len(diffs), "policy"))
				return nil
			}

			var updated int
			for _, d := range diffs {
				payload, mErr := policyUpdatePayload(d.file)
				if mErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: marshal policy for site %s: %v\n", d.file.SiteID, mErr)
					continue
				}
				if _, uErr := c.PolicyUpdateSite(cmd.Context(), d.file.SiteID, payload); uErr != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: update site %s: %v\n", d.file.SiteID, uErr)
					continue
				}
				updated++
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Updated %s\n", pluralize(updated, "policy"))
			return nil
		},
	}
	cmd.Flags().StringVar(&inDir, "dir", "policies", "directory containing policy YAML files")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}
