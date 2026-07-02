package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// FirewallCategory selects a rule category within the shared firewall-control
// endpoint family. The default (empty or "firewall") targets the standard
// firewall; "network-quarantine" targets Network Quarantine, which is the same
// operations addressed under an extra path segment.
type FirewallCategory string

const (
	// FirewallCategoryFirewall is the default category (bare /firewall-control paths).
	FirewallCategoryFirewall FirewallCategory = "firewall"
	// FirewallCategoryNetworkQuarantine addresses Network Quarantine rules.
	FirewallCategoryNetworkQuarantine FirewallCategory = "network-quarantine"
)

// firewallPath builds a firewall-control path for the given category. The
// default/firewall category yields the bare "/firewall-control" paths; any
// other category inserts its segment, e.g. "/firewall-control/network-quarantine".
func firewallPath(cat FirewallCategory, suffix string) string {
	base := "/firewall-control"
	if cat != "" && cat != FirewallCategoryFirewall {
		base += "/" + string(cat)
	}
	return base + suffix
}

// FirewallDirection is the traffic direction of a firewall rule.
type FirewallDirection string

const (
	FirewallDirectionAny      FirewallDirection = "any"
	FirewallDirectionInbound  FirewallDirection = "inbound"
	FirewallDirectionOutbound FirewallDirection = "outbound"
)

// FirewallAction is the action taken by a firewall rule.
type FirewallAction string

const (
	FirewallActionAllow FirewallAction = "Allow"
	FirewallActionBlock FirewallAction = "Block"
)

// FirewallStatus is the status of a firewall rule.
type FirewallStatus string

const (
	FirewallStatusEnabled  FirewallStatus = "Enabled"
	FirewallStatusDisabled FirewallStatus = "Disabled"
)

// FirewallHostType is the type of a host matcher in a firewall rule.
type FirewallHostType string

const (
	FirewallHostAny       FirewallHostType = "any"
	FirewallHostCIDR      FirewallHostType = "cidr"
	FirewallHostRange     FirewallHostType = "range"
	FirewallHostAddresses FirewallHostType = "addresses"
	FirewallHostFQDN      FirewallHostType = "fqdn"
)

// FirewallPortType is the type of a port matcher in a firewall rule.
type FirewallPortType string

const (
	FirewallPortAny   FirewallPortType = "any"
	FirewallPortPorts FirewallPortType = "ports"
	FirewallPortRange FirewallPortType = "range"
)

// FirewallLocationType is the type of a location matcher in a firewall rule.
type FirewallLocationType string

const (
	FirewallLocationAll      FirewallLocationType = "all"
	FirewallLocationSpecific FirewallLocationType = "specific"
	FirewallLocationFallback FirewallLocationType = "fallback"
)

// FirewallAppType is the type of an application matcher in a firewall rule.
type FirewallAppType string

const (
	FirewallAppAny    FirewallAppType = "any"
	FirewallAppPath   FirewallAppType = "path"
	FirewallAppSHA1   FirewallAppType = "sha1"
	FirewallAppSystem FirewallAppType = "system"
)

