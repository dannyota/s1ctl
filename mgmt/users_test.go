package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestUsersUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/users/1000000000000000005" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["fullName"] != "New Name" {
			t.Fatalf("unexpected fullName: %v", body.Data["fullName"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"id": "1000000000000000005", "fullName": "New Name"},
		})
	})
	c := testClient(t, handler)
	u, err := c.UsersUpdate(context.Background(), "1000000000000000005", UserUpdate{FullName: "New Name"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.FullName != "New Name" {
		t.Fatalf("unexpected fullName: %s", u.FullName)
	}
}

func TestUsersGenerateToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/users/generate-api-token" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"token": "PLACEHOLDER_TOKEN_VALUE"},
		})
	})
	c := testClient(t, handler)
	tok, err := c.UsersGenerateToken(context.Background(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "PLACEHOLDER_TOKEN_VALUE" {
		t.Fatalf("unexpected token: %s", tok)
	}
}

func TestUsersRevokeToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/users/revoke-api-token" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data.ID != "1000000000000000005" {
			t.Fatalf("unexpected id: %s", body.Data.ID)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"success": true}})
	})
	c := testClient(t, handler)
	if err := c.UsersRevokeToken(context.Background(), "1000000000000000005"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUsersTokenDetails(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/users/api-token-details" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"createdAt": "2025-01-01T00:00:00Z",
				"expiresAt": "2026-01-01T00:00:00Z",
			},
		})
	})
	c := testClient(t, handler)
	d, err := c.UsersTokenDetails(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.ExpiresAt != "2026-01-01T00:00:00Z" {
		t.Fatalf("unexpected expiresAt: %s", d.ExpiresAt)
	}
	if d.Token != "" {
		t.Fatalf("expected no token in metadata response, got %q", d.Token)
	}
}

func TestUsersTokenDetailsByID(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/users/1000000000000000005/api-token-details" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"createdAt": "2025-01-01T00:00:00Z",
				"expiresAt": "2026-01-01T00:00:00Z",
			},
		})
	})
	c := testClient(t, handler)
	d, err := c.UsersTokenDetailsByID(context.Background(), "1000000000000000005")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.CreatedAt != "2025-01-01T00:00:00Z" {
		t.Fatalf("unexpected createdAt: %s", d.CreatedAt)
	}
}

func TestUsers2FAEnableDisable(t *testing.T) {
	for _, tc := range []struct {
		name string
		path string
		call func(*Client) error
	}{
		{"enable", "/users/2fa/enable", func(c *Client) error {
			return c.Users2FAEnable(context.Background(), "1000000000000000005")
		}},
		{"disable", "/users/2fa/disable", func(c *Client) error {
			return c.Users2FADisable(context.Background(), "1000000000000000005")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatalf("expected POST, got %s", r.Method)
				}
				if r.URL.Path != tc.path {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				var body struct {
					Data struct {
						ID string `json:"id"`
					} `json:"data"`
				}
				json.NewDecoder(r.Body).Decode(&body)
				if body.Data.ID != "1000000000000000005" {
					t.Fatalf("unexpected id: %s", body.Data.ID)
				}
				json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"success": true}})
			})
			c := testClient(t, handler)
			if err := tc.call(c); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
