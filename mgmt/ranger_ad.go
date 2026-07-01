package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// AssessmentStatus is the overall status of a Ranger AD assessment.
type AssessmentStatus string

const (
	AssessmentStatusPending    AssessmentStatus = "PENDING"
	AssessmentStatusInProgress AssessmentStatus = "IN_PROGRESS"
	AssessmentStatusCompleted  AssessmentStatus = "COMPLETED"
)

// ExposureDetectionStatus is the detection status of an AD exposure.
type ExposureDetectionStatus string

const (
	ExposureStatusVulnerable    ExposureDetectionStatus = "Vulnerable"
	ExposureStatusNotVulnerable ExposureDetectionStatus = "Not_Vulnerable"
	ExposureStatusSkipped       ExposureDetectionStatus = "Skipped"
	ExposureStatusInProgress    ExposureDetectionStatus = "In_Progress"
	ExposureStatusPending       ExposureDetectionStatus = "Pending"
	ExposureStatusMitigated     ExposureDetectionStatus = "Mitigated"
)

// ExposureSeverity is the severity level of an AD exposure.
type ExposureSeverity string

const (
	ExposureSeverityCritical ExposureSeverity = "Critical"
	ExposureSeverityHigh     ExposureSeverity = "High"
	ExposureSeverityMedium   ExposureSeverity = "Medium"
	ExposureSeverityLow      ExposureSeverity = "Low"
)

// ExposureSource is the source of an AD exposure detection.
type ExposureSource string

const (
	ExposureSourceOnPremAD ExposureSource = "OnPremAD"
	ExposureSourceAzureAD  ExposureSource = "AzureAD"
)

// AssessmentScanSource is the source type for a triggered assessment.
type AssessmentScanSource string

const (
	AssessmentScanSourceAD    AssessmentScanSource = "AD"
	AssessmentScanSourceAzure AssessmentScanSource = "Azure"
)

// DomainStatus is the assessment status for a single AD domain.
type DomainStatus struct {
	DomainName      string `json:"domainName"`
	ForestName      string `json:"forestName"`
	TotalJobs       int    `json:"totalJobs"`
	CompletedJobs   int    `json:"completedJobs"`
	DomainCompleted bool   `json:"domainCompletedStatus"`
}

// TenantStatus is the assessment status for an Azure tenant.
type TenantStatus struct {
	TenantID        string `json:"tenantId"`
	TotalJobs       int    `json:"totalJobs"`
	CompletedJobs   int    `json:"completedJobs"`
	TenantCompleted bool   `json:"tenantCompletedStatus"`
}

// ADAssessmentStatus is the response from the Ranger AD assessment status endpoint.
type ADAssessmentStatus struct {
	Status  AssessmentStatus `json:"status"`
	Domains []DomainStatus   `json:"domainWiseCurrentStatusList"`
	Tenants []TenantStatus   `json:"tenantWiseCurrentStatusList"`

	Raw json.RawMessage `json:"-"`
}

func (s *ADAssessmentStatus) UnmarshalJSON(b []byte) error {
	type alias ADAssessmentStatus
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

// ADExposure is a Ranger AD exposure finding.
type ADExposure struct {
	ID                  string                  `json:"id"`
	DetectionID         int                     `json:"detectionId"`
	DetectionName       string                  `json:"detectionName"`
	DetectionStatus     ExposureDetectionStatus `json:"detectionStatus"`
	Severity            ExposureSeverity        `json:"severity"`
	Source              ExposureSource          `json:"source"`
	DomainName          string                  `json:"domainName"`
	ForestName          string                  `json:"forestName"`
	VulnerableCount     int                     `json:"vulnerableCount"`
	PrevVulnerableCount *int                    `json:"prevVulnerableCount"`
	Acknowledged        bool                    `json:"acknowledged"`
	Remediable          bool                    `json:"remediable"`
	RunTimestamp        int                     `json:"runTimestamp"`
	SpecialRun          *bool                   `json:"specialRun"`
	SpecialRunTimestamp *int                    `json:"specialRunTimestamp"`
	SkipReason          *string                 `json:"skipReason"`
	SkipCode            *string                 `json:"skipCode"`
	HasExcludableObjs   *bool                   `json:"hasExcludableObjects"`

	Raw json.RawMessage `json:"-"`
}

func (e *ADExposure) UnmarshalJSON(b []byte) error {
	type alias ADExposure
	if err := json.Unmarshal(b, (*alias)(e)); err != nil {
		return err
	}
	e.Raw = append(e.Raw[:0:0], b...)
	return nil
}

// ADAffectedObject is an object affected by a Ranger AD exposure.
type ADAffectedObject struct {
	ID             int     `json:"id"`
	RunID          int     `json:"runId"`
	DN             *string `json:"dn"`
	DisplayName    *string `json:"displayName"`
	SAMAccountName *string `json:"samAccountName"`
	UPN            *string `json:"upn"`
	CommonName     *string `json:"commonName"`
	ObjectType     *string `json:"objectType"`
	ObjectGUID     *string `json:"objectGUID"`
	DNSHostName    *string `json:"dNSHostName"`
	OS             *string `json:"os"`
	AccountStatus  *string `json:"accountStatus"`
	Description    *string `json:"description"`
	LastLogon      *int    `json:"lastLogonTimestamp"`
	WhenCreated    *int    `json:"whenCreated"`
	WhenChanged    *int    `json:"whenChanged"`
	PwdLastSet     *int    `json:"pwdLastSet"`

	Raw json.RawMessage `json:"-"`
}

func (o *ADAffectedObject) UnmarshalJSON(b []byte) error {
	type alias ADAffectedObject
	if err := json.Unmarshal(b, (*alias)(o)); err != nil {
		return err
	}
	o.Raw = append(o.Raw[:0:0], b...)
	return nil
}

// ADAssessmentStatusParams are query parameters for the assessment status endpoint.
type ADAssessmentStatusParams struct {
	SiteIDs    string
	AccountIDs string
}

func (p *ADAssessmentStatusParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "siteIds", p.SiteIDs)
	addString(v, "accountIds", p.AccountIDs)
	return v
}

