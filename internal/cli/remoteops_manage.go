package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// readJSONFile reads a JSON request-body file into dst, surfacing the file name
// on read and parse errors.
func readJSONFile(path string, dst any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func newRemoteOpsUpdateCmd() *cobra.Command {
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update <script-id> --from-file <script.json>",
		Short: "Update a remote script's metadata",
		Long: `Update the metadata of a remote script (name, type, OS types, timeout, and
input requirements) from a JSON file. This changes the script's properties but
not its content.

The file holds the "data" object of the update body, for example:

  {
    "scriptName": "Collect Logs",
    "scriptType": "dataCollection",
    "osTypes": ["linux", "macos"],
    "inputRequired": false,
    "inputExample": "-",
    "inputInstructions": "-",
    "scriptRuntimeTimeoutSeconds": 3600
  }`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			var data mgmt.RemoteScriptUpdateData
			if err := readJSONFile(fromFile, &data); err != nil {
				return err
			}
			action := fmt.Sprintf("update script %s from %s", args[0], fromFile)
			return guard(cmd.OutOrStdout(), "remoteops update", action, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				updated, err := c.RemoteScriptsUpdate(cmd.Context(), args[0], mgmt.RemoteScriptUpdate{Data: data})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), updated)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated script %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "JSON file with the update data object (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newRemoteOpsContentCmd() *cobra.Command {
	var outFile string

	cmd := &cobra.Command{
		Use:   "content <script-id>",
		Short: "Print a remote script's content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			content, err := c.RemoteScriptContent(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outFile != "" {
				if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil { //nolint:gosec
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Wrote script content to %s\n", outFile)
				return nil
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]string{"scriptId": args[0], "scriptContent": content})
			}
			fmt.Fprint(cmd.OutOrStdout(), content)
			if !strings.HasSuffix(content, "\n") {
				fmt.Fprintln(cmd.OutOrStdout())
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&outFile, "out", "", "write script content to file (default: stdout)")
	return markJSON(cmd)
}

func newRemoteOpsUploadLimitsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-limits",
		Short: "Show package upload size limits",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			limits, err := c.RemoteScriptsUploadLimits(cmd.Context())
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), limits)
		},
	}
	return markJSON(cmd)
}

func newRemoteOpsPendingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending",
		Short: "Manage pending remote-script executions awaiting approval",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newRemoteOpsPendingListCmd())
	cmd.AddCommand(newRemoteOpsPendingDecisionCmd(true))
	cmd.AddCommand(newRemoteOpsPendingDecisionCmd(false))
	return cmd
}

func newRemoteOpsPendingListCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var sortBy, sortOrder, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List pending executions awaiting approval",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.RemoteScriptsPendingParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
				Limit:      limit,
				Cursor:     cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.PendingExecution
			var total int
			if all {
				items, total, err = fetchAllREST("pending execution", func(cur string) ([]mgmt.PendingExecution, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.RemoteScriptsPendingList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				items, pag, err = c.RemoteScriptsPendingList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "State", "Script", "Task", "Endpoints", "Creator", "Created"}
			rows := make([][]string, len(items))
			for i, it := range items {
				rows[i] = []string{
					it.PendingExecutionID, string(it.State), orDash(it.ScriptData.ScriptName),
					orDash(it.ExecutionData.TaskDescription), strconv.Itoa(it.TotalEndpoints),
					orDash(it.Creator), orDash(it.CreatedAt),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "pending execution", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field (id, createdAt, state)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return markJSON(cmd)
}

func newRemoteOpsPendingDecisionCmd(approve bool) *cobra.Command {
	var yes bool
	verb, title, done := "decline", "Decline", "Declined"
	if approve {
		verb, title, done = "approve", "Approve", "Approved"
	}

	cmd := &cobra.Command{
		Use:   verb + " <pending-execution-id>",
		Short: title + " a pending execution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			action := fmt.Sprintf("%s pending execution %s", verb, args[0])
			return guard(cmd.OutOrStdout(), "remoteops pending "+verb, action, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.RemoteScriptsPendingDecision(cmd.Context(), args[0], approve); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"pendingExecutionId": args[0], "action": verb})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s pending execution %s\n", done, args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newRemoteOpsGuardrailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guardrails",
		Short: "Manage remote-script execution guardrails",
		Long: `Manage guardrails that require approval before scripts run on large numbers
of endpoints. A guardrail is configured per scope (account, site, or group).`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newRemoteOpsGuardrailsGetCmd())
	cmd.AddCommand(newRemoteOpsGuardrailsSetCmd())
	cmd.AddCommand(newRemoteOpsGuardrailsDeleteCmd())
	cmd.AddCommand(newRemoteOpsGuardrailsCheckCmd())
	return cmd
}

// guardrailScopeFromFlags validates and builds a guardrail scope from flags.
func guardrailScopeFromFlags(scopeID, scopeLevel string) (mgmt.GuardrailScope, error) {
	if scopeID == "" {
		return mgmt.GuardrailScope{}, fmt.Errorf("--scope-id is required")
	}
	if scopeLevel == "" {
		return mgmt.GuardrailScope{}, fmt.Errorf("--scope-level is required (account, site, or group)")
	}
	return mgmt.GuardrailScope{ScopeID: scopeID, ScopeLevel: mgmt.GuardrailScopeLevel(scopeLevel)}, nil
}

func newRemoteOpsGuardrailsGetCmd() *cobra.Command {
	var scopeID, scopeLevel string

	cmd := &cobra.Command{
		Use:   "get --scope-id <id> --scope-level <level>",
		Short: "Get the guardrail configuration for a scope",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := guardrailScopeFromFlags(scopeID, scopeLevel)
			if err != nil {
				return err
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			g, err := c.GuardrailsGet(cmd.Context(), scope)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), g)
			}
			quantity := "-"
			if g.EndpointsQuantity != nil {
				quantity = strconv.Itoa(*g.EndpointsQuantity)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"Enabled", boolIcon(g.Enabled)},
				{"Inherited", boolIcon(g.Inherited)},
				{"EndpointsQuantity", quantity},
				{"ScriptTypes", orDash(strings.Join(g.ScriptTypes, ", "))},
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope ID (required)")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level: account, site, or group (required)")
	return markJSON(cmd)
}

func newRemoteOpsGuardrailsSetCmd() *cobra.Command {
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "set --from-file <guardrail.json>",
		Short: "Create or update a guardrail configuration",
		Long: `Create or update (upsert) a guardrail from a JSON file, for example:

  {
    "scopeId": "000000000000000000",
    "scopeLevel": "site",
    "endpointsQuantity": 100,
    "scriptTypes": ["action"],
    "enabled": true
  }`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			var in mgmt.GuardrailsUpsertInput
			if err := readJSONFile(fromFile, &in); err != nil {
				return err
			}
			target := string(in.ScopeLevel) + "/" + in.ScopeID
			action := fmt.Sprintf("set guardrail for %s from %s", target, fromFile)
			return guard(cmd.OutOrStdout(), "remoteops guardrails set", action, target, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.GuardrailsUpsert(cmd.Context(), in); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "ok", "scopeId": in.ScopeID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Set guardrail for %s\n", target)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "JSON file with the guardrail data (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newRemoteOpsGuardrailsDeleteCmd() *cobra.Command {
	var scopeID, scopeLevel string
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete --scope-id <id> --scope-level <level>",
		Short: "Delete a guardrail configuration for a scope",
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope, err := guardrailScopeFromFlags(scopeID, scopeLevel)
			if err != nil {
				return err
			}
			target := scopeLevel + "/" + scopeID
			action := "delete guardrail for " + target
			return guard(cmd.OutOrStdout(), "remoteops guardrails delete", action, target, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.GuardrailsDelete(cmd.Context(), scope); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "scopeId": scopeID})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted guardrail for %s\n", target)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "scope ID (required)")
	cmd.Flags().StringVar(&scopeLevel, "scope-level", "", "scope level: account, site, or group (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newRemoteOpsGuardrailsCheckCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "check --from-file <check.json>",
		Short: "Check whether a guardrail would require approval for an execution",
		Long: `Read-only guardrail pre-check: report whether running a script on the given
agents would trip a guardrail and require approval. The file holds:

  {
    "scriptId": "3000000000000000001",
    "agentIds": ["4000000000000000001"]
  }`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			var in mgmt.GuardrailCheckInput
			if err := readJSONFile(fromFile, &in); err != nil {
				return err
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			res, err := c.GuardrailsCheck(cmd.Context(), in)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), res)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Requires approval: %s\n", boolIcon(res.RequiresApproval))
			return nil
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "JSON file with scriptId and agentIds (required)")
	return markJSON(cmd)
}
