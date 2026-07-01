// Package graphql is a Go client for the SentinelOne GraphQL APIs
// (UAM Alerts, xSPM, Cloud Security).
//
// The client is pure — HTTP calls and typed structs, no disk I/O.
package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"danny.vn/s1/auth"
)

// Endpoint identifies a GraphQL API domain.
type Endpoint string

const (
	EndpointAlerts            Endpoint = "/web/api/v2.1/unifiedalerts/graphql"
	EndpointMisconfigurations Endpoint = "/web/api/v2.1/xspm/findings/misconfigurations/graphql"
	EndpointVulnerabilities   Endpoint = "/web/api/v2.1/xspm/findings/vulnerabilities/graphql"
	EndpointCloudPolicies     Endpoint = "/web/api/v2.1/cloudsecurity/policies/graphql"
	EndpointCloudOnboarding   Endpoint = "/web/api/v2.1/cloudonboarding/graphql"
	EndpointCloudCompliance   Endpoint = "/web/api/v2.1/cloudsecurity/compliance/graphql"
)

// Client is a SentinelOne GraphQL client. Safe for concurrent use.
type Client struct {
	baseURL string
	http    *http.Client
	limiter *rate.Limiter
}

// Option customizes a Client.
type Option func(*Client)

// WithHTTPClient overrides the underlying *http.Client.
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }

// WithRateLimit overrides the default rate limiter. rps is the sustained
// requests-per-second rate; burst is the maximum burst size.
func WithRateLimit(rps float64, burst int) Option {
	return func(c *Client) { c.limiter = rate.NewLimiter(rate.Limit(rps), burst) }
}

// NewClient builds a GraphQL client.
//
// consoleURL is the console base URL. token is the API token used with
// Bearer auth for all GraphQL endpoints.
func NewClient(consoleURL, token string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(consoleURL, "/"),
		http: &http.Client{
			Timeout:   60 * time.Second,
			Transport: auth.RoundTripper(auth.NewBearer(token), nil),
		},
		limiter: rate.NewLimiter(rate.Limit(10), 20),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

type gqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type gqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []GQLError      `json:"errors,omitempty"`
}

// GQLError is a GraphQL error from the response.
type GQLError struct {
	Message string `json:"message"`
	Path    []any  `json:"path,omitempty"`
}

func (e GQLError) Error() string { return e.Message }

// Do executes a GraphQL query against the given endpoint.
// The result is unmarshalled from response.data into dst.
func (c *Client) Do(ctx context.Context, endpoint Endpoint, query string, vars map[string]any, dst any) error {
	body, err := json.Marshal(gqlRequest{Query: query, Variables: vars})
	if err != nil {
		return fmt.Errorf("graphql: marshal: %w", err)
	}

	if err := c.limiter.Wait(ctx); err != nil {
		return fmt.Errorf("graphql: rate limit: %w", err)
	}

	u := c.baseURL + string(endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("graphql: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("graphql: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("graphql: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &HTTPError{Status: resp.StatusCode, Body: data}
	}

	var gqlResp gqlResponse
	if err := json.Unmarshal(data, &gqlResp); err != nil {
		return fmt.Errorf("graphql: unmarshal: %w", err)
	}
	if len(gqlResp.Errors) > 0 {
		return &QueryError{Errors: gqlResp.Errors}
	}

	if dst != nil {
		if err := json.Unmarshal(gqlResp.Data, dst); err != nil {
			return fmt.Errorf("graphql: unmarshal data: %w", err)
		}
	}
	return nil
}

// HTTPError is a non-2xx HTTP response.
type HTTPError struct {
	Status int
	Body   []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("graphql: HTTP %d", e.Status)
}

// QueryError contains one or more GraphQL errors.
type QueryError struct {
	Errors []GQLError
}

func (e *QueryError) Error() string {
	if len(e.Errors) == 1 {
		return fmt.Sprintf("graphql: %s", e.Errors[0].Message)
	}
	return fmt.Sprintf("graphql: %d errors (first: %s)", len(e.Errors), e.Errors[0].Message)
}
