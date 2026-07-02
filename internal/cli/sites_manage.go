package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newSitesCreateCmd() *cobra.Command {
	var (
		name, accountID, siteType, description, expiration string
		totalLicenses                                      int
		unlimitedLicenses                                  bool
		yes                                                bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a site",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if accountID == "" {
				return fmt.Errorf("--account-id is required")
			}
			action := fmt.Sprintf("create site %q in account %s", name, accountID)
			return guard(cmd.OutOrStdout(), "sites create", action, name, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				s, err := c.SitesCreate(cmd.Context(), mgmt.SiteCreate{
					Name:              name,
					AccountID:         accountID,
					SiteType:          siteType,
					Description:       description,
					Expiration:        expiration,
					UnlimitedLicenses: unlimitedLicenses,
					TotalLicenses:     totalLicenses,
				})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), s)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created site %s (%s)\n", s.Name, s.ID)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "site name (required)")
	cmd.Flags().StringVar(&accountID, "account-id", "", "account ID (required)")
	cmd.Flags().StringVar(&siteType, "site-type", "", "site type")
	cmd.Flags().StringVar(&description, "description", "", "site description")
	cmd.Flags().StringVar(&expiration, "expiration", "", "expiration timestamp (RFC 3339)")
	cmd.Flags().IntVar(&totalLicenses, "total-licenses", 0, "total licenses")
	cmd.Flags().BoolVar(&unlimitedLicenses, "unlimited-licenses", false, "unlimited licenses")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newSitesUpdateCmd() *cobra.Command {
	var (
		name, description, expiration string
		totalLicenses                 int
		unlimitedLicenses             bool
		yes                           bool
	)

	cmd := &cobra.Command{
		Use:   "update <site-id>",
		Short: "Update a site",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var data mgmt.SiteUpdate
			if cmd.Flags().Changed("name") {
				data.Name = &name
			}
			if cmd.Flags().Changed("description") {
				data.Description = &description
			}
			if cmd.Flags().Changed("expiration") {
				data.Expiration = &expiration
			}
			if cmd.Flags().Changed("total-licenses") {
				data.TotalLicenses = &totalLicenses
			}
			if cmd.Flags().Changed("unlimited-licenses") {
				data.UnlimitedLicenses = &unlimitedLicenses
			}
			if data == (mgmt.SiteUpdate{}) {
				return fmt.Errorf("nothing to update: pass at least one field flag")
			}
			return guard(cmd.OutOrStdout(), "sites update", "update site "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				s, err := c.SitesUpdate(cmd.Context(), args[0], data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), s)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated site %s (%s)\n", s.Name, s.ID)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "new site name")
	cmd.Flags().StringVar(&description, "description", "", "new description")
	cmd.Flags().StringVar(&expiration, "expiration", "", "new expiration timestamp (RFC 3339)")
	cmd.Flags().IntVar(&totalLicenses, "total-licenses", 0, "new total licenses")
	cmd.Flags().BoolVar(&unlimitedLicenses, "unlimited-licenses", false, "unlimited licenses")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newSitesDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <site-id>",
		Short: "Delete a site",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "sites delete", "delete site "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.SitesDelete(cmd.Context(), args[0]); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted site %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
