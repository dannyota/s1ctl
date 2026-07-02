package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func decodeReq(t *testing.T, r *http.Request) gqlRequest {
	t.Helper()
	var req gqlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.Fatalf("decode request: %v", err)
	}
	return req
}

func TestCNSRulesList(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointCloudPolicies) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"cnsRules": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"id": "rule-1", "name": "Public bucket", "severity": "HIGH",
							"status": "enabled", "type": "CloudMisconfiguration",
							"providers": []string{"AWS"}, "queryType": "rego",
							"scope": map[string]any{"id": "s1", "level": "SITE", "path": "/root/site"},
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false, "endCursor": "c1"},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	filters := []Filter{{FieldID: "severity", StringIn: &InStr{Values: []string{"HIGH"}}}}
	scope := &Scope{ScopeIDs: []string{"s1"}, ScopeType: "SITE"}
	conn, err := c.CNSRulesList(context.Background(), filters, scope, &ListParams{First: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.TotalCount != 1 || len(conn.Edges) != 1 {
		t.Fatalf("unexpected connection: %+v", conn)
	}
	n := conn.Edges[0].Node
	if n.ID != "rule-1" || n.Severity != "HIGH" || n.Type != "CloudMisconfiguration" {
		t.Errorf("unexpected node: %+v", n)
	}
	if n.Scope.Path != "/root/site" {
		t.Errorf("unexpected scope: %+v", n.Scope)
	}
	if n.Raw == nil {
		t.Error("expected Raw to be populated")
	}
	if gotVars["first"] != float64(10) {
		t.Errorf("expected first=10, got %v", gotVars["first"])
	}
	if _, ok := gotVars["filters"]; !ok {
		t.Error("expected filters variable")
	}
	if _, ok := gotVars["scope"]; !ok {
		t.Error("expected scope variable")
	}
}

func TestCNSRuleGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := decodeReq(t, r).Variables
		if vars["id"] != "rule-9" {
			t.Errorf("expected id=rule-9, got %v", vars["id"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"cnsRule": map[string]any{
				"id": "rule-9", "name": "Custom rego", "queryType": "rego",
				"rawQuery": "package s1\nallow { true }", "status": "disabled",
			}},
		})
	})
	c := testClient(t, handler)
	rule, err := c.CNSRuleGet(context.Background(), "rule-9", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.ID != "rule-9" || rule.RawQuery == "" {
		t.Fatalf("unexpected rule: %+v", rule)
	}
}

func TestCNSRuleGetNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"cnsRule": nil},
		})
	})
	c := testClient(t, handler)
	_, err := c.CNSRuleGet(context.Background(), "missing", nil)
	var qe *QueryError
	if !errors.As(err, &qe) {
		t.Fatalf("expected *QueryError, got %v", err)
	}
}

func TestCNSRuleCreate(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"createCNSRule": map[string]any{"id": "new-rule"}},
		})
	})
	c := testClient(t, handler)
	input := json.RawMessage(`{"name":"My rule","queryType":"rego","severity":"HIGH","type":"CloudMisconfiguration"}`)
	scope := &Scope{ScopeIDs: []string{"s1"}, ScopeType: "SITE"}
	resp, err := c.CNSRuleCreate(context.Background(), input, scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.ID != "new-rule" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	in, ok := gotVars["input"].(map[string]any)
	if !ok || in["name"] != "My rule" || in["severity"] != "HIGH" {
		t.Errorf("unexpected input variable: %v", gotVars["input"])
	}
}

func TestCNSRuleUpdate(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"updateCNSRule": true},
		})
	})
	c := testClient(t, handler)
	input := json.RawMessage(`{"name":"Renamed","queryType":"rego","severity":"LOW","type":"CloudMisconfiguration"}`)
	ok, err := c.CNSRuleUpdate(context.Background(), "rule-3", input, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected update to return true")
	}
	if gotVars["id"] != "rule-3" {
		t.Errorf("expected id=rule-3, got %v", gotVars["id"])
	}
}

// TestCNSRulesActionEmptyIDs verifies the SDK rejects an empty ID list before
// any HTTP request. The API treats an empty ids list as "act on all rules in
// scope", so sending it would be destructive.
func TestCNSRulesActionEmptyIDs(t *testing.T) {
	requests := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"actionOnCNSRules": map[string]any{"ids": []string{"rule-1"}}},
		})
	})
	c := testClient(t, handler)
	ctx := context.Background()

	calls := []struct {
		name string
		call func() (*CloudPoliciesActionResponse, error)
	}{
		{"action nil ids", func() (*CloudPoliciesActionResponse, error) {
			return c.CNSRulesAction(ctx, CNSRuleActionEnable, nil, nil)
		}},
		{"action empty ids", func() (*CloudPoliciesActionResponse, error) {
			return c.CNSRulesAction(ctx, CNSRuleActionDelete, []string{}, nil)
		}},
		{"enable", func() (*CloudPoliciesActionResponse, error) {
			return c.CNSRulesEnable(ctx, nil, nil)
		}},
		{"disable", func() (*CloudPoliciesActionResponse, error) {
			return c.CNSRulesDisable(ctx, nil, nil)
		}},
		{"delete", func() (*CloudPoliciesActionResponse, error) {
			return c.CNSRulesDelete(ctx, nil, nil)
		}},
	}
	for _, tc := range calls {
		resp, err := tc.call()
		if !errors.Is(err, ErrCloudPolicyActionNoIDs) {
			t.Errorf("%s: expected ErrCloudPolicyActionNoIDs, got %v", tc.name, err)
		}
		if resp != nil {
			t.Errorf("%s: expected nil response, got %+v", tc.name, resp)
		}
	}
	if requests != 0 {
		t.Fatalf("expected no HTTP requests for empty IDs, got %d", requests)
	}
}

