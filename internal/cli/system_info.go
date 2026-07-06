package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSystemCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "Show console system information",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newSystemInfoCmd())
	return cmd
}

func newSystemInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show console version, build, and health status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			ctx := cmd.Context()

			info, err := c.SystemInfo(ctx)
			if err != nil {
				return err
			}
			status, err := c.SystemStatus(ctx)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]any{
					"version":            info.Version,
					"release":            info.Release,
					"build":              info.Build,
					"patch":              info.Patch,
					"latestAgentVersion": info.LatestAgentVersion,
					"health":             status.Health,
				})
			}

			w := cmd.OutOrStdout()
			rows := [][]string{
				{"Version", orDash(info.Version)},
				{"Release", orDash(info.Release)},
				{"Build", orDash(info.Build)},
				{"Patch", orDash(info.Patch)},
				{"Latest Agent", orDash(info.LatestAgentVersion)},
				{"Health", orDash(status.Health)},
			}
			printTable([]string{"Field", "Value"}, rows)
			fmt.Fprintln(w)
			return nil
		},
	}
}
