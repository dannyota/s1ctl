package mgmt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
)

func TestRemoteScriptsUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/remote-scripts/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			ConsoleData  string `json:"consoleData"`
			SendActivity *bool  `json:"sendActivity"`
			Data         struct {
				ScriptName                  string   `json:"scriptName"`
				ScriptType                  string   `json:"scriptType"`
				OSTypes                     []string `json:"osTypes"`
				InputRequired               bool     `json:"inputRequired"`
				ScriptRuntimeTimeoutSeconds int      `json:"scriptRuntimeTimeoutSeconds"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Data.ScriptName != "Collect Logs" {
			t.Fatalf("unexpected scriptName: %s", body.Data.ScriptName)
		}
		if body.Data.ScriptType != "dataCollection" {
			t.Fatalf("unexpected scriptType: %s", body.Data.ScriptType)
		}
		if !slices.Equal(body.Data.OSTypes, []string{"linux", "macos"}) {
			t.Fatalf("unexpected osTypes: %v", body.Data.OSTypes)
		}
		if body.Data.ScriptRuntimeTimeoutSeconds != 3600 {
			t.Fatalf("unexpected timeout: %d", body.Data.ScriptRuntimeTimeoutSeconds)
		}
		if body.SendActivity == nil || *body.SendActivity {
			t.Fatalf("expected sendActivity=false, got %v", body.SendActivity)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "1000000000000000001", "fileName": "collect_logs.sh"},
		})
	})
	c := testClient(t, handler)
	no := false
	updated, err := c.RemoteScriptsUpdate(context.Background(), "1000000000000000001", RemoteScriptUpdate{
		SendActivity: &no,
		Data: RemoteScriptUpdateData{
			ScriptName:                  "Collect Logs",
			ScriptType:                  ScriptTypeDataCollection,
			OSTypes:                     []string{"linux", "macos"},
			InputRequired:               false,
			InputExample:                "-",
			InputInstructions:           "-",
			ScriptRuntimeTimeoutSeconds: 3600,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", updated.ID)
	}
}

func TestRemoteScriptsUpdateRequiresID(t *testing.T) {
	c := testClient(t, http.NotFoundHandler())
	if _, err := c.RemoteScriptsUpdate(context.Background(), "", RemoteScriptUpdate{}); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestRemoteScriptsEdit(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/remote-scripts/edit/1000000000000000002" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "multipart/form-data") {
			t.Fatalf("expected multipart, got %s", ct)
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if got := r.FormValue("scriptName"); got != "Collect Logs" {
			t.Fatalf("unexpected scriptName: %s", got)
		}
		if got := r.FormValue("scriptContent"); got != "echo hi" {
			t.Fatalf("unexpected scriptContent: %s", got)
		}
		if got := r.Form["osTypes"]; !slices.Equal(got, []string{"linux"}) {
			t.Fatalf("unexpected osTypes: %v", got)
		}
		if got := r.FormValue("scriptRuntimeTimeoutSeconds"); got != "120" {
			t.Fatalf("unexpected timeout: %s", got)
		}
		if got := r.FormValue("inputRequired"); got != "false" {
			t.Fatalf("unexpected inputRequired: %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "1000000000000000002"},
		})
	})
	c := testClient(t, handler)
	updated, err := c.RemoteScriptsEdit(context.Background(), "1000000000000000002", RemoteScriptEdit{
		ScriptName:                  "Collect Logs",
		ScriptType:                  ScriptTypeDataCollection,
		OSTypes:                     []string{"linux"},
		InputExample:                "-",
		InputInstructions:           "-",
		ScriptRuntimeTimeoutSeconds: 120,
		ScriptContent:               "echo hi",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", updated.ID)
	}
}

func TestRemoteScriptContent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/remote-scripts/script-content" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("scriptId"); got != "1000000000000000003" {
			t.Fatalf("unexpected scriptId: %s", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"scriptContent": "#!/bin/sh\necho hi\n"},
		})
	})
	c := testClient(t, handler)
	content, err := c.RemoteScriptContent(context.Background(), "1000000000000000003")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != "#!/bin/sh\necho hi\n" {
		t.Fatalf("unexpected content: %q", content)
	}
}

func TestRemoteScriptContentRequiresID(t *testing.T) {
	c := testClient(t, http.NotFoundHandler())
	if _, err := c.RemoteScriptContent(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestRemoteScriptsUploadLimits(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/remote-scripts/fetch-upload-limits" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"maxUploadSize": 104857600, "unit": "bytes"},
		})
	})
	c := testClient(t, handler)
	limits, err := c.RemoteScriptsUploadLimits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(limits.Raw), "maxUploadSize") {
		t.Fatalf("expected raw to hold data, got %s", limits.Raw)
	}
	out, err := json.Marshal(limits)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(out), "maxUploadSize") {
		t.Fatalf("expected marshalled limits to hold data, got %s", out)
	}
}

func TestRemoteScriptsPendingList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/remote-scripts/pending-executions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"878572631641628675"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if q.Get("limit") != "25" {
			t.Fatalf("unexpected limit: %s", q.Get("limit"))
		}
		if q.Get("sortBy") != "createdAt" {
			t.Fatalf("unexpected sortBy: %s", q.Get("sortBy"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"pendingExecutionId":  "2000000000000000001",
					"state":               "waiting",
					"createdAt":           "2025-01-01T00:00:00Z",
					"creator":             "Jane",
					"totalEndpoints":      12,
					"canApproveOrDecline": true,
					"scriptData":          map[string]any{"id": "3000000000000000001", "scriptName": "Collect Logs", "scriptType": "dataCollection"},
					"executionData":       map[string]any{"scriptId": "3000000000000000001", "taskDescription": "IR sweep", "outputDestination": "SentinelCloud"},
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "next"},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.RemoteScriptsPendingList(context.Background(), &RemoteScriptsPendingParams{
		SiteIDs: []string{"878572631641628675"},
		Limit:   25,
		SortBy:  "createdAt",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	it := items[0]
	if it.PendingExecutionID != "2000000000000000001" {
		t.Fatalf("unexpected id: %s", it.PendingExecutionID)
	}
	if it.State != PendingStateWaiting {
		t.Fatalf("unexpected state: %s", it.State)
	}
	if it.ScriptData.ScriptName != "Collect Logs" {
		t.Fatalf("unexpected scriptName: %s", it.ScriptData.ScriptName)
	}
	if it.ExecutionData.TaskDescription != "IR sweep" {
		t.Fatalf("unexpected task: %s", it.ExecutionData.TaskDescription)
	}
	if !it.CanApproveOrDecline {
		t.Fatal("expected canApproveOrDecline=true")
	}
	if it.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 || pag.NextCursor != "next" {
		t.Fatalf("unexpected pagination: %+v", pag)
	}
}

func TestRemoteScriptsPendingDecisionApprove(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/remote-scripts/pending-executions/2000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				Action string `json:"action"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Data.Action != "approve" {
			t.Fatalf("unexpected action: %s", body.Data.Action)
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	if err := c.RemoteScriptsPendingDecision(context.Background(), "2000000000000000001", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoteScriptsPendingDecisionDecline(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Data struct {
				Action string `json:"action"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Data.Action != "decline" {
			t.Fatalf("unexpected action: %s", body.Data.Action)
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	if err := c.RemoteScriptsPendingDecision(context.Background(), "2000000000000000001", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoteScriptsPendingDecisionRequiresID(t *testing.T) {
	c := testClient(t, http.NotFoundHandler())
	if err := c.RemoteScriptsPendingDecision(context.Background(), "", true); err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestGuardrailsGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/remote-scripts/guardrails/configuration" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("scopeId") != "878572631641628675" {
			t.Fatalf("unexpected scopeId: %s", q.Get("scopeId"))
		}
		if q.Get("scopeLevel") != "site" {
			t.Fatalf("unexpected scopeLevel: %s", q.Get("scopeLevel"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"endpointsQuantity": 100,
				"scriptTypes":       []string{"action"},
				"inherited":         false,
				"enabled":           true,
			},
		})
	})
	c := testClient(t, handler)
	g, err := c.GuardrailsGet(context.Background(), GuardrailScope{ScopeID: "878572631641628675", ScopeLevel: GuardrailScopeSite})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.EndpointsQuantity == nil || *g.EndpointsQuantity != 100 {
		t.Fatalf("unexpected endpointsQuantity: %v", g.EndpointsQuantity)
	}
	if !g.Enabled {
		t.Fatal("expected enabled=true")
	}
	if !slices.Equal(g.ScriptTypes, []string{"action"}) {
		t.Fatalf("unexpected scriptTypes: %v", g.ScriptTypes)
	}
}

func TestGuardrailsGetRequiresScope(t *testing.T) {
	c := testClient(t, http.NotFoundHandler())
	if _, err := c.GuardrailsGet(context.Background(), GuardrailScope{ScopeLevel: GuardrailScopeSite}); err == nil {
		t.Fatal("expected error for missing scopeId")
	}
}

func TestGuardrailsUpsert(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/remote-scripts/guardrails/configuration" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				ScopeID           string   `json:"scopeId"`
				ScopeLevel        string   `json:"scopeLevel"`
				EndpointsQuantity *int     `json:"endpointsQuantity"`
				ScriptTypes       []string `json:"scriptTypes"`
				Enabled           bool     `json:"enabled"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Data.ScopeID != "878572631641628675" {
			t.Fatalf("unexpected scopeId: %s", body.Data.ScopeID)
		}
		if body.Data.ScopeLevel != "site" {
			t.Fatalf("unexpected scopeLevel: %s", body.Data.ScopeLevel)
		}
		if body.Data.EndpointsQuantity == nil || *body.Data.EndpointsQuantity != 50 {
			t.Fatalf("unexpected endpointsQuantity: %v", body.Data.EndpointsQuantity)
		}
		if !body.Data.Enabled {
			t.Fatal("expected enabled=true")
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	q := 50
	if err := c.GuardrailsUpsert(context.Background(), GuardrailsUpsertInput{
		ScopeID:           "878572631641628675",
		ScopeLevel:        GuardrailScopeSite,
		EndpointsQuantity: &q,
		ScriptTypes:       []string{"action"},
		Enabled:           true,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGuardrailsDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/remote-scripts/guardrails/configuration" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		data, _ := io.ReadAll(r.Body)
		var body struct {
			Data struct {
				ScopeID    string `json:"scopeId"`
				ScopeLevel string `json:"scopeLevel"`
			} `json:"data"`
		}
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Data.ScopeID != "878572631641628675" || body.Data.ScopeLevel != "group" {
			t.Fatalf("unexpected delete body: %s", data)
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	if err := c.GuardrailsDelete(context.Background(), GuardrailScope{ScopeID: "878572631641628675", ScopeLevel: GuardrailScopeGroup}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGuardrailsCheck(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/remote-scripts/guardrails/check" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				ScriptID string   `json:"scriptId"`
				AgentIDs []string `json:"agentIds"`
			} `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Data.ScriptID != "3000000000000000001" {
			t.Fatalf("unexpected scriptId: %s", body.Data.ScriptID)
		}
		if !slices.Equal(body.Data.AgentIDs, []string{"4000000000000000001"}) {
			t.Fatalf("unexpected agentIds: %v", body.Data.AgentIDs)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"requiresApproval": true},
		})
	})
	c := testClient(t, handler)
	res, err := c.GuardrailsCheck(context.Background(), GuardrailCheckInput{
		ScriptID: "3000000000000000001",
		AgentIDs: []string{"4000000000000000001"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.RequiresApproval {
		t.Fatal("expected requiresApproval=true")
	}
}

func TestGuardrailsCheckRequiresInput(t *testing.T) {
	c := testClient(t, http.NotFoundHandler())
	if _, err := c.GuardrailsCheck(context.Background(), GuardrailCheckInput{ScriptID: "3000000000000000001"}); err == nil {
		t.Fatal("expected error for missing agentIds")
	}
}
