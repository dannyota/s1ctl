package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Site is a SentinelOne site.
type Site struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	AccountID         string `json:"accountId"`
	AccountName       string `json:"accountName"`
	State             string `json:"state"`
	SiteType          string `json:"siteType"`
	TotalLicenses     int    `json:"totalLicenses"`
	ActiveLicenses    int    `json:"activeLicenses"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	Expiration        string `json:"expiration"`
	IsDefault         bool   `json:"isDefault"`
	RegistrationToken string `json:"registrationToken"`
	Description       string `json:"description"`
	UnlimitedLicenses bool   `json:"unlimitedLicenses"`

	Raw json.RawMessage `json:"-"`
}

func (s *Site) UnmarshalJSON(b []byte) error {
	type alias Site
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// SiteListParams are query parameters for listing sites.
type SiteListParams struct {
	AccountIDs []string
	States     []string
	SiteType   string
	Query      string
	Limit      int
	Cursor     string
	SortBy     string
	SortOrder  string
	CountOnly  bool
}

func (p *SiteListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "states", p.States)
	addString(v, "siteType", p.SiteType)
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

type sitesResponse struct {
	Data struct {
		Sites      []Site     `json:"sites"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}

// SitesList returns a paginated list of sites.
func (c *Client) SitesList(ctx context.Context, params *SiteListParams) ([]Site, *Pagination, error) {
	var resp sitesResponse
	if err := c.get(ctx, "/sites", params.values(), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data.Sites, &resp.Data.Pagination, nil
}

// SitesCount returns the count of sites matching the filter.
func (c *Client) SitesCount(ctx context.Context, params *SiteListParams) (int, error) {
	if params == nil {
		params = &SiteListParams{}
	}
	params.CountOnly = true
	var resp sitesResponse
	if err := c.get(ctx, "/sites", params.values(), &resp); err != nil {
		return 0, err
	}
	return resp.Data.Pagination.TotalItems, nil
}

// SitesGet returns a single site by ID.
func (c *Client) SitesGet(ctx context.Context, id string) (*Site, error) {
	params := url.Values{}
	params.Set("siteIds", id)
	var resp sitesResponse
	if err := c.get(ctx, "/sites", params, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data.Sites) == 0 {
		return nil, fmt.Errorf("mgmt: site %s not found", id)
	}
	return &resp.Data.Sites[0], nil
}

// SiteCreate is the request body for creating a site.
type SiteCreate struct {
	Name              string `json:"name"`
	AccountID         string `json:"accountId"`
	SiteType          string `json:"siteType,omitempty"`
	Description       string `json:"description,omitempty"`
	Expiration        string `json:"expiration,omitempty"`
	UnlimitedLicenses bool   `json:"unlimitedLicenses"`
	TotalLicenses     int    `json:"totalLicenses"`
}

// SiteUpdate is the request body for updating a site.
type SiteUpdate struct {
	Name              *string `json:"name,omitempty"`
	Description       *string `json:"description,omitempty"`
	Expiration        *string `json:"expiration,omitempty"`
	UnlimitedLicenses *bool   `json:"unlimitedLicenses,omitempty"`
	TotalLicenses     *int    `json:"totalLicenses,omitempty"`
}

// SitesCreate creates a site.
func (c *Client) SitesCreate(ctx context.Context, data SiteCreate) (*Site, error) {
	return create[Site](c, ctx, "/sites", data)
}

// SitesUpdate updates a site.
func (c *Client) SitesUpdate(ctx context.Context, id string, data SiteUpdate) (*Site, error) {
	return update[Site](c, ctx, fmt.Sprintf("/sites/%s", url.PathEscape(id)), data)
}

// SitesDelete deletes a site.
func (c *Client) SitesDelete(ctx context.Context, id string) error {
	return c.delete(ctx, fmt.Sprintf("/sites/%s", url.PathEscape(id)), nil, nil)
}
