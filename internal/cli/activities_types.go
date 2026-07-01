package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newActivitiesTypesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "List available activity types",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			types, err := c.ActivitiesTypes(cmd.Context())
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), types)
			}
			headers := []string{"ID", "Description"}
			rows := make([][]string, len(types))
			for i, t := range types {
				rows[i] = []string{strconv.Itoa(t.ID), truncate(t.Description, 80)}
			}
			printTable(headers, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(len(types), "activity type"))
			return nil
		},
	}
}
