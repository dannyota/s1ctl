package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestDeployCredGroupList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/ranger/cred-groups" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if q.Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":              "1000000000000000001",
					"groupName":       "prod-deploy",
					"groupPassphrase": "encrypted-pass",
					"scopeId":         "225494730938493804",
					"domain":          "CORP",
					"targetOs":        "windows",
					"totalDetails":    3,
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "abc"},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.DeployCredGroupList(context.Background(), &DeployCredGroupListParams{
		SiteIDs: []string{"225494730938493804"},
		Limit:   10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	g := items[0]
	if g.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", g.ID)
	}
	if g.GroupName != "prod-deploy" {
		t.Fatalf("unexpected groupName: %s", g.GroupName)
	}
	if g.TargetOS != DeployTargetOSWindows {
		t.Fatalf("unexpected targetOs: %s", g.TargetOS)
	}
	if g.TotalDetails != 3 {
		t.Fatalf("expected totalDetails=3, got %d", g.TotalDetails)
	}
	if g.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestDeployCredGroupListNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %q", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	items, _, err := c.DeployCredGroupList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestDeployCredGroupCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/ranger/cred-groups" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data map[string]any `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["groupName"] != "new-group" {
			t.Fatalf("unexpected groupName: %v", body.Data["groupName"])
		}
		if body.Data["scopeId"] != "225494730938493804" {
			t.Fatalf("unexpected scopeId: %v", body.Data["scopeId"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":              "1000000000000000002",
				"groupName":       "new-group",
				"groupPassphrase": "enc",
				"scopeId":         "225494730938493804",
				"targetOs":        "windows",
			},
		})
	})
	c := testClient(t, handler)
	g, err := c.DeployCredGroupCreate(context.Background(), DeployCredGroupCreateInput{
		GroupName:       "new-group",
		GroupPassphrase: "enc",
		ScopeID:         "225494730938493804",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", g.ID)
	}
	if g.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestDeployCredGroupDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/ranger/cred-groups/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true},
		})
	})
	c := testClient(t, handler)
	err := c.DeployCredGroupDelete(context.Background(), "1000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeployCredDetailList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/ranger/cred-groups/details" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["credGroupIds"]; !slices.Equal(got, []string{"1000000000000000001"}) {
			t.Fatalf("unexpected credGroupIds: %v", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":          "2000000000000000001",
					"credGroupId": "1000000000000000001",
					"title":       "Admin",
					"credType":    "User/Password",
					"createdAt":   "2025-01-01T00:00:00Z",
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	items, pag, err := c.DeployCredDetailList(context.Background(), &DeployCredDetailListParams{
		CredGroupIDs: []string{"1000000000000000001"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	d := items[0]
	if d.ID != "2000000000000000001" {
		t.Fatalf("unexpected ID: %s", d.ID)
	}
	if d.Title != "Admin" {
		t.Fatalf("unexpected title: %s", d.Title)
	}
	if d.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
}

func TestDeployCredDetailAdd(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/ranger/cred-groups/details" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data struct {
				CredGroupID string `json:"credGroupId"`
				Details     []any  `json:"details"`
			} `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data.CredGroupID != "1000000000000000001" {
			t.Fatalf("unexpected credGroupId: %v", body.Data.CredGroupID)
		}
		if len(body.Data.Details) != 1 {
			t.Fatalf("expected 1 detail, got %d", len(body.Data.Details))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true},
		})
	})
	c := testClient(t, handler)
	err := c.DeployCredDetailAdd(context.Background(), DeployCredDetailAddInput{
		CredGroupID: "1000000000000000001",
		Details: []DeployCredDetailInput{
			{Title: "Admin", CredType: "User/Password", EncryptedKey: "key1", EncryptedCred: "cred1"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeployCredDetailUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/ranger/cred-groups/details/2000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":          "2000000000000000001",
				"credGroupId": "1000000000000000001",
				"title":       "Updated",
				"credType":    "User/Password",
			},
		})
	})
	c := testClient(t, handler)
	d, err := c.DeployCredDetailUpdate(context.Background(), "2000000000000000001", DeployCredDetailInput{
		Title: "Updated", CredType: "User/Password", EncryptedKey: "k", EncryptedCred: "c",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Title != "Updated" {
		t.Fatalf("unexpected title: %s", d.Title)
	}
}

func TestDeployCredDetailDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/ranger/cred-groups/details/2000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true},
		})
	})
	c := testClient(t, handler)
	err := c.DeployCredDetailDelete(context.Background(), "2000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeployCredGroupListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{{"code": 403, "title": "Forbidden"}},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.DeployCredGroupList(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 403 {
		t.Fatalf("expected 403, got %d", ae.Status)
	}
}

func TestDeployEnumValues(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"TargetOSWindows", string(DeployTargetOSWindows), "windows"},
		{"TargetOSOSXLinux", string(DeployTargetOSOSXLinux), "osx_linux"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}
