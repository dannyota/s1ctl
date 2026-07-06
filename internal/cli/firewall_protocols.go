package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newFirewallProtocolsCmd() *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:   "protocols",
		Short: "List available firewall protocols",
		Long:  `Show protocols that can be used in firewall rules.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.FirewallProtocolListParams{
				Query: query,
				Limit: 1000,
			}

			protocols, _, err := c.FirewallProtocolsList(cmd.Context(), params)
			if err != nil {
				return err
			}

			headers := []string{"Value", "Name"}
			rows := make([][]string, len(protocols))
			for i, p := range protocols {
				rows[i] = []string{p.Value, p.Name}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, protocols, len(protocols), len(protocols), "protocol", true)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search protocols")
	return markJSON(cmd)
}
