package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Hyperautomation base path (relative to /web/api/v2.1).
const automationBase = "/hyper-automate/api/public"

// --- Typed enums ---

// WorkflowState is the state of a workflow (active/inactive version state).
type WorkflowState string

const (
	WorkflowStateActive      WorkflowState = "active"
	WorkflowStateInactive    WorkflowState = "inactive"
	WorkflowStateDeactivated WorkflowState = "deactivated"
	WorkflowStateDraft       WorkflowState = "draft"
)

// WorkflowLifecycleState is the lifecycle state of a workflow.
type WorkflowLifecycleState string

const (
	WorkflowLifecycleActive   WorkflowLifecycleState = "active"
	WorkflowLifecycleArchived WorkflowLifecycleState = "archived"
	WorkflowLifecycleDeleted  WorkflowLifecycleState = "deleted"
)

// WorkflowStatus is the run status of a workflow.
type WorkflowStatus string

const (
	WorkflowStatusIdle    WorkflowStatus = "idle"
	WorkflowStatusRunning WorkflowStatus = "running"
)

// ExecutionState is the state of a workflow execution.
type ExecutionState string

const (
	ExecutionStateRunning             ExecutionState = "Running"
	ExecutionStatePending             ExecutionState = "Pending"
	ExecutionStateStuck               ExecutionState = "Stuck"
	ExecutionStateCompleted           ExecutionState = "Completed"
	ExecutionStateError               ExecutionState = "Error"
	ExecutionStateWaiting             ExecutionState = "Waiting"
	ExecutionStateAborted             ExecutionState = "Aborted"
	ExecutionStateCompletedWithErrors ExecutionState = "CompletedWithErrors"
)

// TriggerType is the type of workflow trigger.
type TriggerType string

const (
	TriggerHTTP                TriggerType = "http_trigger"
	TriggerScheduled           TriggerType = "scheduled_trigger"
	TriggerEmail               TriggerType = "email_trigger"
	TriggerManual              TriggerType = "manual_trigger"
	TriggerSingularityResponse TriggerType = "singularity_response_trigger"
	TriggerSnippet             TriggerType = "snippet_trigger"
)

// AutomationScopeLevel is the scope level for a workflow.
type AutomationScopeLevel string

const (
	AutomationScopeTenant  AutomationScopeLevel = "tenant"
	AutomationScopeAccount AutomationScopeLevel = "account"
	AutomationScopeSite    AutomationScopeLevel = "site"
)

// SingularityEventType is the event type for singularity response triggers.
type SingularityEventType string

const (
	SingularityEventAlert            SingularityEventType = "alert"
	SingularityEventIncident         SingularityEventType = "incident"
	SingularityEventMisconfiguration SingularityEventType = "misconfiguration"
	SingularityEventVulnerability    SingularityEventType = "vulnerability"
	SingularityEventActivity         SingularityEventType = "activity"
)

// ExecutionSource is the source of an execution trigger.
type ExecutionSource string

const (
	ExecutionSourceAutomatic ExecutionSource = "automatic"
	ExecutionSourceOnDemand  ExecutionSource = "on_demand"
	ExecutionSourceRerun     ExecutionSource = "rerun"
)

// --- Response structs ---

// WorkflowUser is a user reference in a workflow.
type WorkflowUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// WorkflowDimensions is the canvas dimensions of a workflow.
type WorkflowDimensions struct {
	Width  *float64 `json:"width"`
	Height *float64 `json:"height"`
}

// WorkflowAction is a summary action in a workflow listing.
type WorkflowAction struct {
	ID            string `json:"id"`
	IntegrationID string `json:"integration_id"`
	Type          string `json:"type"`
}

