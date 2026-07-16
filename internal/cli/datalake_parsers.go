package cli

import (
	"fmt"

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
