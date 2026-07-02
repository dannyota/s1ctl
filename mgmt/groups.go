package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
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

// GroupsCount returns the count of groups matching the filter.
func (c *Client) GroupsCount(ctx context.Context, params *GroupListParams) (int, error) {
	if params == nil {
		params = &GroupListParams{}
	}
	params.CountOnly = true
	_, pag, err := list[Group](c, ctx, "/groups", params.values())
	if err != nil {
		return 0, err
	}
	return pag.TotalItems, nil
}

// GroupsGet returns a single group by ID.
func (c *Client) GroupsGet(ctx context.Context, id string) (*Group, error) {
	return getByID[Group](c, ctx, "/groups", "group", id)
}

// GroupCreate is the request body for creating a group.
type GroupCreate struct {
	Name        string `json:"name"`
	SiteID      string `json:"siteId"`
	Description string `json:"description,omitempty"`
}

// GroupUpdate is the request body for updating a group.
type GroupUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// GroupsCreate creates a group.
func (c *Client) GroupsCreate(ctx context.Context, siteID string, data GroupCreate) (*Group, error) {
	data.SiteID = siteID
	return create[Group](c, ctx, "/groups", data)
}

// GroupsUpdate updates a group.
func (c *Client) GroupsUpdate(ctx context.Context, id string, data GroupUpdate) (*Group, error) {
	return update[Group](c, ctx, fmt.Sprintf("/groups/%s", url.PathEscape(id)), data)
}

// GroupsDelete deletes a group.
func (c *Client) GroupsDelete(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/groups/%s", url.PathEscape(id)), nil, nil)
}
