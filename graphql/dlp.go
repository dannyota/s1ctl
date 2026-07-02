package graphql

import (
	"context"
	"encoding/json"
	"errors"
)

// DLPRuleStatus is the status of a data protection rule.
type DLPRuleStatus string

const (
	DLPRuleStatusEnabled  DLPRuleStatus = "ENABLED"
	DLPRuleStatusDisabled DLPRuleStatus = "DISABLED"
)

// DLPClassificationType is the type of a DLP classification.
type DLPClassificationType string

const (
	DLPClassificationTypeAIContextual  DLPClassificationType = "AI_CONTEXTUAL"
	DLPClassificationTypeCodeDetector  DLPClassificationType = "CODE_DETECTOR"
	DLPClassificationTypeRegex         DLPClassificationType = "REGEX"
	DLPClassificationTypeSecrets       DLPClassificationType = "SECRETS"
	DLPClassificationTypeSensitiveData DLPClassificationType = "SENSITIVE_DATA"
)

// ErrDLPScopeRequired is returned when a DLP query that requires a scope (engine
// settings) is called without one. The schema marks the selector as non-null.
var ErrDLPScopeRequired = errors.New("graphql: dlp engine settings require a scope")

// DLPPageInfo is the page-based pagination info returned by DLP queries. Unlike
// the Relay cursor connections used elsewhere, DLP pages are numbered.
type DLPPageInfo struct {
	CurrentPage     int  `json:"currentPage"`
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
	PageSize        int  `json:"pageSize"`
	TotalCount      int  `json:"totalCount"`
	TotalPages      int  `json:"totalPages"`

	Raw json.RawMessage `json:"-"`
}

