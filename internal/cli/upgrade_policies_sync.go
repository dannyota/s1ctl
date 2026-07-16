package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/internal/reconcile"
	"danny.vn/s1/mgmt"
)

func addUpgradePolicySyncCmds(parent *cobra.Command) {
	spec := upgradePoliciesSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// upgradePolicyFile is the YAML representation of an upgrade policy on disk. It
// holds the declarative fields only; server-assigned ID, priority, and
// timestamps are omitted.
type upgradePolicyFile struct {
	Name         string               `yaml:"name"`
	Description  string               `yaml:"description,omitempty"`
	OSType       string               `yaml:"osType"`
	ScopeLevel   string               `yaml:"scopeLevel"`
	ScopeID      string               `yaml:"scopeId,omitempty"`
	IsActive     bool                 `yaml:"isActive"`
	IsScheduled  bool                 `yaml:"isScheduled"`
	AllEndpoints bool                 `yaml:"allEndpoints"`
	MaxRetries   int                  `yaml:"maxRetries"`
	Package      upgradePolicyPkgFile `yaml:"package"`
	Tags         []string             `yaml:"tags,omitempty"`
}

// upgradePolicyPkgFile is the package reference within an upgrade policy file.
type upgradePolicyPkgFile struct {
	FileID string `yaml:"fileId"`
	Major  string `yaml:"major"`
	Minor  string `yaml:"minor"`
	Build  string `yaml:"build"`
}

func upgradePolicyToFile(p mgmt.UpgradePolicy) upgradePolicyFile {
	return upgradePolicyFile{
		Name:         p.Name,
		Description:  p.Description,
		OSType:       p.OSType,
		ScopeLevel:   p.ScopeLevel,
		ScopeID:      p.ScopeID,
		IsActive:     p.IsActive,
		IsScheduled:  p.IsScheduled,
		AllEndpoints: p.AllEndpoints,
		MaxRetries:   p.MaxRetries,
		Package: upgradePolicyPkgFile{
			FileID: p.Package.FileID,
			Major:  p.Package.Major,
			Minor:  p.Package.Minor,
			Build:  p.Package.Build,
		},
		Tags: p.Tags,
	}
}

func (f upgradePolicyFile) toCreate() mgmt.UpgradePolicyCreate {
	return mgmt.UpgradePolicyCreate{
		Name:         f.Name,
		Description:  f.Description,
		OSType:       mgmt.UpgradePolicyOSType(f.OSType),
		ScopeLevel:   mgmt.UpgradePolicyScopeLevel(f.ScopeLevel),
		ScopeID:      f.ScopeID,
		IsActive:     f.IsActive,
		IsScheduled:  f.IsScheduled,
		AllEndpoints: f.AllEndpoints,
		MaxRetries:   f.MaxRetries,
		Package: mgmt.UpgradePolicyPkg{
			FileID: f.Package.FileID,
			Major:  f.Package.Major,
			Minor:  f.Package.Minor,
			Build:  f.Package.Build,
		},
		Tags: f.Tags,
	}
}

// decodeUpgradePolicy maps one local file to a canonical Object. Identity is
// the policy name. The body is re-marshalled so it is byte-equal to what List
// produces.
func decodeUpgradePolicy(data []byte) (reconcile.Object, error) {
	var f upgradePolicyFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.Name == "" {
		return reconcile.Object{}, fmt.Errorf("upgrade policy has no name")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: f.Name, Body: body}, nil
}

// upgradePoliciesSpec adapts upgrade policies to the shared sync builders.
// Identity is the policy name within a scope+OS partition. The list endpoint
// requires --scope-level and --os-type, so pull and push operate on one
// partition at a time.
func upgradePoliciesSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "upgrade policy",
		Command:    "upgrade-policies",
		DefaultDir: "upgrade-policies",
		PullShort:  "Pull upgrade policies to local YAML files",
		PullLong: `Fetch upgrade policies and write them as YAML files.

Each policy produces one file named by its sanitized name. Server-only metadata
(ID, priority, timestamps) is omitted so the files contain only the declarative
definition.

The API requires --scope-level and --os-type to list policies. Pull fetches the
specified partition; use multiple invocations for different OS types or scopes.`,
		PushShort: "Push upgrade policies from local YAML files",
		PushLong: `Read upgrade policy YAML files from a directory and sync them to SentinelOne.

Policies are matched by name: existing policies are updated, new ones are
created, and unchanged ones are skipped. Dry-run by default — pass --yes to
apply changes.

The API requires --scope-level and --os-type to list live policies for matching.
These must match the scope and OS in the local files.`,
		RegisterPullFlags: upgradePolicyScopeFlags,
		RegisterPushFlags: upgradePolicyScopeFlags,
		Build: func(cmd *cobra.Command, scope scopeFlags) (reconcile.Surface, error) {
			var scopeLevel, osType, scopeID string
			if f := cmd.Flags().Lookup("scope-level"); f != nil {
				scopeLevel = f.Value.String()
			}
			if f := cmd.Flags().Lookup("os-type"); f != nil {
				osType = f.Value.String()
			}
			if f := cmd.Flags().Lookup("scope-id"); f != nil {
				scopeID = f.Value.String()
			}
			if scopeLevel == "" {
				return reconcile.Surface{}, fmt.Errorf("--scope-level is required (account, group, site, tenant)")
			}
			if osType == "" {
				return reconcile.Surface{}, fmt.Errorf("--os-type is required (linux, macos, windows)")
			}

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
				Name:    "upgrade policy",
				Command: "upgrade-policies",
				Decode:  decodeUpgradePolicy,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					params := &mgmt.UpgradePolicyListParams{
						ScopeLevel: scopeLevel,
						OSType:     osType,
						Limit:      200,
						SortBy:     "priority",
						SortOrder:  "asc",
					}
					if scopeID != "" {
						params.ScopeID = scopeID
					}
					policies, _, lErr := fetchAllUpgradePoliciesCtx(ctx, c, params)
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(policies))
					for _, p := range policies {
						f := upgradePolicyToFile(p)
						body, mErr := yaml.Marshal(f)
						if mErr != nil {
							return nil, fmt.Errorf("marshal upgrade policy %s: %w", p.Name, mErr)
						}
						objs = append(objs, reconcile.Object{Name: p.Name, ID: p.ID, Body: body})
					}
					return objs, nil
				},
				Create: func(ctx context.Context, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f upgradePolicyFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					return c.UpgradePoliciesCreate(ctx, f.toCreate())
				},
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f upgradePolicyFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					return c.UpgradePoliciesUpdate(ctx, id, f.toCreate())
				},
			}, nil
		},
	}
}

// upgradePolicyScopeFlags registers the --scope-level, --os-type, and --scope-id
// flags shared by upgrade-policies pull and push. --scope-id is a single-value
// flag because upgrade-policies sync operates on one scope at a time.
func upgradePolicyScopeFlags(cmd *cobra.Command, scope *scopeFlags) {
	cmd.Flags().String("scope-level", "", "scope level (account, group, site, tenant) [required]")
	cmd.Flags().String("os-type", "", "OS type (linux, macos, windows) [required]")
	cmd.Flags().String("scope-id", "", "scope ID (site/account/group ID)")
}

// fetchAllUpgradePoliciesCtx pages through upgrade policies using the context.
func fetchAllUpgradePoliciesCtx(ctx context.Context, c *mgmt.Client, params *mgmt.UpgradePolicyListParams) ([]mgmt.UpgradePolicy, int, error) {
	var all []mgmt.UpgradePolicy
	var total int
	params.Skip = 0
	for {
		items, t, err := c.UpgradePoliciesList(ctx, params)
		if err != nil {
			return nil, 0, err
		}
		total = t
		all = append(all, items...)
		if len(all) >= total || len(items) == 0 {
			break
		}
		params.Skip = len(all)
	}
	return all, total, nil
}
