package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newFiltersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filters",
		Short: "Manage saved endpoint filters",
		Long: `Manage saved endpoint filters.

A saved filter pairs a name with a filterFields definition (the set of endpoint
criteria to match). Saved filters back bulk agent actions and dynamic groups.

A filter body can be large, so create and update read a JSON file (--from-file)
with the filter name and its filterFields, rather than taking many flags.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newFiltersListCmd())
	cmd.AddCommand(newFiltersCreateCmd())
	cmd.AddCommand(newFiltersUpdateCmd())
	cmd.AddCommand(newFiltersDeleteCmd())
	return cmd
}

// filterFile is the declarative body read from --from-file: the filter name and
// its filterFields criteria set. JSON only, since filterFields is arbitrary
// nested JSON.
type filterFile struct {
	Name         string          `json:"name"`
	FilterFields json.RawMessage `json:"filterFields"`
}

func readFilterFile(path string) (filterFile, error) {
	var f filterFile
	raw, err := os.ReadFile(path)
	if err != nil {
		return f, fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(raw, &f); err != nil {
		return f, fmt.Errorf("parse %s: %w", path, err)
	}
	if f.Name == "" {
		return f, fmt.Errorf("filter file %s has no name", path)
	}
	return f, nil
}

func newFiltersListCmd() *cobra.Command {
	var siteIDs, accountIDs []string
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List saved filters",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.FilterListParams{
				Query:      query,
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				Limit:      limit,
				Cursor:     cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var filters []mgmt.Filter
			var total int
			if all {
				filters, total, err = fetchAllREST("filter", func(cur string) ([]mgmt.Filter, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.FiltersList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				filters, pag, err = c.FiltersList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Scope", "ScopeID"}
			rows := make([][]string, len(filters))
			for i, f := range filters {
				rows[i] = []string{f.ID, f.Name, f.ScopeLevel, f.ScopeID}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, filters, len(filters), total, "filter", all)
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "free text search on filter name")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}

func newFiltersCreateCmd() *cobra.Command {
	var fromFile string
	var siteIDs, accountIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "create --from-file <filter.json>",
		Short: "Create a saved filter from a JSON file",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			f, err := readFilterFile(fromFile)
			if err != nil {
				return err
			}
			body := mgmt.FilterCreate{
				Data: mgmt.FilterData{Name: f.Name, FilterFields: f.FilterFields},
			}
			if len(siteIDs) > 0 || len(accountIDs) > 0 {
				body.Filter = &mgmt.FilterScope{SiteIDs: siteIDs, AccountIDs: accountIDs}
			}

			action := fmt.Sprintf("create filter %q from %s", f.Name, fromFile)
			return guard(cmd.OutOrStdout(), "filters create", action, fromFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.FiltersCreate(cmd.Context(), body)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created filter %s (%s)\n", created.ID, created.Name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "filter definition JSON file (required)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "create in these site IDs (default: global/tenant)")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "create in these account IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newFiltersUpdateCmd() *cobra.Command {
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update <filter-id> --from-file <filter.json>",
		Short: "Update a saved filter from a JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			f, err := readFilterFile(fromFile)
			if err != nil {
				return err
			}
			body := mgmt.FilterUpdate{Data: mgmt.FilterData{Name: f.Name, FilterFields: f.FilterFields}}

			action := fmt.Sprintf("update filter %s from %s", args[0], fromFile)
			return guard(cmd.OutOrStdout(), "filters update", action, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				updated, err := c.FiltersUpdate(cmd.Context(), args[0], body)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), updated)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated filter %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "filter definition JSON file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newFiltersDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <filter-id>",
		Short: "Delete a saved filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "filters delete", "delete filter "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.FiltersDelete(cmd.Context(), args[0]); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted filter %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
