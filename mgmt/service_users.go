package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// ServiceUserScope is the access scope a service user is bound to.
type ServiceUserScope string

// Service user scopes.
const (
	ServiceUserScopeTenant  ServiceUserScope = "tenant"
	ServiceUserScopeAccount ServiceUserScope = "account"
	ServiceUserScopeSite    ServiceUserScope = "site"
)

// ServiceUserScopeRole binds a scope (account/site) to an RBAC role.
type ServiceUserScopeRole struct {
	ID          string   `json:"id"`
	RoleID      string   `json:"roleId"`
	RoleName    string   `json:"roleName"`
	Roles       []string `json:"roles"`
	Name        string   `json:"name"`
	AccountName string   `json:"accountName"`
}

// ServiceUserRef is a minimal user reference (creator/updater).
type ServiceUserRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ServiceUserAPIToken holds a service user's API token metadata. Value carries
// the secret token and is populated only on create; read responses return
// metadata only (CreatedAt/ExpiresAt).
type ServiceUserAPIToken struct {
	Value     string `json:"value"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
}

// ServiceUser is a SentinelOne service user: a non-interactive identity that
// authenticates with an API token.
type ServiceUser struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Scope          ServiceUserScope       `json:"scope"`
	ScopeRoles     []ServiceUserScopeRole `json:"scopeRoles"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
	LastActivation string                 `json:"lastActivation"`
	CreatedBy      ServiceUserRef         `json:"createdBy"`
	UpdatedBy      ServiceUserRef         `json:"updatedBy"`
	APIToken       ServiceUserAPIToken    `json:"apiToken"`

	Raw json.RawMessage `json:"-"`
}

func (s *ServiceUser) UnmarshalJSON(b []byte) error {
	type alias ServiceUser
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// ServiceUserScopeRoleInput assigns an RBAC role at a scope. ID is required for
// account/site scopes; tenant (global) roles omit it.
type ServiceUserScopeRoleInput struct {
	ID       string   `json:"id,omitempty"`
	RoleID   string   `json:"roleId,omitempty"`
	RoleName string   `json:"roleName,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

// ServiceUserCreate is the payload for creating a service user. Name, Scope, and
// ExpirationDate (RFC3339) are required by the API.
type ServiceUserCreate struct {
	Name           string                      `json:"name"`
	Description    string                      `json:"description,omitempty"`
	ExpirationDate string                      `json:"expirationDate"`
	Scope          ServiceUserScope            `json:"scope"`
	ScopeRoles     []ServiceUserScopeRoleInput `json:"scopeRoles,omitempty"`
}

// ServiceUserUpdate is the payload for updating a service user. All fields are
// optional; only provided fields change.
type ServiceUserUpdate struct {
	Description string                      `json:"description,omitempty"`
	Scope       ServiceUserScope            `json:"scope,omitempty"`
	ScopeRoles  []ServiceUserScopeRoleInput `json:"scopeRoles,omitempty"`
}

// ServiceUserToken is the result of generating an API token for a service user.
// Token is the secret and is returned exactly once.
type ServiceUserToken struct {
	Token     string `json:"token"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`

	Raw json.RawMessage `json:"-"`
}

func (t *ServiceUserToken) UnmarshalJSON(b []byte) error {
	type alias ServiceUserToken
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

// ServiceUserListParams are query parameters for listing service users.
type ServiceUserListParams struct {
	SiteIDs    []string
	AccountIDs []string
	IDs        []string
	RoleIDs    []string
	Query      string
	Limit      int
	Cursor     string
	SortBy     string
	SortOrder  string
}

func (p *ServiceUserListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "ids", p.IDs)
	addCSV(v, "roleIds", p.RoleIDs)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// ServiceUsersList returns a paginated list of service users.
func (c *Client) ServiceUsersList(ctx context.Context, params *ServiceUserListParams) ([]ServiceUser, *Pagination, error) {
	return list[ServiceUser](c, ctx, "/service-users", params.values())
}

// ServiceUsersGet returns a single service user by ID.
func (c *Client) ServiceUsersGet(ctx context.Context, id string) (*ServiceUser, error) {
	var resp singleResponse[ServiceUser]
	if err := c.get(ctx, "/service-users/"+url.PathEscape(id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ServiceUsersCreate creates a service user and returns it. The returned
// APIToken.Value holds the newly generated token (shown only here).
func (c *Client) ServiceUsersCreate(ctx context.Context, data ServiceUserCreate) (*ServiceUser, error) {
	return create[ServiceUser](c, ctx, "/service-users", data)
}

// ServiceUsersUpdate updates a service user and returns the updated resource.
func (c *Client) ServiceUsersUpdate(ctx context.Context, id string, data ServiceUserUpdate) (*ServiceUser, error) {
	return update[ServiceUser](c, ctx, "/service-users/"+url.PathEscape(id), data)
}

// ServiceUsersDelete deletes a single service user by ID.
func (c *Client) ServiceUsersDelete(ctx context.Context, id string) error {
	return c.delete(ctx, "/service-users/"+url.PathEscape(id))
}

// ServiceUsersBulkDelete deletes service users by ID and returns the affected count.
func (c *Client) ServiceUsersBulkDelete(ctx context.Context, ids []string) (int, error) {
	req := map[string]any{
		"filter": map[string]any{"ids": ids},
	}
	var resp affectedResponse
	if err := c.post(ctx, "/service-users/delete-service-users", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// ServiceUsersGenerateToken issues a new API token for a service user, expiring
// on expirationDate (RFC3339). The returned Token is shown only once.
func (c *Client) ServiceUsersGenerateToken(ctx context.Context, id, expirationDate string) (*ServiceUserToken, error) {
	req := map[string]any{
		"data": map[string]any{"expirationDate": expirationDate},
	}
	var resp singleResponse[ServiceUserToken]
	if err := c.post(ctx, "/service-users/"+url.PathEscape(id)+"/generate-api-token", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ServiceUsersExport returns the service users matching the filter as a report.
func (c *Client) ServiceUsersExport(ctx context.Context, params *ServiceUserListParams) ([]byte, error) {
	return c.getRaw(ctx, "/export/service-users", params.values())
}
