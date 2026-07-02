package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestPlatformRulesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/detection-library/platform-rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["platformRuleIds"]; !slices.Equal(got, []string{"1000000000000000001", "1000000000000000002"}) {
			t.Fatalf("unexpected platformRuleIds: %v", got)
		}
		if q.Get("scopeLevel") != "site" {
			t.Fatalf("unexpected scopeLevel: %s", q.Get("scopeLevel"))
		}
		if q.Get("ruleNameSubstring") != "lateral" {
			t.Fatalf("unexpected ruleNameSubstring: %s", q.Get("ruleNameSubstring"))
		}
		if got := q["severities"]; !slices.Equal(got, []string{"High", "Critical"}) {
			t.Fatalf("unexpected severities: %v", got)
		}
		if q.Get("limit") != "50" {
			t.Fatalf("expected limit=50, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "1000000000000000001", "name": "Example Rule",
					"severity": "High", "status": "Active",
					"scopeLevel": "site", "queryType": "events",
					"treatAsThreat": "Suspicious", "generatedAlerts": 3,
					"mitre": []map[string]any{
						{
							"tactic": "Execution",
							"techniques": []map[string]any{
								{"id": "T1059", "title": "Command and Scripting Interpreter"},
							},
						},
					},
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "cur"},
		})
	})
	c := testClient(t, handler)
	rules, pag, err := c.PlatformRulesList(context.Background(), &PlatformRuleListParams{
		IDs:          []string{"1000000000000000001", "1000000000000000002"},
		ScopeLevel:   "site",
		NameContains: "lateral",
		Severities:   []string{"High", "Critical"},
		Limit:        50,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.Severity != PlatformRuleSeverityHigh {
		t.Fatalf("unexpected severity: %s", r.Severity)
	}
	if r.Status != PlatformRuleStatusActive {
		t.Fatalf("unexpected status: %s", r.Status)
	}
	if r.TreatAsThreat != RuleTreatSuspicious {
		t.Fatalf("unexpected treatAsThreat: %s", r.TreatAsThreat)
	}
	if len(r.Mitre) != 1 || r.Mitre[0].Tactic != "Execution" {
		t.Fatalf("unexpected mitre: %+v", r.Mitre)
	}
	if len(r.Mitre[0].Techniques) != 1 || r.Mitre[0].Techniques[0].ID != "T1059" {
		t.Fatalf("unexpected techniques: %+v", r.Mitre[0].Techniques)
	}
	if r.Raw == nil || r.Mitre[0].Raw == nil || r.Mitre[0].Techniques[0].Raw == nil {
		t.Fatal("expected Raw to be populated on rule and nested mitre structs")
	}
	if pag.NextCursor != "cur" {
		t.Fatalf("unexpected cursor: %s", pag.NextCursor)
	}
}

func TestPlatformRulesEnable(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/detection-library/platform-rules/enable" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		ids, _ := body["platformRuleIds"].([]any)
		if len(ids) != 2 || ids[0] != "1000000000000000001" {
			t.Fatalf("unexpected platformRuleIds: %v", body["platformRuleIds"])
		}
		if body["scopeId"] != "225494730938493804" {
			t.Fatalf("unexpected scopeId: %v", body["scopeId"])
		}
		if body["scopeLevel"] != "site" {
			t.Fatalf("unexpected scopeLevel: %v", body["scopeLevel"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 2},
		})
	})
	c := testClient(t, handler)
	affected, err := c.PlatformRulesEnable(context.Background(), PlatformRuleActionFilter{
		PlatformRuleIDs: []string{"1000000000000000001", "1000000000000000002"},
		ScopeID:         "225494730938493804",
		ScopeLevel:      "site",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 2 {
		t.Fatalf("expected 2 affected, got %d", affected)
	}
}

func TestPlatformRulesDisable(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/detection-library/platform-rules/disable" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 1},
		})
	})
	c := testClient(t, handler)
	affected, err := c.PlatformRulesDisable(context.Background(), PlatformRuleActionFilter{
		PlatformRuleIDs: []string{"1000000000000000001"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

// Platform rule enable/disable must refuse an ID-less filter: scopeId and
// scopeLevel only narrow the selection, so an empty ID list would toggle
// every rule in the scope.
func TestPlatformRuleActionRequiresRuleIDs(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	filter := PlatformRuleActionFilter{ScopeID: "225494730938493804", ScopeLevel: "site"}
	tests := []struct {
		name string
		call func() (int, error)
	}{
		{"enable", func() (int, error) { return c.PlatformRulesEnable(context.Background(), filter) }},
		{"disable", func() (int, error) { return c.PlatformRulesDisable(context.Background(), filter) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.call(); err == nil {
				t.Fatal("expected error for empty rule ID list")
			}
		})
	}
}

func TestDetectionSurfacesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/detection-library/surfaces" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"surfaces": []map[string]any{
					{"key": "endpoint", "value": "Endpoint"},
					{"key": "identity", "value": "Identity"},
				},
			},
		})
	})
	c := testClient(t, handler)
	surfaces, err := c.DetectionSurfacesList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(surfaces) != 2 {
		t.Fatalf("expected 2 surfaces, got %d", len(surfaces))
	}
	if surfaces[0].Key != "endpoint" || surfaces[0].Value != "Endpoint" {
		t.Fatalf("unexpected surface: %+v", surfaces[0])
	}
	if surfaces[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestDetectionDataSourcesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/detection-library/data-sources" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"dataSources": []map[string]any{
					{"key": "edr", "value": "EDR"},
				},
			},
		})
	})
	c := testClient(t, handler)
	sources, err := c.DetectionDataSourcesList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sources) != 1 || sources[0].Key != "edr" {
		t.Fatalf("unexpected data sources: %+v", sources)
	}
}

func TestPlatformRulesListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.PlatformRulesList(context.Background(), nil)
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
