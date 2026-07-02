package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestAlertHistory(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointAlerts) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertHistory": map[string]any{
					"edges": []map[string]any{
						{
							"cursor": "h1",
							"node": map[string]any{
								"createdAt": "2024-01-01T00:00:00Z",
								"eventText": "Status changed",
								"eventType": "STATUS_CHANGED",
								"historyItemCreator": map[string]any{
									"userId":   "analyst@example.com",
									"userType": "USER",
								},
							},
						},
						{
							"cursor": "h2",
							"node": map[string]any{
								"createdAt": "2024-01-01T00:01:00Z",
								"eventText": "Mitigation applied",
								"eventType": "MITIGATION",
								"historyItemData": map[string]any{
									"message": map[string]any{"content": "process killed", "type": "TEXT"},
								},
							},
						},
						{
							"cursor": "h3",
							"node": map[string]any{
								"createdAt": "2024-01-01T00:02:00Z",
								"eventText": "Alert enriched",
								"eventType": "ENRICHMENT",
								"historyItemData": map[string]any{
									"description": map[string]any{"content": "enrichment summary", "type": "MARKDOWN"},
								},
							},
						},
					},
					"pageInfo":   map[string]any{"hasNextPage": false, "endCursor": "h3"},
					"totalCount": 3,
				},
			},
		})
	})
	c := testClient(t, handler)
	conn, err := c.AlertHistory(context.Background(), "alert-1", 25, "cur-0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(gotReq.Query, "alertHistory(alertId: $alertId") {
		t.Errorf("query does not target alertHistory: %s", gotReq.Query)
	}
	if gotReq.Variables["alertId"] != "alert-1" {
		t.Errorf("expected alertId=alert-1, got %v", gotReq.Variables["alertId"])
	}
	if gotReq.Variables["first"] != float64(25) {
		t.Errorf("expected first=25, got %v", gotReq.Variables["first"])
	}
	if gotReq.Variables["after"] != "cur-0" {
		t.Errorf("expected after=cur-0, got %v", gotReq.Variables["after"])
	}

	if conn.TotalCount != 3 {
		t.Fatalf("expected totalCount=3, got %d", conn.TotalCount)
	}
	if len(conn.Edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(conn.Edges))
	}

	// Union member UserHistoryItemCreator.
	status := conn.Edges[0].Node
	if status.Creator == nil {
		t.Fatal("expected creator on status item")
	}
	if status.Creator.UserID != "analyst@example.com" || status.Creator.UserType != "USER" {
		t.Errorf("unexpected creator: %+v", status.Creator)
	}
	if status.ActorName() != "analyst@example.com" {
		t.Errorf("unexpected actor name: %s", status.ActorName())
	}
	if status.ActionData != nil {
		t.Errorf("expected nil data on status item, got %+v", status.ActionData)
	}

	// Union member MitigationActionHistoryItemData.
	mitigation := conn.Edges[1].Node
	if mitigation.ActionData == nil || mitigation.ActionData.Message == nil {
		t.Fatal("expected mitigation message data")
	}
	if mitigation.ActionData.Message.Content != "process killed" {
		t.Errorf("unexpected mitigation message: %s", mitigation.ActionData.Message.Content)
	}
	if mitigation.ActorName() != "" {
		t.Errorf("expected empty actor name without creator, got %s", mitigation.ActorName())
	}

	// Union member EnrichmentHistoryItemData.
	enrichment := conn.Edges[2].Node
	if enrichment.ActionData == nil || enrichment.ActionData.Description == nil {
		t.Fatal("expected enrichment description data")
	}
	if enrichment.ActionData.Description.Content != "enrichment summary" {
		t.Errorf("unexpected enrichment description: %s", enrichment.ActionData.Description.Content)
	}
	if enrichment.Raw == nil {
		t.Error("expected Raw to be populated")
	}
}

