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
