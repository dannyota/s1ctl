package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type ThreatTimelineEntry struct {
	ID                   string          `json:"id"`
	ActivityType         int             `json:"activityType"`
	PrimaryDescription   string          `json:"primaryDescription"`
	SecondaryDescription string          `json:"secondaryDescription"`
	Data                 json.RawMessage `json:"data"`
	AccountID            string          `json:"accountId"`
	SiteID               string          `json:"siteId"`
	GroupID              string          `json:"groupId"`
	AgentID              string          `json:"agentId"`
	ThreatID             string          `json:"threatId"`
	UserID               string          `json:"userId"`
	Hash                 string          `json:"hash"`
	OSFamily             string          `json:"osFamily"`
	CreatedAt            string          `json:"createdAt"`
	UpdatedAt            string          `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (e ThreatTimelineEntry) MarshalJSON() ([]byte, error) {
	if e.Raw != nil {
		return e.Raw, nil
	}
	type alias ThreatTimelineEntry
	return json.Marshal(alias(e))
}

func (e *ThreatTimelineEntry) UnmarshalJSON(b []byte) error {
	type alias ThreatTimelineEntry
	if err := json.Unmarshal(b, (*alias)(e)); err != nil {
		return err
	}
	e.Raw = append(e.Raw[:0:0], b...)
	return nil
}

type ThreatTimelineParams struct {
	ActivityTypes []int
	Query         string
	Limit         int
	Cursor        string
	SortBy        string
	SortOrder     string
}

func (p *ThreatTimelineParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addIntCSV(v, "activityTypes", p.ActivityTypes)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// ThreatTimeline returns the activity timeline for a threat.
func (c *Client) ThreatTimeline(ctx context.Context, threatID string, params *ThreatTimelineParams) ([]ThreatTimelineEntry, *Pagination, error) {
	if threatID == "" {
		return nil, nil, fmt.Errorf("mgmt: threat ID is required")
	}
	path := fmt.Sprintf("/threats/%s/timeline", threatID)
	return list[ThreatTimelineEntry](c, ctx, path, params.values())
}