func TestAlertHistoryOmitsOptionalVars(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertHistory": map[string]any{
					"edges":      []any{},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 0,
				},
			},
		})
	})
	c := testClient(t, handler)
	if _, err := c.AlertHistory(context.Background(), "alert-1", 0, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := gotReq.Variables["first"]; ok {
		t.Errorf("expected first to be omitted, got %v", gotReq.Variables["first"])
	}
	if _, ok := gotReq.Variables["after"]; ok {
		t.Errorf("expected after to be omitted, got %v", gotReq.Variables["after"])
	}
}

func TestAlertHistoryGraphQLError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":   nil,
			"errors": []map[string]any{{"message": "alert not found"}},
		})
	})
	c := testClient(t, handler)
	_, err := c.AlertHistory(context.Background(), "missing", 0, "")
	if err == nil {
		t.Fatal("expected error")
	}
	var qe *QueryError
	if !errors.As(err, &qe) {
		t.Fatalf("expected *QueryError, got %T", err)
	}
	if qe.Errors[0].Message != "alert not found" {
		t.Fatalf("unexpected message: %s", qe.Errors[0].Message)
	}
}

func TestAlertGroups(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != string(EndpointAlerts) {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertGroups": map[string]any{
					"edges": []map[string]any{
						{"cursor": "g1", "node": map[string]any{"value": "HIGH", "label": "High", "count": 7}},
						{"cursor": "g2", "node": map[string]any{"value": "LOW", "label": "Low", "count": 2}},
					},
					"pageInfo":   map[string]any{"hasNextPage": false, "endCursor": "g2"},
					"totalCount": 2,
				},
			},
		})
	})
	c := testClient(t, handler)
	conn, err := c.AlertGroups(context.Background(), "severity", &ListParams{
		First: 10,
		Sort:  &SortInput{By: "count", Order: "desc"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Regression: the query must declare $sort and pass it to alertGroups.
	if !strings.Contains(gotReq.Query, "$sort: SortInput") {
		t.Errorf("query missing $sort variable declaration: %s", gotReq.Query)
	}
	if !strings.Contains(gotReq.Query, "sort: $sort") {
		t.Errorf("query does not pass sort through to alertGroups: %s", gotReq.Query)
	}

	if gotReq.Variables["groupByFieldId"] != "severity" {
		t.Errorf("expected groupByFieldId=severity, got %v", gotReq.Variables["groupByFieldId"])
	}
	if gotReq.Variables["first"] != float64(10) {
		t.Errorf("expected first=10, got %v", gotReq.Variables["first"])
	}
	sort, ok := gotReq.Variables["sort"].(map[string]any)
	if !ok {
		t.Fatalf("expected sort variable, got %v", gotReq.Variables["sort"])
	}
	if sort["by"] != "count" || sort["order"] != "desc" {
		t.Errorf("unexpected sort variable: %v", sort)
	}

	if conn.TotalCount != 2 {
		t.Fatalf("expected totalCount=2, got %d", conn.TotalCount)
	}
	if len(conn.Edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(conn.Edges))
	}
	group := conn.Edges[0].Node
	if group.Value != "HIGH" || group.Label != "High" || group.Count != 7 {
		t.Errorf("unexpected group: %+v", group)
	}
	if group.Raw == nil {
		t.Error("expected Raw to be populated")
	}
}

func TestAlertGroupsNilParams(t *testing.T) {
	var gotReq gqlRequest
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotReq)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"alertGroups": map[string]any{
					"edges":      []any{},
					"pageInfo":   map[string]any{"hasNextPage": false},
					"totalCount": 0,
				},
			},
		})
	})
	c := testClient(t, handler)
	conn, err := c.AlertGroups(context.Background(), "status", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotReq.Variables["groupByFieldId"] != "status" {
		t.Errorf("expected groupByFieldId=status, got %v", gotReq.Variables["groupByFieldId"])
	}
	if conn.TotalCount != 0 {
		t.Fatalf("expected totalCount=0, got %d", conn.TotalCount)
	}
}
