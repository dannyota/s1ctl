package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// SettingsParams are query parameters for settings endpoints.
type SettingsParams struct {
	SiteIDs    []string
	AccountIDs []string
}

func (p *SettingsParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	return v
}

type settingsFilter struct {
	SiteIDs    []string `json:"siteIds,omitempty"`
	AccountIDs []string `json:"accountIds,omitempty"`
}

func (p *SettingsParams) filter() settingsFilter {
	if p == nil {
		return settingsFilter{}
	}
	return settingsFilter{SiteIDs: p.SiteIDs, AccountIDs: p.AccountIDs}
}

type settingsRequest struct {
	Data   any            `json:"data"`
	Filter settingsFilter `json:"filter"`
}

// NotificationConfig holds email, SMS, and syslog notification configuration.
type NotificationConfig struct {
	Email  string `json:"email"`
	SMS    string `json:"sms"`
	Syslog string `json:"syslog"`

	Raw json.RawMessage `json:"-"`
}

func (n *NotificationConfig) UnmarshalJSON(b []byte) error {
	type alias NotificationConfig
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// NotificationLastModified tracks who last modified notification settings.
type NotificationLastModified struct {
	UpdatedAt string `json:"updatedAt"`
	UpdatedBy string `json:"updatedBy"`

	Raw json.RawMessage `json:"-"`
}

func (n *NotificationLastModified) UnmarshalJSON(b []byte) error {
	type alias NotificationLastModified
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// NotificationSettings is the full notification settings object.
type NotificationSettings struct {
	Configurations NotificationConfig       `json:"configurations"`
	Notifications  json.RawMessage          `json:"notifications"`
	LastModified   NotificationLastModified `json:"lastModified"`

	Raw json.RawMessage `json:"-"`
}

func (n *NotificationSettings) UnmarshalJSON(b []byte) error {
	type alias NotificationSettings
	if err := json.Unmarshal(b, (*alias)(n)); err != nil {
		return err
	}
	n.Raw = append(n.Raw[:0:0], b...)
	return nil
}

// SSOSettings is the SSO configuration for a site or account.
type SSOSettings struct {
	Enabled                         bool     `json:"enabled"`
	IDPSsoURL                       string   `json:"idpSsoUrl"`
	IDPEntityID                     string   `json:"idpEntityId"`
	IDPCertName                     string   `json:"idpCertName"`
	SPAcsURL                        string   `json:"spAcsUrl"`
	SPEntityID                      string   `json:"spEntityId"`
	DefaultUserRole                 string   `json:"defaultUserRole"`
	DefaultUserRoleID               string   `json:"defaultUserRoleId"`
	AutoProvisioning                bool     `json:"autoProvisioning"`
	Domains                         []string `json:"domains"`
	SSOPropagateDomainsToChildren   bool     `json:"ssoPropagateDomainsToChildren"`
	SSOInheritDomainsFrom           []string `json:"ssoInheritDomainsFrom"`
	SSOElevatedSessionReauthType    string   `json:"ssoElevatedSessionReauthType"`
	SSOElevatedSessionReauthEnabled bool     `json:"ssoElevatedSessionReauthTypeEnabled"`
	SignRequest                     bool     `json:"signRequest"`

	Raw json.RawMessage `json:"-"`
}

func (s *SSOSettings) UnmarshalJSON(b []byte) error {
	type alias SSOSettings
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// SMTPSettings is the SMTP mail configuration for a site or account.
type SMTPSettings struct {
	Inherits     bool   `json:"inherits"`
	Enabled      bool   `json:"enabled"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Encryption   string `json:"encryption"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	NoReplyEmail string `json:"noReplyEmail"`

	Raw json.RawMessage `json:"-"`
}

func (s *SMTPSettings) UnmarshalJSON(b []byte) error {
	type alias SMTPSettings
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// SyslogSettings is the syslog forwarding configuration for a site or account.
type SyslogSettings struct {
	Enabled           bool   `json:"enabled"`
	Host              string `json:"host"`
	Port              int    `json:"port"`
	SSL               bool   `json:"ssl"`
	Format            string `json:"format"`
	ServerCertName    string `json:"serverCertName"`
	ServerCertContent string `json:"serverCertContent"`
	ClientCertName    string `json:"clientCertName"`
	ClientCertContent string `json:"clientCertContent"`
	ClientKeyName     string `json:"clientKeyName"`
	ClientKeyContent  string `json:"clientKeyContent"`
	Token             string `json:"token"`

	Raw json.RawMessage `json:"-"`
}

func (s *SyslogSettings) UnmarshalJSON(b []byte) error {
	type alias SyslogSettings
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// SettingsTestResult is the response from an SMTP or syslog test.
type SettingsTestResult struct {
	Status bool `json:"status"`

	Raw json.RawMessage `json:"-"`
}

func (s *SettingsTestResult) UnmarshalJSON(b []byte) error {
	type alias SettingsTestResult
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// SettingsNotificationsGet returns the notification settings.
func (c *Client) SettingsNotificationsGet(ctx context.Context, params *SettingsParams) (*NotificationSettings, error) {
	var resp singleResponse[NotificationSettings]
	if err := c.get(ctx, "/settings/notifications", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsNotificationsUpdate updates the notification settings.
func (c *Client) SettingsNotificationsUpdate(ctx context.Context, params *SettingsParams, data NotificationSettings) (*NotificationSettings, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[NotificationSettings]
	if err := c.put(ctx, "/settings/notifications", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSSOGet returns the SSO settings.
func (c *Client) SettingsSSOGet(ctx context.Context, params *SettingsParams) (*SSOSettings, error) {
	var resp singleResponse[SSOSettings]
	if err := c.get(ctx, "/settings/sso", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSSOUpdate updates the SSO settings.
func (c *Client) SettingsSSOUpdate(ctx context.Context, params *SettingsParams, data SSOSettings) (*SSOSettings, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[SSOSettings]
	if err := c.put(ctx, "/settings/sso", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSMTPGet returns the SMTP settings.
func (c *Client) SettingsSMTPGet(ctx context.Context, params *SettingsParams) (*SMTPSettings, error) {
	var resp singleResponse[SMTPSettings]
	if err := c.get(ctx, "/settings/smtp", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSMTPUpdate updates the SMTP settings.
func (c *Client) SettingsSMTPUpdate(ctx context.Context, params *SettingsParams, data SMTPSettings) (*SMTPSettings, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[SMTPSettings]
	if err := c.put(ctx, "/settings/smtp", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSMTPTest sends a test email using the provided SMTP settings.
func (c *Client) SettingsSMTPTest(ctx context.Context, params *SettingsParams, data SMTPSettings) (*SettingsTestResult, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[SettingsTestResult]
	if err := c.post(ctx, "/settings/smtp/test", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSyslogGet returns the syslog settings.
func (c *Client) SettingsSyslogGet(ctx context.Context, params *SettingsParams) (*SyslogSettings, error) {
	var resp singleResponse[SyslogSettings]
	if err := c.get(ctx, "/settings/syslog", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSyslogUpdate updates the syslog settings.
func (c *Client) SettingsSyslogUpdate(ctx context.Context, params *SettingsParams, data SyslogSettings) (*SyslogSettings, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[SyslogSettings]
	if err := c.put(ctx, "/settings/syslog", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSyslogTest sends a test message using the provided syslog settings.
func (c *Client) SettingsSyslogTest(ctx context.Context, params *SettingsParams, data SyslogSettings) (*SettingsTestResult, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[SettingsTestResult]
	if err := c.post(ctx, "/settings/syslog/test", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
