package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func addCloudRuleMutations(parent *cobra.Command) {
	parent.AddCommand(newCloudRuleCreateCmd())
	parent.AddCommand(newCloudRuleUpdateCmd())
	parent.AddCommand(newCloudRuleActionCmd("enable", graphql.CNSRuleActionEnable, "Enable CNS custom cloud rules"))
	parent.AddCommand(newCloudRuleActionCmd("disable", graphql.CNSRuleActionDisable, "Disable CNS custom cloud rules"))
	parent.AddCommand(newCloudRuleActionCmd("delete", graphql.CNSRuleActionDelete, "Delete CNS custom cloud rules"))
	parent.AddCommand(newCloudRuleEvaluateCmd())
}

// readRuleJSONFile reads a file and rejects it if it is not valid JSON, so a bad
// rule body fails locally before any API call.
func readRuleJSONFile(path string) (json.RawMessage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if !json.Valid(data) {
		return nil, fmt.Errorf("%s: not valid JSON", path)
	}
	return json.RawMessage(data), nil
}

func newCloudRuleActionCmd(verb string, action graphql.CNSRuleAction, short string) *cobra.Command {
	var yes bool
	var scopeLevel, scopeID string

	cmd := &cobra.Command{
		Use:   verb + " <id> [id...]",
		Short: short,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			return guard(cmd.OutOrStdout(), "cloud-rules "+verb, verb+" "+pluralize(len(args), "cns rule"), strings.Join(args, ","), yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				resp, err := c.CNSRulesAction(cmd.Context(), action, args, scope)
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
					verb, pluralize(affected, "cns rule"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	addCloudRuleScopeFlags(cmd, &scopeLevel, &scopeID)
	return cmd
}

func newCloudRuleCreateCmd() *cobra.Command {
	var fromFile, scopeLevel, scopeID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "create --from-file <rule.json>",
		Short: "Create a CNS custom cloud rule from a JSON file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			return guard(cmd.OutOrStdout(), "cloud-rules create", "create CNS rule from "+fromFile, fromFile, yes, func() error {
				input, err := readRuleJSONFile(fromFile)
				if err != nil {
					return err
				}
				c, err := gqlClient()
				if err != nil {
					return err
				}
				resp, err := c.CNSRuleCreate(cmd.Context(), input, scope)
				if err != nil {
					return err
				}
				id := ""
				if resp != nil {
					id = resp.ID
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{"created": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created CNS rule %s\n", orDash(id))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "path to rule JSON file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply (default: dry-run)")
	addCloudRuleScopeFlags(cmd, &scopeLevel, &scopeID)
	return cmd
}

func newCloudRuleUpdateCmd() *cobra.Command {
	var fromFile, scopeLevel, scopeID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update <id> --from-file <rule.json>",
		Short: "Replace a CNS custom cloud rule from a JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}
			return guard(cmd.OutOrStdout(), "cloud-rules update", "update CNS rule "+args[0]+" from "+fromFile, args[0], yes, func() error {
				input, err := readRuleJSONFile(fromFile)
				if err != nil {
					return err
				}
				c, err := gqlClient()
				if err != nil {
					return err
				}
				ok, err := c.CNSRuleUpdate(cmd.Context(), args[0], input, scope)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{"id": args[0], "updated": ok})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated CNS rule %s (%s)\n", args[0], boolIcon(ok))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "path to rule JSON file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply (default: dry-run)")
	addCloudRuleScopeFlags(cmd, &scopeLevel, &scopeID)
	return cmd
}

func newCloudRuleEvaluateCmd() *cobra.Command {
	var ruleFile, resourceFile, query, config, policyID, scopeLevel, scopeID string

	cmd := &cobra.Command{
		Use:   "evaluate --rule <rule.json> --resource <resource.json>",
		Short: "Evaluate a Rego query against asset JSON (dry-check)",
		Long: `Evaluate a raw Rego query against an asset's JSON before creating or
updating a CNS rule. This is a read-only dry-check: it evaluates only and
mutates nothing.

The Rego query comes from the --rule file's rawQuery field, or from --query
when set. The asset JSON to test against is read from --resource.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if resourceFile == "" {
				return fmt.Errorf("--resource is required")
			}
			if query == "" && ruleFile == "" {
				return fmt.Errorf("--rule or --query is required")
			}
			scope, err := alertsScope(scopeLevel, scopeID)
			if err != nil {
				return err
			}

			regoQuery, ruleConfig, polID := query, config, policyID
			if ruleFile != "" {
				data, rErr := os.ReadFile(ruleFile)
				if rErr != nil {
					return rErr
				}
				var f struct {
					PolicyID             string `json:"policyId"`
					RawQuery             string `json:"rawQuery"`
					RuleConfigParameters string `json:"ruleConfigParameters"`
				}
				if uErr := json.Unmarshal(data, &f); uErr != nil {
					return fmt.Errorf("%s: %w", ruleFile, uErr)
				}
				if regoQuery == "" {
					regoQuery = f.RawQuery
				}
				if ruleConfig == "" {
					ruleConfig = f.RuleConfigParameters
				}
				if polID == "" {
					polID = f.PolicyID
				}
			}
			if regoQuery == "" {
				return fmt.Errorf("no rego query: set --query or provide rawQuery in --rule file")
			}

			resource, err := readRuleJSONFile(resourceFile)
			if err != nil {
				return err
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			resp, err := c.CNSRuleEvaluate(cmd.Context(), polID, regoQuery, resource, ruleConfig, scope)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), resp)
			}
			result, evalErr := "", ""
			if resp != nil {
				result, evalErr = resp.Result, resp.Error
			}
			printTable([]string{"Field", "Value"}, [][]string{
				{"Result", orDash(result)},
				{"Error", orDash(evalErr)},
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&ruleFile, "rule", "", "rule JSON file (rawQuery extracted from it)")
	cmd.Flags().StringVar(&resourceFile, "resource", "", "asset JSON file to evaluate against (required)")
	cmd.Flags().StringVar(&query, "query", "", "inline Rego query (overrides rule file rawQuery)")
	cmd.Flags().StringVar(&config, "config", "", "inline rule config parameters JSON string")
	cmd.Flags().StringVar(&policyID, "policy-id", "", "policy ID to source mandatory parameters")
	addCloudRuleScopeFlags(cmd, &scopeLevel, &scopeID)
	return cmd
}
