package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/graphql"
	"danny.vn/s1/internal/reconcile"
)

func addCloudPolicySyncCmds(parent *cobra.Command) {
	spec := cloudPoliciesSpec()
	parent.AddCommand(newEnginePullCmd(spec))
	parent.AddCommand(newEnginePushCmd(spec))
}

// cloudPolicyFile is the YAML representation of a cloud policy on disk. Cloud
// policies cannot be created through this surface, so the file carries only the
// identity (ID plus name for readability) and the reconciled status.
type cloudPolicyFile struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	Status string `yaml:"status"`
}

func cloudPolicyToFile(p graphql.CloudPolicy) cloudPolicyFile {
	return cloudPolicyFile{
		ID:     p.ID,
		Name:   p.Name,
		Status: p.Status,
	}
}

// decodeCloudPolicy maps one local file to a canonical Object: it validates the
// policy ID and re-marshals through cloudPolicyFile so bodies are byte-equal to
// the ones List produces. Identity is the policy ID.
func decodeCloudPolicy(data []byte) (reconcile.Object, error) {
	var f cloudPolicyFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return reconcile.Object{}, err
	}
	if f.ID == "" {
		return reconcile.Object{}, fmt.Errorf("cloud policy in file has no id")
	}
	body, err := yaml.Marshal(f)
	if err != nil {
		return reconcile.Object{}, err
	}
	return reconcile.Object{Name: f.ID, ID: f.ID, Body: body}, nil
}

// cloudPoliciesSpec adapts cloud security policies to the shared sync builders.
// Identity is the policy ID; the surface is status-reconcile only (NoCreate), so
// an update toggles the policy between enabled and disabled.
func cloudPoliciesSpec() surfaceSpec {
	return surfaceSpec{
		Noun:       "cloud policy",
		Command:    "cloud-policies",
		DefaultDir: "cloud-policies",
		PullShort:  "Pull cloud security policies to local YAML files",
		PullLong: `Fetch all cloud security policies and write them as YAML files.

Each policy produces one file carrying its ID, name, and status. Cloud policies
cannot be created through this surface, so push only reconciles status.`,
		PushShort: "Push cloud security policy status from local YAML files",
		PushLong: `Read cloud policy YAML files from a directory and reconcile their status.

Policies are matched by ID: a status change (enabled/disabled) is applied,
unchanged policies are skipped, and a local file whose ID has no live match
fails per-item since policies cannot be created through this surface. Dry-run by
default — pass --yes to apply changes.`,
		Build: func(_ *cobra.Command, _ scopeFlags) (reconcile.Surface, error) {
			var client *graphql.Client
			getClient := func() (*graphql.Client, error) {
				if client == nil {
					c, err := gqlClient()
					if err != nil {
						return nil, err
					}
					client = c
				}
				return client, nil
			}

			return reconcile.Surface{
				Name:    "cloud policy",
				Command: "cloud-policies",
				Caps:    reconcile.Capabilities{NoCreate: true},
				Decode:  decodeCloudPolicy,
				List: func(ctx context.Context) ([]reconcile.Object, error) {
					c, err := getClient()
					if err != nil {
						return nil, err
					}
					policies, _, lErr := fetchAllGQL("cloud policy", func(after string) (*graphql.Connection[graphql.CloudPolicy], error) {
						return c.CloudPoliciesList(ctx, &graphql.ListParams{First: 100, After: after})
					})
					if lErr != nil {
						return nil, lErr
					}
					objs := make([]reconcile.Object, 0, len(policies))
					for _, p := range policies {
						body, mErr := yaml.Marshal(cloudPolicyToFile(p))
						if mErr != nil {
							return nil, fmt.Errorf("marshal cloud policy %s: %w", p.ID, mErr)
						}
						objs = append(objs, reconcile.Object{Name: p.ID, ID: p.ID, Body: body})
					}
					return objs, nil
				},
				// Create is nil: cloud policies cannot be created here (NoCreate).
				Update: func(ctx context.Context, id string, local reconcile.Object) error {
					c, err := getClient()
					if err != nil {
						return err
					}
					var f cloudPolicyFile
					if uErr := yaml.Unmarshal(local.Body, &f); uErr != nil {
						return uErr
					}
					switch {
					case strings.EqualFold(f.Status, "enabled"):
						_, aErr := c.CloudPoliciesEnable(ctx, []string{id})
						return aErr
					case strings.EqualFold(f.Status, "disabled"):
						_, aErr := c.CloudPoliciesDisable(ctx, []string{id})
						return aErr
					default:
						return fmt.Errorf("unrecognized status %q", f.Status)
					}
				},
			}, nil
		},
	}
}
