package graphql

import (
	"context"
	"encoding/json"
	"errors"
)

// CnappCloudProvider identifies a cloud provider for CNAPP onboarding.
type CnappCloudProvider string

const (
	CnappCloudProviderAWS         CnappCloudProvider = "AWS"
	CnappCloudProviderGCP         CnappCloudProvider = "GCP"
	CnappCloudProviderAzure       CnappCloudProvider = "AZURE"
	CnappCloudProviderOCI         CnappCloudProvider = "OCI"
	CnappCloudProviderAlibaba     CnappCloudProvider = "ALIBABA"
	CnappCloudProviderAWSGovCloud CnappCloudProvider = "AWS_GOVCLOUD"
	CnappCloudProviderAWSChina    CnappCloudProvider = "AWS_CHINA"
)

// CnappCloudEntityType is the type of a cloud entity (organization, individual, member).
type CnappCloudEntityType string

const (
	CnappCloudEntityTypeOrganization         CnappCloudEntityType = "ORGANIZATION"
	CnappCloudEntityTypeIndividual           CnappCloudEntityType = "INDIVIDUAL"
	CnappCloudEntityTypeMember               CnappCloudEntityType = "MEMBER"
	CnappCloudEntityTypeGovCloudOrganization CnappCloudEntityType = "GOVCLOUD_ORGANIZATION"
	CnappCloudEntityTypeGovCloudIndividual   CnappCloudEntityType = "GOVCLOUD_INDIVIDUAL"
	CnappCloudEntityTypeAWSChinaIndividual   CnappCloudEntityType = "AWS_CHINA_INDIVIDUAL"
)

// CnappOperationalStatus is the operational status of an onboarded cloud entity.
type CnappOperationalStatus string

const (
	CnappOperationalStatusOperational       CnappOperationalStatus = "OPERATIONAL"
	CnappOperationalStatusFailed            CnappOperationalStatus = "FAILED"
	CnappOperationalStatusInactive          CnappOperationalStatus = "INACTIVE"
	CnappOperationalStatusResyncing         CnappOperationalStatus = "RESYNCING"
	CnappOperationalStatusInProgress        CnappOperationalStatus = "INPROGRESS"
	CnappOperationalStatusLicenseExpired    CnappOperationalStatus = "LICENSE_EXPIRED"
	CnappOperationalStatusPartiallyLicensed CnappOperationalStatus = "PARTIALLY_LICENSED"
)

// CnappOnboardingType is the type of cloud onboarding (organization or individual).
type CnappOnboardingType string

const (
	CnappOnboardingTypeOrganization CnappOnboardingType = "ORGANIZATION"
	CnappOnboardingTypeIndividual   CnappOnboardingType = "INDIVIDUAL"
)

// CnappScopeType is the scope level for CNAPP operations.
type CnappScopeType string

const (
	CnappScopeTypeTenant  CnappScopeType = "TENANT"
	CnappScopeTypeAccount CnappScopeType = "ACCOUNT"
	CnappScopeTypeSite    CnappScopeType = "SITE"
	CnappScopeTypeGroup   CnappScopeType = "GROUP"
)

// CnappS1Product identifies a SentinelOne product for cloud onboarding.
type CnappS1Product string

const (
	CnappS1ProductCloudVMInventory      CnappS1Product = "CLOUD_VM_INVENTORY"
	CnappS1ProductCloudNativeSecurity   CnappS1Product = "CLOUD_NATIVE_SECURITY"
	CnappS1ProductCloudDataSecurity     CnappS1Product = "CLOUD_DATA_SECURITY"
	CnappS1ProductIdentitySecurity      CnappS1Product = "IDENTITY_SECURITY"
	CnappS1ProductCloudWorkloadSecurity CnappS1Product = "CLOUD_WORKLOAD_SECURITY"
)

// CnappScopeSelector is the scope selector for CNAPP GraphQL operations.
type CnappScopeSelector struct {
	ScopeType CnappScopeType `json:"scopeType"`
	ScopeIDs  []int64        `json:"scopeIds"`
}

