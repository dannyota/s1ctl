package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAlertsResolveCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "resolve <id> [id...]",
		Short: "Resolve alerts (bulk)",
		Long: `Set status to "RESOLVED" on one or more alerts.

Specify one or more alert IDs. Dry-run by default.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would resolve %s. Pass --yes to apply.\n",
					pluralize(len(args), "alert"))
				return nil
			}
			c, err := gqlClient()
			if err != nil {
				return err
			}
			if err := c.AlertsUpdateStatus(cmd.Context(), args, "RESOLVED"); err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]any{
					"status":   "resolved",
					"affected": len(args),
				})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "resolve: %s affected\n", pluralize(len(args), "alert"))
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
