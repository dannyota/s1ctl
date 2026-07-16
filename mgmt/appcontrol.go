package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// AppControlBehavior is the enforcement behavior of an application control rule.
type AppControlBehavior string

const (
	AppControlBehaviorAllow   AppControlBehavior = "ALLOW"
	AppControlBehaviorMonitor AppControlBehavior = "MONITOR"
	AppControlBehaviorBlock   AppControlBehavior = "BLOCK"
)

// AppControlOSType is an OS type for application control rules.
type AppControlOSType string

const (
	AppControlOSMacOS   AppControlOSType = "MACOS"
	AppControlOSWindows AppControlOSType = "WINDOWS"
)

// AppControlScopeLevel is the scope level for application control resources.
type AppControlScopeLevel string

const (
	AppControlScopeAccount AppControlScopeLevel = "ACCOUNT"
	AppControlScopeSite    AppControlScopeLevel = "SITE"
	AppControlScopeGroup   AppControlScopeLevel = "GROUP"
)

// nacBase is the path prefix for NAC endpoints relative to the API v2.1 base.
const nacBase = "/nac/config/api/v1/nac"

// AppControlScope identifies the scope for application control operations.
type AppControlScope struct {
	ScopeType AppControlScopeLevel `json:"scopeType"`
	ScopeIDs  []string             `json:"scopeIds"`
}

// AppControlScopeInfo is server-returned scope metadata on a rule.
type AppControlScopeInfo struct {
	ScopeID    string `json:"scopeId"`
	ScopeLevel string `json:"scopeLevel"`
	ScopeName  string `json:"scopeName"`
	ScopePath  string `json:"scopePath"`

	Raw json.RawMessage `json:"-"`
}

