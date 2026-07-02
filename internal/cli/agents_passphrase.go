package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// newAgentLocalUpgradeStatusCmd reads an agent's current local upgrade/downgrade
// authorization. It is a read: no guard, no audit.
func newAgentLocalUpgradeStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "local-upgrade-status <agent-id>",
		Short: "Show an agent's local upgrade/downgrade authorization",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			auth, err := c.AgentsLocalUpgradeAuthGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), auth)
			}
			printTable([]string{"Field", "Value"}, [][]string{
				{"Agent Authorization", orDash(auth.AgentAuthorization)},
				{"Site Authorization", orDash(auth.SiteAuthorization)},
			})
			return nil
		},
	}
}

// newAgentsPassphrasesCmd lists agent maintenance passphrases. The passphrase
// column is SECRET material; the command is a read (no guard, no audit) and
// emits a sensitive-output notice on stderr.
func newAgentsPassphrasesCmd() *cobra.Command {
	var siteIDs, groupIDs, accountIDs, ids []string
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "passphrases",
		Short: "List agent maintenance passphrases (SECRET output)",
		Long: `List agent maintenance passphrases. The passphrase column is secret
material used to run privileged local agent commands — treat the output
accordingly. Values are never written to the audit log.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AgentPassphraseParams{
				SiteIDs:    siteIDs,
				GroupIDs:   groupIDs,
				AccountIDs: accountIDs,
				IDs:        ids,
				Query:      query,
				Limit:      limit,
				Cursor:     cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []mgmt.AgentPassphrase
			var total int
			if all {
				items, total, err = fetchAllREST("passphrase", func(cur string) ([]mgmt.AgentPassphrase, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.AgentsPassphrases(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				items, pag, err = c.AgentsPassphrases(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Computer", "UUID", "Passphrase", "Created"}
			rows := make([][]string, len(items))
			for i, p := range items {
				rows[i] = []string{p.ID, p.ComputerName, p.UUID, p.Passphrase, orDash(p.CreatedAt)}
			}
			if err := printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "passphrase", all); err != nil {
				return err
			}
			noteSensitiveOutput(cmd.ErrOrStderr())
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&ids, "id", nil, "filter by agent ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	return cmd
}
