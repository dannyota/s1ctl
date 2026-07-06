package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

type cloudPolicyActionFn func(*graphql.Client, *cobra.Command, []string) (*graphql.CloudPoliciesActionResponse, error)

func addCloudPolicyActions(parent *cobra.Command) {
	parent.AddCommand(newCloudPolicyActionCmd("enable", "Enable cloud security policies", func(c *graphql.Client, cmd *cobra.Command, ids []string) (*graphql.CloudPoliciesActionResponse, error) {
		return c.CloudPoliciesEnable(cmd.Context(), ids)
	}))
	parent.AddCommand(newCloudPolicyActionCmd("disable", "Disable cloud security policies", func(c *graphql.Client, cmd *cobra.Command, ids []string) (*graphql.CloudPoliciesActionResponse, error) {
		return c.CloudPoliciesDisable(cmd.Context(), ids)
	}))
	parent.AddCommand(newCloudPolicyActionCmd("delete", "Delete cloud security policies", func(c *graphql.Client, cmd *cobra.Command, ids []string) (*graphql.CloudPoliciesActionResponse, error) {
		return c.CloudPoliciesDelete(cmd.Context(), ids)
	}))
}

func newCloudPolicyActionCmd(verb, short string, fn cloudPolicyActionFn) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <id> [id...]",
		Short: short,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "cloud-policies "+verb, verb+" "+pluralize(len(args), "cloud policy"), strings.Join(args, ","), yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				resp, err := fn(c, cmd, args)
				if err != nil {
					return err
				}
				affectedIDs := []string{}
				if resp != nil && resp.IDs != nil {
					affectedIDs = resp.IDs
				}
				affected := len(affectedIDs)
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{
						"action":   verb,
						"affected": affected,
						"ids":      affectedIDs,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n",
					verb, pluralize(affected, "cloud policy"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}
