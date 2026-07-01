package sdl

import "encoding/json"

// QueryType identifies the kind of SDL query.
type QueryType string

const (
	QueryTypePQ        QueryType = "PQ"
	QueryTypeTopFacets QueryType = "TOP_FACETS"
	QueryTypeLog       QueryType = "LOG"
	QueryTypeGraph     QueryType = "GRAPH"
)

// QueryStatus represents the execution state of a query group.
type QueryStatus string

const (
	QueryStatusRunning QueryStatus = "RUNNING"
	QueryStatusDone    QueryStatus = "DONE"
	QueryStatusError   QueryStatus = "ERROR"
)

// QueryGroupRequest is the input for launchQuery.
type QueryGroupRequest struct {
	Queries   []QueryRequest `json:"queries"`
	PreFilter string         `json:"preFilter,omitempty"`
}

// QueryRequest describes a single query within a group.
type QueryRequest struct {
	ID         string      `json:"id"`
	Type       QueryType   `json:"type"`
	Filter     string      `json:"filter"`
	StartTime  string      `json:"startTime"`
	EndTime    string      `json:"endTime"`
	Origin     string      `json:"origin,omitempty"`
	Tenant     *bool       `json:"tenant,omitempty"`
	PowerQuery *struct{}   `json:"powerQuery,omitempty"`
	FacetQuery *FacetQuery `json:"facetQuery,omitempty"`
}

// FacetQuery configures TOP_FACETS queries.
type FacetQuery struct {
	DetermineNumeric          *bool `json:"determineNumeric,omitempty"`
	IncludeSingleValueFacets  *bool `json:"includeSingleValueFacets,omitempty"`
	NumFacetsToReturn         *int  `json:"numFacetsToReturn,omitempty"`
	NumValuesToReturnPerFacet *int  `json:"numValuesToReturnPerFacet,omitempty"`
}

// QueriesResult is the response from launchQuery/pingQuery.
type QueriesResult struct {
	IDs            []string      `json:"ids"`
	Status         QueryStatus   `json:"status"`
	Token          string        `json:"token"`
	StepsCompleted int           `json:"stepsCompleted"`
	TotalSteps     int           `json:"totalSteps"`
	Results        []QueryResult `json:"results"`

	Raw json.RawMessage `json:"-"`
}

func (r *QueriesResult) UnmarshalJSON(b []byte) error {
	type alias QueriesResult
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// QueryResult holds a single query's result within a group.
type QueryResult struct {
	ID             string          `json:"id"`
	StepsCompleted int             `json:"stepsCompleted"`
	TotalSteps     int             `json:"totalSteps"`
	Error          string          `json:"error"`
	CacheContext   string          `json:"cacheContext"`
	NoResultReason string          `json:"noResultsReason"`
	Data           json.RawMessage `json:"data"`
}

// PQResultData is the typed result for PowerQuery (type: PQ).
// Count fields use float64 because GraphQL JSON encodes all numbers as floats.
type PQResultData struct {
	Columns                      []PQColumn `json:"columns"`
	Cells                        [][]PQCell `json:"cells"`
	MatchCount                   float64    `json:"matchCount"`
	OmittedEvents                float64    `json:"omittedEvents"`
	Outcome                      string     `json:"outcome"`
	PartialResultsDueToTimeLimit bool       `json:"partialResultsDueToTimeLimit"`
}

// PQColumn describes a column in PQ results.
type PQColumn struct {
	Name          string `json:"name"`
	Format        string `json:"format"`
	Type          string `json:"type"`
	DecimalPlaces *int   `json:"decimalPlaces,omitempty"`
}

// PQCell is a single cell value in PQ results.
type PQCell struct {
	Value any    `json:"value"`
	URL   string `json:"url,omitempty"`
}

// FacetResultData is the typed result for TOP_FACETS queries.
type FacetResultData struct {
	Facets            []Facet `json:"facets"`
	MatchCount        float64 `json:"matchCount"`
	SampledEventCount float64 `json:"sampledEventCount"`
	Outcome           string  `json:"outcome"`
}

// Facet is a single facet in TOP_FACETS results.
type Facet struct {
	Name              string       `json:"name"`
	IsNumeric         bool         `json:"isNumeric"`
	MatchCount        float64      `json:"matchCount"`
	SampledMatchCount float64      `json:"sampledMatchCount"`
	UniqueValuesCount float64      `json:"uniqueValuesCount"`
	Values            []FacetValue `json:"values"`
}

// FacetValue is a single value within a facet.
type FacetValue struct {
	Count float64 `json:"count"`
	Value string  `json:"value"`
}

// GraphQL request/response wrappers (unexported).

type graphqlRequest struct {
	Query     string `json:"query"`
	Variables any    `json:"variables,omitempty"`
}

type graphqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphqlError  `json:"errors,omitempty"`
}

type graphqlError struct {
	Message string `json:"message"`
}

type launchQueryData struct {
	LaunchQuery QueriesResult `json:"launchQuery"`
}

type pingQueryData struct {
	PingQuery QueriesResult `json:"pingQuery"`
}
