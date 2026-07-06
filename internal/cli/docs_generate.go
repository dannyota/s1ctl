package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "docs",
		Short:  "Generate documentation",
		Hidden: true,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newDocsGenerateCmd())
	return cmd
}

func newDocsGenerateCmd() *cobra.Command {
	var outDir string

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate command reference docs",
		Long: `Walk the command tree and generate a markdown file per command group
in the docs/commands/ directory. Also updates docs/_sidebar.md.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			root := cmd.Root()
			if err := os.MkdirAll(outDir, 0o750); err != nil {
				return err
			}

			groups := collectGroups(root)

			// Write per-group files.
			for _, g := range groups {
				path := filepath.Join(outDir, g.Filename)
				content := renderGroup(g)
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					return fmt.Errorf("write %s: %w", path, err)
				}
			}

			// Write index page.
			idx := renderIndex(groups)
			idxPath := filepath.Join(outDir, "README.md")
			if err := os.WriteFile(idxPath, []byte(idx), 0o644); err != nil {
				return fmt.Errorf("write %s: %w", idxPath, err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Generated %s + index in %s\n",
				pluralize(len(groups), "file"), outDir)
			return nil
		},
	}
	cmd.Flags().StringVar(&outDir, "out", "docs/commands", "output directory")
	return cmd
}

type cmdGroup struct {
	Name        string
	Short       string
	Filename    string
	Subcommands []*cobra.Command
}

func collectGroups(root *cobra.Command) []cmdGroup {
	var groups []cmdGroup

	// Top-level commands that are groups (have subcommands).
	for _, cmd := range root.Commands() {
		if cmd.Hidden || !cmd.HasSubCommands() {
			continue
		}
		g := cmdGroup{
			Name:     cmd.Name(),
			Short:    cmd.Short,
			Filename: cmd.Name() + ".md",
		}
		for _, sub := range cmd.Commands() {
			if sub.Hidden {
				continue
			}
			g.Subcommands = append(g.Subcommands, sub)
		}
		if len(g.Subcommands) > 0 {
			groups = append(groups, g)
		}
	}

	// Top-level commands without subcommands go into "global.md".
	var globals []*cobra.Command
	for _, cmd := range root.Commands() {
		if cmd.Hidden || cmd.HasSubCommands() {
			continue
		}
		globals = append(globals, cmd)
	}
	if len(globals) > 0 {
		groups = append(groups, cmdGroup{
			Name:        "global",
			Short:       "Top-level commands",
			Filename:    "global.md",
			Subcommands: globals,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}

func renderGroup(g cmdGroup) string {
	var b strings.Builder

	title := g.Name
	if g.Name == "global" {
		title = "Global commands"
	}
	fmt.Fprintf(&b, "# %s\n\n", title)
	fmt.Fprintf(&b, "%s\n", g.Short)

	for _, cmd := range g.Subcommands {
		b.WriteString("\n")
		renderCommand(&b, g.Name, cmd)
	}

	return b.String()
}

func renderCommand(b *strings.Builder, group string, cmd *cobra.Command) {
	name := cmd.Name()
	usage := cmd.UseLine()

	if group != "global" {
		fmt.Fprintf(b, "## %s %s\n\n", group, name)
	} else {
		fmt.Fprintf(b, "## %s\n\n", name)
	}

	fmt.Fprintf(b, "%s\n\n", cmd.Short)

	fmt.Fprintf(b, "```text\n%s\n```\n", usage)

	if cmd.Long != "" && cmd.Long != cmd.Short {
		long := strings.TrimSpace(cmd.Long)
		if strings.Contains(long, "\n  $") || strings.Contains(long, "\n  #") {
			fmt.Fprintf(b, "\n```text\n%s\n```\n", long)
		} else {
			fmt.Fprintf(b, "\n%s\n", long)
		}
	}

	flags := collectFlags(cmd)
	if len(flags) > 0 {
		b.WriteString("\n**Flags**\n\n")
		b.WriteString("| Flag | Type | Default | Description |\n")
		b.WriteString("|------|------|---------|-------------|\n")
		for _, f := range flags {
			b.WriteString(f)
			b.WriteString("\n")
		}
	}

	if cmd.Example != "" {
		fmt.Fprintf(b, "\n**Examples**\n\n```bash\n%s\n```\n", strings.TrimSpace(cmd.Example))
	}
}

func collectFlags(cmd *cobra.Command) []string {
	var rows []string
	cmd.NonInheritedFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		name := "--" + f.Name
		if f.Shorthand != "" {
			name = "-" + f.Shorthand + ", " + name
		}
		typ := f.Value.Type()
		def := f.DefValue
		if def == "" {
			def = "-"
		}
		if def == "[]" {
			def = "-"
		}
		desc := f.Usage
		row := fmt.Sprintf("| `%s` | %s | %s | %s |", name, typ, def, desc)
		rows = append(rows, row)
	})
	return rows
}

func renderIndex(groups []cmdGroup) string {
	var b strings.Builder
	b.WriteString("# Command reference\n\n")
	b.WriteString("Auto-generated from the command tree. Run `s1ctl docs generate` to update.\n\n")

	b.WriteString("| Group | Commands | Description |\n")
	b.WriteString("|-------|----------|-------------|\n")
	for _, g := range groups {
		names := make([]string, len(g.Subcommands))
		for i, cmd := range g.Subcommands {
			names[i] = cmd.Name()
		}
		fmt.Fprintf(&b, "| [%s](%s) | %s | %s |\n",
			g.Name, g.Filename, strings.Join(names, ", "), g.Short)
	}
	return b.String()
}
