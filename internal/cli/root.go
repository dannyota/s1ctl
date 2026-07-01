package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	configFile string
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "s1ctl",
		Short:         "CLI for SentinelOne Singularity Platform",
		Long:          "Operate SentinelOne Singularity Platform as code — pull, diff, push.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output as JSON")
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default ~/.s1ctl/config.yaml)")
	return cmd
}

func requireSubcommand(cmd *cobra.Command) {
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	}
}

// Execute runs the root command and returns an exit code.
func Execute() int {
	root := newRootCmd()
	registerCommands(root)
	if err := root.Execute(); err != nil {
		root.PrintErrln("Error:", err)
		return 1
	}
	return 0
}

func resolveConfig() (consoleURL, token string, err error) {
	cfg, loadErr := loadConfig()
	if loadErr != nil {
		return "", "", loadErr
	}
	return cfg.ConsoleURL, cfg.Token, nil
}

func resolveSDLURL() (sdlURL, token string, err error) {
	cfg, loadErr := loadConfig()
	if loadErr != nil {
		return "", "", loadErr
	}
	if cfg.SDLURL == "" {
		return "", "", fmt.Errorf("SDL URL is required (set S1_SDL_URL or sdl_url in config)\nThe SDL console is separate from the management console (e.g. https://xdr.us1.sentinelone.net)")
	}
	return cfg.SDLURL, cfg.Token, nil
}
