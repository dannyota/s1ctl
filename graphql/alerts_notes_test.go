package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestAlertNotes(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointAlerts) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertNotes": map[string]any{
					"data": []map[string]any{
						{
							"id":        "note-1",
							"alertId":   "alert-1",
							"text":      "looks malicious",
							"type":      "PLAIN_TEXT",
							"createdAt": "2024-01-01T00:00:00Z",
							"updatedAt": "2024-01-01T00:00:00Z",
							"author": map[string]any{
								"userId":   "u1",
								"fullName": "Ana Lyst",
								"email":    "ana@example.com",
							},
						},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	notes, err := c.AlertNotes(context.Background(), "alert-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "alertNotes(alertId: $alertId)") {
		t.Errorf("query does not target alertNotes: %s", gotReq.Query)
	}
	if gotReq.Variables["alertId"] != "alert-1" {
		t.Errorf("expected alertId=alert-1, got %v", gotReq.Variables["alertId"])
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	n := notes[0]
	if n.ID != "note-1" || n.Text != "looks malicious" || n.Type != "PLAIN_TEXT" {
		t.Errorf("unexpected note: %+v", n)
	}
	if n.Author == nil || n.Author.FullName != "Ana Lyst" {
		t.Errorf("unexpected author: %+v", n.Author)
	}
	if n.AuthorName() != "Ana Lyst" {
		t.Errorf("unexpected author name: %s", n.AuthorName())
	}
	if n.Raw == nil {
		t.Error("expected Raw to be populated")
	}
}

func TestAlertsUpdateNote(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"updateAlertNote": map[string]any{
					"data": []map[string]any{
						{"id": "note-1", "alertId": "alert-1", "text": "revised", "createdAt": "t", "updatedAt": "t"},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	notes, err := c.AlertsUpdateNote(context.Background(), "note-1", "revised")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "updateAlertNote(alertNoteId: $alertNoteId") {
		t.Errorf("query does not target updateAlertNote: %s", gotReq.Query)
	}
	if gotReq.Variables["alertNoteId"] != "note-1" {
		t.Errorf("expected alertNoteId=note-1, got %v", gotReq.Variables["alertNoteId"])
	}
	if gotReq.Variables["text"] != "revised" {
		t.Errorf("expected text=revised, got %v", gotReq.Variables["text"])
	}
	if len(notes) != 1 || notes[0].Text != "revised" {
		t.Fatalf("unexpected notes: %+v", notes)
	}
}

func TestAlertsDeleteNote(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"deleteAlertNote": map[string]any{"data": []any{}},
			},
		})
	})
	c := testClient(t, handler)
	notes, err := c.AlertsDeleteNote(context.Background(), "note-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "deleteAlertNote(alertNoteId: $alertNoteId)") {
		t.Errorf("query does not target deleteAlertNote: %s", gotReq.Query)
	}
	if gotReq.Variables["alertNoteId"] != "note-1" {
		t.Errorf("expected alertNoteId=note-1, got %v", gotReq.Variables["alertNoteId"])
	}
	if len(notes) != 0 {
		t.Fatalf("expected 0 notes after delete, got %d", len(notes))
	}
}

func TestAlertsUpdateNoteGraphQLError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":   nil,
			"errors": []map[string]any{{"message": "note not found"}},
		})
	})
	c := testClient(t, handler)
	_, err := c.AlertsUpdateNote(context.Background(), "missing", "x")
	if err == nil {
		t.Fatal("expected error")
	}
	var qe *QueryError
	if !errors.As(err, &qe) {
		t.Fatalf("expected *QueryError, got %T", err)
	}
}

func TestAlertTimeline(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointAlerts) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertTimeline": map[string]any{
					"edges": []map[string]any{
						{
							"cursor": "t1",
							"node": map[string]any{
								"createdAt": "2024-01-01T00:00:00Z",
								"eventText": "note added",
								"eventType": "NOTE",
								"timelineItemCreator": map[string]any{
									"userId":   "u1",
									"userType": "USER",
								},
							},
						},
					},
					"pageInfo":   map[string]any{"hasNextPage": false, "endCursor": "t1"},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	conn, err := c.AlertTimeline(context.Background(), "alert-1", 25, "cur-0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "alertTimeline(alertId: $alertId") {
		t.Errorf("query does not target alertTimeline: %s", gotReq.Query)
	}
	if gotReq.Variables["alertId"] != "alert-1" {
		t.Errorf("expected alertId=alert-1, got %v", gotReq.Variables["alertId"])
	}
	if gotReq.Variables["first"] != float64(25) {
		t.Errorf("expected first=25, got %v", gotReq.Variables["first"])
	}
	if gotReq.Variables["after"] != "cur-0" {
		t.Errorf("expected after=cur-0, got %v", gotReq.Variables["after"])
	}
	if conn.TotalCount != 1 || len(conn.Edges) != 1 {
		t.Fatalf("unexpected connection: total=%d edges=%d", conn.TotalCount, len(conn.Edges))
	}
	entry := conn.Edges[0].Node
	if entry.EventType != "NOTE" || entry.EventText != "note added" {
		t.Errorf("unexpected entry: %+v", entry)
	}
	if entry.ActorName() != "u1" {
		t.Errorf("unexpected actor: %s", entry.ActorName())
	}
	if entry.Raw == nil {
		t.Error("expected Raw to be populated")
	}
}

func TestAlertTimelineOmitsOptionalVars(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertTimeline": map[string]any{
					"edges":      []any{},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 0,
				},
			},
		})
	})
	c := testClient(t, handler)
	if _, err := c.AlertTimeline(context.Background(), "alert-1", 0, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotReq.Variables["first"]; ok {
		t.Errorf("expected first to be omitted, got %v", gotReq.Variables["first"])
	}
	if _, ok := gotReq.Variables["after"]; ok {
		t.Errorf("expected after to be omitted, got %v", gotReq.Variables["after"])
	}
}
