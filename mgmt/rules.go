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
	return update[Rule](c, ctx, fmt.Sprintf("/cloud-detection/rules/%s", id), data)
}
