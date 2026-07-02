package mgmt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAgentsDisconnect(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/agents/actions/disconnect" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req actionRequest
		json.NewDecoder(r.Body).Decode(&req)
		if len(req.Filter.IDs) != 1 || req.Filter.IDs[0] != "A1" {
			t.Fatalf("unexpected filter: %+v", req.Filter)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"affected": 1},
		})
	})
	c := testClient(t, handler)
	affected, err := c.AgentsDisconnect(context.Background(), ActionFilter{IDs: []string{"A1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if affected != 1 {
		t.Fatalf("expected 1 affected, got %d", affected)
	}
}

func TestAgentActionRequiresFilter(t *testing.T) {
	c := NewClient("https://example.sentinelone.net", "tok")
	_, err := c.AgentsDisconnect(context.Background(), ActionFilter{})
	if err == nil {
		t.Fatal("expected error for empty filter")
	}
}

func TestAgentsBroadcast(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/agents/actions/broadcast" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var req actionRequest
		var body struct {
			Data struct {
				Message string `json:"message"`
			} `json:"data"`
		}
		raw, _ := io.ReadAll(r.Body)
		json.Unmarshal(raw, &req)
		json.Unmarshal(raw, &body)
		if body.Data.Message != "heads up" {
			t.Fatalf("unexpected message: %q", body.Data.Message)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 3}})
	})
	c := testClient(t, handler)
	n, err := c.AgentsBroadcast(context.Background(), "heads up", ActionFilter{IDs: []string{"A1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 affected, got %d", n)
	}
}

func TestAgentsResetPassphrase(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/agents/actions/reset-passphrase" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"results": []map[string]any{
					{"agentId": "A1", "attempted": true, "status": "started"},
					{"agentId": "A2", "attempted": false, "status": "skipped"},
				},
				"summary": map[string]any{},
			},
		})
	})
	c := testClient(t, handler)
	n, err := c.AgentsResetPassphrase(context.Background(), ActionFilter{IDs: []string{"A1", "A2"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 attempted, got %d", n)
	}
}

func TestAgentsRanger(t *testing.T) {
	for _, tc := range []struct {
		enable bool
		path   string
	}{
		{true, "/agents/actions/ranger-enable"},
		{false, "/agents/actions/ranger-disable"},
	} {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != tc.path {
				t.Fatalf("unexpected %s %s (want %s)", r.Method, r.URL.Path, tc.path)
			}
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 2}})
		})
		c := testClient(t, handler)
		n, err := c.AgentsRanger(context.Background(), tc.enable, ActionFilter{IDs: []string{"A1"}})
		if err != nil {
			t.Fatalf("enable=%v: unexpected error: %v", tc.enable, err)
		}
		if n != 2 {
			t.Fatalf("enable=%v: expected 2, got %d", tc.enable, n)
		}
	}
}

func TestAgentsFetchInventoryActions(t *testing.T) {
	for _, tc := range []struct {
		path string
		call func(*Client, context.Context, ActionFilter) (int, error)
	}{
		{"/agents/actions/fetch-installed-apps", (*Client).AgentsFetchInstalledApps},
		{"/agents/actions/fetch-firewall-rules", (*Client).AgentsFetchFirewallRules},
	} {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != tc.path {
				t.Fatalf("unexpected %s %s (want %s)", r.Method, r.URL.Path, tc.path)
			}
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 1}})
		})
		c := testClient(t, handler)
		n, err := tc.call(c, context.Background(), ActionFilter{IDs: []string{"A1"}})
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.path, err)
		}
		if n != 1 {
			t.Fatalf("%s: expected 1, got %d", tc.path, n)
		}
	}
}

func TestAgentsFetchFiles(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/agents/A1/actions/fetch-files" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var body struct {
			Data struct {
				Files    []string `json:"files"`
				Password string   `json:"password"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if len(body.Data.Files) != 2 || body.Data.Files[0] != "/etc/hosts" {
			t.Fatalf("unexpected files: %v", body.Data.Files)
		}
		if body.Data.Password != "pw-placeholder" {
			t.Fatalf("unexpected password field: %q", body.Data.Password)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"success": true}})
	})
	c := testClient(t, handler)
	ok, err := c.AgentsFetchFiles(context.Background(), "A1", []string{"/etc/hosts", "/var/log/syslog"}, "pw-placeholder")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected success true")
	}
}

func TestAgentsLocalUpgradeAuthorization(t *testing.T) {
	var gotBody string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/agents/actions/local-upgrade-authorization" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
		raw, _ := io.ReadAll(r.Body)
		gotBody = string(raw)
		json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"affected": 4}})
	})
	c := testClient(t, handler)

	n, err := c.AgentsLocalUpgradeAuthorization(context.Background(), ActionFilter{IDs: []string{"A1"}}, "2030-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 4 {
		t.Fatalf("expected 4, got %d", n)
	}
	if !strings.Contains(gotBody, `"agentAuthorization":"2030-01-01T00:00:00Z"`) {
		t.Fatalf("expected timestamp in body, got %s", gotBody)
	}

	if _, err := c.AgentsLocalUpgradeAuthorization(context.Background(), ActionFilter{IDs: []string{"A1"}}, ""); err != nil {
		t.Fatalf("unexpected error on revoke: %v", err)
	}
	if !strings.Contains(gotBody, `"agentAuthorization":null`) {
		t.Fatalf("expected null in revoke body, got %s", gotBody)
	}
}

func TestAgentsLocalUpgradeAuthGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/agents/A1/local-upgrade-authorization" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"agentAuthorization": "2030-01-01T00:00:00Z",
				"siteAuthorization":  "2031-01-01T00:00:00Z",
			},
		})
	})
	c := testClient(t, handler)
	auth, err := c.AgentsLocalUpgradeAuthGet(context.Background(), "A1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.AgentAuthorization != "2030-01-01T00:00:00Z" || auth.SiteAuthorization != "2031-01-01T00:00:00Z" {
		t.Fatalf("unexpected auth: %+v", auth)
	}
}

func TestAgentsPassphrases(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/agents/passphrases" {
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("siteIds"); got != "S1" {
			t.Fatalf("unexpected siteIds: %q", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "A1", "uuid": "uuid-1", "computerName": "HOST-1", "passphrase": "PASS PHRASE PLACEHOLDER"},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.AgentsPassphrases(context.Background(), &AgentPassphraseParams{SiteIDs: []string{"S1"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].ID != "A1" || items[0].Passphrase != "PASS PHRASE PLACEHOLDER" {
		t.Fatalf("unexpected items: %+v", items)
	}
	if pag == nil || pag.TotalItems != 1 {
		t.Fatalf("unexpected pagination: %+v", pag)
	}
}
