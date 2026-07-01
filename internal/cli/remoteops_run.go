package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRemoteOpsRunCmd() *cobra.Command {
	var (
		agentIDs    []string
		siteIDs     []string
		groupIDs    []string
		output      string
		description string
		inputParams string
		timeout     int
		yes         bool
	)

	cmd := &cobra.Command{
		Use:   "run <script-id>",
		Short: "Execute a remote script on agents",
		Long:  "Run a remote script from the Script Library on targeted agents.\nRequires at least one targeting flag (--agent-id, --site-id, or --group-id).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scriptID := args[0]

			if len(agentIDs) == 0 && len(siteIDs) == 0 && len(groupIDs) == 0 {
				return fmt.Errorf("at least one target is required (--agent-id, --site-id, or --group-id)")
			}

			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would execute script %s. Pass --yes to apply.\n", scriptID)
				return nil
			}

			c, err := mgmtClient()
			if err != nil {
				return err
			}

			dest := mgmt.OutputDestination(output)
			if description == "" {
				description = "s1ctl remote script execution"
			}

			filter := mgmt.RemoteScriptsExecuteFilter{
				IDs:      agentIDs,
				SiteIDs:  siteIDs,
				GroupIDs: groupIDs,
			}
			data := mgmt.RemoteScriptsExecuteParams{
				ScriptID:          scriptID,
				OutputDestination: dest,
				TaskDescription:   description,
				InputParams:       inputParams,
				TimeoutSeconds:    timeout,
			}

			result, err := c.RemoteScriptsExecute(cmd.Context(), filter, data)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), result)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Executed: %s affected\n", pluralize(result.Affected, "agent"))
			if result.ParentTaskID != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Task ID: %s\n", result.ParentTaskID)
				fmt.Fprintf(cmd.OutOrStdout(), "Check status: s1ctl remoteops results %s\n", result.ParentTaskID)
			}
			if result.Pending {
				fmt.Fprintf(cmd.OutOrStdout(), "Pending approval (execution ID: %s)\n", result.PendingExecutionID)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&agentIDs, "agent-id", nil, "target agent IDs")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "target site IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "target group IDs")
	cmd.Flags().StringVar(&output, "output-dest", "SentinelCloud", "output destination (SentinelCloud, Local, None, SingularityXDR)")
	cmd.Flags().StringVar(&description, "description", "", "task description (default: s1ctl remote script execution)")
	cmd.Flags().StringVar(&inputParams, "input-params", "", "script input parameters")
	cmd.Flags().IntVar(&timeout, "timeout", 0, "script runtime timeout in seconds (60-172800)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
