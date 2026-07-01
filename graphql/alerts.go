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

	Raw json.RawMessage `json:"-"`
}

func (a *AlertDetectionSource) UnmarshalJSON(b []byte) error {
	type alias AlertDetectionSource
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// AlertAnalytics holds the rule/analytic that triggered an alert.
type AlertAnalytics struct {
	UID  string `json:"uid"`
	Name string `json:"name"`

	Raw json.RawMessage `json:"-"`
}

func (a *AlertAnalytics) UnmarshalJSON(b []byte) error {
	type alias AlertAnalytics
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// AlertAsset identifies the endpoint associated with an alert.
type AlertAsset struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Raw json.RawMessage `json:"-"`
}

func (a *AlertAsset) UnmarshalJSON(b []byte) error {
	type alias AlertAsset
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
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
	Assets          []AlertAsset         `json:"assets"`
	RealTime        struct {
		Scope ScopeInfo `json:"scope"`
	} `json:"realTime"`

	Raw json.RawMessage `json:"-"`
}

// AgentName returns the name of the first asset, or empty string.
func (a *Alert) AgentName() string {
	for _, asset := range a.Assets {
		if asset.Name != "" {
			return asset.Name
		}
	}
	return ""
}

func (a *Alert) UnmarshalJSON(b []byte) error {
	type alias Alert
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

const alertsQuery = `query Alerts($first: Int, $after: String, $filters: [FilterInput!], $scope: ScopeSelectorInput, $sort: SortInput) {
  alerts(first: $first, after: $after, filters: $filters, scope: $scope, sort: $sort) {
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
        assets { id name }
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
    assets { id name }
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
			"id":      "S1/alert/statusUpdate",
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
			"id":      "S1/alert/analystVerdictUpdate",
			"payload": map[string]any{"analystVerdict": map[string]any{"value": verdict}},
		}},
		"filter": orFilterByIDs(ids),
	}
	return c.doAlertTriggerActions(ctx, vars)
}

// AlertsAddNote adds an investigation note to the specified alerts.
func (c *Client) AlertsAddNote(ctx context.Context, ids []string, text string) error {
	vars := map[string]any{
		"actions": []map[string]any{{
			"id":      "S1/alert/addNote",
			"payload": map[string]any{"note": map[string]any{"value": text}},
		}},
		"filter": orFilterByIDs(ids),
	}
	return c.doAlertTriggerActions(ctx, vars)
}

// AlertGroup is a single group-by bucket from the alertGroups query.
type AlertGroup struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int64  `json:"count"`

	Raw json.RawMessage `json:"-"`
}

func (g *AlertGroup) UnmarshalJSON(b []byte) error {
	type alias AlertGroup
	if err := json.Unmarshal(b, (*alias)(g)); err != nil {
		return err
	}
	g.Raw = append(g.Raw[:0:0], b...)
	return nil
}

const alertGroupsQuery = `query AlertGroups($first: Int, $after: String, $filters: [FilterInput!], $scope: ScopeSelectorInput, $groupByFieldId: String!) {
  alertGroups(first: $first, after: $after, filters: $filters, scope: $scope, groupByFieldId: $groupByFieldId) {
    edges {
      cursor
      node {
        value
        label
        count
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

// AlertGroups returns alert counts grouped by the specified field.
func (c *Client) AlertGroups(ctx context.Context, groupByField string, params *ListParams) (*Connection[AlertGroup], error) {
	vars := listVars(params)
	vars["groupByFieldId"] = groupByField
	var resp struct {
		AlertGroups Connection[AlertGroup] `json:"alertGroups"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertGroupsQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.AlertGroups, nil
}

// AlertHistoryItem is a single entry in an alert's audit trail.
type AlertHistoryItem struct {
	CreatedAt  string               `json:"createdAt"`
	EventText  string               `json:"eventText"`
	EventType  string               `json:"eventType"`
	ReportURL  string               `json:"reportUrl"`
	Creator    *AlertHistoryCreator `json:"historyItemCreator"`
	ActionData *AlertHistoryData    `json:"historyItemData"`

	Raw json.RawMessage `json:"-"`
}

func (h *AlertHistoryItem) UnmarshalJSON(b []byte) error {
	type alias AlertHistoryItem
	if err := json.Unmarshal(b, (*alias)(h)); err != nil {
		return err
	}
	h.Raw = append(h.Raw[:0:0], b...)
	return nil
}

func (h *AlertHistoryItem) ActorName() string {
	if h.Creator != nil {
		return h.Creator.UserID
	}
	return ""
}

type AlertHistoryCreator struct {
	UserID   string `json:"userId"`
	UserType string `json:"userType"`

	Raw json.RawMessage `json:"-"`
}

func (a *AlertHistoryCreator) UnmarshalJSON(b []byte) error {
	type alias AlertHistoryCreator
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

type AlertHistoryData struct {
	Message     *AlertHistoryText `json:"message"`
	Description *AlertHistoryText `json:"description"`

	Raw json.RawMessage `json:"-"`
}

func (a *AlertHistoryData) UnmarshalJSON(b []byte) error {
	type alias AlertHistoryData
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

type AlertHistoryText struct {
	Content string `json:"content"`
	Type    string `json:"type"`

	Raw json.RawMessage `json:"-"`
}

func (a *AlertHistoryText) UnmarshalJSON(b []byte) error {
	type alias AlertHistoryText
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

const alertHistoryQuery = `query AlertHistory($alertId: ID!, $first: Int, $after: String, $filter: AlertHistoryFilterInput) {
  alertHistory(alertId: $alertId, first: $first, after: $after, filter: $filter) {
    edges {
      cursor
      node {
        createdAt
        eventText
        eventType
        reportUrl
        historyItemCreator {
          ... on UserHistoryItemCreator { userId userType }
        }
        historyItemData {
          ... on MitigationActionHistoryItemData {
            message { content type }
          }
          ... on EnrichmentHistoryItemData {
            description { content type }
          }
        }
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

func (c *Client) AlertHistory(ctx context.Context, alertID string, first int, after string) (*Connection[AlertHistoryItem], error) {
	vars := map[string]any{"alertId": alertID}
	if first > 0 {
		vars["first"] = first
	}
	if after != "" {
		vars["after"] = after
	}
	var resp struct {
		AlertHistory Connection[AlertHistoryItem] `json:"alertHistory"`
	}
	if err := c.Do(ctx, EndpointAlerts, alertHistoryQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.AlertHistory, nil
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
