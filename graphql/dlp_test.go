package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestDLPRulesList(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointCloudPolicies) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"dataProtectionRules": map[string]any{
					"nodes": []map[string]any{
						{
							"id": "dlp-1", "name": "Block PII", "status": "ENABLED",
							"rank": 1, "systemPolicy": false,
							"scope":           map[string]any{"id": "s1", "level": "SITE", "path": "/root/site"},
							"classifications": []map[string]any{{"id": "c1", "name": "SSN", "type": "SENSITIVE_DATA"}},
						},
					},
					"pageInfo": map[string]any{
						"currentPage": 1, "hasNextPage": false, "hasPreviousPage": false,
						"pageSize": 20, "totalCount": 1, "totalPages": 1,
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	filter := &DLPRuleFilter{SearchName: "Block", Status: []DLPRuleStatus{DLPRuleStatusEnabled}}
	scope := &Scope{ScopeIDs: []string{"s1"}, ScopeType: "SITE"}
	conn, err := c.DLPRulesList(context.Background(), filter, scope, &DLPPage{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.PageInfo.TotalCount != 1 || len(conn.Nodes) != 1 {
		t.Fatalf("unexpected connection: %+v", conn)
	}
	n := conn.Nodes[0]
	if n.ID != "dlp-1" || n.Status != DLPRuleStatusEnabled || n.Rank != 1 {
		t.Errorf("unexpected node: %+v", n)
	}
	if n.Scope.Path != "/root/site" {
		t.Errorf("unexpected scope: %+v", n.Scope)
	}
	if len(n.Classifications) != 1 || n.Classifications[0].Type != DLPClassificationTypeSensitiveData {
		t.Errorf("unexpected classifications: %+v", n.Classifications)
	}
	if n.Raw == nil {
		t.Error("expected Raw to be populated")
	}
	if _, ok := gotVars["filter"]; !ok {
		t.Error("expected filter variable")
	}
	if _, ok := gotVars["scope"]; !ok {
		t.Error("expected scope variable")
	}
	pag, ok := gotVars["pagination"].(map[string]any)
	if !ok || pag["page"] != float64(1) || pag["pageSize"] != float64(20) {
		t.Errorf("unexpected pagination variable: %v", gotVars["pagination"])
	}
}

func TestDLPRuleGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if vars := decodeReq(t, r).Variables; vars["id"] != "dlp-9" {
			t.Errorf("expected id=dlp-9, got %v", vars["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"dataProtectionRule": map[string]any{
				"id": "dlp-9", "name": "Full rule", "status": "DISABLED", "rank": 3,
				"actions":              map[string]any{"actionTaken": "BLOCK"},
				"impactedEndpoints":    map[string]any{"scopeType": "ALL_ENDPOINTS_IN_SCOPE"},
				"inspectionConditions": map[string]any{"fileTypes": []string{"pdf"}},
			}},
		})
	})
	c := testClient(t, handler)
	rule, err := c.DLPRuleGet(context.Background(), "dlp-9", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.ID != "dlp-9" || rule.Status != DLPRuleStatusDisabled {
		t.Fatalf("unexpected rule: %+v", rule)
	}
	if len(rule.Actions) == 0 || len(rule.InspectionConditions) == 0 {
		t.Errorf("expected raw bodies populated: %+v", rule)
	}
}

func TestDLPRuleGetNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"dataProtectionRule": nil},
		})
	})
	c := testClient(t, handler)
	_, err := c.DLPRuleGet(context.Background(), "missing", nil)
	var qe *QueryError
	if !errors.As(err, &qe) {
		t.Fatalf("expected *QueryError, got %v", err)
	}
}

