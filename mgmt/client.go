// Package mgmt is a Go client for the SentinelOne REST Management API v2.1.
//
// The client is pure — HTTP calls and typed structs, no disk I/O. All on-disk
// layout lives in internal/.
package mgmt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"danny.vn/s1/auth"
)

// Client is a SentinelOne MGMT API client. Safe for concurrent use.
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

// NewClient builds a MGMT API client.
//
// consoleURL is the console base URL (e.g. "https://your-console.sentinelone.net").
// token is the API token. Auth is applied via the ApiToken header format.
func NewClient(consoleURL, token string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(consoleURL, "/") + "/web/api/v2.1",
		http: &http.Client{
			Timeout:   60 * time.Second,
			Transport: auth.RoundTripper(auth.NewApiToken(token), nil),
		},
		limiter: rate.NewLimiter(rate.Limit(10), 20),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// BaseURL returns the resolved API base URL.
func (c *Client) BaseURL() string { return c.baseURL }

func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	return c.queryRequest(ctx, http.MethodGet, path, params, dst)
}

func (c *Client) post(ctx context.Context, path string, body, dst any) error {
	return c.jsonRequest(ctx, http.MethodPost, path, body, dst)
}

func (c *Client) put(ctx context.Context, path string, body, dst any) error {
	return c.jsonRequest(ctx, http.MethodPut, path, body, dst)
}

func (c *Client) delete(ctx context.Context, path string, params url.Values, dst any) error {
	return c.queryRequest(ctx, http.MethodDelete, path, params, dst)
}

// queryRequest sends a request with query parameters and no body (GET, DELETE).
func (c *Client) queryRequest(ctx context.Context, method, path string, params url.Values, dst any) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return fmt.Errorf("mgmt: %w", err)
	}
	return c.do(req, dst)
}

// jsonRequest sends a request with a JSON body (POST, PUT).
func (c *Client) jsonRequest(ctx context.Context, method, path string, body, dst any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("mgmt: marshal: %w", err)
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("mgmt: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, dst)
}

// getRaw sends a GET request and returns the raw response body without
// JSON unmarshalling. Used for export endpoints that return CSV or other
// non-JSON formats.
func (c *Client) getRaw(ctx context.Context, path string, params url.Values) ([]byte, error) {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("mgmt: %w", err)
	}

	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("mgmt: rate limit: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mgmt: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("mgmt: read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, parseError(resp.StatusCode, data)
	}
	return data, nil
}

func (c *Client) do(req *http.Request, dst any) error {
	if err := c.limiter.Wait(req.Context()); err != nil {
		return fmt.Errorf("mgmt: rate limit: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("mgmt: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("mgmt: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return parseError(resp.StatusCode, data)
	}

	if dst != nil {
		if err := json.Unmarshal(data, dst); err != nil {
			return fmt.Errorf("mgmt: unmarshal: %w", err)
		}
	}
	return nil
}
