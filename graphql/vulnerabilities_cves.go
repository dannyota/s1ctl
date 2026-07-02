package graphql

import (
	"context"
	"encoding/json"
)

// Cve holds CVE detail information returned by the cve and cves queries. It is
// a projection of the schema's CveDetail type; richer nested fields
// (riskIndicators, s1BaseValues, timeline) remain available via Raw.
type Cve struct {
	ID                 string  `json:"id"`
	Description        string  `json:"description"`
	NVDBaseScore       float64 `json:"nvdBaseScore"`
	RiskScore          float64 `json:"riskScore"`
	Score              float64 `json:"score"`
	EPSSScore          float64 `json:"epssScore"`
	EPSSPercentile     float64 `json:"epssPercentile"`
	ExploitMaturity    string  `json:"exploitMaturity"`
	ExploitedInTheWild bool    `json:"exploitedInTheWild"`
	RemediationLevel   string  `json:"remediationLevel"`
	ReportConfidence   string  `json:"reportConfidence"`
	KevAvailable       bool    `json:"kevAvailable"`
	PublishedDate      string  `json:"publishedDate"`
	NVDReferenceURL    string  `json:"nvdReferenceUrl"`
	MitreReferenceURL  string  `json:"mitreReferenceUrl"`

	Raw json.RawMessage `json:"-"`
}

