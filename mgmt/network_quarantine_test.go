package mgmt

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"testing"
)

// The network quarantine surface is the firewall-control endpoint family under
// the "network-quarantine" category segment. Each test asserts the segment is
// present and the request body matches the shared firewall shapes.

func TestNetworkQuarantineList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{{"id": "1", "name": "NQ Rule"}},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	rules, pag, err := c.NetworkQuarantineList(context.Background(), &FirewallRuleListParams{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Name != "NQ Rule" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestNetworkQuarantineGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall-control/network-quarantine" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ids"); got != "1000000000000000001" {
			t.Fatalf("unexpected ids query: %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{"id": "1000000000000000001", "name": "NQ Rule"}},
		})
	})
	c := testClient(t, handler)
	rule, err := c.NetworkQuarantineGet(context.Background(), "1000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.ID != "1000000000000000001" {
		t.Fatalf("unexpected rule: %+v", rule)
	}
}

func TestNetworkQuarantineCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter struct {
				SiteIDs []string `json:"siteIds"`
			} `json:"filter"`
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.SiteIDs, []string{"225494730938493804"}) {
			t.Fatalf("unexpected filter siteIds: %v", body.Filter.SiteIDs)
		}
		if body.Data["name"] != "NQ Rule" {
			t.Fatalf("unexpected data name: %v", body.Data["name"])
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": "1", "name": "NQ Rule"}})
	})
	c := testClient(t, handler)
	rule, err := c.NetworkQuarantineCreate(context.Background(),
		FirewallRuleScope{SiteIDs: []string{"225494730938493804"}},
		FirewallRuleCreate{Name: "NQ Rule", Action: FirewallActionBlock, Status: FirewallStatusEnabled, Direction: FirewallDirectionAny})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Name != "NQ Rule" {
		t.Fatalf("unexpected rule: %+v", rule)
	}
}

func TestNetworkQuarantineUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"id": "1000000000000000001", "name": "NQ Rule"}})
	})
	c := testClient(t, handler)
	_, err := c.NetworkQuarantineUpdate(context.Background(), "1000000000000000001",
		FirewallRuleCreate{Name: "NQ Rule", Action: FirewallActionBlock, Status: FirewallStatusEnabled, Direction: FirewallDirectionAny})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNetworkQuarantineDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 2}})
	})
	c := testClient(t, handler)
	affected, err := c.NetworkQuarantineDelete(context.Background(), []string{"1", "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 affected, got %d", affected)
	}
}

func TestNetworkQuarantineSetStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall-control/network-quarantine/enable" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data.Status != "Enabled" {
			t.Fatalf("unexpected status: %s", body.Data.Status)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 1}})
	})
	c := testClient(t, handler)
	affected, err := c.NetworkQuarantineSetStatus(context.Background(), []string{"1"}, FirewallStatusEnabled)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestNetworkQuarantineReorder(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall-control/network-quarantine/reorder" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"success": true}})
	})
	c := testClient(t, handler)
	tenant := true
	err := c.NetworkQuarantineReorder(context.Background(),
		[]RuleOrder{{ID: "1", Order: 1}}, FirewallRuleReorderFilter{Tenant: &tenant})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNetworkQuarantineCopy(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall-control/network-quarantine/copy-rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 3}})
	})
	c := testClient(t, handler)
	group := "225494730938493904"
	affected, err := c.NetworkQuarantineCopy(context.Background(),
		FirewallRuleReorderFilter{SiteIDs: []string{"225494730938493804"}},
		[]FirewallRuleCopyTarget{{GroupID: &group}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 3 {
		t.Fatalf("expected 3 affected, got %d", affected)
	}
}

func TestNetworkQuarantineProtocolsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall-control/network-quarantine/protocols" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{{"value": "tcp", "name": "TCP"}},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	protos, _, err := c.NetworkQuarantineProtocolsList(context.Background(), &FirewallProtocolListParams{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(protos) != 1 || protos[0].Value != "tcp" {
		t.Fatalf("unexpected protocols: %+v", protos)
	}
}

