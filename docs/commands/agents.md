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

Network-isolate agents

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

## agents randomize-uuid

Randomize the agent UUID

```text
s1ctl agents randomize-uuid <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
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
