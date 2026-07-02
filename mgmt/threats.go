package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Threat is a SentinelOne threat.
// The API returns nested objects (threatInfo, agentRealtimeInfo); fields
// are flattened here via a custom UnmarshalJSON.
type Threat struct {
	ID                   string `json:"-"`
	AgentID              string `json:"-"`
	AgentComputerName    string `json:"-"`
	Classification       string `json:"-"`
	ClassificationSource string `json:"-"`
	ConfidenceLevel      string `json:"-"`
	ThreatName           string `json:"-"`
	FilePath             string `json:"-"`
	MitigationStatus     string `json:"-"`
	AnalystVerdict       string `json:"-"`
	IncidentStatus       string `json:"-"`
	CreatedAt            string `json:"-"`
	UpdatedAt            string `json:"-"`

	Raw json.RawMessage `json:"-"`
}

func (t Threat) MarshalJSON() ([]byte, error) {
	if t.Raw != nil {
		return t.Raw, nil
	}
	return []byte("{}"), nil
}

func (t *Threat) UnmarshalJSON(b []byte) error {
	t.Raw = append(t.Raw[:0:0], b...)
	// Threats use nested objects (threatInfo, agentRealtimeInfo).
	// Parse the raw JSON and extract fields from nested paths.
	var flat map[string]json.RawMessage
	if err := json.Unmarshal(b, &flat); err != nil {
		return err
	}
	if v, ok := flat["id"]; ok {
		if err := json.Unmarshal(v, &t.ID); err != nil {
			return err
		}
	}
	if ti, ok := flat["threatInfo"]; ok {
		var info struct {
			Classification       string `json:"classification"`
			ClassificationSource string `json:"classificationSource"`
			ConfidenceLevel      string `json:"confidenceLevel"`
			ThreatName           string `json:"threatName"`
			FilePath             string `json:"filePath"`
			MitigationStatus     string `json:"mitigationStatus"`
			AnalystVerdict       string `json:"analystVerdict"`
			IncidentStatus       string `json:"incidentStatus"`
			CreatedAt            string `json:"createdAt"`
			UpdatedAt            string `json:"updatedAt"`
		}
		if err := json.Unmarshal(ti, &info); err == nil {
			t.Classification = info.Classification
			t.ClassificationSource = info.ClassificationSource
			t.ConfidenceLevel = info.ConfidenceLevel
			t.ThreatName = info.ThreatName
			t.FilePath = info.FilePath
			t.MitigationStatus = info.MitigationStatus
			t.AnalystVerdict = info.AnalystVerdict
			t.IncidentStatus = info.IncidentStatus
			t.CreatedAt = info.CreatedAt
			t.UpdatedAt = info.UpdatedAt
		}
	}
	if ari, ok := flat["agentRealtimeInfo"]; ok {
		var info struct {
			AgentID           string `json:"agentId"`
			AgentComputerName string `json:"agentComputerName"`
		}
		if err := json.Unmarshal(ari, &info); err == nil {
			t.AgentID = info.AgentID
			t.AgentComputerName = info.AgentComputerName
		}
	}
	return nil
}

// ThreatListParams are query parameters for listing threats.
type ThreatListParams struct {
	SiteIDs            []string
	AccountIDs         []string
	GroupIDs           []string
	AgentIDs           []string
	Classifications    []string
	MitigationStatuses []string
	AnalystVerdicts    []string
	IncidentStatuses   []string
	ConfidenceLevels   []string
	Query              string
	Limit              int
	Cursor             string
	SortBy             string
	SortOrder          string
	CountOnly          bool
}

func (p *ThreatListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "agentIds", p.AgentIDs)
	addCSV(v, "classifications", p.Classifications)
	addCSV(v, "mitigationStatuses", p.MitigationStatuses)
	addCSV(v, "analystVerdicts", p.AnalystVerdicts)
	addCSV(v, "incidentStatuses", p.IncidentStatuses)
	addCSV(v, "confidenceLevels", p.ConfidenceLevels)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	if p.CountOnly {
		v.Set("countOnly", "true")
	}
	return v
}

// ThreatsList returns a paginated list of threats.
func (c *Client) ThreatsList(ctx context.Context, params *ThreatListParams) ([]Threat, *Pagination, error) {
	return list[Threat](c, ctx, "/threats", params.values())
}

// ThreatsCount returns the count of threats matching the filter.
func (c *Client) ThreatsCount(ctx context.Context, params *ThreatListParams) (int, error) {
	if params == nil {
		params = &ThreatListParams{}
	}
	params.CountOnly = true
	_, pag, err := list[Threat](c, ctx, "/threats", params.values())
	if err != nil {
		return 0, err
	}
	return pag.TotalItems, nil
}

// ThreatsGet returns a single threat by ID.
func (c *Client) ThreatsGet(ctx context.Context, id string) (*Threat, error) {
	return getByID[Threat](c, ctx, "/threats", "threat", id)
}

// ThreatsExport exports threats matching the filter as CSV bytes.
func (c *Client) ThreatsExport(ctx context.Context, params *ThreatListParams) ([]byte, error) {
	return c.getRaw(ctx, "/threats/export", params.values())
}

// QuarantinedFile is a file quarantined for a threat.
type QuarantinedFile struct {
	FilePath string `json:"filePath"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`

	Raw json.RawMessage `json:"-"`
}

func (f QuarantinedFile) MarshalJSON() ([]byte, error) {
	if f.Raw != nil {
		return f.Raw, nil
	}
	type alias QuarantinedFile
	return json.Marshal(alias(f))
}

func (f *QuarantinedFile) UnmarshalJSON(b []byte) error {
	type alias QuarantinedFile
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// ThreatsQuarantinedFiles returns the files quarantined for a threat.
func (c *Client) ThreatsQuarantinedFiles(ctx context.Context, threatID string) ([]QuarantinedFile, error) {
	if threatID == "" {
		return nil, fmt.Errorf("mgmt: threat ID is required")
	}
	path := fmt.Sprintf("/threats/%s/quarantined-files", url.PathEscape(threatID))
	items, _, err := list[QuarantinedFile](c, ctx, path, nil)
	return items, err
}

// WhiteningOptions describes the exclusion ("whitening") options available
// for a threat.
type WhiteningOptions struct {
	WhiteningOptions []string `json:"whiteningOptions"`
	ThreatType       []string `json:"threatType"`
	ThreatPolicy     string   `json:"threatPolicy"`

	Raw json.RawMessage `json:"-"`
}

func (o WhiteningOptions) MarshalJSON() ([]byte, error) {
	if o.Raw != nil {
		return o.Raw, nil
	}
	type alias WhiteningOptions
	return json.Marshal(alias(o))
}

func (o *WhiteningOptions) UnmarshalJSON(b []byte) error {
	type alias WhiteningOptions
	if err := json.Unmarshal(b, (*alias)(o)); err != nil {
		return err
	}
	o.Raw = append(o.Raw[:0:0], b...)
	return nil
}

// ThreatsWhiteningOptions returns the exclusion options available for a threat.
func (c *Client) ThreatsWhiteningOptions(ctx context.Context, threatID string) (*WhiteningOptions, error) {
	if threatID == "" {
		return nil, fmt.Errorf("mgmt: threat ID is required")
	}
	path := fmt.Sprintf("/threats/%s/whitening-options", url.PathEscape(threatID))
	var resp singleResponse[WhiteningOptions]
	if err := c.get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
