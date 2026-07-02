package mgmt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

// RemoteScriptType is the category of a remote script.
type RemoteScriptType string

const (
	ScriptTypeArtifactCollection RemoteScriptType = "artifactCollection"
	ScriptTypeDataCollection     RemoteScriptType = "dataCollection"
	ScriptTypeAction             RemoteScriptType = "action"
)

// PackageEndpointExpiration controls when a script package is removed from the
// endpoint after execution.
type PackageEndpointExpiration string

const (
	PackageExpirationNone      PackageEndpointExpiration = "None"
	PackageExpirationImmediate PackageEndpointExpiration = "Immediate"
	PackageExpirationOnRestart PackageEndpointExpiration = "OnRestart"
	PackageExpirationTime      PackageEndpointExpiration = "Time"
)

// RemoteScriptUpdateData is the editable metadata of a remote script.
type RemoteScriptUpdateData struct {
	ScriptName                       string                    `json:"scriptName"`
	ScriptType                       RemoteScriptType          `json:"scriptType"`
	OSTypes                          []string                  `json:"osTypes"`
	InputRequired                    bool                      `json:"inputRequired"`
	InputExample                     string                    `json:"inputExample"`
	InputInstructions                string                    `json:"inputInstructions"`
	ScriptDescription                string                    `json:"scriptDescription,omitempty"`
	ScriptRuntimeTimeoutSeconds      int                       `json:"scriptRuntimeTimeoutSeconds"`
	PackageEndpointExpiration        PackageEndpointExpiration `json:"packageEndpointExpiration,omitempty"`
	PackageEndpointExpirationSeconds int                       `json:"packageEndpointExpirationSeconds,omitempty"`
}

// RemoteScriptUpdate is the body of a metadata update (PUT /remote-scripts/{id}).
// It changes a script's properties (name, timeout, input requirements) but not
// its content; use RemoteScriptsEdit to change the script body.
type RemoteScriptUpdate struct {
	ConsoleData  string                 `json:"consoleData,omitempty"`
	Data         RemoteScriptUpdateData `json:"data"`
	SendActivity *bool                  `json:"sendActivity,omitempty"`
}

