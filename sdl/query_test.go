package sdl

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
)

func TestQuery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/query" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			QueryType         string `json:"queryType"`
			Filter            string `json:"filter"`
			StartTime         string `json:"startTime"`
			MaxCount          int    `json:"maxCount"`
			ContinuationToken string `json:"continuationToken"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.QueryType != "log" {
			t.Fatalf("expected queryType=log, got %s", body.QueryType)
		}
		if body.Filter != "serverHost='app-01'" {
			t.Fatalf("unexpected filter: %s", body.Filter)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"matches": []map[string]any{
				{"timestamp": "1700000000", "message": "event-1", "severity": 3, "session": "s1"},
			},
			"sessions": map[string]any{
				"s1": map[string]any{"serverHost": "app-01"},
			},
			"continuationToken": "",
		})
	})
	c := testClient(t, handler)
	resp, err := c.Query(context.Background(), &LogQueryRequest{
		Filter:    "serverHost='app-01'",
		StartTime: "24h",
		MaxCount:  100,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "ok" {
		t.Fatalf("unexpected status: %s", resp.Status)
	}
	if len(resp.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(resp.Matches))
	}
	if resp.Matches[0].Message != "event-1" {
		t.Fatalf("unexpected message: %s", resp.Matches[0].Message)
	}
	if resp.Sessions == nil {
		t.Fatal("expected sessions to be populated")
	}
	if resp.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestQueryError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
	})
	c := testClient(t, handler)
	_, err := c.Query(context.Background(), &LogQueryRequest{Filter: "*"})
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

func TestQueryAllMultiPage(t *testing.T) {
	var callCount atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := callCount.Add(1)
		var body struct {
			ContinuationToken string `json:"continuationToken"`
		}
		json.NewDecoder(r.Body).Decode(&body)

		switch n {
		case 1:
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"matches": []map[string]any{
					{"timestamp": "1", "message": "e1", "session": "s1"},
					{"timestamp": "2", "message": "e2", "session": "s1"},
				},
				"sessions":          map[string]any{"s1": map[string]any{"host": "a"}},
				"continuationToken": "token-page2",
			})
		case 2:
			if body.ContinuationToken != "token-page2" {
				t.Errorf("expected token-page2, got %s", body.ContinuationToken)
			}
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"matches": []map[string]any{
					{"timestamp": "3", "message": "e3", "session": "s2"},
				},
				"sessions":          map[string]any{"s2": map[string]any{"host": "b"}},
				"continuationToken": "",
			})
		default:
			t.Fatalf("unexpected call #%d", n)
		}
	})
	c := testClient(t, handler)
	resp, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*", StartTime: "1h"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(resp.Matches))
	}
	if resp.Matches[0].Message != "e1" {
		t.Fatalf("unexpected first match: %s", resp.Matches[0].Message)
	}
	if resp.Matches[2].Message != "e3" {
		t.Fatalf("unexpected last match: %s", resp.Matches[2].Message)
	}
	if int(callCount.Load()) != 2 {
		t.Fatalf("expected 2 API calls, got %d", callCount.Load())
	}
}

func TestQueryAllTerminatesOnEmptyMatches(t *testing.T) {
	var callCount atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := callCount.Add(1)
		switch n {
		case 1:
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"matches": []map[string]any{
					{"timestamp": "1", "message": "e1"},
				},
				"continuationToken": "token-next",
			})
		case 2:
			// Server returns a token but no matches -- QueryAll must stop.
			json.NewEncoder(w).Encode(map[string]any{
				"status":            "ok",
				"matches":           []any{},
				"continuationToken": "token-stale",
			})
		default:
			t.Fatalf("unexpected call #%d -- QueryAll did not terminate on empty matches", n)
		}
	})
	c := testClient(t, handler)
	resp, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(resp.Matches))
	}
	if int(callCount.Load()) != 2 {
		t.Fatalf("expected 2 API calls, got %d", callCount.Load())
	}
}

func TestQueryAllTerminatesOnNonAdvancingToken(t *testing.T) {
	var callCount atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := callCount.Add(1)
		switch n {
		case 1:
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"matches": []map[string]any{
					{"timestamp": "1", "message": "e1"},
				},
				"continuationToken": "stuck-token",
			})
		case 2:
			// Server returns the same token again -- QueryAll must stop.
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"matches": []map[string]any{
					{"timestamp": "2", "message": "e2"},
				},
				"continuationToken": "stuck-token",
			})
		default:
			t.Fatalf("unexpected call #%d -- QueryAll did not terminate on non-advancing token", n)
		}
	})
	c := testClient(t, handler)
	resp, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Both pages' matches should be merged before stopping.
	if len(resp.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(resp.Matches))
	}
	if int(callCount.Load()) != 2 {
		t.Fatalf("expected 2 API calls, got %d", callCount.Load())
	}
}

func TestQueryAllMergesSessions(t *testing.T) {
	var callCount atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := callCount.Add(1)
		switch n {
		case 1:
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"matches": []map[string]any{
					{"timestamp": "1", "message": "e1", "session": "s1"},
				},
				"sessions":          map[string]any{"s1": map[string]any{"host": "a"}},
				"continuationToken": "tok2",
			})
		case 2:
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"matches": []map[string]any{
					{"timestamp": "2", "message": "e2", "session": "s2"},
				},
				"sessions":          map[string]any{"s2": map[string]any{"host": "b"}},
				"continuationToken": "tok3",
			})
		case 3:
			json.NewEncoder(w).Encode(map[string]any{
				"status":            "ok",
				"matches":           []map[string]any{},
				"continuationToken": "",
			})
		default:
			t.Fatalf("unexpected call #%d", n)
		}
	})
	c := testClient(t, handler)
	resp, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(resp.Sessions))
	}
	if _, ok := resp.Sessions["s1"]; !ok {
		t.Fatal("missing session s1")
	}
	if _, ok := resp.Sessions["s2"]; !ok {
		t.Fatal("missing session s2")
	}
}

func TestQueryAllWithMaxEvents(t *testing.T) {
	var callCount atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := callCount.Add(1)
		switch n {
		case 1:
			json.NewEncoder(w).Encode(map[string]any{
				"status": "ok",
				"matches": []map[string]any{
					{"timestamp": "1", "message": "e1"},
					{"timestamp": "2", "message": "e2"},
					{"timestamp": "3", "message": "e3"},
				},
				"continuationToken": "more",
			})
		default:
			// Should not be reached because maxEvents=2 caps after the first page.
			t.Fatalf("unexpected call #%d -- maxEvents should have stopped pagination", n)
		}
	})
	c := testClient(t, handler)
	resp, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*"}, WithMaxEvents(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Matches) != 2 {
		t.Fatalf("expected 2 matches (capped), got %d", len(resp.Matches))
	}
}

func TestQueryAllWithPageCallback(t *testing.T) {
	var callCount atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := callCount.Add(1)
		switch n {
		case 1:
			json.NewEncoder(w).Encode(map[string]any{
				"status":            "ok",
				"matches":           []map[string]any{{"message": "e1"}},
				"continuationToken": "p2",
			})
		case 2:
			json.NewEncoder(w).Encode(map[string]any{
				"status":            "ok",
				"matches":           []map[string]any{{"message": "e2"}, {"message": "e3"}},
				"continuationToken": "",
			})
		}
	})
	c := testClient(t, handler)
	var callbacks []int
	resp, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*"},
		WithPageCallback(func(fetched int) { callbacks = append(callbacks, fetched) }))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(resp.Matches))
	}
	if len(callbacks) != 2 {
		t.Fatalf("expected 2 callbacks, got %d", len(callbacks))
	}
	if callbacks[0] != 1 {
		t.Fatalf("expected first callback=1, got %d", callbacks[0])
	}
	if callbacks[1] != 3 {
		t.Fatalf("expected second callback=3, got %d", callbacks[1])
	}
}

func TestQueryAllSinglePage(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"matches": []map[string]any{
				{"timestamp": "1", "message": "only"},
			},
			"continuationToken": "",
		})
	})
	c := testClient(t, handler)
	resp, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(resp.Matches))
	}
}

func TestQueryAllError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("bad gateway"))
	})
	c := testClient(t, handler)
	_, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*"})
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

func TestQueryAllErrorOnSecondPage(t *testing.T) {
	var callCount atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := callCount.Add(1)
		if n == 1 {
			json.NewEncoder(w).Encode(map[string]any{
				"status":            "ok",
				"matches":           []map[string]any{{"message": "e1"}},
				"continuationToken": "next",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	})
	c := testClient(t, handler)
	_, err := c.QueryAll(context.Background(), &LogQueryRequest{Filter: "*"})
	if err == nil {
		t.Fatal("expected error on second page")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 500 {
		t.Fatalf("expected 500, got %d", ae.Status)
	}
}
