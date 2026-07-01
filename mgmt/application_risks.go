package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// ApplicationRisk is a CVE risk entry from Application Risk Management.
type ApplicationRisk struct {
	ID                         string `json:"id"`
	CveID                      string `json:"cveId"`
	ApplicationName            string `json:"applicationName"`
	ApplicationVendor          string `json:"applicationVendor"`
	ApplicationVersion         string `json:"applicationVersion"`
	Application                string `json:"application"`
	Severity                   string `json:"severity"`
	BaseScore                  string `json:"baseScore"`
	NvdBaseScore               string `json:"nvdBaseScore"`
	RiskScore                  string `json:"riskScore"`
	CvssVersion                string `json:"cvssVersion"`
	NvdCvssVersion             string `json:"nvdCvssVersion"`
	ExploitCodeMaturity        string `json:"exploitCodeMaturity"`
	ExploitedInTheWild         string `json:"exploitedInTheWild"`
	RemediationLevel           string `json:"remediationLevel"`
	ReportConfidence           string `json:"reportConfidence"`
	MitigationStatus           string `json:"mitigationStatus"`
	MitigationStatusReason     string `json:"mitigationStatusReason"`
	MitigationStatusChangedBy  string `json:"mitigationStatusChangedBy"`
	MitigationStatusChangeTime string `json:"mitigationStatusChangeTime"`
	Status                     string `json:"status"`
	OSType                     string `json:"osType"`
	EndpointID                 string `json:"endpointId"`
	EndpointName               string `json:"endpointName"`
	EndpointType               string `json:"endpointType"`
	DetectionDate              string `json:"detectionDate"`
	PublishedDate              string `json:"publishedDate"`
	DaysDetected               int    `json:"daysDetected"`
	LastScanDate               string `json:"lastScanDate"`
	LastScanResult             string `json:"lastScanResult"`

	Raw json.RawMessage `json:"-"`
}

func (r *ApplicationRisk) UnmarshalJSON(b []byte) error {
	type alias ApplicationRisk
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// ApplicationRiskListParams are query parameters for listing application risks.
type ApplicationRiskListParams struct {
	SiteIDs             []string
	AccountIDs          []string
	Severities          []string
	ApplicationNames    []string
	ApplicationVendor   string
	ExploitCodeMaturity []string
	ExploitedInTheWild  []string
	MitigationStatus    []string
	OSVersions          []string
	AnalystVerdict      []string
	Domains             []string
	IncludeRemovals     *bool
	SortBy              string
	SortOrder           string
	Limit               int
	Cursor              string
}

func (p *ApplicationRiskListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "severities", p.Severities)
	addCSV(v, "applicationNames", p.ApplicationNames)
	addString(v, "applicationVendor__contains", p.ApplicationVendor)
	addCSV(v, "exploitCodeMaturity", p.ExploitCodeMaturity)
	addCSV(v, "exploitedInTheWild", p.ExploitedInTheWild)
	addCSV(v, "mitigationStatus", p.MitigationStatus)
	addCSV(v, "osVersions", p.OSVersions)
	addCSV(v, "analystVerdict", p.AnalystVerdict)
	addCSV(v, "domains", p.Domains)
	addBool(v, "includeRemovals", p.IncludeRemovals)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// ApplicationRisksList returns a paginated list of application CVE risks.
func (c *Client) ApplicationRisksList(ctx context.Context, params *ApplicationRiskListParams) ([]ApplicationRisk, *Pagination, error) {
	return list[ApplicationRisk](c, ctx, "/application-management/risks", params.values())
}

// ApplicationCVE is a CVE entry from the application risk CVEs endpoint.
type ApplicationCVE struct {
	CveID               string `json:"cveId"`
	Severity            string `json:"severity"`
	NvdBaseScore        string `json:"nvdBaseScore"`
	RiskScore           string `json:"riskScore"`
	CvssVersion         string `json:"cvssVersion"`
	Description         string `json:"description"`
	NvdURL              string `json:"nvdUrl"`
	MitreURL            string `json:"mitreUrl"`
	PublishedDate       string `json:"publishedDate"`
	ExploitCodeMaturity string `json:"exploitCodeMaturity"`
	ExploitedInTheWild  string `json:"exploitedInTheWild"`
	RemediationLevel    string `json:"remediationLevel"`
	ReportConfidence    string `json:"reportConfidence"`

	Raw json.RawMessage `json:"-"`
}

func (c *ApplicationCVE) UnmarshalJSON(b []byte) error {
	type alias ApplicationCVE
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// ApplicationCVEListParams are query parameters for listing application CVEs.
type ApplicationCVEListParams struct {
	SiteIDs             []string
	AccountIDs          []string
	GroupIDs            []string
	Severities          []string
	ApplicationName     string
	ApplicationVendor   string
	ApplicationVersions []string
	ApplicationIDs      []string
	CveID               string
	ExploitCodeMaturity []string
	ExploitedInTheWild  []string
	RemediationLevels   []string
	ReportConfidence    []string
	AnalystVerdict      []string
	SortBy              string
	SortOrder           string
	Limit               int
	Cursor              string
}

func (p *ApplicationCVEListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "severities", p.Severities)
	addString(v, "applicationName", p.ApplicationName)
	addString(v, "applicationVendor", p.ApplicationVendor)
	addCSV(v, "applicationVersions", p.ApplicationVersions)
	addCSV(v, "applicationIds", p.ApplicationIDs)
	addString(v, "cveId__contains", p.CveID)
	addCSV(v, "exploitCodeMaturity", p.ExploitCodeMaturity)
	addCSV(v, "exploitedInTheWild", p.ExploitedInTheWild)
	addCSV(v, "remediationLevels", p.RemediationLevels)
	addCSV(v, "reportConfidence", p.ReportConfidence)
	addCSV(v, "analystVerdict", p.AnalystVerdict)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	return v
}

// ApplicationCVEsList returns a paginated list of application CVEs.
func (c *Client) ApplicationCVEsList(ctx context.Context, params *ApplicationCVEListParams) ([]ApplicationCVE, *Pagination, error) {
	return list[ApplicationCVE](c, ctx, "/application-management/risks/cves", params.values())
}
