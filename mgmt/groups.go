package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// Group is a SentinelOne network group.
type Group struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	SiteID            string `json:"siteId"`
	Type              string `json:"type"`
	TotalAgents       int    `json:"totalAgents"`
	FilterID          string `json:"filterId"`
	FilterName        string `json:"filterName"`
	IsDefault         bool   `json:"isDefault"`
	Rank              int    `json:"rank"`
	RegistrationToken string `json:"registrationToken"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	Description       string `json:"description"`

	Raw json.RawMessage `json:"-"`
}

func (g *Group) UnmarshalJSON(b []byte) error {
	type alias Group
	if err := json.Unmarshal(b, (*alias)(g)); err != nil {
		return err
	}
	g.Raw = append(g.Raw[:0:0], b...)
	return nil
}

// GroupListParams are query parameters for listing groups.
type GroupListParams struct {
	SiteIDs   []string
	Types     []string
	Query     string
	Limit     int
	Cursor    string
	SortBy    string
	SortOrder string
	CountOnly bool
}

func (p *GroupListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "types", p.Types)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	if p.CountOnly {
		v.Set("countOnly", "true")
	}
	return v
}

// GroupsList returns a paginated list of groups.
func (c *Client) GroupsList(ctx context.Context, params *GroupListParams) ([]Group, *Pagination, error) {
	return list[Group](c, ctx, "/groups", params.values())
}

// GroupsGet returns a single group by ID.
func (c *Client) GroupsGet(ctx context.Context, id string) (*Group, error) {
	return getByID[Group](c, ctx, "/groups", "group", id)
}
