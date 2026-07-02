package sdl

import (
	"context"
	"encoding/json"
	"fmt"
)

const dashboardsListGQL = `query dashboards {
  dashboardsV2 {
    id
    name
    isBuiltIn
    isEditable
  }
}`

const dashboardGetGQL = `query getDashboard($id: ID) {
  getDashboardV2(id: $id, resolveParameters: true) {
    id
    name
    urlTemplate
    description
    configType
    duration
    access { public users owner }
    tabs { parameters graphs filters options tabName }
  }
}`

const savedSearchDeleteGQL = `mutation deleteSavedSearch($name: String!, $type: SavedSearchType!, $index: Int!) {
  deleteSavedSearchV2(name: $name, type: $type, index: $index) {
    name url index type
  }
}`

// Dashboard is an SDL Data Lake dashboard.
type Dashboard struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsBuiltIn  bool   `json:"isBuiltIn"`
	IsEditable bool   `json:"isEditable"`

	Raw json.RawMessage `json:"-"`
}

func (d *Dashboard) UnmarshalJSON(b []byte) error {
	type alias Dashboard
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// DashboardDetail is the full representation of an SDL dashboard including
// layout, access control, and tab configuration.
type DashboardDetail struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URLTemplate string `json:"urlTemplate"`
	Description string `json:"description"`
	ConfigType  string `json:"configType"`
	Duration    string `json:"duration"`

	Access json.RawMessage `json:"access"`
	Tabs   json.RawMessage `json:"tabs"`

	Raw json.RawMessage `json:"-"`
}

func (d *DashboardDetail) UnmarshalJSON(b []byte) error {
	type alias DashboardDetail
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// DashboardsList returns all dashboards from the SDL console.
func (c *Client) DashboardsList(ctx context.Context) ([]Dashboard, error) {
	var data struct {
		Dashboards []Dashboard `json:"dashboardsV2"`
	}
	if err := c.graphql(ctx, dashboardsListGQL, nil, &data); err != nil {
		return nil, err
	}
	return data.Dashboards, nil
}

// DashboardGet returns a single dashboard by ID with full detail.
func (c *Client) DashboardGet(ctx context.Context, id string) (*DashboardDetail, error) {
	vars := map[string]any{"id": id}
	var data struct {
		Dashboard DashboardDetail `json:"getDashboardV2"`
	}
	if err := c.graphql(ctx, dashboardGetGQL, vars, &data); err != nil {
		return nil, err
	}
	return &data.Dashboard, nil
}

// SavedSearchType is the type of a saved search (PRIVATE or SHARED).
type SavedSearchType string

const (
	SavedSearchTypePrivate SavedSearchType = "PRIVATE"
	SavedSearchTypeShared  SavedSearchType = "SHARED"
)

// SavedSearchDelete deletes a saved search by name, type, and index.
func (c *Client) SavedSearchDelete(ctx context.Context, name string, searchType SavedSearchType, index int) error {
	if name == "" {
		return fmt.Errorf("sdl: saved search name is required")
	}
	vars := map[string]any{
		"name":  name,
		"type":  string(searchType),
		"index": index,
	}
	return c.graphql(ctx, savedSearchDeleteGQL, vars, nil)
}
