package graphql

import (
	"context"
	"encoding/json"
)

// Alert is a UAM unified alert.
type Alert struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Severity        string `json:"severity"`
	Classification  string `json:"classification"`
	ConfidenceLevel string `json:"confidenceLevel"`
	AnalystVerdict  string `json:"analystVerdict"`
	Status          string `json:"status"`
	DetectedAt      string `json:"detectedAt"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
	StorylineID     string `json:"storylineId"`

	Raw json.RawMessage `json:"-"`
}

func (a *Alert) UnmarshalJSON(b []byte) error {
	type alias Alert
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// PageInfo is Relay-style pagination info.
type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	EndCursor       string `json:"endCursor"`
	StartCursor     string `json:"startCursor"`
}

// AlertEdge is a single edge in a Relay connection.
type AlertEdge struct {
	Cursor string `json:"cursor"`
	Node   Alert  `json:"node"`
}

// AlertConnection is the Relay connection response.
type AlertConnection struct {
	Edges      []AlertEdge `json:"edges"`
	PageInfo   PageInfo    `json:"pageInfo"`
	TotalCount int64       `json:"totalCount"`
}

// Filter is a GraphQL filter input.
type Filter struct {
	FieldID     string `json:"fieldId"`
	StringIn    *InStr `json:"stringIn,omitempty"`
	StringEqual *EqStr `json:"stringEqual,omitempty"`
	IsNegated   bool   `json:"isNegated,omitempty"`
}

// InStr is a string "in" filter.
type InStr struct {
	Values []string `json:"values"`
}

// EqStr is a string "equal" filter.
type EqStr struct {
	Value string `json:"value"`
}

// Scope specifies the scope selector.
type Scope struct {
	ScopeIDs  []string `json:"scopeIds"`
	ScopeType string   `json:"scopeType"`
}

// AlertsListParams are parameters for querying alerts.
type AlertsListParams struct {
	First   int      `json:"first,omitempty"`
	After   string   `json:"after,omitempty"`
	Filters []Filter `json:"filters,omitempty"`
	Scope   *Scope   `json:"scope,omitempty"`
}

const alertsQuery = `query Alerts($first: Int, $after: String, $filters: [FilterInput!], $scope: ScopeSelectorInput) {
  alerts(first: $first, after: $after, filters: $filters, scope: $scope) {
    edges {
      cursor
      node {
        id
        name
        description
        severity
        classification
        confidenceLevel
        analystVerdict
        status
        detectedAt
        createdAt
        updatedAt
        storylineId
      }
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      endCursor
      startCursor
    }
    totalCount
  }
}`

// AlertsList queries UAM alerts.
func (c *Client) AlertsList(ctx context.Context, params *AlertsListParams) (*AlertConnection, error) {
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
		Alerts AlertConnection `json:"alerts"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertsQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Alerts, nil
}
