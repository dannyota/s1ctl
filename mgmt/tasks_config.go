package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// TaskType identifies the task category a maintenance-window / concurrency
// configuration applies to.
type TaskType string

// Task types for maintenance-window configuration.
const (
	TaskTypeAgentsUpgrade       TaskType = "agents_upgrade"
	TaskTypeAgentVersionChange  TaskType = "agent_version_change"
	TaskTypeAutoDeploy          TaskType = "auto_deploy"
	TaskTypeScriptExecution     TaskType = "script_execution"
	TaskTypeCISScan             TaskType = "cis_scan"
	TaskTypeGAD                 TaskType = "gad"
	TaskTypeForensicsCollection TaskType = "forensics_collection"
)

// TasksConfig is the task configuration of a scope: concurrency limits and the
// maintenance windows during which the task type may run.
//
// maintenanceWindowsByDay and policyPayload are nested structures captured
// verbatim as raw blobs. maintenanceWindowsByDay is the classic per-day window
// map; policyPayload carries the flexible maintenance-window format and is
// populated only by the flexible endpoints.
type TasksConfig struct {
	InheritParentConcurrencyConfig bool     `json:"inheritParentConcurrencyConfig"`
	InheritParentMaintenanceConfig bool     `json:"inheritParentMaintenanceConfig"`
	MaxConcurrent                  int      `json:"maxConcurrent"`
	ParentMaxConcurrent            int      `json:"parentMaxConcurrent"`
	TimezoneGMT                    string   `json:"timezoneGmt"`
	TaskType                       TaskType `json:"taskType"`
	MaintenanceConfigUpdatedAt     string   `json:"maintenanceConfigUpdatedAt"`
	MaintenanceConfigUpdatedBy     string   `json:"maintenanceConfigUpdatedBy"`
	ConcurrencyConfigUpdatedAt     string   `json:"concurrencyConfigUpdatedAt"`
	ConcurrencyConfigUpdatedBy     string   `json:"concurrencyConfigUpdatedBy"`

	MaintenanceWindowsByDay json.RawMessage `json:"maintenanceWindowsByDay,omitempty"`
	PolicyPayload           json.RawMessage `json:"policyPayload,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (t *TasksConfig) UnmarshalJSON(b []byte) error {
	type alias TasksConfig
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

// TasksConfigParams selects the scope and task type of a task configuration.
// TaskType is required by the API.
type TasksConfigParams struct {
	TaskType   TaskType
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
	Tenant     bool
}

func (p *TasksConfigParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "taskType", string(p.TaskType))
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	if p.Tenant {
		v.Set("tenant", "true")
	}
	return v
}

// TasksConfigData is the declarative concurrency + maintenance-window payload of
// a task configuration write.
type TasksConfigData struct {
	InheritParentConcurrencyConfig bool            `json:"inheritParentConcurrencyConfig"`
	InheritParentMaintenanceConfig bool            `json:"inheritParentMaintenanceConfig"`
	MaxConcurrent                  int             `json:"maxConcurrent"`
	TimezoneGMT                    string          `json:"timezoneGmt"`
	MaintenanceWindowsByDay        json.RawMessage `json:"maintenanceWindowsByDay,omitempty"`
}

// TasksConfigFilter targets the scope and task type a configuration write
// applies to. TaskType is required.
type TasksConfigFilter struct {
	TaskType   TaskType `json:"taskType"`
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	Tenant     bool     `json:"tenant,omitempty"`
}

// TasksConfigWrite is the request body for updating a task configuration.
type TasksConfigWrite struct {
	Data   TasksConfigData   `json:"data"`
	Filter TasksConfigFilter `json:"filter"`
}

type tasksConfigResponse struct {
	Data TasksConfig `json:"data"`
}

// TasksConfigGet returns the task configuration for a scope and task type.
func (c *Client) TasksConfigGet(ctx context.Context, params *TasksConfigParams) (*TasksConfig, error) {
	var resp tasksConfigResponse
	if err := c.get(ctx, "/tasks-configuration", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// TasksConfigUpdate creates or updates a task configuration.
func (c *Client) TasksConfigUpdate(ctx context.Context, body TasksConfigWrite) (*TasksConfig, error) {
	var resp tasksConfigResponse
	if err := c.put(ctx, "/tasks-configuration", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// TasksConfigFlexibleGet returns the task configuration in the flexible
// maintenance-window format (policyPayload). It requires the flexible
// maintenance-window SKU.
func (c *Client) TasksConfigFlexibleGet(ctx context.Context, params *TasksConfigParams) (*TasksConfig, error) {
	var resp tasksConfigResponse
	if err := c.get(ctx, "/tasks-configuration/flexible", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// TasksConfigFlexibleUpdate updates a task configuration in the flexible
// maintenance-window format. The body ({data, filter} with a policy_payload) is
// passed through verbatim because the flexible payload shape is SKU-gated and
// open-ended.
func (c *Client) TasksConfigFlexibleUpdate(ctx context.Context, body json.RawMessage) (*TasksConfig, error) {
	var resp tasksConfigResponse
	if err := c.put(ctx, "/tasks-configuration/flexible", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// MaintenanceWindowsExport exports all maintenance-window occurrences for a
// scope as CSV. Only the flexible (policy_payload) format is supported.
func (c *Client) MaintenanceWindowsExport(ctx context.Context, params *TasksConfigParams) ([]byte, error) {
	return c.getRaw(ctx, "/tasks-configuration/maintenance-windows/export", params.values())
}
