package graphql

import (
	"context"
	"encoding/json"
)

// MisconfigurationNote is an investigation note attached to a misconfiguration.
type MisconfigurationNote struct {
	ID                 string    `json:"id"`
	MisconfigurationID string    `json:"misconfigurationId"`
	Text               string    `json:"text"`
	CreatedAt          string    `json:"createdAt"`
	UpdatedAt          string    `json:"updatedAt"`
	Author             *XSPMUser `json:"author"`

	Raw json.RawMessage `json:"-"`
}

// AuthorName returns the note author's full name, or empty string.
func (n *MisconfigurationNote) AuthorName() string {
	if n.Author != nil {
		return n.Author.FullName
	}
	return ""
}

func (n *MisconfigurationNote) UnmarshalJSON(b []byte) error {
	type alias MisconfigurationNote
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// MisconfigurationRelatedAsset is an asset affected by the same misconfiguration.
type MisconfigurationRelatedAsset struct {
	MisconfigurationID string `json:"misconfigurationId"`
	Organization       string `json:"organization"`
	Asset              Asset  `json:"asset"`

	Raw json.RawMessage `json:"-"`
}

func (a *MisconfigurationRelatedAsset) UnmarshalJSON(b []byte) error {
	type alias MisconfigurationRelatedAsset
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

const misconfigurationNotesQuery = `query MisconfigurationNotes($misconfigurationId: ID!, $first: Int) {
  misconfigurationNotes(misconfigurationId: $misconfigurationId, first: $first) {
    edges {
      node {
        id
        misconfigurationId
        text
        createdAt
        updatedAt
        author { id fullName email deleted }
      }
    }
    totalCount
  }
}`

// MisconfigurationsNotes returns the investigation notes on a misconfiguration.
func (c *Client) MisconfigurationsNotes(ctx context.Context, id string) ([]MisconfigurationNote, error) {
	vars := map[string]any{"misconfigurationId": id, "first": defaultXSPMPageSize}
	var resp struct {
		MisconfigurationNotes Connection[MisconfigurationNote] `json:"misconfigurationNotes"`
	}
	if err := c.Do(ctx, EndpointMisconfigurations, misconfigurationNotesQuery, vars, &resp); err != nil {
		return nil, err
	}
	return edgeNodes(resp.MisconfigurationNotes), nil
}

const misconfigurationAddNoteMutation = `mutation MisconfigurationAddNote($filter: OrFilterSelectionInput, $text: String!) {
  addMisconfigurationNoteV2(filter: $filter, text: $text) {
    updatedFindingIds
  }
}`

// MisconfigurationsAddNote adds an investigation note to the specified
// misconfigurations.
//
// It uses addMisconfigurationNoteV2; the non-V2 variant is deprecated. The
// includeHidden and viewType arguments are left at their server defaults.
func (c *Client) MisconfigurationsAddNote(ctx context.Context, ids []string, text string) error {
	vars := map[string]any{"filter": orFilterByIDs(ids), "text": text}
	return c.Do(ctx, EndpointMisconfigurations, misconfigurationAddNoteMutation, vars, nil)
}

const misconfigurationUpdateNoteMutation = `mutation MisconfigurationUpdateNote($noteId: ID!, $text: String!) {
  updateMisconfigurationNote(noteId: $noteId, text: $text)
}`

// MisconfigurationsUpdateNote updates the text of an existing misconfiguration note.
func (c *Client) MisconfigurationsUpdateNote(ctx context.Context, noteID, text string) error {
	vars := map[string]any{"noteId": noteID, "text": text}
	return c.Do(ctx, EndpointMisconfigurations, misconfigurationUpdateNoteMutation, vars, nil)
}

const misconfigurationDeleteNoteMutation = `mutation MisconfigurationDeleteNote($noteId: ID!) {
  deleteMisconfigurationNote(noteId: $noteId)
}`

// MisconfigurationsDeleteNote deletes a misconfiguration note.
func (c *Client) MisconfigurationsDeleteNote(ctx context.Context, noteID string) error {
	vars := map[string]any{"noteId": noteID}
	return c.Do(ctx, EndpointMisconfigurations, misconfigurationDeleteNoteMutation, vars, nil)
}

const misconfigurationAssignMutation = `mutation MisconfigurationAssign($filter: OrFilterSelectionInput, $userId: ID) {
  misconfigurationUserAssignmentV2(filter: $filter, userId: $userId) {
    updatedFindingIds
  }
}`

// MisconfigurationsAssign assigns the specified misconfigurations to a user. An
// empty userID unassigns them.
//
// It uses misconfigurationUserAssignmentV2; the non-V2 variant is deprecated.
func (c *Client) MisconfigurationsAssign(ctx context.Context, ids []string, userID string) error {
	vars := map[string]any{"filter": orFilterByIDs(ids)}
	if userID != "" {
		vars["userId"] = userID
	}
	return c.Do(ctx, EndpointMisconfigurations, misconfigurationAssignMutation, vars, nil)
}

const misconfigurationHistoryQuery = `query MisconfigurationHistory($misconfigurationId: ID!, $first: Int) {
  misconfigurationHistory(misconfigurationId: $misconfigurationId, first: $first) {
    edges {
      node { createdAt eventText eventType }
    }
    totalCount
  }
}`

// MisconfigurationsHistory returns the history records of a misconfiguration.
func (c *Client) MisconfigurationsHistory(ctx context.Context, id string) ([]XSPMHistoryItem, error) {
	vars := map[string]any{"misconfigurationId": id, "first": defaultXSPMPageSize}
	var resp struct {
		MisconfigurationHistory Connection[XSPMHistoryItem] `json:"misconfigurationHistory"`
	}
	if err := c.Do(ctx, EndpointMisconfigurations, misconfigurationHistoryQuery, vars, &resp); err != nil {
		return nil, err
	}
	return edgeNodes(resp.MisconfigurationHistory), nil
}

const misconfigurationRelatedAssetsQuery = `query MisconfigurationRelatedAssets($filters: [FilterInput!], $first: Int) {
  misconfigurationRelatedAssets(filters: $filters, first: $first) {
    edges {
      node {
        misconfigurationId
        organization
        asset { ` + xspmAssetFields + ` }
      }
    }
    totalCount
  }
}`

// MisconfigurationsRelatedAssets returns the assets related to a misconfiguration.
//
// The schema query selects related assets by FilterInput; the finding ID is
// passed as a filter on fieldId "id" (the deprecated name argument is not used).
func (c *Client) MisconfigurationsRelatedAssets(ctx context.Context, id string) ([]MisconfigurationRelatedAsset, error) {
	vars := map[string]any{
		"filters": []Filter{{FieldID: "id", StringEqual: &EqStr{Value: id}}},
		"first":   defaultXSPMPageSize,
	}
	var resp struct {
		MisconfigurationRelatedAssets Connection[MisconfigurationRelatedAsset] `json:"misconfigurationRelatedAssets"`
	}
	if err := c.Do(ctx, EndpointMisconfigurations, misconfigurationRelatedAssetsQuery, vars, &resp); err != nil {
		return nil, err
	}
	return edgeNodes(resp.MisconfigurationRelatedAssets), nil
}

const misconfigurationsExportQuery = `query MisconfigurationsExport($filters: [FilterInput], $scope: ScopeSelectorInput) {
  misconfigurationsExportToCsv(filters: $filters, scope: $scope) {
    data
  }
}`

// MisconfigurationsExport downloads misconfigurations matching the filters as CSV.
//
// misconfigurationsExportToCsv returns the full CSV inline in CsvResponse.data.
func (c *Client) MisconfigurationsExport(ctx context.Context, filters []Filter, scope *Scope) (string, error) {
	vars := map[string]any{}
	if len(filters) > 0 {
		vars["filters"] = filters
	}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		MisconfigurationsExportToCsv struct {
			Data string `json:"data"`
		} `json:"misconfigurationsExportToCsv"`
	}
	if err := c.Do(ctx, EndpointMisconfigurations, misconfigurationsExportQuery, vars, &resp); err != nil {
		return "", err
	}
	return resp.MisconfigurationsExportToCsv.Data, nil
}
