package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// DVQueryType is the Deep Visibility query type.
type DVQueryType string

const (
	DVQueryTypeEvents       DVQueryType = "events"
	DVQueryTypeProcessState DVQueryType = "processState"
)

// DVResponseState is the state of a Deep Visibility query.
type DVResponseState string

const (
	DVStateRunning        DVResponseState = "RUNNING"
	DVStateProcessRunning DVResponseState = "PROCESS_RUNNING"
	DVStateEventsRunning  DVResponseState = "EVENTS_RUNNING"
	DVStateFinished       DVResponseState = "FINISHED"
	DVStateFailed         DVResponseState = "FAILED"
	DVStateFailedClient   DVResponseState = "FAILED_CLIENT"
	DVStateError          DVResponseState = "ERROR"
	DVStateCancelled      DVResponseState = "QUERY_CANCELLED"
	DVStateTimedOut       DVResponseState = "TIMED_OUT"
	DVStateExpired        DVResponseState = "QUERY_EXPIRED"
)

// IsTerminal reports whether the state is a terminal state (query will not change further).
func (s DVResponseState) IsTerminal() bool {
	switch s {
	case DVStateFinished, DVStateFailed, DVStateFailedClient,
		DVStateError, DVStateCancelled, DVStateTimedOut, DVStateExpired:
		return true
	}
	return false
}

// IsSuccess reports whether the query completed successfully.
func (s DVResponseState) IsSuccess() bool { return s == DVStateFinished }

// DVQueryRequest is the body for POST /dv/init-query.
type DVQueryRequest struct {
	Query      string      `json:"query"`
	FromDate   string      `json:"fromDate"`
	ToDate     string      `json:"toDate"`
	QueryType  DVQueryType `json:"queryType,omitempty"`
	AccountIDs []string    `json:"accountIds,omitempty"`
	SiteIDs    []string    `json:"siteIds,omitempty"`
	IsVerbose  bool        `json:"isVerbose,omitempty"`
	Limit      int         `json:"limit,omitempty"`
}

// DVQueryID is the response from POST /dv/init-query.
type DVQueryID struct {
	QueryID string `json:"queryId"`

	Raw json.RawMessage `json:"-"`
}

func (d *DVQueryID) UnmarshalJSON(b []byte) error {
	type alias DVQueryID
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// DVQueryStatus is the response from GET /dv/query-status.
type DVQueryStatus struct {
	ResponseState  DVResponseState `json:"responseState"`
	ProgressStatus int             `json:"progressStatus"`
	ResponseError  string          `json:"responseError,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (d *DVQueryStatus) UnmarshalJSON(b []byte) error {
	type alias DVQueryStatus
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// DVEvent is a single Deep Visibility event.
type DVEvent struct {
	ID          string `json:"id"`
	EventType   string `json:"eventType"`
	ProcessName string `json:"processName"`
	AgentName   string `json:"agentName"`
	AgentOS     string `json:"agentOs"`
	CreatedAt   string `json:"createdAt"`
	User        string `json:"user"`
	ObjectType  string `json:"objectType"`
	ProcessCmd  string `json:"processCmd"`
	SrcIP       string `json:"agentIp"`
	DstIP       string `json:"dstIp"`
	DstPort     int    `json:"dstPort"`
	FilePath    string `json:"fileFullName"`
	SHA256      string `json:"sha256"`

	Raw json.RawMessage `json:"-"`
}

func (d DVEvent) MarshalJSON() ([]byte, error) {
	if d.Raw != nil {
		return d.Raw, nil
	}
	return []byte("{}"), nil
}

func (d *DVEvent) UnmarshalJSON(b []byte) error {
	type alias DVEvent
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// DVEventsParams are query parameters for GET /dv/events.
type DVEventsParams struct {
	QueryID   string
	Limit     int
	Cursor    string
	SortBy    string
	SortOrder string
	SubQuery  string
}

func (p *DVEventsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	v.Set("queryId", p.QueryID)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addString(v, "subQuery", p.SubQuery)
	return v
}

// dvInitResponse is the envelope for POST /dv/init-query.
type dvInitResponse struct {
	Data DVQueryID `json:"data"`
}

// dvStatusResponse is the envelope for GET /dv/query-status.
type dvStatusResponse struct {
	Data DVQueryStatus `json:"data"`
}

// dvEventsResponse is the envelope for GET /dv/events.
type dvEventsResponse struct {
	Data       []DVEvent  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// DVCreateQuery initiates a Deep Visibility query and returns the query ID.
func (c *Client) DVCreateQuery(ctx context.Context, req *DVQueryRequest) (*DVQueryID, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("mgmt: query is required")
	}
	if req.FromDate == "" || req.ToDate == "" {
		return nil, fmt.Errorf("mgmt: fromDate and toDate are required")
	}
	var resp dvInitResponse
	if err := c.post(ctx, "/dv/init-query", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DVGetQueryStatus checks the status of a Deep Visibility query.
func (c *Client) DVGetQueryStatus(ctx context.Context, queryID string) (*DVQueryStatus, error) {
	params := url.Values{}
	params.Set("queryId", queryID)
	var resp dvStatusResponse
	if err := c.get(ctx, "/dv/query-status", params, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DVGetEvents fetches Deep Visibility events for a completed query.
func (c *Client) DVGetEvents(ctx context.Context, p *DVEventsParams) ([]DVEvent, *Pagination, error) {
	if p == nil || p.QueryID == "" {
		return nil, nil, fmt.Errorf("mgmt: queryId is required")
	}
	if p.Limit == 0 {
		p.Limit = 100
	}
	var resp dvEventsResponse
	if err := c.get(ctx, "/dv/events", p.values(), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// DVCancelQuery cancels a running Deep Visibility query.
func (c *Client) DVCancelQuery(ctx context.Context, queryID string) error {
	body := map[string]string{"queryId": queryID}
	return c.post(ctx, "/dv/cancel-query", body, nil)
}
