package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"danny.vn/s1/mgmt"
)

func newAutomationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "automation",
		Short: "Manage hyperautomation workflows and executions",
		Long: `Manage SentinelOne Hyperautomation workflows — list, inspect,
import, run, and control lifecycle (activate/deactivate).`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAutomationListCmd())
	cmd.AddCommand(newAutomationGetCmd())
	cmd.AddCommand(newAutomationVersionsCmd())
	cmd.AddCommand(newAutomationExportCmd())
	cmd.AddCommand(newAutomationCreateCmd())
	cmd.AddCommand(newAutomationRunCmd())
	cmd.AddCommand(newAutomationActivateCmd())
	cmd.AddCommand(newAutomationDeactivateCmd())
	cmd.AddCommand(newAutomationExecutionsCmd())
	cmd.AddCommand(newAutomationExecutionGetCmd())
	cmd.AddCommand(newAutomationExecutionOutputCmd())
	return cmd
}

// --- list ---

func newAutomationListCmd() *cobra.Command {
	var siteIDs, groupIDs, accountIDs, states, triggerTypes, tags []string
	var nameContains, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List automation workflows",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AutomationListParams{
				SiteIDs:      siteIDs,
				GroupIDs:     groupIDs,
				AccountIDs:   accountIDs,
				States:       states,
				TriggerTypes: triggerTypes,
				Tags:         tags,
				NameContains: nameContains,
				Limit:        limit,
				SortBy:       sortBy,
				SortOrder:    sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.WorkflowListItem
			var total int

			if all {
				// Hyperautomation uses skip-based pagination, not cursor.
				// Fetch all manually.
				var skip int
				for {
					params.Skip = skip
					page, pag, err := c.AutomationList(cmd.Context(), params)
					if err != nil {
						return err
					}
					items = append(items, page...)
					if pag != nil {
						total = pag.TotalItems
					}
					if len(page) < params.Limit {
						break
					}
					skip += len(page)
				}
			} else {
				var pag *mgmt.AutomationPagination
				items, pag, err = c.AutomationList(cmd.Context(), params)
				if err != nil {
					return err
				}
				if pag != nil {
					total = pag.TotalItems
				}
			}

			headers := []string{"ID", "Name", "State", "Status", "Trigger", "Scope", "Updated"}
			rows := make([][]string, len(items))
			for i, item := range items {
				wf := item.Workflow
				trigger := ""
				for _, a := range item.Actions {
					if isTrigger(a.Type) {
						trigger = a.Type
						break
					}
				}
				rows[i] = []string{
					wf.ID,
					truncate(wf.Name, 40),
					string(wf.State),
					string(wf.Status),
					trigger,
					string(wf.ScopeLevel),
					orDash(wf.UpdatedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "workflow", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&states, "state", nil, "filter by state (active, inactive, deactivated, draft)")
	cmd.Flags().StringSliceVar(&triggerTypes, "trigger-type", nil, "filter by trigger type")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "filter by tag")
	cmd.Flags().StringVar(&nameContains, "name", "", "filter by name (contains)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order (asc, desc)")
	return markJSON(cmd)
}

// isTrigger reports whether an action type is a trigger type.
func isTrigger(t string) bool {
	switch t {
	case "http_trigger", "scheduled_trigger", "email_trigger",
		"manual_trigger", "singularity_response_trigger", "snippet_trigger":
		return true
	}
	return false
}

// --- get (export) ---

func newAutomationGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <workflow-id> <version-id>",
		Short: "Get a workflow version (export format)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			wf, err := c.AutomationExport(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), wf)
		},
	}
	return markJSON(cmd)
}

// --- versions ---

func newAutomationVersionsCmd() *cobra.Command {
	var siteIDs, groupIDs, accountIDs []string

	cmd := &cobra.Command{
		Use:   "versions <workflow-id>",
		Short: "List versions of a workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AutomationListParams{
				SiteIDs:    siteIDs,
				GroupIDs:   groupIDs,
				AccountIDs: accountIDs,
			}
			versions, err := c.AutomationVersions(cmd.Context(), args[0], params)
			if err != nil {
				return err
			}

			headers := []string{"Version ID", "Name", "State", "Status", "Exec Status", "Updated"}
			rows := make([][]string, len(versions))
			for i, v := range versions {
				rows[i] = []string{
					v.VersionID,
					truncate(v.Name, 40),
					string(v.State),
					string(v.Status),
					string(v.ExecutionStatus),
					orDash(v.UpdatedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, versions, len(versions), len(versions), "version", false)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	return markJSON(cmd)
}

// --- export ---

func newAutomationExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <workflow-id> <version-id>",
		Short: "Export a workflow version as JSON (suitable for import)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			wf, err := c.AutomationExport(cmd.Context(), args[0], args[1])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), wf)
		},
	}
	return markJSON(cmd)
}

// --- create (import from file) ---

// automationFile is the on-disk format for workflow import.
type automationFile struct {
	Name        string           `yaml:"name" json:"name"`
	Description string           `yaml:"description" json:"description"`
	Actions     []map[string]any `yaml:"actions" json:"actions"`
}

