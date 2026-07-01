package sdl

import (
	"context"
	"encoding/json"
)

// FacetQueryRequest is the request body for a facet aggregation query.
type FacetQueryRequest struct {
	Filter    string `json:"filter,omitempty"`
	Field     string `json:"field"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime,omitempty"`
	MaxCount  int    `json:"maxCount,omitempty"`
	Priority  string `json:"priority,omitempty"`
}

// FacetQueryResponse is the response from a facet aggregation query.
type FacetQueryResponse struct {
	Status     string       `json:"status"`
	Values     []FacetEntry `json:"values"`
	MatchCount int64        `json:"matchCount"`
	CPUUsage   float64      `json:"cpuUsage"`

	Raw json.RawMessage `json:"-"`
}

func (r *FacetQueryResponse) UnmarshalJSON(b []byte) error {
	type alias FacetQueryResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// FacetEntry is a field value and its occurrence count (REST /api/facetQuery).
type FacetEntry struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// FacetQuery gets the most frequent values of a field in matching events.
func (c *Client) FacetQuery(ctx context.Context, req *FacetQueryRequest) (*FacetQueryResponse, error) {
	type body struct {
		QueryType string `json:"queryType"`
		*FacetQueryRequest
	}
	var resp FacetQueryResponse
	if err := c.post(ctx, "/api/facetQuery", &body{QueryType: "facet", FacetQueryRequest: req}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
