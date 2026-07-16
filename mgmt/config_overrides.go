package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ConfigOverrideOSType is the operating system a config override targets.
type ConfigOverrideOSType string

// Config override OS types.
const (
	ConfigOverrideOSLinux         ConfigOverrideOSType = "linux"
	ConfigOverrideOSMacOS         ConfigOverrideOSType = "macos"
	ConfigOverrideOSWindows       ConfigOverrideOSType = "windows"
	ConfigOverrideOSWindowsLegacy ConfigOverrideOSType = "windows_legacy"
)

// ConfigOverrideVersionOption controls which agent versions are targeted.
type ConfigOverrideVersionOption string

// Config override version options.
const (
	ConfigOverrideVersionAll      ConfigOverrideVersionOption = "ALL"
	ConfigOverrideVersionSpecific ConfigOverrideVersionOption = "SPECIFIC"
)

// ConfigOverrideScope is the hierarchy level at which an override applies.
type ConfigOverrideScope string

// Config override scope levels.
const (
	ConfigOverrideScopeGroup   ConfigOverrideScope = "group"
	ConfigOverrideScopeSite    ConfigOverrideScope = "site"
	ConfigOverrideScopeAccount ConfigOverrideScope = "account"
	ConfigOverrideScopeTenant  ConfigOverrideScope = "tenant"
)