func newAutomationCreateCmd() *cobra.Command {
	var fromFile string
	var siteIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "create --from-file <workflow.json>",
		Short: "Create (import) a workflow from a file",
		Long: `Import a workflow definition that was previously exported.
The file should be JSON or YAML matching the export format.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			raw, err := os.ReadFile(fromFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", fromFile, err)
			}
			// Parse to extract name for guard message.
			var f automationFile
			if err := yaml.Unmarshal(raw, &f); err != nil {
				return fmt.Errorf("parse %s: %w", fromFile, err)
			}
			if f.Name == "" {
				return fmt.Errorf("workflow file %s has no name", fromFile)
			}

			// Re-marshal as JSON for the API (yaml.Unmarshal accepts JSON too).
			data, err := json.Marshal(f)
			if err != nil {
				return fmt.Errorf("marshal workflow: %w", err)
			}

			action := fmt.Sprintf("import workflow %q from %s", f.Name, fromFile)
			return guard(cmd.OutOrStdout(), "automation create", action, fromFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				wf, err := c.AutomationImport(cmd.Context(), json.RawMessage(data), siteIDs)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), wf)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Imported workflow %s (%s)\n", wf.ID, wf.Name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "workflow definition file, JSON or YAML (required)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope to site ID")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

// --- run ---

func newAutomationRunCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "run <workflow-id> <version-id>",
		Short: "Trigger a manual workflow execution",
		Long: `Trigger a manual or scheduled workflow execution. This executes
tenant-side automation and may perform actions such as isolating agents,
sending emails, or calling external APIs.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID, versionID := args[0], args[1]
			action := fmt.Sprintf("run workflow %s (version %s)", workflowID, versionID)
			return guard(cmd.OutOrStdout(), "automation run", action, workflowID, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				exec, err := c.AutomationRun(cmd.Context(), workflowID, versionID, nil)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), exec)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Started execution %s (state: %s)\n", exec.ID, exec.State)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

// --- activate ---

func newAutomationActivateCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "activate <workflow-id> <version-id>",
		Short: "Activate a workflow version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID, versionID := args[0], args[1]
			action := fmt.Sprintf("activate workflow %s version %s", workflowID, versionID)
			return guard(cmd.OutOrStdout(), "automation activate", action, workflowID, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.AutomationActivate(cmd.Context(), workflowID, versionID); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{
						"action":      "activate",
						"workflow_id": workflowID,
						"version_id":  versionID,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Activated workflow %s version %s\n", workflowID, versionID)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

// --- deactivate ---

func newAutomationDeactivateCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "deactivate <workflow-id>",
		Short: "Deactivate the active version of a workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workflowID := args[0]
			action := fmt.Sprintf("deactivate workflow %s", workflowID)
			return guard(cmd.OutOrStdout(), "automation deactivate", action, workflowID, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.AutomationDeactivate(cmd.Context(), workflowID); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{
						"action":      "deactivate",
						"workflow_id": workflowID,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deactivated workflow %s\n", workflowID)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

// --- executions ---

func newAutomationExecutionsCmd() *cobra.Command {
	var siteIDs, groupIDs, accountIDs, states, triggerTypes []string
	var workflowID, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "executions",
		Short: "List workflow executions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AutomationExecutionListParams{
				SiteIDs:      siteIDs,
				GroupIDs:     groupIDs,
				AccountIDs:   accountIDs,
				States:       states,
				TriggerTypes: triggerTypes,
				WorkflowID:   workflowID,
				Limit:        limit,
				SortBy:       sortBy,
				SortOrder:    sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.WorkflowExecution
			var total int

			if all {
				var skip int
				for {
					params.Skip = skip
					page, pag, err := c.AutomationExecutions(cmd.Context(), params)
					if err != nil {
						return err
					}
					items = append(items, page...)
					if pag != nil {
						total = pag.TotalItems
					}
					if len(page) < params.Limit {
						break
					}
					skip += len(page)
				}
			} else {
				var pag *mgmt.AutomationPagination
				items, pag, err = c.AutomationExecutions(cmd.Context(), params)
				if err != nil {
					return err
				}
				if pag != nil {
					total = pag.TotalItems
				}
			}

			headers := []string{"ID", "Workflow", "State", "Trigger", "Duration", "Actions", "Created"}
			rows := make([][]string, len(items))
			for i, e := range items {
				rows[i] = []string{
					e.ID,
					truncate(e.WorkflowName, 30),
					string(e.State),
					string(e.Trigger),
					orDash(e.Duration),
					fmt.Sprintf("%d", e.ExecutedActions),
					orDash(e.CreatedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "execution", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&states, "state", nil, "filter by state (Running, Completed, Error, etc.)")
	cmd.Flags().StringSliceVar(&triggerTypes, "trigger-type", nil, "filter by trigger type")
	cmd.Flags().StringVar(&workflowID, "workflow-id", "", "filter by workflow ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort order (asc, desc)")
	return markJSON(cmd)
}

// --- execution get ---

func newAutomationExecutionGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution-get <execution-id>",
		Short: "Get a workflow execution by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			detail, err := c.AutomationExecutionGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), detail)
		},
	}
	return markJSON(cmd)
}

// --- execution output ---

func newAutomationExecutionOutputCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution-output <execution-id>",
		Short: "Get the output of a workflow execution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			out, err := c.AutomationExecutionOutput(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), out)
		},
	}
	return markJSON(cmd)
}