// Workflow is the full workflow object returned by list and import.
type Workflow struct {
	ID                 string                 `json:"id"`
	VersionID          string                 `json:"version_id"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	State              WorkflowState          `json:"state"`
	LifecycleState     WorkflowLifecycleState `json:"lifecycle_state"`
	Status             WorkflowStatus         `json:"status"`
	ScopeID            string                 `json:"scope_id"`
	ScopeLevel         AutomationScopeLevel   `json:"scope_level"`
	MgmtID             string                 `json:"mgmt_id"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
	ActivatedAt        string                 `json:"activated_at"`
	CreatedBy          string                 `json:"created_by"`
	UpdatedBy          string                 `json:"updated_by"`
	CreatedByUser      *WorkflowUser          `json:"created_by_user"`
	UpdatedByUser      *WorkflowUser          `json:"updated_by_user"`
	Tags               []string               `json:"tags"`
	SiteName           string                 `json:"site_name"`
	AccountName        string                 `json:"account_name"`
	ParentScopeID      string                 `json:"parent_scope_id"`
	SiteState          string                 `json:"site_state"`
	AccountState       string                 `json:"account_state"`
	VersionDescription string                 `json:"version_description"`
	VersionCount       int                    `json:"version_count"`
	Timeout            int                    `json:"timeout"`
	DailyMaxExecutions int                    `json:"daily_max_executions"`
	MaxConcurrency     int                    `json:"max_concurrency"`
	NotifyTo           []string               `json:"notify_to"`
	TimeSaved          int                    `json:"time_saved"`
	TimeSavedUnit      string                 `json:"time_saved_unit"`
	IsSnippet          bool                   `json:"is_snippet"`
	Dimensions         *WorkflowDimensions    `json:"dimensions"`
	AvailableUntil     string                 `json:"available_until"`
	Raw                json.RawMessage        `json:"-"`
}

func (w *Workflow) UnmarshalJSON(b []byte) error {
	type alias Workflow
	if err := json.Unmarshal(b, (*alias)(w)); err != nil {
		return err
	}
	w.Raw = append(w.Raw[:0:0], b...)
	return nil
}

// WorkflowListItem wraps a workflow with its actions in list responses.
type WorkflowListItem struct {
	ID       string           `json:"id"`
	Workflow Workflow         `json:"workflow"`
	Actions  []WorkflowAction `json:"actions"`
	Raw      json.RawMessage  `json:"-"`
}

func (w *WorkflowListItem) UnmarshalJSON(b []byte) error {
	type alias WorkflowListItem
	if err := json.Unmarshal(b, (*alias)(w)); err != nil {
		return err
	}
	w.Raw = append(w.Raw[:0:0], b...)
	return nil
}

// WorkflowVersion is a version entry returned by the versions endpoint.
type WorkflowVersion struct {
	ID                 string                 `json:"id"`
	VersionID          string                 `json:"version_id"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	State              WorkflowState          `json:"state"`
	LifecycleState     WorkflowLifecycleState `json:"lifecycle_state"`
	Status             WorkflowStatus         `json:"status"`
	ScopeID            string                 `json:"scope_id"`
	ScopeLevel         AutomationScopeLevel   `json:"scope_level"`
	MgmtID             string                 `json:"mgmt_id"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
	ActivatedAt        string                 `json:"activated_at"`
	CreatedBy          string                 `json:"created_by"`
	UpdatedBy          string                 `json:"updated_by"`
	CreatedByUser      *WorkflowUser          `json:"created_by_user"`
	UpdatedByUser      *WorkflowUser          `json:"updated_by_user"`
	VersionDescription string                 `json:"version_description"`
	VersionCount       int                    `json:"version_count"`
	Timeout            int                    `json:"timeout"`
	DailyMaxExecutions int                    `json:"daily_max_executions"`
	MaxConcurrency     int                    `json:"max_concurrency"`
	NotifyTo           []string               `json:"notify_to"`
	TimeSaved          int                    `json:"time_saved"`
	TimeSavedUnit      string                 `json:"time_saved_unit"`
	IsSnippet          bool                   `json:"is_snippet"`
	Dimensions         *WorkflowDimensions    `json:"dimensions"`
	ExecutionTime      string                 `json:"execution_time"`
	ExecutionStatus    ExecutionState         `json:"execution_status"`
	Raw                json.RawMessage        `json:"-"`
}

func (w *WorkflowVersion) UnmarshalJSON(b []byte) error {
	type alias WorkflowVersion
	if err := json.Unmarshal(b, (*alias)(w)); err != nil {
		return err
	}
	w.Raw = append(w.Raw[:0:0], b...)
	return nil
}

// WorkflowExport is the exported representation of a workflow version.
type WorkflowExport struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Actions     []WorkflowExportEntry `json:"actions"`
	Raw         json.RawMessage       `json:"-"`
}

func (w *WorkflowExport) UnmarshalJSON(b []byte) error {
	type alias WorkflowExport
	if err := json.Unmarshal(b, (*alias)(w)); err != nil {
		return err
	}
	w.Raw = append(w.Raw[:0:0], b...)
	return nil
}

// WorkflowExportEntry is an action entry in an exported workflow.
type WorkflowExportEntry struct {
	Action       json.RawMessage      `json:"action"`
	ExportID     int                  `json:"export_id"`
	ConnectedTo  []WorkflowExportEdge `json:"connected_to"`
	ParentAction *int                 `json:"parent_action"`
}

