package mgmt

import "context"

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

// AgentsUpdateSoftware triggers a software update on agents.
func (c *Client) AgentsUpdateSoftware(ctx context.Context, filter ActionFilter) (int, error) {
	return doAction(c, ctx, "/agents/actions/update-software", filter, nil)
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
