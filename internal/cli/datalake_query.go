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
	var query, startTime, endTime, priority string

	cmd := &cobra.Command{
		Use:   "powerquery",
		Short: "Execute a PowerQuery",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if query == "" {
				return fmt.Errorf("--query is required")
			}
			c, err := sdlClient()
			if err != nil {
				return err
			}
			resp, err := c.PowerQuery(cmd.Context(), &sdl.PowerQueryRequest{
				Query:     query,
				StartTime: startTime,
				EndTime:   endTime,
				Priority:  priority,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(resp)
			}
			if len(resp.Columns) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No results.")
				return nil
			}
			headers := make([]string, len(resp.Columns))
			for i, col := range resp.Columns {
				headers[i] = col.Name
			}
			var rows [][]string
			for _, row := range resp.Values {
				cells := make([]string, len(row))
				for i, v := range row {
					cells[i] = truncate(fmt.Sprint(v), 60)
				}
				rows = append(rows, cells)
			}
			printTable(headers, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(len(resp.Values), "row"))
			return nil
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "PowerQuery expression (required)")
	cmd.Flags().StringVar(&startTime, "start", "24h", "start time (e.g. 24h, 7d)")
	cmd.Flags().StringVar(&endTime, "end", "", "end time")
	cmd.Flags().StringVar(&priority, "priority", "low", "query priority (low, high)")
	return cmd
}

func sdlClient() (*sdl.Client, error) {
	consoleURL, token, err := resolveConfig()
	if err != nil {
		return nil, err
	}
	return sdl.NewClient(consoleURL, token), nil
}
