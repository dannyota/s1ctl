package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestDeviceEventsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/device-control/events" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if got := q["interfaces"]; !slices.Equal(got, []string{"USB", "Bluetooth"}) {
			t.Fatalf("unexpected interfaces: %v", got)
		}
		if q.Get("query") != "example" {
			t.Fatalf("unexpected query: %s", q.Get("query"))
		}
		if q.Get("limit") != "25" {
			t.Fatalf("expected limit=25, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "1000000000000000000", "eventType": "connected",
					"interface": "USB", "deviceName": "Example Drive",
					"vendorId": "1234", "productId": "5678",
					"computerName": "host-1",
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	events, pag, err := c.DeviceEventsList(context.Background(), &DeviceEventListParams{
		SiteIDs:    []string{"225494730938493804"},
		Interfaces: []string{"USB", "Bluetooth"},
		Query:      "example",
		Limit:      25,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].EventType != "connected" {
		t.Fatalf("unexpected eventType: %s", events[0].EventType)
	}
	if events[0].DeviceName != "Example Drive" {
		t.Fatalf("unexpected deviceName: %s", events[0].DeviceName)
	}
	if events[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestDeviceRulesDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/device-control" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("unexpected content type: %s", ct)
		}
		var body struct {
			Filter struct {
				IDs []string `json:"ids"`
			} `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.IDs, []string{"1000000000000000001", "1000000000000000002"}) {
			t.Fatalf("unexpected ids: %v", body.Filter.IDs)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 2},
		})
	})
	c := testClient(t, handler)
	affected, err := c.DeviceRulesDelete(context.Background(), []string{"1000000000000000001", "1000000000000000002"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 affected, got %d", affected)
	}
}

func TestDeviceRulesReorder(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/device-control/reorder" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   []map[string]any `json:"data"`
			Filter map[string]any   `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if len(body.Data) != 2 {
			t.Fatalf("expected 2 orders, got %d", len(body.Data))
		}
		if body.Data[0]["id"] != "1000000000000000001" || body.Data[0]["order"] != float64(1) {
			t.Fatalf("unexpected first order: %v", body.Data[0])
		}
		siteIDs, _ := body.Filter["siteIds"].([]any)
		if len(siteIDs) != 1 || siteIDs[0] != "225494730938493804" {
			t.Fatalf("unexpected filter siteIds: %v", body.Filter["siteIds"])
		}
		if body.Filter["interface"] != "USB" {
			t.Fatalf("unexpected filter interface: %v", body.Filter["interface"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true},
		})
	})
	c := testClient(t, handler)
	iface := DeviceRuleInterfaceUSB
	err := c.DeviceRulesReorder(context.Background(),
		[]RuleOrder{
			{ID: "1000000000000000001", Order: 1},
			{ID: "1000000000000000002", Order: 2},
		},
		DeviceRuleReorderFilter{
			SiteIDs:   []string{"225494730938493804"},
			Interface: &iface,
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeviceRulesCopy(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/device-control/copy-rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter map[string]any   `json:"filter"`
			Data   []map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		siteIDs, _ := body.Filter["siteIds"].([]any)
		if len(siteIDs) != 1 || siteIDs[0] != "225494730938493804" {
			t.Fatalf("unexpected filter siteIds: %v", body.Filter["siteIds"])
		}
		if len(body.Data) != 1 || body.Data[0]["siteId"] != "225494730938493805" {
			t.Fatalf("unexpected targets: %v", body.Data)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 3},
		})
	})
	c := testClient(t, handler)
	target := "225494730938493805"
	affected, err := c.DeviceRulesCopy(context.Background(),
		DeviceRuleScopeFilter{SiteIDs: []string{"225494730938493804"}},
		[]DeviceRuleCopyTarget{{SiteID: &target}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 3 {
		t.Fatalf("expected 3 affected, got %d", affected)
	}
}

func TestDeviceRulesSetStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/device-control/enable" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter struct {
				IDs []string `json:"ids"`
			} `json:"filter"`
			Data struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.IDs, []string{"1000000000000000001"}) {
			t.Fatalf("unexpected ids: %v", body.Filter.IDs)
		}
		if body.Data.Status != "Enabled" {
			t.Fatalf("unexpected status: %s", body.Data.Status)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 1},
		})
	})
	c := testClient(t, handler)
	affected, err := c.DeviceRulesSetStatus(context.Background(),
		[]string{"1000000000000000001"}, DeviceRuleStatusEnabled)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestDeviceRulesDeleteError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 404, "title": "Not Found"},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.DeviceRulesDelete(context.Background(), []string{"1000000000000000001"})
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 404 {
		t.Fatalf("expected 404, got %d", ae.Status)
	}
}
