package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage endpoint tags",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newTagsListCmd())
	return cmd
}

func newTagsListCmd() *cobra.Command {
	var siteIDs []string
	var tagType, query string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tags",
		Long:  "Types: endpoint, firewall, network-quarantine",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if tagType == "" {
				return fmt.Errorf("--type is required (endpoint, firewall, network-quarantine)")
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			tags, pag, err := c.TagsList(cmd.Context(), &mgmt.TagListParams{
				Type:    tagType,
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
			})
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(tags)
			}
			var rows [][]string
			for _, t := range tags {
				rows = append(rows, []string{t.ID, t.Key, t.Value, t.Scope})
			}
			printTable([]string{"ID", "Key", "Value", "Scope"}, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", pluralize(pag.TotalItems, "tag"))
			return nil
		},
	}
	cmd.Flags().StringVar(&tagType, "type", "", "tag type (endpoint, firewall, network-quarantine)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results")
	return cmd
}
