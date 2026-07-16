package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestIdentityADConfigurations(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/adConfigurations") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id": 1, "tenantId": "T1", "cloudlinkId": 10,
					"domainName": "corp.example.com", "domainControllerFqdn": "dc1.corp.example.com",
					"enabled": true, "assessmentStatus": "COMPLETED",
					"encryptionMethod": "LDAPS", "username": "svc-bind",
					"isConnected": true, "ldapReferral": false,
					"featuresOpted": []string{"RANGER_AD"},
					"featureStatusInfo": []map[string]any{
						{"featureType": "AD_ASSESSMENT", "status": "COMPLETED"},
					},
				},
			},
		})
	})
	c := testClient(t, handler)
	configs, err := c.IdentityADConfigurations(context.Background(), &IdentityParams{SiteIDs: "100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(configs) != 1 {
		t.Fatalf("expected 1 config, got %d", len(configs))
	}
	cfg := configs[0]
	if cfg.ID != 1 {
		t.Fatalf("unexpected id: %d", cfg.ID)
	}
	if cfg.DomainName != "corp.example.com" {
		t.Fatalf("unexpected domain: %s", cfg.DomainName)
	}
	if cfg.EncryptionMethod != EncryptionMethodLDAPS {
		t.Fatalf("unexpected encryption: %s", cfg.EncryptionMethod)
	}
	if cfg.Username != "svc-bind" {
		t.Fatalf("unexpected username: %s", cfg.Username)
	}
	if cfg.AssessmentStatus != ADConfigAssessmentCompleted {
		t.Fatalf("unexpected assessment status: %s", cfg.AssessmentStatus)
	}
	if !cfg.IsConnected {
		t.Fatal("expected isConnected=true")
	}
	if len(cfg.FeatureStatusInfo) != 1 {
		t.Fatalf("expected 1 feature status, got %d", len(cfg.FeatureStatusInfo))
	}
	if cfg.FeatureStatusInfo[0].FeatureType != ADFeatureTypeADAssessment {
		t.Fatalf("unexpected feature type: %s", cfg.FeatureStatusInfo[0].FeatureType)
	}
	if cfg.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestIdentityADConfigurationsParams(t *testing.T) {
	var gotQuery string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}})
	})
	c := testClient(t, handler)
	_, err := c.IdentityADConfigurations(context.Background(), &IdentityParams{
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

func TestIdentityADConfigurationsError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{{"code": 403, "title": "Forbidden"}},
		})
	})
	c := testClient(t, handler)
	_, err := c.IdentityADConfigurations(context.Background(), nil)
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

func TestIdentityADConfigurationAdd(t *testing.T) {
	var gotBody map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/addAdConfiguration") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{"data": nil})
	})
	c := testClient(t, handler)
	err := c.IdentityADConfigurationAdd(context.Background(), &IdentityParams{SiteIDs: "100"}, ADConfigurationInput{
		DomainName:           "corp.example.com",
		DomainControllerFqdn: "dc1.corp.example.com",
		UserName:             "svc-bind",
		Password:             "secret",
		EncryptionMethod:     EncryptionMethodLDAPS,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, ok := gotBody["input"].(map[string]any)
	if !ok {
		t.Fatal("expected input in request body")
	}
	if input["domainName"] != "corp.example.com" {
		t.Fatalf("unexpected domainName: %v", input["domainName"])
	}
	if input["password"] != "secret" {
		t.Fatalf("expected password in request body")
	}
}

func TestIdentityADConfigurationDelete(t *testing.T) {
	var gotBody map[string]any
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/deleteAdConfiguration") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&gotBody)
		json.NewEncoder(w).Encode(map[string]any{"data": nil})
	})
	c := testClient(t, handler)
	err := c.IdentityADConfigurationDelete(context.Background(), &IdentityParams{SiteIDs: "100"}, []int64{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, ok := gotBody["input"].([]any)
	if !ok {
		t.Fatal("expected input array in request body")
	}
	if len(input) != 2 {
		t.Fatalf("expected 2 IDs, got %d", len(input))
	}
}

func TestIdentityAvailableFeatures(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"featureName": "RANGER_AD", "available": true},
				{"featureName": "SINGULARITY_IDENTITY", "available": false},
			},
		})
	})
	c := testClient(t, handler)
	features, err := c.IdentityAvailableFeatures(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(features))
	}
	if features[0].FeatureName != ADFeatureRangerAD {
		t.Fatalf("unexpected feature: %s", features[0].FeatureName)
	}
	if !features[0].Available {
		t.Fatal("expected RANGER_AD to be available")
	}
	if features[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestIdentityDomains(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"domain": "corp.example.com", "parentDomain": "example.com", "root": false},
				{"domain": "example.com", "parentDomain": "", "root": true},
			},
		})
	})
	c := testClient(t, handler)
	domains, err := c.IdentityDomains(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(domains) != 2 {
		t.Fatalf("expected 2 domains, got %d", len(domains))
	}
	if domains[0].Domain != "corp.example.com" {
		t.Fatalf("unexpected domain: %s", domains[0].Domain)
	}
	if domains[1].Root != true {
		t.Fatal("expected root=true for second domain")
	}
	if domains[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestIdentityTimezones(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"timeZoneId": "UTC", "displayName": "Coordinated Universal Time"},
				{"timeZoneId": "US/Eastern", "displayName": "Eastern Time"},
			},
		})
	})
	c := testClient(t, handler)
	tzs, err := c.IdentityTimezones(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tzs) != 2 {
		t.Fatalf("expected 2 timezones, got %d", len(tzs))
	}
	if tzs[0].TimeZoneID != "UTC" {
		t.Fatalf("unexpected timezone: %s", tzs[0].TimeZoneID)
	}
	if tzs[0].Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestIdentityOnboardingStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/getOnboardingStatus") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"status":          "COMPLETE",
				"featureSelected": []string{"RANGER_AD"},
				"adConnector":     "CONFIGURED",
				"domainName":      "corp.example.com",
			},
		})
	})
	c := testClient(t, handler)
	status, err := c.IdentityOnboardingStatus(context.Background(), &IdentityParams{SiteIDs: "100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Status != OnboardingStatusComplete {
		t.Fatalf("unexpected status: %s", status.Status)
	}
	if status.ADConnector != ADConnectorConfigured {
		t.Fatalf("unexpected connector: %s", status.ADConnector)
	}
	if status.DomainName != "corp.example.com" {
		t.Fatalf("unexpected domain: %s", status.DomainName)
	}
	if len(status.FeatureSelected) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(status.FeatureSelected))
	}
	if status.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}
