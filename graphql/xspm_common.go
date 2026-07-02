package graphql

import "encoding/json"

// defaultXSPMPageSize is the page size used for the single-page xSPM note,
// history, and related-asset queries.
const defaultXSPMPageSize = 100

// edgeNodes flattens a Relay connection into a slice of its nodes.
func edgeNodes[T any](conn Connection[T]) []T {
	out := make([]T, 0, len(conn.Edges))
	for _, e := range conn.Edges {
		out = append(out, e.Node)
	}
	return out
}

// XSPMUser is a user referenced by an xSPM note author or assignment. It models
// the shared User type in both the misconfigurations and vulnerabilities schemas.
type XSPMUser struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Deleted  bool   `json:"deleted"`

	Raw json.RawMessage `json:"-"`
}

func (u *XSPMUser) UnmarshalJSON(b []byte) error {
	type alias XSPMUser
	if err := json.Unmarshal(b, (*alias)(u)); err != nil {
		return err
	}
	u.Raw = append(u.Raw[:0:0], b...)
	return nil
}

// XSPMHistoryItem is a single entry in an xSPM finding's history. It models the
// shared HistoryItem shape in both the misconfigurations and vulnerabilities
// schemas (createdAt, eventText, eventType).
type XSPMHistoryItem struct {
	CreatedAt string `json:"createdAt"`
	EventText string `json:"eventText"`
	EventType string `json:"eventType"`

	Raw json.RawMessage `json:"-"`
}

func (h *XSPMHistoryItem) UnmarshalJSON(b []byte) error {
	type alias XSPMHistoryItem
	if err := json.Unmarshal(b, (*alias)(h)); err != nil {
		return err
	}
	h.Raw = append(h.Raw[:0:0], b...)
	return nil
}

// xspmAssetFields is the GraphQL selection for an Asset, shared by the related
// assets queries. It matches the Asset fields selected by the list/get queries.
const xspmAssetFields = `id name category subcategory type osType
      cloudInfo { accountId accountName providerName region resourceId }`