// FirewallHost describes a host matcher (local or remote).
type FirewallHost struct {
	Type   FirewallHostType `json:"type"`
	Values []string         `json:"values,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (f *FirewallHost) UnmarshalJSON(b []byte) error {
	type alias FirewallHost
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FirewallPort describes a port matcher (local or remote).
type FirewallPort struct {
	Type   FirewallPortType `json:"type"`
	Values []string         `json:"values,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (f *FirewallPort) UnmarshalJSON(b []byte) error {
	type alias FirewallPort
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FirewallLocation describes a location matcher.
type FirewallLocation struct {
	Type   FirewallLocationType `json:"type"`
	Values []string             `json:"values,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (f *FirewallLocation) UnmarshalJSON(b []byte) error {
	type alias FirewallLocation
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FirewallApplication describes an application matcher.
type FirewallApplication struct {
	Type   FirewallAppType `json:"type"`
	Values []string        `json:"values,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (f *FirewallApplication) UnmarshalJSON(b []byte) error {
	type alias FirewallApplication
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FirewallRule is a SentinelOne firewall rule.
type FirewallRule struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Status       FirewallStatus       `json:"status"`
	Action       FirewallAction       `json:"action"`
	Direction    FirewallDirection    `json:"direction"`
	Protocol     string               `json:"protocol"`
	OSType       string               `json:"osType"`
	OSTypes      []string             `json:"osTypes"`
	Order        int                  `json:"order"`
	Application  *FirewallApplication `json:"application,omitempty"`
	LocalHost    *FirewallHost        `json:"localHost,omitempty"`
	LocalPort    *FirewallPort        `json:"localPort,omitempty"`
	RemoteHosts  []FirewallHost       `json:"remoteHosts,omitempty"`
	RemotePort   *FirewallPort        `json:"remotePort,omitempty"`
	Location     *FirewallLocation    `json:"location,omitempty"`
	Scope        string               `json:"scope"`
	ScopeID      string               `json:"scopeId"`
	Editable     bool                 `json:"editable"`
	RuleCategory string               `json:"ruleCategory"`
	TagIDs       []string             `json:"tagIds"`
	TagNames     []string             `json:"tagNames"`
	Creator      string               `json:"creator"`
	CreatorID    string               `json:"creatorId"`
	CreatedAt    string               `json:"createdAt"`
	UpdatedAt    string               `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (f *FirewallRule) UnmarshalJSON(b []byte) error {
	type alias FirewallRule
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FirewallRuleListParams are query parameters for listing firewall rules.
type FirewallRuleListParams struct {
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
	Query      string
	Limit      int
	Cursor     string
}

func (p *FirewallRuleListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// FirewallRulesList returns a paginated list of firewall rules.
func (c *Client) FirewallRulesList(ctx context.Context, params *FirewallRuleListParams) ([]FirewallRule, *Pagination, error) {
	return c.firewallRulesList(ctx, FirewallCategoryFirewall, params)
}

func (c *Client) firewallRulesList(ctx context.Context, cat FirewallCategory, params *FirewallRuleListParams) ([]FirewallRule, *Pagination, error) {
	return list[FirewallRule](c, ctx, firewallPath(cat, ""), params.values())
}

// FirewallRulesGet returns a single firewall rule by ID.
func (c *Client) FirewallRulesGet(ctx context.Context, id string) (*FirewallRule, error) {
	return c.firewallRulesGet(ctx, FirewallCategoryFirewall, id)
}

func (c *Client) firewallRulesGet(ctx context.Context, cat FirewallCategory, id string) (*FirewallRule, error) {
	return getByID[FirewallRule](c, ctx, firewallPath(cat, ""), "firewall rule", id)
}

// FirewallRuleCreate is the request body for creating or updating a firewall rule.
type FirewallRuleCreate struct {
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	Direction   FirewallDirection    `json:"direction"`
	Protocol    string               `json:"protocol,omitempty"`
	OSTypes     []string             `json:"osTypes,omitempty"`
	Action      FirewallAction       `json:"action"`
	Status      FirewallStatus       `json:"status"`
	Application *FirewallApplication `json:"application,omitempty"`
	LocalHost   *FirewallHost        `json:"localHost,omitempty"`
	LocalPort   *FirewallPort        `json:"localPort,omitempty"`
	RemoteHosts []FirewallHost       `json:"remoteHosts,omitempty"`
	RemotePort  *FirewallPort        `json:"remotePort,omitempty"`
	Location    *FirewallLocation    `json:"location,omitempty"`
	TagIDs      []string             `json:"tagIds,omitempty"`
}

type firewallCreateRequest struct {
	Filter struct {
		SiteIDs    []string `json:"siteIds,omitempty"`
		AccountIDs []string `json:"accountIds,omitempty"`
		GroupIDs   []string `json:"groupIds,omitempty"`
	} `json:"filter"`
	Data FirewallRuleCreate `json:"data"`
}

// FirewallRulesCreate creates a firewall rule.
func (c *Client) FirewallRulesCreate(ctx context.Context, scope FirewallRuleScope, data FirewallRuleCreate) (*FirewallRule, error) {
	return c.firewallRulesCreate(ctx, FirewallCategoryFirewall, scope, data)
}

func (c *Client) firewallRulesCreate(ctx context.Context, cat FirewallCategory, scope FirewallRuleScope, data FirewallRuleCreate) (*FirewallRule, error) {
	req := firewallCreateRequest{Data: data}
	req.Filter.SiteIDs = scope.SiteIDs
	req.Filter.AccountIDs = scope.AccountIDs
	req.Filter.GroupIDs = scope.GroupIDs
	var resp singleResponse[FirewallRule]
	if err := c.post(ctx, firewallPath(cat, ""), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// FirewallRulesUpdate updates a firewall rule by ID.
func (c *Client) FirewallRulesUpdate(ctx context.Context, id string, data FirewallRuleCreate) (*FirewallRule, error) {
	return c.firewallRulesUpdate(ctx, FirewallCategoryFirewall, id, data)
}

func (c *Client) firewallRulesUpdate(ctx context.Context, cat FirewallCategory, id string, data FirewallRuleCreate) (*FirewallRule, error) {
	return update[FirewallRule](c, ctx, firewallPath(cat, "/"+url.PathEscape(id)), data)
}

// FirewallRuleScope identifies the scope for creating a firewall rule.
type FirewallRuleScope struct {
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
}
