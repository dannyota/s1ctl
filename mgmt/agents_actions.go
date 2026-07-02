package mgmt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// AgentsDisconnect network-disconnects (isolates) agents.
func (c *Client) AgentsDisconnect(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/disconnect", filter, nil)
}

// AgentsConnect reconnects previously isolated agents.
func (c *Client) AgentsConnect(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/connect", filter, nil)
}

// AgentsInitiateScan starts a full disk scan on agents.
func (c *Client) AgentsInitiateScan(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/initiate-scan", filter, nil)
}

// AgentsAbortScan aborts a running scan on agents.
func (c *Client) AgentsAbortScan(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/abort-scan", filter, nil)
}

// AgentsDecommission decommissions agents.
func (c *Client) AgentsDecommission(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/decommission", filter, nil)
}

// AgentsShutdown shuts down agents.
func (c *Client) AgentsShutdown(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/shutdown", filter, nil)
}

// AgentsUninstall uninstalls agents.
func (c *Client) AgentsUninstall(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/uninstall", filter, nil)
}

// UpdateSoftwareData specifies which package to use for an agent software update.
// Exactly one of PackageID, FileName, or Path must be set.
type UpdateSoftwareData struct {
	PackageID       string `json:"packageId,omitempty"`
	FileName        string `json:"fileName,omitempty"`
	Path            string `json:"path,omitempty"`
	OSType          string `json:"osType,omitempty"`
	PackageType     string `json:"packageType,omitempty"`
	IsScheduled     *bool  `json:"isScheduled,omitempty"`
	AllowDowngrade  *bool  `json:"allowDowngrade,omitempty"`
	IgnoreConflicts *bool  `json:"ignoreConflicts,omitempty"`
}

// AgentsUpdateSoftware triggers a software update on agents.
func (c *Client) AgentsUpdateSoftware(ctx context.Context, filter ActionFilter, data UpdateSoftwareData) (int, error) {
	return doAction(c, ctx, "/agents/actions/update-software", filter, data)
}

// AgentsMoveToSite moves agents to a different site.
func (c *Client) AgentsMoveToSite(ctx context.Context, siteID string, filter ActionFilter) (int, error) {
	data := map[string]string{"targetSiteId": siteID}
	return doAction(c, ctx, "/agents/actions/move-to-site", filter, data)
}

// AgentsFetchLogs fetches diagnostic logs from agents.
func (c *Client) AgentsFetchLogs(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/fetch-logs", filter, nil)
}

// AgentsRestartMachine restarts the machines running agents.
func (c *Client) AgentsRestartMachine(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/restart-machine", filter, nil)
}

// AgentsEnableAgent enables agents.
func (c *Client) AgentsEnableAgent(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/enable-agent", filter, nil)
}

// AgentsDisableAgent disables agents.
func (c *Client) AgentsDisableAgent(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/disable-agent", filter, nil)
}

// AgentsResetLocalConfig resets local configuration on agents.
func (c *Client) AgentsResetLocalConfig(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/reset-local-config", filter, nil)
}

// AgentsApproveUninstall approves a pending uninstall request on agents.
func (c *Client) AgentsApproveUninstall(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/approve-uninstall", filter, nil)
}

// AgentsRejectUninstall rejects a pending uninstall request on agents.
func (c *Client) AgentsRejectUninstall(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/reject-uninstall", filter, nil)
}

// AgentsMarkUpToDate marks agents as up to date.
func (c *Client) AgentsMarkUpToDate(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/mark-up-to-date", filter, nil)
}

// AgentsSetExternalID sets the external ID on agents.
func (c *Client) AgentsSetExternalID(ctx context.Context, externalID string, filter ActionFilter) (int, error) {
	data := map[string]string{"externalId": externalID}
	return doAction(c, ctx, "/agents/actions/set-external-id", filter, data)
}

// AgentsRandomizeUUID randomizes the UUID on agents.
func (c *Client) AgentsRandomizeUUID(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/randomize-uuid", filter, nil)
}

// AgentsFirewallLogging enables or disables firewall logging on agents.
func (c *Client) AgentsFirewallLogging(ctx context.Context, enable bool, filter ActionFilter) (int, error) {
	data := map[string]bool{"reportLog": enable}
	return doAction(c, ctx, "/agents/actions/firewall-logging", filter, data)
}

// AgentsMoveToGroup moves agents to a different group within the same site.
// The group ID is the target group; the filter selects which agents to move.
func (c *Client) AgentsMoveToGroup(ctx context.Context, groupID string, filter ActionFilter) (int, error) {
	if filter.isEmpty() {
		return 0, fmt.Errorf("mgmt: action requires at least one filter (ids, siteIds, groupIds, or query)")
	}
	req := actionRequest{Filter: filter}
	var resp struct {
		Data struct {
			AgentsMoved int `json:"agentsMoved"`
		} `json:"data"`
	}
	if err := c.put(ctx, fmt.Sprintf("/groups/%s/move-agents", groupID), req, &resp); err != nil {
		return 0, err
	}
	return resp.Data.AgentsMoved, nil
}

// AgentsBroadcast displays a broadcast message on the endpoints of matching
// agents.
func (c *Client) AgentsBroadcast(ctx context.Context, message string, filter ActionFilter) (int, error) {
	data := map[string]string{"message": message}
	return doAction(c, ctx, "/agents/actions/broadcast", filter, data)
}

