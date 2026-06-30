package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSitesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"sites": []map[string]any{
					{
						"id": "S1", "name": "Default Site",
						"state": "active", "siteType": "Trial",
						"accountId": "A1", "totalLicenses": 100,
					},
				},
				"pagination": map[string]any{"totalItems": 1},
			},
		})
	})
	c := testClient(t, handler)
	sites, pag, err := c.SitesList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sites) != 1 {
		t.Fatalf("expected 1 site, got %d", len(sites))
	}
	if sites[0].ID != "S1" {
		t.Fatalf("unexpected id: %s", sites[0].ID)
	}
	if sites[0].Name != "Default Site" {
		t.Fatalf("unexpected name: %s", sites[0].Name)
	}
	if sites[0].TotalLicenses != 100 {
		t.Fatalf("expected 100 licenses, got %d", sites[0].TotalLicenses)
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
	if sites[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSitesGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("siteIds") != "S1" {
			t.Fatalf("expected siteIds=S1, got %s", r.URL.Query().Get("siteIds"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"sites": []map[string]any{
					{"id": "S1", "name": "Default Site"},
				},
			},
		})
	})
	c := testClient(t, handler)
	site, err := c.SitesGet(context.Background(), "S1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site.ID != "S1" {
		t.Fatalf("unexpected id: %s", site.ID)
	}
}

func TestSitesGetNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"sites": []any{},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.SitesGet(context.Background(), "MISSING")
	if err == nil {
		t.Fatal("expected error")
	}
}
