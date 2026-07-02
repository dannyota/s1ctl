package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// LocationOperator is the logical operator applied between a location's
// detection parameters.
type LocationOperator string

// Location detection operators.
const (
	LocationOperatorAll  LocationOperator = "all"
	LocationOperatorAny  LocationOperator = "any"
	LocationOperatorNone LocationOperator = "none"
)

// Location is a firewall location definition. Agents detect their location from
// endpoint network parameters (IP, DNS, NIC, registry key, or management
// connectivity) and apply Location Aware firewall rules that match.
//
// The six detection-parameter groups (dnsLookup, dnsServers, registryKeys,
// serverConnectivity, networkInterfaces, ipAddresses) are captured verbatim as
// raw blobs rather than fully typed: each is a nested object whose shape varies
// by parameter kind, and keeping them raw lets a location round-trip faithfully.
type Location struct {
	ID                 string           `json:"id"`
	Name               string           `json:"name"`
	Description        string           `json:"description"`
	Operator           LocationOperator `json:"operator"`
	Scope              string           `json:"scope"`
	ScopeID            string           `json:"scopeId"`
	ScopeName          string           `json:"scopeName"`
	Editable           bool             `json:"editable"`
	IsFallback         bool             `json:"isFallback"`
	ReportingAgents    int              `json:"reportingAgents"`
	ActiveFirewallRule int              `json:"activeFirewallRules"`
	CreatedAt          string           `json:"createdAt"`
	UpdatedAt          string           `json:"updatedAt"`

	DNSLookup          json.RawMessage `json:"dnsLookup,omitempty"`
	DNSServers         json.RawMessage `json:"dnsServers,omitempty"`
	RegistryKeys       json.RawMessage `json:"registryKeys,omitempty"`
	ServerConnectivity json.RawMessage `json:"serverConnectivity,omitempty"`
	NetworkInterfaces  json.RawMessage `json:"networkInterfaces,omitempty"`
	IPAddresses        json.RawMessage `json:"ipAddresses,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (l *Location) UnmarshalJSON(b []byte) error {
	type alias Location
	if err := json.Unmarshal(b, (*alias)(l)); err != nil {
		return err
	}
	l.Raw = append(l.Raw[:0:0], b...)
	return nil
}

// LocationListParams are query parameters for listing locations.
type LocationListParams struct {
	IDs        []string
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
	SortBy     string
	SortOrder  string
	Limit      int
	Cursor     string
}

func (p *LocationListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "ids", p.IDs)
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// LocationData is the declarative payload of a location: its name, description,
// the logical operator, and the detection-parameter groups.
type LocationData struct {
	Name               string           `json:"name"`
	Description        string           `json:"description,omitempty"`
	Operator           LocationOperator `json:"operator"`
	DNSLookup          json.RawMessage  `json:"dnsLookup,omitempty"`
	DNSServers         json.RawMessage  `json:"dnsServers,omitempty"`
	RegistryKeys       json.RawMessage  `json:"registryKeys,omitempty"`
	ServerConnectivity json.RawMessage  `json:"serverConnectivity,omitempty"`
	NetworkInterfaces  json.RawMessage  `json:"networkInterfaces,omitempty"`
	IPAddresses        json.RawMessage  `json:"ipAddresses,omitempty"`
}

// LocationScope targets the scope a new location is created in.
type LocationScope struct {
	SiteIDs    []string `json:"siteIds,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
}

// LocationCreate is the request body for creating a location.
type LocationCreate struct {
	Data   LocationData  `json:"data"`
	Filter LocationScope `json:"filter"`
}

// LocationUpdate is the request body for updating a location.
type LocationUpdate struct {
	Data LocationData `json:"data"`
}

// LocationsList returns a paginated list of locations.
func (c *Client) LocationsList(ctx context.Context, params *LocationListParams) ([]Location, *Pagination, error) {
	return list[Location](c, ctx, "/locations", params.values())
}

// LocationsCreate creates a location and returns the created object.
func (c *Client) LocationsCreate(ctx context.Context, body LocationCreate) (*Location, error) {
	var resp singleResponse[Location]
	if err := c.post(ctx, "/locations", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// LocationsUpdate updates an existing location.
func (c *Client) LocationsUpdate(ctx context.Context, id string, body LocationUpdate) (*Location, error) {
	var resp singleResponse[Location]
	if err := c.put(ctx, fmt.Sprintf("/locations/%s", url.PathEscape(id)), body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// LocationsDelete deletes one or more locations. The IDs are sent in the body:
// the delete endpoint is on the collection path, not per-location.
func (c *Client) LocationsDelete(ctx context.Context, ids []string) error {
	body := map[string]any{"data": map[string]any{"ids": ids}}
	return c.jsonRequest(ctx, http.MethodDelete, "/locations", body, nil)
}