// ConfigOverrideScopeRef identifies the target scope object (site, group, or
// account) for a config override.
type ConfigOverrideScopeRef struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// ConfigOverride is a SentinelOne agent configuration override.
type ConfigOverride struct {
	ID            string                      `json:"id"`
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	Config        json.RawMessage             `json:"config"`
	OSType        ConfigOverrideOSType        `json:"osType"`
	AgentVersion  string                      `json:"agentVersion"`
	VersionOption ConfigOverrideVersionOption `json:"versionOption"`
	Scope         ConfigOverrideScope         `json:"scope"`
	Site          *ConfigOverrideScopeRef     `json:"site,omitempty"`
	Group         *ConfigOverrideScopeRef     `json:"group,omitempty"`
	Account       *ConfigOverrideScopeRef     `json:"account,omitempty"`
	Agent         *ConfigOverrideScopeRef     `json:"agent,omitempty"`
	CreatedAt     string                      `json:"createdAt"`
	UpdatedAt     string                      `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (c *ConfigOverride) UnmarshalJSON(data []byte) error {
	type alias ConfigOverride
	if err := json.Unmarshal(data, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], data...)
	return nil
}

// ConfigOverrideListParams are query parameters for listing config overrides.
type ConfigOverrideListParams struct {
	SiteIDs       []string
	AccountIDs    []string
	GroupIDs      []string
	IDs           []string
	AgentIDs      []string
	OSTypes       []string
	AgentVersions []string
	VersionOption string
	Query         string
	Tenant        *bool
	Limit         int
	Cursor        string
	SortBy        string
	SortOrder     string
}

func (p *ConfigOverrideListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "ids", p.IDs)
	addCSV(v, "agentIds", p.AgentIDs)
	addCSV(v, "osTypes", p.OSTypes)
	addCSV(v, "agentVersions", p.AgentVersions)
	addString(v, "versionOption", p.VersionOption)
	addString(v, "query", p.Query)
	addBool(v, "tenant", p.Tenant)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// ConfigOverrideCreateInput is the data payload for creating a config override.
type ConfigOverrideCreateInput struct {
	Name          string                       `json:"name"`
	Description   *string                      `json:"description,omitempty"`
	OSType        ConfigOverrideOSType         `json:"osType"`
	Config        json.RawMessage              `json:"config"`
	Scope         ConfigOverrideScope          `json:"scope"`
	AgentVersion  *string                      `json:"agentVersion,omitempty"`
	VersionOption *ConfigOverrideVersionOption `json:"versionOption,omitempty"`
	Site          *ConfigOverrideScopeRef      `json:"site,omitempty"`
	Group         *ConfigOverrideScopeRef      `json:"group,omitempty"`
	Account       *ConfigOverrideScopeRef      `json:"account,omitempty"`
}

// ConfigOverrideUpdateInput is the data payload for updating a config override.
// All fields are optional; only provided fields are changed.
type ConfigOverrideUpdateInput struct {
	Name          *string                      `json:"name,omitempty"`
	Description   *string                      `json:"description,omitempty"`
	OSType        *ConfigOverrideOSType        `json:"osType,omitempty"`
	Config        json.RawMessage              `json:"config,omitempty"`
	Scope         *ConfigOverrideScope         `json:"scope,omitempty"`
	AgentVersion  *string                      `json:"agentVersion,omitempty"`
	VersionOption *ConfigOverrideVersionOption `json:"versionOption,omitempty"`
	Site          *ConfigOverrideScopeRef      `json:"site,omitempty"`
	Group         *ConfigOverrideScopeRef      `json:"group,omitempty"`
	Account       *ConfigOverrideScopeRef      `json:"account,omitempty"`
}

// ConfigOverrideDeleteFilter is the filter for bulk-deleting config overrides.
type ConfigOverrideDeleteFilter struct {
	IDs              []string                     `json:"ids,omitempty"`
	AgentIDs         []string                     `json:"agentIds,omitempty"`
	SiteIDs          []string                     `json:"siteIds,omitempty"`
	AccountIDs       []string                     `json:"accountIds,omitempty"`
	GroupIDs         []string                     `json:"groupIds,omitempty"`
	OSTypes          []string                     `json:"osTypes,omitempty"`
	AgentVersions    []string                     `json:"agentVersions,omitempty"`
	VersionOption    *ConfigOverrideVersionOption `json:"versionOption,omitempty"`
	NameLike         string                       `json:"name__like,omitempty"`
	DescriptionLike  string                       `json:"description__like,omitempty"`
	Query            string                       `json:"query,omitempty"`
	CreatedAtGt      string                       `json:"createdAt__gt,omitempty"`
	CreatedAtGte     string                       `json:"createdAt__gte,omitempty"`
	CreatedAtLt      string                       `json:"createdAt__lt,omitempty"`
	CreatedAtLte     string                       `json:"createdAt__lte,omitempty"`
	CreatedAtBetween string                       `json:"createdAt__between,omitempty"`
	Tenant           *bool                        `json:"tenant,omitempty"`
}

// ConfigOverrideList returns a paginated list of config overrides.
func (c *Client) ConfigOverrideList(ctx context.Context, params *ConfigOverrideListParams) ([]ConfigOverride, *Pagination, error) {
	return list[ConfigOverride](c, ctx, "/config-override", params.values())
}

// ConfigOverrideGet returns a single config override by ID.
func (c *Client) ConfigOverrideGet(ctx context.Context, id string) (*ConfigOverride, error) {
	return getByID[ConfigOverride](c, ctx, "/config-override", "config override", id)
}

// ConfigOverrideCreate creates a new config override and returns the created object.
func (c *Client) ConfigOverrideCreate(ctx context.Context, input ConfigOverrideCreateInput) (*ConfigOverride, error) {
	return create[ConfigOverride](c, ctx, "/config-override", input)
}

// ConfigOverrideUpdate updates an existing config override by ID and returns
// the updated object.
func (c *Client) ConfigOverrideUpdate(ctx context.Context, id string, input ConfigOverrideUpdateInput) (*ConfigOverride, error) {
	return update[ConfigOverride](c, ctx, "/config-override/"+id, input)
}

// ConfigOverrideDelete deletes a single config override by ID.
func (c *Client) ConfigOverrideDelete(ctx context.Context, id string) error {
	var resp struct {
		Data struct {
			Success bool `json:"success"`
		} `json:"data"`
	}
	if err := c.queryRequest(ctx, "DELETE", "/config-override/"+id, nil, &resp); err != nil {
		return err
	}
	if !resp.Data.Success {
		return fmt.Errorf("mgmt: config override delete returned success=false")
	}
	return nil
}

// ConfigOverrideBulkDelete deletes config overrides matching the filter and
// returns the number of affected items.
func (c *Client) ConfigOverrideBulkDelete(ctx context.Context, filter ConfigOverrideDeleteFilter) (int, error) {
	req := struct {
		Filter ConfigOverrideDeleteFilter `json:"filter"`
	}{filter}
	var resp affectedResponse
	if err := c.jsonRequest(ctx, "DELETE", "/config-override", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}