func TestDLPRuleSingleToggle(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"enableDataProtectionRule":  map[string]any{"id": "dlp-1", "name": "R", "status": "ENABLED"},
				"disableDataProtectionRule": map[string]any{"id": "dlp-1", "name": "R", "status": "DISABLED"},
			},
		})
	})
	c := testClient(t, handler)
	scope := &Scope{ScopeIDs: []string{"s1"}, ScopeType: "SITE"}

	rule, err := c.DLPRuleEnable(context.Background(), "dlp-1", scope)
	if err != nil {
		t.Fatalf("enable: %v", err)
	}
	if rule == nil || rule.Status != DLPRuleStatusEnabled {
		t.Fatalf("unexpected enable response: %+v", rule)
	}
	if gotVars["id"] != "dlp-1" {
		t.Errorf("expected id var, got %v", gotVars["id"])
	}
	if _, ok := gotVars["scope"]; !ok {
		t.Error("expected scope variable")
	}
	if _, err := c.DLPRuleDisable(context.Background(), "dlp-1", scope); err != nil {
		t.Fatalf("disable: %v", err)
	}
}

func TestDLPRuleSingleDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if vars := decodeReq(t, r).Variables; vars["id"] != "dlp-7" {
			t.Errorf("expected id=dlp-7, got %v", vars["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"deleteDataProtectionRule": true},
		})
	})
	c := testClient(t, handler)
	ok, err := c.DLPRuleDelete(context.Background(), "dlp-7", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected delete to return true")
	}
}

func TestDLPRulesBulkVerbs(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"bulkEnableDataProtectionRules": []map[string]any{
					{"id": "dlp-1", "status": "ENABLED"}, {"id": "dlp-2", "status": "ENABLED"},
				},
				"bulkDisableDataProtectionRules": []map[string]any{
					{"id": "dlp-1", "status": "DISABLED"}, {"id": "dlp-2", "status": "DISABLED"},
				},
				"bulkDeleteDataProtectionRules": true,
			},
		})
	})
	c := testClient(t, handler)
	ctx := context.Background()
	ids := []string{"dlp-1", "dlp-2"}
	scope := &Scope{ScopeIDs: []string{"s1"}, ScopeType: "SITE"}

	rules, err := c.DLPRulesBulkEnable(ctx, ids, scope)
	if err != nil {
		t.Fatalf("bulk enable: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	sentIDs, ok := gotVars["ids"].([]any)
	if !ok || len(sentIDs) != 2 {
		t.Errorf("unexpected ids var: %v", gotVars["ids"])
	}
	if _, ok := gotVars["scope"]; !ok {
		t.Error("expected scope variable")
	}
	if rules, err := c.DLPRulesBulkDisable(ctx, ids, scope); err != nil || len(rules) != 2 {
		t.Fatalf("bulk disable: rules=%d err=%v", len(rules), err)
	}
	deleted, err := c.DLPRulesBulkDelete(ctx, ids, scope)
	if err != nil || !deleted {
		t.Fatalf("bulk delete: deleted=%v err=%v", deleted, err)
	}
}

// TestDLPRulesBulkEmptyIDs verifies the SDK rejects an empty ID list before any
// HTTP request. bulk*DataProtectionRules take a mandatory ids list; sending an
// empty one is rejected to avoid an unbounded action.
func TestDLPRulesBulkEmptyIDs(t *testing.T) {
	requests := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	})
	c := testClient(t, handler)
	ctx := context.Background()

	rules, err := c.DLPRulesBulkEnable(ctx, nil, nil)
	if !errors.Is(err, ErrCloudPolicyActionNoIDs) || rules != nil {
		t.Errorf("bulk enable: expected guard error, got rules=%v err=%v", rules, err)
	}
	rules, err = c.DLPRulesBulkDisable(ctx, []string{}, nil)
	if !errors.Is(err, ErrCloudPolicyActionNoIDs) || rules != nil {
		t.Errorf("bulk disable: expected guard error, got rules=%v err=%v", rules, err)
	}
	deleted, err := c.DLPRulesBulkDelete(ctx, nil, nil)
	if !errors.Is(err, ErrCloudPolicyActionNoIDs) || deleted {
		t.Errorf("bulk delete: expected guard error, got deleted=%v err=%v", deleted, err)
	}
	if requests != 0 {
		t.Fatalf("expected no HTTP requests for empty IDs, got %d", requests)
	}
}

