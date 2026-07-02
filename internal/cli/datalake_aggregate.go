package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/sdl"
)

func newDatalakeFacetCmd() *cobra.Command {
	var (
		filter, field, start, end string
		maxCount                  int
	)

	cmd := &cobra.Command{
		Use:   "facet",
		Short: "Aggregate the most common values of a field (SDL REST)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if field == "" {
				return fmt.Errorf("--field is required")
			}
			if start == "" {
				return fmt.Errorf("--start is required")
			}
			c, err := sdlClient()
			if err != nil {
				return err
			}
			resp, err := c.FacetQuery(cmd.Context(), &sdl.FacetQueryRequest{
				Filter:    filter,
				Field:     field,
				StartTime: start,
				EndTime:   end,
				MaxCount:  maxCount,
			})
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), resp)
			}
			rows := make([][]string, 0, len(resp.Values))
			for _, v := range resp.Values {
				rows = append(rows, []string{v.Value, strconv.FormatInt(v.Count, 10)})
			}
			printTable([]string{"VALUE", "COUNT"}, rows)
			return nil
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "query filter expression")
	cmd.Flags().StringVar(&field, "field", "", "field to aggregate (required)")
	cmd.Flags().StringVar(&start, "start", "", "start time, e.g. 24h or timestamp (required)")
	cmd.Flags().StringVar(&end, "end", "", "end time")
	cmd.Flags().IntVar(&maxCount, "max-count", 0, "max distinct values to return")
	return cmd
}

func newDatalakeTimeseriesCmd() *cobra.Command {
	var (
		filter, function, start, end string
		buckets                      int
	)

	cmd := &cobra.Command{
		Use:   "timeseries",
		Short: "Run a time-series aggregation (SDL REST)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if filter == "" {
				return fmt.Errorf("--filter is required")
			}
			if start == "" {
				return fmt.Errorf("--start is required")
			}
			c, err := sdlClient()
			if err != nil {
				return err
			}
			resp, err := c.TimeseriesQuery(cmd.Context(), &sdl.TimeseriesQueryRequest{
				Queries: []sdl.TimeseriesQuery{{
					Filter:    filter,
					Function:  function,
					StartTime: start,
					EndTime:   end,
					Buckets:   buckets,
				}},
			})
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), resp)
			}
			var rows [][]string
			for _, res := range resp.Results {
				for i, v := range res.Values {
					val := "-"
					if v != nil {
						val = strconv.FormatFloat(*v, 'f', -1, 64)
					}
					rows = append(rows, []string{strconv.Itoa(i), val})
				}
			}
			printTable([]string{"BUCKET", "VALUE"}, rows)
			return nil
		},
	}
	cmd.Flags().StringVar(&filter, "filter", "", "query filter expression (required)")
	cmd.Flags().StringVar(&function, "function", "", "aggregation function (e.g. count, mean(field))")
	cmd.Flags().StringVar(&start, "start", "", "start time, e.g. 24h or timestamp (required)")
	cmd.Flags().StringVar(&end, "end", "", "end time")
	cmd.Flags().IntVar(&buckets, "buckets", 0, "number of time buckets")
	return cmd
}
