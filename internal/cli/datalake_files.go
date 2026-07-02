package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Manage data lake configuration files",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDatalakeFilesListCmd())
	cmd.AddCommand(newDatalakeFilesGetCmd())
	cmd.AddCommand(newDatalakeFilesPutCmd())
	return cmd
}

func newDatalakeFilesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configuration files",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := sdlClient()
			if err != nil {
				return err
			}
			resp, err := c.ListFiles(cmd.Context())
			if err != nil {
				return err
			}
			rows := make([][]string, 0, len(resp.Paths))
			for _, p := range resp.Paths {
				rows = append(rows, []string{p})
			}
			return printOutput(cmd.OutOrStdout(), []string{"PATH"}, rows, resp, len(rows), len(rows), "file", true)
		},
	}
}

func newDatalakeFilesGetCmd() *cobra.Command {
	var outFile string

	cmd := &cobra.Command{
		Use:   "get <path>",
		Short: "Fetch a configuration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := sdlClient()
			if err != nil {
				return err
			}
			resp, err := c.GetFile(cmd.Context(), &sdl.GetFileRequest{Path: args[0]})
			if err != nil {
				return err
			}
			if resp.Status != "success" && resp.Status != "success/unchanged" {
				return fmt.Errorf("get %s: %s", args[0], resp.Status)
			}
			if outFile != "" {
				if err := os.WriteFile(outFile, []byte(resp.Content), 0o644); err != nil {
					return err
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Wrote %s (version %d) to %s\n", args[0], resp.Version, outFile)
				return nil
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), resp)
			}
			fmt.Fprint(cmd.OutOrStdout(), resp.Content)
			return nil
		},
	}
	cmd.Flags().StringVar(&outFile, "out", "", "write content to a local file instead of stdout")
	return cmd
}

func newDatalakeFilesPutCmd() *cobra.Command {
	var fromFile string
	var deleteFile, yes bool
	var expectedVersion int64

	cmd := &cobra.Command{
		Use:   "put <path> (--from-file <local> | --delete)",
		Short: "Create, update, or delete a configuration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if (fromFile != "") == deleteFile {
				return fmt.Errorf("pass exactly one of --from-file or --delete")
			}
			var content string
			action := "delete data lake file " + args[0]
			if fromFile != "" {
				raw, err := os.ReadFile(fromFile)
				if err != nil {
					return fmt.Errorf("read %s: %w", fromFile, err)
				}
				content = string(raw)
				action = fmt.Sprintf("write %s (%d bytes) to data lake file %s", fromFile, len(raw), args[0])
			}
			return guard(cmd.OutOrStdout(), "datalake files put", action, args[0], yes, func() error {
				c, err := sdlClient()
				if err != nil {
					return err
				}
				resp, err := c.PutFile(cmd.Context(), &sdl.PutFileRequest{
					Path:            args[0],
					Content:         content,
					DeleteFile:      deleteFile,
					ExpectedVersion: expectedVersion,
				})
				if err != nil {
					return err
				}
				if resp.Status != "success" {
					return fmt.Errorf("put %s: %s", args[0], resp.Status)
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), resp)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "put %s: %s\n", args[0], resp.Status)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "local file with the new content")
	cmd.Flags().BoolVar(&deleteFile, "delete", false, "delete the remote file")
	cmd.Flags().Int64Var(&expectedVersion, "expected-version", 0, "fail if the remote version differs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
