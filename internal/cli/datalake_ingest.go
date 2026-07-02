package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeIngestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest events or raw logs into the data lake",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDatalakeIngestEventsCmd())
	cmd.AddCommand(newDatalakeIngestLogsCmd())
	return cmd
}

func newDatalakeIngestEventsCmd() *cobra.Command {
	var file, session string
	var yes bool

	cmd := &cobra.Command{
		Use:   "events --file <events.json> --session <id>",
		Short: "Ingest structured events (addEvents)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if file == "" {
				return fmt.Errorf("--file is required")
			}
			if session == "" {
				return fmt.Errorf("--session is required")
			}
			raw, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read %s: %w", file, err)
			}
			var events []sdl.Event
			if err := json.Unmarshal(raw, &events); err != nil {
				return fmt.Errorf("parse %s: %w", file, err)
			}
			action := fmt.Sprintf("ingest %s from %s (session %s)", pluralize(len(events), "event"), file, session)
			return guard(cmd.OutOrStdout(), "datalake ingest events", action, file, yes, func() error {
				c, err := sdlClient()
				if err != nil {
					return err
				}
				resp, err := c.AddEvents(cmd.Context(), &sdl.AddEventsRequest{Session: session, Events: events})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), resp)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Ingested %s: %s\n", pluralize(len(events), "event"), resp.Status)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "JSON file containing an array of events (required)")
	cmd.Flags().StringVar(&session, "session", "", "unique session ID for this uploader (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newDatalakeIngestLogsCmd() *cobra.Command {
	var file, parser, serverHost, logfile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "logs --file <app.log> --parser <name>",
		Short: "Ingest a plain-text log file (uploadLogs)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if file == "" {
				return fmt.Errorf("--file is required")
			}
			raw, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read %s: %w", file, err)
			}
			action := fmt.Sprintf("upload log file %s (%d bytes)", file, len(raw))
			return guard(cmd.OutOrStdout(), "datalake ingest logs", action, file, yes, func() error {
				c, err := sdlClient()
				if err != nil {
					return err
				}
				resp, err := c.UploadLogs(cmd.Context(), &sdl.UploadLogsRequest{
					Parser:     parser,
					ServerHost: serverHost,
					Logfile:    logfile,
					Body:       string(raw),
				})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), resp)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Uploaded %s: %s\n", file, resp.Status)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "log file to upload (required)")
	cmd.Flags().StringVar(&parser, "parser", "", "parser to apply on ingest")
	cmd.Flags().StringVar(&serverHost, "server-host", "", "serverHost attribute for the events")
	cmd.Flags().StringVar(&logfile, "logfile", "", "logfile attribute for the events")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
