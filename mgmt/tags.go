package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// Tag is a SentinelOne endpoint tag.
type Tag struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
	ScopeID     string `json:"scopeId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (t *Tag) UnmarshalJSON(b []byte) error {
	type alias Tag
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

// TagListParams are query parameters for listing tags.
type TagListParams struct {
	Type       string
	SiteIDs    []string
	AccountIDs []string
	Query      string
	Limit      int
	Cursor     string
}

func (p *TagListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "type", p.Type)
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// TagsList returns a paginated list of tags.
func (c *Client) TagsList(ctx context.Context, params *TagListParams) ([]Tag, *Pagination, error) {
	return list[Tag](c, ctx, "/tags", params.values())
}

// TagsGet returns a single tag by ID.
func (c *Client) TagsGet(ctx context.Context, id string) (*Tag, error) {
	return getByID[Tag](c, ctx, "/tags", "tag", id)
}
