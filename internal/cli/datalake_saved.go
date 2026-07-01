package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeSavedQueriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "saved-queries",
		Short: "List saved PowerQueries",
		Long: `List saved searches from the Singularity Data Lake console.
Shows both private and shared saved queries.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			consoleURL, token, err := resolveConfig()
			if err != nil {
				return err
			}
			c := sdl.NewClient(consoleURL, token)

			var queries []sdl.SavedSearch
			err = runWithSpinner("Fetching saved queries...", func() error {
				var qErr error
				queries, qErr = c.SavedSearches(cmd.Context())
				return qErr
			})
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), queries)
			}
			if len(queries) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No saved queries found.")
				return nil
			}

			headers := []string{"Name", "Type", "Query"}
			rows := make([][]string, len(queries))
			for i, q := range queries {
				rows[i] = []string{
					q.Name,
					q.Type,
					truncate(q.URL, 60),
				}
			}
			printTable(headers, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(len(queries), "saved query"))
			return nil
		},
	}
	return cmd
}
