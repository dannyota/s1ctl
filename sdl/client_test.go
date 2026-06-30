package sdl

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

func TestPowerQuery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/powerQuery" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req PowerQueryRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Query != "* | limit 10" {
			t.Fatalf("unexpected query: %s", req.Query)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"columns": []map[string]any{
				{"name": "timestamp", "type": "long"},
				{"name": "message", "type": "string"},
			},
			"values": [][]any{
				{1234567890, "hello"},
			},
		})
	})
	c := testClient(t, handler)
	resp, err := c.PowerQuery(context.Background(), &PowerQueryRequest{
		Query:     "* | limit 10",
		StartTime: "24h",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("unexpected status: %s", resp.Status)
	}
	if len(resp.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(resp.Columns))
	}
	if len(resp.Values) != 1 {
		t.Fatalf("expected 1 row, got %d", len(resp.Values))
	}
	if resp.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestPowerQueryHTTPError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
	})
	c := testClient(t, handler)
	_, err := c.PowerQuery(context.Background(), &PowerQueryRequest{Query: "*"})
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