// RangerADAssessmentStatus returns the current Ranger AD assessment status.
func (c *Client) RangerADAssessmentStatus(ctx context.Context, params *ADAssessmentStatusParams) (*ADAssessmentStatus, error) {
	var resp singleResponse[ADAssessmentStatus]
	if err := c.get(ctx, "/ranger-ad/assessment-status", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ADExposureFilter is the filter body for listing AD exposures.
type ADExposureFilter struct {
	DetectionStatus []string `json:"detectionStatus,omitempty"`
	DetectionName   []string `json:"detectionName,omitempty"`
	DomainName      []string `json:"domainName,omitempty"`
	ForestName      []string `json:"forestName,omitempty"`
	Severity        []string `json:"severity,omitempty"`
	Source          []string `json:"source,omitempty"`
}

// ADExposureListParams are parameters for listing AD exposures.
type ADExposureListParams struct {
	Limit      int
	Skip       int
	SiteIDs    string
	AccountIDs string
	Filter     ADExposureFilter
}

func (p *ADExposureListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addInt(v, "limit", p.Limit)
	addInt(v, "skip", p.Skip)
	addString(v, "siteIds", p.SiteIDs)
	addString(v, "accountIds", p.AccountIDs)
	return v
}

// RangerADExposures returns a paginated list of AD exposures.
func (c *Client) RangerADExposures(ctx context.Context, params *ADExposureListParams) ([]ADExposure, *Pagination, error) {
	body := struct {
		Filter ADExposureFilter `json:"filter"`
	}{Filter: params.Filter}
	var resp listResponse[ADExposure]
	// This endpoint uses POST with a body + query params for pagination.
	qv := params.values()
	u := "/ranger-ad/get-exposures"
	if len(qv) > 0 {
		u += "?" + qv.Encode()
	}
	if err := c.post(ctx, u, body, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// ADAffectedObjectFilter is the filter body for listing affected objects.
type ADAffectedObjectFilter struct {
	DetectionName []string `json:"detectionName"`
	DomainName    []string `json:"domainName"`
	ForestName    []string `json:"forestName,omitempty"`
	ObjectType    []string `json:"objectType,omitempty"`
}

// ADAffectedObjectListParams are parameters for listing affected objects.
type ADAffectedObjectListParams struct {
	Limit      int
	Skip       int
	SiteIDs    string
	AccountIDs string
	Filter     ADAffectedObjectFilter
}

func (p *ADAffectedObjectListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addInt(v, "limit", p.Limit)
	addInt(v, "skip", p.Skip)
	addString(v, "siteIds", p.SiteIDs)
	addString(v, "accountIds", p.AccountIDs)
	return v
}

// RangerADAffectedObjects returns a paginated list of affected objects for a detection.
func (c *Client) RangerADAffectedObjects(ctx context.Context, params *ADAffectedObjectListParams) ([]ADAffectedObject, *Pagination, error) {
	body := struct {
		Filter ADAffectedObjectFilter `json:"filter"`
	}{Filter: params.Filter}
	var resp listResponse[ADAffectedObject]
	qv := params.values()
	u := "/ranger-ad/get-affected-objects"
	if len(qv) > 0 {
		u += "?" + qv.Encode()
	}
	if err := c.post(ctx, u, body, &resp); err != nil {
		return nil, nil, err
	}
	return resp.Data, &resp.Pagination, nil
}

// ADTriggerAssessmentFilter is the filter for triggering an AD assessment.
type ADTriggerAssessmentFilter struct {
	IsFullScan   bool                        `json:"isFullScan"`
	DomainName   []string                    `json:"domainName,omitempty"`
	ScanSource   *string                     `json:"scanSource,omitempty"`
	ExposureList []ADTriggerExposureListItem `json:"exposureList,omitempty"`
}

// ADTriggerExposureListItem identifies a specific exposure to reassess.
type ADTriggerExposureListItem struct {
	DomainName    string `json:"domainName"`
	DetectionName string `json:"detectionName"`
}

// ADTriggerAssessmentParams are parameters for triggering an AD assessment.
type ADTriggerAssessmentParams struct {
	SiteIDs    string
	AccountIDs string
	Filter     ADTriggerAssessmentFilter
}

func (p *ADTriggerAssessmentParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "siteIds", p.SiteIDs)
	addString(v, "accountIds", p.AccountIDs)
	return v
}

// successResponse is the envelope for success/message responses.
type successResponse struct {
	Data struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	} `json:"data"`
}

// RangerADTriggerAssessment triggers a new AD assessment scan.
func (c *Client) RangerADTriggerAssessment(ctx context.Context, params *ADTriggerAssessmentParams) (bool, string, error) {
	body := struct {
		Filter ADTriggerAssessmentFilter `json:"filter"`
	}{Filter: params.Filter}
	var resp successResponse
	qv := params.values()
	u := "/ranger-ad/trigger-assessment"
	if len(qv) > 0 {
		u += "?" + qv.Encode()
	}
	if err := c.post(ctx, u, body, &resp); err != nil {
		return false, "", err
	}
	return resp.Data.Success, resp.Data.Message, nil
}
