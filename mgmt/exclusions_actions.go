package mgmt

import (
	"context"
	"fmt"
)

// ExclusionCreate is the request body for creating an exclusion.
type ExclusionCreate struct {
	Type              string   `json:"type"`
	Value             string   `json:"value"`
	OSType            string   `json:"osType"`
	Mode              string   `json:"mode,omitempty"`
	Description       string   `json:"description,omitempty"`
	PathExclusionType string   `json:"pathExclusionType,omitempty"`
	GroupIDs          []string `json:"groupIds,omitempty"`
	SiteIDs           []string `json:"siteIds,omitempty"`
}

type exclusionCreateRequest struct {
	Filter struct {
		SiteIDs []string `json:"siteIds,omitempty"`
	} `json:"filter"`
	Data ExclusionCreate `json:"data"`
}

// ExclusionsCreate creates an exclusion.
func (c *Client) ExclusionsCreate(ctx context.Context, siteIDs []string, excl ExclusionCreate) (*Exclusion, error) {
	req := exclusionCreateRequest{Data: excl}
	req.Filter.SiteIDs = siteIDs
	var resp listResponse[Exclusion]
	if err := c.post(ctx, "/exclusions", req, &resp); err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("mgmt: exclusion not returned after create")
	}
	return &resp.Data[0], nil
}

// ExclusionsUpdate updates an exclusion.
func (c *Client) ExclusionsUpdate(ctx context.Context, id string, data ExclusionCreate) (*Exclusion, error) {
	req := map[string]any{"data": data}
	var resp singleResponse[Exclusion]
	if err := c.put(ctx, fmt.Sprintf("/exclusions/%s", id), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ExclusionsDelete deletes exclusions by ID.
func (c *Client) ExclusionsDelete(ctx context.Context, ids []string) (int, error) {
	req := map[string]any{
		"data": map[string]any{},
		"filter": map[string]any{
			"ids": ids,
		},
	}
	var resp affectedResponse
	if err := c.post(ctx, "/exclusions/delete", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}
