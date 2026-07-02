package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"testing"
)

func TestTagRulesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/xdr/assets/tags/rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("status") != "enabled" {
			t.Fatalf("unexpected status: %s", r.URL.Query().Get("status"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":         "400000000000000001",
					"name":       "Tag servers",
					"status":     "enabled",
					"conditions": map[string]any{"op": "and"},
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": nil},
		})
	})
	c := testClient(t, handler)
	rules, pag, err := c.TagRulesList(context.Background(), &TagRuleListParams{Status: "enabled"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Name != "Tag servers" {
		t.Fatalf("unexpected name: %s", rules[0].Name)
	}
	if len(rules[0].Conditions) == 0 {
		t.Fatal("expected conditions to be populated")
	}
	if rules[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestTagRulesCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/xdr/assets/tags/rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Name       string          `json:"name"`
			Conditions json.RawMessage `json:"conditions"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Name != "Tag servers" {
			t.Fatalf("unexpected name: %s", body.Name)
		}
		if len(body.Conditions) == 0 {
			t.Fatal("expected conditions in body")
		}
		json.NewEncoder(w).Encode(map[string]any{"id": "400000000000000002", "name": "Tag servers"})
	})
	c := testClient(t, handler)
	rule, err := c.TagRulesCreate(context.Background(), TagRuleWrite{
		Name:       "Tag servers",
		Conditions: json.RawMessage(`{"op":"and"}`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.ID != "400000000000000002" {
		t.Fatalf("unexpected ID: %s", rule.ID)
	}
}

func TestTagRulesUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/xdr/assets/tags/rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.ID != "400000000000000002" {
			t.Fatalf("expected id in body, got %s", body.ID)
		}
		json.NewEncoder(w).Encode(map[string]any{"id": "400000000000000002", "name": "Tag servers"})
	})
	c := testClient(t, handler)
	rule, err := c.TagRulesUpdate(context.Background(), TagRuleWrite{
		ID:         "400000000000000002",
		Name:       "Tag servers",
		Conditions: json.RawMessage(`{"op":"and"}`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.ID != "400000000000000002" {
		t.Fatalf("unexpected ID: %s", rule.ID)
	}
}

func TestTagRulesDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/xdr/assets/tags/rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["ids"]; !slices.Equal(got, []string{"400000000000000002"}) {
			t.Fatalf("unexpected ids: %v", got)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": nil})
	})
	c := testClient(t, handler)
	if err := c.TagRulesDelete(context.Background(), []string{"400000000000000002"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagRulesTest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/xdr/assets/tags/rules/test" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 42, "nextCursor": nil},
		})
	})
	c := testClient(t, handler)
	count, err := c.TagRulesTest(context.Background(), TagRuleWrite{
		Name:       "Tag servers",
		Conditions: json.RawMessage(`{"op":"and"}`),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 42 {
		t.Fatalf("expected 42 matches, got %d", count)
	}
}
