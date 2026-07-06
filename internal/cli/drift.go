package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/internal/reconcile"
)

// driftResult is one surface's plan summarized for the drift report. Name
// slices are always non-nil so --json renders `[]` rather than `null`.
type driftResult struct {
	Surface   string   `json:"surface"`
	Creates   []string `json:"creates"`
	Updates   []string `json:"updates"`
	LiveOnly  []string `json:"liveOnly"`
	Unchanged int      `json:"unchanged"`
}

// drifted reports whether this surface differs from live state: any create,
// update, or live-only object is drift.
func (r driftResult) drifted() bool {
	return len(r.Creates) > 0 || len(r.Updates) > 0 || len(r.LiveOnly) > 0
}

func newDriftCmd() *cobra.Command {
	var surfaces []string
	var dirRoot string

	cmd := &cobra.Command{
		Use:   "drift",
		Short: "Report drift between committed config and live state",
		Long: `Compare committed config-as-code against the live console for every sync
surface and report the difference without applying anything.

For each surface with a local directory under --dir-root, drift loads the
committed files, lists the live objects, and plans the reconcile: creates
(committed, not live), updates (committed, differs from live), live-only
(live, not committed) and unchanged. Surfaces without a local directory are
skipped — drift checks only what is committed.

The command is read-only: it lists, plans, and reports, and has no apply path.
Exit code is 0 when every checked surface is clean and 1 when any surface has
drift, so a CI job can fail on a non-zero exit.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			specs, err := selectDriftSpecs(surfaces)
			if err != nil {
				return err
			}
			return runDrift(cmd, specs, dirRoot)
		},
	}
	cmd.Flags().StringSliceVar(&surfaces, "surface", nil, "limit to named surfaces (repeatable; default: all)")
	cmd.Flags().StringVar(&dirRoot, "dir-root", ".", "root directory containing per-surface config directories")
	return markJSON(cmd)
}

// selectDriftSpecs returns every registered surface spec, or just the named
// ones. An unknown name is a hard error listing the valid surfaces.
func selectDriftSpecs(names []string) ([]surfaceSpec, error) {
	all := syncSurfaceSpecs()
	if len(names) == 0 {
		return all, nil
	}

	byName := make(map[string]surfaceSpec, len(all))
	valid := make([]string, 0, len(all))
	for _, s := range all {
		byName[s.Command] = s
		valid = append(valid, s.Command)
	}
	sort.Strings(valid)

	out := make([]surfaceSpec, 0, len(names))
	for _, n := range names {
		s, ok := byName[n]
		if !ok {
			return nil, fmt.Errorf("unknown surface %q (valid: %s)", n, strings.Join(valid, ", "))
		}
		out = append(out, s)
	}
	return out, nil
}

// runDrift plans each surface whose local directory exists and reports the
// summary. No API client is constructed until a surface with an existing local
// directory is reached (the surface List closures create clients lazily), so an
// all-directories-missing run stays fully offline.
func runDrift(cmd *cobra.Command, specs []surfaceSpec, dirRoot string) error {
	var results []driftResult
	for _, spec := range specs {
		dir := filepath.Join(dirRoot, spec.DefaultDir)
		info, sErr := os.Stat(dir)
		if sErr != nil || !info.IsDir() {
			continue // drift checks only committed config; skip missing dirs
		}

		result, err := driftSurface(cmd, spec, dir)
		if err != nil {
			return err
		}
		results = append(results, result)
	}

	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no local surface directories found")
		return nil
	}

	return reportDrift(cmd, results)
}

// driftSurface builds the plan for one surface: load committed files, list live
// objects, and classify. BuildPlan warnings go to stderr.
func driftSurface(cmd *cobra.Command, spec surfaceSpec, dir string) (driftResult, error) {
	// push scope lists tenant-wide (no per-directory filters) so the plan
	// compares committed config against the full live set.
	surface, err := spec.Build(cmd, scopeFlags{push: true})
	if err != nil {
		return driftResult{}, err
	}

	local, err := reconcile.LoadDir(dir, surface.Decode)
	if err != nil {
		return driftResult{}, err
	}

	live, err := surface.List(cmd.Context())
	if err != nil {
		return driftResult{}, err
	}

	plan, warnings, err := reconcile.BuildPlan(local, live)
	if err != nil {
		return driftResult{}, err
	}
	for _, w := range warnings {
		fmt.Fprintln(cmd.ErrOrStderr(), "warning:", w)
	}

	return driftResult{
		Surface:   spec.Command,
		Creates:   itemNames(plan.Creates()),
		Updates:   itemNames(plan.Updates()),
		LiveOnly:  itemNames(plan.LiveOnly()),
		Unchanged: len(plan.Unchanged()),
	}, nil
}

// reportDrift prints the per-surface summary and returns a drift error (exit 1)
// when any checked surface differs from live state.
func reportDrift(cmd *cobra.Command, results []driftResult) error {
	headers := []string{"SURFACE", "CREATE", "UPDATE", "LIVE-ONLY", "UNCHANGED"}
	rows := make([][]string, 0, len(results))
	for _, r := range results {
		rows = append(rows, []string{
			r.Surface,
			strconv.Itoa(len(r.Creates)),
			strconv.Itoa(len(r.Updates)),
			strconv.Itoa(len(r.LiveOnly)),
			strconv.Itoa(r.Unchanged),
		})
	}

	switch outputFormat {
	case "json":
		if err := printJSON(cmd.OutOrStdout(), results); err != nil {
			return err
		}
	case "csv":
		if err := printCSV(cmd.OutOrStdout(), headers, rows); err != nil {
			return err
		}
	default:
		printTable(headers, rows)
	}

	drifted := 0
	for _, r := range results {
		if r.drifted() {
			drifted++
		}
	}
	if drifted > 0 {
		return fmt.Errorf("drift detected in %s", pluralize(drifted, "surface"))
	}
	return nil
}

// itemNames extracts the object names from a slice of plan items.
func itemNames(items []reconcile.Item) []string {
	names := make([]string, 0, len(items))
	for _, it := range items {
		names = append(names, it.Name)
	}
	return names
}
