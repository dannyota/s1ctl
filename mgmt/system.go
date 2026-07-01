package mgmt

import (
	"context"
	"encoding/json"
)

type SystemInfo struct {
	Release            string `json:"release"`
	Version            string `json:"version"`
	Build              string `json:"build"`
	Patch              string `json:"patch"`
	LatestAgentVersion string `json:"latestAgentVersion"`

	Raw json.RawMessage `json:"-"`
}

func (s *SystemInfo) UnmarshalJSON(b []byte) error {
	type alias SystemInfo
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

type SystemStatus struct {
	Health string `json:"health"`

	Raw json.RawMessage `json:"-"`
}

func (s *SystemStatus) UnmarshalJSON(b []byte) error {
	type alias SystemStatus
	if err := json.Unmarshal(b, (*alias)(s)); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0:0], b...)
	return nil
}

type systemInfoResponse struct {
	Data SystemInfo `json:"data"`
}

type systemStatusResponse struct {
	Data SystemStatus `json:"data"`
}

// SystemInfo returns console version, build, and patch level.
func (c *Client) SystemInfo(ctx context.Context) (*SystemInfo, error) {
	var resp systemInfoResponse
	if err := c.get(ctx, "/system/info", nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SystemStatus returns the console health status.
func (c *Client) SystemStatus(ctx context.Context) (*SystemStatus, error) {
	var resp systemStatusResponse
	if err := c.get(ctx, "/system/status", nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}
