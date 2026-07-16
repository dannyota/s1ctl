package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAutomationList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflows" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("limit"); got != "10" {
			t.Fatalf("expected limit=10, got %s", got)
		}
		if got := q.Get("states"); got != "active" {
			t.Fatalf("expected states=active, got %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "wf-1",
					"workflow": map[string]any{
						"id":              "wf-1",
						"version_id":      "v-1",
						"name":            "Alert triage",
						"description":     "Auto-triage alerts",
						"state":           "active",
						"lifecycle_state": "active",
						"status":          "idle",
						"scope_id":        "000000",
						"scope_level":     "site",
						"mgmt_id":         "mgmt-1",
						"created_at":      "2025-01-01T00:00:00Z",
						"updated_at":      "2025-01-02T00:00:00Z",
						"created_by":      "user-1",
						"updated_by":      "user-1",
						"version_count":   2,
						"timeout":         86400,
						"is_snippet":      false,
					},
					"actions": []map[string]any{
						{"id": "act-1", "type": "manual_trigger"},
					},
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": ""},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.AutomationList(context.Background(), &AutomationListParams{
		States: []string{"active"},
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	wf := items[0].Workflow
	if wf.Name != "Alert triage" {
		t.Fatalf("unexpected name: %s", wf.Name)
	}
	if wf.State != WorkflowStateActive {
		t.Fatalf("unexpected state: %s", wf.State)
	}
	if wf.ScopeLevel != AutomationScopeSite {
		t.Fatalf("unexpected scope: %s", wf.ScopeLevel)
	}
	if len(items[0].Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(items[0].Actions))
	}
	if items[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestAutomationListNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0, "nextCursor": ""},
		})
	})
	c := testClient(t, handler)
	items, _, err := c.AutomationList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestAutomationVersions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflows/versions/list/wf-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"versions": []map[string]any{
				{
					"id":               "wf-1",
					"version_id":       "v-1",
					"name":             "Alert triage",
					"state":            "active",
					"scope_id":         "000000",
					"scope_level":      "site",
					"mgmt_id":          "mgmt-1",
					"created_by":       "user-1",
					"updated_by":       "user-1",
					"execution_status": "Completed",
				},
				{
					"id":               "wf-1",
					"version_id":       "v-2",
					"name":             "Alert triage",
					"state":            "inactive",
					"scope_id":         "000000",
					"scope_level":      "site",
					"mgmt_id":          "mgmt-1",
					"created_by":       "user-1",
					"updated_by":       "user-1",
					"execution_status": "Error",
				},
			},
		})
	})
	c := testClient(t, handler)
	versions, err := c.AutomationVersions(context.Background(), "wf-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
	if versions[0].VersionID != "v-1" {
		t.Fatalf("unexpected version_id: %s", versions[0].VersionID)
	}
	if versions[0].State != WorkflowStateActive {
		t.Fatalf("unexpected state: %s", versions[0].State)
	}
	if versions[1].ExecutionStatus != ExecutionStateError {
		t.Fatalf("unexpected execution_status: %s", versions[1].ExecutionStatus)
	}
	if versions[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAutomationVersionsEmptyID(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.AutomationVersions(context.Background(), "", nil)
	if err == nil || err.Error() != "mgmt: workflowId is required" {
		t.Fatalf("expected workflowId required error, got %v", err)
	}
}

func TestAutomationExport(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflow-import-export/export/wf-1/v-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"name":        "Alert triage",
			"description": "Auto-triage",
			"actions": []map[string]any{
				{
					"action": map[string]any{
						"type": "manual_trigger",
						"tag":  "core_action",
					},
					"export_id":    1,
					"connected_to": []any{},
				},
			},
		})
	})
	c := testClient(t, handler)
	exp, err := c.AutomationExport(context.Background(), "wf-1", "v-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exp.Name != "Alert triage" {
		t.Fatalf("unexpected name: %s", exp.Name)
	}
	if len(exp.Actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(exp.Actions))
	}
	if exp.Actions[0].ExportID != 1 {
		t.Fatalf("unexpected export_id: %d", exp.Actions[0].ExportID)
	}
	if exp.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAutomationExportEmptyIDs(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.AutomationExport(context.Background(), "", "v-1")
	if err == nil || err.Error() != "mgmt: workflowId is required" {
		t.Fatalf("expected workflowId required error, got %v", err)
	}
	_, err = c.AutomationExport(context.Background(), "wf-1", "")
	if err == nil || err.Error() != "mgmt: versionId is required" {
		t.Fatalf("expected versionId required error, got %v", err)
	}
}

