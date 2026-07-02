package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// PlatformRuleSeverity is the severity of a platform detection rule.
type PlatformRuleSeverity string

const (
	PlatformRuleSeverityInfo     PlatformRuleSeverity = "Info"
	PlatformRuleSeverityLow      PlatformRuleSeverity = "Low"
	PlatformRuleSeverityMedium   PlatformRuleSeverity = "Medium"
	PlatformRuleSeverityHigh     PlatformRuleSeverity = "High"
	PlatformRuleSeverityCritical PlatformRuleSeverity = "Critical"
)

// PlatformRuleStatus is the status of a platform detection rule.
type PlatformRuleStatus string

const (
	PlatformRuleStatusDraft      PlatformRuleStatus = "Draft"
	PlatformRuleStatusActivating PlatformRuleStatus = "Activating"
	PlatformRuleStatusActive     PlatformRuleStatus = "Active"
	PlatformRuleStatusDisabling  PlatformRuleStatus = "Disabling"
	PlatformRuleStatusDisabled   PlatformRuleStatus = "Disabled"
	PlatformRuleStatusDeleted    PlatformRuleStatus = "Deleted"
	PlatformRuleStatusDeleting   PlatformRuleStatus = "Deleting"
)

// PlatformRuleCategory is the category of a platform detection rule.
type PlatformRuleCategory string

const (
	PlatformRuleCategoryEvents        PlatformRuleCategory = "Events"
	PlatformRuleCategoryCorrelation   PlatformRuleCategory = "Correlation"
	PlatformRuleCategoryUEBAFirstSeen PlatformRuleCategory = "UEBAFirstSeen"
	PlatformRuleCategoryScheduled     PlatformRuleCategory = "Scheduled"
)

// PlatformRuleScope is the scope level of a platform detection rule.
type PlatformRuleScope string

const (
	PlatformRuleScopeGlobal  PlatformRuleScope = "global"
	PlatformRuleScopeAccount PlatformRuleScope = "account"
	PlatformRuleScopeSite    PlatformRuleScope = "site"
	PlatformRuleScopeGroup   PlatformRuleScope = "group"
)

