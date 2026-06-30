package cli

import (
	"os"

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
	consoleURL = os.Getenv("S1_CONSOLE_URL")
	token = os.Getenv("S1_TOKEN")
	if consoleURL != "" && token != "" {
		return consoleURL, token, nil
	}
	cfg, loadErr := loadConfig()
	if loadErr != nil {
		return "", "", loadErr
	}
	if consoleURL == "" {
		consoleURL = cfg.ConsoleURL
	}
	if token == "" {
		token = cfg.Token
	}
	return consoleURL, token, nil
}
