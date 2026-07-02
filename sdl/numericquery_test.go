package sdl

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestNumericQuery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/numericQuery" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			QueryType string `json:"queryType"`
			Filter    string `json:"filter"`
			Function  string `json:"function"`
			StartTime string `json:"startTime"`
			Buckets   int    `json:"buckets"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.QueryType != "numeric" {
			t.Fatalf("expected queryType=numeric, got %s", body.QueryType)
		}
		if body.Filter != "serverHost contains \"frontend\"" {
			t.Fatalf("unexpected filter: %s", body.Filter)
		}
		if body.Function != "mean(responseSize)" {
			t.Fatalf("unexpected function: %s", body.Function)
		}
		if body.Buckets != 60 {
			t.Fatalf("expected 60 buckets, got %d", body.Buckets)
		}
		v1, v2, v3 := 42.5, 38.0, 0.0
		json.NewEncoder(w).Encode(map[string]any{
			"status":   "success",
			"values":   []*float64{&v1, &v2, nil, &v3},
			"cpuUsage": 12,
		})
	})
	c := testClient(t, handler)
	resp, err := c.NumericQuery(context.Background(), &NumericQueryRequest{
		Filter:    "serverHost contains \"frontend\"",
		Function:  "mean(responseSize)",
		StartTime: "1h",
		Buckets:   60,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("unexpected status: %s", resp.Status)
	}
	if len(resp.Values) != 4 {
		t.Fatalf("expected 4 values, got %d", len(resp.Values))
	}
	if resp.Values[2] != nil {
		t.Fatalf("expected nil for bucket 2, got %v", *resp.Values[2])
	}
	if *resp.Values[0] != 42.5 {
		t.Fatalf("expected 42.5, got %f", *resp.Values[0])
	}
	if resp.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestNumericQueryError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	})
	c := testClient(t, handler)
	_, err := c.NumericQuery(context.Background(), &NumericQueryRequest{
		StartTime: "1h",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
