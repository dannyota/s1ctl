package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestIOCSeverityString(t *testing.T) {
	tests := []struct {
		sev  IOCSeverity
		want string
	}{
		{IOCSeverityUnknown, "Unknown"},
		{IOCSeverityInformational, "Informational"},
		{IOCSeverityLow, "Low"},
		{IOCSeverityMedium, "Medium"},
		{IOCSeverityHigh, "High"},
		{IOCSeverityCritical, "Critical"},
		{IOCSeverityFatal, "Fatal"},
		{IOCSeverity(7), "7"},
		{IOCSeverity(99), "99"},
	}
	for _, tt := range tests {
		if got := tt.sev.String(); got != tt.want {
			t.Fatalf("IOCSeverity(%d).String() = %q, want %q", int(tt.sev), got, tt.want)
		}
	}
}

func TestIOCsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/threat-intelligence/iocs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if got := q["severity"]; !slices.Equal(got, []string{"3", "4"}) {
			t.Fatalf("unexpected severity: %v", got)
		}
		if q.Get("type") != "IPV4" {
			t.Fatalf("unexpected type: %s", q.Get("type"))
		}
		if got := q["source"]; !slices.Equal(got, []string{"TestFeed"}) {
			t.Fatalf("unexpected source: %v", got)
		}
		if q.Get("value") != "10.0.0.1" {
			t.Fatalf("unexpected value: %s", q.Get("value"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"uuid": "00000000-0000-0000-0000-000000000001",
					"type": "IPV4", "value": "10.0.0.1",
					"source": "TestFeed", "severity": 4,
					"method": "EQUALS", "name": "Test IOC",
					"scope": "site", "scopeId": "225494730938493804",
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	iocs, pag, err := c.IOCsList(context.Background(), &IOCListParams{
		SiteIDs:    []string{"225494730938493804"},
		Severities: []IOCSeverity{IOCSeverityMedium, IOCSeverityHigh},
		Type:       IOCTypeIPv4,
		Sources:    []string{"TestFeed"},
		Value:      "10.0.0.1",
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(iocs) != 1 {
		t.Fatalf("expected 1 IOC, got %d", len(iocs))
	}
	ioc := iocs[0]
	if ioc.UUID != "00000000-0000-0000-0000-000000000001" {
		t.Fatalf("unexpected uuid: %s", ioc.UUID)
	}
	if ioc.Type != IOCTypeIPv4 {
		t.Fatalf("unexpected type: %s", ioc.Type)
	}
	if ioc.Severity != IOCSeverityHigh {
		t.Fatalf("unexpected severity: %d", ioc.Severity)
	}
	if ioc.Scope != IOCScopeSite {
		t.Fatalf("unexpected scope: %s", ioc.Scope)
	}
	if ioc.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestIOCsListNilParams(t *testing.T) {
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
	iocs, _, err := c.IOCsList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(iocs) != 0 {
		t.Fatalf("expected 0 IOCs, got %d", len(iocs))
	}
}

func TestIOCsCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/threat-intelligence/iocs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data []struct {
				Type   string `json:"type"`
				Value  string `json:"value"`
				Source string `json:"source"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if len(body.Data) != 1 {
			t.Fatalf("expected 1 IOC in request, got %d", len(body.Data))
		}
		if body.Data[0].Type != "SHA256" {
			t.Fatalf("unexpected type: %s", body.Data[0].Type)
		}
		if body.Data[0].Source != "TestSource" {
			t.Fatalf("unexpected source: %s", body.Data[0].Source)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"uuid": "00000000-0000-0000-0000-000000000002",
					"type": "SHA256", "value": "abcdef1234567890",
					"source": "TestSource", "severity": 3,
				},
			},
		})
	})
	c := testClient(t, handler)
	sev := IOCSeverityMedium
	created, err := c.IOCsCreate(context.Background(), []IOCCreateInput{
		{
			Type:     IOCTypeSHA256,
			Value:    "abcdef1234567890",
			Source:   "TestSource",
			Severity: &sev,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(created) != 1 {
		t.Fatalf("expected 1 created IOC, got %d", len(created))
	}
	if created[0].UUID != "00000000-0000-0000-0000-000000000002" {
		t.Fatalf("unexpected uuid: %s", created[0].UUID)
	}
	if created[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestIOCsCreateEmptyInput(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	_, err := c.IOCsCreate(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for empty IOC list")
	}
	_, err = c.IOCsCreate(context.Background(), []IOCCreateInput{})
	if err == nil {
		t.Fatal("expected error for empty IOC list")
	}
}

func TestIOCsDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/threat-intelligence/iocs" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("unexpected content type: %s", ct)
		}
		var body struct {
			Filter struct {
				UUIDs []string `json:"uuids"`
			} `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !slices.Equal(body.Filter.UUIDs, []string{"00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"}) {
			t.Fatalf("unexpected uuids: %v", body.Filter.UUIDs)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 2},
		})
	})
	c := testClient(t, handler)
	affected, err := c.IOCsDelete(context.Background(), []string{
		"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 affected, got %d", affected)
	}
}

func TestIOCsDeleteEmptyInput(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	_, err := c.IOCsDelete(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for empty UUID list")
	}
	_, err = c.IOCsDelete(context.Background(), []string{})
	if err == nil {
		t.Fatal("expected error for empty UUID list")
	}
}

func TestThreatIntelConfigs(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/threat-intelligence/user-config" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"scopeId":           "225494730938493804",
					"scopeLevel":        "site",
					"threatMinScore":    50,
					"disableRh":         false,
					"disableThreat":     false,
					"enableXdrMatching": true,
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	configs, err := c.ThreatIntelConfigs(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}
	cfg := configs[0]
	if cfg.ScopeID != "225494730938493804" {
		t.Fatalf("unexpected scopeId: %s", cfg.ScopeID)
	}
	if cfg.ScopeLevel != IOCScopeSite {
		t.Fatalf("unexpected scopeLevel: %s", cfg.ScopeLevel)
	}
	if cfg.ThreatMinScore != 50 {
		t.Fatalf("unexpected threatMinScore: %d", cfg.ThreatMinScore)
	}
	if !cfg.EnableXDRMatching {
		t.Fatal("expected enableXdrMatching to be true")
	}
	if cfg.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestIOCsListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.IOCsList(context.Background(), nil)
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
