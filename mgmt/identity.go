package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// EncryptionMethod is the LDAP encryption method for an AD configuration.
type EncryptionMethod string

const (
	EncryptionMethodLDAP  EncryptionMethod = "LDAP"
	EncryptionMethodLDAPS EncryptionMethod = "LDAPS"
)

// ADConfigAssessmentStatus is the assessment status of an AD configuration.
type ADConfigAssessmentStatus string

const (
	ADConfigAssessmentPending    ADConfigAssessmentStatus = "PENDING"
	ADConfigAssessmentInProgress ADConfigAssessmentStatus = "IN_PROGRESS"
	ADConfigAssessmentCompleted  ADConfigAssessmentStatus = "COMPLETED"
	ADConfigAssessmentFailed     ADConfigAssessmentStatus = "FAILED"
	ADConfigAssessmentCancelled  ADConfigAssessmentStatus = "CANCELLED"
	ADConfigAssessmentNA         ADConfigAssessmentStatus = "NA"
)

// ScopeLevel is the scope level for an AD configuration.
type ScopeLevel string

const (
	ScopeLevelAccount ScopeLevel = "ACCOUNT"
	ScopeLevelSite    ScopeLevel = "SITE"
	ScopeLevelGroup   ScopeLevel = "GROUP"
)

// ADFeatureName is a named feature available for AD configuration.
type ADFeatureName string

const (
	ADFeatureRangerAD        ADFeatureName = "RANGER_AD"
	ADFeatureSingularityID   ADFeatureName = "SINGULARITY_IDENTITY"
	ADFeatureRangerADProtect ADFeatureName = "RANGER_AD_PROTECT"
)

// ADFeatureType is the type of a feature in the feature status info.
type ADFeatureType string

const (
	ADFeatureTypeADSecure        ADFeatureType = "AD_SECURE"
	ADFeatureTypeADAssessment    ADFeatureType = "AD_ASSESSMENT"
	ADFeatureTypeRangerADProtect ADFeatureType = "RANGER_AD_PROTECT"
)

// ADFeatureStatus is the status of a feature in the feature status info.
type ADFeatureStatus string

const (
	ADFeatureStatusNotEnabled            ADFeatureStatus = "NOT_ENABLED"
	ADFeatureStatusCompleted             ADFeatureStatus = "COMPLETED"
	ADFeatureStatusFailed                ADFeatureStatus = "FAILED"
	ADFeatureStatusInProgress            ADFeatureStatus = "IN_PROGRESS"
	ADFeatureStatusSuccess               ADFeatureStatus = "SUCCESS"
	ADFeatureStatusPendingAndInProgress  ADFeatureStatus = "PENDING_AND_IN_PROGRESS"
	ADFeatureStatusMisconfigurationError ADFeatureStatus = "MISCONFIGURATION_ERROR"
)

// OnboardingStatus is the onboarding status of the AD service.
type OnboardingStatus string

const (
	OnboardingStatusComplete   OnboardingStatus = "COMPLETE"
	OnboardingStatusIncomplete OnboardingStatus = "INCOMPLETE"
)

// ConnectorStatus is the status of an AD connector (Cloudlink).
type ConnectorStatus string

const (
	ConnectorStatusActive   ConnectorStatus = "ACTIVE"
	ConnectorStatusInactive ConnectorStatus = "INACTIVE"
)

// ADConnectorConfigStatus is the onboarding state of the AD connector.
type ADConnectorConfigStatus string

const (
	ADConnectorConfigured    ADConnectorConfigStatus = "CONFIGURED"
	ADConnectorConfigPending ADConnectorConfigStatus = "CONFIG_PENDING"
)

// ScopeInfo identifies the scope where an AD configuration was created.
type ScopeInfo struct {
	ScopeID    string     `json:"scopeId"`
	ScopeLevel ScopeLevel `json:"scopeLevel"`
	ScopeName  string     `json:"scopeName"`
	ScopePath  string     `json:"scopePath"`

	Raw json.RawMessage `json:"-"`
}

func (s *ScopeInfo) UnmarshalJSON(b []byte) error {
	type alias ScopeInfo
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// FeatureStatusInfo is the status detail for a single feature on an AD config.
type FeatureStatusInfo struct {
	FeatureType     ADFeatureType   `json:"featureType"`
	Status          ADFeatureStatus `json:"status"`
	StatusMessage   string          `json:"statusMessage"`
	DetailedMessage string          `json:"detailedMessage"`
	StartTime       string          `json:"startTime"`
	EndTime         string          `json:"endTime"`

	Raw json.RawMessage `json:"-"`
}

func (f *FeatureStatusInfo) UnmarshalJSON(b []byte) error {
	type alias FeatureStatusInfo
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// PolicyUsage shows whether a feature's policy is active.
type PolicyUsage struct {
	FeatureName  ADFeatureName `json:"featureName"`
	PolicyActive bool          `json:"policyActive"`

	Raw json.RawMessage `json:"-"`
}

func (p *PolicyUsage) UnmarshalJSON(b []byte) error {
	type alias PolicyUsage
	if err := json.Unmarshal(b, (*alias)(p)); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0:0], b...)
	return nil
}

