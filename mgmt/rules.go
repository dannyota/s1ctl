package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// RuleSeverity is the severity level of a custom detection rule.
type RuleSeverity string

const (
	RuleSeverityInfo     RuleSeverity = "Info"
	RuleSeverityLow      RuleSeverity = "Low"
	RuleSeverityMedium   RuleSeverity = "Medium"
	RuleSeverityHigh     RuleSeverity = "High"
	RuleSeverityCritical RuleSeverity = "Critical"
)

// RuleStatus is the status of a custom detection rule.
type RuleStatus string

const (
	RuleStatusDraft      RuleStatus = "Draft"
	RuleStatusActivating RuleStatus = "Activating"
	RuleStatusActive     RuleStatus = "Active"
	RuleStatusDisabling  RuleStatus = "Disabling"
	RuleStatusDisabled   RuleStatus = "Disabled"
	RuleStatusDeleted    RuleStatus = "Deleted"
	RuleStatusDeleting   RuleStatus = "Deleting"
)

// RuleQueryType is the query type of a custom detection rule.
type RuleQueryType string

const (
	RuleQueryTypeEvents        RuleQueryType = "events"
	RuleQueryTypeCorrelation   RuleQueryType = "correlation"
	RuleQueryTypeUEBAFirstSeen RuleQueryType = "uebafirstseen"
	RuleQueryTypeScheduled     RuleQueryType = "scheduled"
)

// RuleScope is the scope level of a custom detection rule.
type RuleScope string

const (
	RuleScopeGlobal  RuleScope = "global"
	RuleScopeAccount RuleScope = "account"
	RuleScopeSite    RuleScope = "site"
	RuleScopeGroup   RuleScope = "group"
)

// RuleExpirationMode indicates whether a rule is permanent or temporary.
type RuleExpirationMode string

const (
	RuleExpirationPermanent RuleExpirationMode = "Permanent"
	RuleExpirationTemporary RuleExpirationMode = "Temporary"
)

// RuleTreatAsThreat is the auto-response threat classification.
type RuleTreatAsThreat string

const (
	RuleTreatUndefined  RuleTreatAsThreat = "UNDEFINED"
	RuleTreatSuspicious RuleTreatAsThreat = "Suspicious"
	RuleTreatMalicious  RuleTreatAsThreat = "Malicious"
)

// Rule is a SentinelOne custom detection rule (STAR).
type Rule struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Status            RuleStatus         `json:"status"`
	StatusReason      string             `json:"statusReason"`
	Severity          RuleSeverity       `json:"severity"`
	S1QL              string             `json:"s1ql"`
	QueryType         RuleQueryType      `json:"queryType"`
	QueryLang         string             `json:"queryLang"`
	Scope             RuleScope          `json:"scope"`
	ScopeID           []string           `json:"scopeId"`
	ExpirationMode    RuleExpirationMode `json:"expirationMode"`
	Expiration        string             `json:"expiration"`
	Expired           bool               `json:"expired"`
	TreatAsThreat     RuleTreatAsThreat  `json:"treatAsThreat"`
	ActiveResponse    bool               `json:"activeResponse"`
	NetworkQuarantine bool               `json:"networkQuarantine"`
	GeneratedAlerts   int                `json:"generatedAlerts"`
	ReachedLimit      bool               `json:"reachedLimit"`
	Creator           string             `json:"creator"`
	CreatorID         string             `json:"creatorId"`
	AccountID         string             `json:"accountId"`
	AccountName       string             `json:"accountName"`
	SiteID            string             `json:"siteId"`
	SiteName          string             `json:"siteName"`
	CreatedAt         string             `json:"createdAt"`
	UpdatedAt         string             `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (r *Rule) UnmarshalJSON(b []byte) error {
	type alias Rule
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// RuleListParams are query parameters for listing custom detection rules.
type RuleListParams struct {
	SiteIDs      []string
	AccountIDs   []string
	GroupIDs     []string
	IDs          []string
	Status       []string
	Severity     []string
	Scopes       []string
	QueryType    []string
	NameContains string
	Query        string
	Limit        int
	Cursor       string
	SortBy       string
	SortOrder    string
}

func (p *RuleListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "ids", p.IDs)
	addCSV(v, "status", p.Status)
	addCSV(v, "severity", p.Severity)
	addCSV(v, "scopes", p.Scopes)
	addCSV(v, "queryType", p.QueryType)
	addString(v, "name__contains", p.NameContains)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// RulesList returns a paginated list of custom detection rules.
func (c *Client) RulesList(ctx context.Context, params *RuleListParams) ([]Rule, *Pagination, error) {
	return list[Rule](c, ctx, "/cloud-detection/rules", params.values())
}

// RulesGet returns a single custom detection rule by ID.
func (c *Client) RulesGet(ctx context.Context, id string) (*Rule, error) {
	return getByID[Rule](c, ctx, "/cloud-detection/rules", "rule", id)
}

// RuleCreate is the request body for creating or updating a custom detection rule.
type RuleCreate struct {
	Name              string             `json:"name"`
	Description       string             `json:"description,omitempty"`
	S1QL              string             `json:"s1ql"`
	Severity          RuleSeverity       `json:"severity"`
	Status            RuleStatus         `json:"status"`
	QueryType         RuleQueryType      `json:"queryType"`
	QueryLang         string             `json:"queryLang,omitempty"`
	ExpirationMode    RuleExpirationMode `json:"expirationMode"`
	Expiration        string             `json:"expiration,omitempty"`
	TreatAsThreat     RuleTreatAsThreat  `json:"treatAsThreat"`
	NetworkQuarantine bool               `json:"networkQuarantine"`
}

// RulesCreate creates a custom detection rule.
func (c *Client) RulesCreate(ctx context.Context, data RuleCreate) (*Rule, error) {
	return create[Rule](c, ctx, "/cloud-detection/rules", data)
}

// RulesUpdate updates a custom detection rule by ID.
func (c *Client) RulesUpdate(ctx context.Context, id string, data RuleCreate) (*Rule, error) {
	return update[Rule](c, ctx, fmt.Sprintf("/cloud-detection/rules/%s", url.PathEscape(id)), data)
}

// RuleActionFilter selects which rules to enable or disable.
type RuleActionFilter struct {
	IDs        []string `json:"ids,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
}

