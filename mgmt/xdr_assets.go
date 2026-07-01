package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// XDRAssetCounts holds asset counts grouped by category and surface.
type XDRAssetCounts struct {
	Categories XDRCategoryDetails `json:"categories"`
	Surfaces   XDRSurfaceDetails  `json:"surfaces"`

	Raw json.RawMessage `json:"-"`
}

func (x *XDRAssetCounts) UnmarshalJSON(b []byte) error {
	type alias XDRAssetCounts
	if err := json.Unmarshal(b, (*alias)(x)); err != nil {
		return err
	}
	x.Raw = append(x.Raw[:0:0], b...)
	return nil
}

// XDRCategoryDetails holds per-category counts and subcategories.
type XDRCategoryDetails struct {
	Account                XDRCategoryCount `json:"account"`
	AiMl                   XDRCategoryCount `json:"aiMl"`
	ApplicationIntegration XDRCategoryCount `json:"applicationIntegration"`
	CloudApplication       XDRCategoryCount `json:"cloudApplication"`
	Code                   XDRCategoryCount `json:"code"`
	Container              XDRCategoryCount `json:"container"`
	DataAnalysis           XDRCategoryCount `json:"dataAnalysis"`
	DataStore              XDRCategoryCount `json:"dataStore"`
	DeveloperTool          XDRCategoryCount `json:"developerTool"`
	Device                 XDRCategoryCount `json:"device"`
	Function               XDRCategoryCount `json:"function"`
	Governance             XDRCategoryCount `json:"governance"`
	Identity               XDRCategoryCount `json:"identity"`
	Inventory              XDRCategoryCount `json:"inventory"`
	Network                XDRCategoryCount `json:"network"`
	Secrets                XDRCategoryCount `json:"secrets"`
	Server                 XDRCategoryCount `json:"server"`
	Storage                XDRCategoryCount `json:"storage"`
	Workstation            XDRCategoryCount `json:"workstation"`

	Raw json.RawMessage `json:"-"`
}

func (x *XDRCategoryDetails) UnmarshalJSON(b []byte) error {
	type alias XDRCategoryDetails
	if err := json.Unmarshal(b, (*alias)(x)); err != nil {
		return err
	}
	x.Raw = append(x.Raw[:0:0], b...)
	return nil
}

// XDRCategoryCount is a single category's count.
type XDRCategoryCount struct {
	Count int `json:"count"`

	Raw json.RawMessage `json:"-"`
}

func (x *XDRCategoryCount) UnmarshalJSON(b []byte) error {
	type alias XDRCategoryCount
	if err := json.Unmarshal(b, (*alias)(x)); err != nil {
		return err
	}
	x.Raw = append(x.Raw[:0:0], b...)
	return nil
}

// XDRSurfaceDetails holds per-surface counts.
type XDRSurfaceDetails struct {
	Cloud            XDRSurfaceCount `json:"cloud"`
	Endpoint         XDRSurfaceCount `json:"endpoint"`
	Identity         XDRSurfaceCount `json:"identity"`
	Network          XDRSurfaceCount `json:"network"`
	NetworkDiscovery XDRSurfaceCount `json:"networkDiscovery"`

	Raw json.RawMessage `json:"-"`
}

func (x *XDRSurfaceDetails) UnmarshalJSON(b []byte) error {
	type alias XDRSurfaceDetails
	if err := json.Unmarshal(b, (*alias)(x)); err != nil {
		return err
	}
	x.Raw = append(x.Raw[:0:0], b...)
	return nil
}

// XDRSurfaceCount is a single surface's count.
type XDRSurfaceCount struct {
	Count int `json:"count"`

	Raw json.RawMessage `json:"-"`
}

func (x *XDRSurfaceCount) UnmarshalJSON(b []byte) error {
	type alias XDRSurfaceCount
	if err := json.Unmarshal(b, (*alias)(x)); err != nil {
		return err
	}
	x.Raw = append(x.Raw[:0:0], b...)
	return nil
}

// XDRAssetCountsParams are query parameters for asset counts.
type XDRAssetCountsParams struct {
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
}

func (p *XDRAssetCountsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	return v
}

// XDRAssetCounts returns asset counts grouped by category and surface.
func (c *Client) XDRAssetCounts(ctx context.Context, params *XDRAssetCountsParams) (*XDRAssetCounts, error) {
	var resp singleResponse[XDRAssetCounts]
	if err := c.get(ctx, "/xdr/assets/asset-counts", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// XDRAssetCategory holds simple per-category counts (flat).
type XDRAssetCategory struct {
	Account     int `json:"account"`
	Container   int `json:"container"`
	Device      int `json:"device"`
	Identity    int `json:"identity"`
	Inventory   int `json:"inventory"`
	Server      int `json:"server"`
	Storage     int `json:"storage"`
	Workstation int `json:"workstation"`

	Raw json.RawMessage `json:"-"`
}

func (x *XDRAssetCategory) UnmarshalJSON(b []byte) error {
	type alias XDRAssetCategory
	if err := json.Unmarshal(b, (*alias)(x)); err != nil {
		return err
	}
	x.Raw = append(x.Raw[:0:0], b...)
	return nil
}

// XDRAssetCategories returns asset categories with counts.
func (c *Client) XDRAssetCategories(ctx context.Context, params *XDRAssetCountsParams) (*XDRAssetCategory, error) {
	var resp singleResponse[XDRAssetCategory]
	if err := c.get(ctx, "/xdr/assets/categories", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
