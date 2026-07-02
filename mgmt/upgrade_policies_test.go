package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestUpgradePoliciesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/upgrade-policy/policies" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("scopeLevel") != "site" {
			t.Fatalf("unexpected scopeLevel: %s", q.Get("scopeLevel"))
		}
		if q.Get("scopeId") != "225494730938493804" {
			t.Fatalf("unexpected scopeId: %s", q.Get("scopeId"))
		}
		if q.Get("osType") != "windows" {
			t.Fatalf("unexpected osType: %s", q.Get("osType"))
		}
		if q.Get("limit") != "25" {
			t.Fatalf("expected limit=25, got %s", q.Get("limit"))
		}
		if q.Get("skip") != "0" {
			t.Fatalf("expected skip=0, got %s", q.Get("skip"))
		}
		if q.Get("sortBy") != "priority" {
			t.Fatalf("unexpected sortBy: %s", q.Get("sortBy"))
		}
		if q.Get("sortOrder") != "asc" {
			t.Fatalf("unexpected sortOrder: %s", q.Get("sortOrder"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"isInherited":          false,
				"policiesInChildScope": true,
				"policies": []map[string]any{
					{
						"id": "1000000000000000001", "name": "Auto Upgrade GA",
						"description": "Upgrade to latest GA", "osType": "windows",
						"scopeLevel": "site", "scopeId": "225494730938493804",
						"isActive": true, "isScheduled": false,
						"allEndpoints": true, "maxRetries": 3, "priority": 1,
						"package": map[string]any{
							"build": "100", "fileId": "200",
							"major": "24", "minor": "1",
						},
						"tags":      []string{"ga"},
						"createdAt": "2025-01-01T00:00:00Z",
					},
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	policies, total, err := c.UpgradePoliciesList(context.Background(), &UpgradePolicyListParams{
		ScopeLevel: "site",
		ScopeID:    "225494730938493804",
		OSType:     "windows",
		Limit:      25,
		Skip:       0,
		SortBy:     "priority",
		SortOrder:  "asc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(policies) != 1 {
		t.Fatalf("expected 1 policy, got %d", len(policies))
	}
	p := policies[0]
	if p.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", p.ID)
	}
	if p.Name != "Auto Upgrade GA" {
		t.Fatalf("unexpected name: %s", p.Name)
	}
	if p.OSType != "windows" {
		t.Fatalf("unexpected osType: %s", p.OSType)
	}
	if p.ScopeLevel != "site" {
		t.Fatalf("unexpected scopeLevel: %s", p.ScopeLevel)
	}
	if !p.IsActive {
		t.Fatal("expected isActive=true")
	}
	if !p.AllEndpoints {
		t.Fatal("expected allEndpoints=true")
	}
	if p.MaxRetries != 3 {
		t.Fatalf("expected maxRetries=3, got %d", p.MaxRetries)
	}
	if p.Priority != 1 {
		t.Fatalf("expected priority=1, got %d", p.Priority)
	}
	if p.Package.Build != "100" {
		t.Fatalf("unexpected package build: %s", p.Package.Build)
	}
	if p.Package.FileID != "200" {
		t.Fatalf("unexpected package fileId: %s", p.Package.FileID)
	}
	if p.Package.Major != "24" {
		t.Fatalf("unexpected package major: %s", p.Package.Major)
	}
	if p.Package.Raw == nil {
		t.Fatal("expected Package.Raw to be populated")
	}
	if p.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if total != 1 {
		t.Fatalf("expected totalItems=1, got %d", total)
	}
}

func TestUpgradePoliciesListSkipAlwaysSent(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("skip") != "0" {
			t.Fatalf("expected skip=0 even when zero-value, got %q", q.Get("skip"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"policies": []map[string]any{},
			},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.UpgradePoliciesList(context.Background(), &UpgradePolicyListParams{
		ScopeLevel: "account",
		OSType:     "linux",
		SortBy:     "priority",
		SortOrder:  "asc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpgradePoliciesListSkipNonZero(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("skip") != "10" {
			t.Fatalf("expected skip=10, got %q", q.Get("skip"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"policies": []map[string]any{},
			},
			"pagination": map[string]any{"totalItems": 50},
		})
	})
	c := testClient(t, handler)
	_, total, err := c.UpgradePoliciesList(context.Background(), &UpgradePolicyListParams{
		ScopeLevel: "site",
		OSType:     "windows",
		Limit:      10,
		Skip:       10,
		SortBy:     "priority",
		SortOrder:  "asc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 50 {
		t.Fatalf("expected totalItems=50, got %d", total)
	}
}

func TestUpgradePoliciesCreate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/upgrade-policy/policy" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "New Policy" {
			t.Fatalf("unexpected name: %v", body["name"])
		}
		if body["osType"] != "linux" {
			t.Fatalf("unexpected osType: %v", body["osType"])
		}
		if body["scopeLevel"] != "account" {
			t.Fatalf("unexpected scopeLevel: %v", body["scopeLevel"])
		}
		if body["isActive"] != true {
			t.Fatalf("unexpected isActive: %v", body["isActive"])
		}
		pkg, _ := body["package"].(map[string]any)
		if pkg["build"] != "200" {
			t.Fatalf("unexpected package build: %v", pkg["build"])
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	err := c.UpgradePoliciesCreate(context.Background(), UpgradePolicyCreate{
		Name:         "New Policy",
		OSType:       UpgradePolicyOSLinux,
		ScopeLevel:   UpgradePolicyScopeAccount,
		ScopeID:      "225494730938493804",
		IsActive:     true,
		AllEndpoints: true,
		MaxRetries:   3,
		Package:      UpgradePolicyPkg{Build: "200", FileID: "300", Major: "24", Minor: "2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpgradePoliciesUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/upgrade-policy/policy/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "Updated Policy" {
			t.Fatalf("unexpected name: %v", body["name"])
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	err := c.UpgradePoliciesUpdate(context.Background(), "1000000000000000001", UpgradePolicyCreate{
		Name:       "Updated Policy",
		OSType:     UpgradePolicyOSWindows,
		ScopeLevel: UpgradePolicyScopeSite,
		Package:    UpgradePolicyPkg{Build: "100", FileID: "200", Major: "24", Minor: "1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpgradePoliciesDelete(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/upgrade-policy/policy/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["action"] != "delete" {
			t.Fatalf("unexpected action: %v", body["action"])
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	err := c.UpgradePoliciesDelete(context.Background(), "1000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpgradePoliciesActivate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/upgrade-policy/policy/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["action"] != "activate" {
			t.Fatalf("unexpected action: %v", body["action"])
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	err := c.UpgradePoliciesActivate(context.Background(), "1000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpgradePoliciesDeactivate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/upgrade-policy/policy/1000000000000000001" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["action"] != "deactivate" {
			t.Fatalf("unexpected action: %v", body["action"])
		}
		w.WriteHeader(http.StatusOK)
	})
	c := testClient(t, handler)
	err := c.UpgradePoliciesDeactivate(context.Background(), "1000000000000000001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpgradePackagesList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/upgrade-policy/available-packages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("scopeLevel") != "account" {
			t.Fatalf("unexpected scopeLevel: %s", q.Get("scopeLevel"))
		}
		if q.Get("osType") != "macos" {
			t.Fatalf("unexpected osType: %s", q.Get("osType"))
		}
		if q.Get("displayName__contains") != "GA" {
			t.Fatalf("unexpected displayName__contains: %s", q.Get("displayName__contains"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"packages": []map[string]any{
					{
						"build": "100", "major": "24", "minor": "1",
						"displayName": "24.1 GA",
						"fileNames": []map[string]any{
							{"id": "file-1", "name": "SentinelOne_macos_v24_1_100.pkg"},
						},
					},
					{
						"build": "200", "major": "24", "minor": "2",
						"displayName": "24.2 GA",
						"fileNames":   []map[string]any{},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	pkgs, err := c.UpgradePackagesList(context.Background(), &UpgradePackageListParams{
		ScopeLevel:          "account",
		OSType:              "macos",
		DisplayNameContains: "GA",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}
	if pkgs[0].Build != "100" {
		t.Fatalf("unexpected build: %s", pkgs[0].Build)
	}
	if pkgs[0].DisplayName != "24.1 GA" {
		t.Fatalf("unexpected displayName: %s", pkgs[0].DisplayName)
	}
	if len(pkgs[0].FileNames) != 1 {
		t.Fatalf("expected 1 file, got %d", len(pkgs[0].FileNames))
	}
	if pkgs[0].FileNames[0].ID != "file-1" {
		t.Fatalf("unexpected file ID: %s", pkgs[0].FileNames[0].ID)
	}
	if pkgs[0].FileNames[0].Raw == nil {
		t.Fatal("expected FileNames[0].Raw to be populated")
	}
	if pkgs[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestUpgradePoliciesListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.UpgradePoliciesList(context.Background(), nil)
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

func TestUpgradePoliciesCreateError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 400, "title": "Bad Request"},
			},
		})
	})
	c := testClient(t, handler)
	err := c.UpgradePoliciesCreate(context.Background(), UpgradePolicyCreate{Name: "Bad"})
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 400 {
		t.Fatalf("expected 400, got %d", ae.Status)
	}
}

func TestUpgradePolicyEnumValues(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"OSLinux", string(UpgradePolicyOSLinux), "linux"},
		{"OSMacOS", string(UpgradePolicyOSMacOS), "macos"},
		{"OSWindows", string(UpgradePolicyOSWindows), "windows"},
		{"ScopeAccount", string(UpgradePolicyScopeAccount), "account"},
		{"ScopeGroup", string(UpgradePolicyScopeGroup), "group"},
		{"ScopeSite", string(UpgradePolicyScopeSite), "site"},
		{"ScopeTenant", string(UpgradePolicyScopeTenant), "tenant"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}