// TestCNSRulesActionVerbs verifies each wrapper sends the action string matching
// its typed constant, forwards the scope, and decodes the returned rule IDs.
func TestCNSRulesActionVerbs(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"actionOnCNSRules": map[string]any{"ids": []string{"rule-1", "rule-2"}},
			},
		})
	})
	c := testClient(t, handler)
	ctx := context.Background()
	ids := []string{"rule-1", "rule-2"}
	scope := &Scope{ScopeIDs: []string{"s1"}, ScopeType: "SITE"}

	tests := []struct {
		name   string
		action CNSRuleAction
		call   func() (*CloudPoliciesActionResponse, error)
	}{
		{"enable", CNSRuleActionEnable, func() (*CloudPoliciesActionResponse, error) {
			return c.CNSRulesEnable(ctx, ids, scope)
		}},
		{"disable", CNSRuleActionDisable, func() (*CloudPoliciesActionResponse, error) {
			return c.CNSRulesDisable(ctx, ids, scope)
		}},
		{"delete", CNSRuleActionDelete, func() (*CloudPoliciesActionResponse, error) {
			return c.CNSRulesDelete(ctx, ids, scope)
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVars = nil
			resp, err := tt.call()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotVars["action"] != string(tt.action) {
				t.Errorf("expected action=%q, got %v", string(tt.action), gotVars["action"])
			}
			if _, ok := gotVars["scope"]; !ok {
				t.Error("expected scope variable")
			}
			input, ok := gotVars["input"].(map[string]any)
			if !ok {
				t.Fatalf("expected input variable, got %v", gotVars["input"])
			}
			sentIDs, ok := input["ids"].([]any)
			if !ok || len(sentIDs) != 2 {
				t.Errorf("unexpected input ids: %v", input["ids"])
			}
			if resp == nil || len(resp.IDs) != 2 {
				t.Fatalf("unexpected response: %+v", resp)
			}
		})
	}
}

func TestCNSRuleEvaluate(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"evaluateCNSRegoRule": map[string]any{
				"result": "pass", "error": "", "data": map[string]any{"matched": true},
			}},
		})
	})
	c := testClient(t, handler)
	resource := json.RawMessage(`{"bucket":"public"}`)
	resp, err := c.CNSRuleEvaluate(context.Background(), "pol-1", "package s1", resource, `{"k":"v"}`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Result != "pass" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if gotVars["regoQuery"] != "package s1" {
		t.Errorf("expected regoQuery, got %v", gotVars["regoQuery"])
	}
	if gotVars["policyId"] != "pol-1" {
		t.Errorf("expected policyId=pol-1, got %v", gotVars["policyId"])
	}
	if gotVars["ruleConfigParameters"] != `{"k":"v"}` {
		t.Errorf("expected ruleConfigParameters, got %v", gotVars["ruleConfigParameters"])
	}
	rd, ok := gotVars["resourceData"].(map[string]any)
	if !ok || rd["bucket"] != "public" {
		t.Errorf("unexpected resourceData: %v", gotVars["resourceData"])
	}
}

func TestCNSRuleEvaluateOmitsEmptyOptionals(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"evaluateCNSRegoRule": map[string]any{"result": "fail"}},
		})
	})
	c := testClient(t, handler)
	_, err := c.CNSRuleEvaluate(context.Background(), "", "package s1", json.RawMessage(`{}`), "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotVars["policyId"]; ok {
		t.Error("expected policyId to be omitted when empty")
	}
	if _, ok := gotVars["ruleConfigParameters"]; ok {
		t.Error("expected ruleConfigParameters to be omitted when empty")
	}
}

func TestCNSRuleTypes(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"cnsRuleTypes": []map[string]any{
				{"key": "CloudMisconfiguration", "title": "Cloud Misconfiguration"},
				{"key": "KubeMisconfiguration", "title": "Kubernetes Misconfiguration"},
			}},
		})
	})
	c := testClient(t, handler)
	types, err := c.CNSRuleTypes(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(types) != 2 || types[0].Key != "CloudMisconfiguration" {
		t.Fatalf("unexpected types: %+v", types)
	}
}

func TestCNSRuleConfig(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"cnsRuleConfig": map[string]any{"maxRego": 4096}},
		})
	})
	c := testClient(t, handler)
	cfg, err := c.CNSRuleConfig(context.Background(), nil, CNSRuleTypeCloudMisconfiguration)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg) == 0 {
		t.Fatal("expected non-empty config")
	}
	if gotVars["type"] != string(CNSRuleTypeCloudMisconfiguration) {
		t.Errorf("expected type var, got %v", gotVars["type"])
	}
}
