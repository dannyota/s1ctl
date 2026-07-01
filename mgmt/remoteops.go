package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// RemoteScript is a SentinelOne remote ops script.
type RemoteScript struct {
	ID          string   `json:"id"`
	FileName    string   `json:"fileName"`
	FileType    string   `json:"fileType"`
	ScriptType  string   `json:"scriptType"`
	OSTypes     []string `json:"osTypes"`
	ScopeID     string   `json:"scopeId"`
	ScopeLevel  string   `json:"scopeLevel"`
	CreatedAt   string   `json:"createdAt"`
	CreatorID   string   `json:"creatorId"`
	CreatorName string   `json:"creatorName"`

	Raw json.RawMessage `json:"-"`
}

func (r *RemoteScript) UnmarshalJSON(b []byte) error {
	type alias RemoteScript
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// RemoteScriptListParams are query parameters for listing remote scripts.
type RemoteScriptListParams struct {
	SiteIDs    []string
	AccountIDs []string
	OSTypes    []string
	Query      string
	Limit      int
	Cursor     string
}

func (p *RemoteScriptListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "osTypes", p.OSTypes)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// RemoteScriptsList returns a paginated list of remote scripts.
func (c *Client) RemoteScriptsList(ctx context.Context, params *RemoteScriptListParams) ([]RemoteScript, *Pagination, error) {
	return list[RemoteScript](c, ctx, "/remote-scripts", params.values())
}

// RemoteScriptsGet returns a single remote script by ID.
func (c *Client) RemoteScriptsGet(ctx context.Context, id string) (*RemoteScript, error) {
	return getByID[RemoteScript](c, ctx, "/remote-scripts", "remote script", id)
}

// OutputDestination controls where remote script output is sent.
type OutputDestination string

const (
	OutputSentinelCloud  OutputDestination = "SentinelCloud"
	OutputLocal          OutputDestination = "Local"
	OutputNone           OutputDestination = "None"
	OutputSingularityXDR OutputDestination = "SingularityXDR"
)

// RemoteScriptsExecuteParams holds parameters for executing a remote script.
type RemoteScriptsExecuteParams struct {
	ScriptID          string            `json:"scriptId"`
	OutputDestination OutputDestination `json:"outputDestination"`
	TaskDescription   string            `json:"taskDescription"`
	InputParams       string            `json:"inputParams,omitempty"`
	TimeoutSeconds    int               `json:"scriptRuntimeTimeoutSeconds,omitempty"`
}

// RemoteScriptsExecuteFilter identifies which agents to execute the script on.
type RemoteScriptsExecuteFilter struct {
	IDs        []string `json:"ids,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
}

// RemoteScriptsExecuteResult is the response from executing a remote script.
type RemoteScriptsExecuteResult struct {
	Affected           int    `json:"affected"`
	ParentTaskID       string `json:"parentTaskId"`
	Pending            bool   `json:"pending"`
	PendingExecutionID string `json:"pendingExecutionId"`
}

// RemoteScriptsExecute runs a remote script on the specified agents.
func (c *Client) RemoteScriptsExecute(ctx context.Context, filter RemoteScriptsExecuteFilter, data RemoteScriptsExecuteParams) (*RemoteScriptsExecuteResult, error) {
	if data.ScriptID == "" {
		return nil, fmt.Errorf("mgmt: scriptId is required")
	}
	if data.TaskDescription == "" {
		return nil, fmt.Errorf("mgmt: taskDescription is required")
	}
	if len(filter.IDs) == 0 && len(filter.SiteIDs) == 0 && len(filter.AccountIDs) == 0 && len(filter.GroupIDs) == 0 {
		return nil, fmt.Errorf("mgmt: execute requires at least one filter (ids, siteIds, accountIds, or groupIds)")
	}
	req := struct {
		Filter RemoteScriptsExecuteFilter `json:"filter"`
		Data   RemoteScriptsExecuteParams `json:"data"`
	}{Filter: filter, Data: data}
	var resp singleResponse[RemoteScriptsExecuteResult]
	if err := c.post(ctx, "/remote-scripts/execute", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RemoteScriptTask is the status of a single agent's remote script execution.
type RemoteScriptTask struct {
	ID                string `json:"id"`
	ParentTaskID      string `json:"parentTaskId"`
	Type              string `json:"type"`
	Description       string `json:"description"`
	Status            string `json:"status"`
	DetailedStatus    string `json:"detailedStatus"`
	AgentComputerName string `json:"agentComputerName"`
	AgentOSType       string `json:"agentOsType"`
	InitiatedBy       string `json:"initiatedBy"`
	AccountName       string `json:"accountName"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (t *RemoteScriptTask) UnmarshalJSON(b []byte) error {
	type alias RemoteScriptTask
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

// RemoteScriptsStatusParams are query parameters for getting remote script task status.
type RemoteScriptsStatusParams struct {
	ParentTaskID string
	Status       []string
	Limit        int
	Cursor       string
}

func (p *RemoteScriptsStatusParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "parentTaskId", p.ParentTaskID)
	addCSV(v, "status", p.Status)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// RemoteScriptsStatus returns the status of remote script execution tasks.
func (c *Client) RemoteScriptsStatus(ctx context.Context, params *RemoteScriptsStatusParams) ([]RemoteScriptTask, *Pagination, error) {
	if params == nil || params.ParentTaskID == "" {
		return nil, nil, fmt.Errorf("mgmt: parentTaskId is required")
	}
	return list[RemoteScriptTask](c, ctx, "/remote-scripts/status", params.values())
}
