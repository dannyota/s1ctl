package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datalake",
		Short: "Query Singularity Data Lake (SDL)",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDatalakePowerQueryCmd())
	return cmd
}

func newDatalakePowerQueryCmd() *cobra.Command {
	var query, startTime, endTime, priority, protocol string

	cmd := &cobra.Command{
		Use:   "powerquery",
		Short: "Execute a PowerQuery",
		Long: `Execute a PowerQuery against the Singularity Data Lake.

By default, uses the GraphQL protocol which connects through the management
console and does not require a separate SDL URL. Use --protocol rest to use
the REST API, which requires S1_SDL_URL to be configured.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if query == "" {
				return fmt.Errorf("--query is required")
			}
			switch protocol {
			case "graphql":
				return runPowerQueryGraphQL(cmd, query, startTime, endTime)
			case "rest":
				return runPowerQueryREST(cmd, query, startTime, endTime, priority)
			default:
				return fmt.Errorf("unsupported protocol: %s (use graphql or rest)", protocol)
			}
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "PowerQuery expression (required)")
	cmd.Flags().StringVar(&startTime, "start", "24h", "start time (e.g. 24h, 7d)")
	cmd.Flags().StringVar(&endTime, "end", "", "end time")
	cmd.Flags().StringVar(&priority, "priority", "low", "query priority (low, high) [REST only]")
	cmd.Flags().StringVar(&protocol, "protocol", "graphql", "API protocol (graphql, rest)")
	return cmd
}

func runPowerQueryGraphQL(cmd *cobra.Command, query, startTime, endTime string) error {
	consoleURL, token, err := resolveConfig()
	if err != nil {
		return err
	}
	c := sdl.NewClient(consoleURL, token)
	req := &sdl.PowerQueryRequest{
		Query:     query,
		StartTime: startTime,
		EndTime:   endTime,
	}
	var resp *sdl.PowerQueryResponse
	err = runWithSpinner("Running query...", func() error {
		var queryErr error
		resp, queryErr = c.PowerQueryGraphQL(cmd.Context(), req)
		return queryErr
	})
	if err != nil {
		return err
	}
	return printPowerQueryResult(cmd, resp)
}

func runPowerQueryREST(cmd *cobra.Command, query, startTime, endTime, priority string) error {
	c, err := sdlClient()
	if err != nil {
		return err
	}
	req := &sdl.PowerQueryRequest{
		Query:     query,
		StartTime: startTime,
		EndTime:   endTime,
		Priority:  priority,
	}
	var resp *sdl.PowerQueryResponse
	err = runWithSpinner("Running query...", func() error {
		var queryErr error
		resp, queryErr = c.PowerQuery(cmd.Context(), req)
		return queryErr
	})
	if err != nil {
		return err
	}
	return printPowerQueryResult(cmd, resp)
}

func printPowerQueryResult(cmd *cobra.Command, resp *sdl.PowerQueryResponse) error {
	if len(resp.Columns) == 0 {
		if outputFormat == "json" {
			return printJSON(resp)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "No results.")
		return nil
	}
	headers := make([]string, len(resp.Columns))
	for i, col := range resp.Columns {
		headers[i] = col.Name
	}
	rows := make([][]string, len(resp.Values))
	for i, row := range resp.Values {
		cells := make([]string, len(row))
		for j, v := range row {
			cells[j] = truncate(fmt.Sprint(v), 60)
		}
		rows[i] = cells
	}
	return printOutput(cmd.OutOrStdout(), headers, rows, resp, len(resp.Values), len(resp.Values), "row", true)
}

func sdlClient() (*sdl.Client, error) {
	sdlURL, token, err := resolveSDLURL()
	if err != nil {
		return nil, err
	}
	return sdl.NewClient(sdlURL, token), nil
}