func TestDLPClassificationsList(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"dlpClassifications": map[string]any{
					"nodes": []map[string]any{
						{
							"id": "c1", "name": "Credit cards", "type": "REGEX",
							"usedInRulesCount": 2, "systemPolicy": true,
							"patterns": []map[string]any{{"id": "p1", "name": "visa", "pattern": "\\d+"}},
						},
					},
					"pageInfo": map[string]any{
						"currentPage": 1, "hasNextPage": true, "pageSize": 20,
						"totalCount": 5, "totalPages": 1,
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	filter := &DLPClassificationFilter{Type: []DLPClassificationType{DLPClassificationTypeRegex}}
	conn, err := c.DLPClassificationsList(context.Background(), filter, nil, &DLPPage{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conn.Nodes) != 1 || conn.PageInfo.TotalCount != 5 || !conn.PageInfo.HasNextPage {
		t.Fatalf("unexpected connection: %+v", conn)
	}
	n := conn.Nodes[0]
	if n.Type != DLPClassificationTypeRegex || n.UsedInRulesCount != 2 || !n.SystemPolicy {
		t.Errorf("unexpected node: %+v", n)
	}
	if len(n.Patterns) == 0 {
		t.Error("expected patterns raw body populated")
	}
	if _, ok := gotVars["filter"]; !ok {
		t.Error("expected filter variable")
	}
}

func TestDLPClassificationGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if vars := decodeReq(t, r).Variables; vars["id"] != "c9" {
			t.Errorf("expected id=c9, got %v", vars["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"dlpClassification": map[string]any{
				"id": "c9", "name": "Secrets", "type": "SECRETS",
				"secretDetectors": []map[string]any{{"id": "d1", "name": "aws_key", "selected": true}},
			}},
		})
	})
	c := testClient(t, handler)
	cls, err := c.DLPClassificationGet(context.Background(), "c9", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cls.ID != "c9" || cls.Type != DLPClassificationTypeSecrets || len(cls.SecretDetectors) == 0 {
		t.Fatalf("unexpected classification: %+v", cls)
	}
}

func TestDLPClassificationGetNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"dlpClassification": nil},
		})
	})
	c := testClient(t, handler)
	_, err := c.DLPClassificationGet(context.Background(), "missing", nil)
	var qe *QueryError
	if !errors.As(err, &qe) {
		t.Fatalf("expected *QueryError, got %v", err)
	}
}

func TestDLPClassificationDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if vars := decodeReq(t, r).Variables; vars["id"] != "c3" {
			t.Errorf("expected id=c3, got %v", vars["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"deleteDlpClassification": true},
		})
	})
	c := testClient(t, handler)
	ok, err := c.DLPClassificationDelete(context.Background(), "c3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected delete to return true")
	}
}

func TestDLPEngineSettings(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"dlpEngineSettings": map[string]any{
				"blockEncryptedArchive": true, "characterInspectionDepth": "BALANCED",
				"classificationsToInspect": 5, "enableOcr": true, "maskEvidence": false,
				"maxInspectedFileSize": 10485760, "preventAction": "BLOCK",
				"ignoreKeywords": []string{"test"},
				"scope":          map[string]any{"id": "s1", "level": "SITE", "path": "/root/site"},
			}},
		})
	})
	c := testClient(t, handler)
	scope := &Scope{ScopeIDs: []string{"s1"}, ScopeType: "SITE"}
	settings, err := c.DLPEngineSettings(context.Background(), scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !settings.BlockEncryptedArchive || settings.PreventAction != "BLOCK" || settings.MaxInspectedFileSize != 10485760 {
		t.Fatalf("unexpected settings: %+v", settings)
	}
	if _, ok := gotVars["scope"]; !ok {
		t.Error("expected scope variable")
	}
}

// TestDLPEngineSettingsRequiresScope verifies the scope-required guard: the
// schema marks scope as non-null, so the SDK rejects a nil scope up front.
func TestDLPEngineSettingsRequiresScope(t *testing.T) {
	requests := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{}})
	})
	c := testClient(t, handler)
	_, err := c.DLPEngineSettings(context.Background(), nil)
	if !errors.Is(err, ErrDLPScopeRequired) {
		t.Fatalf("expected ErrDLPScopeRequired, got %v", err)
	}
	if requests != 0 {
		t.Fatalf("expected no HTTP requests, got %d", requests)
	}
}
