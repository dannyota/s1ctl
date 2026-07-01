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
	var tagType, query, cursor string
	var limit int
	var all bool

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
			params := &mgmt.TagListParams{
				Type:    tagType,
				SiteIDs: siteIDs,
				Query:   query,
				Limit:   limit,
				Cursor:  cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var tags []mgmt.Tag
			var total int

			if all {
				tags, total, err = fetchAllREST("tag", func(cur string) ([]mgmt.Tag, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.TagsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				tags, pag, err = c.TagsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Key", "Value", "Scope"}
			rows := make([][]string, len(tags))
			for i, t := range tags {
				rows[i] = []string{t.ID, t.Key, t.Value, t.Scope}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, tags, len(tags), total, "tag", all)
		},
	}
	cmd.Flags().StringVar(&tagType, "type", "", "tag type (endpoint, firewall, network-quarantine)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}
