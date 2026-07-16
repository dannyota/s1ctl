package mgmt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestMarketplaceCatalogList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/singularity-marketplace/applications-catalog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "10" {
			t.Fatalf("expected limit=10, got %s", got)
		}
		if got := q.Get("name__contains"); got != "slack" {
			t.Fatalf("expected name__contains=slack, got %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":          "cat-1",
					"name":        "Slack Integration",
					"key":         "slack",
					"category":    "Communication",
					"categoryId":  "cat-comm",
					"description": "Send alerts to Slack",
					"summary":     "Slack connector",
					"type":        "integration",
					"installed":   true,
					"toggleState": "enabled",
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": ""},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.MarketplaceCatalogList(context.Background(), &MarketplaceCatalogListParams{
		NameContains: "slack",
		Limit:        10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Name != "Slack Integration" {
		t.Fatalf("unexpected name: %s", items[0].Name)
	}
	if !items[0].Installed {
		t.Fatal("expected installed=true")
	}
	if items[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestMarketplaceCatalogListNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0, "nextCursor": ""},
		})
	})
	c := testClient(t, handler)
	items, _, err := c.MarketplaceCatalogList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestMarketplaceAppList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/singularity-marketplace/applications" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("siteIds"); got != "site-1" {
			t.Fatalf("expected siteIds=site-1, got %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"applicationCatalogId": "cat-1",
					"name":                 "My Slack",
					"hasAlert":             false,
					"lastInstalledAt":      "2025-06-01T00:00:00Z",
					"scopes": []map[string]any{
						{
							"id":                      "scope-1",
							"applicationInstanceName": "My Slack Instance",
							"status":                  "ACTIVE",
							"scopeLevel":              "site",
							"siteId":                  "site-1",
						},
					},
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": ""},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.MarketplaceAppList(context.Background(), &MarketplaceAppListParams{
		SiteIDs: []string{"site-1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Name != "My Slack" {
		t.Fatalf("unexpected name: %s", items[0].Name)
	}
	if items[0].ApplicationCatalogID != "cat-1" {
		t.Fatalf("unexpected catalogId: %s", items[0].ApplicationCatalogID)
	}
	if len(items[0].Scopes) != 1 {
		t.Fatalf("expected 1 scope, got %d", len(items[0].Scopes))
	}
	if items[0].Scopes[0].ID != "scope-1" {
		t.Fatalf("unexpected scope ID: %s", items[0].Scopes[0].ID)
	}
	if items[0].Scopes[0].Status != "ACTIVE" {
		t.Fatalf("unexpected scope status: %s", items[0].Scopes[0].Status)
	}
	if items[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestMarketplaceCatalogConfig(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/singularity-marketplace/applications-catalog/cat-1/config" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"configurationSchemaFields": []map[string]any{
					{"id": "webhook_url", "label": "Webhook URL", "type": "string"},
				},
			},
		})
	})
	c := testClient(t, handler)
	data, err := c.MarketplaceCatalogConfig(context.Background(), "cat-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil {
		t.Fatal("expected data to be non-nil")
	}
}

func TestMarketplaceCatalogConfigEmptyID(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.MarketplaceCatalogConfig(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty catalogId")
	}
}

func TestMarketplaceAppConfig(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/singularity-marketplace/applications/app-1/config" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"application":               map[string]any{"name": "My Slack"},
				"configurationSchemaFields": []any{},
			},
		})
	})
	c := testClient(t, handler)
	data, err := c.MarketplaceAppConfig(context.Background(), "app-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil {
		t.Fatal("expected data to be non-nil")
	}
}

