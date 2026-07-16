package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestRangerADAssessmentStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/ranger-ad/assessment-status" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"status": "COMPLETED",
				"domainWiseCurrentStatusList": []map[string]any{
					{
						"domainName": "corp.example.com", "forestName": "example.com",
						"totalJobs": 10, "completedJobs": 10, "domainCompletedStatus": true,
					},
				},
				"tenantWiseCurrentStatusList": []map[string]any{
					{
						"tenantId":  "00000000-0000-0000-0000-000000000000",
						"totalJobs": 5, "completedJobs": 5, "tenantCompletedStatus": true,
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	status, err := c.RangerADAssessmentStatus(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Status != AssessmentStatusCompleted {
		t.Fatalf("expected COMPLETED, got %s", status.Status)
	}
	if len(status.Domains) != 1 {
		t.Fatalf("expected 1 domain, got %d", len(status.Domains))
	}
	d := status.Domains[0]
	if d.DomainName != "corp.example.com" {
		t.Fatalf("unexpected domain: %s", d.DomainName)
	}
	if d.ForestName != "example.com" {
		t.Fatalf("unexpected forest: %s", d.ForestName)
	}
	if d.TotalJobs != 10 || d.CompletedJobs != 10 {
		t.Fatalf("unexpected jobs: total=%d completed=%d", d.TotalJobs, d.CompletedJobs)
	}
	if !d.DomainCompleted {
		t.Fatal("expected domainCompleted=true")
	}
	if len(status.Tenants) != 1 {
		t.Fatalf("expected 1 tenant, got %d", len(status.Tenants))
	}
	tn := status.Tenants[0]
	if tn.TenantID != "00000000-0000-0000-0000-000000000000" {
		t.Fatalf("unexpected tenantId: %s", tn.TenantID)
	}
	if !tn.TenantCompleted {
		t.Fatal("expected tenantCompleted=true")
	}
	if status.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestRangerADAssessmentStatusParams(t *testing.T) {
	var gotQuery string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"status": "PENDING"},
		})
	})
	c := testClient(t, handler)
	_, err := c.RangerADAssessmentStatus(context.Background(), &ADAssessmentStatusParams{
		SiteIDs:    "100",
		AccountIDs: "200",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"siteIds=100", "accountIds=200"} {
		if !strings.Contains(gotQuery, want) {
			t.Errorf("query %q missing %q", gotQuery, want)
		}
	}
}

func TestRangerADAssessmentStatusError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden"},
			},
		})
	})
	c := testClient(t, handler)
	_, err := c.RangerADAssessmentStatus(context.Background(), nil)
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

func TestRangerADExposures(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/ranger-ad/get-exposures") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter ADExposureFilter `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if len(body.Filter.Severity) != 1 || body.Filter.Severity[0] != "Critical" {
			t.Fatalf("unexpected severity filter: %v", body.Filter.Severity)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": "E1", "detectionId": 100,
					"detectionName": "Kerberoasting", "detectionStatus": "Vulnerable",
					"severity": "Critical", "source": "OnPremAD",
					"domainName": "corp.example.com", "forestName": "example.com",
					"vulnerableCount": 5, "acknowledged": false, "remediable": true,
					"runTimestamp": 1700000000,
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	exposures, pag, err := c.RangerADExposures(context.Background(), &ADExposureListParams{
		Filter: ADExposureFilter{Severity: []string{"Critical"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exposures) != 1 {
		t.Fatalf("expected 1 exposure, got %d", len(exposures))
	}
	e := exposures[0]
	if e.ID != "E1" {
		t.Fatalf("unexpected id: %s", e.ID)
	}
	if e.DetectionID != 100 {
		t.Fatalf("unexpected detectionId: %d", e.DetectionID)
	}
	if e.DetectionName != "Kerberoasting" {
		t.Fatalf("unexpected detectionName: %s", e.DetectionName)
	}
	if e.DetectionStatus != ExposureStatusVulnerable {
		t.Fatalf("unexpected detectionStatus: %s", e.DetectionStatus)
	}
	if e.Severity != ExposureSeverityCritical {
		t.Fatalf("unexpected severity: %s", e.Severity)
	}
	if e.Source != ExposureSourceOnPremAD {
		t.Fatalf("unexpected source: %s", e.Source)
	}
	if e.VulnerableCount != 5 {
		t.Fatalf("expected vulnerableCount=5, got %d", e.VulnerableCount)
	}
	if !e.Remediable {
		t.Fatal("expected remediable=true")
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
	if e.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestRangerADExposuresNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	exposures, pag, err := c.RangerADExposures(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(exposures) != 0 {
		t.Fatalf("expected 0 exposures, got %d", len(exposures))
	}
	if pag.TotalItems != 0 {
		t.Fatalf("expected totalItems=0, got %d", pag.TotalItems)
	}
}

func TestRangerADExposuresQueryParams(t *testing.T) {
	var gotQuery string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.RangerADExposures(context.Background(), &ADExposureListParams{
		Limit:   25,
		Skip:    10,
		SiteIDs: "100",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"limit=25", "skip=10", "siteIds=100"} {
		if !strings.Contains(gotQuery, want) {
			t.Errorf("query %q missing %q", gotQuery, want)
		}
	}
}

func TestRangerADExposuresError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	})
	c := testClient(t, handler)
	_, _, err := c.RangerADExposures(context.Background(), nil)
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

func TestRangerADAffectedObjects(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/ranger-ad/get-affected-objects") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter ADAffectedObjectFilter `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if len(body.Filter.DetectionName) != 1 || body.Filter.DetectionName[0] != "Kerberoasting" {
			t.Fatalf("unexpected detectionName filter: %v", body.Filter.DetectionName)
		}
		if len(body.Filter.DomainName) != 1 || body.Filter.DomainName[0] != "corp.example.com" {
			t.Fatalf("unexpected domainName filter: %v", body.Filter.DomainName)
		}
		dn := "CN=svc-app,OU=Service Accounts,DC=corp,DC=example,DC=com"
		sam := "svc-app"
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": 1, "runId": 42,
					"dn": dn, "samAccountName": sam,
					"objectType": "user", "accountStatus": "enabled",
				},
			},
			"pagination": map[string]any{"totalItems": 1},
		})
	})
	c := testClient(t, handler)
	objs, pag, err := c.RangerADAffectedObjects(context.Background(), &ADAffectedObjectListParams{
		Filter: ADAffectedObjectFilter{
			DetectionName: []string{"Kerberoasting"},
			DomainName:    []string{"corp.example.com"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(objs) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objs))
	}
	o := objs[0]
	if o.ID != 1 {
		t.Fatalf("unexpected id: %d", o.ID)
	}
	if o.RunID != 42 {
		t.Fatalf("unexpected runId: %d", o.RunID)
	}
	if o.SAMAccountName == nil || *o.SAMAccountName != "svc-app" {
		t.Fatalf("unexpected samAccountName: %v", o.SAMAccountName)
	}
	if pag.TotalItems != 1 {
		t.Fatalf("expected totalItems=1, got %d", pag.TotalItems)
	}
	if o.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestRangerADAffectedObjectsNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":       []any{},
			"pagination": map[string]any{"totalItems": 0},
		})
	})
	c := testClient(t, handler)
	objs, _, err := c.RangerADAffectedObjects(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(objs) != 0 {
		t.Fatalf("expected 0 objects, got %d", len(objs))
	}
}

