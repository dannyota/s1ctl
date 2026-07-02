package graphql

import (
	"context"
	"encoding/json"
)

// VulnerabilityNote is an investigation note attached to a vulnerability.
type VulnerabilityNote struct {
	ID              string    `json:"id"`
	VulnerabilityID string    `json:"vulnerabilityId"`
	Text            string    `json:"text"`
	CreatedAt       string    `json:"createdAt"`
	UpdatedAt       string    `json:"updatedAt"`
	Author          *XSPMUser `json:"author"`

	Raw json.RawMessage `json:"-"`
}

// AuthorName returns the note author's full name, or empty string.
func (n *VulnerabilityNote) AuthorName() string {
	if n.Author != nil {
		return n.Author.FullName
	}
	return ""
}

func (n *VulnerabilityNote) UnmarshalJSON(b []byte) error {
	type alias VulnerabilityNote
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// VulnerabilityRelatedAsset is an asset affected by the same vulnerability.
type VulnerabilityRelatedAsset struct {
	VulnerabilityID string                `json:"vulnerabilityId"`
	Asset           Asset                 `json:"asset"`
	Software        VulnerabilitySoftware `json:"software"`

	Raw json.RawMessage `json:"-"`
}

func (a *VulnerabilityRelatedAsset) UnmarshalJSON(b []byte) error {
	type alias VulnerabilityRelatedAsset
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

const vulnerabilityNotesQuery = `query VulnerabilityNotes($vulnerabilityId: ID!, $first: Int) {
  vulnerabilityNotes(vulnerabilityId: $vulnerabilityId, first: $first) {
    edges {
      node {
        id
        vulnerabilityId
        text
        createdAt
        updatedAt
        author { id fullName email deleted }
      }
    }
    totalCount
  }
}`

// VulnerabilitiesNotes returns the investigation notes on a vulnerability.
func (c *Client) VulnerabilitiesNotes(ctx context.Context, id string) ([]VulnerabilityNote, error) {
	vars := map[string]any{"vulnerabilityId": id, "first": defaultXSPMPageSize}
	var resp struct {
		VulnerabilityNotes Connection[VulnerabilityNote] `json:"vulnerabilityNotes"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, vulnerabilityNotesQuery, vars, &resp); err != nil {
		return nil, err
	}
	return edgeNodes(resp.VulnerabilityNotes), nil
}

const vulnerabilityAddNoteMutation = `mutation VulnerabilityAddNote($filter: OrFilterSelectionInput, $text: String!) {
  addVulnerabilityNoteV2(filter: $filter, text: $text) {
    updatedFindingIds
  }
}`

// VulnerabilitiesAddNote adds an investigation note to the specified
// vulnerabilities.
//
// It uses addVulnerabilityNoteV2; the non-V2 variant is deprecated. The
// includeHidden argument is left at its server default.
func (c *Client) VulnerabilitiesAddNote(ctx context.Context, ids []string, text string) error {
	vars := map[string]any{"filter": orFilterByIDs(ids), "text": text}
	return c.Do(ctx, EndpointVulnerabilities, vulnerabilityAddNoteMutation, vars, nil)
}

const vulnerabilityUpdateNoteMutation = `mutation VulnerabilityUpdateNote($noteId: ID!, $text: String!) {
  updateVulnerabilityNote(noteId: $noteId, text: $text)
}`

// VulnerabilitiesUpdateNote updates the text of an existing vulnerability note.
func (c *Client) VulnerabilitiesUpdateNote(ctx context.Context, noteID, text string) error {
	vars := map[string]any{"noteId": noteID, "text": text}
	return c.Do(ctx, EndpointVulnerabilities, vulnerabilityUpdateNoteMutation, vars, nil)
}

const vulnerabilityDeleteNoteMutation = `mutation VulnerabilityDeleteNote($noteId: ID!) {
  deleteVulnerabilityNote(noteId: $noteId)
}`

// VulnerabilitiesDeleteNote deletes a vulnerability note.
func (c *Client) VulnerabilitiesDeleteNote(ctx context.Context, noteID string) error {
	vars := map[string]any{"noteId": noteID}
	return c.Do(ctx, EndpointVulnerabilities, vulnerabilityDeleteNoteMutation, vars, nil)
}

const vulnerabilityAssignMutation = `mutation VulnerabilityAssign($filter: OrFilterSelectionInput, $userId: ID) {
  vulnerabilityUserAssignmentV2(filter: $filter, userId: $userId) {
    updatedFindingIds
  }
}`

// VulnerabilitiesAssign assigns the specified vulnerabilities to a user. An
// empty userID unassigns them.
//
// It uses vulnerabilityUserAssignmentV2; the non-V2 variant is deprecated.
func (c *Client) VulnerabilitiesAssign(ctx context.Context, ids []string, userID string) error {
	vars := map[string]any{"filter": orFilterByIDs(ids)}
	if userID != "" {
		vars["userId"] = userID
	}
	return c.Do(ctx, EndpointVulnerabilities, vulnerabilityAssignMutation, vars, nil)
}

const vulnerabilityHistoryQuery = `query VulnerabilityHistory($vulnerabilityId: ID!, $first: Int) {
  vulnerabilityHistory(vulnerabilityId: $vulnerabilityId, first: $first) {
    edges {
      node { createdAt eventText eventType }
    }
    totalCount
  }
}`

// VulnerabilitiesHistory returns the history records of a vulnerability.
func (c *Client) VulnerabilitiesHistory(ctx context.Context, id string) ([]XSPMHistoryItem, error) {
	vars := map[string]any{"vulnerabilityId": id, "first": defaultXSPMPageSize}
	var resp struct {
		VulnerabilityHistory Connection[XSPMHistoryItem] `json:"vulnerabilityHistory"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, vulnerabilityHistoryQuery, vars, &resp); err != nil {
		return nil, err
	}
	return edgeNodes(resp.VulnerabilityHistory), nil
}

const vulnerabilityRelatedAssetsQuery = `query VulnerabilityRelatedAssets($cveId: String, $first: Int) {
  vulnerabilityRelatedAssets(cveId: $cveId, first: $first) {
    edges {
      node {
        vulnerabilityId
        asset { ` + xspmAssetFields + ` }
        software { name version fixVersion packageManager type vendor }
      }
    }
    totalCount
  }
}`

// VulnerabilitiesRelatedAssets returns the assets related to a vulnerability.
//
// The finding ID is passed via the dedicated cveId argument, which targets the
// vulnerability directly and unambiguously. The schema marks cveId as
// deprecated-in-favor-of-filters, yet it is preferred here over a guessed
// fieldId "id" filter: an unrecognized filter field could silently return
// every related asset across all findings, whereas cveId cannot mis-target.
func (c *Client) VulnerabilitiesRelatedAssets(ctx context.Context, id string) ([]VulnerabilityRelatedAsset, error) {
	vars := map[string]any{
		"cveId": id,
		"first": defaultXSPMPageSize,
	}
	var resp struct {
		VulnerabilityRelatedAssets Connection[VulnerabilityRelatedAsset] `json:"vulnerabilityRelatedAssets"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, vulnerabilityRelatedAssetsQuery, vars, &resp); err != nil {
		return nil, err
	}
	return edgeNodes(resp.VulnerabilityRelatedAssets), nil
}

const vulnerabilitiesExportQuery = `query VulnerabilitiesExport($filters: [FilterInput], $scope: ScopeSelectorInput) {
  vulnerabilitiesExportToCsv(filters: $filters, scope: $scope) {
    data
  }
}`

// VulnerabilitiesExport downloads vulnerabilities matching the filters as CSV.
//
// vulnerabilitiesExportToCsv returns the full CSV inline in CsvResponse.data.
func (c *Client) VulnerabilitiesExport(ctx context.Context, filters []Filter, scope *Scope) (string, error) {
	vars := map[string]any{}
	if len(filters) > 0 {
		vars["filters"] = filters
	}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		VulnerabilitiesExportToCsv struct {
			Data string `json:"data"`
		} `json:"vulnerabilitiesExportToCsv"`
	}
	if err := c.Do(ctx, EndpointVulnerabilities, vulnerabilitiesExportQuery, vars, &resp); err != nil {
		return "", err
	}
	return resp.VulnerabilitiesExportToCsv.Data, nil
}
