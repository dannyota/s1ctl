package mgmt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"testing"
)

func TestFirewallProtocolsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/protocols" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("query") != "tcp" {
			t.Fatalf("unexpected query: %s", q.Get("query"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"value": "tcp", "name": "TCP"},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	protos, pag, err := c.FirewallProtocolsList(context.Background(), &FirewallProtocolListParams{Query: "tcp", Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(protos) != 1 || protos[0].Value != "tcp" || protos[0].Name != "TCP" {
		t.Fatalf("unexpected protocols: %+v", protos)
	}
	if protos[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestFirewallRulesDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter struct {
				IDs []string `json:"ids"`
			} `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.IDs, []string{"1000000000000000001"}) {
			t.Fatalf("unexpected ids: %v", body.Filter.IDs)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 1},
		})
	})
	c := testClient(t, handler)
	affected, err := c.FirewallRulesDelete(context.Background(), []string{"1000000000000000001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestFirewallRulesReorder(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/reorder" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   []map[string]any `json:"data"`
			Filter map[string]any   `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if len(body.Data) != 1 || body.Data[0]["id"] != "1000000000000000001" || body.Data[0]["order"] != float64(5) {
			t.Fatalf("unexpected data: %v", body.Data)
		}
		if body.Filter["tenant"] != true {
			t.Fatalf("unexpected filter tenant: %v", body.Filter["tenant"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true},
		})
	})
	c := testClient(t, handler)
	tenant := true
	err := c.FirewallRulesReorder(context.Background(),
		[]RuleOrder{{ID: "1000000000000000001", Order: 5}},
		FirewallRuleReorderFilter{Tenant: &tenant})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFirewallRulesCopy(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/copy-rules" {
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
		if len(body.Data) != 1 || body.Data[0]["groupId"] != "225494730938493904" {
			t.Fatalf("unexpected targets: %v", body.Data)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 4},
		})
	})
	c := testClient(t, handler)
	group := "225494730938493904"
	affected, err := c.FirewallRulesCopy(context.Background(),
		FirewallRuleReorderFilter{SiteIDs: []string{"225494730938493804"}},
		[]FirewallRuleCopyTarget{{GroupID: &group}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 4 {
		t.Fatalf("expected 4 affected, got %d", affected)
	}
}

func TestFirewallRulesSetStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/enable" {
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
		if body.Data.Status != "Disabled" {
			t.Fatalf("unexpected status: %s", body.Data.Status)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 1},
		})
	})
	c := testClient(t, handler)
	affected, err := c.FirewallRulesSetStatus(context.Background(),
		[]string{"1000000000000000001"}, FirewallStatusDisabled)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestFirewallRulesExport(t *testing.T) {
	raw := []byte(`[{"name":"Example Rule","action":"Allow"}]`)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/export" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		w.Write(raw)
	})
	c := testClient(t, handler)
	data, err := c.FirewallRulesExport(context.Background(), &FirewallRuleListParams{
		SiteIDs: []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(data, raw) {
		t.Fatalf("unexpected export payload: %s", data)
	}
}

func TestFirewallRulesExportError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden"},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.FirewallRulesExport(context.Background(), nil)
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

func TestFirewallRulesImport(t *testing.T) {
	fileData := []byte(`[{"name":"Example Rule"}]`)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/firewall-control/import" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if got := r.MultipartForm.Value["siteIds"]; !slices.Equal(got, []string{"225494730938493804", "225494730938493805"}) {
			t.Fatalf("unexpected siteIds fields: %v", got)
		}
		if got := r.MultipartForm.Value["tenant"]; !slices.Equal(got, []string{"true"}) {
			t.Fatalf("unexpected tenant field: %v", got)
		}
		file, hdr, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("form file: %v", err)
		}
		defer file.Close()
		if hdr.Filename != "rules.json" {
			t.Fatalf("unexpected filename: %s", hdr.Filename)
		}
		got, _ := io.ReadAll(file)
		if !bytes.Equal(got, fileData) {
			t.Fatalf("unexpected file content: %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"success": true}})
	})
	c := testClient(t, handler)
	err := c.FirewallRulesImport(context.Background(), FirewallImportScope{
		SiteIDs: []string{"225494730938493804", "225494730938493805"},
		Tenant:  true,
	}, "rules.json", fileData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFirewallRulesImportError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 400, "title": "Bad Request", "detail": "invalid file"},
			},
		})
	})
	c := testClient(t, handler)
	err := c.FirewallRulesImport(context.Background(), FirewallImportScope{
		SiteIDs: []string{"225494730938493804"},
	}, "rules.json", []byte("{}"))
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 400 {
		t.Fatalf("expected 400, got %d", ae.Status)
	}
}
