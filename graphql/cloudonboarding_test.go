package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestCnappEntitiesList(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointCloudOnboarding) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"cnappOnboardedCloudEntitiesV2": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"id": "e1", "entityId": "123456789012", "name": "Production",
							"type": "INDIVIDUAL", "onboardingStatus": "OPERATIONAL",
							"activeCoverage":  []string{"CLOUD_NATIVE_SECURITY"},
							"missingCoverage": []string{},
							"scope":           "Account / Site",
							"createdAt":       1700000000000,
							"updatedAt":       1700000001000,
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false, "endCursor": "c1"},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	filters := []CnappFilter{{
		FieldID:  "cloudProvider",
		StringIn: &CnappInStr{Values: []string{"AWS"}},
	}}
	scope := &CnappScopeSelector{ScopeType: CnappScopeTypeAccount, ScopeIDs: []int64{12}}
	conn, err := c.CnappEntitiesList(context.Background(), filters, scope, &ListParams{First: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.TotalCount != 1 || len(conn.Edges) != 1 {
		t.Fatalf("unexpected connection: %+v", conn)
	}
	n := conn.Edges[0].Node
	if n.ID != "e1" || n.EntityID != "123456789012" || n.OnboardingStatus != CnappOperationalStatusOperational {
		t.Errorf("unexpected node: %+v", n)
	}
	if n.Type != CnappCloudEntityTypeIndividual {
		t.Errorf("unexpected type: %v", n.Type)
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

func TestCnappEntitiesListNoFilters(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"cnappOnboardedCloudEntitiesV2": map[string]any{
					"edges":      []any{},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 0,
				},
			},
		})
	})
	c := testClient(t, handler)
	conn, err := c.CnappEntitiesList(context.Background(), []CnappFilter{}, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn.TotalCount != 0 {
		t.Fatalf("expected 0 total, got %d", conn.TotalCount)
	}
	// filters should be present (empty array, required by schema)
	if _, ok := gotVars["filters"]; !ok {
		t.Error("expected filters variable even when empty")
	}
}

func TestCnappEntityGet(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointCloudOnboarding) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"cnappOnboardedCloudEntity": map[string]any{
					"entityId":       "123456789012",
					"entityName":     "My AWS Account",
					"displayName":    "prod-account",
					"onboardingType": "INDIVIDUAL",
					"cloudProvider":  "AWS",
					"activeProducts": []map[string]any{
						{"s1Product": "CLOUD_NATIVE_SECURITY", "features": map[string]any{"detect": true}},
					},
					"extraProperties": map[string]any{"foo": "bar"},
				},
			},
		})
	})
	c := testClient(t, handler)
	detail, err := c.CnappEntityGet(context.Background(), []string{"123456789012"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail.EntityID != "123456789012" || detail.EntityName != "My AWS Account" {
		t.Fatalf("unexpected detail: %+v", detail)
	}
	if detail.OnboardingType != CnappOnboardingTypeIndividual {
		t.Errorf("unexpected onboarding type: %v", detail.OnboardingType)
	}
	if detail.CloudProvider != CnappCloudProviderAWS {
		t.Errorf("unexpected cloud provider: %v", detail.CloudProvider)
	}
	if detail.Raw == nil {
		t.Error("expected Raw to be populated")
	}
	req, ok := gotVars["request"].(map[string]any)
	if !ok {
		t.Fatalf("expected request variable, got %v", gotVars["request"])
	}
	ids, ok := req["accountIds"].([]any)
	if !ok || len(ids) != 1 || ids[0] != "123456789012" {
		t.Errorf("unexpected request accountIds: %v", req["accountIds"])
	}
}

func TestCnappEntityGetNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"cnappOnboardedCloudEntity": nil},
		})
	})
	c := testClient(t, handler)
	_, err := c.CnappEntityGet(context.Background(), []string{"missing"}, nil)
	var qe *QueryError
	if !errors.As(err, &qe) {
		t.Fatalf("expected *QueryError, got %v", err)
	}
}

