package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Singularity Marketplace base path (relative to /web/api/v2.1).
const marketplaceBase = "/singularity-marketplace"

// --- Response structs ---

// MarketplaceCatalogItem is a catalog application in the Singularity Marketplace.
type MarketplaceCatalogItem struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Key         string          `json:"key"`
	Category    string          `json:"category"`
	CategoryID  string          `json:"categoryId"`
	Description string          `json:"description"`
	Summary     string          `json:"summary"`
	Type        string          `json:"type"`
	Installed   bool            `json:"installed"`
	ToggleState string          `json:"toggleState"`
	Raw         json.RawMessage `json:"-"`
}

func (m *MarketplaceCatalogItem) UnmarshalJSON(b []byte) error {
	type alias MarketplaceCatalogItem
	if err := json.Unmarshal(b, (*alias)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0:0], b...)
	return nil
}

// MarketplaceAppScope is a per-installation scope entry for a marketplace app.
type MarketplaceAppScope struct {
	ID                      string          `json:"id"`
	ApplicationInstanceName string          `json:"applicationInstanceName"`
	Status                  string          `json:"status"`
	ScopeLevel              string          `json:"scopeLevel"`
	SiteID                  string          `json:"siteId"`
	Raw                     json.RawMessage `json:"-"`
}

func (m *MarketplaceAppScope) UnmarshalJSON(b []byte) error {
	type alias MarketplaceAppScope
	if err := json.Unmarshal(b, (*alias)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0:0], b...)
	return nil
}

// MarketplaceApp is an installed application in the Singularity Marketplace.
type MarketplaceApp struct {
	ApplicationCatalogID string                `json:"applicationCatalogId"`
	Name                 string                `json:"name"`
	HasAlert             bool                  `json:"hasAlert"`
	LastInstalledAt      string                `json:"lastInstalledAt"`
	Scopes               []MarketplaceAppScope `json:"scopes"`
	Raw                  json.RawMessage       `json:"-"`
}

func (m *MarketplaceApp) UnmarshalJSON(b []byte) error {
	type alias MarketplaceApp
	if err := json.Unmarshal(b, (*alias)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0:0], b...)
	return nil
}

// --- List params ---

// MarketplaceCatalogListParams are the query parameters for listing catalog applications.
type MarketplaceCatalogListParams struct {
	ID                  string
	CategoryContains    string
	NameContains        string
	DescriptionContains string
	Query               string
	CategoryIDs         []string
	Cursor              string
	Limit               int
	SortBy              string
	SortOrder           string
}

func (p *MarketplaceCatalogListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "id", p.ID)
	addString(v, "category__contains", p.CategoryContains)
	addString(v, "name__contains", p.NameContains)
	addString(v, "description__contains", p.DescriptionContains)
	addString(v, "query", p.Query)
	addCSV(v, "categoryIds", p.CategoryIDs)
	addString(v, "cursor", p.Cursor)
	addInt(v, "limit", p.Limit)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// MarketplaceAppListParams are the query parameters for listing installed applications.
type MarketplaceAppListParams struct {
	ApplicationCatalogID string
	ID                   string
	NameContains         string
	CreatorContains      string
	Query                string
	AccountIDs           []string
	SiteIDs              []string
	Cursor               string
	Limit                int
	CountOnly            *bool
	SortBy               string
	SortOrder            string
}

func (p *MarketplaceAppListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "applicationCatalogId", p.ApplicationCatalogID)
	addString(v, "id", p.ID)
	addString(v, "name__contains", p.NameContains)
	addString(v, "creator__contains", p.CreatorContains)
	addString(v, "query", p.Query)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "siteIds", p.SiteIDs)
	addString(v, "cursor", p.Cursor)
	addInt(v, "limit", p.Limit)
	addBool(v, "countOnly", p.CountOnly)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// --- Request structs ---

// MarketplaceConfig is a configuration key-value pair for marketplace apps.
type MarketplaceConfig struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// MarketplaceScopeFilter scopes a marketplace mutation to specific resources.
type MarketplaceScopeFilter struct {
	ApplicationCatalogID string   `json:"applicationCatalogId,omitempty"`
	IDs                  []string `json:"ids,omitempty"`
	ApplicationID        string   `json:"applicationId,omitempty"`
	AccountIDs           []string `json:"accountIds,omitempty"`
	SiteIDs              []string `json:"siteIds,omitempty"`
	GroupIDs             []string `json:"groupIds,omitempty"`
	Tenant               *bool    `json:"tenant,omitempty"`
}

// MarketplaceDeleteFilter scopes a marketplace delete to specific resources.
// The delete endpoint uses different JSON keys than the shared scope filter.
type MarketplaceDeleteFilter struct {
	ID                   []string `json:"id,omitempty"`
	ApplicationCatalogID []string `json:"application_catalog_id,omitempty"`
	NameContains         string   `json:"name__contains,omitempty"`
	CreatorContains      string   `json:"creator__contains,omitempty"`
	Query                string   `json:"query,omitempty"`
	SiteIDs              []string `json:"siteIds,omitempty"`
	AccountIDs           []string `json:"accountIds,omitempty"`
	GroupIDs             []string `json:"groupIds,omitempty"`
	Tenant               *bool    `json:"tenant,omitempty"`
}

