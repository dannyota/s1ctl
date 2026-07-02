package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestVulnerabilitiesNotes(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointVulnerabilities) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"vulnerabilityNotes": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"id":              "note-1",
							"vulnerabilityId": "v-1",
							"text":            "patch pending",
							"createdAt":       "2024-01-01T00:00:00Z",
							"author":          map[string]any{"id": "u1", "fullName": "Ana Lyst", "email": "ana@example.com", "deleted": false},
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	notes, err := c.VulnerabilitiesNotes(context.Background(), "v-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "vulnerabilityNotes(") {
		t.Errorf("query does not target vulnerabilityNotes: %s", gotReq.Query)
	}
	if gotReq.Variables["vulnerabilityId"] != "v-1" {
		t.Errorf("expected vulnerabilityId=v-1, got %v", gotReq.Variables["vulnerabilityId"])
	}
	if len(notes) != 1 || notes[0].Text != "patch pending" || notes[0].AuthorName() != "Ana Lyst" {
		t.Fatalf("unexpected notes: %+v", notes)
	}
}

func TestVulnerabilitiesAddNote(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{"addVulnerabilityNoteV2": map[string]any{"updatedFindingIds": []string{"v-1"}}},
		})
	})
	c := testClient(t, handler)
	if err := c.VulnerabilitiesAddNote(context.Background(), []string{"v-1"}, "note text"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "addVulnerabilityNoteV2(") {
		t.Errorf("query does not target addVulnerabilityNoteV2: %s", gotReq.Query)
	}
	if gotReq.Variables["text"] != "note text" {
		t.Errorf("expected text=note text, got %v", gotReq.Variables["text"])
	}
	if gotReq.Variables["filter"] == nil {
		t.Error("expected filter to be set")
	}
}

func TestVulnerabilitiesUpdateNote(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)                                                            //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"updateVulnerabilityNote": true}}) //nolint:errcheck
	})
	c := testClient(t, handler)
	if err := c.VulnerabilitiesUpdateNote(context.Background(), "note-1", "revised"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "updateVulnerabilityNote(") {
		t.Errorf("query does not target updateVulnerabilityNote: %s", gotReq.Query)
	}
	if gotReq.Variables["noteId"] != "note-1" || gotReq.Variables["text"] != "revised" {
		t.Errorf("unexpected variables: %v", gotReq.Variables)
	}
}

func TestVulnerabilitiesDeleteNote(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)                                                            //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"deleteVulnerabilityNote": true}}) //nolint:errcheck
	})
	c := testClient(t, handler)
	if err := c.VulnerabilitiesDeleteNote(context.Background(), "note-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "deleteVulnerabilityNote(") {
		t.Errorf("query does not target deleteVulnerabilityNote: %s", gotReq.Query)
	}
	if gotReq.Variables["noteId"] != "note-1" {
		t.Errorf("expected noteId=note-1, got %v", gotReq.Variables["noteId"])
	}
}

func TestVulnerabilitiesAssign(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{"vulnerabilityUserAssignmentV2": map[string]any{"updatedFindingIds": []string{"v-1"}}},
		})
	})
	c := testClient(t, handler)
	if err := c.VulnerabilitiesAssign(context.Background(), []string{"v-1"}, "user-9"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "vulnerabilityUserAssignmentV2(") {
		t.Errorf("query does not target vulnerabilityUserAssignmentV2: %s", gotReq.Query)
	}
	if gotReq.Variables["userId"] != "user-9" {
		t.Errorf("expected userId=user-9, got %v", gotReq.Variables["userId"])
	}
}

func TestVulnerabilitiesHistory(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"vulnerabilityHistory": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"createdAt": "2024-01-01T00:00:00Z", "eventText": "verdict set", "eventType": "VERDICT_CHANGE",
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	items, err := c.VulnerabilitiesHistory(context.Background(), "v-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "vulnerabilityHistory(") {
		t.Errorf("query does not target vulnerabilityHistory: %s", gotReq.Query)
	}
	if gotReq.Variables["vulnerabilityId"] != "v-1" {
		t.Errorf("expected vulnerabilityId=v-1, got %v", gotReq.Variables["vulnerabilityId"])
	}
	if len(items) != 1 || items[0].EventType != "VERDICT_CHANGE" {
		t.Fatalf("unexpected history: %+v", items)
	}
}

func TestVulnerabilitiesRelatedAssets(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{
				"vulnerabilityRelatedAssets": map[string]any{
					"edges": []map[string]any{
						{"cursor": "c1", "node": map[string]any{
							"vulnerabilityId": "v-1",
							"asset":           map[string]any{"id": "a1", "name": "host-1", "type": "COMPUTE"},
							"software":        map[string]any{"name": "openssl", "version": "1.0.0", "fixVersion": "1.0.1"},
						}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 1,
				},
			},
		})
	})
	c := testClient(t, handler)
	assets, err := c.VulnerabilitiesRelatedAssets(context.Background(), "v-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "vulnerabilityRelatedAssets(") {
		t.Errorf("query does not target vulnerabilityRelatedAssets: %s", gotReq.Query)
	}
	if !strings.Contains(gotReq.Query, "cveId: $cveId") {
		t.Errorf("query does not pass the id via the cveId argument: %s", gotReq.Query)
	}
	if gotReq.Variables["cveId"] != "v-1" {
		t.Errorf("expected cveId=v-1, got %v", gotReq.Variables["cveId"])
	}
	if _, ok := gotReq.Variables["filters"]; ok {
		t.Errorf("expected no filters variable, got %v", gotReq.Variables["filters"])
	}
	if len(assets) != 1 || assets[0].Asset.Name != "host-1" || assets[0].Software.Name != "openssl" {
		t.Fatalf("unexpected related assets: %+v", assets)
	}
}

func TestVulnerabilitiesExport(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)   //nolint:errcheck
		json.NewEncoder(w).Encode(map[string]any{ //nolint:errcheck
			"data": map[string]any{"vulnerabilitiesExportToCsv": map[string]any{"data": "id,cve\nv-1,CVE-2024-0001\n"}},
		})
	})
	c := testClient(t, handler)
	csv, err := c.VulnerabilitiesExport(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "vulnerabilitiesExportToCsv(") {
		t.Errorf("query does not target vulnerabilitiesExportToCsv: %s", gotReq.Query)
	}
	if !strings.Contains(csv, "CVE-2024-0001") {
		t.Errorf("unexpected csv: %q", csv)
	}
}
