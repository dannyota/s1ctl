package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// RemoteScript is a SentinelOne remote ops script.
type RemoteScript struct {
	ID          string   `json:"id"`
	FileName    string   `json:"fileName"`
	FileType    string   `json:"fileType"`
	ScriptType  string   `json:"scriptType"`
	OSTypes     []string `json:"osTypes"`
	ScopeID     string   `json:"scopeId"`
	ScopeLevel  string   `json:"scopeLevel"`
	CreatedAt   string   `json:"createdAt"`
	CreatorID   string   `json:"creatorId"`
	CreatorName string   `json:"creatorName"`

	Raw json.RawMessage `json:"-"`
}

func (r *RemoteScript) UnmarshalJSON(b []byte) error {
	type alias RemoteScript
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// RemoteScriptListParams are query parameters for listing remote scripts.
type RemoteScriptListParams struct {
	SiteIDs    []string
	AccountIDs []string
	OSTypes    []string
	Query      string
	Limit      int
	Cursor     string
}

func (p *RemoteScriptListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "osTypes", p.OSTypes)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// RemoteScriptsList returns a paginated list of remote scripts.
func (c *Client) RemoteScriptsList(ctx context.Context, params *RemoteScriptListParams) ([]RemoteScript, *Pagination, error) {
	return list[RemoteScript](c, ctx, "/remote-scripts", params.values())
}

// RemoteScriptsGet returns a single remote script by ID.
func (c *Client) RemoteScriptsGet(ctx context.Context, id string) (*RemoteScript, error) {
	return getByID[RemoteScript](c, ctx, "/remote-scripts", "remote script", id)
}
