package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// BlocklistOSType is the operating system a blocklist item targets.
type BlocklistOSType string

// Blocklist OS types.
const (
	BlocklistOSLinux         BlocklistOSType = "linux"
	BlocklistOSMacOS         BlocklistOSType = "macos"
	BlocklistOSWindows       BlocklistOSType = "windows"
	BlocklistOSWindowsLegacy BlocklistOSType = "windows_legacy"
)

// BlocklistType is the restriction type. The API only supports black_hash; any
// other value creates an exclusion rather than a blocklist item.
type BlocklistType string

// Blocklist restriction types.
const (
	BlocklistTypeBlackHash BlocklistType = "black_hash"
)

// BlocklistStatus is the recommendation status returned by BlocklistValidate.
type BlocklistStatus string

// Blocklist validation statuses.
const (
	BlocklistStatusNotRecommended       BlocklistStatus = "Not recommended"
	BlocklistStatusNotAllowed           BlocklistStatus = "Not allowed"
	BlocklistStatusNone                 BlocklistStatus = "NONE"
	BlocklistStatusDuplicatedSHA1       BlocklistStatus = "duplicated_value_sha1"
	BlocklistStatusDuplicatedSHA256     BlocklistStatus = "duplicated_value_sha256"
	BlocklistStatusDuplicatedSHA1SHA256 BlocklistStatus = "duplicated_value_sha1_sha256"
	BlocklistStatusDuplication          BlocklistStatus = "Duplication"
)

// BlocklistItem is a SentinelOne blocklist (restrictions) entry: a SHA1 and/or
// SHA256 hash that agents block from executing.
type BlocklistItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	SHA256Value string `json:"sha256Value"`
	OSType      string `json:"osType"`
	Source      string `json:"source"`
	Description string `json:"description"`
	ScopeName   string `json:"scopeName"`
	ScopePath   string `json:"scopePath"`
	Imported    bool   `json:"imported"`
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (b *BlocklistItem) UnmarshalJSON(data []byte) error {
	type alias BlocklistItem
	if err := json.Unmarshal(data, (*alias)(b)); err != nil {
		return err
	}
	b.Raw = append(b.Raw[:0:0], data...)
	return nil
}

// BlocklistListParams are query parameters for listing blocklist items.
type BlocklistListParams struct {
	SiteIDs    []string
	GroupIDs   []string
	AccountIDs []string
	IDs        []string
	OSTypes    []string
	Types      []string
	Sources    []string
	Query      string
	Value      string
	Tenant     *bool
	Limit      int
	Cursor     string
	SortBy     string
	SortOrder  string
}

func (p *BlocklistListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "ids", p.IDs)
	addCSV(v, "osTypes", p.OSTypes)
	addCSV(v, "types", p.Types)
	addCSV(v, "source", p.Sources)
	addString(v, "query", p.Query)
	addString(v, "value", p.Value)
	addBool(v, "tenant", p.Tenant)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// BlocklistScope identifies the scope a write applies to. Set Tenant for the
// global blocklist, or one or more of the ID slices for account/site/group
// scopes.
type BlocklistScope struct {
	AccountIDs []string
	SiteIDs    []string
	GroupIDs   []string
	Tenant     bool
}

type blocklistFilter struct {
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	Tenant     bool     `json:"tenant,omitempty"`
}

func (s BlocklistScope) filter() blocklistFilter {
	return blocklistFilter(s)
}

// BlocklistCreate is the data payload for creating or updating a blocklist item.
// A blocklist item needs at least one of Value (SHA1) or SHA256Value.
type BlocklistCreate struct {
	Type        BlocklistType   `json:"type"`
	OSType      BlocklistOSType `json:"osType"`
	Value       string          `json:"value,omitempty"`
	SHA256Value string          `json:"sha256Value,omitempty"`
	Description string          `json:"description,omitempty"`
	Source      string          `json:"source,omitempty"`
}

type blocklistWriteRequest struct {
	Filter blocklistFilter `json:"filter"`
	Data   any             `json:"data"`
}

type blocklistUpdateData struct {
	ID string `json:"id"`
	BlocklistCreate
}

// BlocklistValidateInput is the data payload for BlocklistValidate.
type BlocklistValidateInput struct {
	OSType      BlocklistOSType `json:"osType,omitempty"`
	Value       string          `json:"value,omitempty"`
	SHA256Value string          `json:"sha256Value,omitempty"`
}

// BlocklistValidationDetail is one field-level validation error.
type BlocklistValidationDetail struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// BlocklistValidation is the recommendation returned by BlocklistValidate.
type BlocklistValidation struct {
	Status  BlocklistStatus             `json:"status"`
	Details []BlocklistValidationDetail `json:"details"`

	Raw json.RawMessage `json:"-"`
}

func (b *BlocklistValidation) UnmarshalJSON(data []byte) error {
	type alias BlocklistValidation
	if err := json.Unmarshal(data, (*alias)(b)); err != nil {
		return err
	}
	b.Raw = append(b.Raw[:0:0], data...)
	return nil
}

// BlocklistList returns a paginated list of blocklist items.
func (c *Client) BlocklistList(ctx context.Context, params *BlocklistListParams) ([]BlocklistItem, *Pagination, error) {
	return list[BlocklistItem](c, ctx, "/restrictions", params.values())
}

// BlocklistCreate adds a hash to the blocklist at the given scope.
func (c *Client) BlocklistCreate(ctx context.Context, scope BlocklistScope, data BlocklistCreate) (*BlocklistItem, error) {
	req := blocklistWriteRequest{Filter: scope.filter(), Data: data}
	var resp listResponse[BlocklistItem]
	if err := c.post(ctx, "/restrictions", req, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("mgmt: blocklist item not returned after create")
	}
	return &resp.Data[0], nil
}

// BlocklistUpdate changes the properties of a blocklist item identified by id.
func (c *Client) BlocklistUpdate(ctx context.Context, id string, scope BlocklistScope, data BlocklistCreate) (*BlocklistItem, error) {
	req := blocklistWriteRequest{
		Filter: scope.filter(),
		Data:   blocklistUpdateData{ID: id, BlocklistCreate: data},
	}
	var resp listResponse[BlocklistItem]
	if err := c.put(ctx, "/restrictions", req, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("mgmt: blocklist item not returned after update")
	}
	return &resp.Data[0], nil
}

// BlocklistDelete removes blocklist items by ID and returns the affected count.
func (c *Client) BlocklistDelete(ctx context.Context, ids []string) (int, error) {
	req := map[string]any{
		"data": map[string]any{"ids": ids},
	}
	var resp affectedResponse
	if err := c.jsonRequest(ctx, "DELETE", "/restrictions", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// BlocklistValidate checks whether a hash is on SentinelOne's "Not Allowed" or
// "Not Recommended" list before it is added to the blocklist.
func (c *Client) BlocklistValidate(ctx context.Context, scope BlocklistScope, data BlocklistValidateInput) (*BlocklistValidation, error) {
	req := struct {
		Filter blocklistFilter        `json:"filter"`
		Data   BlocklistValidateInput `json:"data"`
	}{scope.filter(), data}
	var resp singleResponse[BlocklistValidation]
	if err := c.post(ctx, "/restrictions/validate", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// BlocklistExport returns a CSV of blocklist items matching the filter.
func (c *Client) BlocklistExport(ctx context.Context, params *BlocklistListParams) ([]byte, error) {
	return c.getRaw(ctx, "/export/restrictions", params.values())
}
