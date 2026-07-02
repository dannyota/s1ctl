package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// User is a SentinelOne user.
type User struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	Scope      string `json:"scope"`
	ScopeRoles []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"scopeRoles"`
	Source       string `json:"source"`
	TwoFaEnabled bool   `json:"twoFaEnabled"`
	DateJoined   string `json:"dateJoined"`
	LastLogin    string `json:"lastLogin"`

	Raw json.RawMessage `json:"-"`
}

func (u *User) UnmarshalJSON(b []byte) error {
	type alias User
	if err := json.Unmarshal(b, (*alias)(u)); err != nil {
		return err
	}
	u.Raw = append(u.Raw[:0:0], b...)
	return nil
}

// UserListParams are query parameters for listing users.
type UserListParams struct {
	SiteIDs    []string
	AccountIDs []string
	Query      string
	Limit      int
	Cursor     string
	SortBy     string
	SortOrder  string
}

func (p *UserListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// UsersList returns a paginated list of users.
func (c *Client) UsersList(ctx context.Context, params *UserListParams) ([]User, *Pagination, error) {
	return list[User](c, ctx, "/users", params.values())
}

// UsersGet returns a single user by ID.
func (c *Client) UsersGet(ctx context.Context, id string) (*User, error) {
	return getByID[User](c, ctx, "/users", "user", id)
}

// UsersDelete deletes a user.
func (c *Client) UsersDelete(ctx context.Context, id string) error {
	req := map[string]any{
		"filter": map[string]any{
			"ids": []string{id},
		},
		"data": map[string]any{},
	}
	return c.post(ctx, "/users/delete", req, nil)
}

// UserUpdate is the payload for updating a user. All fields are optional; only
// provided fields change. Pointer fields distinguish "unset" from "false".
type UserUpdate struct {
	FullName            string `json:"fullName,omitempty"`
	Email               string `json:"email,omitempty"`
	Scope               string `json:"scope,omitempty"`
	CanGenerateAPIToken *bool  `json:"canGenerateApiToken,omitempty"`
	AllowRemoteShell    *bool  `json:"allowRemoteShell,omitempty"`
}

// UsersUpdate updates a user and returns the updated resource.
func (c *Client) UsersUpdate(ctx context.Context, id string, data UserUpdate) (*User, error) {
	return update[User](c, ctx, "/users/"+url.PathEscape(id), data)
}

// UsersGenerateToken generates an API token for the authenticated user and
// returns it. The token is shown only once. forceLegacy requests a legacy token
// even when the auth-tokens switch is on.
func (c *Client) UsersGenerateToken(ctx context.Context, forceLegacy bool) (string, error) {
	data := map[string]any{}
	if forceLegacy {
		data["forceLegacy"] = true
	}
	req := map[string]any{"data": data}
	var resp struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := c.post(ctx, "/users/generate-api-token", req, &resp); err != nil {
		return "", err
	}
	return resp.Data.Token, nil
}

// UsersRevokeToken revokes the API token of the user with the given ID.
func (c *Client) UsersRevokeToken(ctx context.Context, id string) error {
	req := map[string]any{"data": map[string]any{"id": id}}
	return c.post(ctx, "/users/revoke-api-token", req, nil)
}

// UserTokenDetails is API-token metadata. The endpoints return only timestamps;
// Token is defensive — if the API ever echoes the secret it is captured here so
// callers can redact it rather than print it.
type UserTokenDetails struct {
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
	Token     string `json:"token"`

	Raw json.RawMessage `json:"-"`
}

func (d *UserTokenDetails) UnmarshalJSON(b []byte) error {
	type alias UserTokenDetails
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// UsersTokenDetails returns the API-token metadata for the authenticated user.
func (c *Client) UsersTokenDetails(ctx context.Context) (*UserTokenDetails, error) {
	req := map[string]any{"data": map[string]any{}}
	var resp singleResponse[UserTokenDetails]
	if err := c.post(ctx, "/users/api-token-details", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// UsersTokenDetailsByID returns the API-token metadata for a specific user.
func (c *Client) UsersTokenDetailsByID(ctx context.Context, id string) (*UserTokenDetails, error) {
	var resp singleResponse[UserTokenDetails]
	if err := c.get(ctx, "/users/"+url.PathEscape(id)+"/api-token-details", nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Users2FAEnable enables two-factor authentication for the user with the given ID.
func (c *Client) Users2FAEnable(ctx context.Context, id string) error {
	req := map[string]any{"data": map[string]any{"id": id}}
	return c.post(ctx, "/users/2fa/enable", req, nil)
}

// Users2FADisable disables two-factor authentication for the user with the given ID.
func (c *Client) Users2FADisable(ctx context.Context, id string) error {
	req := map[string]any{"data": map[string]any{"id": id}}
	return c.post(ctx, "/users/2fa/disable", req, nil)
}
