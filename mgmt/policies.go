package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Policy is a SentinelOne endpoint policy (at site, group, or account scope).
type Policy struct {
	MitigationMode           string `json:"mitigationMode"`
	MitigationModeSuspicious string `json:"mitigationModeSuspicious"`
	AntiTamperingOn          bool   `json:"antiTamperingOn"`
	NetworkQuarantineOn      bool   `json:"networkQuarantineOn"`
	SnapshotsOn              bool   `json:"snapshotsOn"`
	Ioc                      bool   `json:"ioc"`
	InheritedFrom            string `json:"inheritedFrom"`
	AllowRemoteShell         bool   `json:"allowRemoteShell"`
	ScanNewAgents            bool   `json:"scanNewAgents"`
	AutoDecommissionOn       bool   `json:"autoDecommissionOn"`
	AutoDecommissionDays     int    `json:"autoDecommissionDays"`
	CreatedAt                string `json:"createdAt"`
	UpdatedAt                string `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (p *Policy) UnmarshalJSON(b []byte) error {
	type alias Policy
	if err := json.Unmarshal(b, (*alias)(p)); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

type policyResponse struct {
	Data json.RawMessage `json:"data"`
}

func (c *Client) getPolicy(ctx context.Context, path string) (*Policy, error) {
	var resp policyResponse
	if err := c.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	var p Policy
	if err := json.Unmarshal(resp.Data, &p); err != nil {
		return nil, fmt.Errorf("unmarshal policy: %w", err)
	}
	return &p, nil
}

func (c *Client) putPolicy(ctx context.Context, path string, policy json.RawMessage) (*Policy, error) {
	req := map[string]any{"data": policy}
	var resp policyResponse
	if err := c.put(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	var p Policy
	if err := json.Unmarshal(resp.Data, &p); err != nil {
		return nil, fmt.Errorf("unmarshal policy: %w", err)
	}
	return &p, nil
}

// PolicyGetSite returns the policy for a site.
func (c *Client) PolicyGetSite(ctx context.Context, siteID string) (*Policy, error) {
	return c.getPolicy(ctx, fmt.Sprintf("/sites/%s/policy", url.PathEscape(siteID)))
}

// PolicyGetAccount returns the policy for an account.
func (c *Client) PolicyGetAccount(ctx context.Context, accountID string) (*Policy, error) {
	return c.getPolicy(ctx, fmt.Sprintf("/accounts/%s/policy", url.PathEscape(accountID)))
}

// PolicyGetGroup returns the policy for a group.
func (c *Client) PolicyGetGroup(ctx context.Context, groupID string) (*Policy, error) {
	return c.getPolicy(ctx, fmt.Sprintf("/groups/%s/policy", url.PathEscape(groupID)))
}

// PolicyUpdateSite updates the policy for a site.
func (c *Client) PolicyUpdateSite(ctx context.Context, siteID string, policy json.RawMessage) (*Policy, error) {
	return c.putPolicy(ctx, fmt.Sprintf("/sites/%s/policy", url.PathEscape(siteID)), policy)
}

// PolicyUpdateAccount updates the policy for an account.
func (c *Client) PolicyUpdateAccount(ctx context.Context, accountID string, policy json.RawMessage) (*Policy, error) {
	return c.putPolicy(ctx, fmt.Sprintf("/accounts/%s/policy", url.PathEscape(accountID)), policy)
}

// PolicyUpdateGroup updates the policy for a group.
func (c *Client) PolicyUpdateGroup(ctx context.Context, groupID string, policy json.RawMessage) (*Policy, error) {
	return c.putPolicy(ctx, fmt.Sprintf("/groups/%s/policy", url.PathEscape(groupID)), policy)
}

// PolicyRevertSite reverts a site policy to its parent (account) inherited values.
func (c *Client) PolicyRevertSite(ctx context.Context, siteID string) error {
	return c.put(ctx, fmt.Sprintf("/sites/%s/revert-policy", url.PathEscape(siteID)), map[string]any{}, nil)
}

// PolicyRevertAccount reverts an account policy to the global inherited values.
func (c *Client) PolicyRevertAccount(ctx context.Context, accountID string) error {
	return c.put(ctx, fmt.Sprintf("/accounts/%s/revert-policy", url.PathEscape(accountID)), map[string]any{}, nil)
}

// PolicyRevertGroup reverts a group policy to its parent (site) inherited values.
func (c *Client) PolicyRevertGroup(ctx context.Context, groupID string) error {
	return c.put(ctx, fmt.Sprintf("/groups/%s/revert-policy", url.PathEscape(groupID)), map[string]any{}, nil)
}
