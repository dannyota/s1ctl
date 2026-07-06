package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/spf13/cobra"

	"danny.vn/s1/docs/guides"
	"danny.vn/s1/internal/mcp"
)

func newMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Model Context Protocol server",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newMCPServeCmd())
	cmd.AddCommand(newMCPInstallCmd())
	return cmd
}

func newMCPServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server on stdio",
		Long: `Start a Model Context Protocol (MCP) server that exposes every s1ctl
command as an MCP tool and every docs guide as an MCP resource.

Tools are auto-generated from the command tree — adding a command
automatically creates a tool. Resources are embedded from docs/guides/.

Configure Claude Code to use this server:

  s1ctl mcp install`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			bi := resolveBuildInfo()

			tools := mcp.ToolsFromCobra(cmd.Root())
			resources := mcp.ResourcesFromFS(guides.FS, "guide")
			srv := mcp.NewServer("s1ctl", bi.Version, tools, resources)

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			return srv.Serve(ctx)
		},
	}
}

func newMCPInstallCmd() *cobra.Command {
	var scope string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Configure Claude Code to use s1ctl as an MCP server",
		Long: `Add the s1ctl MCP server to Claude Code's settings so every command is
available as a tool and every guide as a resource.

Scopes:
  project   .claude/settings.json in current directory (default)
  user      ~/.claude/settings.json (global, all projects)`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runMCPInstall(cmd, scope)
		},
	}
	cmd.Flags().StringVar(&scope, "scope", "project", "settings scope (project, user)")
	return cmd
}

func runMCPInstall(cmd *cobra.Command, scope string) error {
	var settingsPath string
	switch scope {
	case "project":
		settingsPath = filepath.Join(".claude", "settings.json")
	case "user":
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		settingsPath = filepath.Join(home, ".claude", "settings.json")
	default:
		return fmt.Errorf("invalid scope %q (valid: project, user)", scope)
	}

	bin, err := s1ctlPath()
	if err != nil {
		return err
	}

	settings := map[string]any{}
	if data, readErr := os.ReadFile(settingsPath); readErr == nil {
		_ = json.Unmarshal(data, &settings)
	}

	servers, _ := settings["mcpServers"].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}

	if _, exists := servers["s1ctl"]; exists {
		status := "unchanged"
		if outputFormat == "json" {
			return printJSON(cmd.OutOrStdout(), map[string]string{
				"path": settingsPath, "status": status,
			})
		}
		fmt.Fprintf(cmd.OutOrStdout(), "s1ctl MCP server already configured in %s\n", settingsPath)
		return nil
	}

	servers["s1ctl"] = map[string]any{
		"command": bin,
		"args":    []string{"mcp", "serve"},
	}
	settings["mcpServers"] = servers

	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o750); err != nil {
		return err
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		return err
	}

	status := "installed"
	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), map[string]string{
			"path": settingsPath, "status": status,
		})
	}
	fmt.Fprintf(cmd.OutOrStdout(), "s1ctl MCP server added to %s\n", settingsPath)
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
