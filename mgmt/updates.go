package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// UpdatePackage is a SentinelOne agent update package.
type UpdatePackage struct {
	ID        string `json:"id"`
	FileName  string `json:"fileName"`
	Version   string `json:"version"`
	OSType    string `json:"osType"`
	Status    string `json:"status"`
	FileSize  int64  `json:"fileSize"`
	ScopeID   string `json:"scopeId"`
	ScopeName string `json:"scopeName"`
	CreatedAt string `json:"createdAt"`

	Raw json.RawMessage `json:"-"`
}

func (u *UpdatePackage) UnmarshalJSON(b []byte) error {
	type alias UpdatePackage
	if err := json.Unmarshal(b, (*alias)(u)); err != nil {
		return err
	}
	u.Raw = append(u.Raw[:0:0], b...)
	return nil
}

// UpdateListParams are query parameters for listing update packages.
type UpdateListParams struct {
	SiteIDs    []string
	AccountIDs []string
	OSTypes    []string
	Status     string
	Query      string
	Limit      int
	Cursor     string
}

func (p *UpdateListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "osTypes", p.OSTypes)
	addString(v, "status", p.Status)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// UpdatesList returns a paginated list of update packages.
func (c *Client) UpdatesList(ctx context.Context, params *UpdateListParams) ([]UpdatePackage, *Pagination, error) {
	return list[UpdatePackage](c, ctx, "/update/agent/packages", params.values())
}

// UpdatesGet returns a single update package by ID.
func (c *Client) UpdatesGet(ctx context.Context, id string) (*UpdatePackage, error) {
	return getByID[UpdatePackage](c, ctx, "/update/agent/packages", "update package", id)
}
