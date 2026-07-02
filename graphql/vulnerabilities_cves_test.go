package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestCvesList(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"cves": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"id": "CVE-2024-0001", "description": "test cve",
							"nvdBaseScore": 9.8, "riskScore": 8.1, "epssScore": 0.5,
							"exploitedInTheWild": true, "publishedDate": "2024-01-01",
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	conn, err := c.CvesList(context.Background(), nil, nil, 50, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "cves(") {
		t.Errorf("query does not target cves: %s", gotReq.Query)
	}
	if gotReq.Variables["first"] != float64(50) {
		t.Errorf("expected first=50, got %v", gotReq.Variables["first"])
	}
	if conn.TotalCount != 1 || len(conn.Edges) != 1 {
		t.Fatalf("unexpected connection: total=%d edges=%d", conn.TotalCount, len(conn.Edges))
	}
	if conn.Edges[0].Node.ID != "CVE-2024-0001" || conn.Edges[0].Node.NVDBaseScore != 9.8 {
		t.Errorf("unexpected cve: %+v", conn.Edges[0].Node)
	}
}

func TestCveGet(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"cve": map[string]any{
					"id": "CVE-2024-0001", "description": "test cve", "nvdBaseScore": 9.8,
					"exploitMaturity": "FUNCTIONAL", "exploitedInTheWild": true,
				},
			},
		})
	})
	c := testClient(t, handler)
	cve, err := c.CveGet(context.Background(), "CVE-2024-0001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "cve(") {
		t.Errorf("query does not target cve: %s", gotReq.Query)
	}
	if gotReq.Variables["id"] != "CVE-2024-0001" {
		t.Errorf("expected id=CVE-2024-0001, got %v", gotReq.Variables["id"])
	}
	if cve.ID != "CVE-2024-0001" || cve.Description != "test cve" || !cve.ExploitedInTheWild {
		t.Errorf("unexpected cve: %+v", cve)
	}
	if cve.Raw == nil {
		t.Error("expected Raw to be populated")
	}
}

func TestUniqueCveCount(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{"uniqueCveCount": map[string]any{"count": 42}},
		})
	})
	c := testClient(t, handler)
	count, err := c.UniqueCveCount(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "uniqueCveCount(") {
		t.Errorf("query does not target uniqueCveCount: %s", gotReq.Query)
	}
	if count != 42 {
		t.Errorf("expected count=42, got %d", count)
	}
}

func TestTopVulnerableApplications(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"topVulnerableApplications": map[string]any{
					"applicationStats": []map[string]any{
						{"name": "openssl", "version": "1.0.0", "assetCount": 10, "cveCount": 5, "highestRiskScore": 9.1},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	stats, err := c.TopVulnerableApplications(context.Background(), nil, nil, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "topVulnerableApplications(") {
		t.Errorf("query does not target topVulnerableApplications: %s", gotReq.Query)
	}
	if gotReq.Variables["limit"] != float64(10) {
		t.Errorf("expected limit=10, got %v", gotReq.Variables["limit"])
	}
	if len(stats) != 1 || stats[0].Name != "openssl" || stats[0].CveCount != 5 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestTopVulnerableAssets(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"topVulnerableAssets": map[string]any{
					"assetStats": []map[string]any{
						{"name": "host-1", "scopeName": "site-a", "cveCount": 12, "highestRiskScore": 9.9},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	stats, err := c.TopVulnerableAssets(context.Background(), nil, nil, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "topVulnerableAssets(") {
		t.Errorf("query does not target topVulnerableAssets: %s", gotReq.Query)
	}
	if len(stats) != 1 || stats[0].Name != "host-1" || stats[0].ScopeName != "site-a" {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestTopVulnerableOsTypes(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"topVulnerableOsTypes": map[string]any{
					"osTypeStats": []map[string]any{
						{"name": "linux", "version": "ubuntu", "assetCount": 20, "cveCount": 30, "averageRiskScore": 6.5},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	stats, err := c.TopVulnerableOsTypes(context.Background(), nil, nil, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "topVulnerableOsTypes(") {
		t.Errorf("query does not target topVulnerableOsTypes: %s", gotReq.Query)
	}
	if len(stats) != 1 || stats[0].Name != "linux" || stats[0].AverageRiskScore != 6.5 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}
