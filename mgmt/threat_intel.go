package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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

// IOCSeverity is the severity level of a threat intelligence IOC.
type IOCSeverity string

const (
	IOCSeverityLow    IOCSeverity = "Low"
	IOCSeverityMedium IOCSeverity = "Medium"
	IOCSeverityHigh   IOCSeverity = "High"
)

// IOC is a SentinelOne threat intelligence indicator of compromise.
type IOC struct {
	ID           string      `json:"id"`
	Type         IOCType     `json:"type"`
	Value        string      `json:"value"`
	Source       string      `json:"source"`
	Severity     IOCSeverity `json:"severity"`
	Method       string      `json:"method"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	ExternalID   string      `json:"externalId"`
	BatchID      string      `json:"batchId"`
	Creator      string      `json:"creator"`
	CreatorID    string      `json:"creatorId"`
	Scope        string      `json:"scope"`
	ScopeID      string      `json:"scopeId"`
	AccountIDs   []string    `json:"accountIds"`
	PatternType  string      `json:"patternType"`
	Pattern      string      `json:"pattern"`
	Reference    []string    `json:"reference"`
	ValidUntil   string      `json:"validUntil"`
	CreationTime string      `json:"creationTime"`
	UpdatedAt    string      `json:"updatedAt"`
	UploadTime   string      `json:"uploadTime"`

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
	Types      []string
	Severities []string
	Sources    []string
	Value      string
	BatchID    string
	Creator    string
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
	addCSV(v, "type", p.Types)
	addCSV(v, "severities", p.Severities)
	addCSV(v, "source", p.Sources)
	addString(v, "value", p.Value)
	addString(v, "batchId", p.BatchID)
	addString(v, "creator", p.Creator)
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
type IOCCreateInput struct {
	Type        IOCType     `json:"type"`
	Value       string      `json:"value"`
	Source      string      `json:"source,omitempty"`
	Severity    IOCSeverity `json:"severity,omitempty"`
	Method      string      `json:"method,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	ExternalID  string      `json:"externalId,omitempty"`
	ValidUntil  string      `json:"validUntil,omitempty"`
}

// IOCsCreate creates one or more threat intelligence IOCs.
func (c *Client) IOCsCreate(ctx context.Context, iocs []IOCCreateInput) (int, error) {
	req := map[string]any{
		"data": iocs,
	}
	var resp affectedResponse
	if err := c.post(ctx, "/threat-intelligence/iocs", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// IOCsDelete deletes threat intelligence IOCs by ID.
func (c *Client) IOCsDelete(ctx context.Context, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("mgmt: at least one IOC ID is required")
	}
	req := map[string]any{
		"data": map[string]any{},
		"filter": map[string]any{
			"ids": ids,
		},
	}
	var resp affectedResponse
	if err := c.post(ctx, "/threat-intelligence/iocs/delete", req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// ThreatIntelConfig is the user's threat intelligence configuration.
type ThreatIntelConfig struct {
	AccountIDs []string `json:"accountIds"`
	SiteIDs    []string `json:"siteIds"`
	TotalIOCs  int      `json:"totalIocs"`
	MaxIOCs    int      `json:"maxIocs"`

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

// ThreatIntelConfig returns the user's threat intelligence configuration.
func (c *Client) ThreatIntelConfig(ctx context.Context) (*ThreatIntelConfig, error) {
	var resp singleResponse[ThreatIntelConfig]
	if err := c.get(ctx, "/threat-intelligence/user-config", nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
