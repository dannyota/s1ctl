package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// Activity is a SentinelOne activity log entry.
type Activity struct {
	ID            string          `json:"id"`
	ActivityType  int             `json:"activityType"`
	PrimaryDesc   string          `json:"primaryDescription"`
	SecondaryDesc string          `json:"secondaryDescription"`
	AccountID     string          `json:"accountId"`
	AccountName   string          `json:"accountName"`
	SiteID        string          `json:"siteId"`
	SiteName      string          `json:"siteName"`
	GroupID       string          `json:"groupId"`
	GroupName     string          `json:"groupName"`
	AgentID       string          `json:"agentId"`
	ThreatID      string          `json:"threatId"`
	UserID        string          `json:"userId"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
	Data          json.RawMessage `json:"data"`

	Raw json.RawMessage `json:"-"`
}

func (a *Activity) UnmarshalJSON(b []byte) error {
	type alias Activity
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// ActivityListParams are query parameters for listing activities.
type ActivityListParams struct {
	SiteIDs       []string
	AccountIDs    []string
	GroupIDs      []string
	AgentIDs      []string
	ThreatIDs     []string
	ActivityTypes []int
	UserIDs       []string
	CreatedAtGt   string
	CreatedAtLt   string
	Limit         int
	Cursor        string
	SortBy        string
	SortOrder     string
	CountOnly     bool
}

func (p *ActivityListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "agentIds", p.AgentIDs)
	addCSV(v, "threatIds", p.ThreatIDs)
	addCSV(v, "userIds", p.UserIDs)
	addIntCSV(v, "activityTypes", p.ActivityTypes)
	addString(v, "createdAt__gt", p.CreatedAtGt)
	addString(v, "createdAt__lt", p.CreatedAtLt)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	if p.CountOnly {
		v.Set("countOnly", "true")
	}
	return v
}

// ActivitiesList returns a paginated list of activities.
func (c *Client) ActivitiesList(ctx context.Context, params *ActivityListParams) ([]Activity, *Pagination, error) {
	return list[Activity](c, ctx, "/activities", params.values())
}
