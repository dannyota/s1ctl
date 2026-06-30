package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestAgentsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/agents" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "1000000000", "computerName": "DESKTOP-01",
					"osType": "windows", "osName": "Windows 11",
					"isActive": true, "agentVersion": "24.1.2.3",
					"siteId": "2000000000", "siteName": "Default",
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": ""},
		})
	})
	c := testClient(t, handler)
	agents, pag, err := c.AgentsList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}
	if agents[0].ID != "1000000000" {
		t.Fatalf("unexpected id: %s", agents[0].ID)
	}
	if agents[0].ComputerName != "DESKTOP-01" {
		t.Fatalf("unexpected computerName: %s", agents[0].ComputerName)
	}
	if !agents[0].IsActive {
		t.Fatal("expected isActive=true")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
	if agents[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAgentsListParams(t *testing.T) {
	var gotQuery string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	active := true
	_, _, err := c.AgentsList(context.Background(), &AgentListParams{
		SiteIDs:  []string{"123"},
		IsActive: &active,
		Limit:    10,
		Query:    "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"siteIds=123", "isActive=true", "limit=10", "query=test"} {
		if !strings.Contains(gotQuery, want) {
			t.Errorf("query %q missing %q", gotQuery, want)
		}
	}
}

func TestAgentsGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ids") != "12345" {
			t.Fatalf("expected ids=12345, got %s", r.URL.Query().Get("ids"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "12345", "computerName": "HOST-1"},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	agent, err := c.AgentsGet(context.Background(), "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agent.ID != "12345" {
		t.Fatalf("unexpected id: %s", agent.ID)
	}
}

func TestAgentsGetNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	_, err := c.AgentsGet(context.Background(), "99999")
	if err == nil {
		t.Fatal("expected not found error")
	}
}

func TestAgentsCount(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("countOnly") != "true" {
			t.Fatal("expected countOnly=true")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 42},
		})
	})
	c := testClient(t, handler)
	count, err := c.AgentsCount(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 42 {
		t.Fatalf("expected 42, got %d", count)
	}
}
