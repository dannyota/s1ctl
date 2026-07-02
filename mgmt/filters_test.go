package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"testing"
)

func TestFiltersList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/filters" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("query") != "infected" {
			t.Fatalf("unexpected query: %s", q.Get("query"))
		}
		if got := q["siteIds"]; !slices.Equal(got, []string{"100000000000000001"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":           "200000000000000001",
					"name":         "Infected endpoints",
					"scopeId":      "100000000000000001",
					"scopeLevel":   "site",
					"filterFields": map[string]any{"infected": true},
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": nil},
		})
	})
	c := testClient(t, handler)
	filters, pag, err := c.FiltersList(context.Background(), &FilterListParams{
		Query:   "infected",
		SiteIDs: []string{"100000000000000001"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(filters))
	}
	if filters[0].Name != "Infected endpoints" {
		t.Fatalf("unexpected name: %s", filters[0].Name)
	}
	if len(filters[0].FilterFields) == 0 {
		t.Fatal("expected filterFields to be populated")
	}
	if filters[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestFiltersCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/filters" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				Name         string          `json:"name"`
				FilterFields json.RawMessage `json:"filterFields"`
			} `json:"data"`
			Filter struct {
				SiteIDs []string `json:"siteIds"`
			} `json:"filter"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data.Name != "Infected endpoints" {
			t.Fatalf("unexpected name: %s", body.Data.Name)
		}
		if !slices.Equal(body.Filter.SiteIDs, []string{"100000000000000001"}) {
			t.Fatalf("unexpected siteIds: %v", body.Filter.SiteIDs)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "200000000000000002", "name": "Infected endpoints"},
		})
	})
	c := testClient(t, handler)
	f, err := c.FiltersCreate(context.Background(), FilterCreate{
		Data:   FilterData{Name: "Infected endpoints", FilterFields: json.RawMessage(`{"infected":true}`)},
		Filter: &FilterScope{SiteIDs: []string{"100000000000000001"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.ID != "200000000000000002" {
		t.Fatalf("unexpected ID: %s", f.ID)
	}
}

func TestFiltersUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/filters/200000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				Name string `json:"name"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data.Name != "Renamed" {
			t.Fatalf("unexpected name: %s", body.Data.Name)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "200000000000000002", "name": "Renamed"},
		})
	})
	c := testClient(t, handler)
	f, err := c.FiltersUpdate(context.Background(), "200000000000000002", FilterUpdate{
		Data: FilterData{Name: "Renamed"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Name != "Renamed" {
		t.Fatalf("unexpected name: %s", f.Name)
	}
}

func TestFiltersDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/filters/200000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"success": true}})
	})
	c := testClient(t, handler)
	if err := c.FiltersDelete(context.Background(), "200000000000000002"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
