package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// DeviceEvent is a device control event from an endpoint.
type DeviceEvent struct {
	ID                   string `json:"id"`
	EventID              string `json:"eventId"`
	Interface            string `json:"interface"`
	DeviceClass          string `json:"deviceClass"`
	ServiceClass         string `json:"serviceClass"`
	RuleID               string `json:"ruleId"`
	VendorID             string `json:"vendorId"`
	ProductID            string `json:"productId"`
	EventTime            string `json:"eventTime"`
	EventType            string `json:"eventType"`
	DeviceName           string `json:"deviceName"`
	UID                  string `json:"uId"`
	AgentID              string `json:"agentId"`
	MinorClass           string `json:"minorClass"`
	ProfileUUIDs         string `json:"profileUuids"`
	LMPVersion           string `json:"lmpVersion"`
	AccessPermission     string `json:"accessPermission"`
	ComputerName         string `json:"computerName"`
	LastLoggedInUserName string `json:"lastLoggedInUserName"`
	DeviceID             string `json:"deviceId"`
	CreatedAt            string `json:"createdAt"`
	UpdatedAt            string `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (d *DeviceEvent) UnmarshalJSON(b []byte) error {
	type alias DeviceEvent
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// DeviceEventListParams are query parameters for listing device control events.
type DeviceEventListParams struct {
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
	Query      string
	Interfaces []string
	Limit      int
	Cursor     string
}

func (p *DeviceEventListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addString(v, "query", p.Query)
	addCSV(v, "interfaces", p.Interfaces)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// DeviceEventsList returns a paginated list of device control events.
func (c *Client) DeviceEventsList(ctx context.Context, params *DeviceEventListParams) ([]DeviceEvent, *Pagination, error) {
	return list[DeviceEvent](c, ctx, "/device-control/events", params.values())
}

// DeviceRulesDelete deletes device control rules by ID.
func (c *Client) DeviceRulesDelete(ctx context.Context, ids []string) (int, error) {
	req := map[string]any{
		"filter": map[string]any{
			"ids": ids,
		},
	}
	var resp affectedResponse
	if err := c.jsonRequest(ctx, http.MethodDelete, "/device-control", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// RuleOrder specifies a rule's desired position.
type RuleOrder struct {
	ID    string `json:"id"`
	Order int    `json:"order"`

	Raw json.RawMessage `json:"-"`
}

func (r *RuleOrder) UnmarshalJSON(b []byte) error {
	type alias RuleOrder
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// DeviceRuleReorderFilter scopes a reorder operation.
type DeviceRuleReorderFilter struct {
	AccountIDs []string             `json:"accountIds,omitempty"`
	SiteIDs    []string             `json:"siteIds,omitempty"`
	GroupIDs   []string             `json:"groupIds,omitempty"`
	Tenant     *bool                `json:"tenant,omitempty"`
	Interface  *DeviceRuleInterface `json:"interface,omitempty"`
}

// DeviceRulesReorder changes the order of device control rules within a scope.
func (c *Client) DeviceRulesReorder(ctx context.Context, orders []RuleOrder, filter DeviceRuleReorderFilter) error {
	req := struct {
		Data   []RuleOrder             `json:"data"`
		Filter DeviceRuleReorderFilter `json:"filter"`
	}{Data: orders, Filter: filter}
	var resp struct {
		Data struct {
			Success bool `json:"success"`
		} `json:"data"`
	}
	return c.put(ctx, "/device-control/reorder", req, &resp)
}

// DeviceRuleCopyTarget specifies a destination scope for copying rules.
type DeviceRuleCopyTarget struct {
	AccountID *string  `json:"accountId,omitempty"`
	SiteID    *string  `json:"siteId,omitempty"`
	GroupIDs  []string `json:"groupIds,omitempty"`
}

// DeviceRulesCopy copies device control rules from a source scope to targets.
func (c *Client) DeviceRulesCopy(ctx context.Context, filter DeviceRuleScopeFilter, targets []DeviceRuleCopyTarget) (int, error) {
	req := struct {
		Filter DeviceRuleScopeFilter  `json:"filter"`
		Data   []DeviceRuleCopyTarget `json:"data"`
	}{Filter: filter, Data: targets}
	var resp affectedResponse
	if err := c.post(ctx, "/device-control/copy-rules", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// DeviceRulesSetStatus enables or disables device control rules by ID.
func (c *Client) DeviceRulesSetStatus(ctx context.Context, ids []string, status DeviceRuleStatus) (int, error) {
	req := struct {
		Filter struct {
			IDs []string `json:"ids"`
		} `json:"filter"`
		Data struct {
			Status DeviceRuleStatus `json:"status"`
		} `json:"data"`
	}{}
	req.Filter.IDs = ids
	req.Data.Status = status
	var resp affectedResponse
	if err := c.put(ctx, "/device-control/enable", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}
