package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"testing"
)

func TestRolesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/rbac/roles" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["accountIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected accountIds: %v", got)
		}
		if q.Get("predefinedRole") != "false" {
			t.Fatalf("unexpected predefinedRole: %s", q.Get("predefinedRole"))
		}
		if q.Get("query") != "admin" {
			t.Fatalf("unexpected query: %s", q.Get("query"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("unexpected limit: %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":             "1000000000000000001",
					"name":           "Custom IT",
					"description":    "IT operators",
					"scope":          "Site",
					"scopeId":        "1000000000000000009",
					"predefinedRole": false,
					"usersInRoles":   3,
					"siteName":       "Default Site",
					"createdAt":      "2025-01-01T00:00:00Z",
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "abc"},
		})
	})
	c := testClient(t, handler)
	pf := false
	roles, pag, err := c.RolesList(context.Background(), &RoleListParams{
		AccountIDs:     []string{"225494730938493804"},
		PredefinedRole: &pf,
		Query:          "admin",
		Limit:          10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}
	role := roles[0]
	if role.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", role.ID)
	}
	if role.Name != "Custom IT" {
		t.Fatalf("unexpected name: %s", role.Name)
	}
	if role.Scope != RoleScopeSite {
		t.Fatalf("unexpected scope: %s", role.Scope)
	}
	if role.UsersInRoles != 3 {
		t.Fatalf("unexpected usersInRoles: %d", role.UsersInRoles)
	}
	if role.PredefinedRole {
		t.Fatal("expected predefinedRole=false")
	}
	if role.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
	if pag.NextCursor != "abc" {
		t.Fatalf("unexpected cursor: %s", pag.NextCursor)
	}
}

func TestRoleGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/rbac/role/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":          "1000000000000000001",
				"name":        "Custom IT",
				"description": "IT operators",
				"scope":       "Account",
				"pages": []map[string]any{
					{
						"name":       "Endpoints",
						"identifier": "endpoints",
						"permissions": []map[string]any{
							{"identifier": "view", "title": "View", "value": true},
						},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	role, err := c.RoleGet(context.Background(), "1000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role.Name != "Custom IT" {
		t.Fatalf("unexpected name: %s", role.Name)
	}
	if role.Scope != RoleScopeAccount {
		t.Fatalf("unexpected scope: %s", role.Scope)
	}
	if len(role.Pages) == 0 {
		t.Fatal("expected pages permission blob to be populated")
	}
	if role.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestRoleTemplate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/rbac/role" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"description": "",
				"pages": []map[string]any{
					{"name": "Endpoints", "identifier": "endpoints", "permissions": []any{}},
				},
			},
		})
	})
	c := testClient(t, handler)
	tmpl, err := c.RoleTemplate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tmpl.Pages) == 0 {
		t.Fatal("expected template pages to be populated")
	}
}

func TestRoleCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/rbac/role" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				Name          string   `json:"name"`
				Description   string   `json:"description"`
				PermissionIDs []string `json:"permissionIds"`
			} `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data.Name != "Custom IT" {
			t.Fatalf("unexpected name: %s", body.Data.Name)
		}
		if body.Data.Description != "IT operators" {
			t.Fatalf("unexpected description: %s", body.Data.Description)
		}
		if !slices.Equal(body.Data.PermissionIDs, []string{"111", "222"}) {
			t.Fatalf("unexpected permissionIds: %v", body.Data.PermissionIDs)
		}
		if body.Filter["tenant"] != true {
			t.Fatalf("expected tenant=true, got %v", body.Filter["tenant"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "1000000000000000002", "name": "Custom IT"},
		})
	})
	c := testClient(t, handler)
	role, err := c.RoleCreate(context.Background(), RoleCreate{
		Data: RoleData{
			Name:          "Custom IT",
			Description:   "IT operators",
			PermissionIDs: []string{"111", "222"},
		},
		Filter: RoleScopeFilter{Tenant: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", role.ID)
	}
	if role.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestRoleUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/rbac/role/1000000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data.Description != "updated" {
			t.Fatalf("unexpected description: %s", body.Data.Description)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "1000000000000000002", "description": "updated"},
		})
	})
	c := testClient(t, handler)
	role, err := c.RoleUpdate(context.Background(), "1000000000000000002", RoleUpdate{
		Data: RoleData{Name: "Custom IT", Description: "updated"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", role.ID)
	}
}

func TestRoleDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/rbac/role/1000000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true},
		})
	})
	c := testClient(t, handler)
	if err := c.RoleDelete(context.Background(), "1000000000000000002"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoleScopeEnumValues(t *testing.T) {
	tests := []struct {
		got  string
		want string
	}{
		{string(RoleScopeGroup), "Group"},
		{string(RoleScopeSite), "Site"},
		{string(RoleScopeAccount), "Account"},
		{string(RoleScopeTenant), "Tenant"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Fatalf("expected %q, got %q", tt.want, tt.got)
		}
	}
}
