package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
)

// Policy is a SentinelOne endpoint policy (at site, group, or account scope).
type Policy struct {
	Raw json.RawMessage `json:"-"`
}

func (p *Policy) UnmarshalJSON(b []byte) error {
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

type policyResponse struct {
	Data json.RawMessage `json:"data"`
}

// PolicyGet returns the policy for a site.
func (c *Client) PolicyGetSite(ctx context.Context, siteID string) (*Policy, error) {
	var resp policyResponse
	if err := c.get(ctx, fmt.Sprintf("/sites/%s/policy", siteID), nil, &resp); err != nil {
		return nil, err
	}
	p := &Policy{}
	p.Raw = append(p.Raw[:0:0], resp.Data...)
	return p, nil
}

// PolicyGetAccount returns the policy for an account.
func (c *Client) PolicyGetAccount(ctx context.Context, accountID string) (*Policy, error) {
	var resp policyResponse
	if err := c.get(ctx, fmt.Sprintf("/accounts/%s/policy", accountID), nil, &resp); err != nil {
		return nil, err
	}
	p := &Policy{}
	p.Raw = append(p.Raw[:0:0], resp.Data...)
	return p, nil
}

// PolicyGetGroup returns the policy for a group.
func (c *Client) PolicyGetGroup(ctx context.Context, siteID, groupID string) (*Policy, error) {
	var resp policyResponse
	path := fmt.Sprintf("/sites/%s/groups/%s/policy", siteID, groupID)
	if err := c.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	p := &Policy{}
	p.Raw = append(p.Raw[:0:0], resp.Data...)
	return p, nil
}

// PolicyUpdateSite updates the policy for a site.
func (c *Client) PolicyUpdateSite(ctx context.Context, siteID string, policy json.RawMessage) (*Policy, error) {
	req := map[string]any{"data": policy}
	var resp policyResponse
	if err := c.put(ctx, fmt.Sprintf("/sites/%s/policy", siteID), req, &resp); err != nil {
		return nil, err
	}
	p := &Policy{}
	p.Raw = append(p.Raw[:0:0], resp.Data...)
	return p, nil
}

// PolicyUpdateAccount updates the policy for an account.
func (c *Client) PolicyUpdateAccount(ctx context.Context, accountID string, policy json.RawMessage) (*Policy, error) {
	req := map[string]any{"data": policy}
	var resp policyResponse
	if err := c.put(ctx, fmt.Sprintf("/accounts/%s/policy", accountID), req, &resp); err != nil {
		return nil, err
	}
	p := &Policy{}
	p.Raw = append(p.Raw[:0:0], resp.Data...)
	return p, nil
}

// PolicyUpdateGroup updates the policy for a group.
func (c *Client) PolicyUpdateGroup(ctx context.Context, siteID, groupID string, policy json.RawMessage) (*Policy, error) {
	req := map[string]any{"data": policy}
	var resp policyResponse
	path := fmt.Sprintf("/sites/%s/groups/%s/policy", siteID, groupID)
	if err := c.put(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	p := &Policy{}
	p.Raw = append(p.Raw[:0:0], resp.Data...)
	return p, nil
}
