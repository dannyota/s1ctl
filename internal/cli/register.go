package cli

import "github.com/spf13/cobra"

func registerCommands(root *cobra.Command) {
	// Foundation
	root.AddCommand(newVersionCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newConfigCmd())
	root.AddCommand(newCommandsCmd())
	root.AddCommand(newCompletionCmd())

	// Endpoint security
	root.AddCommand(newAgentsCmd())
	root.AddCommand(newThreatsCmd())
	root.AddCommand(newAlertsCmd())
	root.AddCommand(newSitesCmd())
	root.AddCommand(newGroupsCmd())
	root.AddCommand(newAccountsCmd())
	root.AddCommand(newPoliciesCmd())
	root.AddCommand(newExclusionsCmd())
	root.AddCommand(newActivitiesCmd())

	// Operations
	root.AddCommand(newUsersCmd())
	root.AddCommand(newTagsCmd())
	root.AddCommand(newRemoteOpsCmd())
	root.AddCommand(newApplicationsCmd())
	root.AddCommand(newDeviceControlCmd())
	root.AddCommand(newFirewallCmd())
	root.AddCommand(newUpdatesCmd())

	// Data lake
	root.AddCommand(newDatalakeCmd())
}
