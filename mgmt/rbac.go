package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// RoleScope is the scope level an RBAC role is defined at.
type RoleScope string

// RBAC role scope levels.
const (
	RoleScopeGroup   RoleScope = "Group"
	RoleScopeSite    RoleScope = "Site"
	RoleScopeAccount RoleScope = "Account"
	RoleScopeTenant  RoleScope = "Tenant"
)

// Role is a SentinelOne RBAC role.
//
// The permission set is modeled two ways by the API, so it is modeled two ways
// here. Reads (list/get/template) return a deeply nested "pages" tree — each
// page carries permissions with per-permission booleans and dependency IDs —
// which is captured verbatim as the Pages raw blob rather than fully typed:
// the tree is large, and writes never consume it. Writes (create/update) take a
// flat PermissionIDs slice on RoleData instead. Keeping Pages as json.RawMessage
// lets a role be pulled and re-serialized faithfully without pinning the SDK to
// an unstable permission-tree shape.
type Role struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Scope          RoleScope `json:"scope"`
	ScopeID        string    `json:"scopeId"`
	PredefinedRole bool      `json:"predefinedRole"`
	UsersInRoles   int       `json:"usersInRoles"`
	AccountName    string    `json:"accountName"`
	SiteName       string    `json:"siteName"`
	Creator        string    `json:"creator"`
	CreatorID      string    `json:"creatorId"`
	UpdatedBy      string    `json:"updatedBy"`
	UpdatedByID    string    `json:"updatedById"`
	CreatedAt      string    `json:"createdAt"`
	UpdatedAt      string    `json:"updatedAt"`

	// Pages is the nested permission tree returned by get/template (absent on
	// list). It round-trips untyped; writes use RoleData.PermissionIDs instead.
	Pages json.RawMessage `json:"pages,omitempty"`

	Raw json.RawMessage `json:"-"`
}

func (r *Role) UnmarshalJSON(b []byte) error {
	type alias Role
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// RoleListParams are query parameters for listing RBAC roles.
type RoleListParams struct {
	AccountIDs     []string
	SiteIDs        []string
	GroupIDs       []string
	IDs            []string
	Query          string
	Name           string
	PredefinedRole *bool // true: system roles only; false: custom roles only
	Tenant         *bool
	SortBy         string
	SortOrder      string
	Limit          int
	Skip           int
	Cursor         string
}

func (p *RoleListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "ids", p.IDs)
	addString(v, "query", p.Query)
	addString(v, "name", p.Name)
	addBool(v, "predefinedRole", p.PredefinedRole)
	addBool(v, "tenant", p.Tenant)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addInt(v, "limit", p.Limit)
	addInt(v, "skip", p.Skip)
	addString(v, "cursor", p.Cursor)
	return v
}

// RoleData is the declarative payload shared by role create and update: the
// role name, description, and the flat list of permission IDs it grants.
// PermissionIDs is omitted from the JSON when empty; whether the API then
// preserves the role's existing permissions or clears them on update is not
// documented, so callers that intend to keep permissions must send them.
type RoleData struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permissionIds,omitempty"`
}

// RoleScopeFilter targets the scope a role write applies to. Set Tenant for the
// global scope, or one or more of the ID slices for account/site/group scopes.
type RoleScopeFilter struct {
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	Tenant     bool     `json:"tenant,omitempty"`
}

// RoleCreate is the request body for creating a role. Both data and filter are
// required: filter names the scope the new role is created in.
type RoleCreate struct {
	Data   RoleData        `json:"data"`
	Filter RoleScopeFilter `json:"filter"`
}

// RoleUpdate is the request body for updating a role. Filter is optional (the
// role is already identified by its ID in the path).
type RoleUpdate struct {
	Data   RoleData         `json:"data"`
	Filter *RoleScopeFilter `json:"filter,omitempty"`
}

// RolesList returns a paginated list of RBAC roles.
func (c *Client) RolesList(ctx context.Context, params *RoleListParams) ([]Role, *Pagination, error) {
	return list[Role](c, ctx, "/rbac/roles", params.values())
}

// RoleGet returns a single role by ID, including its permission tree.
func (c *Client) RoleGet(ctx context.Context, id string) (*Role, error) {
	var resp singleResponse[Role]
	if err := c.get(ctx, fmt.Sprintf("/rbac/role/%s", url.PathEscape(id)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RoleTemplate returns the blank role template (description + permission tree
// with default values) used as a starting point for a new role.
func (c *Client) RoleTemplate(ctx context.Context) (*Role, error) {
	var resp singleResponse[Role]
	if err := c.get(ctx, "/rbac/role", nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RoleCreate creates a new RBAC role.
func (c *Client) RoleCreate(ctx context.Context, body RoleCreate) (*Role, error) {
	var resp singleResponse[Role]
	if err := c.post(ctx, "/rbac/role", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RoleUpdate updates an existing RBAC role.
func (c *Client) RoleUpdate(ctx context.Context, id string, body RoleUpdate) (*Role, error) {
	var resp singleResponse[Role]
	if err := c.put(ctx, fmt.Sprintf("/rbac/role/%s", url.PathEscape(id)), body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// RoleDelete deletes an RBAC role by ID. Users still assigned to the role are
// left without it; use the console to reassign them to a replacement role.
func (c *Client) RoleDelete(ctx context.Context, id string) error {
	req := map[string]any{"data": map[string]any{}}
	return c.jsonRequest(ctx, "DELETE", fmt.Sprintf("/rbac/role/%s", url.PathEscape(id)), req, nil)
}
