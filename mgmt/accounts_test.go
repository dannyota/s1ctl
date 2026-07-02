package mgmt

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAccountsReactivate(t *testing.T) {
	cases := []struct {
		name       string
		unlimited  bool
		expiration string
	}{
		{"unlimited", true, ""},
		{"expiration", false, "2027-01-01T00:00:00Z"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Fatalf("expected PUT, got %s", r.Method)
				}
				if r.URL.Path != "/accounts/A1/reactivate" {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				var req struct {
					Data struct {
						Unlimited  bool   `json:"unlimited"`
						Expiration string `json:"expiration"`
					} `json:"data"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatalf("decode body: %v", err)
				}
				if req.Data.Unlimited != tc.unlimited {
					t.Fatalf("unlimited: got %v want %v", req.Data.Unlimited, tc.unlimited)
				}
				if req.Data.Expiration != tc.expiration {
					t.Fatalf("expiration: got %q want %q", req.Data.Expiration, tc.expiration)
				}
				w.Write([]byte(`{"data":{}}`))
			})
			c := testClient(t, handler)
			if err := c.AccountsReactivate(context.Background(), "A1", tc.unlimited, tc.expiration); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAccountsExpireNow(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/accounts/A1/expire-now" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{"data":{}}`))
	})
	c := testClient(t, handler)
	if err := c.AccountsExpireNow(context.Background(), "A1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAccountsUninstallPasswordMetadata(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/accounts/A1/uninstall-password/metadata" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"expiration": "2027-01-01", "version": 3,
				"createdAt": "2026-01-01", "generatedByName": "placeholder",
			},
		})
	})
	c := testClient(t, handler)
	meta, err := c.AccountsUninstallPasswordMetadata(context.Background(), "A1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.Version != 3 {
		t.Fatalf("unexpected version: %d", meta.Version)
	}
	if meta.Expiration != "2027-01-01" {
		t.Fatalf("unexpected expiration: %s", meta.Expiration)
	}
}

func TestAccountsUninstallPasswordView(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/accounts/A1/uninstall-password/view" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"password": "PW-PLACEHOLDER"},
		})
	})
	c := testClient(t, handler)
	pw, err := c.AccountsUninstallPasswordView(context.Background(), "A1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw.Password != "PW-PLACEHOLDER" {
		t.Fatalf("unexpected password: %q", pw.Password)
	}
}

func TestAccountsUninstallPasswordGenerate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/accounts/A1/uninstall-password/generate" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req struct {
			Data struct {
				Expiration string `json:"expiration"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Data.Expiration != "2027-01-01" {
			t.Fatalf("unexpected expiration in body: %q", req.Data.Expiration)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"expiration": "2027-01-01", "version": 4},
		})
	})
	c := testClient(t, handler)
	meta, err := c.AccountsUninstallPasswordGenerate(context.Background(), "A1", "2027-01-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.Version != 4 {
		t.Fatalf("unexpected version: %d", meta.Version)
	}
}

func TestAccountsUninstallPasswordRevoke(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/accounts/A1/uninstall-password/revoke" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Write([]byte(`{"data":{}}`))
	})
	c := testClient(t, handler)
	if err := c.AccountsUninstallPasswordRevoke(context.Background(), "A1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
