package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := NewClient(srv.URL, "testtoken")
	c.baseURL = srv.URL
	return c
}

func TestDo(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatal("expected application/json content-type")
		}
		var req gqlRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Query != "{ test }" {
			t.Fatalf("unexpected query: %s", req.Query)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"value": "ok"},
		})
	})
	c := testClient(t, handler)
	var dst struct {
		Value string `json:"value"`
	}
	err := c.Do(context.Background(), EndpointAlerts, "{ test }", nil, &dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Value != "ok" {
		t.Fatalf("unexpected value: %s", dst.Value)
	}
}

func TestDoGraphQLError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":   nil,
			"errors": []map[string]any{{"message": "field not found"}},
		})
	})
	c := testClient(t, handler)
	err := c.Do(context.Background(), EndpointAlerts, "{ bad }", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var qe *QueryError
	if !errors.As(err, &qe) {
		t.Fatalf("expected *QueryError, got %T", err)
	}
	if qe.Errors[0].Message != "field not found" {
		t.Fatalf("unexpected message: %s", qe.Errors[0].Message)
	}
}

func TestDoHTTPError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	})
	c := testClient(t, handler)
	err := c.Do(context.Background(), EndpointAlerts, "{ test }", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var he *HTTPError
	if !errors.As(err, &he) {
		t.Fatalf("expected *HTTPError, got %T", err)
	}
	if he.Status != 401 {
		t.Fatalf("expected 401, got %d", he.Status)
	}
}

func TestEndpointPaths(t *testing.T) {
	tests := []struct {
		endpoint Endpoint
		want     string
	}{
		{EndpointAlerts, "/web/api/v2.1/unifiedalerts/graphql"},
		{EndpointMisconfigurations, "/web/api/v2.1/xspm/findings/misconfigurations/graphql"},
		{EndpointVulnerabilities, "/web/api/v2.1/xspm/findings/vulnerabilities/graphql"},
		{EndpointCloudPolicies, "/web/api/v2.1/cloudsecurity/policies/graphql"},
		{EndpointCloudOnboarding, "/web/api/v2.1/cloudonboarding/graphql"},
		{EndpointCloudCompliance, "/web/api/v2.1/cloudsecurity/compliance/graphql"},
	}
	for _, tt := range tests {
		if string(tt.endpoint) != tt.want {
			t.Errorf("endpoint %v = %q, want %q", tt.endpoint, string(tt.endpoint), tt.want)
		}
	}
}
