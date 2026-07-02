package mgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestSettingsNotificationsGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/settings/notifications" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("siteIds") != "225494730938493804" {
			t.Fatalf("unexpected siteIds: %s", q.Get("siteIds"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"configurations": map[string]any{
					"email": "admin@example.com", "sms": "+10000000000", "syslog": "syslog.example.com",
				},
				"notifications": map[string]any{"threatDetected": true},
				"lastModified": map[string]any{
					"updatedAt": "2025-01-01T00:00:00Z", "updatedBy": "admin",
				},
			},
		})
	})
	c := testClient(t, handler)
	ns, err := c.SettingsNotificationsGet(context.Background(), &SettingsParams{
		SiteIDs: []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ns.Configurations.Email != "admin@example.com" {
		t.Fatalf("unexpected email: %s", ns.Configurations.Email)
	}
	if ns.Configurations.SMS != "+10000000000" {
		t.Fatalf("unexpected sms: %s", ns.Configurations.SMS)
	}
	if ns.Configurations.Syslog != "syslog.example.com" {
		t.Fatalf("unexpected syslog: %s", ns.Configurations.Syslog)
	}
	if ns.Configurations.Raw == nil {
		t.Fatal("expected Configurations.Raw to be populated")
	}
	if ns.LastModified.UpdatedBy != "admin" {
		t.Fatalf("unexpected updatedBy: %s", ns.LastModified.UpdatedBy)
	}
	if ns.LastModified.Raw == nil {
		t.Fatal("expected LastModified.Raw to be populated")
	}
	if ns.Notifications == nil {
		t.Fatal("expected Notifications to be populated")
	}
	if ns.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsNotificationsUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/settings/notifications" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   map[string]any `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		cfg, _ := body.Data["configurations"].(map[string]any)
		if cfg["email"] != "new@example.com" {
			t.Fatalf("unexpected email in body: %v", cfg["email"])
		}
		siteIDs, _ := body.Filter["siteIds"].([]any)
		if len(siteIDs) != 1 || siteIDs[0] != "225494730938493804" {
			t.Fatalf("unexpected filter siteIds: %v", body.Filter["siteIds"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"configurations": map[string]any{"email": "new@example.com", "sms": "", "syslog": ""},
				"notifications":  map[string]any{},
				"lastModified":   map[string]any{"updatedAt": "2025-06-01T00:00:00Z", "updatedBy": "admin"},
			},
		})
	})
	c := testClient(t, handler)
	ns, err := c.SettingsNotificationsUpdate(context.Background(),
		&SettingsParams{SiteIDs: []string{"225494730938493804"}},
		NotificationSettings{
			Configurations: NotificationConfig{Email: "new@example.com"},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ns.Configurations.Email != "new@example.com" {
		t.Fatalf("unexpected email: %s", ns.Configurations.Email)
	}
	if ns.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsSSOGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/settings/sso" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"enabled":          true,
				"idpSsoUrl":        "https://idp.example.com/sso",
				"idpEntityId":      "idp-entity-123",
				"spAcsUrl":         "https://sp.example.com/acs",
				"spEntityId":       "sp-entity-456",
				"defaultUserRole":  "viewer",
				"autoProvisioning": true,
				"domains":          []string{"example.com"},
				"signRequest":      true,
			},
		})
	})
	c := testClient(t, handler)
	sso, err := c.SettingsSSOGet(context.Background(), &SettingsParams{
		AccountIDs: []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sso.Enabled {
		t.Fatal("expected enabled=true")
	}
	if sso.IDPSsoURL != "https://idp.example.com/sso" {
		t.Fatalf("unexpected idpSsoUrl: %s", sso.IDPSsoURL)
	}
	if sso.IDPEntityID != "idp-entity-123" {
		t.Fatalf("unexpected idpEntityId: %s", sso.IDPEntityID)
	}
	if sso.SPAcsURL != "https://sp.example.com/acs" {
		t.Fatalf("unexpected spAcsUrl: %s", sso.SPAcsURL)
	}
	if sso.DefaultUserRole != "viewer" {
		t.Fatalf("unexpected defaultUserRole: %s", sso.DefaultUserRole)
	}
	if !sso.AutoProvisioning {
		t.Fatal("expected autoProvisioning=true")
	}
	if len(sso.Domains) != 1 || sso.Domains[0] != "example.com" {
		t.Fatalf("unexpected domains: %v", sso.Domains)
	}
	if !sso.SignRequest {
		t.Fatal("expected signRequest=true")
	}
	if sso.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsSSOUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/settings/sso" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   map[string]any `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["enabled"] != true {
			t.Fatalf("unexpected enabled: %v", body.Data["enabled"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"enabled":     true,
				"idpSsoUrl":   "https://idp.example.com/sso",
				"idpEntityId": "idp-entity-123",
			},
		})
	})
	c := testClient(t, handler)
	sso, err := c.SettingsSSOUpdate(context.Background(),
		&SettingsParams{SiteIDs: []string{"225494730938493804"}},
		SSOSettings{Enabled: true, IDPSsoURL: "https://idp.example.com/sso"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sso.Enabled {
		t.Fatal("expected enabled=true")
	}
	if sso.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsSMTPGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/settings/smtp" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"inherits":     false,
				"enabled":      true,
				"host":         "smtp.example.com",
				"port":         587,
				"encryption":   "tls",
				"username":     "mailuser",
				"noReplyEmail": "noreply@example.com",
			},
		})
	})
	c := testClient(t, handler)
	smtp, err := c.SettingsSMTPGet(context.Background(), &SettingsParams{
		SiteIDs: []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !smtp.Enabled {
		t.Fatal("expected enabled=true")
	}
	if smtp.Host != "smtp.example.com" {
		t.Fatalf("unexpected host: %s", smtp.Host)
	}
	if smtp.Port != 587 {
		t.Fatalf("expected port=587, got %d", smtp.Port)
	}
	if smtp.Encryption != "tls" {
		t.Fatalf("unexpected encryption: %s", smtp.Encryption)
	}
	if smtp.Username != "mailuser" {
		t.Fatalf("unexpected username: %s", smtp.Username)
	}
	if smtp.NoReplyEmail != "noreply@example.com" {
		t.Fatalf("unexpected noReplyEmail: %s", smtp.NoReplyEmail)
	}
	if smtp.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsSMTPUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/settings/smtp" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   map[string]any `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["host"] != "new-smtp.example.com" {
			t.Fatalf("unexpected host: %v", body.Data["host"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"enabled": true, "host": "new-smtp.example.com", "port": 465,
			},
		})
	})
	c := testClient(t, handler)
	smtp, err := c.SettingsSMTPUpdate(context.Background(),
		&SettingsParams{SiteIDs: []string{"225494730938493804"}},
		SMTPSettings{Enabled: true, Host: "new-smtp.example.com", Port: 465},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if smtp.Host != "new-smtp.example.com" {
		t.Fatalf("unexpected host: %s", smtp.Host)
	}
	if smtp.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsSMTPTest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/settings/smtp/test" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   map[string]any `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["host"] != "smtp.example.com" {
			t.Fatalf("unexpected host: %v", body.Data["host"])
		}
		siteIDs, _ := body.Filter["siteIds"].([]any)
		if len(siteIDs) != 1 || siteIDs[0] != "225494730938493804" {
			t.Fatalf("unexpected filter siteIds: %v", body.Filter["siteIds"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"status": true},
		})
	})
	c := testClient(t, handler)
	result, err := c.SettingsSMTPTest(context.Background(),
		&SettingsParams{SiteIDs: []string{"225494730938493804"}},
		SMTPSettings{Host: "smtp.example.com", Port: 587},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Status {
		t.Fatal("expected status=true")
	}
	if result.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsSyslogGet(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/settings/syslog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"enabled": true,
				"host":    "syslog.example.com",
				"port":    514,
				"ssl":     true,
				"format":  "CEF",
				"token":   "tok-123",
			},
		})
	})
	c := testClient(t, handler)
	syslog, err := c.SettingsSyslogGet(context.Background(), &SettingsParams{
		SiteIDs: []string{"225494730938493804"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !syslog.Enabled {
		t.Fatal("expected enabled=true")
	}
	if syslog.Host != "syslog.example.com" {
		t.Fatalf("unexpected host: %s", syslog.Host)
	}
	if syslog.Port != 514 {
		t.Fatalf("expected port=514, got %d", syslog.Port)
	}
	if !syslog.SSL {
		t.Fatal("expected ssl=true")
	}
	if syslog.Format != "CEF" {
		t.Fatalf("unexpected format: %s", syslog.Format)
	}
	if syslog.Token != "tok-123" {
		t.Fatalf("unexpected token: %s", syslog.Token)
	}
	if syslog.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsSyslogUpdate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/settings/syslog" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   map[string]any `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["host"] != "new-syslog.example.com" {
			t.Fatalf("unexpected host: %v", body.Data["host"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"enabled": true, "host": "new-syslog.example.com", "port": 6514,
				"ssl": true, "format": "CEF",
			},
		})
	})
	c := testClient(t, handler)
	syslog, err := c.SettingsSyslogUpdate(context.Background(),
		&SettingsParams{SiteIDs: []string{"225494730938493804"}},
		SyslogSettings{Enabled: true, Host: "new-syslog.example.com", Port: 6514, SSL: true, Format: "CEF"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if syslog.Host != "new-syslog.example.com" {
		t.Fatalf("unexpected host: %s", syslog.Host)
	}
	if syslog.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsSyslogTest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/settings/syslog/test" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var body struct {
			Data   map[string]any `json:"data"`
			Filter map[string]any `json:"filter"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Data["host"] != "syslog.example.com" {
			t.Fatalf("unexpected host: %v", body.Data["host"])
		}
		accountIDs, _ := body.Filter["accountIds"].([]any)
		if len(accountIDs) != 1 || accountIDs[0] != "225494730938493804" {
			t.Fatalf("unexpected filter accountIds: %v", body.Filter["accountIds"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{"status": true},
		})
	})
	c := testClient(t, handler)
	result, err := c.SettingsSyslogTest(context.Background(),
		&SettingsParams{AccountIDs: []string{"225494730938493804"}},
		SyslogSettings{Host: "syslog.example.com", Port: 514},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Status {
		t.Fatal("expected status=true")
	}
	if result.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsNilParams(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("expected no query params, got %q", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"enabled": false, "host": "", "port": 0,
			},
		})
	})
	c := testClient(t, handler)
	smtp, err := c.SettingsSMTPGet(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if smtp.Enabled {
		t.Fatal("expected enabled=false")
	}
	if smtp.Raw == nil {
		t.Fatal("expected Raw to be populated")
	}
}

func TestSettingsGetError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 403, "title": "Forbidden"},
			},
		})
	})
	c := testClient(t, handler)
	tests := []struct {
		name string
		call func() error
	}{
		{"notifications", func() error {
			_, err := c.SettingsNotificationsGet(context.Background(), nil)
			return err
		}},
		{"sso", func() error {
			_, err := c.SettingsSSOGet(context.Background(), nil)
			return err
		}},
		{"smtp", func() error {
			_, err := c.SettingsSMTPGet(context.Background(), nil)
			return err
		}},
		{"syslog", func() error {
			_, err := c.SettingsSyslogGet(context.Background(), nil)
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
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
		})
	}
}

func TestSettingsUpdateError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"code": 400, "title": "Bad Request"},
			},
		})
	})
	c := testClient(t, handler)
	tests := []struct {
		name string
		call func() error
	}{
		{"notifications", func() error {
			_, err := c.SettingsNotificationsUpdate(context.Background(), nil, NotificationSettings{})
			return err
		}},
		{"sso", func() error {
			_, err := c.SettingsSSOUpdate(context.Background(), nil, SSOSettings{})
			return err
		}},
		{"smtp", func() error {
			_, err := c.SettingsSMTPUpdate(context.Background(), nil, SMTPSettings{})
			return err
		}},
		{"syslog", func() error {
			_, err := c.SettingsSyslogUpdate(context.Background(), nil, SyslogSettings{})
			return err
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
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
		})
	}
}
