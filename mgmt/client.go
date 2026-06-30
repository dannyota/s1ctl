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

	"danny.vn/s1/auth"
)

// Client is a SentinelOne MGMT API client. Safe for concurrent use.
type Client struct {
	baseURL string
	http    *http.Client
}

// Option customizes a Client.
type Option func(*Client)

// WithHTTPClient overrides the underlying *http.Client.
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }

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
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// BaseURL returns the resolved API base URL.
func (c *Client) BaseURL() string { return c.baseURL }

func (c *Client) get(ctx context.Context, path string, params url.Values, dst any) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("mgmt: %w", err)
	}
	return c.do(req, dst)
}

func (c *Client) post(ctx context.Context, path string, body, dst any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("mgmt: marshal: %w", err)
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, reader)
	if err != nil {
		return fmt.Errorf("mgmt: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, dst)
}

func (c *Client) do(req *http.Request, dst any) error {
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
