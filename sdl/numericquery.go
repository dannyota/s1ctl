package sdl

import (
	"context"
	"encoding/json"
)

// NumericQueryRequest is the request body for a numeric aggregation query
// (REST POST /api/numericQuery).
type NumericQueryRequest struct {
	Filter    string `json:"filter,omitempty"`
	Function  string `json:"function,omitempty"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime,omitempty"`
	Buckets   int    `json:"buckets,omitempty"`
	Priority  string `json:"priority,omitempty"`
}

// NumericQueryResponse is the response from a numeric aggregation query.
type NumericQueryResponse struct {
	Status   string     `json:"status"`
	Values   []*float64 `json:"values"`
	CPUUsage float64    `json:"cpuUsage"`

	Raw json.RawMessage `json:"-"`
}

func (r *NumericQueryResponse) UnmarshalJSON(b []byte) error {
	type alias NumericQueryResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// NumericQuery executes a numeric aggregation query against the SDL REST API.
//
// The API is effectively deprecated in favour of TimeseriesQuery with
// createSummaries=false, but remains useful for sub-30-second bucket
// granularity and users with limited query permissions.
func (c *Client) NumericQuery(ctx context.Context, req *NumericQueryRequest) (*NumericQueryResponse, error) {
	type body struct {
		QueryType string `json:"queryType"`
		*NumericQueryRequest
	}
	var resp NumericQueryResponse
	if err := c.post(ctx, "/api/numericQuery", &body{QueryType: "numeric", NumericQueryRequest: req}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
