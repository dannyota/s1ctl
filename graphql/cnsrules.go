package graphql

import (
	"context"
	"encoding/json"
)

// CNSRuleType is the type of a CNS (Cloud Native Security) rule.
type CNSRuleType string

const (
	CNSRuleTypeCloudMisconfiguration   CNSRuleType = "CloudMisconfiguration"
	CNSRuleTypeKubeMisconfiguration    CNSRuleType = "KubeMisconfiguration"
	CNSRuleTypeIaCSecurity             CNSRuleType = "IaCSecurity"
	CNSRuleTypeCIEMMisconfiguration    CNSRuleType = "CIEMMisconfiguration"
	CNSRuleTypeAISPMMisconfiguration   CNSRuleType = "AISPMMisconfiguration"
	CNSRuleTypeAdmissionController     CNSRuleType = "AdmissionController"
	CNSRuleTypeAttackPath              CNSRuleType = "AttackPath"
	CNSRuleTypeOffensiveSecurity       CNSRuleType = "OffensiveSecurity"
	CNSRuleTypeVulnerabilityManagement CNSRuleType = "VulnerabilityManagement"
)

// CNSRuleAction is a bulk action applied to CNS rules.
type CNSRuleAction string

const (
	CNSRuleActionEnable  CNSRuleAction = "enable"
	CNSRuleActionDisable CNSRuleAction = "disable"
	CNSRuleActionDelete  CNSRuleAction = "delete"
)

// CNSRule is a CNS (Cloud Native Security) rule: a built-in or user-defined
// cloud policy. The rule body (Rego/graph query plus config parameters) is
// carried as-is; use Raw for the fields not modelled here.
type CNSRule struct {
	ID                   string           `json:"id"`
	Name                 string           `json:"name"`
	Description          string           `json:"description"`
	Severity             string           `json:"severity"`
	Status               string           `json:"status"`
	Type                 string           `json:"type"`
	PolicyCode           string           `json:"policyCode"`
	Providers            []string         `json:"providers"`
	ResourceType         string           `json:"resourceType"`
	Category             string           `json:"category"`
	SubCategory          string           `json:"subCategory"`
	IsSystem             bool             `json:"isSystem"`
	QueryType            string           `json:"queryType"`
	RawQuery             string           `json:"rawQuery"`
	RuleConfigParameters string           `json:"ruleConfigParameters"`
	EnforcementAction    string           `json:"enforcementAction"`
	MgmtID               string           `json:"mgmtId"`
	RecommendedAction    string           `json:"recommendedAction"`
	Impact               string           `json:"impact"`
	IssueMessage         string           `json:"issueMessage"`
	Reference            string           `json:"reference"`
	Scope                CloudPolicyScope `json:"scope"`
	CreatedAt            string           `json:"createdAt"`
	UpdatedAt            string           `json:"updatedAt"`
	CreatedBy            string           `json:"createdBy"`
	UpdatedBy            string           `json:"updatedBy"`

	Raw json.RawMessage `json:"-"`
}

