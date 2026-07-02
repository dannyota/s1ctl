package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestThreatsAddToExclusions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/threats/add-to-exclusions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req struct {
			Filter ActionFilter `json:"filter"`
			Data   struct {
				TargetScope string   `json:"targetScope"`
				Type        string   `json:"type"`
				Mode        string   `json:"mode"`
				Value       string   `json:"value"`
				Actions     []string `json:"actions"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if req.Filter.IDs[0] != "T1" {
			t.Fatalf("unexpected filter: %+v", req.Filter)
		}
		if req.Data.TargetScope != "site" || req.Data.Type != "path" || req.Data.Mode != "suppress" {
			t.Fatalf("unexpected data: %+v", req.Data)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 1}})
	})
	c := testClient(t, handler)
	affected, err := c.ThreatsAddToExclusions(context.Background(),
		ActionFilter{IDs: []string{"T1"}},
		ThreatExclusionOptions{
			TargetScope: ThreatExclusionScopeSite,
			Type:        ThreatExclusionTypePath,
			Mode:        ThreatExclusionModeSuppress,
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestThreatsAddToExclusionsRequiresScopeAndType(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	if _, err := c.ThreatsAddToExclusions(context.Background(), ActionFilter{IDs: []string{"T1"}}, ThreatExclusionOptions{Type: ThreatExclusionTypeHash}); err == nil {
		t.Fatal("expected error for missing target scope")
	}
	if _, err := c.ThreatsAddToExclusions(context.Background(), ActionFilter{IDs: []string{"T1"}}, ThreatExclusionOptions{TargetScope: ThreatExclusionScopeSite}); err == nil {
		t.Fatal("expected error for missing type")
	}
}

func TestThreatsMitigateAlerts(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/threats/mitigate-alerts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req struct {
			Filter json.RawMessage `json:"filter"`
			Data   struct {
				Alerts []ThreatAlert `json:"alerts"`
				Action string        `json:"action"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if req.Filter != nil {
			t.Fatalf("mitigate-alerts must not send a filter, got %s", req.Filter)
		}
		if len(req.Data.Alerts) != 1 || req.Data.Alerts[0].AgentID != "A1" || req.Data.Alerts[0].Storyline != "S1" {
			t.Fatalf("unexpected alerts: %+v", req.Data.Alerts)
		}
		if req.Data.Action != "quarantine" {
			t.Fatalf("unexpected action: %s", req.Data.Action)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 1}})
	})
	c := testClient(t, handler)
	affected, err := c.ThreatsMitigateAlerts(context.Background(),
		[]ThreatAlert{{AgentID: "A1", Storyline: "S1"}}, ThreatMitigationQuarantine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestThreatsMitigateAlertsRequiresAlerts(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	if _, err := c.ThreatsMitigateAlerts(context.Background(), nil, ThreatMitigationKill); err == nil {
		t.Fatal("expected error for empty alerts")
	}
}

func TestThreatsSetExternalTicketID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/threats/external-ticket-id" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req struct {
			Filter ActionFilter `json:"filter"`
			Data   struct {
				ExternalTicketID string `json:"externalTicketId"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if req.Filter.IDs[0] != "T1" {
			t.Fatalf("unexpected filter: %+v", req.Filter)
		}
		if req.Data.ExternalTicketID != "JIRA-42" {
			t.Fatalf("unexpected ticket id: %s", req.Data.ExternalTicketID)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 1}})
	})
	c := testClient(t, handler)
	affected, err := c.ThreatsSetExternalTicketID(context.Background(), ActionFilter{IDs: []string{"T1"}}, "JIRA-42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestThreatsQuarantinedFiles(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/threats/T1/quarantined-files" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"filePath": "/tmp/evil", "fileName": "evil.exe", "fileSize": 1024},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	files, err := c.ThreatsQuarantinedFiles(context.Background(), "T1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 || files[0].FileName != "evil.exe" || files[0].FileSize != 1024 {
		t.Fatalf("unexpected files: %+v", files)
	}
	if files[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestThreatsQuarantinedFilesRequiresID(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	if _, err := c.ThreatsQuarantinedFiles(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestThreatsWhiteningOptions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/threats/T1/whitening-options" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"whiteningOptions": []string{"hash", "path"},
				"threatType":       []string{"static"},
				"threatPolicy":     "P1",
			},
		})
	})
	c := testClient(t, handler)
	opts, err := c.ThreatsWhiteningOptions(context.Background(), "T1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.WhiteningOptions) != 2 || opts.WhiteningOptions[0] != "hash" {
		t.Fatalf("unexpected options: %+v", opts.WhiteningOptions)
	}
	if opts.ThreatPolicy != "P1" {
		t.Fatalf("unexpected policy: %s", opts.ThreatPolicy)
	}
}

func TestThreatsExport(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/threats/export" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("siteIds") != "SITE1" {
			t.Fatalf("expected siteIds=SITE1, got %q", r.URL.Query().Get("siteIds"))
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Write([]byte("id,name\nT1,Eicar\n"))
	})
	c := testClient(t, handler)
	data, err := c.ThreatsExport(context.Background(), &ThreatListParams{SiteIDs: []string{"SITE1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "id,name\nT1,Eicar\n" {
		t.Fatalf("unexpected export body: %q", string(data))
	}
}
