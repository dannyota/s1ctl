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

func TestSitesReactivate(t *testing.T) {
	cases := []struct {
		name       string
		unlimited  bool
		expiration string
	}{
		{"unlimited", true, ""},
		{"expiration", false, "2027-01-01T00:00:00Z"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Fatalf("expected PUT, got %s", r.Method)
				}
				if r.URL.Path != "/sites/S1/reactivate" {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				var req struct {
					Data struct {
						Unlimited  bool   `json:"unlimited"`
						Expiration string `json:"expiration"`
					} `json:"data"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatalf("decode body: %v", err)
				}
				if req.Data.Unlimited != tc.unlimited {
					t.Fatalf("unlimited: got %v want %v", req.Data.Unlimited, tc.unlimited)
				}
				if req.Data.Expiration != tc.expiration {
					t.Fatalf("expiration: got %q want %q", req.Data.Expiration, tc.expiration)
				}
				w.Write([]byte(`{"data":{}}`))
			})
			c := testClient(t, handler)
			if err := c.SitesReactivate(context.Background(), "S1", tc.unlimited, tc.expiration); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSitesExpireNow(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/sites/S1/expire-now" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{"data":{}}`))
	})
	c := testClient(t, handler)
	if err := c.SitesExpireNow(context.Background(), "S1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSitesDuplicate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/sites/duplicate-site" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req struct {
			Data struct {
				Name         string `json:"name"`
				SourceSiteID int64  `json:"sourceSiteId"`
				PolicySource string `json:"policySource"`
				CopyUsers    bool   `json:"copyUsers"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Data.Name != "clone" {
			t.Fatalf("unexpected name: %s", req.Data.Name)
		}
		if req.Data.SourceSiteID != 42 {
			t.Fatalf("unexpected sourceSiteId: %d", req.Data.SourceSiteID)
		}
		if req.Data.PolicySource != "inherit_global" {
			t.Fatalf("unexpected policySource: %s", req.Data.PolicySource)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "S9", "name": "clone", "state": "active"},
		})
	})
	c := testClient(t, handler)
	site, err := c.SitesDuplicate(context.Background(), SiteDuplicate{
		Name:         "clone",
		SourceSiteID: 42,
		PolicySource: PolicySourceInheritGlobal,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site.ID != "S9" {
		t.Fatalf("unexpected id: %s", site.ID)
	}
	if site.Name != "clone" {
		t.Fatalf("unexpected name: %s", site.Name)
	}
}

func TestSitesRegenerateKey(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/sites/S1/regenerate-key" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"registrationToken": "REG-PLACEHOLDER"},
		})
	})
	c := testClient(t, handler)
	tok, err := c.SitesRegenerateKey(context.Background(), "S1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value() != "REG-PLACEHOLDER" {
		t.Fatalf("unexpected token value: %q", tok.Value())
	}
}

func TestSitesToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/sites/S1/token" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"token": "TOK-PLACEHOLDER"},
		})
	})
	c := testClient(t, handler)
	tok, err := c.SitesToken(context.Background(), "S1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value() != "TOK-PLACEHOLDER" {
		t.Fatalf("unexpected token value: %q", tok.Value())
	}
}
