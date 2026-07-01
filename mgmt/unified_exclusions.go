package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// UnifiedExclusion is a SentinelOne unified exclusion entry.
type UnifiedExclusion struct {
	ID                string          `json:"id"`
	ExclusionName     string          `json:"exclusionName"`
	OSType            string          `json:"osType"`
	ThreatType        string          `json:"threatType"`
	ModeType          string          `json:"modeType"`
	InteractionLevel  string          `json:"interactionLevel"`
	Description       string          `json:"description"`
	Reason            string          `json:"reason"`
	Source            string          `json:"source"`
	Type              string          `json:"type"`
	Value             any             `json:"value"`
	PathExclusionType string          `json:"pathExclusionType"`
	Engines           string          `json:"engines"`
	ChildProcess      bool            `json:"childProcess"`
	Recommendation    string          `json:"recommendation"`
	Scope             string          `json:"scope"`
	ScopePath         string          `json:"scopePath"`
	UserName          string          `json:"userName"`
	CreatorName       string          `json:"creatorName"`
	NotRecommended    string          `json:"notRecommended"`
	InAppInventory    bool            `json:"inAppInventory"`
	Imported          bool            `json:"imported"`
	LastHit           string          `json:"lastHit"`
	Hits30d           int             `json:"hits30d"`
	Hits90d           int             `json:"hits90d"`
	HitsAllTime       int             `json:"hitsAllTime"`
	CreatedAt         string          `json:"createdAt"`
	UpdatedAt         string          `json:"updatedAt"`
	Raw               json.RawMessage `json:"-"`
}

func (u *UnifiedExclusion) UnmarshalJSON(b []byte) error {
	type alias UnifiedExclusion
	if err := json.Unmarshal(b, (*alias)(u)); err != nil {
		return err
	}
	u.Raw = append(u.Raw[:0:0], b...)
	return nil
}

// UnifiedExclusionListParams are query parameters for listing unified exclusions.
type UnifiedExclusionListParams struct {
	AccountIDs       []string
	SiteIDs          []string
	GroupIDs         []string
	IDs              []string
	OSTypes          []string
	Source           []string
	ModeType         []string
	ThreatType       []string
	Engines          []string
	InteractionLevel []string
	Conditions       []string
	NameContains     []string
	ValueContains    []string
	IncludeParents   *bool
	IncludeChildren  *bool
	Imported         *bool
	Tenant           *bool
	Limit            int
	Cursor           string
	SortBy           string
	SortOrder        string
	CountOnly        bool
}

func (p *UnifiedExclusionListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "ids", p.IDs)
	addCSV(v, "osTypes", p.OSTypes)
	addCSV(v, "source", p.Source)
	addCSV(v, "modeType", p.ModeType)
	addCSV(v, "threatType", p.ThreatType)
	addCSV(v, "engines", p.Engines)
	addCSV(v, "interactionLevel", p.InteractionLevel)
	addCSV(v, "conditions", p.Conditions)
	addCSV(v, "exclusionName__contains", p.NameContains)
	addCSV(v, "value__contains", p.ValueContains)
	addBool(v, "includeParents", p.IncludeParents)
	addBool(v, "includeChildren", p.IncludeChildren)
	addBool(v, "imported", p.Imported)
	addBool(v, "tenant", p.Tenant)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	if p.CountOnly {
		v.Set("countOnly", "true")
	}
	return v
}

// UnifiedExclusionsList returns a paginated list of unified exclusions.
func (c *Client) UnifiedExclusionsList(ctx context.Context, params *UnifiedExclusionListParams) ([]UnifiedExclusion, *Pagination, error) {
	return list[UnifiedExclusion](c, ctx, "/unified-exclusions", params.values())
}

// UnifiedExclusionsCount returns the count of unified exclusions matching the filter.
func (c *Client) UnifiedExclusionsCount(ctx context.Context, params *UnifiedExclusionListParams) (int, error) {
	if params == nil {
		params = &UnifiedExclusionListParams{}
	}
	params.CountOnly = true
	_, pag, err := list[UnifiedExclusion](c, ctx, "/unified-exclusions", params.values())
	if err != nil {
		return 0, err
	}
	return pag.TotalItems, nil
}

// UnifiedExclusionCreate is the request body for creating a unified exclusion.
type UnifiedExclusionCreate struct {
	ExclusionName     string   `json:"exclusionName"`
	OSType            string   `json:"osType"`
	ThreatType        string   `json:"threatType"`
	ModeType          string   `json:"modeType"`
	Reason            string   `json:"reason"`
	Type              string   `json:"type,omitempty"`
	Description       string   `json:"description,omitempty"`
	InteractionLevel  string   `json:"interactionLevel,omitempty"`
	Source            string   `json:"source,omitempty"`
	Value             any      `json:"value,omitempty"`
	PathExclusionType string   `json:"pathExclusionType,omitempty"`
	Engines           string   `json:"engines,omitempty"`
	ChildProcess      *bool    `json:"childProcess,omitempty"`
	Actions           []string `json:"actions,omitempty"`
	TagIDs            []string `json:"tagIds,omitempty"`
}

// UnifiedExclusionScope defines the scope for a unified exclusion operation.
type UnifiedExclusionScope struct {
	ScopeLevel   string `json:"scopeLevel"`
	ScopeLevelID string `json:"scopeLevelId,omitempty"`
}

type unifiedExclusionCreateRequest struct {
	Data   UnifiedExclusionCreate `json:"data"`
	Filter UnifiedExclusionScope  `json:"filter"`
}

// UnifiedExclusionsCreate creates a unified exclusion.
func (c *Client) UnifiedExclusionsCreate(ctx context.Context, scope UnifiedExclusionScope, data UnifiedExclusionCreate) (*UnifiedExclusion, error) {
	req := unifiedExclusionCreateRequest{Data: data, Filter: scope}
	var resp listResponse[UnifiedExclusion]
	if err := c.post(ctx, "/unified-exclusions", req, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("mgmt: unified exclusion not returned after create")
	}
	return &resp.Data[0], nil
}

// UnifiedExclusionsExport exports unified exclusions as raw JSON.
func (c *Client) UnifiedExclusionsExport(ctx context.Context, params *UnifiedExclusionListParams) (json.RawMessage, error) {
	var raw json.RawMessage
	if err := c.get(ctx, "/unified-exclusions/export", params.values(), &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