func (p *DLPPageInfo) UnmarshalJSON(b []byte) error {
	type alias DLPPageInfo
	if err := json.Unmarshal(b, (*alias)(p)); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

// DLPConnection is a page-based DLP result set (nodes plus page info). DLP
// queries return a numbered-page connection, not a Relay cursor connection.
type DLPConnection[T any] struct {
	Nodes    []T         `json:"nodes"`
	PageInfo DLPPageInfo `json:"pageInfo"`
}

// DLPPage is page-based pagination input for DLP list queries.
type DLPPage struct {
	Page     int
	PageSize int
}

// DLPClassificationSummary is the classification reference embedded in a rule.
type DLPClassificationSummary struct {
	ID   string                `json:"id"`
	Name string                `json:"name"`
	Type DLPClassificationType `json:"type"`

	Raw json.RawMessage `json:"-"`
}

func (s *DLPClassificationSummary) UnmarshalJSON(b []byte) error {
	type alias DLPClassificationSummary
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// DLPRule is a data protection rule. The heavy nested bodies (actions,
// impacted endpoints, inspection conditions) are carried as raw JSON: they are
// only populated by DLPRuleGet, which selects them in full.
type DLPRule struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name"`
	Description     string                     `json:"description"`
	Status          DLPRuleStatus              `json:"status"`
	Rank            int                        `json:"rank"`
	RuleCode        string                     `json:"ruleCode"`
	SystemPolicy    bool                       `json:"systemPolicy"`
	CreatedAt       string                     `json:"createdAt"`
	CreatedBy       string                     `json:"createdBy"`
	UpdatedAt       string                     `json:"updatedAt"`
	Scope           CloudPolicyScope           `json:"scope"`
	Classifications []DLPClassificationSummary `json:"classifications"`

	Actions              json.RawMessage `json:"actions"`
	ImpactedEndpoints    json.RawMessage `json:"impactedEndpoints"`
	InspectionConditions json.RawMessage `json:"inspectionConditions"`

	Raw json.RawMessage `json:"-"`
}

func (r *DLPRule) UnmarshalJSON(b []byte) error {
	type alias DLPRule
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// DLPClassification is a DLP classification entity. Type-specific bodies
// (patterns, data types, secret detectors, code languages) are carried as raw
// JSON and are only populated for the matching classification type.
type DLPClassification struct {
	ID                 string                `json:"id"`
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	Type               DLPClassificationType `json:"type"`
	ClassificationCode string                `json:"classificationCode"`
	SystemPolicy       bool                  `json:"systemPolicy"`
	UsedInRulesCount   int                   `json:"usedInRulesCount"`
	CreatedAt          string                `json:"createdAt"`
	UpdatedAt          string                `json:"updatedAt"`
	Scope              CloudPolicyScope      `json:"scope"`

	DetectionStrictness string          `json:"detectionStrictness"`
	PromptBox           string          `json:"promptBox"`
	ExcludedKeywords    []string        `json:"excludedKeywords"`
	CodeLanguages       json.RawMessage `json:"codeLanguages"`
	DataTypes           json.RawMessage `json:"dataTypes"`
	Patterns            json.RawMessage `json:"patterns"`
	SecretDetectors     json.RawMessage `json:"secretDetectors"`

	Raw json.RawMessage `json:"-"`
}

func (c *DLPClassification) UnmarshalJSON(b []byte) error {
	type alias DLPClassification
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// DLPEngineSettings is the DLP engine configuration for a scope.
type DLPEngineSettings struct {
	BlockEncryptedArchive    bool             `json:"blockEncryptedArchive"`
	BlockUsbModifications    bool             `json:"blockUsbModifications"`
	CharacterInspectionDepth string           `json:"characterInspectionDepth"`
	ClassificationsToInspect int              `json:"classificationsToInspect"`
	CreatedAt                string           `json:"createdAt"`
	EnableOCR                bool             `json:"enableOcr"`
	IgnoreKeywords           []string         `json:"ignoreKeywords"`
	IgnoreRegexes            []string         `json:"ignoreRegexes"`
	InspectionSizeLimit      int              `json:"inspectionSizeLimit"`
	MaskEvidence             bool             `json:"maskEvidence"`
	MaxArchiveLevels         int              `json:"maxArchiveLevels"`
	MaxInspectedFileSize     int64            `json:"maxInspectedFileSize"`
	NotificationMessage      string           `json:"notificationMessage"`
	PreventAction            string           `json:"preventAction"`
	PublishingEnabled        bool             `json:"publishingEnabled"`
	Scope                    CloudPolicyScope `json:"scope"`
	UpdatedAt                string           `json:"updatedAt"`
	UpdatedBy                string           `json:"updatedBy"`

	Raw json.RawMessage `json:"-"`
}

func (s *DLPEngineSettings) UnmarshalJSON(b []byte) error {
	type alias DLPEngineSettings
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// DLPRuleFilter filters dataProtectionRules. Zero-value fields are omitted.
type DLPRuleFilter struct {
	SearchName          string                  `json:"searchName,omitempty"`
	SearchDescription   string                  `json:"searchDescription,omitempty"`
	Status              []DLPRuleStatus         `json:"status,omitempty"`
	Channels            []string                `json:"channels,omitempty"`
	ClassificationTypes []DLPClassificationType `json:"classificationTypes,omitempty"`
}

// DLPClassificationFilter filters dlpClassifications. Zero-value fields are omitted.
type DLPClassificationFilter struct {
	SearchName        string                  `json:"searchName,omitempty"`
	SearchDescription string                  `json:"searchDescription,omitempty"`
	Type              []DLPClassificationType `json:"type,omitempty"`
}

// dlpPaginationVar builds the DlpPaginationInput variable (page/pageSize).
func dlpPaginationVar(p *DLPPage) map[string]any {
	m := map[string]any{}
	if p.Page > 0 {
		m["page"] = p.Page
	}
	if p.PageSize > 0 {
		m["pageSize"] = p.PageSize
	}
	return m
}

// dlpRuleFields is the core rule selection used by list.
const dlpRuleFields = `
    id
    name
    description
    status
    rank
    ruleCode
    systemPolicy
    createdAt
    createdBy
    updatedAt
    scope { id level path }
    classifications { id name type }`

// dlpRuleFullFields adds the heavy nested bodies for the single-rule get.
const dlpRuleFullFields = dlpRuleFields + `
    actions {
      actionTaken
      onScreenNotification { enabled message }
      raiseAlert { enabled severity }
    }
    impactedEndpoints {
      assignedTags { key value }
      osTypes
      scopeType
    }
    inspectionConditions {
      channels {
        usbDevice { enabled includeClipboardContent }
        webDestinations { destinationType enabled exceptions }
      }
      fileOrigin { originType allowedOrigins }
      fileTypes
    }`

const dlpRulesQuery = `query DLPRules($filter: DlpRuleFilterInput, $pagination: DlpPaginationInput, $scope: CloudCommonScopeSelector) {
  dataProtectionRules(filter: $filter, pagination: $pagination, scope: $scope) {
    nodes {` + dlpRuleFields + `
    }
    pageInfo { currentPage hasNextPage hasPreviousPage pageSize totalCount totalPages }
  }
}`

// DLPRulesList queries data protection rules. filter, scope, and page are all
// optional: nil filter/scope means no filter/global scope, nil page means the
// server default page. DLP pages are numbered (page/pageSize), not cursor-based.
func (c *Client) DLPRulesList(ctx context.Context, filter *DLPRuleFilter, scope *Scope, page *DLPPage) (*DLPConnection[DLPRule], error) {
	vars := map[string]any{}
	if filter != nil {
		vars["filter"] = filter
	}
	if scope != nil {
		vars["scope"] = scope
	}
	if page != nil {
		vars["pagination"] = dlpPaginationVar(page)
	}
	var resp struct {
		DataProtectionRules DLPConnection[DLPRule] `json:"dataProtectionRules"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpRulesQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.DataProtectionRules, nil
}

const dlpRuleGetQuery = `query DLPRule($id: ID!, $scope: CloudCommonScopeSelector) {
  dataProtectionRule(id: $id, scope: $scope) {` + dlpRuleFullFields + `
  }
}`

// DLPRuleGet returns a single data protection rule by ID, including its full
// action, endpoint, and inspection bodies. Scope is optional.
func (c *Client) DLPRuleGet(ctx context.Context, id string, scope *Scope) (*DLPRule, error) {
	vars := map[string]any{"id": id}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		DataProtectionRule *DLPRule `json:"dataProtectionRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpRuleGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	if resp.DataProtectionRule == nil {
		return nil, &QueryError{Errors: []GQLError{{Message: "data protection rule not found"}}}
	}
	return resp.DataProtectionRule, nil
}

const dlpRuleEnableMutation = `mutation DLPRuleEnable($id: ID!, $scope: CloudCommonScopeSelector) {
  enableDataProtectionRule(id: $id, scope: $scope) { id name status }
}`

// DLPRuleEnable enables a single data protection rule and returns its new state.
func (c *Client) DLPRuleEnable(ctx context.Context, id string, scope *Scope) (*DLPRule, error) {
	vars := map[string]any{"id": id}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		EnableDataProtectionRule *DLPRule `json:"enableDataProtectionRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpRuleEnableMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.EnableDataProtectionRule, nil
}

const dlpRuleDisableMutation = `mutation DLPRuleDisable($id: ID!, $scope: CloudCommonScopeSelector) {
  disableDataProtectionRule(id: $id, scope: $scope) { id name status }
}`

// DLPRuleDisable disables a single data protection rule and returns its new state.
func (c *Client) DLPRuleDisable(ctx context.Context, id string, scope *Scope) (*DLPRule, error) {
	vars := map[string]any{"id": id}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		DisableDataProtectionRule *DLPRule `json:"disableDataProtectionRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpRuleDisableMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.DisableDataProtectionRule, nil
}

const dlpRuleDeleteMutation = `mutation DLPRuleDelete($id: ID!, $scope: CloudCommonScopeSelector) {
  deleteDataProtectionRule(id: $id, scope: $scope)
}`

// DLPRuleDelete deletes a single data protection rule.
func (c *Client) DLPRuleDelete(ctx context.Context, id string, scope *Scope) (bool, error) {
	vars := map[string]any{"id": id}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		DeleteDataProtectionRule bool `json:"deleteDataProtectionRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpRuleDeleteMutation, vars, &resp); err != nil {
		return false, err
	}
	return resp.DeleteDataProtectionRule, nil
}

const dlpRulesBulkEnableMutation = `mutation DLPRulesBulkEnable($ids: [ID!]!, $scope: CloudCommonScopeSelector) {
  bulkEnableDataProtectionRules(ids: $ids, scope: $scope) { id name status }
}`

// DLPRulesBulkEnable enables multiple data protection rules. At least one ID is
// required; an empty list is rejected with ErrCloudPolicyActionNoIDs.
func (c *Client) DLPRulesBulkEnable(ctx context.Context, ids []string, scope *Scope) ([]DLPRule, error) {
	if len(ids) == 0 {
		return nil, ErrCloudPolicyActionNoIDs
	}
	vars := map[string]any{"ids": ids}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		BulkEnableDataProtectionRules []DLPRule `json:"bulkEnableDataProtectionRules"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpRulesBulkEnableMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.BulkEnableDataProtectionRules, nil
}

const dlpRulesBulkDisableMutation = `mutation DLPRulesBulkDisable($ids: [ID!]!, $scope: CloudCommonScopeSelector) {
  bulkDisableDataProtectionRules(ids: $ids, scope: $scope) { id name status }
}`

// DLPRulesBulkDisable disables multiple data protection rules. At least one ID
// is required; an empty list is rejected with ErrCloudPolicyActionNoIDs.
func (c *Client) DLPRulesBulkDisable(ctx context.Context, ids []string, scope *Scope) ([]DLPRule, error) {
	if len(ids) == 0 {
		return nil, ErrCloudPolicyActionNoIDs
	}
	vars := map[string]any{"ids": ids}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		BulkDisableDataProtectionRules []DLPRule `json:"bulkDisableDataProtectionRules"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpRulesBulkDisableMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.BulkDisableDataProtectionRules, nil
}

const dlpRulesBulkDeleteMutation = `mutation DLPRulesBulkDelete($ids: [ID!]!, $scope: CloudCommonScopeSelector) {
  bulkDeleteDataProtectionRules(ids: $ids, scope: $scope)
}`

// DLPRulesBulkDelete deletes multiple data protection rules. At least one ID is
// required; an empty list is rejected with ErrCloudPolicyActionNoIDs.
func (c *Client) DLPRulesBulkDelete(ctx context.Context, ids []string, scope *Scope) (bool, error) {
	if len(ids) == 0 {
		return false, ErrCloudPolicyActionNoIDs
	}
	vars := map[string]any{"ids": ids}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		BulkDeleteDataProtectionRules bool `json:"bulkDeleteDataProtectionRules"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpRulesBulkDeleteMutation, vars, &resp); err != nil {
		return false, err
	}
	return resp.BulkDeleteDataProtectionRules, nil
}

// dlpClassificationFields is the classification selection used by list and get.
const dlpClassificationFields = `
    id
    name
    description
    type
    classificationCode
    systemPolicy
    usedInRulesCount
    createdAt
    updatedAt
    scope { id level path }
    detectionStrictness
    promptBox
    excludedKeywords
    codeLanguages { id name selected }
    dataTypes { dataType { id name description examples complianceFrameworks } occurrence }
    patterns { id name pattern occurrence complianceFrameworks }
    secretDetectors { id name selected }`

const dlpClassificationsQuery = `query DLPClassifications($filter: DlpClassificationFilterInput, $pagination: DlpPaginationInput, $scope: CloudCommonScopeSelector) {
  dlpClassifications(filter: $filter, pagination: $pagination, scope: $scope) {
    nodes {` + dlpClassificationFields + `
    }
    pageInfo { currentPage hasNextPage hasPreviousPage pageSize totalCount totalPages }
  }
}`

// DLPClassificationsList queries DLP classifications. filter, scope, and page
// are optional; DLP pages are numbered (page/pageSize).
func (c *Client) DLPClassificationsList(ctx context.Context, filter *DLPClassificationFilter, scope *Scope, page *DLPPage) (*DLPConnection[DLPClassification], error) {
	vars := map[string]any{}
	if filter != nil {
		vars["filter"] = filter
	}
	if scope != nil {
		vars["scope"] = scope
	}
	if page != nil {
		vars["pagination"] = dlpPaginationVar(page)
	}
	var resp struct {
		DlpClassifications DLPConnection[DLPClassification] `json:"dlpClassifications"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpClassificationsQuery, vars, &resp); err != nil {
		return nil, err
	}
	return &resp.DlpClassifications, nil
}

const dlpClassificationGetQuery = `query DLPClassification($id: ID!, $scope: CloudCommonScopeSelector) {
  dlpClassification(id: $id, scope: $scope) {` + dlpClassificationFields + `
  }
}`

// DLPClassificationGet returns a single DLP classification by ID. Scope is optional.
func (c *Client) DLPClassificationGet(ctx context.Context, id string, scope *Scope) (*DLPClassification, error) {
	vars := map[string]any{"id": id}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		DlpClassification *DLPClassification `json:"dlpClassification"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpClassificationGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	if resp.DlpClassification == nil {
		return nil, &QueryError{Errors: []GQLError{{Message: "dlp classification not found"}}}
	}
	return resp.DlpClassification, nil
}

const dlpClassificationDeleteMutation = `mutation DLPClassificationDelete($id: ID!) {
  deleteDlpClassification(id: $id)
}`

// DLPClassificationDelete deletes a DLP classification by ID.
func (c *Client) DLPClassificationDelete(ctx context.Context, id string) (bool, error) {
	vars := map[string]any{"id": id}
	var resp struct {
		DeleteDlpClassification bool `json:"deleteDlpClassification"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpClassificationDeleteMutation, vars, &resp); err != nil {
		return false, err
	}
	return resp.DeleteDlpClassification, nil
}

const dlpEngineSettingsQuery = `query DLPEngineSettings($scope: CloudCommonScopeSelector!) {
  dlpEngineSettings(scope: $scope) {
    blockEncryptedArchive
    blockUsbModifications
    characterInspectionDepth
    classificationsToInspect
    createdAt
    enableOcr
    ignoreKeywords
    ignoreRegexes
    inspectionSizeLimit
    maskEvidence
    maxArchiveLevels
    maxInspectedFileSize
    notificationMessage
    preventAction
    publishingEnabled
    scope { id level path }
    updatedAt
    updatedBy
  }
}`

// DLPEngineSettings returns the DLP engine configuration for a scope. The schema
// marks the scope selector as required, so a nil scope is rejected up front.
func (c *Client) DLPEngineSettings(ctx context.Context, scope *Scope) (*DLPEngineSettings, error) {
	if scope == nil {
		return nil, ErrDLPScopeRequired
	}
	vars := map[string]any{"scope": scope}
	var resp struct {
		DlpEngineSettings *DLPEngineSettings `json:"dlpEngineSettings"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, dlpEngineSettingsQuery, vars, &resp); err != nil {
		return nil, err
	}
	if resp.DlpEngineSettings == nil {
		return nil, &QueryError{Errors: []GQLError{{Message: "dlp engine settings not found"}}}
	}
	return resp.DlpEngineSettings, nil
}