// WorkflowExportEdge is an edge between actions in an exported workflow.
type WorkflowExportEdge struct {
	Target       int             `json:"target"`
	CustomHandle json.RawMessage `json:"custom_handle"`
	Payload      json.RawMessage `json:"payload"`
}

// --- Execution structs ---

// WorkflowExecution is a workflow execution record.
type WorkflowExecution struct {
	ID                  string               `json:"id"`
	VersionID           string               `json:"version_id"`
	WorkflowID          string               `json:"workflow_id"`
	State               ExecutionState       `json:"state"`
	OffloadState        string               `json:"offload_state"`
	Duration            string               `json:"duration"`
	TimeSaved           float64              `json:"time_saved"`
	ExecutedActions     int                  `json:"executed_actions"`
	HasExecutionOutput  bool                 `json:"has_execution_output"`
	ScopeID             string               `json:"scope_id"`
	ScopeLevel          AutomationScopeLevel `json:"scope_level"`
	MgmtID              string               `json:"mgmt_id"`
	CreatedAt           string               `json:"created_at"`
	UpdatedAt           string               `json:"updated_at"`
	WorkflowName        string               `json:"workflow_name"`
	WorkflowDescription string               `json:"workflow_description"`
	WorkflowTags        []string             `json:"workflow_tags"`
	VersionCount        int                  `json:"version_count"`
	Trigger             TriggerType          `json:"trigger"`
	SiteName            string               `json:"site_name"`
	AccountName         string               `json:"account_name"`
	ParentScopeID       string               `json:"parent_scope_id"`
	SiteState           string               `json:"site_state"`
	AccountState        string               `json:"account_state"`
	Raw                 json.RawMessage      `json:"-"`
}

func (e *WorkflowExecution) UnmarshalJSON(b []byte) error {
	type alias WorkflowExecution
	if err := json.Unmarshal(b, (*alias)(e)); err != nil {
		return err
	}
	e.Raw = append(e.Raw[:0:0], b...)
	return nil
}

// WorkflowExecutionDetail is the detailed execution returned by get-by-ID.
type WorkflowExecutionDetail struct {
	ID                           string                 `json:"id"`
	VersionID                    string                 `json:"version_id"`
	WorkflowID                   string                 `json:"workflow_id"`
	State                        ExecutionState         `json:"state"`
	Duration                     string                 `json:"duration"`
	TimeSaved                    float64                `json:"time_saved"`
	ExecutedActions              int                    `json:"executed_actions"`
	HasExecutionOutput           bool                   `json:"has_execution_output"`
	ScopeID                      string                 `json:"scope_id"`
	ScopeLevel                   AutomationScopeLevel   `json:"scope_level"`
	MgmtID                       string                 `json:"mgmt_id"`
	CreatedAt                    string                 `json:"created_at"`
	UpdatedAt                    string                 `json:"updated_at"`
	WorkflowState                WorkflowState          `json:"workflow_state"`
	SingularityResponseEventType SingularityEventType   `json:"singularity_response_event_type"`
	SingularityResponseEventID   string                 `json:"singularity_response_event_id"`
	ErrorActions                 []ExecutionErrorAction `json:"error_actions"`
	Raw                          json.RawMessage        `json:"-"`
}

func (e *WorkflowExecutionDetail) UnmarshalJSON(b []byte) error {
	type alias WorkflowExecutionDetail
	if err := json.Unmarshal(b, (*alias)(e)); err != nil {
		return err
	}
	e.Raw = append(e.Raw[:0:0], b...)
	return nil
}

// ExecutionErrorAction describes a failed action within an execution.
type ExecutionErrorAction struct {
	ActionID            string `json:"action_id"`
	ActionExecutionName string `json:"action_execution_name"`
	ActionDisplayName   string `json:"action_display_name"`
	ActionError         string `json:"action_error"`
}

// WorkflowExecutionRun is the response from triggering a workflow execution.
type WorkflowExecutionRun struct {
	ID                 string               `json:"id"`
	VersionID          string               `json:"version_id"`
	WorkflowID         string               `json:"workflow_id"`
	State              ExecutionState       `json:"state"`
	OffloadState       string               `json:"offload_state"`
	Duration           string               `json:"duration"`
	TimeSaved          float64              `json:"time_saved"`
	ExecutedActions    int                  `json:"executed_actions"`
	HasExecutionOutput bool                 `json:"has_execution_output"`
	ScopeID            string               `json:"scope_id"`
	ScopeLevel         AutomationScopeLevel `json:"scope_level"`
	MgmtID             string               `json:"mgmt_id"`
	CreatedAt          string               `json:"created_at"`
	UpdatedAt          string               `json:"updated_at"`
	Raw                json.RawMessage      `json:"-"`
}

