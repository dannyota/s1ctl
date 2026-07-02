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

// SMSSettings is the SMS notification service configuration for a site or
// account. Per the spec the only field the API exposes is the enabled flag.
type SMSSettings struct {
	Enabled bool `json:"enabled"`

	Raw json.RawMessage `json:"-"`
}

func (s *SMSSettings) UnmarshalJSON(b []byte) error {
	type alias SMSSettings
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// SettingsSMSGet returns the SMS settings.
func (c *Client) SettingsSMSGet(ctx context.Context, params *SettingsParams) (*SMSSettings, error) {
	var resp singleResponse[SMSSettings]
	if err := c.get(ctx, "/settings/sms", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSMSUpdate updates the SMS settings.
func (c *Client) SettingsSMSUpdate(ctx context.Context, params *SettingsParams, data SMSSettings) (*SMSSettings, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[SMSSettings]
	if err := c.put(ctx, "/settings/sms", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// NotificationRecipient is a single notification recipient. The GET endpoint
// returns a list of these; PUT sets (creates or updates) one at a time.
type NotificationRecipient struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	SMS       string `json:"sms"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`

	Raw json.RawMessage `json:"-"`
}

func (r *NotificationRecipient) UnmarshalJSON(b []byte) error {
	type alias NotificationRecipient
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// recipientsData is the GET /settings/recipients response data envelope.
type recipientsData struct {
	Recipients []NotificationRecipient `json:"recipients"`
}

// SettingsRecipientsGet returns the configured notification recipients.
func (c *Client) SettingsRecipientsGet(ctx context.Context, params *SettingsParams) ([]NotificationRecipient, error) {
	var resp singleResponse[recipientsData]
	if err := c.get(ctx, "/settings/recipients", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data.Recipients, nil
}

// SettingsRecipientsUpdate sets (creates or updates) a single notification
// recipient and returns the stored recipient.
func (c *Client) SettingsRecipientsUpdate(ctx context.Context, params *SettingsParams, data NotificationRecipient) (*NotificationRecipient, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[NotificationRecipient]
	if err := c.put(ctx, "/settings/recipients", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsRecipientDelete removes a notification recipient by ID.
func (c *Client) SettingsRecipientDelete(ctx context.Context, id string) error {
	return c.delete(ctx, "/settings/recipients/"+url.PathEscape(id))
}

// ADSettings is the Active Directory integration configuration. Password is a
// secret (the bind account credential); the GET endpoint does not echo it, but
// it is accepted on update. Never print Password to --json or the audit log.
type ADSettings struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	RootDN   string `json:"rootDn"`
	SSL      bool   `json:"ssl"`
	Password string `json:"password"`

	Raw json.RawMessage `json:"-"`
}

func (a *ADSettings) UnmarshalJSON(b []byte) error {
	type alias ADSettings
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// SettingsADGet returns the Active Directory settings.
func (c *Client) SettingsADGet(ctx context.Context, params *SettingsParams) (*ADSettings, error) {
	var resp singleResponse[ADSettings]
	if err := c.get(ctx, "/settings/active-directory", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsADUpdate updates the Active Directory settings.
func (c *Client) SettingsADUpdate(ctx context.Context, params *SettingsParams, data ADSettings) (*ADSettings, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[ADSettings]
	if err := c.put(ctx, "/settings/active-directory", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsADTest probes connectivity using the provided Active Directory settings.
func (c *Client) SettingsADTest(ctx context.Context, params *SettingsParams, data ADSettings) (*SettingsTestResult, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[SettingsTestResult]
	if err := c.post(ctx, "/settings/active-directory/test", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ADScopeMapping maps Active Directory groups to admin and viewer scopes.
type ADScopeMapping struct {
	Admin  []string `json:"admin"`
	Viewer []string `json:"viewer"`

	Raw json.RawMessage `json:"-"`
}

func (m *ADScopeMapping) UnmarshalJSON(b []byte) error {
	type alias ADScopeMapping
	if err := json.Unmarshal(b, (*alias)(m)); err != nil {
		return err
	}
	m.Raw = append(m.Raw[:0:0], b...)
	return nil
}

// SettingsADScopeMappingGet returns the Active Directory scope mapping.
func (c *Client) SettingsADScopeMappingGet(ctx context.Context, params *SettingsParams) (*ADScopeMapping, error) {
	var resp singleResponse[ADScopeMapping]
	if err := c.get(ctx, "/settings/active-directory/scope-mapping", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsADScopeMappingUpdate updates the Active Directory scope mapping.
func (c *Client) SettingsADScopeMappingUpdate(ctx context.Context, params *SettingsParams, data ADScopeMapping) (*ADScopeMapping, error) {
	req := settingsRequest{Data: data, Filter: params.filter()}
	var resp singleResponse[ADScopeMapping]
	if err := c.put(ctx, "/settings/active-directory/scope-mapping", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SSOServiceProviderCert is the SAML service-provider signing certificate. The
// PEM is public key material (safe to print), not a secret.
type SSOServiceProviderCert struct {
	FileName  string `json:"fileName"`
	PEM       string `json:"pem"`
	IssuedAt  string `json:"issuedAt"`
	ExpiresAt string `json:"expiresAt"`

	Raw json.RawMessage `json:"-"`
}

func (s *SSOServiceProviderCert) UnmarshalJSON(b []byte) error {
	type alias SSOServiceProviderCert
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// SettingsSSOCert returns the SSO service-provider signing certificate metadata
// and PEM.
func (c *Client) SettingsSSOCert(ctx context.Context, params *SettingsParams) (*SSOServiceProviderCert, error) {
	var resp singleResponse[SSOServiceProviderCert]
	if err := c.get(ctx, "/settings/sso/sp-cert", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SettingsSSOCertDownload returns the raw SSO service-provider certificate file.
func (c *Client) SettingsSSOCertDownload(ctx context.Context, params *SettingsParams) ([]byte, error) {
	return c.getRaw(ctx, "/settings/sso/sp-cert/download", params.values())
}

// CancelPendingEmailsResult reports how many pending emails were cancelled.
type CancelPendingEmailsResult struct {
	Canceled int `json:"canceled"`

	Raw json.RawMessage `json:"-"`
}

func (r *CancelPendingEmailsResult) UnmarshalJSON(b []byte) error {
	type alias CancelPendingEmailsResult
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// SettingsCancelPendingEmails clears queued pending email notifications. The
// request body carries only a filter (no data envelope), per the spec.
func (c *Client) SettingsCancelPendingEmails(ctx context.Context, params *SettingsParams) (*CancelPendingEmailsResult, error) {
	req := struct {
		Filter settingsFilter `json:"filter"`
	}{Filter: params.filter()}
	var resp singleResponse[CancelPendingEmailsResult]
	if err := c.post(ctx, "/settings/notifications/cancel-pending-emails", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
