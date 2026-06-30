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
