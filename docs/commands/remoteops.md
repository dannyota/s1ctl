# remoteops

Manage remote operations and scripts

## remoteops content

Print a remote script's content

```text
s1ctl remoteops content <script-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | - | write script content to file (default: stdout) |

## remoteops get

Get a remote script

```text
s1ctl remoteops get <script-id>
```

## remoteops guardrails

Manage remote-script execution guardrails

```text
s1ctl remoteops guardrails
```

Manage guardrails that require approval before scripts run on large numbers
of endpoints. A guardrail is configured per scope (account, site, or group).

## remoteops list

List remote scripts

```text
s1ctl remoteops list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |

## remoteops pending

Manage pending remote-script executions awaiting approval

```text
s1ctl remoteops pending
```

## remoteops results

Get remote script execution results

```text
s1ctl remoteops results <parent-task-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--status` | stringSlice | - | filter by status (created, pending, in_progress, completed, failed, canceled, expired) |

## remoteops run

Execute a remote script on agents

```text
s1ctl remoteops run <script-id> [flags]
```

Run a remote script from the Script Library on targeted agents.
Requires at least one targeting flag (--agent-id, --site-id, or --group-id).

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--agent-id` | stringSlice | - | target agent IDs |
| `--description` | string | - | task description (default: s1ctl remote script execution) |
| `--group-id` | stringSlice | - | target group IDs |
| `--input-params` | string | - | script input parameters |
| `--output-dest` | string | SentinelCloud | output destination (SentinelCloud, Local, None, SingularityXDR) |
| `--site-id` | stringSlice | - | target site IDs |
| `--timeout` | int | 0 | script runtime timeout in seconds (60-172800) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## remoteops update

Update a remote script's metadata

```text
s1ctl remoteops update <script-id> --from-file <script.json> [flags]
```

Update the metadata of a remote script (name, type, OS types, timeout, and
input requirements) from a JSON file. This changes the script's properties but
not its content.

The file holds the "data" object of the update body, for example:

  {
    "scriptName": "Collect Logs",
    "scriptType": "dataCollection",
    "osTypes": ["linux", "macos"],
    "inputRequired": false,
    "inputExample": "-",
    "inputInstructions": "-",
    "scriptRuntimeTimeoutSeconds": 3600
  }

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | JSON file with the update data object (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## remoteops upload-limits

Show package upload size limits

```text
s1ctl remoteops upload-limits
```
