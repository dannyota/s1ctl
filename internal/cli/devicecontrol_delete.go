package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newDeviceControlDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <rule-id>...",
		Short: "Delete device control rules",
		Long: `Delete one or more device control rules by ID.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "devicecontrol delete",
				"delete "+pluralize(len(args), "device rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.DeviceRulesDelete(cmd.Context(), args)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s\n", pluralize(affected, "device rule"))
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return markJSON(cmd)
}
