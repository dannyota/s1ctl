package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Build metadata. A release sets these via the linker:
//
//	go build -ldflags "-X danny.vn/s1/internal/cli.version=v1.2.3 \
//	  -X danny.vn/s1/internal/cli.commit=$(git rev-parse HEAD) \
//	  -X danny.vn/s1/internal/cli.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
var (
	version = ""
	commit  = ""
	date    = ""
)

// BuildInfo is the resolved version metadata.
type BuildInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date,omitempty"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// resolveBuildInfo merges ldflags-stamped values with the toolchain-embedded VCS
// info, preferring the explicit ldflags values.
func resolveBuildInfo() BuildInfo {
	bi := BuildInfo{
		Version: version, Commit: commit, Date: date,
		GoVersion: runtime.Version(), OS: runtime.GOOS, Arch: runtime.GOARCH,
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if bi.Version == "" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			bi.Version = info.Main.Version
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if bi.Commit == "" {
					bi.Commit = s.Value
				}
			case "vcs.time":
				if bi.Date == "" {
					bi.Date = s.Value
				}
			}
		}
	}
	if bi.Version == "" {
		bi.Version = "dev"
	}
	if bi.Commit == "" {
		bi.Commit = "unknown"
	}
	return bi
}

func shortCommit(c string) string {
	if len(c) > 12 {
		return c[:12]
	}
	return c
}

// versionLine returns a compact one-line version string for doctor/help.
func versionLine() string {
	bi := resolveBuildInfo()
	if bi.Commit == "unknown" {
		return fmt.Sprintf("s1ctl %s (%s %s/%s)",
			bi.Version, bi.GoVersion, bi.OS, bi.Arch)
	}
	return fmt.Sprintf("s1ctl %s (%s, %s %s/%s)",
		bi.Version, shortCommit(bi.Commit), bi.GoVersion, bi.OS, bi.Arch)
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version, commit, and build info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			bi := resolveBuildInfo()
			if outputFormat == "json" {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(bi)
			}
			fmt.Printf("s1ctl %s\n", bi.Version)
			if bi.Commit != "unknown" {
				fmt.Printf("  commit:   %s\n", bi.Commit)
			}
			if bi.Date != "" {
				fmt.Printf("  built:    %s\n", bi.Date)
			}
			fmt.Printf("  go:       %s\n", bi.GoVersion)
			fmt.Printf("  platform: %s/%s\n", bi.OS, bi.Arch)
			return nil
		},
	}
	return markJSON(cmd)
}