// RulesEnable activates custom detection rules matching the filter.
func (c *Client) RulesEnable(ctx context.Context, filter RuleActionFilter) (int, error) {
	return rulesAction(c, ctx, "/cloud-detection/rules/enable", filter)
}

// RulesDisable deactivates custom detection rules matching the filter.
func (c *Client) RulesDisable(ctx context.Context, filter RuleActionFilter) (int, error) {
	return rulesAction(c, ctx, "/cloud-detection/rules/disable", filter)
}

func rulesAction(c *Client, ctx context.Context, path string, filter RuleActionFilter) (int, error) {
	if len(filter.IDs) == 0 && len(filter.SiteIDs) == 0 && len(filter.AccountIDs) == 0 {
		return 0, fmt.Errorf("mgmt: rule action requires at least one filter (ids, siteIds, or accountIds)")
	}
	req := struct {
		Filter RuleActionFilter `json:"filter"`
	}{Filter: filter}
	var resp affectedResponse
	if err := c.put(ctx, path, req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// CloudDetectionAlert is a STAR custom detection alert from the REST API.
type CloudDetectionAlert struct {
	AlertInfo          CDAlertInfo `json:"alertInfo"`
	RuleInfo           CDRuleInfo  `json:"ruleInfo"`
	AgentDetectionInfo CDAgentInfo `json:"agentDetectionInfo"`

	Raw json.RawMessage `json:"-"`
}

func (a *CloudDetectionAlert) UnmarshalJSON(b []byte) error {
	type alias CloudDetectionAlert
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// CDAlertInfo contains the alert-level fields of a cloud detection alert.
type CDAlertInfo struct {
	AlertID        string `json:"alertId"`
	IncidentStatus string `json:"incidentStatus"`
	AnalystVerdict string `json:"analystVerdict"`
	Severity       string `json:"severity"`
	EventType      string `json:"eventType"`
	HitType        string `json:"hitType"`
	Source         string `json:"source"`
	CreatedAt      string `json:"createdAt"`
	ReportedAt     string `json:"reportedAt"`
	UpdatedAt      string `json:"updatedAt"`
	DstIP          string `json:"dstIp"`
	DstPort        string `json:"dstPort"`
	SrcIP          string `json:"srcIp"`
	SrcPort        string `json:"srcPort"`

	Raw json.RawMessage `json:"-"`
}

func (c *CDAlertInfo) UnmarshalJSON(b []byte) error {
	type alias CDAlertInfo
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// CDRuleInfo contains the rule-level fields of a cloud detection alert.
type CDRuleInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Severity      string `json:"severity"`
	QueryType     string `json:"queryType"`
	ScopeLevel    string `json:"scopeLevel"`
	TreatAsThreat string `json:"treatAsThreat"`

	Raw json.RawMessage `json:"-"`
}

func (c *CDRuleInfo) UnmarshalJSON(b []byte) error {
	type alias CDRuleInfo
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// CDAgentInfo contains the agent-level fields of a cloud detection alert.
type CDAgentInfo struct {
	AccountID   string `json:"accountId"`
	SiteID      string `json:"siteId"`
	Name        string `json:"name"`
	MachineType string `json:"machineType"`
	OSFamily    string `json:"osFamily"`
	OSName      string `json:"osName"`
	Version     string `json:"version"`
	UUID        string `json:"uuid"`

	Raw json.RawMessage `json:"-"`
}

func (c *CDAgentInfo) UnmarshalJSON(b []byte) error {
	type alias CDAgentInfo
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// CDAlertListParams are query parameters for listing cloud detection alerts.
type CDAlertListParams struct {
	SiteIDs          []string
	AccountIDs       []string
	RuleNameContains []string
	Severity         []string
	IncidentStatus   []string
	AnalystVerdict   []string
	CreatedAtGt      string
	CreatedAtLt      string
	ReportedAtGt     string
	ReportedAtLt     string
	Query            string
	Limit            int
	Cursor           string
	SortBy           string
	SortOrder        string
}

func (p *CDAlertListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "ruleName__contains", p.RuleNameContains)
	addCSV(v, "severity", p.Severity)
	addCSV(v, "incidentStatus", p.IncidentStatus)
	addCSV(v, "analystVerdict", p.AnalystVerdict)
	addString(v, "createdAt__gt", p.CreatedAtGt)
	addString(v, "createdAt__lt", p.CreatedAtLt)
	addString(v, "reportedAt__gt", p.ReportedAtGt)
	addString(v, "reportedAt__lt", p.ReportedAtLt)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// CloudDetectionAlertsList returns a paginated list of STAR cloud detection alerts.
func (c *Client) CloudDetectionAlertsList(ctx context.Context, params *CDAlertListParams) ([]CloudDetectionAlert, *Pagination, error) {
	return list[CloudDetectionAlert](c, ctx, "/cloud-detection/alerts", params.values())
}
