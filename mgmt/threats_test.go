package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestThreatsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "T1",
					"threatInfo": map[string]any{
						"threatName":       "Eicar",
						"classification":   "Malware",
						"confidenceLevel":  "malicious",
						"mitigationStatus": "mitigated",
						"analystVerdict":   "true_positive",
						"incidentStatus":   "resolved",
						"createdAt":        "2024-01-01T00:00:00Z",
					},
					"agentRealtimeInfo": map[string]any{
						"agentId": "A1",
					},
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	threats, pag, err := c.ThreatsList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(threats) != 1 {
		t.Fatalf("expected 1 threat, got %d", len(threats))
	}
	th := threats[0]
	if th.ID != "T1" {
		t.Fatalf("unexpected id: %s", th.ID)
	}
	if th.ThreatName != "Eicar" {
		t.Fatalf("unexpected threatName: %s", th.ThreatName)
	}
	if th.Classification != "Malware" {
		t.Fatalf("unexpected classification: %s", th.Classification)
	}
	if th.AgentID != "A1" {
		t.Fatalf("unexpected agentId: %s", th.AgentID)
	}
	if th.MitigationStatus != "mitigated" {
		t.Fatalf("unexpected mitigationStatus: %s", th.MitigationStatus)
	}
	if th.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestThreatsGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("ids") != "T1" {
			t.Fatalf("expected ids=T1, got %s", r.URL.Query().Get("ids"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "T1",
					"threatInfo": map[string]any{
						"threatName": "Test",
					},
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	th, err := c.ThreatsGet(context.Background(), "T1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.ID != "T1" {
		t.Fatalf("unexpected id: %s", th.ID)
	}
}

func TestThreatsGetNotFound(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	_, err := c.ThreatsGet(context.Background(), "MISSING")
	if err == nil {
		t.Fatal("expected error")
	}
}
