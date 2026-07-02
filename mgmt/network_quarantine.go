package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Network Quarantine shares the /firewall-control endpoint family with the
// standard firewall; every operation is the firewall operation addressed under
// the "network-quarantine" category segment. These wrappers delegate to the
// category-aware firewall internals, plus a handful of operations the firewall
// CLI does not yet expose (configuration, set-location, move-rules, tags).

// NetworkQuarantineList returns a paginated list of network quarantine rules.
func (c *Client) NetworkQuarantineList(ctx context.Context, params *FirewallRuleListParams) ([]FirewallRule, *Pagination, error) {
	return c.firewallRulesList(ctx, FirewallCategoryNetworkQuarantine, params)
}

// NetworkQuarantineGet returns a single network quarantine rule by ID.
func (c *Client) NetworkQuarantineGet(ctx context.Context, id string) (*FirewallRule, error) {
	return c.firewallRulesGet(ctx, FirewallCategoryNetworkQuarantine, id)
}

// NetworkQuarantineCreate creates a network quarantine rule.
func (c *Client) NetworkQuarantineCreate(ctx context.Context, scope FirewallRuleScope, data FirewallRuleCreate) (*FirewallRule, error) {
	return c.firewallRulesCreate(ctx, FirewallCategoryNetworkQuarantine, scope, data)
}

// NetworkQuarantineUpdate updates a network quarantine rule by ID.
func (c *Client) NetworkQuarantineUpdate(ctx context.Context, id string, data FirewallRuleCreate) (*FirewallRule, error) {
	return c.firewallRulesUpdate(ctx, FirewallCategoryNetworkQuarantine, id, data)
}

// NetworkQuarantineDelete deletes network quarantine rules by ID.
func (c *Client) NetworkQuarantineDelete(ctx context.Context, ids []string) (int, error) {
	return c.firewallRulesDelete(ctx, FirewallCategoryNetworkQuarantine, ids)
}

// NetworkQuarantineSetStatus enables or disables network quarantine rules by ID.
func (c *Client) NetworkQuarantineSetStatus(ctx context.Context, ids []string, status FirewallStatus) (int, error) {
	return c.firewallRulesSetStatus(ctx, FirewallCategoryNetworkQuarantine, ids, status)
}

// NetworkQuarantineReorder changes the order of network quarantine rules within a scope.
func (c *Client) NetworkQuarantineReorder(ctx context.Context, orders []RuleOrder, filter FirewallRuleReorderFilter) error {
	return c.firewallRulesReorder(ctx, FirewallCategoryNetworkQuarantine, orders, filter)
}

// NetworkQuarantineCopy copies network quarantine rules from a source scope to targets.
func (c *Client) NetworkQuarantineCopy(ctx context.Context, filter FirewallRuleReorderFilter, targets []FirewallRuleCopyTarget) (int, error) {
	return c.firewallRulesCopy(ctx, FirewallCategoryNetworkQuarantine, filter, targets)
}

// NetworkQuarantineProtocolsList returns the protocols available for network quarantine rules.
func (c *Client) NetworkQuarantineProtocolsList(ctx context.Context, params *FirewallProtocolListParams) ([]FirewallProtocol, *Pagination, error) {
	return c.firewallProtocolsList(ctx, FirewallCategoryNetworkQuarantine, params)
}

// NetworkQuarantineExport exports network quarantine rules as raw JSON for the given scope.
func (c *Client) NetworkQuarantineExport(ctx context.Context, params *FirewallRuleListParams) ([]byte, error) {
	return c.firewallRulesExport(ctx, FirewallCategoryNetworkQuarantine, params)
}

// NetworkQuarantineImport imports network quarantine rules from a JSON file into the given scope.
func (c *Client) NetworkQuarantineImport(ctx context.Context, scope FirewallImportScope, filename string, fileData []byte) error {
	return c.firewallRulesImport(ctx, FirewallCategoryNetworkQuarantine, scope, filename, fileData)
}

