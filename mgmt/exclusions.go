package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// Exclusion is a SentinelOne exclusion entry.
type Exclusion struct {
	ID                string   `json:"id"`
	Type              string   `json:"type"`
	Value             string   `json:"value"`
	Source            string   `json:"source"`
	OSType            string   `json:"osType"`
	Mode              string   `json:"mode"`
	Description       string   `json:"description"`
	ScopeName         string   `json:"scopeName"`
	ScopePath         string   `json:"scopePath"`
	PathExclusionType string   `json:"pathExclusionType"`
	ApplicationName   string   `json:"applicationName"`
	Actions           []string `json:"actions"`
	Imported          bool     `json:"imported"`
	Inject            bool     `json:"inject"`
	UserID            string   `json:"userId"`
	UserName          string   `json:"userName"`
	CreatedAt         string   `json:"createdAt"`
	UpdatedAt         string   `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (e *Exclusion) UnmarshalJSON(b []byte) error {
	type alias Exclusion
	if err := json.Unmarshal(b, (*alias)(e)); err != nil {
		return err
	}
	e.Raw = append(e.Raw[:0:0], b...)
	return nil
}

// ExclusionListParams are query parameters for listing exclusions.
type ExclusionListParams struct {
	SiteIDs    []string
	GroupIDs   []string
	AccountIDs []string
	Types      []string
	OSTypes    []string
	Query      string
	Limit      int
	Cursor     string
	SortBy     string
	SortOrder  string
}

func (p *ExclusionListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "type", p.Types)
	addCSV(v, "osTypes", p.OSTypes)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// ExclusionsList returns a paginated list of exclusions.
func (c *Client) ExclusionsList(ctx context.Context, params *ExclusionListParams) ([]Exclusion, *Pagination, error) {
	return list[Exclusion](c, ctx, "/exclusions", params.values())
}

// ExclusionsGet returns a single exclusion by ID.
func (c *Client) ExclusionsGet(ctx context.Context, id string) (*Exclusion, error) {
	return getByID[Exclusion](c, ctx, "/exclusions", "exclusion", id)
}