func (e *WorkflowExecutionRun) UnmarshalJSON(b []byte) error {
	type alias WorkflowExecutionRun
	if err := json.Unmarshal(b, (*alias)(e)); err != nil {
		return err
	}
	e.Raw = append(e.Raw[:0:0], b...)
	return nil
}

// --- Pagination ---

// AutomationPagination is the cursor-based pagination for Hyperautomation.
type AutomationPagination struct {
	NextCursor string `json:"nextCursor"`
	TotalItems int    `json:"totalItems"`
}

// --- List params ---

// AutomationListParams are the query parameters for listing workflows.
type AutomationListParams struct {
	SiteIDs      []string
	GroupIDs     []string
	AccountIDs   []string
	Integrations []string
	TriggerTypes []string
	CoreActions  []string
	States       []string
	ScopeIDs     []string
	Tags         []string
	NameContains string
	NameEq       string
	IsSnippet    *bool
	Oversight    *bool
	Limit        int
	Skip         int
	SortBy       string
	SortOrder    string
}

func (p *AutomationListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "integrations", p.Integrations)
	addCSV(v, "trigger_types", p.TriggerTypes)
	addCSV(v, "core_actions", p.CoreActions)
	addCSV(v, "states", p.States)
	addCSV(v, "scope_ids", p.ScopeIDs)
	addCSV(v, "tags", p.Tags)
	addString(v, "name__contains", p.NameContains)
	addString(v, "name__eq", p.NameEq)
	addBool(v, "is_snippet", p.IsSnippet)
	addBool(v, "oversight", p.Oversight)
	addInt(v, "limit", p.Limit)
	addInt(v, "skip", p.Skip)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// AutomationExecutionListParams are the query parameters for listing executions.
type AutomationExecutionListParams struct {
	SiteIDs      []string
	GroupIDs     []string
	AccountIDs   []string
	TriggerTypes []string
	States       []string
	ScopeIDs     []string
	Tags         []string
	WorkflowID   string
	NameContains string
	IsSnippet    *bool
	Limit        int
	Skip         int
	SortBy       string
	SortOrder    string
	CreatedAtGte string
	CreatedAtLt  string
}

func (p *AutomationExecutionListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "trigger_types", p.TriggerTypes)
	addCSV(v, "states", p.States)
	addCSV(v, "scope_ids", p.ScopeIDs)
	addCSV(v, "tags", p.Tags)
	addString(v, "workflow_id", p.WorkflowID)
	addString(v, "workflow_name__contains", p.NameContains)
	addBool(v, "is_snippet", p.IsSnippet)
	addInt(v, "limit", p.Limit)
	addInt(v, "skip", p.Skip)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addString(v, "created_at__gte", p.CreatedAtGte)
	addString(v, "created_at__lt", p.CreatedAtLt)
	return v
}

// --- Request structs ---

// AutomationRunData is the body for triggering a workflow execution.
type AutomationRunData struct {
	Payload                            *string               `json:"payload,omitempty"`
	SingularityResponseEventID         *string               `json:"singularity_response_event_id,omitempty"`
	SingularityResponseEventType       *SingularityEventType `json:"singularity_response_event_type,omitempty"`
	SingularityResponseExecutionSource *ExecutionSource      `json:"singularity_response_execution_source,omitempty"`
}

// --- Client methods ---

