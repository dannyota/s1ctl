package graphql

import (
	"context"
	"encoding/json"
	"fmt"
)

// AlertDetectionSource identifies the product and vendor that detected an alert.
type AlertDetectionSource struct {
	Product string `json:"product"`
	Vendor  string `json:"vendor"`
}

// AlertAnalytics holds the rule/analytic that triggered an alert.
type AlertAnalytics struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

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

	DetectionSource AlertDetectionSource `json:"detectionSource"`
	Analytics       *AlertAnalytics      `json:"analytics"`
	RealTime        struct {
		Scope ScopeInfo `json:"scope"`
	} `json:"realTime"`

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
        detectionSource { product vendor }
        analytics { uid name }
        realTime { scope { account { id name } site { id name } } }
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
func (c *Client) AlertsList(ctx context.Context, params *ListParams) (*Connection[Alert], error) {
	var resp struct {
		Alerts Connection[Alert] `json:"alerts"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertsQuery, listVars(params), &resp); err != nil {
		return nil, err
	}
	return &resp.Alerts, nil
}

const alertGetQuery = `query AlertGet($id: ID!) {
  alert(id: $id) {
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
    detectionSource { product vendor }
    analytics { uid name }
    realTime { scope { account { id name } site { id name } } }
  }
}`

// AlertsGet returns a single alert by ID.
func (c *Client) AlertsGet(ctx context.Context, id string) (*Alert, error) {
	vars := map[string]any{"id": id}
	var resp struct {
		Alert Alert `json:"alert"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Alert, nil
}

const alertTriggerActionsMutation = `mutation AlertTriggerActions($actions: [TriggerActionInput!]!, $filter: OrFilterSelectionInput) {
  alertTriggerActions(actions: $actions, filter: $filter) {
    ... on ActionsTriggered { actions { actionId } }
    ... on TriggerActionsError { errors { errorMessage } }
  }
}`

// AlertsUpdateStatus updates the investigation status of the specified alerts.
func (c *Client) AlertsUpdateStatus(ctx context.Context, ids []string, status string) error {
	vars := map[string]any{
		"actions": []map[string]any{{
			"id":      "status",
			"payload": map[string]any{"status": map[string]any{"value": status}},
		}},
		"filter": orFilterByIDs(ids),
	}
	return c.doAlertTriggerActions(ctx, vars)
}

// AlertsUpdateVerdict updates the analyst verdict of the specified alerts.
func (c *Client) AlertsUpdateVerdict(ctx context.Context, ids []string, verdict string) error {
	vars := map[string]any{
		"actions": []map[string]any{{
			"id":      "analystVerdict",
			"payload": map[string]any{"analystVerdict": map[string]any{"value": verdict}},
		}},
		"filter": orFilterByIDs(ids),
	}
	return c.doAlertTriggerActions(ctx, vars)
}

func (c *Client) doAlertTriggerActions(ctx context.Context, vars map[string]any) error {
	var resp struct {
		AlertTriggerActions struct {
			Errors []struct {
				ErrorMessage string `json:"errorMessage"`
			} `json:"errors"`
		} `json:"alertTriggerActions"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertTriggerActionsMutation, vars, &resp); err != nil {
		return err
	}
	if len(resp.AlertTriggerActions.Errors) > 0 {
		return fmt.Errorf("graphql: %s", resp.AlertTriggerActions.Errors[0].ErrorMessage)
	}
	return nil
}
