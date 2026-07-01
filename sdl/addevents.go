package sdl

import (
	"context"
	"encoding/json"
)

// AddEventsRequest is the request body for event ingestion.
type AddEventsRequest struct {
	Session     string         `json:"session"`
	SessionInfo map[string]any `json:"sessionInfo,omitempty"`
	Events      []Event        `json:"events,omitempty"`
	Threads     []Thread       `json:"threads,omitempty"`
	Logs        []LogMeta      `json:"logs,omitempty"`
}

// Event is a single log event for ingestion.
type Event struct {
	TS     string         `json:"ts"`
	Sev    *int           `json:"sev,omitempty"`
	Thread string         `json:"thread,omitempty"`
	Log    string         `json:"log,omitempty"`
	Attrs  map[string]any `json:"attrs"`
}

// Thread maps a thread ID to a human-readable name.
type Thread struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LogMeta defines constant metadata shared across events in a request.
type LogMeta struct {
	ID    string         `json:"id"`
	Attrs map[string]any `json:"attrs"`
}

// AddEventsResponse is the response from event ingestion.
type AddEventsResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	BytesCharged int64  `json:"bytesCharged"`

	Raw json.RawMessage `json:"-"`
}

func (r *AddEventsResponse) UnmarshalJSON(b []byte) error {
	type alias AddEventsResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// AddEvents ingests one or more structured log events.
func (c *Client) AddEvents(ctx context.Context, req *AddEventsRequest) (*AddEventsResponse, error) {
	var resp AddEventsResponse
	if err := c.post(ctx, "/api/addEvents", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
