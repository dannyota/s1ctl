package graphql

import (
	"context"
	"encoding/json"
)

// CloudPolicyScope is the scope of a cloud policy (CNS rule).
type CloudPolicyScope struct {
	ID    string `json:"id"`
	Level string `json:"level"`
	Path  string `json:"path"`
}

// CloudPolicy is a CNS (Cloud Native Security) rule.
type CloudPolicy struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Severity          string           `json:"severity"`
	Status            string           `json:"status"`
	Type              string           `json:"type"`
	PolicyCode        string           `json:"policyCode"`
	Providers         []string         `json:"providers"`
	ResourceType      string           `json:"resourceType"`
	Category          string           `json:"category"`
	SubCategory       string           `json:"subCategory"`
	IsSystem          bool             `json:"isSystem"`
	CreatedAt         string           `json:"createdAt"`
	UpdatedAt         string           `json:"updatedAt"`
	CreatedBy         string           `json:"createdBy"`
	UpdatedBy         string           `json:"updatedBy"`
	Scope             CloudPolicyScope `json:"scope"`
	RecommendedAction string           `json:"recommendedAction"`
	Impact            string           `json:"impact"`
	IssueMessage      string           `json:"issueMessage"`
	Reference         string           `json:"reference"`

	Raw json.RawMessage `json:"-"`
}

func (p *CloudPolicy) UnmarshalJSON(b []byte) error {
	type alias CloudPolicy
	if err := json.Unmarshal(b, (*alias)(p)); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

const cloudPoliciesQuery = `query CloudPolicies($first: Int, $after: String, $filters: [CloudCommonFilterInput!], $scope: CloudCommonScopeSelector) {
  cnsRules(first: $first, after: $after, filters: $filters, scope: $scope) {
    edges {
      cursor
      node {
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
        createdAt
        updatedAt
        createdBy
        updatedBy
        scope { id level path }
        recommendedAction
        impact
        issueMessage
        reference
      }
    }
    pageInfo { hasNextPage hasPreviousPage endCursor startCursor }
    totalCount
  }
}`

// CloudPoliciesList queries CNS rules (cloud security policies).
func (c *Client) CloudPoliciesList(ctx context.Context, params *ListParams) (*Connection[CloudPolicy], error) {
	var resp struct {
		CnsRules Connection[CloudPolicy] `json:"cnsRules"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cloudPoliciesQuery, listVars(params), &resp); err != nil {
		return nil, err
	}
	return &resp.CnsRules, nil
}

const cloudPolicyGetQuery = `query CloudPolicyGet($id: ID!) {
  cnsRule(id: $id) {
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
    createdAt
    updatedAt
    createdBy
    updatedBy
    scope { id level path }
    recommendedAction
    impact
    issueMessage
    reference
  }
}`

// CloudPoliciesGet returns a single CNS rule by ID.
func (c *Client) CloudPoliciesGet(ctx context.Context, id string) (*CloudPolicy, error) {
	vars := map[string]any{"id": id}
	var resp struct {
		CnsRule *CloudPolicy `json:"cnsRule"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cloudPolicyGetQuery, vars, &resp); err != nil {
		return nil, err
	}
	if resp.CnsRule == nil {
		return nil, &QueryError{Errors: []GQLError{{Message: "cloud policy not found"}}}
	}
	return resp.CnsRule, nil
}

// CloudPoliciesActionResponse is the response from actionOnCNSRules.
type CloudPoliciesActionResponse struct {
	IDs []string `json:"ids"`

	Raw json.RawMessage `json:"-"`
}

func (r *CloudPoliciesActionResponse) UnmarshalJSON(b []byte) error {
	type alias CloudPoliciesActionResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

const cloudPoliciesActionMutation = `mutation CloudPoliciesAction($action: String!, $input: CloudCommonActionInput!) {
  actionOnCNSRules(action: $action, input: $input) {
    ids
  }
}`

// CloudPoliciesAction performs a bulk action (enable, disable, delete) on CNS
// rules by ID.
func (c *Client) CloudPoliciesAction(ctx context.Context, action string, ids []string) (*CloudPoliciesActionResponse, error) {
	vars := map[string]any{
		"action": action,
		"input":  map[string]any{"ids": ids},
	}
	var resp struct {
		ActionOnCNSRules *CloudPoliciesActionResponse `json:"actionOnCNSRules"`
	}
	if err := c.Do(ctx, EndpointCloudPolicies, cloudPoliciesActionMutation, vars, &resp); err != nil {
		return nil, err
	}
	return resp.ActionOnCNSRules, nil
}

// CloudPoliciesEnable enables the specified CNS rules.
func (c *Client) CloudPoliciesEnable(ctx context.Context, ids []string) (*CloudPoliciesActionResponse, error) {
	return c.CloudPoliciesAction(ctx, "enable", ids)
}

// CloudPoliciesDisable disables the specified CNS rules.
func (c *Client) CloudPoliciesDisable(ctx context.Context, ids []string) (*CloudPoliciesActionResponse, error) {
	return c.CloudPoliciesAction(ctx, "disable", ids)
}

// CloudPoliciesDelete deletes the specified CNS rules.
func (c *Client) CloudPoliciesDelete(ctx context.Context, ids []string) (*CloudPoliciesActionResponse, error) {
	return c.CloudPoliciesAction(ctx, "delete", ids)
}
