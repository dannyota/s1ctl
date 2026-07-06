package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newDeviceControlEventsCmd() *cobra.Command {
	var siteIDs, interfaces []string
	var query, cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "events",
		Short: "List device control events",
		Long:  `Show device control events from endpoints with Device Control-enabled Agents.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.DeviceEventListParams{
				SiteIDs:    siteIDs,
				Query:      query,
				Interfaces: interfaces,
				Limit:      limit,
				Cursor:     cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var events []mgmt.DeviceEvent
			var total int

			if all {
				events, total, err = fetchAllREST("device event", func(cur string) ([]mgmt.DeviceEvent, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.DeviceEventsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				events, pag, err = c.DeviceEventsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Event Type", "Interface", "Device", "Agent ID", "Time"}
			rows := make([][]string, len(events))
			for i, e := range events {
				rows[i] = []string{
					e.ID, e.EventType, e.Interface,
					e.DeviceName, e.AgentID, e.EventTime,
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, events, len(events), total, "device event", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&interfaces, "interface", nil, "filter by interface (USB, Bluetooth, Thunderbolt, SDCard)")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return markJSON(cmd)
}
