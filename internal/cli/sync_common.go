package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/internal/reconcile"
)

// scopeFlags carries the flag values shared across every sync surface. Not
// every surface binds every field: pull binds OutDir (+ optional SiteIDs/
// AccountIDs filters), push binds Dir/Yes (+ optional SiteIDs create scope).
// push is set by the push command so a surface's List closure can pick the
// legacy list scope (some surfaces list unfiltered on push, filtered on pull).
type scopeFlags struct {
	OutDir     string
	Dir        string
	SiteIDs    []string
	AccountIDs []string
	Yes        bool
	push       bool
}

// surfaceSpec adapts one config-as-code resource to the shared pull/push
// command builders over the reconcile engine. Build constructs the engine
// Surface, creating any API client(s) lazily inside the returned closures so a
// GraphQL-backed surface fits the same shape (and so an empty/missing local
// directory never needs credentials).
type surfaceSpec struct {
	Noun       string // singular resource noun, e.g. "device rule"
	Command    string // CLI group name for guard/audit strings, e.g. "devicecontrol"
	DefaultDir string // default pull --out / push --dir directory

	PullShort string
	PullLong  string
	PushShort string
	PushLong  string

	// RegisterPullFlags / RegisterPushFlags bind surface-specific flags beyond
	// the common --out (pull) and --dir/--yes (push).
	RegisterPullFlags func(cmd *cobra.Command, scope *scopeFlags)
	RegisterPushFlags func(cmd *cobra.Command, scope *scopeFlags)

	Build func(cmd *cobra.Command, scope scopeFlags) (reconcile.Surface, error)
}

// applySummary is the machine-readable push-apply outcome.
type applySummary struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Failed  int `json:"failed"`
}

// syncSurfaceSpecs is the registry of every engine-backed sync surface. The
// drift command iterates it to plan each surface against its default directory.
func syncSurfaceSpecs() []surfaceSpec {
	return []surfaceSpec{
		blocklistSpec(),
		cloudPoliciesSpec(),
		deviceControlSpec(),
		exclusionsSpec(),
		firewallSpec(),
		networkQuarantineSpec(),
		locationsSpec(),
		groupsSpec(),
		rulesSpec(),
		sitesSpec(),
		tagsSpec(),
	}
}

// newEnginePullCmd renders live objects to per-object YAML files.
func newEnginePullCmd(spec surfaceSpec) *cobra.Command {
	var scope scopeFlags

	cmd := &cobra.Command{
		Use:   "pull",
		Short: spec.PullShort,
		Long:  spec.PullLong,
		RunE: func(cmd *cobra.Command, _ []string) error {
			surface, err := spec.Build(cmd, scope)
			if err != nil {
				return err
			}

			live, err := surface.List(cmd.Context())
			if err != nil {
				return err
			}

			stale, err := reconcile.WriteDir(scope.OutDir, live)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Pulled %s to %s\n",
				pluralize(len(live), spec.Noun), scope.OutDir)
			for _, name := range stale {
				fmt.Fprintf(cmd.ErrOrStderr(),
					"warning: %s has no live %s (delete it or push will re-create the object)\n",
					name, spec.Noun)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&scope.OutDir, "out", spec.DefaultDir, "output directory")
	if spec.RegisterPullFlags != nil {
		spec.RegisterPullFlags(cmd, &scope)
	}
	return cmd
}

// newEnginePushCmd plans local files against live state and applies the diff.
// The client and live List run before the guard so dry-run prints the plan
// (matching the legacy per-object push flow).
func newEnginePushCmd(spec surfaceSpec) *cobra.Command {
	var scope scopeFlags

	cmd := &cobra.Command{
		Use:   "push",
		Short: spec.PushShort,
		Long:  spec.PushLong,
		RunE: func(cmd *cobra.Command, _ []string) error {
			scope.push = true
			surface, err := spec.Build(cmd, scope)
			if err != nil {
				return err
			}

			// A missing directory is a hard error naming it (legacy behavior);
			// an empty-but-present directory is a clean no-op.
			if info, sErr := os.Stat(scope.Dir); sErr != nil {
				return fmt.Errorf("read %s: %w", scope.Dir, sErr)
			} else if !info.IsDir() {
				return fmt.Errorf("read %s: not a directory", scope.Dir)
			}

			local, err := reconcile.LoadDir(scope.Dir, surface.Decode)
			if err != nil {
				return err
			}
			if len(local) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No %s files found.\n", spec.Noun)
				return nil
			}

			live, err := surface.List(cmd.Context())
			if err != nil {
				return err
			}

			plan, warnings, err := reconcile.BuildPlan(local, live)
			if err != nil {
				return err
			}
			for _, w := range warnings {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning:", w)
			}

			action := reconcile.Describe(plan, spec.Noun, scope.Dir)
			return guard(cmd.OutOrStdout(), spec.Command+" push", action, scope.Dir, scope.Yes, func() error {
				sum, aErr := reconcile.Apply(cmd.Context(), surface, plan, cmd.ErrOrStderr())
				if outputFormat == "json" {
					if pErr := printJSON(cmd.OutOrStdout(), applySummary{sum.Created, sum.Updated, sum.Failed}); pErr != nil {
						return pErr
					}
					return aErr
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created %s, updated %s\n",
					pluralize(sum.Created, spec.Noun), pluralize(sum.Updated, spec.Noun))
				if sum.Failed > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "%s failed\n", pluralize(sum.Failed, spec.Noun))
				}
				if n := len(plan.LiveOnly()); n > 0 {
					fmt.Fprintf(cmd.OutOrStdout(),
						"%s live-only in the console, left unchanged\n", pluralize(n, spec.Noun))
				}
				return aErr
			})
		},
	}
	cmd.Flags().StringVar(&scope.Dir, "dir", spec.DefaultDir, "directory containing "+spec.Noun+" YAML files")
	cmd.Flags().BoolVar(&scope.Yes, "yes", false, "apply changes (default: dry-run)")
	if spec.RegisterPushFlags != nil {
		spec.RegisterPushFlags(cmd, &scope)
	}
	return cmd
}