func TestAutomationImport(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflow-import-export/import" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   json.RawMessage `json:"data"`
			Filter *struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"filter"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data == nil {
			t.Fatal("expected data in body")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":              "wf-new",
			"version_id":      "v-1",
			"name":            "Imported workflow",
			"state":           "draft",
			"lifecycle_state": "active",
			"status":          "idle",
			"scope_id":        "000000",
			"scope_level":     "site",
			"mgmt_id":         "mgmt-1",
			"created_by":      "user-1",
			"updated_by":      "user-1",
		})
	})
	c := testClient(t, handler)
	data := json.RawMessage(`{"name":"Imported workflow","actions":[]}`)
	wf, err := c.AutomationImport(context.Background(), data, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wf.ID != "wf-new" {
		t.Fatalf("unexpected ID: %s", wf.ID)
	}
	if wf.Name != "Imported workflow" {
		t.Fatalf("unexpected name: %s", wf.Name)
	}
	if wf.State != WorkflowStateDraft {
		t.Fatalf("unexpected state: %s", wf.State)
	}
	if wf.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAutomationImportEmptyData(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.AutomationImport(context.Background(), nil, nil)
	if err == nil || err.Error() != "mgmt: workflow data is required" {
		t.Fatalf("expected data required error, got %v", err)
	}
}

func TestAutomationActivate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflows/wf-1/v-2/activation" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	c := testClient(t, handler)
	if err := c.AutomationActivate(context.Background(), "wf-1", "v-2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAutomationActivateEmptyIDs(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if err := c.AutomationActivate(context.Background(), "", "v-1"); err == nil {
		t.Fatal("expected error for empty workflowId")
	}
	if err := c.AutomationActivate(context.Background(), "wf-1", ""); err == nil {
		t.Fatal("expected error for empty versionId")
	}
}

func TestAutomationDeactivate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflows/wf-1/deactivate" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	c := testClient(t, handler)
	if err := c.AutomationDeactivate(context.Background(), "wf-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAutomationDeactivateEmptyID(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	if err := c.AutomationDeactivate(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty workflowId")
	}
}

