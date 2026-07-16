package cli

import (
	"github.com/spf13/cobra"

	"danny.vn/s1/graphql"
)

func newCloudOnboardingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud-onboarding",
		Short: "Manage CNAPP cloud account onboarding",
		Long: `Manage cloud account onboarding for Cloud Native Application
Protection Platform (CNAPP). Supports listing and inspecting onboarded
cloud entities (AWS, GCP, Azure, OCI, Alibaba accounts and organizations),
onboarding new entities from a JSON file, and deleting (offboarding) entities.

Uses the cloud onboarding GraphQL API.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newCloudOnboardingListCmd())
	cmd.AddCommand(newCloudOnboardingGetCmd())
	cmd.AddCommand(newCloudOnboardingOnboardCmd())
	cmd.AddCommand(newCloudOnboardingDeleteCmd())
	return cmd
}

func newCloudOnboardingListCmd() *cobra.Command {
	var providers, statuses []string
	var after string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List onboarded cloud entities",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}

			var filters []graphql.CnappFilter
			if len(providers) > 0 {
				filters = append(filters, graphql.CnappFilter{
					FieldID:  "cloudProvider",
					StringIn: &graphql.CnappInStr{Values: providers},
				})
			}
			if len(statuses) > 0 {
				filters = append(filters, graphql.CnappFilter{
					FieldID:  "onboardingStatus",
					StringIn: &graphql.CnappInStr{Values: statuses},
				})
			}
			if filters == nil {
				filters = []graphql.CnappFilter{}
			}

			page := &graphql.ListParams{First: limit, After: after}
			if page.First == 0 {
				page.First = defaultPageSize
			}

			var items []graphql.CnappOnboardedEntity
			var total int64

			if all {
				items, total, err = fetchAllGQL("cloud entity", func(cur string) (*graphql.Connection[graphql.CnappOnboardedEntity], error) {
					page.After = cur
					return c.CnappEntitiesList(cmd.Context(), filters, nil, page)
				})
			} else {
				conn, connErr := c.CnappEntitiesList(cmd.Context(), filters, nil, page)
				if connErr != nil {
					return connErr
				}
				total = conn.TotalCount
				for _, edge := range conn.Edges {
					items = append(items, edge.Node)
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Entity ID", "Name", "Type", "Status", "Coverage", "Scope"}
			rows := make([][]string, len(items))
			for i, e := range items {
				rows[i] = []string{
					orDash(e.ID),
					orDash(e.EntityID),
					truncate(orDash(e.Name), 30),
					string(e.Type),
					string(e.OnboardingStatus),
					joinOrDash(e.ActiveCoverage),
					orDash(e.Scope),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, items, len(items), int(total), "cloud entity", all)
		},
	}
	cmd.Flags().StringSliceVar(&providers, "provider", nil, "filter by cloud provider (AWS, GCP, AZURE, OCI, ALIBABA)")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by operational status")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&after, "after", "", "pagination cursor")
	return markJSON(cmd)
}

func newCloudOnboardingGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <account-id>",
		Short: "Get onboarded cloud entity details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := gqlClient()
			if err != nil {
				return err
			}
			detail, err := c.CnappEntityGet(cmd.Context(), []string{args[0]}, nil)
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), detail)
			}
			rows := [][]string{
				{"Entity ID", orDash(detail.EntityID)},
				{"Entity Name", orDash(detail.EntityName)},
				{"Display Name", orDash(detail.DisplayName)},
				{"Onboarding Type", string(detail.OnboardingType)},
				{"Cloud Provider", string(detail.CloudProvider)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	return markJSON(cmd)
}