func TestNetworkQuarantineExport(t *testing.T) {
	raw := []byte(`[{"name":"NQ Rule"}]`)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall-control/network-quarantine/export" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Write(raw)
	})
	c := testClient(t, handler)
	data, err := c.NetworkQuarantineExport(context.Background(), &FirewallRuleListParams{SiteIDs: []string{"225494730938493804"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, raw) {
		t.Fatalf("unexpected export payload: %s", data)
	}
}

func TestNetworkQuarantineImport(t *testing.T) {
	fileData := []byte(`[{"name":"NQ Rule"}]`)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine/import" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if got := r.MultipartForm.Value["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"success": true}})
	})
	c := testClient(t, handler)
	err := c.NetworkQuarantineImport(context.Background(),
		FirewallImportScope{SiteIDs: []string{"225494730938493804"}}, "rules.json", fileData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNetworkQuarantineConfigurationGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine/configuration" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("siteIds"); got != "225494730938493804" {
			t.Fatalf("unexpected siteIds query: %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"enabled": true, "locationAware": false, "selectedTags": []string{"t1"}},
		})
	})
	c := testClient(t, handler)
	cfg, err := c.NetworkQuarantineConfigurationGet(context.Background(),
		FirewallConfigScope{SiteIDs: []string{"225494730938493804"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled || len(cfg.SelectedTags) != 1 || cfg.SelectedTags[0] != "t1" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestNetworkQuarantineConfigurationUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine/configuration" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter struct {
				SiteIDs []string `json:"siteIds"`
			} `json:"filter"`
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.SiteIDs, []string{"225494730938493804"}) {
			t.Fatalf("unexpected filter siteIds: %v", body.Filter.SiteIDs)
		}
		if body.Data["enabled"] != true {
			t.Fatalf("unexpected data enabled: %v", body.Data["enabled"])
		}
		if _, ok := body.Data["reportBlocked"]; ok {
			t.Fatalf("reportBlocked should be omitted when unset")
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"enabled": true}})
	})
	c := testClient(t, handler)
	enabled := true
	cfg, err := c.NetworkQuarantineConfigurationUpdate(context.Background(),
		FirewallConfigScope{SiteIDs: []string{"225494730938493804"}},
		FirewallConfigurationUpdate{Enabled: &enabled})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestNetworkQuarantineSetLocation(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine/set-location" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter struct {
				IDs []string `json:"ids"`
			} `json:"filter"`
			Data struct {
				Type   string `json:"type"`
				Values []struct {
					ID string `json:"id"`
				} `json:"values"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.IDs, []string{"1"}) {
			t.Fatalf("unexpected filter ids: %v", body.Filter.IDs)
		}
		if body.Data.Type != "specific" || len(body.Data.Values) != 1 || body.Data.Values[0].ID != "loc1" {
			t.Fatalf("unexpected data: %+v", body.Data)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 1}})
	})
	c := testClient(t, handler)
	affected, err := c.NetworkQuarantineSetLocation(context.Background(),
		FirewallActionFilter{IDs: []string{"1"}},
		FirewallLocationTarget{Type: FirewallLocationSpecific, Values: []FirewallLocationValue{{ID: "loc1"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestNetworkQuarantineSetLocationEmptyFilter(t *testing.T) {
	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("no request should be sent for an empty filter")
	})
	c := testClient(t, handler)
	_, err := c.NetworkQuarantineSetLocation(context.Background(),
		FirewallActionFilter{}, FirewallLocationTarget{Type: FirewallLocationAll})
	if err == nil {
		t.Fatal("expected error for empty filter")
	}
}

func TestNetworkQuarantineMoveRules(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine/move-rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter struct {
				IDs []string `json:"ids"`
			} `json:"filter"`
			Data []map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.IDs, []string{"1"}) {
			t.Fatalf("unexpected filter ids: %v", body.Filter.IDs)
		}
		if len(body.Data) != 1 || body.Data[0]["siteId"] != "225494730938493805" {
			t.Fatalf("unexpected targets: %v", body.Data)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 1}})
	})
	c := testClient(t, handler)
	site := "225494730938493805"
	affected, err := c.NetworkQuarantineMoveRules(context.Background(),
		FirewallActionFilter{IDs: []string{"1"}},
		[]FirewallRuleCopyTarget{{SiteID: &site}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestNetworkQuarantineAddTags(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/network-quarantine/add-tags" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				TagIDs []string `json:"tagIds"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Data.TagIDs, []string{"t1", "t2"}) {
			t.Fatalf("unexpected tagIds: %v", body.Data.TagIDs)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 2}})
	})
	c := testClient(t, handler)
	affected, err := c.NetworkQuarantineAddTags(context.Background(),
		FirewallActionFilter{IDs: []string{"1"}}, []string{"t1", "t2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 affected, got %d", affected)
	}
}

func TestNetworkQuarantineRemoveTags(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall-control/network-quarantine/remove-tags" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 1}})
	})
	c := testClient(t, handler)
	affected, err := c.NetworkQuarantineRemoveTags(context.Background(),
		FirewallActionFilter{IDs: []string{"1"}}, []string{"t1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}
