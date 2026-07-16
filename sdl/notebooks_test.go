package sdl

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestNotebooksList(t *testing.T) {
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
		if !strings.Contains(gql.Query, "purpleConversations") {
			t.Fatalf("expected purpleConversations operation, got %s", gql.Query)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"purpleConversations": []map[string]any{
					{
						"id": "nb-1", "name": "Investigation 1",
						"description": "Threat hunt", "isShared": false,
						"isReadOnly": false, "isAppendable": true,
						"accountId": "000000", "notebookSource": "MANUAL",
					},
					{
						"id": "nb-2", "name": "Shared Analysis",
						"description": "Team analysis", "isShared": true,
						"isReadOnly": true, "isAppendable": false,
						"accountId": "000000", "notebookSource": "ALERT",
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	notebooks, err := c.NotebooksList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notebooks) != 2 {
		t.Fatalf("expected 2 notebooks, got %d", len(notebooks))
	}
	if notebooks[0].Name != "Investigation 1" {
		t.Fatalf("unexpected name: %s", notebooks[0].Name)
	}
	if notebooks[0].IsShared {
		t.Fatal("expected IsShared=false")
	}
	if notebooks[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestNotebookGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "purpleConversation") {
			t.Fatalf("expected purpleConversation operation, got %s", gql.Query)
		}
		if gql.Variables["id"] != "nb-1" {
			t.Fatalf("unexpected id: %v", gql.Variables["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"purpleConversation": map[string]any{
					"id": "nb-1", "name": "Investigation 1",
					"description": "Threat hunt", "accountId": "000000",
					"isReadOnly": false, "isAppendable": true, "isShared": false,
					"notebookSource": "MANUAL",
					"entitlements":   map[string]any{"account": "000000"},
				},
			},
		})
	})
	c := testClient(t, handler)
	n, err := c.NotebookGet(context.Background(), "nb-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.Name != "Investigation 1" {
		t.Fatalf("unexpected name: %s", n.Name)
	}
	if n.Description != "Threat hunt" {
		t.Fatalf("unexpected description: %s", n.Description)
	}
	if n.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestNotebookGetValidation(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.NotebookGet(context.Background(), "")
	if err == nil || !strings.Contains(err.Error(), "id is required") {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestNotebookCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "createPurpleConversation") {
			t.Fatalf("expected createPurpleConversation operation, got %s", gql.Query)
		}
		if gql.Variables["name"] != "New Notebook" {
			t.Fatalf("unexpected name: %v", gql.Variables["name"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"createPurpleConversation": map[string]any{
					"id": "nb-new", "name": "New Notebook",
					"description": "Investigation notes",
					"accountId":   "000000", "isReadOnly": false,
					"isAppendable": true, "isShared": false,
					"notebookSource": "MANUAL",
				},
			},
		})
	})
	c := testClient(t, handler)
	n, err := c.NotebookCreate(context.Background(), "New Notebook", "Investigation notes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.Name != "New Notebook" {
		t.Fatalf("unexpected name: %s", n.Name)
	}
	if n.ID != "nb-new" {
		t.Fatalf("unexpected id: %s", n.ID)
	}
}

func TestNotebookCreateValidation(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.NotebookCreate(context.Background(), "", "desc")
	if err == nil || !strings.Contains(err.Error(), "name is required") {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestNotebookUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "updatePurpleConversation") {
			t.Fatalf("expected updatePurpleConversation operation, got %s", gql.Query)
		}
		if gql.Variables["id"] != "nb-1" {
			t.Fatalf("unexpected id: %v", gql.Variables["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"updatePurpleConversation": map[string]any{
					"id": "nb-1", "name": "Updated", "description": "new desc",
				},
			},
		})
	})
	c := testClient(t, handler)
	name := "Updated"
	err := c.NotebookUpdate(context.Background(), "nb-1", &NotebookUpdateInput{Name: &name})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotebookUpdateValidation(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	err := c.NotebookUpdate(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "id is required") {
		t.Fatalf("expected validation error, got %v", err)
	}

	err = c.NotebookUpdate(context.Background(), "nb-1", nil)
	if err == nil || !strings.Contains(err.Error(), "input is required") {
		t.Fatalf("expected input validation error, got %v", err)
	}
}

func TestNotebookDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "deletePurpleConversation") {
			t.Fatalf("expected deletePurpleConversation operation, got %s", gql.Query)
		}
		if gql.Variables["id"] != "nb-1" {
			t.Fatalf("unexpected id: %v", gql.Variables["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"deletePurpleConversation": true,
			},
		})
	})
	c := testClient(t, handler)
	if err := c.NotebookDelete(context.Background(), "nb-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotebookDeleteValidation(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	err := c.NotebookDelete(context.Background(), "")
	if err == nil || !strings.Contains(err.Error(), "id is required") {
		t.Fatalf("expected validation error, got %v", err)
	}
}
