package mgmt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// FirewallProtocol is a protocol available for firewall rules.
type FirewallProtocol struct {
	Value string `json:"value"`
	Name  string `json:"name"`

	Raw json.RawMessage `json:"-"`
}

func (f *FirewallProtocol) UnmarshalJSON(b []byte) error {
	type alias FirewallProtocol
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// FirewallProtocolListParams are query parameters for listing protocols.
type FirewallProtocolListParams struct {
	Query string
	Limit int
}

func (p *FirewallProtocolListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	return v
}

// FirewallProtocolsList returns the protocols available for firewall rules.
func (c *Client) FirewallProtocolsList(ctx context.Context, params *FirewallProtocolListParams) ([]FirewallProtocol, *Pagination, error) {
	return c.firewallProtocolsList(ctx, FirewallCategoryFirewall, params)
}

func (c *Client) firewallProtocolsList(ctx context.Context, cat FirewallCategory, params *FirewallProtocolListParams) ([]FirewallProtocol, *Pagination, error) {
	return list[FirewallProtocol](c, ctx, firewallPath(cat, "/protocols"), params.values())
}

// FirewallRulesDelete deletes firewall rules by ID.
func (c *Client) FirewallRulesDelete(ctx context.Context, ids []string) (int, error) {
	return c.firewallRulesDelete(ctx, FirewallCategoryFirewall, ids)
}

func (c *Client) firewallRulesDelete(ctx context.Context, cat FirewallCategory, ids []string) (int, error) {
	req := map[string]any{
		"filter": map[string]any{
			"ids": ids,
		},
	}
	var resp affectedResponse
	if err := c.jsonRequest(ctx, http.MethodDelete, firewallPath(cat, ""), req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// FirewallRuleReorderFilter scopes a reorder operation.
type FirewallRuleReorderFilter struct {
	AccountIDs []string `json:"accountIds,omitempty"`
	SiteIDs    []string `json:"siteIds,omitempty"`
	GroupIDs   []string `json:"groupIds,omitempty"`
	Tenant     *bool    `json:"tenant,omitempty"`
}

// FirewallRulesReorder changes the order of firewall rules within a scope.
func (c *Client) FirewallRulesReorder(ctx context.Context, orders []RuleOrder, filter FirewallRuleReorderFilter) error {
	return c.firewallRulesReorder(ctx, FirewallCategoryFirewall, orders, filter)
}

func (c *Client) firewallRulesReorder(ctx context.Context, cat FirewallCategory, orders []RuleOrder, filter FirewallRuleReorderFilter) error {
	req := struct {
		Data   []RuleOrder               `json:"data"`
		Filter FirewallRuleReorderFilter `json:"filter"`
	}{Data: orders, Filter: filter}
	var resp struct {
		Data struct {
			Success bool `json:"success"`
		} `json:"data"`
	}
	return c.put(ctx, firewallPath(cat, "/reorder"), req, &resp)
}

// FirewallRuleCopyTarget specifies a destination scope for copying rules.
type FirewallRuleCopyTarget struct {
	AccountID *string `json:"accountId,omitempty"`
	SiteID    *string `json:"siteId,omitempty"`
	GroupID   *string `json:"groupId,omitempty"`
	Tenant    *bool   `json:"tenant,omitempty"`
}

// FirewallRulesCopy copies firewall rules from a source scope to targets.
func (c *Client) FirewallRulesCopy(ctx context.Context, filter FirewallRuleReorderFilter, targets []FirewallRuleCopyTarget) (int, error) {
	return c.firewallRulesCopy(ctx, FirewallCategoryFirewall, filter, targets)
}

func (c *Client) firewallRulesCopy(ctx context.Context, cat FirewallCategory, filter FirewallRuleReorderFilter, targets []FirewallRuleCopyTarget) (int, error) {
	req := struct {
		Filter FirewallRuleReorderFilter `json:"filter"`
		Data   []FirewallRuleCopyTarget  `json:"data"`
	}{Filter: filter, Data: targets}
	var resp affectedResponse
	if err := c.post(ctx, firewallPath(cat, "/copy-rules"), req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// FirewallRulesSetStatus enables or disables firewall rules by ID.
func (c *Client) FirewallRulesSetStatus(ctx context.Context, ids []string, status FirewallStatus) (int, error) {
	return c.firewallRulesSetStatus(ctx, FirewallCategoryFirewall, ids, status)
}

func (c *Client) firewallRulesSetStatus(ctx context.Context, cat FirewallCategory, ids []string, status FirewallStatus) (int, error) {
	req := struct {
		Filter struct {
			IDs []string `json:"ids"`
		} `json:"filter"`
		Data struct {
			Status FirewallStatus `json:"status"`
		} `json:"data"`
	}{}
	req.Filter.IDs = ids
	req.Data.Status = status
	var resp affectedResponse
	if err := c.put(ctx, firewallPath(cat, "/enable"), req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// FirewallRulesExport exports firewall rules as raw JSON for the given scope.
func (c *Client) FirewallRulesExport(ctx context.Context, params *FirewallRuleListParams) ([]byte, error) {
	return c.firewallRulesExport(ctx, FirewallCategoryFirewall, params)
}

func (c *Client) firewallRulesExport(ctx context.Context, cat FirewallCategory, params *FirewallRuleListParams) ([]byte, error) {
	u := c.baseURL + firewallPath(cat, "/export")
	v := params.values()
	if len(v) > 0 {
		u += "?" + v.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("mgmt: %w", err)
	}
	if wErr := c.limiter.Wait(ctx); wErr != nil {
		return nil, fmt.Errorf("mgmt: rate limit: %w", wErr)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mgmt: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("mgmt: read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, parseError(resp.StatusCode, data)
	}
	return data, nil
}

// FirewallImportScope identifies the target scope for importing rules.
type FirewallImportScope struct {
	AccountIDs []string
	SiteIDs    []string
	GroupIDs   []string
	Tenant     bool
}

// FirewallRulesImport imports firewall rules from a JSON file into the given scope.
func (c *Client) FirewallRulesImport(ctx context.Context, scope FirewallImportScope, filename string, fileData []byte) error {
	return c.firewallRulesImport(ctx, FirewallCategoryFirewall, scope, filename, fileData)
}

func (c *Client) firewallRulesImport(ctx context.Context, cat FirewallCategory, scope FirewallImportScope, filename string, fileData []byte) error {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	for _, id := range scope.AccountIDs {
		if err := w.WriteField("accountIds", id); err != nil {
			return fmt.Errorf("mgmt: write field: %w", err)
		}
	}
	for _, id := range scope.SiteIDs {
		if err := w.WriteField("siteIds", id); err != nil {
			return fmt.Errorf("mgmt: write field: %w", err)
		}
	}
	for _, id := range scope.GroupIDs {
		if err := w.WriteField("groupIds", id); err != nil {
			return fmt.Errorf("mgmt: write field: %w", err)
		}
	}
	if scope.Tenant {
		if err := w.WriteField("tenant", "true"); err != nil {
			return fmt.Errorf("mgmt: write field: %w", err)
		}
	}

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("mgmt: create form file: %w", err)
	}
	if _, err := fw.Write(fileData); err != nil {
		return fmt.Errorf("mgmt: write file: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mgmt: close multipart: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+firewallPath(cat, "/import"), &body)
	if err != nil {
		return fmt.Errorf("mgmt: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	if wErr := c.limiter.Wait(ctx); wErr != nil {
		return fmt.Errorf("mgmt: rate limit: %w", wErr)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("mgmt: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mgmt: read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return parseError(resp.StatusCode, data)
	}
	return nil
}
