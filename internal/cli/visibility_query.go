package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newVisibilityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "visibility",
		Short: "Run Deep Visibility queries",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newVisibilityQueryCmd())
	return cmd
}

func newVisibilityQueryCmd() *cobra.Command {
	var (
		query             string
		fromFlag, toFlag  string
		siteIDs           []string
		limit, maxResults int
		sortBy, sortOrder string
		pollInterval      time.Duration
	)

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Run a Deep Visibility query",
		Long: `Run a Deep Visibility query to hunt for endpoint events.

Initiates a query, polls until complete, then fetches and displays results.
The query uses SentinelOne's Deep Visibility query language.

Examples:
  s1ctl visibility query --query "EventType = \"Process Creation\""
  s1ctl visibility query --query "ProcessName contains \"cmd.exe\"" --from 7d
  s1ctl visibility query --query "SHA256 = \"abc123...\"" --json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if query == "" {
				return fmt.Errorf("--query is required")
			}
			return runVisibilityQuery(cmd, visibilityQueryOpts{
				query:        query,
				from:         fromFlag,
				to:           toFlag,
				siteIDs:      siteIDs,
				limit:        limit,
				maxResults:   maxResults,
				sortBy:       sortBy,
				sortOrder:    sortOrder,
				pollInterval: pollInterval,
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "Deep Visibility query expression (required)")
	cmd.Flags().StringVar(&fromFlag, "from", "24h", "start time (duration like 24h/7d, or RFC3339)")
	cmd.Flags().StringVar(&toFlag, "to", "", "end time (default: now)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().IntVar(&limit, "limit", 100, "max events per page (1-1000)")
	cmd.Flags().IntVar(&maxResults, "max-results", 0, "stop after fetching this many events (0 = all)")
	cmd.Flags().StringVar(&sortBy, "sort-by", "createdAt", "sort field (e.g. createdAt, pid)")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "desc", "sort direction (asc, desc)")
	cmd.Flags().DurationVar(&pollInterval, "poll-interval", 2*time.Second, "interval between status polls")
	return cmd
}

type visibilityQueryOpts struct {
	query        string
	from, to     string
	siteIDs      []string
	limit        int
	maxResults   int
	sortBy       string
	sortOrder    string
	pollInterval time.Duration
}

func runVisibilityQuery(cmd *cobra.Command, opts visibilityQueryOpts) error {
	c, err := mgmtClient()
	if err != nil {
		return err
	}

	fromDate, err := resolveTime(opts.from, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("invalid --from: %w", err)
	}
	toDate := time.Now().UTC()
	if opts.to != "" {
		toDate, err = resolveTime(opts.to, 0)
		if err != nil {
			return fmt.Errorf("invalid --to: %w", err)
		}
	}

	req := &mgmt.DVQueryRequest{
		Query:    opts.query,
		FromDate: fromDate.Format(time.RFC3339),
		ToDate:   toDate.Format(time.RFC3339),
		SiteIDs:  opts.siteIDs,
	}

	// Step 1: initiate query.
	var queryID string
	err = runWithSpinner("Initiating query...", func() error {
		resp, initErr := c.DVCreateQuery(cmd.Context(), req)
		if initErr != nil {
			return initErr
		}
		queryID = resp.QueryID
		return nil
	})
	if err != nil {
		return err
	}

	// Step 2: poll until terminal state.
	err = runWithSpinner("Running query...", func() error {
		return pollDVQuery(cmd, c, queryID, opts.pollInterval)
	})
	if err != nil {
		return err
	}

	// Step 3: fetch events.
	events, total, err := fetchDVEvents(c, cmd, queryID, opts)
	if err != nil {
		return err
	}

	return printDVEvents(cmd, events, total)
}

func pollDVQuery(cmd *cobra.Command, c *mgmt.Client, queryID string, interval time.Duration) error {
	for {
		status, err := c.DVGetQueryStatus(cmd.Context(), queryID)
		if err != nil {
			return err
		}
		if status.ResponseState.IsTerminal() {
			if !status.ResponseState.IsSuccess() {
				msg := string(status.ResponseState)
				if status.ResponseError != "" {
					msg += ": " + status.ResponseError
				}
				return fmt.Errorf("query failed: %s", msg)
			}
			return nil
		}
		select {
		case <-cmd.Context().Done():
			// Best-effort cancel.
			_ = c.DVCancelQuery(cmd.Context(), queryID)
			return cmd.Context().Err()
		case <-time.After(interval):
		}
	}
}

func fetchDVEvents(c *mgmt.Client, cmd *cobra.Command, queryID string, opts visibilityQueryOpts) ([]mgmt.DVEvent, int, error) {
	params := &mgmt.DVEventsParams{
		QueryID:   queryID,
		Limit:     opts.limit,
		SortBy:    opts.sortBy,
		SortOrder: opts.sortOrder,
	}

	var all []mgmt.DVEvent
	var total int

	for {
		var events []mgmt.DVEvent
		var pag *mgmt.Pagination
		var err error

		if len(all) == 0 {
			err = runWithSpinner("Fetching events...", func() error {
				events, pag, err = c.DVGetEvents(cmd.Context(), params)
				return err
			})
		} else {
			events, pag, err = c.DVGetEvents(cmd.Context(), params)
		}
		if err != nil {
			clearProgress()
			return nil, 0, err
		}

		all = append(all, events...)
		if pag != nil {
			total = pag.TotalItems
		}
		printProgress("event", len(all), total)

		// Stop if no more pages.
		if pag == nil || pag.NextCursor == "" {
			break
		}
		// Stop if max-results reached.
		if opts.maxResults > 0 && len(all) >= opts.maxResults {
			all = all[:opts.maxResults]
			break
		}
		params.Cursor = pag.NextCursor
	}
	clearProgress()
	return all, total, nil
}

func printDVEvents(cmd *cobra.Command, events []mgmt.DVEvent, total int) error {
	if outputFormat == "json" {
		// Stable shape: always a bare array; truncation is signalled on
		// stderr so it never changes what stdout consumers parse.
		if total > len(events) {
			fmt.Fprintf(cmd.ErrOrStderr(), "Showing %d of %d events. Use --max-results to fetch more.\n", len(events), total)
		}
		if events == nil {
			events = []mgmt.DVEvent{}
		}
		return printJSON(cmd.OutOrStdout(), events)
	}

	if len(events) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No events found.")
		return nil
	}

	headers := []string{"Timestamp", "EventType", "Process", "Agent", "User", "File/Dst"}
	rows := make([][]string, len(events))
	for i, e := range events {
		detail := e.FilePath
		if detail == "" && e.DstIP != "" {
			detail = e.DstIP
			if e.DstPort > 0 {
				detail += fmt.Sprintf(":%d", e.DstPort)
			}
		}
		rows[i] = []string{
			orDash(e.CreatedAt),
			orDash(e.EventType),
			orDash(truncate(e.ProcessName, 30)),
			orDash(truncate(e.AgentName, 20)),
			orDash(truncate(e.User, 20)),
			orDash(truncate(detail, 40)),
		}
	}
	return printOutput(cmd.OutOrStdout(), headers, rows, events, len(events), total, "event", total <= len(events))
}

// resolveTime parses a time string as either a duration (e.g. "24h", "7d")
// subtracted from now, or an RFC3339 timestamp.
func resolveTime(s string, defaultDur time.Duration) (time.Time, error) {
	if s == "" {
		return time.Now().UTC().Add(-defaultDur), nil
	}
	// Try duration-style: "24h", "7d", "30m".
	if d, ok := parseDuration(s); ok {
		return time.Now().UTC().Add(-d), nil
	}
	// Try RFC3339.
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("expected duration (e.g. 24h, 7d) or RFC3339 timestamp, got %q", s)
	}
	return t, nil
}

// parseDuration parses Go-style durations plus "d" for days.
func parseDuration(s string) (time.Duration, bool) {
	if len(s) == 0 {
		return 0, false
	}
	// Handle "Nd" (days) by converting to hours.
	if s[len(s)-1] == 'd' {
		var n int
		if _, err := fmt.Sscanf(s, "%dd", &n); err == nil && n > 0 {
			return time.Duration(n) * 24 * time.Hour, true
		}
		return 0, false
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, false
	}
	return d, true
}
