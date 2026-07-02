package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// IOCType is the indicator type for a threat intelligence IOC.
type IOCType string

const (
	IOCTypeDNS    IOCType = "DNS"
	IOCTypeIPv4   IOCType = "IPV4"
	IOCTypeIPv6   IOCType = "IPV6"
	IOCTypeMD5    IOCType = "MD5"
	IOCTypeSHA1   IOCType = "SHA1"
	IOCTypeSHA256 IOCType = "SHA256"
	IOCTypeURL    IOCType = "URL"
)

// IOCSeverity is the potential impact of a threat intelligence IOC. The API
// represents it as an OCSF-style integer score in the range 0-7.
type IOCSeverity int

// OCSF severity scores. The API accepts 0-7; 7 has no OCSF name and is
// rendered numerically.
const (
	IOCSeverityUnknown       IOCSeverity = 0
	IOCSeverityInformational IOCSeverity = 1
	IOCSeverityLow           IOCSeverity = 2
	IOCSeverityMedium        IOCSeverity = 3
	IOCSeverityHigh          IOCSeverity = 4
	IOCSeverityCritical      IOCSeverity = 5
	IOCSeverityFatal         IOCSeverity = 6
)

// String returns the OCSF severity name, or the numeric score when unnamed.
func (s IOCSeverity) String() string {
	switch s {
	case IOCSeverityUnknown:
		return "Unknown"
	case IOCSeverityInformational:
		return "Informational"
	case IOCSeverityLow:
		return "Low"
	case IOCSeverityMedium:
		return "Medium"
	case IOCSeverityHigh:
		return "High"
	case IOCSeverityCritical:
		return "Critical"
	case IOCSeverityFatal:
		return "Fatal"
	default:
		return strconv.Itoa(int(s))
	}
}

// IOCScope is the scope at which a threat intelligence object is defined.
type IOCScope string

const (
	IOCScopeGlobal  IOCScope = "global"
	IOCScopeAccount IOCScope = "account"
	IOCScopeSite    IOCScope = "site"
	IOCScopeGroup   IOCScope = "group"
)

// IOC is a SentinelOne threat intelligence indicator of compromise.
type IOC struct {
	UUID              string      `json:"uuid"`
	Type              IOCType     `json:"type"`
	Value             string      `json:"value"`
	Source            string      `json:"source"`
	Severity          IOCSeverity `json:"severity"`
	Method            string      `json:"method"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	ExternalID        string      `json:"externalId"`
	BatchID           string      `json:"batchId"`
	Creator           string      `json:"creator"`
	Scope             IOCScope    `json:"scope"`
	ScopeID           string      `json:"scopeId"`
	ParentScopeID     string      `json:"parentScopeId"`
	Category          []string    `json:"category"`
	Labels            []string    `json:"labels"`
	MalwareNames      []string    `json:"malwareNames"`
	CampaignNames     []string    `json:"campaignNames"`
	ThreatActors      []string    `json:"threatActors"`
	ThreatActorTypes  []string    `json:"threatActorTypes"`
	IntrusionSets     []string    `json:"intrusionSets"`
	MitreTactic       []string    `json:"mitreTactic"`
	Metadata          string      `json:"metadata"`
	OriginalRiskScore int         `json:"originalRiskScore"`
	PatternType       string      `json:"patternType"`
	Pattern           string      `json:"pattern"`
	Reference         []string    `json:"reference"`
	ValidUntil        string      `json:"validUntil"`
	CreationTime      string      `json:"creationTime"`
	UpdatedAt         string      `json:"updatedAt"`
	UploadTime        string      `json:"uploadTime"`

	Raw json.RawMessage `json:"-"`
}

func (ioc *IOC) UnmarshalJSON(b []byte) error {
	type alias IOC
	if err := json.Unmarshal(b, (*alias)(ioc)); err != nil {
		return err
	}
	ioc.Raw = append(ioc.Raw[:0:0], b...)
	return nil
}

// IOCListParams are query parameters for listing threat intelligence IOCs.
type IOCListParams struct {
	AccountIDs []string
	SiteIDs    []string
	UUIDs      []string
	Type       IOCType
	Severities []IOCSeverity
	Sources    []string
	Value      string
	ExternalID string
	BatchID    string
	Creators   []string // free-text creator filter (creator__contains)
	Limit      int
	Cursor     string
	SortBy     string
	SortOrder  string
}

func (p *IOCListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "uuids", p.UUIDs)
	addString(v, "type", string(p.Type))
	sevs := make([]string, len(p.Severities))
	for i, s := range p.Severities {
		sevs[i] = strconv.Itoa(int(s))
	}
	addCSV(v, "severity", sevs)
	addCSV(v, "source", p.Sources)
	addString(v, "value", p.Value)
	addString(v, "externalId", p.ExternalID)
	addString(v, "batchId", p.BatchID)
	addCSV(v, "creator__contains", p.Creators)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// IOCsList returns a paginated list of threat intelligence IOCs.
func (c *Client) IOCsList(ctx context.Context, params *IOCListParams) ([]IOC, *Pagination, error) {
	return list[IOC](c, ctx, "/threat-intelligence/iocs", params.values())
}

// IOCCreateInput is the payload for creating a threat intelligence IOC.
// Source, Type, and Value are required. Method defaults to EQUALS when empty.
type IOCCreateInput struct {
	Type        IOCType      `json:"type"`
	Value       string       `json:"value"`
	Source      string       `json:"source"`
	Severity    *IOCSeverity `json:"severity,omitempty"`
	Method      string       `json:"method,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	ExternalID  string       `json:"externalId,omitempty"`
	ValidUntil  string       `json:"validUntil,omitempty"`
}

