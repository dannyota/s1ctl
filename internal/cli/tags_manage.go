package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newTagsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <tag-id>",
		Short: "Get a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			t, err := c.TagsGet(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), t)
			}
			printTable([]string{"FIELD", "VALUE"}, [][]string{
				{"ID", t.ID},
				{"Key", t.Key},
				{"Value", t.Value},
				{"Description", orDash(t.Description)},
			})
			return nil
		},
	}
	return markJSON(cmd)
}

func newTagsCreateCmd() *cobra.Command {
	var (
		key, value, description, scope, scopeID string
		yes                                     bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a tag",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if key == "" {
				return fmt.Errorf("--key is required")
			}
			if value == "" {
				return fmt.Errorf("--value is required")
			}
			action := fmt.Sprintf("create tag %s=%s", key, value)
			return guard(cmd.OutOrStdout(), "tags create", action, key, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				t, err := c.TagsCreate(cmd.Context(), mgmt.TagCreate{
					Key:         key,
					Value:       value,
					Description: description,
					Scope:       scope,
					ScopeID:     scopeID,
				})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), t)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created tag %s=%s (%s)\n", t.Key, t.Value, t.ID)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&key, "key", "", "tag key (required)")
	cmd.Flags().StringVar(&value, "value", "", "tag value (required)")
	cmd.Flags().StringVar(&description, "description", "", "tag description")
	cmd.Flags().StringVar(&scope, "scope", "", "tag scope")
	cmd.Flags().StringVar(&scopeID, "scope-id", "", "tag scope ID")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newTagsUpdateCmd() *cobra.Command {
	var (
		key, value, description string
		yes                     bool
	)

	cmd := &cobra.Command{
		Use:   "update <tag-id>",
		Short: "Update a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var data mgmt.TagUpdate
			if cmd.Flags().Changed("key") {
				data.Key = &key
			}
			if cmd.Flags().Changed("value") {
				data.Value = &value
			}
			if cmd.Flags().Changed("description") {
				data.Description = &description
			}
			if data == (mgmt.TagUpdate{}) {
				return fmt.Errorf("nothing to update: pass at least one field flag")
			}
			return guard(cmd.OutOrStdout(), "tags update", "update tag "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				t, err := c.TagsUpdate(cmd.Context(), args[0], data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), t)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated tag %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&key, "key", "", "new tag key")
	cmd.Flags().StringVar(&value, "value", "", "new tag value")
	cmd.Flags().StringVar(&description, "description", "", "new description")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}

func newTagsDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <tag-id>",
		Short: "Delete a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "tags delete", "delete tag "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.TagsDelete(cmd.Context(), args[0]); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "deleted", "id": args[0]})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deleted tag %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return markJSON(cmd)
}
