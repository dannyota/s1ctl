package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newNetworkMoveCmd() *cobra.Command {
	var targetSiteID, targetAccountID, targetGroupID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "move <rule-id>...",
		Short: "Move network quarantine rules to another scope",
		Long: `Move one or more network quarantine rules to a target scope.

Use --target-site-id, --target-account-id, or --target-group-id for the
destination. At least one target flag is required.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetDesc, err := describeCopyTarget(targetSiteID, targetAccountID, targetGroupID)
			if err != nil {
				return err
			}
			return guard(cmd.OutOrStdout(), "network move",
				"move "+pluralize(len(args), "network quarantine rule")+" to "+targetDesc,
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					target := buildCopyTarget(targetSiteID, targetAccountID, targetGroupID)
					affected, err := c.NetworkQuarantineMoveRules(cmd.Context(),
						mgmt.FirewallActionFilter{IDs: args},
						[]mgmt.FirewallRuleCopyTarget{target})
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Moved %s\n", pluralize(affected, "network quarantine rule"))
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&targetSiteID, "target-site-id", "", "target site ID")
	cmd.Flags().StringVar(&targetAccountID, "target-account-id", "", "target account ID")
	cmd.Flags().StringVar(&targetGroupID, "target-group-id", "", "target group ID")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

func newNetworkSetLocationCmd() *cobra.Command {
	var locType string
	var locationIDs []string
	var yes bool

	cmd := &cobra.Command{
		Use:   "set-location <rule-id>...",
		Short: "Set the location assignment of network quarantine rules",
		Long: `Assign a location matcher to one or more network quarantine rules.

--type is one of all, specific, or fallback. For "specific", pass one or more
--location-id values.
Dry-run by default — pass --yes to apply.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			loc := mgmt.FirewallLocationTarget{Type: mgmt.FirewallLocationType(locType)}
			switch loc.Type {
			case mgmt.FirewallLocationAll, mgmt.FirewallLocationSpecific, mgmt.FirewallLocationFallback:
			default:
				return fmt.Errorf("invalid --type %q: expected all, specific, or fallback", locType)
			}
			if loc.Type == mgmt.FirewallLocationSpecific && len(locationIDs) == 0 {
				return fmt.Errorf("--location-id is required when --type is specific")
			}
			for _, id := range locationIDs {
				loc.Values = append(loc.Values, mgmt.FirewallLocationValue{ID: id})
			}
			return guard(cmd.OutOrStdout(), "network set-location",
				"set location "+locType+" on "+pluralize(len(args), "network quarantine rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.NetworkQuarantineSetLocation(cmd.Context(),
						mgmt.FirewallActionFilter{IDs: args}, loc)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Updated location on %s\n", pluralize(affected, "network quarantine rule"))
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&locType, "type", "all", "location type: all, specific, or fallback")
	cmd.Flags().StringSliceVar(&locationIDs, "location-id", nil, "location IDs (for --type specific)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}

func newNetworkTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Add or remove tags on network quarantine rules",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newNetworkTagsChangeCmd("add"))
	cmd.AddCommand(newNetworkTagsChangeCmd("remove"))
	return cmd
}

// newNetworkTagsChangeCmd builds the tags add/remove subcommands.
func newNetworkTagsChangeCmd(verb string) *cobra.Command {
	var tagIDs []string
	var yes bool

	preposition := "to"
	if verb == "remove" {
		preposition = "from"
	}

	cmd := &cobra.Command{
		Use:   verb + " <rule-id>...",
		Short: strings.ToUpper(verb[:1]) + verb[1:] + " tags " + preposition + " network quarantine rules",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(tagIDs) == 0 {
				return fmt.Errorf("--tag-id is required")
			}
			return guard(cmd.OutOrStdout(), "network tags "+verb,
				verb+" "+pluralize(len(tagIDs), "tag")+" "+preposition+" "+pluralize(len(args), "network quarantine rule"),
				strings.Join(args, ","), yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					filter := mgmt.FirewallActionFilter{IDs: args}
					var affected int
					if verb == "add" {
						affected, err = c.NetworkQuarantineAddTags(cmd.Context(), filter, tagIDs)
					} else {
						affected, err = c.NetworkQuarantineRemoveTags(cmd.Context(), filter, tagIDs)
					}
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Updated tags on %s\n", pluralize(affected, "network quarantine rule"))
					return nil
				})
		},
	}
	cmd.Flags().StringSliceVar(&tagIDs, "tag-id", nil, "tag IDs")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply changes (default: dry-run)")
	return cmd
}
