package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newThreatNotesCmd() *cobra.Command {
	var cursor, sortBy, sortOrder string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "notes <threat-id>",
		Short: "List notes for a threat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ThreatNotesListParams{
				Limit:     limit,
				Cursor:    cursor,
				SortBy:    sortBy,
				SortOrder: sortOrder,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var notes []mgmt.ThreatNote
			var total int

			if all {
				notes, total, err = fetchAllREST("note", func(cur string) ([]mgmt.ThreatNote, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.ThreatNotesList(cmd.Context(), args[0], params)
				})
			} else {
				var pag *mgmt.Pagination
				notes, pag, err = c.ThreatNotesList(cmd.Context(), args[0], params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Text", "Creator", "Created"}
			rows := make([][]string, len(notes))
			for i, n := range notes {
				rows[i] = []string{
					n.ID, truncate(n.Text, 60), n.Creator, n.CreatedAt,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, notes, len(notes), total, "note", all)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	return markJSON(cmd)
}

func newThreatAddNoteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "add-note <threat-id> <text>",
		Short: "Add a note to a threat",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			threatID, text := args[0], args[1]
			return guard(cmd.OutOrStdout(), "threats add-note", "add note to threat "+threatID, threatID, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.ThreatNotesCreate(cmd.Context(), threatID, text)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Note added to %s\n", pluralize(affected, "threat"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}
