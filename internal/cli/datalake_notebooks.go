package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeNotebooksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notebooks",
		Short: "Manage Purple AI notebooks",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDatalakeNotebooksListCmd())
	cmd.AddCommand(newDatalakeNotebooksGetCmd())
	cmd.AddCommand(newDatalakeNotebooksCreateCmd())
	cmd.AddCommand(newDatalakeNotebooksUpdateCmd())
	cmd.AddCommand(newDatalakeNotebooksDeleteCmd())
	return cmd
}

func newDatalakeNotebooksListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List notebooks",
		RunE: func(cmd *cobra.Command, _ []string) error {
			consoleURL, token, err := resolveConfig()
			if err != nil {
				return err
			}
			c := sdl.NewClient(consoleURL, token)

			var notebooks []sdl.Notebook
			err = runWithSpinner("Fetching notebooks...", func() error {
				var nErr error
				notebooks, nErr = c.NotebooksList(cmd.Context())
				return nErr
			})
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), notebooks)
			}
			if len(notebooks) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No notebooks found.")
				return nil
			}

			headers := []string{"ID", "Name", "Source", "Shared", "Read-only"}
			rows := make([][]string, len(notebooks))
			for i, n := range notebooks {
				rows[i] = []string{
					n.ID,
					truncate(n.Name, 40),
					n.NotebookSource,
					fmt.Sprint(n.IsShared),
					fmt.Sprint(n.IsReadOnly),
				}
			}
			printTable(headers, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(len(notebooks), "notebook"))
			return nil
		},
	}
	return markJSON(cmd)
}

func newDatalakeNotebooksGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get notebook details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			consoleURL, token, err := resolveConfig()
			if err != nil {
				return err
			}
			c := sdl.NewClient(consoleURL, token)

			var n *sdl.NotebookDetail
			err = runWithSpinner("Fetching notebook...", func() error {
				var nErr error
				n, nErr = c.NotebookGet(cmd.Context(), args[0])
				return nErr
			})
			if err != nil {
				return err
			}

			return printJSON(cmd.OutOrStdout(), n)
		},
	}
	return markJSON(cmd)
}

func newDatalakeNotebooksCreateCmd() *cobra.Command {
	var (
		name        string
		description string
		yes         bool
	)

	cmd := &cobra.Command{
		Use:   "create --name <name>",
		Short: "Create a notebook",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			return guard(cmd.OutOrStdout(), "notebooks create", fmt.Sprintf("create notebook %q", name), name, yes, func() error {
				consoleURL, token, err := resolveConfig()
				if err != nil {
					return err
				}
				c := sdl.NewClient(consoleURL, token)
				result, err := c.NotebookCreate(cmd.Context(), name, description)
				if err != nil {
					return err
				}
				return printJSON(cmd.OutOrStdout(), result)
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "notebook name (required)")
	cmd.Flags().StringVar(&description, "description", "", "notebook description")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the mutation (default: dry-run)")
	return markJSON(cmd)
}

func newDatalakeNotebooksUpdateCmd() *cobra.Command {
	var (
		name        string
		description string
		yes         bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a notebook's name or description",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			input := &sdl.NotebookUpdateInput{}
			if cmd.Flags().Changed("name") {
				input.Name = &name
			}
			if cmd.Flags().Changed("description") {
				input.Description = &description
			}
			if input.Name == nil && input.Description == nil {
				return fmt.Errorf("at least one of --name or --description is required")
			}
			return guard(cmd.OutOrStdout(), "notebooks update", fmt.Sprintf("update notebook %q", id), id, yes, func() error {
				consoleURL, token, err := resolveConfig()
				if err != nil {
					return err
				}
				c := sdl.NewClient(consoleURL, token)
				return c.NotebookUpdate(cmd.Context(), id, input)
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "new notebook name")
	cmd.Flags().StringVar(&description, "description", "", "new notebook description")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the mutation (default: dry-run)")
	return markJSON(cmd)
}

func newDatalakeNotebooksDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a notebook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return guard(cmd.OutOrStdout(), "notebooks delete", fmt.Sprintf("delete notebook %q", id), id, yes, func() error {
				consoleURL, token, err := resolveConfig()
				if err != nil {
					return err
				}
				c := sdl.NewClient(consoleURL, token)
				return c.NotebookDelete(cmd.Context(), id)
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the mutation (default: dry-run)")
	return markJSON(cmd)
}
