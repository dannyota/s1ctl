package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// FirewallRule is a SentinelOne firewall rule.
type FirewallRule struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Action    string `json:"action"`
	Direction string `json:"direction"`
	Protocol  string `json:"protocol"`
	OSType    string `json:"osType"`
	ScopeID   string `json:"scopeId"`
	ScopeName string `json:"scopeName"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (f *FirewallRule) UnmarshalJSON(b []byte) error {
	type alias FirewallRule
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FirewallRuleListParams are query parameters for listing firewall rules.
type FirewallRuleListParams struct {
	SiteIDs    []string
	AccountIDs []string
	Query      string
	Limit      int
	Cursor     string
}

func (p *FirewallRuleListParams) values() url.Values {
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

// FirewallRulesList returns a paginated list of firewall rules.
func (c *Client) FirewallRulesList(ctx context.Context, params *FirewallRuleListParams) ([]FirewallRule, *Pagination, error) {
	return list[FirewallRule](c, ctx, "/firewall-control", params.values())
}
