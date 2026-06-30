package sdl

import (
	"context"
	"encoding/json"
)

// PowerQueryRequest is the request body for a PowerQuery.
type PowerQueryRequest struct {
	Query      string   `json:"query"`
	StartTime  string   `json:"startTime,omitempty"`
	EndTime    string   `json:"endTime,omitempty"`
	Priority   string   `json:"priority,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
}

// PowerQueryResponse is the response from a PowerQuery.
type PowerQueryResponse struct {
	Status  string             `json:"status"`
	Columns []PowerQueryColumn `json:"columns"`
	Values  [][]any            `json:"values"`

	Raw json.RawMessage `json:"-"`
}

func (r *PowerQueryResponse) UnmarshalJSON(b []byte) error {
	type alias PowerQueryResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// PowerQueryColumn describes a column in the result.
type PowerQueryColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// PowerQuery executes a PowerQuery against the SDL API.
func (c *Client) PowerQuery(ctx context.Context, req *PowerQueryRequest) (*PowerQueryResponse, error) {
	var resp PowerQueryResponse
	if err := c.post(ctx, "/api/powerQuery", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
