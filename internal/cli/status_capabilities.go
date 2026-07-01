package cli

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

func newStatusCapabilitiesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "capabilities",
		Short: "Show s1ctl version, config, and API reachability",
		Args:  cobra.NoArgs,
		RunE:  runStatusCapabilities,
	}
}

type capabilitiesOutput struct {
	Version    string        `json:"version"`
	Commit     string        `json:"commit"`
	Go         string        `json:"go"`
	ConsoleURL string        `json:"console_url"`
	Token      string        `json:"token"`
	SDLURL     string        `json:"sdl_url"`
	ConfigFile string        `json:"config_file"`
	APIs       []checkResult `json:"apis"`
}

func runStatusCapabilities(cmd *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	apis := []checkResult{
		checkMGMT(ctx, cfg.ConsoleURL, cfg.Token),
		checkGraphQL(ctx, cfg.ConsoleURL, cfg.Token),
		checkSDL(ctx, cfg.SDLURL, cfg.Token),
	}

	out := capabilitiesOutput{
		Version:    version,
		Commit:     commit,
		Go:         runtime.Version(),
		ConsoleURL: cfg.ConsoleURL,
		Token:      redactToken(cfg.Token),
		SDLURL:     cfg.SDLURL,
		ConfigFile: cfg.Source(),
		APIs:       apis,
	}

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), out)
	}

	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "s1ctl %s (%s) %s\n\n", version, commit, runtime.Version())
	fmt.Fprintf(w, "Console:  %s\n", orDash(cfg.ConsoleURL))
	fmt.Fprintf(w, "Token:    %s\n", redactToken(cfg.Token))
	sdlDisplay := cfg.SDLURL
	if sdlDisplay == "" {
		sdlDisplay = "not configured"
	}
	fmt.Fprintf(w, "SDL URL:  %s\n", sdlDisplay)
	fmt.Fprintf(w, "Config:   %s\n\n", cfg.Source())

	fmt.Fprintln(w, "API Connectivity")
	for _, a := range apis {
		status := "ok"
		if !a.OK {
			status = "FAIL"
		}
		line := fmt.Sprintf("  %-12s %s (%s)", a.Surface, status, a.Latency)
		if a.Error != "" {
			line += "  " + a.Error
		}
		fmt.Fprintln(w, line)
	}
	return nil
}
