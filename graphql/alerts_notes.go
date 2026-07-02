package graphql

import (
	"context"
	"encoding/json"
)

// AlertUser identifies the author of an alert note or timeline entry.
type AlertUser struct {
	UserID   string `json:"userId"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`

	Raw json.RawMessage `json:"-"`
}

func (u *AlertUser) UnmarshalJSON(b []byte) error {
	type alias AlertUser
	if err := json.Unmarshal(b, (*alias)(u)); err != nil {
		return err
	}
	u.Raw = append(u.Raw[:0:0], b...)
	return nil
}

// AlertNote is an investigation note attached to an alert.
type AlertNote struct {
	ID        string     `json:"id"`
	AlertID   string     `json:"alertId"`
	Text      string     `json:"text"`
	Type      string     `json:"type"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
	Author    *AlertUser `json:"author"`

	Raw json.RawMessage `json:"-"`
}

// AuthorName returns the note author's full name, or empty string.
func (n *AlertNote) AuthorName() string {
	if n.Author != nil {
		return n.Author.FullName
	}
	return ""
}

func (n *AlertNote) UnmarshalJSON(b []byte) error {
	type alias AlertNote
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// alertNotesResponse mirrors AlertNotesListResponse { data: [AlertNote!]! }.
type alertNotesResponse struct {
	Data []AlertNote `json:"data"`
}

const alertNotesQuery = `query AlertNotes($alertId: ID!) {
  alertNotes(alertId: $alertId) {
    data {
      id
      alertId
      text
      type
      createdAt
      updatedAt
      author { userId fullName email }
    }
  }
}`

// AlertNotes returns the investigation notes on an alert.
//
// The alertNotes query takes only alertId; it is not paginated and returns
// the full note list.
func (c *Client) AlertNotes(ctx context.Context, alertID string) ([]AlertNote, error) {
	vars := map[string]any{"alertId": alertID}
	var resp struct {
		AlertNotes alertNotesResponse `json:"alertNotes"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertNotesQuery, vars, &resp); err != nil {
		return nil, err
	}
	return resp.AlertNotes.Data, nil
}

const updateAlertNoteMutation = `mutation UpdateAlertNote($alertNoteId: ID!, $text: String!) {
  updateAlertNote(alertNoteId: $alertNoteId, text: $text) {
    data {
      id
      alertId
      text
      type
      createdAt
      updatedAt
      author { userId fullName email }
    }
  }
}`

// AlertsUpdateNote updates the text of an existing alert note and returns the
// alert's notes after the update.
func (c *Client) AlertsUpdateNote(ctx context.Context, noteID, text string) ([]AlertNote, error) {
	vars := map[string]any{"alertNoteId": noteID, "text": text}
	var resp struct {
		UpdateAlertNote alertNotesResponse `json:"updateAlertNote"`
	}
	if err := c.Do(ctx, EndpointAlerts, updateAlertNoteMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.UpdateAlertNote.Data, nil
}

const deleteAlertNoteMutation = `mutation DeleteAlertNote($alertNoteId: ID!) {
  deleteAlertNote(alertNoteId: $alertNoteId) {
    data {
      id
      alertId
      text
      type
      createdAt
      updatedAt
      author { userId fullName email }
    }
  }
}`

// AlertsDeleteNote deletes an alert note and returns the alert's remaining notes.
func (c *Client) AlertsDeleteNote(ctx context.Context, noteID string) ([]AlertNote, error) {
	vars := map[string]any{"alertNoteId": noteID}
	var resp struct {
		DeleteAlertNote alertNotesResponse `json:"deleteAlertNote"`
	}
	if err := c.Do(ctx, EndpointAlerts, deleteAlertNoteMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.DeleteAlertNote.Data, nil
}

// AlertTimelineCreator identifies the actor behind a timeline entry.
type AlertTimelineCreator struct {
	UserID   string `json:"userId"`
	UserType string `json:"userType"`

	Raw json.RawMessage `json:"-"`
}

func (a *AlertTimelineCreator) UnmarshalJSON(b []byte) error {
	type alias AlertTimelineCreator
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// AlertTimelineEntry is a single item in an alert's timeline.
type AlertTimelineEntry struct {
	CreatedAt string                `json:"createdAt"`
	EventText string                `json:"eventText"`
	EventType string                `json:"eventType"`
	Creator   *AlertTimelineCreator `json:"timelineItemCreator"`

	Raw json.RawMessage `json:"-"`
}

// ActorName returns the timeline entry actor's user ID, or empty string.
func (e *AlertTimelineEntry) ActorName() string {
	if e.Creator != nil {
		return e.Creator.UserID
	}
	return ""
}

func (e *AlertTimelineEntry) UnmarshalJSON(b []byte) error {
	type alias AlertTimelineEntry
	if err := json.Unmarshal(b, (*alias)(e)); err != nil {
		return err
	}
	e.Raw = append(e.Raw[:0:0], b...)
	return nil
}

const alertTimelineQuery = `query AlertTimeline($alertId: ID!, $first: Int, $after: String, $filter: AlertTimelineFilterInput) {
  alertTimeline(alertId: $alertId, first: $first, after: $after, filter: $filter) {
    edges {
      cursor
      node {
        createdAt
        eventText
        eventType
        timelineItemCreator {
          ... on UserTimelineItemCreator { userId userType }
        }
      }
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      endCursor
      startCursor
    }
    totalCount
  }
}`

// AlertTimeline returns the timeline (activity + notes + indicators) for an alert.
func (c *Client) AlertTimeline(ctx context.Context, alertID string, first int, after string) (*Connection[AlertTimelineEntry], error) {
	vars := map[string]any{"alertId": alertID}
	if first > 0 {
		vars["first"] = first
	}
	if after != "" {
		vars["after"] = after
	}
	var resp struct {
		AlertTimeline Connection[AlertTimelineEntry] `json:"alertTimeline"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertTimelineQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.AlertTimeline, nil
}
