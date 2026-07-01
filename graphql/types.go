package graphql

// ScopeEntity is a single scope level (account, site, or group) in responses.
type ScopeEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ScopeInfo is scope information returned in API responses.
type ScopeInfo struct {
	Account ScopeEntity `json:"account"`
	Site    ScopeEntity `json:"site"`
	Group   ScopeEntity `json:"group"`
}

// CloudInfo holds cloud details for an asset.
type CloudInfo struct {
	AccountID    string `json:"accountId"`
	AccountName  string `json:"accountName"`
	ProviderName string `json:"providerName"`
	Region       string `json:"region"`
	ResourceID   string `json:"resourceId"`
}

// Asset is an asset associated with an xSPM finding (misconfiguration or vulnerability).
type Asset struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	Subcategory string     `json:"subcategory"`
	Type        string     `json:"type"`
	OsType      string     `json:"osType"`
	CloudInfo   *CloudInfo `json:"cloudInfo"`
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

// PageInfo is Relay-style pagination info.
type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	EndCursor       string `json:"endCursor"`
	StartCursor     string `json:"startCursor"`
}

// Edge is a single edge in a Relay connection.
type Edge[T any] struct {
	Cursor string `json:"cursor"`
	Node   T      `json:"node"`
}

// Connection is a Relay connection response.
type Connection[T any] struct {
	Edges      []Edge[T] `json:"edges"`
	PageInfo   PageInfo  `json:"pageInfo"`
	TotalCount int64     `json:"totalCount"`
}

// ListParams are parameters for paginated GraphQL list queries.
type ListParams struct {
	First   int      `json:"first,omitempty"`
	After   string   `json:"after,omitempty"`
	Filters []Filter `json:"filters,omitempty"`
	Scope   *Scope   `json:"scope,omitempty"`
}

// listVars builds the GraphQL variables map from list parameters.
func listVars(p *ListParams) map[string]any {
	vars := map[string]any{}
	if p == nil {
		return vars
	}
	if p.First > 0 {
		vars["first"] = p.First
	}
	if p.After != "" {
		vars["after"] = p.After
	}
	if len(p.Filters) > 0 {
		vars["filters"] = p.Filters
	}
	if p.Scope != nil {
		vars["scope"] = p.Scope
	}
	return vars
}

// orFilterByIDs builds an OrFilterSelectionInput that matches any of the given IDs.
func orFilterByIDs(ids []string) map[string]any {
	return map[string]any{
		"or": []map[string]any{{
			"and": []map[string]any{{
				"fieldId":  "id",
				"stringIn": map[string]any{"values": ids},
			}},
		}},
	}
}

// Backward-compatible type aliases for domain-specific Relay types.
type (
	AlertEdge        = Edge[Alert]
	AlertConnection  = Connection[Alert]
	AlertsListParams = ListParams

	MisconfigurationCloudInfo  = CloudInfo
	MisconfigurationAsset      = Asset
	MisconfigurationEdge       = Edge[Misconfiguration]
	MisconfigurationConnection = Connection[Misconfiguration]
	MisconfigurationListParams = ListParams

	VulnerabilityAsset      = Asset
	VulnerabilityEdge       = Edge[Vulnerability]
	VulnerabilityConnection = Connection[Vulnerability]
	VulnerabilityListParams = ListParams

	CloudPolicyEdge       = Edge[CloudPolicy]
	CloudPolicyConnection = Connection[CloudPolicy]
	CloudPolicyListParams = ListParams
)