// CnappOnboardedEntity is an onboarded cloud entity from the list view. Fields
// map the CnappOnboardedCloudEntitiesView schema type. Timestamps are epoch
// milliseconds (Long in GraphQL).
type CnappOnboardedEntity struct {
	ID                    string                 `json:"id"`
	Type                  CnappCloudEntityType   `json:"type"`
	RootEntityID          string                 `json:"rootEntityId"`
	EntityID              string                 `json:"entityId"`
	Name                  string                 `json:"name"`
	OnboardingStatus      CnappOperationalStatus `json:"onboardingStatus"`
	Path                  string                 `json:"path"`
	ActiveCoverage        []string               `json:"activeCoverage"`
	MissingCoverage       []string               `json:"missingCoverage"`
	ResourceLink          string                 `json:"resourceLink"`
	PermissionGranted     string                 `json:"permissionGranted"`
	Scope                 string                 `json:"scope"`
	CreatedAt             int64                  `json:"createdAt"`
	UpdatedAt             int64                  `json:"updatedAt"`
	ErrorCount            int                    `json:"errorCount"`
	HasCoverageGaps       bool                   `json:"hasCoverageGaps"`
	CreatedBy             string                 `json:"createdBy"`
	UpdatedBy             string                 `json:"updatedBy"`
	CloudLogIngestionType string                 `json:"cloudLogIngestionType"`

	Raw json.RawMessage `json:"-"`
}

func (e *CnappOnboardedEntity) UnmarshalJSON(b []byte) error {
	type alias CnappOnboardedEntity
	if err := json.Unmarshal(b, (*alias)(e)); err != nil {
		return err
	}
	e.Raw = append(e.Raw[:0:0], b...)
	return nil
}

// CnappEntityDetail is the response from the single-entity get query. It has
// different fields than the list view — product details and display name
// instead of coverage and error counts.
type CnappEntityDetail struct {
	EntityID        string              `json:"entityId"`
	EntityName      string              `json:"entityName"`
	DisplayName     string              `json:"displayName"`
	OnboardingType  CnappOnboardingType `json:"onboardingType"`
	CloudProvider   CnappCloudProvider  `json:"cloudProvider"`
	ActiveProducts  json.RawMessage     `json:"activeProducts"`
	ExtraProperties json.RawMessage     `json:"extraProperties"`

	Raw json.RawMessage `json:"-"`
}

