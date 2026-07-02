package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestBlocklistList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/restrictions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if got := q["osTypes"]; !slices.Equal(got, []string{"windows", "linux"}) {
			t.Fatalf("unexpected osTypes: %v", got)
		}
		if got := q["types"]; !slices.Equal(got, []string{"black_hash"}) {
			t.Fatalf("unexpected types: %v", got)
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
					"id":          "1000000000000000001",
					"type":        "black_hash",
					"value":       "ffffffffffffffffffffffffffffffffffffffff",
					"sha256Value": "aaaa",
					"osType":      "windows",
					"source":      "user",
					"description": "malware hash",
					"scopeName":   "Default Site",
					"imported":    true,
					"createdAt":   "2025-01-01T00:00:00Z",
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "abc"},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.BlocklistList(context.Background(), &BlocklistListParams{
		SiteIDs: []string{"225494730938493804"},
		OSTypes: []string{"windows", "linux"},
		Types:   []string{"black_hash"},
		SortBy:  "createdAt",
		Limit:   10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	b := items[0]
	if b.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", b.ID)
	}
	if b.Value != "ffffffffffffffffffffffffffffffffffffffff" {
		t.Fatalf("unexpected value: %s", b.Value)
	}
	if b.SHA256Value != "aaaa" {
		t.Fatalf("unexpected sha256Value: %s", b.SHA256Value)
	}
	if b.OSType != "windows" {
		t.Fatalf("unexpected osType: %s", b.OSType)
	}
	if b.Type != "black_hash" {
		t.Fatalf("unexpected type: %s", b.Type)
	}
	if !b.Imported {
		t.Fatal("expected imported=true")
	}
	if b.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
	if pag.NextCursor != "abc" {
		t.Fatalf("unexpected cursor: %s", pag.NextCursor)
	}
}

func TestBlocklistListNilParams(t *testing.T) {
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
	items, _, err := c.BlocklistList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestBlocklistCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/restrictions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   map[string]any `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data["type"] != "black_hash" {
			t.Fatalf("unexpected type: %v", body.Data["type"])
		}
		if body.Data["osType"] != "windows" {
			t.Fatalf("unexpected osType: %v", body.Data["osType"])
		}
		if body.Data["value"] != "ffffffffffffffffffffffffffffffffffffffff" {
			t.Fatalf("unexpected value: %v", body.Data["value"])
		}
		if body.Data["description"] != "bad" {
			t.Fatalf("unexpected description: %v", body.Data["description"])
		}
		sites, _ := body.Filter["siteIds"].([]any)
		if len(sites) != 1 || sites[0] != "225494730938493804" {
			t.Fatalf("unexpected filter siteIds: %v", body.Filter["siteIds"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "1000000000000000002", "type": "black_hash", "value": "ffffffffffffffffffffffffffffffffffffffff", "osType": "windows"},
			},
		})
	})
	c := testClient(t, handler)
	item, err := c.BlocklistCreate(context.Background(),
		BlocklistScope{SiteIDs: []string{"225494730938493804"}},
		BlocklistCreate{
			Type:        BlocklistTypeBlackHash,
			OSType:      BlocklistOSWindows,
			Value:       "ffffffffffffffffffffffffffffffffffffffff",
			Description: "bad",
		},
	)
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

func TestBlocklistCreateTenantScope(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Filter["tenant"] != true {
			t.Fatalf("expected tenant=true, got %v", body.Filter["tenant"])
		}
		if _, ok := body.Filter["siteIds"]; ok {
			t.Fatal("expected siteIds to be omitted for tenant scope")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{"id": "1000000000000000003"}},
		})
	})
	c := testClient(t, handler)
	item, err := c.BlocklistCreate(context.Background(),
		BlocklistScope{Tenant: true},
		BlocklistCreate{Type: BlocklistTypeBlackHash, OSType: BlocklistOSLinux, Value: "ff"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != "1000000000000000003" {
		t.Fatalf("unexpected ID: %s", item.ID)
	}
}

func TestBlocklistUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/restrictions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["id"] != "1000000000000000002" {
			t.Fatalf("unexpected id: %v", body.Data["id"])
		}
		if body.Data["type"] != "black_hash" {
			t.Fatalf("unexpected type: %v", body.Data["type"])
		}
		if body.Data["osType"] != "linux" {
			t.Fatalf("unexpected osType: %v", body.Data["osType"])
		}
		if body.Data["description"] != "updated" {
			t.Fatalf("unexpected description: %v", body.Data["description"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{"id": "1000000000000000002", "description": "updated"}},
		})
	})
	c := testClient(t, handler)
	item, err := c.BlocklistUpdate(context.Background(), "1000000000000000002",
		BlocklistScope{},
		BlocklistCreate{Type: BlocklistTypeBlackHash, OSType: BlocklistOSLinux, Value: "ff", Description: "updated"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", item.ID)
	}
}

func TestBlocklistDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/restrictions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				IDs []string `json:"ids"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Data.IDs, []string{"a", "b"}) {
			t.Fatalf("unexpected ids: %v", body.Data.IDs)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 2},
		})
	})
	c := testClient(t, handler)
	affected, err := c.BlocklistDelete(context.Background(), []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 affected, got %d", affected)
	}
}

