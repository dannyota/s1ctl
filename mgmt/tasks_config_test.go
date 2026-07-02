package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTasksConfigGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tasks-configuration" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("taskType") != "agents_upgrade" {
			t.Fatalf("unexpected taskType: %s", r.URL.Query().Get("taskType"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"maxConcurrent":           10,
				"timezoneGmt":             "GMT+00:00",
				"taskType":                "agents_upgrade",
				"maintenanceWindowsByDay": map[string]any{"monday": []any{}},
			},
		})
	})
	c := testClient(t, handler)
	cfg, err := c.TasksConfigGet(context.Background(), &TasksConfigParams{TaskType: TaskTypeAgentsUpgrade})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxConcurrent != 10 {
		t.Fatalf("unexpected maxConcurrent: %d", cfg.MaxConcurrent)
	}
	if cfg.TaskType != TaskTypeAgentsUpgrade {
		t.Fatalf("unexpected taskType: %s", cfg.TaskType)
	}
	if len(cfg.MaintenanceWindowsByDay) == 0 {
		t.Fatal("expected maintenanceWindowsByDay to be populated")
	}
	if cfg.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestTasksConfigUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/tasks-configuration" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				MaxConcurrent int `json:"maxConcurrent"`
			} `json:"data"`
			Filter struct {
				TaskType string `json:"taskType"`
			} `json:"filter"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data.MaxConcurrent != 5 {
			t.Fatalf("unexpected maxConcurrent: %d", body.Data.MaxConcurrent)
		}
		if body.Filter.TaskType != "agents_upgrade" {
			t.Fatalf("unexpected taskType: %s", body.Filter.TaskType)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"maxConcurrent": 5}})
	})
	c := testClient(t, handler)
	cfg, err := c.TasksConfigUpdate(context.Background(), TasksConfigWrite{
		Data:   TasksConfigData{MaxConcurrent: 5, TimezoneGMT: "GMT+00:00"},
		Filter: TasksConfigFilter{TaskType: TaskTypeAgentsUpgrade, Tenant: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxConcurrent != 5 {
		t.Fatalf("unexpected maxConcurrent: %d", cfg.MaxConcurrent)
	}
}

func TestTasksConfigFlexibleGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks-configuration/flexible" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"policyPayload": map[string]any{"windows": []any{}}},
		})
	})
	c := testClient(t, handler)
	cfg, err := c.TasksConfigFlexibleGet(context.Background(), &TasksConfigParams{TaskType: TaskTypeAgentsUpgrade})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.PolicyPayload) == 0 {
		t.Fatal("expected policyPayload to be populated")
	}
}

func TestTasksConfigFlexibleUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/tasks-configuration/flexible" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body["filter"]; !ok {
			t.Fatal("expected filter in forwarded body")
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"maxConcurrent": 3}})
	})
	c := testClient(t, handler)
	cfg, err := c.TasksConfigFlexibleUpdate(context.Background(),
		json.RawMessage(`{"data":{"maxConcurrent":3},"filter":{"taskType":"agents_upgrade","tenant":true}}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxConcurrent != 3 {
		t.Fatalf("unexpected maxConcurrent: %d", cfg.MaxConcurrent)
	}
}

func TestMaintenanceWindowsExport(t *testing.T) {
	csv := "scope,day,start,end\nGlobal,monday,09:00,17:00\n"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks-configuration/maintenance-windows/export" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("taskType") != "agents_upgrade" {
			t.Fatalf("unexpected taskType: %s", r.URL.Query().Get("taskType"))
		}
		w.Write([]byte(csv)) //nolint:errcheck
	})
	c := testClient(t, handler)
	data, err := c.MaintenanceWindowsExport(context.Background(), &TasksConfigParams{
		TaskType: TaskTypeAgentsUpgrade,
		Tenant:   true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != csv {
		t.Fatalf("unexpected export body: %q", string(data))
	}
}

func TestTaskTypeEnumValues(t *testing.T) {
	if string(TaskTypeAgentsUpgrade) != "agents_upgrade" {
		t.Fatalf("unexpected TaskTypeAgentsUpgrade: %s", TaskTypeAgentsUpgrade)
	}
	if string(TaskTypeScriptExecution) != "script_execution" {
		t.Fatalf("unexpected TaskTypeScriptExecution: %s", TaskTypeScriptExecution)
	}
}
