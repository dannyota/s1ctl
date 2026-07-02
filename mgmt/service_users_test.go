package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"testing"
)

func TestServiceUsersList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/service-users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"1000000000000000001"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if q.Get("query") != "integration" {
			t.Fatalf("unexpected query: %s", q.Get("query"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":          "1000000000000000002",
					"name":        "ci-bot",
					"description": "CI integration",
					"scope":       "tenant",
					"createdAt":   "2025-01-01T00:00:00Z",
					"apiToken": map[string]any{
						"createdAt": "2025-01-01T00:00:00Z",
						"expiresAt": "2026-01-01T00:00:00Z",
					},
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "next"},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.ServiceUsersList(context.Background(), &ServiceUserListParams{
		SiteIDs: []string{"1000000000000000001"},
		Query:   "integration",
		Limit:   10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	su := items[0]
	if su.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", su.ID)
	}
	if su.Name != "ci-bot" {
		t.Fatalf("unexpected name: %s", su.Name)
	}
	if su.Scope != ServiceUserScopeTenant {
		t.Fatalf("unexpected scope: %s", su.Scope)
	}
	// Read responses must never carry the raw token value.
	if su.APIToken.Value != "" {
		t.Fatalf("expected empty apiToken value on read, got %q", su.APIToken.Value)
	}
	if su.APIToken.ExpiresAt != "2026-01-01T00:00:00Z" {
		t.Fatalf("unexpected expiresAt: %s", su.APIToken.ExpiresAt)
	}
	if su.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestServiceUsersGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/service-users/1000000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":    "1000000000000000002",
				"name":  "ci-bot",
				"scope": "account",
				"scopeRoles": []map[string]any{
					{"id": "1000000000000000009", "roleId": "1000000000000000010", "roleName": "Viewer"},
				},
			},
		})
	})
	c := testClient(t, handler)
	su, err := c.ServiceUsersGet(context.Background(), "1000000000000000002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if su.Name != "ci-bot" {
		t.Fatalf("unexpected name: %s", su.Name)
	}
	if len(su.ScopeRoles) != 1 || su.ScopeRoles[0].RoleName != "Viewer" {
		t.Fatalf("unexpected scopeRoles: %v", su.ScopeRoles)
	}
}

func TestServiceUsersCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/service-users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data["name"] != "ci-bot" {
			t.Fatalf("unexpected name: %v", body.Data["name"])
		}
		if body.Data["scope"] != "tenant" {
			t.Fatalf("unexpected scope: %v", body.Data["scope"])
		}
		if body.Data["expirationDate"] != "2026-01-01T00:00:00Z" {
			t.Fatalf("unexpected expirationDate: %v", body.Data["expirationDate"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":    "1000000000000000002",
				"name":  "ci-bot",
				"scope": "tenant",
				"apiToken": map[string]any{
					"value":     "PLACEHOLDER_TOKEN_VALUE",
					"createdAt": "2025-01-01T00:00:00Z",
					"expiresAt": "2026-01-01T00:00:00Z",
				},
			},
		})
	})
	c := testClient(t, handler)
	su, err := c.ServiceUsersCreate(context.Background(), ServiceUserCreate{
		Name:           "ci-bot",
		Scope:          ServiceUserScopeTenant,
		ExpirationDate: "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if su.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", su.ID)
	}
	if su.APIToken.Value != "PLACEHOLDER_TOKEN_VALUE" {
		t.Fatalf("expected token value on create, got %q", su.APIToken.Value)
	}
}

func TestServiceUsersUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/service-users/1000000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["description"] != "updated" {
			t.Fatalf("unexpected description: %v", body.Data["description"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "1000000000000000002", "description": "updated"},
		})
	})
	c := testClient(t, handler)
	su, err := c.ServiceUsersUpdate(context.Background(), "1000000000000000002", ServiceUserUpdate{
		Description: "updated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if su.Description != "updated" {
		t.Fatalf("unexpected description: %s", su.Description)
	}
}

func TestServiceUsersDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/service-users/1000000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"success": true}})
	})
	c := testClient(t, handler)
	if err := c.ServiceUsersDelete(context.Background(), "1000000000000000002"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServiceUsersBulkDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/service-users/delete-service-users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter struct {
				IDs []string `json:"ids"`
			} `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.IDs, []string{"a", "b"}) {
			t.Fatalf("unexpected ids: %v", body.Filter.IDs)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 2}})
	})
	c := testClient(t, handler)
	affected, err := c.ServiceUsersBulkDelete(context.Background(), []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 affected, got %d", affected)
	}
}

func TestServiceUsersGenerateToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/service-users/1000000000000000002/generate-api-token" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["expirationDate"] != "2026-01-01T00:00:00Z" {
			t.Fatalf("unexpected expirationDate: %v", body.Data["expirationDate"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"token":     "PLACEHOLDER_TOKEN_VALUE",
				"createdAt": "2025-01-01T00:00:00Z",
				"expiresAt": "2026-01-01T00:00:00Z",
			},
		})
	})
	c := testClient(t, handler)
	tok, err := c.ServiceUsersGenerateToken(context.Background(), "1000000000000000002", "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Token != "PLACEHOLDER_TOKEN_VALUE" {
		t.Fatalf("unexpected token: %s", tok.Token)
	}
	if tok.ExpiresAt != "2026-01-01T00:00:00Z" {
		t.Fatalf("unexpected expiresAt: %s", tok.ExpiresAt)
	}
}

func TestServiceUsersExport(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/export/service-users" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte("id,name,scope\n1,ci-bot,tenant\n"))
	})
	c := testClient(t, handler)
	data, err := c.ServiceUsersExport(context.Background(), &ServiceUserListParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "id,name,scope\n1,ci-bot,tenant\n" {
		t.Fatalf("unexpected export body: %q", string(data))
	}
}

func TestServiceUserScopeEnumValues(t *testing.T) {
	if string(ServiceUserScopeTenant) != "tenant" {
		t.Fatalf("unexpected tenant: %s", ServiceUserScopeTenant)
	}
	if string(ServiceUserScopeAccount) != "account" {
		t.Fatalf("unexpected account: %s", ServiceUserScopeAccount)
	}
	if string(ServiceUserScopeSite) != "site" {
		t.Fatalf("unexpected site: %s", ServiceUserScopeSite)
	}
}
