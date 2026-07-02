# Global commands

Top-level commands

## commands

List all available commands

```text
s1ctl commands
```

## completion

Generate shell completion scripts

```text
s1ctl completion [bash|zsh|fish|powershell]
```

```text
Generate shell completion scripts for s1ctl.

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
```

## doctor

Verify connectivity to all SentinelOne API surfaces

```text
s1ctl doctor
```

## help

Help about any command

```text
s1ctl help [command]
```

Help provides help for any command in the application.
Simply type s1ctl help [path to command] for full details.

## version

Print version information

```text
s1ctl version
```
