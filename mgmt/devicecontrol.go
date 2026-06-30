package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// DeviceRule is a SentinelOne device control rule.
type DeviceRule struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	RuleName    string `json:"ruleName"`
	DeviceClass string `json:"deviceClass"`
	Action      string `json:"action"`
	Interface   string `json:"interface"`
	ScopeID     string `json:"scopeId"`
	ScopeName   string `json:"scopeName"`
	OSType      string `json:"osType"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (d *DeviceRule) UnmarshalJSON(b []byte) error {
	type alias DeviceRule
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// DeviceRuleListParams are query parameters for listing device rules.
type DeviceRuleListParams struct {
	SiteIDs    []string
	AccountIDs []string
	Query      string
	Limit      int
	Cursor     string
}

func (p *DeviceRuleListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// DeviceRulesList returns a paginated list of device control rules.
func (c *Client) DeviceRulesList(ctx context.Context, params *DeviceRuleListParams) ([]DeviceRule, *Pagination, error) {
	return list[DeviceRule](c, ctx, "/device-control", params.values())
}
