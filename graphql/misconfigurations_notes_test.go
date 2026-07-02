package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestMisconfigurationsNotes(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointMisconfigurations) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"misconfigurationNotes": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"id":                 "note-1",
							"misconfigurationId": "m-1",
							"text":               "review this",
							"createdAt":          "2024-01-01T00:00:00Z",
							"updatedAt":          "2024-01-02T00:00:00Z",
							"author":             map[string]any{"id": "u1", "fullName": "Ana Lyst", "email": "ana@example.com", "deleted": false},
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	notes, err := c.MisconfigurationsNotes(context.Background(), "m-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "misconfigurationNotes(") {
		t.Errorf("query does not target misconfigurationNotes: %s", gotReq.Query)
	}
	if gotReq.Variables["misconfigurationId"] != "m-1" {
		t.Errorf("expected misconfigurationId=m-1, got %v", gotReq.Variables["misconfigurationId"])
	}
	if len(notes) != 1 || notes[0].ID != "note-1" || notes[0].Text != "review this" {
		t.Fatalf("unexpected notes: %+v", notes)
	}
	if notes[0].AuthorName() != "Ana Lyst" {
		t.Errorf("unexpected author name: %s", notes[0].AuthorName())
	}
	if notes[0].Raw == nil {
		t.Error("expected Raw to be populated")
	}
}

func TestMisconfigurationsAddNote(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{"addMisconfigurationNoteV2": map[string]any{"updatedFindingIds": []string{"m-1"}}},
		})
	})
	c := testClient(t, handler)
	if err := c.MisconfigurationsAddNote(context.Background(), []string{"m-1"}, "note text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "addMisconfigurationNoteV2(") {
		t.Errorf("query does not target addMisconfigurationNoteV2: %s", gotReq.Query)
	}
	if gotReq.Variables["text"] != "note text" {
		t.Errorf("expected text=note text, got %v", gotReq.Variables["text"])
	}
	if gotReq.Variables["filter"] == nil {
		t.Error("expected filter to be set")
	}
}

func TestMisconfigurationsUpdateNote(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)                                                               //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"updateMisconfigurationNote": true}}) //nolint:errcheck
	})
	c := testClient(t, handler)
	if err := c.MisconfigurationsUpdateNote(context.Background(), "note-1", "revised"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "updateMisconfigurationNote(") {
		t.Errorf("query does not target updateMisconfigurationNote: %s", gotReq.Query)
	}
	if gotReq.Variables["noteId"] != "note-1" || gotReq.Variables["text"] != "revised" {
		t.Errorf("unexpected variables: %v", gotReq.Variables)
	}
}

func TestMisconfigurationsDeleteNote(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)                                                               //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"deleteMisconfigurationNote": true}}) //nolint:errcheck
	})
	c := testClient(t, handler)
	if err := c.MisconfigurationsDeleteNote(context.Background(), "note-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "deleteMisconfigurationNote(") {
		t.Errorf("query does not target deleteMisconfigurationNote: %s", gotReq.Query)
	}
	if gotReq.Variables["noteId"] != "note-1" {
		t.Errorf("expected noteId=note-1, got %v", gotReq.Variables["noteId"])
	}
}

func TestMisconfigurationsAssign(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{"misconfigurationUserAssignmentV2": map[string]any{"updatedFindingIds": []string{"m-1"}}},
		})
	})
	c := testClient(t, handler)
	if err := c.MisconfigurationsAssign(context.Background(), []string{"m-1"}, "user-9"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "misconfigurationUserAssignmentV2(") {
		t.Errorf("query does not target misconfigurationUserAssignmentV2: %s", gotReq.Query)
	}
	if gotReq.Variables["userId"] != "user-9" {
		t.Errorf("expected userId=user-9, got %v", gotReq.Variables["userId"])
	}
	if gotReq.Variables["filter"] == nil {
		t.Error("expected filter to be set")
	}
}

func TestMisconfigurationsHistory(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"misconfigurationHistory": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"createdAt": "2024-01-01T00:00:00Z", "eventText": "status changed", "eventType": "STATUS_CHANGE",
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	items, err := c.MisconfigurationsHistory(context.Background(), "m-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "misconfigurationHistory(") {
		t.Errorf("query does not target misconfigurationHistory: %s", gotReq.Query)
	}
	if gotReq.Variables["misconfigurationId"] != "m-1" {
		t.Errorf("expected misconfigurationId=m-1, got %v", gotReq.Variables["misconfigurationId"])
	}
	if len(items) != 1 || items[0].EventType != "STATUS_CHANGE" {
		t.Fatalf("unexpected history: %+v", items)
	}
}

func TestMisconfigurationsRelatedAssets(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"misconfigurationRelatedAssets": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"misconfigurationId": "m-1",
							"organization":       "org-a",
							"asset":              map[string]any{"id": "a1", "name": "host-1", "type": "COMPUTE"},
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	assets, err := c.MisconfigurationsRelatedAssets(context.Background(), "m-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "misconfigurationRelatedAssets(") {
		t.Errorf("query does not target misconfigurationRelatedAssets: %s", gotReq.Query)
	}
	if len(assets) != 1 || assets[0].Asset.Name != "host-1" || assets[0].Organization != "org-a" {
		t.Fatalf("unexpected related assets: %+v", assets)
	}
}

func TestMisconfigurationsExport(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{"misconfigurationsExportToCsv": map[string]any{"data": "id,name\nm-1,Test\n"}},
		})
	})
	c := testClient(t, handler)
	csv, err := c.MisconfigurationsExport(context.Background(), []Filter{{FieldID: "severity", StringIn: &InStr{Values: []string{"HIGH"}}}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "misconfigurationsExportToCsv(") {
		t.Errorf("query does not target misconfigurationsExportToCsv: %s", gotReq.Query)
	}
	if !strings.Contains(csv, "m-1,Test") {
		t.Errorf("unexpected csv: %q", csv)
	}
}
