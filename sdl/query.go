package sdl

import (
	"context"
	"encoding/json"
)

// LogQueryRequest is the request body for a log query (REST /api/query).
type LogQueryRequest struct {
	Filter            string `json:"filter,omitempty"`
	StartTime         string `json:"startTime,omitempty"`
	EndTime           string `json:"endTime,omitempty"`
	MaxCount          int    `json:"maxCount,omitempty"`
	PageMode          string `json:"pageMode,omitempty"`
	Columns           string `json:"columns,omitempty"`
	ContinuationToken string `json:"continuationToken,omitempty"`
	Priority          string `json:"priority,omitempty"`
}

// LogQueryResponse is the response from a log query (REST /api/query).
type LogQueryResponse struct {
	Status            string                     `json:"status"`
	Matches           []Match                    `json:"matches"`
	Sessions          map[string]json.RawMessage `json:"sessions"`
	ContinuationToken string                     `json:"continuationToken"`
	CPUUsage          float64                    `json:"cpuUsage"`

	Raw json.RawMessage `json:"-"`
}

func (r *LogQueryResponse) UnmarshalJSON(b []byte) error {
	type alias LogQueryResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// Match is a single log event returned by a query.
type Match struct {
	Timestamp  string         `json:"timestamp"`
	Message    string         `json:"message"`
	Severity   int            `json:"severity"`
	Session    string         `json:"session"`
	Thread     string         `json:"thread"`
	Attributes map[string]any `json:"attributes"`
}

// Query executes a log query against the SDL REST API.
func (c *Client) Query(ctx context.Context, req *LogQueryRequest) (*LogQueryResponse, error) {
	type body struct {
		QueryType string `json:"queryType"`
		*LogQueryRequest
	}
	var resp LogQueryResponse
	if err := c.post(ctx, "/api/query", &body{QueryType: "log", LogQueryRequest: req}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryOption customizes QueryAll behavior.
type QueryOption func(*queryConfig)

type queryConfig struct {
	maxEvents int
	onPage    func(fetched int)
}

// WithMaxEvents caps the total number of events returned across all pages.
// Zero or negative means no limit.
func WithMaxEvents(n int) QueryOption {
	return func(c *queryConfig) { c.maxEvents = n }
}

// WithPageCallback registers a function called after each page with the
// running total of events fetched so far.
func WithPageCallback(fn func(fetched int)) QueryOption {
	return func(c *queryConfig) { c.onPage = fn }
}

// QueryAll executes a log query and follows continuation tokens until all
// results are fetched or the maxEvents cap is reached.
func (c *Client) QueryAll(ctx context.Context, req *LogQueryRequest, opts ...QueryOption) (*LogQueryResponse, error) {
	var cfg queryConfig
	for _, o := range opts {
		o(&cfg)
	}

	resp, err := c.Query(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg.onPage != nil {
		cfg.onPage(len(resp.Matches))
	}

	for resp.ContinuationToken != "" {
		if cfg.maxEvents > 0 && len(resp.Matches) >= cfg.maxEvents {
			resp.Matches = resp.Matches[:cfg.maxEvents]
			break
		}

		page := *req
		page.ContinuationToken = resp.ContinuationToken
		next, err := c.Query(ctx, &page)
		if err != nil {
			return nil, err
		}
		resp.Matches = append(resp.Matches, next.Matches...)
		resp.ContinuationToken = next.ContinuationToken
		if cfg.onPage != nil {
			cfg.onPage(len(resp.Matches))
		}
	}

	if cfg.maxEvents > 0 && len(resp.Matches) > cfg.maxEvents {
		resp.Matches = resp.Matches[:cfg.maxEvents]
	}

	return resp, nil
}
