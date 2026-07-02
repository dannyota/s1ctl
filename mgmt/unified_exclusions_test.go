package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"testing"
)

func TestUnifiedExclusionsList(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/unified-exclusions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		if got := q["osTypes"]; !slices.Equal(got, []string{"windows", "linux"}) {
			t.Fatalf("unexpected osTypes: %v", got)
		}
		if q.Get("sortBy") != "exclusionName" {
			t.Fatalf("unexpected sortBy: %s", q.Get("sortBy"))
		}
		if q.Get("sortOrder") != "asc" {
			t.Fatalf("unexpected sortOrder: %s", q.Get("sortOrder"))
		}
		if q.Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", q.Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":                "1000000000000000001",
					"exclusionName":     "Test Exclusion",
					"osType":            "windows",
					"threatType":        "EDR",
					"modeType":          "suppression",
					"interactionLevel":  "disable_all_monitors",
					"type":              "path",
					"pathExclusionType": "folder",
					"source":            "user",
					"scopeName":         "Default Site",
					"description":       "A test exclusion",
					"childProcess":      true,
					"hits30d":           5,
					"createdAt":         "2025-01-01T00:00:00Z",
				},
			},
			"pagination": map[string]any{"totalItems": 1, "nextCursor": "abc"},
		})
	})
	c := testClient(t, handler)
	excls, pag, err := c.UnifiedExclusionsList(context.Background(), &UnifiedExclusionListParams{
		SiteIDs:   []string{"225494730938493804"},
		OSTypes:   []string{"windows", "linux"},
		SortBy:    "exclusionName",
		SortOrder: "asc",
		Limit:     10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(excls) != 1 {
		t.Fatalf("expected 1 exclusion, got %d", len(excls))
	}
	e := excls[0]
	if e.ID != "1000000000000000001" {
		t.Fatalf("unexpected ID: %s", e.ID)
	}
	if e.ExclusionName != "Test Exclusion" {
		t.Fatalf("unexpected exclusionName: %s", e.ExclusionName)
	}
	if e.OSType != "windows" {
		t.Fatalf("unexpected osType: %s", e.OSType)
	}
	if e.ThreatType != "EDR" {
		t.Fatalf("unexpected threatType: %s", e.ThreatType)
	}
	if e.ModeType != "suppression" {
		t.Fatalf("unexpected modeType: %s", e.ModeType)
	}
	if e.InteractionLevel != "disable_all_monitors" {
		t.Fatalf("unexpected interactionLevel: %s", e.InteractionLevel)
	}
	if e.PathExclusionType != "folder" {
		t.Fatalf("unexpected pathExclusionType: %s", e.PathExclusionType)
	}
	if e.Source != "user" {
		t.Fatalf("unexpected source: %s", e.Source)
	}
	if !e.ChildProcess {
		t.Fatal("expected childProcess=true")
	}
	if e.Hits30d != 5 {
		t.Fatalf("expected hits30d=5, got %d", e.Hits30d)
	}
	if e.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
	if pag.NextCursor != "abc" {
		t.Fatalf("unexpected cursor: %s", pag.NextCursor)
	}
}

func TestUnifiedExclusionsListNilParams(t *testing.T) {
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
	excls, pag, err := c.UnifiedExclusionsList(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(excls) != 0 {
		t.Fatalf("expected 0 exclusions, got %d", len(excls))
	}
	if pag.TotalItems != 0 {
		t.Fatalf("expected totalItems=0, got %d", pag.TotalItems)
	}
}

func TestUnifiedExclusionsListBoolParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("includeParents") != "true" {
			t.Fatalf("expected includeParents=true, got %s", q.Get("includeParents"))
		}
		if q.Get("includeChildren") != "false" {
			t.Fatalf("expected includeChildren=false, got %s", q.Get("includeChildren"))
		}
		if q.Get("imported") != "true" {
			t.Fatalf("expected imported=true, got %s", q.Get("imported"))
		}
		if q.Get("tenant") != "false" {
			t.Fatalf("expected tenant=false, got %s", q.Get("tenant"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	bTrue, bFalse := true, false
	_, _, err := c.UnifiedExclusionsList(context.Background(), &UnifiedExclusionListParams{
		IncludeParents:  &bTrue,
		IncludeChildren: &bFalse,
		Imported:        &bTrue,
		Tenant:          &bFalse,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnifiedExclusionsCount(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		q := r.URL.Query()
		if q.Get("countOnly") != "true" {
			t.Fatalf("expected countOnly=true, got %s", q.Get("countOnly"))
		}
		if got := q["osTypes"]; !slices.Equal(got, []string{"linux"}) {
			t.Fatalf("unexpected osTypes: %v", got)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 42},
		})
	})
	c := testClient(t, handler)
	count, err := c.UnifiedExclusionsCount(context.Background(), &UnifiedExclusionListParams{
		OSTypes: []string{"linux"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 42 {
		t.Fatalf("expected 42, got %d", count)
	}
}

func TestUnifiedExclusionsCountNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("countOnly") != "true" {
			t.Fatal("expected countOnly=true")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []map[string]any{},
			"pagination": map[string]any{"totalItems": 7},
		})
	})
	c := testClient(t, handler)
	count, err := c.UnifiedExclusionsCount(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 7 {
		t.Fatalf("expected 7, got %d", count)
	}
}

func TestUnifiedExclusionsCreate(t *testing.T) {
	var lastFilter map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/unified-exclusions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   map[string]any `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["exclusionName"] != "New Exclusion" {
			t.Fatalf("unexpected exclusionName: %v", body.Data["exclusionName"])
		}
		lastFilter = body.Filter
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":                "1000000000000000002",
					"exclusionName":     "New Exclusion",
					"osType":            "windows",
					"threatType":        "EDR",
					"modeType":          "suppression",
					"type":              "path",
					"pathExclusionType": "folder",
				},
			},
		})
	})
	c := testClient(t, handler)

	excl, err := c.UnifiedExclusionsCreate(context.Background(),
		UnifiedExclusionScope{
			ScopeLevel:   UnifiedExclusionScopeAccount,
			ScopeLevelID: nil,
		},
		UnifiedExclusionCreate{
			ExclusionName:     "New Exclusion",
			OSType:            UnifiedExclusionOSWindows,
			ThreatType:        UnifiedExclusionThreatEDR,
			ModeType:          UnifiedExclusionModeSuppression,
			Type:              UnifiedExclusionTypePath,
			PathExclusionType: UnifiedExclusionPathFolder,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error (account scope): %v", err)
	}
	if lastFilter["scopeLevel"] != "account" {
		t.Fatalf("expected account scope, got %v", lastFilter["scopeLevel"])
	}
	_ = excl

	scopeID := int64(225494730938493804)
	excl, err = c.UnifiedExclusionsCreate(context.Background(),
		UnifiedExclusionScope{
			ScopeLevel:   UnifiedExclusionScopeSite,
			ScopeLevelID: &scopeID,
		},
		UnifiedExclusionCreate{
			ExclusionName:     "New Exclusion",
			OSType:            UnifiedExclusionOSWindows,
			ThreatType:        UnifiedExclusionThreatEDR,
			ModeType:          UnifiedExclusionModeSuppression,
			Type:              UnifiedExclusionTypePath,
			PathExclusionType: UnifiedExclusionPathFolder,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if excl.ID != "1000000000000000002" {
		t.Fatalf("unexpected ID: %s", excl.ID)
	}
	if excl.ExclusionName != "New Exclusion" {
		t.Fatalf("unexpected exclusionName: %s", excl.ExclusionName)
	}
	if excl.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestUnifiedExclusionsCreateScopeLevelIDNil(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body.Filter["scopeLevelId"]; ok {
			t.Fatal("expected scopeLevelId to be omitted when nil")
		}
		if body.Filter["scopeLevel"] != "tenant" {
			t.Fatalf("unexpected scopeLevel: %v", body.Filter["scopeLevel"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": "1000000000000000003", "exclusionName": "Tenant Exclusion"},
			},
		})
	})
	c := testClient(t, handler)
	excl, err := c.UnifiedExclusionsCreate(context.Background(),
		UnifiedExclusionScope{ScopeLevel: UnifiedExclusionScopeTenant},
		UnifiedExclusionCreate{
			ExclusionName: "Tenant Exclusion",
			OSType:        UnifiedExclusionOSLinux,
			ThreatType:    UnifiedExclusionThreatEDR,
			ModeType:      UnifiedExclusionModeAll,
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if excl.ID != "1000000000000000003" {
		t.Fatalf("unexpected ID: %s", excl.ID)
	}
}

func TestUnifiedExclusionsExport(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/unified-exclusions/export" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q["siteIds"]; !slices.Equal(got, []string{"225494730938493804"}) {
			t.Fatalf("unexpected siteIds: %v", got)
		}
		w.Write([]byte(`[{"id":"1000000000000000001","exclusionName":"Exported"}]`))
	})
	c := testClient(t, handler)
	raw, err := c.UnifiedExclusionsExport(context.Background(), &UnifiedExclusionListParams{
		SiteIDs: []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if raw == nil {
		t.Fatal("expected non-nil raw response")
	}
	var items []map[string]any
	if err := json.Unmarshal(raw, &items); err != nil {
		t.Fatalf("failed to unmarshal export: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0]["exclusionName"] != "Exported" {
		t.Fatalf("unexpected exclusionName: %v", items[0]["exclusionName"])
	}
}

func TestUnifiedExclusionsListError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.UnifiedExclusionsList(context.Background(), nil)
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

func TestUnifiedExclusionsCountError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 500, "title": "Internal Server Error"},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.UnifiedExclusionsCount(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var ae *APIError
	if !errors.As(err, &ae) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if ae.Status != 500 {
		t.Fatalf("expected 500, got %d", ae.Status)
	}
}

func TestUnifiedExclusionsCreateError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 400, "title": "Bad Request", "detail": "invalid exclusion"},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.UnifiedExclusionsCreate(context.Background(),
		UnifiedExclusionScope{ScopeLevel: UnifiedExclusionScopeSite},
		UnifiedExclusionCreate{ExclusionName: "Bad"},
	)
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

func TestUnifiedExclusionEnumValues(t *testing.T) {
	// Verify enum constants match expected API values.
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"OSLinux", string(UnifiedExclusionOSLinux), "linux"},
		{"OSMacOS", string(UnifiedExclusionOSMacOS), "macos"},
		{"OSWindows", string(UnifiedExclusionOSWindows), "windows"},
		{"OSWindowsLegacy", string(UnifiedExclusionOSWindowsLegacy), "windows_legacy"},
		{"ThreatEDR", string(UnifiedExclusionThreatEDR), "EDR"},
		{"ThreatIDR", string(UnifiedExclusionThreatIDR), "IDR"},
		{"ModeAll", string(UnifiedExclusionModeAll), "all"},
		{"ModeSuppression", string(UnifiedExclusionModeSuppression), "suppression"},
		{"ModeAgentInterop", string(UnifiedExclusionModeAgentInteroperability), "agent_interoperability"},
		{"ModeBinaryVault", string(UnifiedExclusionModeBinaryVault), "binary_vault"},
		{"TypePath", string(UnifiedExclusionTypePath), "path"},
		{"TypeCertificate", string(UnifiedExclusionTypeCertificate), "certificate"},
		{"TypeBrowser", string(UnifiedExclusionTypeBrowser), "browser"},
		{"TypeFileType", string(UnifiedExclusionTypeFileType), "file_type"},
		{"TypeWhiteHash", string(UnifiedExclusionTypeWhiteHash), "white_hash"},
		{"TypeCommandline", string(UnifiedExclusionTypeCommandline), "commandline"},
		{"TypeContainerNative", string(UnifiedExclusionTypeContainerNative), "container_native"},
		{"InteractionDisableProcess", string(UnifiedExclusionInteractionDisableInProcessMonitor), "disable_in_process_monitor"},
		{"InteractionDisableAll", string(UnifiedExclusionInteractionDisableAllMonitors), "disable_all_monitors"},
		{"InteractionIdentityOnly", string(UnifiedExclusionInteractionIdentityOnly), "identity_only"},
		{"SourceUser", string(UnifiedExclusionSourceUser), "user"},
		{"SourceActionFromThreat", string(UnifiedExclusionSourceActionFromThreat), "action_from_threat"},
		{"SourceCatalog", string(UnifiedExclusionSourceCatalog), "catalog"},
		{"SourcePerformanceInsight", string(UnifiedExclusionSourcePerformanceInsight), "performance_insight"},
		{"PathFile", string(UnifiedExclusionPathFile), "file"},
		{"PathFolder", string(UnifiedExclusionPathFolder), "folder"},
		{"PathSubfolders", string(UnifiedExclusionPathSubfolders), "subfolders"},
		{"ScopeGroup", string(UnifiedExclusionScopeGroup), "group"},
		{"ScopeSite", string(UnifiedExclusionScopeSite), "site"},
		{"ScopeAccount", string(UnifiedExclusionScopeAccount), "account"},
		{"ScopeTenant", string(UnifiedExclusionScopeTenant), "tenant"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, tt.got)
			}
		})
	}
}
