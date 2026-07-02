package cli

import "github.com/spf13/cobra"

func registerCommands(root *cobra.Command) {
	// Foundation
	root.AddCommand(newVersionCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newDoctorCmd())
	root.AddCommand(newConfigCmd())
	root.AddCommand(newCommandsCmd())
	root.AddCommand(newCompletionCmd())

	// Endpoint security
	root.AddCommand(newAgentsCmd())
	root.AddCommand(newThreatsCmd())
	root.AddCommand(newAlertsCmd())
	root.AddCommand(newMisconfigurationsCmd())
	root.AddCommand(newVulnerabilitiesCmd())
	root.AddCommand(newSitesCmd())
	root.AddCommand(newGroupsCmd())
	root.AddCommand(newAccountsCmd())
	root.AddCommand(newPoliciesCmd())
	root.AddCommand(newExclusionsCmd())
	root.AddCommand(newUnifiedExclusionsCmd())
	root.AddCommand(newBlocklistCmd())
	root.AddCommand(newActivitiesCmd())
	root.AddCommand(newCloudPoliciesCmd())
	root.AddCommand(newCloudRulesCmd())
	root.AddCommand(newDLPCmd())
	root.AddCommand(newRulesCmd())

	// Operations
	root.AddCommand(newUsersCmd())
	root.AddCommand(newServiceUsersCmd())
	root.AddCommand(newRolesCmd())
	root.AddCommand(newTagsCmd())
	root.AddCommand(newTagRulesCmd())
	root.AddCommand(newRemoteOpsCmd())
	root.AddCommand(newApplicationsCmd())
	root.AddCommand(newDeviceControlCmd())
	root.AddCommand(newFirewallCmd())
	root.AddCommand(newNetworkCmd())
	root.AddCommand(newLocationsCmd())
	root.AddCommand(newFiltersCmd())
	root.AddCommand(newUpdatesCmd())
	root.AddCommand(newReportsCmd())
	root.AddCommand(newSettingsCmd())
	root.AddCommand(newMaintenanceCmd())
	root.AddCommand(newUpgradePoliciesCmd())

	// Threat intelligence
	root.AddCommand(newIOCsCmd())
	root.AddCommand(newDetectionLibraryCmd())
	root.AddCommand(newRangerADCmd())

	// Asset inventory
	root.AddCommand(newAssetsCmd())

	// System
	root.AddCommand(newSystemCmd())

	// Hunting & data lake
	root.AddCommand(newVisibilityCmd())
	root.AddCommand(newDatalakeCmd())

	// Tooling
	root.AddCommand(newDriftCmd())
	root.AddCommand(newDocsCmd())
	root.AddCommand(newSkillCmd())
	root.AddCommand(newAuditCmd())
}
