package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
	"danny.vn/s1/mgmt"
	"danny.vn/s1/sdl"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Verify connectivity to all SentinelOne API surfaces",
		Args:  cobra.NoArgs,
		RunE:  runDoctor,
	}
}

type checkResult struct {
	Surface string `json:"surface"`
	OK      bool   `json:"ok"`
	Latency string `json:"latency"`
	Error   string `json:"error,omitempty"`
}

func runDoctor(cmd *cobra.Command, _ []string) error {
	consoleURL, token, err := resolveConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	results := make([]checkResult, 3)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() { defer wg.Done(); results[0] = checkMGMT(ctx, consoleURL, token) }()
	go func() { defer wg.Done(); results[1] = checkGraphQL(ctx, consoleURL, token) }()
	go func() { defer wg.Done(); results[2] = checkSDL(ctx, consoleURL, token) }()
	wg.Wait()

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), results)
	}

	allOK := true
	for _, r := range results {
		status := "ok"
		if !r.OK {
			status = "FAIL"
			allOK = false
		}
		line := fmt.Sprintf("  %-12s %s (%s)", r.Surface, status, r.Latency)
		if r.Error != "" {
			line += "  " + r.Error
		}
		fmt.Fprintln(cmd.OutOrStdout(), line)
	}
	if !allOK {
		return fmt.Errorf("one or more API surfaces unreachable")
	}
	return nil
}

func checkMGMT(ctx context.Context, consoleURL, token string) checkResult {
	c := mgmt.NewClient(consoleURL, token)
	start := time.Now()
	_, err := c.SystemStatus(ctx)
	elapsed := time.Since(start).Round(time.Millisecond)
	r := checkResult{Surface: "REST MGMT", Latency: elapsed.String()}
	if err != nil {
		r.Error = err.Error()
	} else {
		r.OK = true
	}
	return r
}

func checkGraphQL(ctx context.Context, consoleURL, token string) checkResult {
	c := graphql.NewClient(consoleURL, token)
	start := time.Now()
	_, err := c.AlertsFiltersCount(ctx, []string{"severity"}, nil, nil)
	elapsed := time.Since(start).Round(time.Millisecond)
	r := checkResult{Surface: "GraphQL", Latency: elapsed.String()}
	if err != nil {
		r.Error = err.Error()
	} else {
		r.OK = true
	}
	return r
}

func checkSDL(ctx context.Context, consoleURL, token string) checkResult {
	c := sdl.NewClient(consoleURL, token)
	start := time.Now()
	_, err := c.DashboardsList(ctx)
	elapsed := time.Since(start).Round(time.Millisecond)
	r := checkResult{Surface: "SDL", Latency: elapsed.String()}
	if err != nil {
		r.Error = err.Error()
	} else {
		r.OK = true
	}
	return r
}
