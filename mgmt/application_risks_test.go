package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestApplicationRisksList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/application-management/risks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804", "225494730938493805"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if got := q["severities"]; !slices.Equal(got, []string{"CRITICAL", "HIGH"}) {
			t.Fatalf("unexpected severities: %v", got)
		}
		if q.Get("applicationVendor__contains") != "Example Corp" {
			t.Fatalf("unexpected applicationVendor__contains: %s", q.Get("applicationVendor__contains"))
		}
		if q.Get("includeRemovals") != "false" {
			t.Fatalf("expected includeRemovals=false, got %q", q.Get("includeRemovals"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", q.Get("limit"))
		}
		if q.Get("cursor") != "next-page" {
			t.Fatalf("expected cursor=next-page, got %s", q.Get("cursor"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "1000000000000000000", "cveId": "CVE-2024-0001",
					"applicationName": "Example App", "severity": "CRITICAL",
					"riskScore": "9.8", "mitigationStatus": "none",
					"endpointName": "host-1", "daysDetected": 12,
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "abc"},
		})
	})
	c := testClient(t, handler)
	include := false
	risks, pag, err := c.ApplicationRisksList(context.Background(), &ApplicationRiskListParams{
		SiteIDs:           []string{"225494730938493804", "225494730938493805"},
		Severities:        []string{"CRITICAL", "HIGH"},
		ApplicationVendor: "Example Corp",
		IncludeRemovals:   &include,
		Limit:             10,
		Cursor:            "next-page",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(risks) != 1 {
		t.Fatalf("expected 1 risk, got %d", len(risks))
	}
	if risks[0].CveID != "CVE-2024-0001" {
		t.Fatalf("unexpected cveId: %s", risks[0].CveID)
	}
	if risks[0].Severity != "CRITICAL" {
		t.Fatalf("unexpected severity: %s", risks[0].Severity)
	}
	if risks[0].DaysDetected != 12 {
		t.Fatalf("expected daysDetected=12, got %d", risks[0].DaysDetected)
	}
	if risks[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 || pag.NextCursor != "abc" {
		t.Fatalf("unexpected pagination: %+v", pag)
	}
}

func TestApplicationRisksListNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	risks, _, err := c.ApplicationRisksList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(risks) != 0 {
		t.Fatalf("expected 0 risks, got %d", len(risks))
	}
}

func TestApplicationCVEsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/application-management/risks/cves" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("cveId__contains") != "CVE-2024" {
			t.Fatalf("unexpected cveId__contains: %s", q.Get("cveId__contains"))
		}
		if q.Get("applicationName") != "Example App" {
			t.Fatalf("unexpected applicationName: %s", q.Get("applicationName"))
		}
		if got := q["groupIds"]; !slices.Equal(got, []string{"225494730938493904"}) {
			t.Fatalf("unexpected groupIds: %v", got)
		}
		if q.Get("limit") != "5" {
			t.Fatalf("expected limit=5, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"cveId": "CVE-2024-0001", "severity": "HIGH",
					"nvdBaseScore": "8.1", "cvssVersion": "3.1",
					"nvdUrl":             "https://nvd.example.com/CVE-2024-0001",
					"exploitedInTheWild": "true",
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	cves, pag, err := c.ApplicationCVEsList(context.Background(), &ApplicationCVEListParams{
		CveID:           "CVE-2024",
		ApplicationName: "Example App",
		GroupIDs:        []string{"225494730938493904"},
		Limit:           5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cves) != 1 {
		t.Fatalf("expected 1 cve, got %d", len(cves))
	}
	if cves[0].CveID != "CVE-2024-0001" {
		t.Fatalf("unexpected cveId: %s", cves[0].CveID)
	}
	if cves[0].NvdBaseScore != "8.1" {
		t.Fatalf("unexpected nvdBaseScore: %s", cves[0].NvdBaseScore)
	}
	if cves[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestApplicationRisksListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 500, "title": "Server Error", "detail": "internal"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.ApplicationRisksList(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 500 {
		t.Fatalf("expected 500, got %d", ae.Status)
	}
}
