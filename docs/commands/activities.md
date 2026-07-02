# activities

View activity log

## activities count

Count activities

```text
s1ctl activities count [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |

## activities export

Export activities as CSV

```text
s1ctl activities export [flags]
```

Bulk export the activity log as CSV. Output goes to stdout by default, or to a file with --out.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--activity-type` | intSlice | - | filter by activity type ID |
| `--end` | string | - | activities before this timestamp (ISO 8601) |
| `--group-id` | stringSlice | - | filter by group ID |
| `--out` | string | - | write to file instead of stdout |
| `--site-id` | stringSlice | - | filter by site ID |
| `--start` | string | - | activities after this timestamp (ISO 8601) |

## activities list

List activities

```text
s1ctl activities list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--activity-type` | intSlice | - | filter by activity type ID |
| `--all` | bool | false | fetch all pages |
| `--created-after` | string | - | filter activities after this date (ISO 8601) |
| `--created-before` | string | - | filter activities before this date (ISO 8601) |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--site-id` | stringSlice | - | filter by site ID |

## activities types

List available activity types

```text
s1ctl activities types
```
