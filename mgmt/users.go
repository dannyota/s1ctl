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