// AgentsResetPassphrase resets the maintenance passphrase on agents. The API
// reports per-agent results; the returned count is the number of agents for
// which a reset was attempted.
func (c *Client) AgentsResetPassphrase(ctx context.Context, filter ActionFilter) (int, error) {
	if filter.isEmpty() {
		return 0, fmt.Errorf("mgmt: action requires at least one filter (ids, siteIds, groupIds, or query)")
	}
	req := actionRequest{Filter: filter}
	var resp struct {
		Data struct {
			Results []struct {
				AgentID   string `json:"agentId"`
				Attempted bool   `json:"attempted"`
				Status    string `json:"status"`
			} `json:"results"`
		} `json:"data"`
	}
	if err := c.post(ctx, "/agents/actions/reset-passphrase", req, &resp); err != nil {
		return 0, err
	}
	n := 0
	for _, r := range resp.Data.Results {
		if r.Attempted {
			n++
		}
	}
	return n, nil
}

// AgentsRanger enables or disables Ranger network discovery on agents.
func (c *Client) AgentsRanger(ctx context.Context, enable bool, filter ActionFilter) (int, error) {
	path := "/agents/actions/ranger-disable"
	if enable {
		path = "/agents/actions/ranger-enable"
	}
	return doAction(c, ctx, path, filter, nil)
}

// AgentsFetchInstalledApps requests the installed-applications inventory from
// agents (surfaced under application management once fetched).
func (c *Client) AgentsFetchInstalledApps(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/fetch-installed-apps", filter, nil)
}

// AgentsFetchFirewallRules requests the current firewall-rules inventory from
// agents. The API requires a data object; an empty object requests the default
// (current, native) configuration.
func (c *Client) AgentsFetchFirewallRules(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/fetch-firewall-rules", filter, struct{}{})
}

// AgentsFetchFiles requests specific files from a single agent. The fetched
// files are uploaded to the console encrypted with password; paths lists up to
// 10 absolute file paths. Returns whether the request was accepted.
func (c *Client) AgentsFetchFiles(ctx context.Context, id string, paths []string, password string) (bool, error) {
	body := map[string]any{
		"data": map[string]any{
			"files":    paths,
			"password": password,
		},
	}
	var resp struct {
		Data struct {
			Success bool `json:"success"`
		} `json:"data"`
	}
	if err := c.post(ctx, fmt.Sprintf("/agents/%s/actions/fetch-files", id), body, &resp); err != nil {
		return false, err
	}
	return resp.Data.Success, nil
}

// AgentsLocalUpgradeAuthorization sets the local upgrade/downgrade approval on
// agents. authorization is the approval-expiration timestamp; an empty string
// clears the authorization (sends null).
func (c *Client) AgentsLocalUpgradeAuthorization(ctx context.Context, filter ActionFilter, authorization string) (int, error) {
	var auth *string
	if authorization != "" {
		auth = &authorization
	}
	data := map[string]*string{"agentAuthorization": auth}
	return doAction(c, ctx, "/agents/actions/local-upgrade-authorization", filter, data)
}

// AgentLocalUpgradeAuth is a single agent's local upgrade/downgrade
// authorization state.
type AgentLocalUpgradeAuth struct {
	AgentAuthorization string `json:"agentAuthorization"`
	SiteAuthorization  string `json:"siteAuthorization"`

	Raw json.RawMessage `json:"-"`
}

func (a *AgentLocalUpgradeAuth) UnmarshalJSON(b []byte) error {
	type alias AgentLocalUpgradeAuth
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// AgentsLocalUpgradeAuthGet returns the local upgrade/downgrade authorization
// for a single agent.
func (c *Client) AgentsLocalUpgradeAuthGet(ctx context.Context, id string) (*AgentLocalUpgradeAuth, error) {
	var resp singleResponse[AgentLocalUpgradeAuth]
	if err := c.get(ctx, fmt.Sprintf("/agents/%s/local-upgrade-authorization", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AgentPassphrase is an agent's maintenance passphrase record. The Passphrase
// field is SECRET material.
type AgentPassphrase struct {
	ID                   string `json:"id"`
	UUID                 string `json:"uuid"`
	ComputerName         string `json:"computerName"`
	Domain               string `json:"domain"`
	LastLoggedInUserName string `json:"lastLoggedInUserName"`
	Passphrase           string `json:"passphrase"`
	CreatedAt            string `json:"createdAt"`
	AcknowledgedAt       string `json:"acknowledgedAt"`
	CreatedByUser        string `json:"createdByUser"`

	Raw json.RawMessage `json:"-"`
}

func (a *AgentPassphrase) UnmarshalJSON(b []byte) error {
	type alias AgentPassphrase
	if err := json.Unmarshal(b, (*alias)(a)); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0:0], b...)
	return nil
}

// AgentPassphraseParams are query parameters for listing agent passphrases.
type AgentPassphraseParams struct {
	SiteIDs    []string
	GroupIDs   []string
	AccountIDs []string
	IDs        []string
	Query      string
	Limit      int
	Cursor     string
	CountOnly  bool
}

func (p *AgentPassphraseParams) values() url.Values {
	v := url.Values{}
	if p == nil {
		return v
	}
	addCSV(v, "siteIds", p.SiteIDs)
	addCSV(v, "groupIds", p.GroupIDs)
	addCSV(v, "accountIds", p.AccountIDs)
	addCSV(v, "ids", p.IDs)
	addString(v, "query", p.Query)
	addInt(v, "limit", p.Limit)
	addString(v, "cursor", p.Cursor)
	if p.CountOnly {
		v.Set("countOnly", "true")
	}
	return v
}

// AgentsPassphrases returns a paginated list of agent passphrases. The
// Passphrase field on each item is SECRET material — handle accordingly.
func (c *Client) AgentsPassphrases(ctx context.Context, params *AgentPassphraseParams) ([]AgentPassphrase, *Pagination, error) {
	return list[AgentPassphrase](c, ctx, "/agents/passphrases", params.values())
}
