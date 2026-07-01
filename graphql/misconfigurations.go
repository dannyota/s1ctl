package graphql

import (
	"context"
	"encoding/json"
)

// MisconfigurationCloudInfo holds cloud details for a misconfiguration asset.
type MisconfigurationCloudInfo struct {
	AccountID    string `json:"accountId"`
	AccountName  string `json:"accountName"`
	ProviderName string `json:"providerName"`
	Region       string `json:"region"`
	ResourceID   string `json:"resourceId"`
}

// MisconfigurationAsset is the asset associated with a misconfiguration.
type MisconfigurationAsset struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Category    string                     `json:"category"`
	Subcategory string                     `json:"subcategory"`
	Type        string                     `json:"type"`
	OsType      string                     `json:"osType"`
	CloudInfo   *MisconfigurationCloudInfo `json:"cloudInfo"`
}

// Misconfiguration is an xSPM misconfiguration finding.
type Misconfiguration struct {
	ID                   string                `json:"id"`
	ExternalID           string                `json:"externalId"`
	Name                 string                `json:"name"`
	Description          string                `json:"description"`
	Severity             string                `json:"severity"`
	Status               string                `json:"status"`
	AnalystVerdict       string                `json:"analystVerdict"`
	Product              string                `json:"product"`
	Vendor               string                `json:"vendor"`
	Environment          string                `json:"environment"`
	DetectedAt           string                `json:"detectedAt"`
	LastSeenAt           string                `json:"lastSeenAt"`
	EventTime            string                `json:"eventTime"`
	ResourceUID          string                `json:"resourceUid"`
	MisconfigurationType string                `json:"misconfigurationType"`
	Organization         string                `json:"organization"`
	ComplianceStandards  []string              `json:"complianceStandards"`
	Asset                MisconfigurationAsset `json:"asset"`
	Scope                ScopeInfo             `json:"scope"`

	Raw json.RawMessage `json:"-"`
}

func (m *Misconfiguration) UnmarshalJSON(b []byte) error {
	type alias Misconfiguration
	if err := json.Unmarshal(b, (*alias)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0:0], b...)
	return nil
}

// MisconfigurationEdge is a single edge in a Relay connection.
type MisconfigurationEdge struct {
	Cursor string           `json:"cursor"`
	Node   Misconfiguration `json:"node"`
}

// MisconfigurationConnection is the Relay connection response for misconfigurations.
type MisconfigurationConnection struct {
	Edges      []MisconfigurationEdge `json:"edges"`
	PageInfo   PageInfo               `json:"pageInfo"`
	TotalCount int64                  `json:"totalCount"`
}

// MisconfigurationListParams are parameters for querying misconfigurations.
type MisconfigurationListParams struct {
	First   int      `json:"first,omitempty"`
	After   string   `json:"after,omitempty"`
	Filters []Filter `json:"filters,omitempty"`
	Scope   *Scope   `json:"scope,omitempty"`
}

const misconfigurationsQuery = `query Misconfigurations($first: Int, $after: String, $filters: [FilterInput!], $scope: ScopeSelectorInput) {
  misconfigurations(first: $first, after: $after, filters: $filters, scope: $scope) {
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
        environment
        detectedAt
        lastSeenAt
        eventTime
        resourceUid
        misconfigurationType
        organization
        complianceStandards
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

// MisconfigurationsList queries xSPM misconfigurations.
func (c *Client) MisconfigurationsList(ctx context.Context, params *MisconfigurationListParams) (*MisconfigurationConnection, error) {
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
		Misconfigurations MisconfigurationConnection `json:"misconfigurations"`
	}
	if err := c.Do(ctx, EndpointMisconfigurations, misconfigurationsQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Misconfigurations, nil
}

const misconfigurationGetQuery = `query MisconfigurationGet($id: ID!) {
  misconfiguration(id: $id) {
    id
    externalId
    name
    description
    severity
    status
    analystVerdict
    product
    vendor
    environment
    detectedAt
    lastSeenAt
    eventTime
    resourceUid
    misconfigurationType
    organization
    asset {
      id name category subcategory type osType
      cloudInfo { accountId accountName providerName region resourceId }
    }
    scope { account { id name } site { id name } group { id name } }
  }
}`

// MisconfigurationsGet returns a single misconfiguration by ID.
func (c *Client) MisconfigurationsGet(ctx context.Context, id string) (*Misconfiguration, error) {
	vars := map[string]any{"id": id}
	var resp struct {
		Misconfiguration Misconfiguration `json:"misconfiguration"`
	}
	if err := c.Do(ctx, EndpointMisconfigurations, misconfigurationGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Misconfiguration, nil
}

const misconfigurationsStatusUpdateMutation = `mutation MisconfigurationsStatusUpdate($filter: OrFilterSelectionInput, $statusUpdate: StatusUpdateInput!) {
  misconfigurationsStatusUpdateV2(filter: $filter, statusUpdate: $statusUpdate) {
    updatedFindingIds
  }
}`

// MisconfigurationsUpdateStatus updates the status of the specified misconfigurations.
func (c *Client) MisconfigurationsUpdateStatus(ctx context.Context, ids []string, status string) error {
	vars := map[string]any{
		"filter":       orFilterByIDs(ids),
		"statusUpdate": map[string]any{"status": status},
	}
	return c.Do(ctx, EndpointMisconfigurations, misconfigurationsStatusUpdateMutation, vars, nil)
}

const misconfigurationsVerdictUpdateMutation = `mutation MisconfigurationsVerdictUpdate($filter: OrFilterSelectionInput, $analystVerdict: AnalystVerdict) {
  misconfigurationsAnalystVerdictUpdateV2(filter: $filter, analystVerdict: $analystVerdict) {
    updatedFindingIds
  }
}`

// MisconfigurationsUpdateVerdict updates the analyst verdict of the specified misconfigurations.
func (c *Client) MisconfigurationsUpdateVerdict(ctx context.Context, ids []string, verdict string) error {
	vars := map[string]any{
		"filter":         orFilterByIDs(ids),
		"analystVerdict": verdict,
	}
	return c.Do(ctx, EndpointMisconfigurations, misconfigurationsVerdictUpdateMutation, vars, nil)
}
