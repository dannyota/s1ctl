package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func addThreatActions(parent *cobra.Command) {
	parent.AddCommand(newThreatMitigateCmd())
	parent.AddCommand(newThreatActionCmd("verdict", "Update analyst verdict on a threat",
		"--verdict", "analyst verdict (true_positive, false_positive, suspicious, undefined)", func(c *mgmt.Client, cmd *cobra.Command, val string, f mgmt.ActionFilter) (int, error) {
			return c.ThreatsUpdateVerdict(cmd.Context(), val, f)
		}))
	parent.AddCommand(newThreatActionCmd("status", "Update incident status on a threat",
		"--status", "incident status (unresolved, in_progress, resolved)", func(c *mgmt.Client, cmd *cobra.Command, val string, f mgmt.ActionFilter) (int, error) {
			return c.ThreatsUpdateStatus(cmd.Context(), val, f)
		}))
	parent.AddCommand(newThreatPlainActionCmd("blacklist", "Add the threat file hash to the blacklist", (*mgmt.Client).ThreatsAddToBlacklist))
	parent.AddCommand(newThreatPlainActionCmd("fetch-file", "Fetch the threat file from the endpoint to the console", (*mgmt.Client).ThreatsFetchFile))
	parent.AddCommand(newThreatAddToExclusionsCmd())
	parent.AddCommand(newThreatMitigateAlertsCmd())
	parent.AddCommand(newThreatSetTicketCmd())
	parent.AddCommand(newThreatQuarantinedFilesCmd())
	parent.AddCommand(newThreatExclusionOptionsCmd())
	parent.AddCommand(newThreatsExportCmd())
}

