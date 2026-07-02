package sdl

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestDashboardsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query string `json:"query"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "dashboardsV2") {
			t.Fatalf("expected dashboardsV2 operation, got %s", gql.Query)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"dashboardsV2": []map[string]any{
					{"id": "d-1", "name": "Overview", "isBuiltIn": true, "isEditable": false},
					{"id": "d-2", "name": "Custom", "isBuiltIn": false, "isEditable": true},
				},
			},
		})
	})
	c := testClient(t, handler)
	dashboards, err := c.DashboardsList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dashboards) != 2 {
		t.Fatalf("expected 2 dashboards, got %d", len(dashboards))
	}
	if dashboards[0].Name != "Overview" {
		t.Fatalf("unexpected name: %s", dashboards[0].Name)
	}
	if !dashboards[0].IsBuiltIn {
		t.Fatal("expected IsBuiltIn=true")
	}
	if dashboards[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestDashboardGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "getDashboardV2") {
			t.Fatalf("expected getDashboardV2 operation, got %s", gql.Query)
		}
		if gql.Variables["id"] != "d-1" {
			t.Fatalf("unexpected id: %v", gql.Variables["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"getDashboardV2": map[string]any{
					"id":          "d-1",
					"name":        "Overview",
					"description": "Main dashboard",
					"configType":  "standard",
					"access":      map[string]any{"public": true},
					"tabs":        []any{map[string]any{"tabName": "Tab1"}},
				},
			},
		})
	})
	c := testClient(t, handler)
	d, err := c.DashboardGet(context.Background(), "d-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Name != "Overview" {
		t.Fatalf("unexpected name: %s", d.Name)
	}
	if d.Description != "Main dashboard" {
		t.Fatalf("unexpected description: %s", d.Description)
	}
	if d.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSavedSearchDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "deleteSavedSearchV2") {
			t.Fatalf("expected deleteSavedSearchV2 operation, got %s", gql.Query)
		}
		if gql.Variables["name"] != "my-query" {
			t.Fatalf("unexpected name: %v", gql.Variables["name"])
		}
		if gql.Variables["type"] != "PRIVATE" {
			t.Fatalf("unexpected type: %v", gql.Variables["type"])
		}
		idx, ok := gql.Variables["index"].(float64)
		if !ok || int(idx) != 0 {
			t.Fatalf("unexpected index: %v", gql.Variables["index"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"deleteSavedSearchV2": map[string]any{
					"name": "my-query", "url": "x", "index": 0, "type": "PRIVATE",
				},
			},
		})
	})
	c := testClient(t, handler)
	if err := c.SavedSearchDelete(context.Background(), "my-query", SavedSearchTypePrivate, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSavedSearchDeleteValidation(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	err := c.SavedSearchDelete(context.Background(), "", SavedSearchTypePrivate, 0)
	if err == nil || !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("expected name validation error, got %v", err)
	}
}