// ADConfiguration is an AD configuration entry returned by the API.
type ADConfiguration struct {
	ID                           int64                    `json:"id"`
	TenantID                     string                   `json:"tenantId"`
	CloudlinkID                  int                      `json:"cloudlinkId"`
	DomainName                   string                   `json:"domainName"`
	DomainControllerFqdn         string                   `json:"domainControllerFqdn"`
	TrustingDomainControllerFqdn string                   `json:"trustingDomainControllerFqdn"`
	TrustingDomainName           string                   `json:"trustingDomainName"`
	PortNumber                   int                      `json:"portNumber"`
	Enabled                      bool                     `json:"enabled"`
	AssessmentStatus             ADConfigAssessmentStatus `json:"assessmentStatus"`
	SyncStatus                   bool                     `json:"syncStatus"`
	EncryptionMethod             EncryptionMethod         `json:"encryptionMethod"`
	CreatedAt                    *ScopeInfo               `json:"createdAt"`
	ScopeBundle                  []ScopeInfo              `json:"scopeBundle"`
	Username                     string                   `json:"username"`
	LDAPReferral                 bool                     `json:"ldapReferral"`
	IsConnected                  bool                     `json:"isConnected"`
	AssessOtherDomains           bool                     `json:"assessOtherDomains"`
	FeaturesOpted                []string                 `json:"featuresOpted"`
	IsPolicyActive               bool                     `json:"isPolicyActive"`
	FeatureStatusInfo            []FeatureStatusInfo      `json:"featureStatusInfo"`
	UseWinRmOverSSL              bool                     `json:"useWinRmOverSsl"`
	ADSync                       bool                     `json:"adSync"`
	PolicyUsage                  []PolicyUsage            `json:"policyUsage"`

	Raw json.RawMessage `json:"-"`
}