func (c *Cve) UnmarshalJSON(b []byte) error {
	type alias Cve
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// CveFilter is a server-side filter for the cves query (CveFilterInput).
// CveFilterInput supports only datetime-range filtering on a CveFilterField;
// filtering by CVSS is applied client-side by callers.
type CveFilter struct {
	Field         string          `json:"field"`
	DateTimeRange json.RawMessage `json:"dateTimeRange,omitempty"`
}

// ApplicationStats is a vulnerable-application entry from topVulnerableApplications.
type ApplicationStats struct {
	Name             string  `json:"name"`
	Version          string  `json:"version"`
	AssetCount       int64   `json:"assetCount"`
	CveCount         int64   `json:"cveCount"`
	HighestRiskScore float64 `json:"highestRiskScore"`
}

// AssetStats is a vulnerable-asset entry from topVulnerableAssets.
type AssetStats struct {
	Name             string  `json:"name"`
	ScopeName        string  `json:"scopeName"`
	CveCount         int64   `json:"cveCount"`
	HighestRiskScore float64 `json:"highestRiskScore"`
}

// OsTypeStats is a vulnerable-OS-type entry from topVulnerableOsTypes.
type OsTypeStats struct {
	Name             string  `json:"name"`
	Version          string  `json:"version"`
	AssetCount       int64   `json:"assetCount"`
	CveCount         int64   `json:"cveCount"`
	AverageRiskScore float64 `json:"averageRiskScore"`
}

const cveFields = `id
        description
        nvdBaseScore
        riskScore
        score
        epssScore
        epssPercentile
        exploitMaturity
        exploitedInTheWild
        remediationLevel
        reportConfidence
        kevAvailable
        publishedDate
        nvdReferenceUrl
        mitreReferenceUrl`

const cvesQuery = `query Cves($first: Int, $after: String, $filters: [CveFilterInput], $scope: ScopeSelectorInput) {
  cves(first: $first, after: $after, filters: $filters, scope: $scope) {
    edges {
      cursor
      node { ` + cveFields + ` }
    }
    pageInfo { hasNextPage hasPreviousPage endCursor startCursor }
    totalCount
  }
}`

// CvesList queries CVEs. filters may be nil; CveFilterInput supports only
// datetime-range filtering, so richer selection (e.g. by CVSS) is done by the
// caller after fetching.
func (c *Client) CvesList(ctx context.Context, filters []CveFilter, scope *Scope, first int, after string) (*Connection[Cve], error) {
	vars := map[string]any{}
	if first > 0 {
		vars["first"] = first
	}
	if after != "" {
		vars["after"] = after
	}
	if len(filters) > 0 {
		vars["filters"] = filters
	}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		Cves Connection[Cve] `json:"cves"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, cvesQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Cves, nil
}

const cveGetQuery = `query Cve($id: String!) {
  cve(id: $id) { ` + cveFields + ` }
}`

// CveGet returns a single CVE by ID.
func (c *Client) CveGet(ctx context.Context, id string) (*Cve, error) {
	vars := map[string]any{"id": id}
	var resp struct {
		Cve Cve `json:"cve"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, cveGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Cve, nil
}

// statsVars builds the shared variables for the top* / uniqueCveCount queries.
func statsVars(filters []Filter, scope *Scope) map[string]any {
	vars := map[string]any{}
	if len(filters) > 0 {
		vars["filters"] = filters
	}
	if scope != nil {
		vars["scope"] = scope
	}
	return vars
}

const uniqueCveCountQuery = `query UniqueCveCount($filters: [FilterInput!], $scope: ScopeSelectorInput) {
  uniqueCveCount(filters: $filters, scope: $scope) {
    count
  }
}`

// UniqueCveCount returns the count of unique CVEs matching the filters.
func (c *Client) UniqueCveCount(ctx context.Context, filters []Filter, scope *Scope) (int64, error) {
	var resp struct {
		UniqueCveCount struct {
			Count int64 `json:"count"`
		} `json:"uniqueCveCount"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, uniqueCveCountQuery, statsVars(filters, scope), &resp); err != nil {
		return 0, err
	}
	return resp.UniqueCveCount.Count, nil
}

const topVulnerableApplicationsQuery = `query TopVulnerableApplications($filters: [FilterInput!], $limit: Int!, $scope: ScopeSelectorInput) {
  topVulnerableApplications(filters: $filters, limit: $limit, scope: $scope) {
    applicationStats { name version assetCount cveCount highestRiskScore }
  }
}`

// TopVulnerableApplications returns the most vulnerable applications.
func (c *Client) TopVulnerableApplications(ctx context.Context, filters []Filter, scope *Scope, limit int) ([]ApplicationStats, error) {
	vars := statsVars(filters, scope)
	vars["limit"] = limit
	var resp struct {
		TopVulnerableApplications struct {
			ApplicationStats []ApplicationStats `json:"applicationStats"`
		} `json:"topVulnerableApplications"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, topVulnerableApplicationsQuery, vars, &resp); err != nil {
		return nil, err
	}
	return resp.TopVulnerableApplications.ApplicationStats, nil
}

const topVulnerableAssetsQuery = `query TopVulnerableAssets($filters: [FilterInput!], $limit: Int!, $scope: ScopeSelectorInput) {
  topVulnerableAssets(filters: $filters, limit: $limit, scope: $scope) {
    assetStats { name scopeName cveCount highestRiskScore }
  }
}`

// TopVulnerableAssets returns the most vulnerable assets.
func (c *Client) TopVulnerableAssets(ctx context.Context, filters []Filter, scope *Scope, limit int) ([]AssetStats, error) {
	vars := statsVars(filters, scope)
	vars["limit"] = limit
	var resp struct {
		TopVulnerableAssets struct {
			AssetStats []AssetStats `json:"assetStats"`
		} `json:"topVulnerableAssets"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, topVulnerableAssetsQuery, vars, &resp); err != nil {
		return nil, err
	}
	return resp.TopVulnerableAssets.AssetStats, nil
}

const topVulnerableOsTypesQuery = `query TopVulnerableOsTypes($filters: [FilterInput!], $limit: Int!, $scope: ScopeSelectorInput) {
  topVulnerableOsTypes(filters: $filters, limit: $limit, scope: $scope) {
    osTypeStats { name version assetCount cveCount averageRiskScore }
  }
}`

// TopVulnerableOsTypes returns the most vulnerable OS types.
func (c *Client) TopVulnerableOsTypes(ctx context.Context, filters []Filter, scope *Scope, limit int) ([]OsTypeStats, error) {
	vars := statsVars(filters, scope)
	vars["limit"] = limit
	var resp struct {
		TopVulnerableOsTypes struct {
			OsTypeStats []OsTypeStats `json:"osTypeStats"`
		} `json:"topVulnerableOsTypes"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, topVulnerableOsTypesQuery, vars, &resp); err != nil {
		return nil, err
	}
	return resp.TopVulnerableOsTypes.OsTypeStats, nil
}
