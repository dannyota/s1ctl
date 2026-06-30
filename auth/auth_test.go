package auth

import (
	"net/http"
	"testing"
)

func TestApiTokenHeader(t *testing.T) {
	creds := NewApiToken("tok_abc")
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err := creds.Apply(req); err != nil {
		t.Fatal(err)
	}
	want := "ApiToken tok_abc"
	if got := req.Header.Get("Authorization"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBearerHeader(t *testing.T) {
	creds := NewBearer("tok_abc")
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err := creds.Apply(req); err != nil {
		t.Fatal(err)
	}
	want := "Bearer tok_abc"
	if got := req.Header.Get("Authorization"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEmptyToken(t *testing.T) {
	for _, creds := range []Credentials{NewApiToken(""), NewBearer("")} {
		req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
		if err := creds.Apply(req); err == nil {
			t.Error("expected error for empty token")
		}
	}
}

func TestRoundTripper(t *testing.T) {
	creds := NewBearer("tok_abc")
	rt := RoundTripper(creds, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		want := "Bearer tok_abc"
		if got := req.Header.Get("Authorization"); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	}))
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestRoundTripperDoesNotMutateOriginal(t *testing.T) {
	creds := NewApiToken("tok_abc")
	rt := RoundTripper(creds, roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody}, nil
	}))
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if got := req.Header.Get("Authorization"); got != "" {
		t.Errorf("original request mutated: Authorization = %q", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
