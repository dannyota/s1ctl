package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

// TestCloudPoliciesActionEmptyIDs verifies the SDK rejects an empty ID list
// before any HTTP request is made. The API treats an empty ids list as "act
// on all rules in scope", so sending it would be destructive.
func TestCloudPoliciesActionEmptyIDs(t *testing.T) {
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
			return c.CloudPoliciesAction(ctx, CloudPolicyActionEnable, nil)
		}},
		{"action empty ids", func() (*CloudPoliciesActionResponse, error) {
			return c.CloudPoliciesAction(ctx, CloudPolicyActionDelete, []string{})
		}},
		{"enable", func() (*CloudPoliciesActionResponse, error) {
			return c.CloudPoliciesEnable(ctx, nil)
		}},
		{"disable", func() (*CloudPoliciesActionResponse, error) {
			return c.CloudPoliciesDisable(ctx, nil)
		}},
		{"delete", func() (*CloudPoliciesActionResponse, error) {
			return c.CloudPoliciesDelete(ctx, nil)
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

// TestCloudPoliciesActionVerbs verifies each wrapper sends the action string
// matching its typed constant and decodes the returned rule IDs.
func TestCloudPoliciesActionVerbs(t *testing.T) {
	var gotVars map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointCloudPolicies) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var req gqlRequest
		json.NewDecoder(r.Body).Decode(&req)
		gotVars = req.Variables
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"actionOnCNSRules": map[string]any{"ids": []string{"rule-1", "rule-2"}},
			},
		})
	})
	c := testClient(t, handler)
	ctx := context.Background()
	ids := []string{"rule-1", "rule-2"}

	tests := []struct {
		name   string
		action CloudPolicyAction
		call   func() (*CloudPoliciesActionResponse, error)
	}{
		{"enable", CloudPolicyActionEnable, func() (*CloudPoliciesActionResponse, error) {
			return c.CloudPoliciesEnable(ctx, ids)
		}},
		{"disable", CloudPolicyActionDisable, func() (*CloudPoliciesActionResponse, error) {
			return c.CloudPoliciesDisable(ctx, ids)
		}},
		{"delete", CloudPolicyActionDelete, func() (*CloudPoliciesActionResponse, error) {
			return c.CloudPoliciesDelete(ctx, ids)
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
			input, ok := gotVars["input"].(map[string]any)
			if !ok {
				t.Fatalf("expected input variable, got %v", gotVars["input"])
			}
			sentIDs, ok := input["ids"].([]any)
			if !ok || len(sentIDs) != 2 || sentIDs[0] != "rule-1" || sentIDs[1] != "rule-2" {
				t.Errorf("unexpected input ids: %v", input["ids"])
			}
			if resp == nil {
				t.Fatal("expected non-nil response")
			}
			if len(resp.IDs) != 2 || resp.IDs[0] != "rule-1" || resp.IDs[1] != "rule-2" {
				t.Errorf("unexpected response ids: %v", resp.IDs)
			}
			if resp.Raw == nil {
				t.Error("expected Raw to be populated")
			}
		})
	}
}

// TestCloudPoliciesActionNullResponse verifies a null actionOnCNSRules payload
// yields a nil response without error; callers are expected to nil-check.
func TestCloudPoliciesActionNullResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"actionOnCNSRules": nil},
		})
	})
	c := testClient(t, handler)
	resp, err := c.CloudPoliciesAction(context.Background(), CloudPolicyActionDisable, []string{"rule-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Fatalf("expected nil response for null payload, got %+v", resp)
	}
}

// TestCloudPoliciesActionGraphQLError verifies a GraphQL errors array surfaces
// as a typed *QueryError.
func TestCloudPoliciesActionGraphQLError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":   nil,
			"errors": []map[string]any{{"message": "permission denied"}},
		})
	})
	c := testClient(t, handler)
	_, err := c.CloudPoliciesEnable(context.Background(), []string{"rule-1"})
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
