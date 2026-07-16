# automation

Manage hyperautomation workflows and executions

## automation activate

Activate a workflow version

```text
s1ctl automation activate <workflow-id> <version-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## automation create

Create (import) a workflow from a file

```text
s1ctl automation create --from-file <workflow.json> [flags]
```

Import a workflow definition that was previously exported.
The file should be JSON or YAML matching the export format.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | workflow definition file, JSON or YAML (required) |
| `--site-id` | stringSlice | - | scope to site ID |
| `--yes` | bool | false | apply the action (default: dry-run) |

## automation deactivate

Deactivate the active version of a workflow

```text
s1ctl automation deactivate <workflow-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## automation execution-get

Get a workflow execution by ID

```text
s1ctl automation execution-get <execution-id>
```

## automation execution-output

Get the output of a workflow execution

```text
s1ctl automation execution-output <execution-id>
```

## automation executions

List workflow executions

```text
s1ctl automation executions [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--group-id` | stringSlice | - | filter by group ID |
| `--limit` | int | 0 | max results per page (default 50) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort order (asc, desc) |
| `--state` | stringSlice | - | filter by state (Running, Completed, Error, etc.) |
| `--trigger-type` | stringSlice | - | filter by trigger type |
| `--workflow-id` | string | - | filter by workflow ID |

## automation export

Export a workflow version as JSON (suitable for import)

```text
s1ctl automation export <workflow-id> <version-id>
```

## automation get

Get a workflow version (export format)

```text
s1ctl automation get <workflow-id> <version-id>
```

## automation list

List automation workflows

```text
s1ctl automation list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--group-id` | stringSlice | - | filter by group ID |
| `--limit` | int | 0 | max results per page (default 50) |
| `--name` | string | - | filter by name (contains) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort order (asc, desc) |
| `--state` | stringSlice | - | filter by state (active, inactive, deactivated, draft) |
| `--tag` | stringSlice | - | filter by tag |
| `--trigger-type` | stringSlice | - | filter by trigger type |

## automation run

Trigger a manual workflow execution

```text
s1ctl automation run <workflow-id> <version-id> [flags]
```

Trigger a manual or scheduled workflow execution. This executes
tenant-side automation and may perform actions such as isolating agents,
sending emails, or calling external APIs.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## automation versions

List versions of a workflow

```text
s1ctl automation versions <workflow-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--group-id` | stringSlice | - | filter by group ID |
| `--site-id` | stringSlice | - | filter by site ID |
