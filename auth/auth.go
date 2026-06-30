// Package auth provides credential types for SentinelOne API surfaces.
//
// SentinelOne uses one token value with two header formats:
//
//   - ApiToken — for the REST MGMT API (Authorization: ApiToken <token>)
//   - Bearer  — for the SDL and GraphQL APIs (Authorization: Bearer <token>)
//
// Both implement Credentials, so each SDK client takes only the format it needs.
// Credentials resolve lazily and are safe for concurrent use.
package auth

import "net/http"

// Credentials applies authentication to an outbound request.
// Implementations must be safe for concurrent use.
type Credentials interface {
	Apply(req *http.Request) error
}

// RoundTripper wraps base so every request carries creds. A nil base uses
// http.DefaultTransport.
func RoundTripper(creds Credentials, base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &authTransport{creds: creds, base: base}
}

type authTransport struct {
	creds Credentials
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r2 := req.Clone(req.Context())
	if err := t.creds.Apply(r2); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(r2)
}
