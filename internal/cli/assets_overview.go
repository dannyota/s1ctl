package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAssetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Manage XDR asset inventory",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAssetsOverviewCmd())
	cmd.AddCommand(newAssetsCategoriesCmd())
	return cmd
}

func newAssetsOverviewCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string

	cmd := &cobra.Command{
		Use:   "overview",
		Short: "Show asset counts by category and surface",
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
			counts, err := c.XDRAssetCounts(cmd.Context(), params)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), counts)
			}

			w := cmd.OutOrStdout()
			cat := counts.Categories

			catRows := []struct {
				name  string
				count int
			}{
				{"Account", cat.Account.Count},
				{"AI/ML", cat.AiMl.Count},
				{"Application Integration", cat.ApplicationIntegration.Count},
				{"Cloud Application", cat.CloudApplication.Count},
				{"Code", cat.Code.Count},
				{"Container", cat.Container.Count},
				{"Data Analysis", cat.DataAnalysis.Count},
				{"Data Store", cat.DataStore.Count},
				{"Developer Tool", cat.DeveloperTool.Count},
				{"Device", cat.Device.Count},
				{"Function", cat.Function.Count},
				{"Governance", cat.Governance.Count},
				{"Identity", cat.Identity.Count},
				{"Inventory", cat.Inventory.Count},
				{"Network", cat.Network.Count},
				{"Secrets", cat.Secrets.Count},
				{"Server", cat.Server.Count},
				{"Storage", cat.Storage.Count},
				{"Workstation", cat.Workstation.Count},
			}

			var rows [][]string
			for _, r := range catRows {
				if r.count > 0 {
					rows = append(rows, []string{r.name, strconv.Itoa(r.count)})
				}
			}

			surf := counts.Surfaces
			surfRows := []struct {
				name  string
				count int
			}{
				{"Cloud", surf.Cloud.Count},
				{"Endpoint", surf.Endpoint.Count},
				{"Identity", surf.Identity.Count},
				{"Network", surf.Network.Count},
				{"Network Discovery", surf.NetworkDiscovery.Count},
			}
			var srows [][]string
			for _, r := range surfRows {
				srows = append(srows, []string{r.name, strconv.Itoa(r.count)})
			}

			if outputFormat == "csv" {
				var crows [][]string
				for _, r := range rows {
					crows = append(crows, append([]string{"category"}, r...))
				}
				for _, r := range srows {
					crows = append(crows, append([]string{"surface"}, r...))
				}
				return printCSV(w, []string{"Type", "Name", "Count"}, crows)
			}

			printTable([]string{"Category", "Count"}, rows)
			fmt.Fprintln(w)
			printTable([]string{"Surface", "Count"}, srows)

			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	return cmd
}

func newAssetsCategoriesCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string

	cmd := &cobra.Command{
		Use:   "categories",
		Short: "List asset categories with counts",
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
			cats, err := c.XDRAssetCategories(cmd.Context(), params)
			if err != nil {
				return err
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), cats)
			}

			catRows := []struct {
				name  string
				count int
			}{
				{"Account", cats.Account},
				{"Container", cats.Container},
				{"Device", cats.Device},
				{"Identity", cats.Identity},
				{"Inventory", cats.Inventory},
				{"Server", cats.Server},
				{"Storage", cats.Storage},
				{"Workstation", cats.Workstation},
			}

			headers := []string{"Category", "Count"}
			rows := make([][]string, len(catRows))
			for i, r := range catRows {
				rows[i] = []string{r.name, strconv.Itoa(r.count)}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, cats, len(catRows), len(catRows), "category", true)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	return cmd
}
