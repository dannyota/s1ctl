package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// UpgradePolicy is a SentinelOne agent auto-upgrade policy.
type UpgradePolicy struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	OSType       string           `json:"osType"`
	ScopeLevel   string           `json:"scopeLevel"`
	ScopeID      string           `json:"scopeId"`
	IsActive     bool             `json:"isActive"`
	IsScheduled  bool             `json:"isScheduled"`
	AllEndpoints bool             `json:"allEndpoints"`
	MaxRetries   int              `json:"maxRetries"`
	Priority     int              `json:"priority"`
	Package      UpgradePolicyPkg `json:"package"`
	Tags         []string         `json:"tags"`
	ActivatedAt  string           `json:"activatedAt"`
	CreatedAt    string           `json:"createdAt"`
	UpdatedAt    string           `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (u *UpgradePolicy) UnmarshalJSON(b []byte) error {
	type alias UpgradePolicy
	if err := json.Unmarshal(b, (*alias)(u)); err != nil {
		return err
	}
	u.Raw = append(u.Raw[:0:0], b...)
	return nil
}

// UpgradePolicyPkg identifies the agent package targeted by an upgrade policy.
type UpgradePolicyPkg struct {
	Build  string `json:"build"`
	FileID string `json:"fileId"`
	Major  string `json:"major"`
	Minor  string `json:"minor"`
}

// UpgradePolicyListParams are query parameters for listing upgrade policies.
type UpgradePolicyListParams struct {
	ScopeLevel string // required: account, group, site, tenant
	ScopeID    string
	OSType     string // required: linux, macos, windows
	Limit      int    // required
	Skip       int
	SortBy     string // required (e.g. priority)
	SortOrder  string // required (asc, desc)
}

func (p *UpgradePolicyListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "scopeLevel", p.ScopeLevel)
	addString(v, "scopeId", p.ScopeID)
	addString(v, "osType", p.OSType)
	addInt(v, "limit", p.Limit)
	if p.Skip > 0 {
		v.Set("skip", strconv.Itoa(p.Skip))
	}
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// upgradePolicyListResponse is the envelope for the upgrade-policy/policies endpoint.
type upgradePolicyListResponse struct {
	Data struct {
		IsInherited     bool            `json:"isInherited"`
		Policies        []UpgradePolicy `json:"policies"`
		PoliciesInChild bool            `json:"policiesInChildScope"`
	} `json:"data"`
	Pagination struct {
		TotalItems int `json:"totalItems"`
	} `json:"pagination"`
}

// UpgradePoliciesList returns upgrade policies for a given scope and OS type.
func (c *Client) UpgradePoliciesList(ctx context.Context, params *UpgradePolicyListParams) ([]UpgradePolicy, int, error) {
	var resp upgradePolicyListResponse
	if err := c.get(ctx, "/upgrade-policy/policies", params.values(), &resp); err != nil {
		return nil, 0, err
	}
	return resp.Data.Policies, resp.Pagination.TotalItems, nil
}

// UpgradePolicyCreate is the request body for creating an upgrade policy.
type UpgradePolicyCreate struct {
	Name         string           `json:"name"`
	Description  string           `json:"description,omitempty"`
	OSType       string           `json:"osType"`
	ScopeLevel   string           `json:"scopeLevel"`
	ScopeID      string           `json:"scopeId,omitempty"`
	IsActive     bool             `json:"isActive"`
	IsScheduled  bool             `json:"isScheduled,omitempty"`
	AllEndpoints bool             `json:"allEndpoints"`
	MaxRetries   int              `json:"maxRetries,omitempty"`
	Package      UpgradePolicyPkg `json:"package"`
	Tags         []string         `json:"tags,omitempty"`
}

// UpgradePoliciesCreate creates an upgrade policy.
func (c *Client) UpgradePoliciesCreate(ctx context.Context, data UpgradePolicyCreate) error {
	return c.post(ctx, "/upgrade-policy/policy", data, nil)
}

// UpgradePoliciesUpdate updates an upgrade policy.
func (c *Client) UpgradePoliciesUpdate(ctx context.Context, id string, data UpgradePolicyCreate) error {
	return c.put(ctx, fmt.Sprintf("/upgrade-policy/policy/%s", id), data, nil)
}

// UpgradePoliciesDelete deletes an upgrade policy by ID.
func (c *Client) UpgradePoliciesDelete(ctx context.Context, id string) error {
	req := map[string]string{"action": "delete"}
	return c.post(ctx, fmt.Sprintf("/upgrade-policy/policy/%s", id), req, nil)
}

// UpgradePoliciesActivate activates an upgrade policy.
func (c *Client) UpgradePoliciesActivate(ctx context.Context, id string) error {
	req := map[string]string{"action": "activate"}
	return c.post(ctx, fmt.Sprintf("/upgrade-policy/policy/%s", id), req, nil)
}

// UpgradePoliciesDeactivate deactivates an upgrade policy.
func (c *Client) UpgradePoliciesDeactivate(ctx context.Context, id string) error {
	req := map[string]string{"action": "deactivate"}
	return c.post(ctx, fmt.Sprintf("/upgrade-policy/policy/%s", id), req, nil)
}

// UpgradePackage is an available agent package for upgrade policies.
type UpgradePackage struct {
	Build       string               `json:"build"`
	Major       string               `json:"major"`
	Minor       string               `json:"minor"`
	DisplayName string               `json:"displayName"`
	FileNames   []UpgradePackageFile `json:"fileNames"`

	Raw json.RawMessage `json:"-"`
}

func (u *UpgradePackage) UnmarshalJSON(b []byte) error {
	type alias UpgradePackage
	if err := json.Unmarshal(b, (*alias)(u)); err != nil {
		return err
	}
	u.Raw = append(u.Raw[:0:0], b...)
	return nil
}

// UpgradePackageFile identifies a downloadable package file.
type UpgradePackageFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UpgradePackageListParams are query parameters for listing available upgrade packages.
type UpgradePackageListParams struct {
	ScopeLevel          string // required: account, group, site, tenant
	ScopeID             string
	OSType              string // required: linux, macos, windows
	DisplayNameContains string
}

func (p *UpgradePackageListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "scopeLevel", p.ScopeLevel)
	addString(v, "scopeId", p.ScopeID)
	addString(v, "osType", p.OSType)
	addString(v, "displayName__contains", p.DisplayNameContains)
	return v
}

type upgradePackageListResponse struct {
	Data struct {
		Packages []UpgradePackage `json:"packages"`
	} `json:"data"`
}

// UpgradePackagesList returns available packages for upgrade policies.
func (c *Client) UpgradePackagesList(ctx context.Context, params *UpgradePackageListParams) ([]UpgradePackage, error) {
	var resp upgradePackageListResponse
	if err := c.get(ctx, "/upgrade-policy/available-packages", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data.Packages, nil
}
