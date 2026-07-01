# agents

Manage endpoint agents

## agents connect

Reconnect an isolated agent

```text
s1ctl agents connect <agent-id> [flags]
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

Network-isolate an agent

```text
s1ctl agents isolate <agent-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
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

## agents scan

Start full disk scan

```text
s1ctl agents scan <agent-id> [flags]
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

Specify agent IDs as arguments, or use --site-id / --group-id / --query
to target agents by filter. Dry-run by default.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--group-id` | stringSlice | - | filter by group ID |
| `--query` | string | - | free text search filter |
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
