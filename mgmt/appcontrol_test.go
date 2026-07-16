package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestAppControlRulesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/nac/config/api/v1/nac/rules/query" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("pageSize") != "10" {
			t.Fatalf("expected pageSize=10, got %s", r.URL.Query().Get("pageSize"))
		}
		var body struct {
			ScopeSelector *struct {
				ScopeType string   `json:"scopeType"`
				ScopeIDs  []string `json:"scopeIds"`
			} `json:"scopeSelector"`
			IncludeParents bool `json:"includeParents"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.ScopeSelector == nil {
			t.Fatal("expected scopeSelector")
		}
		if body.ScopeSelector.ScopeType != "SITE" {
			t.Fatalf("expected SITE scope, got %s", body.ScopeSelector.ScopeType)
		}
		if !body.IncludeParents {
			t.Fatal("expected includeParents=true")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"pageInfo":   map[string]any{"endCursor": "c2", "hasNextPage": true},
			"totalCount": 42,
			"edges": []map[string]any{
				{
					"cursor": "c1",
					"node": map[string]any{
						"id":          "12345",
						"ruleName":    "Block unsigned apps",
						"description": "test rule",
						"behavior":    "BLOCK",
						"osType":      []string{"WINDOWS"},
						"propagation": true,
						"createdAt":   "2024-04-03T12:00:00Z",
						"parameters": map[string]any{
							"publisher": "Evil Corp",
						},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	rules, cursor, total, err := c.AppControlRulesList(context.Background(), &AppControlQueryParams{
		ScopeType:      AppControlScopeSite,
		ScopeIDs:       []string{"225494730938493804"},
		IncludeParents: true,
		PageSize:       10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].ID != "12345" {
		t.Fatalf("unexpected ID: %s", rules[0].ID)
	}
	if rules[0].RuleName != "Block unsigned apps" {
		t.Fatalf("unexpected name: %s", rules[0].RuleName)
	}
	if rules[0].Behavior != AppControlBehaviorBlock {
		t.Fatalf("unexpected behavior: %s", rules[0].Behavior)
	}
	if rules[0].Parameters == nil || rules[0].Parameters.Publisher != "Evil Corp" {
		t.Fatal("expected parameters.publisher to be set")
	}
	if rules[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if cursor != "c2" {
		t.Fatalf("expected cursor c2, got %s", cursor)
	}
	if total != 42 {
		t.Fatalf("expected total 42, got %d", total)
	}
}

func TestAppControlRulesListNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"pageInfo":   map[string]any{"endCursor": "", "hasNextPage": false},
			"totalCount": 0,
			"edges":      []any{},
		})
	})
	c := testClient(t, handler)
	rules, _, total, err := c.AppControlRulesList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(rules))
	}
	if total != 0 {
		t.Fatalf("expected 0 total, got %d", total)
	}
}

func TestAppControlRulesGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/nac/config/api/v1/nac/rules/12345" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":          "12345",
			"ruleName":    "Block unsigned apps",
			"behavior":    "BLOCK",
			"osType":      []string{"WINDOWS"},
			"propagation": false,
		})
	})
	c := testClient(t, handler)
	rule, err := c.AppControlRulesGet(context.Background(), "12345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.ID != "12345" {
		t.Fatalf("unexpected ID: %s", rule.ID)
	}
	if rule.RuleName != "Block unsigned apps" {
		t.Fatalf("unexpected name: %s", rule.RuleName)
	}
	if rule.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAppControlRulesCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/nac/config/api/v1/nac/rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["ruleName"] != "New rule" {
			t.Fatalf("unexpected ruleName: %v", body["ruleName"])
		}
		if body["behavior"] != "MONITOR" {
			t.Fatalf("unexpected behavior: %v", body["behavior"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"success":    true,
			"id":         "99999",
			"statusCode": 200,
		})
	})
	c := testClient(t, handler)
	resp, err := c.AppControlRulesCreate(context.Background(), AppControlRuleInput{
		RuleName: "New rule",
		Behavior: AppControlBehaviorMonitor,
		OSType:   []AppControlOSType{AppControlOSWindows},
		Scope: &AppControlScope{
			ScopeType: AppControlScopeSite,
			ScopeIDs:  []string{"225494730938493804"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	if resp.ID != "99999" {
		t.Fatalf("unexpected ID: %s", resp.ID)
	}
	if resp.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAppControlRulesUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/nac/config/api/v1/nac/rules/12345" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["ruleName"] != "Updated rule" {
			t.Fatalf("unexpected ruleName: %v", body["ruleName"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"success":    true,
			"id":         "12345",
			"statusCode": 200,
		})
	})
	c := testClient(t, handler)
	resp, err := c.AppControlRulesUpdate(context.Background(), "12345", AppControlRuleInput{
		RuleName: "Updated rule",
		Behavior: AppControlBehaviorAllow,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
	if resp.ID != "12345" {
		t.Fatalf("unexpected ID: %s", resp.ID)
	}
}

func TestAppControlRulesDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/nac/config/api/v1/nac/rules" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		ids := r.URL.Query()["ids"]
		if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
			t.Fatalf("unexpected ids: %v", ids)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["scopeType"] != "SITE" {
			t.Fatalf("expected SITE scope, got %v", body["scopeType"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"success":    true,
			"statusCode": 200,
		})
	})
	c := testClient(t, handler)
	resp, err := c.AppControlRulesDelete(context.Background(), []string{"a", "b"}, &AppControlScope{
		ScopeType: AppControlScopeSite,
		ScopeIDs:  []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
}

func TestAppControlLabelsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/nac/config/api/v1/nac/labels" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": "1", "labelName": "productivity"},
			{"id": "2", "labelName": "security"},
		})
	})
	c := testClient(t, handler)
	labels, err := c.AppControlLabelsList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(labels))
	}
	if labels[0].LabelName != "productivity" {
		t.Fatalf("unexpected label name: %s", labels[0].LabelName)
	}
	if labels[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAppControlSettingsGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/nac/config/api/v1/nac/settings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"fallbackBehavior":          "ALLOW",
			"enableApplicationControl":  true,
			"inheritApplicationControl": false,
		})
	})
	c := testClient(t, handler)
	settings, err := c.AppControlSettingsGet(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if settings.FallbackBehavior != AppControlBehaviorAllow {
		t.Fatalf("unexpected fallback: %s", settings.FallbackBehavior)
	}
	if !settings.EnableApplicationControl {
		t.Fatal("expected enabled=true")
	}
	if settings.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAppControlSettingsUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/nac/config/api/v1/nac/settings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"success":    true,
			"statusCode": 200,
		})
	})
	c := testClient(t, handler)
	enabled := true
	resp, err := c.AppControlSettingsUpdate(context.Background(), AppControlSettingsInput{
		FallbackBehavior:         AppControlBehaviorBlock,
		EnableApplicationControl: &enabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatal("expected success=true")
	}
}

func TestAppMgmtSettingsGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/application-management/settings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["siteIds"]; len(got) != 1 || got[0] != "123" {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"extensiveScanEnabled":   true,
				"isDefaultPolicy":        false,
				"hasBreakingInheritance": true,
				"scanSchedule": map[string]any{
					"scanEvery": 2,
					"repeatOn":  "Tuesday",
					"timezone":  "Europe/Berlin",
					"time":      "20:15",
				},
			},
		})
	})
	c := testClient(t, handler)
	s, err := c.AppMgmtSettingsGet(context.Background(), &AppMgmtSettingsListParams{
		SiteIDs: []string{"123"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.ExtensiveScanEnabled {
		t.Fatal("expected extensiveScanEnabled=true")
	}
	if s.ScanSchedule == nil {
		t.Fatal("expected scanSchedule to be present")
	}
	if s.ScanSchedule.ScanEvery != 2 {
		t.Fatalf("expected scanEvery=2, got %d", s.ScanSchedule.ScanEvery)
	}
	if s.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAppMgmtSettingsUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/application-management/settings" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter map[string]any `json:"filter"`
			Data   map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		sites, ok := body.Filter["siteIds"].([]any)
		if !ok || len(sites) != 1 {
			t.Fatalf("expected siteIds filter, got %v", body.Filter)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"extensiveScanEnabled":   false,
				"isDefaultPolicy":        true,
				"hasBreakingInheritance": false,
			},
		})
	})
	c := testClient(t, handler)
	enabled := false
	s, err := c.AppMgmtSettingsUpdate(context.Background(),
		AppMgmtSettingsScope{SiteIDs: []string{"123"}},
		AppMgmtSettingsUpdateData{ExtensiveScanEnabled: &enabled},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ExtensiveScanEnabled {
		t.Fatal("expected extensiveScanEnabled=false")
	}
}

func TestAppControlRulesListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{{"code": 403, "title": "Forbidden"}},
		})
	})
	c := testClient(t, handler)
	_, _, _, err := c.AppControlRulesList(context.Background(), nil)
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

func TestAppControlEnumValues(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"BehaviorAllow", string(AppControlBehaviorAllow), "ALLOW"},
		{"BehaviorMonitor", string(AppControlBehaviorMonitor), "MONITOR"},
		{"BehaviorBlock", string(AppControlBehaviorBlock), "BLOCK"},
		{"OSMacOS", string(AppControlOSMacOS), "MACOS"},
		{"OSWindows", string(AppControlOSWindows), "WINDOWS"},
		{"ScopeAccount", string(AppControlScopeAccount), "ACCOUNT"},
		{"ScopeSite", string(AppControlScopeSite), "SITE"},
		{"ScopeGroup", string(AppControlScopeGroup), "GROUP"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}
