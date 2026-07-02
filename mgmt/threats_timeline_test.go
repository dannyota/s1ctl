package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestThreatTimeline(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/threats/1000000000000000001/timeline" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["activityTypes"]; !slices.Equal(got, []string{"1001", "1002"}) {
			t.Fatalf("unexpected activityTypes: %v", got)
		}
		if q.Get("query") != "example" {
			t.Fatalf("unexpected query: %s", q.Get("query"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", q.Get("limit"))
		}
		if q.Get("sortBy") != "createdAt" {
			t.Fatalf("unexpected sortBy: %s", q.Get("sortBy"))
		}
		if q.Get("sortOrder") != "desc" {
			t.Fatalf("unexpected sortOrder: %s", q.Get("sortOrder"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "1000000000000000050", "activityType": 1001,
					"primaryDescription":   "Threat detected",
					"secondaryDescription": "File quarantined",
					"threatId":             "1000000000000000001",
					"accountId":            "225494730938493804",
					"siteId":               "225494730938493805",
					"agentId":              "225494730938493806",
					"data":                 map[string]any{"confidenceLevel": "malicious"},
					"createdAt":            "2025-01-15T10:00:00Z",
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "abc123"},
		})
	})
	c := testClient(t, handler)
	entries, pag, err := c.ThreatTimeline(context.Background(), "1000000000000000001", &ThreatTimelineParams{
		ActivityTypes: []int{1001, 1002},
		Query:         "example",
		Limit:         10,
		SortBy:        "createdAt",
		SortOrder:     "desc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.ID != "1000000000000000050" {
		t.Fatalf("unexpected id: %s", e.ID)
	}
	if e.ActivityType != 1001 {
		t.Fatalf("unexpected activityType: %d", e.ActivityType)
	}
	if e.PrimaryDescription != "Threat detected" {
		t.Fatalf("unexpected primaryDescription: %s", e.PrimaryDescription)
	}
	if e.SecondaryDescription != "File quarantined" {
		t.Fatalf("unexpected secondaryDescription: %s", e.SecondaryDescription)
	}
	if e.ThreatID != "1000000000000000001" {
		t.Fatalf("unexpected threatId: %s", e.ThreatID)
	}
	if e.AccountID != "225494730938493804" {
		t.Fatalf("unexpected accountId: %s", e.AccountID)
	}
	if e.AgentID != "225494730938493806" {
		t.Fatalf("unexpected agentId: %s", e.AgentID)
	}
	if e.Data == nil {
		t.Fatal("expected Data to be populated")
	}
	if e.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
	if pag.NextCursor != "abc123" {
		t.Fatalf("unexpected cursor: %s", pag.NextCursor)
	}
}

func TestThreatTimelineNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/threats/1000000000000000001/timeline" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %q", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	entries, _, err := c.ThreatTimeline(context.Background(), "1000000000000000001", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestThreatTimelineEmptyID(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	_, _, err := c.ThreatTimeline(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error for empty threat ID")
	}
}

func TestThreatTimelinePathEscape(t *testing.T) {
	var gotRawPath string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRawPath = r.URL.RawPath
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.ThreatTimeline(context.Background(), "t/1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotRawPath != "/threats/t%2F1/timeline" {
		t.Fatalf("unexpected raw path (escape missing): %s", gotRawPath)
	}
}

func TestThreatTimelineMarshalJSON(t *testing.T) {
	// MarshalJSON returns Raw bytes when present.
	raw := json.RawMessage(`{"id":"1000000000000000050","activityType":1001,"primaryDescription":"Test"}`)
	entry := ThreatTimelineEntry{
		ID:                 "1000000000000000050",
		ActivityType:       1001,
		PrimaryDescription: "Test",
		Raw:                raw,
	}
	b, err := json.Marshal(&entry)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	if string(b) != string(raw) {
		t.Fatalf("expected raw passthrough, got %s", string(b))
	}

	// When Raw is nil, MarshalJSON falls back to struct fields.
	entry2 := ThreatTimelineEntry{
		ID:                 "1000000000000000051",
		ActivityType:       1002,
		PrimaryDescription: "Fallback",
	}
	b2, err := json.Marshal(&entry2)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var roundtrip map[string]any
	json.Unmarshal(b2, &roundtrip)
	if roundtrip["id"] != "1000000000000000051" {
		t.Fatalf("unexpected id in roundtrip: %v", roundtrip["id"])
	}
	if roundtrip["primaryDescription"] != "Fallback" {
		t.Fatalf("unexpected primaryDescription in roundtrip: %v", roundtrip["primaryDescription"])
	}
}

func TestThreatTimelineError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.ThreatTimeline(context.Background(), "1000000000000000001", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 403 {
		t.Fatalf("expected 403, got %d", ae.Status)
	}
}
