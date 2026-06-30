package cli

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for s1ctl.

To load completions:

Bash:
  $ source <(s1ctl completion bash)
  # To load completions for each session, execute once:
  $ s1ctl completion bash > /etc/bash_completion.d/s1ctl

Zsh:
  $ source <(s1ctl completion zsh)
  # To load completions for each session, execute once:
  $ s1ctl completion zsh > "${fpath[1]}/_s1ctl"

Fish:
  $ s1ctl completion fish | source
  # To load completions for each session, execute once:
  $ s1ctl completion fish > ~/.config/fish/completions/s1ctl.fish

PowerShell:
  PS> s1ctl completion powershell | Out-String | Invoke-Expression
`,
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletionV2(os.Stdout, true)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}
