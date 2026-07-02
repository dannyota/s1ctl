package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestAlertsFiltersCount(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointAlerts) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertFiltersCount": map[string]any{
					"data": []map[string]any{
						{
							"fieldId":     "severity",
							"hasNextPage": false,
							"values": []map[string]any{
								{"value": "HIGH", "label": "High", "count": 7},
								{"value": "LOW", "label": "Low", "count": 2},
							},
						},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	fields, err := c.AlertsFiltersCount(context.Background(), []string{"severity"},
		[]Filter{{FieldID: "status", StringIn: &InStr{Values: []string{"NEW"}}}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "alertFiltersCount(fieldIds: $fieldIds") {
		t.Errorf("query does not target alertFiltersCount: %s", gotReq.Query)
	}
	fieldIDs, ok := gotReq.Variables["fieldIds"].([]any)
	if !ok || len(fieldIDs) != 1 || fieldIDs[0] != "severity" {
		t.Errorf("unexpected fieldIds: %v", gotReq.Variables["fieldIds"])
	}
	if _, ok := gotReq.Variables["filters"]; !ok {
		t.Errorf("expected filters variable to be set")
	}
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	f := fields[0]
	if f.FieldID != "severity" || len(f.Values) != 2 {
		t.Errorf("unexpected field: %+v", f)
	}
	if f.Values[0].Value != "HIGH" || f.Values[0].Count != 7 {
		t.Errorf("unexpected value: %+v", f.Values[0])
	}
	if f.Raw == nil {
		t.Error("expected Raw to be populated")
	}
}

func TestAlertsGroupByCount(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertGroupByCount": map[string]any{
					"data": []map[string]any{
						{
							"fieldId":     "status",
							"hasNextPage": false,
							"values": []map[string]any{
								{"value": "NEW", "label": "New", "count": 5},
							},
						},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	scope := &Scope{ScopeIDs: []string{"site-1"}, ScopeType: "SITE"}
	fields, err := c.AlertsGroupByCount(context.Background(), []string{"status"}, nil, scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "alertGroupByCount(fieldIds: $fieldIds") {
		t.Errorf("query does not target alertGroupByCount: %s", gotReq.Query)
	}
	sc, ok := gotReq.Variables["scope"].(map[string]any)
	if !ok || sc["scopeType"] != "SITE" {
		t.Errorf("unexpected scope: %v", gotReq.Variables["scope"])
	}
	if len(fields) != 1 || fields[0].FieldID != "status" {
		t.Fatalf("unexpected fields: %+v", fields)
	}
	if fields[0].Values[0].Count != 5 {
		t.Errorf("unexpected count: %d", fields[0].Values[0].Count)
	}
}

func TestAlertsExport(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertsCsvExport": map[string]any{
					"data": "id,name\nalert-1,Test\n",
				},
			},
		})
	})
	c := testClient(t, handler)
	csv, err := c.AlertsExport(context.Background(),
		[]Filter{{FieldID: "severity", StringIn: &InStr{Values: []string{"HIGH"}}}}, nil, ViewTypeAll)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(gotReq.Query, "alertsCsvExport(") {
		t.Errorf("query does not target alertsCsvExport: %s", gotReq.Query)
	}
	if gotReq.Variables["viewType"] != "ALL" {
		t.Errorf("expected viewType=ALL, got %v", gotReq.Variables["viewType"])
	}
	if csv != "id,name\nalert-1,Test\n" {
		t.Errorf("unexpected csv: %q", csv)
	}
}

func TestAlertsExportOmitsEmptyViewType(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"alertsCsvExport": map[string]any{"data": ""}},
		})
	})
	c := testClient(t, handler)
	if _, err := c.AlertsExport(context.Background(), nil, nil, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotReq.Variables["viewType"]; ok {
		t.Errorf("expected viewType to be omitted, got %v", gotReq.Variables["viewType"])
	}
}
