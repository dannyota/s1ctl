package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeDashboardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboards",
		Short: "Manage Data Lake dashboards",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDatalakeDashboardsListCmd())
	cmd.AddCommand(newDatalakeDashboardsGetCmd())
	return cmd
}

func newDatalakeDashboardsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List dashboards",
		RunE: func(cmd *cobra.Command, _ []string) error {
			consoleURL, token, err := resolveConfig()
			if err != nil {
				return err
			}
			c := sdl.NewClient(consoleURL, token)

			var dashboards []sdl.Dashboard
			err = runWithSpinner("Fetching dashboards...", func() error {
				var dErr error
				dashboards, dErr = c.DashboardsList(cmd.Context())
				return dErr
			})
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), dashboards)
			}
			if len(dashboards) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No dashboards found.")
				return nil
			}

			headers := []string{"ID", "Name", "Built-in", "Editable"}
			rows := make([][]string, len(dashboards))
			for i, d := range dashboards {
				rows[i] = []string{
					d.ID,
					d.Name,
					fmt.Sprint(d.IsBuiltIn),
					fmt.Sprint(d.IsEditable),
				}
			}
			printTable(headers, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(len(dashboards), "dashboard"))
			return nil
		},
	}
	return markJSON(cmd)
}

func newDatalakeDashboardsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get dashboard details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			consoleURL, token, err := resolveConfig()
			if err != nil {
				return err
			}
			c := sdl.NewClient(consoleURL, token)

			var d *sdl.DashboardDetail
			err = runWithSpinner("Fetching dashboard...", func() error {
				var dErr error
				d, dErr = c.DashboardGet(cmd.Context(), args[0])
				return dErr
			})
			if err != nil {
				return err
			}

			return printJSON(cmd.OutOrStdout(), d)
		},
	}
	return markJSON(cmd)
}
