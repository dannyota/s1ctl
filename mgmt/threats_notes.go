package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ThreatNote is a note attached to a threat.
type ThreatNote struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Creator   string `json:"creator"`
	CreatorID string `json:"creatorId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Edited    bool   `json:"edited"`

	Raw json.RawMessage `json:"-"`
}

func (n ThreatNote) MarshalJSON() ([]byte, error) {
	if n.Raw != nil {
		return n.Raw, nil
	}
	type alias ThreatNote
	return json.Marshal(alias(n))
}

func (n *ThreatNote) UnmarshalJSON(b []byte) error {
	type alias ThreatNote
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// ThreatNotesListParams are query parameters for listing threat notes.
type ThreatNotesListParams struct {
	Limit     int
	Cursor    string
	SortBy    string
	SortOrder string
}

func (p *ThreatNotesListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// ThreatNotesList returns notes for a threat.
func (c *Client) ThreatNotesList(ctx context.Context, threatID string, params *ThreatNotesListParams) ([]ThreatNote, *Pagination, error) {
	if threatID == "" {
		return nil, nil, fmt.Errorf("mgmt: threat ID is required")
	}
	path := fmt.Sprintf("/threats/%s/notes", threatID)
	return list[ThreatNote](c, ctx, path, params.values())
}

// ThreatNotesCreate adds a note to one or more threats.
func (c *Client) ThreatNotesCreate(ctx context.Context, threatID, text string) (int, error) {
	if threatID == "" {
		return 0, fmt.Errorf("mgmt: threat ID is required")
	}
	if text == "" {
		return 0, fmt.Errorf("mgmt: note text is required")
	}
	return doAction(c, ctx, "/threats/notes", ActionFilter{IDs: []string{threatID}}, map[string]string{"text": text})
}