func (d *CnappEntityDetail) UnmarshalJSON(b []byte) error {
	type alias CnappEntityDetail
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// CnappActionResponse is the response from delete/activate/deactivate mutations.
type CnappActionResponse struct {
	Message   string `json:"message"`
	IsSuccess bool   `json:"isSuccess"`

	Raw json.RawMessage `json:"-"`
}

func (r *CnappActionResponse) UnmarshalJSON(b []byte) error {
	type alias CnappActionResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// CnappOnboardEntityResponse is the response from onboardCnappCloudEntity.
type CnappOnboardEntityResponse struct {
	Message   string `json:"message"`
	IsSuccess bool   `json:"isSuccess"`

	Raw json.RawMessage `json:"-"`
}

func (r *CnappOnboardEntityResponse) UnmarshalJSON(b []byte) error {
	type alias CnappOnboardEntityResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// ErrCnappDeleteNoAccountIDs is returned when a delete is attempted without
// any account IDs. The SDK rejects empty lists to prevent an unbounded action.
var ErrCnappDeleteNoAccountIDs = errors.New("graphql: cloud onboarding delete requires at least one account ID")

// cnappEntityFields is the shared field selection for cloud entity list.
const cnappEntityFields = `
    id
    type
    rootEntityId
    entityId
    name
    onboardingStatus
    path
    activeCoverage
    missingCoverage
    resourceLink
    permissionGranted
    scope
    createdAt
    updatedAt
    errorCount
    hasCoverageGaps
    createdBy
    updatedBy
    cloudLogIngestionType`

const cnappEntitiesListQuery = `query CnappOnboardedEntities($first: Int, $after: String, $filters: [CnappFilterInput]!, $scope: CnappScopeSelector, $sort: CnappSort) {
  cnappOnboardedCloudEntitiesV2(first: $first, after: $after, filters: $filters, scope: $scope, sort: $sort) {
    edges {
      cursor
      node {` + cnappEntityFields + `
      }
    }
    pageInfo { hasNextPage hasPreviousPage endCursor startCursor }
    totalCount
  }
}`

// CnappEntitiesList lists onboarded cloud entities. filters is required by the
// schema (pass an empty slice for no filters). scope and page are optional.
func (c *Client) CnappEntitiesList(ctx context.Context, filters []CnappFilter, scope *CnappScopeSelector, page *ListParams) (*Connection[CnappOnboardedEntity], error) {
	vars := map[string]any{
		"filters": filters,
	}
	if scope != nil {
		vars["scope"] = scope
	}
	if page != nil {
		if page.First > 0 {
			vars["first"] = page.First
		}
		if page.After != "" {
			vars["after"] = page.After
		}
		if page.Sort != nil {
			vars["sort"] = page.Sort
		}
	}
	var resp struct {
		Entities Connection[CnappOnboardedEntity] `json:"cnappOnboardedCloudEntitiesV2"`
	}
	if err := c.Do(ctx, EndpointCloudOnboarding, cnappEntitiesListQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.Entities, nil
}

const cnappEntityGetQuery = `query CnappOnboardedEntity($request: CnappActionRequest!, $scopeSelector: CnappScopeSelector) {
  cnappOnboardedCloudEntity(request: $request, scopeSelector: $scopeSelector) {
    entityId
    entityName
    displayName
    onboardingType
    cloudProvider
    activeProducts { s1Product features }
    extraProperties
  }
}`

// CnappEntityGet returns a single onboarded cloud entity by account ID(s). The
// API takes an action request with accountIds; typically a single ID is passed.
func (c *Client) CnappEntityGet(ctx context.Context, accountIDs []string, scope *CnappScopeSelector) (*CnappEntityDetail, error) {
	vars := map[string]any{
		"request": map[string]any{"accountIds": accountIDs},
	}
	if scope != nil {
		vars["scopeSelector"] = scope
	}
	var resp struct {
		CnappOnboardedCloudEntity *CnappEntityDetail `json:"cnappOnboardedCloudEntity"`
	}
	if err := c.Do(ctx, EndpointCloudOnboarding, cnappEntityGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	if resp.CnappOnboardedCloudEntity == nil {
		return nil, &QueryError{Errors: []GQLError{{Message: "cloud entity not found"}}}
	}
	return resp.CnappOnboardedCloudEntity, nil
}

const cnappOnboardMutation = `mutation CnappOnboard($request: CnappCloudOnboardingRequest!, $scopeSelector: CnappScopeSelector) {
  onboardCnappCloudEntity(request: $request, scopeSelector: $scopeSelector) {
    message
    isSuccess
  }
}`

// CnappOnboard onboards a new cloud entity. request is the full onboarding
// payload (CnappCloudOnboardingRequest); it is sent as raw JSON so the complex
// nested body (credentials, products, leaves/branches) is carried verbatim.
func (c *Client) CnappOnboard(ctx context.Context, request json.RawMessage, scope *CnappScopeSelector) (*CnappOnboardEntityResponse, error) {
	vars := map[string]any{
		"request": request,
	}
	if scope != nil {
		vars["scopeSelector"] = scope
	}
	var resp struct {
		OnboardCnappCloudEntity *CnappOnboardEntityResponse `json:"onboardCnappCloudEntity"`
	}
	if err := c.Do(ctx, EndpointCloudOnboarding, cnappOnboardMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.OnboardCnappCloudEntity, nil
}

const cnappDeleteMutation = `mutation CnappDelete($request: CnappActionRequest!, $scopeSelector: CnappScopeSelector) {
  deleteCnappCloudEntity(request: $request, scopeSelector: $scopeSelector) {
    message
    isSuccess
  }
}`

// CnappDelete deletes (offboards) cloud entities by account ID(s). At least
// one account ID is required; an empty list is rejected by the SDK.
func (c *Client) CnappDelete(ctx context.Context, accountIDs []string, scope *CnappScopeSelector) (*CnappActionResponse, error) {
	if len(accountIDs) == 0 {
		return nil, ErrCnappDeleteNoAccountIDs
	}
	vars := map[string]any{
		"request": map[string]any{"accountIds": accountIDs},
	}
	if scope != nil {
		vars["scopeSelector"] = scope
	}
	var resp struct {
		DeleteCnappCloudEntity *CnappActionResponse `json:"deleteCnappCloudEntity"`
	}
	if err := c.Do(ctx, EndpointCloudOnboarding, cnappDeleteMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.DeleteCnappCloudEntity, nil
}

// CnappFilter is a CNAPP filter input. Unlike the shared Filter type, CNAPP
// uses its own filter structure with typed match variants.
type CnappFilter struct {
	FieldID     string      `json:"fieldId"`
	IsNegated   bool        `json:"isNegated,omitempty"`
	StringIn    *CnappInStr `json:"stringIn,omitempty"`
	StringEqual *CnappEqStr `json:"stringEqual,omitempty"`
}

// CnappInStr is a CNAPP string "in" filter.
type CnappInStr struct {
	Values []string `json:"values"`
}

// CnappEqStr is a CNAPP string "equal" filter.
type CnappEqStr struct {
	Value string `json:"value"`
}
