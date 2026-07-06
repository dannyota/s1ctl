package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeNumericCmd() *cobra.Command {
	var (
		filter, function, start, end, priority string
		buckets                                int
	)

	cmd := &cobra.Command{
		Use:   "numeric",
		Short: "Run a numeric aggregation query (SDL REST)",
		Long: `Run a numeric aggregation query against the Singularity Data Lake.

Counts events, computes event rate, or applies an aggregation function
(e.g. mean, min, max) to a numeric field across one or more time buckets.

Note: numericQuery is effectively deprecated in favour of timeseries with
createSummaries=false, but remains useful for sub-30-second bucket
granularity and users with limited query permissions.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if start == "" {
				return fmt.Errorf("--start is required")
			}
			c, err := sdlClient()
			if err != nil {
				return err
			}
			resp, err := c.NumericQuery(cmd.Context(), &sdl.NumericQueryRequest{
				Filter:    filter,
				Function:  function,
				StartTime: start,
				EndTime:   end,
				Buckets:   buckets,
				Priority:  priority,
			})
			if err != nil {
				return err
			}
			var rows [][]string
			for i, v := range resp.Values {
				val := "-"
				if v != nil {
					val = strconv.FormatFloat(*v, 'f', -1, 64)
				}
				rows = append(rows, []string{strconv.Itoa(i), val})
			}
			return printOutput(cmd.OutOrStdout(), []string{"BUCKET", "VALUE"}, rows, resp, len(rows), len(rows), "bucket", true)
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "query filter expression")
	cmd.Flags().StringVar(&function, "function", "", "aggregation function (e.g. rate, count, mean(field))")
	cmd.Flags().StringVar(&start, "start", "", "start time, e.g. 1h or timestamp (required)")
	cmd.Flags().StringVar(&end, "end", "", "end time")
	cmd.Flags().IntVar(&buckets, "buckets", 0, "number of buckets (1-5000)")
	cmd.Flags().StringVar(&priority, "priority", "", "query priority (low, high)")
	return markJSON(cmd)
}
