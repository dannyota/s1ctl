package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newTagRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag-rules",
		Short: "Manage dynamic asset tag rules",
		Long: `Manage dynamic asset tag rules.

A tag rule automatically applies tags to XDR inventory assets that match a set
of conditions. Rule bodies (conditions, scopes, tags) are nested structures, so
create, update, and test read a JSON file (--from-file) rather than flags.

Use 'tag-rules test --from-file' to preview how many assets a candidate rule
matches before saving it.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newTagRulesListCmd())
	cmd.AddCommand(newTagRulesCreateCmd())
	cmd.AddCommand(newTagRulesUpdateCmd())
	cmd.AddCommand(newTagRulesDeleteCmd())
	cmd.AddCommand(newTagRulesTestCmd())
	return cmd
}

// readTagRuleFile reads a tag rule JSON file into the write body. conditions,
// scopes, and tags are captured as raw JSON.
func readTagRuleFile(path string) (mgmt.TagRuleWrite, error) {
	var body mgmt.TagRuleWrite
	raw, err := os.ReadFile(path)
	if err != nil {
		return body, fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return body, fmt.Errorf("parse %s: %w", path, err)
	}
	if body.Name == "" {
		return body, fmt.Errorf("tag rule file %s has no name", path)
	}
	return body, nil
}

func newTagRulesListCmd() *cobra.Command {
	var siteIDs, accountIDs []string
	var name, status, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List dynamic tag rules",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.TagRuleListParams{
				Name:       name,
				Status:     status,
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				Limit:      limit,
				Cursor:     cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var rules []mgmt.TagRule
			var total int
			if all {
				rules, total, err = fetchAllREST("tag rule", func(cur string) ([]mgmt.TagRule, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.TagRulesList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				rules, pag, err = c.TagRulesList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Status", "Description"}
			rows := make([][]string, len(rules))
			for i, rr := range rules {
				rows[i] = []string{rr.ID, rr.Name, rr.Status, truncate(rr.Description, 40)}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, rules, len(rules), total, "tag rule", all)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "filter by rule name")
	cmd.Flags().StringVar(&status, "status", "", "filter by status (enabled, disabled)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}

func newTagRulesCreateCmd() *cobra.Command {
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "create --from-file <rule.json>",
		Short: "Create a dynamic tag rule from a JSON file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			body, err := readTagRuleFile(fromFile)
			if err != nil {
				return err
			}
			action := fmt.Sprintf("create tag rule %q from %s", body.Name, fromFile)
			return guard(cmd.OutOrStdout(), "tag-rules create", action, fromFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.TagRulesCreate(cmd.Context(), body)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created tag rule %s (%s)\n", created.ID, created.Name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "tag rule definition JSON file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newTagRulesUpdateCmd() *cobra.Command {
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update <rule-id> --from-file <rule.json>",
		Short: "Update a dynamic tag rule from a JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			body, err := readTagRuleFile(fromFile)
			if err != nil {
				return err
			}
			body.ID = args[0]
			action := fmt.Sprintf("update tag rule %s from %s", args[0], fromFile)
			return guard(cmd.OutOrStdout(), "tag-rules update", action, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				updated, err := c.TagRulesUpdate(cmd.Context(), body)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), updated)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated tag rule %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "tag rule definition JSON file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newTagRulesDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <rule-id>",
		Short: "Delete a dynamic tag rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "tag-rules delete", "delete tag rule "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.TagRulesDelete(cmd.Context(), []string{args[0]}); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted tag rule %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newTagRulesTestCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "test --from-file <rule.json>",
		Short: "Report how many assets a candidate tag rule matches",
		Long: `Report how many inventory assets a candidate tag rule would match, without
saving it. This is a read-only dry-check against live inventory.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			body, err := readTagRuleFile(fromFile)
			if err != nil {
				return err
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			count, err := c.TagRulesTest(cmd.Context(), body)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]int{"matches": count})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Tag rule %q matches %d assets\n", body.Name, count)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "tag rule definition JSON file (required)")
	return cmd
}
