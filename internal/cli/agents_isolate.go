package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAgentsIsolateCmd() *cobra.Command {
	var filters []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "isolate [agent-id...]",
		Short: "Isolate agents from the network",
		Long: `Disconnect agents from the network.

Specify agent IDs as arguments, or use --filter to match agents by API
query parameters (e.g. --filter infected=true --filter osTypes=windows).
Both can be combined. Dry-run by default; pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := collectAgentIDs(cmd, args, filters)
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No matching agents found.")
				return nil
			}
			return guard(cmd.OutOrStdout(), "agents isolate", "isolate "+pluralize(len(ids), "agent"), strings.Join(ids, ","), yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsDisconnect(cmd.Context(), mgmt.ActionFilter{IDs: ids})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "isolate: %s affected\n", pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringArrayVar(&filters, "filter", nil, `key=value filter (e.g. --filter infected=true)`)
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newAgentsReconnectCmd() *cobra.Command {
	var filters []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "reconnect [agent-id...]",
		Short: "Reconnect isolated agents",
		Long: `Reconnect previously network-isolated agents.

Specify agent IDs as arguments, or use --filter to match agents by API
query parameters (e.g. --filter networkStatuses=disconnected).
Both can be combined. Dry-run by default; pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := collectAgentIDs(cmd, args, filters)
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No matching agents found.")
				return nil
			}
			return guard(cmd.OutOrStdout(), "agents reconnect", "reconnect "+pluralize(len(ids), "agent"), strings.Join(ids, ","), yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsConnect(cmd.Context(), mgmt.ActionFilter{IDs: ids})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "reconnect: %s affected\n", pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringArrayVar(&filters, "filter", nil, `key=value filter (e.g. --filter networkStatuses=disconnected)`)
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

// collectAgentIDs merges positional IDs with filter-matched agent IDs.
func collectAgentIDs(cmd *cobra.Command, args, filters []string) ([]string, error) {
	hasFilters := len(filters) > 0
	if len(args) == 0 && !hasFilters {
		return nil, fmt.Errorf("specify agent IDs or use --filter")
	}

	var ids []string
	if hasFilters {
		params, err := parseAgentFilters(filters)
		if err != nil {
			return nil, err
		}
		c, err := mgmtClient()
		if err != nil {
			return nil, err
		}
		agents, _, err := fetchAllREST("agent", func(cur string) ([]mgmt.Agent, *mgmt.Pagination, error) {
			params.Cursor = cur
			return c.AgentsList(cmd.Context(), params)
		})
		if err != nil {
			return nil, err
		}
		for _, a := range agents {
			ids = append(ids, a.ID)
		}
	}
	ids = append(ids, args...)
	return ids, nil
}

// parseAgentFilters converts key=value filter strings to AgentListParams.
// Keys match API query parameter names.
func parseAgentFilters(filters []string) (*mgmt.AgentListParams, error) {
	params := &mgmt.AgentListParams{Limit: 1000}
	for _, f := range filters {
		key, val, ok := strings.Cut(f, "=")
		if !ok {
			return nil, fmt.Errorf("invalid filter %q: expected key=value", f)
		}
		switch key {
		case "siteIds":
			params.SiteIDs = append(params.SiteIDs, val)
		case "groupIds":
			params.GroupIDs = append(params.GroupIDs, val)
		case "accountIds":
			params.AccountIDs = append(params.AccountIDs, val)
		case "osTypes":
			params.OSTypes = append(params.OSTypes, val)
		case "networkStatuses":
			params.NetworkStatuses = append(params.NetworkStatuses, val)
		case "machineTypes":
			params.MachineTypes = append(params.MachineTypes, val)
		case "infected":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid filter %s=%s: expected true or false", key, val)
			}
			params.Infected = &b
		case "isActive":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid filter %s=%s: expected true or false", key, val)
			}
			params.IsActive = &b
		case "isDecommissioned":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid filter %s=%s: expected true or false", key, val)
			}
			params.IsDecommissioned = &b
		case "isUninstalled":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid filter %s=%s: expected true or false", key, val)
			}
			params.IsUninstalled = &b
		case "isUpToDate":
			b, err := strconv.ParseBool(val)
			if err != nil {
				return nil, fmt.Errorf("invalid filter %s=%s: expected true or false", key, val)
			}
			params.IsUpToDate = &b
		case "query":
			params.Query = val
		default:
			return nil, fmt.Errorf("unknown agent filter key %q", key)
		}
	}
	return params, nil
}
