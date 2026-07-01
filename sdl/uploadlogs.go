package sdl

import (
	"context"
	"encoding/json"
	"strings"
)

// UploadLogsRequest is the request for plain-text log ingestion.
type UploadLogsRequest struct {
	Parser     string
	ServerHost string
	Logfile    string
	Nonce      string
	Body       string
}

// UploadLogsResponse is the response from plain-text log ingestion.
type UploadLogsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	Raw json.RawMessage `json:"-"`
}

func (r *UploadLogsResponse) UnmarshalJSON(b []byte) error {
	type alias UploadLogsResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// UploadLogs uploads unstructured plain-text log data.
func (c *Client) UploadLogs(ctx context.Context, req *UploadLogsRequest) (*UploadLogsResponse, error) {
	headers := make(map[string]string)
	if req.Parser != "" {
		headers["parser"] = req.Parser
	}
	if req.ServerHost != "" {
		headers["server-host"] = req.ServerHost
	}
	if req.Logfile != "" {
		headers["logfile"] = req.Logfile
	}
	if req.Nonce != "" {
		headers["Nonce"] = req.Nonce
	}

	var resp UploadLogsResponse
	if err := c.postText(ctx, "/api/uploadLogs", "text/plain", headers, strings.NewReader(req.Body), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
