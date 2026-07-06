package cli

import (
	"strings"

	"github.com/spf13/cobra"
)

func newStatusSurfacesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "surfaces",
		Short: "List all API surfaces and supported operations",
		Args:  cobra.NoArgs,
		RunE:  runStatusSurfaces,
	}
	return markJSON(cmd)
}

type surfaceEntry struct {
	Name       string   `json:"name"`
	API        string   `json:"api"`
	Operations []string `json:"operations"`
}

func runStatusSurfaces(cmd *cobra.Command, _ []string) error {
	surfaces := []surfaceEntry{
		{"agents", "REST", []string{"list", "get", "count", "health", "upgrade", "outdated", "versions", "actions"}},
		{"threats", "REST", []string{"list", "get", "count", "resolve", "notes", "add-note", "actions"}},
		{"alerts", "GraphQL", []string{"list", "get", "count", "resolve", "status", "verdict", "add-note"}},
		{"misconfigurations", "GraphQL", []string{"list", "get", "status", "verdict"}},
		{"vulnerabilities", "GraphQL", []string{"list", "get", "health"}},
		{"sites", "REST", []string{"list", "get", "count"}},
		{"groups", "REST", []string{"list", "get", "count", "create", "delete"}},
		{"accounts", "REST", []string{"list", "get", "count"}},
		{"policies", "REST", []string{"list", "get", "diff", "pull", "push"}},
		{"exclusions", "REST", []string{"list", "get", "create", "pull", "push"}},
		{"rules", "REST", []string{"list", "get", "health", "trends", "detections", "diff", "validate", "enable", "disable", "pull", "push"}},
		{"cloud-policies", "GraphQL", []string{"list", "get"}},
		{"activities", "REST", []string{"list", "count"}},
		{"users", "REST", []string{"list", "get"}},
		{"tags", "REST", []string{"list"}},
		{"remote-ops", "REST", []string{"list", "run", "results"}},
		{"applications", "REST", []string{"list"}},
		{"device-control", "REST", []string{"list", "pull", "push"}},
		{"firewall", "REST", []string{"list", "pull", "push"}},
		{"updates", "REST", []string{"list"}},
		{"visibility", "REST", []string{"query"}},
		{"datalake", "SDL", []string{"query", "saved"}},
	}

	headers := []string{"Surface", "API", "Operations"}
	rows := make([][]string, len(surfaces))
	for i, s := range surfaces {
		rows[i] = []string{s.Name, s.API, strings.Join(s.Operations, ", ")}
	}
	return printOutput(cmd.OutOrStdout(), headers, rows, surfaces, len(surfaces), len(surfaces), "surface", true)
}
