package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// Connector is a Cloudlink (AD connector) configuration.
type Connector struct {
	CloudlinkID    int64           `json:"cloudlinkId"`
	MgmtID         int             `json:"mgmtId"`
	Status         ConnectorStatus `json:"status"`
	ComputerName   string          `json:"computerName"`
	AgentType      string          `json:"agentType"`
	OSName         string          `json:"osName"`
	Version        string          `json:"version"`
	GUID           string          `json:"guid"`
	IsUnifiedAgent bool            `json:"isUnifiedAgent"`
	IPAddress      string          `json:"ipAddress"`
	DomainName     string          `json:"domainName"`
	LastSeen       string          `json:"lastSeen"`
	ScopePath      string          `json:"scopePath"`

	Raw json.RawMessage `json:"-"`
}

func (c *Connector) UnmarshalJSON(b []byte) error {
	type alias Connector
	if err := json.Unmarshal(b, (*alias)(c)); err != nil {
		return err
	}
	c.Raw = append(c.Raw[:0:0], b...)
	return nil
}

// WindowsAgent is a Windows unified agent returned by the connector endpoints.
type WindowsAgent struct {
	ID           int    `json:"id"`
	MgmtID       int    `json:"mgmtId"`
	UUID         string `json:"uuid"`
	OSName       string `json:"osName"`
	IPAddress    string `json:"ipAddress"`
	AgentVersion string `json:"agentVersion"`
	AgentType    string `json:"agentType"`
	DomainName   string `json:"domainName"`
	Status       string `json:"status"`
	HostName     string `json:"hostName"`
	ScopePath    string `json:"scopePath"`

	Raw json.RawMessage `json:"-"`
}

func (w *WindowsAgent) UnmarshalJSON(b []byte) error {
	type alias WindowsAgent
	if err := json.Unmarshal(b, (*alias)(w)); err != nil {
		return err
	}
	w.Raw = append(w.Raw[:0:0], b...)
	return nil
}

// IdentityConnectors returns all Cloudlink connector configurations.
func (c *Client) IdentityConnectors(ctx context.Context, params *IdentityParams) ([]Connector, error) {
	var resp singleResponse[[]Connector]
	if err := c.get(ctx, identityBase+"/getCloudlinkConfigurations", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IdentityConnector returns the single Cloudlink connector configuration.
func (c *Client) IdentityConnector(ctx context.Context, params *IdentityParams) (*Connector, error) {
	var resp singleResponse[Connector]
	if err := c.get(ctx, identityBase+"/getCloudlinkConfiguration", params.values(), &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// IdentityConnectorReplace replaces the AD connector with a new agent.
func (c *Client) IdentityConnectorReplace(ctx context.Context, params *IdentityParams, agentUUID string) error {
	qv := params.values()
	qv.Set("agentUuid", agentUUID)
	u := identityBase + "/replaceAdConnector"
	if len(qv) > 0 {
		u += "?" + qv.Encode()
	}
	var resp json.RawMessage
	return c.post(ctx, u, nil, &resp)
}

// WindowsAgentParams are query parameters for listing Windows unified agents.
type WindowsAgentParams struct {
	SiteIDs     string
	AccountIDs  string
	FilterInput string
	RequestID   string
}

func (p *WindowsAgentParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addString(v, "siteIds", p.SiteIDs)
	addString(v, "accountIds", p.AccountIDs)
	addString(v, "filterInput", p.FilterInput)
	addString(v, "requestId", p.RequestID)
	return v
}

// IdentityWindowsAgents returns Windows unified agents matching the filter.
func (c *Client) IdentityWindowsAgents(ctx context.Context, params *WindowsAgentParams) ([]WindowsAgent, error) {
	var resp singleResponse[[]WindowsAgent]
	if err := c.get(ctx, identityBase+"/getWindowsUnifiedAgents", params.values(), &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
