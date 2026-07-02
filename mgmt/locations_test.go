package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"testing"
)

func TestLocationsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/locations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":          "300000000000000001",
					"name":        "HQ Office",
					"operator":    "any",
					"isFallback":  false,
					"ipAddresses": map[string]any{"enabled": true},
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": nil},
		})
	})
	c := testClient(t, handler)
	locs, pag, err := c.LocationsList(context.Background(), &LocationListParams{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locs) != 1 {
		t.Fatalf("expected 1 location, got %d", len(locs))
	}
	if locs[0].Name != "HQ Office" {
		t.Fatalf("unexpected name: %s", locs[0].Name)
	}
	if locs[0].Operator != LocationOperatorAny {
		t.Fatalf("unexpected operator: %s", locs[0].Operator)
	}
	if len(locs[0].IPAddresses) == 0 {
		t.Fatal("expected ipAddresses to be populated")
	}
	if locs[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestLocationsCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/locations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				Name     string `json:"name"`
				Operator string `json:"operator"`
			} `json:"data"`
			Filter struct {
				SiteIDs []string `json:"siteIds"`
			} `json:"filter"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data.Name != "HQ Office" {
			t.Fatalf("unexpected name: %s", body.Data.Name)
		}
		if body.Data.Operator != "any" {
			t.Fatalf("unexpected operator: %s", body.Data.Operator)
		}
		if !slices.Equal(body.Filter.SiteIDs, []string{"100000000000000001"}) {
			t.Fatalf("unexpected siteIds: %v", body.Filter.SiteIDs)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "300000000000000002", "name": "HQ Office"},
		})
	})
	c := testClient(t, handler)
	loc, err := c.LocationsCreate(context.Background(), LocationCreate{
		Data:   LocationData{Name: "HQ Office", Operator: LocationOperatorAny},
		Filter: LocationScope{SiteIDs: []string{"100000000000000001"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.ID != "300000000000000002" {
		t.Fatalf("unexpected ID: %s", loc.ID)
	}
}

func TestLocationsUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/locations/300000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
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
			"data": map[string]any{"id": "300000000000000002", "description": "updated"},
		})
	})
	c := testClient(t, handler)
	loc, err := c.LocationsUpdate(context.Background(), "300000000000000002", LocationUpdate{
		Data: LocationData{Name: "HQ Office", Description: "updated", Operator: LocationOperatorAny},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.ID != "300000000000000002" {
		t.Fatalf("unexpected ID: %s", loc.ID)
	}
}

func TestLocationsDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/locations" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				IDs []string `json:"ids"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if !slices.Equal(body.Data.IDs, []string{"300000000000000002", "300000000000000003"}) {
			t.Fatalf("unexpected ids: %v", body.Data.IDs)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 2}})
	})
	c := testClient(t, handler)
	if err := c.LocationsDelete(context.Background(), []string{"300000000000000002", "300000000000000003"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLocationOperatorEnumValues(t *testing.T) {
	if string(LocationOperatorAll) != "all" || string(LocationOperatorAny) != "any" || string(LocationOperatorNone) != "none" {
		t.Fatal("unexpected location operator enum values")
	}
}
