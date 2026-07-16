package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestConfigOverrideList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/config-override" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if q.Get("sortBy") != "createdAt" {
			t.Fatalf("unexpected sortBy: %s", q.Get("sortBy"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":            "1000000000000000001",
					"name":          "test-override",
					"description":   "test desc",
					"config":        map[string]any{"key": "value"},
					"osType":        "linux",
					"agentVersion":  "2.5.1.1320",
					"versionOption": "ALL",
					"scope":         "site",
					"site":          map[string]any{"id": "225494730938493804", "name": "Default"},
					"createdAt":     "2025-01-01T00:00:00Z",
					"updatedAt":     "2025-01-02T00:00:00Z",
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "abc"},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.ConfigOverrideList(context.Background(), &ConfigOverrideListParams{
		SiteIDs: []string{"225494730938493804"},
		SortBy:  "createdAt",
		Limit:   10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	item := items[0]
	if item.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", item.ID)
	}
	if item.Name != "test-override" {
		t.Fatalf("unexpected name: %s", item.Name)
	}
	if item.OSType != ConfigOverrideOSLinux {
		t.Fatalf("unexpected osType: %s", item.OSType)
	}
	if item.VersionOption != ConfigOverrideVersionAll {
		t.Fatalf("unexpected versionOption: %s", item.VersionOption)
	}
	if item.Scope != ConfigOverrideScopeSite {
		t.Fatalf("unexpected scope: %s", item.Scope)
	}
	if item.Site == nil || item.Site.ID != "225494730938493804" {
		t.Fatal("expected site ref to be populated")
	}
	if item.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestConfigOverrideListNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %q", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	items, _, err := c.ConfigOverrideList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestConfigOverrideGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("ids") != "1000000000000000001" {
			t.Fatalf("unexpected ids: %s", r.URL.Query().Get("ids"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "1000000000000000001", "name": "test", "description": "d", "config": map[string]any{}},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	item, err := c.ConfigOverrideGet(context.Background(), "1000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", item.ID)
	}
	if item.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestConfigOverrideGetNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	_, err := c.ConfigOverrideGet(context.Background(), "999")
	if err == nil {
		t.Fatal("expected error for not found")
	}
}

func TestConfigOverrideCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/config-override" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data["name"] != "my-override" {
			t.Fatalf("unexpected name: %v", body.Data["name"])
		}
		if body.Data["osType"] != "windows" {
			t.Fatalf("unexpected osType: %v", body.Data["osType"])
		}
		if body.Data["scope"] != "site" {
			t.Fatalf("unexpected scope: %v", body.Data["scope"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":          "1000000000000000002",
				"name":        "my-override",
				"description": "",
				"config":      map[string]any{"key": "val"},
				"osType":      "windows",
				"scope":       "site",
			},
		})
	})
	c := testClient(t, handler)
	item, err := c.ConfigOverrideCreate(context.Background(), ConfigOverrideCreateInput{
		Name:   "my-override",
		OSType: ConfigOverrideOSWindows,
		Config: json.RawMessage(`{"key":"val"}`),
		Scope:  ConfigOverrideScopeSite,
		Site:   &ConfigOverrideScopeRef{ID: "225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", item.ID)
	}
	if item.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestConfigOverrideUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/config-override/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["name"] != "updated" {
			t.Fatalf("unexpected name: %v", body.Data["name"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":          "1000000000000000001",
				"name":        "updated",
				"description": "",
				"config":      map[string]any{},
			},
		})
	})
	c := testClient(t, handler)
	name := "updated"
	item, err := c.ConfigOverrideUpdate(context.Background(), "1000000000000000001", ConfigOverrideUpdateInput{
		Name: &name,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", item.ID)
	}
}

func TestConfigOverrideDeleteByID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/config-override/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true},
		})
	})
	c := testClient(t, handler)
	err := c.ConfigOverrideDelete(context.Background(), "1000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfigOverrideBulkDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/config-override" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		ids, _ := body.Filter["ids"].([]any)
		if len(ids) != 2 {
			t.Fatalf("expected 2 ids, got %d", len(ids))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 2},
		})
	})
	c := testClient(t, handler)
	affected, err := c.ConfigOverrideBulkDelete(context.Background(), ConfigOverrideDeleteFilter{
		IDs: []string{"a", "b"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 affected, got %d", affected)
	}
}

func TestConfigOverrideBulkDeleteTenantFalse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		tenant, ok := body.Filter["tenant"]
		if !ok {
			t.Fatal("expected tenant field to be present in filter")
		}
		if tenant != false {
			t.Fatalf("expected tenant=false, got %v", tenant)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 0},
		})
	})
	c := testClient(t, handler)
	f := false
	_, err := c.ConfigOverrideBulkDelete(context.Background(), ConfigOverrideDeleteFilter{
		Tenant: &f,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfigOverrideListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{{"code": 403, "title": "Forbidden"}},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.ConfigOverrideList(context.Background(), nil)
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

func TestConfigOverrideEnumValues(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"OSLinux", string(ConfigOverrideOSLinux), "linux"},
		{"OSMacOS", string(ConfigOverrideOSMacOS), "macos"},
		{"OSWindows", string(ConfigOverrideOSWindows), "windows"},
		{"OSWindowsLegacy", string(ConfigOverrideOSWindowsLegacy), "windows_legacy"},
		{"VersionAll", string(ConfigOverrideVersionAll), "ALL"},
		{"VersionSpecific", string(ConfigOverrideVersionSpecific), "SPECIFIC"},
		{"ScopeGroup", string(ConfigOverrideScopeGroup), "group"},
		{"ScopeSite", string(ConfigOverrideScopeSite), "site"},
		{"ScopeAccount", string(ConfigOverrideScopeAccount), "account"},
		{"ScopeTenant", string(ConfigOverrideScopeTenant), "tenant"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}