func TestMarketplaceAppLog(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/singularity-marketplace/applications/app-1/log" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("only_errors"); got != "true" {
			t.Fatalf("expected only_errors=true, got %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"message": "error occurred", "level": "error"},
			},
		})
	})
	c := testClient(t, handler)
	onlyErrors := true
	entries, err := c.MarketplaceAppLog(context.Background(), "app-1", &onlyErrors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestMarketplaceInstall(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/singularity-marketplace/applications" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to parse body: %v", err)
		}
		data, ok := req["data"].(map[string]any)
		if !ok {
			t.Fatal("expected data field in body")
		}
		if data["applicationInstanceName"] != "My Slack" {
			t.Fatalf("unexpected name: %v", data["applicationInstanceName"])
		}
		filter, ok := req["filter"].(map[string]any)
		if !ok {
			t.Fatal("expected filter field in body")
		}
		if filter["applicationCatalogId"] != "cat-1" {
			t.Fatalf("unexpected catalogId: %v", filter["applicationCatalogId"])
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	c := testClient(t, handler)
	input := &MarketplaceInstallInput{}
	input.Data.Name = "My Slack"
	input.Data.Configurations = []MarketplaceConfig{{ID: "webhook_url", Value: "https://hooks.example.com"}}
	input.Filter.ApplicationCatalogID = "cat-1"
	input.Filter.SiteIDs = []string{"site-1"}
	if err := c.MarketplaceInstall(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarketplaceInstallNilConfigurations(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"configurations":[]`) {
			t.Fatalf("expected configurations:[] in body, got %s", string(body))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	c := testClient(t, handler)
	input := &MarketplaceInstallInput{}
	input.Data.Name = "Test App"
	// Configurations left nil — should marshal as [].
	input.Filter.ApplicationCatalogID = "cat-1"
	if err := c.MarketplaceInstall(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarketplaceUpdateNilConfigurations(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"configurations":[]`) {
			t.Fatalf("expected configurations:[] in body, got %s", string(body))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	c := testClient(t, handler)
	input := &MarketplaceUpdateInput{}
	input.Filter.IDs = []string{"app-1"}
	// Configurations left nil — should marshal as [].
	if err := c.MarketplaceUpdate(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarketplaceInstallNilInput(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if err := c.MarketplaceInstall(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestMarketplaceUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/singularity-marketplace/applications" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	c := testClient(t, handler)
	input := &MarketplaceUpdateInput{}
	input.Data.NameMap = map[string]string{"app-1": "Updated Slack"}
	input.Data.Configurations = []MarketplaceConfig{{ID: "webhook_url", Value: "https://new.example.com"}}
	input.Filter.IDs = []string{"app-1"}
	input.Filter.SiteIDs = []string{"site-1"}
	if err := c.MarketplaceUpdate(context.Background(), input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarketplaceDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/singularity-marketplace/applications" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to parse body: %v", err)
		}
		filter, ok := req["filter"].(map[string]any)
		if !ok {
			t.Fatal("expected filter field in body")
		}
		ids, ok := filter["id"].([]any)
		if !ok || len(ids) != 1 {
			t.Fatalf("expected filter.id with 1 entry, got %v", filter["id"])
		}
		// Verify the delete filter does NOT send the old "ids" key.
		if _, has := filter["ids"]; has {
			t.Fatal("delete filter must not send 'ids' key")
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	c := testClient(t, handler)
	filter := &MarketplaceDeleteFilter{ID: []string{"app-1"}}
	if err := c.MarketplaceDelete(context.Background(), filter); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarketplaceDeleteNilFilter(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if err := c.MarketplaceDelete(context.Background(), nil); err == nil {
		t.Fatal("expected error for nil filter")
	}
}

func TestMarketplaceDeleteEmptyFilter(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if err := c.MarketplaceDelete(context.Background(), &MarketplaceDeleteFilter{}); err == nil {
		t.Fatal("expected error for empty filter")
	}
}

func TestMarketplaceSetMode(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/singularity-marketplace/applications/enable" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	c := testClient(t, handler)
	filter := &MarketplaceScopeFilter{ApplicationID: "app-1"}
	if err := c.MarketplaceSetMode(context.Background(), "enable", filter); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarketplaceSetModeInvalid(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	filter := &MarketplaceScopeFilter{ApplicationID: "app-1"}
	if err := c.MarketplaceSetMode(context.Background(), "restart", filter); err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestMarketplaceSetModeNilFilter(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if err := c.MarketplaceSetMode(context.Background(), "enable", nil); err == nil {
		t.Fatal("expected error for nil filter")
	}
}

func TestMarketplaceAPIError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden", "detail": "Insufficient permissions"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.MarketplaceCatalogList(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