func TestAutomationExecutions(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflow-execution" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("workflow_id"); got != "wf-1" {
			t.Fatalf("expected workflow_id=wf-1, got %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":                   "exec-1",
					"version_id":           "v-1",
					"workflow_id":          "wf-1",
					"state":                "Completed",
					"duration":             "PT5S",
					"executed_actions":     3,
					"has_execution_output": true,
					"scope_id":             "000000",
					"scope_level":          "site",
					"mgmt_id":              "mgmt-1",
					"created_at":           "2025-01-01T00:00:00Z",
					"workflow_name":        "Alert triage",
					"trigger":              "manual_trigger",
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": ""},
		})
	})
	c := testClient(t, handler)
	execs, pag, err := c.AutomationExecutions(context.Background(), &AutomationExecutionListParams{
		WorkflowID: "wf-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(execs) != 1 {
		t.Fatalf("expected 1 execution, got %d", len(execs))
	}
	if execs[0].State != ExecutionStateCompleted {
		t.Fatalf("unexpected state: %s", execs[0].State)
	}
	if execs[0].Trigger != TriggerManual {
		t.Fatalf("unexpected trigger: %s", execs[0].Trigger)
	}
	if !execs[0].HasExecutionOutput {
		t.Fatal("expected has_execution_output=true")
	}
	if execs[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestAutomationExecutionGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflow-execution/exec-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":                              "exec-1",
			"version_id":                      "v-1",
			"workflow_id":                     "wf-1",
			"state":                           "Error",
			"scope_id":                        "000000",
			"scope_level":                     "site",
			"mgmt_id":                         "mgmt-1",
			"created_at":                      "2025-01-01T00:00:00Z",
			"workflow_state":                  "active",
			"singularity_response_event_type": "alert",
			"singularity_response_event_id":   "alert-123",
			"error_actions": []map[string]any{
				{
					"action_id":             "act-5",
					"action_execution_name": "Send Email",
					"action_display_name":   "Notify team",
					"action_error":          "SMTP timeout",
				},
			},
		})
	})
	c := testClient(t, handler)
	detail, err := c.AutomationExecutionGet(context.Background(), "exec-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if detail.State != ExecutionStateError {
		t.Fatalf("unexpected state: %s", detail.State)
	}
	if detail.WorkflowState != WorkflowStateActive {
		t.Fatalf("unexpected workflow_state: %s", detail.WorkflowState)
	}
	if detail.SingularityResponseEventType != SingularityEventAlert {
		t.Fatalf("unexpected event type: %s", detail.SingularityResponseEventType)
	}
	if len(detail.ErrorActions) != 1 {
		t.Fatalf("expected 1 error action, got %d", len(detail.ErrorActions))
	}
	if detail.ErrorActions[0].ActionError != "SMTP timeout" {
		t.Fatalf("unexpected error: %s", detail.ErrorActions[0].ActionError)
	}
	if detail.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAutomationExecutionGetEmptyID(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.AutomationExecutionGet(context.Background(), "")
	if err == nil || err.Error() != "mgmt: executionId is required" {
		t.Fatalf("expected executionId required error, got %v", err)
	}
}

func TestAutomationExecutionOutput(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflow-execution/output/exec-1/raw" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"ExecutionOutput": map[string]any{"result": "ok"},
		})
	})
	c := testClient(t, handler)
	out, err := c.AutomationExecutionOutput(context.Background(), "exec-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatal("expected output to be populated")
	}
	var parsed map[string]string
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if parsed["result"] != "ok" {
		t.Fatalf("unexpected result: %s", parsed["result"])
	}
}

func TestAutomationExecutionOutputEmptyID(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.AutomationExecutionOutput(context.Background(), "")
	if err == nil || err.Error() != "mgmt: executionId is required" {
		t.Fatalf("expected executionId required error, got %v", err)
	}
}

func TestAutomationRun(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/hyper-automate/api/public/workflow-execution/manual/wf-1/v-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":          "exec-new",
			"version_id":  "v-1",
			"workflow_id": "wf-1",
			"state":       "Pending",
			"scope_id":    "000000",
			"scope_level": "site",
			"mgmt_id":     "mgmt-1",
			"created_at":  "2025-01-01T00:00:00Z",
		})
	})
	c := testClient(t, handler)
	exec, err := c.AutomationRun(context.Background(), "wf-1", "v-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.ID != "exec-new" {
		t.Fatalf("unexpected ID: %s", exec.ID)
	}
	if exec.State != ExecutionStatePending {
		t.Fatalf("unexpected state: %s", exec.State)
	}
	if exec.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestAutomationRunEmptyIDs(t *testing.T) {
	c := testClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	_, err := c.AutomationRun(context.Background(), "", "v-1", nil)
	if err == nil || err.Error() != "mgmt: workflowId is required" {
		t.Fatalf("expected workflowId required error, got %v", err)
	}
	_, err = c.AutomationRun(context.Background(), "wf-1", "", nil)
	if err == nil || err.Error() != "mgmt: versionId is required" {
		t.Fatalf("expected versionId required error, got %v", err)
	}
}
