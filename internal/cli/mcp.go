package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/spf13/cobra"

	"danny.vn/s1/docs/guides"
	"danny.vn/s1/internal/mcp"
)

func newMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Run Model Context Protocol server",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newMCPServeCmd())
	cmd.AddCommand(newMCPInstallCmd())
	return cmd
}

func newMCPServeCmd() *cobra.Command {
	var serveReadOnly bool
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server on stdio",
		Long: `Start a Model Context Protocol (MCP) server that exposes every s1ctl
command as an MCP tool and every docs guide as an MCP resource.

Tools are auto-generated from the command tree — adding a command
automatically creates a tool. Resources are embedded from docs/guides/.

Use --read-only to hide mutation tools and block mutations via run.

Configure Claude Code to use this server:

  s1ctl mcp install`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			bi := resolveBuildInfo()

			resources := mcp.ResourcesFromFS(guides.FS, "guide")
			srv := mcp.NewDynamicServer("s1ctl", bi.Version, cmd.Root(), resources,
				mcp.WithReadOnly(serveReadOnly))

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			return srv.Serve(ctx)
		},
	}
	cmd.Flags().BoolVar(&serveReadOnly, "read-only", false, "expose only read-only tools and block mutations")
	return cmd
}

func newMCPInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Register s1ctl in the project .mcp.json",
		Long: `Add s1ctl as an MCP server in the project-level .mcp.json so every
Claude Code session in this directory gets s1ctl tools automatically.
Idempotent — updates the entry if it already exists.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMCPInstall(cmd)
		},
	}
	return markJSON(cmd)
}

func runMCPInstall(cmd *cobra.Command) error {
	bin, err := s1ctlPath()
	if err != nil {
		return err
	}

	const mcpFile = ".mcp.json"

	config := map[string]any{}
	if data, readErr := os.ReadFile(mcpFile); readErr == nil {
		_ = json.Unmarshal(data, &config)
	}

	servers, _ := config["mcpServers"].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}

	servers["s1ctl"] = map[string]any{
		"command": bin,
		"args":    []string{"mcp", "serve"},
	}
	config["mcpServers"] = servers

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(mcpFile, append(data, '\n'), 0o644); err != nil {
		return err
	}

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), map[string]string{
			"path":    mcpFile,
			"binary":  bin,
			"status":  "installed",
			"command": "s1ctl mcp serve",
		})
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Registered s1ctl MCP server in %s\n", mcpFile)
	fmt.Fprintf(cmd.OutOrStdout(), "  command: %s mcp serve\n", bin)
	fmt.Fprintln(cmd.OutOrStdout(), "Restart Claude Code to pick up the new server.")
	return nil
}

func s1ctlPath() (string, error) {
	exe, err := os.Executable()
	if err == nil {
		return exe, nil
	}
	path, err := exec.LookPath("s1ctl")
	if err != nil {
		return "", fmt.Errorf("could not locate s1ctl binary: %w", err)
	}
	return path, nil
}
