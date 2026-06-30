package mgmt

import (
	"encoding/json"
	"fmt"
)

// APIError is a non-2xx response from the SentinelOne API.
type APIError struct {
	Status  int    `json:"-"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
	RawBody []byte `json:"-"`
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("mgmt: HTTP %d: %s: %s", e.Status, e.Title, e.Detail)
	}
	if e.Title != "" {
		return fmt.Sprintf("mgmt: HTTP %d: %s", e.Status, e.Title)
	}
	return fmt.Sprintf("mgmt: HTTP %d", e.Status)
}

func parseError(status int, body []byte) error {
	ae := &APIError{Status: status, RawBody: body}
	// Try parsing the standard error envelope.
	var envelope struct {
		Errors []struct {
			Code   int    `json:"code"`
			Title  string `json:"title"`
			Detail string `json:"detail"`
		} `json:"errors"`
	}
	if json.Unmarshal(body, &envelope) == nil && len(envelope.Errors) > 0 {
		ae.Title = envelope.Errors[0].Title
		ae.Detail = envelope.Errors[0].Detail
	}
	return ae
}
