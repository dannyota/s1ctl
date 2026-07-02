package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// TagRule is a dynamic asset tag rule: a set of conditions that automatically
// applies tags to matching XDR inventory assets.
//
// conditions, scopes, tags, and excludedAssets are nested structures captured
// verbatim as raw blobs rather than fully typed: the condition tree is large
// and open-ended, and keeping it raw lets a rule round-trip faithfully.
type TagRule struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Status         string          `json:"status"`
	SiteID         string          `json:"siteId"`
	AccountID      string          `json:"accountId"`
	MgmtID         string          `json:"mgmtId"`
	Conditions     json.RawMessage `json:"conditions,omitempty"`
	Scopes         json.RawMessage `json:"scopes,omitempty"`
	Tags           json.RawMessage `json:"tags,omitempty"`
	ExcludedAssets json.RawMessage `json:"excludedAssets,omitempty"`
	CreatedByEmail string          `json:"createdByEmail"`
	UpdatedByEmail string          `json:"updatedByEmail"`
	CreatedAt      string          `json:"createdAt"`
	UpdatedAt      string          `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (t *TagRule) UnmarshalJSON(b []byte) error {
	type alias TagRule
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

// TagRuleListParams are query parameters for listing dynamic tag rules.
type TagRuleListParams struct {
	Name       string
	Status     string
	TagIDs     []string
	IDs        []string
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
	SortBy     string
	SortOrder  string
	Limit      int
	Cursor     string
}

func (p *TagRuleListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "name", p.Name)
	addString(v, "status", p.Status)
	addCSV(v, "tagIds", p.TagIDs)
	addCSV(v, "ids", p.IDs)
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// TagRuleWrite is the request body for creating, updating, or testing a tag
// rule. name and conditions are required; ID is set on update to identify the
// rule. The body is sent bare (not wrapped in a data envelope).
type TagRuleWrite struct {
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name"`
	Description    string          `json:"description,omitempty"`
	Status         string          `json:"status,omitempty"`
	Conditions     json.RawMessage `json:"conditions,omitempty"`
	Scopes         json.RawMessage `json:"scopes,omitempty"`
	Tags           json.RawMessage `json:"tags,omitempty"`
	ExcludedAssets json.RawMessage `json:"excludedAssets,omitempty"`
}

// TagRulesList returns a paginated list of dynamic tag rules.
func (c *Client) TagRulesList(ctx context.Context, params *TagRuleListParams) ([]TagRule, *Pagination, error) {
	return list[TagRule](c, ctx, "/xdr/assets/tags/rules", params.values())
}

// TagRulesCreate creates a new dynamic tag rule and returns it.
func (c *Client) TagRulesCreate(ctx context.Context, body TagRuleWrite) (*TagRule, error) {
	var rule TagRule
	if err := c.post(ctx, "/xdr/assets/tags/rules", body, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// TagRulesUpdate updates a dynamic tag rule. The rule ID travels in the body.
func (c *Client) TagRulesUpdate(ctx context.Context, body TagRuleWrite) (*TagRule, error) {
	var rule TagRule
	if err := c.put(ctx, "/xdr/assets/tags/rules", body, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// TagRulesDelete deletes dynamic tag rules by ID. The IDs travel as a query
// parameter (the endpoint is on the collection path).
func (c *Client) TagRulesDelete(ctx context.Context, ids []string) error {
	v := url.Values{}
	addCSV(v, "ids", ids)
	return c.queryRequest(ctx, "DELETE", "/xdr/assets/tags/rules", v, nil)
}

// TagRulesTest reports how many inventory assets a candidate tag rule matches,
// without saving it. The count is the total-items field of the match response.
func (c *Client) TagRulesTest(ctx context.Context, body TagRuleWrite) (int, error) {
	var resp struct {
		Pagination Pagination `json:"pagination"`
	}
	if err := c.post(ctx, "/xdr/assets/tags/rules/test", body, &resp); err != nil {
		return 0, err
	}
	return resp.Pagination.TotalItems, nil
}
