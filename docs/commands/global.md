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

## drift

Report drift between committed config and live state

```text
s1ctl drift [flags]
```

Compare committed config-as-code against the live console for every sync
surface and report the difference without applying anything.

For each surface with a local directory under --dir-root, drift loads the
committed files, lists the live objects, and plans the reconcile: creates
(committed, not live), updates (committed, differs from live), live-only
(live, not committed) and unchanged. Surfaces without a local directory are
skipped — drift checks only what is committed.

The command is read-only: it lists, plans, and reports, and has no apply path.
Exit code is 0 when every checked surface is clean and 1 when any surface has
drift, so a CI job can fail on a non-zero exit.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir-root` | string | . | root directory containing per-surface config directories |
| `--surface` | stringSlice | - | limit to named surfaces (repeatable; default: all) |

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
