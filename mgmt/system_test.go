package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestSystemInfo(t *testing.T) {
	var gotMethod, gotPath, gotQuery string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"release":            "24.2",
				"version":            "24.2.1.100",
				"build":              "100",
				"patch":              "1",
				"latestAgentVersion": "24.1.2.3",
			},
		})
	})
	c := testClient(t, handler)
	info, err := c.SystemInfo(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/system/info" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotQuery != "" {
		t.Fatalf("expected no query params, got %q", gotQuery)
	}
	if info.Release != "24.2" {
		t.Fatalf("unexpected release: %s", info.Release)
	}
	if info.Version != "24.2.1.100" {
		t.Fatalf("unexpected version: %s", info.Version)
	}
	if info.Build != "100" {
		t.Fatalf("unexpected build: %s", info.Build)
	}
	if info.Patch != "1" {
		t.Fatalf("unexpected patch: %s", info.Patch)
	}
	if info.LatestAgentVersion != "24.1.2.3" {
		t.Fatalf("unexpected latestAgentVersion: %s", info.LatestAgentVersion)
	}
	if info.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSystemStatus(t *testing.T) {
	var gotMethod, gotPath string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"health": "ok"},
		})
	})
	c := testClient(t, handler)
	status, err := c.SystemStatus(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/system/status" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if status.Health != "ok" {
		t.Fatalf("unexpected health: %s", status.Health)
	}
	if status.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSystemErrors(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 503, "title": "Service Unavailable"},
			},
		})
	})
	c := testClient(t, handler)
	tests := []struct {
		name string
		call func() error
	}{
		{"info", func() error { _, err := c.SystemInfo(context.Background()); return err }},
		{"status", func() error { _, err := c.SystemStatus(context.Background()); return err }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected error")
			}
			var ae *APIError
			if !errors.As(err, &ae) {
				t.Fatalf("expected *APIError, got %T", err)
			}
			if ae.Status != http.StatusServiceUnavailable {
				t.Fatalf("expected 503, got %d", ae.Status)
			}
		})
	}
}
