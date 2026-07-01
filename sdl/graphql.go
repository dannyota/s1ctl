package sdl

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const graphqlPath = "/sdl/v2/graphql"

const launchQueryGQL = `query launchQuery($queryGroupRequest: QueryGroupRequestInput!) {
  launchQuery(queryGroupRequest: $queryGroupRequest) {
    ids status token stepsCompleted totalSteps
    results {
      id stepsCompleted totalSteps error cacheContext noResultsReason
      data {
        matchCount
        ... on PqResultData {
          columns { name format type decimalPlaces }
          cells { value url }
          matchCount omittedEvents outcome partialResultsDueToTimeLimit
        }
        ... on FacetData {
          facets { name isNumeric matchCount sampledMatchCount uniqueValuesCount
            values { count value }
          }
          matchCount sampledEventCount outcome
        }
      }
    }
  }
}`

const pingQueryGQL = `query pingQuery($ids: [String], $lastStepSeen: Int!, $token: String!) {
  pingQuery(ids: $ids, lastStepSeen: $lastStepSeen, token: $token) {
    ids status token stepsCompleted totalSteps
    results {
      id stepsCompleted totalSteps error cacheContext noResultsReason
      data {
        matchCount
        ... on PqResultData {
          columns { name format type decimalPlaces }
          cells { value url }
          matchCount omittedEvents outcome partialResultsDueToTimeLimit
        }
        ... on FacetData {
          facets { name isNumeric matchCount sampledMatchCount uniqueValuesCount
            values { count value }
          }
          matchCount sampledEventCount outcome
        }
      }
    }
  }
}`

const removeQueryGQL = `mutation removeQuery($token: String!) {
  removeQuery(token: $token)
}`

// graphql sends a GraphQL request and unmarshals the data field into dst.
func (c *Client) graphql(ctx context.Context, query string, variables, dst any) error {
	gqlReq := graphqlRequest{Query: query, Variables: variables}
	var gqlResp graphqlResponse
	if err := c.post(ctx, graphqlPath, gqlReq, &gqlResp); err != nil {
		return err
	}
	if len(gqlResp.Errors) > 0 {
		msgs := make([]string, len(gqlResp.Errors))
		for i, e := range gqlResp.Errors {
			msgs[i] = e.Message
		}
		return fmt.Errorf("sdl graphql: %s", strings.Join(msgs, "; "))
	}
	if dst != nil {
		if err := json.Unmarshal(gqlResp.Data, dst); err != nil {
			return fmt.Errorf("sdl graphql: unmarshal data: %w", err)
		}
	}
	return nil
}

// LaunchQuery starts a query group and returns the initial result.
func (c *Client) LaunchQuery(ctx context.Context, group *QueryGroupRequest) (*QueriesResult, error) {
	vars := map[string]any{"queryGroupRequest": group}
	var data launchQueryData
	if err := c.graphql(ctx, launchQueryGQL, vars, &data); err != nil {
		return nil, err
	}
	return &data.LaunchQuery, nil
}

// PingQuery polls for query results.
func (c *Client) PingQuery(ctx context.Context, ids []string, lastStepSeen int, token string) (*QueriesResult, error) {
	vars := map[string]any{
		"ids":          ids,
		"lastStepSeen": lastStepSeen,
		"token":        token,
	}
	var data pingQueryData
	if err := c.graphql(ctx, pingQueryGQL, vars, &data); err != nil {
		return nil, err
	}
	return &data.PingQuery, nil
}

// RemoveQuery cleans up a completed query token.
func (c *Client) RemoveQuery(ctx context.Context, token string) error {
	vars := map[string]any{"token": token}
	return c.graphql(ctx, removeQueryGQL, vars, nil)
}

