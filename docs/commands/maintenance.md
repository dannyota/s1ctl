# maintenance

Manage task maintenance-window configuration

## maintenance export

Export maintenance-window occurrences as CSV

```text
s1ctl maintenance export --task-type <type> --out <file> [flags]
```

Export all maintenance-window occurrences for a scope as CSV. Only the flexible
(policy_payload) maintenance-window format is supported.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account IDs |
| `--group-id` | stringSlice | - | scope to group IDs |
| `--out` | string | maintenance-windows.csv | output file (use - for stdout) |
| `--site-id` | stringSlice | - | scope to site IDs |
| `--task-type` | string | - | task type, e.g. agents_upgrade (required) |
| `--tenant` | bool | false | scope to the global (tenant) level |

## maintenance get

Get the maintenance-window configuration for a scope

```text
s1ctl maintenance get [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account IDs |
| `--group-id` | stringSlice | - | scope to group IDs |
| `--site-id` | stringSlice | - | scope to site IDs |
| `--task-type` | string | - | task type, e.g. agents_upgrade (required) |
| `--tenant` | bool | false | scope to the global (tenant) level |

## maintenance get-flexible

Get the flexible maintenance-window configuration for a scope

```text
s1ctl maintenance get-flexible [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account IDs |
| `--group-id` | stringSlice | - | scope to group IDs |
| `--site-id` | stringSlice | - | scope to site IDs |
| `--task-type` | string | - | task type, e.g. agents_upgrade (required) |
| `--tenant` | bool | false | scope to the global (tenant) level |

## maintenance set

Set the maintenance-window configuration for a scope

```text
s1ctl maintenance set --task-type <type> --from-file <data.json> [flags]
```

Set the classic per-day maintenance-window configuration for a scope.

--from-file supplies the configuration data payload (maxConcurrent, timezoneGmt,
maintenanceWindowsByDay, inherit flags); the scope and task type come from the
--task-type and scope flags.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account IDs |
| `--from-file` | string | - | configuration data JSON file (required) |
| `--group-id` | stringSlice | - | scope to group IDs |
| `--site-id` | stringSlice | - | scope to site IDs |
| `--task-type` | string | - | task type, e.g. agents_upgrade (required) |
| `--tenant` | bool | false | scope to the global (tenant) level |
| `--yes` | bool | false | apply the action (default: dry-run) |

## maintenance set-flexible

Set the flexible maintenance-window configuration

```text
s1ctl maintenance set-flexible --from-file <body.json> [flags]
```

Set the flexible (policy_payload) maintenance-window configuration.

The flexible format is SKU-gated and open-ended, so --from-file must contain the
full request body: a "data" object with the policy payload and a "filter" object
with the task type and scope. The body is sent verbatim.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | full request body JSON file (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |
