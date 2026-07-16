package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAppControlRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage application control rules",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAppControlRulesListCmd())
	cmd.AddCommand(newAppControlRulesGetCmd())
	cmd.AddCommand(newAppControlRulesCreateCmd())
	cmd.AddCommand(newAppControlRulesUpdateCmd())
	cmd.AddCommand(newAppControlRulesDeleteCmd())
	addAppControlSyncCmds(cmd)
	return cmd
}

func newAppControlRulesListCmd() *cobra.Command {
	var scopeType, cursor string
	var scopeIDs []string
	var includeParents, all bool
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List application control rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AppControlQueryParams{
				IncludeParents: includeParents,
				PageSize:       limit,
				Cursor:         cursor,
			}
			if scopeType != "" {
				params.ScopeType = mgmt.AppControlScopeLevel(strings.ToUpper(scopeType))
				params.ScopeIDs = scopeIDs
			}
			if params.PageSize == 0 {
				params.PageSize = defaultPageSize
			}

			var rules []mgmt.AppControlRule
			var total int

			if all {
				rules, total, err = fetchAllAppControlRules(cmd, c, params)
			} else {
				var t int
				rules, _, t, err = c.AppControlRulesList(cmd.Context(), params)
				total = t
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Behavior", "OS", "Propagation"}
			rows := make([][]string, len(rules))
			for i, r := range rules {
				rows[i] = []string{
					r.ID,
					truncate(r.RuleName, 40),
					string(r.Behavior),
					joinOSTypes(r.OSType),
					fmt.Sprintf("%t", r.Propagation),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, rules, len(rules), total, "application control rule", all)
		},
	}
	cmd.Flags().StringVar(&scopeType, "scope-type", "", "scope type (account, site, group)")
	cmd.Flags().StringSliceVar(&scopeIDs, "scope-id", nil, "scope IDs")
	cmd.Flags().BoolVar(&includeParents, "include-parents", false, "include rules from parent scopes")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return markJSON(cmd)
}

// fetchAllAppControlRules fetches all pages of app control rules using relay
// cursor pagination.
func fetchAllAppControlRules(cmd *cobra.Command, c *mgmt.Client, params *mgmt.AppControlQueryParams) ([]mgmt.AppControlRule, int, error) {
	var all []mgmt.AppControlRule
	var total int
	for {
		page, cursor, t, err := c.AppControlRulesList(cmd.Context(), params)
		if err != nil {
			return nil, 0, err
		}
		total = t
		all = append(all, page...)
		if cursor == "" || len(page) == 0 {
			break
		}
		params.Cursor = cursor
	}
	return all, total, nil
}

func joinOSTypes(oss []mgmt.AppControlOSType) string {
	s := make([]string, len(oss))
	for i, o := range oss {
		s[i] = string(o)
	}
	return strings.Join(s, ",")
}

func newAppControlRulesGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <rule-id>",
		Short: "Get an application control rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			r, err := c.AppControlRulesGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), r)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"ID", r.ID},
				{"Name", r.RuleName},
				{"Behavior", string(r.Behavior)},
				{"OS", joinOSTypes(r.OSType)},
				{"Propagation", fmt.Sprintf("%t", r.Propagation)},
				{"Description", orDash(r.Description)},
				{"Created At", orDash(r.CreatedAt)},
				{"Created By", orDash(r.CreatedBy)},
			})
			return nil
		},
	}
	return markJSON(cmd)
}

func newAppControlRulesCreateCmd() *cobra.Command {
	var name, description, behavior, scopeType string
	var scopeIDs, osTypes []string
	var propagation bool
	var yes bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an application control rule",
		Long: `Create an application control rule.
Dry-run by default — pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if behavior == "" {
				return fmt.Errorf("--behavior is required")
			}

			return guard(cmd.OutOrStdout(), "applications rules create",
				"create application control rule "+name,
				name, yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					input := mgmt.AppControlRuleInput{
						RuleName:    name,
						Description: &description,
						Behavior:    mgmt.AppControlBehavior(strings.ToUpper(behavior)),
						Propagation: &propagation,
					}
					for _, o := range osTypes {
						input.OSType = append(input.OSType, mgmt.AppControlOSType(strings.ToUpper(o)))
					}
					if scopeType != "" && len(scopeIDs) > 0 {
						input.Scope = &mgmt.AppControlScope{
							ScopeType: mgmt.AppControlScopeLevel(strings.ToUpper(scopeType)),
							ScopeIDs:  scopeIDs,
						}
					}
					resp, err := c.AppControlRulesCreate(cmd.Context(), input)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), resp)
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Created application control rule %s\n", resp.ID)
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "rule name (required)")
	cmd.Flags().StringVar(&description, "description", "", "rule description")
	cmd.Flags().StringVar(&behavior, "behavior", "", "rule behavior: allow, monitor, block (required)")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "OS types: windows, macos")
	cmd.Flags().StringVar(&scopeType, "scope-type", "", "scope type: account, site, group")
	cmd.Flags().StringSliceVar(&scopeIDs, "scope-id", nil, "scope IDs")
	cmd.Flags().BoolVar(&propagation, "propagation", false, "enable propagation")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

func newAppControlRulesUpdateCmd() *cobra.Command {
	var name, description, behavior string
	var osTypes []string
	var propagation bool
	var yes bool

	cmd := &cobra.Command{
		Use:   "update <rule-id>",
		Short: "Update an application control rule",
		Long: `Update an application control rule by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" && behavior == "" {
				return fmt.Errorf("at least --name or --behavior is required")
			}

			return guard(cmd.OutOrStdout(), "applications rules update",
				"update application control rule "+args[0],
				args[0], yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					input := mgmt.AppControlRuleInput{
						RuleName:    name,
						Description: &description,
						Propagation: &propagation,
					}
					if behavior != "" {
						input.Behavior = mgmt.AppControlBehavior(strings.ToUpper(behavior))
					}
					for _, o := range osTypes {
						input.OSType = append(input.OSType, mgmt.AppControlOSType(strings.ToUpper(o)))
					}
					resp, err := c.AppControlRulesUpdate(cmd.Context(), args[0], input)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), resp)
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Updated application control rule %s\n", args[0])
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "rule name")
	cmd.Flags().StringVar(&description, "description", "", "rule description")
	cmd.Flags().StringVar(&behavior, "behavior", "", "rule behavior: allow, monitor, block")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "OS types: windows, macos")
	cmd.Flags().BoolVar(&propagation, "propagation", false, "enable propagation")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}

func newAppControlRulesDeleteCmd() *cobra.Command {
	var scopeType string
	var scopeIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <rule-id>...",
		Short: "Delete application control rules",
		Long: `Delete one or more application control rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "applications rules delete",
				"delete "+pluralize(len(args), "application control rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					var scope *mgmt.AppControlScope
					if scopeType != "" && len(scopeIDs) > 0 {
						scope = &mgmt.AppControlScope{
							ScopeType: mgmt.AppControlScopeLevel(strings.ToUpper(scopeType)),
							ScopeIDs:  scopeIDs,
						}
					}
					resp, err := c.AppControlRulesDelete(cmd.Context(), args, scope)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), resp)
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s\n", pluralize(len(args), "application control rule"))
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&scopeType, "scope-type", "", "scope type: account, site, group")
	cmd.Flags().StringSliceVar(&scopeIDs, "scope-id", nil, "scope IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}