// IOCsCreate creates threat intelligence IOCs and returns the created
// indicators.
func (c *Client) IOCsCreate(ctx context.Context, iocs []IOCCreateInput) ([]IOC, error) {
	if len(iocs) == 0 {
		return nil, fmt.Errorf("mgmt: at least one IOC is required")
	}
	req := map[string]any{
		"data": iocs,
	}
	var resp struct {
		Data []IOC `json:"data"`
	}
	if err := c.post(ctx, "/threat-intelligence/iocs", req, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IOCsDelete deletes threat intelligence IOCs by UUID.
func (c *Client) IOCsDelete(ctx context.Context, uuids []string) (int, error) {
	if len(uuids) == 0 {
		return 0, fmt.Errorf("mgmt: at least one IOC UUID is required")
	}
	req := map[string]any{
		"filter": map[string]any{
			"uuids": uuids,
		},
	}
	var resp affectedResponse
	if err := c.jsonRequest(ctx, http.MethodDelete, "/threat-intelligence/iocs", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// ThreatIntelConfig is a threat intelligence user configuration entry.
type ThreatIntelConfig struct {
	ScopeID             string   `json:"scopeId"`
	ScopeLevel          IOCScope `json:"scopeLevel"`
	Description         string   `json:"description"`
	ThreatMinScore      int      `json:"threatMinScore"`
	ThreatExcludeFields []string `json:"threatExcludeFields"`
	ExcludeTII          []string `json:"excludeTii"`
	DisableRH           bool     `json:"disableRh"`
	DisableThreat       bool     `json:"disableThreat"`
	EnableXDRMatching   bool     `json:"enableXdrMatching"`
	CreatedAt           string   `json:"createdAt"`
	UpdatedAt           string   `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (c *ThreatIntelConfig) UnmarshalJSON(b []byte) error {
	type alias ThreatIntelConfig
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// ThreatIntelConfigs returns the threat intelligence user configuration
// entries, one per configured scope.
func (c *Client) ThreatIntelConfigs(ctx context.Context) ([]ThreatIntelConfig, error) {
	configs, _, err := list[ThreatIntelConfig](c, ctx, "/threat-intelligence/user-config", nil)
	return configs, err
}
