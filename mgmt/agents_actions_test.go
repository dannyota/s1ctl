package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAgentsDisconnect(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/agents/actions/disconnect" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req actionRequest
		json.NewDecoder(r.Body).Decode(&req)
		if len(req.Filter.IDs) != 1 || req.Filter.IDs[0] != "A1" {
			t.Fatalf("unexpected filter: %+v", req.Filter)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 1},
		})
	})
	c := testClient(t, handler)
	affected, err := c.AgentsDisconnect(context.Background(), ActionFilter{IDs: []string{"A1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestAgentActionRequiresFilter(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	_, err := c.AgentsDisconnect(context.Background(), ActionFilter{})
	if err == nil {
		t.Fatal("expected error for empty filter")
	}
}
