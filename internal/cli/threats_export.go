package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newThreatQuarantinedFilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "quarantined-files <threat-id>",
		Short: "List files quarantined for a threat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			files, err := c.ThreatsQuarantinedFiles(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			headers := []string{"Name", "Path", "Size"}
			rows := make([][]string, len(files))
			for i, f := range files {
				rows[i] = []string{orDash(f.FileName), orDash(f.FilePath), fmt.Sprintf("%d", f.FileSize)}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, files, len(files), len(files), "file", true)
		},
	}
}

func newThreatExclusionOptionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exclusion-options <threat-id>",
		Short: "Show the exclusion (whitening) options available for a threat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			opts, err := c.ThreatsWhiteningOptions(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), opts)
			}
			rows := [][]string{
				{"Exclusion options", orDash(strings.Join(opts.WhiteningOptions, ", "))},
				{"Threat type", orDash(strings.Join(opts.ThreatType, ", "))},
				{"Threat policy", orDash(opts.ThreatPolicy)},
			}
			printTable([]string{"Field", "Value"}, rows)
			return nil
		},
	}
}

func newThreatsExportCmd() *cobra.Command {
	var siteIDs, classifications, statuses, verdicts, mitigationStatuses []string
	var query, outFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export threats to a CSV file",
		Long:  "Export threats matching the filters as CSV. Writes to --out, or stdout when --out is omitted.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			params := &mgmt.ThreatListParams{
				SiteIDs:            siteIDs,
				Classifications:    classifications,
				IncidentStatuses:   statuses,
				AnalystVerdicts:    verdicts,
				MitigationStatuses: mitigationStatuses,
				Query:              query,
			}
			data, err := c.ThreatsExport(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outFile == "" || outFile == "-" {
				_, err = cmd.OutOrStdout().Write(data)
				return err
			}
			if err := os.WriteFile(outFile, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Exported threats to %s\n", outFile)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	cmd.Flags().StringSliceVar(&classifications, "classification", nil, "filter by classification")
	cmd.Flags().StringSliceVar(&statuses, "status", nil, "filter by incident status")
	cmd.Flags().StringSliceVar(&verdicts, "verdict", nil, "filter by analyst verdict")
	cmd.Flags().StringSliceVar(&mitigationStatuses, "mitigation-status", nil, "filter by mitigation status")
	cmd.Flags().StringVar(&query, "query", "", "free text search")
	cmd.Flags().StringVar(&outFile, "out", "", "output file (default: stdout)")
	return cmd
}
