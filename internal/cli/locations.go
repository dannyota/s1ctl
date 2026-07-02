package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newLocationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "locations",
		Short: "Manage firewall locations",
		Long: `Manage firewall location definitions.

A location is detected by endpoint network parameters (IP, DNS, NIC, registry
key, or management connectivity) and drives Location Aware firewall rules.

'create' and 'update' set the location name, description, and match operator via
flags. Detection parameters are round-tripped through 'locations pull' /
'locations push', which preserves the full definition as YAML.`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newLocationsListCmd())
	cmd.AddCommand(newLocationsCreateCmd())
	cmd.AddCommand(newLocationsUpdateCmd())
	cmd.AddCommand(newLocationsDeleteCmd())
	addLocationSyncCmds(cmd)
	return cmd
}

// validOperator maps a CLI --operator value to the typed enum, or errors.
func validOperator(op string) (mgmt.LocationOperator, error) {
	switch mgmt.LocationOperator(op) {
	case mgmt.LocationOperatorAll, mgmt.LocationOperatorAny, mgmt.LocationOperatorNone:
		return mgmt.LocationOperator(op), nil
	default:
		return "", fmt.Errorf("--operator must be one of: all, any, none")
	}
}

func newLocationsListCmd() *cobra.Command {
	var siteIDs, accountIDs, groupIDs []string
	var cursor string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List firewall locations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.LocationListParams{
				SiteIDs:    siteIDs,
				AccountIDs: accountIDs,
				GroupIDs:   groupIDs,
				Limit:      limit,
				Cursor:     cursor,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var locs []mgmt.Location
			var total int
			if all {
				locs, total, err = fetchAllREST("location", func(cur string) ([]mgmt.Location, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.LocationsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				locs, pag, err = c.LocationsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Operator", "Scope", "Fallback"}
			rows := make([][]string, len(locs))
			for i, l := range locs {
				rows[i] = []string{l.ID, l.Name, string(l.Operator), l.ScopeName, boolIcon(l.IsFallback)}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, locs, len(locs), total, "location", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "filter by account ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	cmd.Flags().StringVar(&cursor, "cursor", "", "pagination cursor")
	return cmd
}

func newLocationsCreateCmd() *cobra.Command {
	var name, description, operator string
	var siteIDs, accountIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "create --name <name>",
		Short: "Create a firewall location",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			op, err := validOperator(operator)
			if err != nil {
				return err
			}
			body := mgmt.LocationCreate{
				Data:   mgmt.LocationData{Name: name, Description: description, Operator: op},
				Filter: mgmt.LocationScope{SiteIDs: siteIDs, AccountIDs: accountIDs},
			}

			action := fmt.Sprintf("create location %q", name)
			return guard(cmd.OutOrStdout(), "locations create", action, name, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.LocationsCreate(cmd.Context(), body)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created location %s (%s)\n", created.ID, created.Name)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "location name (required)")
	cmd.Flags().StringVar(&description, "description", "", "location description")
	cmd.Flags().StringVar(&operator, "operator", "any", "match operator: all, any, none")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "create in these site IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "create in these account IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newLocationsUpdateCmd() *cobra.Command {
	var name, description, operator string
	var yes bool

	cmd := &cobra.Command{
		Use:   "update <location-id> --name <name>",
		Short: "Update a firewall location's name, description, or operator",
		Long: `Update a location's name, description, and match operator.

Note: this replaces the location definition with the supplied fields; detection
parameters not expressible as flags are managed through 'locations push'.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			op, err := validOperator(operator)
			if err != nil {
				return err
			}
			body := mgmt.LocationUpdate{
				Data: mgmt.LocationData{Name: name, Description: description, Operator: op},
			}

			action := fmt.Sprintf("update location %s", args[0])
			return guard(cmd.OutOrStdout(), "locations update", action, args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if _, err := c.LocationsUpdate(cmd.Context(), args[0], body); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated location %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "location name (required)")
	cmd.Flags().StringVar(&description, "description", "", "location description")
	cmd.Flags().StringVar(&operator, "operator", "any", "match operator: all, any, none")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newLocationsDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <location-id>",
		Short: "Delete a firewall location",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "locations delete", "delete location "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.LocationsDelete(cmd.Context(), []string{args[0]}); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted location %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
