package sdl

import (
	"context"
	"encoding/json"
)

// GetFileRequest is the request body for retrieving a configuration file.
type GetFileRequest struct {
	Path            string `json:"path"`
	ExpectedVersion int64  `json:"expectedVersion,omitempty"`
	PrettyPrint     bool   `json:"prettyprint,omitempty"`
}

// GetFileResponse is the response from retrieving a configuration file.
// Status is "success", "success/unchanged", or "success/noSuchFile".
type GetFileResponse struct {
	Status        string `json:"status"`
	Path          string `json:"path"`
	Version       int64  `json:"version"`
	CreateDate    int64  `json:"createDate"`
	ModDate       int64  `json:"modDate"`
	Content       string `json:"content"`
	StalenessSlop int64  `json:"stalenessSlop"`
	Message       string `json:"message"`

	Raw json.RawMessage `json:"-"`
}

func (r *GetFileResponse) UnmarshalJSON(b []byte) error {
	type alias GetFileResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// PutFileRequest is the request body for creating, updating, or deleting a
// configuration file. Set DeleteFile to true (and omit Content) to delete.
type PutFileRequest struct {
	Path            string `json:"path"`
	Content         string `json:"content"`
	DeleteFile      bool   `json:"deleteFile,omitempty"`
	ExpectedVersion int64  `json:"expectedVersion,omitempty"`
	PrettyPrint     bool   `json:"prettyprint,omitempty"`
}

// PutFileResponse is the response from creating, updating, or deleting a
// configuration file. Status is "success" or "error/client/versionMismatch".
type PutFileResponse struct {
	Status string `json:"status"`

	Raw json.RawMessage `json:"-"`
}

func (r *PutFileResponse) UnmarshalJSON(b []byte) error {
	type alias PutFileResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// ListFilesResponse is the response from listing all configuration files.
type ListFilesResponse struct {
	Status string   `json:"status"`
	Paths  []string `json:"paths"`

	Raw json.RawMessage `json:"-"`
}

func (r *ListFilesResponse) UnmarshalJSON(b []byte) error {
	type alias ListFilesResponse
	if err := json.Unmarshal(b, (*alias)(r)); err != nil {
		return err
	}
	r.Raw = append(r.Raw[:0:0], b...)
	return nil
}

// GetFile retrieves a configuration file by path.
func (c *Client) GetFile(ctx context.Context, req *GetFileRequest) (*GetFileResponse, error) {
	var resp GetFileResponse
	if err := c.post(ctx, "/api/getFile", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PutFile creates, updates, or deletes a configuration file.
func (c *Client) PutFile(ctx context.Context, req *PutFileRequest) (*PutFileResponse, error) {
	var resp PutFileResponse
	if err := c.post(ctx, "/api/putFile", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListFiles lists all configuration file paths.
func (c *Client) ListFiles(ctx context.Context) (*ListFilesResponse, error) {
	var resp ListFilesResponse
	if err := c.post(ctx, "/api/listFiles", struct{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
