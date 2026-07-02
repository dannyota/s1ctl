package graphql

import (
	"context"
	"encoding/json"
)

// ViewType is a predefined alerts view selector (ScopeSelectorInput viewType).
type ViewType string

const (
	ViewTypeAll          ViewType = "ALL"
	ViewTypeCloud        ViewType = "CLOUD"
	ViewTypeCustomAlerts ViewType = "CUSTOM_ALERTS"
	ViewTypeEndpoint     ViewType = "ENDPOINT"
	ViewTypeIdentity     ViewType = "IDENTITY"
	ViewTypeThirdParty   ViewType = "THIRD_PARTY"
)

// AlertCountValue is a single value bucket in a filters/group-by count result.
type AlertCountValue struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int64  `json:"count"`

	Raw json.RawMessage `json:"-"`
}

func (v *AlertCountValue) UnmarshalJSON(b []byte) error {
	type alias AlertCountValue
	if err := json.Unmarshal(b, (*alias)(v)); err != nil {
		return err
	}
	v.Raw = append(v.Raw[:0:0], b...)
	return nil
}

// AlertFieldCount holds the counted values for a single field. It models both
// FilterCount (alertFiltersCount) and GroupByCount (alertGroupByCount), which
// share the same shape.
type AlertFieldCount struct {
	FieldID     string            `json:"fieldId"`
	HasNextPage bool              `json:"hasNextPage"`
	Values      []AlertCountValue `json:"values"`

	Raw json.RawMessage `json:"-"`
}

func (f *AlertFieldCount) UnmarshalJSON(b []byte) error {
	type alias AlertFieldCount
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// countVars builds the shared variable map for the count queries.
func countVars(fieldIDs []string, filters []Filter, scope *Scope) map[string]any {
	vars := map[string]any{"fieldIds": fieldIDs}
	if len(filters) > 0 {
		vars["filters"] = filters
	}
	if scope != nil {
		vars["scope"] = scope
	}
	return vars
}

const alertFiltersCountQuery = `query AlertFiltersCount($fieldIds: [String!]!, $filters: [FilterInput], $scope: ScopeSelectorInput) {
  alertFiltersCount(fieldIds: $fieldIds, filters: $filters, scope: $scope) {
    data {
      fieldId
      hasNextPage
      values { value label count }
    }
  }
}`

// AlertsFiltersCount returns, per field, the distinct filter values available
// for the current selection along with their alert cardinality.
func (c *Client) AlertsFiltersCount(ctx context.Context, fieldIDs []string, filters []Filter, scope *Scope) ([]AlertFieldCount, error) {
	var resp struct {
		AlertFiltersCount struct {
			Data []AlertFieldCount `json:"data"`
		} `json:"alertFiltersCount"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertFiltersCountQuery, countVars(fieldIDs, filters, scope), &resp); err != nil {
		return nil, err
	}
	return resp.AlertFiltersCount.Data, nil
}

const alertGroupByCountQuery = `query AlertGroupByCount($fieldIds: [String!]!, $filters: [FilterInput], $scope: ScopeSelectorInput) {
  alertGroupByCount(fieldIds: $fieldIds, filters: $filters, scope: $scope) {
    data {
      fieldId
      hasNextPage
      values { value label count }
    }
  }
}`

// AlertsGroupByCount returns alert counts grouped by the specified fields.
//
// The alertGroupByCount query is deprecated server-side in favor of alertGroups
// (exposed via AlertGroups); it is retained here for completeness.
func (c *Client) AlertsGroupByCount(ctx context.Context, fieldIDs []string, filters []Filter, scope *Scope) ([]AlertFieldCount, error) {
	var resp struct {
		AlertGroupByCount struct {
			Data []AlertFieldCount `json:"data"`
		} `json:"alertGroupByCount"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertGroupByCountQuery, countVars(fieldIDs, filters, scope), &resp); err != nil {
		return nil, err
	}
	return resp.AlertGroupByCount.Data, nil
}

const alertsCsvExportQuery = `query AlertsCsvExport($filters: [FilterInput!], $scope: ScopeSelectorInput, $viewType: ViewType) {
  alertsCsvExport(filters: $filters, scope: $scope, viewType: $viewType) {
    data
  }
}`

// AlertsExport downloads alerts matching the filters as CSV.
//
// alertsCsvExport returns the full CSV inline in CsvResponse.data (a String);
// it is not a download URL or an async task. The returned value is the raw CSV
// text.
func (c *Client) AlertsExport(ctx context.Context, filters []Filter, scope *Scope, viewType ViewType) (string, error) {
	vars := map[string]any{}
	if len(filters) > 0 {
		vars["filters"] = filters
	}
	if scope != nil {
		vars["scope"] = scope
	}
	if viewType != "" {
		vars["viewType"] = string(viewType)
	}
	var resp struct {
		AlertsCsvExport struct {
			Data string `json:"data"`
		} `json:"alertsCsvExport"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertsCsvExportQuery, vars, &resp); err != nil {
		return "", err
	}
	return resp.AlertsCsvExport.Data, nil
}
