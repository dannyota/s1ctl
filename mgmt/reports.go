package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// Report is a generated SentinelOne report.
type Report struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Scope           string          `json:"scope"`
	Frequency       string          `json:"frequency"`
	Interval        string          `json:"interval"`
	ScheduleType    string          `json:"scheduleType"`
	CreatorID       string          `json:"creatorId"`
	CreatorName     string          `json:"creatorName"`
	CreatedAt       string          `json:"createdAt"`
	FromDate        string          `json:"fromDate"`
	ToDate          string          `json:"toDate"`
	InsightTypes    json.RawMessage `json:"insightTypes"`
	AttachmentTypes []string        `json:"attachmentTypes"`
	Status          string          `json:"status"`
	Sites           string          `json:"sites"`

	Raw json.RawMessage `json:"-"`
}

func (r *Report) UnmarshalJSON(b []byte) error {
	type alias Report
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// ReportListParams are query parameters for listing reports.
type ReportListParams struct {
	SiteIDs      []string
	AccountIDs   []string
	IDs          []string
	Name         string
	Scope        string
	Frequency    string
	ScheduleType string
	Query        string
	TaskID       string
	Limit        int
	Cursor       string
	SortBy       string
	SortOrder    string
}

func (p *ReportListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "ids", p.IDs)
	addString(v, "name", p.Name)
	addString(v, "scope", p.Scope)
	addString(v, "frequency", p.Frequency)
	addString(v, "scheduleType", p.ScheduleType)
	addString(v, "query", p.Query)
	addString(v, "taskId", p.TaskID)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// ReportsList returns a paginated list of generated reports.
func (c *Client) ReportsList(ctx context.Context, params *ReportListParams) ([]Report, *Pagination, error) {
	return list[Report](c, ctx, "/reports", params.values())
}

// ReportTask is a SentinelOne report task or schedule.
type ReportTask struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Scope           string          `json:"scope"`
	Frequency       string          `json:"frequency"`
	Day             string          `json:"day"`
	ScheduleType    string          `json:"scheduleType"`
	CreatorID       string          `json:"creatorId"`
	CreatorName     string          `json:"creatorName"`
	InsightTypes    json.RawMessage `json:"insightTypes"`
	AttachmentTypes []string        `json:"attachmentTypes"`
	Sites           string          `json:"sites"`
	FromDate        string          `json:"fromDate"`
	ToDate          string          `json:"toDate"`
	Recipients      []string        `json:"recipients"`
	IsTrend         bool            `json:"isTrend"`

	Raw json.RawMessage `json:"-"`
}

func (t *ReportTask) UnmarshalJSON(b []byte) error {
	type alias ReportTask
	if err := json.Unmarshal(b, (*alias)(t)); err != nil {
		return err
	}
	t.Raw = append(t.Raw[:0:0], b...)
	return nil
}

// ReportTaskListParams are query parameters for listing report tasks.
type ReportTaskListParams struct {
	SiteIDs      []string
	AccountIDs   []string
	IDs          []string
	Name         string
	Scope        string
	Frequency    string
	ScheduleType string
	Query        string
	Limit        int
	Cursor       string
	SortBy       string
	SortOrder    string
}

func (p *ReportTaskListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "ids", p.IDs)
	addString(v, "name", p.Name)
	addString(v, "scope", p.Scope)
	addString(v, "frequency", p.Frequency)
	addString(v, "scheduleType", p.ScheduleType)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	return v
}

// ReportTasksList returns a paginated list of report tasks.
func (c *Client) ReportTasksList(ctx context.Context, params *ReportTaskListParams) ([]ReportTask, *Pagination, error) {
	return list[ReportTask](c, ctx, "/report-tasks", params.values())
}

// ReportTaskCreate is the input for creating a report task.
type ReportTaskCreate struct {
	Name            string          `json:"name"`
	ScheduleType    string          `json:"scheduleType"`
	InsightTypes    json.RawMessage `json:"insightTypes"`
	Frequency       string          `json:"frequency,omitempty"`
	Day             string          `json:"day,omitempty"`
	FromDate        string          `json:"fromDate,omitempty"`
	ToDate          string          `json:"toDate,omitempty"`
	AttachmentTypes []string        `json:"attachmentTypes,omitempty"`
	Recipients      []string        `json:"recipients,omitempty"`
	IsTrend         *bool           `json:"isTrend,omitempty"`
}

type reportTaskCreateRequest struct {
	Filter struct {
		SiteIDs    []string `json:"siteIds,omitempty"`
		AccountIDs []string `json:"accountIds,omitempty"`
		Scope      string   `json:"scope,omitempty"`
	} `json:"filter"`
	Data ReportTaskCreate `json:"data"`
}

// ReportTasksCreate creates a new report task.
func (c *Client) ReportTasksCreate(ctx context.Context, siteIDs, accountIDs []string, scope string, task ReportTaskCreate) error {
	req := reportTaskCreateRequest{Data: task}
	req.Filter.SiteIDs = siteIDs
	req.Filter.AccountIDs = accountIDs
	req.Filter.Scope = scope
	var resp json.RawMessage
	return c.post(ctx, "/report-tasks", req, &resp)
}

// insightTypesResponse is the response envelope for insight types.
type insightTypesResponse struct {
	Data struct {
		InsightTypes json.RawMessage `json:"insightTypes"`
	} `json:"data"`
}

// InsightTypesParams are query parameters for listing insight types.
type InsightTypesParams struct {
	SiteIDs    []string
	AccountIDs []string
	GroupIDs   []string
}

func (p *InsightTypesParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	return v
}

// ReportsInsightTypes returns available report insight types.
func (c *Client) ReportsInsightTypes(ctx context.Context, params *InsightTypesParams) (json.RawMessage, error) {
	var resp insightTypesResponse
	if err := c.get(ctx, "/reports/insights/types", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data.InsightTypes, nil
}

// ReportDownload downloads a report in the specified format (pdf or html).
func (c *Client) ReportDownload(ctx context.Context, reportID, format string) ([]byte, error) {
	path := fmt.Sprintf("/reports/%s/%s", reportID, format)
	return c.getRaw(ctx, path, nil)
}