// MarketplaceInstallInput is the request body for installing a marketplace application.
type MarketplaceInstallInput struct {
	Data struct {
		Name           string              `json:"applicationInstanceName"`
		Configurations []MarketplaceConfig `json:"configurations"`
	} `json:"data"`
	Filter MarketplaceScopeFilter `json:"filter"`
}

// MarketplaceUpdateInput is the request body for updating a marketplace application.
type MarketplaceUpdateInput struct {
	Data struct {
		NameMap        map[string]string   `json:"applicationIdToNameMap,omitempty"`
		Configurations []MarketplaceConfig `json:"configurations"`
	} `json:"data"`
	Filter MarketplaceScopeFilter `json:"filter"`
}

// --- Client methods ---

// MarketplaceCatalogList lists available applications in the Singularity Marketplace catalog.
func (c *Client) MarketplaceCatalogList(ctx context.Context, params *MarketplaceCatalogListParams) ([]MarketplaceCatalogItem, *Pagination, error) {
	return list[MarketplaceCatalogItem](c, ctx, marketplaceBase+"/applications-catalog", params.values())
}

// MarketplaceCatalogConfig returns the configuration schema fields for a catalog application.
func (c *Client) MarketplaceCatalogConfig(ctx context.Context, catalogID string) (json.RawMessage, error) {
	if catalogID == "" {
		return nil, fmt.Errorf("mgmt: catalogId is required")
	}
	path := marketplaceBase + "/applications-catalog/" + url.PathEscape(catalogID) + "/config"
	var resp struct {
		Data json.RawMessage `json:"data"`
	}
	if err := c.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// MarketplaceAppList lists installed marketplace applications.
func (c *Client) MarketplaceAppList(ctx context.Context, params *MarketplaceAppListParams) ([]MarketplaceApp, *Pagination, error) {
	return list[MarketplaceApp](c, ctx, marketplaceBase+"/applications", params.values())
}

// MarketplaceAppConfig returns the configuration for an installed marketplace application.
func (c *Client) MarketplaceAppConfig(ctx context.Context, appID string) (json.RawMessage, error) {
	if appID == "" {
		return nil, fmt.Errorf("mgmt: applicationId is required")
	}
	path := marketplaceBase + "/applications/" + url.PathEscape(appID) + "/config"
	var resp struct {
		Data json.RawMessage `json:"data"`
	}
	if err := c.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// MarketplaceAppLog returns log entries for an installed marketplace application.
func (c *Client) MarketplaceAppLog(ctx context.Context, appID string, onlyErrors *bool) ([]json.RawMessage, error) {
	if appID == "" {
		return nil, fmt.Errorf("mgmt: applicationId is required")
	}
	path := marketplaceBase + "/applications/" + url.PathEscape(appID) + "/log"
	params := url.Values{}
	addBool(params, "only_errors", onlyErrors)
	var resp struct {
		Data []json.RawMessage `json:"data"`
	}
	if err := c.get(ctx, path, params, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// MarketplaceInstall installs a marketplace application.
func (c *Client) MarketplaceInstall(ctx context.Context, input *MarketplaceInstallInput) error {
	if input == nil {
		return fmt.Errorf("mgmt: install input is required")
	}
	if input.Data.Configurations == nil {
		input.Data.Configurations = []MarketplaceConfig{}
	}
	return c.post(ctx, marketplaceBase+"/applications", input, nil)
}

// MarketplaceUpdate updates the configuration of an installed marketplace application.
func (c *Client) MarketplaceUpdate(ctx context.Context, input *MarketplaceUpdateInput) error {
	if input == nil {
		return fmt.Errorf("mgmt: update input is required")
	}
	if input.Data.Configurations == nil {
		input.Data.Configurations = []MarketplaceConfig{}
	}
	return c.put(ctx, marketplaceBase+"/applications", input, nil)
}

// MarketplaceDelete deletes an installed marketplace application.
func (c *Client) MarketplaceDelete(ctx context.Context, filter *MarketplaceDeleteFilter) error {
	if filter == nil {
		return fmt.Errorf("mgmt: delete filter is required")
	}
	if len(filter.ID) == 0 && len(filter.ApplicationCatalogID) == 0 &&
		len(filter.SiteIDs) == 0 && len(filter.AccountIDs) == 0 &&
		len(filter.GroupIDs) == 0 && filter.Query == "" {
		return fmt.Errorf("mgmt: delete filter must specify at least one selector")
	}
	body := struct {
		Filter *MarketplaceDeleteFilter `json:"filter"`
	}{Filter: filter}
	return c.jsonRequest(ctx, "DELETE", marketplaceBase+"/applications", body, nil)
}

// MarketplaceSetMode enables or disables installed marketplace applications.
// mode must be "enable" or "disable".
func (c *Client) MarketplaceSetMode(ctx context.Context, mode string, filter *MarketplaceScopeFilter) error {
	if mode != "enable" && mode != "disable" {
		return fmt.Errorf("mgmt: mode must be \"enable\" or \"disable\", got %q", mode)
	}
	if filter == nil {
		return fmt.Errorf("mgmt: filter is required")
	}
	body := struct {
		Filter *MarketplaceScopeFilter `json:"filter"`
	}{Filter: filter}
	return c.post(ctx, marketplaceBase+"/applications/"+url.PathEscape(mode), body, nil)
}