// pollUntilDone polls a running query until it reaches a terminal state.
// On context cancellation, it removes the query token before returning.
func (c *Client) pollUntilDone(ctx context.Context, result *QueriesResult) (*QueriesResult, error) {
	for result.Status == QueryStatusRunning {
		select {
		case <-ctx.Done():
			_ = c.RemoveQuery(context.Background(), result.Token)
			return nil, ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
		var err error
		result, err = c.PingQuery(ctx, result.IDs, result.StepsCompleted, result.Token)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// PowerQueryGraphQL executes a PowerQuery via the SDL GraphQL API.
//
// This endpoint lives on the management console ({consoleURL}/sdl/v2/graphql)
// and does not require a separate SDL URL. The client must be created with
// the management console URL.
//
// The result is converted to the same PowerQueryResponse format as the REST
// API, so callers can switch protocols transparently.
func (c *Client) PowerQueryGraphQL(ctx context.Context, req *PowerQueryRequest) (*PowerQueryResponse, error) {
	startNano, err := resolveTime(req.StartTime, true)
	if err != nil {
		return nil, fmt.Errorf("sdl graphql: start time: %w", err)
	}
	endNano, err := resolveTime(req.EndTime, false)
	if err != nil {
		return nil, fmt.Errorf("sdl graphql: end time: %w", err)
	}

	queryID := newQueryID()
	group := &QueryGroupRequest{
		Queries: []QueryRequest{{
			ID:         queryID,
			Type:       QueryTypePQ,
			Filter:     req.Query,
			StartTime:  startNano,
			EndTime:    endNano,
			Origin:     "SEARCH",
			PowerQuery: &struct{}{},
		}},
	}

	result, err := c.LaunchQuery(ctx, group)
	if err != nil {
		return nil, err
	}

	result, err = c.pollUntilDone(ctx, result)
	if err != nil {
		return nil, err
	}
	defer func() { _ = c.RemoveQuery(context.Background(), result.Token) }()

	if result.Status == QueryStatusError {
		var errs []string
		for _, r := range result.Results {
			if r.Error != "" {
				errs = append(errs, r.Error)
			}
		}
		if len(errs) > 0 {
			detail := strings.Join(errs, "; ")
			if result.TotalSteps > 1 {
				return nil, fmt.Errorf("sdl graphql: query failed at step %d/%d: %s (try --protocol rest for complex queries)", result.StepsCompleted, result.TotalSteps, detail)
			}
			return nil, fmt.Errorf("sdl graphql: query error: %s", detail)
		}
		if result.TotalSteps > 1 {
			return nil, fmt.Errorf("sdl graphql: query failed at step %d/%d (try --protocol rest for complex queries)", result.StepsCompleted, result.TotalSteps)
		}
		return nil, fmt.Errorf("sdl graphql: query failed")
	}

	resp, err := convertPQResult(result, queryID)
	if err != nil {
		if result.TotalSteps > 1 {
			return nil, fmt.Errorf("%w (multi-step query, try --protocol rest)", err)
		}
		return nil, err
	}
	return resp, nil
}

// convertPQResult converts a GraphQL QueriesResult into the REST-compatible
// PowerQueryResponse format.
func convertPQResult(result *QueriesResult, queryID string) (*PowerQueryResponse, error) {
	for _, r := range result.Results {
		if r.ID != queryID || len(r.Data) == 0 {
			continue
		}
		var pq PQResultData
		if err := json.Unmarshal(r.Data, &pq); err != nil {
			return nil, fmt.Errorf("sdl graphql: unmarshal PQ data: %w", err)
		}
		resp := &PowerQueryResponse{
			Status: string(result.Status),
		}
		resp.Columns = make([]PowerQueryColumn, len(pq.Columns))
		for i, col := range pq.Columns {
			resp.Columns[i] = PowerQueryColumn{
				Name: col.Name,
				Type: col.Type,
			}
		}
		resp.Values = make([][]any, len(pq.Cells))
		for i, row := range pq.Cells {
			vals := make([]any, len(row))
			for j, cell := range row {
				vals[j] = cell.Value
			}
			resp.Values[i] = vals
		}
		return resp, nil
	}
	return &PowerQueryResponse{Status: string(result.Status)}, nil
}

// resolveTime converts a time specification to a nanosecond epoch string.
//
// Accepted formats:
//   - Relative duration: "24h", "7d", "1h30m"
//   - Epoch seconds (integer): "1719792000"
//   - Nanosecond epoch (>= 1e15): "1719792000000000000"
//   - Empty string: returns "now" (end time) or "24h ago" (start time)
func resolveTime(s string, isStart bool) (string, error) {
	if s == "" {
		if isStart {
			return nanoStr(time.Now().Add(-24 * time.Hour)), nil
		}
		return nanoStr(time.Now()), nil
	}

	if d, ok := parseDuration(s); ok {
		return nanoStr(time.Now().Add(-d)), nil
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		if n > 1e15 {
			return s, nil
		}
		return strconv.FormatInt(n*1e9, 10), nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return nanoStr(t), nil
	}

	return "", fmt.Errorf("unsupported time format: %q", s)
}

// parseDuration extends time.ParseDuration with support for "d" (days).
func parseDuration(s string) (time.Duration, bool) {
	if days, ok := strings.CutSuffix(s, "d"); ok {
		n, err := strconv.Atoi(days)
		if err != nil {
			return 0, false
		}
		return time.Duration(n) * 24 * time.Hour, true
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, false
	}
	return d, true
}

func nanoStr(t time.Time) string {
	return strconv.FormatInt(t.UnixNano(), 10)
}

func newQueryID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
