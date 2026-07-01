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
