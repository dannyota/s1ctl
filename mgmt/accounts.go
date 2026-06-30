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

// AccountsGet returns a single account by ID (uses path param, not ?ids=).
func (c *Client) AccountsGet(ctx context.Context, id string) (*Account, error) {
	items, _, err := list[Account](c, ctx, fmt.Sprintf("/accounts/%s", id), nil)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("mgmt: account %s not found", id)
	}
	return &items[0], nil
}
