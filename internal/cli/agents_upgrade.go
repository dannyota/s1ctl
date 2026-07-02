package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newAgentsUpgradeCmd() *cobra.Command {
	var siteIDs, groupIDs []string
	var query string
	var packageID, fileName, path, osType, packageType string
	var allowDowngrade, ignoreConflicts, scheduled bool
	var yes bool

	cmd := &cobra.Command{
		Use:   "upgrade [agent-id...]",
		Short: "Trigger agent software upgrade",
		Long: `Trigger a software update on one or more agents.

Exactly one of --package-id, --file-name, or --path is required to
identify the upgrade package. The --file-name option also requires
--os-type.

Specify agent IDs as arguments, or use --site-id / --group-id / --query
to target agents by filter. Dry-run by default.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := mgmt.ActionFilter{
				IDs:     args,
				SiteIDs: siteIDs,
				Query:   query,
			}
			if len(filter.IDs) == 0 && len(filter.SiteIDs) == 0 && filter.Query == "" {
				return fmt.Errorf("specify agent IDs or --site-id / --query")
			}
			set := 0
			for _, s := range []string{packageID, fileName, path} {
				if s != "" {
					set++
				}
			}
			if set == 0 {
				return fmt.Errorf("one of --package-id, --file-name, or --path is required")
			}
			if set > 1 {
				return fmt.Errorf("only one of --package-id, --file-name, or --path may be specified")
			}
			data := mgmt.UpdateSoftwareData{
				PackageID:   packageID,
				FileName:    fileName,
				Path:        path,
				OSType:      osType,
				PackageType: packageType,
			}
			if cmd.Flags().Changed("allow-downgrade") {
				data.AllowDowngrade = &allowDowngrade
			}
			if cmd.Flags().Changed("ignore-conflicts") {
				data.IgnoreConflicts = &ignoreConflicts
			}
			if cmd.Flags().Changed("scheduled") {
				data.IsScheduled = &scheduled
			}
			_ = groupIDs
			return guard(cmd.OutOrStdout(), "agents upgrade", "trigger upgrade on "+describeFilter(filter), describeFilter(filter), yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				affected, err := c.AgentsUpdateSoftware(cmd.Context(), filter, data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]int{"affected": affected})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "upgrade: %s affected\n", pluralize(affected, "agent"))
				return nil
			})
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&groupIDs, "group-id", nil, "filter by group ID")
	cmd.Flags().StringVar(&query, "query", "", "free text search filter")
	cmd.Flags().StringVar(&packageID, "package-id", "", "upgrade package ID")
	cmd.Flags().StringVar(&fileName, "file-name", "", "upgrade package file name")
	cmd.Flags().StringVar(&path, "path", "", "local path to upgrade package on the endpoint")
	cmd.Flags().StringVar(&osType, "os-type", "", "target OS type (linux, macos, windows)")
	cmd.Flags().StringVar(&packageType, "package-type", "", "package type (Agent, Ranger, AgentAndRanger)")
	cmd.Flags().BoolVar(&allowDowngrade, "allow-downgrade", false, "allow downgrading the agent version")
	cmd.Flags().BoolVar(&ignoreConflicts, "ignore-conflicts", false, "ignore conflicts with active upgrade policies")
	cmd.Flags().BoolVar(&scheduled, "scheduled", false, "upgrade according to agent upgrade schedule")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newAgentsOutdatedCmd() *cobra.Command {
	var siteIDs []string
	var limit int
	var all bool

	cmd := &cobra.Command{
		Use:   "outdated",
		Short: "List agents not on the latest version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			isUpToDate := false
			params := &mgmt.AgentListParams{
				SiteIDs:    siteIDs,
				IsUpToDate: &isUpToDate,
				Limit:      limit,
			}
			if params.Limit == 0 {
				params.Limit = defaultPageSize
			}

			var agents []mgmt.Agent
			var total int

			if all {
				agents, total, err = fetchAllREST("agent", func(cur string) ([]mgmt.Agent, *mgmt.Pagination, error) {
					params.Cursor = cur
					return c.AgentsList(cmd.Context(), params)
				})
			} else {
				var pag *mgmt.Pagination
				agents, pag, err = c.AgentsList(cmd.Context(), params)
				if pag != nil {
					total = pag.TotalItems
				}
			}
			if err != nil {
				return err
			}

			headers := []string{"ID", "Name", "Version", "OS", "Site", "Last Active"}
			rows := make([][]string, len(agents))
			for i, a := range agents {
				rows[i] = []string{
					a.ID, a.ComputerName, a.AgentVersion,
					a.OSType, a.SiteName, orDash(a.LastActiveDate),
				}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, agents, len(agents), total, "agent", all)
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "max results per page (default 50)")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all pages")
	return cmd
}

func newAgentsVersionsCmd() *cobra.Command {
	var siteIDs []string

	cmd := &cobra.Command{
		Use:   "versions",
		Short: "Show agent version distribution",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.AgentListParams{
				SiteIDs: siteIDs,
				Limit:   1000,
			}
			agents, _, err := fetchAllREST("agent", func(cur string) ([]mgmt.Agent, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.AgentsList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			type versionCount struct {
				Version string `json:"version"`
				Linux   int    `json:"linux"`
				Windows int    `json:"windows"`
				MacOS   int    `json:"macos"`
				Total   int    `json:"total"`
			}
			counts := make(map[string]*versionCount)
			for _, a := range agents {
				vc, ok := counts[a.AgentVersion]
				if !ok {
					vc = &versionCount{Version: a.AgentVersion}
					counts[a.AgentVersion] = vc
				}
				vc.Total++
				switch strings.ToLower(a.OSType) {
				case "linux":
					vc.Linux++
				case "windows":
					vc.Windows++
				case "macos", "osx":
					vc.MacOS++
				}
			}
			var versions []*versionCount
			for _, vc := range counts {
				versions = append(versions, vc)
			}
			sort.Slice(versions, func(i, j int) bool {
				return versions[i].Total > versions[j].Total
			})

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), versions)
			}

			headers := []string{"Version", "Linux", "Windows", "macOS", "Total"}
			rows := make([][]string, len(versions))
			for i, v := range versions {
				rows[i] = []string{
					v.Version,
					fmt.Sprintf("%d", v.Linux),
					fmt.Sprintf("%d", v.Windows),
					fmt.Sprintf("%d", v.MacOS),
					fmt.Sprintf("%d", v.Total),
				}
			}
			printTable(headers, rows)
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s across %s\n",
				pluralize(len(versions), "version"), pluralize(len(agents), "agent"))
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}

func describeFilter(f mgmt.ActionFilter) string {
	var parts []string
	if len(f.IDs) > 0 {
		parts = append(parts, pluralize(len(f.IDs), "agent"))
	}
	if len(f.SiteIDs) > 0 {
		parts = append(parts, fmt.Sprintf("site %s", strings.Join(f.SiteIDs, ",")))
	}
	if f.Query != "" {
		parts = append(parts, fmt.Sprintf("query %q", f.Query))
	}
	return strings.Join(parts, ", ")
}
