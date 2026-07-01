package sdl

import (
	"context"
	"encoding/json"
)

// TimeseriesQueryRequest is the request body for a time-series aggregation query.
type TimeseriesQueryRequest struct {
	Queries []TimeseriesQuery `json:"queries"`
}

// TimeseriesQuery specifies a single time-series query within the request.
type TimeseriesQuery struct {
	Filter           string `json:"filter"`
	Function         string `json:"function,omitempty"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime,omitempty"`
	Buckets          int    `json:"buckets,omitempty"`
	CreateSummaries  *bool  `json:"createSummaries,omitempty"`
	OnlyUseSummaries *bool  `json:"onlyUseSummaries,omitempty"`
	Priority         string `json:"priority,omitempty"`
}

// TimeseriesQueryResponse is the response from a time-series aggregation query.
type TimeseriesQueryResponse struct {
	Status  string             `json:"status"`
	Results []TimeseriesResult `json:"results"`

	Raw json.RawMessage `json:"-"`
}

func (r *TimeseriesQueryResponse) UnmarshalJSON(b []byte) error {
	type alias TimeseriesQueryResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// TimeseriesResult holds the values for a single query in the request.
type TimeseriesResult struct {
	Values              []*float64 `json:"values"`
	CPUUsage            float64    `json:"cpuUsage"`
	FoundExistingSeries bool       `json:"foundExistingSeries"`

	Raw json.RawMessage `json:"-"`
}

func (t *TimeseriesResult) UnmarshalJSON(b []byte) error {
	type alias TimeseriesResult
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

// TimeseriesQuery executes a time-series aggregation query.
func (c *Client) TimeseriesQuery(ctx context.Context, req *TimeseriesQueryRequest) (*TimeseriesQueryResponse, error) {
	var resp TimeseriesQueryResponse
	if err := c.post(ctx, "/api/timeseriesQuery", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
