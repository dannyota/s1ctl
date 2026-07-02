package mgmt

import (
	"context"
	"fmt"
	"net/url"
)

// ThreatsMitigate applies a mitigation action to threats.
func (c *Client) ThreatsMitigate(ctx context.Context, action string, filter ActionFilter) (int, error) {
	if action == "" {
		return 0, fmt.Errorf("mgmt: mitigation action is required")
	}
	return doAction(c, ctx, fmt.Sprintf("/threats/mitigate/%s", url.PathEscape(action)), filter, nil)
}

// ThreatsUpdateVerdict updates the analyst verdict on threats.
func (c *Client) ThreatsUpdateVerdict(ctx context.Context, verdict string, filter ActionFilter) (int, error) {
	data := map[string]string{"analystVerdict": verdict}
	return doAction(c, ctx, "/threats/analyst-verdict", filter, data)
}

// ThreatsUpdateStatus updates the incident status on threats.
func (c *Client) ThreatsUpdateStatus(ctx context.Context, status string, filter ActionFilter) (int, error) {
	data := map[string]string{"incidentStatus": status}
	return doAction(c, ctx, "/threats/incident", filter, data)
}

// ThreatsAddToBlacklist adds threat hashes to the blacklist.
func (c *Client) ThreatsAddToBlacklist(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/threats/add-to-blacklist", filter, nil)
}

// ThreatsFetchFile fetches threat files for further analysis.
func (c *Client) ThreatsFetchFile(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/threats/fetch-file", filter, nil)
}

// ThreatExclusionScope is the scope in which an exclusion is created.
type ThreatExclusionScope string

const (
	ThreatExclusionScopeGroup   ThreatExclusionScope = "group"
	ThreatExclusionScopeSite    ThreatExclusionScope = "site"
	ThreatExclusionScopeAccount ThreatExclusionScope = "account"
	ThreatExclusionScopeTenant  ThreatExclusionScope = "tenant"
)

// ThreatExclusionType is the exclusion type created from a threat.
type ThreatExclusionType string

const (
	ThreatExclusionTypeHash        ThreatExclusionType = "hash"
	ThreatExclusionTypePath        ThreatExclusionType = "path"
	ThreatExclusionTypeCertificate ThreatExclusionType = "certificate"
	ThreatExclusionTypeBrowser     ThreatExclusionType = "browser"
	ThreatExclusionTypeFileType    ThreatExclusionType = "file_type"
)

// ThreatExclusionMode is the exclusion mode (path exclusions only).
type ThreatExclusionMode string

const (
	ThreatExclusionModeSuppress                 ThreatExclusionMode = "suppress"
	ThreatExclusionModeSuppressDynamicOnly      ThreatExclusionMode = "suppress_dynamic_only"
	ThreatExclusionModeSuppressDFIOnly          ThreatExclusionMode = "suppress_dfi_only"
	ThreatExclusionModeDisableInProcMonitor     ThreatExclusionMode = "disable_in_process_monitor"
	ThreatExclusionModeDisableInProcMonitorDeep ThreatExclusionMode = "disable_in_process_monitor_deep"
	ThreatExclusionModeDisableAllMonitors       ThreatExclusionMode = "disable_all_monitors"
	ThreatExclusionModeDisableAllMonitorsDeep   ThreatExclusionMode = "disable_all_monitors_deep"
	ThreatExclusionModeSuppressAppControl       ThreatExclusionMode = "suppress_app_control"
	ThreatExclusionModeSuppressDriftDetection   ThreatExclusionMode = "suppress_drift_detection"
)

// ThreatExclusionOptions configure how a threat is added to exclusions.
// TargetScope and Type are required; Value defaults to the relevant value
// from the threat when omitted.
type ThreatExclusionOptions struct {
	TargetScope       ThreatExclusionScope `json:"targetScope"`
	Type              ThreatExclusionType  `json:"type"`
	Value             string               `json:"value,omitempty"`
	Description       string               `json:"description,omitempty"`
	Mode              ThreatExclusionMode  `json:"mode,omitempty"`
	PathExclusionType string               `json:"pathExclusionType,omitempty"`
	Note              string               `json:"note,omitempty"`
	ExternalTicketID  string               `json:"externalTicketId,omitempty"`
	Actions           []string             `json:"actions,omitempty"`
}

// ThreatsAddToExclusions creates an exclusion from the selected threats.
func (c *Client) ThreatsAddToExclusions(ctx context.Context, filter ActionFilter, opts ThreatExclusionOptions) (int, error) {
	if opts.TargetScope == "" {
		return 0, fmt.Errorf("mgmt: exclusion target scope is required")
	}
	if opts.Type == "" {
		return 0, fmt.Errorf("mgmt: exclusion type is required")
	}
	return doAction(c, ctx, "/threats/add-to-exclusions", filter, opts)
}

// ThreatMitigationAction is a mitigation action applied to threats or alerts.
type ThreatMitigationAction string

const (
	ThreatMitigationKill                ThreatMitigationAction = "kill"
	ThreatMitigationRemediate           ThreatMitigationAction = "remediate"
	ThreatMitigationRollbackRemediation ThreatMitigationAction = "rollback-remediation"
	ThreatMitigationQuarantine          ThreatMitigationAction = "quarantine"
	ThreatMitigationUnQuarantine        ThreatMitigationAction = "un-quarantine"
	ThreatMitigationRemoveMacros        ThreatMitigationAction = "remove_macros"
	ThreatMitigationRestoreMacros       ThreatMitigationAction = "restore_macros"
)

// ThreatAlert identifies a Deep Visibility alert to mark as a threat and
// mitigate. Both AgentID and Storyline are required by the API.
type ThreatAlert struct {
	AgentID   string `json:"agentId"`
	Storyline string `json:"storyline"`
}

// ThreatsMitigateAlerts marks the given alerts as threats and runs a
// mitigation action. Unlike the other threat actions this endpoint takes an
// explicit list of alerts (agent ID + storyline) rather than a filter.
func (c *Client) ThreatsMitigateAlerts(ctx context.Context, alerts []ThreatAlert, action ThreatMitigationAction) (int, error) {
	if len(alerts) == 0 {
		return 0, fmt.Errorf("mgmt: at least one alert (agent ID + storyline) is required")
	}
	data := map[string]any{"alerts": alerts}
	if action != "" {
		data["action"] = action
	}
	var resp affectedResponse
	if err := c.post(ctx, "/threats/mitigate-alerts", map[string]any{"data": data}, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// ThreatsSetExternalTicketID sets the external ticket ID on the selected
// threats.
func (c *Client) ThreatsSetExternalTicketID(ctx context.Context, filter ActionFilter, ticketID string) (int, error) {
	return doAction(c, ctx, "/threats/external-ticket-id", filter, map[string]string{"externalTicketId": ticketID})
}
