package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeParsersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parsers",
		Short: "Manage Data Lake parsers (configuration files)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDatalakeParsersListCmd())
	cmd.AddCommand(newDatalakeParsersGetCmd())
	cmd.AddCommand(newDatalakeParsersCreateCmd())
	cmd.AddCommand(newDatalakeParsersDeleteCmd())
	return cmd
}

func newDatalakeParsersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List parsers",
		RunE: func(cmd *cobra.Command, _ []string) error {
			consoleURL, token, err := resolveConfig()
			if err != nil {
				return err
			}
			c := sdl.NewClient(consoleURL, token)

			var parsers []sdl.Parser
			err = runWithSpinner("Fetching parsers...", func() error {
				var pErr error
				parsers, pErr = c.ParsersList(cmd.Context())
				return pErr
			})
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), parsers)
			}
			if len(parsers) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No parsers found.")
				return nil
			}

			headers := []string{"UDO ID", "Name", "Read-only", "Version"}
			rows := make([][]string, len(parsers))
			for i, p := range parsers {
				rows[i] = []string{
					p.UdoID,
					p.Name,
					fmt.Sprint(p.ReadOnly),
					fmt.Sprint(p.Version),
				}
			}
			printTable(headers, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(len(parsers), "parser"))
			return nil
		},
	}
	return markJSON(cmd)
}

func newDatalakeParsersGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <udo-id>",
		Short: "Get parser details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			consoleURL, token, err := resolveConfig()
			if err != nil {
				return err
			}
			c := sdl.NewClient(consoleURL, token)

			var p *sdl.ParserDetail
			err = runWithSpinner("Fetching parser...", func() error {
				var pErr error
				p, pErr = c.ParserGet(cmd.Context(), args[0])
				return pErr
			})
			if err != nil {
				return err
			}

			return printJSON(cmd.OutOrStdout(), p)
		},
	}
	return markJSON(cmd)
}

func newDatalakeParsersCreateCmd() *cobra.Command {
	var (
		name     string
		udoID    string
		fromFile string
		yes      bool
	)

	cmd := &cobra.Command{
		Use:   "create --from-file <path> --name <name>",
		Short: "Create or update a parser from a file",
		Long: `Create or update a Data Lake parser (configuration file) from a local file.

The parser content is read from --from-file. If the parser already exists
(matched by --udo-id), it is updated; otherwise a new parser is created.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			content, err := os.ReadFile(fromFile)
			if err != nil {
				return fmt.Errorf("read parser file: %w", err)
			}
			contentStr := string(content)
			input := &sdl.ParserCreateInput{
				Name:    &name,
				Content: &contentStr,
			}
			if udoID != "" {
				input.UdoID = &udoID
			}
			return guard(cmd.OutOrStdout(), "parsers create", fmt.Sprintf("create parser %q from %s", name, fromFile), name, yes, func() error {
				consoleURL, token, cErr := resolveConfig()
				if cErr != nil {
					return cErr
				}
				c := sdl.NewClient(consoleURL, token)
				result, cErr := c.ParserCreate(cmd.Context(), input)
				if cErr != nil {
					return cErr
				}
				return printJSON(cmd.OutOrStdout(), result)
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "parser name (required)")
	cmd.Flags().StringVar(&udoID, "udo-id", "", "UDO ID (update existing parser)")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "path to parser content file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the mutation (default: dry-run)")
	return markJSON(cmd)
}

func newDatalakeParsersDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <udo-id>",
		Short: "Delete a parser",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			udoID := args[0]
			return guard(cmd.OutOrStdout(), "parsers delete", fmt.Sprintf("delete parser %q", udoID), udoID, yes, func() error {
				consoleURL, token, err := resolveConfig()
				if err != nil {
					return err
				}
				c := sdl.NewClient(consoleURL, token)
				return c.ParserDelete(cmd.Context(), udoID, nil)
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the mutation (default: dry-run)")
	return markJSON(cmd)
}
