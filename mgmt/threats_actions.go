package mgmt

import (
	"context"
	"fmt"
)

// ThreatsMitigate applies a mitigation action to threats.
func (c *Client) ThreatsMitigate(ctx context.Context, action string, filter ActionFilter) (int, error) {
	if action == "" {
		return 0, fmt.Errorf("mgmt: mitigation action is required")
	}
	return doAction(c, ctx, fmt.Sprintf("/threats/mitigate/%s", action), filter, nil)
}

// ThreatsUpdateVerdict updates the analyst verdict on threats.
func (c *Client) ThreatsUpdateVerdict(ctx context.Context, verdict string, filter ActionFilter) (int, error) {
	data := map[string]string{"analystVerdict": verdict}
	return doAction(c, ctx, "/threats/analyst-verdict", filter, data)
}

// ThreatsUpdateStatus updates the incident status on threats.
func (c *Client) ThreatsUpdateStatus(ctx context.Context, status string, filter ActionFilter) (int, error) {
	data := map[string]string{"incidentStatus": status}
	return doAction(c, ctx, "/threats/incident", filter, data)
}