func (r *CNSRule) UnmarshalJSON(b []byte) error {
	type alias CNSRule
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// cnsRuleFields is the shared field selection for cnsRules and cnsRule.
const cnsRuleFields = `
    id
    name
    description
    severity
    status
    type
    policyCode
    providers
    resourceType
    category
    subCategory
    isSystem
    queryType
    rawQuery
    ruleConfigParameters
    enforcementAction
    mgmtId
    recommendedAction
    impact
    issueMessage
    reference
    scope { id level path }
    createdAt
    updatedAt
    createdBy
    updatedBy`

const cnsRulesQuery = `query CNSRules($first: Int, $after: String, $filters: [CloudCommonFilterInput!], $scope: CloudCommonScopeSelector, $sort: CloudCommonSortInput) {
  cnsRules(first: $first, after: $after, filters: $filters, scope: $scope, sort: $sort) {
    edges {
      cursor
      node {` + cnsRuleFields + `
      }
    }
    pageInfo { hasNextPage hasPreviousPage endCursor startCursor }
    totalCount
  }
}`

// CNSRulesList queries CNS rules (custom and built-in cloud policies). filters
// and scope override any set on page; page carries pagination (first/after) and
// sort. Scope is optional: absent means global scope.
func (c *Client) CNSRulesList(ctx context.Context, filters []Filter, scope *Scope, page *ListParams) (*Connection[CNSRule], error) {
	p := ListParams{Filters: filters, Scope: scope}
	if page != nil {
		p.First = page.First
		p.After = page.After
		p.Sort = page.Sort
	}
	var resp struct {
		CnsRules Connection[CNSRule] `json:"cnsRules"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cnsRulesQuery, listVars(&p), &resp); err != nil {
		return nil, err
	}
	return &resp.CnsRules, nil
}

const cnsRuleGetQuery = `query CNSRuleGet($id: ID!, $scope: CloudCommonScopeSelector) {
  cnsRule(id: $id, scope: $scope) {` + cnsRuleFields + `
  }
}`

// CNSRuleGet returns a single CNS rule by ID. Scope is optional.
func (c *Client) CNSRuleGet(ctx context.Context, id string, scope *Scope) (*CNSRule, error) {
	vars := map[string]any{"id": id}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		CnsRule *CNSRule `json:"cnsRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cnsRuleGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	if resp.CnsRule == nil {
		return nil, &QueryError{Errors: []GQLError{{Message: "cns rule not found"}}}
	}
	return resp.CnsRule, nil
}

// CNSRuleCreateResponse is the response from createCNSRule.
type CNSRuleCreateResponse struct {
	ID string `json:"id"`

	Raw json.RawMessage `json:"-"`
}

func (r *CNSRuleCreateResponse) UnmarshalJSON(b []byte) error {
	type alias CNSRuleCreateResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

const cnsRuleCreateMutation = `mutation CNSRuleCreate($input: CNSRuleInput!, $scope: CloudCommonScopeSelector) {
  createCNSRule(input: $input, scope: $scope) {
    id
  }
}`

// CNSRuleCreate creates a new CNS rule. input is a CNSRuleInput payload (name,
// queryType, rawQuery, severity, type, ruleConfigParameters, etc.); it is sent
// verbatim so the large Rego/config body is carried without lossy re-typing.
// Scope is required by the API: a rule cannot be created in global scope.
func (c *Client) CNSRuleCreate(ctx context.Context, input json.RawMessage, scope *Scope) (*CNSRuleCreateResponse, error) {
	vars := map[string]any{"input": input}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		CreateCNSRule *CNSRuleCreateResponse `json:"createCNSRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cnsRuleCreateMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.CreateCNSRule, nil
}

const cnsRuleUpdateMutation = `mutation CNSRuleUpdate($id: ID!, $input: CNSRuleInput!, $scope: CloudCommonScopeSelector) {
  updateCNSRule(id: $id, input: $input, scope: $scope)
}`

// CNSRuleUpdate replaces an existing CNS rule by ID. The update is a full
// replacement (not a patch): input must carry the complete CNSRuleInput. A few
// fields (type, providers, queryType, resourceType) cannot change once created.
func (c *Client) CNSRuleUpdate(ctx context.Context, id string, input json.RawMessage, scope *Scope) (bool, error) {
	vars := map[string]any{"id": id, "input": input}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		UpdateCNSRule bool `json:"updateCNSRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cnsRuleUpdateMutation, vars, &resp); err != nil {
		return false, err
	}
	return resp.UpdateCNSRule, nil
}

const cnsRulesActionMutation = `mutation CNSRulesAction($action: String!, $input: CloudCommonActionInput!, $scope: CloudCommonScopeSelector) {
  actionOnCNSRules(action: $action, input: $input, scope: $scope) {
    ids
  }
}`

// CNSRulesAction performs a bulk action (enable, disable, delete) on CNS rules
// by ID. At least one ID is required: the API treats an empty ids list as "act
// on all rules in scope", so the SDK rejects it with ErrCloudPolicyActionNoIDs.
func (c *Client) CNSRulesAction(ctx context.Context, action CNSRuleAction, ids []string, scope *Scope) (*CloudPoliciesActionResponse, error) {
	if len(ids) == 0 {
		return nil, ErrCloudPolicyActionNoIDs
	}
	vars := map[string]any{
		"action": string(action),
		"input":  map[string]any{"ids": ids},
	}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		ActionOnCNSRules *CloudPoliciesActionResponse `json:"actionOnCNSRules"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cnsRulesActionMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.ActionOnCNSRules, nil
}

// CNSRulesEnable enables the specified CNS rules.
func (c *Client) CNSRulesEnable(ctx context.Context, ids []string, scope *Scope) (*CloudPoliciesActionResponse, error) {
	return c.CNSRulesAction(ctx, CNSRuleActionEnable, ids, scope)
}

// CNSRulesDisable disables the specified CNS rules.
func (c *Client) CNSRulesDisable(ctx context.Context, ids []string, scope *Scope) (*CloudPoliciesActionResponse, error) {
	return c.CNSRulesAction(ctx, CNSRuleActionDisable, ids, scope)
}

// CNSRulesDelete deletes the specified CNS rules.
func (c *Client) CNSRulesDelete(ctx context.Context, ids []string, scope *Scope) (*CloudPoliciesActionResponse, error) {
	return c.CNSRulesAction(ctx, CNSRuleActionDelete, ids, scope)
}

// CNSRegoEvaluateResponse is the response from evaluateCNSRegoRule.
type CNSRegoEvaluateResponse struct {
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
	Result string          `json:"result"`

	Raw json.RawMessage `json:"-"`
}

func (r *CNSRegoEvaluateResponse) UnmarshalJSON(b []byte) error {
	type alias CNSRegoEvaluateResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

const cnsRuleEvaluateMutation = `mutation CNSRuleEvaluate($policyId: String, $regoQuery: String!, $resourceData: JSON!, $ruleConfigParameters: String, $scope: CloudCommonScopeSelector) {
  evaluateCNSRegoRule(policyId: $policyId, regoQuery: $regoQuery, resourceData: $resourceData, ruleConfigParameters: $ruleConfigParameters, scope: $scope) {
    data
    error
    result
  }
}`

// CNSRuleEvaluate runs a raw Rego query against an asset's JSON before a rule is
// created or updated. resourceData is the asset JSON; policyID and
// ruleConfigParameters are optional. This is a dry-check: it evaluates only and
// mutates nothing.
func (c *Client) CNSRuleEvaluate(ctx context.Context, policyID, regoQuery string, resourceData json.RawMessage, ruleConfigParameters string, scope *Scope) (*CNSRegoEvaluateResponse, error) {
	vars := map[string]any{
		"regoQuery":    regoQuery,
		"resourceData": resourceData,
	}
	if policyID != "" {
		vars["policyId"] = policyID
	}
	if ruleConfigParameters != "" {
		vars["ruleConfigParameters"] = ruleConfigParameters
	}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		EvaluateCNSRegoRule *CNSRegoEvaluateResponse `json:"evaluateCNSRegoRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cnsRuleEvaluateMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.EvaluateCNSRegoRule, nil
}

// CNSRuleTypeInfo maps a CNS rule type to its display title.
type CNSRuleTypeInfo struct {
	Key   string `json:"key"`
	Title string `json:"title"`

	Raw json.RawMessage `json:"-"`
}

func (t *CNSRuleTypeInfo) UnmarshalJSON(b []byte) error {
	type alias CNSRuleTypeInfo
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

const cnsRuleTypesQuery = `query CNSRuleTypes($scope: CloudCommonScopeSelector) {
  cnsRuleTypes(scope: $scope) {
    key
    title
  }
}`

// CNSRuleTypes returns the supported CNS rule types with their titles.
func (c *Client) CNSRuleTypes(ctx context.Context, scope *Scope) ([]CNSRuleTypeInfo, error) {
	vars := map[string]any{}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		CnsRuleTypes []CNSRuleTypeInfo `json:"cnsRuleTypes"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cnsRuleTypesQuery, vars, &resp); err != nil {
		return nil, err
	}
	return resp.CnsRuleTypes, nil
}

const cnsRuleConfigQuery = `query CNSRuleConfig($type: CNSRuleType!, $scope: CloudCommonScopeSelector) {
  cnsRuleConfig(type: $type, scope: $scope)
}`

// CNSRuleConfig returns the rule config (limits and allowed values) for a rule
// type as raw JSON. Scope is optional.
func (c *Client) CNSRuleConfig(ctx context.Context, scope *Scope, ruleType CNSRuleType) (json.RawMessage, error) {
	vars := map[string]any{"type": string(ruleType)}
	if scope != nil {
		vars["scope"] = scope
	}
	var resp struct {
		CnsRuleConfig json.RawMessage `json:"cnsRuleConfig"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cnsRuleConfigQuery, vars, &resp); err != nil {
		return nil, err
	}
	return resp.CnsRuleConfig, nil
}
