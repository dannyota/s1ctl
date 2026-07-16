package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// DeployTargetOS is the operating system for a Sentinel Deploy credential group.
type DeployTargetOS string

// Deploy target OS types.
const (
	DeployTargetOSWindows  DeployTargetOS = "windows"
	DeployTargetOSOSXLinux DeployTargetOS = "osx_linux"
)

// DeployCredGroup is a Sentinel Deploy (Ranger) credential group used for
// deploying agents to unprotected endpoints.
type DeployCredGroup struct {
	ID              string         `json:"id"`
	GroupName       string         `json:"groupName"`
	GroupPassphrase string         `json:"groupPassphrase"`
	ScopeID         string         `json:"scopeId"`
	Domain          string         `json:"domain"`
	TargetOS        DeployTargetOS `json:"targetOs"`
	TotalDetails    int            `json:"totalDetails"`

	Raw json.RawMessage `json:"-"`
}

func (d *DeployCredGroup) UnmarshalJSON(data []byte) error {
	type alias DeployCredGroup
	if err := json.Unmarshal(data, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], data...)
	return nil
}

// DeployCredGroupListParams are query parameters for listing credential groups.
type DeployCredGroupListParams struct {
	SiteIDs       []string
	AccountIDs    []string
	IDs           []string
	GroupName     string
	GroupNameLike string
	TargetOS      string
	Limit         int
	Cursor        string
	SortBy        string
	SortOrder     string
}

func (p *DeployCredGroupListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "ids", p.IDs)
	addString(v, "groupName", p.GroupName)
	addString(v, "groupNameLike", p.GroupNameLike)
	addString(v, "targetOs", p.TargetOS)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// DeployCredGroupCreateInput is the data payload for creating a credential group.
type DeployCredGroupCreateInput struct {
	GroupName       string          `json:"groupName"`
	GroupPassphrase string          `json:"groupPassphrase"`
	ScopeID         string          `json:"scopeId"`
	Domain          *string         `json:"domain,omitempty"`
	TargetOS        *DeployTargetOS `json:"targetOs,omitempty"`
}

// DeployCredDetail is a single credential entry within a credential group.
type DeployCredDetail struct {
	ID          string `json:"id"`
	CredGroupID string `json:"credGroupId"`
	Title       string `json:"title"`
	CredType    string `json:"credType"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	CreatedBy   string `json:"createdBy"`
	UpdatedBy   string `json:"updatedBy"`

	Raw json.RawMessage `json:"-"`
}

func (d *DeployCredDetail) UnmarshalJSON(data []byte) error {
	type alias DeployCredDetail
	if err := json.Unmarshal(data, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], data...)
	return nil
}

// DeployCredDetailListParams are query parameters for listing credential
// group details.
type DeployCredDetailListParams struct {
	SiteIDs      []string
	AccountIDs   []string
	IDs          []string
	CredGroupIDs []string
	Title        string
	TitleLike    string
	CredTypeLike string
	Limit        int
	Cursor       string
	SortBy       string
	SortOrder    string
}

func (p *DeployCredDetailListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "ids", p.IDs)
	addCSV(v, "credGroupIds", p.CredGroupIDs)
	addString(v, "title", p.Title)
	addString(v, "titleLike", p.TitleLike)
	addString(v, "credTypeLike", p.CredTypeLike)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// DeployCredDetailInput is a single credential detail for creation.
type DeployCredDetailInput struct {
	Title         string `json:"title"`
	CredType      string `json:"credType"`
	EncryptedKey  string `json:"encryptedKey"`
	EncryptedCred string `json:"encryptedCred"`
}

// DeployCredDetailAddInput is the data payload for adding credential details
// to a group.
type DeployCredDetailAddInput struct {
	CredGroupID string                  `json:"credGroupId"`
	Details     []DeployCredDetailInput `json:"details"`
}

// DeployCredGroupList returns a paginated list of credential groups.
func (c *Client) DeployCredGroupList(ctx context.Context, params *DeployCredGroupListParams) ([]DeployCredGroup, *Pagination, error) {
	return list[DeployCredGroup](c, ctx, "/ranger/cred-groups", params.values())
}

// DeployCredGroupCreate creates a new credential group.
func (c *Client) DeployCredGroupCreate(ctx context.Context, input DeployCredGroupCreateInput) (*DeployCredGroup, error) {
	return create[DeployCredGroup](c, ctx, "/ranger/cred-groups", input)
}

// DeployCredGroupDelete deletes a credential group by ID.
func (c *Client) DeployCredGroupDelete(ctx context.Context, id string) error {
	var resp struct {
		Data struct {
			Success bool `json:"success"`
		} `json:"data"`
	}
	if err := c.queryRequest(ctx, "DELETE", "/ranger/cred-groups/"+id, nil, &resp); err != nil {
		return err
	}
	if !resp.Data.Success {
		return fmt.Errorf("mgmt: deploy cred group delete returned success=false")
	}
	return nil
}

// DeployCredDetailList returns a paginated list of credential group details.
func (c *Client) DeployCredDetailList(ctx context.Context, params *DeployCredDetailListParams) ([]DeployCredDetail, *Pagination, error) {
	return list[DeployCredDetail](c, ctx, "/ranger/cred-groups/details", params.values())
}

// DeployCredDetailAdd adds credential details to a credential group.
func (c *Client) DeployCredDetailAdd(ctx context.Context, input DeployCredDetailAddInput) error {
	var resp struct {
		Data struct {
			Success bool `json:"success"`
		} `json:"data"`
	}
	req := map[string]any{"data": input}
	if err := c.post(ctx, "/ranger/cred-groups/details", req, &resp); err != nil {
		return err
	}
	return nil
}

// DeployCredDetailUpdate updates a credential detail by ID.
func (c *Client) DeployCredDetailUpdate(ctx context.Context, detailID string, input DeployCredDetailInput) (*DeployCredDetail, error) {
	return update[DeployCredDetail](c, ctx, "/ranger/cred-groups/details/"+detailID, input)
}

// DeployCredDetailDelete deletes a credential detail by ID.
func (c *Client) DeployCredDetailDelete(ctx context.Context, detailID string) error {
	var resp struct {
		Data struct {
			Success bool `json:"success"`
		} `json:"data"`
	}
	if err := c.queryRequest(ctx, "DELETE", "/ranger/cred-groups/details/"+detailID, nil, &resp); err != nil {
		return err
	}
	if !resp.Data.Success {
		return fmt.Errorf("mgmt: deploy cred detail delete returned success=false")
	}
	return nil
}
