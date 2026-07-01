package mgmt

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// listResponse is the standard SentinelOne list envelope.
type listResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// singleResponse is the standard SentinelOne single-object envelope.
type singleResponse[T any] struct {
	Data T `json:"data"`
}

// list fetches a paginated list of resources.
func list[T any](c *Client, ctx context.Context, path string, params url.Values) ([]T, *Pagination, error) {
	var resp listResponse[T]
	if err := c.get(ctx, path, params, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// create posts a resource creation wrapped in {"data": ...} and returns the created resource.
func create[T any](c *Client, ctx context.Context, path string, data any) (*T, error) {
	req := map[string]any{"data": data}
	var resp singleResponse[T]
	if err := c.post(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// update puts a resource update wrapped in {"data": ...} and returns the updated resource.
func update[T any](c *Client, ctx context.Context, path string, data any) (*T, error) {
	req := map[string]any{"data": data}
	var resp singleResponse[T]
	if err := c.put(ctx, path, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// getByID fetches a single resource by ID using the ?ids= query param.
func getByID[T any](c *Client, ctx context.Context, path, entity, id string) (*T, error) {
	params := url.Values{}
	params.Set("ids", id)
	items, _, err := list[T](c, ctx, path, params)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("mgmt: %s %s not found", entity, id)
	}
	return &items[0], nil
}

// affectedResponse is the envelope for mutation actions.
type affectedResponse struct {
	Data struct {
		Affected int `json:"affected"`
	} `json:"data"`
}

// ActionFilter identifies which resources to act on.
type ActionFilter struct {
	IDs     []string `json:"ids,omitempty"`
	SiteIDs []string `json:"siteIds,omitempty"`
	Query   string   `json:"query,omitempty"`
}

type actionRequest struct {
	Filter ActionFilter `json:"filter"`
	Data   any          `json:"data,omitempty"`
}

// doAction posts a mutation action and returns the affected count.
func doAction(c *Client, ctx context.Context, path string, filter ActionFilter, data any) (int, error) { //nolint:unparam
	if len(filter.IDs) == 0 && len(filter.SiteIDs) == 0 && filter.Query == "" {
		return 0, fmt.Errorf("mgmt: action requires at least one filter (ids, siteIds, or query)")
	}
	req := actionRequest{Filter: filter, Data: data}
	var resp affectedResponse
	if err := c.post(ctx, path, req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.Affected, nil
}

// query param helpers

func addCSV(v url.Values, key string, vals []string) {
	for _, s := range vals {
		v.Add(key, s)
	}
}

func addString(v url.Values, key, val string) {
	if val != "" {
		v.Set(key, val)
	}
}

func addInt(v url.Values, key string, val int) { //nolint:unparam
	if val > 0 {
		v.Set(key, strconv.Itoa(val))
	}
}

func addBool(v url.Values, key string, val *bool) {
	if val != nil {
		v.Set(key, strconv.FormatBool(*val))
	}
}

func addIntCSV(v url.Values, key string, vals []int) {
	for _, i := range vals {
		v.Add(key, strconv.Itoa(i))
	}
}