func (a *ADConfiguration) UnmarshalJSON(b []byte) error {
	type alias ADConfiguration
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// ADConfigScope identifies an allowed scope for an AD configuration.
type ADConfigScope struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// AccessOverTrustInfo holds trust-domain settings for AD configuration.
type AccessOverTrustInfo struct {
	AccessOverTrustEnabled       bool   `json:"accessOverTrustEnabled"`
	TrustingDomainName           string `json:"trustingDomainName,omitempty"`
	TrustingDomainControllerFqdn string `json:"trustingDomainControllerFqdn,omitempty"`
}

// ADConfigurationInput is the request body for creating an AD configuration.
// Password is a secret field — never log or echo it.
type ADConfigurationInput struct {
	CloudlinkID                *int64               `json:"cloudlinkId,omitempty"`
	FeaturesOpted              []string             `json:"featuresOpted,omitempty"`
	AllowedScopes              []ADConfigScope      `json:"allowedScopes,omitempty"`
	IsAllScopesAllowed         *bool                `json:"isAllScopesAllowed,omitempty"`
	DomainName                 string               `json:"domainName"`
	AssessOtherDomainsInForest *bool                `json:"assessOtherDomainsInForest,omitempty"`
	DomainControllerFqdn       string               `json:"domainControllerFqdn"`
	UserName                   string               `json:"userName"`
	Password                   string               `json:"password"`
	EncryptionMethod           EncryptionMethod     `json:"encryptionMethod"`
	AccessOverTrustInfo        *AccessOverTrustInfo `json:"accessOverTrustInfo,omitempty"`
	EnableThreatDetection      *bool                `json:"enableThreatDetection,omitempty"`
	LDAPReferral               *bool                `json:"ldapReferral,omitempty"`
	UseWinRmOverSSL            *bool                `json:"useWinRmOverSsl,omitempty"`
	ADSync                     *bool                `json:"adSync,omitempty"`
}

// ADFeature is an available feature returned by the available-features endpoint.
type ADFeature struct {
	FeatureName ADFeatureName `json:"featureName"`
	Available   bool          `json:"available"`

	Raw json.RawMessage `json:"-"`
}

func (f *ADFeature) UnmarshalJSON(b []byte) error {
	type alias ADFeature
	if err := json.Unmarshal(b, (*alias)(f)); err != nil {
		return err
	}
	f.Raw = append(f.Raw[:0:0], b...)
	return nil
}

// DomainInfo describes an AD domain.
type DomainInfo struct {
	Domain       string `json:"domain"`
	ParentDomain string `json:"parentDomain"`
	Root         bool   `json:"root"`

	Raw json.RawMessage `json:"-"`
}

func (d *DomainInfo) UnmarshalJSON(b []byte) error {
	type alias DomainInfo
	if err := json.Unmarshal(b, (*alias)(d)); err != nil {
		return err
	}
	d.Raw = append(d.Raw[:0:0], b...)
	return nil
}

// TimeZoneInfo is a timezone pair returned by the timezones endpoint.
type TimeZoneInfo struct {
	TimeZoneID  string `json:"timeZoneId"`
	DisplayName string `json:"displayName"`

	Raw json.RawMessage `json:"-"`
}

func (t *TimeZoneInfo) UnmarshalJSON(b []byte) error {
	type alias TimeZoneInfo
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

// ADOnboardingStatus is the onboarding status response.
type ADOnboardingStatus struct {
	Status          OnboardingStatus        `json:"status"`
	FeatureSelected []string                `json:"featureSelected"`
	ADConnector     ADConnectorConfigStatus `json:"adConnector"`
	DomainName      string                  `json:"domainName"`

	Raw json.RawMessage `json:"-"`
}

func (o *ADOnboardingStatus) UnmarshalJSON(b []byte) error {
	type alias ADOnboardingStatus
	if err := json.Unmarshal(b, (*alias)(o)); err != nil {
		return err
	}
	o.Raw = append(o.Raw[:0:0], b...)
	return nil
}

// IdentityParams are common query parameters for Identity AD Service endpoints.
type IdentityParams struct {
	SiteIDs    string
	AccountIDs string
}

func (p *IdentityParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "siteIds", p.SiteIDs)
	addString(v, "accountIds", p.AccountIDs)
	return v
}

const identityBase = "/identity/adservice/api"

// IdentityADConfigurations returns all AD configurations.
func (c *Client) IdentityADConfigurations(ctx context.Context, params *IdentityParams) ([]ADConfiguration, error) {
	var resp singleResponse[[]ADConfiguration]
	if err := c.get(ctx, identityBase+"/adConfigurations", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IdentityADConfigurationAdd creates a new AD configuration.
// The input contains credentials (userName, password) — handle as secrets.
func (c *Client) IdentityADConfigurationAdd(ctx context.Context, params *IdentityParams, input ADConfigurationInput) error {
	body := struct {
		Input ADConfigurationInput `json:"input"`
	}{Input: input}
	qv := params.values()
	u := identityBase + "/addAdConfiguration"
	if len(qv) > 0 {
		u += "?" + qv.Encode()
	}
	var resp json.RawMessage
	return c.post(ctx, u, body, &resp)
}

// IdentityADConfigurationDelete deletes AD configurations by ID.
func (c *Client) IdentityADConfigurationDelete(ctx context.Context, params *IdentityParams, ids []int64) error {
	body := struct {
		Input []int64 `json:"input"`
	}{Input: ids}
	qv := params.values()
	u := identityBase + "/deleteAdConfiguration"
	if len(qv) > 0 {
		u += "?" + qv.Encode()
	}
	var resp json.RawMessage
	return c.post(ctx, u, body, &resp)
}

// IdentityAvailableFeatures returns the list of available AD features.
func (c *Client) IdentityAvailableFeatures(ctx context.Context, params *IdentityParams) ([]ADFeature, error) {
	var resp singleResponse[[]ADFeature]
	if err := c.get(ctx, identityBase+"/availableFeatures", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IdentityDomains returns AD domain information.
func (c *Client) IdentityDomains(ctx context.Context, params *IdentityParams) ([]DomainInfo, error) {
	var resp singleResponse[[]DomainInfo]
	if err := c.get(ctx, identityBase+"/domains", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IdentityADDomains returns AD domain information via the getAdDomains endpoint.
func (c *Client) IdentityADDomains(ctx context.Context, params *IdentityParams) ([]DomainInfo, error) {
	var resp singleResponse[[]DomainInfo]
	if err := c.get(ctx, identityBase+"/getAdDomains", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IdentityTimezones returns the list of available timezone pairs.
func (c *Client) IdentityTimezones(ctx context.Context, params *IdentityParams) ([]TimeZoneInfo, error) {
	var resp singleResponse[[]TimeZoneInfo]
	if err := c.get(ctx, identityBase+"/timezones", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IdentityConfigFeatures returns the list of AD configuration feature names.
func (c *Client) IdentityConfigFeatures(ctx context.Context, params *IdentityParams) ([]string, error) {
	var resp singleResponse[[]string]
	if err := c.get(ctx, identityBase+"/adConfigurationFeatures", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IdentityOnboardingStatus returns the current AD service onboarding status.
func (c *Client) IdentityOnboardingStatus(ctx context.Context, params *IdentityParams) (*ADOnboardingStatus, error) {
	var resp singleResponse[ADOnboardingStatus]
	if err := c.get(ctx, identityBase+"/getOnboardingStatus", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