func newThreatAddToExclusionsCmd() *cobra.Command {
	var scope, exclType, mode, value, description, note, ticketID, pathExclusionType string
	var yes bool

	cmd := &cobra.Command{
		Use:   "add-to-exclusions <threat-id>",
		Short: "Create an exclusion from a threat",
		Long: `Create an exclusion from a threat, overriding the malicious verdict.
Scopes: group, site, account, tenant. Types: hash, path, certificate, browser, file_type.
Mode applies to path exclusions only (e.g. suppress, disable_all_monitors).`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if scope == "" {
				return fmt.Errorf("--scope is required (group, site, account, tenant)")
			}
			if exclType == "" {
				return fmt.Errorf("--type is required (hash, path, certificate, browser, file_type)")
			}
			return guard(cmd.OutOrStdout(), "threats add-to-exclusions",
				fmt.Sprintf("add threat %s to %s exclusions (%s)", args[0], scope, exclType),
				args[0], yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					opts := mgmt.ThreatExclusionOptions{
						TargetScope:       mgmt.ThreatExclusionScope(scope),
						Type:              mgmt.ThreatExclusionType(exclType),
						Mode:              mgmt.ThreatExclusionMode(mode),
						Value:             value,
						Description:       description,
						Note:              note,
						ExternalTicketID:  ticketID,
						PathExclusionType: pathExclusionType,
					}
					affected, err := c.ThreatsAddToExclusions(cmd.Context(), mgmt.ActionFilter{IDs: []string{args[0]}}, opts)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "add-to-exclusions: %s affected\n", pluralize(affected, "threat"))
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&scope, "scope", "", "exclusion scope (group, site, account, tenant)")
	cmd.Flags().StringVar(&exclType, "type", "", "exclusion type (hash, path, certificate, browser, file_type)")
	cmd.Flags().StringVar(&mode, "mode", "", "exclusion mode (path exclusions only, e.g. suppress)")
	cmd.Flags().StringVar(&value, "value", "", "exclusion value (defaults to the threat's value)")
	cmd.Flags().StringVar(&description, "description", "", "exclusion description")
	cmd.Flags().StringVar(&note, "note", "", "note to add to the threat")
	cmd.Flags().StringVar(&ticketID, "ticket-id", "", "external ticket ID to set on the threat")
	cmd.Flags().StringVar(&pathExclusionType, "path-exclusion-type", "", "excluded path type (path exclusions only)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newThreatMitigateAlertsCmd() *cobra.Command {
	var agentID, storyline, action string
	var yes bool

	cmd := &cobra.Command{
		Use:   "mitigate-alerts",
		Short: "Mark an alert as a threat and run a mitigation action",
		Long: `Mark a Deep Visibility alert (identified by agent ID and storyline) as a
threat and run a mitigation action.
Actions: kill, remediate, rollback-remediation, quarantine, un-quarantine, remove_macros, restore_macros.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if agentID == "" {
				return fmt.Errorf("--agent-id is required")
			}
			if storyline == "" {
				return fmt.Errorf("--storyline is required")
			}
			target := agentID + "/" + storyline
			return guard(cmd.OutOrStdout(), "threats mitigate-alerts",
				fmt.Sprintf("mitigate alert %s (%s)", target, orDash(action)),
				target, yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					alerts := []mgmt.ThreatAlert{{AgentID: agentID, Storyline: storyline}}
					affected, err := c.ThreatsMitigateAlerts(cmd.Context(), alerts, mgmt.ThreatMitigationAction(action))
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "mitigate-alerts: %s affected\n", pluralize(affected, "alert"))
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&agentID, "agent-id", "", "agent ID that reported the alert (required)")
	cmd.Flags().StringVar(&storyline, "storyline", "", "storyline of the alert (required)")
	cmd.Flags().StringVar(&action, "action", "", "mitigation action (kill, remediate, quarantine, etc.)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newThreatSetTicketCmd() *cobra.Command {
	var ticketID string
	var yes bool

	cmd := &cobra.Command{
		Use:   "set-ticket <threat-id>",
		Short: "Set the external ticket ID on a threat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if ticketID == "" {
				return fmt.Errorf("--ticket-id is required")
			}
			return guard(cmd.OutOrStdout(), "threats set-ticket",
				fmt.Sprintf("set ticket %s on threat %s", ticketID, args[0]),
				args[0], yes, func() error {
					c, err := mgmtClient()
					if err != nil {
						return err
					}
					affected, err := c.ThreatsSetExternalTicketID(cmd.Context(), mgmt.ActionFilter{IDs: []string{args[0]}}, ticketID)
					if err != nil {
						return err
					}
					if outputFormat == "json" {
						return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
					}
					fmt.Fprintf(cmd.OutOrStdout(), "set-ticket: %s affected\n", pluralize(affected, "threat"))
					return nil
				})
		},
	}
	cmd.Flags().StringVar(&ticketID, "ticket-id", "", "external ticket ID (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newThreatPlainActionCmd(verb, short string, call func(*mgmt.Client, context.Context, mgmt.ActionFilter) (int, error)) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <threat-id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "threats "+verb, verb+" threat "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := call(c, cmd.Context(), mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", verb, pluralize(affected, "threat"))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

type threatActionFn func(*mgmt.Client, *cobra.Command, string, mgmt.ActionFilter) (int, error)

func newThreatActionCmd(verb, short, flagName, flagDesc string, fn threatActionFn) *cobra.Command {
	var val string
	var yes bool

	cmd := &cobra.Command{
		Use:   verb + " <threat-id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if val == "" {
				return fmt.Errorf("%s is required", flagName)
			}
			return guard(cmd.OutOrStdout(), "threats "+verb, fmt.Sprintf("set %s=%s on threat %s", verb, val, args[0]), args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := fn(c, cmd, val, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", verb, pluralize(affected, "threat"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&val, verb, "", flagDesc)
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newThreatMitigateCmd() *cobra.Command {
	var action string
	var yes bool

	cmd := &cobra.Command{
		Use:   "mitigate <threat-id>",
		Short: "Apply mitigation action to a threat",
		Long:  "Actions: kill, quarantine, remediate, rollback-remediation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if action == "" {
				return fmt.Errorf("--action is required (kill, quarantine, remediate, rollback-remediation)")
			}
			return guard(cmd.OutOrStdout(), "threats mitigate", action+" threat "+args[0], args[0], yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.ThreatsMitigate(cmd.Context(), action, mgmt.ActionFilter{IDs: []string{args[0]}})
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s affected\n", action, pluralize(affected, "threat"))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&action, "action", "", "mitigation action (kill, quarantine, remediate, rollback-remediation)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}
