package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// blocklistScope builds the write scope from CLI flags: the given scope IDs, or
// the global (tenant) blocklist when no scope flag is set.
func blocklistScope(siteIDs, groupIDs, accountIDs []string) mgmt.BlocklistScope {
	if len(siteIDs) == 0 && len(groupIDs) == 0 && len(accountIDs) == 0 {
		return mgmt.BlocklistScope{Tenant: true}
	}
	return mgmt.BlocklistScope{
		SiteIDs:    siteIDs,
		GroupIDs:   groupIDs,
		AccountIDs: accountIDs,
	}
}

func newBlocklistCreateCmd() *cobra.Command {
	var (
		value       string
		sha256Value string
		osType      string
		blockType   string
		description string
		source      string
		siteIDs     []string
		groupIDs    []string
		accountIDs  []string
		yes         bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Add a hash to the blocklist",
		Long: `Add a SHA1 (--value) and/or SHA256 (--sha256) hash to the blocklist.

OS types: windows, linux, macos, windows_legacy
Type must be black_hash (any other value creates an exclusion instead).

New items are added to the scope given by --site-id/--group-id/--account-id, or
to the global (tenant) blocklist when no scope flag is set.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if value == "" {
				return fmt.Errorf("--value is required")
			}
			if osType == "" {
				return fmt.Errorf("--os-type is required")
			}

			data := mgmt.BlocklistCreate{
				Type:        mgmt.BlocklistType(blockType),
				OSType:      mgmt.BlocklistOSType(osType),
				Value:       value,
				SHA256Value: sha256Value,
				Description: description,
				Source:      source,
			}
			scope := blocklistScope(siteIDs, groupIDs, accountIDs)

			action := fmt.Sprintf("add %s hash %q (%s) to the blocklist", blockType, value, osType)
			return guard(cmd.OutOrStdout(), "blocklist create", action, value, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				created, err := c.BlocklistCreate(cmd.Context(), scope, data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), created)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created blocklist item %s (%s)\n", created.ID, created.Value)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&value, "value", "", "SHA1 hash to block (required)")
	cmd.Flags().StringVar(&sha256Value, "sha256", "", "SHA256 hash to block")
	cmd.Flags().StringVar(&osType, "os-type", "", "target OS (windows, linux, macos, windows_legacy) (required)")
	cmd.Flags().StringVar(&blockType, "type", string(mgmt.BlocklistTypeBlackHash), "restriction type (black_hash)")
	cmd.Flags().StringVar(&description, "description", "", "blocklist item description")
	cmd.Flags().StringVar(&source, "source", "", "blocklist item source")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "target site IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "target group IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "target account IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newBlocklistUpdateCmd() *cobra.Command {
	var (
		value       string
		sha256Value string
		osType      string
		blockType   string
		description string
		source      string
		yes         bool
	)

	cmd := &cobra.Command{
		Use:   "update <blocklist-id>",
		Short: "Update a blocklist item (full replacement)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if value == "" || osType == "" {
				return fmt.Errorf("--value and --os-type are required (update replaces the item)")
			}
			data := mgmt.BlocklistCreate{
				Type:        mgmt.BlocklistType(blockType),
				OSType:      mgmt.BlocklistOSType(osType),
				Value:       value,
				SHA256Value: sha256Value,
				Description: description,
				Source:      source,
			}
			return guard(cmd.OutOrStdout(), "blocklist update", "update blocklist item "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				updated, err := c.BlocklistUpdate(cmd.Context(), args[0], mgmt.BlocklistScope{}, data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), updated)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated blocklist item %s\n", args[0])
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&value, "value", "", "SHA1 hash (required)")
	cmd.Flags().StringVar(&sha256Value, "sha256", "", "SHA256 hash")
	cmd.Flags().StringVar(&osType, "os-type", "", "target OS (required)")
	cmd.Flags().StringVar(&blockType, "type", string(mgmt.BlocklistTypeBlackHash), "restriction type (black_hash)")
	cmd.Flags().StringVar(&description, "description", "", "blocklist item description")
	cmd.Flags().StringVar(&source, "source", "", "blocklist item source")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newBlocklistDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <blocklist-id>",
		Short: "Delete a blocklist item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "blocklist delete",
				"delete blocklist item "+args[0], args[0], yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.BlocklistDelete(cmd.Context(), []string{args[0]})
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s\n", pluralize(affected, "blocklist item"))
					return nil
				})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newBlocklistValidateCmd() *cobra.Command {
	var (
		value       string
		sha256Value string
		osType      string
		siteIDs     []string
		groupIDs    []string
		accountIDs  []string
	)

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Check whether a hash is Not Allowed or Not Recommended",
		Long: `Check whether a hash is on SentinelOne's "Not Allowed" or "Not Recommended"
list before adding it to the blocklist. This is a read-only check.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if value == "" && sha256Value == "" {
				return fmt.Errorf("--value or --sha256 is required")
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			res, err := c.BlocklistValidate(cmd.Context(),
				blocklistScope(siteIDs, groupIDs, accountIDs),
				mgmt.BlocklistValidateInput{
					OSType:      mgmt.BlocklistOSType(osType),
					Value:       value,
					SHA256Value: sha256Value,
				})
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), res)
			}
			rows := [][]string{{"Status", string(res.Status)}}
			for _, d := range res.Details {
				rows = append(rows, []string{d.Field, d.Error})
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
	cmd.Flags().StringVar(&value, "value", "", "SHA1 hash to validate")
	cmd.Flags().StringVar(&sha256Value, "sha256", "", "SHA256 hash to validate")
	cmd.Flags().StringVar(&osType, "os-type", "", "target OS (windows, linux, macos, windows_legacy)")
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "scope site IDs")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "scope group IDs")
	cmd.Flags().StringSliceVar(&accountIDs, "account-id", nil, "scope account IDs")
	return cmd
}
