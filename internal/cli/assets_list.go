package cli

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAssetsListCmd() *cobra.Command {
	var (
		assetType  string
		filters    []string
		limit      int
		skip       int
		cursor     string
		sortBy     string
		sortOrder  string
		siteIDs    []string
		accountIDs []string
		groupIDs   []string
		all        bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List assets from the XDR inventory",
		Long: `List assets from the XDR asset inventory.

When --type is omitted, lists assets across all types.
Use --filter key=value to pass type-specific API query parameters.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.XDRAssetListParams{
				Limit:      limit,
				Skip:       skip,
				Cursor:     cursor,
				SortBy:     sortBy,
				SortOrder:  sortOrder,
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
			}
			extra, err := parseFilterKV(filters)
			if err != nil {
				return err
			}
			params.Extra = extra
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var items []json.RawMessage
			var total int

			if all {
				items, total, err = fetchAllREST("asset", func(cur string) ([]json.RawMessage, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.XDRAssetList(cmd.Context(), mgmt.AssetType(assetType), params)
				})
				if err != nil {
					return err
				}
			} else {
				page, pag, fetchErr := c.XDRAssetList(cmd.Context(), mgmt.AssetType(assetType), params)
				if fetchErr != nil {
					return fetchErr
				}
				items = page
				if pag != nil {
					total = pag.TotalItems
				}
			}

			if outputFormat == "json" {
				env := struct {
					Data       []json.RawMessage `json:"data"`
					Returned   int               `json:"returned"`
					Total      int               `json:"total"`
					NextCursor string            `json:"nextCursor,omitempty"`
				}{
					Data:       items,
					Returned:   len(items),
					Total:      total,
					NextCursor: params.Cursor,
				}
				if all {
					env.NextCursor = ""
				}
				return printJSON(cmd.OutOrStdout(), env)
			}

			headers, rows := assetRows(items, assetType)
			if total == 0 {
				total = len(items)
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), total, "asset", all)
		},
	}
	cmd.Flags().StringVar(&assetType, "type", "", "asset type slug (e.g. device, server, surface/cloud)")
	cmd.Flags().StringArrayVar(&filters, "filter", nil, `key=value filter (e.g. --filter osTypes=windows)`)
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page")
	cmd.Flags().IntVar(&skip, "skip", 0, "number of results to skip")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort field")
	cmd.Flags().StringVar(&sortOrder, "sort-order", "", "sort direction (asc, desc)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	return markJSON(cmd)
}

// assetRows extracts table rows from raw JSON asset items.
// Without a type, shows ID/Name/Type/Category/Surface; with a type,
// shows ID/Name from the JSON object.
func assetRows(items []json.RawMessage, assetType string) ([]string, [][]string) {
	if assetType == "" {
		headers := []string{"ID", "Name", "Type", "Category", "Surface"}
		rows := make([][]string, len(items))
		for i, raw := range items {
			m := jsonMap(raw)
			rows[i] = []string{
				jsonStr(m, "id"),
				truncate(jsonStr(m, "name"), 50),
				jsonStr(m, "assetType"),
				jsonStr(m, "category"),
				jsonStr(m, "surface"),
			}
		}
		return headers, rows
	}
	headers := []string{"ID", "Name"}
	rows := make([][]string, len(items))
	for i, raw := range items {
		m := jsonMap(raw)
		rows[i] = []string{
			jsonStr(m, "id"),
			truncate(jsonStr(m, "name"), 60),
		}
	}
	return headers, rows
}

func newAssetsExportCmd() *cobra.Command {
	var (
		assetType  string
		filters    []string
		outputFile string
		siteIDs    []string
		accountIDs []string
		groupIDs   []string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export assets from the XDR inventory",
		Long: `Export assets as raw CSV or JSON from the API.

Streams the raw export response to a file or stdout.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.XDRAssetListParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
			}
			extra, err := parseFilterKV(filters)
			if err != nil {
				return err
			}
			params.Extra = extra

			data, err := c.XDRAssetExport(cmd.Context(), mgmt.AssetType(assetType), params)
			if err != nil {
				return err
			}
			if outputFile != "" {
				if err := os.WriteFile(outputFile, data, 0o644); err != nil {
					return fmt.Errorf("write export file: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Exported to %s (%d bytes)\n", outputFile, len(data))
				return nil
			}
			_, err = cmd.OutOrStdout().Write(data)
			return err
		},
	}
	cmd.Flags().StringVar(&assetType, "type", "", "asset type slug (e.g. device, server)")
	cmd.Flags().StringArrayVar(&filters, "filter", nil, `key=value filter`)
	cmd.Flags().StringVar(&outputFile, "output-file", "", "write export to file instead of stdout")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	return cmd
}

func newAssetsFilterOptionsCmd() *cobra.Command {
	var (
		assetType  string
		siteIDs    []string
		accountIDs []string
		groupIDs   []string
	)

	cmd := &cobra.Command{
		Use:   "filter-options",
		Short: "Show available filter fields for an asset type",
		Long: `Show available filters for the given asset type.

Combines autocomplete and free-text filter information.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.XDRAssetListParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
			}
			at := mgmt.AssetType(assetType)

			autocomplete, err := c.XDRAssetFilterAutocomplete(cmd.Context(), at, params)
			if err != nil {
				return err
			}
			freeText, err := c.XDRAssetFilterFreeText(cmd.Context(), at, params)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]any{
					"autocomplete": autocomplete,
					"freeText":     freeText,
				})
			}

			w := cmd.OutOrStdout()
			fmt.Fprintln(w, "Autocomplete filters:")
			for _, raw := range autocomplete {
				fmt.Fprintf(w, "  %s\n", string(raw))
			}
			fmt.Fprintln(w)
			fmt.Fprintln(w, "Free-text filters:")
			for _, raw := range freeText {
				fmt.Fprintf(w, "  %s\n", string(raw))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&assetType, "type", "", "asset type slug (required)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	_ = cmd.MarkFlagRequired("type")
	return markJSON(cmd)
}

func newAssetsSubCategoriesCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string

	cmd := &cobra.Command{
		Use:   "sub-categories",
		Short: "List asset sub-categories",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.XDRAssetCountsParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
			}
			data, err := c.XDRAssetSubCategories(cmd.Context(), params)
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), data)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	return markJSON(cmd)
}

// parseFilterKV converts --filter key=value pairs into url.Values.
func parseFilterKV(filters []string) (url.Values, error) {
	if len(filters) == 0 {
		return nil, nil
	}
	v := url.Values{}
	for _, f := range filters {
		key, val, ok := strings.Cut(f, "=")
		if !ok {
			return nil, fmt.Errorf("invalid filter %q: expected key=value", f)
		}
		v.Add(key, val)
	}
	return v, nil
}

// jsonMap unmarshals raw JSON into a map for table rendering.
func jsonMap(raw json.RawMessage) map[string]any {
	var m map[string]any
	_ = json.Unmarshal(raw, &m)
	return m
}

// jsonStr extracts a string value from a JSON map, returning "" for missing
// or non-string values.
func jsonStr(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return fmt.Sprint(v)
}
