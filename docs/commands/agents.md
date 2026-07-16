# agents

Manage endpoint agents

## agents abort-scan

Abort a running disk scan

```text
s1ctl agents abort-scan <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents approve-uninstall

Approve a pending uninstall request

```text
s1ctl agents approve-uninstall <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents broadcast

Display a broadcast message on an agent's endpoint

```text
s1ctl agents broadcast <agent-id> --message <text> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--message` | string | - | message text to broadcast (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents count

Count agents

```text
s1ctl agents count [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--site-id` | stringSlice | - | filter by site ID |

## agents decommission

Decommission an agent

```text
s1ctl agents decommission <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents disable

Disable an agent

```text
s1ctl agents disable <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents enable

Enable a disabled agent

```text
s1ctl agents enable <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents fetch-files

Fetch specific files from an agent to the console

```text
s1ctl agents fetch-files <agent-id> --path <file> [--path <file>...] [--password <pw>] [flags]
```

Fetch up to 10 files from a single agent. The files are uploaded to the
console encrypted with --password (required by the platform to open the
resulting archive). The password is never written to the audit log.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--password` | string | - | archive encryption password |
| `--path` | stringArray | - | absolute file path to fetch (repeatable, up to 10) (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents fetch-firewall-rules

Fetch the current firewall-rules inventory

```text
s1ctl agents fetch-firewall-rules <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents fetch-installed-apps

Fetch the installed-applications inventory

```text
s1ctl agents fetch-installed-apps <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents fetch-logs

Fetch agent logs to the console

```text
s1ctl agents fetch-logs <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents firewall-logging

Enable or disable firewall logging on an agent

```text
s1ctl agents firewall-logging <agent-id> --state on|off [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state` | string | - | "on" or "off" (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents get

Get agent details

```text
s1ctl agents get <agent-id>
```

## agents health

Classify agents by operational state

```text
s1ctl agents health [flags]
```

Fetch all agents and classify them as active, offline (disconnected),
decommissioned, or infected. Helps identify endpoints that need attention.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--site-id` | stringSlice | - | filter by site ID |

## agents isolate

Isolate agents from the network

```text
s1ctl agents isolate [agent-id...] [flags]
```

Disconnect agents from the network.

Specify agent IDs as arguments, or use --filter to match agents by API
query parameters (e.g. --filter infected=true --filter osTypes=windows).
Both can be combined. Dry-run by default; pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--filter` | stringArray | - | key=value filter (e.g. --filter infected=true) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents list

List agents

```text
s1ctl agents list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--active` | bool | false | filter by active status |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--group-id` | stringSlice | - | filter by group ID |
| `--infected` | bool | false | filter by infection status |
| `--limit` | int | 0 | max results per page (default 50) |
| `--machine-type` | stringSlice | - | filter by machine type (server, desktop, laptop) |
| `--network-status` | stringSlice | - | filter by network status (connected, disconnected) |
| `--os-type` | stringSlice | - | filter by OS type |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. computerName, lastActiveDate) |
| `--sort-order` | string | - | sort direction (asc, desc) |

## agents local-upgrade

Authorize or revoke local upgrade/downgrade on an agent

```text
s1ctl agents local-upgrade <agent-id> --state on|off [--until <timestamp>] [flags]
```

Set an agent's local upgrade/downgrade authorization.

--state on authorizes local upgrades until the --until expiration timestamp
(RFC3339, e.g. 2030-01-01T00:00:00Z), which is required. --state off revokes
the authorization.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state` | string | - | "on" or "off" (required) |
| `--until` | string | - | authorization expiration timestamp (RFC3339); required with --state on |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents local-upgrade-status

Show an agent's local upgrade/downgrade authorization

```text
s1ctl agents local-upgrade-status <agent-id>
```

## agents mark-up-to-date

Mark an agent as up to date

```text
s1ctl agents mark-up-to-date <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents move

Move an agent to a different group

```text
s1ctl agents move <agent-id> --group-id <target-group-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--group-id` | string | - | target group ID (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents move-to-site

Move an agent to a different site

```text
s1ctl agents move-to-site <agent-id> --site-id <target-site-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--site-id` | string | - | target site ID (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents outdated

List agents not on the latest version

```text
s1ctl agents outdated [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--site-id` | stringSlice | - | filter by site ID |

## agents passphrases

List agent maintenance passphrases (SECRET output)

```text
s1ctl agents passphrases [flags]
```

List agent maintenance passphrases. The passphrase column is secret
material used to run privileged local agent commands — treat the output
accordingly. Values are never written to the audit log.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--group-id` | stringSlice | - | filter by group ID |
| `--id` | stringSlice | - | filter by agent ID |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |

## agents randomize-uuid

Randomize the agent UUID

```text
s1ctl agents randomize-uuid <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents ranger

Enable or disable Ranger network discovery on an agent

```text
s1ctl agents ranger <agent-id> --state on|off [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state` | string | - | "on" or "off" (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents reconnect

Reconnect isolated agents

```text
s1ctl agents reconnect [agent-id...] [flags]
```

Reconnect previously network-isolated agents.

Specify agent IDs as arguments, or use --filter to match agents by API
query parameters (e.g. --filter networkStatuses=disconnected).
Both can be combined. Dry-run by default; pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--filter` | stringArray | - | key=value filter (e.g. --filter networkStatuses=disconnected) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents reject-uninstall

Reject a pending uninstall request

```text
s1ctl agents reject-uninstall <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents reset-config

Reset agent local configuration

```text
s1ctl agents reset-config <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents reset-passphrase

Reset the agent maintenance passphrase

```text
s1ctl agents reset-passphrase <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents restart

Restart the endpoint

```text
s1ctl agents restart <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents scan

Start full disk scan

```text
s1ctl agents scan <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents set-external-id

Set the external ID on an agent

```text
s1ctl agents set-external-id <agent-id> --external-id <value> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--external-id` | string | - | external ID value (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents shutdown

Shut down the endpoint

```text
s1ctl agents shutdown <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents uninstall

Uninstall an agent

```text
s1ctl agents uninstall <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents upgrade

Trigger agent software upgrade

```text
s1ctl agents upgrade [agent-id...] [flags]
```

Trigger a software update on one or more agents.

Exactly one of --package-id, --file-name, or --path is required to
identify the upgrade package. The --file-name option also requires
--os-type.

Specify agent IDs as arguments, or use --site-id / --group-id / --query
to target agents by filter. Dry-run by default.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--allow-downgrade` | bool | false | allow downgrading the agent version |
| `--file-name` | string | - | upgrade package file name |
| `--group-id` | stringSlice | - | filter by group ID |
| `--ignore-conflicts` | bool | false | ignore conflicts with active upgrade policies |
| `--os-type` | string | - | target OS type (linux, macos, windows) |
| `--package-id` | string | - | upgrade package ID |
| `--package-type` | string | - | package type (Agent, Ranger, AgentAndRanger) |
| `--path` | string | - | local path to upgrade package on the endpoint |
| `--query` | string | - | free text search filter |
| `--scheduled` | bool | false | upgrade according to agent upgrade schedule |
| `--site-id` | stringSlice | - | filter by site ID |
| `--yes` | bool | false | apply the action (default: dry-run) |

## agents versions

Show agent version distribution

```text
s1ctl agents versions [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--site-id` | stringSlice | - | filter by site ID |
