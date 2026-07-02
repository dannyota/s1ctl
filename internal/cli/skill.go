package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	skill "danny.vn/s1/skills/s1ctl"
)

func newSkillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Print the agent operating guide embedded in the binary",
		Long: `Emit the s1ctl agent operating guide baked into the binary.

go install danny.vn/s1/cmd/s1ctl@latest ships only the binary, so this is
how an agent retrieves the guide without the repo:
  s1ctl skill            print the guide (--json wraps {name,description,body})
  s1ctl skill install    write it into an agent skills dir for auto-detection`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), skill.Parse())
			}
			fmt.Fprint(cmd.OutOrStdout(), skill.Markdown())
			return nil
		},
	}
	cmd.AddCommand(newSkillInstallCmd())
	return cmd
}

func newSkillInstallCmd() *cobra.Command {
	var dir string
	var force bool

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the operating guide into an agent skills directory",
		Long: `Write the embedded guide to <skills-dir>/s1ctl/SKILL.md so an agent
harness detects it as a first-class skill. Default directory is
~/.claude/skills; --dir overrides it. Idempotent: no-op when the file
already matches; if it exists with different content, left untouched
unless --force.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			base := dir
			if base == "" {
				base = defaultSkillsDir()
			}
			dest := filepath.Join(base, "s1ctl", "SKILL.md")
			doc := skill.Parse()
			name := doc.Name
			if name == "" {
				name = "s1ctl"
			}

			status := "installed"
			if existing, err := os.ReadFile(dest); err == nil {
				switch {
				case string(existing) == skill.Markdown():
					status = "unchanged"
				case !force:
					return fmt.Errorf("%s already exists with different content; pass --force to overwrite", dest)
				default:
					status = "updated"
				}
			}

			if status != "unchanged" {
				if err := os.MkdirAll(filepath.Dir(dest), 0o750); err != nil {
					return fmt.Errorf("create skills dir: %w", err)
				}
				if err := os.WriteFile(dest, []byte(skill.Markdown()), 0o644); err != nil {
					return fmt.Errorf("write skill: %w", err)
				}
			}

			if outputFormat == "json" {
				return printJSON(cmd.OutOrStdout(), map[string]string{"name": name, "path": dest, "status": status})
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Skill %q %s -> %s\n", name, status, dest)
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "", "skills directory (default ~/.claude/skills)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing SKILL.md with different content")
	return cmd
}

func defaultSkillsDir() string {
	if cfg := os.Getenv("CLAUDE_CONFIG_DIR"); cfg != "" {
		return filepath.Join(cfg, "skills")
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return filepath.Join(".claude", "skills")
	}
	return filepath.Join(home, ".claude", "skills")
}