// RemoteScriptsUpdate changes the metadata of an existing remote script.
func (c *Client) RemoteScriptsUpdate(ctx context.Context, id string, upd RemoteScriptUpdate) (*RemoteScript, error) {
	if id == "" {
		return nil, fmt.Errorf("mgmt: script id is required")
	}
	var resp singleResponse[RemoteScript]
	if err := c.put(ctx, "/remote-scripts/"+url.PathEscape(id), upd, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RemoteScriptEdit is the full editable definition of a remote script, used to
// change the script content along with its metadata via the multipart
// /remote-scripts/edit/{id} endpoint. Content is supplied inline via
// ScriptContent (set ScriptContentEncoded when it is base64).
type RemoteScriptEdit struct {
	ScriptName                       string
	ScriptType                       RemoteScriptType
	OSTypes                          []string
	InputRequired                    bool
	InputExample                     string
	InputInstructions                string
	ScriptRuntimeTimeoutSeconds      int
	ScriptContent                    string
	ScriptContentEncoded             bool
	ScriptDescription                string
	ConsoleData                      string
	SendActivity                     *bool
	PackageRemoved                   bool
	PackageMaxSize                   string
	PackageEndpointExpiration        PackageEndpointExpiration
	PackageEndpointExpirationSeconds int
}

// RemoteScriptsEdit changes a script's content and metadata. The endpoint takes
// multipart/form-data; content is sent inline as the scriptContent field.
func (c *Client) RemoteScriptsEdit(ctx context.Context, id string, edit RemoteScriptEdit) (*RemoteScript, error) {
	if id == "" {
		return nil, fmt.Errorf("mgmt: script id is required")
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	fields := []struct {
		name, val string
	}{
		{"scriptName", edit.ScriptName},
		{"scriptType", string(edit.ScriptType)},
		{"inputExample", edit.InputExample},
		{"inputInstructions", edit.InputInstructions},
		{"inputRequired", strconv.FormatBool(edit.InputRequired)},
		{"scriptRuntimeTimeoutSeconds", strconv.Itoa(edit.ScriptRuntimeTimeoutSeconds)},
	}
	for _, f := range fields {
		if err := w.WriteField(f.name, f.val); err != nil {
			return nil, fmt.Errorf("mgmt: write field %s: %w", f.name, err)
		}
	}
	for _, os := range edit.OSTypes {
		if err := w.WriteField("osTypes", os); err != nil {
			return nil, fmt.Errorf("mgmt: write field osTypes: %w", err)
		}
	}
	optional := map[string]string{
		"scriptContent":             edit.ScriptContent,
		"scriptDescription":         edit.ScriptDescription,
		"consoleData":               edit.ConsoleData,
		"packageMaxSize":            edit.PackageMaxSize,
		"packageEndpointExpiration": string(edit.PackageEndpointExpiration),
	}
	for name, val := range optional {
		if val == "" {
			continue
		}
		if err := w.WriteField(name, val); err != nil {
			return nil, fmt.Errorf("mgmt: write field %s: %w", name, err)
		}
	}
	if edit.ScriptContent != "" {
		if err := w.WriteField("isScriptContentEncoded", strconv.FormatBool(edit.ScriptContentEncoded)); err != nil {
			return nil, fmt.Errorf("mgmt: write field isScriptContentEncoded: %w", err)
		}
	}
	if edit.PackageRemoved {
		if err := w.WriteField("packageRemoved", "true"); err != nil {
			return nil, fmt.Errorf("mgmt: write field packageRemoved: %w", err)
		}
	}
	if edit.PackageEndpointExpirationSeconds > 0 {
		if err := w.WriteField("packageEndpointExpirationSeconds", strconv.Itoa(edit.PackageEndpointExpirationSeconds)); err != nil {
			return nil, fmt.Errorf("mgmt: write field packageEndpointExpirationSeconds: %w", err)
		}
	}
	if edit.SendActivity != nil {
		if err := w.WriteField("sendActivity", strconv.FormatBool(*edit.SendActivity)); err != nil {
			return nil, fmt.Errorf("mgmt: write field sendActivity: %w", err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("mgmt: close multipart: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseURL+"/remote-scripts/edit/"+url.PathEscape(id), &body)
	if err != nil {
		return nil, fmt.Errorf("mgmt: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	var resp singleResponse[RemoteScript]
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RemoteScriptContent returns the raw text content of a remote script.
func (c *Client) RemoteScriptContent(ctx context.Context, id string) (string, error) {
	if id == "" {
		return "", fmt.Errorf("mgmt: script id is required")
	}
	v := url.Values{}
	v.Set("scriptId", id)
	var resp singleResponse[struct {
		ScriptContent string `json:"scriptContent"`
	}]
	if err := c.get(ctx, "/remote-scripts/script-content", v, &resp); err != nil {
		return "", err
	}
	return resp.Data.ScriptContent, nil
}

// UploadLimits holds the package upload limits for remote scripts. The API
// response shape is not fixed by the spec, so the data object is preserved as
// raw JSON.
type UploadLimits struct {
	Raw json.RawMessage
}

// MarshalJSON emits the preserved upload-limits data object.
func (u UploadLimits) MarshalJSON() ([]byte, error) {
	if len(u.Raw) == 0 {
		return []byte("null"), nil
	}
	return u.Raw, nil
}

// RemoteScriptsUploadLimits returns the package upload size limits.
func (c *Client) RemoteScriptsUploadLimits(ctx context.Context) (*UploadLimits, error) {
	var resp struct {
		Data json.RawMessage `json:"data"`
	}
	if err := c.get(ctx, "/remote-scripts/fetch-upload-limits", nil, &resp); err != nil {
		return nil, err
	}
	return &UploadLimits{Raw: resp.Data}, nil
}

// PendingExecutionState is the review state of a pending remote-script execution.
type PendingExecutionState string

const (
	PendingStateWaiting  PendingExecutionState = "waiting"
	PendingStateApproved PendingExecutionState = "approved"
	PendingStateDeclined PendingExecutionState = "declined"
	PendingStateExpired  PendingExecutionState = "expired"
)

// PendingExecution is a remote-script execution awaiting approval.
type PendingExecution struct {
	PendingExecutionID  string                `json:"pendingExecutionId"`
	State               PendingExecutionState `json:"state"`
	CreatedAt           string                `json:"createdAt"`
	Creator             string                `json:"creator"`
	CreatorID           string                `json:"creatorId"`
	Reviewer            string                `json:"reviewer"`
	TotalEndpoints      int                   `json:"totalEndpoints"`
	CanApproveOrDecline bool                  `json:"canApproveOrDecline"`
	ScriptData          struct {
		ID         string `json:"id"`
		ScriptName string `json:"scriptName"`
		ScriptType string `json:"scriptType"`
	} `json:"scriptData"`
	ExecutionData struct {
		ScriptID          string `json:"scriptId"`
		TaskDescription   string `json:"taskDescription"`
		OutputDestination string `json:"outputDestination"`
	} `json:"executionData"`

	Raw json.RawMessage `json:"-"`
}

func (p *PendingExecution) UnmarshalJSON(b []byte) error {
	type alias PendingExecution
	if err := json.Unmarshal(b, (*alias)(p)); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

// RemoteScriptsPendingParams are query parameters for listing pending executions.
type RemoteScriptsPendingParams struct {
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
	SortBy     string
	SortOrder  string
	Limit      int
	Cursor     string
}

func (p *RemoteScriptsPendingParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// RemoteScriptsPendingList returns a paginated list of pending executions.
func (c *Client) RemoteScriptsPendingList(ctx context.Context, params *RemoteScriptsPendingParams) ([]PendingExecution, *Pagination, error) {
	return list[PendingExecution](c, ctx, "/remote-scripts/pending-executions", params.values())
}

// PendingExecutionAction is the decision applied to a pending execution.
type PendingExecutionAction string

const (
	PendingActionApprove PendingExecutionAction = "approve"
	PendingActionDecline PendingExecutionAction = "decline"
)

// RemoteScriptsPendingDecision approves or declines a pending execution.
func (c *Client) RemoteScriptsPendingDecision(ctx context.Context, id string, approve bool) error {
	if id == "" {
		return fmt.Errorf("mgmt: pending execution id is required")
	}
	action := PendingActionDecline
	if approve {
		action = PendingActionApprove
	}
	req := map[string]any{"data": map[string]any{"action": action}}
	return c.put(ctx, "/remote-scripts/pending-executions/"+url.PathEscape(id), req, nil)
}

// GuardrailScopeLevel is the scope at which a guardrail applies.
type GuardrailScopeLevel string

const (
	GuardrailScopeAccount GuardrailScopeLevel = "account"
	GuardrailScopeSite    GuardrailScopeLevel = "site"
	GuardrailScopeGroup   GuardrailScopeLevel = "group"
)

// GuardrailScope identifies a guardrail configuration by scope.
type GuardrailScope struct {
	ScopeID    string
	ScopeLevel GuardrailScopeLevel
}

func (s GuardrailScope) validate() error {
	if s.ScopeID == "" {
		return fmt.Errorf("mgmt: scopeId is required")
	}
	if s.ScopeLevel == "" {
		return fmt.Errorf("mgmt: scopeLevel is required")
	}
	return nil
}

// Guardrails is a remote-script guardrail configuration for a scope. A guardrail
// requires approval before scripts of the listed types run on more than
// EndpointsQuantity endpoints.
type Guardrails struct {
	EndpointsQuantity *int     `json:"endpointsQuantity"`
	ScriptTypes       []string `json:"scriptTypes"`
	Inherited         bool     `json:"inherited"`
	Enabled           bool     `json:"enabled"`

	Raw json.RawMessage `json:"-"`
}

func (g *Guardrails) UnmarshalJSON(b []byte) error {
	type alias Guardrails
	if err := json.Unmarshal(b, (*alias)(g)); err != nil {
		return err
	}
	g.Raw = append(g.Raw[:0:0], b...)
	return nil
}

// GuardrailsGet returns the guardrail configuration for a scope.
func (c *Client) GuardrailsGet(ctx context.Context, scope GuardrailScope) (*Guardrails, error) {
	if err := scope.validate(); err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("scopeId", scope.ScopeID)
	v.Set("scopeLevel", string(scope.ScopeLevel))
	var resp singleResponse[Guardrails]
	if err := c.get(ctx, "/remote-scripts/guardrails/configuration", v, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GuardrailsUpsertInput is the body for creating or updating a guardrail. The
// endpointsQuantity threshold is required by the spec and may be null.
type GuardrailsUpsertInput struct {
	ScopeID           string              `json:"scopeId"`
	ScopeLevel        GuardrailScopeLevel `json:"scopeLevel"`
	EndpointsQuantity *int                `json:"endpointsQuantity"`
	ScriptTypes       []string            `json:"scriptTypes"`
	Enabled           bool                `json:"enabled"`
}

// GuardrailsUpsert creates or updates (if it does not exist) a guardrail
// configuration for a scope.
func (c *Client) GuardrailsUpsert(ctx context.Context, in GuardrailsUpsertInput) error {
	if in.ScopeID == "" || in.ScopeLevel == "" {
		return fmt.Errorf("mgmt: scopeId and scopeLevel are required")
	}
	req := map[string]any{"data": in}
	return c.post(ctx, "/remote-scripts/guardrails/configuration", req, nil)
}

// GuardrailsDelete removes a guardrail configuration for a scope.
func (c *Client) GuardrailsDelete(ctx context.Context, scope GuardrailScope) error {
	if err := scope.validate(); err != nil {
		return err
	}
	req := map[string]any{"data": map[string]any{
		"scopeId":    scope.ScopeID,
		"scopeLevel": scope.ScopeLevel,
	}}
	return c.jsonRequest(ctx, http.MethodDelete, "/remote-scripts/guardrails/configuration", req, nil)
}

// GuardrailCheckInput is the body for a guardrail check.
type GuardrailCheckInput struct {
	ScriptID string   `json:"scriptId"`
	AgentIDs []string `json:"agentIds"`
}

// GuardrailCheckResult reports whether a guardrail requires approval for an
// execution.
type GuardrailCheckResult struct {
	RequiresApproval bool `json:"requiresApproval"`

	Raw json.RawMessage `json:"-"`
}

func (r *GuardrailCheckResult) UnmarshalJSON(b []byte) error {
	type alias GuardrailCheckResult
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// GuardrailsCheck reports whether running a script on the given agents would
// trip a guardrail and require approval.
func (c *Client) GuardrailsCheck(ctx context.Context, in GuardrailCheckInput) (*GuardrailCheckResult, error) {
	if in.ScriptID == "" {
		return nil, fmt.Errorf("mgmt: scriptId is required")
	}
	if len(in.AgentIDs) == 0 {
		return nil, fmt.Errorf("mgmt: at least one agentId is required")
	}
	req := map[string]any{"data": in}
	var resp singleResponse[GuardrailCheckResult]
	if err := c.post(ctx, "/remote-scripts/guardrails/check", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