func TestRangerADTriggerAssessment(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/ranger-ad/trigger-assessment") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter ADTriggerAssessmentFilter `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !body.Filter.IsFullScan {
			t.Fatal("expected isFullScan=true")
		}
		if len(body.Filter.DomainName) != 1 || body.Filter.DomainName[0] != "corp.example.com" {
			t.Fatalf("unexpected domainName: %v", body.Filter.DomainName)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"success": true,
				"message": "Assessment triggered",
			},
		})
	})
	c := testClient(t, handler)
	ok, msg, err := c.RangerADTriggerAssessment(context.Background(), &ADTriggerAssessmentParams{
		SiteIDs: "100",
		Filter: ADTriggerAssessmentFilter{
			IsFullScan: true,
			DomainName: []string{"corp.example.com"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected success=true")
	}
	if msg != "Assessment triggered" {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestRangerADTriggerAssessmentNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"success": true,
				"message": "OK",
			},
		})
	})
	c := testClient(t, handler)
	ok, _, err := c.RangerADTriggerAssessment(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected success=true")
	}
}

func TestRangerADTriggerAssessmentError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 400, "title": "Bad Request", "detail": "invalid scan source"},
			},
		})
	})
	c := testClient(t, handler)
	_, _, err := c.RangerADTriggerAssessment(context.Background(), nil)
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

func TestRangerADSetSkippedExposures(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/ranger-ad/set-skipped-exposures") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter ADSkipExposuresFilter `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !body.Filter.Skip {
			t.Fatal("expected skip=true")
		}
		if len(body.Filter.DetectionName) != 1 || body.Filter.DetectionName[0] != "Kerberoasting" {
			t.Fatalf("unexpected detectionName: %v", body.Filter.DetectionName)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true, "message": "Exposures skipped"},
		})
	})
	c := testClient(t, handler)
	ok, msg, err := c.RangerADSetSkippedExposures(context.Background(), &ADSkipExposuresParams{
		SiteIDs: "100",
		Filter: ADSkipExposuresFilter{
			DetectionName: []string{"Kerberoasting"},
			DomainName:    []string{"corp.example.com"},
			Skip:          true,
			SkipReason:    "accepted risk",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected success=true")
	}
	if msg != "Exposures skipped" {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestRangerADSetAckStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/ranger-ad/set-ack-status") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Filter ADAckExposuresFilter `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if !body.Filter.Acknowledged {
			t.Fatal("expected acknowledged=true")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"success": true, "message": "Status updated"},
		})
	})
	c := testClient(t, handler)
	ok, msg, err := c.RangerADSetAckStatus(context.Background(), &ADAckExposuresParams{
		SiteIDs: "100",
		Filter: ADAckExposuresFilter{
			DetectionName: []string{"Kerberoasting"},
			DomainName:    []string{"corp.example.com"},
			Acknowledged:  true,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected success=true")
	}
	if msg != "Status updated" {
		t.Fatalf("unexpected message: %s", msg)
	}
}
