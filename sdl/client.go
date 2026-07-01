// Package sdl is a Go client for the SentinelOne Singularity Data Lake (SDL) API.
//
// The client is pure — HTTP calls and typed structs, no disk I/O.
package sdl

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

// Client is a SentinelOne SDL API client. Safe for concurrent use.
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

// NewClient builds an SDL API client.
//
// consoleURL is the console base URL. token is the API token used with
// Bearer auth for SDL endpoints.
func NewClient(consoleURL, token string, opts ...Option) *Client {
	c := &Client{
		baseURL: strings.TrimRight(consoleURL, "/"),
		http: &http.Client{
			Timeout:   120 * time.Second,
			Transport: auth.RoundTripper(auth.NewBearer(token), nil),
		},
		limiter: rate.NewLimiter(rate.Limit(10), 20),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// do executes req and decodes the JSON response body into dst.
func (c *Client) do(req *http.Request, dst any) error {
	if err := c.limiter.Wait(req.Context()); err != nil {
		return fmt.Errorf("sdl: rate limit: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("sdl: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("sdl: read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &APIError{Status: resp.StatusCode, Body: respBody}
	}

	if dst != nil {
		if err := json.Unmarshal(respBody, dst); err != nil {
			return fmt.Errorf("sdl: unmarshal: %w", err)
		}
	}
	return nil
}

func (c *Client) post(ctx context.Context, path string, body, dst any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("sdl: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("sdl: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, dst)
}

func (c *Client) postText(ctx context.Context, path, contentType string, headers map[string]string, body io.Reader, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("sdl: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.do(req, dst)
}

// APIError is a non-2xx response from the SDL API.
type APIError struct {
	Status int
	Body   []byte
}

func (e *APIError) Error() string {
	if len(e.Body) > 0 {
		return fmt.Sprintf("sdl: HTTP %d: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("sdl: HTTP %d", e.Status)
}
