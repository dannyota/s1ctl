package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// Application is a SentinelOne application inventory entry.
type Application struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Version       string `json:"version"`
	Publisher     string `json:"publisher"`
	Size          int64  `json:"size"`
	InstalledDate string `json:"installedDate"`
	OSType        string `json:"osType"`
	AgentID       string `json:"agentId"`

	Raw json.RawMessage `json:"-"`
}

func (a *Application) UnmarshalJSON(b []byte) error {
	type alias Application
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// ApplicationListParams are query parameters for listing applications.
type ApplicationListParams struct {
	AgentIDs []string
	SiteIDs  []string
	Query    string
	Limit    int
	Cursor   string
}

func (p *ApplicationListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "agentIds", p.AgentIDs)
	addCSV(v, "siteIds", p.SiteIDs)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// ApplicationsList returns a paginated list of installed applications.
func (c *Client) ApplicationsList(ctx context.Context, params *ApplicationListParams) ([]Application, *Pagination, error) {
	return list[Application](c, ctx, "/installed-applications", params.values())
}
