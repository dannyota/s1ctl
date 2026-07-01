package graphql

import (
	"context"
	"encoding/json"
)

// VulnerabilityCVE holds CVE information for a vulnerability.
type VulnerabilityCVE struct {
	ID              string  `json:"id"`
	NVDBaseScore    float64 `json:"nvdBaseScore"`
	RiskScore       float64 `json:"riskScore"`
	EPSSScore       float64 `json:"epssScore"`
	ExploitMaturity string  `json:"exploitMaturity"`
	ExploitedInWild bool    `json:"exploitedInTheWild"`
	PublishedDate   string  `json:"publishedDate"`
}

// VulnerabilitySoftware holds the affected software for a vulnerability.
type VulnerabilitySoftware struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	FixVersion     string `json:"fixVersion"`
	PackageManager string `json:"packageManager"`
	Type           string `json:"type"`
	Vendor         string `json:"vendor"`
}

// VulnerabilityAsset is the asset associated with a vulnerability.
type VulnerabilityAsset struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Category    string                     `json:"category"`
	Subcategory string                     `json:"subcategory"`
	Type        string                     `json:"type"`
	OsType      string                     `json:"osType"`
	CloudInfo   *MisconfigurationCloudInfo `json:"cloudInfo"`
}

// Vulnerability is an xSPM vulnerability finding.
type Vulnerability struct {
	ID             string                `json:"id"`
	ExternalID     string                `json:"externalId"`
	Name           string                `json:"name"`
	Severity       string                `json:"severity"`
	Status         string                `json:"status"`
	AnalystVerdict string                `json:"analystVerdict"`
	Product        string                `json:"product"`
	Vendor         string                `json:"vendor"`
	DetectedAt     string                `json:"detectedAt"`
	LastSeenAt     string                `json:"lastSeenAt"`
	UpdatedAt      string                `json:"updatedAt"`
	ResolvedAt     string                `json:"resolvedAt"`
	CVE            VulnerabilityCVE      `json:"cve"`
	Software       VulnerabilitySoftware `json:"software"`
	Asset          VulnerabilityAsset    `json:"asset"`
	Scope          ScopeInfo             `json:"scope"`

	Raw json.RawMessage `json:"-"`
}

func (v *Vulnerability) UnmarshalJSON(b []byte) error {
	type alias Vulnerability
	if err := json.Unmarshal(b, (*alias)(v)); err != nil {
		return err
	}
	v.Raw = append(v.Raw[:0:0], b...)
	return nil
}

// VulnerabilityEdge is a single edge in a Relay connection.
type VulnerabilityEdge struct {
	Cursor string        `json:"cursor"`
	Node   Vulnerability `json:"node"`
}

// VulnerabilityConnection is the Relay connection response for vulnerabilities.
type VulnerabilityConnection struct {
	Edges      []VulnerabilityEdge `json:"edges"`
	PageInfo   PageInfo            `json:"pageInfo"`
	TotalCount int64               `json:"totalCount"`
}

// VulnerabilityListParams are parameters for querying vulnerabilities.
type VulnerabilityListParams struct {
	First   int      `json:"first,omitempty"`
	After   string   `json:"after,omitempty"`
	Filters []Filter `json:"filters,omitempty"`
	Scope   *Scope   `json:"scope,omitempty"`
}

const vulnerabilitiesQuery = `query Vulnerabilities($first: Int, $after: String, $filters: [FilterInput!], $scope: ScopeSelectorInput) {
  vulnerabilities(first: $first, after: $after, filters: $filters, scope: $scope) {
    edges {
      cursor
      node {
        id
        externalId
        name
        severity
        status
        analystVerdict
        product
        vendor
        detectedAt
        lastSeenAt
        updatedAt
        resolvedAt
        cve {
          id nvdBaseScore riskScore epssScore
          exploitMaturity exploitedInTheWild publishedDate
        }
        software { name version fixVersion packageManager type vendor }
        asset {
          id name category subcategory type osType
          cloudInfo { accountId accountName providerName region resourceId }
        }
        scope { account { id name } site { id name } group { id name } }
      }
    }
    pageInfo { hasNextPage hasPreviousPage endCursor startCursor }
    totalCount
  }
}`

// VulnerabilitiesList queries xSPM vulnerabilities.
func (c *Client) VulnerabilitiesList(ctx context.Context, params *VulnerabilityListParams) (*VulnerabilityConnection, error) {
	vars := map[string]any{}
	if params != nil {
		if params.First > 0 {
			vars["first"] = params.First
		}
		if params.After != "" {
			vars["after"] = params.After
		}
		if len(params.Filters) > 0 {
			vars["filters"] = params.Filters
		}
		if params.Scope != nil {
			vars["scope"] = params.Scope
		}
	}
	var resp struct {
		Vulnerabilities VulnerabilityConnection `json:"vulnerabilities"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, vulnerabilitiesQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Vulnerabilities, nil
}

const vulnerabilityGetQuery = `query VulnerabilityGet($id: ID!) {
  vulnerability(id: $id) {
    id
    externalId
    name
    severity
    status
    analystVerdict
    product
    vendor
    detectedAt
    lastSeenAt
    updatedAt
    resolvedAt
    cve {
      id nvdBaseScore riskScore epssScore
      exploitMaturity exploitedInTheWild publishedDate
    }
    software { name version fixVersion packageManager type vendor }
    asset {
      id name category subcategory type osType
      cloudInfo { accountId accountName providerName region resourceId }
    }
    scope { account { id name } site { id name } group { id name } }
  }
}`

// VulnerabilitiesGet returns a single vulnerability by ID.
func (c *Client) VulnerabilitiesGet(ctx context.Context, id string) (*Vulnerability, error) {
	vars := map[string]any{"id": id}
	var resp struct {
		Vulnerability Vulnerability `json:"vulnerability"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, vulnerabilityGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Vulnerability, nil
}

const vulnerabilitiesStatusUpdateMutation = `mutation VulnerabilitiesStatusUpdate($filter: OrFilterSelectionInput, $statusUpdate: StatusUpdateInput!) {
  vulnerabilitiesStatusUpdateV2(filter: $filter, statusUpdate: $statusUpdate) {
    updatedFindingIds
  }
}`

// VulnerabilitiesUpdateStatus updates the status of the specified vulnerabilities.
func (c *Client) VulnerabilitiesUpdateStatus(ctx context.Context, ids []string, status string) error {
	vars := map[string]any{
		"filter":       orFilterByIDs(ids),
		"statusUpdate": map[string]any{"status": status},
	}
	return c.Do(ctx, EndpointVulnerabilities, vulnerabilitiesStatusUpdateMutation, vars, nil)
}

const vulnerabilitiesVerdictUpdateMutation = `mutation VulnerabilitiesVerdictUpdate($filter: OrFilterSelectionInput, $analystVerdict: AnalystVerdict) {
  vulnerabilitiesAnalystVerdictUpdateV2(filter: $filter, analystVerdict: $analystVerdict) {
    updatedFindingIds
  }
}`

// VulnerabilitiesUpdateVerdict updates the analyst verdict of the specified vulnerabilities.
func (c *Client) VulnerabilitiesUpdateVerdict(ctx context.Context, ids []string, verdict string) error {
	vars := map[string]any{
		"filter":         orFilterByIDs(ids),
		"analystVerdict": verdict,
	}
	return c.Do(ctx, EndpointVulnerabilities, vulnerabilitiesVerdictUpdateMutation, vars, nil)
}
