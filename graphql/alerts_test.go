package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAlertsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointAlerts) {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alerts": map[string]any{
					"edges": []map[string]any{
						{
							"cursor": "c1",
							"node": map[string]any{
								"id":       "alert-1",
								"name":     "Test Alert",
								"severity": "HIGH",
								"status":   "IN_PROGRESS",
							},
						},
					},
					"pageInfo": map[string]any{
						"hasNextPage": false,
						"endCursor":   "c1",
					},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	conn, err := c.AlertsList(context.Background(), &AlertsListParams{First: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.TotalCount != 1 {
		t.Fatalf("expected totalCount=1, got %d", conn.TotalCount)
	}
	if len(conn.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(conn.Edges))
	}
	alert := conn.Edges[0].Node
	if alert.ID != "alert-1" {
		t.Fatalf("unexpected id: %s", alert.ID)
	}
	if alert.Name != "Test Alert" {
		t.Fatalf("unexpected name: %s", alert.Name)
	}
	if alert.Severity != "HIGH" {
		t.Fatalf("unexpected severity: %s", alert.Severity)
	}
	if alert.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAlertsListWithFilters(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req gqlRequest
		json.NewDecoder(r.Body).Decode(&req)
		gotVars = req.Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alerts": map[string]any{
					"edges":      []any{},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 0,
				},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.AlertsList(context.Background(), &AlertsListParams{
		First: 5,
		Filters: []Filter{
			{FieldID: "severity", StringIn: &InStr{Values: []string{"HIGH", "CRITICAL"}}},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotVars["first"] != float64(5) {
		t.Fatalf("expected first=5, got %v", gotVars["first"])
	}
	filters, ok := gotVars["filters"].([]any)
	if !ok || len(filters) != 1 {
		t.Fatalf("expected 1 filter, got %v", gotVars["filters"])
	}
}

func TestAlertsListNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alerts": map[string]any{
					"edges":      []any{},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 0,
				},
			},
		})
	})
	c := testClient(t, handler)
	conn, err := c.AlertsList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.TotalCount != 0 {
		t.Fatalf("expected totalCount=0, got %d", conn.TotalCount)
	}
}
