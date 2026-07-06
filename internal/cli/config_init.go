package cli

import (
	"fmt"
	"os"

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
		Short: "Configure s1ctl interactively",
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
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show resolved configuration",
		Args:  cobra.NoArgs,
		RunE:  runConfigShow,
	}
	return markJSON(cmd)
}

func runConfigShow(cmd *cobra.Command, _ []string) error {
	// Load without validation so show works as a diagnostic even with
	// partial configuration (e.g. missing token).
	inst, err := config.Load(configFile)
	if err != nil {
		return err
	}

	token := redactToken(inst.Token)
	sdlDisplay := inst.SDLURL
	if sdlDisplay == "" {
		sdlDisplay = "not configured"
	}
	consoleDisplay := inst.ConsoleURL
	if consoleDisplay == "" {
		consoleDisplay = "not configured"
	}

	type envVar struct {
		name string
		set  bool
	}
	envVars := []envVar{
		{"S1_CONSOLE_URL", os.Getenv("S1_CONSOLE_URL") != ""},
		{"S1_TOKEN", os.Getenv("S1_TOKEN") != ""},
		{"S1_SDL_URL", os.Getenv("S1_SDL_URL") != ""},
	}

	if outputFormat == "json" {
		envMap := make(map[string]bool, len(envVars))
		for _, v := range envVars {
			envMap[v.name] = v.set
		}
		out := map[string]any{
			"console_url": inst.ConsoleURL,
			"token":       token,
			"sdl_url":     inst.SDLURL,
			"config_file": inst.Source(),
			"env_vars":    envMap,
		}
		return printJSON(cmd.OutOrStdout(), out)
	}

	rows := [][]string{
		{"Console URL", consoleDisplay},
		{"API Token", token},
		{"SDL URL", sdlDisplay},
		{"Config File", inst.Source()},
	}
	for _, v := range envVars {
		status := "not set"
		if v.set {
			status = "set"
		}
		rows = append(rows, []string{v.name, status})
	}

	printTable([]string{"Setting", "Value"}, rows)
	return nil
}

// redactToken masks all but the last 4 characters of a token.
func redactToken(token string) string {
	if token == "" {
		return "not configured"
	}
	if len(token) <= 4 {
		return "****"
	}
	return "****" + token[len(token)-4:]
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
