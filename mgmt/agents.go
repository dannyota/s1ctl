package mgmt

import (
	"context"
	"encoding/json"
	"net/url"
)

// Agent is a SentinelOne endpoint agent.
type Agent struct {
	ID                      string `json:"id"`
	ComputerName            string `json:"computerName"`
	Domain                  string `json:"domain"`
	OSType                  string `json:"osType"`
	OSName                  string `json:"osName"`
	OSArch                  string `json:"osArch"`
	AgentVersion            string `json:"agentVersion"`
	IsActive                bool   `json:"isActive"`
	Infected                bool   `json:"infected"`
	ActiveThreats           int    `json:"activeThreats"`
	NetworkStatus           string `json:"networkStatus"`
	MachineType             string `json:"machineType"`
	AccountID               string `json:"accountId"`
	AccountName             string `json:"accountName"`
	SiteID                  string `json:"siteId"`
	SiteName                string `json:"siteName"`
	GroupID                 string `json:"groupId"`
	GroupName               string `json:"groupName"`
	ExternalIP              string `json:"externalIp"`
	LastActiveDate          string `json:"lastActiveDate"`
	RegisteredAt            string `json:"registeredAt"`
	CreatedAt               string `json:"createdAt"`
	UpdatedAt               string `json:"updatedAt"`
	LastLoggedInUserName    string `json:"lastLoggedInUserName"`
	IsDecommissioned        bool   `json:"isDecommissioned"`
	IsUninstalled           bool   `json:"isUninstalled"`
	IsUpToDate              bool   `json:"isUpToDate"`
	ScanStatus              string `json:"scanStatus"`
	MitigationMode          string `json:"mitigationMode"`
	FirewallEnabled         bool   `json:"firewallEnabled"`
	TotalMemory             int    `json:"totalMemory"`
	CPUCount                int    `json:"cpuCount"`
	CoreCount               int    `json:"coreCount"`
	UUID                    string `json:"uuid"`
	SerialNumber            string `json:"serialNumber"`
	ModelName               string `json:"modelName"`
	AppsVulnerabilityStatus string `json:"appsVulnerabilityStatus"`
	OperationalState        string `json:"operationalState"`

	Raw json.RawMessage `json:"-"`
}

func (a *Agent) UnmarshalJSON(b []byte) error {
	type alias Agent
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// AgentListParams are query parameters for listing agents.
type AgentListParams struct {
	SiteIDs          []string
	GroupIDs         []string
	AccountIDs       []string
	OSTypes          []string
	IsActive         *bool
	Infected         *bool
	IsDecommissioned *bool
	IsUninstalled    *bool
	IsUpToDate       *bool
	NetworkStatuses  []string
	MachineTypes     []string
	Query            string
	Limit            int
	Cursor           string
	SortBy           string
	SortOrder        string
	CountOnly        bool
}

func (p *AgentListParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "osTypes", p.OSTypes)
	addCSV(v, "networkStatuses", p.NetworkStatuses)
	addCSV(v, "machineTypes", p.MachineTypes)
	addBool(v, "isActive", p.IsActive)
	addBool(v, "infected", p.Infected)
	addBool(v, "isDecommissioned", p.IsDecommissioned)
	addBool(v, "isUninstalled", p.IsUninstalled)
	addBool(v, "isUpToDate", p.IsUpToDate)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	addString(v, "sortBy", p.SortBy)
	addString(v, "sortOrder", p.SortOrder)
	if p.CountOnly {
		v.Set("countOnly", "true")
	}
	return v
}

// AgentsList returns a paginated list of agents.
func (c *Client) AgentsList(ctx context.Context, params *AgentListParams) ([]Agent, *Pagination, error) {
	return list[Agent](c, ctx, "/agents", params.values())
}

// AgentsCount returns the count of agents matching the filter.
func (c *Client) AgentsCount(ctx context.Context, params *AgentListParams) (int, error) {
	if params == nil {
		params = &AgentListParams{}
	}
	params.CountOnly = true
	_, pag, err := list[Agent](c, ctx, "/agents", params.values())
	if err != nil {
		return 0, err
	}
	return pag.TotalItems, nil
}

// AgentsGet returns a single agent by ID.
func (c *Client) AgentsGet(ctx context.Context, id string) (*Agent, error) {
	return getByID[Agent](c, ctx, "/agents", "agent", id)
}
