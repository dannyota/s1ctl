package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestReportsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/reports" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("scope") != "site" {
			t.Fatalf("unexpected scope: %s", q.Get("scope"))
		}
		if q.Get("frequency") != "weekly" {
			t.Fatalf("unexpected frequency: %s", q.Get("frequency"))
		}
		if q.Get("scheduleType") != "scheduled" {
			t.Fatalf("unexpected scheduleType: %s", q.Get("scheduleType"))
		}
		if q.Get("limit") != "5" {
			t.Fatalf("expected limit=5, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "1000000000000000001", "name": "Weekly Report",
					"scope": "site", "frequency": "weekly",
					"scheduleType": "scheduled", "status": "done",
					"creatorId":   "1000000000000000099",
					"creatorName": "admin",
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	reports, pag, err := c.ReportsList(context.Background(), &ReportListParams{
		Scope:        ReportScopeSite,
		Frequency:    ReportFrequencyWeekly,
		ScheduleType: ReportScheduleScheduled,
		Limit:        5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(reports))
	}
	r := reports[0]
	if r.ID != "1000000000000000001" {
		t.Fatalf("unexpected id: %s", r.ID)
	}
	if r.Name != "Weekly Report" {
		t.Fatalf("unexpected name: %s", r.Name)
	}
	if r.Scope != ReportScopeSite {
		t.Fatalf("unexpected scope: %s", r.Scope)
	}
	if r.Frequency != ReportFrequencyWeekly {
		t.Fatalf("unexpected frequency: %s", r.Frequency)
	}
	if r.ScheduleType != ReportScheduleScheduled {
		t.Fatalf("unexpected scheduleType: %s", r.ScheduleType)
	}
	if r.Status != "done" {
		t.Fatalf("unexpected status: %s", r.Status)
	}
	if r.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestReportsListNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %q", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	reports, _, err := c.ReportsList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reports) != 0 {
		t.Fatalf("expected 0 reports, got %d", len(reports))
	}
}

func TestReportTasksList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/report-tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("name") != "Monthly" {
			t.Fatalf("unexpected name: %s", q.Get("name"))
		}
		if q.Get("scope") != "account" {
			t.Fatalf("unexpected scope: %s", q.Get("scope"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "1000000000000000010", "name": "Monthly Summary",
					"scope": "account", "frequency": "monthly",
					"scheduleType": "scheduled", "day": "1",
					"creatorId":    "1000000000000000099",
					"creatorName":  "admin",
					"insightTypes": []string{"threat"},
					"isTrend":      false,
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	tasks, pag, err := c.ReportTasksList(context.Background(), &ReportTaskListParams{
		Name:  "Monthly",
		Scope: ReportScopeAccount,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	task := tasks[0]
	if task.ID != "1000000000000000010" {
		t.Fatalf("unexpected id: %s", task.ID)
	}
	if task.Scope != ReportScopeAccount {
		t.Fatalf("unexpected scope: %s", task.Scope)
	}
	if task.Frequency != ReportFrequencyMonthly {
		t.Fatalf("unexpected frequency: %s", task.Frequency)
	}
	if task.Day != "1" {
		t.Fatalf("unexpected day: %s", task.Day)
	}
	if task.IsTrend {
		t.Fatal("expected isTrend to be false")
	}
	if task.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestReportTasksCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/report-tasks" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter struct {
				SiteIDs    []string `json:"siteIds"`
				AccountIDs []string `json:"accountIds"`
				Scope      string   `json:"scope"`
			} `json:"filter"`
			Data struct {
				Name         string   `json:"name"`
				ScheduleType string   `json:"scheduleType"`
				Frequency    string   `json:"frequency"`
				Recipients   []string `json:"recipients"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Filter.Scope != "site" {
			t.Fatalf("unexpected filter scope: %s", body.Filter.Scope)
		}
		if len(body.Filter.SiteIDs) != 1 || body.Filter.SiteIDs[0] != "225494730938493804" {
			t.Fatalf("unexpected siteIds: %v", body.Filter.SiteIDs)
		}
		if body.Data.Name != "Test Task" {
			t.Fatalf("unexpected name: %s", body.Data.Name)
		}
		if body.Data.ScheduleType != "scheduled" {
			t.Fatalf("unexpected scheduleType: %s", body.Data.ScheduleType)
		}
		if body.Data.Frequency != "weekly" {
			t.Fatalf("unexpected frequency: %s", body.Data.Frequency)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "1000000000000000020"},
		})
	})
	c := testClient(t, handler)
	err := c.ReportTasksCreate(context.Background(),
		[]string{"225494730938493804"}, nil,
		ReportScopeSite,
		ReportTaskCreate{
			Name:         "Test Task",
			ScheduleType: ReportScheduleScheduled,
			Frequency:    ReportFrequencyWeekly,
			InsightTypes: json.RawMessage(`["threat"]`),
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReportsInsightTypes(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/reports/insights/types" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if len(q["siteIds"]) != 1 || q["siteIds"][0] != "225494730938493804" {
			t.Fatalf("unexpected siteIds: %v", q["siteIds"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"insightTypes": []string{"threat", "application", "risk"},
			},
		})
	})
	c := testClient(t, handler)
	raw, err := c.ReportsInsightTypes(context.Background(), &InsightTypesParams{
		SiteIDs: []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if raw == nil {
		t.Fatal("expected non-nil insight types")
	}
	var types []string
	if err := json.Unmarshal(raw, &types); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(types) != 3 {
		t.Fatalf("expected 3 insight types, got %d", len(types))
	}
}

func TestReportsInsightTypesNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %q", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"insightTypes": []string{"threat"},
			},
		})
	})
	c := testClient(t, handler)
	raw, err := c.ReportsInsightTypes(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if raw == nil {
		t.Fatal("expected non-nil insight types")
	}
}

func TestReportDownload(t *testing.T) {
	pdfContent := []byte("%PDF-1.4 fake content")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/reports/1000000000000000001/pdf" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Write(pdfContent)
	})
	c := testClient(t, handler)
	data, err := c.ReportDownload(context.Background(), "1000000000000000001", "pdf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(pdfContent) {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func TestReportDownloadHTML(t *testing.T) {
	htmlContent := []byte("<html><body>Report</body></html>")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reports/1000000000000000002/html" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlContent)
	})
	c := testClient(t, handler)
	data, err := c.ReportDownload(context.Background(), "1000000000000000002", "html")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(htmlContent) {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func TestReportDownloadError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 404, "title": "Not Found"},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.ReportDownload(context.Background(), "0000000000000000000", "pdf")
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 404 {
		t.Fatalf("expected 404, got %d", ae.Status)
	}
}

func TestReportsListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 503, "title": "Service Unavailable"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.ReportsList(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 503 {
		t.Fatalf("expected 503, got %d", ae.Status)
	}
}