// AutomationList lists all workflows.
func (c *Client) AutomationList(ctx context.Context, params *AutomationListParams) ([]WorkflowListItem, *AutomationPagination, error) {
	var resp struct {
		Data       []WorkflowListItem   `json:"data"`
		Pagination AutomationPagination `json:"pagination"`
	}
	if err := c.get(ctx, automationBase+"/workflows", params.values(), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// AutomationVersions lists all versions of a workflow.
func (c *Client) AutomationVersions(ctx context.Context, workflowID string, params *AutomationListParams) ([]WorkflowVersion, error) {
	if workflowID == "" {
		return nil, fmt.Errorf("mgmt: workflowId is required")
	}
	var resp struct {
		Versions []WorkflowVersion `json:"versions"`
	}
	path := automationBase + "/workflows/versions/list/" + url.PathEscape(workflowID)
	if err := c.get(ctx, path, params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Versions, nil
}

// AutomationExport exports a specific workflow version.
func (c *Client) AutomationExport(ctx context.Context, workflowID, versionID string) (*WorkflowExport, error) {
	if workflowID == "" {
		return nil, fmt.Errorf("mgmt: workflowId is required")
	}
	if versionID == "" {
		return nil, fmt.Errorf("mgmt: versionId is required")
	}
	path := automationBase + "/workflow-import-export/export/" +
		url.PathEscape(workflowID) + "/" + url.PathEscape(versionID)
	var resp WorkflowExport
	if err := c.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AutomationImport imports a workflow from an exported definition.
func (c *Client) AutomationImport(ctx context.Context, data json.RawMessage, siteIDs []string) (*Workflow, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("mgmt: workflow data is required")
	}
	// The API expects {data: ..., filter: ...}. The data is the full
	// export body (name, description, actions). We wrap it with the
	// filter envelope.
	body := struct {
		Data   json.RawMessage `json:"data"`
		Filter *scopeFilter    `json:"filter,omitempty"`
	}{Data: data}
	if len(siteIDs) > 0 {
		body.Filter = &scopeFilter{
			Type:  "JsonPath",
			Value: fmt.Sprintf("$.siteIds anyOneOf %s", toFilterCSV(siteIDs)),
		}
	}
	var resp Workflow
	if err := c.post(ctx, automationBase+"/workflow-import-export/import", body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AutomationActivate activates a specific workflow version.
func (c *Client) AutomationActivate(ctx context.Context, workflowID, versionID string) error {
	if workflowID == "" {
		return fmt.Errorf("mgmt: workflowId is required")
	}
	if versionID == "" {
		return fmt.Errorf("mgmt: versionId is required")
	}
	path := automationBase + "/workflows/" +
		url.PathEscape(workflowID) + "/" + url.PathEscape(versionID) + "/activation"
	return c.post(ctx, path, nil, nil)
}

// AutomationDeactivate deactivates the active version of a workflow.
func (c *Client) AutomationDeactivate(ctx context.Context, workflowID string) error {
	if workflowID == "" {
		return fmt.Errorf("mgmt: workflowId is required")
	}
	path := automationBase + "/workflows/" + url.PathEscape(workflowID) + "/deactivate"
	return c.post(ctx, path, nil, nil)
}

// AutomationExecutions lists workflow executions.
func (c *Client) AutomationExecutions(ctx context.Context, params *AutomationExecutionListParams) ([]WorkflowExecution, *AutomationPagination, error) {
	var resp struct {
		Data       []WorkflowExecution  `json:"data"`
		Pagination AutomationPagination `json:"pagination"`
	}
	if err := c.get(ctx, automationBase+"/workflow-execution", params.values(), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// AutomationExecutionGet gets a workflow execution by ID.
func (c *Client) AutomationExecutionGet(ctx context.Context, executionID string) (*WorkflowExecutionDetail, error) {
	if executionID == "" {
		return nil, fmt.Errorf("mgmt: executionId is required")
	}
	path := automationBase + "/workflow-execution/" + url.PathEscape(executionID)
	var resp WorkflowExecutionDetail
	if err := c.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AutomationExecutionOutput gets the raw output of a workflow execution.
func (c *Client) AutomationExecutionOutput(ctx context.Context, executionID string) (json.RawMessage, error) {
	if executionID == "" {
		return nil, fmt.Errorf("mgmt: executionId is required")
	}
	path := automationBase + "/workflow-execution/output/" + url.PathEscape(executionID) + "/raw"
	var resp struct {
		ExecutionOutput json.RawMessage `json:"ExecutionOutput"`
	}
	if err := c.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.ExecutionOutput, nil
}

// AutomationRun triggers a manual workflow execution.
func (c *Client) AutomationRun(ctx context.Context, workflowID, versionID string, data *AutomationRunData) (*WorkflowExecutionRun, error) {
	if workflowID == "" {
		return nil, fmt.Errorf("mgmt: workflowId is required")
	}
	if versionID == "" {
		return nil, fmt.Errorf("mgmt: versionId is required")
	}
	body := struct {
		Data *AutomationRunData `json:"data,omitempty"`
	}{Data: data}
	path := automationBase + "/workflow-execution/manual/" +
		url.PathEscape(workflowID) + "/" + url.PathEscape(versionID)
	var resp WorkflowExecutionRun
	if err := c.post(ctx, path, body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// scopeFilter is the JsonPath filter used by Hyperautomation endpoints.
type scopeFilter struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func toFilterCSV(ids []string) string {
	return strings.Join(ids, ",")
}
