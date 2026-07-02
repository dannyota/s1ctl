# remoteops

Remote operations and scripts

## remoteops get

Get a remote script

```text
s1ctl remoteops get <script-id>
```

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