func TestCnappOnboard(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointCloudOnboarding) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"onboardCnappCloudEntity": map[string]any{
					"message":   "onboarding initiated",
					"isSuccess": true,
				},
			},
		})
	})
	c := testClient(t, handler)
	input := json.RawMessage(`{"onBoardingType":"INDIVIDUAL","cloudProvider":"AWS","products":[{"s1Product":"CLOUD_NATIVE_SECURITY"}]}`)
	scope := &CnappScopeSelector{ScopeType: CnappScopeTypeAccount, ScopeIDs: []int64{42}}
	resp, err := c.CnappOnboard(context.Background(), input, scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || !resp.IsSuccess {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Message != "onboarding initiated" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
	if resp.Raw == nil {
		t.Error("expected Raw to be populated")
	}
	req, ok := gotVars["request"].(map[string]any)
	if !ok || req["cloudProvider"] != "AWS" {
		t.Errorf("unexpected request variable: %v", gotVars["request"])
	}
	if _, ok := gotVars["scopeSelector"]; !ok {
		t.Error("expected scopeSelector variable")
	}
}

func TestCnappOnboardGraphQLError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":   nil,
			"errors": []map[string]any{{"message": "permission denied"}},
		})
	})
	c := testClient(t, handler)
	_, err := c.CnappOnboard(context.Background(), json.RawMessage(`{}`), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var qe *QueryError
	if !errors.As(err, &qe) {
		t.Fatalf("expected *QueryError, got %T", err)
	}
	if qe.Errors[0].Message != "permission denied" {
		t.Fatalf("unexpected message: %s", qe.Errors[0].Message)
	}
}

func TestCnappDelete(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointCloudOnboarding) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		gotVars = decodeReq(t, r).Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"deleteCnappCloudEntity": map[string]any{
					"message":   "deleted",
					"isSuccess": true,
				},
			},
		})
	})
	c := testClient(t, handler)
	resp, err := c.CnappDelete(context.Background(), []string{"acc-1", "acc-2"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || !resp.IsSuccess {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if resp.Raw == nil {
		t.Error("expected Raw to be populated")
	}
	req, ok := gotVars["request"].(map[string]any)
	if !ok {
		t.Fatalf("expected request variable, got %v", gotVars["request"])
	}
	ids, ok := req["accountIds"].([]any)
	if !ok || len(ids) != 2 || ids[0] != "acc-1" || ids[1] != "acc-2" {
		t.Errorf("unexpected accountIds: %v", req["accountIds"])
	}
}

func TestCnappDeleteEmptyIDs(t *testing.T) {
	requests := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"deleteCnappCloudEntity": map[string]any{"isSuccess": true}},
		})
	})
	c := testClient(t, handler)
	ctx := context.Background()

	resp, err := c.CnappDelete(ctx, nil, nil)
	if !errors.Is(err, ErrCnappDeleteNoAccountIDs) {
		t.Errorf("nil ids: expected ErrCnappDeleteNoAccountIDs, got %v", err)
	}
	if resp != nil {
		t.Errorf("nil ids: expected nil response, got %+v", resp)
	}

	resp, err = c.CnappDelete(ctx, []string{}, nil)
	if !errors.Is(err, ErrCnappDeleteNoAccountIDs) {
		t.Errorf("empty ids: expected ErrCnappDeleteNoAccountIDs, got %v", err)
	}
	if resp != nil {
		t.Errorf("empty ids: expected nil response, got %+v", resp)
	}

	if requests != 0 {
		t.Fatalf("expected no HTTP requests for empty IDs, got %d", requests)
	}
}

func TestCnappDeleteNullResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"deleteCnappCloudEntity": nil},
		})
	})
	c := testClient(t, handler)
	resp, err := c.CnappDelete(context.Background(), []string{"acc-1"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatalf("expected nil response for null payload, got %+v", resp)
	}
}
