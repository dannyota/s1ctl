package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Filter is a saved endpoint filter. A saved filter pairs a name with a
// filterFields definition (the set of endpoint criteria to match) and can be
// used to run bulk agent actions or to back a dynamic group.
//
// filterFields is an open-ended set of endpoint criteria whose keys track the
// agents query surface; it is captured verbatim as a raw blob rather than fully
// typed so a filter round-trips faithfully without pinning the SDK to an
// unstable field set.
type Filter struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	ScopeID      string          `json:"scopeId"`
	ScopeLevel   string          `json:"scopeLevel"`
	FilterFields json.RawMessage `json:"filterFields,omitempty"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (f *Filter) UnmarshalJSON(b []byte) error {
	type alias Filter
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FilterListParams are query parameters for listing saved filters.
type FilterListParams struct {
	Query      string
	IDs        []string
	SiteIDs    []string
	AccountIDs []string
	SortBy     string
	SortOrder  string
	Limit      int
	Cursor     string
}

func (p *FilterListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "query", p.Query)
	addCSV(v, "ids", p.IDs)
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// FilterData is the declarative payload of a saved filter: its name and the
// filterFields criteria set.
type FilterData struct {
	Name         string          `json:"name"`
	FilterFields json.RawMessage `json:"filterFields,omitempty"`
}

// FilterScope targets the scope a new saved filter is created in. Leave both
// empty to create a global (tenant) filter.
type FilterScope struct {
	SiteIDs    []string `json:"siteIds,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
}

// FilterCreate is the request body for creating a saved filter.
type FilterCreate struct {
	Data   FilterData   `json:"data"`
	Filter *FilterScope `json:"filter,omitempty"`
}

// FilterUpdate is the request body for updating a saved filter. Supplying
// filterFields replaces the existing criteria set.
type FilterUpdate struct {
	Data FilterData `json:"data"`
}

// FiltersList returns a paginated list of saved filters.
func (c *Client) FiltersList(ctx context.Context, params *FilterListParams) ([]Filter, *Pagination, error) {
	return list[Filter](c, ctx, "/filters", params.values())
}

// FiltersCreate saves a new filter and returns the created object.
func (c *Client) FiltersCreate(ctx context.Context, body FilterCreate) (*Filter, error) {
	var resp singleResponse[Filter]
	if err := c.post(ctx, "/filters", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// FiltersUpdate updates an existing saved filter.
func (c *Client) FiltersUpdate(ctx context.Context, id string, body FilterUpdate) (*Filter, error) {
	var resp singleResponse[Filter]
	if err := c.put(ctx, fmt.Sprintf("/filters/%s", url.PathEscape(id)), body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// FiltersDelete deletes a saved filter by ID.
func (c *Client) FiltersDelete(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/filters/%s", url.PathEscape(id)))
}
