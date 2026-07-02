package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestXDRAssetCounts(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/xdr/assets/asset-counts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"categories": map[string]any{
					"device":      map[string]any{"count": 10},
					"server":      map[string]any{"count": 5},
					"workstation": map[string]any{"count": 3},
					"container":   map[string]any{"count": 2},
					"identity":    map[string]any{"count": 1},
				},
				"surfaces": map[string]any{
					"cloud":            map[string]any{"count": 8},
					"endpoint":         map[string]any{"count": 7},
					"identity":         map[string]any{"count": 4},
					"network":          map[string]any{"count": 6},
					"networkDiscovery": map[string]any{"count": 2},
				},
			},
		})
	})
	c := testClient(t, handler)
	counts, err := c.XDRAssetCounts(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if counts.Categories.Device.Count != 10 {
		t.Fatalf("expected device=10, got %d", counts.Categories.Device.Count)
	}
	if counts.Categories.Server.Count != 5 {
		t.Fatalf("expected server=5, got %d", counts.Categories.Server.Count)
	}
	if counts.Categories.Workstation.Count != 3 {
		t.Fatalf("expected workstation=3, got %d", counts.Categories.Workstation.Count)
	}
	if counts.Surfaces.Cloud.Count != 8 {
		t.Fatalf("expected cloud=8, got %d", counts.Surfaces.Cloud.Count)
	}
	if counts.Surfaces.Endpoint.Count != 7 {
		t.Fatalf("expected endpoint=7, got %d", counts.Surfaces.Endpoint.Count)
	}
	if counts.Surfaces.NetworkDiscovery.Count != 2 {
		t.Fatalf("expected networkDiscovery=2, got %d", counts.Surfaces.NetworkDiscovery.Count)
	}
	if counts.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestXDRAssetCountsParams(t *testing.T) {
	var gotQuery string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"categories": map[string]any{},
				"surfaces":   map[string]any{},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.XDRAssetCounts(context.Background(), &XDRAssetCountsParams{
		SiteIDs:    []string{"100", "200"},
		AccountIDs: []string{"300"},
		GroupIDs:   []string{"400"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"siteIds=100", "siteIds=200", "accountIds=300", "groupIds=400"} {
		if !strings.Contains(gotQuery, want) {
			t.Errorf("query %q missing %q", gotQuery, want)
		}
	}
}

func TestXDRAssetCountsError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden", "detail": "Insufficient permissions"},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.XDRAssetCounts(context.Background(), nil)
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

func TestXDRAssetCategories(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/xdr/assets/categories" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"device":      12,
				"server":      8,
				"workstation": 4,
				"container":   3,
				"identity":    2,
				"account":     1,
				"inventory":   6,
				"storage":     5,
			},
		})
	})
	c := testClient(t, handler)
	cat, err := c.XDRAssetCategories(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.Device != 12 {
		t.Fatalf("expected device=12, got %d", cat.Device)
	}
	if cat.Server != 8 {
		t.Fatalf("expected server=8, got %d", cat.Server)
	}
	if cat.Workstation != 4 {
		t.Fatalf("expected workstation=4, got %d", cat.Workstation)
	}
	if cat.Inventory != 6 {
		t.Fatalf("expected inventory=6, got %d", cat.Inventory)
	}
	if cat.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestXDRAssetCategoriesParams(t *testing.T) {
	var gotQuery string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{},
		})
	})
	c := testClient(t, handler)
	_, err := c.XDRAssetCategories(context.Background(), &XDRAssetCountsParams{
		SiteIDs: []string{"500"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotQuery, "siteIds=500") {
		t.Errorf("query %q missing siteIds=500", gotQuery)
	}
}

func TestXDRAssetCategoriesError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	})
	c := testClient(t, handler)
	_, err := c.XDRAssetCategories(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 500 {
		t.Fatalf("expected 500, got %d", ae.Status)
	}
}
