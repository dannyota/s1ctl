package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Account is a SentinelOne account.
type Account struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	State          string `json:"state"`
	AccountType    string `json:"accountType"`
	TotalLicenses  int    `json:"totalLicenses"`
	ActiveLicenses int    `json:"activeLicenses"`
	ActiveAgents   int    `json:"activeAgents"`
	NumberOfSites  int    `json:"numberOfSites"`
	Expiration     string `json:"expiration"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	UsageType      string `json:"usageType"`
	BillingMode    string `json:"billingMode"`

	Raw json.RawMessage `json:"-"`
}

func (a *Account) UnmarshalJSON(b []byte) error {
	type alias Account
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// AccountListParams are query parameters for listing accounts.
type AccountListParams struct {
	States    []string
	IDs       []string
	Query     string
	Limit     int
	Cursor    string
	SortBy    string
	SortOrder string
	CountOnly bool
}

func (p *AccountListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "states", p.States)
	addCSV(v, "ids", p.IDs)
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

// AccountsList returns a paginated list of accounts.
func (c *Client) AccountsList(ctx context.Context, params *AccountListParams) ([]Account, *Pagination, error) {
	return list[Account](c, ctx, "/accounts", params.values())
}

// AccountsCount returns the count of accounts matching the filter.
func (c *Client) AccountsCount(ctx context.Context, params *AccountListParams) (int, error) {
	if params == nil {
		params = &AccountListParams{}
	}
	params.CountOnly = true
	_, pag, err := list[Account](c, ctx, "/accounts", params.values())
	if err != nil {
		return 0, err
	}
	return pag.TotalItems, nil
}

// AccountsGet returns a single account by ID (uses path param, not ?ids=).
func (c *Client) AccountsGet(ctx context.Context, id string) (*Account, error) {
	items, _, err := list[Account](c, ctx, fmt.Sprintf("/accounts/%s", url.PathEscape(id)), nil)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("mgmt: account %s not found", id)
	}
	return &items[0], nil
}

// UninstallPassword is the agent uninstall password for an account. The
// password value is sensitive secret material.
type UninstallPassword struct {
	Password string `json:"password"`

	Raw json.RawMessage `json:"-"`
}

func (p *UninstallPassword) UnmarshalJSON(b []byte) error {
	type alias UninstallPassword
	if err := json.Unmarshal(b, (*alias)(p)); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

// UninstallPasswordMeta describes an account's uninstall password without
// exposing the secret itself.
type UninstallPasswordMeta struct {
	Expiration      string `json:"expiration"`
	Version         int    `json:"version"`
	CreatedAt       string `json:"createdAt"`
	LastRevoked     string `json:"lastRevoked"`
	RevokedByID     int    `json:"revokedById"`
	RevokedByName   string `json:"revokedByName"`
	GeneratedByID   int    `json:"generatedById"`
	GeneratedByName string `json:"generatedByName"`

	Raw json.RawMessage `json:"-"`
}

func (m *UninstallPasswordMeta) UnmarshalJSON(b []byte) error {
	type alias UninstallPasswordMeta
	if err := json.Unmarshal(b, (*alias)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0:0], b...)
	return nil
}

// AccountsReactivate reactivates an expired account. Pass unlimited=true to
// reactivate with no expiration, or a non-empty RFC3339 expiration to bound the
// license window. The caller chooses one; the spec requires the data wrapper.
func (c *Client) AccountsReactivate(ctx context.Context, id string, unlimited bool, expiration string) error {
	body := reactivateBody{Data: reactivateData{Unlimited: unlimited, Expiration: expiration}}
	return c.put(ctx, fmt.Sprintf("/accounts/%s/reactivate", url.PathEscape(id)), body, nil)
}

// AccountsExpireNow expires an account immediately.
func (c *Client) AccountsExpireNow(ctx context.Context, id string) error {
	return c.post(ctx, fmt.Sprintf("/accounts/%s/expire-now", url.PathEscape(id)), nil, nil)
}

// AccountsUninstallPasswordMetadata returns metadata about an account's
// uninstall password (no secret material).
func (c *Client) AccountsUninstallPasswordMetadata(ctx context.Context, id string) (*UninstallPasswordMeta, error) {
	var resp singleResponse[UninstallPasswordMeta]
	if err := c.get(ctx, fmt.Sprintf("/accounts/%s/uninstall-password/metadata", url.PathEscape(id)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AccountsUninstallPasswordView returns the account's current uninstall
// password. The returned value is sensitive.
func (c *Client) AccountsUninstallPasswordView(ctx context.Context, id string) (*UninstallPassword, error) {
	var resp singleResponse[UninstallPassword]
	if err := c.get(ctx, fmt.Sprintf("/accounts/%s/uninstall-password/view", url.PathEscape(id)), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AccountsUninstallPasswordGenerate generates (or regenerates) the account's
// uninstall password and returns the resulting metadata. Per the spec this
// endpoint returns metadata, not the password itself — read it back with
// AccountsUninstallPasswordView. The spec requires data.expiration
// (yyyy-mm-dd), so it is always sent; callers must supply a non-empty value.
func (c *Client) AccountsUninstallPasswordGenerate(ctx context.Context, id, expiration string) (*UninstallPasswordMeta, error) {
	body := map[string]any{"data": map[string]any{"expiration": expiration}}
	var resp singleResponse[UninstallPasswordMeta]
	if err := c.post(ctx, fmt.Sprintf("/accounts/%s/uninstall-password/generate", url.PathEscape(id)), body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AccountsUninstallPasswordRevoke revokes the account's uninstall password.
func (c *Client) AccountsUninstallPasswordRevoke(ctx context.Context, id string) error {
	return c.post(ctx, fmt.Sprintf("/accounts/%s/uninstall-password/revoke", url.PathEscape(id)), nil, nil)
}
