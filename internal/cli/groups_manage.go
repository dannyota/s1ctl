package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newGroupsCreateCmd() *cobra.Command {
	var siteID, name, description string
	var yes bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a group",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if siteID == "" {
				return fmt.Errorf("--site-id is required")
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would create group %q in site %s. Pass --yes to apply.\n",
					name, siteID)
				return nil
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			g, err := c.GroupsCreate(cmd.Context(), siteID, mgmt.GroupCreate{
				Name:        name,
				Description: description,
			})
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), g)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created group %s (%s)\n", g.Name, g.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&siteID, "site-id", "", "site ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "group name (required)")
	cmd.Flags().StringVar(&description, "description", "", "group description")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newGroupsDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <group-id>",
		Short: "Delete a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			if !yes {
				fmt.Fprintf(cmd.OutOrStdout(), "Would delete group %s. Pass --yes to apply.\n", id)
				return nil
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			if err := c.GroupsDelete(cmd.Context(), id); err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": id})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted group %s\n", id)
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