func (s *AppControlScopeInfo) UnmarshalJSON(b []byte) error {
	type alias AppControlScopeInfo
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// AppControlConditions describes rule match conditions.
type AppControlConditions struct {
	Publisher          string `json:"publisher,omitempty"`
	Path               string `json:"path,omitempty"`
	Signer             string `json:"signer,omitempty"`
	SHA256             string `json:"sha256,omitempty"`
	Process            string `json:"process,omitempty"`
	ParentProcess      string `json:"parentProcess,omitempty"`
	ApplicationVersion string `json:"applicationVersion,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (c *AppControlConditions) UnmarshalJSON(b []byte) error {
	type alias AppControlConditions
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// AppControlRule is a SentinelOne application control (NAC) rule.
type AppControlRule struct {
	ID          string                 `json:"id"`
	RuleName    string                 `json:"ruleName"`
	Description string                 `json:"description"`
	Scope       *AppControlScopeInfo   `json:"scope,omitempty"`
	OSType      []AppControlOSType     `json:"osType"`
	Parameters  *AppControlConditions  `json:"parameters,omitempty"`
	Behavior    AppControlBehavior     `json:"behavior"`
	Propagation bool                   `json:"propagation"`
	Exceptions  []AppControlConditions `json:"exceptions"`
	CreatedAt   string                 `json:"createdAt"`
	CreatedBy   string                 `json:"createdBy"`

	Raw json.RawMessage `json:"-"`
}

func (r *AppControlRule) UnmarshalJSON(b []byte) error {
	type alias AppControlRule
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// AppControlRuleInput is the request body for creating or updating an
// application control rule.
type AppControlRuleInput struct {
	ID          string                 `json:"id,omitempty"`
	RuleName    string                 `json:"ruleName"`
	Description string                 `json:"description,omitempty"`
	Scope       *AppControlScope       `json:"scope,omitempty"`
	OSType      []AppControlOSType     `json:"osType,omitempty"`
	Propagation *bool                  `json:"propagation,omitempty"`
	Parameters  *AppControlConditions  `json:"parameters,omitempty"`
	Exceptions  []AppControlConditions `json:"exceptions,omitempty"`
	Behavior    AppControlBehavior     `json:"behavior,omitempty"`
}

// AppControlQueryParams are parameters for querying application control rules.
type AppControlQueryParams struct {
	ScopeType      AppControlScopeLevel
	ScopeIDs       []string
	IncludeParents bool
	PageSize       int
	Cursor         string
}

// appControlQueryRequest is the POST body for the rules query endpoint.
type appControlQueryRequest struct {
	ScopeSelector  *AppControlScope `json:"scopeSelector,omitempty"`
	IncludeParents bool             `json:"includeParents"`
}

// appControlRuleEdge is one item in the relay-style paginated response.
type appControlRuleEdge struct {
	Node   AppControlRule `json:"node"`
	Cursor string         `json:"cursor"`
}

// appControlPageInfo carries relay-style pagination cursors.
type appControlPageInfo struct {
	EndCursor   string `json:"endCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}

// appControlConnection is the relay-style paginated response envelope.
type appControlConnection struct {
	PageInfo   appControlPageInfo   `json:"pageInfo"`
	Edges      []appControlRuleEdge `json:"edges"`
	TotalCount int                  `json:"totalCount"`
}

// AppControlRulesList queries application control rules using relay-style
// cursor pagination.
func (c *Client) AppControlRulesList(ctx context.Context, params *AppControlQueryParams) ([]AppControlRule, string, int, error) {
	q := url.Values{}
	if params != nil {
		addInt(q, "pageSize", params.PageSize)
		addString(q, "cursor", params.Cursor)
	}

	body := appControlQueryRequest{}
	if params != nil {
		body.IncludeParents = params.IncludeParents
		if len(params.ScopeIDs) > 0 && params.ScopeType != "" {
			body.ScopeSelector = &AppControlScope{
				ScopeType: params.ScopeType,
				ScopeIDs:  params.ScopeIDs,
			}
		}
	}

	path := nacBase + "/rules/query"
	u := path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	var conn appControlConnection
	if err := c.jsonRequest(ctx, "POST", u, body, &conn); err != nil {
		return nil, "", 0, err
	}

	rules := make([]AppControlRule, len(conn.Edges))
	for i, e := range conn.Edges {
		rules[i] = e.Node
	}
	return rules, conn.PageInfo.EndCursor, conn.TotalCount, nil
}

// AppControlRulesGet returns a single application control rule by ID.
func (c *Client) AppControlRulesGet(ctx context.Context, id string) (*AppControlRule, error) {
	path := nacBase + "/rules/" + url.PathEscape(id)
	var rule AppControlRule
	if err := c.get(ctx, path, nil, &rule); err != nil {
		return nil, err
	}
	if rule.ID == "" {
		return nil, fmt.Errorf("mgmt: application control rule %s not found", id)
	}
	return &rule, nil
}

// AppControlCommonResponse is the response from create/update/delete operations.
type AppControlCommonResponse struct {
	Success          bool   `json:"success"`
	ID               string `json:"id"`
	StatusCode       int    `json:"statusCode"`
	StatusMessage    string `json:"statusMessage"`
	ValidationErrors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Value   string `json:"value"`
	} `json:"validationErrors"`

	Raw json.RawMessage `json:"-"`
}

func (r *AppControlCommonResponse) UnmarshalJSON(b []byte) error {
	type alias AppControlCommonResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// AppControlRulesCreate creates an application control rule.
func (c *Client) AppControlRulesCreate(ctx context.Context, input AppControlRuleInput) (*AppControlCommonResponse, error) {
	path := nacBase + "/rules"
	var resp AppControlCommonResponse
	if err := c.post(ctx, path, input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppControlRulesUpdate updates an application control rule by ID.
func (c *Client) AppControlRulesUpdate(ctx context.Context, id string, input AppControlRuleInput) (*AppControlCommonResponse, error) {
	path := nacBase + "/rules/" + url.PathEscape(id)
	var resp AppControlCommonResponse
	if err := c.put(ctx, path, input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppControlRulesDelete deletes application control rules by ID.
func (c *Client) AppControlRulesDelete(ctx context.Context, ids []string, scope *AppControlScope) (*AppControlCommonResponse, error) {
	q := url.Values{}
	for _, id := range ids {
		q.Add("ids", id)
	}
	path := nacBase + "/rules"
	u := path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	var resp AppControlCommonResponse
	if err := c.jsonRequest(ctx, "DELETE", u, scope, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppControlLabel is a label in the application control system.
type AppControlLabel struct {
	ID        string `json:"id"`
	LabelName string `json:"labelName"`

	Raw json.RawMessage `json:"-"`
}

func (l *AppControlLabel) UnmarshalJSON(b []byte) error {
	type alias AppControlLabel
	if err := json.Unmarshal(b, (*alias)(l)); err != nil {
		return err
	}
	l.Raw = append(l.Raw[:0:0], b...)
	return nil
}

// AppControlLabelsList returns all application control labels.
func (c *Client) AppControlLabelsList(ctx context.Context) ([]AppControlLabel, error) {
	path := nacBase + "/labels"
	var labels []AppControlLabel
	if err := c.get(ctx, path, nil, &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

// AppControlSettings is the application control (NAC) settings.
type AppControlSettings struct {
	FallbackBehavior          AppControlBehavior `json:"fallbackBehavior"`
	EnableApplicationControl  bool               `json:"enableApplicationControl"`
	InheritApplicationControl bool               `json:"inheritApplicationControl"`

	Raw json.RawMessage `json:"-"`
}

func (s *AppControlSettings) UnmarshalJSON(b []byte) error {
	type alias AppControlSettings
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// AppControlSettingsInput is the request body for updating NAC settings.
type AppControlSettingsInput struct {
	Scope                     *AppControlScope   `json:"scope,omitempty"`
	FallbackBehavior          AppControlBehavior `json:"fallbackBehavior,omitempty"`
	EnableApplicationControl  *bool              `json:"enableApplicationControl,omitempty"`
	InheritApplicationControl *bool              `json:"inheritApplicationControl,omitempty"`
}

// AppControlSettingsGet returns application control (NAC) settings.
func (c *Client) AppControlSettingsGet(ctx context.Context) (*AppControlSettings, error) {
	path := nacBase + "/settings"
	var settings AppControlSettings
	if err := c.get(ctx, path, nil, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// AppControlSettingsUpdate updates application control (NAC) settings.
func (c *Client) AppControlSettingsUpdate(ctx context.Context, input AppControlSettingsInput) (*AppControlCommonResponse, error) {
	path := nacBase + "/settings"
	var resp AppControlCommonResponse
	if err := c.put(ctx, path, input, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AppMgmtScanSchedule is the scan schedule in application management settings.
type AppMgmtScanSchedule struct {
	ScanEvery int    `json:"scanEvery"`
	RepeatOn  string `json:"repeatOn"`
	Timezone  string `json:"timezone"`
	Time      string `json:"time"`

	Raw json.RawMessage `json:"-"`
}

func (s *AppMgmtScanSchedule) UnmarshalJSON(b []byte) error {
	type alias AppMgmtScanSchedule
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// AppMgmtSettings is the application management settings.
type AppMgmtSettings struct {
	ExtensiveScanEnabled   bool                 `json:"extensiveScanEnabled"`
	IsDefaultPolicy        bool                 `json:"isDefaultPolicy"`
	ScanSchedule           *AppMgmtScanSchedule `json:"scanSchedule,omitempty"`
	HasBreakingInheritance bool                 `json:"hasBreakingInheritance"`

	Raw json.RawMessage `json:"-"`
}

func (s *AppMgmtSettings) UnmarshalJSON(b []byte) error {
	type alias AppMgmtSettings
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// AppMgmtSettingsListParams are query parameters for getting application
// management settings.
type AppMgmtSettingsListParams struct {
	SiteIDs    []string
	GroupIDs   []string
	AccountIDs []string
}

func (p *AppMgmtSettingsListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	return v
}

// AppMgmtSettingsGet returns application management settings.
func (c *Client) AppMgmtSettingsGet(ctx context.Context, params *AppMgmtSettingsListParams) (*AppMgmtSettings, error) {
	var resp singleResponse[AppMgmtSettings]
	if err := c.get(ctx, "/application-management/settings", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AppMgmtSettingsScope identifies the scope for application management settings updates.
type AppMgmtSettingsScope struct {
	Tenant     bool     `json:"tenant,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
}

// AppMgmtSettingsUpdateData is the data payload for updating application
// management settings.
type AppMgmtSettingsUpdateData struct {
	ExtensiveScanEnabled *bool                `json:"extensiveScanEnabled,omitempty"`
	IsDefaultPolicy      *bool                `json:"isDefaultPolicy,omitempty"`
	ScanSchedule         *AppMgmtScanSchedule `json:"scanSchedule,omitempty"`
}

// AppMgmtSettingsUpdate updates application management settings.
func (c *Client) AppMgmtSettingsUpdate(ctx context.Context, scope AppMgmtSettingsScope, data AppMgmtSettingsUpdateData) (*AppMgmtSettings, error) {
	req := struct {
		Filter AppMgmtSettingsScope      `json:"filter"`
		Data   AppMgmtSettingsUpdateData `json:"data"`
	}{scope, data}
	var resp singleResponse[AppMgmtSettings]
	if err := c.post(ctx, "/application-management/settings", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
