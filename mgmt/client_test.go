package mgmt

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

func TestNewClient(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	if c.BaseURL() != "https://example.sentinelone.net/web/api/v2.1" {
		t.Fatalf("unexpected baseURL: %s", c.BaseURL())
	}
}

func TestNewClientTrailingSlash(t *testing.T) {
	c := NewClient("https://example.sentinelone.net/", "tok")
	if c.BaseURL() != "https://example.sentinelone.net/web/api/v2.1" {
		t.Fatalf("unexpected baseURL: %s", c.BaseURL())
	}
}

func TestAPIErrorParsing(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden", "detail": "Insufficient permissions"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.AgentsList(context.Background(), nil)
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
	if ae.Title != "Forbidden" {
		t.Fatalf("expected Forbidden, got %s", ae.Title)
	}
	if ae.Detail != "Insufficient permissions" {
		t.Fatalf("expected detail, got %s", ae.Detail)
	}
}

func TestAPIErrorNonJSON(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("bad gateway"))
	})
	c := testClient(t, handler)
	_, _, err := c.AgentsList(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 502 {
		t.Fatalf("expected 502, got %d", ae.Status)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c := NewClient("https://example.sentinelone.net", "tok", WithHTTPClient(custom))
	if c.http != custom {
		t.Fatal("custom HTTP client not applied")
	}
}
