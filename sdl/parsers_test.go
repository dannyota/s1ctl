package sdl

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestParsersList(t *testing.T) {
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
		if !strings.Contains(gql.Query, "configFiles") {
			t.Fatalf("expected configFiles operation, got %s", gql.Query)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"configFiles": []map[string]any{
					{"udoId": "p-1", "name": "apache", "readOnly": true, "version": 3},
					{"udoId": "p-2", "name": "custom-app", "readOnly": false, "version": 1},
				},
			},
		})
	})
	c := testClient(t, handler)
	parsers, err := c.ParsersList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(parsers) != 2 {
		t.Fatalf("expected 2 parsers, got %d", len(parsers))
	}
	if parsers[0].Name != "apache" {
		t.Fatalf("unexpected name: %s", parsers[0].Name)
	}
	if !parsers[0].ReadOnly {
		t.Fatal("expected ReadOnly=true")
	}
	if parsers[0].Version != 3 {
		t.Fatalf("expected version 3, got %d", parsers[0].Version)
	}
	if parsers[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestParserGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "configFile") {
			t.Fatalf("expected configFile operation, got %s", gql.Query)
		}
		if gql.Variables["udoId"] != "p-1" {
			t.Fatalf("unexpected udoId: %v", gql.Variables["udoId"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"configFile": map[string]any{
					"udoId":        "p-1",
					"name":         "apache",
					"content":      "parser content here",
					"createdDate":  "2024-01-01T00:00:00Z",
					"modifiedDate": "2024-06-15T12:00:00Z",
					"readOnly":     true,
					"version":      3,
				},
			},
		})
	})
	c := testClient(t, handler)
	p, err := c.ParserGet(context.Background(), "p-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "apache" {
		t.Fatalf("unexpected name: %s", p.Name)
	}
	if p.Content != "parser content here" {
		t.Fatalf("unexpected content: %s", p.Content)
	}
	if p.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestParserGetValidation(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.ParserGet(context.Background(), "")
	if err == nil || !strings.Contains(err.Error(), "udoId is required") {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestParserCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "addConfigFile") {
			t.Fatalf("expected addConfigFile operation, got %s", gql.Query)
		}
		if gql.Variables["name"] != "my-parser" {
			t.Fatalf("unexpected name: %v", gql.Variables["name"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"addConfigFile": map[string]any{
					"udoId":        "p-new",
					"name":         "my-parser",
					"content":      "new content",
					"createdDate":  "2024-06-15T00:00:00Z",
					"modifiedDate": "2024-06-15T00:00:00Z",
					"readOnly":     false,
					"version":      1,
				},
			},
		})
	})
	c := testClient(t, handler)
	name := "my-parser"
	content := "new content"
	p, err := c.ParserCreate(context.Background(), &ParserCreateInput{
		Name:    &name,
		Content: &content,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "my-parser" {
		t.Fatalf("unexpected name: %s", p.Name)
	}
	if p.UdoID != "p-new" {
		t.Fatalf("unexpected udoId: %s", p.UdoID)
	}
}

func TestParserCreateValidation(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.ParserCreate(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "input is required") {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestParserDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdl/v2/graphql" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var gql struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&gql)
		if !strings.Contains(gql.Query, "deleteConfigFile") {
			t.Fatalf("expected deleteConfigFile operation, got %s", gql.Query)
		}
		if gql.Variables["udoId"] != "p-1" {
			t.Fatalf("unexpected udoId: %v", gql.Variables["udoId"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"deleteConfigFile": map[string]any{
					"udoId": "p-1",
					"name":  "apache",
				},
			},
		})
	})
	c := testClient(t, handler)
	if err := c.ParserDelete(context.Background(), "p-1", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParserDeleteValidation(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	err := c.ParserDelete(context.Background(), "", nil)
	if err == nil || !strings.Contains(err.Error(), "udoId is required") {
		t.Fatalf("expected validation error, got %v", err)
	}
}
