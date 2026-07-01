package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"danny.vn/s1/config"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage s1ctl configuration",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newConfigInitCmd())
	cmd.AddCommand(newConfigShowCmd())
	return cmd
}

func newConfigInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Interactive configuration wizard",
		Args:  cobra.NoArgs,
		RunE:  runConfigInit,
	}
}

func runConfigInit(cmd *cobra.Command, _ []string) error {
	path := configFile
	if path == "" {
		path = config.DefaultPath()
	}
	existing := config.ReadForEdit(path)

	var consoleURL, token, sdlURL string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Console URL").
				Description("e.g. https://your-console.sentinelone.net").
				Value(&consoleURL).
				Placeholder(existing.ConsoleURL),
			huh.NewInput().
				Title("API Token").
				Description("From Settings > Users > API Token").
				Value(&token).
				EchoMode(huh.EchoModePassword).
				Placeholder("(unchanged)"),
			huh.NewInput().
				Title("SDL URL (optional)").
				Description("Data Lake console, e.g. https://xdr.us1.sentinelone.net").
				Value(&sdlURL).
				Placeholder(existing.SDLURL),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}

	if consoleURL != "" {
		existing.ConsoleURL = consoleURL
	}
	if token != "" {
		existing.Token = token
	}
	if sdlURL != "" {
		existing.SDLURL = sdlURL
	}

	if err := existing.Validate(); err != nil {
		return err
	}
	if err := config.Save(path, existing); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Config saved to %s\n", path)
	return nil
}

func newConfigShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show resolved configuration",
		Args:  cobra.NoArgs,
		RunE:  runConfigShow,
	}
}

func runConfigShow(cmd *cobra.Command, _ []string) error {
	inst, err := loadConfig()
	if err != nil {
		return err
	}

	if jsonOutput {
		out := map[string]string{
			"console_url": inst.ConsoleURL,
			"token":       "(redacted)",
			"source":      inst.Source(),
		}
		if inst.SDLURL != "" {
			out["sdl_url"] = inst.SDLURL
		}
		return printJSON(out)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Console URL: %s\n", inst.ConsoleURL)
	fmt.Fprintf(cmd.OutOrStdout(), "Token:       %s\n", "(redacted)")
	if inst.SDLURL != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "SDL URL:     %s\n", inst.SDLURL)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Source:      %s\n", inst.Source())
	return nil
}

func loadConfig() (*config.Instance, error) {
	inst, err := config.Load(configFile)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if err := inst.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w\nRun 's1ctl config init' to configure", err)
	}
	return inst, nil
}