func TestBlocklistValidate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/restrictions/validate" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["osType"] != "linux" {
			t.Fatalf("unexpected osType: %v", body.Data["osType"])
		}
		if body.Data["value"] != "ff" {
			t.Fatalf("unexpected value: %v", body.Data["value"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"status": "Not recommended",
				"details": []map[string]any{
					{"field": "value", "error": "decreases security"},
				},
			},
		})
	})
	c := testClient(t, handler)
	res, err := c.BlocklistValidate(context.Background(),
		BlocklistScope{},
		BlocklistValidateInput{OSType: BlocklistOSLinux, Value: "ff"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != BlocklistStatusNotRecommended {
		t.Fatalf("unexpected status: %s", res.Status)
	}
	if len(res.Details) != 1 || res.Details[0].Field != "value" {
		t.Fatalf("unexpected details: %v", res.Details)
	}
}

func TestBlocklistExport(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/export/restrictions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		w.Write([]byte("value,osType,type\nff,linux,black_hash\n"))
	})
	c := testClient(t, handler)
	data, err := c.BlocklistExport(context.Background(), &BlocklistListParams{
		SiteIDs: []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "value,osType,type\nff,linux,black_hash\n" {
		t.Fatalf("unexpected export body: %q", string(data))
	}
}

func TestBlocklistListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{{"code": 403, "title": "Forbidden"}},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.BlocklistList(context.Background(), nil)
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

func TestBlocklistEnumValues(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"OSLinux", string(BlocklistOSLinux), "linux"},
		{"OSMacOS", string(BlocklistOSMacOS), "macos"},
		{"OSWindows", string(BlocklistOSWindows), "windows"},
		{"OSWindowsLegacy", string(BlocklistOSWindowsLegacy), "windows_legacy"},
		{"TypeBlackHash", string(BlocklistTypeBlackHash), "black_hash"},
		{"StatusNotRecommended", string(BlocklistStatusNotRecommended), "Not recommended"},
		{"StatusNotAllowed", string(BlocklistStatusNotAllowed), "Not allowed"},
		{"StatusNone", string(BlocklistStatusNone), "NONE"},
		{"StatusDupSHA1", string(BlocklistStatusDuplicatedSHA1), "duplicated_value_sha1"},
		{"StatusDupSHA256", string(BlocklistStatusDuplicatedSHA256), "duplicated_value_sha256"},
		{"StatusDupBoth", string(BlocklistStatusDuplicatedSHA1SHA256), "duplicated_value_sha1_sha256"},
		{"StatusDuplication", string(BlocklistStatusDuplication), "Duplication"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}
