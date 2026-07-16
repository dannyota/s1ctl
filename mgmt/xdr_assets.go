package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// AssetType identifies an XDR asset category or surface subtype.
// The string value is the URL slug used in API paths.
type AssetType string

const (
	AssetTypeAccount                AssetType = "account"
	AssetTypeAiMl                   AssetType = "ai-ml"
	AssetTypeApplicationIntegration AssetType = "application-integration"
	AssetTypeCloudApplication       AssetType = "cloud-application"
	AssetTypeContainer              AssetType = "container"
	AssetTypeDataAnalysis           AssetType = "data-analysis"
	AssetTypeDataStore              AssetType = "data-store"
	AssetTypeDeveloperTool          AssetType = "developer-tool"
	AssetTypeDevice                 AssetType = "device"
	AssetTypeFunction               AssetType = "function"
	AssetTypeGovernance             AssetType = "governance"
	AssetTypeIdentity               AssetType = "identity"
	AssetTypeNetwork                AssetType = "network"
	AssetTypeServer                 AssetType = "server"
	AssetTypeStorage                AssetType = "storage"
	AssetTypeWorkstation            AssetType = "workstation"

	AssetTypeSurfaceCloud            AssetType = "surface/cloud"
	AssetTypeSurfaceEndpoint         AssetType = "surface/endpoint"
	AssetTypeSurfaceIdentity         AssetType = "surface/identity"
	AssetTypeSurfaceNetworkDiscovery AssetType = "surface/networkDiscovery"
)

// assetPath builds the API path for an asset type. When t is empty the
// base /xdr/assets path is returned (cross-type operations).
func assetPath(t AssetType, suffix string) string {
	base := "/xdr/assets"
	if t != "" {
		base += "/" + string(t)
	}
	if suffix != "" {
		base += "/" + suffix
	}
	return base
}

// XDRAssetListParams holds common query parameters for asset listing endpoints.
type XDRAssetListParams struct {
	Limit      int
	Skip       int
	Cursor     string
	SortBy     string
	SortOrder  string
	CountOnly  *bool
	SkipCount  *bool
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
	Extra      url.Values // type-specific filter passthrough
}

func (p *XDRAssetListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addInt(v, "limit", p.Limit)
	addInt(v, "skip", p.Skip)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addBool(v, "countOnly", p.CountOnly)
	addBool(v, "skipCount", p.SkipCount)
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	for k, vals := range p.Extra {
		for _, val := range vals {
			v.Add(k, val)
		}
	}
	return v
}

// XDRAssetList returns assets of the given type. When assetType is empty
// it lists across all types.
func (c *Client) XDRAssetList(ctx context.Context, assetType AssetType, params *XDRAssetListParams) ([]json.RawMessage, *Pagination, error) {
	var resp listResponse[json.RawMessage]
	if err := c.get(ctx, assetPath(assetType, ""), params.values(), &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// XDRAssetExport returns the raw export response (CSV or JSON) for the
// given asset type.
func (c *Client) XDRAssetExport(ctx context.Context, assetType AssetType, params *XDRAssetListParams) ([]byte, error) {
	return c.getRaw(ctx, assetPath(assetType, "export"), params.values())
}

// XDRAssetFilterAutocomplete returns filter autocomplete suggestions.
func (c *Client) XDRAssetFilterAutocomplete(ctx context.Context, assetType AssetType, params *XDRAssetListParams) ([]json.RawMessage, error) {
	var resp listResponse[json.RawMessage]
	if err := c.get(ctx, assetPath(assetType, "filters/autocomplete"), params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// XDRAssetFilterCount returns per-filter-value counts.
func (c *Client) XDRAssetFilterCount(ctx context.Context, assetType AssetType, params *XDRAssetListParams) (json.RawMessage, error) {
	var resp singleResponse[json.RawMessage]
	if err := c.get(ctx, assetPath(assetType, "filters/count"), params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// XDRAssetFilterFreeText returns free-text filter fields.
func (c *Client) XDRAssetFilterFreeText(ctx context.Context, assetType AssetType, params *XDRAssetListParams) ([]json.RawMessage, error) {
	var resp listResponse[json.RawMessage]
	if err := c.get(ctx, assetPath(assetType, "filters/free-text"), params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// XDRAssetActionInput is the request body for asset actions.
type XDRAssetActionInput struct {
	ActionName string   `json:"actionName"`
	IDIn       []string `json:"id__in,omitempty"`
	IDNin      []string `json:"id__nin,omitempty"`
}

// XDRAssetAction performs an action on assets of the given type and returns
// the affected count.
func (c *Client) XDRAssetAction(ctx context.Context, assetType AssetType, body *XDRAssetActionInput) (int, error) {
	if body == nil || body.ActionName == "" {
		return 0, fmt.Errorf("mgmt: action name is required")
	}
	var resp affectedResponse
	if err := c.post(ctx, assetPath(assetType, "action"), body, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// XDRAssetNoteInput is the request body for asset note operations.
type XDRAssetNoteInput struct {
	ID         string `json:"id,omitempty"`
	ResourceID string `json:"resourceId,omitempty"`
	Note       string `json:"note,omitempty"`
}

// XDRAssetNoteCreate creates or updates a note on an asset.
func (c *Client) XDRAssetNoteCreate(ctx context.Context, input *XDRAssetNoteInput) error {
	return c.post(ctx, "/xdr/assets/notes", input, nil)
}

// XDRAssetNoteDelete deletes a note from an asset.
func (c *Client) XDRAssetNoteDelete(ctx context.Context, input *XDRAssetNoteInput) error {
	return c.jsonRequest(ctx, "DELETE", "/xdr/assets/notes", input, nil)
}

// XDRAssetTags returns asset tags.
func (c *Client) XDRAssetTags(ctx context.Context, params *XDRAssetListParams) ([]json.RawMessage, error) {
	var resp listResponse[json.RawMessage]
	if err := c.get(ctx, "/xdr/assets/tags", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// XDRAssetTagCount returns tag counts.
func (c *Client) XDRAssetTagCount(ctx context.Context, body any) (json.RawMessage, error) {
	var resp singleResponse[json.RawMessage]
	if err := c.post(ctx, "/xdr/assets/tags/count", body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// XDRAssetSubCategories returns sub-category information.
func (c *Client) XDRAssetSubCategories(ctx context.Context, params *XDRAssetCountsParams) (json.RawMessage, error) {
	var resp singleResponse[json.RawMessage]
	if err := c.get(ctx, "/xdr/assets/sub-categories", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

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
