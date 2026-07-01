package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

func newStatusEnumsCmd() *cobra.Command {
	var group string

	cmd := &cobra.Command{
		Use:   "enums",
		Short: "List known enum values used across the CLI",
		Long: `Show all known enum values grouped by domain. Use --group to filter
to a specific domain (e.g. alerts, threats, agents, rules, exclusions, policies).`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStatusEnums(cmd, group)
		},
	}
	cmd.Flags().StringVar(&group, "group", "", "filter to a specific enum group")
	return cmd
}

type enumGroup struct {
	Group  string   `json:"group"`
	Field  string   `json:"field"`
	Values []string `json:"values"`
}

func runStatusEnums(cmd *cobra.Command, filter string) error {
	groups := allEnumGroups()

	if filter != "" {
		upper := strings.ToUpper(filter)
		filtered := groups[:0]
		for _, g := range groups {
			if strings.ToUpper(g.Group) == upper {
				filtered = append(filtered, g)
			}
		}
		groups = filtered
	}

	if outputFormat == "json" {
		return printJSON(cmd.OutOrStdout(), groups)
	}

	headers := []string{"Group", "Field", "Values"}
	rows := make([][]string, len(groups))
	for i, g := range groups {
		rows[i] = []string{g.Group, g.Field, strings.Join(g.Values, ", ")}
	}
	printTable(headers, rows)
	return nil
}

func allEnumGroups() []enumGroup {
	return []enumGroup{
		// Alerts (GraphQL UAM)
		{"alerts", "severity", []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}},
		{"alerts", "status", []string{"NEW", "IN_PROGRESS", "RESOLVED", "UNRESOLVED"}},
		{"alerts", "analystVerdict", []string{
			"UNDEFINED", "TRUE_POSITIVE", "FALSE_POSITIVE", "SUSPICIOUS",
		}},
		{"alerts", "source", []string{"STAR", "EDR", "CWS"}},

		// Threats (REST MGMT)
		{"threats", "classification", []string{
			"Malware", "PUP", "Ransomware", "Trojan", "Worm", "Exploit",
			"Hack Tool", "Downloader", "Backdoor", "Infostealer",
		}},
		{"threats", "incidentStatus", []string{"unresolved", "resolved", "in_progress"}},
		{"threats", "analystVerdict", []string{
			"undefined", "true_positive", "false_positive", "suspicious",
		}},
		{"threats", "mitigationStatus", []string{
			"not_mitigated", "mitigated", "marked_as_benign",
		}},
		{"threats", "confidenceLevel", []string{"malicious", "suspicious"}},
		{"threats", "mitigationAction", []string{
			"kill", "quarantine", "remediate", "rollback-remediation",
		}},

		// Agents (REST MGMT)
		{"agents", "osType", []string{"windows", "linux", "macos", "windows_legacy"}},
		{"agents", "networkStatus", []string{"connected", "connecting", "disconnected", "disconnecting"}},
		{"agents", "machineType", []string{"desktop", "laptop", "server", "kubernetes node", "storage", "unknown"}},
		{"agents", "operationalState", []string{
			"na", "fully_disabled", "partially_disabled",
			"disabled_error", "db_corruption",
		}},
		{"agents", "scanStatus", []string{"none", "started", "aborted", "finished"}},
		{"agents", "mitigationMode", []string{"detect", "protect"}},

		// Rules (REST MGMT / STAR)
		{"rules", "status", []string{"Draft", "Activating", "Active", "Disabling", "Disabled", "Deleted", "Deleting"}},
		{"rules", "severity", []string{"Info", "Low", "Medium", "High", "Critical"}},
		{"rules", "scope", []string{"global", "account", "site", "group"}},
		{"rules", "queryType", []string{"events", "correlation", "uebafirstseen", "scheduled"}},
		{"rules", "expirationMode", []string{"Permanent", "Temporary"}},
		{"rules", "treatAsThreat", []string{"UNDEFINED", "Suspicious", "Malicious"}},

		// Exclusions (REST MGMT)
		{"exclusions", "type", []string{
			"path", "file_type", "white_hash", "browser",
			"certificate", "document_type",
		}},
		{"exclusions", "osType", []string{"windows", "linux", "macos", "windows_legacy"}},
		{"exclusions", "mode", []string{"suppress", "suppress_dynamic_only", "suppress_app_control"}},
		{"exclusions", "pathExclusionType", []string{"subfolders", "file", "glob"}},

		// Policies (REST MGMT)
		{"policies", "mitigationMode", []string{"detect", "protect"}},
		{"policies", "mitigationModeSuspicious", []string{"detect", "protect"}},

		// Misconfigurations (GraphQL xSPM)
		{"misconfigurations", "severity", []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}},
		{"misconfigurations", "status", []string{"NEW", "IN_PROGRESS", "RESOLVED"}},
		{"misconfigurations", "analystVerdict", []string{
			"UNDEFINED", "TRUE_POSITIVE", "FALSE_POSITIVE", "SUSPICIOUS",
		}},

		// Vulnerabilities (GraphQL xSPM)
		{"vulnerabilities", "severity", []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}},
		{"vulnerabilities", "status", []string{"NEW", "IN_PROGRESS", "RESOLVED"}},

		// Cloud Policies (GraphQL)
		{"cloud-policies", "severity", []string{"LOW", "MEDIUM", "HIGH", "CRITICAL"}},
		{"cloud-policies", "status", []string{"ENABLED", "DISABLED"}},
	}
}