// MitreTechnique is a MITRE ATT&CK technique reference.
type MitreTechnique struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`

	Raw json.RawMessage `json:"-"`
}

func (m *MitreTechnique) UnmarshalJSON(b []byte) error {
	type alias MitreTechnique
	if err := json.Unmarshal(b, (*alias)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0:0], b...)
	return nil
}

// MitreTactic is a MITRE ATT&CK tactic with associated techniques.
type MitreTactic struct {
	Tactic     string           `json:"tactic"`
	Techniques []MitreTechnique `json:"techniques"`

	Raw json.RawMessage `json:"-"`
}

func (m *MitreTactic) UnmarshalJSON(b []byte) error {
	type alias MitreTactic
	if err := json.Unmarshal(b, (*alias)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0:0], b...)
	return nil
}

// PlatformRule is a SentinelOne platform (pre-built) detection rule.
type PlatformRule struct {
	ID                    string               `json:"id"`
	Name                  string               `json:"name"`
	Description           string               `json:"description"`
	Severity              PlatformRuleSeverity `json:"severity"`
	Status                PlatformRuleStatus   `json:"status"`
	ScopeLevel            PlatformRuleScope    `json:"scopeLevel"`
	HighestInheritedScope PlatformRuleScope    `json:"highestInheritedScopeLevel"`
	QueryType             string               `json:"queryType"`
	S1QL                  string               `json:"s1ql"`
	CreatedBy             string               `json:"createdBy"`
	AttackSurfaces        []string             `json:"attackSurfaces"`
	Sources               []string             `json:"sources"`
	Tags                  []string             `json:"tags"`
	ActiveResponse        bool                 `json:"activeResponse"`
	NetworkQuarantine     bool                 `json:"networkQuarantine"`
	TreatAsThreat         RuleTreatAsThreat    `json:"treatAsThreat"`
	GeneratedAlerts       int                  `json:"generatedAlerts"`
	Mitre                 []MitreTactic        `json:"mitre"`
	CreatedAt             string               `json:"createdAt"`
	UpdatedAt             string               `json:"updatedAt"`
	LastAlertTime         string               `json:"lastAlertTime"`

	Raw json.RawMessage `json:"-"`
}

func (r *PlatformRule) UnmarshalJSON(b []byte) error {
	type alias PlatformRule
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// PlatformRuleListParams are query parameters for listing platform detection rules.
type PlatformRuleListParams struct {
	IDs            []string
	ScopeID        string
	ScopeLevel     string
	Severities     []string
	Statuses       []string
	AttackSurfaces []string
	Sources        []string
	Categories     []string
	Tags           []string
	MitreTactics   []string
	NameContains   string
	S1QLContains   string
	DescContains   string
	Limit          int
	Cursor         string
}

func (p *PlatformRuleListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "platformRuleIds", p.IDs)
	addString(v, "scopeId", p.ScopeID)
	addString(v, "scopeLevel", p.ScopeLevel)
	addCSV(v, "severities", p.Severities)
	addCSV(v, "statuses", p.Statuses)
	addCSV(v, "attackSurfaces", p.AttackSurfaces)
	addCSV(v, "sources", p.Sources)
	addCSV(v, "categories", p.Categories)
	addCSV(v, "tags", p.Tags)
	addCSV(v, "mitreTactics", p.MitreTactics)
	addString(v, "ruleNameSubstring", p.NameContains)
	addString(v, "s1ql__contains", p.S1QLContains)
	addString(v, "description__contains", p.DescContains)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// PlatformRulesList returns a paginated list of platform detection rules.
func (c *Client) PlatformRulesList(ctx context.Context, params *PlatformRuleListParams) ([]PlatformRule, *Pagination, error) {
	return list[PlatformRule](c, ctx, "/detection-library/platform-rules", params.values())
}

// PlatformRuleActionFilter selects which platform rules to enable or disable.
type PlatformRuleActionFilter struct {
	PlatformRuleIDs []string `json:"platformRuleIds,omitempty"`
	ScopeID         string   `json:"scopeId,omitempty"`
	ScopeLevel      string   `json:"scopeLevel,omitempty"`
}

// PlatformRulesEnable enables platform detection rules matching the filter.
func (c *Client) PlatformRulesEnable(ctx context.Context, filter PlatformRuleActionFilter) (int, error) {
	return platformRuleAction(c, ctx, "/detection-library/platform-rules/enable", filter)
}

// PlatformRulesDisable disables platform detection rules matching the filter.
func (c *Client) PlatformRulesDisable(ctx context.Context, filter PlatformRuleActionFilter) (int, error) {
	return platformRuleAction(c, ctx, "/detection-library/platform-rules/disable", filter)
}

func platformRuleAction(c *Client, ctx context.Context, path string, filter PlatformRuleActionFilter) (int, error) {
	// Explicit rule IDs are required by design: ScopeID/ScopeLevel only narrow
	// the ID list, and allowing an ID-less filter would toggle every rule in
	// the scope in one call.
	if len(filter.PlatformRuleIDs) == 0 {
		return 0, fmt.Errorf("mgmt: platform rule action requires at least one rule ID")
	}
	var resp affectedResponse
	if err := c.put(ctx, path, filter, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// DetectionKeyValue is a key-value pair returned by detection library lookups.
type DetectionKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`

	Raw json.RawMessage `json:"-"`
}

func (d *DetectionKeyValue) UnmarshalJSON(b []byte) error {
	type alias DetectionKeyValue
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// DetectionSurfacesList returns the available detection surfaces.
func (c *Client) DetectionSurfacesList(ctx context.Context) ([]DetectionKeyValue, error) {
	var resp struct {
		Data struct {
			Surfaces []DetectionKeyValue `json:"surfaces"`
		} `json:"data"`
	}
	if err := c.get(ctx, "/detection-library/surfaces", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data.Surfaces, nil
}

// DetectionDataSourcesList returns the available detection data sources.
func (c *Client) DetectionDataSourcesList(ctx context.Context) ([]DetectionKeyValue, error) {
	var resp struct {
		Data struct {
			DataSources []DetectionKeyValue `json:"dataSources"`
		} `json:"data"`
	}
	if err := c.get(ctx, "/detection-library/data-sources", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data.DataSources, nil
}
