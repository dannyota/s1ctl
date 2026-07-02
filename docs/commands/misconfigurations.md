# misconfigurations

Manage xSPM misconfigurations

## misconfigurations assign

Assign a misconfiguration to a user

```text
s1ctl misconfigurations assign <id> --user-id <user-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--user-id` | string | - | assignee user ID (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## misconfigurations export

Export misconfigurations to a CSV file

```text
s1ctl misconfigurations export [flags]
```

Export misconfigurations matching the filters as CSV via
misconfigurationsExportToCsv. The API returns the full CSV inline; it is written
to --out, or to stdout when --out is omitted.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | - | output file (default: stdout) |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status |

## misconfigurations get

Get misconfiguration details

```text
s1ctl misconfigurations get <id>
```

## misconfigurations history

Show the history of a misconfiguration

```text
s1ctl misconfigurations history <id>
```

## misconfigurations list

List misconfigurations

```text
s1ctl misconfigurations list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--after` | string | - | pagination cursor |
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status |

## misconfigurations note-add

Add an investigation note to a misconfiguration

```text
s1ctl misconfigurations note-add <id> --text <text> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--text` | string | - | note text (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## misconfigurations note-delete

Delete a misconfiguration note

```text
s1ctl misconfigurations note-delete <note-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## misconfigurations note-update

Update the text of a misconfiguration note

```text
s1ctl misconfigurations note-update <note-id> --text <text> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--text` | string | - | new note text (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## misconfigurations notes

List investigation notes on a misconfiguration

```text
s1ctl misconfigurations notes <id>
```

## misconfigurations related-assets

List assets related to a misconfiguration

```text
s1ctl misconfigurations related-assets <id>
```

## misconfigurations status

Update misconfiguration status

```text
s1ctl misconfigurations status <id> <status> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## misconfigurations verdict

Update misconfiguration analyst verdict (TRUE_POSITIVE, FALSE_POSITIVE, SUSPICIOUS, UNDEFINED)

```text
s1ctl misconfigurations verdict <id> <verdict> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |
