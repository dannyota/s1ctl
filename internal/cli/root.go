package cli

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

var (
	outputFormat string
	verbose      bool
	noProgress   bool
	configFile   string
)

const defaultPageSize = 50

func newRootCmd() *cobra.Command {
	var jsonFlag bool

	cmd := &cobra.Command{
		Use:           "s1ctl",
		Short:         "CLI for SentinelOne Singularity Platform",
		Long:          "Operate SentinelOne Singularity Platform as code — pull, diff, push.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if jsonFlag {
				outputFormat = "json"
			}
			validFormats := []string{"table", "json", "csv"}
			if !slices.Contains(validFormats, outputFormat) {
				return fmt.Errorf("invalid output format %q (valid: table, json, csv)", outputFormat)
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVar(&outputFormat, "output", "table", "output format (table, json, csv)")
	cmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "shorthand for --output json")
	cmd.MarkFlagsMutuallyExclusive("output", "json")
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "show detailed error information")
	cmd.PersistentFlags().BoolVar(&noProgress, "no-progress", false, "disable spinners and progress output")
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
		printError(root.ErrOrStderr(), err)
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