// --- Operations shared by both categories but only surfaced for NQ today. ---

// FirewallConfiguration is the firewall-control configuration for a scope.
type FirewallConfiguration struct {
	Enabled                 bool     `json:"enabled"`
	LocationAware           bool     `json:"locationAware"`
	ReportBlocked           bool     `json:"reportBlocked"`
	Inherits                bool     `json:"inherits"`
	InheritedFrom           string   `json:"inheritedFrom"`
	SelectedTags            []string `json:"selectedTags"`
	InheritSettings         bool     `json:"inheritSettings"`
	InheritAllFirewallRules bool     `json:"inheritAllFirewallRules"`

	Raw json.RawMessage `json:"-"`
}

func (f *FirewallConfiguration) UnmarshalJSON(b []byte) error {
	type alias FirewallConfiguration
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FirewallConfigurationUpdate is the mutable firewall-control configuration.
// Pointer fields are only sent when set, so callers patch individual toggles.
type FirewallConfigurationUpdate struct {
	Enabled                 *bool    `json:"enabled,omitempty"`
	LocationAware           *bool    `json:"locationAware,omitempty"`
	ReportBlocked           *bool    `json:"reportBlocked,omitempty"`
	Inherits                *bool    `json:"inherits,omitempty"`
	InheritedFrom           *string  `json:"inheritedFrom,omitempty"`
	SelectedTags            []string `json:"selectedTags,omitempty"`
	InheritSettings         *bool    `json:"inheritSettings,omitempty"`
	InheritAllFirewallRules *bool    `json:"inheritAllFirewallRules,omitempty"`
}

// FirewallConfigScope scopes a configuration read/write to an account, site,
// group, or the whole tenant.
type FirewallConfigScope struct {
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	Tenant     *bool    `json:"tenant,omitempty"`
}

func (s FirewallConfigScope) values() url.Values {
	v := url.Values{}
	addCSV(v, "accountIds", s.AccountIDs)
	addCSV(v, "siteIds", s.SiteIDs)
	addCSV(v, "groupIds", s.GroupIDs)
	addBool(v, "tenant", s.Tenant)
	return v
}

func (c *Client) firewallConfigurationGet(ctx context.Context, cat FirewallCategory, scope FirewallConfigScope) (*FirewallConfiguration, error) {
	var resp singleResponse[FirewallConfiguration]
	if err := c.get(ctx, firewallPath(cat, "/configuration"), scope.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) firewallConfigurationUpdate(ctx context.Context, cat FirewallCategory, scope FirewallConfigScope, data FirewallConfigurationUpdate) (*FirewallConfiguration, error) {
	req := struct {
		Filter FirewallConfigScope         `json:"filter"`
		Data   FirewallConfigurationUpdate `json:"data"`
	}{Filter: scope, Data: data}
	var resp singleResponse[FirewallConfiguration]
	if err := c.put(ctx, firewallPath(cat, "/configuration"), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// NetworkQuarantineConfigurationGet returns the network quarantine configuration for a scope.
func (c *Client) NetworkQuarantineConfigurationGet(ctx context.Context, scope FirewallConfigScope) (*FirewallConfiguration, error) {
	return c.firewallConfigurationGet(ctx, FirewallCategoryNetworkQuarantine, scope)
}

// NetworkQuarantineConfigurationUpdate updates the network quarantine configuration for a scope.
func (c *Client) NetworkQuarantineConfigurationUpdate(ctx context.Context, scope FirewallConfigScope, data FirewallConfigurationUpdate) (*FirewallConfiguration, error) {
	return c.firewallConfigurationUpdate(ctx, FirewallCategoryNetworkQuarantine, scope, data)
}

// FirewallActionFilter selects the rules a firewall-control action applies to.
type FirewallActionFilter struct {
	IDs        []string `json:"ids,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	Query      string   `json:"query,omitempty"`
	Tenant     *bool    `json:"tenant,omitempty"`
}

func (f FirewallActionFilter) isEmpty() bool {
	return len(f.IDs) == 0 && len(f.AccountIDs) == 0 && len(f.SiteIDs) == 0 &&
		len(f.GroupIDs) == 0 && f.Query == "" && f.Tenant == nil
}

// FirewallLocationValue identifies a location a rule is scoped to.
type FirewallLocationValue struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Scope string `json:"scope,omitempty"`
}

// FirewallLocationTarget is the desired location assignment for matched rules.
type FirewallLocationTarget struct {
	Type   FirewallLocationType    `json:"type"`
	Values []FirewallLocationValue `json:"values,omitempty"`
}

func (c *Client) firewallSetLocation(ctx context.Context, cat FirewallCategory, filter FirewallActionFilter, loc FirewallLocationTarget) (int, error) {
	if filter.isEmpty() {
		return 0, fmt.Errorf("mgmt: set-location requires at least one filter (ids, siteIds, groupIds, accountIds, query, or tenant)")
	}
	req := struct {
		Filter FirewallActionFilter   `json:"filter"`
		Data   FirewallLocationTarget `json:"data"`
	}{Filter: filter, Data: loc}
	var resp affectedResponse
	if err := c.post(ctx, firewallPath(cat, "/set-location"), req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// NetworkQuarantineSetLocation assigns a location to matched network quarantine rules.
func (c *Client) NetworkQuarantineSetLocation(ctx context.Context, filter FirewallActionFilter, loc FirewallLocationTarget) (int, error) {
	return c.firewallSetLocation(ctx, FirewallCategoryNetworkQuarantine, filter, loc)
}

func (c *Client) firewallMoveRules(ctx context.Context, cat FirewallCategory, filter FirewallActionFilter, targets []FirewallRuleCopyTarget) (int, error) {
	if filter.isEmpty() {
		return 0, fmt.Errorf("mgmt: move-rules requires at least one filter (ids, siteIds, groupIds, accountIds, query, or tenant)")
	}
	req := struct {
		Filter FirewallActionFilter     `json:"filter"`
		Data   []FirewallRuleCopyTarget `json:"data"`
	}{Filter: filter, Data: targets}
	var resp affectedResponse
	if err := c.post(ctx, firewallPath(cat, "/move-rules"), req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// NetworkQuarantineMoveRules moves matched network quarantine rules to target scopes.
func (c *Client) NetworkQuarantineMoveRules(ctx context.Context, filter FirewallActionFilter, targets []FirewallRuleCopyTarget) (int, error) {
	return c.firewallMoveRules(ctx, FirewallCategoryNetworkQuarantine, filter, targets)
}

func (c *Client) firewallChangeTags(ctx context.Context, cat FirewallCategory, suffix string, filter FirewallActionFilter, tagIDs []string) (int, error) {
	if filter.isEmpty() {
		return 0, fmt.Errorf("mgmt: tag change requires at least one filter (ids, siteIds, groupIds, accountIds, query, or tenant)")
	}
	req := struct {
		Filter FirewallActionFilter `json:"filter"`
		Data   struct {
			TagIDs []string `json:"tagIds"`
		} `json:"data"`
	}{Filter: filter}
	req.Data.TagIDs = tagIDs
	var resp affectedResponse
	if err := c.post(ctx, firewallPath(cat, suffix), req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// NetworkQuarantineAddTags adds tags to matched network quarantine rules.
func (c *Client) NetworkQuarantineAddTags(ctx context.Context, filter FirewallActionFilter, tagIDs []string) (int, error) {
	return c.firewallChangeTags(ctx, FirewallCategoryNetworkQuarantine, "/add-tags", filter, tagIDs)
}

// NetworkQuarantineRemoveTags removes tags from matched network quarantine rules.
func (c *Client) NetworkQuarantineRemoveTags(ctx context.Context, filter FirewallActionFilter, tagIDs []string) (int, error) {
	return c.firewallChangeTags(ctx, FirewallCategoryNetworkQuarantine, "/remove-tags", filter, tagIDs)
}
