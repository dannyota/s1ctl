package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeSavedQueriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "saved-queries",
		Short: "Manage saved PowerQueries",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDatalakeSavedQueriesListCmd())
	cmd.AddCommand(newDatalakeSavedQueriesDeleteCmd())
	return cmd
}

func newDatalakeSavedQueriesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
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

			if len(queries) == 0 && outputFormat == "table" {
				fmt.Fprintln(cmd.OutOrStdout(), "No saved queries found.")
				return nil
			}

			// Build full (untruncated) rows for CSV; truncated rows for table.
			headers := []string{"Name", "Type", "Index", "Query"}
			fullRows := make([][]string, len(queries))
			tableRows := make([][]string, len(queries))
			for i, q := range queries {
				fullRows[i] = []string{
					q.Name,
					q.Type,
					strconv.Itoa(q.Index),
					q.URL,
				}
				tableRows[i] = []string{
					q.Name,
					q.Type,
					strconv.Itoa(q.Index),
					truncate(q.URL, 60),
				}
			}

			switch outputFormat {
			case "json":
				return printJSON(cmd.OutOrStdout(), queries)
			case "csv":
				return printCSV(cmd.OutOrStdout(), headers, fullRows)
			default:
				printTable(headers, tableRows)
				printFooter(cmd.OutOrStdout(), len(queries), len(queries), "saved query", true)
				return nil
			}
		},
	}
	return markJSON(cmd)
}

func newDatalakeSavedQueriesDeleteCmd() *cobra.Command {
	var (
		searchType string
		index      int
		yes        bool
	)

	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a saved query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			st := sdl.SavedSearchType(searchType)
			switch st {
			case sdl.SavedSearchTypePrivate, sdl.SavedSearchTypeShared:
			default:
				return fmt.Errorf("--type must be PRIVATE or SHARED")
			}
			return guard(cmd.OutOrStdout(), "saved-queries delete", fmt.Sprintf("delete saved query %q (type=%s, index=%d)", name, searchType, index), name, yes, func() error {
				consoleURL, token, err := resolveConfig()
				if err != nil {
					return err
				}
				c := sdl.NewClient(consoleURL, token)
				return c.SavedSearchDelete(cmd.Context(), name, st, index)
			})
		},
	}
	cmd.Flags().StringVar(&searchType, "type", "PRIVATE", "saved search type (PRIVATE, SHARED)")
	cmd.Flags().IntVar(&index, "index", 0, "saved search index")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the mutation (default: dry-run)")
	return markJSON(cmd)
}
