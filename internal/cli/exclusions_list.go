package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newExclusionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exclusions",
		Short: "Manage exclusions and blocklist",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newExclusionsListCmd())
	cmd.AddCommand(newExclusionsGetCmd())
	addExclusionSyncCmds(cmd)
	return cmd
}

func newExclusionsListCmd() *cobra.Command {
	var siteIDs, types, osTypes []string
	var query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List exclusions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			exclusions, pag, err := c.ExclusionsList(cmd.Context(), &mgmt.ExclusionListParams{
				SiteIDs: siteIDs,
				Types:   types,
				OSTypes: osTypes,
				Query:   query,
				Limit:   limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(exclusions)
			}
			var rows [][]string
			for _, e := range exclusions {
				rows = append(rows, []string{
					e.ID, e.Type, truncate(e.Value, 50), e.OSType, e.Mode,
				})
			}
			printTable([]string{"ID", "Type", "Value", "OS", "Mode"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "exclusion"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&types, "type", nil, "filter by exclusion type")
	cmd.Flags().StringSliceVar(&osTypes, "os-type", nil, "filter by OS type")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}

func newExclusionsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <exclusion-id>",
		Short: "Get exclusion details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			e, err := c.ExclusionsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(e)
			}
			rows := [][]string{
				{"ID", e.ID},
				{"Type", e.Type},
				{"Value", e.Value},
				{"OS", e.OSType},
				{"Mode", e.Mode},
				{"Description", e.Description},
				{"Scope", e.ScopeName},
				{"User", e.UserName},
				{"Created", e.CreatedAt},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}
