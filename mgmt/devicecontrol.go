package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// DeviceRuleAction is the action of a device control rule.
type DeviceRuleAction string

const (
	DeviceRuleActionAllow DeviceRuleAction = "Allow"
	DeviceRuleActionBlock DeviceRuleAction = "Block"
)

// DeviceRuleStatus is the status of a device control rule.
type DeviceRuleStatus string

const (
	DeviceRuleStatusEnabled  DeviceRuleStatus = "Enabled"
	DeviceRuleStatusDisabled DeviceRuleStatus = "Disabled"
)

// DeviceRuleInterface is the physical bus type of a device.
type DeviceRuleInterface string

const (
	DeviceRuleInterfaceUSB         DeviceRuleInterface = "USB"
	DeviceRuleInterfaceBluetooth   DeviceRuleInterface = "Bluetooth"
	DeviceRuleInterfaceThunderbolt DeviceRuleInterface = "Thunderbolt"
	DeviceRuleInterfaceSDCard      DeviceRuleInterface = "SDCard"
)

// DeviceRuleType is the rule type that determines which fields are required.
type DeviceRuleType string

const (
	DeviceRuleTypeClass            DeviceRuleType = "class"
	DeviceRuleTypeProductID        DeviceRuleType = "productId"
	DeviceRuleTypeVendorID         DeviceRuleType = "vendorId"
	DeviceRuleTypeDeviceID         DeviceRuleType = "deviceId"
	DeviceRuleTypeUID              DeviceRuleType = "uid"
	DeviceRuleTypeHWIdentifiers    DeviceRuleType = "hwIdentifiers"
	DeviceRuleTypeBluetoothVersion DeviceRuleType = "bluetoothVersion"
	DeviceRuleTypeSDCard           DeviceRuleType = "sdCard"
)

// DeviceRuleAccessPermission is the access permission for a device rule.
type DeviceRuleAccessPermission string

const (
	DeviceRuleAccessReadOnly      DeviceRuleAccessPermission = "Read-Only"
	DeviceRuleAccessReadWrite     DeviceRuleAccessPermission = "Read-Write"
	DeviceRuleAccessNotApplicable DeviceRuleAccessPermission = "Not-Applicable"
)

// DeviceRuleScope is the scope level of a device control rule.
type DeviceRuleScope string

const (
	DeviceRuleScopeGlobal  DeviceRuleScope = "global"
	DeviceRuleScopeAccount DeviceRuleScope = "account"
	DeviceRuleScopeSite    DeviceRuleScope = "site"
	DeviceRuleScopeGroup   DeviceRuleScope = "group"
)

// DeviceRule is a SentinelOne device control rule.
type DeviceRule struct {
	ID               string                     `json:"id"`
	RuleName         string                     `json:"ruleName"`
	Status           DeviceRuleStatus           `json:"status"`
	Action           DeviceRuleAction           `json:"action"`
	Interface        DeviceRuleInterface        `json:"interface"`
	RuleType         DeviceRuleType             `json:"ruleType"`
	AccessPermission DeviceRuleAccessPermission `json:"accessPermission"`
	DeviceClass      string                     `json:"deviceClass"`
	DeviceID         string                     `json:"deviceId"`
	VendorID         string                     `json:"vendorId"`
	ProductID        string                     `json:"productId"`
	UID              string                     `json:"uid"`
	Version          string                     `json:"version"`
	Order            int                        `json:"order"`
	Scope            DeviceRuleScope            `json:"scope"`
	ScopeID          string                     `json:"scopeId"`
	ScopeName        string                     `json:"scopeName"`
	OSType           string                     `json:"osType"`
	MinorClasses     []string                   `json:"minorClasses"`
	BluetoothAddress string                     `json:"bluetoothAddress"`
	GattService      []string                   `json:"gattService"`
	ManufacturerName string                     `json:"manufacturerName"`
	DeviceName       string                     `json:"deviceName"`
	CreatedAt        string                     `json:"createdAt"`
	UpdatedAt        string                     `json:"updatedAt"`

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

// DeviceRulesGet returns a single device control rule by ID.
func (c *Client) DeviceRulesGet(ctx context.Context, id string) (*DeviceRule, error) {
	return getByID[DeviceRule](c, ctx, "/device-control", "device rule", id)
}

// DeviceRuleCreate is the request body for creating a device control rule.
type DeviceRuleCreate struct {
	RuleName         string                     `json:"ruleName"`
	Interface        DeviceRuleInterface        `json:"interface"`
	RuleType         DeviceRuleType             `json:"ruleType"`
	Action           DeviceRuleAction           `json:"action"`
	Status           DeviceRuleStatus           `json:"status"`
	AccessPermission DeviceRuleAccessPermission `json:"accessPermission"`
	DeviceClass      string                     `json:"deviceClass,omitempty"`
	DeviceID         string                     `json:"deviceId,omitempty"`
	VendorID         string                     `json:"vendorId,omitempty"`
	ProductID        string                     `json:"productId,omitempty"`
	UID              string                     `json:"uid,omitempty"`
	Version          string                     `json:"version,omitempty"`
	MinorClasses     []string                   `json:"minorClasses,omitempty"`
	BluetoothAddress string                     `json:"bluetoothAddress,omitempty"`
	GattService      []string                   `json:"gattService,omitempty"`
	ManufacturerName string                     `json:"manufacturerName,omitempty"`
	DeviceName       string                     `json:"deviceName,omitempty"`
}

// deviceRuleCreateRequest is the POST body for creating a device control rule.
// The SentinelOne API requires both data and filter in the body.
type deviceRuleCreateRequest struct {
	Data   DeviceRuleCreate      `json:"data"`
	Filter DeviceRuleScopeFilter `json:"filter"`
}

// DeviceRuleScopeFilter sets the scope for a new device control rule.
type DeviceRuleScopeFilter struct {
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	Tenant     *bool    `json:"tenant,omitempty"`
}

// DeviceRulesCreate creates a device control rule at the specified scope.
func (c *Client) DeviceRulesCreate(ctx context.Context, data DeviceRuleCreate, filter DeviceRuleScopeFilter) (*DeviceRule, error) {
	req := deviceRuleCreateRequest{Data: data, Filter: filter}
	var resp singleResponse[DeviceRule]
	if err := c.post(ctx, "/device-control", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeviceRulesUpdate updates a device control rule by ID.
func (c *Client) DeviceRulesUpdate(ctx context.Context, id string, data DeviceRuleCreate) (*DeviceRule, error) {
	return update[DeviceRule](c, ctx, fmt.Sprintf("/device-control/%s", url.PathEscape(id)), data)
}
